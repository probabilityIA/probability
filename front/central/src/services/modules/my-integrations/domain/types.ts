/**
 * Códigos de categorías por nivel jerárquico.
 * Canales: donde se originan las órdenes (paralelos)
 * Servicios: donde se procesan (independientes desde el hub)
 */
export const CHANNEL_CODES = ['platform', 'ecommerce'] as const;
export const SERVICE_CODES = ['messaging', 'invoicing', 'shipping', 'payment'] as const;

export const CATEGORY_ICONS: Record<string, string> = {
    platform: '🧩',
    ecommerce: '🛒',
    invoicing: '🧾',
    messaging: '💬',
    payment: '💳',
    shipping: '🚚',
};
