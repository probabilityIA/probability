'use client';

import { useState, useEffect } from 'react';
import { Alert, SecretInput } from '@/shared/ui';
import { TokenStorage } from '@/shared/utils';
import {
    BeakerIcon,
    TruckIcon,
    KeyIcon,
    ShoppingBagIcon,
    InformationCircleIcon,
    ChevronLeftIcon,
    ChevronRightIcon,
    ClipboardDocumentIcon,
    CheckIcon,
} from '@heroicons/react/24/outline';
import ShopifyWebhookManager from './ShopifyWebhookManager';
import { ShopifyInventorySection, ShopifyInventoryConfig } from './ShopifyInventorySection';
import { ShopifyLocationMappingSection, ShopifyLocationMapping } from './ShopifyLocationMappingSection';
import { ShopifyInventorySyncModal } from './ShopifyInventorySyncModal';
import { ShopifyProductSyncModal } from './ShopifyProductSyncModal';
import { getActionError } from '@/shared/utils/action-result';
import { useToast } from '@/shared/providers/toast-provider';
import {
    enableShopifyCarrierServiceAction,
    disableShopifyCarrierServiceAction,
    getActiveIntegrationTypesAction,
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

const GREEN = 'var(--color-primary)';
const GREEN_DARK = 'color-mix(in srgb, var(--color-primary) 85%, black)';
const GREEN_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
const GREEN_BORDER = 'color-mix(in srgb, var(--color-primary) 25%, white)';
const CARD_BG = '#fafafd';
const CARD_BORDER = '#eceaf3';
const INPUT_BORDER = '#e9e9f0';

const SHOPIFY_ASSETS_BASE = 'https://probability-media-assets.s3.us-east-1.amazonaws.com/docs/shopify-integration';
const SHOPIFY_REQUIRED_SCOPES = 'read_customers,write_customers,read_fulfillments,write_fulfillments,read_orders,write_orders,read_products,write_products';
const SHOPIFY_REDIRECT_URL = 'https://www.probabilityia.com.co/api/v1/shopify/callback';
const SHOPIFY_APP_URL = 'https://www.probabilityia.com.co';

const SHOPIFY_CONNECTION_STEPS = [
    {
        step: 1,
        title: 'Entra a Integraciones en tu panel de Probability',
        description: 'En el menu lateral, abre "Integraciones" y selecciona la pestana "Plataforma".',
        image: `${SHOPIFY_ASSETS_BASE}/01-panel-probability-integraciones-plataforma.png`,
    },
    {
        step: 2,
        title: 'Cambia a la pestana "E-commerce"',
        description: 'Ahi veras las integraciones de tienda online que ya tengas conectadas.',
        image: `${SHOPIFY_ASSETS_BASE}/02-panel-probability-integraciones-ecommerce.png`,
    },
    {
        step: 3,
        title: 'Crea una nueva integracion',
        description: 'Clic en "Crear Integracion" y selecciona la categoria "E-commerce".',
        image: `${SHOPIFY_ASSETS_BASE}/03-modal-nueva-integracion-seleccionar-categoria.png`,
    },
    {
        step: 4,
        title: 'Selecciona Shopify como proveedor',
        description: 'Dentro de E-commerce, elige la tarjeta "Shopify".',
        image: `${SHOPIFY_ASSETS_BASE}/04-modal-nueva-integracion-seleccionar-shopify.png`,
    },
    {
        step: 5,
        title: 'Entra a Shopify Partners',
        description: 'Inicia sesion en partners.shopify.com y selecciona la organizacion/tienda de desarrollo donde crearas la Custom App.',
        image: `${SHOPIFY_ASSETS_BASE}/05-shopify-partners-selector-tiendas.png`,
    },
    {
        step: 6,
        title: 'Crea una app nueva',
        description: 'En el Dev Dashboard, ve a "Apps" y clic en "Crear app".',
        image: `${SHOPIFY_ASSETS_BASE}/06-shopify-dev-dashboard-apps-crear-app.png`,
    },
    {
        step: 7,
        title: 'Nombra tu app',
        description: 'Elige "Empezar desde Dev Dashboard", escribe un nombre para la app (ej: el nombre del negocio) y clic en "Crear".',
        image: `${SHOPIFY_ASSETS_BASE}/07-shopify-dev-dashboard-crear-app-nombre.png`,
    },
    {
        step: 8,
        title: 'Configura los alcances (scopes)',
        description: 'En "Versiones" > "Acceso a la API", pega los Alcances necesarios en el campo "Alcances". Estos son obligatorios para que la app pueda leer/escribir clientes, pedidos, productos y fulfillments.',
        image: `${SHOPIFY_ASSETS_BASE}/08-shopify-app-scopes-configurados-urls-vacias-INCOMPLETO.png`,
        copyItems: [
            { label: 'Alcances requeridos', text: SHOPIFY_REQUIRED_SCOPES },
        ],
    },
    {
        step: 9,
        title: 'Configura las URLs (paso critico)',
        description: 'Justo debajo, llena "URLs de redireccionamiento" y "URL de la app" con los valores de abajo. Si dejas estos campos vacios, Shopify rechazara la conexion con un error de "hosts must match". Clic en "Publicar".',
        image: `${SHOPIFY_ASSETS_BASE}/14-shopify-app-urls-redireccionamiento-correctas-FIX.png`,
        copyItems: [
            { label: 'URL de redireccionamiento', text: SHOPIFY_REDIRECT_URL },
            { label: 'URL de la app', text: SHOPIFY_APP_URL },
        ],
    },
    {
        step: 10,
        title: 'Copia el Client ID y Client Secret',
        description: 'Ve a "Configuracion" y copia el "ID de cliente" y el "Secreto" (clic en el icono del ojo para revelarlo). Los necesitaras en el paso 13.',
        image: `${SHOPIFY_ASSETS_BASE}/10-shopify-app-credenciales-client-id-secret.png`,
    },
    {
        step: 11,
        title: 'Publica la primera version',
        description: 'Shopify crea automaticamente una version de la app. La veras activa en la seccion "Versiones".',
        image: `${SHOPIFY_ASSETS_BASE}/09-shopify-app-versiones-publicada.png`,
    },
    {
        step: 12,
        title: 'Vuelve al formulario de Probability',
        description: 'Pega el Client ID y Client Secret que copiaste en el paso 10, junto con el dominio de la tienda, y clic en "Conectar con Shopify".',
        image: `${SHOPIFY_ASSETS_BASE}/11-probability-form-conexion-custom-app.png`,
    },
    {
        step: 13,
        title: 'Autoriza el acceso en Shopify',
        description: 'Shopify te mostrara los permisos que la app necesita. Clic en "Instalar" para autorizar.',
        image: `${SHOPIFY_ASSETS_BASE}/12-shopify-instalar-app-permisos.png`,
    },
    {
        step: 14,
        title: 'Listo, la integracion queda activa',
        description: 'Vuelve al panel de Probability: la integracion Shopify aparecera con estado "Activo".',
        image: `${SHOPIFY_ASSETS_BASE}/13-panel-probability-integracion-shopify-activa.png`,
    },
];

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1 flex items-start gap-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

interface ToggleRowProps {
    icon: React.ReactNode;
    title: string;
    subtitle: string;
    checked: boolean;
    onToggle: () => void;
    disabled?: boolean;
    activeColor?: 'indigo' | 'orange';
}

function ToggleRow({ icon, title, subtitle, checked, onToggle, disabled }: ToggleRowProps) {
    return (
        <div className="flex items-center justify-between gap-3 px-3 py-2.5">
            <div className="flex items-center gap-2.5 min-w-0">
                <span
                    className="flex h-8 w-8 items-center justify-center rounded-lg shrink-0"
                    style={{ backgroundColor: 'color-mix(in srgb, var(--color-primary) 10%, white)' }}
                >
                    {icon}
                </span>
                <div className="min-w-0">
                    <p className="text-[13px] font-semibold text-gray-800 dark:text-gray-100 leading-tight">{title}</p>
                    <p className="text-[11px] text-gray-500 dark:text-gray-400 leading-tight mt-0.5">{subtitle}</p>
                </div>
            </div>
            <button
                type="button"
                role="switch"
                aria-checked={checked}
                onClick={onToggle}
                disabled={disabled}
                className="relative inline-flex h-7 w-12 items-center rounded-full transition-colors focus:outline-none shrink-0 disabled:opacity-50"
                style={{ backgroundColor: checked ? 'var(--color-primary)' : '#e5e7eb' }}
            >
                <span className={`inline-block h-5 w-5 transform rounded-full bg-white shadow-md transition-transform ${checked ? 'translate-x-6' : 'translate-x-1'}`} />
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
    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);
    const [guideStepIndex, setGuideStepIndex] = useState(0);
    const [copiedGuideLabel, setCopiedGuideLabel] = useState<string | null>(null);

    const handleCopyGuideText = async (label: string, text: string) => {
        await navigator.clipboard.writeText(text);
        setCopiedGuideLabel(label);
        setTimeout(() => setCopiedGuideLabel((prev) => (prev === label ? null : prev)), 2000);
    };

    const { showToast } = useToast();
    const [carrierEnabled, setCarrierEnabled] = useState<boolean>(
        initialData?.config?.carrier_calculated_shipping_enabled === true
    );
    const [carrierLoading, setCarrierLoading] = useState(false);
    const [inventorySyncOpen, setInventorySyncOpen] = useState(false);
    const [productSyncOpen, setProductSyncOpen] = useState(false);
    const [inventorySync, setInventorySync] = useState<ShopifyInventoryConfig>(() => {
        const c: any = initialData?.config || {};
        return {
            enabled: !!c.inventory_sync_enabled,
            mode: c.inventory_warehouse_mode === 'single' ? 'single' : 'sum',
            single_warehouse_id: Number(c.inventory_single_warehouse_id) || 0,
            warehouse_ids: Array.isArray(c.inventory_warehouse_ids) ? c.inventory_warehouse_ids.map(Number) : [],
        };
    });
    const [defaultLocationId, setDefaultLocationId] = useState<string>(() => {
        const c: any = initialData?.config || {};
        return c.shopify_default_location_id ? String(c.shopify_default_location_id) : '';
    });
    const [locationMappings, setLocationMappings] = useState<ShopifyLocationMapping[]>(() => {
        const c: any = initialData?.config || {};
        return Array.isArray(c.shopify_location_mappings)
            ? c.shopify_location_mappings.map((m: any) => ({
                internal_warehouse_id: Number(m.internal_warehouse_id) || 0,
                shopify_location_id: String(m.shopify_location_id ?? ''),
            }))
            : [];
    });

    useEffect(() => {
        let cancelled = false;
        getActiveIntegrationTypesAction()
            .then((res: any) => {
                if (cancelled) return;
                const types = res?.data || [];
                const shopify = types.find((t: any) => t.id === 1 || /shopify/i.test(t.code || ''));
                if (shopify?.image_url) setLogoUrl(shopify.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

    const handleToggleCarrierService = async () => {
        if (!integrationId || carrierLoading) return;

        setCarrierLoading(true);
        const next = !carrierEnabled;
        const result: any = next
            ? await enableShopifyCarrierServiceAction(integrationId)
            : await disableShopifyCarrierServiceAction(integrationId);

        if (!result || result.success === false) {
            showToast(result?.message || 'No se pudo actualizar la cotizacion en checkout', 'error');
        } else {
            setCarrierEnabled(next);
            showToast(
                result.message || (next
                    ? 'Cotizacion en checkout activada'
                    : 'Cotizacion en checkout desactivada'),
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

            const cleanMappings = locationMappings
                .filter((m) => m.internal_warehouse_id > 0 && m.shopify_location_id.trim() !== '')
                .map((m) => ({
                    internal_warehouse_id: m.internal_warehouse_id,
                    shopify_location_id: m.shopify_location_id.trim(),
                }));

            const mergedConfig = {
                ...(initialData?.config || {}),
                inventory_sync_enabled: inventorySync.enabled,
                inventory_warehouse_mode: inventorySync.mode,
                inventory_single_warehouse_id: inventorySync.single_warehouse_id,
                inventory_warehouse_ids: inventorySync.warehouse_ids,
                shopify_default_location_id: defaultLocationId.trim() ? Number(defaultLocationId.trim()) || 0 : 0,
                shopify_location_mappings: cleanMappings,
            };

            onSubmit({
                name: formData.name,
                store_id: formData.shop_domain,
                config: mergedConfig,
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
                throw new Error('No se recibio URL de autorizacion');
            }
        } catch (err: any) {
            console.error('Error al conectar con Shopify:', err);
            setError(getActionError(err, 'Error al conectar con Shopify'));
            setLoading(false);
        }
    };

    const connected = isEdit && !!accessToken;

    return (
        <form onSubmit={handleConnectShopify} className="space-y-3 w-full">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div
                className="flex flex-col gap-3 rounded-xl p-4 sm:flex-row sm:items-center sm:justify-between dark:bg-gray-800/60"
                style={{ backgroundColor: GREEN_SOFT, border: `1px solid ${GREEN_BORDER}` }}
            >
                <div className="flex items-center gap-3">
                    <span
                        className="flex h-11 w-11 items-center justify-center rounded-xl overflow-hidden shrink-0 bg-white dark:bg-gray-900"
                        style={{ border: `1px solid ${GREEN_BORDER}`, ...(logoUrl && !logoFailed ? {} : { backgroundColor: GREEN }) }}
                    >
                        {logoUrl && !logoFailed ? (
                            <img
                                src={logoUrl}
                                alt="Shopify"
                                className="h-8 w-8 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <ShoppingBagIcon className="h-6 w-6 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-base font-bold text-gray-900 dark:text-white leading-tight">Conexion Shopify Custom App</h2>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-0.5">
                            {isEdit
                                ? 'Datos de tu Custom App de Shopify. Puedes modificar las credenciales si es necesario.'
                                : 'Ingresa las credenciales de tu Custom App creada en el Shopify Partner Dashboard. Seras redirigido a Shopify para autorizar.'}
                        </p>
                    </div>
                </div>
                <span
                    className="inline-flex items-center gap-2 self-start rounded-full px-3 py-1 text-[11px] font-semibold shrink-0 bg-white dark:bg-gray-900"
                    style={connected
                        ? { border: `1px solid ${GREEN_BORDER}`, color: GREEN_DARK }
                        : { border: '1px solid #e5e7eb', color: '#6b7280' }}
                >
                    <span className="h-2 w-2 rounded-full" style={{ backgroundColor: connected ? GREEN : '#9ca3af' }} />
                    {connected ? 'Conectado' : 'Sin conectar'}
                </span>
            </div>

            <div
                className="rounded-xl p-4 dark:bg-gray-800/60"
                style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex items-center gap-2 mb-3">
                    <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                        <KeyIcon style={{ color: GREEN, width: 16, height: 16 }} />
                    </span>
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Datos de conexion</h3>
                </div>

                <div className="space-y-3">
                    <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                        <div>
                            <label className={fieldLabel}>
                                Nombre de la Integracion <span style={{ color: GREEN }}>*</span>
                            </label>
                            <input
                                type="text"
                                required
                                placeholder="Ej: Tienda Principal"
                                value={formData.name}
                                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                className={inputCls}
                                style={{ borderColor: INPUT_BORDER }}
                            />
                        </div>
                        <div>
                            <label className={fieldLabel}>
                                Dominio de la tienda <span style={{ color: GREEN }}>*</span>
                            </label>
                            <input
                                type="text"
                                required
                                placeholder="tienda.myshopify.com"
                                value={formData.shop_domain}
                                onChange={(e) => setFormData({ ...formData, shop_domain: e.target.value })}
                                className={inputCls}
                                style={{ borderColor: INPUT_BORDER }}
                            />
                        </div>
                    </div>

                    {isEdit && accessToken && (
                        <div>
                            <label className={`${fieldLabel} flex items-center gap-2`}>
                                Access Token (OAuth)
                                <span
                                    className="rounded px-1.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide text-gray-500 dark:text-gray-400 bg-gray-100 dark:bg-gray-700"
                                    style={{ border: `1px solid ${INPUT_BORDER}` }}
                                >
                                    solo lectura
                                </span>
                            </label>
                            <SecretInput
                                value={accessToken}
                                readOnly
                                className="w-full bg-gray-50 dark:bg-gray-700 font-mono text-sm rounded-lg"
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Obtenido automaticamente durante el flujo OAuth</span>
                            </p>
                        </div>
                    )}

                    <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                        <div>
                            <label className={fieldLabel}>
                                Client ID (API Key) {!isEdit && <span style={{ color: GREEN }}>*</span>}
                            </label>
                            <input
                                type="text"
                                required={!isEdit}
                                placeholder={isEdit ? '' : 'Pegar Client ID aqui'}
                                value={formData.client_id}
                                onChange={(e) => setFormData({ ...formData, client_id: e.target.value })}
                                autoComplete="off"
                                className={`${inputCls} font-mono`}
                                style={{ borderColor: INPUT_BORDER }}
                            />
                        </div>
                        <div>
                            <label className={fieldLabel}>
                                Client Secret {!isEdit && <span style={{ color: GREEN }}>*</span>}
                            </label>
                            <SecretInput
                                value={formData.client_secret}
                                onChange={(e) => setFormData({ ...formData, client_secret: e.target.value })}
                                placeholder={isEdit ? '' : 'Pegar Client Secret aqui'}
                                required={!isEdit}
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                        </div>
                    </div>
                </div>
            </div>

            {!isEdit && (
                <div
                    className="rounded-xl p-4 dark:bg-gray-800/60"
                    style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
                >
                    <div className="flex items-center gap-2 mb-1">
                        <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                            <InformationCircleIcon style={{ color: GREEN, width: 16, height: 16 }} />
                        </span>
                        <h3 className="text-sm font-bold text-gray-900 dark:text-white">Sigue los pasos para hacer la conexion a Shopify</h3>
                    </div>
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-4 ml-9">
                        Crea la Custom App en el Shopify Partner Dashboard y copia aqui el Client ID y Client Secret que genera.
                    </p>

                    {(() => {
                        const totalSteps = SHOPIFY_CONNECTION_STEPS.length;
                        const currentStep = SHOPIFY_CONNECTION_STEPS[guideStepIndex];
                        const goPrev = () => setGuideStepIndex((i) => Math.max(0, i - 1));
                        const goNext = () => setGuideStepIndex((i) => Math.min(totalSteps - 1, i + 1));

                        return (
                            <div>
                                <div className="flex items-center justify-between mb-2">
                                    <span
                                        className="flex h-6 w-6 shrink-0 items-center justify-center rounded-full text-[11px] font-bold text-white"
                                        style={{ backgroundColor: GREEN }}
                                    >
                                        {currentStep.step}
                                    </span>
                                    <span className="text-[11px] font-semibold text-gray-400 dark:text-gray-500">
                                        Paso {guideStepIndex + 1} de {totalSteps}
                                    </span>
                                </div>

                                <p className="text-[13px] font-semibold text-gray-800 dark:text-gray-100">{currentStep.title}</p>
                                <p className="text-[12px] text-gray-500 dark:text-gray-400 mt-0.5">{currentStep.description}</p>

                                <div className="relative mt-3 flex items-center gap-2">
                                    <button
                                        type="button"
                                        onClick={goPrev}
                                        disabled={guideStepIndex === 0}
                                        aria-label="Paso anterior"
                                        className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-white dark:bg-gray-800 shadow-sm disabled:opacity-30 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-700"
                                        style={{ border: `1px solid ${INPUT_BORDER}` }}
                                    >
                                        <ChevronLeftIcon className="h-5 w-5 text-gray-600 dark:text-gray-300" />
                                    </button>

                                    <div className="min-w-0 flex-1 flex justify-center">
                                        <img
                                            src={currentStep.image}
                                            alt={currentStep.title}
                                            className="max-h-[28rem] w-auto max-w-full rounded-lg object-contain"
                                            style={{ border: `1px solid ${INPUT_BORDER}` }}
                                            loading="lazy"
                                        />
                                    </div>

                                    <button
                                        type="button"
                                        onClick={goNext}
                                        disabled={guideStepIndex === totalSteps - 1}
                                        aria-label="Paso siguiente"
                                        className="flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-white dark:bg-gray-800 shadow-sm disabled:opacity-30 disabled:cursor-not-allowed hover:bg-gray-50 dark:hover:bg-gray-700"
                                        style={{ border: `1px solid ${INPUT_BORDER}` }}
                                    >
                                        <ChevronRightIcon className="h-5 w-5 text-gray-600 dark:text-gray-300" />
                                    </button>
                                </div>

                                {currentStep.copyItems && currentStep.copyItems.length > 0 && (
                                    <div className="mt-3 space-y-3">
                                        {currentStep.copyItems.map((item) => (
                                            <div key={item.label}>
                                                <div className="flex items-center justify-between mb-1">
                                                    <span className="text-[11px] font-semibold text-gray-500 dark:text-gray-400">
                                                        {item.label}
                                                    </span>
                                                    <button
                                                        type="button"
                                                        onClick={() => handleCopyGuideText(item.label, item.text)}
                                                        className="inline-flex items-center gap-1 text-[11px] font-semibold px-2 py-1 rounded-md bg-white dark:bg-gray-800 hover:bg-gray-50 dark:hover:bg-gray-700"
                                                        style={{ border: `1px solid ${INPUT_BORDER}`, color: copiedGuideLabel === item.label ? GREEN : '#6b7280' }}
                                                    >
                                                        {copiedGuideLabel === item.label ? (
                                                            <>
                                                                <CheckIcon className="h-3.5 w-3.5" /> Copiado
                                                            </>
                                                        ) : (
                                                            <>
                                                                <ClipboardDocumentIcon className="h-3.5 w-3.5" /> Copiar
                                                            </>
                                                        )}
                                                    </button>
                                                </div>
                                                <code className="block w-full rounded-lg bg-gray-50 dark:bg-gray-900 px-3 py-2 text-[11px] font-mono text-gray-700 dark:text-gray-300 break-all" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                                                    {item.text}
                                                </code>
                                            </div>
                                        ))}
                                    </div>
                                )}

                                <div className="flex items-center justify-center gap-1.5 mt-3">
                                    {SHOPIFY_CONNECTION_STEPS.map((step, idx) => (
                                        <button
                                            key={step.step}
                                            type="button"
                                            onClick={() => setGuideStepIndex(idx)}
                                            aria-label={`Ir al paso ${idx + 1}`}
                                            className="h-1.5 rounded-full transition-all"
                                            style={{
                                                width: idx === guideStepIndex ? '18px' : '6px',
                                                backgroundColor: idx === guideStepIndex ? GREEN : INPUT_BORDER,
                                            }}
                                        />
                                    ))}
                                </div>
                            </div>
                        );
                    })()}
                </div>
            )}

            {isEdit && integrationId && (
                <>
                    <div
                        className="rounded-xl p-4 dark:bg-gray-800/60"
                        style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
                    >
                        <ShopifyWebhookManager integrationId={integrationId} />
                    </div>

                    <div
                        className="rounded-xl p-4 dark:bg-gray-800/60"
                        style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
                    >
                        <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">Sincronizar productos</h4>
                        <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                            Cruza los productos por SKU; crea en Shopify o en Probability los que falten y asocia los que coinciden.
                        </p>
                        <button
                            type="button"
                            onClick={() => setProductSyncOpen(true)}
                            className="mt-3 w-full inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors"
                            style={{ backgroundColor: GREEN }}
                        >
                            Sincronizar productos
                        </button>
                    </div>

                    <ShopifyProductSyncModal
                        isOpen={productSyncOpen}
                        onClose={() => setProductSyncOpen(false)}
                        integrationId={integrationId}
                        businessId={initialData?.business_id ?? null}
                    />

                    <ShopifyInventorySection
                        value={inventorySync}
                        onChange={setInventorySync}
                        businessId={initialData?.business_id ?? null}
                        integrationId={integrationId}
                        onSyncNow={() => setInventorySyncOpen(true)}
                        canSyncNow={inventorySync.enabled}
                    />

                    {inventorySync.enabled && (
                        <ShopifyLocationMappingSection
                            mappings={locationMappings}
                            onChangeMappings={setLocationMappings}
                            defaultLocationId={defaultLocationId}
                            onChangeDefaultLocation={setDefaultLocationId}
                            businessId={initialData?.business_id ?? null}
                        />
                    )}

                    <ShopifyInventorySyncModal
                        isOpen={inventorySyncOpen}
                        onClose={() => setInventorySyncOpen(false)}
                        integrationId={integrationId}
                        businessId={initialData?.business_id ?? null}
                    />

                    <div
                        className="rounded-xl p-4 dark:bg-gray-800/60"
                        style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
                    >
                        <div className="flex items-center gap-2 mb-2">
                            <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                                <TruckIcon style={{ color: GREEN, width: 16, height: 16 }} />
                            </span>
                            <h3 className="text-sm font-bold text-gray-900 dark:text-white">Configuracion de envio</h3>
                        </div>
                        <div
                            className="rounded-lg bg-white dark:bg-gray-800 divide-y divide-gray-100 dark:divide-gray-700"
                            style={{ border: `1px solid ${INPUT_BORDER}` }}
                        >
                            <ToggleRow
                                icon={<TruckIcon className="w-4 h-4" style={{ color: 'var(--color-primary)' }} />}
                                title="Cotizacion en checkout"
                                subtitle="Tarifas en tiempo real con varias transportadoras al pagar"
                                checked={carrierEnabled}
                                onToggle={handleToggleCarrierService}
                                disabled={carrierLoading}
                                activeColor="indigo"
                            />
                            <div>
                                <ToggleRow
                                    icon={<BeakerIcon className="w-4 h-4" style={{ color: 'var(--color-primary)' }} />}
                                    title="Modo de pruebas"
                                    subtitle="Redirige las peticiones a la URL de pruebas"
                                    checked={isTesting}
                                    onToggle={() => setIsTesting(!isTesting)}
                                    activeColor="orange"
                                />
                                {isTesting && initialData?.base_url_test && (
                                    <p className="px-3 pb-2.5 -mt-1 text-[11px] font-mono text-orange-700 dark:text-orange-400 break-all">
                                        Sandbox: {initialData.base_url_test}
                                    </p>
                                )}
                            </div>
                        </div>
                    </div>
                </>
            )}

            <div className="flex flex-col-reverse gap-2.5 pt-3 border-t border-gray-100 dark:border-gray-700 sm:flex-row sm:justify-end sm:items-center">
                {onCancel && (
                    <button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="px-5 py-2 text-[13px] font-semibold rounded-lg bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                        style={{ border: `1px solid ${INPUT_BORDER}` }}
                    >
                        Cancelar
                    </button>
                )}
                <button
                    type="submit"
                    disabled={loading || !formData.name || !formData.shop_domain || (!isEdit && (!formData.client_id || !formData.client_secret))}
                    className="px-5 py-2 text-[13px] font-semibold rounded-lg text-white flex items-center justify-center gap-2 transition-colors disabled:opacity-60"
                    style={{ backgroundColor: GREEN }}
                    onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                    onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                >
                    {loading && (
                        <svg className="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                    )}
                    {isEdit
                        ? (loading ? 'Guardando...' : 'Guardar integracion')
                        : (loading ? 'Conectando...' : 'Conectar con Shopify')}
                </button>
            </div>
        </form>
    );
}
