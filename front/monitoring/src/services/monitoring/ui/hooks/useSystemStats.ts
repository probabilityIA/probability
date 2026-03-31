'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import type { SystemStats } from '../../domain/types';
import { getTokenAction } from '../../infra/actions';

export function useSystemStats(interval = 5000) {
    const [stats, setStats] = useState<SystemStats | null>(null);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);
    const tokenRef = useRef<string | null>(null);

    const fetchStats = useCallback(async () => {
        if (!tokenRef.current) {
            tokenRef.current = await getTokenAction();
        }
        if (!tokenRef.current) return;
        try {
            const res = await fetch(`/api/system?token=${encodeURIComponent(tokenRef.current)}`);
            if (res.ok) setStats(await res.json());
        } catch { /* non-critical */ }
    }, []);

    useEffect(() => {
        fetchStats();
        timerRef.current = setInterval(fetchStats, interval);
        return () => { if (timerRef.current) clearInterval(timerRef.current); };
    }, [fetchStats, interval]);

    return stats;
}
