'use server';

import { cookies } from 'next/headers';
import type {
    OrdersComparePage,
    OrdersCompareQuery,
    OrdersApplyResult,
} from '../../domain/orders-compare-types';

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

async function readBody(response: Response) {
    const text = await response.text();
    try {
        return text ? JSON.parse(text) : {};
    } catch {
        return { message: text };
    }
}

export async function compareChannelOrdersAction(query: OrdersCompareQuery): Promise<OrdersComparePage> {
    const params = new URLSearchParams({ integration_id: String(query.integrationId) });
    if (query.businessId) params.set('business_id', String(query.businessId));
    if (query.from) params.set('from', query.from);
    if (query.to) params.set('to', query.to);
    if (query.page) params.set('page', String(query.page));
    if (query.pageSize) params.set('page_size', String(query.pageSize));
    if (query.limit) params.set('limit', String(query.limit));
    if (query.onlyDiff) params.set('only_diff', 'true');
    if (query.search) params.set('q', query.search);

    const response = await fetch(`${API_BASE_URL}/orders-compare?${params.toString()}`, {
        method: 'GET',
        headers: await authHeaders(),
        cache: 'no-store',
    });

    const body = await readBody(response);
    if (!response.ok || !body?.success) {
        throw new Error(body?.error || body?.message || 'No se pudieron leer las ordenes del canal');
    }
    return body.data as OrdersComparePage;
}

export async function applyChannelOrdersAction(
    integrationId: number,
    externalIds: string[],
    businessId?: number,
): Promise<OrdersApplyResult> {
    const params = new URLSearchParams();
    if (businessId) params.set('business_id', String(businessId));
    const suffix = params.toString() ? `?${params.toString()}` : '';

    const response = await fetch(`${API_BASE_URL}/orders-compare/apply${suffix}`, {
        method: 'POST',
        headers: await authHeaders(),
        body: JSON.stringify({ integration_id: integrationId, external_ids: externalIds }),
        cache: 'no-store',
    });

    const body = await readBody(response);
    if (!response.ok || !body?.success) {
        throw new Error(body?.error || body?.message || 'No se pudieron crear las ordenes en Probability');
    }
    return body.data as OrdersApplyResult;
}
