'use client';

import { useState, useMemo } from 'react';
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

export default function DashboardCharts({ stats }: DashboardChartsProps) {
  const [activeTab, setActiveTab] = useState<'forecast' | 'monthly' | 'demand' | 'carrier'>('forecast');

  // Tab 1: Pronóstico de Órdenes
  const forecastData = useMemo(() => {
    if (!stats?.orders_by_week || stats.orders_by_week.length === 0) return null;

    const historicalWeeks = stats.orders_by_week.map(w => ({
      week: w.week,
      orders: w.count,
      type: 'historical' as const,
    }));

    const ordersArray = stats.orders_by_week.map(w => w.count);
    const emaValues = calculateEMA(ordersArray);
    const forecastedValues = forecastEMA(emaValues);

    const forecastWeeks = forecastedValues.map((orders, idx) => {
      const weekNum = stats.orders_by_week!.length + idx + 1;
      return {
        week: `Sem ${weekNum}`,
        orders: Math.round(orders),
        forecast: Math.round(orders),
        upper: Math.round(orders * 1.15),
        lower: Math.round(orders * 0.85),
        type: 'forecast' as const,
      };
    });

    return [...historicalWeeks, ...forecastWeeks];
  }, [stats?.orders_by_week]);

  // Tab 2: Órdenes por Mes
  const monthlyData = useMemo(() => {
    if (!stats?.orders_by_month || stats.orders_by_month.length === 0) return null;

    const data = stats.orders_by_month.map((m: OrdersByMonth) => ({
      month: m.month?.split(' ')[0] || '',
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

  // Tab 3: Días de Mayor Demanda
  const demandData = useMemo(() => {
    if (!stats?.shipments_by_day_of_week) return null;

    return (stats.shipments_by_day_of_week as ShipmentsByDayOfWeek[])
      .map(d => ({
        day: d.day_name,
        orders: d.count,
      }))
      .sort((a, b) => b.orders - a.orders);
  }, [stats?.shipments_by_day_of_week]);

  // Tab 4: Por Transportadora
  const carrierData = useMemo(() => {
    if (!stats?.shipments_by_carrier) return null;

    const carriers: { [key: string]: any } = {};

    (stats.shipments_by_carrier as ShipmentsByCarrier[]).forEach(c => {
      const cleanCarrier = c.carrier?.toLowerCase() || 'unknown';
      carriers[cleanCarrier] = {
        name: cleanCarrier,
        value: c.count,
        fill: cleanCarrier.includes('enviame')
          ? COLORS.primary
          : cleanCarrier.includes('envioclick')
            ? COLORS.secondary
            : COLORS.tertiary,
      };
    });

    return Object.values(carriers);
  }, [stats?.shipments_by_carrier]);

  // Custom Tooltip
  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;
      return (
        <div className="bg-white p-3 border border-gray-200 rounded-lg shadow-md">
          <p className="text-sm font-semibold text-gray-900">{data.week || data.month || data.day}</p>
          <p className="text-sm text-gray-700">
            Órdenes: <span className="font-bold">{data.orders || data.value}</span>
          </p>
          {data.type && <p className="text-xs text-gray-500 mt-1">{data.type === 'forecast' ? '📊 Predicción' : '📈 Histórico'}</p>}
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
              <XAxis dataKey="week" tick={{ fontSize: 12, fill: '#9CA3AF' }} angle={-45} textAnchor="end" height={80} />
              <YAxis tick={{ fontSize: 12, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <Tooltip content={<CustomTooltip />} />
              <ReferenceLine
                x={forecastData[forecastData.length - 5]?.week}
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

        {/* TAB 3: Días de Mayor Demanda */}
        {activeTab === 'demand' && demandData && (
          <ResponsiveContainer width="100%" height={280}>
            <BarChart
              data={demandData}
              layout="vertical"
              margin={{ top: 5, right: 30, left: 100, bottom: 5 }}
            >
              <CartesianGrid strokeDasharray="3 3" stroke="#f5f5f5" />
              <XAxis type="number" tick={{ fontSize: 12, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <YAxis dataKey="day" type="category" tick={{ fontSize: 12, fill: '#9CA3AF' }} width={95} />
              <Tooltip content={<CustomTooltip />} />
              <Bar
                dataKey="orders"
                fill={COLORS.primary}
                radius={[0, 4, 4, 0]}
              >
                <LabelList dataKey="orders" position="right" fontSize={12} fill="#6B7280" />
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        )}

        {/* TAB 4: Por Transportadora */}
        {activeTab === 'carrier' && carrierData && (
          <div className="flex gap-6 justify-center items-start">
            <div className="w-64">
              <ResponsiveContainer width="100%" height={280}>
                <PieChart>
                  <Pie
                    data={carrierData}
                    cx="50%"
                    cy="50%"
                    innerRadius={50}
                    outerRadius={90}
                    paddingAngle={2}
                    dataKey="value"
                    label={({ name, percent }) => `${name} ${(percent * 100).toFixed(0)}%`}
                  >
                    {carrierData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={entry.fill} />
                    ))}
                  </Pie>
                  <Tooltip formatter={(value) => `${value} envíos`} />
                </PieChart>
              </ResponsiveContainer>
            </div>

            <div className="flex-1">
              <div className="text-sm font-semibold text-gray-900 mb-4">Estado de Envíos</div>
              <div className="space-y-3">
                {stats?.shipments_by_status?.map((status: any) => (
                  <div key={status.status} className="flex items-center gap-3">
                    <div
                      className="w-3 h-3 rounded-full"
                      style={{
                        backgroundColor:
                          status.status === 'delivered'
                            ? COLORS.success
                            : status.status === 'in_transit'
                              ? COLORS.primary
                              : COLORS.danger,
                      }}
                    />
                    <span className="text-sm text-gray-700">
                      {status.status === 'delivered'
                        ? 'Entregado'
                        : status.status === 'in_transit'
                          ? 'En tránsito'
                          : 'Fallido'}
                      : {status.count}
                    </span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
