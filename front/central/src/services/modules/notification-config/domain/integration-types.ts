/**
 * Domain - Tipos para configuraciones de notificación por integración
 * (Migrado desde services/integrations/notification-config)
 */

export type IntegrationNotifType = 'whatsapp' | 'email' | 'sms';
export type TriggerType = 'order.created' | 'order.updated' | 'order.status_changed';
export type RecipientType = 'customer' | 'business';

export interface NotificationConditions {
  trigger: TriggerType;
  statuses: string[];
  payment_methods: number[];
  source_integration_id?: number | null;
}

export interface IntegrationNotifDetails {
  template_name: string;
  recipient_type: RecipientType;
  language: string;
}

export interface IntegrationNotificationConfig {
  id: number;
  integration_id: number;
  notification_type: IntegrationNotifType;
  is_active: boolean;
  conditions: NotificationConditions;
  config: IntegrationNotifDetails;
  description: string;
  priority: number;
  created_at: string;
  updated_at: string;
}

export interface CreateNotificationConfigDTO {
  integration_id: number;
  notification_type: IntegrationNotifType;
  is_active: boolean;
  conditions: NotificationConditions;
  config: IntegrationNotifDetails;
  description: string;
  priority: number;
}

export interface UpdateNotificationConfigDTO {
  notification_type?: IntegrationNotifType;
  is_active?: boolean;
  conditions?: NotificationConditions;
  config?: IntegrationNotifDetails;
  description?: string;
  priority?: number;
}

export interface FilterNotificationConfigDTO {
  integration_id?: number;
  notification_type?: IntegrationNotifType;
  is_active?: boolean;
  trigger?: TriggerType;
}

export interface PaymentMethod {
  id: number;
  name: string;
  code: string;
  description: string;
  is_active: boolean;
}

export interface WhatsAppTemplate {
  name: string;
  language: string;
  status: string;
  category: string;
  components: any[];
}
