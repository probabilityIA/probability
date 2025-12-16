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
        store_name: initialData?.config?.store_name || '',
        api_version: initialData?.config?.api_version || '2024-01',
        access_token: initialData?.credentials?.access_token || '',
        business_id: initialData?.business_id || null,
    });

    // Función para generar el código automáticamente desde el nombre
    const generateCode = (name: string): string => {
        if (!name) return '';
        return name
            .toLowerCase()
            .trim()
            .replace(/\s+/g, '_')
            .replace(/[^a-z0-9_]/g, '')
            .replace(/_+/g, '_')
            .replace(/^_|_$/g, '');
    };

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
            setError('Store Name y Access Token son requeridos para probar la conexi?n');
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
            setError(err.message || 'Error al probar la conexi?n');
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

            // Generar código automáticamente desde el nombre (solo si no estamos editando o no hay código inicial)
            const generatedCode = isEdit && initialData?.code 
                ? initialData.code 
                : generateCode(formData.name);

            await onSubmit({
                name: formData.name,
                code: generatedCode,
                config,
                credentials,
                business_id: formData.business_id,
            });
        } catch (err: any) {
            console.error('Error saving Shopify integration:', err);
            setError(err.message || 'Error al guardar la integraci?n');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6 w-full">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {testSuccess && (
                <Alert type="success" onClose={() => setTestSuccess(false)}>
                    ? Conexi?n exitosa con Shopify
                </Alert>
            )}

            {/* Formulario en una sola tarjeta - 2 columnas, 2 filas */}
            <div className="p-6 rounded-lg border border-gray-200 bg-white">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Fila 1, Columna 1: Nombre de la Integración */}
                    <div className="min-w-0">
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Nombre de la Integración *
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="Ej: Tienda Principal"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            className="w-full"
                        />
                    </div>

                    {/* Fila 1, Columna 2: Store Name */}
                    <div className="min-w-0">
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Store Name *
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="mystore.myshopify.com"
                            value={formData.store_name}
                            onChange={(e) => setFormData({ ...formData, store_name: e.target.value })}
                            className="w-full"
                        />
                    </div>

                    {/* Fila 2, Columna 1: API Version */}
                    <div className="min-w-0">
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            API Version
                        </label>
                        <Select
                            value={formData.api_version}
                            onChange={(e) => setFormData({ ...formData, api_version: e.target.value })}
                            options={apiVersions}
                            className="w-full"
                        />
                    </div>

                    {/* Fila 2, Columna 2: Access Token */}
                    <div className="min-w-0">
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Access Token *
                        </label>
                        <Input
                            type="password"
                            required
                            placeholder="shpat_xxxxxxxxxxxxx"
                            value={formData.access_token}
                            onChange={(e) => setFormData({ ...formData, access_token: e.target.value })}
                            className="w-full"
                        />
                    </div>
                </div>
            </div>

            {/* Action Buttons */}
            <div className="flex flex-row justify-end gap-3 pt-4 border-t">
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
                    type="button"
                    onClick={handleTestConnection}
                    disabled={testing || !formData.store_name || !formData.access_token}
                    loading={testing}
                    variant="outline"
                >
                    {testing ? 'Probando conexión...' : '🔌 Probar Conexión'}
                </Button>
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
