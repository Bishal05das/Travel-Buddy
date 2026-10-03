"use client";

import Link from "next/link";
import { useState } from "react";
import { useAsync } from "@/hooks/useAsync";
import { CUSTOMER_URL } from "@/lib/config";
import { createAgency, listPermissions } from "@/lib/endpoints";
import { MemberForm } from "./members/MemberForm";
import { RequireRole } from "./RequireRole";
import { Button } from "./ui/Button";
import { Alert, ErrorState, LoadingState } from "./ui/Feedback";
import { Input } from "./ui/Field";
import { Card, PageHeader } from "./ui/PageHeader";

export function AdminOnboarding() {
  return (
    <RequireRole roles={["super"]}>
      {() => (
        <div>
          <PageHeader title="Onboard an agency" description="Create the agency, then its owner account. The owner can manage the agency and its staff." />
          <Onboarding />
        </div>
      )}
    </RequireRole>
  );
}

function Onboarding() {
  const [agency, setAgency] = useState<{ id: string; name: string } | null>(null);
  const [memberEmail, setMemberEmail] = useState<string | null>(null);

  return (
    <div className="max-w-2xl space-y-6">
      <Card>
        <h2 className="mb-4 font-semibold">1. Agency details</h2>
        {agency ? (
          <Alert tone="success">
            <strong>{agency.name}</strong> was created.{" "}
            <Link href={`${CUSTOMER_URL}/agencies/${agency.id}`} className="underline">
              View its page
            </Link>
          </Alert>
        ) : (
          <AgencyForm onCreated={setAgency} />
        )}
      </Card>
      <Card className={agency ? "" : "opacity-60"}>
        <h2 className="mb-4 font-semibold">2. Agency owner account</h2>
        {!agency ? (
          <p className="text-sm text-slate-500">Create the agency first.</p>
        ) : memberEmail ? (
          <div className="space-y-3">
            <Alert tone="success">
              {memberEmail} can now log in as the agency owner and manage {agency.name}.
            </Alert>
            <Button
              variant="secondary"
              onClick={() => {
                setAgency(null);
                setMemberEmail(null);
              }}
            >
              Onboard another agency
            </Button>
          </div>
        ) : (
          <FirstMember agencyId={agency.id} onCreated={setMemberEmail} />
        )}
      </Card>
    </div>
  );
}

function AgencyForm({ onCreated }: { onCreated: (a: { id: string; name: string }) => void }) {
  const [v, setV] = useState({ name: "", address: "", registration_id: "" });
  const [image, setImage] = useState<File | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const set = (key: keyof typeof v) => (e: React.ChangeEvent<HTMLInputElement>) => setV({ ...v, [key]: e.target.value });

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    const errs: Record<string, string> = {};
    if (v.name.trim().length < 3) errs.name = "Use at least 3 characters.";
    if (v.address.trim().length < 3) errs.address = "Enter the agency's address.";
    if (v.registration_id.trim().length < 3) errs.registration_id = "Enter the registration number.";
    if (!image) errs.image = "Add a logo or photo.";
    else if (!["image/jpeg", "image/png", "image/webp"].includes(image.type)) errs.image = "Use a JPG, PNG or WebP image.";
    setErrors(errs);
    if (Object.keys(errs).length) return;

    setSubmitting(true);
    setSubmitError(null);
    try {
      const form = new FormData();
      form.set("name", v.name.trim());
      form.set("address", v.address.trim());
      form.set("registration_id", v.registration_id.trim());
      form.set("image", image!);
      const res = await createAgency(form);
      onCreated({ id: res.agency_id, name: v.name.trim() });
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Could not create the agency.");
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={submit} className="space-y-4" noValidate>
      <Input label="Agency name" value={v.name} error={errors.name} onChange={set("name")} />
      <Input label="Address" value={v.address} error={errors.address} onChange={set("address")} />
      <Input label="Registration number" value={v.registration_id} error={errors.registration_id} onChange={set("registration_id")} />
      <Input
        label="Logo or photo"
        type="file"
        accept="image/jpeg,image/png,image/webp"
        error={errors.image}
        onChange={(e) => setImage(e.target.files?.[0] ?? null)}
      />
      {submitError && <Alert tone="error">{submitError}</Alert>}
      <Button type="submit" loading={submitting}>
        Create agency
      </Button>
    </form>
  );
}

function FirstMember({ agencyId, onCreated }: { agencyId: string; onCreated: (email: string) => void }) {
  const permissions = useAsync(listPermissions, []);
  if (permissions.loading && !permissions.data) return <LoadingState label="Loading permissions" />;
  if (permissions.error) return <ErrorState message={permissions.error} onRetry={permissions.reload} />;
  return (
    <div className="space-y-3">
      <p className="text-sm text-slate-600">All permissions are selected so this person can run the agency and add colleagues.</p>
      <MemberForm agencyId={agencyId} permissions={permissions.data ?? []} owner onCreated={onCreated} />
    </div>
  );
}
