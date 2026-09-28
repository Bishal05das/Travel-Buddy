import { AgencyView } from "@/components/AgencyView";

export default async function AgencyPage(props: PageProps<"/agencies/[id]">) {
  const { id } = await props.params;
  return <AgencyView id={id} />;
}
