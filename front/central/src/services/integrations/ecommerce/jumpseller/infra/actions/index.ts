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

async function getWithAuth(path: string) {
    const cookieStore = await cookies();
    const sessionToken = cookieStore.get('session_token')?.value;
    const businessToken = cookieStore.get('business_token')?.value;

    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
    if (businessToken) headers['X-Business-Token'] = businessToken;

    const response = await fetch(`${API_BASE_URL}${path}`, {
        method: 'GET',
        headers,
        cache: 'no-store',
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

export async function getJumpsellerLocationsAction(integrationId: number, businessId?: number) {
    const params = new URLSearchParams({ integration_id: String(integrationId) });
    if (businessId) params.set('business_id', String(businessId));
    return getWithAuth(`/jumpseller/locations?${params.toString()}`);
}

export async function syncJumpsellerProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/jumpseller/products/sync', body);
}

export async function reconcileJumpsellerProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/jumpseller/products/reconcile', body);
}

export async function applyJumpsellerProductsAction(integrationId: number, direction: 'to_jumpseller' | 'to_probability', businessId?: number, mode: 'create' | 'update' = 'create') {
    const body: Record<string, unknown> = { integration_id: integrationId, direction, mode };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/jumpseller/products/apply', body);
}

export async function associateJumpsellerProductsAction(integrationId: number, businessId?: number, skus?: string[]) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    if (skus && skus.length > 0) body.skus = skus;
    return postWithAuth('/jumpseller/products/associate', body);
}

export async function syncJumpsellerInventoryAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/jumpseller/inventory/sync', body);
}
