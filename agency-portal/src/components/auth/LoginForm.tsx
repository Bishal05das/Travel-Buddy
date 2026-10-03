"use client";

import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { useAuthActions, useSession } from "@/lib/auth";
import { loginMember, loginUser } from "@/lib/endpoints";
import { decodeSession, type Session } from "@/lib/session";
import { Button } from "../ui/Button";
import { Alert } from "../ui/Feedback";
import { Input } from "../ui/Field";
import { Card } from "../ui/PageHeader";

type Mode = "staff" | "admin";

function destinationFor(role: Session["role"], next?: string): string {
  if (role === "super") return next?.startsWith("/admin") ? next : "/admin";
  if (role === "member") {
    return next?.startsWith("/dashboard") || next?.startsWith("/tours/") ? next : "/dashboard";
  }
  return "/login";
}

export function LoginForm({ next, initialMode }: { next?: string; initialMode: Mode }) {
  const router = useRouter();
  const { session, ready } = useSession();
  const { login } = useAuthActions();
  const [mode, setMode] = useState<Mode>(initialMode);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (ready && session && session.role !== "user" && !submitting) router.replace(destinationFor(session.role, next));
  }, [ready, session, submitting, next, router]);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!email || password.length < 8) {
      setError("Enter your email and password (at least 8 characters).");
      return;
    }
    setSubmitting(true);
    setError(null);
    try {
      const token = await (mode === "staff" ? loginMember : loginUser)(email.trim(), password);
      const account = decodeSession(token);
      if (!account) throw new Error("The server returned an invalid session.");
      if (mode === "staff" && account.role !== "member") throw new Error("This is not an agency staff account.");
      if (mode === "admin" && account.role !== "super") throw new Error("This account is not a platform admin.");
      login(token);
      router.replace(destinationFor(account.role, next));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed.");
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <Card className="space-y-5 p-6">
        <div>
          <h1 className="text-xl font-semibold">Portal login</h1>
          <p className="text-sm text-slate-600">Manage your agency or onboard a new one.</p>
        </div>
        <div role="tablist" aria-label="Account type" className="grid grid-cols-2 rounded-lg bg-slate-100 p-1 text-sm">
          {(["staff", "admin"] as const).map((accountType) => (
            <button
              key={accountType}
              type="button"
              role="tab"
              aria-selected={mode === accountType}
              onClick={() => { setMode(accountType); setError(null); }}
              className={`rounded-md py-1.5 font-medium ${mode === accountType ? "bg-white shadow-sm" : "text-slate-600"}`}
            >
              {accountType === "staff" ? "Agency staff" : "Platform admin"}
            </button>
          ))}
        </div>
        <form onSubmit={submit} className="space-y-4" noValidate>
          <Input label="Email" type="email" autoComplete="email" value={email} onChange={(e) => setEmail(e.target.value)} />
          <Input label="Password" type="password" autoComplete="current-password" value={password} onChange={(e) => setPassword(e.target.value)} />
          {error && <Alert tone="error">{error}</Alert>}
          <Button type="submit" loading={submitting} className="w-full">Log in</Button>
        </form>
        <p className="text-center text-sm text-slate-500">Agency staff accounts are created by an agency manager.</p>
      </Card>
    </div>
  );
}
