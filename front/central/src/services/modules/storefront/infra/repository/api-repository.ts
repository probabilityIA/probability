import { env } from '@/shared/config/env';
import { IStorefrontRepository } from '../../domain/ports';
import {
    StorefrontProduct,
    StorefrontOrder,
    CreateStorefrontOrderDTO,
    RegisterDTO,
    PaginatedResponse,
} from '../../domain/types';

export class StorefrontApiRepository implements IStorefrontRepository {
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
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(url, { ...options, headers });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.message || data.error || 'Error en la solicitud');
        }

        return data;
    }

    async getCatalog(params?: { page?: number; page_size?: number; search?: string; category?: string; business_id?: number }): Promise<PaginatedResponse<StorefrontProduct>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<StorefrontProduct>>(`/storefront/catalog?${searchParams.toString()}`);
    }

    async getProduct(id: string, businessId?: number): Promise<StorefrontProduct> {
        const params = businessId ? `?business_id=${businessId}` : '';
        return this.fetch<StorefrontProduct>(`/storefront/catalog/${id}${params}`);
    }

    async createOrder(data: CreateStorefrontOrderDTO, businessId?: number): Promise<{ message: string }> {
        const params = businessId ? `?business_id=${businessId}` : '';
        return this.fetch<{ message: string }>(`/storefront/orders${params}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async getOrders(params?: { page?: number; page_size?: number; business_id?: number }): Promise<PaginatedResponse<StorefrontOrder>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<StorefrontOrder>>(`/storefront/orders?${searchParams.toString()}`);
    }

    async getOrder(id: string, businessId?: number): Promise<StorefrontOrder> {
        const params = businessId ? `?business_id=${businessId}` : '';
        return this.fetch<StorefrontOrder>(`/storefront/orders/${id}${params}`);
    }

    async register(data: RegisterDTO): Promise<{ message: string }> {
        return this.fetch<{ message: string }>('/storefront/register', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }
}
