/**
 * UI Layer - Custom Hooks
 * Conectan los componentes React con los casos de uso
 */

'use client';

import { useState, useEffect, useCallback } from 'react';
import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  OrderStatus,
  WhatsAppTemplate,
} from '../../domain/types';
import {
  CreateNotificationConfigUseCase,
  UpdateNotificationConfigUseCase,
  GetNotificationConfigUseCase,
  ListNotificationConfigsUseCase,
  DeleteNotificationConfigUseCase,
  ToggleNotificationConfigUseCase,
  GetPaymentMethodsUseCase,
  GetOrderStatusesUseCase,
  GetWhatsAppTemplatesUseCase,
} from '../../app/use-cases';
import {
  NotificationConfigApiRepository,
  PaymentMethodApiRepository,
  OrderStatusApiRepository,
  WhatsAppTemplateApiRepository,
} from '../../infra/repository/api-repository';

// Instancias de repositorios (singleton pattern)
const notificationConfigRepo = new NotificationConfigApiRepository();
const paymentMethodRepo = new PaymentMethodApiRepository();
const orderStatusRepo = new OrderStatusApiRepository();
const whatsappTemplateRepo = new WhatsAppTemplateApiRepository();

// Instancias de casos de uso
const createUseCase = new CreateNotificationConfigUseCase(notificationConfigRepo);
const updateUseCase = new UpdateNotificationConfigUseCase(notificationConfigRepo);
const getUseCase = new GetNotificationConfigUseCase(notificationConfigRepo);
const listUseCase = new ListNotificationConfigsUseCase(notificationConfigRepo);
const deleteUseCase = new DeleteNotificationConfigUseCase(notificationConfigRepo);
const toggleUseCase = new ToggleNotificationConfigUseCase(notificationConfigRepo);
const getPaymentMethodsUseCase = new GetPaymentMethodsUseCase(paymentMethodRepo);
const getOrderStatusesUseCase = new GetOrderStatusesUseCase(orderStatusRepo);
const getWhatsAppTemplatesUseCase = new GetWhatsAppTemplatesUseCase(whatsappTemplateRepo);

/**
 * useNotificationConfigs - Hook principal para gestionar configuraciones
 */
export function useNotificationConfigs(filters?: FilterNotificationConfigDTO) {
  const [configs, setConfigs] = useState<IntegrationNotificationConfig[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchConfigs = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await listUseCase.execute(filters || {});
      setConfigs(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    fetchConfigs();
  }, [fetchConfigs]);

  const createConfig = async (dto: CreateNotificationConfigDTO) => {
    setLoading(true);
    setError(null);
    try {
      const newConfig = await createUseCase.execute(dto);
      setConfigs((prev) => [...prev, newConfig]);
      return newConfig;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error creando configuración');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const updateConfig = async (id: number, dto: UpdateNotificationConfigDTO) => {
    setLoading(true);
    setError(null);
    try {
      const updatedConfig = await updateUseCase.execute(id, dto);
      setConfigs((prev) => prev.map((c) => (c.id === id ? updatedConfig : c)));
      return updatedConfig;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error actualizando configuración');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const deleteConfig = async (id: number) => {
    setLoading(true);
    setError(null);
    try {
      await deleteUseCase.execute(id);
      setConfigs((prev) => prev.filter((c) => c.id !== id));
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error eliminando configuración');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  const toggleConfig = async (id: number, isActive: boolean) => {
    setLoading(true);
    setError(null);
    try {
      const updatedConfig = await toggleUseCase.execute(id, isActive);
      setConfigs((prev) => prev.map((c) => (c.id === id ? updatedConfig : c)));
      return updatedConfig;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error cambiando estado');
      throw err;
    } finally {
      setLoading(false);
    }
  };

  return {
    configs,
    loading,
    error,
    refetch: fetchConfigs,
    createConfig,
    updateConfig,
    deleteConfig,
    toggleConfig,
  };
}

/**
 * useNotificationConfig - Hook para una configuración específica
 */
export function useNotificationConfig(id: number | null) {
  const [config, setConfig] = useState<IntegrationNotificationConfig | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    const fetchConfig = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await getUseCase.execute(id);
        setConfig(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error obteniendo configuración');
      } finally {
        setLoading(false);
      }
    };

    fetchConfig();
  }, [id]);

  return { config, loading, error };
}

/**
 * usePaymentMethods - Hook para obtener métodos de pago
 */
export function usePaymentMethods() {
  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchPaymentMethods = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await getPaymentMethodsUseCase.execute();
        setPaymentMethods(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error obteniendo métodos de pago');
      } finally {
        setLoading(false);
      }
    };

    fetchPaymentMethods();
  }, []);

  return { paymentMethods, loading, error };
}

/**
 * useOrderStatuses - Hook para obtener estados de orden
 */
export function useOrderStatuses() {
  const [orderStatuses, setOrderStatuses] = useState<OrderStatus[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchOrderStatuses = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await getOrderStatusesUseCase.execute();
        setOrderStatuses(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error obteniendo estados de orden');
      } finally {
        setLoading(false);
      }
    };

    fetchOrderStatuses();
  }, []);

  return { orderStatuses, loading, error };
}

/**
 * useWhatsAppTemplates - Hook para obtener plantillas de WhatsApp
 */
export function useWhatsAppTemplates(integrationId: number | null) {
  const [templates, setTemplates] = useState<WhatsAppTemplate[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!integrationId) return;

    const fetchTemplates = async () => {
      setLoading(true);
      setError(null);
      try {
        const data = await getWhatsAppTemplatesUseCase.execute(integrationId);
        setTemplates(data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Error obteniendo plantillas de WhatsApp');
      } finally {
        setLoading(false);
      }
    };

    fetchTemplates();
  }, [integrationId]);

  return { templates, loading, error };
}
