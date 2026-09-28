"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { homeFor, useAuthActions, useSession } from "@/lib/auth";
import { loginMember, loginUser } from "@/lib/endpoints";
import { decodeSession } from "@/lib/session";
import { Button } from "../ui/Button";
import { Alert } from "../ui/Feedback";
import { Input } from "../ui/Field";
import { Card } from "../ui/PageHeader";

type Mode = "user" | "member";

export function LoginForm({ next, initialMode }: { next?: string; initialMode: Mode }) {
  const router = useRouter();
  const { session, ready } = useSession();
  const { login } = useAuthActions();
  const [mode, setMode] = useState<Mode>(initialMode);
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  // Already logged in (e.g. via the back button): go where they were heading.
  useEffect(() => {
    if (ready && session && !submitting) router.replace(next ?? homeFor(session));
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
      const token = await (mode === "user" ? loginUser : loginMember)(email.trim(), password);
      const s = decodeSession(token);
      if (!s) throw new Error("The server returned an invalid session.");
      login(token);
      router.replace(next ?? homeFor(s));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed.");
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <Card className="space-y-5 p-6">
        <div>
          <h1 className="text-xl font-semibold">Log in</h1>
          <p className="text-sm text-slate-600">Welcome back to Travel Buddy.</p>
        </div>
        <div role="tablist" aria-label="Account type" className="grid grid-cols-2 rounded-lg bg-slate-100 p-1 text-sm">
          {(["user", "member"] as const).map((m) => (
            <button
              key={m}
              type="button"
              role="tab"
              aria-selected={mode === m}
              onClick={() => {
                setMode(m);
                setError(null);
              }}
              className={`rounded-md py-1.5 font-medium ${mode === m ? "bg-white shadow-sm" : "text-slate-600"}`}
            >
              {m === "user" ? "Customer" : "Agency staff"}
            </button>
          ))}
        </div>
        <form onSubmit={submit} className="space-y-4" noValidate>
          <Input label="Email" type="email" autoComplete="email" value={email} onChange={(e) => setEmail(e.target.value)} />
          <Input
            label="Password"
            type="password"
            autoComplete="current-password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          {error && <Alert tone="error">{error}</Alert>}
          <Button type="submit" loading={submitting} className="w-full">
            Log in
          </Button>
        </form>
        {mode === "user" ? (
          <p className="text-center text-sm text-slate-600">
            No account yet?{" "}
            <Link href={`/register${next ? `?next=${encodeURIComponent(next)}` : ""}`} className="font-medium text-teal-700 hover:underline">
              Sign up
            </Link>
          </p>
        ) : (
          <p className="text-center text-sm text-slate-500">Staff accounts are created by your agency.</p>
        )}
      </Card>
    </div>
  );
}
