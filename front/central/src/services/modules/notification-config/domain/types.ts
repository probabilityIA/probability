export interface NotificationConfig {
    id: number;
    created_at: string;
    updated_at: string;
    deleted_at?: string;

    business_id: number;
    event_type: string;
    enabled: boolean;
    channels: string[];
    filters: Record<string, any>;
    description: string;
}

export interface CreateConfigDTO {
    business_id: number;
    event_type: string;
    enabled?: boolean;
    channels?: string[];
    filters?: Record<string, any>;
    description?: string;
}

export interface UpdateConfigDTO {
    enabled?: boolean;
    channels?: string[];
    filters?: Record<string, any>;
    description?: string;
}

export interface ConfigFilter {
    business_id?: number;
    event_type?: string;
}
