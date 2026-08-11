'use client';

import { useCallback, useEffect, useState } from 'react';
import { getProspectsStatsAction } from '../../infra/actions';
import { ProspectStats } from '../../domain/types';

export function useProspectsStats() {
  const [stats, setStats] = useState<ProspectStats | null>(null);
  const [loading, setLoading] = useState(true);

  const refetch = useCallback(() => {
    getProspectsStatsAction()
      .then((res) => setStats(res.data))
      .catch((error) => console.error('Error al cargar estadisticas de prospectos:', error));
  }, []);

  useEffect(() => {
    let cancelled = false;
    getProspectsStatsAction()
      .then((res) => {
        if (!cancelled) setStats(res.data);
      })
      .catch((error) => {
        console.error('Error al cargar estadisticas de prospectos:', error);
      })
      .finally(() => {
        if (!cancelled) setLoading(false);
      });
    return () => { cancelled = true; };
  }, []);

  return { stats, loading, refetch };
}
