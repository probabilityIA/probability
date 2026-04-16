import { env } from '@/shared/config/env';
import { IPublicSiteRepository } from '../../domain/ports';
import { PublicBusiness, PublicProduct, PaginatedResponse, ContactFormDTO } from '../../domain/types';

export class PublicSiteApiRepository implements IPublicSiteRepository {
    private baseUrl: string;

    constructor() {
        this.baseUrl = env.API_BASE_URL;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        const res = await fetch(url, { ...options, headers, cache: 'no-store' });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.message || data.error || 'Error en la solicitud');
        }

        return data;
    }

    async getBusinessPage(slug: string): Promise<PublicBusiness> {
        return this.fetch<PublicBusiness>(`/public/tienda/${slug}`);
    }

    async getCatalog(slug: string, params?: { page?: number; page_size?: number; search?: string; category?: string }): Promise<PaginatedResponse<PublicProduct>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<PublicProduct>>(`/public/tienda/${slug}/catalog?${searchParams.toString()}`);
    }

    async getProduct(slug: string, productId: string): Promise<PublicProduct> {
        return this.fetch<PublicProduct>(`/public/tienda/${slug}/product/${productId}`);
    }

    async submitContact(slug: string, data: ContactFormDTO): Promise<{ message: string }> {
        return this.fetch<{ message: string }>(`/public/tienda/${slug}/contact`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }
}
