import { env } from '@/shared/config/env';
import { IPublicSiteRepository } from '../../domain/ports';
import { PublicBusiness, PublicProduct, PaginatedResponse, ContactFormDTO, CreateCheckoutInput, CheckoutSession, TiendaSession, TiendaRegisterDTO, TiendaOrderSummary, TiendaAddress, TiendaAddressInput, TiendaProfileUpdate } from '../../domain/types';

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

    async getCategories(slug: string): Promise<string[]> {
        const res = await this.fetch<{ data: string[] }>(`/public/tienda/${slug}/categories`, {}, true);
        return res.data || [];
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

    async createCheckoutSession(slug: string, data: CreateCheckoutInput, token?: string): Promise<CheckoutSession> {
        const res = await this.fetch<{ data: CheckoutSession }>(`/public/tienda/${slug}/checkout/bold/signature`, {
            method: 'POST',
            body: JSON.stringify(data),
            headers: token ? { Authorization: `Bearer ${token}` } : undefined,
        });
        return res.data;
    }

    async getCheckoutStatus(slug: string, reference: string): Promise<{ status: string }> {
        const res = await this.fetch<{ data: { status: string } }>(`/public/tienda/${slug}/checkout/${reference}/status`);
        return res.data;
    }

    async registerCustomer(slug: string, data: TiendaRegisterDTO): Promise<{ success: boolean; message?: string }> {
        return this.fetch<{ success: boolean; message?: string }>(`/storefront/register`, {
            method: 'POST',
            body: JSON.stringify({ ...data, business_code: slug }),
        });
    }

    async getSession(slug: string, token: string): Promise<TiendaSession> {
        const res = await this.fetch<{ data: TiendaSession }>(`/public/tienda/${slug}/session`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        return res.data;
    }

    async getMyOrders(slug: string, token: string, page = 1): Promise<{ data: TiendaOrderSummary[]; total: number }> {
        return this.fetch<{ data: TiendaOrderSummary[]; total: number }>(`/public/tienda/${slug}/my-orders?page=${page}&page_size=10`, {
            headers: { Authorization: `Bearer ${token}` },
        });
    }

    async updateProfile(slug: string, token: string, data: TiendaProfileUpdate): Promise<{ success: boolean }> {
        return this.fetch<{ success: boolean }>(`/public/tienda/${slug}/profile`, {
            method: 'PUT',
            body: JSON.stringify(data),
            headers: { Authorization: `Bearer ${token}` },
        });
    }

    async uploadAvatar(slug: string, token: string, formData: FormData): Promise<{ success: boolean; avatar_url: string }> {
        const url = `${this.baseUrl}/public/tienda/${slug}/profile/avatar`;
        const res = await fetch(url, {
            method: 'POST',
            body: formData,
            headers: { Authorization: `Bearer ${token}` },
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error || 'Error al subir imagen');
        return data;
    }

    async listAddresses(slug: string, token: string): Promise<TiendaAddress[]> {
        const res = await this.fetch<{ data: TiendaAddress[] }>(`/public/tienda/${slug}/addresses`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        return res.data || [];
    }

    async addAddress(slug: string, token: string, data: TiendaAddressInput): Promise<{ success: boolean }> {
        return this.fetch<{ success: boolean }>(`/public/tienda/${slug}/addresses`, {
            method: 'POST',
            body: JSON.stringify(data),
            headers: { Authorization: `Bearer ${token}` },
        });
    }

    async setPrimaryAddress(slug: string, token: string, id: number): Promise<{ success: boolean }> {
        return this.fetch<{ success: boolean }>(`/public/tienda/${slug}/addresses/${id}/primary`, {
            method: 'PUT',
            headers: { Authorization: `Bearer ${token}` },
        });
    }

    async deleteAddress(slug: string, token: string, id: number): Promise<{ success: boolean }> {
        return this.fetch<{ success: boolean }>(`/public/tienda/${slug}/addresses/${id}`, {
            method: 'DELETE',
            headers: { Authorization: `Bearer ${token}` },
        });
    }

    async getMyPrices(slug: string, token: string): Promise<Record<string, number>> {
        const res = await this.fetch<{ data: { prices: Record<string, number> } }>(`/public/tienda/${slug}/my-prices`, {
            headers: { Authorization: `Bearer ${token}` },
        });
        return res.data.prices || {};
    }
}
