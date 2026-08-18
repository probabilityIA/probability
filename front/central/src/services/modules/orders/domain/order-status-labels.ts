import { getStatusByCode } from './order-status-transitions';

const PLATFORM_STATUS_LABELS: Record<string, string> = {
    authorized: 'Autorizada',
    confirmed: 'Confirmada',
    failed: 'Fallida',
    fulfilled: 'Despachada',
    paid: 'Pagada',
    partially_paid: 'Pago parcial',
    partially_refunded: 'Reembolso parcial',
    processing: 'En proceso',
    ready_for_pickup: 'Lista para recoger',
    unfulfilled: 'Sin despachar',
    unpaid: 'Sin pagar',
    voided: 'Anulada',
};

export function getOrderStatusLabel(code?: string | null): string {
    if (!code) return '';
    const trimmed = code.trim();
    if (!trimmed) return '';

    const fromFlow = getStatusByCode(trimmed);
    if (fromFlow) return fromFlow.name;

    return PLATFORM_STATUS_LABELS[trimmed] || trimmed;
}

const PAYMENT_STATUSES = new Set([
    'authorized',
    'paid',
    'partially_paid',
    'unpaid',
    'refunded',
    'partially_refunded',
    'voided',
    'pending_payment',
]);

export function isPaymentStatus(status?: string | null): boolean {
    if (!status) return false;
    return PAYMENT_STATUSES.has(status.trim().toLowerCase().replace(/^wc-/, ''));
}
