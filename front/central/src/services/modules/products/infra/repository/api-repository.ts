import { env } from '@/shared/config/env';
import { IProductRepository } from '../../domain/ports';
import {
    Product,
    PaginatedResponse,
    GetProductsParams,
    SingleResponse,
    CreateProductDTO,
    UpdateProductDTO,
    ActionResponse,
    AddProductIntegrationDTO,
    ProductIntegrationsResponse
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
            throw new Error(data.message || data.error || 'An error occurred');
        }

        return data;
    }

    /** Agrega ?business_id=X a la url si se provee (para super admin) */
    private withBusinessId(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
    }

    async getProducts(params?: GetProductsParams): Promise<PaginatedResponse<Product>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const query = searchParams.toString();
        return this.fetch<PaginatedResponse<Product>>(`/products${query ? `?${query}` : ''}`);
    }

    async getProductById(id: string, businessId?: number): Promise<SingleResponse<Product>> {
        return this.fetch<SingleResponse<Product>>(this.withBusinessId(`/products/${id}`, businessId));
    }

    async createProduct(data: CreateProductDTO, businessId?: number): Promise<SingleResponse<Product>> {
        return this.fetch<SingleResponse<Product>>(this.withBusinessId('/products', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateProduct(id: string, data: UpdateProductDTO, businessId?: number): Promise<SingleResponse<Product>> {
        return this.fetch<SingleResponse<Product>>(this.withBusinessId(`/products/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteProduct(id: string, businessId?: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(this.withBusinessId(`/products/${id}`, businessId), {
            method: 'DELETE',
        });
    }

    // ═══════════════════════════════════════════
    // Product-Integration Management
    // ═══════════════════════════════════════════

    async addProductIntegration(
        productId: string,
        data: AddProductIntegrationDTO,
        businessId?: number
    ): Promise<SingleResponse<any>> {
        return this.fetch<SingleResponse<any>>(this.withBusinessId(`/products/${productId}/integrations`, businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async removeProductIntegration(
        productId: string,
        integrationId: number,
        businessId?: number
    ): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(this.withBusinessId(`/products/${productId}/integrations/${integrationId}`, businessId), {
            method: 'DELETE',
        });
    }

    async getProductIntegrations(productId: string, businessId?: number): Promise<ProductIntegrationsResponse> {
        return this.fetch<ProductIntegrationsResponse>(this.withBusinessId(`/products/${productId}/integrations`, businessId));
    }
}
