'use client';
/* eslint-disable react-hooks/preserve-manual-memoization */

import React, { useMemo, useState } from 'react';
import { Settings } from 'lucide-react';

interface FamilyData {
  carrier: string;
  income: number;
  percentage: number;
  color: string;
  gradient: { start: string; end: string };
}

interface CarrierDistributionCardProps {
  data: Array<{
    carrier: string;
    income?: number;
    count?: number;
  }>;
  currency?: string;
  title?: string;
  subtitle?: string;
  valueLabel?: 'income' | 'count';
}

const CARRIER_COLORS: Record<
  string,
  { color: string; lightBg: string; gradient: { start: string; end: string } }
> = {
  INTERRAPIDISIMO: {
    color: '#1F2937',
    lightBg: '#F3F4F6',
    gradient: { start: '#374151', end: '#111827' },
  },
  ENVIA: {
    color: '#EF4444',
    lightBg: '#FEF2F2',
    gradient: { start: '#EF4444', end: '#7F1D1D' },
  },
  COORDINADORA: {
    color: '#3B82F6',
    lightBg: '#EFF6FF',
    gradient: { start: '#3B82F6', end: '#1E40AF' },
  },
  SERVIENTREGA: {
    color: '#F59E0B',
    lightBg: '#FFFBEB',
    gradient: { start: '#F59E0B', end: '#78350F' },
  },
  TODOCARGO: {
    color: '#8B5CF6',
    lightBg: '#FAF5FF',
    gradient: { start: '#8B5CF6', end: '#4C1D95' },
  },
  DHLEXPRESS: {
    color: '#EC4899',
    lightBg: '#FDF2F8',
    gradient: { start: '#EC4899', end: '#831843' },
  },
  FEDEX: {
    color: '#06B6D4',
    lightBg: '#F0F9FA',
    gradient: { start: '#06B6D4', end: '#164E63' },
  },
  MELON: {
    color: '#D4A574',
    lightBg: '#FAF6F1',
    gradient: { start: '#D4A574', end: '#8B6F47' },
  },
  MANUAL: {
    color: '#EC4899',
    lightBg: '#FDF2F8',
    gradient: { start: '#EC4899', end: '#831843' },
  },
  OTROS: {
    color: '#6B7280',
    lightBg: '#F9FAFB',
    gradient: { start: '#6B7280', end: '#1F2937' },
  },
  melon: {
    color: '#D4A574',
    lightBg: '#FAF6F1',
    gradient: { start: '#D4A574', end: '#8B6F47' },
  },
  manual: {
    color: '#EC4899',
    lightBg: '#FDF2F8',
    gradient: { start: '#EC4899', end: '#831843' },
  },
  otros: {
    color: '#6B7280',
    lightBg: '#F9FAFB',
    gradient: { start: '#6B7280', end: '#1F2937' },
  },
};

const DEFAULT_COLOR = { color: '#6B7280', lightBg: '#F9FAFB', gradient: { start: '#6B7280', end: '#1F2937' } };

