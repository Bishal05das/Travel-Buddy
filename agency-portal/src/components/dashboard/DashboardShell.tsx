"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { createContext, useContext, type ReactNode } from "react";
import { useAsync } from "@/hooks/useAsync";
import { getAgency } from "@/lib/endpoints";
import type { Session } from "@/lib/session";
import type { Agency } from "@/lib/types";
import { RequireRole } from "../RequireRole";
import { EmptyState } from "../ui/Feedback";

interface DashboardContextValue {
  agencyId: string;
  session: Session;
  agency?: Agency;
  reloadAgency: () => void;
}

const DashboardContext = createContext<DashboardContextValue | null>(null);

/** The logged-in member's agency, available to every dashboard page. */
export function useDashboard(): DashboardContextValue {
  const ctx = useContext(DashboardContext);
  if (!ctx) throw new Error("useDashboard must be used inside the dashboard layout");
  return ctx;
}

const tabs = [
  { href: "/dashboard", label: "Overview" },
  { href: "/dashboard/tours", label: "Tours" },
  { href: "/dashboard/bookings", label: "Bookings" },
  { href: "/dashboard/members", label: "Team" },
  { href: "/dashboard/settings", label: "Settings" },
];

export function DashboardShell({ children }: { children: ReactNode }) {
  return <RequireRole roles={["member"]}>{(session) => <Shell session={session}>{children}</Shell>}</RequireRole>;
}

function Shell({ session, children }: { session: Session; children: ReactNode }) {
  const pathname = usePathname();
  const agencyId = session.agencyId;
  const agency = useAsync(() => (agencyId ? getAgency(agencyId) : Promise.resolve(undefined)), [agencyId]);

  if (!agencyId) return <EmptyState title="Your account isn't linked to an agency." description="Log out and back in, or contact your agency." />;

  return (
    <DashboardContext.Provider value={{ agencyId, session, agency: agency.data, reloadAgency: agency.reload }}>
      <div className="mb-6">
        <p className="text-sm text-slate-500">Agency dashboard</p>
        <h1 className="text-2xl font-semibold">{agency.data?.name ?? " "}</h1>
        <nav className="mt-4 flex flex-wrap gap-1 border-b border-slate-200" aria-label="Dashboard">
          {tabs.map((t) => {
            const active = t.href === "/dashboard" ? pathname === t.href : pathname.startsWith(t.href);
            return (
              <Link
                key={t.href}
                href={t.href}
                aria-current={active ? "page" : undefined}
                className={`-mb-px border-b-2 px-4 py-2 text-sm font-medium ${active ? "border-teal-700 text-teal-800" : "border-transparent text-slate-600 hover:text-slate-900"}`}
              >
                {t.label}
              </Link>
            );
          })}
        </nav>
      </div>
      {children}
    </DashboardContext.Provider>
  );
}

/** Shown instead of a section the member's role doesn't grant. */
export function NoAccess({ what, permission }: { what: string; permission: string }) {
  return (
    <EmptyState
      title={`You don't have access to ${what}.`}
      description={`Ask an agency manager to grant you the "${permission}" permission.`}
    />
  );
}
