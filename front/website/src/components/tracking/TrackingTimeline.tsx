/** @jsxImportSource react */
import { useEffect, useMemo, useRef } from 'react';
import type { WheelEvent } from 'react';
import type { TrackingHistory } from '../../types/tracking';

interface TrackingTimelineProps {
  history?: TrackingHistory[];
  isLoading?: boolean;
  error?: string;
}

function IconMapPin() {
  return (
    <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
      <path stroke-linecap="round" stroke-linejoin="round" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
    </svg>
  );
}

function IconAlert() {
  return (
    <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
    </svg>
  );
}

function IconClock() {
  return (
    <svg class="w-8 h-8" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
      <path stroke-linecap="round" stroke-linejoin="round" d="M12 6v6l4 2m6-2a10 10 0 11-20 0 10 10 0 0120 0z" />
    </svg>
  );
}

function IconDocument() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M6 2a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8.828a2 2 0 00-.586-1.414l-4.828-4.828A2 2 0 0012.172 2H6z" />
    </svg>
  );
}

function IconCheck() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="3">
      <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
    </svg>
  );
}

function IconTruck() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M3 16V6a1 1 0 011-1h9v11M3 16h10m0 0h4m-4 0V9h4l3 3.5V16m0 0h-3m-9 0a2 2 0 11-4 0 2 2 0 014 0zm10 0a2 2 0 11-4 0 2 2 0 014 0z" />
    </svg>
  );
}

function IconFlag() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M5 21V4m0 0h11l-1.5 4L16 12H5" />
    </svg>
  );
}

function IconUndo() {
  return (
    <svg class="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
      <path stroke-linecap="round" stroke-linejoin="round" d="M9 15L4 10m0 0l5-5m-5 5h11a4 4 0 010 8h-1" />
    </svg>
  );
}

const STATUS_LABEL_ES: Record<string, string> = {
  pending: 'Pendiente',
  picked_up: 'Recogido',
  in_transit: 'En tránsito',
  out_for_delivery: 'En reparto',
  delivered: 'Entregado',
  on_hold: 'Con novedad',
  returned: 'Devuelto',
  cancelled: 'Cancelado',
  failed: 'Fallido',
};

const STATUS_ICON: Record<string, () => JSX.Element> = {
  pending: IconDocument,
  picked_up: IconCheck,
  in_transit: IconTruck,
  out_for_delivery: IconMapPin,
  delivered: IconFlag,
  on_hold: IconAlert,
  cancelled: IconAlert,
  failed: IconAlert,
  returned: IconUndo,
};

function translateStatus(status: string): string {
  return STATUS_LABEL_ES[status.trim().toLowerCase()] ?? status;
}

function getStatusIcon(status: string) {
  return STATUS_ICON[status.trim().toLowerCase()] ?? IconMapPin;
}

function dedupeEvents(events: TrackingHistory[]): TrackingHistory[] {
  const seen = new Set<string>();
  const result: TrackingHistory[] = [];
  for (const event of events) {
    const key = `${event.status}|${event.description}`;
    if (seen.has(key)) continue;
    seen.add(key);
    result.push(event);
  }
  return result;
}

function formatEventDate(raw: string): string {
  const d = new Date(raw);
  if (Number.isNaN(d.getTime())) return raw;
  const formatted = new Intl.DateTimeFormat('es-CO', {
    day: '2-digit',
    month: 'short',
    hour: '2-digit',
    minute: '2-digit',
    hour12: true,
    timeZone: 'America/Bogota',
  }).format(d);
  return formatted.replace('.', '');
}

