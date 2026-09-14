import { useEffect, useRef, useState, useCallback } from 'react';
import { envPublic } from '@/shared/config/env';

interface UseSSEOptions {
    onMessage?: (event: MessageEvent) => void;
    onError?: (event: Event) => void;
    onOpen?: (event: Event) => void;
    eventTypes?: string[];
    integrationId?: number;
    businessId?: number;
    orderIds?: string[];
    enabled?: boolean;
}

export const useSSE = (options: UseSSEOptions = {}) => {
    const [isConnected, setIsConnected] = useState(false);
    const eventSourceRef = useRef<EventSource | null>(null);

    const onMessageRef = useRef(options.onMessage);
    const onErrorRef = useRef(options.onError);
    const onOpenRef = useRef(options.onOpen);

    useEffect(() => {
        onMessageRef.current = options.onMessage;
        onErrorRef.current = options.onError;
        onOpenRef.current = options.onOpen;
    });

    const enabled = options.enabled ?? true;
    const connectionParams = JSON.stringify({
        eventTypes: options.eventTypes,
        integrationId: options.integrationId,
        businessId: options.businessId,
        orderIds: options.orderIds,
        enabled,
    });

    const connect = useCallback(() => {
        const { eventTypes, integrationId, businessId, orderIds, enabled: paramsEnabled } = JSON.parse(connectionParams);

        if (eventSourceRef.current) {
            eventSourceRef.current.close();
            eventSourceRef.current = null;
        }
        if (!paramsEnabled) {
            return;
        }

        const params = new URLSearchParams();
        if (eventTypes && eventTypes.length > 0) {
            params.append('event_types', eventTypes.join(','));
        }
        if (integrationId) {
            params.append('integration_id', integrationId.toString());
        }
        if (businessId) {
            params.append('business_id', businessId.toString());
        }
        if (orderIds && orderIds.length > 0) {
            params.append('order_ids', orderIds.join(','));
        }

        const baseUrl = `${envPublic.SSE_BASE_URL}/notify/sse/order-notify`;
        const url = `${baseUrl}?${params.toString()}`;

        const eventSource = new EventSource(url, { withCredentials: true });

        eventSource.onopen = (event) => {
            setIsConnected(true);
            if (onOpenRef.current) onOpenRef.current(event);
        };

        eventSource.onmessage = (event) => {
            if (onMessageRef.current) onMessageRef.current(event);
        };

        eventSource.onerror = (event) => {
            setIsConnected(false);
            if (onErrorRef.current) onErrorRef.current(event);
        };

        if (eventTypes) {
            eventTypes.forEach((type: string) => {
                eventSource.addEventListener(type, (event) => {
                    if (onMessageRef.current) onMessageRef.current(event);
                });
            });
        }

        eventSourceRef.current = eventSource;
    }, [connectionParams]);

    const disconnect = useCallback(() => {
        if (eventSourceRef.current) {
            eventSourceRef.current.close();
            eventSourceRef.current = null;
            setIsConnected(false);
        }
    }, []);

    useEffect(() => {
        connect();
        return () => {
            disconnect();
        };
    }, [connect, disconnect]);

    return { isConnected, disconnect, connect };
};
