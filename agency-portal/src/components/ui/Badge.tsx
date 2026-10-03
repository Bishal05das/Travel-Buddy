import type { BookingStatus, TourStatus } from "@/lib/types";

type Color = "green" | "amber" | "red" | "slate" | "sky";
const colors: Record<Color, string> = {
  green: "bg-emerald-100 text-emerald-800",
  amber: "bg-amber-100 text-amber-900",
  red: "bg-red-100 text-red-800",
  slate: "bg-slate-100 text-slate-700",
  sky: "bg-sky-100 text-sky-800",
};

export function Badge({ color, children }: { color: Color; children: React.ReactNode }) {
  return <span className={`inline-flex rounded-full px-2.5 py-0.5 text-xs font-medium capitalize ${colors[color]}`}>{children}</span>;
}

const bookingColors: Record<BookingStatus, Color> = { pending: "amber", confirmed: "green", cancelled: "red", completed: "sky" };
export function BookingStatusBadge({ status }: { status: BookingStatus }) {
  return <Badge color={bookingColors[status] ?? "slate"}>{status}</Badge>;
}

const tourColors: Record<TourStatus, Color> = { open: "green", closed: "slate", cancelled: "red" };
export function TourStatusBadge({ status }: { status: TourStatus }) {
  return <Badge color={tourColors[status] ?? "slate"}>{status}</Badge>;
}

export function PaymentBadge({ status }: { status: string }) {
  if (!status) return null;
  const color: Color = status === "success" ? "green" : status === "failed" ? "red" : "amber";
  return <Badge color={color}>{status === "success" ? "paid" : status === "pending" ? "payment pending" : "payment failed"}</Badge>;
}
