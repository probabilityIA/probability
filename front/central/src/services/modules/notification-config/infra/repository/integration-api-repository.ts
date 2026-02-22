/**
 * Infrastructure Layer - Repositorios API para configuraciones de notificación por integración
 * (Migrado desde services/integrations/notification-config)
 */

import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  WhatsAppTemplate,
} from '../../domain/integration-types';
import type { OrderStatus } from '../../domain/types';
import type {
  IIntegrationNotificationConfigRepository,
  IPaymentMethodRepository,
  IIntegrationOrderStatusRepository,
  IWhatsAppTemplateRepository,
} from '../../domain/integration-ports';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

export class IntegrationNotificationConfigApiRepository
  implements IIntegrationNotificationConfigRepository
{
  private baseUrl = `${API_BASE_URL}/integrations/notification-configs`;

  async create(dto: CreateNotificationConfigDTO): Promise<IntegrationNotificationConfig> {
    const response = await fetch(this.baseUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${this.getToken()}`,
      },
      body: JSON.stringify(dto),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Error creando configuración de notificación');
    }

    return await response.json();
  }

  async update(id: number, dto: UpdateNotificationConfigDTO): Promise<IntegrationNotificationConfig> {
    const response = await fetch(`${this.baseUrl}/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${this.getToken()}`,
      },
      body: JSON.stringify(dto),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Error actualizando configuración de notificación');
    }

    return await response.json();
  }

  async getById(id: number): Promise<IntegrationNotificationConfig> {
    const response = await fetch(`${this.baseUrl}/${id}`, {
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Error obteniendo configuración de notificación');
    }

    return await response.json();
  }

  async list(filters: FilterNotificationConfigDTO): Promise<IntegrationNotificationConfig[]> {
    const params = new URLSearchParams();

    if (filters.integration_id) params.append('integration_id', filters.integration_id.toString());
    if (filters.notification_type) params.append('notification_type', filters.notification_type);
    if (filters.is_active !== undefined) params.append('is_active', filters.is_active.toString());
    if (filters.trigger) params.append('trigger', filters.trigger);

    const url = `${this.baseUrl}?${params.toString()}`;
    const response = await fetch(url, {
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Error listando configuraciones de notificación');
    }

    return await response.json();
  }

  async delete(id: number): Promise<void> {
    const response = await fetch(`${this.baseUrl}/${id}`, {
      method: 'DELETE',
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Error eliminando configuración de notificación');
    }
  }

  private getToken(): string {
    return localStorage.getItem('token') || '';
  }
}

export class PaymentMethodApiRepository implements IPaymentMethodRepository {
  private baseUrl = `${API_BASE_URL}/payment-methods`;

  async list(): Promise<PaymentMethod[]> {
    const response = await fetch(this.baseUrl, {
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error('Error obteniendo métodos de pago');
    }

    const data = await response.json();
    return data.data || data;
  }

  private getToken(): string {
    return localStorage.getItem('token') || '';
  }
}

export class IntegrationOrderStatusApiRepository implements IIntegrationOrderStatusRepository {
  private baseUrl = `${API_BASE_URL}/order-statuses`;

  async list(): Promise<OrderStatus[]> {
    const response = await fetch(this.baseUrl, {
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

    if (!response.ok) {
      throw new Error('Error obteniendo estados de orden');
    }

    const data = await response.json();
    return data.data || data;
  }

  private getToken(): string {
    return localStorage.getItem('token') || '';
  }
}

export class WhatsAppTemplateApiRepository implements IWhatsAppTemplateRepository {
  private baseUrl = `${API_BASE_URL}/integrations/whatsapp`;

  async list(integrationId: number): Promise<WhatsAppTemplate[]> {
    const response = await fetch(
      `${this.baseUrl}/templates?integration_id=${integrationId}`,
      {
        headers: {
          Authorization: `Bearer ${this.getToken()}`,
        },
      }
    );

    if (!response.ok) {
      throw new Error('Error obteniendo plantillas de WhatsApp');
    }

    const data = await response.json();
    return data.templates || data.data || data;
  }

  private getToken(): string {
    return localStorage.getItem('token') || '';
  }
}
