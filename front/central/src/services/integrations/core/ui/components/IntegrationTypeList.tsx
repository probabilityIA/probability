'use client';

import { useState } from 'react';
import { useIntegrationTypes } from '../hooks/useIntegrationTypes';
import { IntegrationType } from '../../domain/types';
import { Button, Badge, Spinner, Table, Alert, ConfirmModal } from '@/shared/ui';

interface IntegrationTypeListProps {
    onEdit?: (integrationType: IntegrationType) => void;
}

export default function IntegrationTypeList({ onEdit }: IntegrationTypeListProps) {
    const {
        integrationTypes,
        loading,
        error,
        setError,
        deleteIntegrationType,
        refresh
    } = useIntegrationTypes();

    const [deleteModal, setDeleteModal] = useState<{ show: boolean; id: number | null }>({
        show: false,
        id: null
    });

    const handleDeleteClick = (id: number) => {
        setDeleteModal({ show: true, id });
    };

    const handleDeleteConfirm = async () => {
        if (deleteModal.id) {
            const success = await deleteIntegrationType(deleteModal.id);
            if (success) {
                setDeleteModal({ show: false, id: null });
            }
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
        { key: 'logo', label: 'Logo' },
        { key: 'name', label: 'Nombre' },
        { key: 'category', label: 'Categoría' },
        { key: 'status', label: 'Estado' },
        { key: 'actions', label: 'Acciones' }
    ];

    const renderRow = (type: IntegrationType) => ({
        id: type.id,
        logo: (
            <div className="flex items-center justify-center">
                {type.image_url ? (
                    <img
                        src={type.image_url}
                        alt={type.name}
                        className="w-12 h-12 object-contain border border-gray-200 rounded-lg p-1 bg-white"
                        onError={(e) => {
                            // Si la imagen falla al cargar, mostrar un placeholder
                            (e.target as HTMLImageElement).style.display = 'none';
                            const parent = (e.target as HTMLImageElement).parentElement;
                            if (parent) {
                                parent.innerHTML = '<div class="w-12 h-12 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-xs">Sin logo</div>';
                            }
                        }}
                    />
                ) : (
                    <div className="w-12 h-12 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-xs">
                        Sin logo
                    </div>
                )}
            </div>
        ),
        name: (
            <div>
                <div className="text-sm font-medium text-gray-900">{type.name}</div>
                {type.description && (
                    <div className="text-sm text-gray-500">{type.description}</div>
                )}
            </div>
        ),
        category: type.category ? (
            <span
                className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                style={{ backgroundColor: type.category.color }}
            >
                {type.category.name}
            </span>
        ) : (
            <span className="text-gray-400 text-sm">Sin categoría</span>
        ),
        status: (
            <Badge type={type.is_active ? 'success' : 'error'}>
                {type.is_active ? 'Activo' : 'Inactivo'}
            </Badge>
        ),
        actions: (
            <div className="flex gap-2">
                {onEdit && (
                    <Button
                        variant="outline"
                        size="sm"
                        onClick={() => onEdit(type)}
                    >
                        Editar
                    </Button>
                )}
                <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handleDeleteClick(type.id)}
                >
                    Eliminar
                </Button>
            </div>
        )
    });

    return (
        <div className="space-y-4">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <Table
                columns={columns}
                data={integrationTypes.map(renderRow)}
                emptyMessage="No hay tipos de integración disponibles"
            />

            <ConfirmModal
                isOpen={deleteModal.show}
                onClose={() => setDeleteModal({ show: false, id: null })}
                onConfirm={handleDeleteConfirm}
                title="Eliminar Tipo de Integración"
                message="¿Estás seguro de que deseas eliminar este tipo de integración? Esta acción no se puede deshacer y podría afectar a las integraciones existentes."
            />
        </div>
    );
}
