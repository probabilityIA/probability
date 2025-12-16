import { env } from '@/shared/config/env';
import { IIntegrationRepository } from '../../domain/ports';
import {
    Integration,
    PaginatedResponse,
    GetIntegrationsParams,
    SingleResponse,
    CreateIntegrationDTO,
    UpdateIntegrationDTO,
    ActionResponse,
    IntegrationType,
    CreateIntegrationTypeDTO,
    UpdateIntegrationTypeDTO
} from '../../domain/types';

export class IntegrationApiRepository implements IIntegrationRepository {
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

        // Si el body es FormData, no establecer Content-Type (el navegador lo hará automáticamente)
        const isFormData = options.body instanceof FormData;
        
        const headers: Record<string, string> = {
            'Accept': 'application/json',
            ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
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
                throw new Error(data.error || data.message || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    async getIntegrations(params?: GetIntegrationsParams): Promise<PaginatedResponse<Integration>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const response = await this.fetch<PaginatedResponse<Integration>>(`/integrations?${searchParams.toString()}`);
        return {
            ...response,
            data: response.data || []
        };
    }

    async getIntegrationById(id: number): Promise<SingleResponse<Integration>> {
        return this.fetch<SingleResponse<Integration>>(`/integrations/${id}`);
    }

    async getIntegrationByType(type: string, businessId?: number): Promise<SingleResponse<Integration>> {
        const searchParams = new URLSearchParams();
        if (businessId) searchParams.append('business_id', String(businessId));
        return this.fetch<SingleResponse<Integration>>(`/integrations/type/${type}?${searchParams.toString()}`);
    }

    async createIntegration(data: CreateIntegrationDTO): Promise<SingleResponse<Integration>> {
        return this.fetch<SingleResponse<Integration>>('/integrations', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateIntegration(id: number, data: UpdateIntegrationDTO): Promise<SingleResponse<Integration>> {
        return this.fetch<SingleResponse<Integration>>(`/integrations/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteIntegration(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/integrations/${id}`, {
            method: 'DELETE',
        });
    }

    async testConnection(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/integrations/${id}/test`, {
            method: 'POST',
        });
    }

    async activateIntegration(id: number): Promise<SingleResponse<Integration>> {
        return this.fetch<SingleResponse<Integration>>(`/integrations/${id}/activate`, {
            method: 'PUT',
        });
    }

    async deactivateIntegration(id: number): Promise<SingleResponse<Integration>> {
        return this.fetch<SingleResponse<Integration>>(`/integrations/${id}/deactivate`, {
            method: 'PUT',
        });
    }

    async setAsDefault(id: number): Promise<SingleResponse<Integration>> {
        return this.fetch<SingleResponse<Integration>>(`/integrations/${id}/set-default`, {
            method: 'PUT',
        });
    }

    async syncOrders(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/shopify/sync/${id}`, {
            method: 'POST',
        });
    }

    async testIntegration(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/integrations/${id}/test`, {
            method: 'POST',
        });
    }

    async testConnectionRaw(typeCode: string, config: Record<string, any>, credentials: Record<string, any>): Promise<ActionResponse> {
        return this.fetch<ActionResponse>('/integrations/test', {
            method: 'POST',
            body: JSON.stringify({
                type_code: typeCode,
                config,
                credentials
            })
        });
    }

    // Integration Types
    async getIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>> {
        return this.fetch<SingleResponse<IntegrationType[]>>('/integration-types');
    }

    async getActiveIntegrationTypes(): Promise<SingleResponse<IntegrationType[]>> {
        return this.fetch<SingleResponse<IntegrationType[]>>('/integration-types/active');
    }

    async getIntegrationTypeById(id: number): Promise<SingleResponse<IntegrationType>> {
        return this.fetch<SingleResponse<IntegrationType>>(`/integration-types/${id}`);
    }

    async getIntegrationTypeByCode(code: string): Promise<SingleResponse<IntegrationType>> {
        return this.fetch<SingleResponse<IntegrationType>>(`/integration-types/code/${code}`);
    }

    async createIntegrationType(data: CreateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>> {
        // Si hay imagen, usar FormData, sino JSON
        if (data.image_file) {
            const formData = new FormData();
            formData.append('name', data.name);
            if (data.code) formData.append('code', data.code);
            if (data.description) formData.append('description', data.description);
            if (data.icon) formData.append('icon', data.icon);
            formData.append('category', data.category);
            if (data.is_active !== undefined) formData.append('is_active', String(data.is_active));
            if (data.config_schema) formData.append('credentials_schema', JSON.stringify(data.config_schema));
            if (data.credentials_schema) formData.append('credentials_schema', JSON.stringify(data.credentials_schema));
            if (data.setup_instructions) formData.append('setup_instructions', data.setup_instructions);
            formData.append('image_file', data.image_file);

            return this.fetch<SingleResponse<IntegrationType>>('/integration-types', {
                method: 'POST',
                body: formData,
                headers: {} // No establecer Content-Type, el navegador lo hará automáticamente con FormData
            });
        }

        return this.fetch<SingleResponse<IntegrationType>>('/integration-types', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateIntegrationType(id: number, data: UpdateIntegrationTypeDTO): Promise<SingleResponse<IntegrationType>> {
        // Si hay imagen o remove_image, usar FormData, sino JSON
        if (data.image_file || data.remove_image !== undefined) {
            const formData = new FormData();
            if (data.name) formData.append('name', data.name);
            if (data.code) formData.append('code', data.code);
            if (data.description) formData.append('description', data.description);
            if (data.icon) formData.append('icon', data.icon);
            if (data.category) formData.append('category', data.category);
            if (data.is_active !== undefined) formData.append('is_active', String(data.is_active));
            if (data.config_schema) formData.append('credentials_schema', JSON.stringify(data.config_schema));
            if (data.credentials_schema) formData.append('credentials_schema', JSON.stringify(data.credentials_schema));
            if (data.setup_instructions) formData.append('setup_instructions', data.setup_instructions);
            if (data.image_file) formData.append('image_file', data.image_file);
            if (data.remove_image !== undefined) formData.append('remove_image', String(data.remove_image));

            return this.fetch<SingleResponse<IntegrationType>>(`/integration-types/${id}`, {
                method: 'PUT',
                body: formData,
                headers: {} // No establecer Content-Type, el navegador lo hará automáticamente con FormData
            });
        }

        return this.fetch<SingleResponse<IntegrationType>>(`/integration-types/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteIntegrationType(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/integration-types/${id}`, {
            method: 'DELETE',
        });
    }
}
