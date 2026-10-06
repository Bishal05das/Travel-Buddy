import type { Metadata } from "next";
import { SearchView } from "@/components/SearchView";

export const metadata: Metadata = { title: "Search" };

function first(v: string | string[] | undefined): string {
  return (Array.isArray(v) ? v[0] : v) ?? "";
}

function positiveInteger(v: string | string[] | undefined, fallback: number): number {
  const value = first(v);
  const number = Number(value);
  return /^\d+$/.test(value) && Number.isSafeInteger(number) && number > 0 ? number : fallback;
}

export default async function SearchPage(props: PageProps<"/search">) {
  const sp = await props.searchParams;
  const filters = {
    q: first(sp.q),
    min_price: first(sp.min_price),
    max_price: first(sp.max_price),
    start_date: first(sp.start_date),
    end_date: first(sp.end_date),
  };
  const page = positiveInteger(sp.page, 1);
  const limit = Math.min(positiveInteger(sp.limit, 20), 50);
  // key: remount the form when the URL changes (e.g. from the navbar search).
  return <SearchView key={JSON.stringify({ ...filters, page, limit })} filters={filters} page={page} limit={limit} />;
}
