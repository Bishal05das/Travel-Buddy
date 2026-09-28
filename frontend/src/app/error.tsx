"use client";

import { ErrorState } from "@/components/ui/Feedback";

// Catches unexpected rendering errors in any route below the root layout.
export default function Error({ reset }: { error: Error & { digest?: string }; reset: () => void }) {
  return <ErrorState message="Something went wrong while showing this page." onRetry={reset} />;
}
