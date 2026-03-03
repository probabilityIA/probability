"use server";

interface LoginResult {
  token: string;
  user: { id: number; name: string; email: string };
  is_super_admin: boolean;
  scope: string;
}

export async function loginAction(
  email: string,
  password: string
): Promise<LoginResult> {
  const centralUrl =
    process.env.CENTRAL_API_URL || "http://localhost:3050";

  const res = await fetch(`${centralUrl}/api/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });

  const body = await res.json();

  if (!res.ok) {
    throw new Error(body.error || body.message || "Login failed");
  }

  const data = body.data;

  // Extract JWT from Set-Cookie header (server-side only)
  let token = "";
  const setCookie = res.headers.get("set-cookie");
  if (setCookie) {
    const match = setCookie.match(/session_token=([^;]+)/);
    if (match) {
      token = match[1];
    }
  }

  if (!token) {
    throw new Error(
      "Could not extract session token. Make sure the central backend is running."
    );
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
