/**
 * IntegrationSourceSelector - Selector de integración origen
 * (Migrado desde services/integrations/notification-config)
 */
'use client';

import { useState, useEffect } from 'react';
import { useIntegrations } from '@/services/integrations/core/ui/hooks/useIntegrations';
import type { Integration } from '@/services/integrations/core/domain/types';

interface IntegrationSourceSelectorProps {
  businessId: number;
  value: number | null | undefined;
  onChange: (value: number | null) => void;
  disabled?: boolean;
}

export function IntegrationSourceSelector({
  businessId,
  value,
  onChange,
  disabled = false,
}: IntegrationSourceSelectorProps) {
  const { integrations, loading } = useIntegrations();
  const [externalIntegrations, setExternalIntegrations] = useState<Integration[]>([]);

  useEffect(() => {
    const filtered = integrations.filter(
      (integration) =>
        integration.category === 'external' &&
        integration.is_active &&
        integration.business_id === businessId
    );
    setExternalIntegrations(filtered);
  }, [integrations, businessId]);

  if (loading) {
    return (
      <div className="animate-pulse">
        <div className="h-10 bg-gray-200 rounded"></div>
      </div>
    );
  }

  return (
    <div>
      <label className="block text-sm font-medium text-gray-700 mb-2">
        Integración Origen
        <span className="text-gray-500 text-xs ml-2">(opcional)</span>
      </label>
      <select
        value={value ?? ''}
        onChange={(e) => {
          const val = e.target.value;
          onChange(val === '' ? null : parseInt(val, 10));
        }}
        disabled={disabled}
        className="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-blue-500 focus:border-blue-500 disabled:bg-gray-100"
      >
        <option value="">Todas las integraciones</option>
        {externalIntegrations.map((integration) => (
          <option key={integration.id} value={integration.id}>
            {integration.name} ({integration.type})
          </option>
        ))}
      </select>
      <p className="text-xs text-gray-500 mt-1">
        Filtra notificaciones solo para órdenes de esta integración
      </p>
    </div>
  );
}
