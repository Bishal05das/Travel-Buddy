// The access token lives in localStorage and is exposed as a small external
// store, so the API client and React components read the same value.
import type { Role } from "./types";

const KEY = "travelbuddy.token";
const listeners = new Set<() => void>();

export interface Session {
  token: string;
  userId: string; // user_id for customers and super users, member_id for agency members
  role: Role;
  agencyId?: string;
  expiresAt: number; // ms since epoch
}

export function getToken(): string | null {
  try {
    return localStorage.getItem(KEY);
  } catch {
    return null;
  }
}

function notify() {
  listeners.forEach((l) => l());
}

export function setToken(token: string) {
  try {
    localStorage.setItem(KEY, token);
  } catch {
    // storage unavailable (private mode); the session lasts until reload
  }
  notify();
}

export function clearToken() {
  try {
    localStorage.removeItem(KEY);
  } catch {
    // ignore
  }
  notify();
}

export function subscribe(listener: () => void) {
  listeners.add(listener);
  const onStorage = (e: StorageEvent) => e.key === KEY && listener();
  window.addEventListener("storage", onStorage); // logins/logouts in other tabs
  return () => {
    listeners.delete(listener);
    window.removeEventListener("storage", onStorage);
  };
}

/** Reads the claims from a JWT. The server verifies the signature; this is only for display and routing. */
export function decodeSession(token: string | null): Session | null {
  if (!token) return null;
  try {
    const part = token.split(".")[1];
    const json = atob(part.replace(/-/g, "+").replace(/_/g, "/"));
    const claims = JSON.parse(json) as { sub: string; role: Role; agency_id?: string; exp: number };
    const expiresAt = claims.exp * 1000;
    if (!claims.sub || !claims.role || expiresAt <= Date.now()) return null;
    return { token, userId: claims.sub, role: claims.role, agencyId: claims.agency_id, expiresAt };
  } catch {
    return null;
  }
}
