import { env } from '@/shared/config/env';
import { IProductRepository } from '../../domain/ports';
import {
    Product,
    PaginatedResponse,
    GetProductsParams,
    SingleResponse,
    CreateProductDTO,
    UpdateProductDTO,
    ActionResponse
} from '../../domain/types';

export class ProductApiRepository implements IProductRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        console.log(`[API Request] ${options.method || 'GET'} ${url}`, {
            headers: options.headers,
            body: options.body
        });

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

            console.log(`[API Response] ${res.status} ${url}`, data);

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

    async getProducts(params?: GetProductsParams): Promise<PaginatedResponse<Product>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<Product>>(`/products?${searchParams.toString()}`);
    }

    async getProductById(id: string): Promise<SingleResponse<Product>> {
        return this.fetch<SingleResponse<Product>>(`/products/${id}`);
    }

    async createProduct(data: CreateProductDTO): Promise<SingleResponse<Product>> {
        return this.fetch<SingleResponse<Product>>('/products', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateProduct(id: string, data: UpdateProductDTO): Promise<SingleResponse<Product>> {
        return this.fetch<SingleResponse<Product>>(`/products/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteProduct(id: string): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/products/${id}`, {
            method: 'DELETE',
        });
    }
}
