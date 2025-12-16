'use client';

import React, { useState, useCallback, useMemo } from 'react';
import { Button } from '@/shared/ui/button';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { Resource } from '../../domain/types';
import { ResourceForm } from './ResourceForm';
import { useResources } from '../hooks/useResources';

export const ResourceList: React.FC = () => {
    const {
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
        refresh,
        setError
    } = useResources();

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingResource, setEditingResource] = useState<Resource | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = [
        {
            key: 'name',
            label: 'Nombre',
            type: 'text',
            placeholder: 'Buscar por nombre...',
        },
        {
            key: 'description',
            label: 'Descripción',
            type: 'text',
            placeholder: 'Buscar por descripción...',
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

        if (filters.description) {
            active.push({
                key: 'description',
                label: 'Descripción',
                value: filters.description,
                type: 'text',
            });
        }

        return active;
    }, [filters]);

    // Manejar adición de filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev };
            (newFilters as any)[filterKey] = value;
            return newFilters;
        });
        setPage(1);
    }, [setFilters, setPage]);

    // Manejar eliminación de filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev };
            delete (newFilters as any)[filterKey];
            return newFilters;
        });
        setPage(1);
    }, [setFilters, setPage]);

    const handleDelete = async () => {
        if (deleteId) {
            const success = await deleteResource(deleteId);
            if (success) setDeleteId(null);
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingResource(null);
        refresh();
    };

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900">Recursos</h1>
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
                            onClick={() => { setEditingResource(null); setShowCreateModal(true); }}
                            className="flex items-center justify-center flex-shrink-0"
                            title="Crear recurso"
                            aria-label="Crear recurso"
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
                                        Descripción
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
                                        <td colSpan={5} className="px-6 py-4 text-center text-sm text-gray-500">
                                            Cargando recursos...
                                        </td>
                                    </tr>
                                ) : resources.length === 0 ? (
                                    <tr>
                                        <td colSpan={5} className="px-6 py-4 text-center text-sm text-gray-500">
                                            No hay recursos disponibles
                                        </td>
                                    </tr>
                                ) : (
                                    resources.map((resource) => (
                                        <tr key={resource.id} className="hover:bg-gray-50">
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                                {resource.id}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                                {resource.name}
                                            </td>
                                            <td className="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">
                                                {resource.description || '-'}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                                {resource.business_type_name || '-'}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                                <div className="flex justify-end gap-2">
                                                    <button
                                                        onClick={() => { setEditingResource(resource); setShowCreateModal(true); }}
                                                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                                                        title="Editar recurso"
                                                        aria-label="Editar recurso"
                                                    >
                                                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                        </svg>
                                                    </button>
                                                    <button
                                                        onClick={() => setDeleteId(resource.id)}
                                                        className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                                                        title="Eliminar recurso"
                                                        aria-label="Eliminar recurso"
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
                    {!loading && resources.length > 0 && (
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
                onClose={() => { setShowCreateModal(false); setEditingResource(null); }}
                title={editingResource ? "Editar Recurso" : "Crear Recurso"}
                size="sm"
            >
                <ResourceForm
                    initialData={editingResource || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingResource(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Eliminar Recurso"
                message="¿Estás seguro de que deseas eliminar este recurso? Esta acción no se puede deshacer."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />
        </div>
    );
};
