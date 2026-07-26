
export interface NotificationType {
    id: number;
    name: string;
    code: string;
    description?: string;
    icon?: string;
    is_active: boolean;
    config_schema?: Record<string, any>;
    created_at: string;
    updated_at: string;
}

export interface CreateNotificationTypeDTO {
    name: string;
    code: string;
    description?: string;
    icon?: string;
    is_active?: boolean;
    config_schema?: Record<string, any>;
}

export interface UpdateNotificationTypeDTO {
    name?: string;
    description?: string;
    icon?: string;
    is_active?: boolean;
    config_schema?: Record<string, any>;
}

export interface NotificationEventType {
    id: number;
    notification_type_id: number;
    event_code: string;
    event_name: string;
    description?: string;
    template_config?: Record<string, any>;
    is_active: boolean;
    allowed_order_status_ids?: number[];
    created_at: string;
    updated_at: string;

    notification_type?: NotificationType;
}

export interface CreateNotificationEventTypeDTO {
    notification_type_id: number;
    event_code: string;
    event_name: string;
    description?: string;
    template_config?: Record<string, any>;
    is_active?: boolean;
    allowed_order_status_ids?: number[];
}

export interface UpdateNotificationEventTypeDTO {
    event_name?: string;
    description?: string;
    template_config?: Record<string, any>;
    is_active?: boolean;
    allowed_order_status_ids?: number[];
}

export interface NotificationConfig {
    id: number;
    business_id: number;
    integration_id: number;
    notification_type_id: number;
    notification_event_type_id: number;
    enabled: boolean;
    filters?: Record<string, any>;
    description?: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;

    notification_type?: NotificationType;
    notification_event_type?: NotificationEventType;
    order_status_ids?: number[];

    notification_type_name?: string;
    notification_event_name?: string;

    event_type?: string;
    channels?: string[];
}

export interface CreateConfigDTO {
    business_id: number;
    integration_id: number;
    notification_type_id: number;
    notification_event_type_id: number;
    enabled?: boolean;
    filters?: Record<string, any>;
    description?: string;
    order_status_ids?: number[];
}

export interface UpdateConfigDTO {
    integration_id?: number;
    notification_type_id?: number;
    notification_event_type_id?: number;
    enabled?: boolean;
    filters?: Record<string, any>;
    description?: string;
    order_status_ids?: number[];
}

export interface ConfigFilter {
    business_id?: number;
    integration_id?: number;
    notification_type_id?: number;
    notification_event_type_id?: number;
}

export interface OrderStatus {
    id: number;
    code: string;
    name: string;
    description?: string;
    category: string;
    is_active: boolean;
    icon?: string;
    color?: string;
    created_at: string;
    updated_at: string;
}

export interface Integration {
    id: number;
    name: string;
    code: string;
    type: string;
    business_id?: number;
    is_active: boolean;
    integration_type_id?: number;
    integration_type_name?: string;
    integration_type_icon?: string;
    created_at: string;
    updated_at: string;
}

export interface SyncRule {
    id?: number;
    notification_type_id: number;
    notification_event_type_id: number;
    enabled: boolean;
    description: string;
    order_status_ids: number[];
}

export interface SyncConfigsDTO {
    integration_id: number;
    rules: SyncRule[];
}

export interface SyncConfigsResponse {
    created: number;
    updated: number;
    deleted: number;
    configs: NotificationConfig[];
}

export interface MessageAuditLog {
    id: string;
    conversation_id: string;
    message_id: string;
    direction: 'outbound' | 'inbound';
    template_name: string;
    content: string;
    status: 'sent' | 'delivered' | 'read' | 'failed';
    delivered_at?: string;
    read_at?: string;
    created_at: string;
    phone_number: string;
    order_number: string;
    business_id: number;
}

export interface MessageAuditStats {
    total_sent: number;
    total_delivered: number;
    total_read: number;
    total_failed: number;
    success_rate: number;
}

export interface MessageAuditFilter {
    business_id: number;
    status?: string;
    direction?: string;
    template_name?: string;
    date_from?: string;
    date_to?: string;
    page?: number;
    page_size?: number;
}

export interface PaginatedMessageAuditResponse {
    data: MessageAuditLog[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface ConversationSummary {
    id: string;
    phone_number: string;
    order_number: string;
    current_state: string;
    message_count: number;
    last_message_content: string;
    last_message_direction: string;
    last_message_status: string;
    last_activity: string;
    created_at: string;
}

export interface ConversationListFilter {
    business_id: number;
    state?: string;
    phone?: string;
    date_from?: string;
    date_to?: string;
    page?: number;
    page_size?: number;
}

export interface PaginatedConversationListResponse {
    data: ConversationSummary[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface ConversationMessage {
    id: string;
    direction: string;
    message_id: string;
    template_name: string;
    content: string;
    status: string;
    delivered_at?: string;
    read_at?: string;
    created_at: string;
}

export interface ConversationDetailResponse {
    conversation_id: string;
    phone_number: string;
    order_number: string;
    current_state: string;
    messages: ConversationMessage[];
}
