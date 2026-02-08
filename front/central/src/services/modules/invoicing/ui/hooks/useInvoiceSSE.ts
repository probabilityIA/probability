'use client';

import { useCallback, useRef, useEffect } from 'react';
import { useSSE } from '@/shared/hooks/use-sse';
import type {
  InvoiceSSEEvent,
  InvoiceSSEEventData,
  InvoiceSSEEventType,
} from '../../domain/types';

const INVOICE_EVENT_TYPES: InvoiceSSEEventType[] = [
  'invoice.created',
  'invoice.failed',
  'invoice.cancelled',
  'credit_note.created',
  'bulk_job.progress',
  'bulk_job.completed',
];

interface UseInvoiceSSEOptions {
  businessId: number;
  onInvoiceCreated?: (data: InvoiceSSEEventData) => void;
  onInvoiceFailed?: (data: InvoiceSSEEventData) => void;
  onInvoiceCancelled?: (data: InvoiceSSEEventData) => void;
  onCreditNoteCreated?: (data: InvoiceSSEEventData) => void;
  onBulkJobProgress?: (data: InvoiceSSEEventData) => void;
  onBulkJobCompleted?: (data: InvoiceSSEEventData) => void;
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
  } = options;

  // Use refs for callbacks to avoid reconnecting when they change
  const callbacksRef = useRef({
    onInvoiceCreated,
    onInvoiceFailed,
    onInvoiceCancelled,
    onCreditNoteCreated,
    onBulkJobProgress,
    onBulkJobCompleted,
  });

  useEffect(() => {
    callbacksRef.current = {
      onInvoiceCreated,
      onInvoiceFailed,
      onInvoiceCancelled,
      onCreditNoteCreated,
      onBulkJobProgress,
      onBulkJobCompleted,
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
