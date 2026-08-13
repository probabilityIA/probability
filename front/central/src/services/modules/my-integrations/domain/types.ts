export const CHANNEL_CODES = ['platform', 'ecommerce'] as const;
export const SERVICE_CODES = ['messaging', 'invoicing'] as const;
export const INTERNAL_CODES = ['internal'] as const;

export const CHANNELS_COLOR = '#3b82f6';

export const CATEGORY_COLORS: Record<string, string> = {
    platform: '#8b5cf6',
    ecommerce: '#3b82f6',
    messaging: '#a855f7',
    invoicing: '#10b981',
    internal: '#6366f1',
};

export interface ChannelBrand {
    chip: string;
    dot: string;
}

const NEUTRAL_BRAND: ChannelBrand = {
    chip: 'border-gray-200 bg-white text-gray-600 dark:border-gray-600 dark:bg-gray-800 dark:text-gray-300',
    dot: '#9ca3af',
};

const CHANNEL_BRANDS: Record<string, ChannelBrand> = {
    'mercado libre': {
        chip: 'border-amber-300 bg-amber-50 text-amber-800 dark:border-amber-500/50 dark:bg-amber-900/30 dark:text-amber-200',
        dot: '#ffe600',
    },
    woocommerce: {
        chip: 'border-purple-300 bg-purple-50 text-purple-800 dark:border-purple-500/50 dark:bg-purple-900/30 dark:text-purple-200',
        dot: '#7f54b3',
    },
    siigo: {
        chip: 'border-sky-300 bg-sky-50 text-sky-800 dark:border-sky-500/50 dark:bg-sky-900/30 dark:text-sky-200',
        dot: '#0ea5e9',
    },
    shopify: {
        chip: 'border-lime-300 bg-lime-50 text-lime-800 dark:border-lime-500/50 dark:bg-lime-900/30 dark:text-lime-200',
        dot: '#95bf47',
    },
    vtex: {
        chip: 'border-pink-300 bg-pink-50 text-pink-800 dark:border-pink-500/50 dark:bg-pink-900/30 dark:text-pink-200',
        dot: '#f71963',
    },
    jumpseller: {
        chip: 'border-orange-300 bg-orange-50 text-orange-800 dark:border-orange-500/50 dark:bg-orange-900/30 dark:text-orange-200',
        dot: '#f97316',
    },
};

export function channelBrand(channelCode?: string): ChannelBrand {
    if (!channelCode) return NEUTRAL_BRAND;
    return CHANNEL_BRANDS[channelCode.trim().toLowerCase()] ?? NEUTRAL_BRAND;
}

export type SyncRunKind = 'inventory' | 'products';

export interface SyncRunDetail {
    sku: string;
    label: string;
    tone: 'ok' | 'warn' | 'error';
    group?: string;
    matched_by?: string;
    matched_value?: string;
    parent_ref?: string;
    parent_label?: string;
    variant_label?: string;
}

export interface SyncRunDetailQuery {
    integration_id: number;
    kind: SyncRunKind;
    group?: string;
    q?: string;
    page?: number;
    page_size?: number;
    business_id?: number;
}

export interface SyncRunDetailPage {
    items: SyncRunDetail[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface SyncRunRecord {
    integration_id: number;
    kind: SyncRunKind;
    status: string;
    message?: string;
    finished_at?: string;
    total: number;
    updated: number;
    unchanged: number;
    skipped: number;
    failed: number;
    matched: number;
    not_associated: number;
    only_in_probability: number;
    only_in_channel: number;
    channel_no_sku?: number;
    sku_changed?: number;
    sku_typo?: number;
    detail: SyncRunDetail[];
}

export interface SyncRunPayload {
    integration_id: number;
    business_id?: number;
    kind: SyncRunKind;
    status?: string;
    message?: string;
    total?: number;
    updated?: number;
    unchanged?: number;
    skipped?: number;
    failed?: number;
    matched?: number;
    not_associated?: number;
    only_in_probability?: number;
    only_in_channel?: number;
    channel_no_sku?: number;
    sku_changed?: number;
    sku_typo?: number;
    detail?: SyncRunDetail[];
}

export const INTERNAL_MODULE_RESOURCE_NAME: Record<string, string> = {
    inventory: 'Inventario',
    delivery: 'Ultima Milla',
    notifications: 'Notificaciones',
    customers: 'Clientes',
    storefront_module: 'Storefront',
    invoicing_module: 'Facturacion',
};

export type FindingSeverity = 'error' | 'warn' | 'info';

export interface Finding {
    code: string;
    severity: FindingSeverity;
    title: string;
    detail: string;
    count: number;
    channels?: string[];
}

export interface FindingChannelSummary {
    integration_id: number;
    integration_name: string;
    channel_code: string;
    matched: number;
    not_associated: number;
    only_in_channel: number;
    channel_no_sku: number;
    sku_changed: number;
    sku_typo: number;
    compared_at?: string;
}

export interface FindingsReport {
    total: number;
    findings: Finding[];
    channels: FindingChannelSummary[];
}
