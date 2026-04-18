import { env } from '@/shared/config/env';
import {
    ChangeInventoryStateDTO,
    ConvertUoMInput,
    ConvertUoMResult,
    CreateLotDTO,
    CreateProductUoMDTO,
    CreateSerialDTO,
    GetLotsParams,
    GetSerialsParams,
    InventoryLot,
    InventorySerial,
    InventoryState,
    LotListResponse,
    ProductUoM,
    SerialListResponse,
    UnitOfMeasure,
    UpdateLotDTO,
    UpdateSerialDTO,
} from '../../domain/traceability-types';

export class TraceabilityApiRepository {
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

    async listLots(params: GetLotsParams = {}): Promise<LotListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.product_id) qs.set('product_id', params.product_id);
        if (params.status) qs.set('status', params.status);
        if (params.expiring_in_days) qs.set('expiring_in_days', String(params.expiring_in_days));
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<LotListResponse>(`/inventory/lots${suffix}`);
    }

    async createLot(data: CreateLotDTO, businessId?: number): Promise<InventoryLot> {
        return this.request<InventoryLot>(`/inventory/lots${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateLot(id: number, data: UpdateLotDTO, businessId?: number): Promise<InventoryLot> {
        return this.request<InventoryLot>(`/inventory/lots/${id}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteLot(id: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/inventory/lots/${id}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async listSerials(params: GetSerialsParams = {}): Promise<SerialListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.product_id) qs.set('product_id', params.product_id);
        if (params.lot_id) qs.set('lot_id', String(params.lot_id));
        if (params.state_id) qs.set('state_id', String(params.state_id));
        if (params.location_id) qs.set('location_id', String(params.location_id));
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<SerialListResponse>(`/inventory/serials${suffix}`);
    }

    async createSerial(data: CreateSerialDTO, businessId?: number): Promise<InventorySerial> {
        return this.request<InventorySerial>(`/inventory/serials${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateSerial(id: number, data: UpdateSerialDTO, businessId?: number): Promise<InventorySerial> {
        return this.request<InventorySerial>(`/inventory/serials/${id}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async listStates(): Promise<{ data: InventoryState[] }> {
        return this.request<{ data: InventoryState[] }>(`/inventory/states`);
    }

    async changeState(data: ChangeInventoryStateDTO, businessId?: number) {
        return this.request(`/inventory/state-transitions${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async listUoMs(): Promise<{ data: UnitOfMeasure[] }> {
        return this.request<{ data: UnitOfMeasure[] }>(`/inventory/uoms`);
    }

    async listProductUoMs(productId: string, businessId?: number): Promise<{ data: ProductUoM[] }> {
        return this.request<{ data: ProductUoM[] }>(`/inventory/products/${productId}/uoms${this.businessQuery(businessId)}`);
    }

    async createProductUoM(productId: string, data: CreateProductUoMDTO, businessId?: number): Promise<ProductUoM> {
        return this.request<ProductUoM>(`/inventory/products/${productId}/uoms${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async deleteProductUoM(id: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/inventory/product-uoms/${id}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async convertUoM(data: ConvertUoMInput, businessId?: number): Promise<ConvertUoMResult> {
        return this.request<ConvertUoMResult>(`/inventory/uoms/convert${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }
}
