'use client';

import { useState, useEffect } from 'react';
import { createPortal } from 'react-dom';
import { MessageAudit } from './MessageAudit';
import { WhatsAppConversations } from './WhatsAppConversations';
import { IntegrationRulesForm, RULES_TABS_SLOT_ID } from './IntegrationRulesForm';
import { CampaignsSection } from './CampaignsSection';
import { NOTIFICATION_STATS_REFRESH_EVENT } from './NotificationSummaryKpis';
import { NOTIFICATIONS_TABS_SLOT_ID } from '@/shared/ui/notifications-subnavbar';
import { Modal } from '@/shared/ui/modal';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNotificationBusiness } from '@/shared/contexts/notification-business-context';

const tabs = [
  { key: 'conversations' as const, label: 'Conversaciones' },
  { key: 'audit' as const, label: 'Auditoria' },
  { key: 'campaigns' as const, label: 'Campañas' },
];

export function NotificationDashboard() {
  const { isSuperAdmin } = usePermissions();
  const { selectedBusinessId } = useNotificationBusiness();

  const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

  const [activeTab, setActiveTab] = useState<'audit' | 'conversations' | 'campaigns'>('conversations');

  const [isRulesModalOpen, setIsRulesModalOpen] = useState(false);
  const [configRefreshKey, setConfigRefreshKey] = useState(0);
  const [tabsSlot, setTabsSlot] = useState<HTMLElement | null>(null);

  useEffect(() => {
    setTabsSlot(document.getElementById(NOTIFICATIONS_TABS_SLOT_ID));
  }, []);

  useEffect(() => {
    setConfigRefreshKey((prev) => prev + 1);
    setIsRulesModalOpen(false);
  }, [selectedBusinessId]);

  const handleRulesSuccess = () => {
    setIsRulesModalOpen(false);
    setConfigRefreshKey((prev) => prev + 1);
    window.dispatchEvent(new Event(NOTIFICATION_STATS_REFRESH_EVENT));
  };

  const tabsBar = (
    <div className="flex items-center justify-between gap-3">
      <nav className="flex items-center gap-1.5" aria-label="Tabs">
        {tabs.map((tab) => (
          <button
            key={tab.key}
            onClick={() => setActiveTab(tab.key)}
            style={
              activeTab === tab.key
                ? { backgroundColor: 'var(--color-secondary-500)', color: 'var(--color-on-secondary, white)' }
                : {}
            }
            className={`shrink-0 rounded-lg px-3 py-2 text-sm font-medium whitespace-nowrap transition-all ${
              activeTab === tab.key
                ? ''
                : 'text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-700 hover:text-gray-900 dark:hover:text-gray-100'
            }`}
          >
            {tab.label}
          </button>
        ))}
      </nav>

      <button
        type="button"
        onClick={() => setIsRulesModalOpen(true)}
        className="flex items-center gap-2 rounded-lg px-3 py-1.5 text-sm font-medium transition-opacity hover:opacity-90"
        style={{ backgroundColor: 'var(--color-primary)', color: 'var(--color-on-primary, white)' }}
      >
        <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h7" />
        </svg>
        Reglas
      </button>
    </div>
  );

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6">
      {requiresBusinessSelection ? (
        <div className="text-center py-16 text-gray-500 dark:text-gray-400">
          {'Selecciona un negocio para ver las configuraciones de notificación'}
        </div>
      ) : (
        <div key={configRefreshKey} className="space-y-6">
          {tabsSlot && createPortal(tabsBar, tabsSlot)}

          {activeTab === 'audit' && (
            <MessageAudit businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined} />
          )}

          {activeTab === 'conversations' && (
            <WhatsAppConversations businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined} />
          )}

          {activeTab === 'campaigns' && (
            <CampaignsSection businessId={isSuperAdmin ? selectedBusinessId ?? undefined : undefined} />
          )}

          <Modal
            isOpen={isRulesModalOpen}
            onClose={() => setIsRulesModalOpen(false)}
            title={(
              <span className="flex w-full items-center justify-between gap-4 pr-8">
                <span className="text-lg font-semibold">{'Reglas de Notificación'}</span>
                <span id={RULES_TABS_SLOT_ID} className="inline-flex items-center empty:hidden" />
              </span>
            )}
            size="4xl"
            noPadding
            noBodyScroll
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
