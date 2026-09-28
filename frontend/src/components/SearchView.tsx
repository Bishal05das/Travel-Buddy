"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { useState } from "react";
import { useAsync } from "@/hooks/useAsync";
import { searchTours, type SearchParams } from "@/lib/endpoints";
import { unitPrice } from "@/lib/format";
import { TourCard, TourGrid } from "./tours/TourCard";
import { Button } from "./ui/Button";
import { EmptyState, ErrorState, LoadingState } from "./ui/Feedback";
import { Input } from "./ui/Field";
import { PageHeader } from "./ui/PageHeader";

type Filters = Required<SearchParams>;

export function SearchView({ filters }: { filters: Filters }) {
  const router = useRouter();
  const [form, setForm] = useState<Filters>(filters);
  const [formError, setFormError] = useState<string | null>(null);
  const { data, error, loading, reload } = useAsync(() => searchTours(filters), [JSON.stringify(filters)]);

  function apply(e: React.FormEvent) {
    e.preventDefault();
    if (form.min_price && form.max_price && Number(form.min_price) > Number(form.max_price)) {
      setFormError("Minimum price can't be higher than maximum price.");
      return;
    }
    if (form.start_date && form.end_date && form.start_date > form.end_date) {
      setFormError("'From' date must be before 'to' date.");
      return;
    }
    setFormError(null);
    const q = new URLSearchParams(Object.entries(form).filter(([, v]) => v.trim() !== ""));
    router.push(`/search${q.size ? `?${q}` : ""}`);
  }

  const set = (key: keyof Filters) => (e: React.ChangeEvent<HTMLInputElement>) => setForm({ ...form, [key]: e.target.value });
  const hasFilters = Object.values(filters).some(Boolean);

  return (
    <div>
      <PageHeader title="Search tours" description="Filter by name, destination, price and dates." />
      <form onSubmit={apply} className="mb-8 grid gap-4 rounded-xl border border-slate-200 bg-white p-4 sm:grid-cols-2 lg:grid-cols-6">
        <div className="lg:col-span-2">
          <Input label="Keyword" value={form.q} onChange={set("q")} placeholder="Tour, place or agency" />
        </div>
        <Input label="Min price (৳)" type="number" min={0} value={form.min_price} onChange={set("min_price")} />
        <Input label="Max price (৳)" type="number" min={0} value={form.max_price} onChange={set("max_price")} />
        <Input label="Starts from" type="date" value={form.start_date} onChange={set("start_date")} />
        <Input label="Ends by" type="date" value={form.end_date} onChange={set("end_date")} />
        <div className="flex items-end gap-2 sm:col-span-2 lg:col-span-6">
          <Button type="submit">Apply filters</Button>
          {hasFilters && (
            <Button type="button" variant="ghost" onClick={() => router.push("/search")}>
              Clear
            </Button>
          )}
          {formError && <p className="text-sm text-red-600">{formError}</p>}
        </div>
      </form>

      {loading && !data ? (
        <LoadingState label="Searching" />
      ) : error ? (
        <ErrorState message={error} onRetry={reload} />
      ) : data ? (
        <div className="space-y-8" aria-busy={loading}>
          {data.Agencies.length > 0 && filters.q && (
            <section>
              <h2 className="mb-3 text-lg font-semibold">Agencies</h2>
              <div className="flex flex-wrap gap-2">
                {data.Agencies.map((a) => (
                  <Link
                    key={a.agency_id}
                    href={`/agencies/${a.agency_id}`}
                    className="rounded-full border border-slate-300 bg-white px-4 py-1.5 text-sm hover:border-teal-600"
                  >
                    {a.name} · ★ {a.rating.toFixed(1)}
                  </Link>
                ))}
              </div>
            </section>
          )}
          <section>
            <h2 className="mb-3 text-lg font-semibold">
              Tours <span className="text-sm font-normal text-slate-500">({data.Tours.length})</span>
            </h2>
            {data.Tours.length === 0 ? (
              <EmptyState
                title="No tours match your search"
                description={hasFilters ? "Try a different keyword or widen the price and date range." : undefined}
              />
            ) : (
              <TourGrid>
                {data.Tours.map((t) => (
                  <TourCard
                    key={t.tour_id}
                    tour={{
                      id: t.tour_id,
                      name: t.name,
                      agencyName: t.agency_name,
                      startDate: t.start_date,
                      endDate: t.end_date,
                      price: unitPrice(t),
                      originalPrice: t.price,
                      seats: t.available_seat,
                      status: t.status,
                    }}
                  />
                ))}
              </TourGrid>
            )}
          </section>
        </div>
      ) : null}
    </div>
  );
}
