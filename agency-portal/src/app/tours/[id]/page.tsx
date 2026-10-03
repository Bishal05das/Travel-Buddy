import { GuestBookingTour } from "@/components/dashboard/GuestBookingTour";

export default async function GuestBookingPage(props: PageProps<"/tours/[id]">) {
  const { id } = await props.params;
  return <GuestBookingTour id={id} />;
}
