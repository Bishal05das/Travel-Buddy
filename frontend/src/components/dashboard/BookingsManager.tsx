"use client";

import { useState } from "react";
import { useAsync } from "@/hooks/useAsync";
import { agencyBookings, listAgencyTours, setBookingStatus } from "@/lib/endpoints";
import type { Booking, BookingStatus } from "@/lib/types";
import { BookingCard, StatusTabs } from "../bookings/BookingCard";
import { Button } from "../ui/Button";
import { Alert, EmptyState, ErrorState, LoadingState } from "../ui/Feedback";
import { Select } from "../ui/Field";
import { Pagination } from "../ui/Pagination";
import { NoAccess, useDashboard } from "./DashboardShell";

type Action = "confirmed" | "cancelled" | "completed";

export function BookingsManager({ initialStatus }: { initialStatus: BookingStatus | "" }) {
  const { agencyId } = useDashboard();
  const [status, setStatus] = useState<BookingStatus | "">(initialStatus);
  const [tourId, setTourId] = useState("");
  const [page, setPage] = useState(1);
  const bookings = useAsync(() => agencyBookings(agencyId, { status, tour_id: tourId, page, limit: 10 }), [agencyId, status, tourId, page]);
  const tours = useAsync(() => listAgencyTours(agencyId, 1, 100), [agencyId]);
  const [busy, setBusy] = useState<string | null>(null);
  const [notice, setNotice] = useState<{ tone: "success" | "error"; text: string } | null>(null);

  async function act(b: Booking, to: Action) {
    const prompts: Record<Action, string> = {
      confirmed: `Confirm this booking? Only confirm after you've received ${b.payment_method} payment with transaction ID ${b.transaction_id}.`,
      cancelled: `Cancel ${b.customer_name}'s booking? The seats return to the tour.${b.payment_status === "success" ? " The payment was verified, so you'll need to refund it." : ""}`,
      completed: "Mark this booking as completed?",
    };
    if (!window.confirm(prompts[to])) return;
    setBusy(b.booking_id);
    setNotice(null);
    try {
      await setBookingStatus(agencyId, b.booking_id, to);
      setNotice({ tone: "success", text: `Booking for ${b.customer_name} is now ${to}.` });
      bookings.reload();
    } catch (err) {
      setNotice({ tone: "error", text: err instanceof Error ? err.message : "The update failed." });
    } finally {
      setBusy(null);
    }
  }

  if (bookings.status === 403) return <NoAccess what="bookings" permission="booking:read" />;

  return (
    <div className="space-y-4">
      <div className="flex flex-wrap items-end justify-between gap-4">
        <StatusTabs
          value={status}
          onChange={(v) => {
            setStatus(v);
            setPage(1);
          }}
        />
        <div className="w-64">
          <Select
            label="Tour"
            value={tourId}
            onChange={(e) => {
              setTourId(e.target.value);
              setPage(1);
            }}
          >
            <option value="">All tours</option>
            {tours.data?.Data.map((t) => (
              <option key={t.tour_id} value={t.tour_id}>
                {t.name}
              </option>
            ))}
          </Select>
        </div>
      </div>
      {notice && <Alert tone={notice.tone}>{notice.text}</Alert>}
      {bookings.loading && !bookings.data ? (
        <LoadingState label="Loading bookings" />
      ) : bookings.error ? (
        <ErrorState message={bookings.error} onRetry={bookings.reload} />
      ) : bookings.data && bookings.data.Data.length === 0 ? (
        <EmptyState title={status || tourId ? "No bookings match these filters" : "No bookings yet"} />
      ) : bookings.data ? (
        <div className="space-y-3" aria-busy={bookings.loading}>
          {bookings.data.Data.map((b) => (
            <BookingCard
              key={b.booking_id}
              booking={b}
              showCustomer
              actions={
                (b.status === "pending" || b.status === "confirmed") && (
                  <>
                    {b.status === "pending" && (
                      <Button size="sm" loading={busy === b.booking_id} onClick={() => act(b, "confirmed")}>
                        Confirm payment
                      </Button>
                    )}
                    {b.status === "confirmed" && (
                      <Button size="sm" variant="secondary" loading={busy === b.booking_id} onClick={() => act(b, "completed")}>
                        Mark completed
                      </Button>
                    )}
                    <Button size="sm" variant="danger-ghost" disabled={busy === b.booking_id} onClick={() => act(b, "cancelled")}>
                      Cancel booking
                    </Button>
                  </>
                )
              }
            />
          ))}
          <Pagination page={page} totalPages={bookings.data.Meta.TotalPage} onChange={setPage} />
        </div>
      ) : null}
    </div>
  );
}
