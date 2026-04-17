'use client';

import { useState, useMemo, useEffect } from 'react';
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
  const defaultTabs = [
    { id: 'forecast', label: 'Pronóstico de Órdenes', icon: '📈' },
    { id: 'monthly', label: 'Órdenes por Mes', icon: '📊' },
    { id: 'demand', label: 'Días de Mayor Demanda', icon: '🔥' },
    { id: 'carrier', label: 'Por Transportadora', icon: '🚚' },
  ];

  const [activeTab, setActiveTab] = useState<'forecast' | 'monthly' | 'demand' | 'carrier'>('forecast');
  const [tabs, setTabs] = useState(defaultTabs);
  const [draggedTab, setDraggedTab] = useState<string | null>(null);
  const [topSellingDays, setTopSellingDays] = useState<any[]>([]);
  const [hoveredValue, setHoveredValue] = useState<number | null>(null);
  const [hoverX, setHoverX] = useState<string | null>(null);

  // Cargar orden de tabs desde localStorage
  useEffect(() => {
    const saved = localStorage.getItem('dashboardTabsOrder');
    if (saved) {
      try {
        const order = JSON.parse(saved);
        setTabs(order);
      } catch (e) {
        // Si hay error, usar orden por defecto
      }
    }
  }, []);

  // Guardar orden de tabs en localStorage
  const saveTabs = (newTabs: any) => {
    setTabs(newTabs);
    localStorage.setItem('dashboardTabsOrder', JSON.stringify(newTabs));
  };

  const handleDragStart = (e: React.DragEvent, tabId: string) => {
    setDraggedTab(tabId);
    e.dataTransfer.effectAllowed = 'move';
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  const handleDrop = (e: React.DragEvent, targetTabId: string) => {
    e.preventDefault();
    if (!draggedTab || draggedTab === targetTabId) {
      setDraggedTab(null);
      return;
    }

    const draggedIndex = tabs.findIndex(t => t.id === draggedTab);
    const targetIndex = tabs.findIndex(t => t.id === targetTabId);

    const newTabs = [...tabs];
    [newTabs[draggedIndex], newTabs[targetIndex]] = [newTabs[targetIndex], newTabs[draggedIndex]];
    saveTabs(newTabs);
    setDraggedTab(null);
  };

  const handleDragEnd = () => {
    setDraggedTab(null);
  };

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

  // KPIs para Pronóstico
  const forecastKPIs = useMemo(() => {
    if (!forecastData || forecastData.length === 0) return null;

    const historicalData = forecastData.filter(d => d.type === 'historical');
    const forecastDataPoints = forecastData.filter(d => d.type === 'forecast');

    const historicalOrders = historicalData.map(d => d.orders);
    const avgHistorical = historicalOrders.length > 0
      ? Math.round(historicalOrders.reduce((a, b) => a + b, 0) / historicalOrders.length)
      : 0;

    const lastHistoricalWeek = historicalData[historicalData.length - 1]?.orders || 0;
    const projectedLastWeek = forecastDataPoints[0]?.forecast || 0;

    return {
      avgHistorical,
      projectedLastWeek,
      margin: 15,
      historicalCount: historicalData.length,
      splitIndex: historicalData.length,
    };
  }, [forecastData]);

  // Tab 2: Órdenes por Mes
  const monthlyData = useMemo(() => {
    if (!stats?.orders_by_month || stats.orders_by_month.length === 0) return null;

    const data = stats.orders_by_month.map((m: OrdersByMonth) => ({
      month: m.month?.split(' ')[0] || m.month || '',
      orders: m.count,
      percentage: m.percentage || 0,
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

  // Custom Tooltip con Dark Theme
  const CustomTooltip = ({ active, payload }: any) => {
    if (active && payload && payload.length) {
      const data = payload[0].payload;

      // Filtrar items sin valor
      const filteredPayload = payload.filter((p: any) => p.value !== null && p.value !== undefined);

      if (data.week && data.dateRange) {
        setHoveredValue(data.orders);
        setHoverX(data.weekLabel);

        return (
          <div style={{ backgroundColor: '#1e1e2e', padding: '16px 20px', borderRadius: '12px', border: '1.5px solid #3a3a4a', minWidth: '220px' }}>
            <p style={{ fontSize: '14px', fontWeight: 700, color: '#ffffff', margin: '0 0 12px 0' }}>
              {data.week} · {data.dateRange}
            </p>
            {filteredPayload.map((p: any, idx: number) => (
              p.value !== null && (
                <p key={idx} style={{ fontSize: '13px', color: '#ffffff', margin: '6px 0', fontWeight: 500 }}>
                  {p.name}: <span style={{ fontWeight: 700, color: p.color || '#3B82F6' }}>{(typeof p.value === 'number' ? p.value.toLocaleString() : p.value)}</span>
                </p>
              )
            ))}
            {data.upper && (
              <p style={{ fontSize: '12px', color: '#a0a0b0', marginTop: '10px', fontStyle: 'italic' }}>
                Margen: ±{((data.upper - data.lower) / 2 / data.orders * 100).toFixed(0)}%
              </p>
            )}
          </div>
        );
      }

      return (
        <div style={{ backgroundColor: '#1e1e2e', padding: '16px 20px', borderRadius: '12px', border: '1.5px solid #3a3a4a', minWidth: '200px' }}>
          <p style={{ fontSize: '14px', fontWeight: 700, color: '#ffffff', margin: 0 }}>
            {data.week || data.month || data.day}
          </p>
          <p style={{ fontSize: '13px', color: '#8B5CF6', margin: '8px 0 0 0', fontWeight: 600 }}>
            {data.orders !== undefined ? `${data.orders.toLocaleString()} órdenes` : `${data.value} envíos`}
          </p>
        </div>
      );
    }
    return null;
  };


  return (
    <div className="mt-6 bg-white rounded-xl border border-gray-200 shadow-sm overflow-hidden">
      {/* Folder Tabs */}
      <div className="flex gap-0 overflow-x-auto border-b border-gray-200 bg-gray-50 px-2 pt-2">
        {tabs.map((tab, idx) => (
          <button
            key={tab.id}
            draggable
            onDragStart={(e) => handleDragStart(e, tab.id)}
            onDragOver={handleDragOver}
            onDrop={(e) => handleDrop(e, tab.id)}
            onDragEnd={handleDragEnd}
            onClick={() => setActiveTab(tab.id as any)}
            className={`relative px-6 py-3 font-medium text-sm whitespace-nowrap transition-all cursor-move ${
              draggedTab === tab.id ? 'opacity-50' : 'opacity-100'
            } ${
              activeTab === tab.id
                ? 'bg-white text-purple-600 border-l border-r border-t border-gray-200 rounded-t-lg shadow-sm'
                : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
            }`}
            style={
              activeTab === tab.id
                ? {
                    borderBottomColor: 'white',
                    marginBottom: '-1px',
                  }
                : {}
            }
          >
            ⋮⋮ {tab.icon} {tab.label}
          </button>
        ))}
      </div>

      {/* Tab Content */}
      <div className="p-6 animate-fadeIn">
        {/* TAB 1: Pronóstico Mejorado */}
        {activeTab === 'forecast' && forecastData && forecastKPIs && (
          <div>
            {/* KPI Cards */}
            <div className="grid grid-cols-3 gap-4 mb-6">
              {/* Promedio Histórico */}
              <div className="bg-gradient-to-br from-blue-50 to-blue-100 p-4 rounded-lg border border-blue-200">
                <p className="text-xs font-medium text-blue-600 mb-1">Promedio Histórico</p>
                <p className="text-2xl font-bold text-blue-900">{forecastKPIs.avgHistorical.toLocaleString()}</p>
                <p className="text-xs text-blue-700 mt-1">{forecastKPIs.historicalCount} semanas</p>
              </div>

              {/* Valor Proyectado */}
              <div className="bg-gradient-to-br from-purple-50 to-purple-100 p-4 rounded-lg border border-purple-200">
                <p className="text-xs font-medium text-purple-600 mb-1">Valor Proyectado</p>
                <p className="text-2xl font-bold text-purple-900">{forecastKPIs.projectedLastWeek.toLocaleString()}</p>
                <p className="text-xs text-purple-700 mt-1">Próxima semana</p>
              </div>

              {/* Margen de Error */}
              <div className="bg-gradient-to-br from-amber-50 to-amber-100 p-4 rounded-lg border border-amber-200">
                <p className="text-xs font-medium text-amber-600 mb-1">Margen de Error</p>
                <p className="text-2xl font-bold text-amber-900">±{forecastKPIs.margin}%</p>
                <p className="text-xs text-amber-700 mt-1">Rango de confianza</p>
              </div>
            </div>

            {/* Badge con Hover Value + Gráfico */}
            <div className="relative">
              {hoveredValue !== null && (
                <div className="absolute top-2 right-2 bg-purple-600 text-white px-3 py-1 rounded-full text-xs font-semibold z-10 shadow-lg">
                  {hoveredValue.toLocaleString()} órdenes
                </div>
              )}

              <ResponsiveContainer width="100%" height={320}>
                <ComposedChart
                  data={forecastData}
                  margin={{ top: 40, right: 30, left: 0, bottom: 60 }}
                  onMouseMove={(state: any) => {
                    if (state.isTooltipActive && state.tooltipPayload && state.tooltipPayload[0]) {
                      setHoverX(state.tooltipPayload[0].payload.weekLabel);
                    }
                  }}
                >
                  <defs>
                    <linearGradient id="confidenceGradient" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#3B82F6" stopOpacity={0.25} />
                      <stop offset="95%" stopColor="#3B82F6" stopOpacity={0.02} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={true} />
                  <XAxis dataKey="weekLabel" tick={{ fontSize: 11, fill: '#9CA3AF' }} angle={-45} textAnchor="end" height={80} />
                  <YAxis tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
                  <Tooltip content={<CustomTooltip />} cursor={{ stroke: '#D1D5DB', strokeDasharray: '5 5' }} />

                  {/* Split Line: Línea vertical que separa histórico de pronóstico */}
                  <ReferenceLine
                    x={forecastData[forecastKPIs.splitIndex - 1]?.weekLabel}
                    stroke="#9CA3AF"
                    strokeDasharray="5 5"
                    strokeWidth={2}
                    label={{ value: 'Pronóstico →', position: 'insideTopRight', offset: -10, fill: '#6B7280', fontSize: 11, fontWeight: 500 }}
                  />

                  {/* Área de confianza (rango ±15%) */}
                  <Area
                    type="monotone"
                    dataKey="upper"
                    fill="url(#confidenceGradient)"
                    stroke="none"
                    isAnimationActive={false}
                    dot={false}
                    name="Rango superior"
                  />

                  {/* Línea de órdenes históricas (sólida) */}
                  <Line
                    type="monotone"
                    dataKey="orders"
                    stroke={COLORS.primary}
                    strokeWidth={3}
                    dot={false}
                    isAnimationActive={false}
                    name="Histórico"
                  />

                  {/* Línea de pronóstico (punteada) */}
                  <Line
                    type="monotone"
                    dataKey="forecast"
                    stroke={COLORS.primary}
                    strokeWidth={2}
                    strokeDasharray="6 3"
                    dot={false}
                    isAnimationActive={false}
                    name="Pronóstico"
                  />
                </ComposedChart>
              </ResponsiveContainer>
            </div>
          </div>
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
                <LabelList
                  dataKey="percentage"
                  position="bottom"
                  fontSize={11}
                  fill="#9CA3AF"
                  formatter={(value: number) => `${value > 0 ? '+' : ''}${value.toFixed(1)}%`}
                />
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
                    label={({ index, percent }) => {
                      const carrier = carrierData.carriers[index || 0];
                      return `${carrier?.displayName || 'Unknown'} ${((percent || 0) * 100).toFixed(0)}%`;
                    }}
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
