import { BookingsManager } from "@/components/dashboard/BookingsManager";
import type { BookingStatus } from "@/lib/types";

const statuses = ["pending", "confirmed", "cancelled", "completed"];

export default async function DashboardBookingsPage(props: PageProps<"/dashboard/bookings">) {
  const sp = await props.searchParams;
  const status = typeof sp.status === "string" && statuses.includes(sp.status) ? (sp.status as BookingStatus) : "";
  return <BookingsManager initialStatus={status} />;
}
