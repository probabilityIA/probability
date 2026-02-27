'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';

// Define types locally if not yet available in shared types, or import if they exist.
// Based on the existing page.tsx:
export interface Wallet {
    ID: string;
    BusinessID: number;
    Balance: number;
}

export interface WalletTransactionRequest {
    ID: string;
    WalletID: string;
    Amount: number;
    CreatedAt: string;
    // Add other fields as necessary from the backend response
}

/**
 * Helper to get the auth header from cookies
 */
async function getAuthHeader() {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value;
    return {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
    };
}

/**
 * Fetch all wallets (Admin only)
 */
export async function getWalletsAction() {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/all`, {
            headers,
            cache: 'no-store'
        });

        if (!res.ok) {
            throw new Error(`Failed to fetch wallets: ${res.status} ${res.statusText}`);
        }

        const data = await res.json();
        return { success: true, data: data as Wallet[] };
    } catch (error: any) {
        console.error('getWalletsAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Fetch pending recharge requests (Admin only)
 */
export async function getPendingRequestsAction() {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/admin/pending-requests`, {
            headers,
            cache: 'no-store'
        });

        if (!res.ok) {
            if (res.status === 404) return { success: true, data: [] };
            throw new Error(`Failed to fetch pending requests: ${res.status}`);
        }

        const data = await res.json();
        return { success: true, data: data || [] };
    } catch (error: any) {
        console.error('getPendingRequestsAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Fetch processed (approved/rejected) recharge requests (Admin only)
 */
export async function getProcessedRequestsAction() {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/admin/processed-requests`, {
            headers,
            cache: 'no-store'
        });

        if (!res.ok) {
            if (res.status === 404) return { success: true, data: [] };
            throw new Error(`Failed to fetch processed requests: ${res.status}`);
        }

        const data = await res.json();
        return { success: true, data: data || [] };
    } catch (error: any) {
        console.error('getProcessedRequestsAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Process a recharge request (Admin only)
 *
 */
export async function processRequestAction(id: string, action: 'approve' | 'reject') {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/admin/requests/${id}/${action}`, {
            method: 'POST',
            headers,
        });

        if (!res.ok) {
            throw new Error(`Failed to process request: ${res.status}`);
        }

        return { success: true };
    } catch (error: any) {
        console.error('processRequestAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Get current business wallet balance
 * @param businessId Optional - super admin can specify a business_id to view another business's balance
 */
export async function getWalletBalanceAction(businessId?: number) {
    try {
        const headers = await getAuthHeader();
        const url = businessId
            ? `${env.API_BASE_URL}/pay/wallet/balance?business_id=${businessId}`
            : `${env.API_BASE_URL}/pay/wallet/balance`;
        const res = await fetch(url, {
            headers,
            cache: 'no-store'
        });

        if (!res.ok) {
            throw new Error(`Failed to fetch balance: ${res.status}`);
        }

        const data = await res.json();
        return { success: true, data: data as Wallet };
    } catch (error: any) {
        console.error('getWalletBalanceAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Request a recharge
 * @param businessId Optional - super admin can recharge on behalf of a specific business
 */
export async function rechargeWalletAction(amount: number, businessId?: number) {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/recharge`, {
            method: 'POST',
            headers,
            body: JSON.stringify({ amount, ...(businessId ? { business_id: businessId } : {}) })
        });

        if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || `Failed to recharge: ${res.status}`);
        }

        const data = await res.json();
        return { success: true, data };
    } catch (error: any) {
        console.error('rechargeWalletAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Manual debit from a business wallet (Admin only)
 */
export async function manualDebitAction(businessId: number, amount: number, reference: string) {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/admin/manual-debit`, {
            method: 'POST',
            headers,
            body: JSON.stringify({ business_id: businessId, amount, reference })
        });

        if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || `Failed to debit wallet: ${res.status}`);
        }

        return { success: true };
    } catch (error: any) {
        console.error('manualDebitAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Get business transaction history
 * @param businessId Optional - super admin can view history for a specific business
 */
export async function getWalletHistoryAction(businessId?: number) {
    try {
        const headers = await getAuthHeader();
        const url = businessId
            ? `${env.API_BASE_URL}/pay/wallet/history?business_id=${businessId}`
            : `${env.API_BASE_URL}/pay/wallet/history`;
        const res = await fetch(url, {
            headers,
            cache: 'no-store'
        });

        if (!res.ok) {
            if (res.status === 404) return { success: true, data: [] };
            throw new Error(`Failed to fetch wallet history: ${res.status}`);
        }

        const data = await res.json();
        return { success: true, data: data || [] };
    } catch (error: any) {
        console.error('getWalletHistoryAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Clear recharge history for a business (Admin only)
 */
export async function clearRechargeHistoryAction(businessId: number) {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/admin/history/${businessId}`, {
            method: 'DELETE',
            headers,
        });

        if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || `Failed to clear history: ${res.status}`);
        }

        return { success: true };
    } catch (error: any) {
        console.error('clearRechargeHistoryAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Debit from wallet for guide generation
 */
export async function debitForGuideAction(amount: number, trackingNumber: string) {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/pay/wallet/debit-guide`, {
            method: 'POST',
            headers,
            body: JSON.stringify({ amount, tracking_number: trackingNumber })
        });

        if (!res.ok) {
            const errData = await res.json().catch(() => ({}));
            throw new Error(errData.error || `Failed to debit wallet: ${res.status}`);
        }

        const data = await res.json();
        return { success: true, data };
    } catch (error: any) {
        console.error('debitForGuideAction error:', error);
        return { success: false, error: error.message };
    }
}
