import { env } from '@/shared/config/env';
import { IInventoryRepository } from '../../domain/ports';
import {
    InventoryLevel,
    InventoryListResponse,
    MovementListResponse,
    MovementTypeListResponse,
    StockMovement,
    GetInventoryParams,
    GetMovementsParams,
    AdjustStockDTO,
    TransferStockDTO,
} from '../../domain/types';

export class InventoryApiRepository implements IInventoryRepository {
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

    async getProductInventory(productId: string, businessId?: number): Promise<InventoryLevel[]> {
        return this.fetch<InventoryLevel[]>(this.withBusinessId(`/inventory/product/${productId}`, businessId));
    }

    async getWarehouseInventory(warehouseId: number, params?: GetInventoryParams): Promise<InventoryListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<InventoryListResponse>(`/inventory/warehouse/${warehouseId}${query ? `?${query}` : ''}`);
    }

    async adjustStock(data: AdjustStockDTO, businessId?: number): Promise<StockMovement> {
        return this.fetch<StockMovement>(this.withBusinessId('/inventory/adjust', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async transferStock(data: TransferStockDTO, businessId?: number): Promise<{ message: string }> {
        return this.fetch<{ message: string }>(this.withBusinessId('/inventory/transfer', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async getMovements(params?: GetMovementsParams): Promise<MovementListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<MovementListResponse>(`/inventory/movements${query ? `?${query}` : ''}`);
    }

    async getMovementTypes(params?: { page?: number; page_size?: number; active_only?: boolean; business_id?: number }): Promise<MovementTypeListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null) {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<MovementTypeListResponse>(`/inventory/movement-types${query ? `?${query}` : ''}`);
    }
}
