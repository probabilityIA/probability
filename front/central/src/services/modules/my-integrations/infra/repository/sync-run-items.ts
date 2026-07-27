import type { SyncRunDetail, SyncRunDetailPage, SyncRunDetailQuery } from '../../domain/types';

const EMPTY_PAGE: SyncRunDetailPage = { items: [], total: 0, page: 1, page_size: 50, total_pages: 0 };

export async function fetchSyncRunItems(
    query: SyncRunDetailQuery,
    signal?: AbortSignal,
): Promise<SyncRunDetailPage> {
    const params = new URLSearchParams({
        integration_id: String(query.integration_id),
        kind: query.kind,
        page: String(query.page ?? 1),
        page_size: String(query.page_size ?? 50),
    });
    if (query.group) params.set('group', query.group);
    if (query.q) params.set('q', query.q);
    if (query.business_id) params.set('business_id', String(query.business_id));

    const response = await fetch(`/internal/sync-run-items?${params.toString()}`, {
        signal,
        cache: 'no-store',
    });
    if (!response.ok) return EMPTY_PAGE;

    const body = await response.json();
    return {
        items: Array.isArray(body?.data) ? (body.data as SyncRunDetail[]) : [],
        total: Number(body?.total) || 0,
        page: Number(body?.page) || 1,
        page_size: Number(body?.page_size) || 50,
        total_pages: Number(body?.total_pages) || 0,
    };
}
