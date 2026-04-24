'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    listPutawayRulesAction,
    listPutawaySuggestionsAction,
    listReplenishmentTasksAction,
    listCrossDockLinksAction,
    listVelocitiesAction,
} from '../../infra/actions/operations';
import {
    CrossDockLink,
    GetCrossDockLinksParams,
    GetPutawayRulesParams,
    GetPutawaySuggestionsParams,
    GetReplenishmentTasksParams,
    GetVelocitiesParams,
    ProductVelocity,
    PutawayRule,
    PutawaySuggestion,
    ReplenishmentTask,
} from '../../domain/operations-types';
import { getActionError } from '@/shared/utils/action-result';

export function usePutawayRules(initial: GetPutawayRulesParams = {}) {
    const [rules, setRules] = useState<PutawayRule[]>([]);
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
            const r = await listPutawayRulesAction({ page, page_size: pageSize, active_only: initial.active_only, business_id: initial.business_id });
            setRules(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar reglas'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, initial.active_only, initial.business_id]);

    useEffect(() => { fetch(); }, [fetch]);

    return { rules, loading, error, page, pageSize, total, totalPages, setPage, setPageSize, refresh: fetch };
}

export function usePutawaySuggestions(initial: GetPutawaySuggestionsParams = {}) {
    const [suggestions, setSuggestions] = useState<PutawaySuggestion[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [status, setStatus] = useState(initial.status || '');

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listPutawaySuggestionsAction({ page, page_size: pageSize, status, business_id: initial.business_id });
            setSuggestions(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar sugerencias'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, status, initial.business_id]);

    useEffect(() => { fetch(); }, [fetch]);

    return { suggestions, loading, error, page, pageSize, total, totalPages, status, setStatus, setPage, setPageSize, refresh: fetch };
}

export function useReplenishmentTasks(initial: GetReplenishmentTasksParams = {}) {
    const [tasks, setTasks] = useState<ReplenishmentTask[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetReplenishmentTasksParams, 'page' | 'page_size'>>({
        warehouse_id: initial.warehouse_id,
        status: initial.status,
        assigned_to: initial.assigned_to,
        business_id: initial.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listReplenishmentTasksAction({ page, page_size: pageSize, ...filters });
            setTasks(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar tareas'));
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

export function useCrossDockLinks(initial: GetCrossDockLinksParams = {}) {
    const [links, setLinks] = useState<CrossDockLink[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(initial.page || 1);
    const [pageSize, setPageSize] = useState(initial.page_size || 10);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [filters, setFilters] = useState<Omit<GetCrossDockLinksParams, 'page' | 'page_size'>>({
        outbound_order_id: initial.outbound_order_id,
        status: initial.status,
        business_id: initial.business_id,
    });

    const fetch = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const r = await listCrossDockLinksAction({ page, page_size: pageSize, ...filters });
            setLinks(r.data || []);
            setTotal(r.total || 0);
            setTotalPages(r.total_pages || 1);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar enlaces cross-dock'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, filters]);

    useEffect(() => { fetch(); }, [fetch]);

    return { links, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh: fetch };
}

export function useVelocities(initial: GetVelocitiesParams) {
    const [velocities, setVelocities] = useState<ProductVelocity[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetch = useCallback(async () => {
        if (!initial.warehouse_id) return;
        setLoading(true);
        setError(null);
        try {
            const r = await listVelocitiesAction(initial);
            setVelocities(r.data || []);
        } catch (e: any) {
            setError(getActionError(e, 'Error al cargar rotacion'));
        } finally {
            setLoading(false);
        }
    }, [initial.warehouse_id, initial.period, initial.rank, initial.limit, initial.business_id]);

    useEffect(() => { fetch(); }, [fetch]);

    return { velocities, loading, error, refresh: fetch };
}
