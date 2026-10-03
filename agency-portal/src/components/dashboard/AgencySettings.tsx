"use client";

import { useAsync } from "@/hooks/useAsync";
import { getAgency, updateAgencyImage } from "@/lib/endpoints";
import { ImageEditor } from "../ImageEditor";
import { ErrorState, LoadingState } from "../ui/Feedback";
import { useDashboard } from "./DashboardShell";

export function AgencySettings() {
  const { agencyId, reloadAgency } = useDashboard();
  const { data: agency, error, loading, reload } = useAsync(() => getAgency(agencyId), [agencyId]);
  if (loading && !agency) return <LoadingState label="Loading agency settings" />;
  if (error) return <ErrorState message={error} onRetry={reload} />;
  if (!agency) return null;

  return (
    <div className="space-y-4">
      <div>
        <h2 className="text-lg font-semibold">Agency settings</h2>
        <p className="text-sm text-slate-600">Update the image customers see for {agency.name}. Owners and staff with agency:update permission can save changes.</p>
      </div>
      <ImageEditor
        key={agencyId}
        label="Agency image"
        imagePath={agency.image_path}
        permission="agency:update"
        onSave={async (form) => {
          const result = await updateAgencyImage(agencyId, form);
          reloadAgency();
          return result;
        }}
      />
    </div>
  );
}
