import { TourDetail } from "@/components/tours/TourDetail";

export default async function TourPage(props: PageProps<"/tours/[id]">) {
  const { id } = await props.params;
  return <TourDetail id={id} />;
}
