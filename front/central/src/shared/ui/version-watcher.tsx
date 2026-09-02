'use client';

import { useCallback, useEffect, useRef, useState } from 'react';

const CHECK_INTERVAL_MS = 5 * 60 * 1000;
const HIDDEN_RELOAD_DELAY_MS = 30 * 1000;

const CURRENT_VERSION = process.env.NEXT_PUBLIC_APP_VERSION || 'dev';

export function VersionWatcher() {
    const [stale, setStale] = useState(false);
    const staleRef = useRef(false);
    const hiddenTimer = useRef<ReturnType<typeof setTimeout> | null>(null);

    const reload = useCallback(() => {
        window.location.reload();
    }, []);

    const check = useCallback(async () => {
        if (staleRef.current || CURRENT_VERSION === 'dev') return;
        try {
            const res = await fetch('/api/app-version', { cache: 'no-store' });
            if (!res.ok) return;
            const data = (await res.json()) as { version?: string };
            if (data.version && data.version !== 'dev' && data.version !== CURRENT_VERSION) {
                staleRef.current = true;
                setStale(true);
            }
        } catch {
            return;
        }
    }, []);

    useEffect(() => {
        if (CURRENT_VERSION === 'dev') return;

        check();
        const interval = setInterval(check, CHECK_INTERVAL_MS);

        const onVisibility = () => {
            if (document.visibilityState === 'visible') {
                if (hiddenTimer.current) {
                    clearTimeout(hiddenTimer.current);
                    hiddenTimer.current = null;
                }
                check();
                return;
            }
            if (staleRef.current && !hiddenTimer.current) {
                hiddenTimer.current = setTimeout(reload, HIDDEN_RELOAD_DELAY_MS);
            }
        };

        document.addEventListener('visibilitychange', onVisibility);
        window.addEventListener('focus', check);

        return () => {
            clearInterval(interval);
            document.removeEventListener('visibilitychange', onVisibility);
            window.removeEventListener('focus', check);
            if (hiddenTimer.current) clearTimeout(hiddenTimer.current);
        };
    }, [check, reload]);

    if (!stale) return null;

    return (
        <div className="fixed bottom-4 left-1/2 z-[100] -translate-x-1/2 px-4">
            <div className="flex items-center gap-3 rounded-xl border border-gray-200 bg-white px-4 py-3 shadow-lg dark:border-gray-700 dark:bg-gray-800">
                <span className="text-sm text-gray-700 dark:text-gray-200">
                    {'Hay una versión nueva de la aplicación.'}
                </span>
                <button
                    type="button"
                    onClick={reload}
                    style={{ backgroundColor: 'var(--color-primary)' }}
                    className="inline-flex h-8 items-center rounded-lg px-3 text-sm font-semibold text-white transition hover:opacity-90"
                >
                    Actualizar
                </button>
            </div>
        </div>
    );
}
