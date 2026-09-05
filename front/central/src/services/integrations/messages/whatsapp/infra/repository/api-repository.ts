import { env } from '@/shared/config/env';
import { IWhatsAppRepository } from '../../domain/ports';
import { WhatsAppProvisionResponse, WhatsAppTemplatesResponse } from '../../domain/types';

export class WhatsAppApiRepository implements IWhatsAppRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async request<T>(path: string, options: RequestInit = {}): Promise<T> {
        const response = await fetch(`${this.baseUrl}${path}`, {
            ...options,
            headers: {
                Accept: 'application/json',
                'Content-Type': 'application/json',
                ...(this.token ? { Authorization: `Bearer ${this.token}` } : {}),
                ...(options.headers || {}),
            },
            cache: 'no-store',
        });

        const body = await response.json().catch(() => ({}));

        if (!response.ok) {
            throw new Error(body?.error || body?.message || 'Error consultando las plantillas de WhatsApp');
        }

        return body as T;
    }

    private query(businessId?: number, refresh?: boolean): string {
        const params = new URLSearchParams();
        if (businessId) params.set('business_id', String(businessId));
        if (refresh) params.set('refresh', 'true');
        const qs = params.toString();
        return qs ? `?${qs}` : '';
    }

    async getTemplatesStatus(businessId?: number, refresh?: boolean): Promise<WhatsAppTemplatesResponse> {
        return this.request<WhatsAppTemplatesResponse>(
            `/integrations/whatsapp/templates/status${this.query(businessId, refresh)}`
        );
    }

    async provisionTemplates(businessId?: number): Promise<WhatsAppProvisionResponse> {
        return this.request<WhatsAppProvisionResponse>(
            `/integrations/whatsapp/templates/provision${this.query(businessId)}`,
            { method: 'POST' }
        );
    }
}
