'use client';

import { useState, FormEvent } from 'react';
import type {
  InvoicingConfig,
  InvoicingFilters,
  CreateConfigDTO,
} from '@/services/modules/invoicing/domain/types';
import { InvoicingFilterBuilder } from './InvoicingFilterBuilder';
import { useInvoicingConfig } from '@/services/modules/invoicing/ui/hooks/useInvoicingConfig';

interface InvoicingConfigFormProps {
  integrationIds: number[]; // Array de IDs de integraciones de e-commerce
  invoicingIntegrationId: number;
  businessId: number;
  onSuccess?: () => void;
  onCancel?: () => void;
  initialData?: InvoicingConfig;
}

/**
 * Formulario de configuración de facturación electrónica
 */
export function InvoicingConfigForm({
  integrationIds,
  invoicingIntegrationId,
  businessId,
  onSuccess,
  onCancel,
  initialData,
}: InvoicingConfigFormProps) {
  const { createConfig, updateConfig, loading } = useInvoicingConfig(businessId);

  const [formData, setFormData] = useState<Partial<InvoicingConfig>>({
    business_id: businessId,
    enabled: initialData?.enabled ?? true,
    auto_invoice: initialData?.auto_invoice ?? false,
    filters: initialData?.filters ?? {},
  });

  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);

    try {
      if (initialData?.id) {
        // Actualizar configuración existente
        const result = await updateConfig(initialData.id, {
          enabled: formData.enabled,
          auto_invoice: formData.auto_invoice,
          filters: formData.filters,
          invoicing_integration_id: invoicingIntegrationId,
        });

        if (result.success) {
          onSuccess?.();
        } else {
          setError(result.error || 'Error al actualizar configuración');
        }
      } else {
        // Crear múltiples configuraciones (una por cada tienda seleccionada)
        const results = await Promise.all(
          integrationIds.map((integrationId) => {
            const createData: CreateConfigDTO = {
              business_id: businessId,
              integration_id: integrationId,
              invoicing_integration_id: invoicingIntegrationId,
              enabled: formData.enabled,
              auto_invoice: formData.auto_invoice,
              filters: formData.filters,
            };

            return createConfig(createData);
          })
        );

        // Verificar si todas fueron exitosas
        const failedResults = results.filter((r) => !r.success);

        if (failedResults.length === 0) {
          onSuccess?.();
        } else if (failedResults.length === results.length) {
          setError(failedResults[0]?.error || 'Error al crear las configuraciones');
        } else {
          const failMessages = failedResults.map((r) => r.error).filter(Boolean).join('; ');
          setError(
            `Se crearon ${results.length - failedResults.length} de ${results.length} configuraciones.${failMessages ? ` Errores: ${failMessages}` : ' Algunas fallaron.'}`
          );
        }
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
          disabled={loading}
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
