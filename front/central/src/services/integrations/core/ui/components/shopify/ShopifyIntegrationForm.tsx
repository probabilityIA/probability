'use client';

import { useState, useEffect } from 'react';
import { Input, Button, Alert, Select } from '@/shared/ui';
import { WebhookInfo, ShopifyWebhookInfo } from '../../../domain/types';
import { listWebhooksAction, deleteWebhookAction, verifyWebhooksAction, createWebhookAction } from '../../../infra/actions';

interface ShopifyConfig {
    store_name: string;
    api_version?: string;
    webhook_url?: string;
    webhook_configured?: boolean;
    webhook_ids?: string[];
}

interface ShopifyCredentials {
    access_token: string;
}

interface ShopifyIntegrationFormProps {
    onSubmit: (data: {
        name: string;
        code: string;
        store_id: string;
        config: ShopifyConfig;
        credentials: ShopifyCredentials;
        business_id?: number | null;
    }) => Promise<void>;
    onCancel?: () => void;
    onTestConnection?: (config: ShopifyConfig, credentials: ShopifyCredentials) => Promise<boolean>;
    onGetWebhook?: () => Promise<WebhookInfo | null>;
    initialData?: {
        name?: string;
        code?: string;
        store_id?: string;
        config?: ShopifyConfig;
        credentials?: ShopifyCredentials;
        business_id?: number | null;
    };
    isEdit?: boolean;
    integrationId?: number;
}

