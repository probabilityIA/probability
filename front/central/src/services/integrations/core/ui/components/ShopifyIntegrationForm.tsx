'use client';

import { useState } from 'react';
import { Input, Button, Alert, Select } from '@/shared/ui';

interface ShopifyConfig {
    store_name: string;
    api_version?: string;
}

interface ShopifyCredentials {
    access_token: string;
}

interface ShopifyIntegrationFormProps {
    onSubmit: (data: {
        name: string;
        code: string;
        config: ShopifyConfig;
        credentials: ShopifyCredentials;
        business_id?: number | null;
    }) => Promise<void>;
    onCancel?: () => void;
    onTestConnection?: (config: ShopifyConfig, credentials: ShopifyCredentials) => Promise<boolean>;
    initialData?: {
        name?: string;
        code?: string;
        config?: ShopifyConfig;
        credentials?: ShopifyCredentials;
        business_id?: number | null;
    };
    isEdit?: boolean;
}

export default function ShopifyIntegrationForm({
    onSubmit,
    onCancel,
    onTestConnection,
    initialData,
    isEdit = false
}: ShopifyIntegrationFormProps) {
    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        code: initialData?.code || '',
        store_name: initialData?.config?.store_name || '',
        api_version: initialData?.config?.api_version || '2024-01',
        access_token: initialData?.credentials?.access_token || '',
        business_id: initialData?.business_id || null,
    });

    const [loading, setLoading] = useState(false);
    const [testing, setTesting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [testSuccess, setTestSuccess] = useState(false);

    const apiVersions = [
        { value: '2024-01', label: '2024-01' },
        { value: '2024-04', label: '2024-04' },
        { value: '2024-07', label: '2024-07' },
        { value: '2024-10', label: '2024-10' },
    ];

    const handleTestConnection = async () => {
        if (!formData.store_name || !formData.access_token) {
            setError('Store Name y Access Token son requeridos para probar la conexión');
            return;
        }

        setTesting(true);
        setError(null);
        setTestSuccess(false);

        try {
            const config: ShopifyConfig = {
                store_name: formData.store_name,
                api_version: formData.api_version,
            };

            const credentials: ShopifyCredentials = {
                access_token: formData.access_token,
            };

            if (onTestConnection) {
                const success = await onTestConnection(config, credentials);
                if (success) {
                    setTestSuccess(true);
                    setError(null);
                } else {
                    setError('No se pudo conectar con Shopify. Verifica tus credenciales.');
                }
            }
        } catch (err: any) {
            console.error('Test connection error:', err);
            setError(err.message || 'Error al probar la conexión');
        } finally {
            setTesting(false);
        }
    };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            const config: ShopifyConfig = {
                store_name: formData.store_name,
                api_version: formData.api_version,
            };

            const credentials: ShopifyCredentials = {
                access_token: formData.access_token,
            };

            await onSubmit({
                name: formData.name,
                code: formData.code,
                config,
                credentials,
                business_id: formData.business_id,
            });
        } catch (err: any) {
            console.error('Error saving Shopify integration:', err);
            setError(err.message || 'Error al guardar la integración');
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

            {testSuccess && (
                <Alert type="success" onClose={() => setTestSuccess(false)}>
                    ✓ Conexión exitosa con Shopify
                </Alert>
            )}

            {/* Basic Info - 2 columns */}
            <div className="p-6 bg-gray-50 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Información Básica</h3>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Nombre de la Integración *
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="Ej: Tienda Principal"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        />
                        <p className="mt-1 text-xs text-gray-500">Nombre descriptivo para identificar esta integración</p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Código Único *
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="Ej: shopify_main"
                            value={formData.code}
                            onChange={(e) => setFormData({ ...formData, code: e.target.value.toLowerCase().replace(/\s+/g, '_') })}
                            disabled={isEdit}
                        />
                        <p className="mt-1 text-xs text-gray-500">Identificador único (letras, números y guiones bajos)</p>
                    </div>
                </div>
            </div>

            {/* Shopify Configuration - 2 columns */}
            <div className="p-6 bg-blue-50 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Configuración de Shopify</h3>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Store Name * (Nombre de la tienda)
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="mystore.myshopify.com"
                            value={formData.store_name}
                            onChange={(e) => setFormData({ ...formData, store_name: e.target.value })}
                        />
                        <p className="mt-1 text-xs text-gray-500">Nombre completo de tu tienda Shopify</p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            API Version
                        </label>
                        <Select
                            value={formData.api_version}
                            onChange={(e) => setFormData({ ...formData, api_version: e.target.value })}
                            options={apiVersions}
                        />
                        <p className="mt-1 text-xs text-gray-500">Versión de la API de Shopify a utilizar</p>
                    </div>
                </div>
            </div>

            {/* Shopify Credentials - Full width */}
            <div className="p-6 bg-yellow-50 rounded-lg">
                <h3 className="text-base font-semibold text-gray-800 mb-4">Credenciales de Shopify</h3>

                <div className="max-w-2xl">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">
                            Access Token *
                        </label>
                        <Input
                            type="password"
                            required
                            placeholder="shpat_xxxxxxxxxxxxx"
                            value={formData.access_token}
                            onChange={(e) => setFormData({ ...formData, access_token: e.target.value })}
                        />
                        <p className="mt-1 text-xs text-gray-500">Token de acceso de la API de Shopify Admin</p>
                    </div>

                    {/* Test Connection Button */}
                    <div className="mt-4">
                        <Button
                            type="button"
                            onClick={handleTestConnection}
                            disabled={testing || !formData.store_name || !formData.access_token}
                            loading={testing}
                            variant="outline"
                            className="w-full"
                        >
                            {testing ? 'Probando conexión...' : '🔌 Probar Conexión'}
                        </Button>
                    </div>
                </div>
            </div>

            {/* Action Buttons */}
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
                    {isEdit ? 'Actualizar Integración' : 'Crear Integración'}
                </Button>
            </div>
        </form>
    );
}
