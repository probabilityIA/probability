import { env } from '@/shared/config/env';
import { IPermissionRepository } from '../../domain/ports';
import {
    Permission,
    PaginatedResponse,
    GetPermissionsParams,
    SingleResponse,
    CreatePermissionDTO,
    UpdatePermissionDTO,
    ActionResponse
} from '../../domain/types';

export class PermissionApiRepository implements IPermissionRepository {
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
                throw new Error(data.message || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    async getPermissions(params?: GetPermissionsParams): Promise<PaginatedResponse<Permission>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const response = await this.fetch<PaginatedResponse<Permission>>(`/permissions?${searchParams.toString()}`);
        return {
            ...response,
            data: response.data || []
        };
    }

    async getPermissionById(id: number): Promise<SingleResponse<Permission>> {
        return this.fetch<SingleResponse<Permission>>(`/permissions/${id}`);
    }

    async getPermissionsByScope(scopeId: number): Promise<PaginatedResponse<Permission>> {
        return this.fetch<PaginatedResponse<Permission>>(`/permissions/scope/${scopeId}`);
    }

    async getPermissionsByResource(resource: string): Promise<PaginatedResponse<Permission>> {
        return this.fetch<PaginatedResponse<Permission>>(`/permissions/resource/${resource}`);
    }

    async createPermission(data: CreatePermissionDTO): Promise<SingleResponse<Permission>> {
        return this.fetch<SingleResponse<Permission>>('/permissions', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updatePermission(id: number, data: UpdatePermissionDTO): Promise<SingleResponse<Permission>> {
        return this.fetch<SingleResponse<Permission>>(`/permissions/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deletePermission(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/permissions/${id}`, {
            method: 'DELETE',
        });
    }
}
