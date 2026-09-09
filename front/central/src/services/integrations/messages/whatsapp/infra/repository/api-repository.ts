import { env } from '@/shared/config/env';
import { IWhatsAppRepository } from '../../domain/ports';
import {
    WhatsAppAddNumberValues,
    WhatsAppConnectionResponse,
    WhatsAppConnectionValues,
    WhatsAppEmbeddedSignupConfigResponse,
    WhatsAppEmbeddedSignupPayload,
    WhatsAppEmbeddedSignupResponse,
    WhatsAppNumberResponse,
    WhatsAppBusinessProfileResponse,
    WhatsAppBusinessProfileValues,
    WhatsAppProvisionResponse,
    WhatsAppTemplatesResponse,
} from '../../domain/types';

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

    async saveConnection(
        values: WhatsAppConnectionValues,
        businessId?: number
    ): Promise<WhatsAppConnectionResponse> {
        return this.request<WhatsAppConnectionResponse>(
            `/integrations/whatsapp/connection${this.query(businessId)}`,
            {
                method: 'PUT',
                body: JSON.stringify({
                    use_platform_token: values.use_platform_token,
                    waba_id: values.waba_id,
                    phone_number_id: values.phone_number_id,
                    access_token: values.access_token,
                }),
            }
        );
    }

    async getNumberState(businessId?: number): Promise<WhatsAppNumberResponse> {
        return this.request<WhatsAppNumberResponse>(`/integrations/whatsapp/numbers${this.query(businessId)}`);
    }

    async getBusinessProfile(businessId?: number): Promise<WhatsAppBusinessProfileResponse> {
        return this.request<WhatsAppBusinessProfileResponse>(
            `/integrations/whatsapp/profile${this.query(businessId)}`
        );
    }

    async updateBusinessProfile(
        values: WhatsAppBusinessProfileValues,
        businessId?: number
    ): Promise<WhatsAppBusinessProfileResponse> {
        return this.request<WhatsAppBusinessProfileResponse>(
            `/integrations/whatsapp/profile${this.query(businessId)}`,
            { method: 'PUT', body: JSON.stringify(values) }
        );
    }

    async updateBusinessProfilePhoto(
        file: File,
        businessId?: number
    ): Promise<WhatsAppBusinessProfileResponse> {
        const form = new FormData();
        form.append('photo', file);

        const response = await fetch(
            `${this.baseUrl}/integrations/whatsapp/profile/photo${this.query(businessId)}`,
            {
                method: 'POST',
                headers: {
                    Accept: 'application/json',
                    ...(this.token ? { Authorization: `Bearer ${this.token}` } : {}),
                },
                body: form,
                cache: 'no-store',
            }
        );

        const body = await response.json().catch(() => ({}));

        if (!response.ok) {
            throw new Error(body?.error || body?.message || 'No se pudo actualizar la foto');
        }

        return body as WhatsAppBusinessProfileResponse;
    }

    async addNumber(values: WhatsAppAddNumberValues, businessId?: number): Promise<WhatsAppNumberResponse> {
        return this.request<WhatsAppNumberResponse>(`/integrations/whatsapp/numbers${this.query(businessId)}`, {
            method: 'POST',
            body: JSON.stringify(values),
        });
    }

    async requestNumberCode(method: string, businessId?: number): Promise<WhatsAppNumberResponse> {
        return this.request<WhatsAppNumberResponse>(`/integrations/whatsapp/numbers/code${this.query(businessId)}`, {
            method: 'POST',
            body: JSON.stringify({ method }),
        });
    }

    async verifyNumberCode(code: string, businessId?: number): Promise<WhatsAppNumberResponse> {
        return this.request<WhatsAppNumberResponse>(`/integrations/whatsapp/numbers/verify${this.query(businessId)}`, {
            method: 'POST',
            body: JSON.stringify({ code }),
        });
    }

    async registerNumber(businessId?: number): Promise<WhatsAppNumberResponse> {
        return this.request<WhatsAppNumberResponse>(`/integrations/whatsapp/numbers/register${this.query(businessId)}`, {
            method: 'POST',
        });
    }

    async getEmbeddedSignupConfig(): Promise<WhatsAppEmbeddedSignupConfigResponse> {
        return this.request<WhatsAppEmbeddedSignupConfigResponse>(
            '/integrations/whatsapp/embedded-signup/config'
        );
    }

    async completeEmbeddedSignup(
        payload: WhatsAppEmbeddedSignupPayload,
        businessId?: number
    ): Promise<WhatsAppEmbeddedSignupResponse> {
        return this.request<WhatsAppEmbeddedSignupResponse>(
            `/integrations/whatsapp/embedded-signup${this.query(businessId)}`,
            { method: 'POST', body: JSON.stringify(payload) }
        );
    }
}
