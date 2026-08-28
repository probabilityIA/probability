import { env } from '@/shared/config/env';
import { ILegalRepository } from '../../domain/ports';
import { AcceptLegalResult, PendingLegalDocuments } from '../../domain/types';

export class LegalApiRepository implements ILegalRepository {
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
            ...((options.headers as Record<string, string>) || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(`${this.baseUrl}${path}`, { ...options, headers, cache: 'no-store' });
        const data = await res.json();

        if (!res.ok || data?.success === false) {
            throw new Error(data?.message || data?.error || 'Error al consultar los documentos legales');
        }

        return data.data as T;
    }

    async getPendingDocuments(): Promise<PendingLegalDocuments> {
        return this.request<PendingLegalDocuments>('/legal/pending');
    }

    async acceptDocuments(documentIds: number[]): Promise<AcceptLegalResult> {
        return this.request<AcceptLegalResult>('/legal/accept', {
            method: 'POST',
            body: JSON.stringify({ document_ids: documentIds }),
        });
    }
}
