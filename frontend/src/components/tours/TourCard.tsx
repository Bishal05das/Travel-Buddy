import Link from "next/link";
import { formatDate, formatMoney } from "@/lib/format";
import type { TourStatus } from "@/lib/types";
import { TourStatusBadge } from "../ui/Badge";
import { TourImage } from "./TourImage";

export interface TourCardData {
  id: string;
  name: string;
  agencyName?: string;
  startDate: string;
  endDate: string;
  price: number; // per person, after discount
  originalPrice?: number; // before discount, shown struck through when higher
  seats: number;
  status?: TourStatus;
  imagePath?: string;
}

export function TourCard({ tour }: { tour: TourCardData }) {
  const discounted = tour.originalPrice !== undefined && tour.originalPrice > tour.price;
  return (
    <Link
      href={`/tours/${tour.id}`}
      className="group flex flex-col overflow-hidden rounded-xl border border-slate-200 bg-white shadow-sm transition hover:shadow-md"
    >
      <TourImage path={tour.imagePath} alt={tour.name} className="h-36 w-full" />
      <div className="flex flex-1 flex-col gap-2 p-4">
        <div className="flex items-start justify-between gap-2">
          <h3 className="font-semibold text-slate-900 group-hover:text-teal-800">{tour.name}</h3>
          {tour.status && tour.status !== "open" && <TourStatusBadge status={tour.status} />}
        </div>
        {tour.agencyName && <p className="text-sm text-slate-500">by {tour.agencyName}</p>}
        <p className="text-sm text-slate-600">
          {formatDate(tour.startDate)} – {formatDate(tour.endDate)}
        </p>
        <div className="mt-auto flex items-end justify-between pt-2">
          <div>
            <span className="text-lg font-semibold text-slate-900">{formatMoney(tour.price)}</span>
            {discounted && <span className="ml-2 text-sm text-slate-400 line-through">{formatMoney(tour.originalPrice!)}</span>}
            <span className="text-xs text-slate-500"> / person</span>
          </div>
          <span className={`text-xs ${tour.seats > 0 ? "text-slate-500" : "font-medium text-red-600"}`}>
            {tour.seats > 0 ? `${tour.seats} seats left` : "Sold out"}
          </span>
        </div>
      </div>
    </Link>
  );
}

export function TourGrid({ children }: { children: React.ReactNode }) {
  return <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">{children}</div>;
}
