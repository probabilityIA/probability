'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';

export interface SavedQuoteRate {
    carrier?: string;
    product?: string;
    flete?: number;
    deliveryDays?: number;
    [key: string]: any;
}

export interface SavedQuote {
    id: number;
    business_id: number;
    integration_id: number;
    source: string;
    order_uuid?: string | null;
    external_order_ref?: string;
    rates: SavedQuoteRate[];
    selected_carrier?: string;
    selected_service_code?: string;
    status: string;
    expires_at?: string | null;
    created_at: string;
}

export interface SavedQuotesResponse {
    data: SavedQuote[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export async function getSavedQuotesAction(params: {
    businessId?: number | null;
    page?: number;
    pageSize?: number;
    source?: string;
    status?: string;
}): Promise<SavedQuotesResponse> {
    const empty: SavedQuotesResponse = { data: [], total: 0, page: 1, page_size: 10, total_pages: 0 };
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value || '';

        const qs = new URLSearchParams();
        if (params.businessId) qs.set('business_id', String(params.businessId));
        qs.set('page', String(params.page || 1));
        qs.set('page_size', String(params.pageSize || 10));
        if (params.source) qs.set('source', params.source);
        if (params.status) qs.set('status', params.status);

        const res = await fetch(`${env.API_BASE_URL}/shipments/quotes?${qs.toString()}`, {
            headers: {
                'Authorization': `Bearer ${token}`,
                'Content-Type': 'application/json',
            },
            cache: 'no-store',
        });

        if (!res.ok) return empty;
        const data = await res.json();
        return {
            data: data.data || [],
            total: data.total || 0,
            page: data.page || 1,
            page_size: data.page_size || 10,
            total_pages: data.total_pages || 0,
        };
    } catch (e) {
        console.error('getSavedQuotesAction error:', e);
        return empty;
    }
}
