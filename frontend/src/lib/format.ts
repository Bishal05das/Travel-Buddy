import type { Tour } from "./types";

const money = new Intl.NumberFormat("en-US", { maximumFractionDigits: 0 });

export function formatMoney(amount: number): string {
  return `৳${money.format(amount)}`;
}

/** Formats an API date (UTC midnight for tour dates) as e.g. "10 Feb 2031". */
export function formatDate(value: string): string {
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleDateString("en-GB", { day: "numeric", month: "short", year: "numeric", timeZone: "UTC" });
}

export function formatDateTime(value: string): string {
  const d = new Date(value);
  if (Number.isNaN(d.getTime())) return "—";
  return d.toLocaleString("en-GB", { day: "numeric", month: "short", year: "numeric", hour: "2-digit", minute: "2-digit" });
}

/** "2031-02-10T00:00:00Z" -> "2031-02-10" for <input type="date">. */
export function toDateInput(value: string): string {
  return value ? value.slice(0, 10) : "";
}

/** "2031-02-10" -> "2031-02-10T00:00:00Z", the format the API expects in JSON. */
export function fromDateInput(value: string): string {
  return `${value}T00:00:00Z`;
}

/**
 * Per-person price after discount. Mirrors the API's integer formula
 * (price - price*discount/100 with integer division) so the total we send
 * matches the server's check exactly.
 */
export function unitPrice(tour: Pick<Tour, "price" | "discount">): number {
  return tour.price - Math.floor((tour.price * tour.discount) / 100);
}

/** Whether a tour can take bookings right now, and if not, why. */
export function bookability(tour: Tour): { ok: true } | { ok: false; reason: string } {
  if (tour.status === "cancelled") return { ok: false, reason: "This tour was cancelled by the agency." };
  if (tour.status === "closed") return { ok: false, reason: "The agency has closed bookings for this tour." };
  if (tour.available_seat <= 0) return { ok: false, reason: "This tour is sold out." };
  // The API closes enrollment at 00:00 UTC on the last enrollment date.
  if (new Date(tour.last_enrollment_date).getTime() < Date.now()) {
    return { ok: false, reason: "The enrollment deadline has passed." };
  }
  return { ok: true };
}
