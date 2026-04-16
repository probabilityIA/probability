// Constantes compartidas del flujo de estados de ordenes v2
// Mirror de: back/central/services/modules/orders/internal/domain/entities/order_status.go

export interface StatusStep {
    code: string;
    name: string;
    description: string;
    color: string;
    category: 'initial' | 'warehouse' | 'assignment' | 'transit' | 'delivery' | 'return' | 'final' | 'issue';
}

export interface MetadataField {
    key: string;
    label: string;
    type: 'text' | 'textarea';
    required: boolean;
    placeholder?: string;
}

// ═══════════════════════════════════════════════════════════════
// STATUS_FLOW — todos los estados del sistema con sus propiedades
// ═══════════════════════════════════════════════════════════════

export const STATUS_FLOW: StatusStep[] = [
    // Iniciales
    { code: 'pending', name: 'Pendiente', description: 'Orden recibida, pendiente de procesamiento', color: '#F59E0B', category: 'initial' },
    { code: 'on_hold', name: 'En Espera', description: 'Orden pausada temporalmente (ej: esperando confirmacion de pago). Puede volver a Pendiente.', color: '#6B7280', category: 'initial' },

    // Almacen
    { code: 'picking', name: 'Seleccionando productos', description: 'Se estan seleccionando los productos del inventario', color: '#3B82F6', category: 'warehouse' },
    { code: 'inventory_issue', name: 'Novedad de inventario', description: 'Problema con el inventario (sin stock, producto danado). Puede volver a Picking.', color: '#FB923C', category: 'issue' },
    { code: 'packing', name: 'Empacando', description: 'El pedido esta siendo empacado para despacho', color: '#6366F1', category: 'warehouse' },
    { code: 'ready_to_ship', name: 'Listo para despacho', description: 'El pedido esta empacado y listo para ser recogido por el transportista', color: '#8B5CF6', category: 'warehouse' },

    // Asignacion
    { code: 'assigned_to_driver', name: 'Asignado a piloto', description: 'Un conductor/piloto ha sido asignado para recoger el pedido', color: '#A855F7', category: 'assignment' },
    { code: 'picked_up', name: 'Recogido', description: 'El piloto recogio el pedido del almacen', color: '#D946EF', category: 'assignment' },

    // Transito
    { code: 'in_transit', name: 'En camino', description: 'El pedido esta en camino hacia el destino del cliente', color: '#EC4899', category: 'transit' },
    { code: 'out_for_delivery', name: 'En reparto final', description: 'El pedido esta en la ultima milla, proximo a ser entregado', color: '#F43F5E', category: 'transit' },

    // Resultado de entrega
    { code: 'delivered', name: 'Entregada', description: 'El pedido fue entregado exitosamente al cliente', color: '#22C55E', category: 'delivery' },
    { code: 'delivery_novelty', name: 'Novedad de entrega', description: 'Hubo un problema en la entrega (cliente ausente, direccion incorrecta). Se puede reintentar.', color: '#F97316', category: 'issue' },
    { code: 'delivery_failed', name: 'Entrega fallida', description: 'No se pudo entregar el pedido. Se inicia proceso de devolucion.', color: '#EF4444', category: 'issue' },
    { code: 'rejected', name: 'Rechazado', description: 'El cliente rechazo el pedido. Se inicia proceso de devolucion.', color: '#DC2626', category: 'issue' },

    // Devoluciones
    { code: 'return_in_transit', name: 'Devolucion en camino', description: 'El pedido esta siendo devuelto al almacen', color: '#F59E0B', category: 'return' },
    { code: 'returned', name: 'Devuelto', description: 'El pedido fue devuelto al almacen exitosamente', color: '#EAB308', category: 'return' },

    // Finales
    { code: 'completed', name: 'Completada', description: 'Orden completada exitosamente. Entregada y sin novedades.', color: '#16A34A', category: 'final' },
    { code: 'cancelled', name: 'Cancelada', description: 'Orden cancelada. Se puede cancelar desde cualquier estado no terminal.', color: '#EF4444', category: 'final' },
    { code: 'refunded', name: 'Reembolsada', description: 'Se realizo reembolso al cliente. Posible desde Entregada, Completada o Devuelto.', color: '#7C3AED', category: 'final' },
];

// ═══════════════════════════════════════════════════════════════
// CATEGORY_LABELS — labels para agrupar estados por fase
// ═══════════════════════════════════════════════════════════════

export const CATEGORY_LABELS: Record<string, { label: string; icon: string }> = {
    initial: { label: 'Inicio', icon: '1' },
    warehouse: { label: 'Almacen', icon: '2' },
    issue: { label: 'Novedades', icon: '!' },
    assignment: { label: 'Asignacion', icon: '3' },
    transit: { label: 'Transito', icon: '4' },
    delivery: { label: 'Entrega', icon: '5' },
    return: { label: 'Devoluciones', icon: '6' },
    final: { label: 'Estados Finales', icon: '7' },
};

