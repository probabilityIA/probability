import { useState, useEffect, useCallback } from 'react';
import { getUsersAction, deleteUserAction } from '../../infra/actions';
import { User, Pagination } from '../../domain/types';

export const useUsers = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [pagination, setPagination] = useState<Pagination | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Filters
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(10);
    const [searchName, setSearchName] = useState('');
    const [searchEmail, setSearchEmail] = useState('');
    const [filterIsActive, setFilterIsActive] = useState<string>('');
    const [filterRoleId, setFilterRoleId] = useState<string>('');
    const [filterBusinessId, setFilterBusinessId] = useState<string>('');

    const fetchUsers = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getUsersAction({
                page,
                page_size: pageSize,
                name: searchName || undefined,
                email: searchEmail || undefined,
                is_active: filterIsActive === 'true' ? true : filterIsActive === 'false' ? false : undefined,
                role_id: filterRoleId ? Number(filterRoleId) : undefined,
                business_id: filterBusinessId ? Number(filterBusinessId) : undefined,
            });
            setUsers(response.data || []);
            setPagination(response.pagination);
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error fetching users';
            setError(errorMessage);
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, searchName, searchEmail, filterIsActive, filterRoleId, filterBusinessId]);

    const deleteUser = async (id: number) => {
        try {
            await deleteUserAction(id);
            fetchUsers();
            return true;
        } catch (err: unknown) {
            const errorMessage = err instanceof Error ? err.message : 'Error deleting user';
            setError(errorMessage);
            return false;
        }
    };

    useEffect(() => {
        fetchUsers();
    }, [fetchUsers]);

    return {
        users,
        pagination,
        loading,
        error,
        page,
        setPage,
        pageSize,
        setPageSize,
        searchName,
        setSearchName,
        searchEmail,
        setSearchEmail,
        filterIsActive,
        setFilterIsActive,
        filterRoleId,
        setFilterRoleId,
        filterBusinessId,
        setFilterBusinessId,
        deleteUser,
        refresh: fetchUsers,
        setError
    };
};
