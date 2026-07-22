export const CHANNEL_CODES = ['platform', 'ecommerce'] as const;
export const SERVICE_CODES = ['messaging', 'invoicing'] as const;
export const INTERNAL_CODES = ['internal'] as const;

export const CATEGORY_COLORS: Record<string, string> = {
    platform: '#8b5cf6',
    ecommerce: '#3b82f6',
    messaging: '#a855f7',
    invoicing: '#10b981',
    internal: '#6366f1',
};

export const INTERNAL_MODULE_RESOURCE_NAME: Record<string, string> = {
    inventory: 'Inventario',
    delivery: 'Ultima Milla',
    notifications: 'Notificaciones',
    customers: 'Clientes',
    storefront_module: 'Storefront',
    invoicing_module: 'Facturacion',
};
