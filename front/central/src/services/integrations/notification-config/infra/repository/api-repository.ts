/**
 * Infrastructure Layer - Repositorio API
 * Implementa los ports del dominio usando llamadas HTTP
 */

import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  OrderStatus,
  WhatsAppTemplate,
} from '../../domain/types';
import type {
  INotificationConfigRepository,
  IPaymentMethodRepository,
  IOrderStatusRepository,
  IWhatsAppTemplateRepository,
} from '../../domain/ports';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

/**
 * NotificationConfigApiRepository - Implementación HTTP del repositorio
 */
export class NotificationConfigApiRepository implements INotificationConfigRepository {
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
    // TODO: Obtener token del contexto de autenticación
    return localStorage.getItem('token') || '';
  }
}

/**
 * PaymentMethodApiRepository
 */
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
    return data.data || data; // Manejar ambos formatos de respuesta
  }

  private getToken(): string {
    return localStorage.getItem('token') || '';
  }
}

/**
 * OrderStatusApiRepository
 */
export class OrderStatusApiRepository implements IOrderStatusRepository {
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

/**
 * WhatsAppTemplateApiRepository
 */
export class WhatsAppTemplateApiRepository implements IWhatsAppTemplateRepository {
  private baseUrl = `${API_BASE_URL}/integrations/whatsapp`;

  async list(integrationId: number): Promise<WhatsAppTemplate[]> {
    const response = await fetch(`${this.baseUrl}/templates?integration_id=${integrationId}`, {
      headers: {
        Authorization: `Bearer ${this.getToken()}`,
      },
    });

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
