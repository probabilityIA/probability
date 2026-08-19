'use server';

import { cookies } from 'next/headers';
import type { CompareInventoryOptions } from '@/services/modules/my-integrations/domain/types';

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

export async function syncTiendanubeProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/tiendanube/products/sync', body);
}

export async function reconcileTiendanubeProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/tiendanube/products/reconcile', body);
}

export async function applyTiendanubeProductsAction(
    integrationId: number,
    direction: 'to_tiendanube' | 'to_probability',
    businessId?: number,
    mode: 'create' | 'update' = 'create',
    skus?: string[],
) {
    const body: Record<string, unknown> = { integration_id: integrationId, direction, mode };
    if (skus && skus.length > 0) body.skus = skus;
    if (businessId) body.business_id = businessId;
    return postWithAuth('/tiendanube/products/apply', body);
}

export async function associateTiendanubeProductsAction(integrationId: number, businessId?: number, skus?: string[]) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    if (skus && skus.length > 0) body.skus = skus;
    return postWithAuth('/tiendanube/products/associate', body);
}

export async function syncTiendanubeInventoryAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/tiendanube/inventory/sync', body);
}

export async function compareTiendanubeInventoryAction(
    integrationId: number,
    businessId?: number,
    page = 1,
    pageSize = 100,
    skus?: string[],
    opciones?: CompareInventoryOptions,
) {
    const body: Record<string, unknown> = { integration_id: integrationId, page, page_size: pageSize };
    if (businessId) body.business_id = businessId;
    if (skus && skus.length > 0) body.skus = skus;
    if (opciones?.source) body.source = opciones.source;
    if (opciones?.only_diff) body.only_diff = true;
    if (opciones?.q) body.q = opciones.q;
    return postWithAuth('/tiendanube/inventory/compare', body);
}
