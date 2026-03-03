'use client';

import { useState, useEffect, useCallback } from 'react';
import { SummaryStatsRow } from './SummaryStatsRow';
import { ConfigListTable } from './ConfigListTable';
import { MessageAudit } from './MessageAudit';
import { IntegrationPicker } from './IntegrationPicker';
import { IntegrationRulesForm } from './IntegrationRulesForm';
import { Modal } from '@/shared/ui/modal';
import {
  getConfigsAction,
  getNotificationTypesAction,
  getNotificationEventTypesAction,
} from '../../infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNotificationBusiness } from '@/shared/contexts/notification-business-context';
import type { IntegrationSimple } from '@/services/integrations/core/domain/types';

interface Stats {
  integrationCount: number;
  activeRulesCount: number;
  channelCount: number;
  eventTypeCount: number;
}

export function NotificationDashboard() {
  const { isSuperAdmin } = usePermissions();
  const { selectedBusinessId } = useNotificationBusiness();

  const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

  // Stats
  const [stats, setStats] = useState<Stats>({
    integrationCount: 0,
    activeRulesCount: 0,
    channelCount: 0,
    eventTypeCount: 0,
  });
  const [statsLoading, setStatsLoading] = useState(true);

  // Config flow: picker -> rules form
  const [isPickerModalOpen, setIsPickerModalOpen] = useState(false);
  const [isRulesModalOpen, setIsRulesModalOpen] = useState(false);
  const [selectedIntegration, setSelectedIntegration] = useState<IntegrationSimple | undefined>(undefined);
  const [configRefreshKey, setConfigRefreshKey] = useState(0);

  // Reset on business change
  useEffect(() => {
    setConfigRefreshKey((prev) => prev + 1);
    setSelectedIntegration(undefined);
    setIsPickerModalOpen(false);
    setIsRulesModalOpen(false);
  }, [selectedBusinessId]);

  // Fetch stats in parallel
  const fetchStats = useCallback(async () => {
    if (requiresBusinessSelection) return;
    setStatsLoading(true);
    try {
      const [configsRes, typesRes, eventsRes] = await Promise.all([
        getConfigsAction({
          ...(selectedBusinessId ? { business_id: selectedBusinessId } : {}),
        }),
        getNotificationTypesAction(),
        getNotificationEventTypesAction(),
      ]);

      const configs = configsRes.data || [];
      const uniqueIntegrations = new Set(configs.map((c) => c.integration_id));

      setStats({
        integrationCount: uniqueIntegrations.size,
        activeRulesCount: configs.filter((c) => c.enabled).length,
        channelCount: typesRes.success ? typesRes.data.length : 0,
        eventTypeCount: eventsRes.success ? eventsRes.data.length : 0,
      });
    } catch {
      // Stats are non-critical, keep defaults
    } finally {
      setStatsLoading(false);
    }
  }, [requiresBusinessSelection, selectedBusinessId]);

  useEffect(() => {
    fetchStats();
  }, [fetchStats, configRefreshKey]);

  // --- Config handlers ---
  const handleCreateConfig = () => setIsPickerModalOpen(true);

  const handlePickIntegration = (integration: IntegrationSimple) => {
    setSelectedIntegration(integration);
    setIsPickerModalOpen(false);
    setIsRulesModalOpen(true);
  };

  const handleConfigureIntegration = (integration: IntegrationSimple) => {
    setSelectedIntegration(integration);
    setIsRulesModalOpen(true);
  };

  const handleRulesSuccess = () => {
    setIsRulesModalOpen(false);
    setSelectedIntegration(undefined);
    setConfigRefreshKey((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-50 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900">Centro de Notificaciones</h1>
      </div>

      {/* Gate: require business selection */}
      {requiresBusinessSelection ? (
        <div className="text-center py-16 text-gray-500">
          Selecciona un negocio para ver las configuraciones de notificación
        </div>
      ) : (
        <div className="space-y-8">
          {/* Summary Stats */}
          <SummaryStatsRow
            integrationCount={stats.integrationCount}
            activeRulesCount={stats.activeRulesCount}
            channelCount={stats.channelCount}
            eventTypeCount={stats.eventTypeCount}
            loading={statsLoading}
          />

          {/* Main section: Rules by Integration */}
          <ConfigListTable
            onConfigure={handleConfigureIntegration}
            onCreate={handleCreateConfig}
            refreshKey={configRefreshKey}
            selectedBusinessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
          />

          {/* Message Audit */}
          <MessageAudit businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined} />

          {/* --- Modals --- */}

          {/* Integration Picker */}
          <Modal
            isOpen={isPickerModalOpen}
            onClose={() => setIsPickerModalOpen(false)}
            title="Seleccionar Integración"
          >
            <IntegrationPicker
              businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined}
              onSelect={handlePickIntegration}
              onCancel={() => setIsPickerModalOpen(false)}
            />
          </Modal>

          {/* Integration Rules Form */}
          <Modal
            isOpen={isRulesModalOpen}
            onClose={() => setIsRulesModalOpen(false)}
            title={`Reglas de Notificación — ${selectedIntegration?.name || ''}`}
            size="5xl"
          >
            {selectedIntegration && (
              <IntegrationRulesForm
                integration={selectedIntegration}
                businessId={isSuperAdmin ? (selectedBusinessId ?? 0) : 0}
                onSuccess={handleRulesSuccess}
                onCancel={() => setIsRulesModalOpen(false)}
              />
            )}
          </Modal>

        </div>
      )}
    </div>
  );
}
