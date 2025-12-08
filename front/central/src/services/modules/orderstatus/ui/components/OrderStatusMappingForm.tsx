'use client';

import { useState, useEffect } from 'react';
import { OrderStatusMapping, CreateOrderStatusMappingDTO, UpdateOrderStatusMappingDTO } from '../../domain/types';
import { createOrderStatusMappingAction, updateOrderStatusMappingAction } from '../../infra/actions';
import { Button, Alert, Input, Select } from '@/shared/ui';

interface OrderStatusMappingFormProps {
    mapping?: OrderStatusMapping;
    onSuccess: () => void;
    onCancel: () => void;
}

export default function OrderStatusMappingForm({ mapping, onSuccess, onCancel }: OrderStatusMappingFormProps) {
    const [formData, setFormData] = useState<CreateOrderStatusMappingDTO>({
        integration_type: mapping?.integration_type || '',
        original_status: mapping?.original_status || '',
        mapped_status: mapping?.mapped_status || '',
        priority: mapping?.priority || 0,
        description: mapping?.description || '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [success, setSuccess] = useState<string | null>(null);

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
                    mapped_status: formData.mapped_status,
                    priority: formData.priority,
                    description: formData.description,
                };
                response = await updateOrderStatusMappingAction(mapping.id, updateData);
            } else {
                // Create
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
                        value={formData.integration_type}
                        onChange={(e) => handleChange('integration_type', e.target.value)}
                        required
                        disabled={!!mapping}
                        placeholder="Seleccionar tipo"
                        options={[
                            { value: 'shopify', label: 'Shopify' },
                            { value: 'whatsapp', label: 'WhatsApp' },
                            { value: 'woocommerce', label: 'WooCommerce' },
                        ]}
                    />
                </div>

                {/* Original Status */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Estado Original <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.original_status}
                        onChange={(e) => handleChange('original_status', e.target.value)}
                        placeholder="ej: pending, fulfilled, etc."
                        required
                    />
                </div>

                {/* Mapped Status */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Estado Mapeado <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.mapped_status}
                        onChange={(e) => handleChange('mapped_status', e.target.value)}
                        placeholder="ej: pending, processing, etc."
                        required
                    />
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
