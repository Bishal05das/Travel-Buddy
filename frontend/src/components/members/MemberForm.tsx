"use client";

import { useState } from "react";
import { createMember } from "@/lib/endpoints";
import type { Permission } from "@/lib/types";
import { Button } from "../ui/Button";
import { Alert } from "../ui/Feedback";
import { Input } from "../ui/Field";
import { PermissionPicker } from "./PermissionPicker";

const blank = { name: "", email: "", phone: "", password: "", role_name: "" };

/** Adds a staff member to an agency. */
export function MemberForm({
  agencyId,
  permissions,
  defaultAll = false,
  onCreated,
}: {
  agencyId: string;
  permissions: Permission[];
  defaultAll?: boolean;
  onCreated: (email: string) => void;
}) {
  const [v, setV] = useState(blank);
  const [selected, setSelected] = useState<number[]>(defaultAll ? permissions.map((p) => p.permission_id) : []);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const set = (key: keyof typeof blank) => (e: React.ChangeEvent<HTMLInputElement>) => setV({ ...v, [key]: e.target.value });

  function validate() {
    const e: Record<string, string> = {};
    if (v.name.trim().length < 3) e.name = "Use at least 3 characters.";
    if (!/^\S+@\S+\.\S+$/.test(v.email)) e.email = "Enter a valid email address.";
    if (!/^\+[1-9]\d{6,14}$/.test(v.phone)) e.phone = "Use international format, e.g. +8801712345678.";
    if (v.password.length < 8 || v.password.length > 64) e.password = "Use 8 to 64 characters.";
    if (v.role_name.trim().length < 2) e.role_name = "Describe the role, e.g. Manager.";
    if (selected.length === 0) e.permissions = "Choose at least one permission.";
    setErrors(e);
    return Object.keys(e).length === 0;
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!validate()) return;
    setSubmitting(true);
    setSubmitError(null);
    try {
      await createMember(agencyId, { ...v, name: v.name.trim(), email: v.email.trim(), role_name: v.role_name.trim(), permissions: selected });
      onCreated(v.email.trim());
      setV(blank);
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Could not add the member.");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={submit} className="space-y-4" noValidate>
      <div className="grid gap-4 sm:grid-cols-2">
        <Input label="Full name" value={v.name} error={errors.name} onChange={set("name")} />
        <Input label="Role title" placeholder="e.g. Manager, Sales" value={v.role_name} error={errors.role_name} onChange={set("role_name")} />
        <Input label="Email" type="email" value={v.email} error={errors.email} onChange={set("email")} />
        <Input label="Phone" type="tel" placeholder="+8801712345678" value={v.phone} error={errors.phone} onChange={set("phone")} />
        <Input
          label="Temporary password"
          type="password"
          autoComplete="new-password"
          value={v.password}
          error={errors.password}
          hint="Share it with the member so they can log in."
          onChange={set("password")}
        />
      </div>
      <PermissionPicker permissions={permissions} selected={selected} onChange={setSelected} error={errors.permissions} />
      {submitError && <Alert tone="error">{submitError}</Alert>}
      <Button type="submit" loading={submitting}>
        Add member
      </Button>
    </form>
  );
}
