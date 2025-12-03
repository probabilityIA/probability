import { useState, useEffect, useCallback } from 'react';
import { getPermissionsAction, deletePermissionAction } from '../../infra/actions';
import { Permission } from '../../domain/types';

export const usePermissions = () => {
    const [permissions, setPermissions] = useState<Permission[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Filters
    const [searchName, setSearchName] = useState('');
    const [filterScope, setFilterScope] = useState<string>('');
    const [filterBusinessType, setFilterBusinessType] = useState<string>('');

    const fetchPermissions = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getPermissionsAction({
                name: searchName || undefined,
                scope_id: filterScope ? Number(filterScope) : undefined,
                business_type_id: filterBusinessType ? Number(filterBusinessType) : undefined,
            });
            setPermissions(response.data || []);
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error fetching permissions';
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    }, [searchName, filterScope, filterBusinessType]);

    const deletePermission = async (id: number) => {
        try {
            await deletePermissionAction(id);
            fetchPermissions();
            return true;
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error deleting permission';
            setError(errorMessage);
            return false;
        }
    };

    useEffect(() => {
        fetchPermissions();
    }, [fetchPermissions]);

    return {
        permissions,
        loading,
        error,
        searchName,
        setSearchName,
        filterScope,
        setFilterScope,
        filterBusinessType,
        setFilterBusinessType,
        deletePermission,
        refresh: fetchPermissions,
        setError
    };
};
