/**
 * Domain - Ports para configuraciones de notificación por integración
 * (Migrado desde services/integrations/notification-config)
 */

import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  WhatsAppTemplate,
} from './integration-types';
import type { OrderStatus } from './types';

export interface IIntegrationNotificationConfigRepository {
  create(dto: CreateNotificationConfigDTO): Promise<IntegrationNotificationConfig>;
  update(id: number, dto: UpdateNotificationConfigDTO): Promise<IntegrationNotificationConfig>;
  getById(id: number): Promise<IntegrationNotificationConfig>;
  list(filters: FilterNotificationConfigDTO): Promise<IntegrationNotificationConfig[]>;
  delete(id: number): Promise<void>;
}

export interface IPaymentMethodRepository {
  list(): Promise<PaymentMethod[]>;
}

export interface IIntegrationOrderStatusRepository {
  list(): Promise<OrderStatus[]>;
}

export interface IWhatsAppTemplateRepository {
  list(integrationId: number): Promise<WhatsAppTemplate[]>;
}
