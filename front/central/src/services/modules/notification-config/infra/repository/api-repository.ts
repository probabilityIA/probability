import { INotificationConfigRepository } from "../../domain/ports";
import { NotificationConfig, CreateConfigDTO, UpdateConfigDTO, ConfigFilter } from "../../domain/types";

export class NotificationConfigApiRepository implements INotificationConfigRepository {
    private baseUrl: string;
    private token: string;

    constructor(baseUrl: string, token: string) {
        this.baseUrl = baseUrl;
        this.token = token;
    }

    async create(dto: CreateConfigDTO): Promise<NotificationConfig> {
        const response = await fetch(`${this.baseUrl}/notification-configs`, {
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

    async getById(id: number): Promise<NotificationConfig> {
        const response = await fetch(`${this.baseUrl}/notification-configs/${id}`, {
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

    async update(id: number, dto: UpdateConfigDTO): Promise<NotificationConfig> {
        const response = await fetch(`${this.baseUrl}/notification-configs/${id}`, {
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

    async delete(id: number): Promise<void> {
        const response = await fetch(`${this.baseUrl}/notification-configs/${id}`, {
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

    async list(filter?: ConfigFilter): Promise<NotificationConfig[]> {
        const params = new URLSearchParams();
        if (filter) {
            if (filter.business_id) params.append("business_id", filter.business_id.toString());
            if (filter.integration_id) params.append("integration_id", filter.integration_id.toString());
            if (filter.notification_type_id) params.append("notification_type_id", filter.notification_type_id.toString());
            if (filter.notification_event_type_id) params.append("notification_event_type_id", filter.notification_event_type_id.toString());
        }

        const response = await fetch(`${this.baseUrl}/notification-configs?${params.toString()}`, {
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
