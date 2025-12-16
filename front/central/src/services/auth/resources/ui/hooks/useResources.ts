'use client';

import { useState, useEffect, useCallback } from 'react';
import { Resource, GetResourcesParams } from '../../domain/types';
import { getResourcesAction, deleteResourceAction } from '../../infra/actions';

export const useResources = () => {
    const [resources, setResources] = useState<Resource[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const [filters, setFilters] = useState<GetResourcesParams>({});

    const loadResources = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getResourcesAction({
                ...filters,
                page,
                page_size: pageSize,
            });

            if (response.success && response.data) {
                setResources(response.data.resources || []);
                setTotal(response.data.total || 0);
                setTotalPages(response.data.total_pages || 1);
            } else {
                setError(response.message || 'Error al cargar recursos');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar recursos');
        } finally {
            setLoading(false);
        }
    }, [filters, page, pageSize]);

    useEffect(() => {
        loadResources();
    }, [loadResources]);

    const deleteResource = async (id: number): Promise<boolean> => {
        try {
            const response = await deleteResourceAction(id);
            if (response.success) {
                loadResources();
                return true;
            } else {
                setError(response.message || 'Error al eliminar recurso');
                return false;
            }
        } catch (err: any) {
            setError(err.message || 'Error al eliminar recurso');
            return false;
        }
    };

    return {
        resources,
        loading,
        error,
        page,
        setPage,
        pageSize,
        setPageSize,
        totalPages,
        total,
        filters,
        setFilters,
        deleteResource,
        refresh: loadResources,
        setError,
    };
};
