import { useEffect, useRef, useState, useCallback } from 'react';
import { envPublic } from '@/shared/config/env';

interface UseSSEOptions {
    onMessage?: (event: MessageEvent) => void;
    onError?: (event: Event) => void;
    onOpen?: (event: Event) => void;
    eventTypes?: string[];
    integrationId?: number;
    businessId?: number;
}

export const useSSE = (options: UseSSEOptions = {}) => {
    const [isConnected, setIsConnected] = useState(false);
    const eventSourceRef = useRef<EventSource | null>(null);

    // Use refs for callbacks to avoid reconnecting when they change (e.g. inline functions)
    const onMessageRef = useRef(options.onMessage);
    const onErrorRef = useRef(options.onError);
    const onOpenRef = useRef(options.onOpen);

    // Update refs on every render
    useEffect(() => {
        onMessageRef.current = options.onMessage;
        onErrorRef.current = options.onError;
        onOpenRef.current = options.onOpen;
    });

    // Memoize connection parameters to avoid unnecessary reconnects
    // We use JSON.stringify to compare arrays/objects by value
    const connectionParams = JSON.stringify({
        eventTypes: options.eventTypes,
        integrationId: options.integrationId,
        businessId: options.businessId
    });

    const connect = useCallback(() => {
        // Parse params inside callback to use them
        const { eventTypes, integrationId, businessId } = JSON.parse(connectionParams);

        if (eventSourceRef.current) {
            eventSourceRef.current.close();
        }

        // Construct URL with query params
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

        const baseUrl = `${envPublic.API_BASE_URL}/notify/sse/order-notify`;
        // If businessId is provided in options, we might want to use the /sse/:businessID endpoint
        // But based on routes.go, /sse/:businessID is also supported.
        // Let's use the query param approach for flexibility as per the handler logic.

        const url = `${baseUrl}?${params.toString()}`;

        const eventSource = new EventSource(url);

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
            // EventSource automatically attempts to reconnect, but we can handle custom logic here if needed
        };

        // Add custom event listeners if eventTypes are specified
        if (eventTypes) {
            eventTypes.forEach((type: string) => {
                eventSource.addEventListener(type, (event) => {
                    if (onMessageRef.current) onMessageRef.current(event);
                });
            });
        }

        eventSourceRef.current = eventSource;
    }, [connectionParams]); // Only reconnect if connection parameters change

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
