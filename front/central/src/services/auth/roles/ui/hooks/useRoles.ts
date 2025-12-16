import { useState, useEffect, useCallback, useMemo } from 'react';
import { getRolesAction, deleteRoleAction } from '../../infra/actions';
import { Role, GetRolesParams } from '../../domain/types';

export const useRoles = () => {
    const [allRoles, setAllRoles] = useState<Role[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);

    // Filters
    const [filters, setFilters] = useState<GetRolesParams>({
        page: 1,
        page_size: 20,
    });

    const fetchRoles = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getRolesAction({
                name: filters.name,
                scope_id: filters.scope_id,
                business_type_id: filters.business_type_id,
                level: filters.level,
                is_system: filters.is_system,
            });
            setAllRoles(response.data || []);
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error fetching roles';
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    }, [filters.name, filters.scope_id, filters.business_type_id, filters.level, filters.is_system]);

    // Paginación del lado del cliente
    const { roles, totalPages, total } = useMemo(() => {
        const startIndex = (page - 1) * pageSize;
        const endIndex = startIndex + pageSize;
        const paginatedRoles = allRoles.slice(startIndex, endIndex);
        const total = allRoles.length;
        const totalPages = Math.ceil(total / pageSize);

        return {
            roles: paginatedRoles,
            totalPages,
            total,
        };
    }, [allRoles, page, pageSize]);

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
        allRoles,
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
        deleteRole,
        refresh: fetchRoles,
        setError
    };
};
