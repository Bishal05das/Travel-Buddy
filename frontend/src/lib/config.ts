// Base URL of the Travel Buddy API, e.g. http://localhost:3000.
export const API_URL = (process.env.NEXT_PUBLIC_API_URL ?? "http://localhost:3000").replace(/\/$/, "");
export const AGENCY_PORTAL_URL = (process.env.NEXT_PUBLIC_AGENCY_PORTAL_URL ?? "http://localhost:3003").replace(/\/$/, "");

/** Absolute URL for an image path returned by the API. */
export function imageUrl(path?: string | null): string | null {
  if (!path) return null;
  // Agency images come as "images/agencies/x.png", tour images as "tours/x.png".
  return path.startsWith("images/") ? `${API_URL}/${path}` : `${API_URL}/images/${path}`;
}