export default function ShopifyIntegrationForm({
    onSubmit,
    onCancel,
    onTestConnection,
    onGetWebhook,
    initialData,
    isEdit = false,
    integrationId
}: ShopifyIntegrationFormProps) {
    // Debug logs
    console.log('🛒 ShopifyIntegrationForm initialData:', initialData);
    console.log('🛒 ShopifyIntegrationForm store_id:', initialData?.store_id);
    console.log('🛒 ShopifyIntegrationForm config:', initialData?.config);
    console.log('🛒 ShopifyIntegrationForm credentials:', initialData?.credentials);
    
    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        store_name: initialData?.store_id || initialData?.config?.store_name || '',
        api_version: initialData?.config?.api_version || '2024-01',
        access_token: initialData?.credentials?.access_token || '',
        business_id: initialData?.business_id || null,
    });
    
    // Flag para indicar si las credenciales originales están disponibles
    const hasExistingCredentials = isEdit && !initialData?.credentials?.access_token;

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
    const [webhookInfo, setWebhookInfo] = useState<WebhookInfo | null>(null);
    const [loadingWebhook, setLoadingWebhook] = useState(false);
    const [showWebhook, setShowWebhook] = useState(false);
    const [copied, setCopied] = useState(false);
    const [webhooks, setWebhooks] = useState<ShopifyWebhookInfo[]>([]);
    const [loadingWebhooks, setLoadingWebhooks] = useState(false);
    const [webhookConfigured, setWebhookConfigured] = useState(false);
    const [webhookUrl, setWebhookUrl] = useState<string>('');
    const [showCreateWebhookModal, setShowCreateWebhookModal] = useState(false);
    const [webhookUrlToCreate, setWebhookUrlToCreate] = useState<string>('');
    const [creatingWebhook, setCreatingWebhook] = useState(false);
    const [verifyingWebhooks, setVerifyingWebhooks] = useState(false);

    const apiVersions = [
        { value: '2024-01', label: '2024-01' },
        { value: '2024-04', label: '2024-04' },
        { value: '2024-07', label: '2024-07' },
        { value: '2024-10', label: '2024-10' },
    ];

    // Cargar webhooks y estado al montar el componente en modo edición
    useEffect(() => {
        if (isEdit && integrationId) {
            loadWebhooks();
            // Obtener estado del webhook desde el config
            if (initialData?.config) {
                const config = initialData.config;
                setWebhookConfigured(config.webhook_configured === true);
                setWebhookUrl(config.webhook_url || '');
            }
        }
    }, [isEdit, integrationId]);

    const loadWebhooks = async () => {
        if (!integrationId) return;

        setLoadingWebhooks(true);
        try {
            const response = await listWebhooksAction(integrationId);
            if (response.success && response.data) {
                setWebhooks(response.data);
            }
        } catch (err: any) {
            console.error('Error loading webhooks:', err);
            // No mostrar error si no hay webhooks, es normal
        } finally {
            setLoadingWebhooks(false);
        }
    };

    const handleDeleteWebhook = async (webhookId: string) => {
        if (!integrationId) return;
        
        if (!confirm('¿Estás seguro de que deseas eliminar este webhook?')) {
            return;
        }

        try {
            await deleteWebhookAction(integrationId, webhookId);
            // Recargar la lista de webhooks
            await loadWebhooks();
        } catch (err: any) {
            console.error('Error deleting webhook:', err);
            setError(err.message || 'Error al eliminar el webhook');
        }
    };

    const handleGetWebhook = async () => {
        if (!onGetWebhook) return;
        
        setLoadingWebhook(true);
        setError(null);
        
        try {
            const info = await onGetWebhook();
            if (info) {
                setWebhookInfo(info);
                setShowWebhook(true);
            }
        } catch (err: any) {
            console.error('Error getting webhook:', err);
            setError(err.message || 'Error al obtener el webhook');
        } finally {
            setLoadingWebhook(false);
        }
    };

    const handleCopyWebhook = async () => {
        if (webhookInfo?.url) {
            try {
                await navigator.clipboard.writeText(webhookInfo.url);
                setCopied(true);
                setTimeout(() => setCopied(false), 2000);
            } catch (err) {
                console.error('Error copying to clipboard:', err);
            }
        }
    };

    const handleVerifyAndShowCreateModal = async () => {
        if (!integrationId) return;

        setVerifyingWebhooks(true);
        setError(null);

        try {
            // Primero obtener la URL del webhook
            if (onGetWebhook) {
                const info = await onGetWebhook();
                if (info && info.url) {
                    setWebhookUrlToCreate(info.url);
                    // Verificar si existen webhooks con esta URL
                    const verifyResponse = await verifyWebhooksAction(integrationId);
                    if (verifyResponse.success && verifyResponse.data && verifyResponse.data.length > 0) {
                        // Si ya existen webhooks, mostrar mensaje
                        setError(`Ya existen ${verifyResponse.data.length} webhook(s) con esta URL. Se eliminarán antes de crear nuevos.`);
                    }
                    setShowCreateWebhookModal(true);
                } else {
                    setError('No se pudo obtener la URL del webhook');
                }
            }
        } catch (err: any) {
            console.error('Error verifying webhooks:', err);
            setError(err.message || 'Error al verificar webhooks');
        } finally {
            setVerifyingWebhooks(false);
        }
    };

    const handleCreateWebhook = async () => {
        if (!integrationId) return;

        setCreatingWebhook(true);
        setError(null);

        try {
            const response = await createWebhookAction(integrationId);
            if (response.success) {
                setWebhookConfigured(true);
                setWebhookUrl(response.data.webhook_url);
                setShowCreateWebhookModal(false);
                // Recargar webhooks
                await loadWebhooks();
                setError(null);
                setTestSuccess(true);
                setTimeout(() => setTestSuccess(false), 5000);
            } else {
                setError('No se pudo crear el webhook');
            }
        } catch (err: any) {
            console.error('Error creating webhook:', err);
            const errorMessage = err.message || 'Error al crear el webhook';
            // Si el error menciona localhost, mostrar mensaje más claro
            if (errorMessage.toLowerCase().includes('localhost') || errorMessage.toLowerCase().includes('pruebas')) {
                setError('⚠️ No se pueden crear webhooks en entorno de pruebas (localhost). Los webhooks solo se pueden crear en producción.');
            } else {
                setError(errorMessage);
            }
        } finally {
            setCreatingWebhook(false);
        }
    };

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
                store_id: formData.store_name, // El store_name es el store_id para Shopify
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
                            Access Token {!isEdit && '*'}
                        </label>
                        <Input
                            type="password"
                            required={!isEdit}
                            placeholder={hasExistingCredentials ? "Dejar vacío para mantener el actual" : "shpat_xxxxxxxxxxxxx"}
                            value={formData.access_token}
                            onChange={(e) => setFormData({ ...formData, access_token: e.target.value })}
                            className="w-full"
                        />
                        {hasExistingCredentials && (
                            <p className="text-xs text-gray-500 mt-1">
                                🔒 El token actual está protegido. Solo ingresa un valor si deseas cambiarlo.
                            </p>
                        )}
                    </div>
                </div>
            </div>

            {/* Webhook Section - Solo visible en modo edición */}
            {isEdit && integrationId && (
                <div className="p-6 rounded-lg border border-gray-200 bg-white space-y-6">
                    <div className="flex items-center justify-between">
                        <div>
                            <h3 className="text-lg font-medium text-gray-900">🔗 Configuración de Webhooks</h3>
                            <p className="text-sm text-gray-500 mt-1">
                                Estado de los webhooks configurados en Shopify
                            </p>
                        </div>
                        <div className="flex gap-2">
                            {!webhookConfigured && (
                                <Button
                                    type="button"
                                    onClick={handleVerifyAndShowCreateModal}
                                    disabled={verifyingWebhooks}
                                    loading={verifyingWebhooks}
                                    variant="primary"
                                    size="sm"
                                >
                                    {verifyingWebhooks ? 'Verificando...' : '➕ Crear Webhook'}
                                </Button>
                            )}
                            <Button
                                type="button"
                                onClick={handleGetWebhook}
                                disabled={loadingWebhook}
                                loading={loadingWebhook}
                                variant="outline"
                                size="sm"
                            >
                                {loadingWebhook ? 'Cargando...' : showWebhook ? '🔄 Actualizar URL' : '👁️ Ver URL'}
                            </Button>
                            <Button
                                type="button"
                                onClick={loadWebhooks}
                                disabled={loadingWebhooks}
                                loading={loadingWebhooks}
                                variant="outline"
                                size="sm"
                            >
                                🔄 Actualizar Lista
                            </Button>
                        </div>
                    </div>

                    {/* Estado del Webhook */}
                    <div className="p-4 rounded-lg border-2" style={{
                        backgroundColor: webhookConfigured ? '#f0fdf4' : '#fef2f2',
                        borderColor: webhookConfigured ? '#86efac' : '#fca5a5'
                    }}>
                        <div className="flex items-center justify-between">
                            <div className="flex items-center gap-3">
                                <span className="text-2xl">{webhookConfigured ? '✅' : '⚠️'}</span>
                                <div>
                                    <p className="font-medium text-gray-900">
                                        {webhookConfigured ? 'Webhooks Configurados Correctamente' : 'Webhooks No Configurados'}
                                    </p>
                                    <p className="text-sm text-gray-600 mt-1">
                                        {webhookConfigured 
                                            ? 'Los webhooks están activos y funcionando en Shopify' 
                                            : 'Los webhooks no se han configurado correctamente en Shopify'}
                                    </p>
                                </div>
                            </div>
                        </div>
                        {webhookUrl && (
                            <div className="mt-3 p-2 bg-white rounded border border-gray-200">
                                <span className="text-xs font-medium text-gray-600">URL del Webhook:</span>
                                <code className="block text-xs text-gray-800 break-all mt-1">{webhookUrl}</code>
                            </div>
                        )}
                    </div>

                    {/* Lista de Webhooks */}
                    <div>
                        <h4 className="text-md font-medium text-gray-900 mb-3">📋 Webhooks Registrados en Shopify</h4>
                        {loadingWebhooks ? (
                            <p className="text-sm text-gray-500">Cargando webhooks...</p>
                        ) : webhooks.length > 0 ? (
                            <div className="space-y-3">
                                {webhooks.map((webhook) => (
                                    <div key={webhook.id} className="p-4 bg-gray-50 rounded-lg border border-gray-200">
                                        <div className="flex items-start justify-between">
                                            <div className="flex-1">
                                                <div className="flex items-center gap-2 mb-2">
                                                    <span className="text-sm font-medium text-gray-900">ID: {webhook.id}</span>
                                                    <span className="text-xs px-2 py-1 bg-blue-100 text-blue-800 rounded">
                                                        {webhook.topic}
                                                    </span>
                                                </div>
                                                <div className="text-sm text-gray-600 mb-2">
                                                    <span className="font-medium">URL:</span>
                                                    <code className="ml-2 text-xs break-all">{webhook.address}</code>
                                                </div>
                                                <div className="flex gap-4 text-xs text-gray-500">
                                                    <span>Creado: {new Date(webhook.created_at).toLocaleDateString()}</span>
                                                    <span>Actualizado: {new Date(webhook.updated_at).toLocaleDateString()}</span>
                                                </div>
                                            </div>
                                            <Button
                                                type="button"
                                                onClick={() => handleDeleteWebhook(webhook.id)}
                                                variant="outline"
                                                size="sm"
                                                className="ml-4 text-red-600 hover:bg-red-50"
                                            >
                                                🗑️ Eliminar
                                            </Button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="p-4 bg-yellow-50 border border-yellow-200 rounded-lg">
                                <p className="text-sm text-yellow-800">
                                    ⚠️ No se encontraron webhooks registrados en Shopify. Los webhooks se crean automáticamente al crear la integración.
                                </p>
                            </div>
                        )}
                    </div>

                    {/* Información de URL del Webhook (opcional) */}
                    {showWebhook && webhookInfo && (
                        <div className="p-4 bg-blue-50 rounded-lg border border-blue-200">
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-sm font-medium text-blue-900">URL del Webhook para Configurar</span>
                                <Button
                                    type="button"
                                    onClick={handleCopyWebhook}
                                    variant="outline"
                                    size="sm"
                                >
                                    {copied ? '✅ Copiado!' : '📋 Copiar'}
                                </Button>
                            </div>
                            <code className="block p-3 bg-white border border-blue-300 rounded text-sm text-gray-800 break-all">
                                {webhookInfo.url}
                            </code>
                            {webhookInfo.events && webhookInfo.events.length > 0 && (
                                <div className="mt-3">
                                    <span className="text-xs font-medium text-blue-900">Eventos:</span>
                                    <div className="mt-1 flex flex-wrap gap-1">
                                        {webhookInfo.events.map((event, idx) => (
                                            <span key={idx} className="text-xs bg-blue-100 text-blue-800 px-2 py-0.5 rounded">
                                                {event}
                                            </span>
                                        ))}
                                    </div>
                                </div>
                            )}
                        </div>
                    )}
                </div>
            )}

            {/* Modal de Confirmación para Crear Webhook */}
            {showCreateWebhookModal && (
                <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
                    <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
                        <h3 className="text-lg font-semibold text-gray-900 mb-4">
                            ⚠️ Confirmar Creación de Webhook
                        </h3>
                        <div className="space-y-4">
                            <div>
                                <p className="text-sm text-gray-700 mb-2">
                                    Se creará un webhook en Shopify con la siguiente URL:
                                </p>
                                <div className="p-3 bg-gray-50 rounded border border-gray-200">
                                    <code className="text-sm text-gray-800 break-all">
                                        {webhookUrlToCreate}
                                    </code>
                                </div>
                            </div>
                            <div className="p-3 bg-blue-50 rounded border border-blue-200">
                                <p className="text-sm text-blue-800">
                                    <strong>Eventos que se registrarán:</strong>
                                </p>
                                <ul className="text-xs text-blue-700 mt-2 list-disc list-inside space-y-1">
                                    <li>orders/create</li>
                                    <li>orders/updated</li>
                                    <li>orders/paid</li>
                                    <li>orders/cancelled</li>
                                    <li>orders/fulfilled</li>
                                    <li>orders/partially_fulfilled</li>
                                </ul>
                            </div>
                            {error && (
                                <Alert type="error" onClose={() => setError(null)}>
                                    {error}
                                </Alert>
                            )}
                        </div>
                        <div className="flex gap-3 mt-6 justify-end">
                            <Button
                                type="button"
                                onClick={() => {
                                    setShowCreateWebhookModal(false);
                                    setWebhookUrlToCreate('');
                                    setError(null);
                                }}
                                variant="outline"
                                disabled={creatingWebhook}
                            >
                                Cancelar
                            </Button>
                            <Button
                                type="button"
                                onClick={handleCreateWebhook}
                                loading={creatingWebhook}
                                disabled={creatingWebhook}
                                variant="primary"
                            >
                                {creatingWebhook ? 'Creando...' : '✅ Confirmar y Crear'}
                            </Button>
                        </div>
                    </div>
                </div>
            )}

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
