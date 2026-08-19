'use server';

import type { CompareInventoryOptions } from '@/services/modules/my-integrations/domain/types';

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

export async function syncSiigoInventoryAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/invoicing/inventory/sync', body);
}

export async function compareSiigoInventoryAction(integrationId: number, businessId?: number, page = 1, pageSize = 100, skus?: string[], opciones?: CompareInventoryOptions) {
    const body: Record<string, unknown> = { integration_id: integrationId, page, page_size: pageSize };
    if (businessId) body.business_id = businessId;
    if (skus && skus.length > 0) body.skus = skus;
    if (opciones?.source) body.source = opciones.source;
    if (opciones?.only_diff) body.only_diff = true;
    if (opciones?.q) body.q = opciones.q;
    return postWithAuth('/siigo/inventory/compare', body);
}

export async function listSiigoWarehousesAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/invoicing/inventory/siigo-warehouses', body);
}

export async function reconcileSiigoProductsAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/siigo/products/reconcile', body);
}

export async function startSiigoProductsReconcileAction(integrationId: number, businessId?: number) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    return postWithAuth('/siigo/products/reconcile/start', body);
}

export async function applySiigoProductsAction(integrationId: number, businessId?: number, skus?: string[]) {
    const body: Record<string, unknown> = { integration_id: integrationId };
    if (businessId) body.business_id = businessId;
    if (skus && skus.length > 0) body.skus = skus;
    return postWithAuth('/siigo/products/apply', body);
}

export interface SiigoCatalogItem {
    id: number;
    code?: string;
    name: string;
    detail?: string;
    percent?: string;
}

export interface SiigoCatalogs {
    document_types_fv: SiigoCatalogItem[];
    document_types_nc: SiigoCatalogItem[];
    document_types_rc: SiigoCatalogItem[];
    document_types_cc: SiigoCatalogItem[];
    payment_types_fv: SiigoCatalogItem[];
    payment_types_rc: SiigoCatalogItem[];
    sellers: SiigoCatalogItem[];
    taxes: SiigoCatalogItem[];
    cost_centers: SiigoCatalogItem[];
    warehouses: SiigoCatalogItem[];
    errors?: string[];
}

export async function getSiigoCatalogsAction(integrationId: number): Promise<{ success: boolean; message?: string; data?: SiigoCatalogs }> {
    const cookieStore = await cookies();
    const sessionToken = cookieStore.get('session_token')?.value;
    const businessToken = cookieStore.get('business_token')?.value;

    const headers: Record<string, string> = { 'Content-Type': 'application/json' };
    if (sessionToken) headers['Authorization'] = `Bearer ${sessionToken}`;
    if (businessToken) headers['X-Business-Token'] = businessToken;

    const response = await fetch(`${API_BASE_URL}/siigo/catalogs?integration_id=${integrationId}`, {
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
    return { success: true, data: data.data as SiigoCatalogs };
}
