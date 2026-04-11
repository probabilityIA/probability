'use client';

import { useState, useMemo } from 'react';
import { format } from 'date-fns';
import { es } from 'date-fns/locale';
import { DashboardStats, OrdersByWeek, OrdersByMonth, ShipmentsByDayOfWeek, ShipmentsByCarrier } from '../../domain/types';
import {
  ComposedChart,
  BarChart,
  PieChart,
  LineChart,
  Bar,
  Line,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  Area,
  AreaChart,
  ReferenceLine,
  ResponsiveContainer,
  LabelList,
} from 'recharts';

interface DashboardChartsProps {
  stats: DashboardStats | null;
  selectedBusinessId?: number;
}

// Colores
const COLORS = {
  primary: '#8B5CF6',
  secondary: '#A78BFA',
  tertiary: '#DDD6FE',
  success: '#10B981',
  warning: '#F59E0B',
  danger: '#EF4444',
};

// Funciones matemáticas
const calculateEMA = (data: number[], alpha: number = 0.4): number[] => {
  if (data.length === 0) return [];
  const ema: number[] = [];
  ema[0] = data[0];
  for (let i = 1; i < data.length; i++) {
    ema[i] = data[i] * alpha + ema[i - 1] * (1 - alpha);
  }
  return ema;
};

const forecastEMA = (ema: number[], periods: number = 4): number[] => {
  const forecast: number[] = [];
  const lastEMA = ema[ema.length - 1];
  for (let i = 0; i < periods; i++) {
    // Proyectamos el último EMA con una ligera tendencia decreciente
    const trend = (ema[ema.length - 1] - ema[Math.max(0, ema.length - 4)]) / 4;
    forecast.push(lastEMA + trend * (i + 1));
  }
  return forecast;
};

const calculateLinearRegression = (data: number[]): { slope: number; intercept: number } => {
  const n = data.length;
  let sumX = 0, sumY = 0, sumXY = 0, sumX2 = 0;

  for (let i = 0; i < n; i++) {
    sumX += i;
    sumY += data[i];
    sumXY += i * data[i];
    sumX2 += i * i;
  }

  const slope = (n * sumXY - sumX * sumY) / (n * sumX2 - sumX * sumX);
  const intercept = (sumY - slope * sumX) / n;

  return { slope, intercept };
};

