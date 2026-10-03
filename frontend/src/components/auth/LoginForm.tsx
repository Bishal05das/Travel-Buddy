"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useEffect, useState } from "react";
import { homeFor, useAuthActions, useSession } from "@/lib/auth";
import { loginUser } from "@/lib/endpoints";
import { decodeSession } from "@/lib/session";
import { Button } from "../ui/Button";
import { Alert } from "../ui/Feedback";
import { Input } from "../ui/Field";
import { Card } from "../ui/PageHeader";

export function LoginForm({ next }: { next?: string }) {
  const router = useRouter();
  const { session, ready } = useSession();
  const { login } = useAuthActions();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    if (ready && session?.role === "user" && !submitting) router.replace(next ?? homeFor(session));
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
      const token = await loginUser(email.trim(), password);
      const account = decodeSession(token);
      if (!account) throw new Error("The server returned an invalid session.");
      if (account.role !== "user") throw new Error("Use the agency portal to sign in as staff or admin.");
      login(token);
      router.replace(next ?? homeFor(account));
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed.");
      setSubmitting(false);
    }
  }

  return (
    <div className="mx-auto max-w-md">
      <Card className="space-y-5 p-6">
        <div>
          <h1 className="text-xl font-semibold">Customer login</h1>
          <p className="text-sm text-slate-600">Log in to book tours and manage your bookings.</p>
        </div>
        <form onSubmit={submit} className="space-y-4" noValidate>
          <Input label="Email" type="email" autoComplete="email" value={email} onChange={(e) => setEmail(e.target.value)} />
          <Input label="Password" type="password" autoComplete="current-password" value={password} onChange={(e) => setPassword(e.target.value)} />
          {error && <Alert tone="error">{error}</Alert>}
          <Button type="submit" loading={submitting} className="w-full">Log in</Button>
        </form>
        <p className="text-center text-sm text-slate-600">
          No account yet?{" "}
          <Link href={`/register${next ? `?next=${encodeURIComponent(next)}` : ""}`} className="font-medium text-teal-700 hover:underline">Sign up</Link>
        </p>
      </Card>
    </div>
  );
}
