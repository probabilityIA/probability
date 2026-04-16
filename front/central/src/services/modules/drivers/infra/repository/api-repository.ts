import { env } from '@/shared/config/env';
import { IDriverRepository } from '../../domain/ports';
import {
    DriverInfo,
    DriversListResponse,
    GetDriversParams,
    CreateDriverDTO,
    UpdateDriverDTO,
    DeleteDriverResponse,
} from '../../domain/types';

export class DriverApiRepository implements IDriverRepository {
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

    async getDrivers(params?: GetDriversParams): Promise<DriversListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<DriversListResponse>(`/drivers${query ? `?${query}` : ''}`);
    }

    async getDriverById(id: number, businessId?: number): Promise<DriverInfo> {
        return this.fetch<DriverInfo>(this.withBusinessId(`/drivers/${id}`, businessId));
    }

    async createDriver(data: CreateDriverDTO, businessId?: number): Promise<DriverInfo> {
        return this.fetch<DriverInfo>(this.withBusinessId('/drivers', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateDriver(id: number, data: UpdateDriverDTO, businessId?: number): Promise<DriverInfo> {
        return this.fetch<DriverInfo>(this.withBusinessId(`/drivers/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteDriver(id: number, businessId?: number): Promise<DeleteDriverResponse> {
        return this.fetch<DeleteDriverResponse>(this.withBusinessId(`/drivers/${id}`, businessId), {
            method: 'DELETE',
        });
    }
}
