'use client';

import { useState } from 'react';
import { Input, Button, Alert } from '@/shared/ui';
import { TokenStorage } from '@/shared/utils';
import { BeakerIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import ShopifyWebhookManager from './ShopifyWebhookManager';
import { getActionError } from '@/shared/utils/action-result';

interface ShopifyOAuthFormProps {
    onCancel?: () => void;
    onSubmit?: (data: any) => void;
    onTestConnection?: (config: any, credentials: any) => Promise<boolean>;
    onGetWebhook?: () => Promise<any>;
    initialData?: {
        name?: string;
        code?: string;
        store_id?: string;
        config?: any;
        credentials?: any;
        business_id?: number | null;
        is_testing?: boolean;
        base_url_test?: string;
    };
    isEdit?: boolean;
    integrationId?: number;
}

export default function ShopifyOAuthForm({
    onCancel,
    onSubmit,
    onTestConnection,
    onGetWebhook,
    initialData,
    isEdit,
    integrationId,
}: ShopifyOAuthFormProps) {
    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        shop_domain: initialData?.store_id || '',
        client_id: initialData?.credentials?.client_id || '',
        client_secret: initialData?.credentials?.client_secret || '',
    });

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isTesting, setIsTesting] = useState(initialData?.is_testing || false);
    const [showSecrets, setShowSecrets] = useState(false);

    const accessToken = initialData?.credentials?.access_token || '';

    const handleConnectShopify = async (e: React.FormEvent) => {
        e.preventDefault();

        // In edit mode, call onSubmit with updated data including is_testing
        if (isEdit && onSubmit) {
            const credentials: any = {};
            if (formData.client_id) credentials.client_id = formData.client_id;
            if (formData.client_secret) credentials.client_secret = formData.client_secret;
            if (accessToken) credentials.access_token = accessToken;

            onSubmit({
                name: formData.name,
                store_id: formData.shop_domain,
                config: initialData?.config || {},
                credentials: Object.keys(credentials).length > 0 ? credentials : undefined,
                is_testing: isTesting,
            });
            return;
        }

        if (!formData.name || !formData.shop_domain || !formData.client_id || !formData.client_secret) {
            setError('Por favor completa todos los campos');
            return;
        }

        setLoading(true);
        setError(null);

        try {
            const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
            const response = await fetch(`${apiBaseUrl}/integrations/shopify/connect/custom`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${TokenStorage.getSessionToken()}`
                },
                credentials: 'include',
                body: JSON.stringify({
                    shop_domain: formData.shop_domain,
                    integration_name: formData.name,
                    client_id: formData.client_id,
                    client_secret: formData.client_secret
                })
            });

            const data = await response.json();

            if (!response.ok || !data.success) {
                throw new Error(data.error || data.message || 'Error al iniciar OAuth');
            }

            if (data.authorization_url) {
                window.location.href = data.authorization_url;
            } else {
                throw new Error('No se recibió URL de autorización');
            }
        } catch (err: any) {
            console.error('Error al conectar con Shopify:', err);
            setError(getActionError(err, 'Error al conectar con Shopify'));
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleConnectShopify} className="space-y-6 w-full">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="p-6 rounded-lg border border-gray-200 bg-white">
                <div className="space-y-4">
                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
                        <div className="flex items-start gap-3">
                            <span className="text-2xl">&#8505;&#65039;</span>
                            <div>
                                <p className="text-sm font-medium text-blue-900 mb-1">
                                    Conexión Shopify Custom App
                                </p>
                                <p className="text-xs text-blue-700">
                                    {isEdit
                                        ? 'Datos de tu Custom App de Shopify. Puedes modificar las credenciales si es necesario.'
                                        : 'Ingresa las credenciales de tu Custom App creada en el Shopify Partner Dashboard. Serás redirigido a Shopify para autorizar.'}
                                </p>
                            </div>
                        </div>
                    </div>

                    {/* Nombre de la Integración */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
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
                        <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                            Un nombre descriptivo para identificar esta integración
                        </p>
                    </div>

                    {/* Store Domain */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                            Dominio de la Tienda *
                        </label>
                        <Input
                            type="text"
                            required
                            placeholder="tienda.myshopify.com"
                            value={formData.shop_domain}
                            onChange={(e) => setFormData({ ...formData, shop_domain: e.target.value })}
                            className="w-full"
                        />
                    </div>

                    {/* Access Token - Solo en modo edición, read-only */}
                    {isEdit && accessToken && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                                Access Token (OAuth)
                            </label>
                            <div className="relative">
                                <Input
                                    type={showSecrets ? 'text' : 'password'}
                                    value={accessToken}
                                    readOnly
                                    className="w-full bg-gray-50 pr-10"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowSecrets(!showSecrets)}
                                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:text-gray-300"
                                >
                                    {showSecrets
                                        ? <EyeSlashIcon className="w-5 h-5" />
                                        : <EyeIcon className="w-5 h-5" />}
                                </button>
                            </div>
                            <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                                Obtenido automáticamente durante el flujo OAuth (solo lectura)
                            </p>
                        </div>
                    )}

                    {/* Client ID */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                            Client ID (API Key) {!isEdit && '*'}
                        </label>
                        <Input
                            type="text"
                            required={!isEdit}
                            placeholder={isEdit ? '' : 'Pegar Client ID aquí'}
                            value={formData.client_id}
                            onChange={(e) => setFormData({ ...formData, client_id: e.target.value })}
                            autoComplete="off"
                            className="w-full"
                        />
                    </div>

                    {/* Client Secret */}
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                            Client Secret {!isEdit && '*'}
                        </label>
                        <div className="relative">
                            <Input
                                type={showSecrets ? 'text' : 'password'}
                                required={!isEdit}
                                placeholder={isEdit ? '' : 'Pegar Client Secret aquí'}
                                value={formData.client_secret}
                                onChange={(e) => setFormData({ ...formData, client_secret: e.target.value })}
                                autoComplete="new-password"
                                className="w-full pr-10"
                            />
                            <button
                                type="button"
                                onClick={() => setShowSecrets(!showSecrets)}
                                className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 dark:text-gray-300"
                            >
                                {showSecrets
                                    ? <EyeSlashIcon className="w-5 h-5" />
                                    : <EyeIcon className="w-5 h-5" />}
                            </button>
                        </div>
                    </div>
                </div>
            </div>

            {/* Webhook Management - Solo en modo edición */}
            {isEdit && integrationId && (
                <div className="p-6 rounded-lg border border-gray-200 bg-white">
                    <h3 className="text-base font-semibold text-gray-900 dark:text-white mb-4 pb-3 border-b border-gray-200">
                        Webhooks
                    </h3>
                    <ShopifyWebhookManager integrationId={integrationId} />
                </div>
            )}

            {/* Modo de Pruebas - Solo en modo edición */}
            {isEdit && (
                <div className="bg-orange-50 rounded-xl p-6 space-y-4 border border-orange-200">
                    <div className="flex items-center gap-2 mb-2">
                        <BeakerIcon className="w-5 h-5 text-orange-600" />
                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                            Modo de Pruebas
                        </h3>
                    </div>
                    <div className="flex items-center justify-between p-3 bg-white rounded-lg border border-orange-200">
                        <div className="flex-1">
                            <p className="text-sm font-medium text-gray-800 dark:text-gray-100">Activar modo testing</p>
                            <p className="text-xs text-gray-500 dark:text-gray-400 mt-0.5">
                                Las peticiones a Shopify se redirigirán a la URL de pruebas configurada.
                            </p>
                        </div>
                        <button
                            type="button"
                            onClick={() => setIsTesting(!isTesting)}
                            className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none ml-4 flex-shrink-0 ${isTesting ? 'bg-orange-500' : 'bg-gray-200'}`}
                        >
                            <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${isTesting ? 'translate-x-6' : 'translate-x-1'}`} />
                        </button>
                    </div>
                    {isTesting && (
                        <Alert type="warning">
                            Modo de pruebas activado. Las peticiones de sincronización y webhooks se enviarán al servidor de pruebas en lugar de Shopify.
                            {initialData?.base_url_test && (
                                <p className="mt-2 text-xs font-mono text-orange-800 break-all">
                                    URL sandbox: {initialData.base_url_test}
                                </p>
                            )}
                        </Alert>
                    )}
                </div>
            )}

            {/* Action Buttons */}
            <div className="flex flex-row justify-end gap-3 pt-4 border-t">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        variant="outline"
                        disabled={loading}
                    >
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    disabled={loading || !formData.name || !formData.shop_domain || (!isEdit && (!formData.client_id || !formData.client_secret))}
                    loading={loading}
                    variant="primary"
                >
                    {isEdit
                        ? (loading ? 'Guardando...' : 'Guardar Cambios')
                        : (loading ? 'Conectando...' : 'Conectar con Shopify')}
                </Button>
            </div>
        </form>
    );
}
