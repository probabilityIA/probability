export async function register() {
    if (process.env.NEXT_RUNTIME !== 'nodejs') return;

    const { infoFromBody, isPublicAuthRequest, ACCESS_DENIED_COOKIE } = await import('./shared/utils/access-denied');
    const apiBase = process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || '';
    const originalFetch = globalThis.fetch;

    globalThis.fetch = async (input: RequestInfo | URL, init?: RequestInit) => {
        const response = await originalFetch(input, init);
        if (response.status !== 401 && response.status !== 402 && response.status !== 403) {
            return response;
        }

        const url = typeof input === 'string' ? input : input instanceof URL ? input.toString() : input.url;
        if (!apiBase || !url.startsWith(apiBase) || isPublicAuthRequest(url)) {
            return response;
        }

        try {
            const body = await response.clone().json().catch(() => ({}));
            const info = infoFromBody(response.status, body);
            if (!info) return response;
            const { cookies } = await import('next/headers');
            const store = await cookies();
            store.set(ACCESS_DENIED_COOKIE, JSON.stringify(info), {
                path: '/',
                maxAge: 20,
                httpOnly: false,
                sameSite: 'lax',
            });
        } catch {
            return response;
        }
        return response;
    };
}
