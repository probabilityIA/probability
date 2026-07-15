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

export async function syncWooProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/woocommerce/products/sync', body);
}

export async function reconcileWooProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/woocommerce/products/reconcile', body);
}

export async function applyWooProductsAction(integrationId: number, direction: 'to_woo' | 'to_probability', businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId, direction };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/woocommerce/products/apply', body);
}

export async function syncWooInventoryAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/woocommerce/inventory/sync', body);
}

export async function getWooPluginZipAction(): Promise<{ success: boolean; data?: string; message?: string }> {
    try {
        const response = await fetch(`${API_BASE_URL}/woocommerce/plugin-download`, { cache: 'no-store' });
        if (!response.ok) {
            return { success: false, message: `No se pudo descargar el plugin (Error ${response.status})` };
        }
        const buf = await response.arrayBuffer();
        return { success: true, data: Buffer.from(buf).toString('base64') };
    } catch {
        return { success: false, message: 'No se pudo conectar con el servidor para descargar el plugin' };
    }
}
