'use client';

import { useState, useEffect, FormEvent } from 'react';
import type {
  InvoicingConfig,
  CreateConfigDTO,
} from '@/services/modules/invoicing/domain/types';
import type { Integration } from '@/services/integrations/core/domain/types';
import { getIntegrationsAction } from '@/services/integrations/core/infra/actions';
import { useInvoicingConfig } from '@/services/modules/invoicing/ui/hooks/useInvoicingConfig';

interface InvoicingConfigFormProps {
  integrationIds: number[];
  invoicingIntegrationId: number;
  businessId: number;
  onSuccess?: () => void;
  onCancel?: () => void;
  initialData?: InvoicingConfig;
}

export function InvoicingConfigForm({
  integrationIds,
  invoicingIntegrationId,
  businessId,
  onSuccess,
  onCancel,
  initialData,
}: InvoicingConfigFormProps) {
  const { createConfig, updateConfig, loading } = useInvoicingConfig(businessId);

  const [formData, setFormData] = useState({
    enabled: initialData?.enabled ?? true,
    auto_invoice: initialData?.auto_invoice ?? false,
    payment_status: (initialData?.filters?.payment_status as string) ?? '',
  });

  // Selección de integraciones de origen
  const initialSelected = initialData?.integration_ids?.length
    ? initialData.integration_ids
    : integrationIds;
  const [selectedIntegrationIds, setSelectedIntegrationIds] = useState<number[]>(initialSelected);
  const [availableIntegrations, setAvailableIntegrations] = useState<Integration[]>([]);
  const [loadingIntegrations, setLoadingIntegrations] = useState(false);

  const [error, setError] = useState<string | null>(null);

  // Cargar integraciones de origen disponibles del negocio (ecommerce + platform)
  useEffect(() => {
    if (!businessId) return;
    setLoadingIntegrations(true);
    Promise.all([
      getIntegrationsAction({ business_id: businessId, category: 'ecommerce', page_size: 100 }),
      getIntegrationsAction({ business_id: businessId, category: 'platform', page_size: 100 }),
    ])
      .then(([ecommerce, platform]) => {
        setAvailableIntegrations([...(ecommerce.data ?? []), ...(platform.data ?? [])]);
      })
      .catch(() => setAvailableIntegrations([]))
      .finally(() => setLoadingIntegrations(false));
  }, [businessId]);

  const toggleIntegration = (id: number) => {
    setSelectedIntegrationIds((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
    );
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);

    if (selectedIntegrationIds.length === 0) {
      setError('Debes seleccionar al menos una integración de origen de órdenes');
      return;
    }

    const filters = formData.payment_status
      ? { payment_status: formData.payment_status as 'paid' | 'unpaid' | 'partial' }
      : {};

    try {
      if (initialData?.id) {
        const result = await updateConfig(initialData.id, {
          enabled: formData.enabled,
          auto_invoice: formData.auto_invoice,
          filters,
          integration_ids: selectedIntegrationIds,
        });

        if (result.success) {
          onSuccess?.();
        } else {
          setError(result.error || 'Error al actualizar configuración');
        }
      } else {
        const createData: CreateConfigDTO = {
          business_id: businessId,
          integration_ids: selectedIntegrationIds,
          invoicing_integration_id: invoicingIntegrationId,
          enabled: formData.enabled,
          auto_invoice: formData.auto_invoice,
          filters,
        };

        const result = await createConfig(createData);

        if (result.success) {
          onSuccess?.();
        } else {
          setError(result.error || 'Error al crear la configuración');
        }
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      {/* Habilitar facturación */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <label className="flex items-center gap-3 cursor-pointer">
          <input
            type="checkbox"
            checked={formData.enabled}
            onChange={(e) => setFormData({ ...formData, enabled: e.target.checked })}
            disabled={loading}
            className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
          />
          <div>
            <span className="text-sm font-medium text-gray-900">Habilitar facturación</span>
            <p className="text-xs text-gray-500">Permite que esta integración genere facturas electrónicas</p>
          </div>
        </label>
      </div>

      {/* Facturación automática */}
      {formData.enabled && (
        <div className="bg-white p-4 rounded-lg border border-gray-200">
          <label className="flex items-center gap-3 cursor-pointer">
            <input
              type="checkbox"
              checked={formData.auto_invoice}
              onChange={(e) => setFormData({ ...formData, auto_invoice: e.target.checked })}
              disabled={loading}
              className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
            />
            <div>
              <span className="text-sm font-medium text-gray-900">Facturación automática</span>
              <p className="text-xs text-gray-500">Las órdenes que cumplan los filtros se facturarán automáticamente</p>
            </div>
          </label>
        </div>
      )}

      {/* Filtros de Pago */}
      {formData.enabled && formData.auto_invoice && (
        <div className="bg-white p-4 rounded-lg border border-gray-200">
          <h4 className="text-sm font-medium text-gray-900 mb-3">Filtros de Pago</h4>
          <div>
            <label className="block text-sm text-gray-700 mb-1">Estado de pago</label>
            <select
              value={formData.payment_status}
              onChange={(e) => setFormData({ ...formData, payment_status: e.target.value })}
              disabled={loading}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
            >
              <option value="">Sin filtro (todas las órdenes)</option>
              <option value="paid">Solo pagadas</option>
              <option value="unpaid">Solo sin pagar</option>
              <option value="partial">Pago parcial</option>
            </select>
          </div>
        </div>
      )}

      {/* Integraciones de origen de órdenes */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <h4 className="text-sm font-medium text-gray-900 mb-1">
          Fuentes de órdenes
        </h4>
        <p className="text-xs text-gray-500 mb-3">
          Selecciona las integraciones desde las cuales se facturarán las órdenes
        </p>

        {loadingIntegrations ? (
          <p className="text-xs text-gray-400">Cargando integraciones...</p>
        ) : availableIntegrations.length === 0 ? (
          <p className="text-xs text-gray-400">No hay integraciones e-commerce disponibles</p>
        ) : (
          <div className="space-y-2">
            {availableIntegrations.map((integration) => (
              <label
                key={integration.id}
                className="flex items-center gap-3 cursor-pointer p-2 rounded-md hover:bg-gray-50"
              >
                <input
                  type="checkbox"
                  checked={selectedIntegrationIds.includes(integration.id)}
                  onChange={() => toggleIntegration(integration.id)}
                  disabled={loading}
                  className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
                />
                <div className="flex items-center gap-2 min-w-0">
                  {integration.integration_type?.image_url && (
                    <img
                      src={integration.integration_type.image_url}
                      alt={integration.name}
                      className="w-5 h-5 object-contain flex-shrink-0"
                    />
                  )}
                  <span className="text-sm text-gray-800 truncate">{integration.name}</span>
                  {integration.type && (
                    <span className="text-xs text-gray-400 flex-shrink-0">({integration.type})</span>
                  )}
                </div>
              </label>
            ))}
          </div>
        )}
      </div>

      {/* Acciones */}
      <div className="flex items-center gap-3 pt-2 border-t">
        {onCancel && (
          <button
            type="button"
            onClick={onCancel}
            disabled={loading}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50"
          >
            Cancelar
          </button>
        )}
        <button
          type="submit"
          disabled={loading}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? 'Guardando...' : initialData?.id ? 'Actualizar' : 'Crear Configuración'}
        </button>
      </div>
    </form>
  );
}
