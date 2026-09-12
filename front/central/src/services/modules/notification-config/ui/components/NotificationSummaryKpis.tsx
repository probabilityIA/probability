'use client';

import { useCallback, useEffect, useState } from 'react';
import {
  getConfigsAction,
  getNotificationTypesAction,
  getNotificationEventTypesAction,
} from '../../infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useNotificationBusiness } from '@/shared/contexts/notification-business-context';

export const NOTIFICATION_STATS_REFRESH_EVENT = 'notification-stats-refresh';

interface Stats {
  integrations: number;
  activeRules: number;
  channels: number;
  eventTypes: number;
}

const cards = [
  { key: 'integrations' as const, label: 'Integraciones', dot: '#a855f7', text: 'text-purple-600 dark:text-purple-400' },
  { key: 'activeRules' as const, label: 'Reglas activas', dot: '#22c55e', text: 'text-emerald-600 dark:text-emerald-400' },
  { key: 'channels' as const, label: 'Canales', dot: '#6366f1', text: 'text-indigo-600 dark:text-indigo-400' },
  { key: 'eventTypes' as const, label: 'Eventos', dot: '#f59e0b', text: 'text-amber-600 dark:text-amber-400' },
];

export function NotificationSummaryKpis() {
  const { isSuperAdmin } = usePermissions();
  const { selectedBusinessId } = useNotificationBusiness();

  const [stats, setStats] = useState<Stats | null>(null);

  const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

  const load = useCallback(async () => {
    if (requiresBusinessSelection) {
      setStats(null);
      return;
    }

    try {
      const [configsRes, typesRes, eventsRes] = await Promise.all([
        getConfigsAction({
          ...(selectedBusinessId ? { business_id: selectedBusinessId } : {}),
        }),
        getNotificationTypesAction(),
        getNotificationEventTypesAction(),
      ]);

      const configs = configsRes.data || [];

      setStats({
        integrations: new Set(configs.map((c) => c.integration_id)).size,
        activeRules: configs.filter((c) => c.enabled).length,
        channels: typesRes.success ? typesRes.data.length : 0,
        eventTypes: eventsRes.success ? eventsRes.data.length : 0,
      });
    } catch {
      setStats(null);
    }
  }, [requiresBusinessSelection, selectedBusinessId]);

  useEffect(() => {
    load();
    window.addEventListener(NOTIFICATION_STATS_REFRESH_EVENT, load);
    return () => window.removeEventListener(NOTIFICATION_STATS_REFRESH_EVENT, load);
  }, [load]);

  if (!stats) return null;

  return (
    <div className="hidden xl:flex items-center gap-3">
      {cards.map((card) => (
        <div key={card.key} className="flex items-center gap-2 leading-tight">
          <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full" style={{ backgroundColor: card.dot }} />
          <div>
            <p className="text-[8px] font-semibold uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
              {card.label}
            </p>
            <p className={`text-sm font-bold tabular-nums ${card.text}`}>{stats[card.key]}</p>
          </div>
        </div>
      ))}
    </div>
  );
}
