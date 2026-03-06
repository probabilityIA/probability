import { env } from '@/shared/config/env';
import { IRouteRepository } from '../../domain/ports';
import {
    RouteInfo,
    RouteDetail,
    RouteStopInfo,
    RoutesListResponse,
    GetRoutesParams,
    CreateRouteDTO,
    UpdateRouteDTO,
    AddStopDTO,
    UpdateStopDTO,
    UpdateStopStatusDTO,
    ReorderStopsDTO,
    DeleteRouteResponse,
    DriverOption,
    VehicleOption,
    AssignableOrder,
} from '../../domain/types';

export class RouteApiRepository implements IRouteRepository {
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

    // ============================================
    // Route CRUD
    // ============================================

    async getRoutes(params?: GetRoutesParams): Promise<RoutesListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<RoutesListResponse>(`/routes${query ? `?${query}` : ''}`);
    }

    async getRouteById(id: number, businessId?: number): Promise<RouteDetail> {
        return this.fetch<RouteDetail>(this.withBusinessId(`/routes/${id}`, businessId));
    }

    async createRoute(data: CreateRouteDTO, businessId?: number): Promise<RouteInfo> {
        return this.fetch<RouteInfo>(this.withBusinessId('/routes', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateRoute(id: number, data: UpdateRouteDTO, businessId?: number): Promise<RouteInfo> {
        return this.fetch<RouteInfo>(this.withBusinessId(`/routes/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteRoute(id: number, businessId?: number): Promise<DeleteRouteResponse> {
        return this.fetch<DeleteRouteResponse>(this.withBusinessId(`/routes/${id}`, businessId), {
            method: 'DELETE',
        });
    }

    // ============================================
    // Route lifecycle
    // ============================================

    async startRoute(id: number, businessId?: number): Promise<RouteDetail> {
        return this.fetch<RouteDetail>(this.withBusinessId(`/routes/${id}/start`, businessId), {
            method: 'POST',
        });
    }

    async completeRoute(id: number, businessId?: number): Promise<RouteDetail> {
        return this.fetch<RouteDetail>(this.withBusinessId(`/routes/${id}/complete`, businessId), {
            method: 'POST',
        });
    }

    // ============================================
    // Stop management
    // ============================================

    async addStop(routeId: number, data: AddStopDTO, businessId?: number): Promise<RouteStopInfo> {
        return this.fetch<RouteStopInfo>(this.withBusinessId(`/routes/${routeId}/stops`, businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateStop(routeId: number, stopId: number, data: UpdateStopDTO, businessId?: number): Promise<RouteStopInfo> {
        return this.fetch<RouteStopInfo>(this.withBusinessId(`/routes/${routeId}/stops/${stopId}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteStop(routeId: number, stopId: number, businessId?: number): Promise<DeleteRouteResponse> {
        return this.fetch<DeleteRouteResponse>(this.withBusinessId(`/routes/${routeId}/stops/${stopId}`, businessId), {
            method: 'DELETE',
        });
    }

    async updateStopStatus(routeId: number, stopId: number, data: UpdateStopStatusDTO, businessId?: number): Promise<RouteStopInfo> {
        return this.fetch<RouteStopInfo>(this.withBusinessId(`/routes/${routeId}/stops/${stopId}/status`, businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async reorderStops(routeId: number, data: ReorderStopsDTO, businessId?: number): Promise<RouteDetail> {
        return this.fetch<RouteDetail>(this.withBusinessId(`/routes/${routeId}/stops/reorder`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    // ============================================
    // Form options
    // ============================================

    async getAvailableDrivers(businessId?: number): Promise<DriverOption[]> {
        return this.fetch<DriverOption[]>(this.withBusinessId('/routes/available-drivers', businessId));
    }

    async getAvailableVehicles(businessId?: number): Promise<VehicleOption[]> {
        return this.fetch<VehicleOption[]>(this.withBusinessId('/routes/available-vehicles', businessId));
    }

    async getAssignableOrders(businessId?: number): Promise<AssignableOrder[]> {
        return this.fetch<AssignableOrder[]>(this.withBusinessId('/routes/assignable-orders', businessId));
    }
}
