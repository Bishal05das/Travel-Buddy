"use client";

import Link from "next/link";
import { useAsync } from "@/hooks/useAsync";
import { getHome } from "@/lib/endpoints";
import { SearchBox } from "./SearchBox";
import { TourCard, TourGrid } from "./tours/TourCard";
import { EmptyState, ErrorState, LoadingState } from "./ui/Feedback";

export function HomeView() {
  const { data, error, loading, reload } = useAsync(getHome, []);

  return (
    <div className="space-y-10">
      <section className="rounded-2xl bg-gradient-to-br from-teal-700 to-sky-800 px-6 py-12 text-white sm:px-10">
        <h1 className="text-3xl font-semibold sm:text-4xl">Find your next trip</h1>
        <p className="mt-2 max-w-xl text-teal-50">Browse tours from trusted travel agencies and book in a few clicks.</p>
        <div className="mt-6">
          <SearchBox />
        </div>
      </section>

      {loading && !data ? (
        <LoadingState label="Loading tours" />
      ) : error ? (
        <ErrorState message={error} onRetry={reload} />
      ) : data ? (
        <>
          <section>
            <h2 className="mb-4 text-xl font-semibold">Popular tours</h2>
            {data.top_tours.length === 0 ? (
              <EmptyState title="No open tours right now" description="Check back soon, agencies add new trips regularly." />
            ) : (
              <TourGrid>
                {data.top_tours.map((t) => (
                  <TourCard
                    key={t.tour_id}
                    tour={{
                      id: t.tour_id,
                      name: t.name,
                      agencyName: t.agency_name,
                      startDate: t.start_date,
                      endDate: t.end_date,
                      price: t.final_price,
                      originalPrice: t.price,
                      seats: t.available_seat,
                      imagePath: t.image_path,
                    }}
                  />
                ))}
              </TourGrid>
            )}
          </section>

          <section>
            <h2 className="mb-4 text-xl font-semibold">Top agencies</h2>
            {data.top_agencies.length === 0 ? (
              <EmptyState title="No agencies yet" />
            ) : (
              <div className="grid gap-3 sm:grid-cols-2 lg:grid-cols-3">
                {data.top_agencies.map((a) => (
                  <Link
                    key={a.agency_id}
                    href={`/agencies/${a.agency_id}`}
                    className="rounded-xl border border-slate-200 bg-white p-4 shadow-sm transition hover:shadow-md"
                  >
                    <p className="font-semibold text-slate-900">{a.name}</p>
                    {a.address && <p className="text-sm text-slate-500">{a.address}</p>}
                    <p className="mt-2 text-sm text-slate-600">
                      ★ {a.rating.toFixed(1)} · {a.total_tours} open {a.total_tours === 1 ? "tour" : "tours"}
                    </p>
                  </Link>
                ))}
              </div>
            )}
          </section>
        </>
      ) : null}
    </div>
  );
}
