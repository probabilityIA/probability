import { env } from '@/shared/config/env';
import {
    CodSummary,
    CodOrder,
    CarrierConfig,
    PaymentCut,
    ReportFilters,
    CodOrdersParams,
    SaveCarrierConfigInput,
    Paginated,
    SingleResult,
    CutsResult,
} from '../../domain/types';

export class CodReportApiRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };
        if (this.token) headers['Authorization'] = `Bearer ${this.token}`;
        const res = await fetch(`${this.baseUrl}${path}`, { ...options, headers, cache: 'no-store' });
        const data = await res.json().catch(() => ({}));
        if (!res.ok) throw new Error(data.message || data.error || 'Error en la solicitud');
        return data as T;
    }

    private rangeParams(f: ReportFilters): URLSearchParams {
        const sp = new URLSearchParams();
        if (f.range === 'custom' && f.start_date && f.end_date) {
            sp.append('start_date', f.start_date);
            sp.append('end_date', f.end_date);
        } else if (f.range) {
            sp.append('range', f.range);
        }
        if (f.carrier) sp.append('carrier', f.carrier);
        if (f.business_id) sp.append('business_id', String(f.business_id));
        return sp;
    }

    async getSummary(f: ReportFilters): Promise<SingleResult<CodSummary>> {
        const sp = this.rangeParams(f);
        return this.request<SingleResult<CodSummary>>(`/cod-report/summary?${sp.toString()}`);
    }

    async getOrders(p: CodOrdersParams): Promise<Paginated<CodOrder>> {
        const sp = this.rangeParams(p);
        if (p.page) sp.append('page', String(p.page));
        if (p.page_size) sp.append('page_size', String(p.page_size));
        if (p.collected !== undefined) sp.append('collected', String(p.collected));
        if (p.search) sp.append('search', p.search);
        return this.request<Paginated<CodOrder>>(`/cod-report/orders?${sp.toString()}`);
    }

    async getCuts(businessId?: number): Promise<CutsResult> {
        const sp = new URLSearchParams();
        if (businessId) sp.append('business_id', String(businessId));
        return this.request<CutsResult>(`/cod-report/cuts?${sp.toString()}`);
    }

    async confirmCut(periodStart: string, periodEnd: string, businessId?: number): Promise<SingleResult<PaymentCut>> {
        const sp = new URLSearchParams();
        if (businessId) sp.append('business_id', String(businessId));
        return this.request<SingleResult<PaymentCut>>(`/cod-report/cuts/confirm?${sp.toString()}`, {
            method: 'POST',
            body: JSON.stringify({ period_start: periodStart, period_end: periodEnd }),
        });
    }

    async getCarrierConfigs(businessId?: number): Promise<SingleResult<CarrierConfig[]>> {
        const sp = new URLSearchParams();
        if (businessId) sp.append('business_id', String(businessId));
        return this.request<SingleResult<CarrierConfig[]>>(`/cod-report/carrier-config?${sp.toString()}`);
    }

    async saveCarrierConfig(input: SaveCarrierConfigInput, businessId?: number): Promise<SingleResult<CarrierConfig>> {
        const sp = new URLSearchParams();
        if (businessId) sp.append('business_id', String(businessId));
        return this.request<SingleResult<CarrierConfig>>(`/cod-report/carrier-config?${sp.toString()}`, {
            method: 'PUT',
            body: JSON.stringify(input),
        });
    }
}
