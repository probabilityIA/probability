'use client';

import { useState } from 'react';
import { useIntegrations } from '../hooks/useIntegrations';
import { Integration } from '../../domain/types';
import { Input, Button, Badge, Spinner, Table, Alert, ConfirmModal } from '@/shared/ui';

interface IntegrationListProps {
    onEdit?: (integration: Integration) => void;
}

export default function IntegrationList({ onEdit }: IntegrationListProps) {
    const {
        integrations,
        loading,
        error,
        page,
        setPage,
        totalPages,
        search,
        setSearch,
        filterType,
        setFilterType,
        filterCategory,
        setFilterCategory,
        deleteIntegration,
        toggleActive,
        setAsDefault,
        testConnection,
        syncOrders,
        setError
    } = useIntegrations();

    const [deleteModal, setDeleteModal] = useState<{ show: boolean; id: number | null }>({
        show: false,
        id: null
    });

    const handleDeleteClick = (id: number) => {
        setDeleteModal({ show: true, id });
    };

    const handleDeleteConfirm = async () => {
        if (deleteModal.id) {
            const success = await deleteIntegration(deleteModal.id);
            if (success) {
                setDeleteModal({ show: false, id: null });
            }
        }
    };

    const handleTest = async (id: number) => {
        const result = await testConnection(id);
        if (result.success) {
            alert('✅ Conexión exitosa');
        } else {
            alert(`❌ Error: ${result.message}`);
        }
    };

    const handleSync = async (id: number) => {
        // Implement logic to show loading state if needed, or just toast
        // Since we don't have toast in this component (it seems), we use alert for now or try to get toast provider
        // Wait, OrderList had toast provider. Let's see imports.
        // It imports Alert from ui but not useToast.
        // I'll stick to alert for simplicity as requested, or just call the function.
        // The user said "better move that button".
        const result = await syncOrders(id);
        if (result.success) {
            alert('✅ Sincronización iniciada correctamente');
        } else {
            alert(`❌ Error al iniciar sincronización: ${result.message}`);
        }
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    const columns = [
        { key: 'id', label: 'ID' },
        { key: 'name', label: 'Nombre' },
        { key: 'type', label: 'Tipo' },
        { key: 'category', label: 'Categoría' },
        { key: 'status', label: 'Estado' },
        { key: 'actions', label: 'Acciones' }
    ];

    const renderRow = (integration: Integration) => ({
        id: integration.id,
        name: (
            <div>
                <div className="text-sm font-medium text-gray-900">{integration.name}</div>
                {integration.description && (
                    <div className="text-sm text-gray-500">{integration.description}</div>
                )}
            </div>
        ),
        type: integration.type,
        category: integration.category,
        status: (
            <div className="flex items-center gap-2">
                <Badge type={integration.is_active ? 'success' : 'error'}>
                    {integration.is_active ? 'Activo' : 'Inactivo'}
                </Badge>
                {integration.is_default && (
                    <Badge type="primary">Por defecto</Badge>
                )}
            </div>
        ),
        actions: (
            <div className="flex gap-2">
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleTest(integration.id)}
                >
                    Probar
                </Button>
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleSync(integration.id)}
                >
                    ↻ Sincronizar
                </Button>
                {onEdit && (
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => onEdit(integration)}
                    >
                        Editar
                    </Button>
                )}
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => toggleActive(integration.id, integration.is_active)}
                >
                    {integration.is_active ? 'Desactivar' : 'Activar'}
                </Button>
                {!integration.is_default && (
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => setAsDefault(integration.id)}
                    >
                        Por defecto
                    </Button>
                )}
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleDeleteClick(integration.id)}
                >
                    Eliminar
                </Button>
            </div>
        )
    });

    return (
        <div className="space-y-4">
            {/* Filters */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                <Input
                    type="text"
                    placeholder="Buscar por nombre..."
                    value={search}
                    onChange={(e) => setSearch(e.target.value)}
                />
                <Input
                    type="text"
                    placeholder="Filtrar por tipo..."
                    value={filterType}
                    onChange={(e) => setFilterType(e.target.value)}
                />
                <Input
                    type="text"
                    placeholder="Filtrar por categoría..."
                    value={filterCategory}
                    onChange={(e) => setFilterCategory(e.target.value)}
                />
            </div>

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <Table
                columns={columns}
                data={integrations.map(renderRow)}
                emptyMessage="No hay integraciones disponibles"
            />

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex justify-center items-center gap-2">
                    <Button
                        onClick={() => setPage(page - 1)}
                        disabled={page === 1}
                        variant="primary"
                    >
                        Anterior
                    </Button>
                    <span className="text-sm text-gray-700">
                        Página {page} de {totalPages}
                    </span>
                    <Button
                        onClick={() => setPage(page + 1)}
                        disabled={page === totalPages}
                        variant="primary"
                    >
                        Siguiente
                    </Button>
                </div>
            )}

            <ConfirmModal
                isOpen={deleteModal.show}
                onClose={() => setDeleteModal({ show: false, id: null })}
                onConfirm={handleDeleteConfirm}
                title="Eliminar Integración"
                message="¿Estás seguro de que deseas eliminar esta integración? Esta acción no se puede deshacer."
            />
        </div>
    );
}
