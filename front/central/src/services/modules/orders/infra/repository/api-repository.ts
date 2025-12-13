import { env } from '@/shared/config/env';
import { IOrderRepository } from '../../domain/ports';
import {
    Order,
    PaginatedResponse,
    GetOrdersParams,
    SingleResponse,
    CreateOrderDTO,
    UpdateOrderDTO,
    ActionResponse
} from '../../domain/types';

export class OrderApiRepository implements IOrderRepository {
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

    async getOrders(params?: GetOrdersParams): Promise<PaginatedResponse<Order>> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) searchParams.append(key, String(value));
            });
        }
        return this.fetch<PaginatedResponse<Order>>(`/orders?${searchParams.toString()}`);
    }

    async getOrderById(id: string): Promise<SingleResponse<Order>> {
        return this.fetch<SingleResponse<Order>>(`/orders/${id}`);
    }

    async createOrder(data: CreateOrderDTO): Promise<SingleResponse<Order>> {
        return this.fetch<SingleResponse<Order>>('/orders', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateOrder(id: string, data: UpdateOrderDTO): Promise<SingleResponse<Order>> {
        return this.fetch<SingleResponse<Order>>(`/orders/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    deleteOrder(id: string): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/orders/${id}`, {
            method: 'DELETE',
        });
    }

    async getOrderRaw(id: string): Promise<SingleResponse<any>> {
        return this.fetch<SingleResponse<any>>(`/orders/${id}/raw`);
    }

    async getAIRecommendation(origin: string, destination: string): Promise<any> {
        const searchParams = new URLSearchParams({ origin, destination });
        return this.fetch<any>(`/ai/recommendation?${searchParams.toString()}`);
    }
}
