export const PLATFORM_CODE = 'platform';
export const CHANNEL_CODES = ['ecommerce'] as const;
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

export { channelBrand } from '@/shared/utils/channel-brand';
export type { ChannelBrand } from '@/shared/utils/channel-brand';

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
