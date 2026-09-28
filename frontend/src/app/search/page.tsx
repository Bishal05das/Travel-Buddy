import type { Metadata } from "next";
import { SearchView } from "@/components/SearchView";

export const metadata: Metadata = { title: "Search" };

function first(v: string | string[] | undefined): string {
  return (Array.isArray(v) ? v[0] : v) ?? "";
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
  // key: remount the form when the URL changes (e.g. from the navbar search).
  return <SearchView key={JSON.stringify(filters)} filters={filters} />;
}
