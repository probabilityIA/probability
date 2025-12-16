import { env } from '@/shared/config/env';
import { IUserRepository } from '../../domain/ports';
import {
    User,
    PaginatedResponse,
    GetUsersParams,
    SingleResponse,
    CreateUserDTO,
    UpdateUserDTO,
    ActionResponse,
    AssignRolesDTO
} from '../../domain/types';

export class UserApiRepository implements IUserRepository {
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

    async getUsers(params?: GetUsersParams): Promise<PaginatedResponse<User>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const response = await this.fetch<PaginatedResponse<User>>(`/users?${searchParams.toString()}`);
        return {
            ...response,
            data: response.data || []
        };
    }

    async getUserById(id: number): Promise<SingleResponse<User>> {
        return this.fetch<SingleResponse<User>>(`/users/${id}`);
    }

    async createUser(data: CreateUserDTO): Promise<SingleResponse<User>> {
        const formData = new FormData();
        formData.append('name', data.name);
        formData.append('email', data.email);
        if (data.phone) formData.append('phone', data.phone);
        if (data.is_active !== undefined) formData.append('is_active', String(data.is_active));
        if (data.avatarFile) formData.append('avatarFile', data.avatarFile);
        if (data.scope_id !== undefined) formData.append('scope_id', String(data.scope_id));
        if (data.business_ids) formData.append('business_ids', data.business_ids.join(','));

        return this.fetch<SingleResponse<User>>('/users', {
            method: 'POST',
            body: formData,
            // Content-Type header is automatically set by browser for FormData with boundary
        });
    }

    async updateUser(id: number, data: UpdateUserDTO): Promise<SingleResponse<User>> {
        const formData = new FormData();
        if (data.name) formData.append('name', data.name);
        if (data.email) formData.append('email', data.email);
        if (data.phone) formData.append('phone', data.phone);
        if (data.is_active !== undefined) formData.append('is_active', String(data.is_active));
        if (data.remove_avatar !== undefined) formData.append('remove_avatar', String(data.remove_avatar));
        if (data.avatarFile) formData.append('avatarFile', data.avatarFile);
        if (data.business_ids) formData.append('business_ids', data.business_ids.join(','));

        return this.fetch<SingleResponse<User>>(`/users/${id}`, {
            method: 'PUT',
            body: formData,
        });
    }

    async deleteUser(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/users/${id}`, {
            method: 'DELETE',
        });
    }

    async assignRoles(id: number, data: AssignRolesDTO): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/users/${id}/assign-role`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(data),
        });
    }
}
