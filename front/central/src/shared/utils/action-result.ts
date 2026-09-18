const SANITIZED_PATTERNS = [
    'An error occurred in the Server',
    'specific message is omitted in production',
    'digest property is included',
    'Minified React error',
];

export function getActionError(error: unknown, fallback?: string): string {
    const message = error instanceof Error ? error.message : String(error || '');

    const isSanitized = SANITIZED_PATTERNS.some((p) => message.includes(p));
    if (isSanitized || !message) {
        return fallback || 'Ocurrió un error. Por favor intenta de nuevo.';
    }

    return message;
}

export function unwrapAction<T>(result: { success: true; data: T } | { success: false; error: string }): T {
    if (!result.success) {
        throw new Error(result.error);
    }
    return result.data;
}
