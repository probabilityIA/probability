import { env } from '@/shared/config/env';
import { IMarketingLeadsRepository } from '../../domain/ports';
import { GetMarketingLeadsParams, PaginatedLeadsResponse } from '../../domain/types';

export class MarketingLeadsApiRepository implements IMarketingLeadsRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    async getLeads(params?: GetMarketingLeadsParams): Promise<PaginatedLeadsResponse> {
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

        const res = await fetch(`${this.baseUrl}/marketing-leads?${searchParams.toString()}`, { headers });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.error || data.message || 'Error al obtener los leads');
        }

        return data;
    }
}
