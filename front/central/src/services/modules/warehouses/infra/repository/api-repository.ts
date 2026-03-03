import { env } from '@/shared/config/env';
import { IWarehouseRepository } from '../../domain/ports';
import {
    Warehouse,
    WarehouseDetail,
    WarehouseLocation,
    WarehousesListResponse,
    GetWarehousesParams,
    CreateWarehouseDTO,
    UpdateWarehouseDTO,
    CreateLocationDTO,
    UpdateLocationDTO,
} from '../../domain/types';

export class WarehouseApiRepository implements IWarehouseRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(url, { ...options, headers });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.error || data.message || 'An error occurred');
        }

        return data;
    }

    private withBusinessId(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
    }

    async getWarehouses(params?: GetWarehousesParams): Promise<WarehousesListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<WarehousesListResponse>(`/warehouses${query ? `?${query}` : ''}`);
    }

    async getWarehouseById(id: number, businessId?: number): Promise<WarehouseDetail> {
        return this.fetch<WarehouseDetail>(this.withBusinessId(`/warehouses/${id}`, businessId));
    }

    async createWarehouse(data: CreateWarehouseDTO, businessId?: number): Promise<Warehouse> {
        return this.fetch<Warehouse>(this.withBusinessId('/warehouses', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateWarehouse(id: number, data: UpdateWarehouseDTO, businessId?: number): Promise<Warehouse> {
        return this.fetch<Warehouse>(this.withBusinessId(`/warehouses/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteWarehouse(id: number, businessId?: number): Promise<void> {
        await this.fetch<void>(this.withBusinessId(`/warehouses/${id}`, businessId), {
            method: 'DELETE',
        });
    }

    async getLocations(warehouseId: number, businessId?: number): Promise<WarehouseLocation[]> {
        return this.fetch<WarehouseLocation[]>(this.withBusinessId(`/warehouses/${warehouseId}/locations`, businessId));
    }

    async createLocation(warehouseId: number, data: CreateLocationDTO, businessId?: number): Promise<WarehouseLocation> {
        return this.fetch<WarehouseLocation>(this.withBusinessId(`/warehouses/${warehouseId}/locations`, businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateLocation(warehouseId: number, locationId: number, data: UpdateLocationDTO, businessId?: number): Promise<WarehouseLocation> {
        return this.fetch<WarehouseLocation>(this.withBusinessId(`/warehouses/${warehouseId}/locations/${locationId}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteLocation(warehouseId: number, locationId: number, businessId?: number): Promise<void> {
        await this.fetch<void>(this.withBusinessId(`/warehouses/${warehouseId}/locations/${locationId}`, businessId), {
            method: 'DELETE',
        });
    }
}
