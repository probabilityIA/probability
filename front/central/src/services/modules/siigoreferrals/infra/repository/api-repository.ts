import { env } from '@/shared/config/env';
import { ISiigoReferralsRepository } from '../../domain/ports';
import { GetSiigoReferralsParams, PaginatedSiigoReferralsResponse } from '../../domain/types';

export class SiigoReferralsApiRepository implements ISiigoReferralsRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    async getReferrals(params?: GetSiigoReferralsParams): Promise<PaginatedSiigoReferralsResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }

        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
        };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(`${this.baseUrl}/siigo-referrals?${searchParams.toString()}`, { headers });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.error || data.message || 'Error al obtener los referidos');
        }

        return data;
    }
}
