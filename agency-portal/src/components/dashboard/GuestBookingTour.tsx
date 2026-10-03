"use client";

import { useAsync } from "@/hooks/useAsync";
import { getTour } from "@/lib/endpoints";
import { bookability, formatMoney, unitPrice } from "@/lib/format";
import type { Session } from "@/lib/session";
import { BookingForm } from "../bookings/BookingForm";
import { RequireRole } from "../RequireRole";
import { ButtonLink } from "../ui/Button";
import { Alert, EmptyState, ErrorState, LoadingState } from "../ui/Feedback";
import { Card } from "../ui/PageHeader";

export function GuestBookingTour({ id }: { id: string }) {
  return <RequireRole roles={["member"]}>{(session) => <BookingPage id={id} session={session} />}</RequireRole>;
}

function BookingPage({ id, session }: { id: string; session: Session }) {
  const tour = useAsync(() => getTour(id), [id]);
  if (tour.loading && !tour.data) return <LoadingState label="Loading tour" />;
  if (tour.error) return <ErrorState message={tour.error} onRetry={tour.reload} />;
  if (!tour.data || tour.data.agency_id !== session.agencyId) {
    return <EmptyState title="Tour not found in your agency" action={<ButtonLink href="/dashboard/tours">Back to tours</ButtonLink>} />;
  }

  const status = bookability(tour.data);
  return (
    <div className="max-w-xl space-y-5">
      <ButtonLink href="/dashboard/tours" variant="ghost">← Back to tours</ButtonLink>
      <div>
        <h1 className="text-2xl font-semibold">Book a guest on {tour.data.name}</h1>
        <p className="mt-1 text-sm text-slate-600">
          {formatMoney(unitPrice(tour.data))} per person · {tour.data.available_seat} seats left
        </p>
      </div>
      <Card>
        {status.ok ? <BookingForm tour={tour.data} onBooked={tour.reload} /> : <Alert tone="warning">{status.reason}</Alert>}
      </Card>
    </div>
  );
}
