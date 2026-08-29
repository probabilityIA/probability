import { env } from '@/shared/config/env';
import { IAccountingRepository } from '../../domain/ports';
import {
    AccountingReport,
    AccountingService,
    ClientProfile,
    Concept,
    CreateConceptDTO,
    CreateEntryDTO,
    CreateInvoiceDTO,
    CreateTaxDTO,
    DianConfig,
    EmitDianDTO,
    EntriesListResponse,
    Entry,
    GetEntriesParams,
    GetInvoicesParams,
    Invoice,
    InvoicesListResponse,
    SaveServiceDTO,
    SyncResult,
    Tax,
    UpdateConceptDTO,
    UpdateDianConfigDTO,
    UpdateInvoiceDTO,
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
            throw new Error(data?.error || data?.message || 'Ocurrió un error');
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

    async getServices(): Promise<AccountingService[]> {
        const res = await this.fetch<AccountingService[]>('/accounting/services');
        return res.data || [];
    }

    async createService(data: SaveServiceDTO): Promise<AccountingService> {
        const res = await this.fetch<AccountingService>('/accounting/services', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async updateService(id: number, data: SaveServiceDTO): Promise<AccountingService> {
        const res = await this.fetch<AccountingService>(`/accounting/services/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async deleteService(id: number): Promise<void> {
        await this.fetch<unknown>(`/accounting/services/${id}`, {
            method: 'DELETE',
        });
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

    async getInvoices(params?: GetInvoicesParams): Promise<InvoicesListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        const res = await this.fetch<Invoice[]>(`/accounting/invoices${query ? `?${query}` : ''}`);
        return {
            data: res.data || [],
            total: res.total ?? 0,
            page: res.page ?? 1,
            page_size: res.page_size ?? 10,
            total_pages: res.total_pages ?? 0,
        };
    }

    async getInvoice(id: number): Promise<Invoice> {
        const res = await this.fetch<Invoice>(`/accounting/invoices/${id}`);
        return res.data;
    }

    async createInvoice(data: CreateInvoiceDTO): Promise<Invoice> {
        const res = await this.fetch<Invoice>('/accounting/invoices', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async updateInvoice(id: number, data: UpdateInvoiceDTO): Promise<Invoice> {
        const res = await this.fetch<Invoice>(`/accounting/invoices/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async deleteInvoice(id: number): Promise<void> {
        await this.fetch<unknown>(`/accounting/invoices/${id}`, {
            method: 'DELETE',
        });
    }

    async sendInvoice(id: number, emailTo?: string): Promise<Invoice> {
        const res = await this.fetch<Invoice>(`/accounting/invoices/${id}/send`, {
            method: 'POST',
            body: JSON.stringify({ email_to: emailTo || '' }),
        });
        return res.data;
    }

    async payInvoice(id: number): Promise<Invoice> {
        const res = await this.fetch<Invoice>(`/accounting/invoices/${id}/pay`, {
            method: 'POST',
        });
        return res.data;
    }

    async cancelInvoice(id: number): Promise<Invoice> {
        const res = await this.fetch<Invoice>(`/accounting/invoices/${id}/cancel`, {
            method: 'POST',
        });
        return res.data;
    }

    async getDianConfig(): Promise<DianConfig> {
        const res = await this.fetch<DianConfig>('/accounting/dian-config');
        return res.data;
    }

    async updateDianConfig(data: UpdateDianConfigDTO): Promise<DianConfig> {
        const res = await this.fetch<DianConfig>('/accounting/dian-config', {
            method: 'PUT',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async emitInvoiceDian(id: number, data: EmitDianDTO): Promise<Invoice> {
        const res = await this.fetch<Invoice>(`/accounting/invoices/${id}/emit-dian`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async getClientProfile(businessId: number): Promise<ClientProfile> {
        const res = await this.fetch<ClientProfile>(`/accounting/client-profile?business_id=${businessId}`);
        return res.data;
    }
}
