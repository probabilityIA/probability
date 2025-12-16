'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Button } from '@/shared/ui/button';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { Select } from '@/shared/ui/select';
import { Spinner } from '@/shared/ui/spinner';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { User, GetUsersParams } from '../../domain/types';
import { UserForm } from './UserForm';
import { getUsersAction, deleteUserAction, getUserByIdAction, assignRolesAction } from '../../infra/actions';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { getRolesAction } from '@/services/auth/roles/infra/actions';
import { getBusinessesAction } from '@/services/auth/business/infra/actions';
import { Role } from '@/services/auth/roles/domain/types';

export const UserList: React.FC = () => {
    const [users, setUsers] = useState<User[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const [pageSize, setPageSize] = useState(20);

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingUser, setEditingUser] = useState<User | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);
    const [loadingUser, setLoadingUser] = useState(false);

    // Estados para asignar rol
    const [showAssignRoleModal, setShowAssignRoleModal] = useState(false);
    const [assigningRoleUser, setAssigningRoleUser] = useState<User | null>(null);
    const [roles, setRoles] = useState<Role[]>([]);
    const [businesses, setBusinesses] = useState<{ id: number; name: string }[]>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<string>('');
    const [selectedRoleId, setSelectedRoleId] = useState<string>('');
    const [loadingRoles, setLoadingRoles] = useState(false);
    const [assigningRole, setAssigningRole] = useState(false);

    // Filters
    const [filters, setFilters] = useState<GetUsersParams>({
        page: 1,
        page_size: 20,
    });

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = [
        {
            key: 'name',
            label: 'Nombre',
            type: 'text',
            placeholder: 'Buscar por nombre...',
        },
        {
            key: 'email',
            label: 'Email',
            type: 'text',
            placeholder: 'Buscar por email...',
        },
        {
            key: 'phone',
            label: 'Teléfono',
            type: 'text',
            placeholder: 'Buscar por teléfono...',
        },
        {
            key: 'is_active',
            label: 'Estado',
            type: 'select',
            options: [
                { value: 'true', label: 'Activo' },
                { value: 'false', label: 'Inactivo' },
            ],
        },
        {
            key: 'role_id',
            label: 'ID de Rol',
            type: 'text',
            placeholder: 'Filtrar por ID de rol...',
        },
        {
            key: 'business_id',
            label: 'ID de Negocio',
            type: 'text',
            placeholder: 'Filtrar por ID de negocio...',
        },
    ];

    // Convertir filtros a ActiveFilter[]
    const activeFilters: ActiveFilter[] = useMemo(() => {
        const active: ActiveFilter[] = [];

        if (filters.name) {
            active.push({
                key: 'name',
                label: 'Nombre',
                value: filters.name,
                type: 'text',
            });
        }

        if (filters.email) {
            active.push({
                key: 'email',
                label: 'Email',
                value: filters.email,
                type: 'text',
            });
        }

        if (filters.phone) {
            active.push({
                key: 'phone',
                label: 'Teléfono',
                value: filters.phone,
                type: 'text',
            });
        }

        if (filters.is_active !== undefined) {
            active.push({
                key: 'is_active',
                label: 'Estado',
                value: filters.is_active ? 'Activo' : 'Inactivo',
                type: 'select',
            });
        }

        if (filters.role_id) {
            active.push({
                key: 'role_id',
                label: 'ID de Rol',
                value: String(filters.role_id),
                type: 'text',
            });
        }

        if (filters.business_id) {
            active.push({
                key: 'business_id',
                label: 'ID de Negocio',
                value: String(filters.business_id),
                type: 'text',
            });
        }

        return active;
    }, [filters]);

    // Manejar adición de filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };

            if (filterKey === 'is_active') {
                newFilters.is_active = value === 'true' || value === true;
            } else if (filterKey === 'role_id') {
                newFilters.role_id = value ? Number(value) : undefined;
            } else if (filterKey === 'business_id') {
                newFilters.business_id = value ? Number(value) : undefined;
            } else {
                (newFilters as any)[filterKey] = value;
            }

            return newFilters;
        });
    }, []);

    // Manejar eliminación de filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };
            delete (newFilters as any)[filterKey];
            return newFilters;
        });
    }, []);

    // Cargar usuarios
    const loadUsers = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getUsersAction(filters);
            if (response.success) {
                setUsers(response.data || []);
                if (response.pagination) {
                    setPage(response.pagination.current_page);
                    setTotalPages(response.pagination.last_page);
                    setTotal(response.pagination.total);
                    setPageSize(response.pagination.per_page);
                }
            } else {
                setError(response.message || 'Error al cargar usuarios');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar usuarios');
        } finally {
            setLoading(false);
        }
    }, [filters]);

    useEffect(() => {
        loadUsers();
    }, [loadUsers]);

    const handleDelete = async () => {
        if (deleteId) {
            try {
                const response = await deleteUserAction(deleteId);
                if (response.success) {
                    setDeleteId(null);
                    loadUsers();
                } else {
                    setError(response.message || 'Error al eliminar usuario');
                }
            } catch (err: any) {
                setError(err.message || 'Error al eliminar usuario');
            }
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingUser(null);
        loadUsers();
    };

    const handleEdit = async (user: User) => {
        setLoadingUser(true);
        try {
            const response = await getUserByIdAction(user.id);
            if (response.success && response.data) {
                setEditingUser(response.data);
                setShowCreateModal(true);
            } else {
                setEditingUser(user);
                setShowCreateModal(true);
            }
        } catch (error) {
            console.error('Error loading user for edit:', error);
            setEditingUser(user);
            setShowCreateModal(true);
        } finally {
            setLoadingUser(false);
        }
    };

    // Cargar roles y negocios cuando se abre el modal de asignar rol
    const handleOpenAssignRole = async (user: User) => {
        setAssigningRoleUser(user);
        setShowAssignRoleModal(true);
        setLoadingRoles(true);

        // Preseleccionar valores si el usuario ya tiene asignaciones
        const firstAssignment = user.business_role_assignments?.[0];
        if (firstAssignment) {
            setSelectedBusinessId(firstAssignment.business_id ? String(firstAssignment.business_id) : '');
            setSelectedRoleId(String(firstAssignment.role_id));
        } else {
            setSelectedBusinessId('');
            setSelectedRoleId('');
        }

        try {
            const [rolesResponse, businessesResponse] = await Promise.all([
                getRolesAction({ page_size: 100 }),
                getBusinessesAction({ page: 1, per_page: 100 })
            ]);

            if (rolesResponse.success && rolesResponse.data) {
                setRoles(rolesResponse.data);
            }

            if (businessesResponse.success && businessesResponse.data) {
                setBusinesses(businessesResponse.data.map((b: any) => ({ id: b.id, name: b.name })));
            }
        } catch (err: any) {
            console.error('Error loading roles/businesses:', err);
            setError('Error al cargar roles y negocios');
        } finally {
            setLoadingRoles(false);
        }
    };

    // Asignar rol al usuario
    const handleAssignRole = async () => {
        if (!assigningRoleUser || !selectedRoleId) {
            setError('Debe seleccionar un rol');
            return;
        }

        // Para usuarios platform no se requiere business_id
        const isPlatformUser = assigningRoleUser.scope_code === 'platform';
        if (!isPlatformUser && !selectedBusinessId) {
            setError('Debe seleccionar un negocio y un rol');
            return;
        }

        setAssigningRole(true);
        try {
            const response = await assignRolesAction(assigningRoleUser.id, {
                assignments: [{
                    business_id: isPlatformUser ? 0 : parseInt(selectedBusinessId),
                    role_id: parseInt(selectedRoleId)
                }]
            });

            if (response.success) {
                setShowAssignRoleModal(false);
                setAssigningRoleUser(null);
                loadUsers(); // Recargar usuarios
            } else {
                setError(response.message || 'Error al asignar rol');
            }
        } catch (err: any) {
            console.error('Error assigning role:', err);
            setError(err.message || 'Error al asignar rol');
        } finally {
            setAssigningRole(false);
        }
    };

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900">Usuarios</h1>
            </div>

            {error && <Alert type="error" onClose={() => setError(null)}>{error}</Alert>}

            {/* Filtros dinámicos y Tabla */}
            <div>
                <div className="bg-white rounded-t-lg shadow-sm border border-gray-200 border-b-0">
                    <div className="flex items-center justify-between p-4 sm:p-6 border-b border-gray-200 gap-4">
                        <div className="flex-1 min-w-0">
                            <DynamicFilters
                                availableFilters={availableFilters}
                                activeFilters={activeFilters}
                                onAddFilter={handleAddFilter}
                                onRemoveFilter={handleRemoveFilter}
                                className="!p-0 !border-0 !shadow-none"
                            />
                        </div>
                        <Button
                            variant="primary"
                            size="sm"
                            onClick={() => { setEditingUser(null); setShowCreateModal(true); }}
                            className="flex items-center justify-center flex-shrink-0"
                            title="Crear usuario"
                            aria-label="Crear usuario"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                        </Button>
                    </div>
                </div>
                {/* Tabla */}
                <div className="bg-white rounded-b-lg rounded-t-none shadow-sm border border-gray-200 border-t-0 overflow-hidden -mt-px">
                    <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    ID
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Avatar
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Nombre
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Email
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Teléfono
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Scope
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Rol
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Activo
                                </th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {loading ? (
                                <tr>
                                    <td colSpan={9} className="px-6 py-4 text-center text-sm text-gray-500">
                                        Cargando usuarios...
                                    </td>
                                </tr>
                            ) : users.length === 0 ? (
                                <tr>
                                    <td colSpan={9} className="px-6 py-4 text-center text-sm text-gray-500">
                                        No hay usuarios disponibles
                                    </td>
                                </tr>
                            ) : (
                                users.map((user) => (
                                    <tr key={user.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {user.id}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            {user.avatar_url ? (
                                                <img
                                                    src={user.avatar_url}
                                                    alt={user.name}
                                                    className="w-10 h-10 rounded-full object-cover"
                                                />
                                            ) : (
                                                <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center">
                                                    <span className="text-sm font-medium text-gray-600">
                                                        {user.name.charAt(0).toUpperCase()}
                                                    </span>
                                                </div>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {user.name}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {user.email}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {user.phone || '-'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            {user.scope_code ? (
                                                <span
                                                    className={`px-2 py-1 text-xs font-medium rounded-full ${
                                                        user.scope_code === 'platform'
                                                            ? 'bg-purple-100 text-purple-800'
                                                            : 'bg-blue-100 text-blue-800'
                                                    }`}
                                                >
                                                    {user.scope_code === 'platform' ? 'Platform' : 'Business'}
                                                </span>
                                            ) : (
                                                <span className="px-2 py-1 text-xs font-medium rounded-full bg-gray-100 text-gray-600">
                                                    Sin scope
                                                </span>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {user.business_role_assignments && user.business_role_assignments.length > 0 ? (
                                                <div className="flex flex-wrap gap-1">
                                                    {user.business_role_assignments.map((assignment, idx) => (
                                                        <span
                                                            key={idx}
                                                            className="px-2 py-0.5 text-xs bg-gray-100 text-gray-700 rounded"
                                                            title={`Negocio: ${assignment.business_name || assignment.business_id}`}
                                                        >
                                                            {assignment.role_name || `Rol ${assignment.role_id}`}
                                                        </span>
                                                    ))}
                                                </div>
                                            ) : (
                                                <span className="text-gray-400">-</span>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span
                                                className={`px-2 py-1 text-xs font-medium rounded-full ${
                                                    user.is_active
                                                        ? 'bg-green-100 text-green-800'
                                                        : 'bg-red-100 text-red-800'
                                                }`}
                                            >
                                                {user.is_active ? 'Sí' : 'No'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex justify-end gap-2">
                                                <button
                                                    onClick={() => handleOpenAssignRole(user)}
                                                    className="p-2 bg-indigo-500 hover:bg-indigo-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                                                    title="Asignar rol"
                                                    aria-label="Asignar rol"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                                                    </svg>
                                                </button>
                                                <button
                                                    onClick={() => handleEdit(user)}
                                                    className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                                                    title="Editar usuario"
                                                    aria-label="Editar usuario"
                                                    disabled={loadingUser}
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                    </svg>
                                                </button>
                                                <button
                                                    onClick={() => setDeleteId(user.id)}
                                                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                                                    title="Eliminar usuario"
                                                    aria-label="Eliminar usuario"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                                    </svg>
                                                </button>
                                            </div>
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {/* Paginación */}
                {!loading && users.length > 0 && (
                    <div className="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
                        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
                            {/* Desktop: Full pagination */}
                            <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                                <div>
                                    <p className="text-sm text-gray-700">
                                        Mostrando{' '}
                                        <span className="font-medium">
                                            {(page - 1) * pageSize + 1}
                                        </span>{' '}
                                        a{' '}
                                        <span className="font-medium">
                                            {Math.min(page * pageSize, total)}
                                        </span>{' '}
                                        de <span className="font-medium">{total}</span> resultados
                                    </p>
                                </div>
                                <nav className="flex items-center gap-2">
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page - 1 })}
                                        disabled={page === 1}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-l-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Anterior
                                    </button>
                                    <span className="relative inline-flex items-center px-3 sm:px-4 py-2 border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-700">
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => setFilters({ ...filters, page: page + 1 })}
                                        disabled={page === totalPages}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-r-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Siguiente
                                    </button>
                                </nav>
                            </div>

                            {/* Mobile: Page size selector */}
                            <div className="flex items-center justify-between w-full sm:hidden pt-2 border-t border-gray-200">
                                <div className="flex items-center gap-2">
                                    <label className="text-xs text-gray-700 whitespace-nowrap">
                                        Mostrar:
                                    </label>
                                    <select
                                        value={pageSize}
                                        onChange={(e) => {
                                            const newPageSize = parseInt(e.target.value);
                                            setFilters({ ...filters, page_size: newPageSize, page: 1 });
                                        }}
                                        className="px-2 py-1.5 text-xs border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900 bg-white"
                                    >
                                        <option value="10">10</option>
                                        <option value="20">20</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                    </select>
                                </div>
                                <p className="text-xs text-gray-500">
                                    Página {page} de {totalPages}
                                </p>
                            </div>
                        </div>
                    </div>
                )}
                </div>
            </div>

            <Modal
                isOpen={showCreateModal}
                onClose={() => { setShowCreateModal(false); setEditingUser(null); }}
                title={editingUser ? "Editar Usuario" : "Crear Usuario"}
                size="sm"
            >
                <UserForm
                    initialData={editingUser || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingUser(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Eliminar Usuario"
                message="¿Estás seguro de que deseas eliminar este usuario? Esta acción no se puede deshacer."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />

            {/* Modal para asignar rol */}
            <Modal
                isOpen={showAssignRoleModal}
                onClose={() => { setShowAssignRoleModal(false); setAssigningRoleUser(null); }}
                title={`Asignar Rol a ${assigningRoleUser?.name || 'Usuario'}`}
                size="sm"
            >
                {loadingRoles ? (
                    <div className="flex items-center justify-center py-8">
                        <Spinner size="lg" />
                    </div>
                ) : (
                    <div className="space-y-4">
                        {/* Indicador de tipo de usuario */}
                        {assigningRoleUser?.scope_code === 'platform' && (
                            <div className="p-3 bg-purple-50 border border-purple-200 rounded-lg">
                                <p className="text-sm text-purple-700">
                                    <span className="font-medium">Usuario Platform:</span> Los roles se asignan a nivel de plataforma (sin negocio específico).
                                </p>
                            </div>
                        )}

                        {/* Roles actuales del usuario */}
                        {assigningRoleUser?.business_role_assignments && assigningRoleUser.business_role_assignments.length > 0 && (
                            <div className="p-3 bg-gray-50 rounded-lg">
                                <p className="text-sm font-medium text-gray-700 mb-2">Roles actuales:</p>
                                <div className="flex flex-wrap gap-2">
                                    {assigningRoleUser.business_role_assignments.map((assignment, idx) => (
                                        <span
                                            key={idx}
                                            className="px-2 py-1 text-xs bg-blue-100 text-blue-800 rounded-full"
                                        >
                                            {assignment.role_name}{assignment.business_name ? ` - ${assignment.business_name}` : ''}
                                        </span>
                                    ))}
                                </div>
                            </div>
                        )}

                        {/* Selector de negocio - Solo para usuarios NO platform */}
                        {assigningRoleUser?.scope_code !== 'platform' && (
                            <Select
                                label="Negocio *"
                                value={selectedBusinessId}
                                onChange={(e) => setSelectedBusinessId(e.target.value)}
                                options={[
                                    { value: '', label: 'Seleccionar negocio...' },
                                    ...businesses.map(b => ({ value: String(b.id), label: b.name }))
                                ]}
                            />
                        )}

                        <Select
                            label="Rol *"
                            value={selectedRoleId}
                            onChange={(e) => setSelectedRoleId(e.target.value)}
                            options={[
                                { value: '', label: 'Seleccionar rol...' },
                                ...roles
                                    .filter(r => assigningRoleUser?.scope_code === 'platform' 
                                        ? r.scope_code === 'platform' 
                                        : r.scope_code === 'business')
                                    .map(r => ({ value: String(r.id), label: r.name }))
                            ]}
                        />

                        {assigningRoleUser?.scope_code !== 'platform' && (
                            <p className="text-xs text-gray-500">
                                El usuario debe estar previamente asociado al negocio para poder asignarle un rol.
                            </p>
                        )}

                        <div className="flex justify-end gap-2 pt-4">
                            <Button
                                type="button"
                                variant="secondary"
                                onClick={() => { setShowAssignRoleModal(false); setAssigningRoleUser(null); }}
                            >
                                Cancelar
                            </Button>
                            <Button
                                type="button"
                                onClick={handleAssignRole}
                                disabled={assigningRole || !selectedRoleId || (assigningRoleUser?.scope_code !== 'platform' && !selectedBusinessId)}
                            >
                                {assigningRole ? <Spinner size="sm" /> : 'Asignar Rol'}
                            </Button>
                        </div>
                    </div>
                )}
            </Modal>
        </div>
    );
};
