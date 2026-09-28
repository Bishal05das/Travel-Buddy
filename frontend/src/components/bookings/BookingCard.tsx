import Link from "next/link";
import type { ReactNode } from "react";
import { formatDate, formatDateTime, formatMoney } from "@/lib/format";
import type { Booking, BookingStatus } from "@/lib/types";
import { BookingStatusBadge, PaymentBadge, TourStatusBadge } from "../ui/Badge";

const reasons: Record<NonNullable<Booking["cancellation_reason"]>, string> = {
  customer: "Cancelled by the customer.",
  agency: "Cancelled by the agency.",
  tour_cancelled: "The agency cancelled this tour, so the booking was cancelled.",
};

export function BookingCard({ booking: b, showCustomer = false, actions }: { booking: Booking; showCustomer?: boolean; actions?: ReactNode }) {
  const refundDue = b.status === "cancelled" && b.payment_status === "success";
  return (
    <article className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm">
      <div className="flex flex-wrap items-start justify-between gap-3">
        <div>
          <Link href={`/tours/${b.tour_id}`} className="font-semibold text-slate-900 hover:text-teal-800">
            {b.tour_name}
          </Link>
          <p className="text-sm text-slate-500">
            {b.agency_name} · starts {formatDate(b.tour_start_date)}
          </p>
        </div>
        <div className="flex flex-wrap gap-2">
          <BookingStatusBadge status={b.status} />
          <PaymentBadge status={b.payment_status} />
          {b.tour_status !== "open" && (
            <span className="inline-flex items-center gap-1 text-xs text-slate-500">
              tour <TourStatusBadge status={b.tour_status} />
            </span>
          )}
        </div>
      </div>

      <dl className="mt-3 grid grid-cols-2 gap-x-4 gap-y-1 text-sm sm:grid-cols-4">
        <div>
          <dt className="text-slate-500">People</dt>
          <dd className="font-medium">{b.number_of_people}</dd>
        </div>
        <div>
          <dt className="text-slate-500">Total</dt>
          <dd className="font-medium">{formatMoney(b.total_price)}</dd>
        </div>
        <div>
          <dt className="text-slate-500">Payment</dt>
          <dd className="font-medium capitalize">
            {b.payment_method} · <span className="font-mono text-xs">{b.transaction_id}</span>
          </dd>
        </div>
        <div>
          <dt className="text-slate-500">Booked</dt>
          <dd className="font-medium">{formatDateTime(b.booking_date)}</dd>
        </div>
        {showCustomer && (
          <div className="col-span-2 sm:col-span-4">
            <dt className="text-slate-500">Customer {b.created_by === "agency_member" && "(booked at the agency)"}</dt>
            <dd className="font-medium">
              {b.customer_name} · {b.customer_email} · {b.customer_phone}
            </dd>
          </div>
        )}
      </dl>

      {b.status === "cancelled" && b.cancellation_reason && (
        <p className={`mt-3 text-sm ${b.cancellation_reason === "tour_cancelled" ? "font-medium text-red-700" : "text-slate-600"}`}>
          {reasons[b.cancellation_reason]}
          {refundDue && " The payment was verified, so a refund is due from the agency."}
        </p>
      )}
      {b.status === "pending" && !showCustomer && (
        <p className="mt-3 text-sm text-slate-600">Waiting for the agency to verify your payment.</p>
      )}

      {actions && <div className="mt-4 flex flex-wrap gap-2 border-t border-slate-100 pt-3">{actions}</div>}
    </article>
  );
}

const tabs: { value: BookingStatus | ""; label: string }[] = [
  { value: "", label: "All" },
  { value: "pending", label: "Pending" },
  { value: "confirmed", label: "Confirmed" },
  { value: "completed", label: "Completed" },
  { value: "cancelled", label: "Cancelled" },
];

export function StatusTabs({ value, onChange }: { value: BookingStatus | ""; onChange: (v: BookingStatus | "") => void }) {
  return (
    <div role="tablist" aria-label="Filter by status" className="flex flex-wrap gap-1">
      {tabs.map((t) => (
        <button
          key={t.value}
          role="tab"
          aria-selected={value === t.value}
          onClick={() => onChange(t.value)}
          className={`rounded-full px-3 py-1 text-sm font-medium ${value === t.value ? "bg-teal-700 text-white" : "bg-white text-slate-600 ring-1 ring-slate-200 hover:bg-slate-50"}`}
        >
          {t.label}
        </button>
      ))}
    </div>
  );
}
