"use client";

import { useState } from "react";
import { useAsync } from "@/hooks/useAsync";
import { getAgency, listAgencyTours } from "@/lib/endpoints";
import { unitPrice } from "@/lib/format";
import { TourCard, TourGrid } from "./tours/TourCard";
import { TourImage } from "./tours/TourImage";
import { ButtonLink } from "./ui/Button";
import { EmptyState, ErrorState, LoadingState } from "./ui/Feedback";
import { Pagination } from "./ui/Pagination";

export function AgencyView({ id }: { id: string }) {
  const [page, setPage] = useState(1);
  const agency = useAsync(() => getAgency(id), [id]);
  const tours = useAsync(() => listAgencyTours(id, page, 9), [id, page]);

  if (agency.loading && !agency.data) return <LoadingState label="Loading agency" />;
  if (agency.error) {
    if (agency.status === 404 || agency.status === 400) {
      return <EmptyState title="Agency not found" action={<ButtonLink href="/">Back to tours</ButtonLink>} />;
    }
    return <ErrorState message={agency.error} onRetry={agency.reload} />;
  }
  const a = agency.data!;

  return (
    <div className="space-y-8">
      <header className="flex items-center gap-5">
        <TourImage path={a.image_path} alt={a.name} className="h-20 w-20 shrink-0 rounded-xl" />
        <div>
          <h1 className="text-2xl font-semibold">{a.name}</h1>
          <p className="text-slate-600">
            {a.address && <>{a.address} · </>}★ {a.rating.toFixed(1)}
          </p>
        </div>
      </header>

      <section>
        <h2 className="mb-4 text-lg font-semibold">Tours</h2>
        {tours.loading && !tours.data ? (
          <LoadingState label="Loading tours" />
        ) : tours.error ? (
          <ErrorState message={tours.error} onRetry={tours.reload} />
        ) : tours.data && tours.data.Data.length === 0 ? (
          <EmptyState title="This agency has no tours yet" />
        ) : tours.data ? (
          <div className="space-y-4" aria-busy={tours.loading}>
            <TourGrid>
              {tours.data.Data.map((t) => (
                <TourCard
                  key={t.tour_id}
                  tour={{
                    id: t.tour_id,
                    name: t.name,
                    startDate: t.start_date,
                    endDate: t.end_date,
                    price: unitPrice(t),
                    originalPrice: t.price,
                    seats: t.available_seat,
                    status: t.status,
                    imagePath: t.image_path,
                  }}
                />
              ))}
            </TourGrid>
            <Pagination page={page} totalPages={tours.data.Meta.TotalPage} onChange={setPage} />
          </div>
        ) : null}
      </section>
    </div>
  );
}
