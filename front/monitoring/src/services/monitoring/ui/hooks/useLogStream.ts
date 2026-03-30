'use client';

import { useState, useEffect, useRef, useCallback } from 'react';

interface UseLogStreamOptions {
    containerId: string;
    enabled?: boolean;
}

export function useLogStream({ containerId, enabled = true }: UseLogStreamOptions) {
    const [lines, setLines] = useState<string[]>([]);
    const [connected, setConnected] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const eventSourceRef = useRef<EventSource | null>(null);

    const connect = useCallback(() => {
        if (!containerId || !enabled) return;

        // Close existing connection
        eventSourceRef.current?.close();
        setError(null);

        const es = new EventSource(`/api/logs/${containerId}`);
        eventSourceRef.current = es;

        es.onopen = () => {
            setConnected(true);
            setError(null);
        };

        es.onmessage = (event) => {
            setLines(prev => {
                const next = [...prev, event.data];
                // Keep last 1000 lines to avoid memory issues
                if (next.length > 1000) return next.slice(-1000);
                return next;
            });
        };

        es.onerror = () => {
            setConnected(false);
            setError('Connection lost');
            es.close();
        };
    }, [containerId, enabled]);

    const disconnect = useCallback(() => {
        eventSourceRef.current?.close();
        eventSourceRef.current = null;
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
