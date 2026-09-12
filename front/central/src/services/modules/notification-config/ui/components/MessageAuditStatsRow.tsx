'use client';

import { useEffect, useState } from 'react';
import { createPortal } from 'react-dom';
import { NOTIFICATIONS_STATS_SLOT_ID } from '@/shared/ui/notifications-subnavbar';
import type { MessageAuditStats } from '../../domain/types';

interface MessageAuditStatsRowProps {
  stats: MessageAuditStats | null;
  loading?: boolean;
}

export function MessageAuditStatsRow({ stats, loading = false }: MessageAuditStatsRowProps) {
  const [slot, setSlot] = useState<HTMLElement | null>(null);

  useEffect(() => {
    setSlot(document.getElementById(NOTIFICATIONS_STATS_SLOT_ID));
  }, []);

  if (!slot) return null;

  const totalSent = stats
    ? stats.total_sent + stats.total_delivered + stats.total_read
    : 0;

  const cards = [
    {
      label: 'Enviados',
      value: totalSent,
      dot: '#22c55e',
      text: 'text-emerald-600 dark:text-emerald-400',
    },
    {
      label: 'Fallidos',
      value: stats?.total_failed ?? 0,
      dot: '#ef4444',
      text: 'text-red-600 dark:text-red-400',
    },
    {
      label: 'Tasa de éxito',
      value: stats ? `${stats.success_rate.toFixed(1)}%` : '0%',
      dot: '#a855f7',
      text: 'text-purple-600 dark:text-purple-400',
    },
  ];

  return createPortal(
    <div className="hidden xl:flex items-center gap-3 border-l border-gray-200 dark:border-gray-700 pl-3">
      {cards.map((card) => (
        <div key={card.label} className="flex items-center gap-2 leading-tight">
          <span className="h-1.5 w-1.5 flex-shrink-0 rounded-full" style={{ backgroundColor: card.dot }} />
          <div>
            <p className="text-[8px] font-semibold uppercase tracking-[0.12em] text-gray-400 dark:text-gray-500">
              {card.label}
            </p>
            {loading ? (
              <div className="mt-0.5 h-3 w-8 animate-pulse rounded bg-gray-200 dark:bg-gray-700" />
            ) : (
              <p className={`text-sm font-bold tabular-nums ${card.text}`}>{card.value}</p>
            )}
          </div>
        </div>
      ))}
    </div>,
    slot,
  );
}
