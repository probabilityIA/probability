'use client';

import { useState, useEffect } from 'react';
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import { useIntegrationTypes } from '../hooks/useIntegrationTypes';
import { IntegrationType, CreateIntegrationTypeDTO, UpdateIntegrationTypeDTO, IntegrationCategory } from '../../domain/types';
import { Input, Select, Button, Alert, FileInput } from '@/shared/ui';
import { getIntegrationCategoriesAction, getIntegrationTypePlatformCredentialsAction } from '../../infra/actions';

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
        category_id: 0,
        is_active: true,
        config_schema: '{}',
        credentials_schema: '{}',
        setup_instructions: '',
        base_url: '',
        base_url_test: '',
        platform_credentials: '{}',
    });

    const [categories, setCategories] = useState<IntegrationCategory[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [imageFile, setImageFile] = useState<File | null>(null);
    const [removeImage, setRemoveImage] = useState(false);
    const [imagePreview, setImagePreview] = useState<string | null>(null);
    const [showPlatformCredentials, setShowPlatformCredentials] = useState(false);
    // Campos estructurados para WhatsApp
    const [whatsappCredentials, setWhatsappCredentials] = useState({
        phone_number_id: '',
        access_token: '',
        verify_token: '',
    });
    const [showWhatsappSecrets, setShowWhatsappSecrets] = useState(false);

    useEffect(() => {
        getIntegrationCategoriesAction()
            .then((res) => {
                if (res.success && res.data.length > 0) {
                    setCategories(res.data);
                    // Si no hay tipo existente, usar la primera categoría como default
                    if (!integrationType) {
                        setFormData((prev) => ({ ...prev, category_id: res.data[0].id }));
                    }
                }
            })
            .catch(() => {});
    }, []);

    useEffect(() => {
        if (integrationType) {
            setFormData({
                name: integrationType.name,
                code: integrationType.code,
                description: integrationType.description || '',
                category_id: integrationType.category_id || 0,
                is_active: integrationType.is_active,
                config_schema: JSON.stringify(integrationType.config_schema || {}, null, 2),
                credentials_schema: JSON.stringify(integrationType.credentials_schema || {}, null, 2),
                setup_instructions: integrationType.setup_instructions || '',
                base_url: integrationType.base_url || '',
                base_url_test: integrationType.base_url_test || '',
                platform_credentials: '{}',
            });
            // Cargar preview de imagen existente si hay
            if (integrationType.image_url) {
                setImagePreview(integrationType.image_url);
            }
            // Cargar credenciales de plataforma desencriptadas si existen
            if (integrationType.has_platform_credentials) {
                getIntegrationTypePlatformCredentialsAction(integrationType.id)
                    .then((res) => {
                        if (res.success && res.data && Object.keys(res.data).length > 0) {
                            if (integrationType.code === 'whatsapp') {
                                // Poblar campos estructurados de WhatsApp
                                setWhatsappCredentials({
                                    phone_number_id: res.data.phone_number_id || '',
                                    access_token: res.data.access_token || '',
                                    verify_token: res.data.verify_token || '',
                                });
                            } else {
                                setFormData((prev) => ({
                                    ...prev,
                                    platform_credentials: JSON.stringify(res.data, null, 2),
                                }));
                            }
                        }
                    })
                    .catch(() => {});
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
            // Parse platform_credentials — only send if non-empty
            let platformCredentials: Record<string, string> | undefined;
            if (formData.code === 'whatsapp') {
                // Para WhatsApp, construir desde campos estructurados
                const wa: Record<string, string> = {};
                if (whatsappCredentials.phone_number_id.trim()) wa.phone_number_id = whatsappCredentials.phone_number_id.trim();
                if (whatsappCredentials.access_token.trim()) wa.access_token = whatsappCredentials.access_token.trim();
                if (whatsappCredentials.verify_token.trim()) wa.verify_token = whatsappCredentials.verify_token.trim();
                if (Object.keys(wa).length > 0) platformCredentials = wa;
            } else {
                try {
                    const parsed = formData.platform_credentials ? JSON.parse(formData.platform_credentials) : {};
                    if (parsed && typeof parsed === 'object' && Object.keys(parsed).length > 0) {
                        platformCredentials = parsed;
                    }
                } catch {
                    throw new Error('Las credenciales de plataforma no son un JSON válido');
                }
            }

            if (integrationType) {
                // Update
                const updateData: UpdateIntegrationTypeDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    category_id: formData.category_id,
                    is_active: formData.is_active,
                    config_schema: formData.config_schema ? JSON.parse(formData.config_schema) : undefined,
                    credentials_schema: formData.credentials_schema ? JSON.parse(formData.credentials_schema) : undefined,
                    setup_instructions: formData.setup_instructions,
                    image_file: imageFile || undefined,
                    remove_image: removeImage || undefined,
                    base_url: formData.base_url || undefined,
                    base_url_test: formData.base_url_test || undefined,
                    platform_credentials: platformCredentials,
                };
                success = await updateIntegrationType(integrationType.id, updateData);
            } else {
                // Create
                const createData: CreateIntegrationTypeDTO = {
                    name: formData.name,
                    code: formData.code,
                    description: formData.description,
                    category_id: formData.category_id,
                    is_active: formData.is_active,
                    config_schema: formData.config_schema ? JSON.parse(formData.config_schema) : undefined,
                    credentials_schema: formData.credentials_schema ? JSON.parse(formData.credentials_schema) : undefined,
                    setup_instructions: formData.setup_instructions,
                    image_file: imageFile || undefined,
                    base_url: formData.base_url || undefined,
                    base_url_test: formData.base_url_test || undefined,
                    platform_credentials: platformCredentials,
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
                        value={String(formData.category_id)}
                        onChange={(e) => setFormData({ ...formData, category_id: Number(e.target.value) })}
                        options={categories.map((cat) => ({ value: String(cat.id), label: cat.name }))}
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

            {/* URLs de la API - 2 columns */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        URL de Producción
                    </label>
                    <Input
                        type="url"
                        value={formData.base_url}
                        onChange={(e) => setFormData({ ...formData, base_url: e.target.value })}
                        placeholder="https://api.ejemplo.com/v1"
                        className="font-mono text-sm"
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        URL base del API en producción
                    </p>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        URL de Pruebas (Sandbox)
                    </label>
                    <Input
                        type="url"
                        value={formData.base_url_test}
                        onChange={(e) => setFormData({ ...formData, base_url_test: e.target.value })}
                        placeholder="https://sandbox.ejemplo.com/v1"
                        className="font-mono text-sm"
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        URL del entorno sandbox para modo de pruebas
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

            {/* Platform Credentials (encrypted) */}
            {formData.code === 'whatsapp' ? (
                /* Campos estructurados para WhatsApp */
                <div className="space-y-4">
                    <div className="flex items-center justify-between">
                        <h3 className="text-sm font-medium text-gray-700">
                            Credenciales de WhatsApp (se encriptarán)
                        </h3>
                        <button
                            type="button"
                            onClick={() => setShowWhatsappSecrets((v) => !v)}
                            className="flex items-center gap-1 text-xs text-gray-500 hover:text-gray-700"
                        >
                            {showWhatsappSecrets ? (
                                <><EyeSlashIcon className="w-4 h-4" /> Ocultar</>
                            ) : (
                                <><EyeIcon className="w-4 h-4" /> Mostrar</>
                            )}
                        </button>
                    </div>
                    <div className="grid grid-cols-1 gap-4 p-4 border border-gray-200 rounded-lg bg-gray-50">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Phone Number ID
                            </label>
                            <Input
                                type="text"
                                value={whatsappCredentials.phone_number_id}
                                onChange={(e) => setWhatsappCredentials({ ...whatsappCredentials, phone_number_id: e.target.value })}
                                placeholder="123456789012345"
                                className="font-mono"
                            />
                            <p className="mt-1 text-xs text-gray-500">
                                ID del número de teléfono en Meta Business Manager
                            </p>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Access Token
                            </label>
                            <Input
                                type={showWhatsappSecrets ? 'text' : 'password'}
                                value={whatsappCredentials.access_token}
                                onChange={(e) => setWhatsappCredentials({ ...whatsappCredentials, access_token: e.target.value })}
                                placeholder="EAAxxxxxxxxx..."
                                className="font-mono"
                            />
                            <p className="mt-1 text-xs text-gray-500">
                                Token de acceso permanente de la app en Meta
                            </p>
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">
                                Verify Token (Webhook)
                            </label>
                            <Input
                                type={showWhatsappSecrets ? 'text' : 'password'}
                                value={whatsappCredentials.verify_token}
                                onChange={(e) => setWhatsappCredentials({ ...whatsappCredentials, verify_token: e.target.value })}
                                placeholder="mi_token_secreto"
                                className="font-mono"
                            />
                            <p className="mt-1 text-xs text-gray-500">
                                Token de verificación para el webhook (Meta → WhatsApp → Configuration)
                            </p>
                        </div>
                    </div>
                    <p className="text-xs text-gray-500">
                        Deja los campos vacíos para no modificar los valores actuales.
                    </p>
                </div>
            ) : (
                /* Editor JSON genérico para otros tipos */
                <div>
                    <div className="flex items-center justify-between mb-1">
                        <label className="block text-sm font-medium text-gray-700">
                            Credenciales de Plataforma (JSON)
                        </label>
                        <button
                            type="button"
                            onClick={() => setShowPlatformCredentials((v) => !v)}
                            className="flex items-center gap-1 text-xs text-gray-500 hover:text-gray-700"
                        >
                            {showPlatformCredentials ? (
                                <>
                                    <EyeSlashIcon className="w-4 h-4" />
                                    Ocultar
                                </>
                            ) : (
                                <>
                                    <EyeIcon className="w-4 h-4" />
                                    Mostrar
                                </>
                            )}
                        </button>
                    </div>
                    <textarea
                        value={showPlatformCredentials
                            ? formData.platform_credentials
                            : formData.platform_credentials.replace(/:\s*"([^"]*)"/g, ': "••••••••"')
                        }
                        onChange={(e) => {
                            if (showPlatformCredentials) {
                                setFormData({ ...formData, platform_credentials: e.target.value });
                            }
                        }}
                        readOnly={!showPlatformCredentials}
                        rows={6}
                        className={`w-full px-3 py-2 bg-gray-900 border border-gray-700 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-xs ${showPlatformCredentials ? 'text-green-400' : 'text-gray-500 cursor-default'}`}
                        placeholder={showPlatformCredentials ? '{\n  "api_key": "tu-api-key-aqui"\n}' : ''}
                        spellCheck={false}
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        Credenciales globales del proveedor (se encriptarán). Deja <code>{'{}'}</code> para no cambiarlas.
                    </p>
                </div>
            )}

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
