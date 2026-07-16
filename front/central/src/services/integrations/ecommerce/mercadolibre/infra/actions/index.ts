'use server';

import { cookies } from 'next/headers';

const API_BASE_URL = process.env.API_BASE_URL || 'http://localhost:3050/api/v1';

async function postWithAuth(path: string, body: Record<string, unknown>) {
    const cookieStore = await cookies();
    const sessionToken = cookieStore.get('session_token')?.value;
    const businessToken = cookieStore.get('business_token')?.value;

    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
    if (businessToken) headers['X-Business-Token'] = businessToken;

    const response = await fetch(`${API_BASE_URL}${path}`, {
        method: 'POST',
        headers,
        body: JSON.stringify(body),
    });

    const text = await response.text();
    let data: any = {};
    try {
        data = text ? JSON.parse(text) : {};
    } catch {
        data = { message: text };
    }

    if (!response.ok) {
        return { success: false, message: data.error || data.message || `Error ${response.status}` };
    }
    return { success: true, ...data };
}

export async function reconcileMeliProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/integrations/meli/products/reconcile', body);
}

export async function associateMeliProductsAction(integrationId: number, businessId?: number, skus?: string[]) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    if (skus && skus.length > 0) body.skus = skus;
    return postWithAuth('/integrations/meli/products/associate', body);
}

export async function applyMeliProductsAction(integrationId: number, direction: 'to_meli' | 'to_probability', businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId, direction };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/integrations/meli/products/apply', body);
}

export async function syncMeliInventoryAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/integrations/meli/inventory/sync', body);
}
