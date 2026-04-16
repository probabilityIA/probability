import { env } from '@/shared/config/env';
import { IAnnouncementRepository } from '../../domain/ports';
import {
    AnnouncementInfo,
    AnnouncementStats,
    AnnouncementCategory,
    AnnouncementsListResponse,
    GetAnnouncementsParams,
    CreateAnnouncementDTO,
    UpdateAnnouncementDTO,
    RegisterViewDTO,
    ChangeStatusDTO,
    DeleteAnnouncementResponse,
    UploadImageResponse,
    DeleteImageResponse,
} from '../../domain/types';

export class AnnouncementApiRepository implements IAnnouncementRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(url, { ...options, headers });
        const data = await res.json();

        if (!res.ok) {
            throw new Error(data.error || data.message || 'An error occurred');
        }

        return data;
    }

    async getAnnouncements(params?: GetAnnouncementsParams): Promise<AnnouncementsListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<AnnouncementsListResponse>(`/announcements${query ? `?${query}` : ''}`);
    }

    async getAnnouncementById(id: number): Promise<AnnouncementInfo> {
        const res = await this.fetch<{ data: AnnouncementInfo }>(`/announcements/${id}`);
        return res.data;
    }

    async createAnnouncement(data: CreateAnnouncementDTO): Promise<AnnouncementInfo> {
        const res = await this.fetch<{ data: AnnouncementInfo }>('/announcements', {
            method: 'POST',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async updateAnnouncement(id: number, data: UpdateAnnouncementDTO): Promise<AnnouncementInfo> {
        const res = await this.fetch<{ data: AnnouncementInfo }>(`/announcements/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
        return res.data;
    }

    async deleteAnnouncement(id: number): Promise<DeleteAnnouncementResponse> {
        return this.fetch<DeleteAnnouncementResponse>(`/announcements/${id}`, {
            method: 'DELETE',
        });
    }

    async getActiveAnnouncements(businessId?: number): Promise<AnnouncementInfo[]> {
        const query = businessId ? `?business_id=${businessId}` : '';
        const res = await this.fetch<{ data: AnnouncementInfo[] }>(`/announcements/active${query}`);
        return res.data || [];
    }

    async registerView(announcementId: number, data: RegisterViewDTO): Promise<void> {
        await this.fetch<void>(`/announcements/${announcementId}/view`, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async getStats(announcementId: number): Promise<AnnouncementStats> {
        const res = await this.fetch<{ data: AnnouncementStats }>(`/announcements/${announcementId}/stats`);
        return res.data;
    }

    async listCategories(): Promise<AnnouncementCategory[]> {
        const res = await this.fetch<{ data: AnnouncementCategory[] }>('/announcements/categories');
        return res.data || [];
    }

    async changeStatus(id: number, data: ChangeStatusDTO): Promise<void> {
        await this.fetch<void>(`/announcements/${id}/status`, {
            method: 'PATCH',
            body: JSON.stringify(data),
        });
    }

    async forceRedisplay(id: number): Promise<void> {
        await this.fetch<void>(`/announcements/${id}/force-redisplay`, {
            method: 'POST',
        });
    }

    async uploadImage(announcementId: number, formData: FormData): Promise<UploadImageResponse> {
        const url = `${this.baseUrl}/announcements/${announcementId}/image`;
        const headers: Record<string, string> = {
            'Accept': 'application/json',
        };
        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }
        const res = await fetch(url, {
            method: 'POST',
            headers,
            body: formData,
        });
        const data = await res.json();
        if (!res.ok) {
            throw new Error(data.message || data.error || 'Error uploading image');
        }
        return data;
    }

    async deleteImage(announcementId: number, imageId: number): Promise<DeleteImageResponse> {
        return this.fetch<DeleteImageResponse>(`/announcements/${announcementId}/image/${imageId}`, {
            method: 'DELETE',
        });
    }
}
