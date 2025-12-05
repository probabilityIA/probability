import { env } from '@/shared/config/env';
import { IOrderStatusMappingRepository } from '../../domain/ports';
import {
    OrderStatusMapping,
    PaginatedResponse,
    GetOrderStatusMappingsParams,
    SingleResponse,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO,
    ActionResponse
} from '../../domain/types';

export class OrderStatusMappingApiRepository implements IOrderStatusMappingRepository {
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

    async getOrderStatusMappings(params?: GetOrderStatusMappingsParams): Promise<PaginatedResponse<OrderStatusMapping>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<OrderStatusMapping>>(`/order-status-mappings?${searchParams.toString()}`);
    }

    async getOrderStatusMappingById(id: number): Promise<SingleResponse<OrderStatusMapping>> {
        return this.fetch<SingleResponse<OrderStatusMapping>>(`/order-status-mappings/${id}`);
    }

    async createOrderStatusMapping(data: CreateOrderStatusMappingDTO): Promise<SingleResponse<OrderStatusMapping>> {
        return this.fetch<SingleResponse<OrderStatusMapping>>('/order-status-mappings', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateOrderStatusMapping(id: number, data: UpdateOrderStatusMappingDTO): Promise<SingleResponse<OrderStatusMapping>> {
        return this.fetch<SingleResponse<OrderStatusMapping>>(`/order-status-mappings/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteOrderStatusMapping(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/order-status-mappings/${id}`, {
            method: 'DELETE',
        });
    }

    async toggleOrderStatusMappingActive(id: number): Promise<SingleResponse<OrderStatusMapping>> {
        return this.fetch<SingleResponse<OrderStatusMapping>>(`/order-status-mappings/${id}/toggle`, {
            method: 'PATCH',
        });
    }
}
