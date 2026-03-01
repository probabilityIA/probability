'use client';

import type { MessageAuditStats } from '../../domain/types';

interface MessageAuditStatsRowProps {
  stats: MessageAuditStats | null;
  loading?: boolean;
}

export function MessageAuditStatsRow({ stats, loading = false }: MessageAuditStatsRowProps) {
  const totalSent = stats
    ? stats.total_sent + stats.total_delivered + stats.total_read
    : 0;

  const cards = [
    {
      label: 'Enviados',
      value: totalSent,
      colorClasses: 'bg-green-50 text-green-700',
      iconColor: 'text-green-500',
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8" />
        </svg>
      ),
    },
    {
      label: 'Fallidos',
      value: stats?.total_failed ?? 0,
      colorClasses: 'bg-red-50 text-red-700',
      iconColor: 'text-red-500',
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L4.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>
      ),
    },
    {
      label: 'Tasa de Exito',
      value: stats ? `${stats.success_rate.toFixed(1)}%` : '0%',
      colorClasses: 'bg-blue-50 text-blue-700',
      iconColor: 'text-blue-500',
      icon: (
        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
        </svg>
      ),
    },
  ];

  return (
    <div className="grid grid-cols-3 gap-3">
      {cards.map((card) => (
        <div
          key={card.label}
          className={`rounded-lg p-3 ${card.colorClasses} border border-opacity-20`}
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium opacity-75">{card.label}</p>
              {loading ? (
                <div className="h-7 w-10 bg-current opacity-10 rounded mt-1 animate-pulse" />
              ) : (
                <p className="text-xl font-bold mt-1">{card.value}</p>
              )}
            </div>
            <div className={`${card.iconColor} opacity-60`}>{card.icon}</div>
          </div>
        </div>
      ))}
    </div>
  );
}
