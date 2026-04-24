'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    listLPNsAction,
    scanAction,
    listSyncLogsAction,
} from '../../infra/actions/capture';
import {
    GetLPNsParams,
    GetSyncLogsParams,
    InventorySyncLog,
    LicensePlate,
    ScanInput,
    ScanResult,
} from '../../domain/capture-types';
import { getActionError } from '@/shared/utils/action-result';

export function useLPNs(initial: GetLPNsParams = {}) {
    const [lpns, setLPNs] = useState<LicensePlate[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetLPNsParams, 'page' | 'page_size'>>({
        lpn_type: initial.lpn_type,
        status: initial.status,
        location_id: initial.location_id,
        business_id: initial.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listLPNsAction({ page, page_size: pageSize, ...filters });
            setLPNs(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar LPNs'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => { fetch(); }, [fetch]);

    return {
        lpns, loading, error, page, pageSize, total, totalPages, filters,
        setPage,
        setPageSize: (n: number) => { setPageSize(n); setPage(1); },
        setFilters: (next: Partial<typeof filters>) => { setFilters((f) => ({ ...f, ...next })); setPage(1); },
        refresh: fetch,
    };
}

interface UseScanState {
    result: ScanResult | null;
    loading: boolean;
    error: string | null;
}

export function useScan(businessId?: number) {
    const [state, setState] = useState<UseScanState>({ result: null, loading: false, error: null });

    const scan = useCallback(async (input: ScanInput) => {
        setState({ result: null, loading: true, error: null });
        const response = await scanAction(input, businessId);
        if (response.success) {
            setState({ result: response.data, loading: false, error: null });
            return response.data;
        }
        setState({ result: null, loading: false, error: response.error });
        return null;
    }, [businessId]);

    const reset = useCallback(() => setState({ result: null, loading: false, error: null }), []);

    return { ...state, scan, reset };
}

export function useSyncLogs(initial: GetSyncLogsParams = {}) {
    const [logs, setLogs] = useState<InventorySyncLog[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetSyncLogsParams, 'page' | 'page_size'>>({
        integration_id: initial.integration_id,
        direction: initial.direction,
        status: initial.status,
        business_id: initial.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listSyncLogsAction({ page, page_size: pageSize, ...filters });
            setLogs(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar logs de sincronizacion'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => { fetch(); }, [fetch]);

    return {
        logs, loading, error, page, pageSize, total, totalPages, filters,
        setPage,
        setPageSize: (n: number) => { setPageSize(n); setPage(1); },
        setFilters: (next: Partial<typeof filters>) => { setFilters((f) => ({ ...f, ...next })); setPage(1); },
        refresh: fetch,
    };
}
