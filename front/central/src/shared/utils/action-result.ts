/**
 * Utility for handling Server Action errors in production.
 *
 * Next.js sanitizes error messages thrown from Server Actions in production,
 * replacing them with: "An error occurred in the Server Components render..."
 *
 * This helper detects sanitized messages and returns a user-friendly fallback.
 */

const SANITIZED_PATTERNS = [
    'An error occurred in the Server',
    'specific message is omitted in production',
    'digest property is included',
];

/**
 * Extract a usable error message from a caught Server Action error.
 * In production, Next.js sanitizes thrown error messages — this detects that
 * and returns the provided fallback instead.
 *
 * @example
 * catch (err: any) {
 *     setError(getActionError(err, 'Error al crear la ruta'));
 * }
 */
export function getActionError(error: unknown, fallback?: string): string {
    const message = error instanceof Error ? error.message : String(error || '');

    const isSanitized = SANITIZED_PATTERNS.some((p) => message.includes(p));
    if (isSanitized || !message) {
        return fallback || 'Ocurrio un error. Por favor intenta de nuevo.';
    }

    return message;
}
