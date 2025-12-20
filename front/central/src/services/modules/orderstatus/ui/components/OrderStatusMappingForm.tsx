'use client';

import { useState, useEffect } from 'react';
import { OrderStatusMapping, CreateOrderStatusMappingDTO, UpdateOrderStatusMappingDTO } from '../../domain/types';
import { createOrderStatusMappingAction, updateOrderStatusMappingAction, getOrderStatusMappingsAction } from '../../infra/actions';
import { getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { Button, Alert, Input, Select } from '@/shared/ui';

interface OrderStatusMappingFormProps {
    mapping?: OrderStatusMapping;
    onSuccess: () => void;
    onCancel: () => void;
}

export default function OrderStatusMappingForm({ mapping, onSuccess, onCancel }: OrderStatusMappingFormProps) {
    const [formData, setFormData] = useState<CreateOrderStatusMappingDTO>({
        integration_type_id: mapping?.integration_type_id || 0,
        original_status: mapping?.original_status || '',
        order_status_id: mapping?.order_status_id || 0,
        priority: mapping?.priority || 0,
        description: mapping?.description || '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);
    const [integrationTypes, setIntegrationTypes] = useState<Array<{ value: number; label: string }>>([]);
    const [orderStatuses, setOrderStatuses] = useState<Array<{ value: number; label: string }>>([]);
    const [loadingData, setLoadingData] = useState(true);

    // Cargar tipos de integración y estados de orden
    useEffect(() => {
        const loadFormData = async () => {
            setLoadingData(true);
            try {
                // Cargar tipos de integración
                const integrationTypesResponse = await getActiveIntegrationTypesAction();
                if (integrationTypesResponse.success && integrationTypesResponse.data) {
                    const types = integrationTypesResponse.data.map(type => ({
                        value: type.id,
                        label: type.name
                    }));
                    setIntegrationTypes(types);
                }

                // Cargar todos los mapeos para extraer los order_statuses únicos
                const mappingsResponse = await getOrderStatusMappingsAction({});
                const mappings = (mappingsResponse as any).data || mappingsResponse.data || [];
                
                // Extraer order_statuses únicos desde los mappings
                const statusMap = new Map<number, { id: number; name: string }>();
                mappings.forEach((m: OrderStatusMapping) => {
                    if (m.order_status && !statusMap.has(m.order_status.id)) {
                        statusMap.set(m.order_status.id, {
                            id: m.order_status.id,
                            name: m.order_status.name
                        });
                    }
                });
                
                const statuses = Array.from(statusMap.values()).map(status => ({
                    value: status.id,
                    label: status.name
                }));
                
                // Ordenar por nombre
                statuses.sort((a, b) => a.label.localeCompare(b.label));
                setOrderStatuses(statuses);

                // Si estamos editando, establecer valores iniciales
                if (mapping) {
                    setFormData({
                        integration_type_id: mapping.integration_type_id,
                        original_status: mapping.original_status,
                        order_status_id: mapping.order_status_id,
                        priority: mapping.priority,
                        description: mapping.description,
                    });
                }
            } catch (err: any) {
                console.error('Error loading form data:', err);
                setError('Error al cargar los datos del formulario');
            } finally {
                setLoadingData(false);
            }
        };

        loadFormData();
    }, [mapping]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setSuccess(null);

        try {
            let response;
            if (mapping) {
                // Update
                const updateData: UpdateOrderStatusMappingDTO = {
                    original_status: formData.original_status,
                    order_status_id: formData.order_status_id,
                    priority: formData.priority,
                    description: formData.description,
                };
                response = await updateOrderStatusMappingAction(mapping.id, updateData);
            } else {
                // Create
                if (!formData.integration_type_id || !formData.order_status_id) {
                    setError('Por favor completa todos los campos requeridos');
                    setLoading(false);
                    return;
                }
                response = await createOrderStatusMappingAction(formData);
            }

            if (response.success) {
                setSuccess(mapping ? 'Mapping actualizado exitosamente' : 'Mapping creado exitosamente');
                setTimeout(() => {
                    onSuccess();
                }, 1000);
            } else {
                setError(response.message || 'Error al guardar el mapping');
            }
        } catch (err: any) {
            setError(err.message || 'Error al guardar el mapping');
        } finally {
            setLoading(false);
        }
    };

    const handleChange = (field: keyof CreateOrderStatusMappingDTO, value: any) => {
        setFormData({ ...formData, [field]: value });
    };

    if (loadingData) {
        return (
            <div className="text-center py-8">
                <div className="text-gray-500">Cargando datos del formulario...</div>
            </div>
        );
    }

    return (
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {success && (
                <Alert type="success" onClose={() => setSuccess(null)}>
                    {success}
                </Alert>
            )}

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Integration Type */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Tipo de Integración <span className="text-red-500">*</span>
                    </label>
                    <Select
                        value={String(formData.integration_type_id)}
                        onChange={(e) => handleChange('integration_type_id', parseInt(e.target.value) || 0)}
                        required
                        disabled={!!mapping}
                        placeholder="Seleccionar tipo de integración"
                        options={integrationTypes.map(type => ({
                            value: String(type.value),
                            label: type.label
                        }))}
                    />
                </div>

                {/* Original Status */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Estado Original (de la Integración) <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.original_status}
                        onChange={(e) => handleChange('original_status', e.target.value)}
                        placeholder="ej: pending, fulfilled, paid, etc."
                        required
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        El estado tal como lo usa la integración externa
                    </p>
                </div>

                {/* Order Status (Mapeado) */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Estado de Probability <span className="text-red-500">*</span>
                    </label>
                    <Select
                        value={String(formData.order_status_id)}
                        onChange={(e) => handleChange('order_status_id', parseInt(e.target.value) || 0)}
                        required
                        placeholder="Seleccionar estado de Probability"
                        options={orderStatuses.map(status => ({
                            value: String(status.value),
                            label: status.label
                        }))}
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        El estado unificado de Probability al que se mapea
                    </p>
                </div>

                {/* Priority */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Prioridad
                    </label>
                    <Input
                        type="number"
                        value={formData.priority}
                        onChange={(e) => handleChange('priority', parseInt(e.target.value) || 0)}
                        placeholder="0"
                        min="0"
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        Mayor prioridad = mayor número
                    </p>
                </div>
            </div>

            {/* Description */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    Descripción
                </label>
                <textarea
                    value={formData.description}
                    onChange={(e) => handleChange('description', e.target.value)}
                    placeholder="Descripción opcional del mapping"
                    rows={3}
                    className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
            </div>

            {/* Actions */}
            <div className="flex justify-end gap-3 pt-4 border-t">
                <Button
                    type="button"
                    variant="outline"
                    onClick={onCancel}
                    disabled={loading}
                >
                    Cancelar
                </Button>
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading}
                >
                    {loading ? 'Guardando...' : mapping ? 'Actualizar' : 'Crear'}
                </Button>
            </div>
        </form>
    );
}
