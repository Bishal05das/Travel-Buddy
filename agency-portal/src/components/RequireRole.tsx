"use client";

import { usePathname, useRouter } from "next/navigation";
import { useEffect, type ReactNode } from "react";
import { homeFor, useSession } from "@/lib/auth";
import type { Session } from "@/lib/session";
import type { Role } from "@/lib/types";
import { ButtonLink } from "./ui/Button";
import { EmptyState, LoadingState } from "./ui/Feedback";

const roleNames: Record<Role, string> = { user: "customers", member: "agency staff", super: "platform admins" };

/**
 * Renders children only for the given roles. Logged-out visitors are sent
 * to the login page and come back here afterwards.
 */
export function RequireRole({ roles, children }: { roles: Role[]; children: (session: Session) => ReactNode }) {
  const { session, ready } = useSession();
  const router = useRouter();
  const pathname = usePathname();

  useEffect(() => {
    if (ready && !session) {
      const loginAs = roles.includes("super") ? "&as=admin" : "&as=staff";
      router.replace(`/login?next=${encodeURIComponent(pathname)}${loginAs}`);
    }
  }, [ready, session, router, pathname, roles]);

  if (!ready || !session) return <LoadingState label="Checking your session" />;
  if (!roles.includes(session.role)) {
    return (
      <EmptyState
        title={`This page is for ${roles.map((r) => roleNames[r]).join(" and ")}.`}
        description="You're logged in with a different kind of account."
        action={<ButtonLink href={homeFor(session)}>Go to your home page</ButtonLink>}
      />
    );
  }
  return <>{children(session)}</>;
}
