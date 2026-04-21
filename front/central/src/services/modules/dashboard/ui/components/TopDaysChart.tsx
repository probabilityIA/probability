'use client';

import React, { useState, useEffect } from 'react';
import { OrdersByDate } from '../../domain/types';

interface TopDaysChartProps {
  data: OrdersByDate[];
  dateRange?: { from: string; to: string };
}

/**
 * TopDaysChart - Componente agnóstico de presentación
 *
 * NOTA: Los cálculos (orden, top 5, opacidad, patrón de días) deben hacerse en el BACKEND
 * para optimizar rendimiento en alta concurrencia.
 *
 * El backend debe retornar datos ya enriquecidos con:
 * - heightPercent: (count / maxCount) * 100
 * - opacity: valor según posición
 * - dayName: nombre completo del día
 * - dayShort: abreviado
 * - rank: posición en top 5
 * - maxCount: máximo del período
 * - weekdayCount: cuántas veces aparece en top 5
 *
 * Este componente SOLO renderiza los datos sin hacer cálculos.
 */
export function TopDaysChart({ data, dateRange }: TopDaysChartProps) {
  const [animated, setAnimated] = useState(false);
  const [hoveredIndex, setHoveredIndex] = useState<number | null>(null);

  useEffect(() => {
    const t = setTimeout(() => setAnimated(true), 200);
    return () => clearTimeout(t);
  }, []);

  if (!data || data.length === 0) {
    return (
      <div className="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6 shadow-sm text-center py-12">
        <p className="text-gray-500 dark:text-gray-400">
          No hay datos disponibles para mostrar los días con más órdenes
        </p>
      </div>
    );
  }

  // Usar datos ya calculados por el backend
  const topDays = data.slice(0, 5);

  // Extraer patrones de días que se repiten
  const repeatedDays = data
    .filter((day) => day.weekdayCount && day.weekdayCount > 1)
    .slice(0, 2);

  const subtitleText = dateRange
    ? `${new Date(dateRange.from).toLocaleDateString('es-CO')} - ${new Date(dateRange.to).toLocaleDateString('es-CO')}`
    : 'Período actual';

  return (
    <div className="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-900 p-6 shadow-sm">
      {/* Header */}
      <div className="mb-6 flex items-start justify-between">
        <div>
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
            Top 5 días con más órdenes
          </h3>
          <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">
            {subtitleText}
          </p>
        </div>
        <div className="inline-flex items-center gap-2 rounded-full bg-purple-50 dark:bg-purple-950 px-3 py-1">
          <div className="h-2 w-2 rounded-full bg-purple-600"></div>
          <span className="text-xs font-medium text-purple-600 dark:text-purple-400">
            Barras verticales
          </span>
        </div>
      </div>

      {/* Barras Verticales */}
      <div className="mb-6">
        <div className="flex items-end gap-2 h-32 px-2">
          {topDays.map((day, index) => {
            // Usar datos pre-calculados del backend, con fallback para compatibilidad
            const heightPercent = day.heightPercent || ((day.count / (day.maxCount || 100)) * 100);
            const opacity = day.opacity ?? [1, 0.82, 0.66, 0.5, 0.36][index];
            const isHovered = hoveredIndex === index;

            return (
              <div
                key={day.date}
                className="flex-1 flex flex-col items-center justify-end gap-1 cursor-pointer"
                onMouseEnter={() => setHoveredIndex(index)}
                onMouseLeave={() => setHoveredIndex(null)}
              >
                {/* Valor numérico */}
                <div className="text-xs font-medium text-gray-700 dark:text-gray-300 h-5">
                  {day.count.toLocaleString('es-CO')}
                </div>

                {/* Barra */}
                <div
                  className="w-full rounded-t-md transition-all duration-300 max-h-16"
                  style={{
                    background: isHovered
                      ? 'linear-gradient(180deg, #A78BFA, #C4B5FD)'
                      : 'linear-gradient(180deg, #7C3AED, #A78BFA)',
                    height: animated ? `${Math.min(heightPercent, 128)}px` : '0px',
                    opacity: isHovered ? 1 : opacity,
                    minHeight: animated && heightPercent > 0 ? '4px' : '0px',
                    transitionTimingFunction: 'cubic-bezier(0.4, 0, 0.2, 1)',
                    filter: isHovered ? 'drop-shadow(0 4px 12px rgba(124, 58, 237, 0.4))' : 'none',
                    transform: isHovered ? 'scaleY(1.05)' : 'scaleY(1)',
                    transformOrigin: 'bottom',
                  }}
                />

                {/* Fecha */}
                <div className="text-xs text-gray-500 dark:text-gray-400 mt-2">
                  {new Date(day.date).toLocaleDateString('es-CO', {
                    day: 'numeric',
                    month: 'short',
                  })}
                </div>

                {/* Día de semana */}
                <div className="text-xs text-gray-500 dark:text-gray-400">
                  {day.dayShort
                    ? day.dayShort.charAt(0).toUpperCase() + day.dayShort.slice(1)
                    : new Date(day.date)
                        .toLocaleDateString('es-CO', { weekday: 'short' })
                        .charAt(0)
                        .toUpperCase() +
                      new Date(day.date)
                        .toLocaleDateString('es-CO', { weekday: 'short' })
                        .slice(1)}
                </div>
              </div>
            );
          })}
        </div>

        {/* Línea divisora */}
        <div className="border-t border-gray-200 dark:border-gray-700 mt-4" />
      </div>

      {/* Badges de Patrón */}
      {repeatedDays.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {repeatedDays.map((day, idx) => {
            const colors =
              idx === 0
                ? { bg: '#EDE9FE', text: '#5B21B6' }
                : { bg: '#FEF3C7', text: '#92400E' };

            return (
              <div
                key={day.date}
                className="inline-flex items-center gap-2 px-3 py-1 rounded-full text-xs font-medium"
                style={{
                  backgroundColor: colors.bg,
                  color: colors.text,
                }}
              >
                {day.weekdayCount} de 5 son {day.dayName}
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
