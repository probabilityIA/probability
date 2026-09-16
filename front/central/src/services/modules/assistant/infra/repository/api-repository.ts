import { env } from '@/shared/config/env';
import type { IAssistantRepository } from '../../domain/ports';
import type {
    AssistantAlert,
    AssistantAlertsUnread,
    AssistantErrorCode,
    AssistantHistoryMessage,
    AssistantReply,
    AssistantState,
    FeedbackValue,
    PaginatedResponse,
    ReviewFilters,
    ReviewMessage,
    ReviewSummary,
} from '../../domain/types';

export class AssistantApiError extends Error {
    constructor(
        message: string,
        readonly code: AssistantErrorCode,
        readonly resetAt: string | null = null,
    ) {
        super(message);
    }
}

function toQuery(filters: ReviewFilters): string {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') params.append(key, String(value));
    });
    const query = params.toString();
    return query ? `?${query}` : '';
}

export class AssistantApiRepository implements IAssistantRepository {
    constructor(private readonly token: string) {}

    private async requestBody(path: string, init: RequestInit = {}): Promise<Record<string, unknown>> {
        const res = await fetch(`${env.API_BASE_URL}${path}`, {
            ...init,
            headers: {
                'Content-Type': 'application/json',
                Authorization: `Bearer ${this.token}`,
                ...(init.headers || {}),
            },
            cache: 'no-store',
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok || !body?.success) {
            throw new AssistantApiError(
                body?.message || body?.error || `Error ${res.status}`,
                (body?.code as AssistantErrorCode) || (res.status >= 500 ? 'assistant_unavailable' : 'internal_error'),
                body?.reset_at ?? null,
            );
        }
        return body;
    }

    private async request<T>(path: string, init: RequestInit = {}): Promise<T> {
        const body = await this.requestBody(path, init);
        return body.data as T;
    }

    chat(messages: AssistantHistoryMessage[], conversationId: string, pathname: string, businessId?: number | null): Promise<AssistantReply> {
        const query = businessId ? `?business_id=${businessId}` : '';
        return this.request<AssistantReply>(`/ai/assistant/chat${query}`, {
            method: 'POST',
            body: JSON.stringify({ messages, conversation_id: conversationId, pathname }),
        });
    }

    getState(): Promise<AssistantState> {
        return this.request<AssistantState>('/ai/assistant/state', { method: 'GET' });
    }

    async markIntroSeen(): Promise<void> {
        await this.requestBody('/ai/assistant/intro-seen', { method: 'POST' });
    }

    async sendFeedback(messageId: string, value: FeedbackValue): Promise<void> {
        await this.requestBody(`/ai/assistant/messages/${encodeURIComponent(messageId)}/feedback`, {
            method: 'POST',
            body: JSON.stringify({ value }),
        });
    }

    async markClick(messageId: string): Promise<void> {
        await this.requestBody(`/ai/assistant/messages/${encodeURIComponent(messageId)}/click`, { method: 'POST' });
    }

    async listReviewMessages(filters: ReviewFilters): Promise<PaginatedResponse<ReviewMessage>> {
        const body = await this.requestBody(`/ai/assistant/admin/messages${toQuery(filters)}`, { method: 'GET' });
        return {
            data: (body.data as ReviewMessage[]) || [],
            total: Number(body.total) || 0,
            page: Number(body.page) || 1,
            page_size: Number(body.page_size) || 20,
            total_pages: Number(body.total_pages) || 1,
        };
    }

    getReviewSummary(filters: ReviewFilters): Promise<ReviewSummary> {
        return this.request<ReviewSummary>(`/ai/assistant/admin/summary${toQuery(filters)}`, { method: 'GET' });
    }

    private alertsQuery(businessId: number | null | undefined, extra: Record<string, string | number> = {}): string {
        const params = new URLSearchParams();
        if (businessId) params.append('business_id', String(businessId));
        Object.entries(extra).forEach(([key, value]) => params.append(key, String(value)));
        const query = params.toString();
        return query ? `?${query}` : '';
    }

    async listAlerts(businessId: number | null | undefined, page: number, pageSize: number): Promise<PaginatedResponse<AssistantAlert>> {
        const data = await this.request<PaginatedResponse<AssistantAlert>>(
            `/ai/assistant/alerts${this.alertsQuery(businessId, { page, page_size: pageSize })}`,
            { method: 'GET' },
        );
        return {
            data: data?.data || [],
            total: Number(data?.total) || 0,
            page: Number(data?.page) || page,
            page_size: Number(data?.page_size) || pageSize,
            total_pages: Number(data?.total_pages) || 1,
        };
    }

    getAlertsUnread(businessId: number | null | undefined): Promise<AssistantAlertsUnread> {
        return this.request<AssistantAlertsUnread>(`/ai/assistant/alerts/unread${this.alertsQuery(businessId)}`, { method: 'GET' });
    }

    async markAlertsSeen(businessId: number | null | undefined): Promise<void> {
        await this.requestBody(`/ai/assistant/alerts/seen${this.alertsQuery(businessId)}`, { method: 'POST' });
    }
}
