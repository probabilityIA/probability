import { env } from '@/shared/config/env';
import { IActionRepository } from '../../domain/ports';
import {
    Action,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    GetActionsParams,
    CreateActionDTO,
    UpdateActionDTO
} from '../../domain/types';

export class ActionApiRepository implements IActionRepository {
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

        try {
            const res = await fetch(url, {
                ...options,
                headers,
            });

            const data = await res.json();

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

    async getActions(params?: GetActionsParams): Promise<PaginatedResponse<Action>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        const response = await this.fetch<PaginatedResponse<Action>>(`/actions?${searchParams.toString()}`);
        return {
            ...response,
            data: {
                ...response.data,
                actions: response.data.actions || []
            }
        };
    }

    async getActionById(id: number): Promise<SingleResponse<Action>> {
        return this.fetch<SingleResponse<Action>>(`/actions/${id}`);
    }

    async createAction(data: CreateActionDTO): Promise<SingleResponse<Action>> {
        return this.fetch<SingleResponse<Action>>('/actions', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateAction(id: number, data: UpdateActionDTO): Promise<SingleResponse<Action>> {
        return this.fetch<SingleResponse<Action>>(`/actions/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteAction(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/actions/${id}`, {
            method: 'DELETE',
        });
    }
}
