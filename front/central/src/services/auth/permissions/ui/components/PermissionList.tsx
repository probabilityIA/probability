'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Button } from '@/shared/ui/button';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { Permission } from '../../domain/types';
import { PermissionForm } from './PermissionForm';
import { getPermissionsAction, deletePermissionAction } from '../../infra/actions';
import { ConfirmModal } from '@/shared/ui/confirm-modal';

export const PermissionList: React.FC = () => {
    const [permissions, setPermissions] = useState<Permission[]>([]);
    const [allPermissions, setAllPermissions] = useState<Permission[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const [pageSize, setPageSize] = useState(20);

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingPermission, setEditingPermission] = useState<Permission | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    // Filters
    const [filters, setFilters] = useState<{
        name?: string;
        scope_id?: number;
        business_type_id?: number;
    }>({});

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

        return active;
    }, [filters]);

    // Manejar adición de filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev };

            if (filterKey === 'scope_id') {
                newFilters.scope_id = value ? Number(value) : undefined;
            } else if (filterKey === 'business_type_id') {
                newFilters.business_type_id = value ? Number(value) : undefined;
            } else {
                (newFilters as any)[filterKey] = value;
            }

            return newFilters;
        });
        setPage(1);
    }, []);

    // Manejar eliminación de filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };
            delete (newFilters as any)[filterKey];
            return newFilters;
        });
        setPage(1);
    }, []);

    // Cargar permisos
    const loadPermissions = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getPermissionsAction({
                name: filters.name,
                scope_id: filters.scope_id,
                business_type_id: filters.business_type_id,
            });
            if (response.success) {
                setAllPermissions(response.data || []);
                setTotal(response.data?.length || 0);
            } else {
                setError('Error al cargar permisos');
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar permisos');
        } finally {
            setLoading(false);
        }
    }, [filters]);

    useEffect(() => {
        loadPermissions();
    }, [loadPermissions]);

    // Paginación del lado del cliente
    useEffect(() => {
        const startIndex = (page - 1) * pageSize;
        const endIndex = startIndex + pageSize;
        const paginated = allPermissions.slice(startIndex, endIndex);
        setPermissions(paginated);
        setTotalPages(Math.ceil(allPermissions.length / pageSize));
    }, [allPermissions, page, pageSize]);

    const handleDelete = async () => {
        if (deleteId) {
            try {
                await deletePermissionAction(deleteId);
                setDeleteId(null);
                loadPermissions();
            } catch (err: any) {
                setError(err.message || 'Error al eliminar permiso');
            }
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingPermission(null);
        loadPermissions();
    };

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900">Permisos</h1>
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
                            onClick={() => { setEditingPermission(null); setShowCreateModal(true); }}
                            className="flex items-center justify-center flex-shrink-0"
                            title="Crear permiso"
                            aria-label="Crear permiso"
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
                                    Código
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Recurso
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acción
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Scope
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Tipo de Negocio
                                </th>
                                <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Acciones
                                </th>
                            </tr>
                        </thead>
                        <tbody className="bg-white divide-y divide-gray-200">
                            {loading ? (
                                <tr>
                                    <td colSpan={8} className="px-6 py-4 text-center text-sm text-gray-500">
                                        Cargando permisos...
                                    </td>
                                </tr>
                            ) : permissions.length === 0 ? (
                                <tr>
                                    <td colSpan={8} className="px-6 py-4 text-center text-sm text-gray-500">
                                        No hay permisos disponibles
                                    </td>
                                </tr>
                            ) : (
                                permissions.map((permission) => (
                                    <tr key={permission.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {permission.id}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {permission.name}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {permission.code}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {permission.resource}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {permission.action}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {permission.scope_name}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {permission.business_type_name || '-'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex justify-end gap-2">
                                                <button
                                                    onClick={() => { setEditingPermission(permission); setShowCreateModal(true); }}
                                                    className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                                                    title="Editar permiso"
                                                    aria-label="Editar permiso"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                    </svg>
                                                </button>
                                                <button
                                                    onClick={() => setDeleteId(permission.id)}
                                                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                                                    title="Eliminar permiso"
                                                    aria-label="Eliminar permiso"
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
                {!loading && permissions.length > 0 && (
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

            <Modal
                isOpen={showCreateModal}
                onClose={() => { setShowCreateModal(false); setEditingPermission(null); }}
                title={editingPermission ? "Editar Permiso" : "Crear Permiso"}
                size="sm"
            >
                <PermissionForm
                    initialData={editingPermission || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingPermission(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Eliminar Permiso"
                message="¿Estás seguro de que deseas eliminar este permiso? Esta acción no se puede deshacer."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
