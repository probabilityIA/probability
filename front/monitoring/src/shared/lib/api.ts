const API_URL = process.env.MONITORING_API_URL || 'http://localhost:3070';

export async function apiFetch<T>(path: string, options?: RequestInit & { token?: string }): Promise<T> {
    const { token, ...fetchOptions } = options || {};

    const headers: Record<string, string> = {
        'Content-Type': 'application/json',
        ...(fetchOptions.headers as Record<string, string> || {}),
    };

    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const res = await fetch(`${API_URL}${path}`, {
        ...fetchOptions,
        headers,
        cache: 'no-store',
    });

    if (!res.ok) {
        const body = await res.json().catch(() => ({ error: res.statusText }));
        throw new Error(body.error || `API error: ${res.status}`);
    }

    return res.json();
}
