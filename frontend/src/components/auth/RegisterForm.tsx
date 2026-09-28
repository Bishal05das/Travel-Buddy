"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useAuthActions } from "@/lib/auth";
import { loginUser, register } from "@/lib/endpoints";
import { Button } from "../ui/Button";
import { Alert } from "../ui/Feedback";
import { Input } from "../ui/Field";
import { Card } from "../ui/PageHeader";

export function RegisterForm({ next }: { next?: string }) {
  const router = useRouter();
  const { login } = useAuthActions();
  const [form, setForm] = useState({ name: "", email: "", phone: "", password: "", confirm: "" });
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const set = (key: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement>) => setForm({ ...form, [key]: e.target.value });

  function validate() {
    const e: Record<string, string> = {};
    if (form.name.trim().length < 2) e.name = "Enter your name.";
    if (!/^\S+@\S+\.\S+$/.test(form.email)) e.email = "Enter a valid email address.";
    if (!/^\+[1-9]\d{6,14}$/.test(form.phone)) e.phone = "Use international format, e.g. +8801712345678.";
    if (form.password.length < 8 || form.password.length > 64) e.password = "Use 8 to 64 characters.";
    if (form.confirm !== form.password) e.confirm = "Passwords don't match.";
    setErrors(e);
    return Object.keys(e).length === 0;
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!validate()) return;
    setSubmitting(true);
    setError(null);
    try {
      await register({ name: form.name.trim(), email: form.email.trim(), phone: form.phone.trim(), password: form.password });
      // Log straight in so the user can continue where they were.
      login(await loginUser(form.email.trim(), form.password));
      router.replace(next ?? "/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Sign-up failed.");
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <Card className="space-y-5 p-6">
        <div>
          <h1 className="text-xl font-semibold">Create an account</h1>
          <p className="text-sm text-slate-600">Book tours and keep track of your trips.</p>
        </div>
        <form onSubmit={submit} className="space-y-4" noValidate>
          <Input label="Full name" autoComplete="name" value={form.name} error={errors.name} onChange={set("name")} />
          <Input label="Email" type="email" autoComplete="email" value={form.email} error={errors.email} onChange={set("email")} />
          <Input label="Phone" type="tel" autoComplete="tel" placeholder="+8801712345678" value={form.phone} error={errors.phone} onChange={set("phone")} />
          <Input label="Password" type="password" autoComplete="new-password" value={form.password} error={errors.password} onChange={set("password")} />
          <Input label="Confirm password" type="password" autoComplete="new-password" value={form.confirm} error={errors.confirm} onChange={set("confirm")} />
          {error && <Alert tone="error">{error}</Alert>}
          <Button type="submit" loading={submitting} className="w-full">
            Sign up
          </Button>
        </form>
        <p className="text-center text-sm text-slate-600">
          Already have an account?{" "}
          <Link href={`/login${next ? `?next=${encodeURIComponent(next)}` : ""}`} className="font-medium text-teal-700 hover:underline">
            Log in
          </Link>
        </p>
      </Card>
    </div>
  );
}
