import { getToken } from "./auth";

const CENTRAL_API = "/api/central/v1";
const TESTING_API = "/api/testing/v1";

async function fetchWithAuth(url: string, options: RequestInit = {}): Promise<Response> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string> || {}),
  };

  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const res = await fetch(url, { ...options, headers });
  return res;
}

export async function centralAPI<T>(path: string, options: RequestInit = {}): Promise<T> {
  let res: Response;
  try {
    res = await fetchWithAuth(`${CENTRAL_API}${path}`, options);
  } catch {
    throw new Error("Central backend not reachable. Make sure it's running (make run-backend)");
  }
  let data;
  try {
    data = await res.json();
  } catch {
    throw new Error(`Central API returned invalid response (status ${res.status})`);
  }
  if (!res.ok) {
    throw new Error(data.error || data.message || `Request failed: ${res.status}`);
  }
  return data;
}

export async function testingAPI<T>(path: string, options: RequestInit = {}): Promise<T> {
  let res: Response;
  try {
    res = await fetchWithAuth(`${TESTING_API}${path}`, options);
  } catch {
    throw new Error("Testing backend not reachable. Make sure it's running (make run-testing)");
  }
  let data;
  try {
    data = await res.json();
  } catch {
    throw new Error(`Testing API returned invalid response (status ${res.status})`);
  }
  if (!res.ok) {
    throw new Error(data.error || data.message || `Request failed: ${res.status}`);
  }
  return data;
}
