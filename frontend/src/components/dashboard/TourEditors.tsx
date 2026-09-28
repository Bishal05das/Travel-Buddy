"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useAsync } from "@/hooks/useAsync";
import { createTour, getTour, updateTour } from "@/lib/endpoints";
import { formToUpdate, TourForm, tourToForm } from "../tours/TourForm";
import { ButtonLink } from "../ui/Button";
import { Alert, EmptyState, ErrorState, LoadingState } from "../ui/Feedback";
import { useDashboard } from "./DashboardShell";

function BackLink() {
  return (
    <Link href="/dashboard/tours" className="text-sm text-teal-700 hover:underline">
      ← Back to tours
    </Link>
  );
}

export function NewTour() {
  const { agencyId } = useDashboard();
  const router = useRouter();
  return (
    <div className="space-y-4">
      <BackLink />
      <h2 className="text-lg font-semibold">New tour</h2>
      <TourForm
        mode="create"
        onSubmit={async (v, image) => {
          const form = new FormData();
          const values = formToUpdate(v);
          form.set("name", values.name);
          form.set("description", values.description);
          form.set("start_date", v.start_date);
          form.set("end_date", v.end_date);
          form.set("last_enrollment_date", v.last_enrollment_date);
          form.set("total_seat", String(values.total_seat));
          form.set("price", String(values.price));
          form.set("discount", String(values.discount));
          if (image) form.set("image", image);
          await createTour(agencyId, form);
          router.push("/dashboard/tours");
        }}
      />
    </div>
  );
}

export function EditTour({ id }: { id: string }) {
  const { agencyId } = useDashboard();
  const router = useRouter();
  const { data: tour, error, loading, reload } = useAsync(() => getTour(id), [id]);

  if (loading && !tour) return <LoadingState label="Loading tour" />;
  if (error) return <ErrorState message={error} onRetry={reload} />;
  if (!tour || tour.agency_id !== agencyId) {
    return <EmptyState title="Tour not found in your agency" action={<ButtonLink href="/dashboard/tours">Back to tours</ButtonLink>} />;
  }
  if (tour.status === "cancelled") {
    return (
      <div className="space-y-4">
        <BackLink />
        <Alert tone="warning">This tour was cancelled and can no longer be edited.</Alert>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      <BackLink />
      <h2 className="text-lg font-semibold">Edit {tour.name}</h2>
      <TourForm
        mode="edit"
        initial={tourToForm(tour)}
        bookedSeats={tour.total_seat - tour.available_seat}
        onSubmit={async (v) => {
          await updateTour(agencyId, tour.tour_id, formToUpdate(v));
          router.push("/dashboard/tours");
        }}
      />
    </div>
  );
}
