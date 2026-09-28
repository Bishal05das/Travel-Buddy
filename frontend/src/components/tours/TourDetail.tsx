"use client";

import Link from "next/link";
import { useAsync } from "@/hooks/useAsync";
import { useSession } from "@/lib/auth";
import { getAgency, getTour } from "@/lib/endpoints";
import { bookability, formatDate, formatMoney, unitPrice } from "@/lib/format";
import type { Tour } from "@/lib/types";
import { BookingForm } from "../bookings/BookingForm";
import { TourStatusBadge } from "../ui/Badge";
import { ButtonLink } from "../ui/Button";
import { Alert, EmptyState, ErrorState, LoadingState } from "../ui/Feedback";
import { Card } from "../ui/PageHeader";
import { TourImage } from "./TourImage";

export function TourDetail({ id }: { id: string }) {
  const tour = useAsync(() => getTour(id), [id]);
  const agencyId = tour.data?.agency_id;
  const agency = useAsync(() => (agencyId ? getAgency(agencyId) : Promise.resolve(undefined)), [agencyId]);

  if (tour.loading && !tour.data) return <LoadingState label="Loading tour" />;
  if (tour.error) {
    if (/not found|invalid/i.test(tour.error)) {
      return <EmptyState title="Tour not found" description="It may have been removed by the agency." action={<ButtonLink href="/search">Browse tours</ButtonLink>} />;
    }
    return <ErrorState message={tour.error} onRetry={tour.reload} />;
  }
  const t = tour.data!;
  const discounted = t.discount > 0;

  return (
    <div className="grid gap-8 lg:grid-cols-[1fr_380px]">
      <article className="space-y-6">
        <TourImage path={t.image_path} alt={t.name} className="h-64 w-full rounded-2xl sm:h-80" />
        <div>
          <div className="flex flex-wrap items-center gap-3">
            <h1 className="text-3xl font-semibold text-slate-900">{t.name}</h1>
            <TourStatusBadge status={t.status} />
          </div>
          {agency.data && (
            <p className="mt-1 text-slate-600">
              by{" "}
              <Link href={`/agencies/${t.agency_id}`} className="font-medium text-teal-700 hover:underline">
                {agency.data.name}
              </Link>{" "}
              · ★ {agency.data.rating.toFixed(1)}
            </p>
          )}
        </div>
        <dl className="grid grid-cols-2 gap-4 rounded-xl border border-slate-200 bg-white p-5 sm:grid-cols-4">
          <Fact label="Starts" value={formatDate(t.start_date)} />
          <Fact label="Ends" value={formatDate(t.end_date)} />
          <Fact label="Book by" value={formatDate(t.last_enrollment_date)} />
          <Fact label="Seats left" value={`${t.available_seat} of ${t.total_seat}`} />
        </dl>
        <section>
          <h2 className="mb-2 text-lg font-semibold">About this tour</h2>
          <p className="whitespace-pre-line text-slate-700">{t.description}</p>
        </section>
      </article>

      <aside>
        <Card className="space-y-4 lg:sticky lg:top-6">
          <div>
            <span className="text-2xl font-semibold">{formatMoney(unitPrice(t))}</span>
            <span className="text-sm text-slate-500"> / person</span>
            {discounted && (
              <p className="text-sm text-slate-500">
                <span className="line-through">{formatMoney(t.price)}</span> · {t.discount}% off
              </p>
            )}
          </div>
          <BookingPanel tour={t} onBooked={tour.reload} />
        </Card>
      </aside>
    </div>
  );
}

function Fact({ label, value }: { label: string; value: string }) {
  return (
    <div>
      <dt className="text-xs uppercase tracking-wide text-slate-500">{label}</dt>
      <dd className="mt-1 font-medium text-slate-900">{value}</dd>
    </div>
  );
}

function BookingPanel({ tour, onBooked }: { tour: Tour; onBooked: () => void }) {
  const { session, ready } = useSession();
  const status = bookability(tour);

  if (!status.ok) return <Alert tone="warning">{status.reason}</Alert>;
  if (!ready) return null;
  if (!session) {
    return (
      <div className="space-y-2">
        <ButtonLink href={`/login?next=/tours/${tour.tour_id}`} className="w-full">
          Log in to book
        </ButtonLink>
        <p className="text-center text-sm text-slate-500">
          New here?{" "}
          <Link href={`/register?next=/tours/${tour.tour_id}`} className="text-teal-700 hover:underline">
            Create an account
          </Link>
        </p>
      </div>
    );
  }
  if (session.role === "user") return <BookingForm tour={tour} onBooked={onBooked} />;
  if (session.role === "member" && session.agencyId === tour.agency_id) {
    return (
      <div className="space-y-3">
        <p className="text-sm text-slate-600">Book this tour for a walk-in customer.</p>
        <BookingForm tour={tour} guest onBooked={onBooked} />
      </div>
    );
  }
  return <Alert tone="info">Log in with a customer account to book this tour.</Alert>;
}
