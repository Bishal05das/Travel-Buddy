"use client";

import type { Permission } from "@/lib/types";

/** Checkboxes for permissions, grouped by resource (tour, booking, member, agency). */
export function PermissionPicker({
  permissions,
  selected,
  onChange,
  error,
}: {
  permissions: Permission[];
  selected: number[];
  onChange: (ids: number[]) => void;
  error?: string;
}) {
  const groups = new Map<string, Permission[]>();
  for (const p of permissions) groups.set(p.resource, [...(groups.get(p.resource) ?? []), p]);

  const toggle = (id: number) => onChange(selected.includes(id) ? selected.filter((x) => x !== id) : [...selected, id]);

  return (
    <fieldset className="space-y-2">
      <legend className="text-sm font-medium text-slate-700">Permissions</legend>
      <div className="grid gap-3 sm:grid-cols-2">
        {[...groups.entries()].map(([resource, perms]) => (
          <div key={resource} className="rounded-lg border border-slate-200 p-3">
            <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-slate-500">{resource}s</p>
            <div className="flex flex-wrap gap-x-4 gap-y-1">
              {perms.map((p) => (
                <label key={p.permission_id} className="flex items-center gap-2 text-sm">
                  <input
                    type="checkbox"
                    className="h-4 w-4 accent-teal-700"
                    checked={selected.includes(p.permission_id)}
                    onChange={() => toggle(p.permission_id)}
                  />
                  {p.action}
                </label>
              ))}
            </div>
          </div>
        ))}
      </div>
      {error && <p className="text-xs text-red-600">{error}</p>}
    </fieldset>
  );
}

export function permissionNames(ids: number[], all: Permission[]): string {
  const names = ids.map((id) => all.find((p) => p.permission_id === id)?.name ?? `#${id}`);
  return names.length ? names.join(", ") : "No permissions";
}
