'use client';

import { useState, useEffect, useRef, useCallback } from 'react';
import type { SystemStats } from '../../domain/types';
import { getToken } from '@/shared/lib/token';

export function useSystemStats(interval = 5000) {
    const [stats, setStats] = useState<SystemStats | null>(null);
    const timerRef = useRef<ReturnType<typeof setInterval> | null>(null);

    const fetchStats = useCallback(async () => {
        const token = getToken();
        if (!token) return;
        try {
            const res = await fetch(`/api/system?token=${encodeURIComponent(token)}`);
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
