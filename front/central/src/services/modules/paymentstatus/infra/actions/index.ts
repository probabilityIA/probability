'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';
import { PaymentStatusInfo } from '../../domain/types';

async function fetchPaymentStatuses(isActive?: boolean): Promise<{ success: boolean; data: PaymentStatusInfo[]; message?: string }> {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    
    const params = new URLSearchParams();
    if (isActive !== undefined) {
        params.append('is_active', String(isActive));
    }
    const url = `${env.API_BASE_URL}/payment-statuses${params.toString() ? `?${params.toString()}` : ''}`;
    
    const response = await fetch(url, {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    });
    
    if (!response.ok) {
        throw new Error(`Failed to fetch payment statuses: ${response.statusText}`);
    }
    
    return response.json();
}

export const getPaymentStatusesAction = async (isActive?: boolean): Promise<{ success: boolean; data: PaymentStatusInfo[]; message?: string }> => {
    try {
        return await fetchPaymentStatuses(isActive);
    } catch (error: any) {
        console.error('Get Payment Statuses Action Error:', error.message);
        throw new Error(error.message);
    }
};
