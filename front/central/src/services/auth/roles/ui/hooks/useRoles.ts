import { useState, useEffect, useCallback } from 'react';
import { getRolesAction, deleteRoleAction } from '../../infra/actions';
import { Role } from '../../domain/types';

export const useRoles = () => {
    const [roles, setRoles] = useState<Role[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Filters
    const [searchName, setSearchName] = useState('');
    const [filterScope, setFilterScope] = useState<string>('');
    const [filterBusinessType, setFilterBusinessType] = useState<string>('');
    const [filterLevel, setFilterLevel] = useState<string>('');
    const [filterIsSystem, setFilterIsSystem] = useState<string>('');

    const fetchRoles = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getRolesAction({
                name: searchName || undefined,
                scope_id: filterScope ? Number(filterScope) : undefined,
                business_type_id: filterBusinessType ? Number(filterBusinessType) : undefined,
                level: filterLevel ? Number(filterLevel) : undefined,
                is_system: filterIsSystem === 'true' ? true : filterIsSystem === 'false' ? false : undefined,
            });
            setRoles(response.data || []);
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error fetching roles';
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    }, [searchName, filterScope, filterBusinessType, filterLevel, filterIsSystem]);

    const deleteRole = async (id: number) => {
        try {
            await deleteRoleAction(id);
            fetchRoles();
            return true;
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error deleting role';
            setError(errorMessage);
            return false;
        }
    };

    useEffect(() => {
        fetchRoles();
    }, [fetchRoles]);

    return {
        roles,
        loading,
        error,
        searchName,
        setSearchName,
        filterScope,
        setFilterScope,
        filterBusinessType,
        setFilterBusinessType,
        filterLevel,
        setFilterLevel,
        filterIsSystem,
        setFilterIsSystem,
        deleteRole,
        refresh: fetchRoles,
        setError
    };
};
