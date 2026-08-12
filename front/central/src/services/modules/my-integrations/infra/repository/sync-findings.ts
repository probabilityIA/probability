import type { FindingsReport } from '../../domain/types';

const EMPTY: FindingsReport = { total: 0, findings: [], channels: [] };

export async function fetchSyncFindings(businessId?: number, signal?: AbortSignal): Promise<FindingsReport> {
    const params = new URLSearchParams();
    if (businessId) params.set('business_id', String(businessId));

    const response = await fetch(`/internal/sync-findings?${params.toString()}`, { signal, cache: 'no-store' });
    if (!response.ok) return EMPTY;

    const body = await response.json();
    const data = body?.data;
    if (!data) return EMPTY;

    return {
        total: Number(data.total) || 0,
        findings: Array.isArray(data.findings) ? data.findings : [],
        channels: Array.isArray(data.channels) ? data.channels : [],
    };
}

export interface FindingItem {
    sku: string;
    name?: string;
    detail?: string;
    channels?: string[];
}

export interface FindingItemsPage {
    items: FindingItem[];
    total: number;
    page: number;
    total_pages: number;
}

const EMPTY_ITEMS: FindingItemsPage = { items: [], total: 0, page: 1, total_pages: 0 };

export async function fetchFindingItems(
    code: string,
    businessId?: number,
    page = 1,
    search = '',
    signal?: AbortSignal,
): Promise<FindingItemsPage> {
    const params = new URLSearchParams({ code, page: String(page), page_size: '50' });
    if (businessId) params.set('business_id', String(businessId));
    if (search) params.set('q', search);

    const response = await fetch(`/internal/sync-finding-items?${params.toString()}`, { signal, cache: 'no-store' });
    if (!response.ok) return EMPTY_ITEMS;

    const body = await response.json();
    return {
        items: Array.isArray(body?.data) ? body.data : [],
        total: Number(body?.total) || 0,
        page: Number(body?.page) || 1,
        total_pages: Number(body?.total_pages) || 0,
    };
}

export function findingItemsCsvUrl(code: string, businessId?: number, search = ''): string {
    const params = new URLSearchParams({ code, format: 'csv' });
    if (businessId) params.set('business_id', String(businessId));
    if (search) params.set('q', search);
    return `/internal/sync-finding-items?${params.toString()}`;
}
