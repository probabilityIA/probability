'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';
import { WooStoreState } from '@/services/woostore/domain/types';

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

export async function getWooStoreStatusAction(): Promise<WooStoreState> {
    return call('/status', 'GET');
}

export async function startWooStoreAction(): Promise<WooStoreState> {
    return call('/start', 'POST');
}

export async function stopWooStoreAction(): Promise<WooStoreState> {
    return call('/stop', 'POST');
}
