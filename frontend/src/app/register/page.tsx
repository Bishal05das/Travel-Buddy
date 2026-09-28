import type { Metadata } from "next";
import { RegisterForm } from "@/components/auth/RegisterForm";

export const metadata: Metadata = { title: "Create an account" };

export default async function RegisterPage(props: PageProps<"/register">) {
  const sp = await props.searchParams;
  const next = typeof sp.next === "string" && sp.next.startsWith("/") ? sp.next : undefined;
  return <RegisterForm next={next} />;
}
