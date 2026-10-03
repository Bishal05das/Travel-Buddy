"use client";

import Link from "next/link";
import { useAsync } from "@/hooks/useAsync";
import { CUSTOMER_URL, imageUrl } from "@/lib/config";
import { agencyBookings, listAgencyTours } from "@/lib/endpoints";
import { ButtonLink } from "../ui/Button";
import { Spinner } from "../ui/Feedback";
import { useDashboard } from "./DashboardShell";

export function Overview() {
  const { agencyId, agency } = useDashboard();
  const tours = useAsync(() => listAgencyTours(agencyId, 1, 1), [agencyId]);
  const pending = useAsync(() => agencyBookings(agencyId, { status: "pending", limit: 1 }), [agencyId]);
  const confirmed = useAsync(() => agencyBookings(agencyId, { status: "confirmed", limit: 1 }), [agencyId]);
  const agencyImage = imageUrl(agency?.image_path);

  return (
    <div className="space-y-6">
      {agency && (
        <section aria-labelledby="agency-image-heading" className="flex flex-wrap items-center gap-5 rounded-xl border border-slate-200 bg-white p-5 shadow-sm">
          {agencyImage ? (
            // eslint-disable-next-line @next/next/no-img-element -- uploaded image served by the API
            <img src={agencyImage} alt={`${agency.name} agency image`} className="h-28 w-40 rounded-lg bg-slate-100 object-cover" />
          ) : (
            <div className="flex h-28 w-40 items-center justify-center rounded-lg bg-slate-100 text-sm text-slate-500">No agency image</div>
          )}
          <div className="flex-1 space-y-2">
            <h2 id="agency-image-heading" className="text-lg font-semibold">Agency image</h2>
            <p className="text-sm text-slate-600">Choose the image customers see on your agency page.</p>
            <ButtonLink href="/dashboard/settings" variant="secondary">
              {agencyImage ? "Edit agency image" : "Set agency image"}
            </ButtonLink>
          </div>
        </section>
      )}
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
