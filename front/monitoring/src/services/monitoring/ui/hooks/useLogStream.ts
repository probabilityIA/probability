'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import { getTokenAction } from '../../infra/actions';

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

        // Close existing
        abortRef.current?.abort();
        setError(null);

        const token = await getTokenAction();
        if (!token) {
            setError('Not authenticated');
            return;
        }

        const controller = new AbortController();
        abortRef.current = controller;

        try {
            // Fetch logs via Next.js proxy with token in query (avoids SSE proxy issues)
            const res = await fetch(`/api/logs/${containerId}?token=${encodeURIComponent(token)}`, {
                signal: controller.signal,
            });

            if (!res.ok || !res.body) {
                setError(`HTTP ${res.status}`);
                return;
            }

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
                    if (part.startsWith('data: ')) {
                        newLines.push(part.slice(6));
                    }
                }

                if (newLines.length > 0) {
                    setLines(prev => {
                        const next = [...prev, ...newLines];
                        if (next.length > 1000) return next.slice(-1000);
                        return next;
                    });
                }
            }
        } catch (err) {
            if (err instanceof Error && err.name !== 'AbortError') {
                setError('Connection lost');
            }
        } finally {
            setConnected(false);
        }
    }, [containerId, enabled]);

    const disconnect = useCallback(() => {
        abortRef.current?.abort();
        abortRef.current = null;
        setConnected(false);
    }, []);

    const clear = useCallback(() => {
        setLines([]);
    }, []);

    useEffect(() => {
        if (enabled) {
            connect();
        }
        return () => disconnect();
    }, [connect, disconnect, enabled]);

    return { lines, connected, error, clear, reconnect: connect, disconnect };
}
