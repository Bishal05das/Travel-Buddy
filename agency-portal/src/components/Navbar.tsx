"use client";

import Link from "next/link";
import { useAsync } from "@/hooks/useAsync";
import { useAuthActions, useSession } from "@/lib/auth";
import { CUSTOMER_URL } from "@/lib/config";
import { getMyProfile } from "@/lib/endpoints";
import { Button, ButtonLink } from "./ui/Button";

export function Navbar() {
  const { session, ready } = useSession();
  const { logout } = useAuthActions();

  return (
    <header className="border-b border-slate-200 bg-white">
      <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-3 px-4 py-3">
        <Link href="/" className="text-lg font-semibold text-teal-800">
          Travel Buddy · Agency Portal
        </Link>
        <nav className="flex flex-wrap items-center gap-2">
          {session?.role === "member" && <ButtonLink href="/dashboard" variant="ghost">Dashboard</ButtonLink>}
          {session?.role === "super" && <ButtonLink href="/admin" variant="ghost">Admin</ButtonLink>}
          <a href={CUSTOMER_URL} className="rounded-md px-3 py-2 text-sm font-medium text-slate-600 hover:text-slate-900">
            Customer site
          </a>
        </nav>
        <div className="flex items-center gap-3">
          {session?.role === "member" && <MemberIdentity key={session.userId} />}
          {session?.role === "super" && <span className="text-sm text-slate-600">Platform admin</span>}
          {ready && (session ? (
            <Button
              variant="secondary"
              size="sm"
              onClick={() => {
                logout();
                window.location.assign(new URL("/login", window.location.origin).toString());
              }}
            >
              Log out
            </Button>
          ) : (
            <ButtonLink href="/login" size="sm">Log in</ButtonLink>
          ))}
        </div>
      </div>
    </header>
  );
}

function MemberIdentity() {
  const profile = useAsync(getMyProfile, []);

  return (
    <Link href="/profile" className="rounded-md px-2 py-1 text-sm hover:bg-slate-100 focus-visible:outline-2 focus-visible:outline-teal-700">
      <span className="block max-w-56 truncate font-medium text-slate-900">{profile.data?.name ?? "My profile"}</span>
      <span className="block max-w-56 truncate text-xs text-slate-500">
        {profile.data ? `${profile.data.is_owner ? "Agency owner" : profile.data.role_name} · View profile` : profile.loading ? "Loading account…" : "View account details"}
      </span>
    </Link>
  );
}
