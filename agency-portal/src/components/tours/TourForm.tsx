"use client";

import { useState } from "react";
import { formatMoney, fromDateInput, toDateInput, unitPrice } from "@/lib/format";
import type { Tour } from "@/lib/types";
import { Button } from "../ui/Button";
import { Alert } from "../ui/Feedback";
import { Input, TextArea } from "../ui/Field";

export interface TourFormValues {
  name: string;
  description: string;
  start_date: string; // YYYY-MM-DD
  end_date: string;
  last_enrollment_date: string;
  total_seat: string;
  price: string;
  discount: string;
}

const MAX_IMAGE_BYTES = 10 * 1024 * 1024;

export function tourToForm(t: Tour): TourFormValues {
  return {
    name: t.name,
    description: t.description,
    start_date: toDateInput(t.start_date),
    end_date: toDateInput(t.end_date),
    last_enrollment_date: toDateInput(t.last_enrollment_date),
    total_seat: String(t.total_seat),
    price: String(t.price),
    discount: String(t.discount),
  };
}

export function formToUpdate(v: TourFormValues) {
  return {
    name: v.name.trim(),
    description: v.description.trim(),
    start_date: fromDateInput(v.start_date),
    end_date: fromDateInput(v.end_date),
    last_enrollment_date: fromDateInput(v.last_enrollment_date),
    total_seat: Number(v.total_seat),
    price: Number(v.price),
    discount: Number(v.discount || 0),
  };
}

const empty: TourFormValues = {
  name: "",
  description: "",
  start_date: "",
  end_date: "",
  last_enrollment_date: "",
  total_seat: "",
  price: "",
  discount: "0",
};

/**
 * Create (with image upload) or edit a tour. `bookedSeats` is set when
 * editing, since capacity can't go below the seats already booked.
 */
export function TourForm({
  initial,
  mode,
  bookedSeats = 0,
  onSubmit,
}: {
  initial?: TourFormValues;
  mode: "create" | "edit";
  bookedSeats?: number;
  onSubmit: (values: TourFormValues, image: File | null) => Promise<void>;
}) {
  const [v, setV] = useState<TourFormValues>(initial ?? empty);
  const [image, setImage] = useState<File | null>(null);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const set = (key: keyof TourFormValues) => (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) =>
    setV({ ...v, [key]: e.target.value });

  function validate() {
    const e: Record<string, string> = {};
    const seats = Number(v.total_seat);
    const price = Number(v.price);
    const discount = Number(v.discount || 0);
    if (v.name.trim().length < 3) e.name = "Use at least 3 characters.";
    if (v.description.trim().length < 10) e.description = "Use at least 10 characters.";
    if (!v.start_date) e.start_date = "Pick a start date.";
    if (!v.end_date) e.end_date = "Pick an end date.";
    else if (v.start_date && v.end_date <= v.start_date) e.end_date = "Must be after the start date.";
    if (!v.last_enrollment_date) e.last_enrollment_date = "Pick the last day to book.";
    else if (v.start_date && v.last_enrollment_date > v.start_date) e.last_enrollment_date = "Must be on or before the start date.";
    if (!Number.isInteger(seats) || seats < 1) e.total_seat = "Enter a whole number of at least 1.";
    else if (seats < bookedSeats) e.total_seat = `${bookedSeats} seats are already booked; capacity can't be lower.`;
    if (!Number.isInteger(price) || price < 1) e.price = "Enter a price in whole taka.";
    if (!Number.isInteger(discount) || discount < 0 || discount > 100) e.discount = "Enter 0 to 100.";
    if (mode === "create") {
      if (!image) e.image = "Add a cover image.";
      else if (!["image/jpeg", "image/png", "image/webp"].includes(image.type)) e.image = "Use a JPG, PNG or WebP image.";
      else if (image.size > MAX_IMAGE_BYTES) e.image = "The image must be under 10 MB.";
    }
    setErrors(e);
    return Object.keys(e).length === 0;
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!validate()) return;
    setSubmitting(true);
    setSubmitError(null);
    try {
      await onSubmit(v, image);
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Saving failed.");
      setSubmitting(false);
    }
  }

  const preview = Number(v.price) > 0 ? unitPrice({ price: Number(v.price), discount: Number(v.discount || 0) }) : null;

  return (
    <form onSubmit={submit} className="max-w-2xl space-y-4" noValidate>
      <Input label="Tour name" value={v.name} error={errors.name} onChange={set("name")} />
      <TextArea label="Description" value={v.description} error={errors.description} onChange={set("description")} />
      <div className="grid gap-4 sm:grid-cols-3">
        <Input label="Start date" type="date" value={v.start_date} error={errors.start_date} onChange={set("start_date")} />
        <Input label="End date" type="date" value={v.end_date} error={errors.end_date} onChange={set("end_date")} />
        <Input
          label="Last day to book"
          type="date"
          value={v.last_enrollment_date}
          error={errors.last_enrollment_date}
          onChange={set("last_enrollment_date")}
        />
      </div>
      <div className="grid gap-4 sm:grid-cols-3">
        <Input
          label="Total seats"
          type="number"
          min={Math.max(1, bookedSeats)}
          value={v.total_seat}
          error={errors.total_seat}
          hint={mode === "edit" ? `${bookedSeats} already booked` : undefined}
          onChange={set("total_seat")}
        />
        <Input label="Price per person (৳)" type="number" min={1} value={v.price} error={errors.price} onChange={set("price")} />
        <Input
          label="Discount (%)"
          type="number"
          min={0}
          max={100}
          value={v.discount}
          error={errors.discount}
          hint={preview !== null ? `Customers pay ${formatMoney(preview)}` : undefined}
          onChange={set("discount")}
        />
      </div>
      {mode === "create" && (
        <Input
          label="Cover image"
          type="file"
          accept="image/jpeg,image/png,image/webp"
          error={errors.image}
          hint="JPG, PNG or WebP, up to 10 MB."
          onChange={(e) => setImage(e.target.files?.[0] ?? null)}
        />
      )}
      {submitError && <Alert tone="error">{submitError}</Alert>}
      <Button type="submit" loading={submitting}>
        {mode === "create" ? "Create tour" : "Save changes"}
      </Button>
    </form>
  );
}
