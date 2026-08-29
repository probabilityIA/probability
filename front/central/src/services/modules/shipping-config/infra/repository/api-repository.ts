import { IShippingConfigRepository } from '../../domain/ports';
import { ShippingConfigOverview, ShippingConfig, SaveShippingConfigRequest } from '../../domain/types';
import { env } from '@/shared/config/env';

export class ShippingConfigApiRepository implements IShippingConfigRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
        const headers: Record<string, string> = {
            'Content-Type': 'application/json',
            ...options.headers as Record<string, string>,
        };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const response = await fetch(`${this.baseUrl}${endpoint}`, { ...options, headers, cache: 'no-store' });
        const data = await response.json();
        if (!response.ok) {
            throw new Error(data.error || data.message || 'Error en la configuración de envíos');
        }
        return data.data as T;
    }

    private scope(businessId?: number): string {
        return businessId ? `?business_id=${businessId}` : '';
    }

    async getOverview(businessId?: number): Promise<ShippingConfigOverview> {
        return this.fetch<ShippingConfigOverview>(`/shipping-config${this.scope(businessId)}`);
    }

    async saveBusinessConfig(req: SaveShippingConfigRequest, businessId?: number): Promise<ShippingConfig> {
        return this.fetch<ShippingConfig>(`/shipping-config${this.scope(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(req),
        });
    }

    async saveWarehouseConfig(warehouseId: number, req: SaveShippingConfigRequest, businessId?: number): Promise<ShippingConfig> {
        return this.fetch<ShippingConfig>(`/shipping-config/warehouses/${warehouseId}${this.scope(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(req),
        });
    }

    async deleteWarehouseConfig(warehouseId: number, businessId?: number): Promise<void> {
        await this.fetch<unknown>(`/shipping-config/warehouses/${warehouseId}${this.scope(businessId)}`, {
            method: 'DELETE',
        });
    }

    async setDefaultWarehouse(warehouseId: number, businessId?: number): Promise<void> {
        await this.fetch<unknown>(`/shipping-config/warehouses/${warehouseId}/default${this.scope(businessId)}`, {
            method: 'PUT',
        });
    }
}
