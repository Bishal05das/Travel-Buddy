"use client";

import { useRouter } from "next/navigation";
import { useEffect } from "react";
import { homeFor, useSession } from "@/lib/auth";
import { LoadingState } from "@/components/ui/Feedback";

export default function HomePage() {
  const { session, ready } = useSession();
  const router = useRouter();
  useEffect(() => {
    if (ready) router.replace(session ? homeFor(session) : "/login");
  }, [ready, session, router]);
  return <LoadingState label="Opening agency portal" />;
}
