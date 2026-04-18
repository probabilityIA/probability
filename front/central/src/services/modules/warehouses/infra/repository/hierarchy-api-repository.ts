import { env } from '@/shared/config/env';
import {
    Aisle,
    CreateAisleDTO,
    CreateRackDTO,
    CreateRackLevelDTO,
    CreateZoneDTO,
    CubingCheckResult,
    PaginatedResponse,
    Rack,
    RackLevel,
    UpdateAisleDTO,
    UpdateRackDTO,
    UpdateRackLevelDTO,
    UpdateZoneDTO,
    ValidateCubingInput,
    WarehouseTree,
    Zone,
} from '../../domain/hierarchy-types';

export class HierarchyApiRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
        const headers: Record<string, string> = {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };
        if (this.token) headers['Authorization'] = `Bearer ${this.token}`;

        const res = await fetch(`${this.baseUrl}${path}`, { ...options, headers });
        const data = await res.json().catch(() => ({}));
        if (!res.ok) {
            throw new Error((data && (data.error || data.message)) || `HTTP ${res.status}`);
        }
        return data as T;
    }

    private businessQuery(businessId?: number): string {
        return businessId ? `?business_id=${businessId}` : '';
    }

    async getTree(warehouseId: number, businessId?: number): Promise<WarehouseTree> {
        return this.request<WarehouseTree>(`/warehouses/${warehouseId}/tree${this.businessQuery(businessId)}`);
    }

    async listZones(warehouseId: number, params: { page?: number; page_size?: number; active_only?: boolean } = {}, businessId?: number): Promise<PaginatedResponse<Zone>> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.active_only) qs.set('active_only', 'true');
        if (businessId) qs.set('business_id', String(businessId));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<PaginatedResponse<Zone>>(`/warehouses/${warehouseId}/zones${suffix}`);
    }

    async createZone(data: CreateZoneDTO, businessId?: number): Promise<Zone> {
        return this.request<Zone>(`/zones${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async getZone(zoneId: number, businessId?: number): Promise<Zone> {
        return this.request<Zone>(`/zones/${zoneId}${this.businessQuery(businessId)}`);
    }

    async updateZone(zoneId: number, data: UpdateZoneDTO, businessId?: number): Promise<Zone> {
        return this.request<Zone>(`/zones/${zoneId}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteZone(zoneId: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/zones/${zoneId}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async listAisles(zoneId: number, params: { page?: number; page_size?: number } = {}, businessId?: number): Promise<PaginatedResponse<Aisle>> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (businessId) qs.set('business_id', String(businessId));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<PaginatedResponse<Aisle>>(`/zones/${zoneId}/aisles${suffix}`);
    }

    async createAisle(data: CreateAisleDTO, businessId?: number): Promise<Aisle> {
        return this.request<Aisle>(`/aisles${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateAisle(aisleId: number, data: UpdateAisleDTO, businessId?: number): Promise<Aisle> {
        return this.request<Aisle>(`/aisles/${aisleId}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteAisle(aisleId: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/aisles/${aisleId}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async listRacks(aisleId: number, params: { page?: number; page_size?: number } = {}, businessId?: number): Promise<PaginatedResponse<Rack>> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (businessId) qs.set('business_id', String(businessId));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<PaginatedResponse<Rack>>(`/aisles/${aisleId}/racks${suffix}`);
    }

    async createRack(data: CreateRackDTO, businessId?: number): Promise<Rack> {
        return this.request<Rack>(`/racks${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateRack(rackId: number, data: UpdateRackDTO, businessId?: number): Promise<Rack> {
        return this.request<Rack>(`/racks/${rackId}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteRack(rackId: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/racks/${rackId}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async listRackLevels(rackId: number, params: { page?: number; page_size?: number } = {}, businessId?: number): Promise<PaginatedResponse<RackLevel>> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (businessId) qs.set('business_id', String(businessId));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<PaginatedResponse<RackLevel>>(`/racks/${rackId}/levels${suffix}`);
    }

    async createRackLevel(data: CreateRackLevelDTO, businessId?: number): Promise<RackLevel> {
        return this.request<RackLevel>(`/rack-levels${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateRackLevel(levelId: number, data: UpdateRackLevelDTO, businessId?: number): Promise<RackLevel> {
        return this.request<RackLevel>(`/rack-levels/${levelId}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteRackLevel(levelId: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/rack-levels/${levelId}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async validateCubing(input: ValidateCubingInput, businessId?: number): Promise<CubingCheckResult> {
        return this.request<CubingCheckResult>(`/inventory/positions/validate-cubing${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(input),
        });
    }
}
