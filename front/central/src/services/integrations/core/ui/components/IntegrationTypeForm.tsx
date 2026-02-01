'use client';

import { useState, useEffect } from 'react';
import { useIntegrationTypes } from '../hooks/useIntegrationTypes';
import { IntegrationType, CreateIntegrationTypeDTO, UpdateIntegrationTypeDTO } from '../../domain/types';
import { Input, Select, Button, Alert, FileInput } from '@/shared/ui';

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
        category: 'internal',
        is_active: true,
        config_schema: '{}',
        credentials_schema: '{}',
        setup_instructions: '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [removeImage, setRemoveImage] = useState(false);
    const [imagePreview, setImagePreview] = useState<string | null>(null);

    useEffect(() => {
        if (integrationType) {
            setFormData({
                name: integrationType.name,
                code: integrationType.code,
                description: integrationType.description || '',
                category: integrationType.category?.code || integrationType.integration_category?.code || 'internal',
                is_active: integrationType.is_active,
                config_schema: JSON.stringify(integrationType.config_schema || {}, null, 2),
                credentials_schema: JSON.stringify(integrationType.credentials_schema || {}, null, 2),
                setup_instructions: integrationType.setup_instructions || '',
            });
            // Cargar preview de imagen existente si hay
            if (integrationType.image_url) {
                setImagePreview(integrationType.image_url);
            }
        }
    }, [integrationType]);

    const handleImageChange = (file: File | null) => {
        setImageFile(file);
        setRemoveImage(false);
        if (file) {
            // Crear preview de la nueva imagen
            const reader = new FileReader();
            reader.onloadend = () => {
                setImagePreview(reader.result as string);
            };
            reader.readAsDataURL(file);
        } else {
            // Si se elimina el archivo seleccionado, volver a la imagen original
            setImagePreview(integrationType?.image_url || null);
        }
    };

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
                    category: formData.category,
                    is_active: formData.is_active,
                    config_schema: formData.config_schema ? JSON.parse(formData.config_schema) : undefined,
                    credentials_schema: formData.credentials_schema ? JSON.parse(formData.credentials_schema) : undefined,
                    setup_instructions: formData.setup_instructions,
                    image_file: imageFile || undefined,
                    remove_image: removeImage || undefined,
                };
                success = await updateIntegrationType(integrationType.id, updateData);
            } else {
                // Create
                const createData: CreateIntegrationTypeDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    category: formData.category,
                    is_active: formData.is_active,
                    config_schema: formData.config_schema ? JSON.parse(formData.config_schema) : undefined,
                    credentials_schema: formData.credentials_schema ? JSON.parse(formData.credentials_schema) : undefined,
                    setup_instructions: formData.setup_instructions,
                    image_file: imageFile || undefined,
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
        <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {/* Basic Info - 2 columns */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Nombre *
                    </label>
                    <Input
                        type="text"
                        required
                        value={formData.name}
                        onChange={(e) => {
                            const name = e.target.value;
                            setFormData({
                                ...formData,
                                name,
                                // Auto-generate code from name if creating new
                                code: integrationType ? formData.code : name.toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_]/g, '')
                            });
                        }}
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
            </div>

            {/* Image Upload Section */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                    Logo del Tipo de Integración
                </label>
                <div className="space-y-4">
                    {/* Image Preview */}
                    {imagePreview && (
                        <div className="flex items-center gap-4">
                            <img
                                src={imagePreview}
                                alt="Preview"
                                className="w-24 h-24 object-contain border border-gray-300 rounded-lg p-2 bg-gray-50"
                            />
                            <div className="flex-1">
                                <p className="text-sm text-gray-600">
                                    {imageFile ? 'Nueva imagen seleccionada' : 'Imagen actual'}
                                </p>
                            </div>
                        </div>
                    )}

                    {/* File Input */}
                    <FileInput
                        accept="image/*"
                        onChange={handleImageChange}
                        buttonText="Seleccionar imagen"
                        helperText="Formatos soportados: JPG, PNG, GIF, WEBP. Tamaño máximo: 10MB"
                    />

                    {/* Remove Image Option (only when editing and has existing image) */}
                    {integrationType && integrationType.image_url && (
                        <div className="flex items-center">
                            <label className="flex items-center">
                                <input
                                    type="checkbox"
                                    checked={removeImage}
                                    onChange={(e) => {
                                        setRemoveImage(e.target.checked);
                                        if (e.target.checked) {
                                            setImageFile(null);
                                            setImagePreview(null);
                                        } else {
                                            setImagePreview(integrationType.image_url || null);
                                        }
                                    }}
                                    className="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                                />
                                <span className="text-sm text-gray-700">Eliminar imagen actual</span>
                            </label>
                        </div>
                    )}
                </div>
            </div>

            {/* Description - Full width */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Descripción
                </label>
                <textarea
                    value={formData.description}
                    onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                    rows={2}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900"
                />
            </div>

            {/* JSON Editors - 2 columns */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {/* Config Schema JSON Editor */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Config Schema (JSON)
                    </label>
                    <textarea
                        value={formData.config_schema}
                        onChange={(e) => setFormData({ ...formData, config_schema: e.target.value })}
                        rows={12}
                        className="w-full px-3 py-2 bg-gray-900 text-green-400 border border-gray-700 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-xs"
                        placeholder='{"type": "object", "properties": {...}}'
                        spellCheck={false}
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        Campos de configuración (no sensibles)
                    </p>
                </div>

                {/* Credentials Schema JSON Editor */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Credentials Schema (JSON)
                    </label>
                    <textarea
                        value={formData.credentials_schema}
                        onChange={(e) => setFormData({ ...formData, credentials_schema: e.target.value })}
                        rows={12}
                        className="w-full px-3 py-2 bg-gray-900 text-green-400 border border-gray-700 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-xs"
                        placeholder='{"type": "object", "properties": {...}}'
                        spellCheck={false}
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        Campos de credenciales (tokens, keys, etc.)
                    </p>
                </div>
            </div>

            {/* Setup Instructions - Full width */}
            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                    Instrucciones de Configuración
                </label>
                <textarea
                    value={formData.setup_instructions}
                    onChange={(e) => setFormData({ ...formData, setup_instructions: e.target.value })}
                    rows={6}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-gray-900"
                    placeholder="Pasos para configurar esta integración:&#10;&#10;1. Ve a...&#10;2. Configura...&#10;3. Copia..."
                />
                <p className="mt-1 text-xs text-gray-500">
                    Instrucciones paso a paso para el usuario
                </p>
            </div>

            {/* Active Checkbox */}
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

            <div className="flex justify-end space-x-3 pt-4 border-t">
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
