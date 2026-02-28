/**
 * ConfigListTable - Configuraciones de notificación agrupadas por integración
 */
'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { Badge, Button } from '@/shared/ui';
import { NotificationConfig } from '../../domain/types';
import { getConfigsAction } from '../../infra/actions';
import { useIntegrationsSimple } from '@/services/integrations/core/ui/hooks/useIntegrationsSimple';
import { useToast } from '@/shared/providers/toast-provider';
import type { IntegrationSimple } from '@/services/integrations/core/domain/types';

interface IntegrationGroup {
  integration: IntegrationSimple;
  configs: NotificationConfig[];
  activeCount: number;
  channels: string[];
}

interface ConfigListTableProps {
  onConfigure: (integration: IntegrationSimple) => void;
  onCreate: () => void;
  refreshKey?: number;
  selectedBusinessId?: number;
}

export function ConfigListTable({ onConfigure, onCreate, refreshKey = 0, selectedBusinessId }: ConfigListTableProps) {
  const { showToast } = useToast();
  const { integrations, loading: loadingIntegrations } = useIntegrationsSimple(
    selectedBusinessId ? { businessId: selectedBusinessId } : undefined
  );

  const [configs, setConfigs] = useState<NotificationConfig[]>([]);
  const [loading, setLoading] = useState(true);

  // Fetch all configs
  const fetchConfigs = useCallback(async () => {
    setLoading(true);
    try {
      const response = await getConfigsAction({
        ...(selectedBusinessId ? { business_id: selectedBusinessId } : {}),
      });
      setConfigs(response.data || []);
    } catch (error) {
      showToast('Error al cargar configuraciones', 'error');
    } finally {
      setLoading(false);
    }
  }, [selectedBusinessId, showToast]);

  useEffect(() => {
    fetchConfigs();
  }, [fetchConfigs, refreshKey]);

  // Group configs by integration_id and enrich with integration info
  const groups: IntegrationGroup[] = useMemo(() => {
    const groupMap = new Map<number, NotificationConfig[]>();

    for (const config of configs) {
      const key = config.integration_id;
      if (!groupMap.has(key)) groupMap.set(key, []);
      groupMap.get(key)!.push(config);
    }

    const result: IntegrationGroup[] = [];

    for (const [integrationId, groupConfigs] of groupMap) {
      const integration = integrations.find((i) => i.id === integrationId);
      if (!integration) {
        // Create a fallback integration object
        result.push({
          integration: {
            id: integrationId,
            name: `Integración #${integrationId}`,
            type: 'unknown',
            category: '',
            category_name: '',
            business_id: selectedBusinessId || null,
            is_active: true,
          },
          configs: groupConfigs,
          activeCount: groupConfigs.filter((c) => c.enabled).length,
          channels: [...new Set(groupConfigs.map((c) => c.notification_type_name || '').filter(Boolean))],
        });
        continue;
      }

      result.push({
        integration,
        configs: groupConfigs,
        activeCount: groupConfigs.filter((c) => c.enabled).length,
        channels: [...new Set(groupConfigs.map((c) => c.notification_type_name || '').filter(Boolean))],
      });
    }

    return result.sort((a, b) => a.integration.name.localeCompare(b.integration.name));
  }, [configs, integrations, selectedBusinessId]);

  const isLoading = loading || loadingIntegrations;

  return (
    <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
      {/* Header */}
      <div className="p-4 border-b flex items-center justify-between">
        <div>
          <h3 className="text-sm font-medium text-gray-900">Reglas por Integración</h3>
          <p className="text-xs text-gray-500 mt-0.5">
            {groups.length} integración(es) con reglas configuradas
          </p>
        </div>
        <Button onClick={onCreate} size="sm">
          <svg className="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Agregar Integración
        </Button>
      </div>

      {/* Content */}
      {isLoading ? (
        <div className="text-center py-12 text-gray-500">Cargando...</div>
      ) : groups.length === 0 ? (
        <div className="text-center py-12">
          <div className="text-gray-400 mb-2">
            <svg className="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
            </svg>
          </div>
          <p className="text-gray-500 text-sm">No hay reglas de notificación configuradas</p>
          <p className="text-gray-400 text-xs mt-1">Haz clic en &quot;Agregar Integración&quot; para empezar</p>
        </div>
      ) : (
        <div className="divide-y divide-gray-100">
          {groups.map((group) => (
            <div
              key={group.integration.id}
              className="flex items-center justify-between px-4 py-3 hover:bg-gray-50 transition-colors"
            >
              {/* Integration info */}
              <div className="flex items-center gap-3 min-w-0 flex-1">
                {group.integration.image_url ? (
                  <img
                    src={group.integration.image_url}
                    alt={group.integration.name}
                    className="w-8 h-8 object-contain rounded shrink-0"
                  />
                ) : (
                  <div className="w-8 h-8 rounded bg-gray-100 flex items-center justify-center shrink-0">
                    <span className="text-xs font-bold text-gray-400">
                      {group.integration.type?.charAt(0).toUpperCase() || "?"}
                    </span>
                  </div>
                )}
                <div className="min-w-0">
                  <p className="text-sm font-medium text-gray-900 truncate">
                    {group.integration.name}
                  </p>
                  <p className="text-xs text-gray-500">
                    {group.integration.category_name || group.integration.type}
                  </p>
                </div>
              </div>

              {/* Rules count */}
              <div className="text-center px-4 shrink-0">
                <p className="text-sm font-medium text-gray-900">
                  {group.configs.length} regla(s)
                </p>
                <p className="text-xs text-gray-500">
                  {group.activeCount} activa(s)
                </p>
              </div>

              {/* Channels badges */}
              <div className="flex flex-wrap gap-1 px-4 shrink-0 max-w-48">
                {group.channels.map((channel) => {
                  const name = channel.toLowerCase();
                  const badgeType = name.includes('whatsapp') ? 'success' as const
                    : name.includes('email') ? 'warning' as const
                    : name.includes('sms') ? 'primary' as const
                    : 'secondary' as const;
                  return (
                    <Badge key={channel} type={badgeType}>
                      {channel}
                    </Badge>
                  );
                })}
              </div>

              {/* Actions */}
              <div className="shrink-0">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => onConfigure(group.integration)}
                >
                  Configurar
                </Button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
