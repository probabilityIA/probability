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
