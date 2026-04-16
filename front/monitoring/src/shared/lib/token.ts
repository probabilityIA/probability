export function getToken(): string | null {
    if (typeof document === 'undefined') return null;
    const match = document.cookie.match(/(?:^|;\s*)monitoring_token=([^;]*)/);
    return match ? decodeURIComponent(match[1]) : null;
}
