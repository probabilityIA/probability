'use client';

import React, { useState, useEffect, useCallback, useMemo } from 'react';
import { Button } from '@/shared/ui/button';
import { Alert } from '@/shared/ui/alert';
import { Modal } from '@/shared/ui/modal';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui';
import { Business, GetBusinessesParams, ConfiguredResource, BusinessConfiguredResources } from '../../domain/types';
import { BusinessForm } from './BusinessForm';
import { 
    getBusinessesAction, 
    deleteBusinessAction, 
    getBusinessTypesAction,
    getBusinessConfiguredResourcesAction,
    activateResourceAction,
    deactivateResourceAction
} from '../../infra/actions';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { BusinessType } from '../../domain/types';
import { Spinner } from '@/shared/ui/spinner';

export const BusinessList: React.FC = () => {
    const [businesses, setBusinesses] = useState<Business[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const [pageSize, setPageSize] = useState(20);
    const [businessTypes, setBusinessTypes] = useState<BusinessType[]>([]);

    const [showCreateModal, setShowCreateModal] = useState(false);
    const [editingBusiness, setEditingBusiness] = useState<Business | null>(null);
    const [deleteId, setDeleteId] = useState<number | null>(null);

    // Estado para modal de recursos
    const [showResourcesModal, setShowResourcesModal] = useState(false);
    const [selectedBusiness, setSelectedBusiness] = useState<Business | null>(null);
    const [businessResources, setBusinessResources] = useState<ConfiguredResource[]>([]);
    const [loadingResources, setLoadingResources] = useState(false);
    const [savingResource, setSavingResource] = useState<number | null>(null);

    // Filters
    const [filters, setFilters] = useState<GetBusinessesParams>({
        page: 1,
        per_page: 20,
    });

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = useMemo(() => [
        {
            key: 'name',
            label: 'Nombre',
            type: 'text',
            placeholder: 'Buscar por nombre...',
        },
        {
            key: 'business_type_id',
            label: 'Tipo',
            type: 'select',
            options: [
                { value: '', label: 'Todos los tipos' },
                ...businessTypes.map(t => ({ value: String(t.id), label: t.name }))
            ],
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
    ], [businessTypes]);

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

        if (filters.business_type_id) {
            const type = businessTypes.find(t => t.id === filters.business_type_id);
            active.push({
                key: 'business_type_id',
                label: 'Tipo',
                value: String(filters.business_type_id),
                type: 'select',
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

        return active;
    }, [filters, businessTypes]);

    // Manejar adición de filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev, page: 1 };

            if (filterKey === 'is_active') {
                newFilters.is_active = value === 'true' || value === true;
            } else if (filterKey === 'business_type_id') {
                newFilters.business_type_id = value ? Number(value) : undefined;
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

    // Cargar tipos de negocio
    const loadBusinessTypes = useCallback(async () => {
        try {
            const res = await getBusinessTypesAction();
            setBusinessTypes(res.data);
        } catch (e) {
            console.error("Failed to load business types", e);
        }
    }, []);

    // Cargar negocios
    const loadBusinesses = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const response = await getBusinessesAction(filters);
            setBusinesses(response.data || []);
            if (response.pagination) {
                setPage(response.pagination.current_page);
                setTotalPages(response.pagination.last_page);
                setTotal(response.pagination.total);
                setPageSize(response.pagination.per_page);
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar negocios');
        } finally {
            setLoading(false);
        }
    }, [filters]);

    useEffect(() => {
        loadBusinessTypes();
    }, [loadBusinessTypes]);

    useEffect(() => {
        loadBusinesses();
    }, [loadBusinesses]);

    const handleDelete = async () => {
        if (deleteId) {
            try {
                const response = await deleteBusinessAction(deleteId);
                if (response.success) {
                    setDeleteId(null);
                    loadBusinesses();
                } else {
                    setError(response.message || 'Error al eliminar negocio');
                }
            } catch (err: any) {
                setError(err.message || 'Error al eliminar negocio');
            }
        }
    };

    const handleSave = () => {
        setShowCreateModal(false);
        setEditingBusiness(null);
        loadBusinesses();
    };

    // Abrir modal de recursos
    const handleOpenResources = async (business: Business) => {
        setSelectedBusiness(business);
        setShowResourcesModal(true);
        setLoadingResources(true);
        try {
            const response = await getBusinessConfiguredResourcesAction(business.id);
            if (response.data) {
                setBusinessResources(response.data.resources || []);
            }
        } catch (err: any) {
            setError(err.message || 'Error al cargar recursos');
            setBusinessResources([]);
        } finally {
            setLoadingResources(false);
        }
    };

    // Toggle recurso activo/inactivo
    const handleToggleResource = async (resource: ConfiguredResource) => {
        if (!selectedBusiness) return;
        
        setSavingResource(resource.resource_id);
        try {
            if (resource.is_active) {
                await deactivateResourceAction(resource.resource_id, selectedBusiness.id);
            } else {
                await activateResourceAction(resource.resource_id, selectedBusiness.id);
            }
            
            // Actualizar estado local
            setBusinessResources(prev => 
                prev.map(r => 
                    r.resource_id === resource.resource_id 
                        ? { ...r, is_active: !r.is_active }
                        : r
                )
            );
        } catch (err: any) {
            setError(err.message || 'Error al actualizar recurso');
        } finally {
            setSavingResource(null);
        }
    };

    // Cerrar modal de recursos
    const handleCloseResourcesModal = () => {
        setShowResourcesModal(false);
        setSelectedBusiness(null);
        setBusinessResources([]);
    };

    return (
        <div className="p-6 space-y-6">
            <div className="flex justify-between items-center">
                <h1 className="text-2xl font-bold text-gray-900">Negocios</h1>
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
                            onClick={() => { setEditingBusiness(null); setShowCreateModal(true); }}
                            className="flex items-center justify-center flex-shrink-0"
                            title="Crear negocio"
                            aria-label="Crear negocio"
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
                                    Logo
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Nombre
                                </th>
                                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                                    Tipo
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
                                    <td colSpan={6} className="px-6 py-4 text-center text-sm text-gray-500">
                                        Cargando negocios...
                                    </td>
                                </tr>
                            ) : businesses.length === 0 ? (
                                <tr>
                                    <td colSpan={6} className="px-6 py-4 text-center text-sm text-gray-500">
                                        No hay negocios disponibles
                                    </td>
                                </tr>
                            ) : (
                                businesses.map((business) => (
                                    <tr key={business.id} className="hover:bg-gray-50">
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                                            {business.id}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            {business.logo_url ? (
                                                <img
                                                    src={business.logo_url}
                                                    alt={business.name}
                                                    className="w-10 h-10 rounded-full object-cover"
                                                />
                                            ) : (
                                                <div className="w-10 h-10 rounded-full bg-gray-200 flex items-center justify-center">
                                                    <span className="text-sm font-medium text-gray-600">
                                                        {business.name.charAt(0).toUpperCase()}
                                                    </span>
                                                </div>
                                            )}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                                            {business.name}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                                            {business.business_type?.name || business.business_type_id || '-'}
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap">
                                            <span
                                                className={`px-2 py-1 text-xs font-medium rounded-full ${
                                                    business.is_active
                                                        ? 'bg-green-100 text-green-800'
                                                        : 'bg-red-100 text-red-800'
                                                }`}
                                            >
                                                {business.is_active ? 'Sí' : 'No'}
                                            </span>
                                        </td>
                                        <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                                            <div className="flex justify-end gap-2">
                                                <button
                                                    onClick={() => handleOpenResources(business)}
                                                    className="p-2 bg-purple-500 hover:bg-purple-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-purple-500 focus:ring-offset-2"
                                                    title="Configurar recursos"
                                                    aria-label="Configurar recursos"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                                                    </svg>
                                                </button>
                                                <button
                                                    onClick={() => { setEditingBusiness(business); setShowCreateModal(true); }}
                                                    className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                                                    title="Editar negocio"
                                                    aria-label="Editar negocio"
                                                >
                                                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                                                    </svg>
                                                </button>
                                                <button
                                                    onClick={() => setDeleteId(business.id)}
                                                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                                                    title="Eliminar negocio"
                                                    aria-label="Eliminar negocio"
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
                {!loading && businesses.length > 0 && (
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
                                            setFilters({ ...filters, per_page: newPageSize, page: 1 });
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
                onClose={() => { setShowCreateModal(false); setEditingBusiness(null); }}
                title={editingBusiness ? "Editar Negocio" : "Crear Negocio"}
                size="4xl"
            >
                <BusinessForm
                    initialData={editingBusiness || undefined}
                    onSuccess={handleSave}
                    onCancel={() => { setShowCreateModal(false); setEditingBusiness(null); }}
                />
            </Modal>

            <ConfirmModal
                isOpen={!!deleteId}
                title="Eliminar Negocio"
                message="¿Estás seguro de que deseas eliminar este negocio? Esta acción no se puede deshacer."
                onConfirm={handleDelete}
                onClose={() => setDeleteId(null)}
            />

            {/* Modal de configuración de recursos */}
            <Modal
                isOpen={showResourcesModal}
                onClose={handleCloseResourcesModal}
                title={`Configurar Recursos - ${selectedBusiness?.name || ''}`}
                size="lg"
            >
                <div className="space-y-4">
                    {loadingResources ? (
                        <div className="flex justify-center py-8">
                            <Spinner />
                        </div>
                    ) : businessResources.length === 0 ? (
                        <div className="text-center py-8 text-gray-500">
                            <svg className="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                            </svg>
                            <p className="mt-2">No hay recursos configurados para este negocio.</p>
                            <p className="text-sm">Contacta al administrador para asignar recursos.</p>
                        </div>
                    ) : (
                        <>
                            <p className="text-sm text-gray-600 mb-4">
                                Activa o desactiva los recursos disponibles para este negocio. Los usuarios solo podrán acceder a los recursos activos.
                            </p>
                            <div className="divide-y divide-gray-200 border rounded-lg overflow-hidden">
                                {businessResources.map((resource) => (
                                    <div 
                                        key={resource.resource_id} 
                                        className="flex items-center justify-between p-4 hover:bg-gray-50 transition-colors"
                                    >
                                        <div className="flex items-center gap-3">
                                            <div className={`w-2 h-2 rounded-full ${resource.is_active ? 'bg-green-500' : 'bg-gray-300'}`} />
                                            <span className="font-medium text-gray-900">{resource.resource_name}</span>
                                        </div>
                                        <button
                                            onClick={() => handleToggleResource(resource)}
                                            disabled={savingResource === resource.resource_id}
                                            className={`
                                                relative inline-flex h-6 w-11 flex-shrink-0 cursor-pointer rounded-full border-2 border-transparent 
                                                transition-colors duration-200 ease-in-out focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
                                                ${resource.is_active ? 'bg-blue-600' : 'bg-gray-200'}
                                                ${savingResource === resource.resource_id ? 'opacity-50 cursor-wait' : ''}
                                            `}
                                        >
                                            {savingResource === resource.resource_id ? (
                                                <span className="absolute inset-0 flex items-center justify-center">
                                                    <svg className="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                                    </svg>
                                                </span>
                                            ) : (
                                                <span
                                                    className={`
                                                        pointer-events-none inline-block h-5 w-5 transform rounded-full bg-white shadow ring-0 
                                                        transition duration-200 ease-in-out
                                                        ${resource.is_active ? 'translate-x-5' : 'translate-x-0'}
                                                    `}
                                                />
                                            )}
                                        </button>
                                    </div>
                                ))}
                            </div>
                            <div className="flex items-center justify-between pt-4 text-sm text-gray-500">
                                <span>
                                    {businessResources.filter(r => r.is_active).length} de {businessResources.length} recursos activos
                                </span>
                            </div>
                        </>
                    )}
                    <div className="flex justify-end pt-4 border-t">
                        <Button variant="secondary" onClick={handleCloseResourcesModal}>
                            Cerrar
                        </Button>
                    </div>
                </div>
            </Modal>
        </div>
    );
};
