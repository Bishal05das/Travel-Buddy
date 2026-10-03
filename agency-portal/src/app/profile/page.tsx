"use client";

import { RequireRole } from "@/components/RequireRole";
import { ButtonLink } from "@/components/ui/Button";
import { ErrorState, LoadingState } from "@/components/ui/Feedback";
import { Card, PageHeader } from "@/components/ui/PageHeader";
import { useAsync } from "@/hooks/useAsync";
import { getMyProfile } from "@/lib/endpoints";

export default function ProfilePage() {
  return <RequireRole roles={["member"]}>{(session) => <Profile key={session.userId} />}</RequireRole>;
}

function Profile() {
  const profile = useAsync(getMyProfile, []);

  if (profile.loading) return <LoadingState label="Loading your profile" />;
  if (profile.error) return <ErrorState message={profile.error} onRetry={profile.reload} />;
  if (!profile.data) return null;

  const member = profile.data;
  const details = [
    ["Full name", member.name],
    ["Email", member.email],
    ["Phone", member.phone],
    ["Agency", member.agency_name],
    ["Role", member.is_owner ? "Agency owner" : member.role_name],
  ];

  return (
    <div className="mx-auto max-w-3xl">
      <PageHeader
        title="My profile"
        description="Your signed-in agency account."
        actions={<ButtonLink href="/dashboard" variant="secondary">Dashboard</ButtonLink>}
      />
      <Card>
        <div className="mb-6 flex flex-wrap items-center gap-3">
          <h2 className="text-xl font-semibold text-slate-900">{member.name}</h2>
          <span className="rounded-full bg-teal-50 px-3 py-1 text-xs font-medium text-teal-800">
            {member.is_owner ? "Agency owner" : "Staff"}
          </span>
        </div>
        <dl className="grid gap-5 sm:grid-cols-2">
          {details.map(([label, value]) => (
            <div key={label} className="min-w-0">
              <dt className="text-sm text-slate-500">{label}</dt>
              <dd className="mt-1 break-words font-medium text-slate-900">{value}</dd>
            </div>
          ))}
        </dl>
        {member.is_owner && (
          <p className="mt-6 border-t border-slate-200 pt-4 text-sm text-slate-600">
            You have full access to your agency and can add or remove staff. Your owner account cannot be removed, and its permissions cannot be changed.
          </p>
        )}
      </Card>
    </div>
  );
}
