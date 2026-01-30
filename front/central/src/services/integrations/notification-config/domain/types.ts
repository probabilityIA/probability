/**
 * Domain - Tipos del núcleo de negocio
 * No tienen dependencias externas
 */

export type NotificationType = 'whatsapp' | 'email' | 'sms';
export type TriggerType = 'order.created' | 'order.updated' | 'order.status_changed';
export type RecipientType = 'customer' | 'business';

export interface NotificationConditions {
  trigger: TriggerType;
  statuses: string[];
  payment_methods: number[];
}

export interface NotificationConfig {
  template_name: string;
  recipient_type: RecipientType;
  language: string;
}

export interface IntegrationNotificationConfig {
  id: number;
  integration_id: number;
  notification_type: NotificationType;
  is_active: boolean;
  conditions: NotificationConditions;
  config: NotificationConfig;
  description: string;
  priority: number;
  created_at: string;
  updated_at: string;
}

export interface CreateNotificationConfigDTO {
  integration_id: number;
  notification_type: NotificationType;
  is_active: boolean;
  conditions: NotificationConditions;
  config: NotificationConfig;
  description: string;
  priority: number;
}

export interface UpdateNotificationConfigDTO {
  notification_type?: NotificationType;
  is_active?: boolean;
  conditions?: NotificationConditions;
  config?: NotificationConfig;
  description?: string;
  priority?: number;
}

export interface FilterNotificationConfigDTO {
  integration_id?: number;
  notification_type?: NotificationType;
  is_active?: boolean;
  trigger?: TriggerType;
}

// Tipos auxiliares
export interface PaymentMethod {
  id: number;
  name: string;
  code: string;
  description: string;
  is_active: boolean;
}

export interface OrderStatus {
  id: number;
  name: string;
  code: string;
  description: string;
  color: string;
}

export interface WhatsAppTemplate {
  name: string;
  language: string;
  status: string;
  category: string;
  components: any[];
}
