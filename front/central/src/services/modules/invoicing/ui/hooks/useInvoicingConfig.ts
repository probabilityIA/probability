'use client';

import { useState, useEffect } from 'react';
import type {
  InvoicingConfig,
  InvoicingProvider,
  CreateConfigDTO,
  UpdateConfigDTO,
} from '@/services/modules/invoicing/domain/types';
import {
  getConfigsAction,
  createConfigAction,
  updateConfigAction,
  deleteConfigAction,
  getProvidersAction,
} from '@/services/modules/invoicing/infra/actions';

/**
 * Hook para gestionar configuraciones de facturación
 */
export function useInvoicingConfig(businessId?: number) {
  const [configs, setConfigs] = useState<InvoicingConfig[]>([]);
  const [providers, setProviders] = useState<InvoicingProvider[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  /**
   * Cargar configuraciones de facturación
   */
  const fetchConfigs = async (bId?: number) => {
    if (!bId && !businessId) return;

    const targetBusinessId = bId || businessId!;

    setLoading(true);
    setError(null);

    try {
      const response = await getConfigsAction({ business_id: targetBusinessId });
      setConfigs(response.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Cargar proveedores de facturación
   */
  const fetchProviders = async (countryCode?: string) => {
    setLoading(true);
    setError(null);

    try {
      const response = await getProvidersAction({});
      setProviders(response.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    } finally {
      setLoading(false);
    }
  };

  /**
   * Crear configuración
   */
  const createConfig = async (data: CreateConfigDTO) => {
    setLoading(true);
    setError(null);

    try {
      await createConfigAction(data);

      // Recargar configuraciones
      await fetchConfigs();
      return { success: true };
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Error desconocido';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  };

  /**
   * Actualizar configuración
   */
  const updateConfig = async (id: number, data: UpdateConfigDTO) => {
    setLoading(true);
    setError(null);

    try {
      await updateConfigAction(id, data);

      // Recargar configuraciones
      await fetchConfigs();
      return { success: true };
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Error desconocido';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  };

  /**
   * Eliminar configuración
   */
  const deleteConfig = async (id: number) => {
    setLoading(true);
    setError(null);

    try {
      await deleteConfigAction(id);

      // Recargar configuraciones
      await fetchConfigs();
      return { success: true };
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Error desconocido';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  };

  /**
   * Activar/desactivar configuración
   */
  const toggleConfig = async (id: number, enabled: boolean) => {
    setLoading(true);
    setError(null);

    try {
      await updateConfigAction(id, { enabled });

      // Actualizar estado local
      setConfigs((prev) =>
        prev.map((config) =>
          config.id === id ? { ...config, enabled } : config
        )
      );
      return { success: true };
    } catch (err) {
      const errorMessage =
        err instanceof Error ? err.message : 'Error desconocido';
      setError(errorMessage);
      return { success: false, error: errorMessage };
    } finally {
      setLoading(false);
    }
  };

  // Cargar datos iniciales
  useEffect(() => {
    if (businessId) {
      fetchConfigs();
      // fetchProviders() eliminado - ya no se usan providers deprecados
    }
  }, [businessId]);

  return {
    configs,
    providers,
    loading,
    error,
    fetchConfigs,
    fetchProviders,
    createConfig,
    updateConfig,
    deleteConfig,
    toggleConfig,
  };
}
