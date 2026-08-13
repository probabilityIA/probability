export interface DataSummaryCell {
    integration_id: number;
    integration_name: string;
    channel_code: string;
    can_fill: number;
    can_overwrite: number;
}

export interface DataSummaryRow {
    field: string;
    label: string;
    note?: string;
    cells: DataSummaryCell[];
}

export interface DataSummary {
    rows: DataSummaryRow[];
    snapshot_at?: string;
}

export type DataMode = 'fill_empty' | 'overwrite';

export interface DataConflict {
    source: string;
    integration_id?: number;
    integration_name?: string;
    count: number;
    last_change_at?: string;
}

export interface DataSample {
    product_id: string;
    sku: string;
    name?: string;
    current?: string;
    incoming?: string;
}

export interface DataPreview {
    field: string;
    mode: DataMode;
    would_fill: number;
    would_replace: number;
    conflicts: DataConflict[];
    samples: DataSample[];
}

export interface DataApplyResult {
    batch_id: string;
    field: string;
    applied: number;
    applied_at: string;
}

const EMPTY_SUMMARY: DataSummary = { rows: [] };

function url(accion: string, businessId?: number): string {
    const params = new URLSearchParams({ accion });
    if (businessId) params.set('business_id', String(businessId));
    return `/internal/channel-data?${params.toString()}`;
}

export async function fetchDataSummary(businessId?: number, signal?: AbortSignal): Promise<DataSummary> {
    try {
        const response = await fetch(url('summary', businessId), { signal, cache: 'no-store' });
        if (!response.ok) return EMPTY_SUMMARY;
        const body = await response.json();
        if (!body?.success) return EMPTY_SUMMARY;
        return {
            rows: Array.isArray(body.data) ? body.data : [],
            snapshot_at: body.snapshot_at ?? undefined,
        };
    } catch {
        return EMPTY_SUMMARY;
    }
}

export async function fetchDataPreview(
    businessId: number | undefined,
    integrationId: number,
    field: string,
    mode: DataMode,
    signal?: AbortSignal,
): Promise<DataPreview | null> {
    try {
        const response = await fetch(url('preview', businessId), {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ business_id: businessId, integration_id: integrationId, field, mode }),
            signal,
            cache: 'no-store',
        });
        if (!response.ok) return null;
        const body = await response.json();
        if (!body?.success) return null;
        return {
            field: body.field,
            mode: body.mode,
            would_fill: Number(body.would_fill) || 0,
            would_replace: Number(body.would_replace) || 0,
            conflicts: Array.isArray(body.conflicts) ? body.conflicts : [],
            samples: Array.isArray(body.samples) ? body.samples : [],
        };
    } catch {
        return null;
    }
}

export async function applyChannelData(
    businessId: number | undefined,
    integrationId: number,
    field: string,
    mode: DataMode,
): Promise<DataApplyResult | null> {
    const response = await fetch(url('apply', businessId), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ integration_id: integrationId, field, mode }),
        cache: 'no-store',
    });
    if (!response.ok) return null;
    const body = await response.json();
    if (!body?.success) return null;
    return body.data as DataApplyResult;
}

export async function undoChannelData(businessId: number | undefined, batchId: string): Promise<number | null> {
    const response = await fetch(url('undo', businessId), {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ batch_id: batchId }),
        cache: 'no-store',
    });
    if (!response.ok) return null;
    const body = await response.json();
    if (!body?.success) return null;
    return Number(body.data?.reverted) || 0;
}
