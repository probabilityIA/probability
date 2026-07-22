'use server';

import { getAuthToken } from '@/shared/utils/server-auth';

import { env } from '@/shared/config/env';

async function buildHeaders(): Promise<Record<string, string>> {
    const token = await getAuthToken();
    const headers: Record<string, string> = {};
    if (token) headers['Authorization'] = `Bearer ${token}`;
    return headers;
}

export interface SubscriptionType {
    id: number;
    name: string;
    code: string;
    description: string;
    price: number;
    billing_period: string;
    active: boolean;
    module_codes: string[];
    max_ecommerce_channels: number;
    created_at?: string;
    updated_at?: string;
}

export interface BusinessSubscription {
    id?: number;
    business_id: number;
    subscription_type_id: number;
    subscription_type_name: string;
    months: number;
    amount: number;
    start_date?: string;
    end_date?: string;
    status: string;
    payment_reference?: string;
    notes?: string;
    created_at?: string;
}

export interface BusinessModuleOverride {
    id: number;
    business_id: number;
    module_code: string;
    granted_by_user_id: number;
    notes?: string;
    created_at?: string;
}

export async function getMySubscriptionAction(businessId?: number): Promise<{ success: boolean; data?: BusinessSubscription; error?: string }> {
    try {
        const headers = await buildHeaders();
        const url = businessId
            ? `${env.API_BASE_URL}/subscriptions/me?business_id=${businessId}`
            : `${env.API_BASE_URL}/subscriptions/me`;
        const res = await fetch(url, {
            headers,
            cache: 'no-store',
        });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function listSubscriptionTypesAction(activeOnly = false): Promise<{ success: boolean; data?: SubscriptionType[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const url = `${env.API_BASE_URL}/subscriptions/types${activeOnly ? '?active_only=true' : ''}`;
        const res = await fetch(url, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function createSubscriptionTypeAction(payload: {
    name: string;
    code: string;
    description?: string;
    price: number;
    billing_period?: string;
    module_codes: string[];
    max_ecommerce_channels?: number;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/types`, {
            method: 'POST',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err?.error || `Error ${res.status}`);
        }
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function updateSubscriptionTypeAction(id: number, payload: {
    name: string;
    description?: string;
    price: number;
    billing_period?: string;
    active: boolean;
    module_codes: string[];
    max_ecommerce_channels?: number;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/types/${id}`, {
            method: 'PUT',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err?.error || `Error ${res.status}`);
        }
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function deleteSubscriptionTypeAction(id: number): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/types/${id}`, {
            method: 'DELETE',
            headers,
        });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function getModuleCodesAction(): Promise<{ success: boolean; data?: string[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/module-codes`, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export interface ModuleInfo {
    code: string;
    name: string;
}

export async function getModuleCatalogAction(): Promise<{ success: boolean; data?: ModuleInfo[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/module-catalog`, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function getMyModulesAction(businessId?: number): Promise<{ success: boolean; data?: string[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const url = businessId
            ? `${env.API_BASE_URL}/subscriptions/my-modules?business_id=${businessId}`
            : `${env.API_BASE_URL}/subscriptions/my-modules`;
        const res = await fetch(url, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function purchaseSubscriptionAction(payload: {
    subscriptionTypeId: number;
    months: number;
}, businessId?: number): Promise<{ success: boolean; data?: BusinessSubscription; error?: string }> {
    try {
        const headers = await buildHeaders();
        const url = businessId
            ? `${env.API_BASE_URL}/subscriptions/purchase?business_id=${businessId}`
            : `${env.API_BASE_URL}/subscriptions/purchase`;
        const res = await fetch(url, {
            method: 'POST',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify({ subscription_type_id: payload.subscriptionTypeId, months: payload.months }),
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err?.error || `Error ${res.status}`);
        }
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function registerSubscriptionPaymentAction(payload: {
    businessId: number;
    subscriptionTypeId: number;
    monthsToAdd: number;
    paymentReference?: string;
    notes?: string;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/register-payment`, {
            method: 'POST',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify({
                business_id: payload.businessId,
                subscription_type_id: payload.subscriptionTypeId,
                months: payload.monthsToAdd,
                payment_reference: payload.paymentReference,
                notes: payload.notes,
            }),
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err?.error || `Error ${res.status}`);
        }
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function disableSubscriptionAction(businessId: number): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/disable?business_id=${businessId}`, {
            method: 'POST',
            headers,
        });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function listOverridesAction(businessId: number): Promise<{ success: boolean; data?: BusinessModuleOverride[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/overrides/${businessId}`, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function grantOverrideAction(payload: {
    businessId: number;
    moduleCode: string;
    notes?: string;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/overrides`, {
            method: 'POST',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify({ business_id: payload.businessId, module_code: payload.moduleCode, notes: payload.notes }),
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err?.error || `Error ${res.status}`);
        }
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function revokeOverrideAction(businessId: number, moduleCode: string): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/overrides/${businessId}/${moduleCode}`, {
            method: 'DELETE',
            headers,
        });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}