export default function DashboardCharts({ stats, selectedBusinessId }: DashboardChartsProps) {
  const [activeTab, setActiveTab] = useState<'forecast' | 'monthly' | 'demand' | 'carrier'>('forecast');
  const [topSellingDays, setTopSellingDays] = useState<any[]>([]);

  // Tab 1: Pronóstico de Órdenes
  const forecastData = useMemo(() => {
    if (!stats?.orders_by_week || stats.orders_by_week.length === 0) return null;

    const historicalWeeks = stats.orders_by_week.map((w, idx) => {
      const startDate = new Date(w.start_date);
      const endDate = new Date(w.end_date);
      const dateRange = `${format(startDate, 'dd MMM', { locale: es })} – ${format(endDate, 'dd MMM', { locale: es })}`;

      return {
        weekLabel: `Sem ${idx + 1}`,
        week: `Sem ${idx + 1}`,
        orders: w.count,
        dateRange,
        type: 'historical' as const,
      };
    });

    const ordersArray = stats.orders_by_week.map(w => w.count);
    const emaValues = calculateEMA(ordersArray);
    const forecastedValues = forecastEMA(emaValues);

    const forecastWeeks = forecastedValues.map((orders, idx) => {
      const weekNum = stats.orders_by_week!.length + idx + 1;
      return {
        weekLabel: `Sem ${weekNum}`,
        week: `Sem ${weekNum}`,
        orders: Math.round(orders),
        forecast: Math.round(orders),
        upper: Math.round(orders * 1.15),
        lower: Math.round(orders * 0.85),
        dateRange: 'Pronóstico',
        type: 'forecast' as const,
      };
    });

    return [...historicalWeeks, ...forecastWeeks];
  }, [stats?.orders_by_week]);

  // Tab 2: Órdenes por Mes
  const monthlyData = useMemo(() => {
    if (!stats?.orders_by_month || stats.orders_by_month.length === 0) return null;

    const data = stats.orders_by_month.map((m: OrdersByMonth) => ({
      month: m.month?.split(' ')[0] || m.month || '',
      orders: m.count,
    }));

    // Calcular línea de tendencia
    const ordersArray = data.map(d => d.orders);
    const { slope, intercept } = calculateLinearRegression(ordersArray);

    const trendData = data.map((d, i) => ({
      ...d,
      trend: Math.round(slope * i + intercept),
    }));

    return trendData;
  }, [stats?.orders_by_month]);

  // Tab 3: TOP 5 Días de Mayor Demanda (fechas específicas)
  // Cargar desde endpoint backend GET /api/v1/dashboard/top-selling-days
  useMemo(() => {
    const fetchTopDays = async () => {
      try {
        const url = new URL('/api/v1/dashboard/top-selling-days', window.location.origin);
        if (selectedBusinessId) {
          url.searchParams.append('business_id', selectedBusinessId.toString());
        }
        url.searchParams.append('limit', '5');

        const response = await fetch(url.toString());
        if (response.ok) {
          const result = await response.json();
          if (result.data && Array.isArray(result.data)) {
            setTopSellingDays(result.data);
          }
        }
      } catch (error) {
        console.error('Error fetching top selling days:', error);
      }
    };

    fetchTopDays();
  }, [selectedBusinessId]);

  const demandData = useMemo(() => {
    if (topSellingDays.length === 0) return null;

    return topSellingDays.map((d, idx) => ({
      date: d.date,
      label: d.formatted,
      orders: d.total,
      isTop: idx === 0,
    }));
  }, [topSellingDays]);

  // Tab 4: Por Transportadora
  const CARRIER_COLORS = {
    enviame: '#7C3AED',
    envioclick: '#06B6D4',
    mipaquete: '#F59E0B',
  };

  const carrierData = useMemo(() => {
    if (!stats?.shipments_by_carrier) return null;

    const carriers: { [key: string]: any } = {};
    let totalShipments = 0;

    (stats.shipments_by_carrier as ShipmentsByCarrier[]).forEach(c => {
      const cleanCarrier = c.carrier?.toLowerCase() || 'unknown';
      const carrierKey = cleanCarrier.includes('enviame')
        ? 'enviame'
        : cleanCarrier.includes('envioclick')
          ? 'envioclick'
          : cleanCarrier.includes('mipaquete')
            ? 'mipaquete'
            : 'other';

      carriers[cleanCarrier] = {
        name: c.carrier || 'Unknown',
        displayName: c.carrier || 'Unknown',
        value: c.count,
        fill: CARRIER_COLORS[carrierKey as keyof typeof CARRIER_COLORS] || CARRIER_COLORS.mipaquete,
      };
      totalShipments += c.count;
    });

    const carrierList = Object.values(carriers);
    const carrierMetrics = carrierList.map(carrier => ({
      ...carrier,
      percentage: totalShipments > 0 ? ((carrier.value / totalShipments) * 100).toFixed(1) : '0',
    }));

    return { carriers: carrierList, carrierMetrics, totalShipments };
  }, [stats?.shipments_by_carrier]);

  // Custom Tooltip
  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;

      // Para gráficas de pronóstico y meses
      if (data.week && data.dateRange) {
        return (
          <div className="bg-white p-3 border border-gray-200 rounded-lg shadow-md">
            <p className="text-sm font-semibold text-gray-900">
              {data.week} · {data.dateRange}
            </p>
            <p className="text-sm text-gray-700">
              {data.orders.toLocaleString()} órdenes
            </p>
            {data.upper && (
              <p className="text-xs text-gray-500 mt-1">±{((data.upper - data.lower) / 2 / data.orders * 100).toFixed(0)}%</p>
            )}
          </div>
        );
      }

      // Para otros gráficos
      return (
        <div className="bg-white p-3 border border-gray-200 rounded-lg shadow-md">
          <p className="text-sm font-semibold text-gray-900">{data.week || data.month || data.day}</p>
          <p className="text-sm text-gray-700">
            {data.orders !== undefined ? `${data.orders.toLocaleString()} órdenes` : `${data.value} envíos`}
          </p>
        </div>
      );
    }
    return null;
  };

  const tabs = [
    { id: 'forecast', label: 'Pronóstico de Órdenes', icon: '📈' },
    { id: 'monthly', label: 'Órdenes por Mes', icon: '📊' },
    { id: 'demand', label: 'Días de Mayor Demanda', icon: '🔥' },
    { id: 'carrier', label: 'Por Transportadora', icon: '🚚' },
  ];

  return (
    <div className="mt-6 p-6 bg-white rounded-xl border border-gray-200 shadow-sm">
      {/* Tabs */}
      <div className="flex gap-2 mb-6 overflow-x-auto pb-2">
        {tabs.map(tab => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id as any)}
            className={`px-4 py-2 rounded-full font-medium text-sm whitespace-nowrap transition-all ${
              activeTab === tab.id
                ? 'bg-purple-500 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            }`}
          >
            {tab.icon} {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="animate-fadeIn">
        {/* TAB 1: Pronóstico */}
        {activeTab === 'forecast' && forecastData && (
          <ResponsiveContainer width="100%" height={280}>
            <ComposedChart data={forecastData} margin={{ top: 20, right: 30, left: 0, bottom: 60 }}>
              <defs>
                <linearGradient id="confidenceGradient" x1="0" y1="0" x2="0" y2="1">
                  <stop offset="5%" stopColor={COLORS.primary} stopOpacity={0.1} />
                  <stop offset="95%" stopColor={COLORS.primary} stopOpacity={0} />
                </linearGradient>
              </defs>
              <CartesianGrid strokeDasharray="3 3" stroke="#f5f5f5" />
              <XAxis dataKey="weekLabel" tick={{ fontSize: 12, fill: '#9CA3AF' }} angle={-45} textAnchor="end" height={80} />
              <YAxis tick={{ fontSize: 12, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <Tooltip content={<CustomTooltip />} />
              <ReferenceLine
                x={forecastData[forecastData.length - 4]?.weekLabel}
                stroke="#D1D5DB"
                strokeDasharray="5 5"
                label={{ value: 'Pronóstico →', position: 'top', fill: '#6B7280', fontSize: 12 }}
              />
              <Area
                type="monotone"
                dataKey="upper"
                fill="url(#confidenceGradient)"
                stroke="none"
                isAnimationActive={false}
              />
              <Area
                type="monotone"
                dataKey="lower"
                fill="white"
                stroke="none"
                isAnimationActive={false}
              />
              <Line
                type="monotone"
                dataKey="orders"
                stroke={COLORS.primary}
                strokeWidth={2}
                dot={false}
                isAnimationActive={false}
              />
              <Line
                type="monotone"
                dataKey="forecast"
                stroke={COLORS.primary}
                strokeWidth={2}
                strokeDasharray="5 4"
                dot={false}
                isAnimationActive={false}
              />
            </ComposedChart>
          </ResponsiveContainer>
        )}

        {/* TAB 2: Órdenes por Mes */}
        {activeTab === 'monthly' && monthlyData && (
          <ResponsiveContainer width="100%" height={280}>
            <ComposedChart data={monthlyData} margin={{ top: 20, right: 30, left: 0, bottom: 60 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#f5f5f5" />
              <XAxis dataKey="month" tick={{ fontSize: 12, fill: '#9CA3AF' }} angle={-45} textAnchor="end" height={80} />
              <YAxis tick={{ fontSize: 12, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <Tooltip content={<CustomTooltip />} />
              <Bar dataKey="orders" fill={COLORS.primary} fillOpacity={0.15} radius={[4, 4, 0, 0]}>
                <LabelList dataKey="orders" position="top" fontSize={12} fill="#6B7280" />
              </Bar>
              <Line
                type="monotone"
                dataKey="trend"
                stroke={COLORS.primary}
                strokeWidth={2}
                dot={false}
                isAnimationActive={false}
              />
            </ComposedChart>
          </ResponsiveContainer>
        )}

        {/* TAB 3: TOP 5 Días de Mayor Demanda */}
        {activeTab === 'demand' && demandData && (
          <ResponsiveContainer width="100%" height={280}>
            <BarChart
              data={demandData}
              layout="vertical"
              margin={{ top: 5, right: 30, left: 140, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="#f5f5f5" />
              <XAxis type="number" tick={{ fontSize: 12, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <YAxis
                dataKey="label"
                type="category"
                tick={{ fontSize: 12, fill: '#9CA3AF' }}
                width={135}
              />
              <Tooltip
                contentStyle={{ backgroundColor: '#fff', border: '1px solid #e5e7eb', borderRadius: '8px' }}
                formatter={(value) => value ? `${(value as number).toLocaleString()} órdenes` : '0 órdenes'}
              />
              <Bar
                dataKey="orders"
                radius={[0, 4, 4, 0]}
              >
                {demandData?.map((entry, index) => (
                  <Cell
                    key={`cell-${index}`}
                    fill={entry.isTop ? COLORS.primary : 'rgba(139, 92, 246, 0.2)'}
                  />
                ))}
                <LabelList dataKey="orders" position="right" fontSize={12} fill="#6B7280" />
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        )}

        {/* TAB 4: Por Transportadora */}
        {activeTab === 'carrier' && carrierData && (
          <div className="flex flex-col items-center gap-8">
            {/* PieChart Donut con Centro Customizado */}
            <div className="w-80 relative">
              <ResponsiveContainer width="100%" height={320}>
                <PieChart>
                  <Pie
                    data={carrierData.carriers}
                    cx="50%"
                    cy="50%"
                    innerRadius={65}
                    outerRadius={110}
                    paddingAngle={2}
                    dataKey="value"
                    isAnimationActive={true}
                    animationBegin={0}
                    animationDuration={800}
                    label={({ displayName, percent }) => `${carrierData.carriers.find(c => c.displayName === displayName)?.displayName || 'Unknown'} ${(percent * 100).toFixed(0)}%`}
                    labelLine={true}
                  >
                    {carrierData.carriers.map((entry, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={entry.fill}
                        stroke="white"
                        strokeWidth={3}
                      />
                    ))}
                  </Pie>
                  <Tooltip
                    formatter={(value) => value ? `${(value as number).toLocaleString()} envíos` : '0 envíos'}
                    contentStyle={{ backgroundColor: '#fff', border: '1px solid #e5e7eb', borderRadius: '8px' }}
                  />
                </PieChart>
              </ResponsiveContainer>
              {/* Centro del Donut */}
              <div className="absolute inset-0 flex flex-col items-center justify-center">
                <div className="text-xs font-medium text-gray-600">Total</div>
                <div className="text-3xl font-bold text-gray-900">{carrierData.totalShipments.toLocaleString()}</div>
              </div>
            </div>

            {/* Metric Cards por Transportadora */}
            <div className="w-full grid grid-cols-3 gap-4">
              {carrierData.carrierMetrics?.map((carrier) => (
                <div
                  key={carrier.displayName}
                  className="p-4 bg-white rounded-lg border border-gray-200"
                  style={{
                    borderTop: `3px solid ${carrier.fill}`,
                  }}
                >
                  <div className="text-sm font-semibold text-gray-900 mb-2">{carrier.displayName}</div>
                  <div className="text-2xl font-bold text-gray-900">{carrier.value.toLocaleString()}</div>
                  <div
                    className="text-xs font-medium mt-2 px-2 py-1 rounded-full text-white w-fit"
                    style={{ backgroundColor: carrier.fill }}
                  >
                    {carrier.percentage}%
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
