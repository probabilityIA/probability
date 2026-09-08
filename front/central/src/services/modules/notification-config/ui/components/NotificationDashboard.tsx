'use client';

import { useState, useEffect, useCallback } from 'react';
import { SummaryStatsRow } from './SummaryStatsRow';
import { MessageAudit } from './MessageAudit';
import { WhatsAppConversations } from './WhatsAppConversations';
import { IntegrationRulesForm } from './IntegrationRulesForm';
import { CampaignsSection } from './CampaignsSection';
import { Modal } from '@/shared/ui/modal';
import {
  getConfigsAction,
  getNotificationTypesAction,
  getNotificationEventTypesAction,
} from '../../infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNotificationBusiness } from '@/shared/contexts/notification-business-context';

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

  // Tab state
  const [activeTab, setActiveTab] = useState<'audit' | 'conversations' | 'campaigns'>('audit');

  const [isRulesModalOpen, setIsRulesModalOpen] = useState(false);
  const [configRefreshKey, setConfigRefreshKey] = useState(0);

  // Reset on business change
  useEffect(() => {
    setConfigRefreshKey((prev) => prev + 1);
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

  const handleRulesSuccess = () => {
    setIsRulesModalOpen(false);
    setConfigRefreshKey((prev) => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
      {/* Header */}
      <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
        <h1 className="text-2xl sm:text-3xl font-bold text-gray-900 dark:text-white">Centro de Notificaciones</h1>
      </div>

      {/* Gate: require business selection */}
      {requiresBusinessSelection ? (
        <div className="text-center py-16 text-gray-500 dark:text-gray-400">
          Selecciona un negocio para ver las configuraciones de notificación
        </div>
      ) : (
        <div className="space-y-6">
          {/* Summary Stats */}
          <SummaryStatsRow
            integrationCount={stats.integrationCount}
            activeRulesCount={stats.activeRulesCount}
            channelCount={stats.channelCount}
            eventTypeCount={stats.eventTypeCount}
            loading={statsLoading}
          />

          {/* Tabs */}
          <div className="flex items-end justify-between border-b border-gray-200 dark:border-gray-700">
            <nav className="flex gap-6" aria-label="Tabs">
              {([
                { key: 'audit' as const, label: 'Auditoria' },
                { key: 'conversations' as const, label: 'Conversaciones' },
                { key: 'campaigns' as const, label: 'Campañas' },
              ]).map((tab) => (
                <button
                  key={tab.key}
                  onClick={() => setActiveTab(tab.key)}
                  className={`pb-3 text-sm font-medium border-b-2 transition-colors ${
                    activeTab === tab.key
                      ? 'border-[var(--color-primary)] text-[var(--color-primary)]'
                      : 'border-transparent text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300 hover:border-gray-300'
                  }`}
                >
                  {tab.label}
                </button>
              ))}
            </nav>

            <button
              type="button"
              onClick={() => setIsRulesModalOpen(true)}
              className="mb-2 flex items-center gap-2 rounded-lg px-3 py-1.5 text-sm font-medium text-white transition-opacity hover:opacity-90 disabled:opacity-50"
              style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}
            >
              <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h7" />
              </svg>
              Reglas
            </button>
          </div>

          {activeTab === 'audit' && (
            <MessageAudit businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined} />
          )}

          {activeTab === 'conversations' && (
            <WhatsAppConversations businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined} />
          )}

          {activeTab === 'campaigns' && (
            <CampaignsSection businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined} />
          )}

          {/* Integration Rules Form */}
          <Modal
            isOpen={isRulesModalOpen}
            onClose={() => setIsRulesModalOpen(false)}
            title="Reglas de Notificación"
            size="4xl"
          >
            {isRulesModalOpen && (
              <IntegrationRulesForm
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
