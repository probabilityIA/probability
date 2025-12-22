import { useEffect, useRef, useState, useCallback } from 'react';
import { envPublic } from '@/shared/config/env';

interface IntegrationEventData {
    id: string;
    type: string;
    integration_id: number;
    business_id?: number;
    timestamp: string;
    data: any;
    metadata?: Record<string, any>;
}

interface UseIntegrationEventsOptions {
    businessId?: number;
    integrationId?: number;
    eventTypes?: string[];
    onOrderCreated?: (event: IntegrationEventData) => void;
    onOrderUpdated?: (event: IntegrationEventData) => void;
    onOrderRejected?: (event: IntegrationEventData) => void;
    onSyncStarted?: (event: IntegrationEventData) => void;
    onSyncCompleted?: (event: IntegrationEventData) => void;
    onSyncFailed?: (event: IntegrationEventData) => void;
    onError?: (error: Event) => void;
    onOpen?: (event: Event) => void;
}

export const useIntegrationEvents = (options: UseIntegrationEventsOptions = {}) => {
    const [isConnected, setIsConnected] = useState(false);
    const eventSourceRef = useRef<EventSource | null>(null);

    // Use refs for callbacks to avoid reconnecting when they change
    const callbacksRef = useRef({
        onOrderCreated: options.onOrderCreated,
        onOrderUpdated: options.onOrderUpdated,
        onOrderRejected: options.onOrderRejected,
        onSyncStarted: options.onSyncStarted,
        onSyncCompleted: options.onSyncCompleted,
        onSyncFailed: options.onSyncFailed,
        onError: options.onError,
        onOpen: options.onOpen,
    });

    // Update refs on every render
    useEffect(() => {
        callbacksRef.current = {
            onOrderCreated: options.onOrderCreated,
            onOrderUpdated: options.onOrderUpdated,
            onOrderRejected: options.onOrderRejected,
            onSyncStarted: options.onSyncStarted,
            onSyncCompleted: options.onSyncCompleted,
            onSyncFailed: options.onSyncFailed,
            onError: options.onError,
            onOpen: options.onOpen,
        };
    });

    // Memoize connection parameters
    const connectionParams = JSON.stringify({
        eventTypes: options.eventTypes,
        integrationId: options.integrationId,
        businessId: options.businessId,
    });

    const connect = useCallback(() => {
        const { eventTypes, integrationId, businessId } = JSON.parse(connectionParams);

        if (eventSourceRef.current) {
            eventSourceRef.current.close();
        }

        // Construir URL para eventos de integraciones
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

        // Usar el endpoint de eventos de integraciones
        const baseUrl = businessId 
            ? `${envPublic.API_BASE_URL}/integrations/events/sse/${businessId}`
            : `${envPublic.API_BASE_URL}/integrations/events/sse`;
        
        const url = params.toString() ? `${baseUrl}?${params.toString()}` : baseUrl;

        const eventSource = new EventSource(url);

        eventSource.onopen = (event) => {
            setIsConnected(true);
            if (callbacksRef.current.onOpen) {
                callbacksRef.current.onOpen(event);
            }
        };

        eventSource.onerror = (event) => {
            setIsConnected(false);
            if (callbacksRef.current.onError) {
                callbacksRef.current.onError(event);
            }
        };

        // Manejar eventos específicos de integraciones
        const eventTypeHandlers: Record<string, (data: IntegrationEventData) => void> = {
            'integration.sync.order.created': (data) => {
                if (callbacksRef.current.onOrderCreated) {
                    callbacksRef.current.onOrderCreated(data);
                }
            },
            'integration.sync.order.updated': (data) => {
                if (callbacksRef.current.onOrderUpdated) {
                    callbacksRef.current.onOrderUpdated(data);
                }
            },
            'integration.sync.order.rejected': (data) => {
                if (callbacksRef.current.onOrderRejected) {
                    callbacksRef.current.onOrderRejected(data);
                }
            },
            'integration.sync.started': (data) => {
                if (callbacksRef.current.onSyncStarted) {
                    callbacksRef.current.onSyncStarted(data);
                }
            },
            'integration.sync.completed': (data) => {
                if (callbacksRef.current.onSyncCompleted) {
                    callbacksRef.current.onSyncCompleted(data);
                }
            },
            'integration.sync.failed': (data) => {
                if (callbacksRef.current.onSyncFailed) {
                    callbacksRef.current.onSyncFailed(data);
                }
            },
        };

        // Agregar listeners para cada tipo de evento
        Object.keys(eventTypeHandlers).forEach((eventType) => {
            eventSource.addEventListener(eventType, (event: MessageEvent) => {
                // #region agent log
                fetch('http://127.0.0.1:7246/ingest/75c49945-688b-42a7-9abc-0a24de34a930', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        location: 'useIntegrationEvents.ts:addEventListener',
                        message: `Event received: ${eventType}`,
                        data: { eventType, eventData: event.data, hasHandler: !!eventTypeHandlers[eventType] },
                        timestamp: Date.now(),
                        sessionId: 'debug-session',
                        runId: 'run1',
                        hypothesisId: 'F'
                    })
                }).catch(() => {});
                // #endregion
                try {
                    const data = JSON.parse(event.data);
                    // #region agent log
                    fetch('http://127.0.0.1:7246/ingest/75c49945-688b-42a7-9abc-0a24de34a930', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            location: 'useIntegrationEvents.ts:addEventListener',
                            message: `Event parsed successfully: ${eventType}`,
                            data: { eventType, parsedData: data, hasCallback: !!callbacksRef.current.onOrderUpdated },
                            timestamp: Date.now(),
                            sessionId: 'debug-session',
                            runId: 'run1',
                            hypothesisId: 'F'
                        })
                    }).catch(() => {});
                    // #endregion
                    eventTypeHandlers[eventType](data);
                    // #region agent log
                    fetch('http://127.0.0.1:7246/ingest/75c49945-688b-42a7-9abc-0a24de34a930', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            location: 'useIntegrationEvents.ts:addEventListener',
                            message: `Handler called for: ${eventType}`,
                            data: { eventType },
                            timestamp: Date.now(),
                            sessionId: 'debug-session',
                            runId: 'run1',
                            hypothesisId: 'F'
                        })
                    }).catch(() => {});
                    // #endregion
                } catch (err) {
                    console.error(`Error parsing ${eventType} event:`, err);
                    // #region agent log
                    fetch('http://127.0.0.1:7246/ingest/75c49945-688b-42a7-9abc-0a24de34a930', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            location: 'useIntegrationEvents.ts:addEventListener',
                            message: `Error parsing event: ${eventType}`,
                            data: { eventType, error: String(err) },
                            timestamp: Date.now(),
                            sessionId: 'debug-session',
                            runId: 'run1',
                            hypothesisId: 'F'
                        })
                    }).catch(() => {});
                    // #endregion
                }
            });
        });

        // También manejar mensajes genéricos
        eventSource.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                // Si el tipo está en los handlers, ya fue procesado arriba
                // Si no, intentar procesarlo según el tipo
                if (data.type && eventTypeHandlers[data.type]) {
                    eventTypeHandlers[data.type](data);
                }
            } catch (err) {
                console.error('Error parsing integration event:', err);
            }
        };

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
