'use client';

interface SummaryStatsRowProps {
  integrationCount: number;
  activeRulesCount: number;
  channelCount: number;
  eventTypeCount: number;
  loading?: boolean;
}

const stats = [
  {
    key: 'integrations',
    label: 'Integraciones',
    colorClasses: 'bg-purple-50 text-purple-700',
    iconColor: 'text-purple-500',
    icon: (
      <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
      </svg>
    ),
  },
  {
    key: 'activeRules',
    label: 'Reglas Activas',
    colorClasses: 'bg-green-50 text-green-700',
    iconColor: 'text-green-500',
    icon: (
      <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
    ),
  },
  {
    key: 'channels',
    label: 'Canales',
    colorClasses: 'bg-blue-50 text-blue-700',
    iconColor: 'text-blue-500',
    icon: (
      <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
      </svg>
    ),
  },
  {
    key: 'eventTypes',
    label: 'Eventos',
    colorClasses: 'bg-orange-50 text-orange-700',
    iconColor: 'text-orange-500',
    icon: (
      <svg className="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13 10V3L4 14h7v7l9-11h-7z" />
      </svg>
    ),
  },
] as const;

export function SummaryStatsRow({
  integrationCount,
  activeRulesCount,
  channelCount,
  eventTypeCount,
  loading = false,
}: SummaryStatsRowProps) {
  const values: Record<string, number> = {
    integrations: integrationCount,
    activeRules: activeRulesCount,
    channels: channelCount,
    eventTypes: eventTypeCount,
  };

  return (
    <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
      {stats.map((stat) => (
        <div
          key={stat.key}
          className={`rounded-xl p-4 ${stat.colorClasses} border border-opacity-20`}
        >
          <div className="flex items-center justify-between">
            <div>
              <p className="text-xs font-medium opacity-75">{stat.label}</p>
              {loading ? (
                <div className="h-8 w-12 bg-current opacity-10 rounded mt-1 animate-pulse" />
              ) : (
                <p className="text-2xl font-bold mt-1">{values[stat.key]}</p>
              )}
            </div>
            <div className={`${stat.iconColor} opacity-60`}>{stat.icon}</div>
          </div>
        </div>
      ))}
    </div>
  );
}
