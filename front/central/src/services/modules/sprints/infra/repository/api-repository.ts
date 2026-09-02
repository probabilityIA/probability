import { env } from '@/shared/config/env';
import { ISprintRepository } from '../../domain/ports';
import {
    Sprint,
    PaginatedSprints,
    ListSprintsParams,
    CreateSprintDTO,
    UpdateSprintDTO,
    SprintStatus,
} from '../../domain/types';

export class SprintApiRepository implements ISprintRepository {
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

    private withQuery(path: string, params?: Record<string, unknown>): string {
        if (!params) return path;
        const sp = new URLSearchParams();
        Object.entries(params).forEach(([k, v]) => {
            if (v === undefined || v === null || v === '') return;
            sp.append(k, String(v));
        });
        const q = sp.toString();
        return q ? `${path}?${q}` : path;
    }

    list(params?: ListSprintsParams): Promise<PaginatedSprints> {
        return this.request<PaginatedSprints>(this.withQuery('/sprints', params as Record<string, unknown>), {
            method: 'GET',
            headers: this.headers(),
        });
    }

    get(id: number): Promise<Sprint> {
        return this.request<Sprint>(`/sprints/${id}`, {
            method: 'GET',
            headers: this.headers(),
        });
    }

    create(data: CreateSprintDTO): Promise<Sprint> {
        return this.request<Sprint>('/sprints', {
            method: 'POST',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(data),
        });
    }

    update(id: number, data: UpdateSprintDTO): Promise<Sprint> {
        return this.request<Sprint>(`/sprints/${id}`, {
            method: 'PUT',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify(data),
        });
    }

    async remove(id: number): Promise<void> {
        await this.request<unknown>(`/sprints/${id}`, { method: 'DELETE', headers: this.headers() });
    }

    changeStatus(id: number, status: SprintStatus): Promise<Sprint> {
        return this.request<Sprint>(`/sprints/${id}/status`, {
            method: 'PATCH',
            headers: this.headers({ 'Content-Type': 'application/json' }),
            body: JSON.stringify({ status }),
        });
    }
}
