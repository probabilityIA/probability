'use server';

import { getAuthToken } from '@/shared/utils/server-auth';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:3050';

async function buildHeaders(): Promise<Record<string, string>> {
    const token = await getAuthToken();
    const headers: Record<string, string> = {};
    if (token) headers['Authorization'] = `Bearer ${token}`;
    return headers;
}

export interface BusinessSubscription {
    id?: number;
    businessId: number;
    amount: number;
    startDate?: string;
    endDate?: string;
    status: string; // 'paid' | 'pending' | 'rejected'
    paymentReference?: string;
    notes?: string;
    createdAt?: string;
}

export interface BusinessSubscriptionStatus {
    subscriptionStatus: string; // 'active' | 'expired' | 'cancelled'
    subscriptionEndDate?: string;
    businessName?: string;
}

/** Obtiene el estado de suscripción actual del negocio autenticado.
 *  Super admins pueden pasar businessId para ver el estado de otro negocio. */
export async function getMySubscriptionAction(businessId?: number): Promise<{ success: boolean; data?: BusinessSubscription; status?: BusinessSubscriptionStatus; error?: string }> {
    try {
        const headers = await buildHeaders();
        const url = businessId
            ? `${API_BASE}/api/v1/subscriptions/me?businessId=${businessId}`
            : `${API_BASE}/api/v1/subscriptions/me`;
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

/** Super Admin: Registra un pago de suscripción y activa al negocio */
export async function registerSubscriptionPaymentAction(payload: {
    businessId: number;
    amount: number;
    monthsToAdd: number;
    paymentReference?: string;
    notes?: string;
}): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${API_BASE}/api/v1/subscriptions/register-payment`, {
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

/** Super Admin: Deshabilita manualmente la suscripción de un negocio */
export async function disableSubscriptionAction(businessId: number): Promise<{ success: boolean; error?: string }> {
    try {
        const headers = await buildHeaders();
        const res = await fetch(`${API_BASE}/api/v1/subscriptions/disable?businessId=${businessId}`, {
            method: 'POST',
            headers,
        });
        if (!res.ok) throw new Error(`Error ${res.status}`);
        return { success: true };
    } catch (err: any) {
        return { success: false, error: err.message };
    }
}
