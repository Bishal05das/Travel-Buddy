"use client";

import { useState } from "react";
import { useAsync } from "@/hooks/useAsync";
import { deleteMember, listMembers, listPermissions, updateMemberPermissions } from "@/lib/endpoints";
import type { Member, Permission } from "@/lib/types";
import { MemberForm } from "../members/MemberForm";
import { PermissionPicker, permissionNames } from "../members/PermissionPicker";
import { Button } from "../ui/Button";
import { Alert, EmptyState, ErrorState, LoadingState } from "../ui/Feedback";
import { Card } from "../ui/PageHeader";
import { NoAccess, useDashboard } from "./DashboardShell";

export function TeamManager() {
  const { agencyId, session } = useDashboard();
  const members = useAsync(() => listMembers(agencyId), [agencyId]);
  const permissions = useAsync(listPermissions, []);
  const [adding, setAdding] = useState(false);
  const [notice, setNotice] = useState<{ tone: "success" | "error"; text: string } | null>(null);

  if (members.status === 403) return <NoAccess what="the team list" permission="member:read" />;
  if ((members.loading && !members.data) || (permissions.loading && !permissions.data)) return <LoadingState label="Loading team" />;
  if (members.error) return <ErrorState message={members.error} onRetry={members.reload} />;
  if (permissions.error) return <ErrorState message={permissions.error} onRetry={permissions.reload} />;
  const all = permissions.data ?? [];

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Team</h2>
        {!adding && <Button onClick={() => setAdding(true)}>Add member</Button>}
      </div>
      {notice && <Alert tone={notice.tone}>{notice.text}</Alert>}
      {adding && (
        <Card>
          <div className="mb-4 flex items-center justify-between">
            <h3 className="font-semibold">New team member</h3>
            <Button variant="ghost" size="sm" onClick={() => setAdding(false)}>
              Close
            </Button>
          </div>
          <p className="mb-4 text-sm text-slate-600">You can only grant permissions you have yourself.</p>
          <MemberForm
            agencyId={agencyId}
            permissions={all}
            onCreated={(email) => {
              setAdding(false);
              setNotice({ tone: "success", text: `${email} was added. They can log in as agency staff.` });
              members.reload();
            }}
          />
        </Card>
      )}
      {members.data && members.data.length === 0 ? (
        <EmptyState title="No team members yet" />
      ) : (
        <div className="space-y-3">
          {members.data?.map((m) => (
            <MemberRow
              key={m.member_id}
              member={m}
              all={all}
              isSelf={m.member_id === session.userId}
              onChanged={(text, tone = "success") => {
                setNotice({ tone, text });
                if (tone === "success") members.reload();
              }}
            />
          ))}
        </div>
      )}
    </div>
  );
}

function MemberRow({
  member,
  all,
  isSelf,
  onChanged,
}: {
  member: Member;
  all: Permission[];
  isSelf: boolean;
  onChanged: (text: string, tone?: "success" | "error") => void;
}) {
  const [editing, setEditing] = useState(false);
  const [selected, setSelected] = useState(member.permissions);
  const [busy, setBusy] = useState(false);

  async function save() {
    setBusy(true);
    try {
      await updateMemberPermissions(member.member_id, selected);
      setEditing(false);
      onChanged(`Permissions updated for ${member.name}.`);
    } catch (err) {
      onChanged(err instanceof Error ? err.message : "Could not update permissions.", "error");
    } finally {
      setBusy(false);
    }
  }

  async function remove() {
    if (!window.confirm(`Remove ${member.name} from the agency? They will no longer be able to log in.`)) return;
    setBusy(true);
    try {
      await deleteMember(member.member_id);
      onChanged(`${member.name} was removed.`);
    } catch (err) {
      onChanged(err instanceof Error ? err.message : "Could not remove the member.", "error");
      setBusy(false);
    }
  }

  return (
    <Card>
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <p className="font-medium">
            {member.name} {isSelf && <span className="text-xs text-slate-500">(you)</span>}
          </p>
          <p className="text-sm text-slate-500">
            {member.email} · {member.phone}
          </p>
          {!editing && <p className="mt-2 text-xs text-slate-600">{permissionNames(member.permissions, all)}</p>}
        </div>
        {!editing && !isSelf && (
          <div className="flex gap-1">
            <Button variant="ghost" size="sm" onClick={() => setEditing(true)}>
              Edit permissions
            </Button>
            <Button variant="danger-ghost" size="sm" loading={busy} onClick={remove}>
              Remove
            </Button>
          </div>
        )}
      </div>
      {editing && (
        <div className="mt-4 space-y-3">
          <PermissionPicker permissions={all} selected={selected} onChange={setSelected} />
          <div className="flex gap-2">
            <Button size="sm" loading={busy} onClick={save}>
              Save
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={() => {
                setSelected(member.permissions);
                setEditing(false);
              }}
            >
              Cancel
            </Button>
          </div>
        </div>
      )}
    </Card>
  );
}
