import { env } from '@/shared/config/env';
import type { Access } from '../../domain/types';

export class AccessApiRepository {
    constructor(private readonly token: string) {}

    async getMyAccess(businessId?: number): Promise<Access> {
        const query = businessId ? `?business_id=${businessId}` : '';
        const res = await fetch(`${env.API_BASE_URL}/auth/me/access${query}`, {
            method: 'GET',
            headers: { Authorization: `Bearer ${this.token}` },
            cache: 'no-store',
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok || !body?.success || !body?.data) {
            throw new Error(body?.error || body?.message || `No se pudo obtener el acceso (${res.status})`);
        }
        return body.data as Access;
    }
}
