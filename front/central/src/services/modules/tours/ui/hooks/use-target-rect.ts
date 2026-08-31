'use client';

import { useEffect, useState } from 'react';

export interface TargetRect {
    top: number;
    left: number;
    width: number;
    height: number;
}

const TIMEOUT_MS = 1500;
const POLL_MS = 120;

export function useTargetRect(selector: string | null, enabled: boolean) {
    const [rect, setRect] = useState<TargetRect | null>(null);
    const [notFound, setNotFound] = useState(false);

    useEffect(() => {
        if (!enabled || !selector) {
            setRect(null);
            setNotFound(false);
            return;
        }

        setNotFound(false);
        let cancelled = false;
        let intervalId: ReturnType<typeof setInterval> | null = null;
        const inicio = Date.now();

        const medir = () => {
            const el = document.querySelector(selector) as HTMLElement | null;
            if (!el) {
                if (Date.now() - inicio > TIMEOUT_MS && !cancelled) setNotFound(true);
                return false;
            }
            const box = el.getBoundingClientRect();
            if (box.width === 0 && box.height === 0) return false;
            if (!cancelled) {
                setRect({ top: box.top, left: box.left, width: box.width, height: box.height });
            }
            return true;
        };

        const encontrado = medir();
        if (!encontrado) {
            intervalId = setInterval(() => {
                if (medir() && intervalId) {
                    clearInterval(intervalId);
                    intervalId = null;
                }
            }, POLL_MS);
        }

        const actualizar = () => { medir(); };
        window.addEventListener('resize', actualizar);
        window.addEventListener('scroll', actualizar, true);

        return () => {
            cancelled = true;
            if (intervalId) clearInterval(intervalId);
            window.removeEventListener('resize', actualizar);
            window.removeEventListener('scroll', actualizar, true);
        };
    }, [selector, enabled]);

    return { rect, notFound };
}
