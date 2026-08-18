const STATUS_LABELS: Record<string, string> = {
    pending: 'Pendiente',
    processing: 'En procesamiento',
    confirmed: 'Confirmada',
    on_hold: 'En espera',
    'on-hold': 'En espera',
    picking: 'Seleccionando productos',
    packing: 'Empacando',
    ready_to_ship: 'Listo para despacho',
    ready_for_pickup: 'Listo para recoger',
    assigned_to_driver: 'Asignado a piloto',
    picked_up: 'Recogido',
    shipped: 'Enviada',
    in_transit: 'En camino',
    out_for_delivery: 'En reparto final',
    delivered: 'Entregada',
    completed: 'Completada',
    fulfilled: 'Preparada',
    unfulfilled: 'Sin preparar',
    delivery_novelty: 'Novedad de entrega',
    delivery_failed: 'Entrega fallida',
    failed: 'Fallida',
    returned: 'Devuelta',
    return_in_transit: 'Devoluci\u00f3n en camino',
    rejected: 'Rechazada',
    cancelled: 'Cancelada',
    canceled: 'Cancelada',
    inventory_issue: 'Novedad de inventario',
    authorized: 'Pago autorizado',
    paid: 'Pagada',
    partially_paid: 'Pago parcial',
    unpaid: 'Sin pagar',
    refunded: 'Reembolsada',
    partially_refunded: 'Reembolso parcial',
    voided: 'Pago anulado',
    pending_payment: 'Pago pendiente',
};

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

export function orderStatusLabel(status?: string | null): string {
    if (!status) return '-';
    const key = status.trim().toLowerCase().replace(/^wc-/, '');
    if (STATUS_LABELS[key]) return STATUS_LABELS[key];
    const words = key.replace(/[_-]+/g, ' ').trim();
    return words.charAt(0).toUpperCase() + words.slice(1);
}

export function isPaymentStatus(status?: string | null): boolean {
    if (!status) return false;
    return PAYMENT_STATUSES.has(status.trim().toLowerCase().replace(/^wc-/, ''));
}
