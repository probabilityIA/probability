'use client';

import { useState, useEffect, useCallback } from 'react';
import {
    getIntegrationsAction,
    deleteIntegrationAction,
    activateIntegrationAction,
    deactivateIntegrationAction,
    setAsDefaultAction,
    testConnectionAction,
    syncOrdersAction
} from '../../infra/actions';
import { Integration, SyncOrdersParams } from '../../domain/types';
import { TokenStorage } from '@/shared/utils/token-storage';
import { getActionError } from '@/shared/utils/action-result';

const PAGE_SIZE = 10;

export const useIntegrations = (initialCategory: string = '', businessId: number | null = null) => {
    const [integrations, setIntegrations] = useState<Integration[]>([]);
    const [loading, setLoading] = useState(true);
    const [loadingMore, setLoadingMore] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);

    const [search, setSearch] = useState('');
    const [filterType, setFilterType] = useState<string>('');
    const [filterCategory, setFilterCategory] = useState<string>(initialCategory);

    const buildParams = useCallback((targetPage: number) => ({
        page: targetPage,
        page_size: PAGE_SIZE,
        search: search || undefined,
        type: filterType || undefined,
        category: filterCategory || undefined,
        business_id: businessId || undefined,
    }), [search, filterType, filterCategory, businessId]);

    const fetchIntegrations = useCallback(async () => {
        const isFirstPage = page === 1;
        if (isFirstPage) {
            setLoading(true);
        } else {
            setLoadingMore(true);
        }
        setError(null);
        try {
            const token = TokenStorage.getSessionToken();
            const response = await getIntegrationsAction(buildParams(page), token);
            const data = response.data || [];
            setIntegrations(prev => (isFirstPage ? data : [...prev, ...data]));
            setTotalPages(response.total_pages);
            setTotal(response.total ?? 0);
        } catch (err: any) {
            console.error('Error fetching integrations:', err);
            setError(getActionError(err, 'Error fetching integrations'));
        } finally {
            if (isFirstPage) {
                setLoading(false);
            } else {
                setLoadingMore(false);
            }
        }
    }, [page, buildParams]);

    const reloadAll = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const token = TokenStorage.getSessionToken();
            const accumulated: Integration[] = [];
            let lastTotalPages = 1;
            for (let p = 1; p <= page; p++) {
                const response = await getIntegrationsAction(buildParams(p), token);
                accumulated.push(...(response.data || []));
                lastTotalPages = response.total_pages;
                setTotal(response.total ?? 0);
                if (p >= response.total_pages) break;
            }
            setIntegrations(accumulated);
            setTotalPages(lastTotalPages);
        } catch (err: any) {
            console.error('Error fetching integrations:', err);
            setError(getActionError(err, 'Error fetching integrations'));
        } finally {
            setLoading(false);
        }
    }, [page, buildParams]);

    const hasMore = page < totalPages;

    const loadMore = useCallback(() => {
        setPage(prev => (prev < totalPages ? prev + 1 : prev));
    }, [totalPages]);

    const deleteIntegration = async (id: number) => {
        try {
            const token = TokenStorage.getSessionToken();
            const result = await deleteIntegrationAction(id, token);
            if (!result?.success) {
                const message = (result as any)?.message || 'Error deleting integration';
                setError(message);
                return false;
            }
            reloadAll();
            return true;
        } catch (err: any) {
            console.error('Error deleting integration:', err);
            setError(getActionError(err, 'Error deleting integration'));
            return false;
        }
    };

    const toggleActive = async (id: number, isActive: boolean) => {
        try {
            const token = TokenStorage.getSessionToken();
            const result = isActive
                ? await deactivateIntegrationAction(id, token)
                : await activateIntegrationAction(id, token);
            if (!result?.success) {
                const message = (result as any)?.message || 'Error updating status';
                setError(message);
                return false;
            }
            setIntegrations(prev => prev.map(i => (i.id === id ? { ...i, is_active: !isActive } : i)));
            return true;
        } catch (err: any) {
            console.error('Error toggling integration status:', err);
            setError(getActionError(err, 'Error updating status'));
            return false;
        }
    };

    const setAsDefault = async (id: number) => {
        try {
            const token = TokenStorage.getSessionToken();
            const result = await setAsDefaultAction(id, token);
            if (!result?.success) {
                const message = (result as any)?.message || 'Error setting default';
                setError(message);
                return false;
            }
            reloadAll();
            return true;
        } catch (err: any) {
            console.error('Error setting default integration:', err);
            setError(getActionError(err, 'Error setting default'));
            return false;
        }
    };

    const testConnection = async (id: number) => {
        try {
            const token = TokenStorage.getSessionToken();
            const res = await testConnectionAction(id, token);
            return res;
        } catch (err: any) {
            console.error('Error testing connection:', err);
            return { success: false, message: err.message || 'Error testing connection' };
        }
    };

    const syncOrders = async (id: number, params?: SyncOrdersParams) => {
        try {
            const token = TokenStorage.getSessionToken();
            const res = await syncOrdersAction(id, params, token);
            return res;
        } catch (err: any) {
            console.error('Error syncing orders:', err);
            return { success: false, message: err.message || 'Error syncing orders' };
        }
    };

    useEffect(() => {
        fetchIntegrations();
    }, [fetchIntegrations]);

    return {
        integrations,
        loading,
        loadingMore,
        error,
        page,
        setPage,
        totalPages,
        total,
        hasMore,
        loadMore,
        search,
        setSearch,
        filterType,
        setFilterType,
        filterCategory,
        setFilterCategory,
        deleteIntegration,
        toggleActive,
        setAsDefault,
        testConnection,
        syncOrders,
        refresh: reloadAll,

        setError
    };
};
