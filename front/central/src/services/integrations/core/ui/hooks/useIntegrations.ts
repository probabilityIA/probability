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
import { Integration } from '../../domain/types';

export const useIntegrations = () => {
    const [integrations, setIntegrations] = useState<Integration[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);

    // Filters
    const [search, setSearch] = useState('');
    const [filterType, setFilterType] = useState<string>('');
    const [filterCategory, setFilterCategory] = useState<string>('');

    const fetchIntegrations = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getIntegrationsAction({
                page,
                page_size: 10,
                search: search || undefined,
                type: filterType || undefined,
                category: filterCategory || undefined,
            });
            setIntegrations(response.data || []);
            setTotalPages(response.total_pages);
        } catch (err: any) {
            console.error('Error fetching integrations:', err);
            setError(err.message || 'Error fetching integrations');
        } finally {
            setLoading(false);
        }
    }, [page, search, filterType, filterCategory]);

    const deleteIntegration = async (id: number) => {
        try {
            await deleteIntegrationAction(id);
            fetchIntegrations();
            return true;
        } catch (err: any) {
            console.error('Error deleting integration:', err);
            setError(err.message || 'Error deleting integration');
            return false;
        }
    };

    const toggleActive = async (id: number, isActive: boolean) => {
        try {
            if (isActive) {
                await deactivateIntegrationAction(id);
            } else {
                await activateIntegrationAction(id);
            }
            fetchIntegrations();
            return true;
        } catch (err: any) {
            console.error('Error toggling integration status:', err);
            setError(err.message || 'Error updating status');
            return false;
        }
    };

    const setAsDefault = async (id: number) => {
        try {
            await setAsDefaultAction(id);
            fetchIntegrations();
            return true;
        } catch (err: any) {
            console.error('Error setting default integration:', err);
            setError(err.message || 'Error setting default');
            return false;
        }
    };

    const testConnection = async (id: number) => {
        try {
            const res = await testConnectionAction(id);
            return res;
        } catch (err: any) {
            console.error('Error testing connection:', err);
            return { success: false, message: err.message || 'Error testing connection' };
        }
    };

    const syncOrders = async (id: number) => {
        try {
            const res = await syncOrdersAction(id);
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
        error,
        page,
        setPage,
        totalPages,
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
        refresh: fetchIntegrations,

        setError
    };
};
