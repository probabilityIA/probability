import { env } from '@/shared/config/env';
import { IGeozoneRepository } from '../../domain/ports';
import {
    Geozone,
    GeozonesListResponse,
    GetGeozonesParams,
    CreateGeozoneDTO,
    LookupParams,
    LookupResponse,
    BulkImportRequest,
    BulkImportResponse,
    DisplayFeatureCollection,
    GeozoneType,
} from '../../domain/types';

export class GeozoneApiRepository implements IGeozoneRepository {
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
            ...((options.headers as Record<string, string>) || {}),
        };
        if (this.token) headers['Authorization'] = `Bearer ${this.token}`;
        const res = await fetch(url, { ...options, headers, cache: 'no-store' });
        const text = await res.text();
        const data = text ? JSON.parse(text) : null;
        if (!res.ok) {
            const msg = (data && (data.error || data.message)) || `HTTP ${res.status}`;
            throw new Error(msg);
        }
        return data as T;
    }

    private withBusinessId(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
    }

    async list(params?: GetGeozonesParams): Promise<GeozonesListResponse> {
        const sp = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([k, v]) => {
                if (v !== undefined && v !== null && v !== '') sp.append(k, String(v));
            });
        }
        const q = sp.toString();
        return this.fetch<GeozonesListResponse>(`/geozones${q ? `?${q}` : ''}`);
    }

    async getById(id: number, includeGeom: boolean, businessId?: number): Promise<Geozone> {
        const sp = new URLSearchParams({ include_geometry: String(includeGeom) });
        if (businessId) sp.append('business_id', String(businessId));
        return this.fetch<Geozone>(`/geozones/${id}?${sp.toString()}`);
    }

    async create(data: CreateGeozoneDTO, businessId?: number): Promise<Geozone> {
        return this.fetch<Geozone>(this.withBusinessId('/geozones', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async bulkImport(data: BulkImportRequest, businessId?: number): Promise<BulkImportResponse> {
        return this.fetch<BulkImportResponse>(this.withBusinessId('/geozones/bulk', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async lookup(params: LookupParams): Promise<LookupResponse> {
        const sp = new URLSearchParams({
            lat: String(params.lat),
            lng: String(params.lng),
        });
        if (params.type) sp.append('type', params.type);
        if (params.business_id) sp.append('business_id', String(params.business_id));
        return this.fetch<LookupResponse>(`/geozones/lookup?${sp.toString()}`);
    }

    async getForDisplay(geozoneType: GeozoneType | '', zoom: number, bbox?: string): Promise<DisplayFeatureCollection> {
        const sp = new URLSearchParams({ zoom: String(zoom) });
        if (geozoneType) sp.append('type', geozoneType);
        if (bbox) sp.append('bbox', bbox);
        return this.fetch<DisplayFeatureCollection>(`/geozones/display?${sp.toString()}`);
    }

    async remove(id: number, businessId?: number): Promise<void> {
        await this.fetch<void>(this.withBusinessId(`/geozones/${id}`, businessId), {
            method: 'DELETE',
        });
    }

    async probabilityByCarrier(orderId: string, businessId: number): Promise<import('../../domain/types').ProbabilityResult[]> {
        const sp = new URLSearchParams({ order_id: orderId, business_id: String(businessId) });
        const res = await this.fetch<{ success: boolean; data: import('../../domain/types').ProbabilityResult[] }>(`/geozones/probability/by-carrier?${sp.toString()}`);
        return res.data || [];
    }

    async getOrderZone(orderId: string, businessId: number): Promise<Geozone | null> {
        const sp = new URLSearchParams({ order_id: orderId, business_id: String(businessId) });
        const res = await this.fetch<{ success: boolean; data: Geozone | null }>(`/geozones/order-zone?${sp.toString()}`);
        return res.data;
    }

    async probability(req: import('../../domain/types').ProbabilityRequest): Promise<import('../../domain/types').ProbabilityResult> {
        const sp = new URLSearchParams();
        sp.append('business_id', String(req.business_id));
        if (req.order_id) sp.append('order_id', req.order_id);
        if (req.lat !== undefined) sp.append('lat', String(req.lat));
        if (req.lng !== undefined) sp.append('lng', String(req.lng));
        if (req.carrier) sp.append('carrier', req.carrier);
        const res = await this.fetch<{ success: boolean; data: import('../../domain/types').ProbabilityResult }>(`/geozones/probability?${sp.toString()}`);
        return res.data;
    }
}
