'use server';

import { cookies } from 'next/headers';
import type { SyncRunRecord } from '../../domain/types';

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
