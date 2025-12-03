import { env } from '@/shared/config/env';
import { IRoleRepository } from '../../domain/ports';
import {
    Role,
    PaginatedResponse,
    GetRolesParams,
    SingleResponse,
    CreateRoleDTO,
    UpdateRoleDTO,
    ActionResponse,
    RolePermissionsResponse,
    AssignPermissionsDTO,
    AssignPermissionsResponse
} from '../../domain/types';

export class RoleApiRepository implements IRoleRepository {
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

            console.log(`[API Response] ${res.status} ${url} `, data);

            if (!res.ok) {
                console.error(`[API Error] ${res.status} ${url} `, data);
                throw new Error(data.message || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url} `, error);
            throw error;
        }
    }

    async getRoles(params?: GetRolesParams): Promise<PaginatedResponse<Role>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const response = await this.fetch<PaginatedResponse<Role>>(`/roles?${searchParams.toString()}`);
        return {
            ...response,
            data: response.data || []
        };
    }

    async getRoleById(id: number): Promise<SingleResponse<Role>> {
        return this.fetch<SingleResponse<Role>>(`/roles /${id} `);
    }

    async getRolesByScope(scopeId: number): Promise<PaginatedResponse<Role>> {
        return this.fetch<PaginatedResponse<Role>>(`/roles /scope /${scopeId} `);
    }

    async getRolesByLevel(level: number): Promise<PaginatedResponse<Role>> {
        return this.fetch<PaginatedResponse<Role>>(`/roles /level /${level} `);
    }

    async getSystemRoles(): Promise<PaginatedResponse<Role>> {
        return this.fetch<PaginatedResponse<Role>>('/roles/system');
    }

    async createRole(data: CreateRoleDTO): Promise<SingleResponse<Role>> {
        return this.fetch<SingleResponse<Role>>('/roles', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateRole(id: number, data: UpdateRoleDTO): Promise<SingleResponse<Role>> {
        return this.fetch<SingleResponse<Role>>(`/roles /${id} `, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    // Assuming DELETE exists for consistency, though not in docs
    async deleteRole(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/roles /${id} `, {
            method: 'DELETE',
        });
    }

    async assignPermissions(id: number, data: AssignPermissionsDTO): Promise<AssignPermissionsResponse> {
        return this.fetch<AssignPermissionsResponse>(`/roles /${id}/permissions`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async getRolePermissions(id: number): Promise<RolePermissionsResponse> {
        return this.fetch<RolePermissionsResponse>(`/roles/${id}/permissions`);
    }

    async removePermissionFromRole(roleId: number, permissionId: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/roles/${roleId}/permissions/${permissionId}`, {
            method: 'DELETE',
        });
    }
}
