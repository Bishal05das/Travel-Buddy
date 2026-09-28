"use client";

import Link from "next/link";
import { useState } from "react";
import { useAsync } from "@/hooks/useAsync";
import { deleteTour, listAgencyTours, setTourStatus } from "@/lib/endpoints";
import { formatDate, formatMoney, unitPrice } from "@/lib/format";
import type { Tour, TourStatus } from "@/lib/types";
import { TourStatusBadge } from "../ui/Badge";
import { Button, ButtonLink } from "../ui/Button";
import { Alert, EmptyState, ErrorState, LoadingState } from "../ui/Feedback";
import { Pagination } from "../ui/Pagination";
import { useDashboard } from "./DashboardShell";

export function ToursManager() {
  const { agencyId } = useDashboard();
  const [page, setPage] = useState(1);
  const { data, error, loading, reload } = useAsync(() => listAgencyTours(agencyId, page, 10), [agencyId, page]);
  const [busy, setBusy] = useState<string | null>(null);
  const [notice, setNotice] = useState<{ tone: "success" | "error"; text: string } | null>(null);

  async function run(tour: Tour, action: () => Promise<string>) {
    setBusy(tour.tour_id);
    setNotice(null);
    try {
      setNotice({ tone: "success", text: await action() });
      reload();
    } catch (err) {
      setNotice({ tone: "error", text: err instanceof Error ? err.message : "The action failed." });
    } finally {
      setBusy(null);
    }
  }

  function changeStatus(tour: Tour, status: TourStatus) {
    if (status === "cancelled") {
      const booked = tour.total_seat - tour.available_seat;
      const ok = window.confirm(
        `Cancel "${tour.name}"?\n\nThis can't be undone. ${
          booked > 0 ? `All active bookings (${booked} seats) will be cancelled and customers will see the tour as cancelled. ` : ""
        }Verified payments must be refunded outside Travel Buddy.`,
      );
      if (!ok) return;
    }
    run(tour, async () => {
      const r = await setTourStatus(tour.tour_id, status);
      if (status === "cancelled") {
        return `"${tour.name}" was cancelled${r.cancelled_bookings ? ` along with ${r.cancelled_bookings} booking${r.cancelled_bookings === 1 ? "" : "s"}` : ""}.`;
      }
      return `"${tour.name}" is now ${status}.`;
    });
  }

  function remove(tour: Tour) {
    if (!window.confirm(`Delete "${tour.name}" permanently? Its bookings will be deleted too.`)) return;
    run(tour, async () => {
      await deleteTour(tour.tour_id);
      return `"${tour.name}" was deleted.`;
    });
  }

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-lg font-semibold">Tours</h2>
        <ButtonLink href="/dashboard/tours/new">New tour</ButtonLink>
      </div>
      {notice && <Alert tone={notice.tone}>{notice.text}</Alert>}
      {loading && !data ? (
        <LoadingState label="Loading tours" />
      ) : error ? (
        <ErrorState message={error} onRetry={reload} />
      ) : data && data.Data.length === 0 ? (
        <EmptyState title="No tours yet" description="Create your first tour to start taking bookings." action={<ButtonLink href="/dashboard/tours/new">Create a tour</ButtonLink>} />
      ) : data ? (
        <>
          <div className="overflow-x-auto rounded-xl border border-slate-200 bg-white" aria-busy={loading}>
            <table className="w-full text-left text-sm">
              <thead className="bg-slate-50 text-xs uppercase text-slate-500">
                <tr>
                  <th className="px-4 py-3">Tour</th>
                  <th className="px-4 py-3">Dates</th>
                  <th className="px-4 py-3">Seats</th>
                  <th className="px-4 py-3">Price</th>
                  <th className="px-4 py-3">Status</th>
                  <th className="px-4 py-3 text-right">Actions</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-100">
                {data.Data.map((t) => (
                  <tr key={t.tour_id}>
                    <td className="px-4 py-3 font-medium">
                      <Link href={`/tours/${t.tour_id}`} className="hover:text-teal-800">
                        {t.name}
                      </Link>
                    </td>
                    <td className="whitespace-nowrap px-4 py-3 text-slate-600">
                      {formatDate(t.start_date)} – {formatDate(t.end_date)}
                    </td>
                    <td className="px-4 py-3 text-slate-600">
                      {t.total_seat - t.available_seat} / {t.total_seat} booked
                    </td>
                    <td className="px-4 py-3 text-slate-600">
                      {formatMoney(unitPrice(t))}
                      {t.discount > 0 && <span className="text-xs text-slate-400"> (−{t.discount}%)</span>}
                    </td>
                    <td className="px-4 py-3">
                      <TourStatusBadge status={t.status} />
                    </td>
                    <td className="px-4 py-3">
                      <div className="flex justify-end gap-1">
                        {t.status === "cancelled" ? (
                          <span className="text-xs text-slate-400">Cancelled tours can&apos;t be changed</span>
                        ) : (
                          <>
                            <ButtonLink href={`/dashboard/tours/${t.tour_id}/edit`} variant="ghost" size="sm">
                              Edit
                            </ButtonLink>
                            {t.status === "open" ? (
                              <Button variant="ghost" size="sm" disabled={busy === t.tour_id} onClick={() => changeStatus(t, "closed")}>
                                Close
                              </Button>
                            ) : (
                              <Button variant="ghost" size="sm" disabled={busy === t.tour_id} onClick={() => changeStatus(t, "open")}>
                                Reopen
                              </Button>
                            )}
                            <Button variant="danger-ghost" size="sm" disabled={busy === t.tour_id} onClick={() => changeStatus(t, "cancelled")}>
                              Cancel tour
                            </Button>
                          </>
                        )}
                        <Button variant="danger-ghost" size="sm" disabled={busy === t.tour_id} onClick={() => remove(t)}>
                          Delete
                        </Button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          <Pagination page={page} totalPages={data.Meta.TotalPage} onChange={setPage} />
        </>
      ) : null}
    </div>
  );
}
