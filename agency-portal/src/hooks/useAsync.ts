"use client";

import { useCallback, useEffect, useRef, useState, type DependencyList } from "react";
import { ApiError } from "@/lib/api";

export interface AsyncState<T> {
  data: T | undefined;
  error: string | null;
  status: number | null; // HTTP status of the error, if any
  loading: boolean;
  reload: () => void;
}

/**
 * Runs `load` when `deps` change and tracks loading / error / data. Stale
 * responses from earlier runs are ignored. Data from the previous run stays
 * visible while reloading, so lists don't flash empty.
 */
export function useAsync<T>(load: () => Promise<T>, deps: DependencyList): AsyncState<T> {
  const [data, setData] = useState<T>();
  const [error, setError] = useState<{ message: string; status: number | null } | null>(null);
  const [loading, setLoading] = useState(true);
  const [nonce, setNonce] = useState(0);
  const loadRef = useRef(load);

  useEffect(() => {
    loadRef.current = load;
  });

  useEffect(() => {
    let active = true;
    // eslint-disable-next-line react-hooks/set-state-in-effect -- starting a fetch is the effect
    setLoading(true);
    setError(null);
    loadRef
      .current()
      .then((result) => {
        if (active) setData(result);
      })
      .catch((e: unknown) => {
        if (!active) return;
        setError({
          message: e instanceof Error ? e.message : "Something went wrong.",
          status: e instanceof ApiError ? e.status : null,
        });
      })
      .finally(() => {
        if (active) setLoading(false);
      });
    return () => {
      active = false;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps -- callers pass the inputs of `load` as deps
  }, [...deps, nonce]);

  const reload = useCallback(() => setNonce((n) => n + 1), []);
  return { data, error: error?.message ?? null, status: error?.status ?? null, loading, reload };
}
