// ===================================================
// NOTIFICATION TYPES - Tipos/canales de notificación
// ===================================================

export interface NotificationType {
    id: number;
    name: string; // "WhatsApp", "SSE", "Email", "SMS"
    code: string; // "whatsapp", "sse", "email", "sms"
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

// ===================================================
// NOTIFICATION EVENT TYPES - Eventos por tipo
// ===================================================

export interface NotificationEventType {
    id: number;
    notification_type_id: number;
    event_code: string; // "order.created", "order.shipped", etc.
    event_name: string; // "Confirmación de Pedido", "Pedido Enviado"
    description?: string;
    template_config?: Record<string, any>;
    is_active: boolean;
    allowed_order_status_ids?: number[]; // Estados permitidos (vacío = todos)
    created_at: string;
    updated_at: string;
    // Relación expandida (opcional)
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

// ===================================================
// NOTIFICATION CONFIG - Configuraciones de negocio
// ===================================================

export interface NotificationConfig {
    id: number;
    business_id: number;
    integration_id: number; // NUEVO - La integración origen del evento
    notification_type_id: number; // NUEVO - Canal de salida (WhatsApp, SSE, etc.)
    notification_event_type_id: number; // NUEVO - Tipo de evento
    enabled: boolean;
    filters?: Record<string, any>;
    description?: string;
    created_at: string;
    updated_at: string;
    deleted_at?: string;

    // Relaciones expandidas (opcionales)
    notification_type?: NotificationType;
    notification_event_type?: NotificationEventType;
    order_status_ids?: number[]; // IDs de estados de orden a notificar

    // Campos adicionales enviados por el backend (para evitar cargar relaciones completas)
    notification_type_name?: string; // Nombre del canal (ej: "WhatsApp")
    notification_event_name?: string; // Nombre del evento (ej: "Nueva Orden")

    // Campos DEPRECADOS (mantener para compatibilidad)
    event_type?: string; // DEPRECATED
    channels?: string[]; // DEPRECATED
}

export interface CreateConfigDTO {
    business_id: number;
    integration_id: number; // La integración que genera el evento
    notification_type_id: number; // Canal de salida
    notification_event_type_id: number; // Tipo de evento
    enabled?: boolean;
    filters?: Record<string, any>;
    description?: string;
    order_status_ids?: number[]; // Estados de orden a notificar
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

// ===================================================
// ORDER STATUS - Estados de orden (relación M2M)
// ===================================================

export interface OrderStatus {
    id: number;
    code: string; // "pending", "processing", "completed", etc.
    name: string; // "Pendiente", "En Procesamiento", "Completada"
    description?: string;
    category: string; // "active", "completed", "cancelled", "refunded"
    is_active: boolean;
    icon?: string;
    color?: string;
    created_at: string;
    updated_at: string;
}

// ===================================================
// INTEGRATION - Integraciones (para selector)
// ===================================================

export interface Integration {
    id: number;
    name: string; // "Shopify - Mi Tiendita", "WhatsApp Business"
    code: string; // "shopify_store_1", "whatsapp_platform"
    type: string; // Tipo de integración (del IntegrationType)
    business_id?: number;
    is_active: boolean;
    integration_type_id?: number;
    integration_type_name?: string;
    integration_type_icon?: string;
    created_at: string;
    updated_at: string;
}

// ===================================================
// SYNC - Sincronización batch de reglas
// ===================================================

export interface SyncRule {
    id?: number; // undefined = crear, valor = actualizar
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

// ===================================================
// MESSAGE AUDIT - Auditoría de mensajes
// ===================================================

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
