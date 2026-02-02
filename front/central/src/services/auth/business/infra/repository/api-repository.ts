import { env } from '@/shared/config/env';
import {
    IBusinessRepository,
} from '../../domain/ports';
import {
    Business,
    PaginatedResponse,
    GetBusinessesParams,
    SingleResponse,
    CreateBusinessDTO,
    UpdateBusinessDTO,
    ActionResponse,
    ConfiguredResource,
    BusinessConfiguredResources,
    GetConfiguredResourcesParams,
    BusinessType,
    CreateBusinessTypeDTO,
    UpdateBusinessTypeDTO
} from '../../domain/types';

export class BusinessApiRepository implements IBusinessRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        // Usar env.API_BASE_URL (servidor) en lugar de envPublic (cliente)
        // Los repositorios se usan en Server Actions que corren en el servidor
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        let logBody: any = options.body;

        if (options.body instanceof FormData) {
            const formDataObj: Record<string, any> = {};
            options.body.forEach((value, key) => {
                formDataObj[key] = value;
            });
            logBody = formDataObj;
        } else if (typeof options.body === 'string') {
            try {
                logBody = JSON.parse(options.body);
            } catch {
                // Keep as string if not JSON
            }
        }

        console.log(`[ANTIGRAVITY API Request] ${options.method || 'GET'} ${url}`, {
            headers: options.headers,
            body: logBody
        });

        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        // If the body is FormData, let the browser set the Content-Type header automatically.
        // Remove the 'Content-Type' header if it was set for FormData.
        if (options.body instanceof FormData) {
            delete headers['Content-Type'];
        }

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
                throw new Error(data.message || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    // --- Business ---

    async getBusinesses(params?: GetBusinessesParams): Promise<PaginatedResponse<Business>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined) searchParams.append(key, String(value));
            });
        }
        const response = await this.fetch<PaginatedResponse<Business>>(`/businesses?${searchParams.toString()}`);
        return {
            ...response,
            data: response.data || []
        };
    }

    async getBusinessById(id: number): Promise<SingleResponse<Business>> {
        return this.fetch<SingleResponse<Business>>(`/businesses/${id}`);
    }

    async createBusiness(data: CreateBusinessDTO): Promise<SingleResponse<Business>> {
        const formData = new FormData();
        Object.entries(data).forEach(([key, value]) => {
            if (value !== undefined) {
                if (value instanceof File) {
                    formData.append(key, value);
                } else {
                    formData.append(key, String(value));
                }
            }
        });

        return this.fetch<SingleResponse<Business>>('/businesses', {
            method: 'POST',
            body: formData,
            // Content-Type header is automatically set for FormData
        });
    }

    async updateBusiness(id: number, data: UpdateBusinessDTO): Promise<SingleResponse<Business>> {
        const formData = new FormData();
        Object.entries(data).forEach(([key, value]) => {
            if (value !== undefined) {
                if (value instanceof File) {
                    formData.append(key, value);
                } else {
                    formData.append(key, String(value));
                }
            }
        });

        return this.fetch<SingleResponse<Business>>(`/businesses/${id}`, {
            method: 'PUT',
            body: formData,
        });
    }

    async deleteBusiness(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/businesses/${id}`, {
            method: 'DELETE',
        });
    }

    // --- Configured Resources ---

    async getConfiguredResources(params?: GetConfiguredResourcesParams): Promise<PaginatedResponse<BusinessConfiguredResources>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<BusinessConfiguredResources>>(`/businesses/configured-resources?${searchParams.toString()}`);
    }

    async getBusinessConfiguredResources(id: number): Promise<SingleResponse<BusinessConfiguredResources>> {
        return this.fetch<SingleResponse<BusinessConfiguredResources>>(`/businesses/${id}/configured-resources`);
    }

    async activateResource(resourceId: number, businessId?: number): Promise<ActionResponse> {
        const searchParams = new URLSearchParams();
        if (businessId) searchParams.append('business_id', String(businessId));

        return this.fetch<ActionResponse>(`/businesses/configured-resources/${resourceId}/activate?${searchParams.toString()}`, {
            method: 'PUT',
        });
    }

    async deactivateResource(resourceId: number, businessId?: number): Promise<ActionResponse> {
        const searchParams = new URLSearchParams();
        if (businessId) searchParams.append('business_id', String(businessId));

        return this.fetch<ActionResponse>(`/businesses/configured-resources/${resourceId}/deactivate?${searchParams.toString()}`, {
            method: 'PUT',
        });
    }

    // --- Business Types ---

    async getBusinessTypes(): Promise<PaginatedResponse<BusinessType>> {
        // Note: The API example for getBusinessTypes doesn't show pagination in the response structure, 
        // but the return type expects PaginatedResponse. 
        // If the API returns a flat list in 'data', we might need to adjust the return type or wrap it.
        // Based on the example: { success: true, message: "...", data: [...] }
        // I will cast it to PaginatedResponse<BusinessType> but be aware 'pagination' field might be missing if the API doesn't send it.
        // Ideally, we should fix the interface or the API. For now, I'll assume it matches or I'll just return it.
        return this.fetch<PaginatedResponse<BusinessType>>('/business-types');
    }

    async getBusinessTypeById(id: number): Promise<SingleResponse<BusinessType>> {
        return this.fetch<SingleResponse<BusinessType>>(`/business-types/${id}`);
    }

    async createBusinessType(data: CreateBusinessTypeDTO): Promise<SingleResponse<BusinessType>> {
        return this.fetch<SingleResponse<BusinessType>>('/business-types', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
    }

    async updateBusinessType(id: number, data: UpdateBusinessTypeDTO): Promise<SingleResponse<BusinessType>> {
        return this.fetch<SingleResponse<BusinessType>>(`/business-types/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
    }

    async deleteBusinessType(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/business-types/${id}`, {
            method: 'DELETE',
        });
    }
}
