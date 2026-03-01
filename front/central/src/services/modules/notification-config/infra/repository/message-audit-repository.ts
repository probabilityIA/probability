import type { MessageAuditFilter, PaginatedMessageAuditResponse, MessageAuditStats } from "../../domain/types";

export class MessageAuditApiRepository {
    private baseUrl: string;
    private token: string;

    constructor(baseUrl: string, token: string) {
        this.baseUrl = baseUrl;
        this.token = token;
    }

    async list(filter: MessageAuditFilter): Promise<PaginatedMessageAuditResponse> {
        const params = new URLSearchParams();
        params.append("business_id", filter.business_id.toString());
        if (filter.status) params.append("status", filter.status);
        if (filter.direction) params.append("direction", filter.direction);
        if (filter.template_name) params.append("template_name", filter.template_name);
        if (filter.date_from) params.append("date_from", filter.date_from);
        if (filter.date_to) params.append("date_to", filter.date_to);
        if (filter.page) params.append("page", filter.page.toString());
        if (filter.page_size) params.append("page_size", filter.page_size.toString());

        const response = await fetch(
            `${this.baseUrl}/notification-configs/message-audit?${params.toString()}`,
            {
                headers: { Authorization: `Bearer ${this.token}` },
                cache: "no-store",
            }
        );

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to list message audit logs");
        }

        return response.json();
    }

    async getStats(businessId: number, dateFrom?: string, dateTo?: string): Promise<MessageAuditStats> {
        const params = new URLSearchParams();
        params.append("business_id", businessId.toString());
        if (dateFrom) params.append("date_from", dateFrom);
        if (dateTo) params.append("date_to", dateTo);

        const response = await fetch(
            `${this.baseUrl}/notification-configs/message-audit/stats?${params.toString()}`,
            {
                headers: { Authorization: `Bearer ${this.token}` },
                cache: "no-store",
            }
        );

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || "Failed to get message audit stats");
        }

        return response.json();
    }
}
