/**
 * Application Layer - Casos de uso para configuraciones de notificación por integración
 * (Migrado desde services/integrations/notification-config)
 */

import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  WhatsAppTemplate,
} from '../domain/integration-types';
import type { OrderStatus } from '../domain/types';
import type {
  IIntegrationNotificationConfigRepository,
  IPaymentMethodRepository,
  IIntegrationOrderStatusRepository,
  IWhatsAppTemplateRepository,
} from '../domain/integration-ports';

export class CreateIntegrationNotificationConfigUseCase {
  constructor(private repository: IIntegrationNotificationConfigRepository) {}

  async execute(dto: CreateNotificationConfigDTO): Promise<IntegrationNotificationConfig> {
    if (dto.priority < 0) {
      throw new Error('La prioridad debe ser un número positivo');
    }

    if (dto.conditions.statuses.length === 0 && dto.conditions.payment_methods.length === 0) {
      console.warn('Configuración sin filtros: se aplicará a todas las órdenes con este trigger');
    }

    return await this.repository.create(dto);
  }
}

export class UpdateIntegrationNotificationConfigUseCase {
  constructor(private repository: IIntegrationNotificationConfigRepository) {}

  async execute(id: number, dto: UpdateNotificationConfigDTO): Promise<IntegrationNotificationConfig> {
    if (dto.priority !== undefined && dto.priority < 0) {
      throw new Error('La prioridad debe ser un número positivo');
    }

    return await this.repository.update(id, dto);
  }
}

export class GetIntegrationNotificationConfigUseCase {
  constructor(private repository: IIntegrationNotificationConfigRepository) {}

  async execute(id: number): Promise<IntegrationNotificationConfig> {
    return await this.repository.getById(id);
  }
}

export class ListIntegrationNotificationConfigsUseCase {
  constructor(private repository: IIntegrationNotificationConfigRepository) {}

  async execute(filters: FilterNotificationConfigDTO): Promise<IntegrationNotificationConfig[]> {
    return await this.repository.list(filters);
  }
}

export class DeleteIntegrationNotificationConfigUseCase {
  constructor(private repository: IIntegrationNotificationConfigRepository) {}

  async execute(id: number): Promise<void> {
    return await this.repository.delete(id);
  }
}

export class ToggleIntegrationNotificationConfigUseCase {
  constructor(private repository: IIntegrationNotificationConfigRepository) {}

  async execute(id: number, isActive: boolean): Promise<IntegrationNotificationConfig> {
    return await this.repository.update(id, { is_active: isActive });
  }
}

export class GetPaymentMethodsUseCase {
  constructor(private repository: IPaymentMethodRepository) {}

  async execute(): Promise<PaymentMethod[]> {
    return await this.repository.list();
  }
}

export class GetIntegrationOrderStatusesUseCase {
  constructor(private repository: IIntegrationOrderStatusRepository) {}

  async execute(): Promise<OrderStatus[]> {
    return await this.repository.list();
  }
}

export class GetWhatsAppTemplatesUseCase {
  constructor(private repository: IWhatsAppTemplateRepository) {}

  async execute(integrationId: number): Promise<WhatsAppTemplate[]> {
    return await this.repository.list(integrationId);
  }
}
