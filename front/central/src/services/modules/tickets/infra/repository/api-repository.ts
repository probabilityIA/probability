import { env } from '@/shared/config/env';
import { ITicketRepository } from '../../domain/ports';
import {
    Ticket,
    TicketComment,
    TicketAttachment,
    TicketHistoryEntry,
    PaginatedTickets,
    ListTicketsParams,
    CreateTicketDTO,
    UpdateTicketDTO,
} from '../../domain/types';

export class TicketApiRepository implements ITicketRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private headers(extra: Record<string, string> = {}): Record<string, string> {
        const h: Record<string, string> = { Accept: 'application/json', ...extra };
        if (this.token) h['Authorization'] = `Bearer ${this.token}`;
        return h;
    }

    private async request<T>(path: string, init: RequestInit = {}): Promise<T> {
        const res = await fetch(`${this.baseUrl}${path}`, { ...init, cache: 'no-store' });
        const text = await res.text();
        const data = text ? JSON.parse(text) : null;
        if (!res.ok) throw new Error(data?.error || data?.message || `HTTP ${res.status}`);
        return data as T;
    }

    private withQuery(path: string, params?: Record<string, any>): string {
        if (!params) return path;
        const sp = new URLSearchParams();
        Object.entries(params).forEach(([k, v]) => {
            if (v === undefined || v === null || v === '') return;
            sp.append(k, String(v));
        });
        const q = sp.toString();
        return q ? `${path}?${q}` : path;
    }

    list(params?: ListTicketsParams): Promise<PaginatedTickets> {
        return this.request<PaginatedTickets>(this.withQuery('/tickets', params), {
            method: 'GET',
            headers: this.headers(),
        });
    }

    get(id: number, businessId?: number): Promise<Ticket> {
        return this.request<Ticket>(this.withQuery(`/tickets/${id}`, { business_id: businessId }), {
            method: 'GET',
            headers: this.headers(),
        });
    }

    create(data: CreateTicketDTO): Promise<Ticket> {
        return this.request<Ticket>('/tickets', {
            method: 'POST',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(data),
        });
    }

    update(id: number, data: UpdateTicketDTO): Promise<Ticket> {
        return this.request<Ticket>(`/tickets/${id}`, {
            method: 'PUT',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(data),
        });
    }

    async remove(id: number): Promise<void> {
        await this.request<unknown>(`/tickets/${id}`, { method: 'DELETE', headers: this.headers() });
    }

    changeStatus(id: number, status: string, note?: string): Promise<Ticket> {
        return this.request<Ticket>(`/tickets/${id}/status`, {
            method: 'PATCH',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify({ status, note: note ?? '' }),
        });
    }

    assign(id: number, assignedToId: number | null): Promise<Ticket> {
        return this.request<Ticket>(`/tickets/${id}/assign`, {
            method: 'PATCH',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify({ assigned_to_id: assignedToId }),
        });
    }

    escalate(id: number, note?: string): Promise<Ticket> {
        return this.request<Ticket>(`/tickets/${id}/escalate`, {
            method: 'PATCH',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify({ note: note ?? '' }),
        });
    }

    async listComments(id: number, businessId?: number): Promise<TicketComment[]> {
        const r = await this.request<{ data: TicketComment[] }>(this.withQuery(`/tickets/${id}/comments`, { business_id: businessId }), {
            method: 'GET',
            headers: this.headers(),
        });
        return r.data || [];
    }

    addComment(id: number, body: string, isInternal: boolean): Promise<TicketComment> {
        return this.request<TicketComment>(`/tickets/${id}/comments`, {
            method: 'POST',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify({ body, is_internal: isInternal }),
        });
    }

    async listAttachments(id: number, businessId?: number): Promise<TicketAttachment[]> {
        const r = await this.request<{ data: TicketAttachment[] }>(this.withQuery(`/tickets/${id}/attachments`, { business_id: businessId }), {
            method: 'GET',
            headers: this.headers(),
        });
        return r.data || [];
    }

    async uploadAttachment(id: number, file: File, commentId?: number): Promise<TicketAttachment> {
        const fd = new FormData();
        fd.append('file', file);
        if (commentId) fd.append('comment_id', String(commentId));
        return this.request<TicketAttachment>(`/tickets/${id}/attachments`, {
            method: 'POST',
            headers: this.headers(),
            body: fd,
        });
    }

    async deleteAttachment(attachmentId: number): Promise<void> {
        await this.request<unknown>(`/tickets/attachments/${attachmentId}`, { method: 'DELETE', headers: this.headers() });
    }

    async listHistory(id: number, businessId?: number): Promise<TicketHistoryEntry[]> {
        const r = await this.request<{ data: TicketHistoryEntry[] }>(this.withQuery(`/tickets/${id}/history`, { business_id: businessId }), {
            method: 'GET',
            headers: this.headers(),
        });
        return r.data || [];
    }
}
