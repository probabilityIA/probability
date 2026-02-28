"use client";

import { useIntegrationsSimple } from "@/services/integrations/core/ui/hooks/useIntegrationsSimple";
import type { IntegrationSimple } from "@/services/integrations/core/domain/types";

interface IntegrationPickerProps {
  businessId?: number;
  onSelect: (integration: IntegrationSimple) => void;
  onCancel: () => void;
}

export function IntegrationPicker({ businessId, onSelect, onCancel }: IntegrationPickerProps) {
  const { integrations, loading } = useIntegrationsSimple(
    businessId ? { businessId } : undefined
  );

  // Filter to ecommerce integrations that are active
  const ecommerceIntegrations = integrations.filter(
    (i) => i.is_active && i.category === "ecommerce"
  );

  return (
    <div className="space-y-4">
      <p className="text-sm text-gray-600">
        Selecciona la integración de origen para configurar sus reglas de notificación.
      </p>

      {loading ? (
        <div className="text-center py-8 text-gray-500">Cargando integraciones...</div>
      ) : ecommerceIntegrations.length === 0 ? (
        <div className="text-center py-8">
          <p className="text-gray-500 text-sm">No hay integraciones ecommerce disponibles</p>
          <p className="text-gray-400 text-xs mt-1">Crea una integración primero desde el módulo de Integraciones</p>
        </div>
      ) : (
        <div className="grid gap-2">
          {ecommerceIntegrations.map((integration) => (
            <button
              key={integration.id}
              type="button"
              onClick={() => onSelect(integration)}
              className="flex items-center gap-3 px-4 py-3 rounded-lg border border-gray-200 bg-white hover:bg-blue-50 hover:border-blue-300 transition-all text-left group"
            >
              {integration.image_url ? (
                <img
                  src={integration.image_url}
                  alt={integration.name}
                  className="w-9 h-9 object-contain rounded shrink-0"
                />
              ) : (
                <div className="w-9 h-9 rounded bg-gray-100 flex items-center justify-center shrink-0 group-hover:bg-blue-100">
                  <span className="text-sm font-bold text-gray-400 group-hover:text-blue-500">
                    {integration.type?.charAt(0).toUpperCase() || "?"}
                  </span>
                </div>
              )}
              <div className="min-w-0 flex-1">
                <p className="text-sm font-medium text-gray-900 truncate group-hover:text-blue-700">
                  {integration.name}
                </p>
                <p className="text-xs text-gray-500">
                  {integration.category_name || integration.type}
                </p>
              </div>
              <svg className="w-5 h-5 text-gray-300 group-hover:text-blue-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
              </svg>
            </button>
          ))}
        </div>
      )}

      <div className="flex justify-end pt-2 border-t">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-sm text-gray-600 hover:text-gray-900 transition-colors"
        >
          Cancelar
        </button>
      </div>
    </div>
  );
}
