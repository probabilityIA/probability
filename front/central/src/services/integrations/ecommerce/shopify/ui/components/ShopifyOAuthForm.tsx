'use client';

import { useState } from 'react';
import { Input, Button, Alert, SecretInput } from '@/shared/ui';
import { TokenStorage } from '@/shared/utils';
import {
    BeakerIcon,
    TruckIcon,
    KeyIcon,
    BoltIcon,
    InformationCircleIcon,
} from '@heroicons/react/24/outline';
import ShopifyWebhookManager from './ShopifyWebhookManager';
import { getActionError } from '@/shared/utils/action-result';
import { useToast } from '@/shared/providers/toast-provider';
import {
    enableShopifyCarrierServiceAction,
    disableShopifyCarrierServiceAction,
} from '@/services/integrations/core/infra/actions';

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

interface ToggleRowProps {
    icon: React.ReactNode;
    title: string;
    subtitle: string;
    checked: boolean;
    onToggle: () => void;
    disabled?: boolean;
    activeColor?: 'indigo' | 'orange';
}

function ToggleRow({ icon, title, subtitle, checked, onToggle, disabled, activeColor = 'indigo' }: ToggleRowProps) {
    const onColor = activeColor === 'orange' ? 'bg-orange-500' : 'bg-indigo-500';
    return (
        <div className="flex items-center justify-between gap-3 px-4 py-2.5">
            <div className="flex items-center gap-2.5 min-w-0">
                <span className="shrink-0">{icon}</span>
                <div className="min-w-0">
                    <p className="text-sm font-medium text-gray-800 dark:text-gray-100 leading-tight">{title}</p>
                    <p className="text-xs text-gray-500 dark:text-gray-400 leading-tight mt-0.5">{subtitle}</p>
                </div>
            </div>
            <button
                type="button"
                onClick={onToggle}
                disabled={disabled}
                className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none shrink-0 disabled:opacity-50 ${checked ? onColor : 'bg-gray-300 dark:bg-gray-600'}`}
            >
                <span className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow-sm transition-transform ${checked ? 'translate-x-5' : 'translate-x-0.5'}`} />
            </button>
        </div>
    );
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
    const accessToken = initialData?.credentials?.access_token || '';

    const { showToast } = useToast();
    const [carrierEnabled, setCarrierEnabled] = useState<boolean>(
        initialData?.config?.carrier_calculated_shipping_enabled === true
    );
    const [carrierLoading, setCarrierLoading] = useState(false);

    const handleToggleCarrierService = async () => {
        if (!integrationId || carrierLoading) return;

        setCarrierLoading(true);
        const next = !carrierEnabled;
        const result: any = next
            ? await enableShopifyCarrierServiceAction(integrationId)
            : await disableShopifyCarrierServiceAction(integrationId);

        if (!result || result.success === false) {
            showToast(result?.message || 'No se pudo actualizar la cotización en checkout', 'error');
        } else {
            setCarrierEnabled(next);
            showToast(
                result.message || (next
                    ? 'Cotización en checkout activada'
                    : 'Cotización en checkout desactivada'),
                'success'
            );
        }
        setCarrierLoading(false);
    };

    const handleConnectShopify = async (e: React.FormEvent) => {
        e.preventDefault();

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

    const sectionHeader = (icon: React.ReactNode, title: string) => (
        <div className="flex items-center gap-2 px-5 py-3 border-b border-gray-100 dark:border-gray-700 bg-gray-50/70 dark:bg-gray-800/60">
            {icon}
            <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100">{title}</h3>
        </div>
    );

    return (
        <form onSubmit={handleConnectShopify} className="space-y-5 w-full">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="flex items-start gap-3 rounded-xl border border-emerald-200 dark:border-emerald-900/50 bg-gradient-to-r from-emerald-50 to-teal-50 dark:from-emerald-950/30 dark:to-teal-950/30 px-4 py-3">
                <InformationCircleIcon className="w-5 h-5 text-emerald-600 dark:text-emerald-400 shrink-0 mt-0.5" />
                <div>
                    <p className="text-sm font-semibold text-emerald-900 dark:text-emerald-200">Conexión Shopify Custom App</p>
                    <p className="text-xs text-emerald-700 dark:text-emerald-300/80 mt-0.5">
                        {isEdit
                            ? 'Datos de tu Custom App de Shopify. Puedes modificar las credenciales si es necesario.'
                            : 'Ingresa las credenciales de tu Custom App creada en el Shopify Partner Dashboard. Serás redirigido a Shopify para autorizar.'}
                    </p>
                </div>
            </div>

            <div className={isEdit ? 'grid grid-cols-1 lg:grid-cols-2 gap-5 items-start' : 'max-w-2xl mx-auto'}>
                <section className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm overflow-hidden">
                    {sectionHeader(<KeyIcon className="w-4 h-4 text-gray-500 dark:text-gray-400" />, 'Datos de conexión')}
                    <div className="p-5 space-y-4">
                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1.5">
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
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1.5">
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
                        </div>

                        {isEdit && accessToken && (
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1.5">
                                    Access Token (OAuth)
                                </label>
                                <SecretInput
                                    value={accessToken}
                                    readOnly
                                    className="w-full bg-gray-50 dark:bg-gray-700"
                                />
                                <p className="text-[11px] text-gray-400 dark:text-gray-500 mt-1">
                                    Obtenido automáticamente durante el flujo OAuth (solo lectura)
                                </p>
                            </div>
                        )}

                        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1.5">
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
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1.5">
                                    Client Secret {!isEdit && '*'}
                                </label>
                                <SecretInput
                                    value={formData.client_secret}
                                    onChange={(e) => setFormData({ ...formData, client_secret: e.target.value })}
                                    placeholder={isEdit ? '' : 'Pegar Client Secret aquí'}
                                    required={!isEdit}
                                    className="w-full"
                                />
                            </div>
                        </div>
                    </div>
                </section>

                {isEdit && integrationId && (
                    <div className="space-y-5">
                        <section className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm overflow-hidden">
                            {sectionHeader(<BoltIcon className="w-4 h-4 text-amber-500" />, 'Webhooks')}
                            <div className="p-5">
                                <ShopifyWebhookManager integrationId={integrationId} />
                            </div>
                        </section>

                        <section className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm overflow-hidden">
                            {sectionHeader(<TruckIcon className="w-4 h-4 text-indigo-500" />, 'Configuración de envío')}
                            <div className="divide-y divide-gray-100 dark:divide-gray-700">
                                <ToggleRow
                                    icon={<TruckIcon className="w-4 h-4 text-indigo-500" />}
                                    title="Cotización en checkout"
                                    subtitle="Tarifas en tiempo real con varias transportadoras al pagar"
                                    checked={carrierEnabled}
                                    onToggle={handleToggleCarrierService}
                                    disabled={carrierLoading}
                                    activeColor="indigo"
                                />
                                <div>
                                    <ToggleRow
                                        icon={<BeakerIcon className="w-4 h-4 text-orange-500" />}
                                        title="Modo de pruebas"
                                        subtitle="Redirige las peticiones a la URL de pruebas"
                                        checked={isTesting}
                                        onToggle={() => setIsTesting(!isTesting)}
                                        activeColor="orange"
                                    />
                                    {isTesting && initialData?.base_url_test && (
                                        <p className="px-4 pb-2.5 -mt-1 text-[11px] font-mono text-orange-700 dark:text-orange-400 break-all">
                                            Sandbox: {initialData.base_url_test}
                                        </p>
                                    )}
                                </div>
                            </div>
                        </section>
                    </div>
                )}
            </div>

            <div className="flex flex-row justify-end gap-3 pt-4 border-t border-gray-200 dark:border-gray-700">
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
