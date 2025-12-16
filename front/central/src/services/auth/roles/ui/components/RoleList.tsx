'use client';

import React, { useState, useCallback, useMemo } from 'react';
import { Button } from '@/shared/ui/button';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { Spinner } from '@/shared/ui/spinner';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { Role, GetRolesParams } from '../../domain/types';
import { Permission } from '@/services/auth/permissions/domain/types';
import { RoleForm } from './RoleForm';
import { useRoles } from '../hooks/useRoles';
import { getPermissionsAction } from '@/services/auth/permissions/infra/actions';
import { getRolePermissionsAction, assignPermissionsAction } from '../../infra/actions';

export const RoleList: React.FC = () => {
    const {
        roles,
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
        refresh,
        setError
    } = useRoles();

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingRole, setEditingRole] = useState<Role | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    // Estados para asignar permisos
    const [showPermissionsModal, setShowPermissionsModal] = useState(false);
    const [selectedRole, setSelectedRole] = useState<Role | null>(null);
    const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
    const [rolePermissionIds, setRolePermissionIds] = useState<number[]>([]);
    const [loadingPermissions, setLoadingPermissions] = useState(false);
    const [savingPermissions, setSavingPermissions] = useState(false);

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = [
        {
            key: 'name',
            label: 'Nombre',
            type: 'text',
            placeholder: 'Buscar por nombre...',
        },
        {
            key: 'scope_id',
            label: 'ID de Scope',
            type: 'text',
            placeholder: 'Filtrar por ID de scope...',
        },
        {
            key: 'business_type_id',
            label: 'ID de Tipo de Negocio',
            type: 'text',
            placeholder: 'Filtrar por ID de tipo de negocio...',
        },
        {
            key: 'level',
            label: 'Nivel',
            type: 'text',
            placeholder: 'Filtrar por nivel...',
        },
        {
            key: 'is_system',
            label: 'Tipo',
            type: 'select',
            options: [
                { value: 'true', label: 'Sistema' },
                { value: 'false', label: 'No Sistema' },
            ],
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

        if (filters.scope_id) {
            active.push({
                key: 'scope_id',
                label: 'ID de Scope',
                value: String(filters.scope_id),
                type: 'text',
            });
        }

        if (filters.business_type_id) {
            active.push({
                key: 'business_type_id',
                label: 'ID de Tipo de Negocio',
                value: String(filters.business_type_id),
                type: 'text',
            });
        }

        if (filters.level) {
            active.push({
                key: 'level',
                label: 'Nivel',
                value: String(filters.level),
                type: 'text',
            });
        }

        if (filters.is_system !== undefined) {
            active.push({
                key: 'is_system',
                label: 'Tipo',
                value: filters.is_system ? 'Sistema' : 'No Sistema',
                type: 'select',
            });
        }

        return active;
    }, [filters]);

    // Manejar adición de filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };

            if (filterKey === 'is_system') {
                newFilters.is_system = value === 'true' || value === true;
            } else if (filterKey === 'scope_id') {
                newFilters.scope_id = value ? Number(value) : undefined;
            } else if (filterKey === 'business_type_id') {
                newFilters.business_type_id = value ? Number(value) : undefined;
            } else if (filterKey === 'level') {
                newFilters.level = value ? Number(value) : undefined;
            } else {
                (newFilters as any)[filterKey] = value;
            }

            return newFilters;
        });
        setPage(1);
    }, [setFilters, setPage]);

    // Manejar eliminación de filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };
            delete (newFilters as any)[filterKey];
            return newFilters;
        });
        setPage(1);
    }, [setFilters, setPage]);

    const handleDelete = async () => {
        if (deleteId) {
            const success = await deleteRole(deleteId);
            if (success) setDeleteId(null);
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingRole(null);
        refresh();
    };

    // Abrir modal de asignar permisos
    const handleOpenPermissions = async (role: Role) => {
        setSelectedRole(role);
        setShowPermissionsModal(true);
        setLoadingPermissions(true);

        try {
            // Cargar todos los permisos y los permisos actuales del rol en paralelo
            const [permissionsRes, rolePermissionsRes] = await Promise.all([
                getPermissionsAction({ scope_id: role.scope_id }),
                getRolePermissionsAction(role.id)
            ]);

            if (permissionsRes.success && permissionsRes.data) {
                setAllPermissions(permissionsRes.data);
            }

            if (rolePermissionsRes.success && rolePermissionsRes.permissions) {
                setRolePermissionIds(rolePermissionsRes.permissions.map((p: any) => p.id));
            } else {
                setRolePermissionIds([]);
            }
        } catch (err: any) {
            console.error('Error loading permissions:', err);
            setError('Error al cargar permisos');
        } finally {
            setLoadingPermissions(false);
        }
    };

    // Toggle permiso
    const handleTogglePermission = (permissionId: number) => {
        setRolePermissionIds(prev => {
            if (prev.includes(permissionId)) {
                return prev.filter(id => id !== permissionId);
            } else {
                return [...prev, permissionId];
            }
        });
    };

    // Guardar permisos
    const handleSavePermissions = async () => {
        if (!selectedRole) return;

        setSavingPermissions(true);
        try {
            const response = await assignPermissionsAction(selectedRole.id, {
                permission_ids: rolePermissionIds
            });

            if (response.success) {
                setShowPermissionsModal(false);
                setSelectedRole(null);
            } else {
                setError(response.message || 'Error al asignar permisos');
            }
        } catch (err: any) {
            console.error('Error saving permissions:', err);
            setError('Error al guardar permisos');
        } finally {
            setSavingPermissions(false);
        }
    };

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900">Roles</h1>
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
                            onClick={() => { setEditingRole(null); setShowCreateModal(true); }}
                            className="flex items-center justify-center flex-shrink-0"
                            title="Crear rol"
                            aria-label="Crear rol"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                        </Button>
                    </div>
                </div>
                {/* Tabla */}
                <div className="bg-white rounded-b-lg rounded-t-none shadow-sm border border-gray-200 border-t-0 overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    ID
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Nombre
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Nivel
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Sistema
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Scope
                                </th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {loading ? (
                                <tr>
                                    <td colSpan={6} className="px-6 py-4 text-center text-sm text-gray-500">
                                        Cargando roles...
                                    </td>
                                </tr>
                            ) : roles.length === 0 ? (
                                <tr>
                                    <td colSpan={6} className="px-6 py-4 text-center text-sm text-gray-500">
                                        No hay roles disponibles
                                    </td>
                                </tr>
                            ) : (
                                roles.map((role) => (
                                    <tr key={role.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {role.id}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {role.name}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {role.level}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span
                                                className={`px-2 py-1 text-xs font-medium rounded-full ${
                                                    role.is_system
                                                        ? 'bg-blue-100 text-blue-800'
                                                        : 'bg-gray-100 text-gray-800'
                                                }`}
                                            >
                                                {role.is_system ? 'Sí' : 'No'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {role.scope_name || '-'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex justify-end gap-2">
                                                {/* Botón Asignar Permisos */}
                                                <button
                                                    onClick={() => handleOpenPermissions(role)}
                                                    className="p-2 bg-purple-500 hover:bg-purple-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
                                                    title="Asignar permisos"
                                                    aria-label="Asignar permisos"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
                                                    </svg>
                                                </button>
                                                {/* Botón Editar */}
                                                <button
                                                    onClick={() => { setEditingRole(role); setShowCreateModal(true); }}
                                                    className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                                                    title="Editar rol"
                                                    aria-label="Editar rol"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                    </svg>
                                                </button>
                                                {/* Botón Eliminar */}
                                                <button
                                                    onClick={() => setDeleteId(role.id)}
                                                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                                                    title="Eliminar rol"
                                                    aria-label="Eliminar rol"
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
                {!loading && roles.length > 0 && (
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
                                        onClick={() => setPage(page - 1)}
                                        disabled={page === 1}
                                        className="relative inline-flex items-center px-2 sm:px-3 py-2 rounded-l-md border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                                    >
                                        Anterior
                                    </button>
                                    <span className="relative inline-flex items-center px-3 sm:px-4 py-2 border border-gray-300 bg-white text-xs sm:text-sm font-medium text-gray-700">
                                        Página {page} de {totalPages}
                                    </span>
                                    <button
                                        onClick={() => setPage(page + 1)}
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
                                            setPageSize(newPageSize);
                                            setPage(1);
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

            {/* Modal Crear/Editar Rol */}
            <Modal
                isOpen={showCreateModal}
                onClose={() => { setShowCreateModal(false); setEditingRole(null); }}
                title={editingRole ? "Editar Rol" : "Crear Rol"}
                size="sm"
            >
                <RoleForm
                    initialData={editingRole || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingRole(null); }}
                />
            </Modal>

            {/* Modal Asignar Permisos */}
            <Modal
                isOpen={showPermissionsModal}
                onClose={() => { setShowPermissionsModal(false); setSelectedRole(null); }}
                title={`Asignar Permisos - ${selectedRole?.name || ''}`}
                size="lg"
            >
                {loadingPermissions ? (
                    <div className="flex justify-center items-center py-8">
                        <Spinner size="lg" />
                        <span className="ml-2 text-gray-600">Cargando permisos...</span>
                    </div>
                ) : (
                    <div className="space-y-4">
                        {/* Info del rol */}
                        <div className="p-3 bg-gray-50 rounded-lg border border-gray-200">
                            <p className="text-sm text-gray-600">
                                <span className="font-medium">Scope:</span> {selectedRole?.scope_name || '-'}
                            </p>
                            <p className="text-sm text-gray-500 mt-1">
                                Solo se muestran permisos del mismo scope que el rol.
                            </p>
                        </div>

                        {/* Lista de permisos */}
                        {allPermissions.length === 0 ? (
                            <div className="text-center py-8 text-gray-500">
                                No hay permisos disponibles para este scope.
                            </div>
                        ) : (
                            <div className="max-h-96 overflow-y-auto border border-gray-200 rounded-lg">
                                <table className="min-w-full divide-y divide-gray-200">
                                    <thead className="bg-gray-50 sticky top-0">
                                        <tr>
                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">
                                                Seleccionar
                                            </th>
                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">
                                                Permiso
                                            </th>
                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">
                                                Recurso
                                            </th>
                                            <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">
                                                Acción
                                            </th>
                                        </tr>
                                    </thead>
                                    <tbody className="bg-white divide-y divide-gray-200">
                                        {allPermissions.map((permission) => (
                                            <tr 
                                                key={permission.id} 
                                                className={`hover:bg-gray-50 cursor-pointer ${
                                                    rolePermissionIds.includes(permission.id) ? 'bg-purple-50' : ''
                                                }`}
                                                onClick={() => handleTogglePermission(permission.id)}
                                            >
                                                <td className="px-4 py-2">
                                                    <input
                                                        type="checkbox"
                                                        checked={rolePermissionIds.includes(permission.id)}
                                                        onChange={() => handleTogglePermission(permission.id)}
                                                        className="h-4 w-4 text-purple-600 focus:ring-purple-500 border-gray-300 rounded"
                                                        onClick={(e) => e.stopPropagation()}
                                                    />
                                                </td>
                                                <td className="px-4 py-2 text-sm font-medium text-gray-900">
                                                    {permission.name}
                                                </td>
                                                <td className="px-4 py-2 text-sm text-gray-500">
                                                    {permission.resource || '-'}
                                                </td>
                                                <td className="px-4 py-2 text-sm text-gray-500">
                                                    {permission.action || '-'}
                                                </td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        )}

                        {/* Contador de seleccionados */}
                        <div className="flex justify-between items-center text-sm text-gray-600">
                            <span>{rolePermissionIds.length} permiso(s) seleccionado(s)</span>
                        </div>

                        {/* Botones */}
                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200">
                            <Button
                                variant="secondary"
                                onClick={() => { setShowPermissionsModal(false); setSelectedRole(null); }}
                            >
                                Cancelar
                            </Button>
                            <Button
                                variant="primary"
                                onClick={handleSavePermissions}
                                disabled={savingPermissions}
                            >
                                {savingPermissions ? <Spinner size="sm" /> : 'Guardar Permisos'}
                            </Button>
                        </div>
                    </div>
                )}
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Eliminar Rol"
                message="¿Estás seguro de que deseas eliminar este rol? Esta acción no se puede deshacer."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
