import { env } from '@/shared/config/env';
import {
    AddToLPNInput,
    CreateLPNDTO,
    GetLPNsParams,
    GetSyncLogsParams,
    InboundSyncInput,
    InboundSyncResult,
    InventorySyncLogListResponse,
    LicensePlate,
    LicensePlateLine,
    LicensePlateListResponse,
    MergeLPNInput,
    MoveLPNInput,
    ScanInput,
    ScanResult,
    UpdateLPNDTO,
} from '../../domain/capture-types';

export class CaptureApiRepository {
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

    async listLPNs(params: GetLPNsParams = {}): Promise<LicensePlateListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.lpn_type) qs.set('lpn_type', params.lpn_type);
        if (params.status) qs.set('status', params.status);
        if (params.location_id) qs.set('location_id', String(params.location_id));
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<LicensePlateListResponse>(`/inventory/lpn${suffix}`);
    }

    async createLPN(data: CreateLPNDTO, businessId?: number): Promise<LicensePlate> {
        return this.request<LicensePlate>(`/inventory/lpn${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async getLPN(id: number, businessId?: number): Promise<LicensePlate> {
        return this.request<LicensePlate>(`/inventory/lpn/${id}${this.businessQuery(businessId)}`);
    }

    async updateLPN(id: number, data: UpdateLPNDTO, businessId?: number): Promise<LicensePlate> {
        return this.request<LicensePlate>(`/inventory/lpn/${id}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteLPN(id: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/inventory/lpn/${id}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async addToLPN(id: number, data: AddToLPNInput, businessId?: number): Promise<LicensePlateLine> {
        return this.request<LicensePlateLine>(`/inventory/lpn/${id}/lines${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async moveLPN(id: number, data: MoveLPNInput, businessId?: number): Promise<LicensePlate> {
        return this.request<LicensePlate>(`/inventory/lpn/${id}/move${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async dissolveLPN(id: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/inventory/lpn/${id}/dissolve${this.businessQuery(businessId)}`, { method: 'POST' });
    }

    async mergeLPN(id: number, data: MergeLPNInput, businessId?: number): Promise<LicensePlate> {
        return this.request<LicensePlate>(`/inventory/lpn/${id}/merge${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async scan(data: ScanInput, businessId?: number): Promise<ScanResult> {
        return this.request<ScanResult>(`/inventory/scan${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async inboundSync(integrationId: number, data: InboundSyncInput, businessId?: number): Promise<InboundSyncResult> {
        return this.request<InboundSyncResult>(`/inventory/sync/inbound/${integrationId}${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async listSyncLogs(params: GetSyncLogsParams = {}): Promise<InventorySyncLogListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.integration_id) qs.set('integration_id', String(params.integration_id));
        if (params.direction) qs.set('direction', params.direction);
        if (params.status) qs.set('status', params.status);
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<InventorySyncLogListResponse>(`/inventory/sync/logs${suffix}`);
    }
}
