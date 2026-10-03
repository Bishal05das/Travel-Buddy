import type { Metadata } from "next";
import { AdminOnboarding } from "@/components/AdminOnboarding";

export const metadata: Metadata = { title: "Admin" };

export default function AdminPage() {
  return <AdminOnboarding />;
}
