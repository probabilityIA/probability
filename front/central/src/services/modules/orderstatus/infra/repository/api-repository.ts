import { env } from '@/shared/config/env';
import { IOrderStatusMappingRepository } from '../../domain/ports';
import {
    OrderStatusMapping,
    PaginatedResponse,
    GetOrderStatusMappingsParams,
    SingleResponse,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO,
    ActionResponse,
    OrderStatusInfo,
    CreateOrderStatusDTO,
    UpdateOrderStatusDTO,
    EcommerceIntegrationType,
    ChannelStatusInfo,
    CreateChannelStatusDTO,
    UpdateChannelStatusDTO
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

    async getOrderStatuses(isActive?: boolean): Promise<{ success: boolean; data: OrderStatusInfo[]; message?: string }> {
        const params = new URLSearchParams();
        if (isActive !== undefined) {
            params.append('is_active', String(isActive));
        }
        const url = `/order-statuses${params.toString() ? `?${params.toString()}` : ''}`;
        return this.fetch<{ success: boolean; data: OrderStatusInfo[]; message?: string }>(url);
    }

    async toggleOrderStatusMappingActive(id: number): Promise<SingleResponse<OrderStatusMapping>> {
        return this.fetch<SingleResponse<OrderStatusMapping>>(`/order-status-mappings/${id}/toggle`, {
            method: 'PATCH',
        });
    }

    async createOrderStatus(data: CreateOrderStatusDTO): Promise<SingleResponse<OrderStatusInfo>> {
        return this.fetch<SingleResponse<OrderStatusInfo>>('/order-statuses', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async getOrderStatusById(id: number): Promise<SingleResponse<OrderStatusInfo>> {
        return this.fetch<SingleResponse<OrderStatusInfo>>(`/order-statuses/${id}`);
    }

    async updateOrderStatus(id: number, data: UpdateOrderStatusDTO): Promise<SingleResponse<OrderStatusInfo>> {
        return this.fetch<SingleResponse<OrderStatusInfo>>(`/order-statuses/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteOrderStatus(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/order-statuses/${id}`, {
            method: 'DELETE',
        });
    }

    async getEcommerceIntegrationTypes(): Promise<{ success: boolean; data: EcommerceIntegrationType[]; message?: string }> {
        return this.fetch<{ success: boolean; data: EcommerceIntegrationType[]; message?: string }>('/ecommerce-integration-types');
    }

    async getChannelStatuses(integrationTypeId: number, isActive?: boolean): Promise<{ success: boolean; data: ChannelStatusInfo[]; message?: string }> {
        const params = new URLSearchParams({ integration_type_id: String(integrationTypeId) });
        if (isActive !== undefined) params.append('is_active', String(isActive));
        return this.fetch<{ success: boolean; data: ChannelStatusInfo[]; message?: string }>(`/channel-statuses?${params.toString()}`);
    }

    async createChannelStatus(data: CreateChannelStatusDTO): Promise<SingleResponse<ChannelStatusInfo>> {
        return this.fetch<SingleResponse<ChannelStatusInfo>>('/channel-statuses', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateChannelStatus(id: number, data: UpdateChannelStatusDTO): Promise<SingleResponse<ChannelStatusInfo>> {
        return this.fetch<SingleResponse<ChannelStatusInfo>>(`/channel-statuses/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteChannelStatus(id: number): Promise<ActionResponse> {
        return this.fetch<ActionResponse>(`/channel-statuses/${id}`, {
            method: 'DELETE',
        });
    }
}
