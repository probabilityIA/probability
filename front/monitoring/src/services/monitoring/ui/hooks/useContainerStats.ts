'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import type { ContainerStats } from '../../domain/types';

interface UseContainerStatsOptions {
    containerId: string;
    interval?: number; // ms, default 5000
    enabled?: boolean;
}

export function useContainerStats({ containerId, interval = 5000, enabled = true }: UseContainerStatsOptions) {
    const [stats, setStats] = useState<ContainerStats | null>(null);
    const [loading, setLoading] = useState(true);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const fetchStats = useCallback(async () => {
        if (!containerId || !enabled) return;
        try {
            const res = await fetch(`/api/stats/${containerId}`);
            if (res.ok) {
                const data = await res.json();
                setStats(data);
            }
        } catch {
            // Silently fail - stats are non-critical
        } finally {
            setLoading(false);
        }
    }, [containerId, enabled]);

    useEffect(() => {
        if (!enabled) return;

        fetchStats();
        timerRef.current = setInterval(fetchStats, interval);

        return () => {
            if (timerRef.current) clearInterval(timerRef.current);
        };
    }, [fetchStats, interval, enabled]);

    return { stats, loading };
}