export function CarrierDistributionCard({
  data,
  currency = 'COP',
  title = 'Ingresos por Familia',
  subtitle = 'Distribución de ingresos por transportadora',
  valueLabel = 'income',
}: CarrierDistributionCardProps) {
  const [hoveredCarrier, setHoveredCarrier] = useState<string | null>(null);

  const getValue = (item: typeof data[0]): number => {
    if (valueLabel === 'income') {
      return item.income || 0;
    }
    return item.count || 0;
  };

  const getCarrierColor = (carrierName: string) => {
    if (!carrierName) return DEFAULT_COLOR;

    const trimmed = carrierName.trim();
    const normalized = trimmed.toUpperCase();

    // Try exact match first
    if (CARRIER_COLORS[trimmed]) {
      return CARRIER_COLORS[trimmed];
    }

    // Try case-insensitive exact match
    const exactMatch = Object.entries(CARRIER_COLORS).find(
      ([key]) => key.toUpperCase() === normalized
    );
    if (exactMatch) return exactMatch[1];

    // Try partial match (contains)
    const partialMatch = Object.entries(CARRIER_COLORS).find(
      ([key]) => normalized.includes(key.toUpperCase()) || key.toUpperCase().includes(normalized)
    );
    if (partialMatch) return partialMatch[1];

    console.warn(`Color no encontrado para carrier: "${carrierName}"`);
    return DEFAULT_COLOR;
  };

  const processedData = useMemo(() => {
    const total = data.reduce((sum, item) => sum + getValue(item), 0);

    const sorted = [...data]
      .sort((a, b) => getValue(b) - getValue(a))
      .slice(0, 4)
      .map((item) => {
        const colors = getCarrierColor(item.carrier);
        return {
          carrier: item.carrier,
          value: getValue(item),
          percentage: (getValue(item) / total) * 100,
          color: colors.color,
          lightBg: colors.lightBg,
          gradient: colors.gradient,
        };
      });

    const others = data
      .slice(4)
      .reduce((sum, item) => sum + getValue(item), 0);

    if (others > 0) {
      sorted.push({
        carrier: 'Otros',
        value: others,
        percentage: (others / total) * 100,
        color: DEFAULT_COLOR.color,
        lightBg: DEFAULT_COLOR.lightBg,
        gradient: DEFAULT_COLOR.gradient,
      });
    }

    return { data: sorted, total };
  }, [data, valueLabel]);

  const donutSlices = useMemo(() => {
    const radius = 90;
    const innerRadius = 54;
    const slices: Array<{ carrier: string; path: string; color: string; opacity: number }> = [];
    let currentAngle = -Math.PI / 2;

    processedData.data.forEach((item) => {
      const sliceAngle = (item.percentage / 100) * Math.PI * 2;
      const startAngle = currentAngle;
      const endAngle = currentAngle + sliceAngle;

      const x1 = 120 + radius * Math.cos(startAngle);
      const y1 = 120 + radius * Math.sin(startAngle);
      const x2 = 120 + radius * Math.cos(endAngle);
      const y2 = 120 + radius * Math.sin(endAngle);

      const ix1 = 120 + innerRadius * Math.cos(startAngle);
      const iy1 = 120 + innerRadius * Math.sin(startAngle);
      const ix2 = 120 + innerRadius * Math.cos(endAngle);
      const iy2 = 120 + innerRadius * Math.sin(endAngle);

      const largeArc = sliceAngle > Math.PI ? 1 : 0;

      const pathData = [
        `M ${x1} ${y1}`,
        `A ${radius} ${radius} 0 ${largeArc} 1 ${x2} ${y2}`,
        `L ${ix2} ${iy2}`,
        `A ${innerRadius} ${innerRadius} 0 ${largeArc} 0 ${ix1} ${iy1}`,
        'Z',
      ].join(' ');

      const isHovered = hoveredCarrier === item.carrier;
      const opacity = hoveredCarrier && !isHovered ? 0.4 : 1;

      slices.push({
        carrier: item.carrier,
        path: pathData,
        color: item.gradient.start,
        opacity,
      });

      currentAngle = endAngle;
    });

    return slices;
  }, [processedData.data, hoveredCarrier]);

  const formatValue = (value: number) => {
    if (valueLabel === 'income') {
      return new Intl.NumberFormat('es-CO', {
        style: 'currency',
        currency,
        minimumFractionDigits: 0,
        maximumFractionDigits: 0,
      }).format(value);
    }
    return value.toLocaleString('es-CO');
  };

  return (
    <div className="rounded-lg border border-gray-200 bg-white p-6 shadow-sm">
      {/* Header */}
      <div className="mb-6 flex items-start justify-between">
        <div className="flex items-center gap-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg bg-blue-50">
            <Settings className="h-5 w-5 text-blue-600" />
          </div>
          <div>
            <h3 className="font-semibold text-gray-900" style={{ fontFamily: 'Plus Jakarta Sans' }}>
              {title}
            </h3>
            <p className="text-sm text-gray-500">{subtitle}</p>
          </div>
        </div>
        <div className="flex items-center gap-2 rounded-full bg-blue-50 px-3 py-1">
          <span
            className="text-sm font-medium text-blue-600"
            style={{ fontFamily: 'JetBrains Mono' }}
          >
            {formatValue(processedData.total)}
          </span>
        </div>
      </div>

      {/* Content */}
      <div className="flex gap-12">
        {/* Donut Chart - Left Side */}
        <div className="flex-1 flex items-center justify-center">
          <svg viewBox="0 0 240 240" className="h-full w-full" style={{ minHeight: '400px', maxWidth: '400px' }}>
            <defs>
              {processedData.data.map((item) => {
                const gradId = `gradient-${item.carrier.toUpperCase().replace(/\s+/g, '-')}`;
                return (
                  <linearGradient
                    key={gradId}
                    id={gradId}
                    x1="0%"
                    y1="0%"
                    x2="100%"
                    y2="100%"
                  >
                    <stop offset="0%" stopColor={item.gradient.start} />
                    <stop offset="100%" stopColor={item.gradient.end} />
                  </linearGradient>
                );
              })}
            </defs>

            {donutSlices.map((slice) => {
              const gradId = `gradient-${slice.carrier.toUpperCase().replace(/\s+/g, '-')}`;
              return (
                <path
                  key={slice.carrier}
                  d={slice.path}
                  fill={`url(#${gradId})`}
                  opacity={slice.opacity}
                  className="cursor-pointer transition-opacity duration-200"
                  onMouseEnter={() => setHoveredCarrier(slice.carrier)}
                  onMouseLeave={() => setHoveredCarrier(null)}
                />
              );
            })}

            {/* Center circle for donut effect */}
            <circle cx="120" cy="120" r="54" fill="white" />
          </svg>
        </div>

        {/* Metrics Grid - Right Side */}
        <div className="flex-1 flex flex-col justify-center">
          <div className="grid grid-cols-2 gap-4">
            {processedData.data.map((item) => {
              const isHovered = hoveredCarrier === item.carrier;
              return (
                <div
                  key={item.carrier}
                  className="rounded-lg border-2 p-4 transition-all duration-200"
                  onMouseEnter={() => setHoveredCarrier(item.carrier)}
                  onMouseLeave={() => setHoveredCarrier(null)}
                  style={{
                    backgroundColor: isHovered ? item.lightBg : '#FFFFFF',
                    borderColor: item.color,
                    opacity: hoveredCarrier && !isHovered ? 0.5 : 1,
                    boxShadow: isHovered ? `0 4px 12px ${item.color}20` : 'none',
                  }}
                >
                  <div className="mb-2 flex items-center justify-between">
                    <span
                      className="text-sm font-semibold"
                      style={{ color: item.color }}
                    >
                      {item.carrier.toUpperCase()}
                    </span>
                    <span
                      className="text-xs font-bold px-2 py-1 rounded-full text-white"
                      style={{ backgroundColor: item.color, fontFamily: 'JetBrains Mono' }}
                    >
                      {item.percentage.toFixed(1)}%
                    </span>
                  </div>
                  <p
                    className="text-lg font-bold"
                    style={{ color: item.color, fontFamily: 'JetBrains Mono' }}
                  >
                    {formatValue(item.value)}
                  </p>
                </div>
              );
            })}
          </div>
        </div>
      </div>
    </div>
  );
}
