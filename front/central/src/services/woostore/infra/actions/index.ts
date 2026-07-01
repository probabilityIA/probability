'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';

export interface WooStoreState {
    instance_id?: string;
    state?: string;
    public_ip?: string;
    store_url?: string;
    error?: string;
}

async function call(path: string, method: 'GET' | 'POST'): Promise<WooStoreState> {
    const token = (await cookies()).get('session_token')?.value;
    if (!token) {
        return { error: 'No hay sesion activa' };
    }
    try {
        const res = await fetch(`${env.API_BASE_URL}/woo-store${path}`, {
            method,
            headers: { Authorization: `Bearer ${token}` },
            cache: 'no-store',
        });
        const data = await res.json().catch(() => ({}));
        if (!res.ok) {
            return { error: data.error || `Error ${res.status}` };
        }
        return data as WooStoreState;
    } catch (e: any) {
        return { error: e?.message || 'Error de red' };
    }
}

export const getWooStoreStatusAction = () => call('/status', 'GET');
export const startWooStoreAction = () => call('/start', 'POST');
export const stopWooStoreAction = () => call('/stop', 'POST');
