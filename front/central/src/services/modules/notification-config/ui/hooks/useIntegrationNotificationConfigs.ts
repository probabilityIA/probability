/**
 * UI Layer - Custom Hooks para configuraciones de notificación por integración
 * (Migrado desde services/integrations/notification-config)
 */

'use client';

import { useState, useEffect, useCallback } from 'react';
import type {
  IntegrationNotificationConfig,
  CreateNotificationConfigDTO,
  UpdateNotificationConfigDTO,
  FilterNotificationConfigDTO,
  PaymentMethod,
  WhatsAppTemplate,
} from '../../domain/integration-types';
import type { OrderStatus } from '../../domain/types';
import {
  CreateIntegrationNotificationConfigUseCase,
  UpdateIntegrationNotificationConfigUseCase,
  GetIntegrationNotificationConfigUseCase,
  ListIntegrationNotificationConfigsUseCase,
  DeleteIntegrationNotificationConfigUseCase,
  ToggleIntegrationNotificationConfigUseCase,
  GetPaymentMethodsUseCase,
  GetIntegrationOrderStatusesUseCase,
  GetWhatsAppTemplatesUseCase,
} from '../../app/integration-use-cases';
import {
  IntegrationNotificationConfigApiRepository,
  PaymentMethodApiRepository,
  IntegrationOrderStatusApiRepository,
  WhatsAppTemplateApiRepository,
} from '../../infra/repository/integration-api-repository';

// Instancias singleton de repositorios
const notificationConfigRepo = new IntegrationNotificationConfigApiRepository();
const paymentMethodRepo = new PaymentMethodApiRepository();
const orderStatusRepo = new IntegrationOrderStatusApiRepository();
const whatsappTemplateRepo = new WhatsAppTemplateApiRepository();

// Instancias singleton de casos de uso
const createUseCase = new CreateIntegrationNotificationConfigUseCase(notificationConfigRepo);
const updateUseCase = new UpdateIntegrationNotificationConfigUseCase(notificationConfigRepo);
const getUseCase = new GetIntegrationNotificationConfigUseCase(notificationConfigRepo);
const listUseCase = new ListIntegrationNotificationConfigsUseCase(notificationConfigRepo);
const deleteUseCase = new DeleteIntegrationNotificationConfigUseCase(notificationConfigRepo);
const toggleUseCase = new ToggleIntegrationNotificationConfigUseCase(notificationConfigRepo);
const getPaymentMethodsUseCase = new GetPaymentMethodsUseCase(paymentMethodRepo);
const getOrderStatusesUseCase = new GetIntegrationOrderStatusesUseCase(orderStatusRepo);
const getWhatsAppTemplatesUseCase = new GetWhatsAppTemplatesUseCase(whatsappTemplateRepo);

export function useIntegrationNotificationConfigs(filters?: FilterNotificationConfigDTO) {
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

export function useIntegrationNotificationConfig(id: number | null) {
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

export function useIntegrationOrderStatuses() {
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
