"use client";

import { useState } from "react";
import { useAsync } from "@/hooks/useAsync";
import { cancelMyBooking, myBookings } from "@/lib/endpoints";
import type { Booking, BookingStatus } from "@/lib/types";
import { RequireRole } from "../RequireRole";
import { Button, ButtonLink } from "../ui/Button";
import { Alert, EmptyState, ErrorState, LoadingState } from "../ui/Feedback";
import { PageHeader } from "../ui/PageHeader";
import { Pagination } from "../ui/Pagination";
import { BookingCard, StatusTabs } from "./BookingCard";

export function MyBookings() {
  return <RequireRole roles={["user"]}>{() => <MyBookingsList />}</RequireRole>;
}

function canCancel(b: Booking) {
  return (b.status === "pending" || b.status === "confirmed") && new Date(b.tour_start_date).getTime() > Date.now();
}

function MyBookingsList() {
  const [status, setStatus] = useState<BookingStatus | "">("");
  const [page, setPage] = useState(1);
  const { data, error, loading, reload } = useAsync(() => myBookings({ status, page, limit: 10 }), [status, page]);
  const [busy, setBusy] = useState<string | null>(null);
  const [notice, setNotice] = useState<{ tone: "success" | "error"; text: string } | null>(null);

  async function cancel(b: Booking) {
    const verified = b.payment_status === "success";
    const msg = `Cancel your booking for ${b.tour_name}?${verified ? "\n\nYour payment was verified; contact the agency for the refund." : ""}`;
    if (!window.confirm(msg)) return;
    setBusy(b.booking_id);
    setNotice(null);
    try {
      await cancelMyBooking(b.booking_id);
      setNotice({ tone: "success", text: `Your booking for ${b.tour_name} was cancelled.` });
      reload();
    } catch (err) {
      setNotice({ tone: "error", text: err instanceof Error ? err.message : "Could not cancel the booking." });
    } finally {
      setBusy(null);
    }
  }

  return (
    <div>
      <PageHeader title="My bookings" description="Track your bookings, payments and trip status." />
      <div className="mb-4">
        <StatusTabs
          value={status}
          onChange={(v) => {
            setStatus(v);
            setPage(1);
          }}
        />
      </div>
      {notice && (
        <div className="mb-4">
          <Alert tone={notice.tone}>{notice.text}</Alert>
        </div>
      )}
      {loading && !data ? (
        <LoadingState label="Loading your bookings" />
      ) : error ? (
        <ErrorState message={error} onRetry={reload} />
      ) : data && data.Data.length === 0 ? (
        <EmptyState
          title={status ? `No ${status} bookings` : "You haven't booked any tours yet"}
          action={!status && <ButtonLink href="/search">Find a tour</ButtonLink>}
        />
      ) : data ? (
        <div className="space-y-3" aria-busy={loading}>
          {data.Data.map((b) => (
            <BookingCard
              key={b.booking_id}
              booking={b}
              actions={
                canCancel(b) && (
                  <Button variant="secondary" size="sm" loading={busy === b.booking_id} onClick={() => cancel(b)}>
                    Cancel booking
                  </Button>
                )
              }
            />
          ))}
          <Pagination page={page} totalPages={data.Meta.TotalPage} onChange={setPage} />
        </div>
      ) : null}
    </div>
  );
}
