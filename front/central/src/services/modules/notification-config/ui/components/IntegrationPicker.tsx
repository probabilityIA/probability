"use client";

import { useIntegrationsSimple } from "@/services/integrations/core/ui/hooks/useIntegrationsSimple";
import type { IntegrationSimple } from "@/services/integrations/core/domain/types";

interface IntegrationPickerProps {
  businessId?: number;
  onSelect: (integration: IntegrationSimple) => void;
  onCancel: () => void;
}

function IntegrationButton({ integration, onSelect }: { integration: IntegrationSimple; onSelect: (i: IntegrationSimple) => void }) {
  return (
    <button
      type="button"
      onClick={() => onSelect(integration)}
      className="flex items-center gap-3 px-4 py-3 rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 hover:bg-purple-50 dark:hover:bg-purple-900/20 hover:border-purple-300 dark:hover:border-purple-600 transition-all text-left group"
    >
      {integration.image_url ? (
        <img
          src={integration.image_url}
          alt={integration.name}
          className="w-9 h-9 object-contain rounded shrink-0"
        />
      ) : (
        <div className="w-9 h-9 rounded bg-gray-100 flex items-center justify-center shrink-0 group-hover:bg-purple-100 dark:group-hover:bg-purple-900/30">
          <span className="text-sm font-bold text-gray-400 group-hover:text-purple-600 dark:group-hover:text-purple-400">
            {integration.name?.charAt(0).toUpperCase() || "?"}
          </span>
        </div>
      )}
      <div className="min-w-0 flex-1">
        <p className="text-sm font-medium text-gray-900 dark:text-white truncate group-hover:text-purple-700 dark:group-hover:text-purple-400">
          {integration.name}
        </p>
        <p className="text-xs text-gray-500 dark:text-gray-400">
          {integration.category_name || integration.type}
        </p>
      </div>
      <svg className="w-5 h-5 text-gray-300 group-hover:text-purple-500 dark:group-hover:text-purple-400 shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
      </svg>
    </button>
  );
}

export function IntegrationPicker({ businessId, onSelect, onCancel }: IntegrationPickerProps) {
  const { integrations, loading } = useIntegrationsSimple(
    businessId ? { businessId } : undefined
  );

  const platformIntegrations = integrations.filter(
    (i) => i.is_active && i.category === "platform"
  );
  const ecommerceIntegrations = integrations.filter(
    (i) => i.is_active && i.category === "ecommerce"
  );

  const hasAny = platformIntegrations.length > 0 || ecommerceIntegrations.length > 0;

  return (
    <div className="space-y-4">
      <p className="text-sm text-gray-600 dark:text-gray-300">
        Selecciona la integración de origen para configurar sus reglas de notificación.
      </p>

      {loading ? (
        <div className="text-center py-8 text-gray-500 dark:text-gray-400">Cargando integraciones...</div>
      ) : !hasAny ? (
        <div className="text-center py-8">
          <p className="text-gray-500 dark:text-gray-400 text-sm">No hay integraciones disponibles</p>
          <p className="text-gray-400 text-xs mt-1">Crea una integración primero desde el módulo de Integraciones</p>
        </div>
      ) : (
        <div className="space-y-3">
          {platformIntegrations.length > 0 && (
            <div>
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">Plataforma</p>
              <div className="grid gap-2">
                {platformIntegrations.map((integration) => (
                  <IntegrationButton key={integration.id} integration={integration} onSelect={onSelect} />
                ))}
              </div>
            </div>
          )}
          {ecommerceIntegrations.length > 0 && (
            <div>
              <p className="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-2">E-commerce</p>
              <div className="grid gap-2">
                {ecommerceIntegrations.map((integration) => (
                  <IntegrationButton key={integration.id} integration={integration} onSelect={onSelect} />
                ))}
              </div>
            </div>
          )}
        </div>
      )}

      <div className="flex justify-end pt-2 border-t">
        <button
          type="button"
          onClick={onCancel}
          className="px-4 py-2 text-sm text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:text-white transition-colors"
        >
          Cancelar
        </button>
      </div>
    </div>
  );
}
