/**
 * Domain - Ports (Interfaces/Contratos)
 * Definen cómo interactuar con el mundo exterior
 */

import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  OrderStatus,
  WhatsAppTemplate,
} from './types';

/**
 * INotificationConfigRepository - Puerto para acceso a datos
 * Implementado por la capa de infraestructura
 */
export interface INotificationConfigRepository {
  create(dto: CreateNotificationConfigDTO): Promise<IntegrationNotificationConfig>;
  update(id: number, dto: UpdateNotificationConfigDTO): Promise<IntegrationNotificationConfig>;
  getById(id: number): Promise<IntegrationNotificationConfig>;
  list(filters: FilterNotificationConfigDTO): Promise<IntegrationNotificationConfig[]>;
  delete(id: number): Promise<void>;
}

/**
 * IPaymentMethodRepository - Puerto para obtener métodos de pago
 */
export interface IPaymentMethodRepository {
  list(): Promise<PaymentMethod[]>;
}

/**
 * IOrderStatusRepository - Puerto para obtener estados de orden
 */
export interface IOrderStatusRepository {
  list(): Promise<OrderStatus[]>;
}

/**
 * IWhatsAppTemplateRepository - Puerto para obtener plantillas de WhatsApp
 */
export interface IWhatsAppTemplateRepository {
  list(integrationId: number): Promise<WhatsAppTemplate[]>;
}
