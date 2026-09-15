export type AccessDeniedCode = 'forbidden' | 'subscription_suspended' | 'unauthenticated';

export interface AccessDeniedInfo {
    code: AccessDeniedCode;
    message: string;
    permission?: string;
}

export const ACCESS_DENIED_EVENT = 'probability:access-denied';
export const ACCESS_DENIED_COOKIE = 'pb_access_denied';

const PUBLIC_AUTH_PATHS = ['/auth/login', '/auth/forgot-password', '/auth/reset-password', '/auth/verify-otp', '/auth/recovery-channels', '/auth/demo-'];

export function notifyAccessDenied(info: AccessDeniedInfo): void {
    if (typeof window === 'undefined') return;
    window.dispatchEvent(new CustomEvent<AccessDeniedInfo>(ACCESS_DENIED_EVENT, { detail: info }));
}

export function isAccessDeniedMessage(message: string): AccessDeniedInfo | null {
    if (!message) return null;
    if (/^Tu suscripci[oó]n est[aá] vencida/i.test(message)) {
        return { code: 'subscription_suspended', message };
    }
    if (/^Tu sesi[oó]n expir[oó]/i.test(message)) {
        return { code: 'unauthenticated', message };
    }
    if (/^No tienes permiso (para|de)/i.test(message) || message === 'forbidden') {
        return { code: 'forbidden', message: message === 'forbidden' ? 'No tienes permiso para esta acción.' : message };
    }
    return null;
}

export function isPublicAuthRequest(url: string): boolean {
    return PUBLIC_AUTH_PATHS.some((path) => url.includes(path));
}

export function codeForStatus(status: number): AccessDeniedCode | null {
    if (status === 401) return 'unauthenticated';
    if (status === 402) return 'subscription_suspended';
    if (status === 403) return 'forbidden';
    return null;
}

export function defaultMessageFor(code: AccessDeniedCode): string {
    switch (code) {
        case 'unauthenticated':
            return 'Tu sesión expiró. Vuelve a iniciar sesión.';
        case 'subscription_suspended':
            return 'Tu suscripción está vencida. Ponte al día en el módulo de Suscripción para seguir usando esta función.';
        default:
            return 'No tienes permiso para esta acción.';
    }
}

export function infoFromBody(status: number, body: any): AccessDeniedInfo | null {
    const code = (body?.code as AccessDeniedCode | undefined) ?? codeForStatus(status);
    if (!code || !['forbidden', 'subscription_suspended', 'unauthenticated'].includes(code)) return null;
    const message = typeof body?.message === 'string' && body.message ? body.message : defaultMessageFor(code);
    return { code, message, permission: typeof body?.permission === 'string' ? body.permission : undefined };
}

export function readAccessDeniedCookie(): AccessDeniedInfo | null {
    if (typeof document === 'undefined') return null;
    const raw = document.cookie.split('; ').find((c) => c.startsWith(`${ACCESS_DENIED_COOKIE}=`));
    if (!raw) return null;
    document.cookie = `${ACCESS_DENIED_COOKIE}=; Max-Age=0; path=/`;
    let value = raw.slice(ACCESS_DENIED_COOKIE.length + 1);
    for (let i = 0; i < 3 && !value.trim().startsWith('{'); i++) {
        try {
            value = decodeURIComponent(value);
        } catch {
            return null;
        }
    }
    try {
        return JSON.parse(value) as AccessDeniedInfo;
    } catch {
        return null;
    }
}
