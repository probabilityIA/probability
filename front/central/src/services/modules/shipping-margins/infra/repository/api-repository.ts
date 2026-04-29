import { env } from '@/shared/config/env';
import { IShippingMarginRepository } from '../../domain/ports';
import {
    ShippingMargin,
    ShippingMarginsListResponse,
    GetShippingMarginsParams,
    CreateShippingMarginDTO,
    UpdateShippingMarginDTO,
    ProfitReportParams,
    ProfitReportResponse,
    ProfitReportDetailParams,
    ProfitReportDetailResponse,
} from '../../domain/types';

export class ShippingMarginApiRepository implements IShippingMarginRepository {
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

    async list(params?: GetShippingMarginsParams): Promise<ShippingMarginsListResponse> {
        const sp = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([k, v]) => {
                if (v !== undefined && v !== null && v !== '') sp.append(k, String(v));
            });
        }
        const q = sp.toString();
        return this.fetch<ShippingMarginsListResponse>(`/shipping-margins${q ? `?${q}` : ''}`);
    }

    async getById(id: number, businessId?: number): Promise<ShippingMargin> {
        return this.fetch<ShippingMargin>(this.withBusinessId(`/shipping-margins/${id}`, businessId));
    }

    async create(data: CreateShippingMarginDTO, businessId?: number): Promise<ShippingMargin> {
        return this.fetch<ShippingMargin>(this.withBusinessId('/shipping-margins', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async update(id: number, data: UpdateShippingMarginDTO, businessId?: number): Promise<ShippingMargin> {
        return this.fetch<ShippingMargin>(this.withBusinessId(`/shipping-margins/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async profitReport(params: ProfitReportParams): Promise<ProfitReportResponse> {
        const sp = new URLSearchParams();
        Object.entries(params).forEach(([k, v]) => {
            if (v !== undefined && v !== null && v !== '') sp.append(k, String(v));
        });
        const q = sp.toString();
        return this.fetch<ProfitReportResponse>(`/shipping-margins/profit-report${q ? `?${q}` : ''}`);
    }

    async profitReportDetail(params: ProfitReportDetailParams): Promise<ProfitReportDetailResponse> {
        const sp = new URLSearchParams();
        Object.entries(params).forEach(([k, v]) => {
            if (v !== undefined && v !== null && v !== '') sp.append(k, String(v));
        });
        const q = sp.toString();
        return this.fetch<ProfitReportDetailResponse>(`/shipping-margins/profit-report/detail${q ? `?${q}` : ''}`);
    }
}
