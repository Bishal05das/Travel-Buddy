import { EditTour } from "@/components/dashboard/TourEditors";

export default async function EditTourPage(props: PageProps<"/dashboard/tours/[id]/edit">) {
  const { id } = await props.params;
  return <EditTour id={id} />;
}
