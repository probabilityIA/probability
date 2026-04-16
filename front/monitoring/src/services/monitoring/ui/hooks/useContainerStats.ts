'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import type { ContainerStats } from '../../domain/types';
import { getToken } from '@/shared/lib/token';

interface UseContainerStatsOptions {
    containerId: string;
    interval?: number;
    enabled?: boolean;
}

export function useContainerStats({ containerId, interval = 5000, enabled = true }: UseContainerStatsOptions) {
    const [stats, setStats] = useState<ContainerStats | null>(null);
    const [loading, setLoading] = useState(true);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const fetchStats = useCallback(async () => {
        if (!containerId || !enabled) return;
        const token = getToken();
        if (!token) return;
        try {
            const res = await fetch(`/api/stats/${containerId}?token=${encodeURIComponent(token)}`);
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
