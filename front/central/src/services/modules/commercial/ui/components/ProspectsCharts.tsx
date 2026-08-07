'use client';

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
} from 'recharts';
import { ProspectStats } from '../../domain/types';

const CATEGORICAL = ['#7C3AED', '#06B6D4', '#F59E0B', '#10B981', '#EF4444', '#6366F1', '#EC4899', '#84CC16'];
const SEQUENTIAL = '#7C3AED';
const STATUS_COLORS: Record<string, string> = {
  calificado: '#10B981',
  descartado: '#9CA3AF',
  nuevo: '#7C3AED',
  contactado: '#06B6D4',
  cliente: '#F59E0B',
};

function StatTile({ label, value, hint }: { label: string; value: string; hint?: string }) {
  return (
    <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
      <p className="text-xs font-medium text-gray-500 dark:text-gray-400">{label}</p>
      <p className="mt-1 text-2xl font-bold text-gray-900 dark:text-white">{value}</p>
      {hint && <p className="mt-0.5 text-xs text-gray-400 dark:text-gray-500 truncate">{hint}</p>}
    </div>
  );
}

function ChartCard({ title, children }: { title: string; children: React.ReactNode }) {
  return (
    <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-4">
      <h3 className="text-sm font-semibold text-gray-900 dark:text-white mb-3">{title}</h3>
      {children}
    </div>
  );
}

function CustomTooltip({ active, payload, label }: any) {
  if (!active || !payload || !payload.length) return null;
  return (
    <div className="rounded-lg border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 px-3 py-2 shadow-sm text-xs">
      <p className="font-semibold text-gray-900 dark:text-white">{label}</p>
      <p className="text-gray-600 dark:text-gray-300">{payload[0].value} prospectos</p>
    </div>
  );
}

export function ProspectsCharts({ stats }: { stats: ProspectStats }) {
  const platformData = stats.by_platform.map((item) => ({ name: item.key, count: item.count }));
  const cityData = stats.by_city.map((item) => ({ name: item.key, count: item.count }));
  const scoreData = stats.score_buckets.map((item) => ({ name: item.range_label, count: item.count }));
  const statusData = stats.by_status.map((item) => ({ name: item.key, count: item.count }));

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-6 gap-3">
        <StatTile label="Total prospectos" value={stats.total.toLocaleString('es-CO')} />
        <StatTile label="Calificados" value={(stats.by_status.find(s => s.key === 'calificado')?.count ?? 0).toLocaleString('es-CO')} />
        <StatTile label="Con pago contraentrega" value={stats.has_cod_count.toLocaleString('es-CO')} />
        <StatTile label="Multi-canal" value={stats.multi_channel_count.toLocaleString('es-CO')} />
        <StatTile label="Score promedio" value={stats.average_score.toFixed(1)} />
        <StatTile label="Mejor prospecto" value={stats.top_score.toFixed(0)} hint={stats.top_score_domain} />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <ChartCard title="Por plataforma">
          <ResponsiveContainer width="100%" height={240}>
            <BarChart data={platformData} margin={{ top: 8, right: 8, left: 0, bottom: 8 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
              <XAxis dataKey="name" tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} allowDecimals={false} />
              <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(124,58,237,0.06)' }} />
              <Bar dataKey="count" radius={[4, 4, 0, 0]}>
                {platformData.map((_, index) => (
                  <Cell key={index} fill={CATEGORICAL[index % CATEGORICAL.length]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>

        <ChartCard title="Por estado">
          <ResponsiveContainer width="100%" height={240}>
            <BarChart data={statusData} margin={{ top: 8, right: 8, left: 0, bottom: 8 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
              <XAxis dataKey="name" tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} allowDecimals={false} />
              <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(124,58,237,0.06)' }} />
              <Bar dataKey="count" radius={[4, 4, 0, 0]}>
                {statusData.map((entry, index) => (
                  <Cell key={index} fill={STATUS_COLORS[entry.name] || CATEGORICAL[index % CATEGORICAL.length]} />
                ))}
              </Bar>
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>

        <ChartCard title="Top ciudades">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={cityData} layout="vertical" margin={{ top: 8, right: 16, left: 8, bottom: 8 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" horizontal={false} />
              <XAxis type="number" tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} allowDecimals={false} />
              <YAxis type="category" dataKey="name" tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} width={90} />
              <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(124,58,237,0.06)' }} />
              <Bar dataKey="count" fill={SEQUENTIAL} radius={[0, 4, 4, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>

        <ChartCard title="Distribucion de score">
          <ResponsiveContainer width="100%" height={260}>
            <BarChart data={scoreData} margin={{ top: 8, right: 8, left: 0, bottom: 8 }}>
              <CartesianGrid strokeDasharray="3 3" stroke="#f0f0f0" vertical={false} />
              <XAxis dataKey="name" tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} />
              <YAxis tick={{ fontSize: 11, fill: '#9CA3AF' }} axisLine={false} tickLine={false} allowDecimals={false} />
              <Tooltip content={<CustomTooltip />} cursor={{ fill: 'rgba(124,58,237,0.06)' }} />
              <Bar dataKey="count" fill={SEQUENTIAL} radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
        </ChartCard>
      </div>
    </div>
  );
}
