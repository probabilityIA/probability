import type { IBackfillRepository } from '../../domain/ports';
import type { BackfillEvent, JobState, PreviewRequest, PreviewResponse, RunRequest, RunResponse } from '../../domain/types';

export class BackfillApiRepository implements IBackfillRepository {
    private baseUrl: string;
    private token: string;

    constructor(baseUrl: string, token: string) {
        this.baseUrl = baseUrl;
        this.token = token;
    }

    private headers() {
        return {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${this.token}`,
        };
    }

    async listEvents(): Promise<BackfillEvent[]> {
        const res = await fetch(`${this.baseUrl}/notification-backfill/events`, {
            headers: this.headers(),
            cache: 'no-store',
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err.error || 'Failed to list backfill events');
        }
        const body = await res.json();
        return body.events ?? [];
    }

    async preview(req: PreviewRequest): Promise<PreviewResponse> {
        const res = await fetch(`${this.baseUrl}/notification-backfill/preview`, {
            method: 'POST',
            headers: this.headers(),
            body: JSON.stringify(req),
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err.error || 'Preview failed');
        }
        return res.json();
    }

    async run(req: RunRequest): Promise<RunResponse> {
        const res = await fetch(`${this.baseUrl}/notification-backfill/run`, {
            method: 'POST',
            headers: this.headers(),
            body: JSON.stringify(req),
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err.error || 'Run failed');
        }
        return res.json();
    }

    async getJob(jobId: string): Promise<JobState> {
        const res = await fetch(`${this.baseUrl}/notification-backfill/jobs/${jobId}`, {
            headers: this.headers(),
            cache: 'no-store',
        });
        if (!res.ok) {
            const err = await res.json().catch(() => ({}));
            throw new Error(err.error || 'Failed to get job');
        }
        return res.json();
    }
}
