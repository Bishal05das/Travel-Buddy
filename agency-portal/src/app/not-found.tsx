import { ButtonLink } from "@/components/ui/Button";
import { EmptyState } from "@/components/ui/Feedback";

export default function NotFound() {
  return <EmptyState title="Page not found" description="The page you're looking for doesn't exist." action={<ButtonLink href="/">Back to portal</ButtonLink>} />;
}
