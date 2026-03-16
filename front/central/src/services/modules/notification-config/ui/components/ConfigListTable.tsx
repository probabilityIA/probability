/**
 * ConfigListTable - Configuraciones de notificación agrupadas por integración
 */
'use client';

import { useState, useEffect, useCallback, useMemo } from 'react';
import { NotificationConfig } from '../../domain/types';
import { getConfigsAction, testIntegrationConnectionAction } from '../../infra/actions';
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

      // Skip configs for inactive/unknown integrations
      if (!integration) continue;

      result.push({
        configs: groupConfigs,
        activeCount: groupConfigs.filter((c) => c.enabled).length,
        channels: [...new Set(groupConfigs.map((c) => c.notification_type_name || '').filter(Boolean))],
        events: buildEvents(groupConfigs),
        integration,
      });
    }

    return result.sort((a, b) => a.integration.name.localeCompare(b.integration.name));
  }, [configs, integrations, selectedBusinessId]);

  // WhatsApp integration for testing
  const whatsAppIntegration = useMemo(
    () => integrations.find((i) => i.category === 'messaging' && i.is_active),
    [integrations]
  );

  const [testingWhatsApp, setTestingWhatsApp] = useState(false);

  const handleTestWhatsApp = async () => {
    if (!whatsAppIntegration) {
      showToast('No se encontró integración WhatsApp activa', 'error');
      return;
    }
    setTestingWhatsApp(true);
    try {
      const result = await testIntegrationConnectionAction(whatsAppIntegration.id);
      if (result.success) {
        showToast(result.message || 'Mensaje de prueba enviado', 'success');
      } else {
        showToast(result.error || 'Error al probar conexión', 'error');
      }
    } catch {
      showToast('Error inesperado al probar conexión', 'error');
    } finally {
      setTestingWhatsApp(false);
    }
  };

  const isLoading = loading || loadingIntegrations;

  // Check if any group has WhatsApp rules
  const hasWhatsAppRules = groups.some((g) =>
    g.events.some((e) => e.channelCode.includes('whatsapp'))
  );

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
        <div className="flex items-center gap-2">
          {hasWhatsAppRules && whatsAppIntegration && (
            <button
              type="button"
              onClick={handleTestWhatsApp}
              disabled={testingWhatsApp}
              className="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-green-50 text-green-700 hover:bg-green-100 transition-colors text-xs font-medium disabled:opacity-50"
              title="Enviar mensaje de prueba WhatsApp"
            >
              {testingWhatsApp ? (
                <svg className="w-3.5 h-3.5 animate-spin" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                </svg>
              ) : (
                <svg className="w-3.5 h-3.5" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347z" />
                  <path d="M12 0C5.373 0 0 5.373 0 12c0 2.625.846 5.059 2.284 7.034L.789 23.492a.5.5 0 00.611.611l4.458-1.495A11.934 11.934 0 0012 24c6.627 0 12-5.373 12-12S18.627 0 12 0zm0 22c-2.287 0-4.405-.744-6.122-2.003l-.427-.32-2.645.887.887-2.645-.32-.427A9.935 9.935 0 012 12C2 6.486 6.486 2 12 2s10 4.486 10 10-4.486 10-10 10z" />
                </svg>
              )}
              Probar WhatsApp
            </button>
          )}
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
