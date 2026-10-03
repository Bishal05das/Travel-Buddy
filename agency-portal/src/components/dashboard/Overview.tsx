"use client";

import Link from "next/link";
import { useAsync } from "@/hooks/useAsync";
import { CUSTOMER_URL } from "@/lib/config";
import { agencyBookings, listAgencyTours } from "@/lib/endpoints";
import { ButtonLink } from "../ui/Button";
import { Spinner } from "../ui/Feedback";
import { useDashboard } from "./DashboardShell";

export function Overview() {
  const { agencyId, agency } = useDashboard();
  const tours = useAsync(() => listAgencyTours(agencyId, 1, 1), [agencyId]);
  const pending = useAsync(() => agencyBookings(agencyId, { status: "pending", limit: 1 }), [agencyId]);
  const confirmed = useAsync(() => agencyBookings(agencyId, { status: "confirmed", limit: 1 }), [agencyId]);

  return (
    <div className="space-y-6">
      <div className="grid gap-4 sm:grid-cols-3">
        <Stat label="Tours" state={tours} href="/dashboard/tours" />
        <Stat label="Bookings awaiting payment check" state={pending} href="/dashboard/bookings?status=pending" highlight />
        <Stat label="Confirmed bookings" state={confirmed} href="/dashboard/bookings?status=confirmed" />
      </div>
      <div className="flex flex-wrap gap-2">
        <ButtonLink href="/dashboard/tours/new">Create a tour</ButtonLink>
        <ButtonLink href={`${CUSTOMER_URL}/agencies/${agencyId}`} variant="secondary">
          View public agency page
        </ButtonLink>
      </div>
      {agency && (
        <p className="text-sm text-slate-500">
          {agency.address} · registration {agency.reg_id} · rating ★ {agency.rating.toFixed(1)}
        </p>
      )}
    </div>
  );
}

function Stat({
  label,
  state,
  href,
  highlight,
}: {
  label: string;
  state: { data?: { Meta: { TotalCount: number } }; error: string | null; status: number | null; loading: boolean };
  href: string;
  highlight?: boolean;
}) {
  const count = state.data?.Meta.TotalCount;
  return (
    <Link href={href} className="rounded-xl border border-slate-200 bg-white p-5 shadow-sm hover:shadow-md">
      <p className="text-sm text-slate-500">{label}</p>
      <div className="mt-2 text-3xl font-semibold">
        {state.loading && count === undefined ? (
          <Spinner label="" />
        ) : state.status === 403 ? (
          <span className="text-sm font-normal text-slate-400">No access</span>
        ) : state.error ? (
          <span className="text-sm font-normal text-red-600">Unavailable</span>
        ) : (
          <span className={highlight && count ? "text-amber-600" : ""}>{count}</span>
        )}
      </div>
    </Link>
  );
}
