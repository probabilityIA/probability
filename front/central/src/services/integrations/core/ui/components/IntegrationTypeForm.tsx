'use client';

import { useState, useEffect } from 'react';
import { useIntegrationTypes } from '../hooks/useIntegrationTypes';
import { IntegrationType, CreateIntegrationTypeDTO, UpdateIntegrationTypeDTO } from '../../domain/types';
import { Input, Select, Button, Alert } from '@/shared/ui';

interface IntegrationTypeFormProps {
    integrationType?: IntegrationType;
    onSuccess?: () => void;
    onCancel?: () => void;
}

export default function IntegrationTypeForm({ integrationType, onSuccess, onCancel }: IntegrationTypeFormProps) {
    const { createIntegrationType, updateIntegrationType } = useIntegrationTypes();

    const [formData, setFormData] = useState({
        name: '',
        code: '',
        description: '',
        icon: '',
        category: 'internal',
        is_active: true,
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (integrationType) {
            setFormData({
                name: integrationType.name,
                code: integrationType.code,
                description: integrationType.description || '',
                icon: integrationType.icon || '',
                category: integrationType.category,
                is_active: integrationType.is_active,
            });
        }
    }, [integrationType]);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            let success = false;
            if (integrationType) {
                // Update
                const updateData: UpdateIntegrationTypeDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    icon: formData.icon,
                    category: formData.category,
                    is_active: formData.is_active,
                };
                success = await updateIntegrationType(integrationType.id, updateData);
            } else {
                // Create
                const createData: CreateIntegrationTypeDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    icon: formData.icon,
                    category: formData.category,
                    is_active: formData.is_active,
                };
                success = await createIntegrationType(createData);
            }

            if (success) {
                onSuccess?.();
            }
        } catch (err: any) {
            console.error('Error saving integration type:', err);
            setError(err.message || 'Error al guardar el tipo de integración');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Nombre *
                </label>
                <Input
                    type="text"
                    required
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Código *
                </label>
                <Input
                    type="text"
                    required
                    value={formData.code}
                    onChange={(e) => setFormData({ ...formData, code: e.target.value })}
                    disabled={!!integrationType} // Code usually shouldn't change
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Categoría *
                </label>
                <Select
                    required
                    value={formData.category}
                    onChange={(e) => setFormData({ ...formData, category: e.target.value })}
                    options={[
                        { value: 'internal', label: 'Interna' },
                        { value: 'external', label: 'Externa' }
                    ]}
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Descripción
                </label>
                <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    rows={3}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Icono
                </label>
                <Input
                    type="text"
                    value={formData.icon}
                    onChange={(e) => setFormData({ ...formData, icon: e.target.value })}
                    placeholder="Nombre del icono (ej: whatsapp-icon)"
                />
            </div>

            <div className="flex items-center space-x-4">
                <label className="flex items-center">
                    <input
                        type="checkbox"
                        checked={formData.is_active}
                        onChange={(e) => setFormData({ ...formData, is_active: e.target.checked })}
                        className="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <span className="text-sm font-medium text-gray-700">Activo</span>
                </label>
            </div>

            <div className="flex justify-end space-x-3 pt-4">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        variant="outline"
                    >
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    disabled={loading}
                    loading={loading}
                    variant="primary"
                >
                    {integrationType ? 'Actualizar' : 'Crear'}
                </Button>
            </div>
        </form>
    );
}
