'use client';

import { OrderStatusMapping } from '../../domain/types';
import { Badge } from '@/shared/ui';

interface OrderStatusMappingDetailsProps {
    mapping: OrderStatusMapping;
}

export default function OrderStatusMappingDetails({ mapping }: OrderStatusMappingDetailsProps) {
    const formatDate = (dateString: string) => {
        return new Date(dateString).toLocaleDateString('es-CO', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit',
        });
    };

    return (
        <div className="space-y-6">
            {/* Header */}
            <div className="flex items-center justify-between pb-4 border-b">
                <h3 className="text-lg font-semibold text-gray-900">
                    Detalles del Mapping
                </h3>
                <Badge type={mapping.is_active ? 'success' : 'secondary'}>
                    {mapping.is_active ? 'Activo' : 'Inactivo'}
                </Badge>
            </div>

            {/* Details Grid */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Integration Type */}
                <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                        Tipo de Integración
                    </label>
                    <p className="text-base font-medium text-gray-900 capitalize">
                        {mapping.integration_type}
                    </p>
                </div>

                {/* Priority */}
                <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                        Prioridad
                    </label>
                    <p className="text-base font-medium text-gray-900">
                        {mapping.priority}
                    </p>
                </div>

                {/* Original Status */}
                <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                        Estado Original
                    </label>
                    <p className="text-base text-gray-900">
                        {mapping.original_status}
                    </p>
                </div>

                {/* Mapped Status */}
                <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                        Estado Mapeado
                    </label>
                    <p className="text-base font-medium text-gray-900">
                        {mapping.mapped_status}
                    </p>
                </div>

                {/* Created At */}
                <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                        Fecha de Creación
                    </label>
                    <p className="text-base text-gray-900">
                        {formatDate(mapping.created_at)}
                    </p>
                </div>

                {/* Updated At */}
                <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                        Última Actualización
                    </label>
                    <p className="text-base text-gray-900">
                        {formatDate(mapping.updated_at)}
                    </p>
                </div>
            </div>

            {/* Description */}
            {mapping.description && (
                <div>
                    <label className="block text-sm font-medium text-gray-500 mb-1">
                        Descripción
                    </label>
                    <p className="text-base text-gray-900 bg-gray-50 p-4 rounded-lg">
                        {mapping.description}
                    </p>
                </div>
            )}

            {/* ID */}
            <div className="pt-4 border-t">
                <label className="block text-sm font-medium text-gray-500 mb-1">
                    ID
                </label>
                <p className="text-sm text-gray-600 font-mono">
                    {mapping.id}
                </p>
            </div>
        </div>
    );
}
