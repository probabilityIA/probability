'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import type { SystemStats } from '../../domain/types';

export function useSystemStats(interval = 5000) {
    const [stats, setStats] = useState<SystemStats | null>(null);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const fetchStats = useCallback(async () => {
        try {
            const res = await fetch('/api/system');
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
