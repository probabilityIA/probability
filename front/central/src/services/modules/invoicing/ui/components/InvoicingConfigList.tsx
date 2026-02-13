'use client';

import { useState } from 'react';
import type { InvoicingConfig } from '@/services/modules/invoicing/domain/types';
import { useInvoicingConfig } from '@/services/modules/invoicing/ui/hooks/useInvoicingConfig';
import {
  TrashIcon,
  PencilIcon,
  CheckCircleIcon,
  XCircleIcon,
} from '@heroicons/react/24/outline';

interface InvoicingConfigListProps {
  configs: InvoicingConfig[];
  onEdit?: (config: InvoicingConfig) => void;
  onRefresh?: () => void;
}

/**
 * Lista de configuraciones de facturación
 */
export function InvoicingConfigList({
  configs,
  onEdit,
  onRefresh,
}: InvoicingConfigListProps) {
  const { deleteConfig, toggleConfig, toggleAutoInvoice, loading } = useInvoicingConfig();
  const [deletingId, setDeletingId] = useState<number | null>(null);

  const handleDelete = async (id: number) => {
    if (!confirm('¿Está seguro de eliminar esta configuración?')) {
      return;
    }

    setDeletingId(id);
    const result = await deleteConfig(id);
    setDeletingId(null);

    if (result.success) {
      onRefresh?.();
    }
  };

  const handleToggleEnabled = async (id: number, currentState: boolean) => {
    const result = await toggleConfig(id, !currentState);

    if (result.success) {
      onRefresh?.();
    }
  };

  const handleToggleAutoInvoice = async (id: number, currentState: boolean) => {
    const result = await toggleAutoInvoice(id, !currentState);

    if (result.success) {
      onRefresh?.();
    }
  };

  if (configs.length === 0) {
    return (
      <div className="text-center py-12 bg-gray-50 rounded-lg border-2 border-dashed border-gray-300">
        <p className="text-gray-500">No hay configuraciones de facturación.</p>
        <p className="text-sm text-gray-400 mt-1">
          Crea una nueva configuración para empezar.
        </p>
      </div>
    );
  }

  return (
    <div className="space-y-4">
      {configs.map((config) => {
        const isDeleting = deletingId === config.id;
        const activeFilters = Object.entries(config.filters || {}).filter(
          ([_, value]) =>
            value !== undefined &&
            value !== null &&
            (Array.isArray(value) ? value.length > 0 : true)
        );

        return (
          <div
            key={config.id}
            className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
          >
            <div className="flex items-start justify-between">
              {/* Info principal */}
              <div className="flex-1">
                <div className="flex items-center gap-3 mb-2">
                  {/* Estado */}
                  {config.enabled ? (
                    <CheckCircleIcon className="w-5 h-5 text-green-500" />
                  ) : (
                    <XCircleIcon className="w-5 h-5 text-gray-400" />
                  )}

                  <h3 className="font-medium text-gray-900">
                    Integración #{config.integration_id}
                  </h3>

                  {/* Badge de auto-facturación - CLICKEABLE */}
                  <button
                    onClick={() =>
                      config.id && handleToggleAutoInvoice(config.id, config.auto_invoice)
                    }
                    disabled={loading || !config.id}
                    className={`px-2 py-1 text-xs rounded-full transition-colors disabled:opacity-50 hover:opacity-80 ${
                      config.auto_invoice
                        ? 'bg-blue-100 text-blue-800 hover:bg-blue-200'
                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                    }`}
                    title={config.auto_invoice ? 'Click para desactivar auto-facturación' : 'Click para activar auto-facturación'}
                  >
                    {config.auto_invoice ? 'Automático' : 'Manual'}
                  </button>

                  {/* Badge de estado - CLICKEABLE */}
                  <button
                    onClick={() =>
                      config.id && handleToggleEnabled(config.id, config.enabled)
                    }
                    disabled={loading || !config.id}
                    className={`px-2 py-1 text-xs rounded-full transition-colors disabled:opacity-50 hover:opacity-80 ${
                      config.enabled
                        ? 'bg-green-100 text-green-800 hover:bg-green-200'
                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                    }`}
                    title={config.enabled ? 'Click para desactivar' : 'Click para activar'}
                  >
                    {config.enabled ? 'Activo' : 'Inactivo'}
                  </button>
                </div>

                {/* Detalles */}
                <dl className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
                  <div>
                    <dt className="text-gray-500">Proveedor ID:</dt>
                    <dd className="text-gray-900 font-medium">
                      {config.invoicing_provider_id}
                    </dd>
                  </div>

                  <div>
                    <dt className="text-gray-500">Facturación automática:</dt>
                    <dd className="text-gray-900">
                      {config.auto_invoice ? 'Sí' : 'No'}
                    </dd>
                  </div>

                  {config.auto_invoice && (
                    <div className="col-span-2">
                      <dt className="text-gray-500">Filtros activos:</dt>
                      <dd className="text-gray-900">
                        {activeFilters.length === 0 ? (
                          <span className="text-gray-400">
                            Sin filtros (todas las órdenes)
                          </span>
                        ) : (
                          <div className="flex flex-wrap gap-1 mt-1">
                            {activeFilters.map(([key, value]) => (
                              <span
                                key={key}
                                className="px-2 py-0.5 text-xs bg-gray-100 text-gray-700 rounded"
                              >
                                {key}
                              </span>
                            ))}
                          </div>
                        )}
                      </dd>
                    </div>
                  )}

                  {config.config?.include_shipping && (
                    <div>
                      <dt className="text-gray-500">Incluye envío:</dt>
                      <dd className="text-gray-900">Sí</dd>
                    </div>
                  )}

                  {config.config?.default_tax_rate && (
                    <div>
                      <dt className="text-gray-500">IVA:</dt>
                      <dd className="text-gray-900">
                        {config.config.default_tax_rate}%
                      </dd>
                    </div>
                  )}
                </dl>
              </div>

              {/* Acciones */}
              <div className="flex items-center gap-2 ml-4">
                {/* Toggle estado */}
                <button
                  onClick={() =>
                    config.id && handleToggleEnabled(config.id, config.enabled)
                  }
                  disabled={loading || !config.id}
                  className="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded disabled:opacity-50"
                  title={
                    config.enabled
                      ? 'Desactivar configuración'
                      : 'Activar configuración'
                  }
                >
                  {config.enabled ? (
                    <XCircleIcon className="w-5 h-5" />
                  ) : (
                    <CheckCircleIcon className="w-5 h-5" />
                  )}
                </button>

                {/* Editar */}
                {onEdit && (
                  <button
                    onClick={() => onEdit(config)}
                    disabled={loading}
                    className="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded disabled:opacity-50"
                    title="Editar configuración"
                  >
                    <PencilIcon className="w-5 h-5" />
                  </button>
                )}

                {/* Eliminar */}
                <button
                  onClick={() => config.id && handleDelete(config.id)}
                  disabled={loading || isDeleting || !config.id}
                  className="p-2 text-gray-600 hover:text-red-600 hover:bg-red-50 rounded disabled:opacity-50"
                  title="Eliminar configuración"
                >
                  {isDeleting ? (
                    <div className="w-5 h-5 border-2 border-red-500 border-t-transparent rounded-full animate-spin" />
                  ) : (
                    <TrashIcon className="w-5 h-5" />
                  )}
                </button>
              </div>
            </div>

            {/* Notas (si existen) */}
            {config.config?.notes && (
              <div className="mt-3 pt-3 border-t border-gray-100">
                <p className="text-xs text-gray-500">Notas:</p>
                <p className="text-sm text-gray-700 mt-1">
                  {config.config.notes}
                </p>
              </div>
            )}
          </div>
        );
      })}
    </div>
  );
}
