import { INotificationConfigRepository } from "../../domain/ports";
import { NotificationConfig, CreateConfigDTO, UpdateConfigDTO, ConfigFilter, SyncConfigsDTO, SyncConfigsResponse } from "../../domain/types";

export class NotificationConfigApiRepository implements INotificationConfigRepository {
    private baseUrl: string;
    private token: string;

    constructor(baseUrl: string, token: string) {
        this.baseUrl = baseUrl;
        this.token = token;
    }

    private withBusinessId(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
    }

    async create(dto: CreateConfigDTO, businessId?: number): Promise<NotificationConfig> {
        const url = this.withBusinessId(`${this.baseUrl}/notification-configs`, businessId || dto.business_id);
        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${this.token}`,
            },
            body: JSON.stringify(dto),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to create notification config");
        }

        return response.json();
    }

    async getById(id: number, businessId?: number): Promise<NotificationConfig> {
        const url = this.withBusinessId(`${this.baseUrl}/notification-configs/${id}`, businessId);
        const response = await fetch(url, {
            headers: {
                Authorization: `Bearer ${this.token}`,
            },
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to get notification config");
        }

        return response.json();
    }

    async update(id: number, dto: UpdateConfigDTO, businessId?: number): Promise<NotificationConfig> {
        const url = this.withBusinessId(`${this.baseUrl}/notification-configs/${id}`, businessId);
        const response = await fetch(url, {
            method: "PATCH",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${this.token}`,
            },
            body: JSON.stringify(dto),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to update notification config");
        }

        return response.json();
    }

    async delete(id: number, businessId?: number): Promise<void> {
        const url = this.withBusinessId(`${this.baseUrl}/notification-configs/${id}`, businessId);
        const response = await fetch(url, {
            method: "DELETE",
            headers: {
                Authorization: `Bearer ${this.token}`,
            },
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to delete notification config");
        }
    }

    async syncByIntegration(dto: SyncConfigsDTO, businessId?: number): Promise<SyncConfigsResponse> {
        const url = this.withBusinessId(`${this.baseUrl}/notification-configs/sync`, businessId);
        const response = await fetch(url, {
            method: "PUT",
            headers: {
                "Content-Type": "application/json",
                Authorization: `Bearer ${this.token}`,
            },
            body: JSON.stringify(dto),
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to sync notification configs");
        }

        return response.json();
    }

    async list(filter?: ConfigFilter): Promise<NotificationConfig[]> {
        const params = new URLSearchParams();
        if (filter) {
            if (filter.business_id) params.append("business_id", filter.business_id.toString());
            if (filter.integration_id) params.append("integration_id", filter.integration_id.toString());
            if (filter.notification_type_id) params.append("notification_type_id", filter.notification_type_id.toString());
            if (filter.notification_event_type_id) params.append("notification_event_type_id", filter.notification_event_type_id.toString());
        }

        const url = `${this.baseUrl}/notification-configs?${params.toString()}`;

        const response = await fetch(url, {
            headers: {
                Authorization: `Bearer ${this.token}`,
            },
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to list notification configs");
        }

        return response.json();
    }
}
