import { env } from '@/shared/config/env';
import { IDashboardRepository } from '../../domain/ports';
import { DashboardStatsResponse } from '../../domain/types';

export class DashboardApiRepository implements IDashboardRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        if (this.token) {
            (headers as any)['Authorization'] = `Bearer ${this.token}`;
        }

        try {
            const res = await fetch(url, {
                ...options,
                headers,
            });

            const data = await res.json();

            if (!res.ok) {
                console.error(`[API Error] ${res.status} ${url}`, data);
                throw new Error(data.message || data.error || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    async getStats(businessId?: number): Promise<DashboardStatsResponse> {
        const params = businessId ? `?business_id=${businessId}` : '';
        return this.fetch<DashboardStatsResponse>(`/dashboard/stats${params}`);
    }
}
