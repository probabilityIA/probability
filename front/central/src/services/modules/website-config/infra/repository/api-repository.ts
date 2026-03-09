import { env } from '@/shared/config/env';
import { IWebsiteConfigRepository } from '../../domain/ports';
import { WebsiteConfigData, UpdateWebsiteConfigDTO } from '../../domain/types';

export class WebsiteConfigApiRepository implements IWebsiteConfigRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(url, { ...options, headers });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.message || data.error || 'Error en la solicitud');
        }

        return data;
    }

    private withBusinessId(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
    }

    async getConfig(businessId?: number): Promise<WebsiteConfigData> {
        return this.fetch<WebsiteConfigData>(this.withBusinessId('/website-config', businessId));
    }

    async updateConfig(data: UpdateWebsiteConfigDTO, businessId?: number): Promise<WebsiteConfigData> {
        return this.fetch<WebsiteConfigData>(this.withBusinessId('/website-config', businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }
}
