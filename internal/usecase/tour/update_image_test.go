package tourusecase_test

import (
	"context"
	"errors"
	"github.com/bishal05das/travelbuddy/internal/domain"
	mocks "github.com/bishal05das/travelbuddy/internal/mocks/repository"
	tourusecase "github.com/bishal05das/travelbuddy/internal/usecase/tour"
	"github.com/google/uuid"
	"testing"
)

func TestTourImageReplacementAgencyScope(t *testing.T) {
	ctx := context.Background()
	agencyID, otherAgencyID := uuid.New(), uuid.New()
	repo := mocks.NewMockTourRepository()
	tour := &domain.Tour{AgencyID: agencyID, ImagePath: "tours/old.png", Status: "open"}
	if err := repo.CreateTour(ctx, tour); err != nil {
		t.Fatal(err)
	}
	uc := tourusecase.NewUpdateTourImageUseCase(repo)
	for name, actor := range map[string]domain.Actor{
		"customer":       {Role: domain.RoleUser, AgencyID: &agencyID},
		"other agency":   {Role: domain.RoleMember, AgencyID: &otherAgencyID},
		"missing agency": {Role: domain.RoleMember},
	} {
		if _, err := uc.Execute(ctx, actor, agencyID, tour.TourID, "new.png"); !errors.Is(err, domain.ErrImageAccessDenied) {
			t.Errorf("%s error: %v", name, err)
		}
	}
	foreign := domain.Actor{Role: domain.RoleMember, AgencyID: &otherAgencyID}
	if _, err := uc.Execute(ctx, foreign, otherAgencyID, tour.TourID, "new.png"); !errors.Is(err, domain.ErrImageTargetNotFound) {
		t.Fatalf("foreign tour: %v", err)
	}
	owner := domain.Actor{Role: domain.RoleMember, AgencyID: &agencyID}
	old, err := uc.Execute(ctx, owner, agencyID, tour.TourID, "images/tours/new.png")
	if err != nil || old != "tours/old.png" {
		t.Fatalf("replace: %q %v", old, err)
	}
	if err := repo.UpdateTourStatus(ctx, tour.TourID, "cancelled", &agencyID); err != nil {
		t.Fatal(err)
	}
	if _, err := uc.Execute(ctx, owner, agencyID, tour.TourID, "next.png"); !errors.Is(err, domain.ErrTourCancelled) {
		t.Fatalf("cancelled tour: %v", err)
	}
}
