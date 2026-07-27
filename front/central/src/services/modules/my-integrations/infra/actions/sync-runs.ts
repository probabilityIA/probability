'use server';

import { cookies } from 'next/headers';
import type { SyncRunDetail, SyncRunDetailPage, SyncRunDetailQuery, SyncRunRecord } from '../../domain/types';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:3050/api/v1';

async function authHeaders(): Promise<Record<string, string>> {
    const cookieStore = await cookies();
    const sessionToken = cookieStore.get('session_token')?.value;
    const businessToken = cookieStore.get('business_token')?.value;

    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
    if (businessToken) headers['X-Business-Token'] = businessToken;
    return headers;
}

export async function listSyncRunsAction(businessId?: number): Promise<SyncRunRecord[]> {
    try {
        const query = businessId ? `?business_id=${businessId}` : '';
        const response = await fetch(`${API_BASE_URL}/integrations/sync-runs${query}`, {
            method: 'GET',
            headers: await authHeaders(),
            cache: 'no-store',
        });
        if (!response.ok) return [];
        const body = await response.json();
        return Array.isArray(body?.data) ? (body.data as SyncRunRecord[]) : [];
    } catch {
        return [];
    }
}

const EMPTY_DETAIL_PAGE: SyncRunDetailPage = { items: [], total: 0, page: 1, page_size: 50, total_pages: 0 };

export async function listSyncRunItemsAction(query: SyncRunDetailQuery): Promise<SyncRunDetailPage> {
    try {
        const params = new URLSearchParams({
            integration_id: String(query.integration_id),
            kind: query.kind,
            page: String(query.page ?? 1),
            page_size: String(query.page_size ?? 50),
        });
        if (query.group) params.set('group', query.group);
        if (query.q) params.set('q', query.q);
        if (query.business_id) params.set('business_id', String(query.business_id));

        const response = await fetch(`${API_BASE_URL}/integrations/sync-runs/items?${params.toString()}`, {
            method: 'GET',
            headers: await authHeaders(),
            cache: 'no-store',
        });
        if (!response.ok) return EMPTY_DETAIL_PAGE;
        const body = await response.json();
        return {
            items: Array.isArray(body?.data) ? (body.data as SyncRunDetail[]) : [],
            total: Number(body?.total) || 0,
            page: Number(body?.page) || 1,
            page_size: Number(body?.page_size) || 50,
            total_pages: Number(body?.total_pages) || 0,
        };
    } catch {
        return EMPTY_DETAIL_PAGE;
    }
}
