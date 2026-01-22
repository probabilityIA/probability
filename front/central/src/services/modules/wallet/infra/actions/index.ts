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
        // env.API_BASE_URL already includes /api/v1
        const res = await fetch(`${env.API_BASE_URL}/wallet/all`, {
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
        const res = await fetch(`${env.API_BASE_URL}/wallet/admin/pending-requests`, {
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
        const res = await fetch(`${env.API_BASE_URL}/wallet/admin/processed-requests`, {
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
        const res = await fetch(`${env.API_BASE_URL}/wallet/admin/requests/${id}/${action}`, {
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
 */
export async function getWalletBalanceAction() {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/wallet/balance`, {
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
 */
export async function rechargeWalletAction(amount: number) {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/wallet/recharge`, {
            method: 'POST',
            headers,
            body: JSON.stringify({ amount })
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
 * Report that a payment has been made (Business only)
 */
export async function reportPaymentAction(requestId: string) {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/wallet/report-payment/${requestId}`, {
            method: 'POST',
            headers,
        });

        if (!res.ok) {
            throw new Error(`Failed to report payment: ${res.status}`);
        }

        return { success: true };
    } catch (error: any) {
        console.error('reportPaymentAction error:', error);
        return { success: false, error: error.message };
    }
}

/**
 * Manual debit from a business wallet (Admin only)
 */
export async function manualDebitAction(businessId: number, amount: number, reference: string) {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/wallet/admin/manual-debit`, {
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
 * Get current business transaction history
 */
export async function getWalletHistoryAction() {
    try {
        const headers = await getAuthHeader();
        const res = await fetch(`${env.API_BASE_URL}/wallet/history`, {
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
