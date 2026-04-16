'use client';

import { useCallback, useRef, useEffect } from 'react';
import { useSSE } from '@/shared/hooks/use-sse';

export type OrderSSEEventType = 'order.score_calculated' | 'order.created' | 'order.updated' | 'order.status_changed';

interface OrderSSEEventData {
    order_id?: string;
    order_number?: string;
    delivery_probability?: number;
    negative_factors?: string[];
    [key: string]: unknown;
}

interface OrderSSEEvent {
    type?: string;
    data?: OrderSSEEventData;
    metadata?: { event_type?: string };
    order?: Record<string, unknown>;
}

const ORDER_EVENT_TYPES: OrderSSEEventType[] = [
    'order.score_calculated',
    'order.created',
    'order.updated',
    'order.status_changed',
];

interface UseOrderSSEOptions {
    businessId: number;
    onScoreCalculated?: (data: OrderSSEEvent) => void;
    onOrderCreated?: (data: OrderSSEEvent) => void;
    onOrderUpdated?: (data: OrderSSEEvent) => void;
    onStatusChanged?: (data: OrderSSEEvent) => void;
}

export function useOrderSSE(options: UseOrderSSEOptions) {
    const { businessId, onScoreCalculated, onOrderCreated, onOrderUpdated, onStatusChanged } = options;

    const callbacksRef = useRef({ onScoreCalculated, onOrderCreated, onOrderUpdated, onStatusChanged });

    useEffect(() => {
        callbacksRef.current = { onScoreCalculated, onOrderCreated, onOrderUpdated, onStatusChanged };
    });

    const handleMessage = useCallback((event: MessageEvent) => {
        try {
            const parsed: OrderSSEEvent = JSON.parse(event.data);
            const eventType = parsed.type || parsed.metadata?.event_type;

            switch (eventType) {
                case 'order.score_calculated':
                    callbacksRef.current.onScoreCalculated?.(parsed);
                    break;
                case 'order.created':
                    callbacksRef.current.onOrderCreated?.(parsed);
                    break;
                case 'order.updated':
                    callbacksRef.current.onOrderUpdated?.(parsed);
                    break;
                case 'order.status_changed':
                    callbacksRef.current.onStatusChanged?.(parsed);
                    break;
            }
        } catch {
            // Ignore non-JSON messages (heartbeats, etc.)
        }
    }, []);

    const { isConnected } = useSSE({
        businessId,
        eventTypes: ORDER_EVENT_TYPES,
        onMessage: handleMessage,
    });

    return { isConnected };
}
