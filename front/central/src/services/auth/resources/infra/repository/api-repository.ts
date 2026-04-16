import { env } from '@/shared/config/env';
import { IResourceRepository } from '../../domain/ports';
import {
    Resource,
    PaginatedResponse,
    GetResourcesParams,
    SingleResponse,
    CreateResourceDTO,
    UpdateResourceDTO,
    ActionResponse
} from '../../domain/types';

export class ResourceApiRepository implements IResourceRepository {
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

    async getResources(params?: GetResourcesParams): Promise<PaginatedResponse<Resource>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const response = await this.fetch<PaginatedResponse<Resource>>(`/resources?${searchParams.toString()}`);
        return {
            ...response,
            data: {
                ...response.data,
                resources: response.data.resources || []
            }
        };
    }

    async getResourceById(id: number): Promise<SingleResponse<Resource>> {
        return this.fetch<SingleResponse<Resource>>(`/resources/${id}`);
    }

    async createResource(data: CreateResourceDTO): Promise<SingleResponse<Resource>> {
        return this.fetch<SingleResponse<Resource>>('/resources', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateResource(id: number, data: UpdateResourceDTO): Promise<SingleResponse<Resource>> {
        return this.fetch<SingleResponse<Resource>>(`/resources/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteResource(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/resources/${id}`, {
            method: 'DELETE',
        });
    }
}
