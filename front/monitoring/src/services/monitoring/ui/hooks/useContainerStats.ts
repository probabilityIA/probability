'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import type { ContainerStats } from '../../domain/types';
import { getTokenAction } from '../../infra/actions';

interface UseContainerStatsOptions {
    containerId: string;
    interval?: number;
    enabled?: boolean;
}

export function useContainerStats({ containerId, interval = 5000, enabled = true }: UseContainerStatsOptions) {
    const [stats, setStats] = useState<ContainerStats | null>(null);
    const [loading, setLoading] = useState(true);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
    const tokenRef = useRef<string | null>(null);

    const fetchStats = useCallback(async () => {
        if (!containerId || !enabled) return;
        if (!tokenRef.current) {
            tokenRef.current = await getTokenAction();
        }
        if (!tokenRef.current) return;
        try {
            const res = await fetch(`/api/stats/${containerId}?token=${encodeURIComponent(tokenRef.current)}`);
            if (res.ok) setStats(await res.json());
        } catch { /* non-critical */ }
        finally { setLoading(false); }
    }, [containerId, enabled]);

    useEffect(() => {
        if (!enabled) return;
        fetchStats();
        timerRef.current = setInterval(fetchStats, interval);
        return () => { if (timerRef.current) clearInterval(timerRef.current); };
    }, [fetchStats, interval, enabled]);

    return { stats, loading };
}