export const CATEGORY_ORDER = ['initial', 'warehouse', 'assignment', 'transit', 'delivery', 'return', 'final', 'issue'];

// ═══════════════════════════════════════════════════════════════
// VALID_TRANSITIONS — mirror exacto de order_status.go validTransitions
// ═══════════════════════════════════════════════════════════════

const VALID_TRANSITIONS: Record<string, string[]> = {
    pending:            ['picking', 'on_hold'],
    on_hold:            ['pending', 'picking'],
    picking:            ['packing', 'inventory_issue', 'on_hold'],
    inventory_issue:    ['picking'],
    packing:            ['ready_to_ship', 'on_hold'],
    ready_to_ship:      ['assigned_to_driver', 'on_hold'],
    assigned_to_driver: ['picked_up'],
    picked_up:          ['in_transit'],
    in_transit:         ['out_for_delivery'],
    out_for_delivery:   ['delivered', 'delivery_novelty', 'rejected', 'delivery_failed'],
    delivered:          ['completed', 'refunded', 'return_in_transit'],
    delivery_novelty:   ['assigned_to_driver', 'out_for_delivery', 'delivery_failed', 'return_in_transit'],
    delivery_failed:    ['return_in_transit'],
    rejected:           ['return_in_transit'],
    return_in_transit:  ['returned'],
    returned:           ['refunded'],
    completed:          ['refunded'],
    failed:             [],
    // Terminales — sin transiciones
    cancelled:          [],
    refunded:           [],
};

// ═══════════════════════════════════════════════════════════════
// TERMINAL_STATUSES
// ═══════════════════════════════════════════════════════════════

const TERMINAL_STATUSES = new Set(['cancelled', 'refunded']);

// ═══════════════════════════════════════════════════════════════
// STATUS_METADATA_FIELDS — campos de metadata requeridos por estado destino
// ═══════════════════════════════════════════════════════════════

export const STATUS_METADATA_FIELDS: Record<string, MetadataField[]> = {
    assigned_to_driver: [
        { key: 'driver_name', label: 'Nombre del piloto', type: 'text', required: true, placeholder: 'Ej: Carlos Ramirez' },
    ],
    in_transit: [
        { key: 'tracking_number', label: 'Numero de seguimiento', type: 'text', required: false, placeholder: 'Ej: TRK-001' },
        { key: 'tracking_link', label: 'Link de seguimiento', type: 'text', required: false, placeholder: 'https://...' },
    ],
    cancelled: [
        { key: 'reason', label: 'Razon de cancelacion', type: 'textarea', required: true, placeholder: 'Explica por que se cancela esta orden...' },
    ],
    on_hold: [
        { key: 'reason', label: 'Razon de pausa', type: 'textarea', required: true, placeholder: 'Explica por que se pausa esta orden...' },
    ],
    delivery_novelty: [
        { key: 'reason', label: 'Descripcion de la novedad', type: 'textarea', required: true, placeholder: 'Describe la novedad de entrega...' },
    ],
    delivery_failed: [
        { key: 'reason', label: 'Razon de fallo', type: 'textarea', required: false, placeholder: 'Describe por que fallo la entrega...' },
    ],
    rejected: [
        { key: 'reason', label: 'Razon del rechazo', type: 'textarea', required: false, placeholder: 'Describe por que el cliente rechazo...' },
    ],
    inventory_issue: [
        { key: 'notes', label: 'Notas del problema', type: 'textarea', required: false, placeholder: 'Describe el problema de inventario...' },
    ],
    return_in_transit: [
        { key: 'tracking_number', label: 'Numero de seguimiento devolucion', type: 'text', required: false, placeholder: 'Ej: RET-001' },
    ],
};

// ═══════════════════════════════════════════════════════════════
// HELPERS
// ═══════════════════════════════════════════════════════════════

/** Retorna true si el estado es terminal (cancelled, refunded) */
export function isTerminalStatus(status: string): boolean {
    return TERMINAL_STATUSES.has(status);
}

/** Obtiene un StatusStep por codigo */
export function getStatusByCode(code: string): StatusStep | undefined {
    return STATUS_FLOW.find(s => s.code === code);
}

/**
 * Retorna los estados validos a los que se puede transicionar desde el estado actual.
 * Incluye 'cancelled' para cualquier estado no-terminal (regla especial del backend).
 * Retorna [] para estados terminales.
 */
export function getValidTransitions(currentStatus: string): StatusStep[] {
    if (isTerminalStatus(currentStatus)) return [];

    const directTargets = VALID_TRANSITIONS[currentStatus] || [];

    // cancelled es accesible desde cualquier estado no-terminal
    const allTargets = directTargets.includes('cancelled')
        ? directTargets
        : [...directTargets, 'cancelled'];

    return allTargets
        .map(code => STATUS_FLOW.find(s => s.code === code))
        .filter((s): s is StatusStep => s !== undefined);
}
