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
}

export interface UpdateNotificationEventTypeDTO {
    event_name?: string;
    description?: string;
    template_config?: Record<string, any>;
    is_active?: boolean;
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
    created_at: string;
    updated_at: string;
}
