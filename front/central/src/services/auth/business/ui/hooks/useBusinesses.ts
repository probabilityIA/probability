import { useState, useEffect, useCallback } from 'react';
import { getBusinessesAction, deleteBusinessAction, getBusinessTypesAction } from '../../infra/actions';
import { Business, BusinessType } from '../../domain/types';

export const useBusinesses = () => {
    const [businesses, setBusinesses] = useState<Business[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);

    // Filters
    const [searchName, setSearchName] = useState('');
    const [filterType, setFilterType] = useState<string>('');
    const [businessTypes, setBusinessTypes] = useState<BusinessType[]>([]);

    const fetchBusinesses = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getBusinessesAction({
                page,
                per_page: 10,
                name: searchName || undefined,
                business_type_id: filterType ? Number(filterType) : undefined,
            });
            setBusinesses(response.data || []);
            setTotalPages(response.pagination.last_page);
        } catch (err: any) {
            setError(err.message || 'Error fetching businesses');
        } finally {
            setLoading(false);
        }
    }, [page, searchName, filterType]);

    const fetchTypes = useCallback(async () => {
        try {
            const res = await getBusinessTypesAction();
            setBusinessTypes(res.data);
        } catch (e) {
            console.error("Failed to load business types", e);
        }
    }, []);

    const deleteBusiness = async (id: number) => {
        try {
            await deleteBusinessAction(id);
            fetchBusinesses();
            return true;
        } catch (err: any) {
            setError(err.message || 'Error deleting business');
            return false;
        }
    };

    useEffect(() => {
        fetchTypes();
    }, [fetchTypes]);

    useEffect(() => {
        fetchBusinesses();
    }, [fetchBusinesses]);

    return {
        businesses,
        loading,
        error,
        page,
        setPage,
        totalPages,
        searchName,
        setSearchName,
        filterType,
        setFilterType,
        businessTypes,
        deleteBusiness,
        refresh: fetchBusinesses,
        setError
    };
};
