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
    business_id?: number;
    included_shipments?: number;
    shipment_overage_price?: number;
    included_invoices?: number;
    invoice_overage_price?: number;
    included_orders?: number;
    order_overage_price?: number;
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
    payment_method?: string;
    payment_reference?: string;
    notes?: string;
    created_at?: string;
    subscription_type?: SubscriptionType;
    business_subscription_status?: string;
    auto_payment_enabled?: boolean;
    payment_window_start?: string;
    payment_window_end?: string;
}

export interface BusinessModuleOverride {
    id: number;
    business_id: number;
    module_code: string;
    granted_by_user_id: number;
    notes?: string;
    expires_at?: string;
    created_at?: string;
}

export interface SubscriptionAuditLog {
    id: number;
    business_id: number;
    actor_user_id?: number;
    actor_label: string;
    action: string;
    description: string;
    created_at: string;
}

export interface AdminBusinessRow {
    id: number;
    name: string;
    code: string;
    plan_name?: string;
    status: string;
    cycle_start_date?: string;
    cycle_end_date?: string;
    last_payment_amount?: number;
    last_payment_date?: string;
    forecasted_payment?: number;
    cutoff_day?: number;
}

export interface AdminKPIs {
    active_count: number;
    expiring_soon_count: number;
    expired_or_suspended_count: number;
    mrr: number;
}

export interface SubscriptionUsage {
    plan_name: string;
    plan_price: number;
    billing_period: string;
    module_codes: string[];
    max_ecommerce_channels: number;
    cycle_start_date: string;
    cycle_end_date: string;
    included_shipments?: number;
    shipment_overage_price?: number;
    shipments_used: number;
    included_invoices?: number;
    invoice_overage_price?: number;
    invoices_used: number;
    included_orders?: number;
    order_overage_price?: number;
    orders_used: number;
    forecasted_payment?: number;
}

export async function getMySubscriptionUsageAction(businessId?: number): Promise<{ success: boolean; data?: SubscriptionUsage; error?: string }> {
    try {
        const headers = await buildHeaders();
        const url = businessId
            ? `${env.API_BASE_URL}/subscriptions/me/usage?business_id=${businessId}`
            : `${env.API_BASE_URL}/subscriptions/me/usage`;
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

export async function setAutoPaymentAction(enabled: boolean, businessId?: number): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        headers['Content-Type'] = 'application/json';
        const url = businessId
            ? `${env.API_BASE_URL}/subscriptions/me/auto-payment?business_id=${businessId}`
            : `${env.API_BASE_URL}/subscriptions/me/auto-payment`;
        const res = await fetch(url, {
            method: 'PUT',
            headers,
            body: JSON.stringify({ enabled }),
        });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        return { success: true };
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
    included_shipments?: number;
    shipment_overage_price?: number;
    included_invoices?: number;
    invoice_overage_price?: number;
    included_orders?: number;
    order_overage_price?: number;
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
    included_shipments?: number;
    shipment_overage_price?: number;
    included_invoices?: number;
    invoice_overage_price?: number;
    included_orders?: number;
    order_overage_price?: number;
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

export async function listCustomPlansAction(businessId?: number): Promise<{ success: boolean; data?: SubscriptionType[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const url = businessId
            ? `${env.API_BASE_URL}/subscriptions/custom-plans?business_id=${businessId}`
            : `${env.API_BASE_URL}/subscriptions/custom-plans`;
        const res = await fetch(url, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function createCustomPlanAction(payload: {
    name: string;
    code: string;
    description?: string;
    price: number;
    billing_period?: string;
    module_codes: string[];
    max_ecommerce_channels?: number;
    business_id: number;
    months: number;
    payment_reference?: string;
    notes?: string;
    included_shipments?: number;
    shipment_overage_price?: number;
    included_invoices?: number;
    invoice_overage_price?: number;
    included_orders?: number;
    order_overage_price?: number;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/custom-plans`, {
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

export async function updateCustomPlanAction(id: number, payload: {
    name: string;
    description?: string;
    price: number;
    billing_period?: string;
    active: boolean;
    module_codes: string[];
    max_ecommerce_channels?: number;
    included_shipments?: number;
    shipment_overage_price?: number;
    included_invoices?: number;
    invoice_overage_price?: number;
    included_orders?: number;
    order_overage_price?: number;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/custom-plans/${id}`, {
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

export async function deleteCustomPlanAction(id: number): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/custom-plans/${id}`, {
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
    paymentMethod?: string;
    paymentReference?: string;
    notes?: string;
    startDate?: string;
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
                payment_method: payload.paymentMethod,
                payment_reference: payload.paymentReference,
                notes: payload.notes,
                start_date: payload.startDate,
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

export async function editSubscriptionDatesAction(payload: {
    businessId: number;
    startDate: string;
    endDate: string;
    cutoffDay?: number;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/edit-dates`, {
            method: 'PUT',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify({
                business_id: payload.businessId,
                start_date: payload.startDate,
                end_date: payload.endDate,
                cutoff_day: payload.cutoffDay,
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
    expiresAt?: string;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/overrides`, {
            method: 'POST',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify({
                business_id: payload.businessId,
                module_code: payload.moduleCode,
                notes: payload.notes,
                expires_at: payload.expiresAt,
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

export async function reactivateSubscriptionAction(businessId: number): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/reactivate`, {
            method: 'POST',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify({ business_id: businessId }),
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

export async function extendCourtesyAction(payload: {
    businessId: number;
    days: number;
    reason: string;
}): Promise<{ success: boolean; data?: BusinessSubscription; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/extend-days`, {
            method: 'POST',
            headers: { ...headers, 'Content-Type': 'application/json' },
            body: JSON.stringify({ business_id: payload.businessId, days: payload.days, reason: payload.reason }),
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

export async function revertPaymentAction(subscriptionId: number): Promise<{ success: boolean; data?: BusinessSubscription; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/payments/${subscriptionId}/revert`, {
            method: 'POST',
            headers,
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

export async function listPaymentHistoryAction(businessId: number): Promise<{ success: boolean; data?: BusinessSubscription[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/payments/${businessId}`, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function listAuditLogsAction(businessId: number, limit = 50): Promise<{ success: boolean; data?: SubscriptionAuditLog[]; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/audit-logs/${businessId}?limit=${limit}`, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function listAdminBusinessesAction(params: {
    page?: number;
    pageSize?: number;
    search?: string;
    status?: string;
}): Promise<{ success: boolean; data?: AdminBusinessRow[]; total?: number; page?: number; pageSize?: number; totalPages?: number; error?: string }> {
    try {
        const headers = await buildHeaders();
        const query = new URLSearchParams();
        query.set('page', String(params.page || 1));
        query.set('page_size', String(params.pageSize || 10));
        if (params.search) query.set('search', params.search);
        if (params.status) query.set('status', params.status);
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/admin/businesses?${query.toString()}`, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return {
            success: true,
            data: json.data,
            total: json.total,
            page: json.page,
            pageSize: json.page_size,
            totalPages: json.total_pages,
        };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}

export async function getAdminKPIsAction(): Promise<{ success: boolean; data?: AdminKPIs; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${env.API_BASE_URL}/subscriptions/admin/kpis`, { headers, cache: 'no-store' });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        const json = await res.json();
        return { success: true, data: json.data };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}
