"use server";

// Single backend URL — all modules use their own routes under this base
const API_URL =
  process.env.API_URL || "http://localhost:9092";

async function fetchServer(
  url: string,
  token: string,
  options: RequestInit = {}
) {
  const res = await fetch(url, {
    ...options,
    cache: "no-store",
    headers: {
      "Content-Type": "application/json",
      Authorization: `Bearer ${token}`,
      ...((options.headers as Record<string, string>) || {}),
    },
  });

  const data = await res.json();

  if (!res.ok) {
    throw new Error(
      data.error || data.message || `Request failed: ${res.status}`
    );
  }

  return data;
}

// ── Auth ─────────────────────────────────────────────
export async function loginAction(email: string, password: string) {
  const res = await fetch(`${API_URL}/api/v1/auth/login`, {
    method: "POST",
    cache: "no-store",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const body = await res.json();
  if (!res.ok) {
    throw new Error(body.error || body.message || "Login failed");
  }

  const data = body.data;

  // Extract JWT from Set-Cookie header
  let token = "";
  const setCookie = res.headers.get("set-cookie");
  if (setCookie) {
    const match = setCookie.match(/session_token=([^;]+)/);
    if (match) token = match[1];
  }

  if (!token) {
    throw new Error("Could not extract session token");
  }

  return {
    token,
    user: {
      id: data.user?.id ?? 0,
      name: data.user?.name ?? "",
      email: data.user?.email ?? email,
    },
    is_super_admin: data.is_super_admin ?? false,
    scope: data.scope ?? "",
  };
}

// ── Businesses ───────────────────────────────────────
export async function fetchBusinessesAction(token: string) {
  const data = await fetchServer(`${API_URL}/api/v1/businesses`, token);
  return data.data as { id: number; name: string }[];
}

// ── Orders ───────────────────────────────────────────
export async function fetchReferenceDataAction(
  token: string,
  businessId: number
) {
  const data = await fetchServer(
    `${API_URL}/api/v1/orders/reference-data?business_id=${businessId}`,
    token
  );
  return data.data;
}

export async function buildPayloadsAction(
  token: string,
  businessId: number,
  dto: {
    count: number;
    integration_id?: number;
    random_products: boolean;
    max_items_per_order: number;
    topic: string;
    customer_name?: string;
    customer_phone?: string;
  }
) {
  const data = await fetchServer(
    `${API_URL}/api/v1/orders/generate?business_id=${businessId}`,
    token,
    {
      method: "POST",
      body: JSON.stringify(dto),
    }
  );
  return data.data;
}

// Send a single pre-built payload to its target URL
export async function sendPayloadAction(payload: {
  url: string;
  method: string;
  headers: Record<string, string>;
  body: Record<string, unknown>;
  raw_body?: string;
}): Promise<{
  status_code: number;
  body: string;
  duration_ms: number;
}> {
  const start = Date.now();

  try {
    // Use raw_body when available (preserves exact bytes for HMAC validation)
    const requestBody = payload.raw_body || JSON.stringify(payload.body);

    const res = await fetch(payload.url, {
      method: payload.method,
      cache: "no-store",
      headers: payload.headers,
      body: requestBody,
    });

    const text = await res.text();
    const durationMs = Date.now() - start;

    return {
      status_code: res.status,
      body: text,
      duration_ms: durationMs,
    };
  } catch (err: unknown) {
    const durationMs = Date.now() - start;
    return {
      status_code: 0,
      body: err instanceof Error ? err.message : "Network error",
      duration_ms: durationMs,
    };
  }
}
