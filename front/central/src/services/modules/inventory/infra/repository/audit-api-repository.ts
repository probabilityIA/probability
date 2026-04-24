import { env } from '@/shared/config/env';
import {
    ApproveDiscrepancyInput,
    CreateCycleCountPlanDTO,
    CycleCountLineListResponse,
    CycleCountPlan,
    CycleCountPlanListResponse,
    CycleCountTask,
    CycleCountTaskListResponse,
    GenerateCountTaskInput,
    GenerateCountTaskResult,
    GetCountLinesParams,
    GetCountPlansParams,
    GetCountTasksParams,
    GetDiscrepanciesParams,
    InventoryDiscrepancy,
    InventoryDiscrepancyListResponse,
    KardexExportResult,
    KardexQueryInput,
    RejectDiscrepancyInput,
    StartCountTaskInput,
    SubmitCountLineInput,
    SubmitCountLineResult,
    UpdateCycleCountPlanDTO,
} from '../../domain/audit-types';

export class AuditApiRepository {
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

    async listCountPlans(params: GetCountPlansParams = {}): Promise<CycleCountPlanListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.warehouse_id) qs.set('warehouse_id', String(params.warehouse_id));
        if (params.active_only) qs.set('active_only', 'true');
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<CycleCountPlanListResponse>(`/inventory/cycle-count-plans${suffix}`);
    }

    async createCountPlan(data: CreateCycleCountPlanDTO, businessId?: number): Promise<CycleCountPlan> {
        return this.request<CycleCountPlan>(`/inventory/cycle-count-plans${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateCountPlan(id: number, data: UpdateCycleCountPlanDTO, businessId?: number): Promise<CycleCountPlan> {
        return this.request<CycleCountPlan>(`/inventory/cycle-count-plans/${id}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteCountPlan(id: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/inventory/cycle-count-plans/${id}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async listCountTasks(params: GetCountTasksParams = {}): Promise<CycleCountTaskListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.warehouse_id) qs.set('warehouse_id', String(params.warehouse_id));
        if (params.plan_id) qs.set('plan_id', String(params.plan_id));
        if (params.status) qs.set('status', params.status);
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<CycleCountTaskListResponse>(`/inventory/cycle-count-tasks${suffix}`);
    }

    async generateCountTask(data: GenerateCountTaskInput, businessId?: number): Promise<GenerateCountTaskResult> {
        return this.request<GenerateCountTaskResult>(`/inventory/cycle-count-tasks/generate${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async startCountTask(id: number, data: StartCountTaskInput, businessId?: number): Promise<CycleCountTask> {
        return this.request<CycleCountTask>(`/inventory/cycle-count-tasks/${id}/start${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async finishCountTask(id: number, businessId?: number): Promise<CycleCountTask> {
        return this.request<CycleCountTask>(`/inventory/cycle-count-tasks/${id}/finish${this.businessQuery(businessId)}`, {
            method: 'POST',
        });
    }

    async listCountLines(params: GetCountLinesParams): Promise<CycleCountLineListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.status) qs.set('status', params.status);
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<CycleCountLineListResponse>(`/inventory/cycle-count-tasks/${params.task_id}/lines${suffix}`);
    }

    async submitCountLine(id: number, data: SubmitCountLineInput, businessId?: number): Promise<SubmitCountLineResult> {
        return this.request<SubmitCountLineResult>(`/inventory/cycle-count-lines/${id}/submit${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async listDiscrepancies(params: GetDiscrepanciesParams = {}): Promise<InventoryDiscrepancyListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.task_id) qs.set('task_id', String(params.task_id));
        if (params.status) qs.set('status', params.status);
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<InventoryDiscrepancyListResponse>(`/inventory/discrepancies${suffix}`);
    }

    async approveDiscrepancy(id: number, data: ApproveDiscrepancyInput, businessId?: number): Promise<InventoryDiscrepancy> {
        return this.request<InventoryDiscrepancy>(`/inventory/discrepancies/${id}/approve${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async rejectDiscrepancy(id: number, data: RejectDiscrepancyInput, businessId?: number): Promise<InventoryDiscrepancy> {
        return this.request<InventoryDiscrepancy>(`/inventory/discrepancies/${id}/reject${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async exportKardex(data: KardexQueryInput): Promise<KardexExportResult> {
        const qs = new URLSearchParams();
        qs.set('product_id', data.product_id);
        qs.set('warehouse_id', String(data.warehouse_id));
        if (data.from) qs.set('from', data.from);
        if (data.to) qs.set('to', data.to);
        if (data.business_id) qs.set('business_id', String(data.business_id));
        return this.request<KardexExportResult>(`/inventory/kardex/export?${qs.toString()}`);
    }
}
