/**
 * Application Layer - Casos de Uso
 * Orquesta la lógica de negocio usando los ports del dominio
 */

import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  OrderStatus,
  WhatsAppTemplate,
} from '../domain/types';
import type {
  INotificationConfigRepository,
  IPaymentMethodRepository,
  IOrderStatusRepository,
  IWhatsAppTemplateRepository,
} from '../domain/ports';

/**
 * CreateNotificationConfigUseCase
 */
export class CreateNotificationConfigUseCase {
  constructor(private repository: INotificationConfigRepository) {}

  async execute(dto: CreateNotificationConfigDTO): Promise<IntegrationNotificationConfig> {
    // Validaciones de negocio
    if (dto.priority < 0) {
      throw new Error('La prioridad debe ser un número positivo');
    }

    if (dto.conditions.statuses.length === 0 && dto.conditions.payment_methods.length === 0) {
      console.warn('Configuración sin filtros: se aplicará a todas las órdenes con este trigger');
    }

    return await this.repository.create(dto);
  }
}

/**
 * UpdateNotificationConfigUseCase
 */
export class UpdateNotificationConfigUseCase {
  constructor(private repository: INotificationConfigRepository) {}

  async execute(id: number, dto: UpdateNotificationConfigDTO): Promise<IntegrationNotificationConfig> {
    if (dto.priority !== undefined && dto.priority < 0) {
      throw new Error('La prioridad debe ser un número positivo');
    }

    return await this.repository.update(id, dto);
  }
}

/**
 * GetNotificationConfigUseCase
 */
export class GetNotificationConfigUseCase {
  constructor(private repository: INotificationConfigRepository) {}

  async execute(id: number): Promise<IntegrationNotificationConfig> {
    return await this.repository.getById(id);
  }
}

/**
 * ListNotificationConfigsUseCase
 */
export class ListNotificationConfigsUseCase {
  constructor(private repository: INotificationConfigRepository) {}

  async execute(filters: FilterNotificationConfigDTO): Promise<IntegrationNotificationConfig[]> {
    return await this.repository.list(filters);
  }
}

/**
 * DeleteNotificationConfigUseCase
 */
export class DeleteNotificationConfigUseCase {
  constructor(private repository: INotificationConfigRepository) {}

  async execute(id: number): Promise<void> {
    return await this.repository.delete(id);
  }
}

/**
 * GetPaymentMethodsUseCase
 */
export class GetPaymentMethodsUseCase {
  constructor(private repository: IPaymentMethodRepository) {}

  async execute(): Promise<PaymentMethod[]> {
    return await this.repository.list();
  }
}

/**
 * GetOrderStatusesUseCase
 */
export class GetOrderStatusesUseCase {
  constructor(private repository: IOrderStatusRepository) {}

  async execute(): Promise<OrderStatus[]> {
    return await this.repository.list();
  }
}

/**
 * GetWhatsAppTemplatesUseCase
 */
export class GetWhatsAppTemplatesUseCase {
  constructor(private repository: IWhatsAppTemplateRepository) {}

  async execute(integrationId: number): Promise<WhatsAppTemplate[]> {
    return await this.repository.list(integrationId);
  }
}

/**
 * ToggleNotificationConfigUseCase - Caso de uso helper
 */
export class ToggleNotificationConfigUseCase {
  constructor(private repository: INotificationConfigRepository) {}

  async execute(id: number, isActive: boolean): Promise<IntegrationNotificationConfig> {
    return await this.repository.update(id, { is_active: isActive });
  }
}
