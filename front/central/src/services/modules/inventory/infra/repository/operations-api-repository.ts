import { env } from '@/shared/config/env';
import {
    AssignReplenishmentInput,
    CompleteReplenishmentInput,
    ConfirmPutawayInput,
    CreateCrossDockLinkDTO,
    CreatePutawayRuleDTO,
    CreateReplenishmentTaskDTO,
    CrossDockLink,
    CrossDockLinkListResponse,
    GetCrossDockLinksParams,
    GetPutawayRulesParams,
    GetPutawaySuggestionsParams,
    GetReplenishmentTasksParams,
    GetVelocitiesParams,
    ProductVelocity,
    PutawayRule,
    PutawayRuleListResponse,
    PutawaySuggestResult,
    PutawaySuggestion,
    PutawaySuggestionListResponse,
    ReplenishmentDetectResult,
    ReplenishmentTask,
    ReplenishmentTaskListResponse,
    RunSlottingInput,
    SlottingRunResult,
    SuggestPutawayInput,
    UpdatePutawayRuleDTO,
} from '../../domain/operations-types';

export class OperationsApiRepository {
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

    async listPutawayRules(params: GetPutawayRulesParams = {}): Promise<PutawayRuleListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.active_only) qs.set('active_only', 'true');
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<PutawayRuleListResponse>(`/inventory/putaway-rules${suffix}`);
    }

    async createPutawayRule(data: CreatePutawayRuleDTO, businessId?: number): Promise<PutawayRule> {
        return this.request<PutawayRule>(`/inventory/putaway-rules${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updatePutawayRule(id: number, data: UpdatePutawayRuleDTO, businessId?: number): Promise<PutawayRule> {
        return this.request<PutawayRule>(`/inventory/putaway-rules/${id}${this.businessQuery(businessId)}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deletePutawayRule(id: number, businessId?: number): Promise<void> {
        await this.request<{ message: string }>(`/inventory/putaway-rules/${id}${this.businessQuery(businessId)}`, { method: 'DELETE' });
    }

    async suggestPutaway(data: SuggestPutawayInput, businessId?: number): Promise<PutawaySuggestResult> {
        return this.request<PutawaySuggestResult>(`/inventory/putaway/suggest${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async confirmPutaway(id: number, data: ConfirmPutawayInput, businessId?: number): Promise<PutawaySuggestion> {
        return this.request<PutawaySuggestion>(`/inventory/putaway/suggestions/${id}/confirm${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async listPutawaySuggestions(params: GetPutawaySuggestionsParams = {}): Promise<PutawaySuggestionListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.status) qs.set('status', params.status);
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<PutawaySuggestionListResponse>(`/inventory/putaway/suggestions${suffix}`);
    }

    async listReplenishmentTasks(params: GetReplenishmentTasksParams = {}): Promise<ReplenishmentTaskListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.warehouse_id) qs.set('warehouse_id', String(params.warehouse_id));
        if (params.status) qs.set('status', params.status);
        if (params.assigned_to) qs.set('assigned_to', String(params.assigned_to));
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<ReplenishmentTaskListResponse>(`/inventory/replenishment/tasks${suffix}`);
    }

    async createReplenishmentTask(data: CreateReplenishmentTaskDTO, businessId?: number): Promise<ReplenishmentTask> {
        return this.request<ReplenishmentTask>(`/inventory/replenishment/tasks${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async assignReplenishment(id: number, data: AssignReplenishmentInput, businessId?: number): Promise<ReplenishmentTask> {
        return this.request<ReplenishmentTask>(`/inventory/replenishment/tasks/${id}/assign${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async completeReplenishment(id: number, data: CompleteReplenishmentInput, businessId?: number): Promise<ReplenishmentTask> {
        return this.request<ReplenishmentTask>(`/inventory/replenishment/tasks/${id}/complete${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async cancelReplenishment(id: number, reason: string, businessId?: number): Promise<ReplenishmentTask> {
        return this.request<ReplenishmentTask>(`/inventory/replenishment/tasks/${id}/cancel${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify({ reason }),
        });
    }

    async detectReplenishment(businessId?: number): Promise<ReplenishmentDetectResult> {
        return this.request<ReplenishmentDetectResult>(`/inventory/replenishment/detect${this.businessQuery(businessId)}`, {
            method: 'POST',
        });
    }

    async listCrossDockLinks(params: GetCrossDockLinksParams = {}): Promise<CrossDockLinkListResponse> {
        const qs = new URLSearchParams();
        if (params.page) qs.set('page', String(params.page));
        if (params.page_size) qs.set('page_size', String(params.page_size));
        if (params.outbound_order_id) qs.set('outbound_order_id', params.outbound_order_id);
        if (params.status) qs.set('status', params.status);
        if (params.business_id) qs.set('business_id', String(params.business_id));
        const suffix = qs.toString() ? `?${qs.toString()}` : '';
        return this.request<CrossDockLinkListResponse>(`/inventory/cross-dock/links${suffix}`);
    }

    async createCrossDockLink(data: CreateCrossDockLinkDTO, businessId?: number): Promise<CrossDockLink> {
        return this.request<CrossDockLink>(`/inventory/cross-dock/links${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async executeCrossDock(id: number, businessId?: number): Promise<CrossDockLink> {
        return this.request<CrossDockLink>(`/inventory/cross-dock/links/${id}/execute${this.businessQuery(businessId)}`, {
            method: 'POST',
        });
    }

    async runSlotting(data: RunSlottingInput, businessId?: number): Promise<SlottingRunResult> {
        return this.request<SlottingRunResult>(`/inventory/slotting/run${this.businessQuery(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async listVelocities(params: GetVelocitiesParams): Promise<{ data: ProductVelocity[] }> {
        const qs = new URLSearchParams();
        qs.set('warehouse_id', String(params.warehouse_id));
        if (params.period) qs.set('period', params.period);
        if (params.rank) qs.set('rank', params.rank);
        if (params.limit) qs.set('limit', String(params.limit));
        if (params.business_id) qs.set('business_id', String(params.business_id));
        return this.request<{ data: ProductVelocity[] }>(`/inventory/slotting/velocities?${qs.toString()}`);
    }
}
