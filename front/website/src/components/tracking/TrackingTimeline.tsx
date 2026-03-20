import type { TrackingHistory } from '../../types/tracking';

interface TrackingTimelineProps {
  history?: TrackingHistory[];
  isLoading?: boolean;
  error?: string;
}

export default function TrackingTimeline({ history, isLoading = false, error }: TrackingTimelineProps) {
  if (isLoading) {
    return (
      <div class="flex items-center justify-center py-8 text-gray-500">
        <svg class="w-5 h-5 animate-spin mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
        <span>Cargando historial de rastreo...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div class="flex items-start gap-3 bg-amber-50 border border-amber-200 rounded-lg p-4 text-amber-700">
        <span class="text-2xl">⚠️</span>
        <div class="text-sm">
          <p class="font-semibold">No se pudo cargar el historial</p>
          <p class="text-xs mt-1">{error}</p>
        </div>
      </div>
    );
  }

  if (!history || history.length === 0) {
    return (
      <div class="text-center py-8 text-gray-500">
        <div class="text-4xl mb-2">⏱️</div>
        <p class="text-sm">Aún no hay actualizaciones de rastreo disponibles.</p>
        <p class="text-xs text-gray-400 mt-1">
          El transportista pronto proporcionará información.
        </p>
      </div>
    );
  }

  return (
    <div class="space-y-4">
      <p class="text-sm font-bold text-gray-600 uppercase tracking-wider mb-4">
        Historial de Eventos
      </p>

      <div class="relative">
        {/* Gradient line */}
        <div class="absolute left-5 top-0 bottom-0 w-1 bg-gradient-to-b from-blue-400 via-blue-300 to-gray-200 rounded-full" />

        {/* Timeline events */}
        <div class="space-y-0 pl-16">
          {history.map((event, idx) => {
            const isFirst = idx === 0;
            return (
              <div
                key={idx}
                class="relative pb-6"
                style={{ animationDelay: `${idx * 50}ms`, animationFillMode: 'both' }}
              >
                {/* Circle dot */}
                <div
                  class={`
                    absolute -left-12 top-1 w-5 h-5 rounded-full ring-2 ring-white
                    transition-all duration-300
                    ${isFirst ? 'bg-blue-500 scale-110 shadow-md shadow-blue-500/50' : 'bg-gray-300'}
                  `}
                />

                {/* Content */}
                <div class="bg-white rounded-lg border border-gray-200 p-4 hover:shadow-md transition-shadow">
                  <div class="flex items-baseline justify-between gap-2 mb-1">
                    <p class={`font-bold text-sm ${isFirst ? 'text-blue-700' : 'text-gray-800'}`}>
                      {event.status}
                    </p>
                    <p class="text-xs text-gray-400 flex-shrink-0">{event.date}</p>
                  </div>

                  {event.description && (
                    <p class="text-sm text-gray-600 mb-2">{event.description}</p>
                  )}

                  {event.location && (
                    <div class="flex items-center gap-2 text-xs text-gray-500">
                      <span>📍</span>
                      <span>{event.location}</span>
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
