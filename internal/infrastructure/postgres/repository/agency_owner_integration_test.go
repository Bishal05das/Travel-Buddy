package repository_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/bishal05das/travelbuddy/internal/domain"
	"github.com/bishal05das/travelbuddy/internal/infrastructure/postgres/repository"
	memberusecase "github.com/bishal05das/travelbuddy/internal/usecase/agencyMember"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// Run only against a disposable PostgreSQL database. Every table and fixture is
// created in a unique schema; the normal application database is never selected
// implicitly from .env or DATABASE_URL.
func TestAgencyOwnerIntegration(t *testing.T) {
	dsn := os.Getenv("OWNER_TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("set OWNER_TEST_DATABASE_URL to a disposable PostgreSQL database")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	defer db.Close()
	// SET search_path must apply to every operation, including transactions.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	schema := "owner_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	exec := func(query string, args ...any) {
		t.Helper()
		if _, err := db.ExecContext(ctx, query, args...); err != nil {
			t.Fatalf("fixture SQL failed: %v", err)
		}
	}
	exec("CREATE SCHEMA " + pq.QuoteIdentifier(schema))
	defer func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		if _, err := db.ExecContext(cleanupCtx, "DROP SCHEMA "+pq.QuoteIdentifier(schema)+" CASCADE"); err != nil {
			t.Errorf("remove isolated test schema: %v", err)
		}
	}()
	exec("SET search_path TO " + pq.QuoteIdentifier(schema) + ", public")

	agencyID, otherAgencyID := uuid.New(), uuid.New()
	ownerID, staffID, otherOwnerID, otherStaffID := uuid.New(), uuid.New(), uuid.New(), uuid.New()
	var ownerRole, staffRole, otherOwnerRole, otherStaffRole int
	insertRole := func(agency uuid.UUID, name string) int {
		t.Helper()
		var id int
		if err := db.QueryRowContext(ctx, `INSERT INTO roles (agency_id, role_name) VALUES ($1, $2) RETURNING role_id`, agency, name).Scan(&id); err != nil {
			t.Fatalf("insert role: %v", err)
		}
		return id
	}
	insertLegacyMember := func(id, agency uuid.UUID, role int, name, email, joined string) {
		t.Helper()
		exec(`INSERT INTO agency_members (member_id, agency_id, role_id, name, email, phone, password, joined_at)
              VALUES ($1, $2, $3, $4, $5, '+8801700000000', 'private-password-hash', $6)`,
			id, agency, role, name, email, joined)
	}
	paths, err := filepath.Glob("../../../../migrations/*.up.sql")
	if err != nil || len(paths) == 0 {
		t.Fatalf("find migrations: count=%d err=%v", len(paths), err)
	}
	sort.Strings(paths)
	seeded := false
	for _, path := range paths {
		if filepath.Base(path) == "000019_agency_owner.up.sql" {
			// An explicit Owner role wins even when a staff member joined first.
			exec(`INSERT INTO agency (agency_id, name) VALUES ($1, 'First Agency'), ($2, 'Other Agency')`, agencyID, otherAgencyID)
			ownerRole = insertRole(agencyID, " Owner ")
			staffRole = insertRole(agencyID, "Staff")
			otherOwnerRole = insertRole(otherAgencyID, "Manager")
			otherStaffRole = insertRole(otherAgencyID, "Staff")
			insertLegacyMember(staffID, agencyID, staffRole, "First Staff", "staff@example.test", "2024-01-01")
			insertLegacyMember(ownerID, agencyID, ownerRole, "Agency Owner", "owner@example.test", "2024-02-01")
			// An agency without an Owner role retains its earliest member as owner.
			insertLegacyMember(otherOwnerID, otherAgencyID, otherOwnerRole, "Original Manager", "manager@example.test", "2024-01-01")
			insertLegacyMember(otherStaffID, otherAgencyID, otherStaffRole, "Other Staff", "other@example.test", "2024-02-01")
			seeded = true
		}
		migration, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read migration %s: %v", path, err)
		}
		if _, err := db.ExecContext(ctx, string(migration)); err != nil {
			t.Fatalf("apply migration %s: %v", path, err)
		}
	}
	if !seeded {
		t.Fatal("owner migration was not found")
	}

	members := repository.NewAgencyMemberRepositoryDB(db)
	permissions := repository.NewPermissionRepositoryDB(db)
	roles := repository.NewRoleRepository(db)
	tx := repository.NewTxManager(db)
	create := memberusecase.NewCreateAgencyMemberUseCase(tx, members, roles)
	remove := memberusecase.NewDeleteAgencyMemberUseCase(members)
	update := memberusecase.NewUpdatePermissionUseCase(tx, members, roles)
	owner := domain.Actor{ID: ownerID, Role: domain.RoleMember, AgencyID: &agencyID}
	staff := domain.Actor{ID: staffID, Role: domain.RoleMember, AgencyID: &agencyID}
	super := domain.Actor{ID: uuid.New(), Role: domain.RoleSuper}

	t.Run("legacy_owner_backfill", func(t *testing.T) {
		for _, member := range []struct {
			id        uuid.UUID
			wantOwner bool
		}{{ownerID, true}, {staffID, false}, {otherOwnerID, true}, {otherStaffID, false}} {
			var got bool
			if err := db.GetContext(ctx, &got, `SELECT is_owner FROM agency_members WHERE member_id=$1`, member.id); err != nil {
				t.Fatal(err)
			}
			if got != member.wantOwner {
				t.Errorf("member %s is_owner=%t, want %t", member.id, got, member.wantOwner)
			}
		}
	})

	t.Run("owner_permissions_without_role_grants", func(t *testing.T) {
		var grantCount int
		if err := db.GetContext(ctx, &grantCount, `SELECT COUNT(*) FROM role_permissions WHERE role_id=$1`, ownerRole); err != nil || grantCount != 0 {
			t.Fatalf("owner role grants=%d, err=%v", grantCount, err)
		}
		var allPermissions []domain.Permission
		if err := db.SelectContext(ctx, &allPermissions, `SELECT permission_id, name, resource, action FROM permissions`); err != nil {
			t.Fatal(err)
		}
		if len(allPermissions) == 0 {
			t.Fatal("permission seeds are missing")
		}
		wantIDs := make([]int, 0, len(allPermissions))
		for _, permission := range allPermissions {
			allowed, err := permissions.MemberHasPermission(ctx, ownerID, agencyID, permission.Resource, permission.Action)
			if err != nil || !allowed {
				t.Errorf("owner denied %s:%s: allowed=%t err=%v", permission.Resource, permission.Action, allowed, err)
			}
			wantIDs = append(wantIDs, permission.PermissionID)
		}
		gotIDs, err := members.GetPermissionIDs(ctx, ownerID)
		if err != nil {
			t.Fatal(err)
		}
		sort.Ints(gotIDs)
		sort.Ints(wantIDs)
		if !reflect.DeepEqual(gotIDs, wantIDs) {
			t.Errorf("owner permission IDs=%v, want %v", gotIDs, wantIDs)
		}
		allowed, err := permissions.MemberHasPermission(ctx, staffID, agencyID, "member", "delete")
		if err != nil || allowed {
			t.Errorf("staff without grants may delete: allowed=%t err=%v", allowed, err)
		}
	})

	t.Run("owner_creates_and_removes_staff", func(t *testing.T) {
		ids, err := members.GetPermissionIDs(ctx, ownerID)
		if err != nil {
			t.Fatal(err)
		}
		req := &domain.CreateMemberRequest{
			AgencyID: agencyID, Name: "New Staff", Email: "new-staff@example.test",
			Phone: "+8801700000001", Password: "StaffTestPassword42", RoleName: "Agent", Permissions: ids,
		}
		if err := create.Execute(ctx, owner, req); err != nil {
			t.Fatalf("owner could not add staff: %v", err)
		}
		added, err := members.FindMember(ctx, req.Email)
		if err != nil || added == nil {
			t.Fatalf("find added staff: member=%v err=%v", added, err)
		}
		profile, err := members.GetMemberProfile(ctx, added.MemberID, agencyID)
		if err != nil || profile.IsOwner {
			t.Fatalf("new staff incorrectly owns agency: profile=%+v err=%v", profile, err)
		}
		if added.Password == req.Password {
			t.Error("staff password was stored in plaintext")
		}
		if err := remove.Execute(ctx, owner, added.MemberID); err != nil {
			t.Fatalf("owner could not remove staff: %v", err)
		}
		if _, err := members.GetMemberProfile(ctx, added.MemberID, agencyID); !errors.Is(err, domain.ErrMemberNotFound) {
			t.Errorf("removed staff still has profile: %v", err)
		}
	})

	t.Run("owner_provisioning_restricted_to_platform_admin", func(t *testing.T) {
		for name, actor := range map[string]domain.Actor{"owner": owner, "staff": staff} {
			if err := create.Execute(ctx, actor, &domain.CreateMemberRequest{AgencyID: agencyID, RoleName: " OwNeR "}); !errors.Is(err, domain.ErrOwnerCreationForbidden) {
				t.Errorf("%s can provision another owner: %v", name, err)
			}
		}
		freshAgencyID := uuid.New()
		exec(`INSERT INTO agency (agency_id, name) VALUES ($1, 'Fresh Agency')`, freshAgencyID)
		ids, err := members.GetPermissionIDs(ctx, ownerID)
		if err != nil {
			t.Fatal(err)
		}
		req := &domain.CreateMemberRequest{
			AgencyID: freshAgencyID, Name: "Fresh Owner", Email: "fresh-owner@example.test",
			Phone: "+8801700000003", Password: "OwnerTestPassword42", RoleName: "Owner", Permissions: ids,
		}
		if err := create.Execute(ctx, super, req); err != nil {
			t.Fatalf("platform admin could not provision agency owner: %v", err)
		}
		created, err := members.FindMember(ctx, req.Email)
		if err != nil || created == nil {
			t.Fatalf("find provisioned owner: member=%v err=%v", created, err)
		}
		profile, err := members.GetMemberProfile(ctx, created.MemberID, freshAgencyID)
		if err != nil || !profile.IsOwner {
			t.Fatalf("provisioned account is not owner: profile=%+v err=%v", profile, err)
		}
		if err := remove.Execute(ctx, super, created.MemberID); !errors.Is(err, domain.ErrOwnerProtected) {
			t.Errorf("platform admin can remove newly provisioned owner: %v", err)
		}
	})

	t.Run("all_callers_cannot_remove_or_restrict_owner", func(t *testing.T) {
		// Give staff every permission so the owner guard is independently tested.
		exec(`INSERT INTO role_permissions (role_id, permission_id) SELECT $1, permission_id FROM permissions`, staffRole)
		for name, actor := range map[string]domain.Actor{"owner": owner, "staff": staff, "platform_admin": super} {
			t.Run(name, func(t *testing.T) {
				if err := remove.Execute(ctx, actor, ownerID); !errors.Is(err, domain.ErrOwnerProtected) {
					t.Errorf("delete owner error=%v, want ErrOwnerProtected", err)
				}
				if err := update.Execute(ctx, actor, ownerID, &domain.UpdatePermissionRequest{Permissions: []int{}}); !errors.Is(err, domain.ErrOwnerProtected) {
					t.Errorf("restrict owner error=%v, want ErrOwnerProtected", err)
				}
			})
		}
		shared := &domain.AgencyMember{AgencyID: agencyID, RoleID: ownerRole, Name: "Shared Role Staff", Email: "shared@example.test", Phone: "+8801700000002", Password: "fixture-hash"}
		if err := members.CreateMember(ctx, shared); err != nil {
			t.Fatal(err)
		}
		if err := update.Execute(ctx, super, shared.MemberID, &domain.UpdatePermissionRequest{Permissions: []int{}}); !errors.Is(err, domain.ErrOwnerProtected) {
			t.Errorf("shared owner role could be modified: %v", err)
		}
	})

	t.Run("database_blocks_deletion_cascades_and_demotion", func(t *testing.T) {
		cases := []struct {
			name, query string
			args        []any
		}{
			{"direct_delete", `DELETE FROM agency_members WHERE member_id=$1`, []any{ownerID}},
			{"agency_cascade", `DELETE FROM agency WHERE agency_id=$1`, []any{agencyID}},
			{"role_cascade", `DELETE FROM roles WHERE role_id=$1`, []any{ownerRole}},
			{"demotion", `UPDATE agency_members SET is_owner=FALSE WHERE member_id=$1`, []any{ownerID}},
			{"agency_transfer", `UPDATE agency_members SET agency_id=$1 WHERE member_id=$2`, []any{otherAgencyID, ownerID}},
			{"role_reassignment", `UPDATE agency_members SET role_id=$1 WHERE member_id=$2`, []any{staffRole, ownerID}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := db.ExecContext(ctx, tc.query, tc.args...)
				var postgresError *pq.Error
				if !errors.As(err, &postgresError) || postgresError.Constraint != "agency_owner_protected" {
					t.Errorf("owner protection error=%v, want agency_owner_protected constraint", err)
				}
			})
		}
		if err := repository.NewAgencyRepositoryDB(db).DeleteAgency(ctx, agencyID); !errors.Is(err, domain.ErrOwnerProtected) {
			t.Errorf("agency cascade error not mapped to owner protection: %v", err)
		}
		var intact bool
		if err := db.GetContext(ctx, &intact, `SELECT is_owner AND agency_id=$2 AND role_id=$3 FROM agency_members WHERE member_id=$1`, ownerID, agencyID, ownerRole); err != nil || !intact {
			t.Errorf("owner changed after rejected mutations: intact=%t err=%v", intact, err)
		}
		_, err := db.ExecContext(ctx, `UPDATE agency_members SET is_owner=TRUE WHERE member_id=$1`, staffID)
		var postgresError *pq.Error
		if !errors.As(err, &postgresError) || postgresError.Code != "23505" {
			t.Errorf("second owner allowed: %v", err)
		}
	})

	t.Run("agency_isolation", func(t *testing.T) {
		allowed, err := permissions.MemberHasPermission(ctx, ownerID, otherAgencyID, "member", "delete")
		if err != nil || allowed {
			t.Errorf("owner allowed into another agency: allowed=%t err=%v", allowed, err)
		}
		for _, id := range []uuid.UUID{otherOwnerID, otherStaffID} {
			if err := remove.Execute(ctx, owner, id); err == nil || errors.Is(err, domain.ErrOwnerProtected) {
				t.Errorf("foreign member delete should return not-found without disclosing ownership: %v", err)
			}
			if _, err := members.GetMemberProfile(ctx, id, agencyID); !errors.Is(err, domain.ErrMemberNotFound) {
				t.Errorf("foreign profile exposed: %v", err)
			}
		}
		if err := create.Execute(ctx, owner, &domain.CreateMemberRequest{AgencyID: otherAgencyID, RoleName: "Staff"}); err == nil {
			t.Error("owner created staff in another agency")
		}
		if err := create.Execute(ctx, staff, &domain.CreateMemberRequest{AgencyID: agencyID, RoleName: " OWNER "}); !errors.Is(err, domain.ErrOwnerCreationForbidden) {
			t.Errorf("staff can create another owner: %v", err)
		}
	})

	t.Run("profile_identity_and_password_exclusion", func(t *testing.T) {
		profile, err := members.GetMemberProfile(ctx, ownerID, agencyID)
		if err != nil {
			t.Fatal(err)
		}
		if profile.MemberID != ownerID || profile.AgencyID != agencyID || profile.AgencyName != "First Agency" || profile.Name != "Agency Owner" || profile.Email != "owner@example.test" || profile.Phone != "+8801700000000" || profile.RoleName != " Owner " || !profile.IsOwner {
			t.Errorf("incorrect profile: %+v", profile)
		}
		payload, err := json.Marshal(profile)
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(strings.ToLower(string(payload)), "password") || strings.Contains(string(payload), "private-password-hash") {
			t.Errorf("profile exposes password: %s", payload)
		}
		staffProfile, err := members.GetMemberProfile(ctx, staffID, agencyID)
		if err != nil || staffProfile.IsOwner || staffProfile.Name != "First Staff" {
			t.Errorf("staff profile incorrect: %+v err=%v", staffProfile, err)
		}
		if _, err := members.GetMemberProfile(ctx, uuid.New(), agencyID); !errors.Is(err, domain.ErrMemberNotFound) {
			t.Errorf("unknown member profile error=%v", err)
		}
	})
}
