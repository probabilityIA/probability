'use client';

import { useCallback, useRef, useEffect } from 'react';
import { useSSE } from '@/shared/hooks/use-sse';

interface StorefrontSSEEvent {
    type?: string;
    data?: { order_id?: string; order_number?: string; [key: string]: unknown };
    metadata?: { event_type?: string };
}

interface UseStorefrontSSEOptions {
    businessId: number;
    onOrderCreated?: (data: StorefrontSSEEvent) => void;
    onOrderUpdated?: (data: StorefrontSSEEvent) => void;
}

export function useStorefrontSSE(options: UseStorefrontSSEOptions) {
    const { businessId, onOrderCreated, onOrderUpdated } = options;

    const callbacksRef = useRef({ onOrderCreated, onOrderUpdated });

    useEffect(() => {
        callbacksRef.current = { onOrderCreated, onOrderUpdated };
    });

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed: StorefrontSSEEvent = JSON.parse(event.data);
            const eventType = parsed.type || parsed.metadata?.event_type;

            switch (eventType) {
                case 'order.created':
                    callbacksRef.current.onOrderCreated?.(parsed);
                    break;
                case 'order.updated':
                case 'order.status_changed':
                    callbacksRef.current.onOrderUpdated?.(parsed);
                    break;
            }
        } catch {
            // Ignore non-JSON messages (heartbeats)
        }
    }, []);

    const { isConnected } = useSSE({
        businessId,
        eventTypes: ['order.created', 'order.updated', 'order.status_changed'],
        onMessage: handleMessage,
    });

    return { isConnected };
}
