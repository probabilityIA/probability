import { env } from '@/shared/config/env';
import type { ITourRepository } from '../../domain/ports';
import type { SaveTourProgressInput, SkipAllTour, TourProgress } from '../../domain/types';

export class TourApiRepository implements ITourRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private buildPath(path: string, businessId?: number): string {
        if (!businessId) return path;
        const separator = path.includes('?') ? '&' : '?';
        return `${path}${separator}business_id=${businessId}`;
    }

    private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
        const headers: Record<string, string> = {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...((options.headers as Record<string, string>) || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(`${this.baseUrl}${path}`, { ...options, headers, cache: 'no-store' });
        const data = await res.json();

        if (!res.ok || data?.success === false) {
            throw new Error(data?.message || data?.error || 'Error en el servicio de tours');
        }

        return data.data as T;
    }

    async listProgress(businessId?: number): Promise<TourProgress[]> {
        const data = await this.request<TourProgress[] | null>(this.buildPath('/tours/progress', businessId));
        return data ?? [];
    }

    async saveProgress(input: SaveTourProgressInput, businessId?: number): Promise<TourProgress> {
        return this.request<TourProgress>(this.buildPath('/tours/progress', businessId), {
            method: 'PUT',
            body: JSON.stringify(input),
        });
    }

    async resetTour(tourKey: string, businessId?: number): Promise<void> {
        await this.request<null>(this.buildPath(`/tours/progress/${encodeURIComponent(tourKey)}`, businessId), {
            method: 'DELETE',
        });
    }

    async resetAll(businessId?: number): Promise<void> {
        await this.request<null>(this.buildPath('/tours/progress/reset', businessId), { method: 'POST' });
    }

    async skipAll(tours: SkipAllTour[], businessId?: number): Promise<void> {
        await this.request<null>(this.buildPath('/tours/progress/skip-all', businessId), {
            method: 'POST',
            body: JSON.stringify({ tours }),
        });
    }
}
