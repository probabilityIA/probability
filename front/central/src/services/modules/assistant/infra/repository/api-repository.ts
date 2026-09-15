import { env } from '@/shared/config/env';
import type { IAssistantRepository } from '../../domain/ports';
import type {
    AssistantErrorCode,
    AssistantHistoryMessage,
    AssistantReply,
    AssistantState,
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

export class AssistantApiRepository implements IAssistantRepository {
    constructor(private readonly token: string) {}

    private async request<T>(path: string, init: RequestInit = {}): Promise<T> {
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
        return body.data as T;
    }

    chat(messages: AssistantHistoryMessage[]): Promise<AssistantReply> {
        return this.request<AssistantReply>('/ai/assistant/chat', {
            method: 'POST',
            body: JSON.stringify({ messages }),
        });
    }

    getState(): Promise<AssistantState> {
        return this.request<AssistantState>('/ai/assistant/state', { method: 'GET' });
    }

    async markIntroSeen(): Promise<void> {
        await this.request<unknown>('/ai/assistant/intro-seen', { method: 'POST' });
    }
}
