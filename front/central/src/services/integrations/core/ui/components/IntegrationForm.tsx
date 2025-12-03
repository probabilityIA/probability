'use client';

import { useState } from 'react';
import { createIntegrationAction, updateIntegrationAction } from '../../infra/actions';
import { Integration, CreateIntegrationDTO, UpdateIntegrationDTO } from '../../domain/types';
import { Input, Select, Button, Alert } from '@/shared/ui';

interface IntegrationFormProps {
    integration?: Integration;
    onSuccess?: () => void;
    onCancel?: () => void;
}

export default function IntegrationForm({ integration, onSuccess, onCancel }: IntegrationFormProps) {
    const [formData, setFormData] = useState({
        name: integration?.name || '',
        code: integration?.code || '',
        type: integration?.type || 'whatsapp',
        category: integration?.category || 'external',
        business_id: integration?.business_id || null,
        description: integration?.description || '',
        is_active: integration?.is_active ?? true,
        is_default: integration?.is_default ?? false,
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            if (integration) {
                // Update
                const updateData: UpdateIntegrationDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    is_active: formData.is_active,
                    is_default: formData.is_default,
                };
                await updateIntegrationAction(integration.id, updateData);
            } else {
                // Create
                const createData: CreateIntegrationDTO = {
                    ...formData,
                    business_id: formData.business_id || null,
                };
                await createIntegrationAction(createData);
            }
            onSuccess?.();
        } catch (err: any) {
            console.error('Error saving integration:', err);
            setError(err.message || 'Error al guardar la integración');
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
                    disabled={!!integration}
                />
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Tipo *
                </label>
                <Select
                    required
                    value={formData.type}
                    onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                    disabled={!!integration}
                    options={[
                        { value: 'whatsapp', label: 'WhatsApp' },
                        { value: 'shopify', label: 'Shopify' },
                        { value: 'mercado_libre', label: 'Mercado Libre' }
                    ]}
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
                    disabled={!!integration}
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

                <label className="flex items-center">
                    <input
                        type="checkbox"
                        checked={formData.is_default}
                        onChange={(e) => setFormData({ ...formData, is_default: e.target.checked })}
                        className="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <span className="text-sm font-medium text-gray-700">Por defecto</span>
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
                    {integration ? 'Actualizar' : 'Crear'}
                </Button>
            </div>
        </form>
    );
}
