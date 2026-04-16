'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';
import { FulfillmentStatusInfo } from '../../domain/types';

async function fetchFulfillmentStatuses(isActive?: boolean): Promise<{ success: boolean; data: FulfillmentStatusInfo[]; message?: string }> {
    const cookieStore = await cookies();
    const token = cookieStore.get('session_token')?.value || null;
    
    const params = new URLSearchParams();
    if (isActive !== undefined) {
        params.append('is_active', String(isActive));
    }
    const url = `${env.API_BASE_URL}/fulfillment-statuses${params.toString() ? `?${params.toString()}` : ''}`;
    
    const response = await fetch(url, {
        headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
        },
    });
    
    if (!response.ok) {
        throw new Error(`Failed to fetch fulfillment statuses: ${response.statusText}`);
    }
    
    return response.json();
}

export const getFulfillmentStatusesAction = async (isActive?: boolean): Promise<{ success: boolean; data: FulfillmentStatusInfo[]; message?: string }> => {
    try {
        return await fetchFulfillmentStatuses(isActive);
    } catch (error: any) {
        console.error('Get Fulfillment Statuses Action Error:', error.message);
        throw new Error(error.message);
    }
};
