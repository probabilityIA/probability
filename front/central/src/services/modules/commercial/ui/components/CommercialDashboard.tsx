'use client';

import { useProspectsStats } from '../hooks/useProspectsStats';
import { ProspectsCharts } from './ProspectsCharts';
import { ProspectsList } from './ProspectsList';

export function CommercialDashboard() {
  const { stats, loading } = useProspectsStats();

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-xl font-bold text-gray-900 dark:text-white">Comercial</h1>
        <p className="text-sm text-gray-500 dark:text-gray-400">
          Prospectos comerciales detectados para el equipo de ventas.
        </p>
      </div>

      {loading && (
        <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-6 text-center text-sm text-gray-400">
          Cargando estadisticas...
        </div>
      )}
      {!loading && stats && <ProspectsCharts stats={stats} />}

      <ProspectsList />
    </div>
  );
}
