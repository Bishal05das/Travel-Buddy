import { API_URL } from "./config";
import { clearToken, getToken } from "./session";

/** An API failure with the HTTP status (0 when the server was unreachable). */
export class ApiError extends Error {
  constructor(
    public status: number,
    message: string,
  ) {
    super(message);
    this.name = "ApiError";
  }
}


type Body = FormData | object | string | undefined;

async function request<T>(method: string, path: string, body?: Body): Promise<T> {
  const headers: Record<string, string> = {};
  let payload: BodyInit | undefined;
  if (body instanceof FormData) {
    payload = body; // the browser sets the multipart boundary
  } else if (body !== undefined) {
    headers["Content-Type"] = "application/json";
    payload = JSON.stringify(body);
  }
  const token = getToken();
  if (token) headers.Authorization = `Bearer ${token}`;

  let res: Response;
  try {
    res = await fetch(`${API_URL}${path}`, { method, headers, body: payload });
  } catch {
    throw new ApiError(0, "Cannot reach the server. Check your connection and try again.");
  }

  // The API answers with JSON, or plain text for some errors.
  const text = await res.text();
  let data: unknown = text;
  try {
    data = text ? JSON.parse(text) : null;
  } catch {
    // keep the plain text
  }

  if (!res.ok) {
    // An expired or revoked token: drop it so the UI returns to logged-out.
    if (res.status === 401 && token) clearToken();
    throw new ApiError(res.status, errorMessage(res.status, data));
  }
  return data as T;
}

function errorMessage(status: number, data: unknown): string {
  if (status === 429) return "Too many requests. Please wait a moment and try again.";
  if (status === 401) return typeof data === "string" && /invalid email or password/i.test(data) ? "Invalid email or password." : "Please log in again.";
  if (status === 403) return "You don't have permission to do that.";
  const raw = typeof data === "string" ? data.trim() : "";
  if (!raw) return status >= 500 ? "Something went wrong on the server. Please try again." : "The request failed.";
  return friendlyValidation(raw);
}

/**
 * Turns Go validator output such as
 *   Key: 'CreateUserReq.Phone' Error:Field validation for 'Phone' failed on the 'e164' tag
 * into "Phone must be in international format, e.g. +8801712345678".
 */
function friendlyValidation(message: string): string {
  const matches = [...message.matchAll(/Field validation for '(\w+)' failed on the '(\w+)' tag/g)];
  if (matches.length === 0) return message.charAt(0).toUpperCase() + message.slice(1);
  return matches
    .map(([, field, tag]) => {
      const name = field.replace(/([a-z])([A-Z])/g, "$1 $2").replace(/Id$/, "ID");
      switch (tag) {
        case "required":
          return `${name} is required`;
        case "email":
          return `${name} must be a valid email address`;
        case "e164":
          return `${name} must be in international format, e.g. +8801712345678`;
        case "min":
          return `${name} is too short`;
        case "max":
          return `${name} is too long`;
        case "gtfield":
          return `${name} must be after the start date`;
        case "ltefield":
          return `${name} must not be after the start date`;
        default:
          return `${name} is invalid`;
      }
    })
    .join(". ");
}

export const api = {
  get: <T>(path: string) => request<T>("GET", path),
  post: <T>(path: string, body?: Body) => request<T>("POST", path, body),
  put: <T>(path: string, body?: Body) => request<T>("PUT", path, body),
  patch: <T>(path: string, body?: Body) => request<T>("PATCH", path, body),
  delete: <T>(path: string) => request<T>("DELETE", path),
};
