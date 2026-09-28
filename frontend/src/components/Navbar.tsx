"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { useAuthActions, useSession } from "@/lib/auth";
import { Button, ButtonLink } from "./ui/Button";

function NavLink({ href, children }: { href: string; children: React.ReactNode }) {
  const pathname = usePathname();
  const active = href === "/" ? pathname === "/" : pathname.startsWith(href);
  return (
    <Link
      href={href}
      aria-current={active ? "page" : undefined}
      className={`rounded-md px-3 py-2 text-sm font-medium ${active ? "bg-teal-50 text-teal-800" : "text-slate-600 hover:text-slate-900"}`}
    >
      {children}
    </Link>
  );
}

export function Navbar() {
  const { session, ready } = useSession();
  const { logout } = useAuthActions();

  return (
    <header className="border-b border-slate-200 bg-white">
      <div className="mx-auto flex max-w-6xl flex-wrap items-center justify-between gap-3 px-4 py-3">
        <Link href="/" className="text-lg font-semibold text-teal-800">
          Travel Buddy
        </Link>
        <nav className="flex flex-wrap items-center gap-1">
          <NavLink href="/">Explore</NavLink>
          <NavLink href="/search">Search</NavLink>
          {session?.role === "user" && <NavLink href="/bookings">My bookings</NavLink>}
          {session?.role === "member" && <NavLink href="/dashboard">Agency dashboard</NavLink>}
          {session?.role === "super" && <NavLink href="/admin">Admin</NavLink>}
        </nav>
        <div className="flex items-center gap-2">
          {!ready ? null : session ? (
            <Button
              variant="secondary"
              size="sm"
              onClick={() => {
                logout();
                // A full navigation: a client-side push would race with route
                // guards redirecting to /login, and it also drops any cached
                // data from the previous user.
                // eslint-disable-next-line @next/next/no-location-assign-relative-destination -- intentional hard navigation, see above
                window.location.assign("/");
              }}
            >
              Log out
            </Button>
          ) : (
            <>
              <ButtonLink href="/login" variant="ghost" size="sm">
                Log in
              </ButtonLink>
              <ButtonLink href="/register" size="sm">
                Sign up
              </ButtonLink>
            </>
          )}
        </div>
      </div>
    </header>
  );
}
