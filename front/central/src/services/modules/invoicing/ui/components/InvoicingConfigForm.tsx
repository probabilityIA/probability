'use client';

import { useState, FormEvent } from 'react';
import type {
  InvoicingConfig,
  InvoicingFilters,
  InvoicingProvider,
  CreateConfigDTO,
} from '@/services/modules/invoicing/domain/types';
import { InvoicingProviderSelector } from './InvoicingProviderSelector';
import { InvoicingFilterBuilder } from './InvoicingFilterBuilder';
import { useInvoicingConfig } from '@/services/modules/invoicing/ui/hooks/useInvoicingConfig';

interface InvoicingConfigFormProps {
  integrationId: number;
  businessId: number;
  providers: InvoicingProvider[];
  onSuccess?: () => void;
  onCancel?: () => void;
  initialData?: InvoicingConfig;
}

/**
 * Formulario de configuración de facturación electrónica
 */
export function InvoicingConfigForm({
  integrationId,
  businessId,
  providers,
  onSuccess,
  onCancel,
  initialData,
}: InvoicingConfigFormProps) {
  const { createConfig, updateConfig, loading } = useInvoicingConfig(businessId);

  const [formData, setFormData] = useState<Partial<InvoicingConfig>>({
    integration_id: integrationId,
    business_id: businessId,
    enabled: initialData?.enabled ?? true,
    auto_invoice: initialData?.auto_invoice ?? false,
    invoicing_provider_id: initialData?.invoicing_provider_id,
    filters: initialData?.filters ?? {},
    config: initialData?.config ?? {},
  });

  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);

    // Validaciones
    if (!formData.invoicing_provider_id) {
      setError('Debe seleccionar un proveedor de facturación');
      return;
    }

    try {
      let result;

      if (initialData?.id) {
        // Actualizar configuración existente
        result = await updateConfig(initialData.id, {
          enabled: formData.enabled,
          auto_invoice: formData.auto_invoice,
          filters: formData.filters,
          config: formData.config,
          invoicing_provider_id: formData.invoicing_provider_id,
        });
      } else {
        // Crear nueva configuración
        const createData: CreateConfigDTO = {
          business_id: businessId,
          integration_id: integrationId,
          invoicing_provider_id: formData.invoicing_provider_id,
          enabled: formData.enabled,
          auto_invoice: formData.auto_invoice,
          filters: formData.filters,
          config: formData.config,
        };

        result = await createConfig(createData);
      }

      if (result.success) {
        onSuccess?.();
      } else {
        setError(result.error || 'Error al guardar configuración');
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Error desconocido');
    }
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-6">
      {/* Error global */}
      {error && (
        <div className="bg-red-50 border border-red-200 rounded-lg p-4">
          <p className="text-sm text-red-800">{error}</p>
        </div>
      )}

      {/* Selector de proveedor */}
      <InvoicingProviderSelector
        providers={providers}
        value={formData.invoicing_provider_id}
        onChange={(providerId) =>
          setFormData({ ...formData, invoicing_provider_id: providerId })
        }
        disabled={loading}
      />

      {/* Toggle de habilitado */}
      <div className="bg-white p-4 rounded-lg border border-gray-200">
        <label className="flex items-center gap-3 cursor-pointer">
          <input
            type="checkbox"
            checked={formData.enabled}
            onChange={(e) =>
              setFormData({ ...formData, enabled: e.target.checked })
            }
            disabled={loading}
            className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
          />
          <div>
            <span className="text-sm font-medium text-gray-900">
              Habilitar facturación
            </span>
            <p className="text-xs text-gray-500">
              Permite que esta integración genere facturas electrónicas
            </p>
          </div>
        </label>
      </div>

      {/* Toggle de auto-facturación */}
      {formData.enabled && (
        <div className="bg-white p-4 rounded-lg border border-gray-200">
          <label className="flex items-center gap-3 cursor-pointer">
            <input
              type="checkbox"
              checked={formData.auto_invoice}
              onChange={(e) =>
                setFormData({ ...formData, auto_invoice: e.target.checked })
              }
              disabled={loading}
              className="w-5 h-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500 disabled:opacity-50"
            />
            <div>
              <span className="text-sm font-medium text-gray-900">
                Facturación automática
              </span>
              <p className="text-xs text-gray-500">
                Las órdenes que cumplan los filtros se facturarán
                automáticamente
              </p>
            </div>
          </label>
        </div>
      )}

      {/* Constructor de filtros (solo si auto-facturación está activa) */}
      {formData.enabled && formData.auto_invoice && (
        <InvoicingFilterBuilder
          filters={formData.filters || {}}
          onChange={(filters) => setFormData({ ...formData, filters })}
          disabled={loading}
        />
      )}

      {/* Configuración adicional */}
      {formData.enabled && (
        <div className="bg-white p-4 rounded-lg border border-gray-200">
          <h4 className="font-medium text-gray-900 mb-4">
            Configuración Adicional
          </h4>

          <div className="space-y-4">
            {/* Incluir envío */}
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={formData.config?.include_shipping ?? true}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    config: {
                      ...formData.config,
                      include_shipping: e.target.checked,
                    },
                  })
                }
                disabled={loading}
                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <span className="text-sm text-gray-700">
                Incluir costo de envío en factura
              </span>
            </label>

            {/* Aplicar descuento */}
            <label className="flex items-center gap-2">
              <input
                type="checkbox"
                checked={formData.config?.apply_discount ?? true}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    config: {
                      ...formData.config,
                      apply_discount: e.target.checked,
                    },
                  })
                }
                disabled={loading}
                className="rounded border-gray-300 text-blue-600 focus:ring-blue-500"
              />
              <span className="text-sm text-gray-700">
                Aplicar descuentos automáticamente
              </span>
            </label>

            {/* Tasa de impuesto */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Tasa de impuesto por defecto (%)
              </label>
              <input
                type="number"
                value={formData.config?.default_tax_rate ?? 19}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    config: {
                      ...formData.config,
                      default_tax_rate: Number(e.target.value),
                    },
                  })
                }
                placeholder="19"
                min="0"
                max="100"
                step="0.01"
                disabled={loading}
                className="w-full md:w-32 px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
              />
              <p className="text-xs text-gray-500 mt-1">
                IVA aplicable (ej: 19% en Colombia)
              </p>
            </div>

            {/* Notas */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Notas adicionales
              </label>
              <textarea
                value={formData.config?.notes ?? ''}
                onChange={(e) =>
                  setFormData({
                    ...formData,
                    config: {
                      ...formData.config,
                      notes: e.target.value,
                    },
                  })
                }
                placeholder="Notas que aparecerán en la factura..."
                rows={3}
                disabled={loading}
                className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
              />
            </div>
          </div>
        </div>
      )}

      {/* Botones de acción */}
      <div className="flex items-center gap-3 pt-4 border-t">
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
          disabled={loading || !formData.invoicing_provider_id}
          className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading
            ? 'Guardando...'
            : initialData?.id
            ? 'Actualizar Configuración'
            : 'Crear Configuración'}
        </button>
      </div>
    </form>
  );
}
