import { env } from '@/shared/config/env';
import { IAccountingRepository } from '../../domain/ports';
import {
    AccountingReport,
    Concept,
    CreateConceptDTO,
    CreateEntryDTO,
    CreateTaxDTO,
    EntriesListResponse,
    Entry,
    GetEntriesParams,
    SyncResult,
    Tax,
    UpdateConceptDTO,
    UpdateTaxDTO,
} from '../../domain/types';

interface Envelope<T> {
    success: boolean;
    data: T;
    error?: string;
    message?: string;
    total?: number;
    page?: number;
    page_size?: number;
    total_pages?: number;
}

export class AccountingApiRepository implements IAccountingRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<Envelope<T>> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(url, { cache: 'no-store', ...options, headers });
        const data = await res.json();

        if (!res.ok || data?.success === false) {
            throw new Error(data?.error || data?.message || 'Ocurrio un error');
        }

        return data as Envelope<T>;
    }

    async getConcepts(): Promise<Concept[]> {
        const res = await this.fetch<Concept[]>('/accounting/concepts');
        return res.data || [];
    }

    async createConcept(data: CreateConceptDTO): Promise<Concept> {
        const res = await this.fetch<Concept>('/accounting/concepts', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async updateConcept(id: number, data: UpdateConceptDTO): Promise<Concept> {
        const res = await this.fetch<Concept>(`/accounting/concepts/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async setConceptTax(conceptId: number, taxId: number, isActive: boolean): Promise<void> {
        await this.fetch<unknown>(`/accounting/concepts/${conceptId}/taxes/${taxId}`, {
            method: 'PUT',
            body: JSON.stringify({ is_active: isActive }),
        });
    }

    async getTaxes(): Promise<Tax[]> {
        const res = await this.fetch<Tax[]>('/accounting/taxes');
        return res.data || [];
    }

    async createTax(data: CreateTaxDTO): Promise<Tax> {
        const res = await this.fetch<Tax>('/accounting/taxes', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async updateTax(id: number, data: UpdateTaxDTO): Promise<Tax> {
        const res = await this.fetch<Tax>(`/accounting/taxes/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async getEntries(params?: GetEntriesParams): Promise<EntriesListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        const res = await this.fetch<Entry[]>(`/accounting/entries${query ? `?${query}` : ''}`);
        return {
            data: res.data || [],
            total: res.total ?? 0,
            page: res.page ?? 1,
            page_size: res.page_size ?? 10,
            total_pages: res.total_pages ?? 0,
        };
    }

    async createEntry(data: CreateEntryDTO): Promise<Entry> {
        const res = await this.fetch<Entry>('/accounting/entries', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async deleteEntry(id: number): Promise<void> {
        await this.fetch<unknown>(`/accounting/entries/${id}`, {
            method: 'DELETE',
        });
    }

    async getReport(from: string, to: string): Promise<AccountingReport> {
        const query = new URLSearchParams({ from, to }).toString();
        const res = await this.fetch<AccountingReport>(`/accounting/report?${query}`);
        return res.data;
    }

    async syncNow(): Promise<SyncResult> {
        const res = await this.fetch<SyncResult>('/accounting/sync', {
            method: 'POST',
        });
        return res.data;
    }
}
