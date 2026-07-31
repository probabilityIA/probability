import { env } from '@/shared/config/env';
import { IPublicSiteRepository } from '../../domain/ports';
import { PublicBusiness, PublicProduct, PaginatedResponse, ContactFormDTO, CreateCheckoutInput, CheckoutSession } from '../../domain/types';

export class PublicSiteApiRepository implements IPublicSiteRepository {
    private baseUrl: string;

    constructor() {
        this.baseUrl = env.API_BASE_URL;
    }

    private async fetch<T>(path: string, options: RequestInit = {}, cacheable = false): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        const cachePolicy: RequestInit = cacheable
            ? { next: { revalidate: 60, tags: ['public-tienda'] } }
            : { cache: 'no-store' };

        const res = await fetch(url, { ...options, ...cachePolicy, headers });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.message || data.error || 'Error en la solicitud');
        }

        return data;
    }

    async getBusinessPage(slug: string): Promise<PublicBusiness> {
        return this.fetch<PublicBusiness>(`/public/tienda/${slug}`, {}, true);
    }

    async getCatalog(slug: string, params?: { page?: number; page_size?: number; search?: string; category?: string }): Promise<PaginatedResponse<PublicProduct>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<PublicProduct>>(`/public/tienda/${slug}/catalog?${searchParams.toString()}`, {}, true);
    }

    async getProduct(slug: string, productId: string): Promise<PublicProduct> {
        return this.fetch<PublicProduct>(`/public/tienda/${slug}/product/${productId}`, {}, true);
    }

    async submitContact(slug: string, data: ContactFormDTO): Promise<{ message: string }> {
        return this.fetch<{ message: string }>(`/public/tienda/${slug}/contact`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async createCheckoutSession(slug: string, data: CreateCheckoutInput): Promise<CheckoutSession> {
        const res = await this.fetch<{ data: CheckoutSession }>(`/public/tienda/${slug}/checkout/bold/signature`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async getCheckoutStatus(slug: string, reference: string): Promise<{ status: string }> {
        const res = await this.fetch<{ data: { status: string } }>(`/public/tienda/${slug}/checkout/${reference}/status`);
        return res.data;
    }
}
