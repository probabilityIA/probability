'use client';

import { useCallback, useRef, useEffect } from 'react';
import { useSSE } from '@/shared/hooks/use-sse';
import type {
  ShipmentSSEEvent,
  ShipmentSSEEventData,
  ShipmentSSEEventType,
} from '../../domain/types';

const SHIPMENT_EVENT_TYPES: ShipmentSSEEventType[] = [
  'shipment.quote_received',
  'shipment.quote_failed',
  'shipment.guide_generated',
  'shipment.guide_failed',
  'shipment.tracking_updated',
  'shipment.tracking_failed',
  'shipment.cancelled',
  'shipment.cancel_failed',
];

interface UseShipmentSSEOptions {
  businessId: number;
  onQuoteReceived?: (data: ShipmentSSEEventData) => void;
  onQuoteFailed?: (data: ShipmentSSEEventData) => void;
  onGuideGenerated?: (data: ShipmentSSEEventData) => void;
  onGuideFailed?: (data: ShipmentSSEEventData) => void;
  onTrackingUpdated?: (data: ShipmentSSEEventData) => void;
  onTrackingFailed?: (data: ShipmentSSEEventData) => void;
  onShipmentCancelled?: (data: ShipmentSSEEventData) => void;
  onCancelFailed?: (data: ShipmentSSEEventData) => void;
}

export function useShipmentSSE(options: UseShipmentSSEOptions) {
  const {
    businessId,
    onQuoteReceived,
    onQuoteFailed,
    onGuideGenerated,
    onGuideFailed,
    onTrackingUpdated,
    onTrackingFailed,
    onShipmentCancelled,
    onCancelFailed,
  } = options;

  // Use refs for callbacks to avoid reconnecting when they change
  const callbacksRef = useRef({
    onQuoteReceived,
    onQuoteFailed,
    onGuideGenerated,
    onGuideFailed,
    onTrackingUpdated,
    onTrackingFailed,
    onShipmentCancelled,
    onCancelFailed,
  });

  useEffect(() => {
    callbacksRef.current = {
      onQuoteReceived,
      onQuoteFailed,
      onGuideGenerated,
      onGuideFailed,
      onTrackingUpdated,
      onTrackingFailed,
      onShipmentCancelled,
      onCancelFailed,
    };
  });

  const handleMessage = useCallback((event: MessageEvent) => {
    try {
      const parsed: ShipmentSSEEvent = JSON.parse(event.data);
      const eventType = parsed.type || parsed.metadata?.event_type;
      const data = parsed.data;

      if (!data) return;

      switch (eventType) {
        case 'shipment.quote_received':
          callbacksRef.current.onQuoteReceived?.(data);
          break;
        case 'shipment.quote_failed':
          callbacksRef.current.onQuoteFailed?.(data);
          break;
        case 'shipment.guide_generated':
          callbacksRef.current.onGuideGenerated?.(data);
          break;
        case 'shipment.guide_failed':
          callbacksRef.current.onGuideFailed?.(data);
          break;
        case 'shipment.tracking_updated':
          callbacksRef.current.onTrackingUpdated?.(data);
          break;
        case 'shipment.tracking_failed':
          callbacksRef.current.onTrackingFailed?.(data);
          break;
        case 'shipment.cancelled':
          callbacksRef.current.onShipmentCancelled?.(data);
          break;
        case 'shipment.cancel_failed':
          callbacksRef.current.onCancelFailed?.(data);
          break;
      }
    } catch {
      // Ignore non-JSON messages (heartbeats, etc.)
    }
  }, []);

  const { isConnected } = useSSE({
    businessId,
    eventTypes: SHIPMENT_EVENT_TYPES,
    onMessage: handleMessage,
  });

  return { isConnected };
}
