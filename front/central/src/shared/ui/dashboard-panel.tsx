'use client';

import { useMemo } from 'react';
import {
  BarChart,
  Bar,
  LineChart,
  Line,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from 'recharts';

interface DashboardPanelProps {
  isDark: boolean;
}

export const DashboardPanel = ({ isDark }: DashboardPanelProps) => {
  // Mock data para órdenes por estado
  const ordersByStatusData = [
    { status: 'Pendiente', orders: 24, fill: isDark ? '#f472b6' : '#ec4899' },
    { status: 'Enviada', orders: 42, fill: isDark ? '#8b5cf6' : '#8b5cf6' },
    { status: 'Entregada', orders: 89, fill: isDark ? '#14f5a0' : '#34d399' },
  ];

  // Mock data para últimos 7 días
  const last7DaysData = [
    { day: 'Lun', orders: 12 },
    { day: 'Mar', orders: 19 },
    { day: 'Mié', orders: 15 },
    { day: 'Jue', orders: 28 },
    { day: 'Vie', orders: 35 },
    { day: 'Sáb', orders: 42 },
    { day: 'Dom', orders: 24 },
  ];

  // Mock data para pedidos entregados vs no entregados
  const deliveryStatusData = [
    { name: 'Entregados', value: 89, fill: isDark ? '#14f5a0' : '#34d399' },
    { name: 'No Entregados', value: 66, fill: isDark ? '#f472b6' : '#ec4899' },
  ];

  // Calcular métricas
  const metrics = useMemo(() => {
    const totalOrders = ordersByStatusData.reduce((sum, item) => sum + item.orders, 0);
    const totalRevenue = totalOrders * 285000; // Promedio $285k por orden
    const avgOrderValue = totalRevenue / totalOrders;
    const last7DaysTotal = last7DaysData.reduce((sum, item) => sum + item.orders, 0);
    const avgDaily = Math.round(last7DaysTotal / 7);

    return {
      totalOrders,
      totalRevenue,
      avgOrderValue,
      avgDaily,
    };
  }, []);

  const colors = {
    text: isDark ? '#ffffff' : '#1a1a2e',
    textSecondary: isDark ? '#a1a1aa' : '#71717a',
    cardBg: isDark ? 'rgba(255, 255, 255, 0.08)' : 'rgba(255, 255, 255, 0.15)',
    cardBorder: isDark ? 'rgba(255, 255, 255, 0.12)' : 'rgba(255, 255, 255, 0.2)',
    gridColor: isDark ? 'rgba(255, 255, 255, 0.1)' : 'rgba(255, 255, 255, 0.15)',
    lineColor: isDark ? '#8b5cf6' : '#7c3aed',
  };

  return (
    <div className={`w-full h-full flex flex-col p-8 overflow-y-auto ${isDark ? 'bg-[#0f0a1e]' : 'bg-gradient-to-br from-[#5b21b6] via-[#7c3aed] to-[#4c1d95]'}`}>
      {/* Header */}
      <div className="mb-8">
        <h2 className={`text-2xl font-bold mb-2 ${isDark ? 'text-white' : 'text-white'}`}>
          {isDark ? 'Panel de Control' : 'Dashboard Logístico'}
        </h2>
        <p className={`text-sm ${isDark ? 'text-[#a1a1aa]' : 'text-white/70'}`}>
          {isDark ? 'Gestiona tus envíos en tiempo real' : 'Últimas 24 horas'}
        </p>
      </div>

      {/* Metrics Grid */}
      <div className="grid grid-cols-2 gap-4 mb-8">
        {/* Total Orders */}
        <div
          className="px-4 py-4 rounded-lg border backdrop-blur-sm"
          style={{
            backgroundColor: colors.cardBg,
            borderColor: colors.cardBorder,
          }}
        >
          <p className={`text-xs font-medium mb-2 ${isDark ? 'text-[#a1a1aa]' : 'text-white/60'}`}>
            Órdenes Totales
          </p>
          <p className={`text-2xl font-bold ${isDark ? 'text-white' : 'text-white'}`}>
            {metrics.totalOrders}
          </p>
          <p className={`text-xs mt-1 ${isDark ? 'text-[#14f5a0]' : 'text-[#34d399]'}`}>
            ↑ 12% vs semana anterior
          </p>
        </div>

        {/* Revenue */}
        <div
          className="px-4 py-4 rounded-lg border backdrop-blur-sm"
          style={{
            backgroundColor: colors.cardBg,
            borderColor: colors.cardBorder,
          }}
        >
          <p className={`text-xs font-medium mb-2 ${isDark ? 'text-[#a1a1aa]' : 'text-white/60'}`}>
            Ingresos
          </p>
          <p className={`text-2xl font-bold ${isDark ? 'text-white' : 'text-white'}`}>
            ${(metrics.totalRevenue / 1000000).toFixed(1)}M
          </p>
          <p className={`text-xs mt-1 ${isDark ? 'text-[#f472b6]' : 'text-[#ec4899]'}`}>
            ↑ 8% vs mes anterior
          </p>
        </div>

        {/* Avg Daily */}
        <div
          className="px-4 py-4 rounded-lg border backdrop-blur-sm"
          style={{
            backgroundColor: colors.cardBg,
            borderColor: colors.cardBorder,
          }}
        >
          <p className={`text-xs font-medium mb-2 ${isDark ? 'text-[#a1a1aa]' : 'text-white/60'}`}>
            Promedio Diario
          </p>
          <p className={`text-2xl font-bold ${isDark ? 'text-white' : 'text-white'}`}>
            {metrics.avgDaily}
          </p>
          <p className={`text-xs mt-1 ${isDark ? 'text-[#8b5cf6]' : 'text-[#c4b5fd]'}`}>
            órdenes/día
          </p>
        </div>

        {/* Avg Order Value */}
        <div
          className="px-4 py-4 rounded-lg border backdrop-blur-sm"
          style={{
            backgroundColor: colors.cardBg,
            borderColor: colors.cardBorder,
          }}
        >
          <p className={`text-xs font-medium mb-2 ${isDark ? 'text-[#a1a1aa]' : 'text-white/60'}`}>
            Ticket Promedio
          </p>
          <p className={`text-2xl font-bold ${isDark ? 'text-white' : 'text-white'}`}>
            ${(metrics.avgOrderValue / 1000).toFixed(0)}k
          </p>
          <p className={`text-xs mt-1 ${isDark ? 'text-[#c4b5fd]' : 'text-[#ddd6fe]'}`}>
            por orden
          </p>
        </div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 flex-1">
        {/* Bar Chart - Orders by Status */}
        <div
          className="p-4 rounded-lg border backdrop-blur-sm"
          style={{
            backgroundColor: colors.cardBg,
            borderColor: colors.cardBorder,
          }}
        >
          <h3 className={`text-sm font-semibold mb-4 ${isDark ? 'text-white' : 'text-white'}`}>
            Órdenes por Estado
          </h3>
          <ResponsiveContainer width="100%" height={220}>
            <BarChart data={ordersByStatusData}>
              <CartesianGrid
                strokeDasharray="3 3"
                stroke={colors.gridColor}
              />
              <XAxis
                dataKey="status"
                tick={{ fill: colors.textSecondary, fontSize: 12 }}
              />
              <YAxis
                tick={{ fill: colors.textSecondary, fontSize: 12 }}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: isDark ? '#1a1430' : '#5b21b6',
                  border: `1px solid ${colors.cardBorder}`,
                  borderRadius: '8px',
                  color: colors.text,
                }}
                labelStyle={{ color: colors.text }}
              />
              <Bar
                dataKey="orders"
                fill={colors.lineColor}
                radius={[8, 8, 0, 0]}
              />
            </BarChart>
          </ResponsiveContainer>
        </div>

        {/* Pie Chart - Delivered vs Not Delivered */}
        <div
          className="p-4 rounded-lg border backdrop-blur-sm"
          style={{
            backgroundColor: colors.cardBg,
            borderColor: colors.cardBorder,
          }}
        >
          <h3 className={`text-sm font-semibold mb-4 ${isDark ? 'text-white' : 'text-white'}`}>
            Entregados vs No Entregados
          </h3>
          <ResponsiveContainer width="100%" height={220}>
            <PieChart>
              <Pie
                data={deliveryStatusData}
                cx="50%"
                cy="50%"
                labelLine={false}
                label={(entry) => `${entry.name}: ${entry.value}`}
                outerRadius={80}
                fill="#8884d8"
                dataKey="value"
              >
                {deliveryStatusData.map((entry, index) => (
                  <Cell key={`cell-${index}`} fill={entry.fill} />
                ))}
              </Pie>
              <Tooltip
                contentStyle={{
                  backgroundColor: isDark ? '#1a1430' : '#5b21b6',
                  border: `1px solid ${colors.cardBorder}`,
                  borderRadius: '8px',
                  color: colors.text,
                }}
                labelStyle={{ color: colors.text }}
              />
            </PieChart>
          </ResponsiveContainer>
        </div>

        {/* Line Chart - Last 7 Days (Full width on small screens) */}
        <div
          className="p-4 rounded-lg border backdrop-blur-sm lg:col-span-2"
          style={{
            backgroundColor: colors.cardBg,
            borderColor: colors.cardBorder,
          }}
        >
          <h3 className={`text-sm font-semibold mb-4 ${isDark ? 'text-white' : 'text-white'}`}>
            Órdenes Últimos 7 Días
          </h3>
          <ResponsiveContainer width="100%" height={220}>
            <LineChart data={last7DaysData}>
              <CartesianGrid
                strokeDasharray="3 3"
                stroke={colors.gridColor}
              />
              <XAxis
                dataKey="day"
                tick={{ fill: colors.textSecondary, fontSize: 12 }}
              />
              <YAxis
                tick={{ fill: colors.textSecondary, fontSize: 12 }}
              />
              <Tooltip
                contentStyle={{
                  backgroundColor: isDark ? '#1a1430' : '#5b21b6',
                  border: `1px solid ${colors.cardBorder}`,
                  borderRadius: '8px',
                  color: colors.text,
                }}
                labelStyle={{ color: colors.text }}
              />
              <Line
                type="monotone"
                dataKey="orders"
                stroke={colors.lineColor}
                strokeWidth={3}
                dot={{ fill: colors.lineColor, r: 5 }}
                activeDot={{ r: 7 }}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>

      {/* Footer */}
      <div className={`mt-8 pt-4 border-t text-center text-xs ${isDark ? 'border-white/10 text-[#6d5ba0]' : 'border-white/20 text-white/60'}`}>
        Datos actualizados hace 2 minutos
      </div>
    </div>
  );
};
