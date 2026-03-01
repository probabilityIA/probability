/**
 * ConfigListTable - Configuraciones de notificación agrupadas por integración
 */
'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { NotificationConfig } from '../../domain/types';
import { getConfigsAction } from '../../infra/actions';
import { useIntegrationsSimple } from '@/services/integrations/core/ui/hooks/useIntegrationsSimple';
import { useToast } from '@/shared/providers/toast-provider';
import type { IntegrationSimple } from '@/services/integrations/core/domain/types';

interface EventSummary {
  eventName: string;
  channelName: string;
  channelCode: string;
  enabled: boolean;
}

interface IntegrationGroup {
  integration: IntegrationSimple;
  configs: NotificationConfig[];
  activeCount: number;
  channels: string[];
  events: EventSummary[];
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

    const buildEvents = (cfgs: NotificationConfig[]): EventSummary[] =>
      cfgs.map((c) => ({
        eventName: c.notification_event_name || `Evento #${c.notification_event_type_id}`,
        channelName: c.notification_type_name || '',
        channelCode: (c.notification_type_name || '').toLowerCase(),
        enabled: c.enabled,
      }));

    for (const [integrationId, groupConfigs] of groupMap) {
      const integration = integrations.find((i) => i.id === integrationId);
      const common = {
        configs: groupConfigs,
        activeCount: groupConfigs.filter((c) => c.enabled).length,
        channels: [...new Set(groupConfigs.map((c) => c.notification_type_name || '').filter(Boolean))],
        events: buildEvents(groupConfigs),
      };

      if (!integration) {
        result.push({
          ...common,
          integration: {
            id: integrationId,
            name: `Integración #${integrationId}`,
            type: 'unknown',
            category: '',
            category_name: '',
            business_id: selectedBusinessId || null,
            is_active: true,
          },
        });
        continue;
      }

      result.push({ ...common, integration });
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
        <button
          type="button"
          onClick={onCreate}
          className="p-2 rounded-lg bg-blue-50 text-blue-600 hover:bg-blue-100 transition-colors"
          title="Agregar Integración"
        >
          <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
        </button>
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
              className="flex items-start gap-4 px-4 py-3 hover:bg-gray-50 transition-colors"
            >
              {/* Integration info */}
              <div className="flex items-center gap-3 shrink-0 w-[180px] pt-0.5">
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

              {/* Events list */}
              <div className="flex-1 min-w-0">
                <div className="flex flex-wrap gap-1.5">
                  {group.events.map((ev, i) => {
                    const chCode = ev.channelCode;
                    const channelBg = chCode.includes('whatsapp') ? 'bg-green-100 text-green-700'
                      : chCode.includes('email') ? 'bg-orange-100 text-orange-700'
                      : chCode.includes('sms') ? 'bg-purple-100 text-purple-700'
                      : 'bg-blue-100 text-blue-700';
                    return (
                      <span
                        key={i}
                        className={`inline-flex items-center gap-1.5 px-2 py-1 rounded-md text-[11px] font-medium border ${
                          ev.enabled
                            ? `${channelBg} border-transparent`
                            : 'bg-gray-50 text-gray-400 border-gray-200 line-through'
                        }`}
                      >
                        <span className="font-semibold">{ev.channelName}</span>
                        <span className="text-[9px] opacity-60">/</span>
                        {ev.eventName}
                      </span>
                    );
                  })}
                </div>
              </div>

              {/* Count summary */}
              <div className="text-right shrink-0 pt-0.5">
                <p className="text-xs text-gray-500">
                  {group.activeCount}/{group.configs.length} activas
                </p>
              </div>

              {/* Configure icon */}
              <div className="shrink-0 pt-0.5">
                <button
                  type="button"
                  onClick={() => onConfigure(group.integration)}
                  className="p-1.5 rounded-md text-gray-400 hover:text-blue-600 hover:bg-blue-50 transition-colors"
                  title="Configurar"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.066 2.573c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.573 1.066c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.066-2.573c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
