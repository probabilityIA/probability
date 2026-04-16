import { env } from '@/shared/config/env';
import { IVehicleRepository } from '../../domain/ports';
import {
    VehicleInfo,
    VehiclesListResponse,
    GetVehiclesParams,
    CreateVehicleDTO,
    UpdateVehicleDTO,
    DeleteVehicleResponse,
} from '../../domain/types';

export class VehicleApiRepository implements IVehicleRepository {
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

        try {
            const res = await fetch(url, { ...options, headers });
            const data = await res.json();

            if (!res.ok) {
                throw new Error(data.error || data.message || 'An error occurred');
            }

            return data;
        } catch (error) {
            throw error;
        }
    }

    /** Agrega ?business_id=X a la url si se provee (para super admin) */
    private withBusinessId(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
    }

    async getVehicles(params?: GetVehiclesParams): Promise<VehiclesListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<VehiclesListResponse>(`/vehicles${query ? `?${query}` : ''}`);
    }

    async getVehicleById(id: number, businessId?: number): Promise<VehicleInfo> {
        return this.fetch<VehicleInfo>(this.withBusinessId(`/vehicles/${id}`, businessId));
    }

    async createVehicle(data: CreateVehicleDTO, businessId?: number): Promise<VehicleInfo> {
        return this.fetch<VehicleInfo>(this.withBusinessId('/vehicles', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateVehicle(id: number, data: UpdateVehicleDTO, businessId?: number): Promise<VehicleInfo> {
        return this.fetch<VehicleInfo>(this.withBusinessId(`/vehicles/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteVehicle(id: number, businessId?: number): Promise<DeleteVehicleResponse> {
        return this.fetch<DeleteVehicleResponse>(this.withBusinessId(`/vehicles/${id}`, businessId), {
            method: 'DELETE',
        });
    }
}
