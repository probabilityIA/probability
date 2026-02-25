'use client';

import { useCallback, useRef, useEffect } from 'react';
import { useSSE } from '@/shared/hooks/use-sse';
import type {
  InvoiceSSEEvent,
  InvoiceSSEEventData,
  InvoiceSSEEventType,
  CompareResponseData,
} from '../../domain/types';

const INVOICE_EVENT_TYPES: InvoiceSSEEventType[] = [
  'invoice.created',
  'invoice.failed',
  'invoice.cancelled',
  'credit_note.created',
  'bulk_job.progress',
  'bulk_job.completed',
  'invoice.compare_ready',
];

interface UseInvoiceSSEOptions {
  businessId: number;
  onInvoiceCreated?: (data: InvoiceSSEEventData) => void;
  onInvoiceFailed?: (data: InvoiceSSEEventData) => void;
  onInvoiceCancelled?: (data: InvoiceSSEEventData) => void;
  onCreditNoteCreated?: (data: InvoiceSSEEventData) => void;
  onBulkJobProgress?: (data: InvoiceSSEEventData) => void;
  onBulkJobCompleted?: (data: InvoiceSSEEventData) => void;
  onCompareReady?: (data: CompareResponseData) => void;
}

export function useInvoiceSSE(options: UseInvoiceSSEOptions) {
  const {
    businessId,
    onInvoiceCreated,
    onInvoiceFailed,
    onInvoiceCancelled,
    onCreditNoteCreated,
    onBulkJobProgress,
    onBulkJobCompleted,
    onCompareReady,
  } = options;

  // Use refs for callbacks to avoid reconnecting when they change
  const callbacksRef = useRef({
    onInvoiceCreated,
    onInvoiceFailed,
    onInvoiceCancelled,
    onCreditNoteCreated,
    onBulkJobProgress,
    onBulkJobCompleted,
    onCompareReady,
  });

  useEffect(() => {
    callbacksRef.current = {
      onInvoiceCreated,
      onInvoiceFailed,
      onInvoiceCancelled,
      onCreditNoteCreated,
      onBulkJobProgress,
      onBulkJobCompleted,
      onCompareReady,
    };
  });

  const handleMessage = useCallback((event: MessageEvent) => {
    try {
      const parsed: InvoiceSSEEvent = JSON.parse(event.data);
      const eventType = parsed.type || parsed.metadata?.event_type;
      const data = parsed.data;

      if (!data) return;

      switch (eventType) {
        case 'invoice.created':
          callbacksRef.current.onInvoiceCreated?.(data);
          break;
        case 'invoice.failed':
          callbacksRef.current.onInvoiceFailed?.(data);
          break;
        case 'invoice.cancelled':
          callbacksRef.current.onInvoiceCancelled?.(data);
          break;
        case 'credit_note.created':
          callbacksRef.current.onCreditNoteCreated?.(data);
          break;
        case 'bulk_job.progress':
          callbacksRef.current.onBulkJobProgress?.(data);
          break;
        case 'bulk_job.completed':
          callbacksRef.current.onBulkJobCompleted?.(data);
          break;
        case 'invoice.compare_ready':
          callbacksRef.current.onCompareReady?.(data as unknown as CompareResponseData);
          break;
      }
    } catch {
      // Ignore non-JSON messages (heartbeats, etc.)
    }
  }, []);

  const { isConnected } = useSSE({
    businessId,
    eventTypes: INVOICE_EVENT_TYPES,
    onMessage: handleMessage,
  });

  return { isConnected };
}
