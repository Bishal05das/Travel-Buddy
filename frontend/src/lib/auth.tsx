"use client";

import { useCallback, useEffect, useMemo, useSyncExternalStore } from "react";
import { clearToken, decodeSession, getToken, setToken, subscribe, type Session } from "./session";

const noopSubscribe = () => () => {};

/**
 * The current session. `ready` is false during server rendering and
 * hydration, before localStorage can be read; route guards must wait for it.
 */
export function useSession(): { session: Session | null; ready: boolean } {
  const token = useSyncExternalStore(subscribe, getToken, () => null);
  const ready = useSyncExternalStore(noopSubscribe, () => true, () => false);
  const session = useMemo(() => decodeSession(token), [token]);

  // Log out exactly when the token expires.
  useEffect(() => {
    if (!session) {
      if (token) clearToken(); // stored token is expired or malformed
      return;
    }
    const timer = setTimeout(clearToken, session.expiresAt - Date.now());
    return () => clearTimeout(timer);
  }, [session, token]);

  return { session, ready };
}

export function useAuthActions() {
  const login = useCallback((token: string) => setToken(token), []);
  const logout = useCallback(() => clearToken(), []);
  return { login, logout };
}

/** Where each role lands after logging in. */
export function homeFor(session: Session): string {
  switch (session.role) {
    case "member":
      return "/dashboard";
    case "super":
      return "/admin";
    default:
      return "/bookings";
  }
}
