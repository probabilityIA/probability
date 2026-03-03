const TOKEN_KEY = "testing_token";
const USER_KEY = "testing_user";

export interface UserInfo {
  userId: number;
  businessId: number;
  roleId: number;
  name: string;
  email: string;
}

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return sessionStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  sessionStorage.setItem(TOKEN_KEY, token);
}

export function getUser(): UserInfo | null {
  if (typeof window === "undefined") return null;
  const data = sessionStorage.getItem(USER_KEY);
  return data ? JSON.parse(data) : null;
}

export function setUser(user: UserInfo): void {
  sessionStorage.setItem(USER_KEY, JSON.stringify(user));
}

export function clearAuth(): void {
  sessionStorage.removeItem(TOKEN_KEY);
  sessionStorage.removeItem(USER_KEY);
}

export function isAuthenticated(): boolean {
  return !!getToken();
}

export function parseJWT(token: string): { user_id: number; business_id: number; role_id: number } | null {
  try {
    const base64Url = token.split(".")[1];
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const payload = JSON.parse(atob(base64));
    return {
      user_id: payload.user_id,
      business_id: payload.business_id,
      role_id: payload.role_id,
    };
  } catch {
    return null;
  }
}
