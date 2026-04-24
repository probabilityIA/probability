'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    listCountPlansAction,
    listCountTasksAction,
    listCountLinesAction,
    listDiscrepanciesAction,
    exportKardexAction,
} from '../../infra/actions/audit';
import {
    CycleCountLine,
    CycleCountPlan,
    CycleCountTask,
    GetCountLinesParams,
    GetCountPlansParams,
    GetCountTasksParams,
    GetDiscrepanciesParams,
    InventoryDiscrepancy,
    KardexExportResult,
    KardexQueryInput,
} from '../../domain/audit-types';
import { getActionError } from '@/shared/utils/action-result';

export function useCountPlans(initial: GetCountPlansParams = {}) {
    const [plans, setPlans] = useState<CycleCountPlan[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listCountPlansAction({ page, page_size: pageSize, warehouse_id: initial.warehouse_id, active_only: initial.active_only, business_id: initial.business_id });
            setPlans(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar planes de conteo'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, initial.warehouse_id, initial.active_only, initial.business_id]);

    useEffect(() => { fetch(); }, [fetch]);

    return { plans, loading, error, page, pageSize, total, totalPages, setPage, setPageSize, refresh: fetch };
}

export function useCountTasks(initial: GetCountTasksParams = {}) {
    const [tasks, setTasks] = useState<CycleCountTask[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetCountTasksParams, 'page' | 'page_size'>>({
        warehouse_id: initial.warehouse_id,
        plan_id: initial.plan_id,
        status: initial.status,
        business_id: initial.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listCountTasksAction({ page, page_size: pageSize, ...filters });
            setTasks(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar tareas de conteo'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => { fetch(); }, [fetch]);

    return {
        tasks, loading, error, page, pageSize, total, totalPages, filters,
        setPage,
        setPageSize: (n: number) => { setPageSize(n); setPage(1); },
        setFilters: (next: Partial<typeof filters>) => { setFilters((f) => ({ ...f, ...next })); setPage(1); },
        refresh: fetch,
    };
}

export function useCountLines(params: GetCountLinesParams) {
    const [lines, setLines] = useState<CycleCountLine[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [total, setTotal] = useState(0);

    const fetch = useCallback(async () => {
        if (!params.task_id) return;
        setLoading(true);
        setError(null);
        try {
            const r = await listCountLinesAction(params);
            setLines(r.data || []);
            setTotal(r.total || 0);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar lineas de conteo'));
        } finally {
            setLoading(false);
        }
    }, [params.task_id, params.page, params.page_size, params.status, params.business_id]);

    useEffect(() => { fetch(); }, [fetch]);

    return { lines, loading, error, total, refresh: fetch };
}

export function useDiscrepancies(initial: GetDiscrepanciesParams = {}) {
    const [discrepancies, setDiscrepancies] = useState<InventoryDiscrepancy[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetDiscrepanciesParams, 'page' | 'page_size'>>({
        task_id: initial.task_id,
        status: initial.status,
        business_id: initial.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listDiscrepanciesAction({ page, page_size: pageSize, ...filters });
            setDiscrepancies(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar discrepancias'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => { fetch(); }, [fetch]);

    return {
        discrepancies, loading, error, page, pageSize, total, totalPages, filters,
        setPage,
        setPageSize: (n: number) => { setPageSize(n); setPage(1); },
        setFilters: (next: Partial<typeof filters>) => { setFilters((f) => ({ ...f, ...next })); setPage(1); },
        refresh: fetch,
    };
}

interface UseKardexState {
    data: KardexExportResult | null;
    loading: boolean;
    error: string | null;
}

export function useKardex() {
    const [state, setState] = useState<UseKardexState>({ data: null, loading: false, error: null });

    const load = useCallback(async (input: KardexQueryInput) => {
        setState({ data: null, loading: true, error: null });
        try {
            const data = await exportKardexAction(input);
            setState({ data, loading: false, error: null });
            return data;
        } catch (e: any) {
            setState({ data: null, loading: false, error: getActionError(e, 'Error al cargar kardex') });
            return null;
        }
    }, []);

    return { ...state, load };
}
