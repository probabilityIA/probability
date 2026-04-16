'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { getToken } from '@/shared/lib/token';

interface UseLogStreamOptions {
    containerId: string;
    enabled?: boolean;
}

export function useLogStream({ containerId, enabled = true }: UseLogStreamOptions) {
    const [lines, setLines] = useState<string[]>([]);
    const [connected, setConnected] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const abortRef = useRef<AbortController | null>(null);

    const connect = useCallback(async () => {
        if (!containerId || !enabled) return;

        abortRef.current?.abort();
        setError(null);

        const token = getToken();
        if (!token) { setError('Not authenticated'); return; }

        const controller = new AbortController();
        abortRef.current = controller;

        try {
            const res = await fetch(`/api/logs/${containerId}?token=${encodeURIComponent(token)}`, {
                signal: controller.signal,
            });

            if (!res.ok || !res.body) { setError(`HTTP ${res.status}`); return; }

            setConnected(true);
            const reader = res.body.getReader();
            const decoder = new TextDecoder();
            let buffer = '';

            while (true) {
                const { done, value } = await reader.read();
                if (done) break;

                buffer += decoder.decode(value, { stream: true });
                const parts = buffer.split('\n');
                buffer = parts.pop() || '';

                const newLines: string[] = [];
                for (const part of parts) {
                    if (part.startsWith('data: ')) newLines.push(part.slice(6));
                }

                if (newLines.length > 0) {
                    setLines(prev => {
                        const next = [...prev, ...newLines];
                        return next.length > 1000 ? next.slice(-1000) : next;
                    });
                }
            }
        } catch (err) {
            if (err instanceof Error && err.name !== 'AbortError') setError('Connection lost');
        } finally {
            setConnected(false);
        }
    }, [containerId, enabled]);

    const disconnect = useCallback(() => {
        abortRef.current?.abort();
        abortRef.current = null;
        setConnected(false);
    }, []);

    const clear = useCallback(() => setLines([]), []);

    useEffect(() => {
        if (enabled) connect();
        return () => disconnect();
    }, [connect, disconnect, enabled]);

    return { lines, connected, error, clear, reconnect: connect, disconnect };
}