export default function TrackingTimeline({ history, isLoading = false, error }: TrackingTimelineProps) {
  const scrollRef = useRef<HTMLDivElement>(null);
  const pauseUntilRef = useRef(0);

  const orderedHistory = useMemo(() => {
    if (!history || history.length === 0) return [];
    // El orden del API no siempre viene cronológico (algunos eventos llegan
    // fuera de secuencia), así que se ordena por fecha real: antiguo a la
    // izquierda, reciente a la derecha. Luego se quitan repeticiones del
    // mismo evento (misma transportadora suele reportar el mismo estado
    // varias veces).
    const sortedByDate = [...history].sort((a, b) => new Date(a.date).getTime() - new Date(b.date).getTime());
    return dedupeEvents(sortedByDate);
  }, [history]);

  useEffect(() => {
    const el = scrollRef.current;
    if (!el || orderedHistory.length <= 1) return;

    const SPEED_PX_PER_SEC = 40;
    const EDGE_HOLD_MS = 1200;
    let rafId: number;
    let lastTime: number | null = null;
    let direction: 1 | -1 = 1;
    let holdUntil = 0;

    const step = (time: number) => {
      if (lastTime === null) lastTime = time;
      const delta = time - lastTime;
      lastTime = time;

      const now = Date.now();
      if (now >= pauseUntilRef.current && now >= holdUntil) {
        const maxScroll = el.scrollWidth - el.clientWidth;
        if (maxScroll > 0) {
          const next = el.scrollLeft + direction * ((SPEED_PX_PER_SEC * delta) / 1000);
          el.scrollLeft = Math.min(maxScroll, Math.max(0, next));

          if (el.scrollLeft >= maxScroll - 0.5) {
            direction = -1;
            holdUntil = now + EDGE_HOLD_MS;
          } else if (el.scrollLeft <= 0.5) {
            direction = 1;
            holdUntil = now + EDGE_HOLD_MS;
          }
        }
      }

      rafId = requestAnimationFrame(step);
    };

    rafId = requestAnimationFrame(step);
    return () => cancelAnimationFrame(rafId);
  }, [orderedHistory.length]);

  const pauseAutoScroll = () => {
    pauseUntilRef.current = Date.now() + 4000;
  };

  const handleWheel = (e: WheelEvent<HTMLDivElement>) => {
    // Solo pausa con intención de scroll horizontal (trackpad/shift+rueda).
    // Un scroll vertical normal de la página no debe detener el avance.
    if (Math.abs(e.deltaX) > Math.abs(e.deltaY)) pauseAutoScroll();
  };

  if (isLoading) {
    return (
      <div class="flex items-center justify-center py-8 text-[#8B85A0]">
        <svg class="w-5 h-5 animate-spin mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
        </svg>
        <span>Cargando historial de rastreo...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div class="flex items-start gap-3 bg-[#F6F3FD] border border-[#E7E2F5] rounded-2xl p-5 text-[#7C3AED]">
        <IconAlert />
        <div class="text-sm">
          <p class="font-semibold">No se pudo cargar el historial</p>
          <p class="text-xs mt-1 text-[#8B85A0]">{error}</p>
        </div>
      </div>
    );
  }

  if (!history || history.length === 0) {
    return (
      <div class="text-center py-8 text-[#8B85A0]">
        <div class="flex justify-center mb-3 text-[#C4B5FD]">
          <IconClock />
        </div>
        <p class="text-sm">Aún no hay actualizaciones de rastreo disponibles.</p>
        <p class="text-xs text-[#B0A9C2] mt-1">El transportista pronto proporcionará información.</p>
      </div>
    );
  }

  return (
    <>
    <div class="space-y-4">
      <div class="flex items-center justify-between mb-2">
        <p class="font-space-grotesk text-sm font-bold text-[#181225] uppercase tracking-wider">Historial de eventos</p>
        {orderedHistory.length > 3 && (
          <p class="text-xs text-[#B0A9C2] hidden sm:block">Se mueve solo · desliza para retroceder</p>
        )}
      </div>

      <div
        ref={scrollRef}
        class="timeline-scroll overflow-x-auto -mx-1 px-1"
        onPointerDown={pauseAutoScroll}
        onWheel={handleWheel}
        onTouchStart={pauseAutoScroll}
      >
        <div class="flex items-start w-max">
          {orderedHistory.map((event, idx) => {
            const isLatest = idx === orderedHistory.length - 1;
            const StatusIcon = getStatusIcon(event.status);
            return (
              <div key={idx} class="contents">
                <div class="flex flex-col items-center w-[200px] sm:w-[228px] shrink-0 px-2">
                  <div
                    class={`w-8 h-8 rounded-full flex items-center justify-center ring-4 ring-white transition-all duration-300 ${
                      isLatest
                        ? 'bg-gradient-to-br from-[#A855F7] to-[#7C3AED] text-white scale-110 shadow-[0_8px_18px_-6px_rgba(124,58,237,0.6)]'
                        : 'bg-[#F0EAFD] text-[#B9AED0]'
                    }`}
                  >
                    <StatusIcon />
                  </div>

                  <div
                    class={`w-full mt-4 bg-white rounded-2xl border p-4 transition-shadow ${
                      isLatest
                        ? 'border-[#D8B4FE] shadow-[0_14px_30px_-16px_rgba(124,58,237,0.45)]'
                        : 'border-[#EFEAFB] hover:shadow-[0_10px_24px_-16px_rgba(124,58,237,0.4)]'
                    }`}
                  >
                    <p class={`font-bold text-sm mb-1 ${isLatest ? 'text-[#7C3AED]' : 'text-[#3A2E5C]'}`}>
                      {translateStatus(event.status)}
                    </p>
                    <p class="text-[11px] text-[#B0A9C2] mb-2">{formatEventDate(event.date)}</p>

                    {event.description && <p class="text-xs text-[#6B6480] mb-2 leading-snug">{event.description}</p>}

                    {event.location && (
                      <div class="flex items-center gap-1.5 text-xs text-[#8B85A0]">
                        <IconMapPin />
                        <span>{event.location}</span>
                      </div>
                    )}
                  </div>
                </div>

                {idx < orderedHistory.length - 1 && (
                  <div class="h-[3px] w-8 sm:w-12 shrink-0 mt-[18px] rounded-full bg-gradient-to-r from-[#7C3AED] to-[#A855F7]" />
                )}
              </div>
            );
          })}
        </div>
      </div>
    </div>
    <style>{`
      .timeline-scroll {
        -webkit-mask-image: linear-gradient(90deg, transparent 0, black 20px, black calc(100% - 20px), transparent 100%);
        mask-image: linear-gradient(90deg, transparent 0, black 20px, black calc(100% - 20px), transparent 100%);
      }
    `}</style>
    </>
  );
}
