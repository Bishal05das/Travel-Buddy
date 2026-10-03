"use client";

import Link from "next/link";
import { useState } from "react";
import { createBooking } from "@/lib/endpoints";
import { formatMoney, unitPrice } from "@/lib/format";
import type { Booking, PaymentMethod, Tour } from "@/lib/types";
import { Button } from "../ui/Button";
import { Alert } from "../ui/Feedback";
import { Input, Select } from "../ui/Field";

const MAX_PEOPLE_PER_BOOKING = 20;

/** Books a tour for the logged-in customer. */
export function BookingForm({ tour, onBooked }: { tour: Tour; onBooked: () => void }) {
  const maxPeople = Math.min(tour.available_seat, MAX_PEOPLE_PER_BOOKING);
  const [people, setPeople] = useState(1);
  const [method, setMethod] = useState<PaymentMethod>("bkash");
  const [transactionId, setTransactionId] = useState("");
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);
  const [booked, setBooked] = useState<Booking | null>(null);

  const total = unitPrice(tour) * people;

  function validate() {
    const e: Record<string, string> = {};
    if (!Number.isInteger(people) || people < 1 || people > maxPeople) e.people = `Choose between 1 and ${maxPeople} people.`;
    if (transactionId.trim().length < 5) e.transactionId = "Enter the transaction ID from your payment (at least 5 characters).";
    setErrors(e);
    return Object.keys(e).length === 0;
  }

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    if (!validate()) return;
    setSubmitting(true);
    setSubmitError(null);
    try {
      const input = { number_of_people: people, total_price: total, method, transaction_id: transactionId.trim() };
      const result = await createBooking(tour.tour_id, input);
      setBooked(result);
      onBooked();
    } catch (err) {
      setSubmitError(err instanceof Error ? err.message : "Booking failed.");
    } finally {
      setSubmitting(false);
    }
  }

  if (booked) {
    return (
      <div className="space-y-3">
        <Alert tone="success">
          <p className="font-medium">Booking received for {booked.number_of_people} {booked.number_of_people === 1 ? "person" : "people"}.</p>
          <p className="mt-1">
            It&apos;s <strong>pending</strong> until the agency verifies the payment of {formatMoney(booked.total_price)}.
          </p>
        </Alert>
        <Link href="/bookings" className="text-sm font-medium text-teal-700 hover:underline">
          View my bookings →
        </Link>
      </div>
    );
  }

  return (
    <form onSubmit={submit} className="space-y-4" noValidate>
      <Input
        label="Number of people"
        type="number"
        min={1}
        max={maxPeople}
        value={people}
        error={errors.people}
        hint={`${tour.available_seat} seats left`}
        onChange={(e) => setPeople(Number(e.target.value))}
      />
      <Select label="Payment method" value={method} onChange={(e) => setMethod(e.target.value as PaymentMethod)}>
        <option value="bkash">bKash</option>
        <option value="nagad">Nagad</option>
        <option value="bank">Bank transfer</option>
      </Select>
      <Input
        label="Transaction ID"
        value={transactionId}
        error={errors.transactionId}
        hint="Pay the total first, then enter the transaction ID from your receipt."
        onChange={(e) => setTransactionId(e.target.value)}
      />
      <div className="flex items-center justify-between rounded-lg bg-slate-50 px-4 py-3 text-sm">
        <span className="text-slate-600">
          {formatMoney(unitPrice(tour))} × {Number.isFinite(people) ? people : 0}
        </span>
        <span className="text-base font-semibold text-slate-900">Total {formatMoney(Number.isFinite(total) ? total : 0)}</span>
      </div>
      {submitError && <Alert tone="error">{submitError}</Alert>}
      <Button type="submit" loading={submitting} className="w-full">
        Book now
      </Button>
    </form>
  );
}
