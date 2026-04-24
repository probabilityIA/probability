import { env } from '@/shared/config/env';
import { IProductRepository } from '../../domain/ports';
import {
    Product,
    ProductFamily,
    PaginatedResponse,
    GetProductsParams,
    GetFamiliesParams,
    SingleResponse,
    CreateProductDTO,
    UpdateProductDTO,
    CreateProductFamilyDTO,
    UpdateProductFamilyDTO,
    ActionResponse,
    AddProductIntegrationDTO,
    UpdateProductIntegrationDTO,
    ProductIntegrationsResponse,
    UploadImageResponse
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
        if (this.token) headers['Authorization'] = `Bearer ${this.token}`;
        const res = await fetch(url, { ...options, headers });
        const data = await res.json();
        if (!res.ok) throw new Error(data.message || data.error || 'An error occurred');
        return data;
    }

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

    async uploadProductImage(productId: string, formData: FormData, businessId?: number): Promise<UploadImageResponse> {
        const url = `${this.baseUrl}${this.withBusinessId(`/products/${productId}/image`, businessId)}`;
        const headers: Record<string, string> = { 'Accept': 'application/json' };
        if (this.token) headers['Authorization'] = `Bearer ${this.token}`;
        const res = await fetch(url, { method: 'POST', headers, body: formData });
        const data = await res.json();
        if (!res.ok) throw new Error(data.message || data.error || 'Error uploading image');
        return data;
    }

    async getProductFamilies(params?: GetFamiliesParams): Promise<PaginatedResponse<ProductFamily>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const query = searchParams.toString();
        return this.fetch<PaginatedResponse<ProductFamily>>(`/products/families${query ? `?${query}` : ''}`);
    }

    async getProductFamilyById(familyId: number, businessId?: number): Promise<SingleResponse<ProductFamily>> {
        return this.fetch<SingleResponse<ProductFamily>>(this.withBusinessId(`/products/families/${familyId}`, businessId));
    }

    async getFamilyVariants(familyId: number, businessId?: number): Promise<{ success: boolean; data: Product[] }> {
        return this.fetch<{ success: boolean; data: Product[] }>(this.withBusinessId(`/products/families/${familyId}/variants`, businessId));
    }

    async createProductFamily(data: CreateProductFamilyDTO, businessId?: number): Promise<SingleResponse<ProductFamily>> {
        return this.fetch<SingleResponse<ProductFamily>>(this.withBusinessId('/products/families', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateProductFamily(familyId: number, data: UpdateProductFamilyDTO, businessId?: number): Promise<SingleResponse<ProductFamily>> {
        return this.fetch<SingleResponse<ProductFamily>>(this.withBusinessId(`/products/families/${familyId}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteProductFamily(familyId: number, businessId?: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(this.withBusinessId(`/products/families/${familyId}`, businessId), {
            method: 'DELETE',
        });
    }

    async addProductIntegration(productId: string, data: AddProductIntegrationDTO, businessId?: number): Promise<SingleResponse<any>> {
        return this.fetch<SingleResponse<any>>(this.withBusinessId(`/products/${productId}/integrations`, businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateProductIntegration(productId: string, integrationId: number, data: UpdateProductIntegrationDTO, businessId?: number): Promise<SingleResponse<any>> {
        return this.fetch<SingleResponse<any>>(this.withBusinessId(`/products/${productId}/integrations/${integrationId}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async removeProductIntegration(productId: string, integrationId: number, businessId?: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(this.withBusinessId(`/products/${productId}/integrations/${integrationId}`, businessId), {
            method: 'DELETE',
        });
    }

    async getProductIntegrations(productId: string, businessId?: number): Promise<ProductIntegrationsResponse> {
        return this.fetch<ProductIntegrationsResponse>(this.withBusinessId(`/products/${productId}/integrations`, businessId));
    }

    async lookupProductByExternalRef(
        integrationId: number,
        params: { external_variant_id?: string; external_sku?: string; external_product_id?: string; external_barcode?: string },
        businessId?: number
    ): Promise<SingleResponse<Product>> {
        const searchParams = new URLSearchParams({ integration_id: String(integrationId) });
        Object.entries(params).forEach(([key, value]) => {
            if (value) searchParams.append(key, value);
        });
        if (businessId) searchParams.append('business_id', String(businessId));
        return this.fetch<SingleResponse<Product>>(`/products/lookup-by-external?${searchParams.toString()}`);
    }
}
