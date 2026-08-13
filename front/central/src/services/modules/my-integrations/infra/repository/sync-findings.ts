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

export type FixSide = 'channel' | 'probability';
export type FindingPattern = 'spacing' | 'transposed' | 'one_char';

export interface FindingItem {
    sku: string;
    name?: string;
    detail?: string;
    channels?: string[];
    counterpart_sku?: string;
    counterpart_name?: string;
    channel_qty?: number;
    own_qty?: number;
    fix_side?: FixSide;
    pattern?: FindingPattern;
    present_in?: string[];
    missing_in?: string[];
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

export interface MatrixColumn {
    integration_id: number;
    name: string;
    code: string;
    image_url?: string;
    is_sales: boolean;
}

export interface MatrixCell {
    integration_id: number;
    present: boolean;
    sku?: string;
    barcode?: string;
    external_id?: string;
    variant_id?: string;
    sku_matches: boolean;
    sku_known: boolean;
}

export interface MatrixRow {
    product_id: string;
    sku: string;
    barcode?: string;
    name?: string;
    image_url?: string;
    cells: MatrixCell[];
}

export interface MatrixFilters {
    search?: string;
    presentIn?: number[];
    missingIn?: number[];
}

export interface MatrixPage {
    columns: MatrixColumn[];
    rows: MatrixRow[];
    total: number;
    page: number;
    total_pages: number;
}

const EMPTY_MATRIX: MatrixPage = { columns: [], rows: [], total: 0, page: 1, total_pages: 0 };

function matrixParams(businessId?: number, filters: MatrixFilters = {}): URLSearchParams {
    const params = new URLSearchParams();
    if (businessId) params.set('business_id', String(businessId));
    if (filters.search) params.set('q', filters.search);
    if (filters.presentIn?.length) params.set('present_in', filters.presentIn.join(','));
    if (filters.missingIn?.length) params.set('missing_in', filters.missingIn.join(','));
    return params;
}

export async function fetchMatchMatrix(
    businessId?: number,
    page = 1,
    filters: MatrixFilters = {},
    signal?: AbortSignal,
    pageSize = 10,
): Promise<MatrixPage> {
    const params = matrixParams(businessId, filters);
    params.set('page', String(page));
    params.set('page_size', String(pageSize));

    const response = await fetch(`/internal/sync-matrix?${params.toString()}`, { signal, cache: 'no-store' });
    if (!response.ok) return EMPTY_MATRIX;

    const body = await response.json();
    return {
        columns: Array.isArray(body?.columns) ? body.columns : [],
        rows: Array.isArray(body?.data) ? body.data : [],
        total: Number(body?.total) || 0,
        page: Number(body?.page) || 1,
        total_pages: Number(body?.total_pages) || 0,
    };
}

export function matchMatrixCsvUrl(businessId?: number, filters: MatrixFilters = {}): string {
    const params = matrixParams(businessId, filters);
    params.set('format', 'csv');
    return `/internal/sync-matrix?${params.toString()}`;
}
