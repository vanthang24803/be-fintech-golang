import { getDefaultApiBaseUrl, readPublicEnv } from "@/config/env";
import { ApiResponse } from "./types";
import { toast } from "sonner";

export function getApiBaseUrl(): string {
  return (readPublicEnv("NEXT_PUBLIC_API_BASE_URL") ?? getDefaultApiBaseUrl()).replace(
    /\/+$/,
    "",
  );
}

export function buildApiUrl(path: string): string {
  return `${getApiBaseUrl()}/${path.replace(/^\/+/, "")}`;
}

let refreshPromise: Promise<string | null> | null = null;

async function refreshAccessToken(): Promise<string | null> {
  if (refreshPromise) return refreshPromise;

  refreshPromise = (async () => {
    const refreshToken = typeof window !== "undefined" ? localStorage.getItem("refresh_token") : null;
    if (!refreshToken) return null;
    try {
      const res = await fetch(buildApiUrl("auth/refresh"), {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ refresh_token: refreshToken }),
      });
      const json = await res.json();
      if (res.ok && json.data?.access_token) {
        localStorage.setItem("access_token", json.data.access_token);
        if (json.data.refresh_token) {
          localStorage.setItem("refresh_token", json.data.refresh_token);
        }
        return json.data.access_token as string;
      }
      return null;
    } catch {
      return null;
    } finally {
      refreshPromise = null;
    }
  })();

  return refreshPromise;
}

async function doFetch<T>(url: string, options: RequestInit, token: string | null): Promise<T> {
  const headers = new Headers(options.headers);
  if (token && !headers.has("Authorization")) {
    headers.set("Authorization", `Bearer ${token}`);
  }
  if (!headers.has("Content-Type") && !(options.body instanceof FormData)) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(url, { ...options, headers, method: options.method || "POST" });

  const contentType = response.headers.get("content-type");
  if (contentType?.includes("application/json")) {
    const result: ApiResponse<T> = await response.json();
    if (!response.ok || (result.code && result.code >= 4000)) {
      const err = new Error(result.message || "An error occurred") as any;
      err.code = result.code;
      throw err;
    }
    return result.data;
  }

  if (!response.ok) {
    throw new Error(`Error ${response.status}: ${response.statusText}`);
  }
  return {} as T;
}

const AUTH_ERROR_CODES = [4010, 4011, 4012];

export async function apiRequest<T>(
  path: string,
  options: RequestInit = {},
): Promise<T> {
  const url = buildApiUrl(path);
  const token = typeof window !== "undefined" ? localStorage.getItem("access_token") : null;

  try {
    return await doFetch<T>(url, options, token);
  } catch (error: any) {
    // Token expired — try refresh once
    if (AUTH_ERROR_CODES.includes(error.code)) {
      const newToken = await refreshAccessToken();
      if (newToken) {
        try {
          return await doFetch<T>(url, options, newToken);
        } catch (retryErr: any) {
          // Refresh worked but request still fails — show error
          toast.error(retryErr.message);
          throw retryErr;
        }
      }
      // Refresh failed — clear tokens silently, let AuthProvider handle redirect
      localStorage.removeItem("access_token");
      localStorage.removeItem("refresh_token");
      throw error;
    }

    toast.error(error.message);
    throw error;
  }
}
