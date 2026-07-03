'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Select, Modal, Alert, SecretInput, ConfirmModal } from '@/shared/ui';
import { WooCommerceCredentials, WooCommerceConfig } from '../../domain/types';
import { createIntegrationAction, updateIntegrationAction, testConnectionRawAction, getActiveIntegrationTypesAction, getWooCommerceConnectionInfoAction, rotateWooCommerceTokenAction, revokeWooCommerceTokenAction } from '@/services/integrations/core/infra/actions';
import { WooCommerceConnectionInfo } from '@/services/integrations/core/domain/types';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import { WooProductSyncModal } from './WooProductSyncModal';
import { WooWebhookManager } from './WooWebhookManager';
import { getWooPluginZipAction } from '../../infra/actions';
import {
    KeyIcon,
    Cog6ToothIcon,
    ShoppingBagIcon,
    InformationCircleIcon,
    BoltIcon,
    TruckIcon,
    ArrowDownTrayIcon,
    ClipboardDocumentIcon,
    ClipboardDocumentCheckIcon,
    ArrowPathIcon,
    NoSymbolIcon,
    ChevronDownIcon,
    PhotoIcon,
    BeakerIcon,
} from '@heroicons/react/24/outline';

const HELP_IMAGES = [
    {
        src: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/manuals/woocommerce/step-1-rest-api-keys.png',
        caption: 'Paso 1: en WordPress ve a WooCommerce -> Ajustes -> Avanzado -> REST API y haz clic en "Agregar clave".',
    },
    {
        src: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/manuals/woocommerce/step-2-crear-key.png',
        caption: 'Paso 2: escribe una descripcion, elige permisos "Lectura/Escritura" y genera la clave. Copia el Consumer Key y el Consumer Secret.',
    },
];

interface WooCommerceConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    isEdit?: boolean;
    integrationId?: number;
    initialData?: {
        name?: string;
        store_id?: string;
        config?: any;
        credentials?: any;
        business_id?: number | null;
        is_testing?: boolean;
    };
}

const GREEN = 'var(--color-primary)';
const GREEN_DARK = 'color-mix(in srgb, var(--color-primary) 85%, black)';
const GREEN_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
const GREEN_BORDER = 'color-mix(in srgb, var(--color-primary) 25%, white)';
const CARD_BG = '#fafafd';
const CARD_BORDER = '#eceaf3';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1 flex items-start gap-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]';

const GUIDE_STEPS = [
    'Ingresa al panel de WordPress',
    'Ve a Ajustes / Avanzado / REST API',
    'Haz clic en "Agregar clave"',
    'Permisos de Lectura/Escritura',
    'Copia el Consumer Key y Secret',
];

const PLUGIN_STEPS = [
    'Descarga el plugin con el boton de abajo',
    'En WordPress: Plugins → Anadir nuevo → Subir plugin',
    'Sube el .zip y activa el plugin',
    'WooCommerce → Ajustes → Envio → tu zona → Anadir "Probability (Transportadoras)"',
    'Pega la Clave de conexion y guarda',
];

export function WooCommerceConfigForm({ onSuccess, onCancel, isEdit, integrationId, initialData }: WooCommerceConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(initialData?.business_id ?? null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);
    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);
    const [connInfo, setConnInfo] = useState<WooCommerceConnectionInfo | null>(null);
    const [loadingConn, setLoadingConn] = useState(false);
    const [copiedKey, setCopiedKey] = useState(false);

    useEffect(() => {
        if (!isEdit || !integrationId) return;
        let active = true;
        setLoadingConn(true);
        getWooCommerceConnectionInfoAction(integrationId)
            .then((res: any) => {
                if (active && res && res.connection_key) {
                    setConnInfo(res as WooCommerceConnectionInfo);
                }
            })
            .catch(() => { })
            .finally(() => { if (active) setLoadingConn(false); });
        return () => { active = false; };
    }, [isEdit, integrationId]);

    const [rotating, setRotating] = useState(false);
    const [revoking, setRevoking] = useState(false);
    const [showRevokeConfirm, setShowRevokeConfirm] = useState(false);
    const [productSyncOpen, setProductSyncOpen] = useState(false);
    const [showHelpImages, setShowHelpImages] = useState(false);
    const [isTesting, setIsTesting] = useState<boolean>(!!initialData?.is_testing);
    const [downloadingPlugin, setDownloadingPlugin] = useState(false);

    const handleDownloadPlugin = async () => {
        setDownloadingPlugin(true);
        try {
            const res = await getWooPluginZipAction();
            if (!res?.success || !res?.data) {
                setErrorModal(res?.message || 'No se pudo descargar el plugin');
                return;
            }
            const bytes = Uint8Array.from(atob(res.data), (c) => c.charCodeAt(0));
            const blob = new Blob([bytes], { type: 'application/zip' });
            const url = URL.createObjectURL(blob);
            const a = document.createElement('a');
            a.href = url;
            a.download = 'probability-shipping.zip';
            document.body.appendChild(a);
            a.click();
            a.remove();
            URL.revokeObjectURL(url);
        } catch {
            setErrorModal('No se pudo descargar el plugin');
        } finally {
            setDownloadingPlugin(false);
        }
    };

    const handleCopyKey = async () => {
        if (!connInfo?.connection_key) return;
        await navigator.clipboard.writeText(connInfo.connection_key);
        setCopiedKey(true);
        setTimeout(() => setCopiedKey(false), 2000);
    };

    const handleRotateKey = async () => {
        if (!integrationId || rotating) return;
        setRotating(true);
        try {
            const res: any = await rotateWooCommerceTokenAction(integrationId);
            if (res && res.connection_key) {
                setConnInfo(res as WooCommerceConnectionInfo);
                showToast('Clave rotada. La clave anterior dejo de funcionar.', 'success');
            } else {
                showToast(res?.message || 'No se pudo rotar la clave', 'error');
            }
        } finally {
            setRotating(false);
        }
    };

    const doRevokeKey = async () => {
        if (!integrationId || revoking) return;
        setRevoking(true);
        try {
            const res: any = await revokeWooCommerceTokenAction(integrationId);
            if (res && res.revoked) {
                setConnInfo(res as WooCommerceConnectionInfo);
                showToast('Clave revocada.', 'success');
            } else {
                showToast(res?.message || 'No se pudo revocar la clave', 'error');
            }
        } finally {
            setRevoking(false);
        }
    };

    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        store_url: initialData?.config?.store_url || initialData?.store_id || '',
        consumer_key: initialData?.credentials?.consumer_key || '',
        consumer_secret: initialData?.credentials?.consumer_secret || '',
    });

    useEffect(() => {
        const checkUserAndLoadBusinesses = async () => {
            const permissions = TokenStorage.getPermissions();
            const isSuperUser = permissions?.is_super || false;
            setIsSuperAdmin(isSuperUser);

            if (isSuperUser) {
                setLoadingBusinesses(true);
                try {
                    const response = await getBusinessesSimpleAction();
                    if (response.success && response.data) {
                        setBusinesses(response.data);
                    }
                } catch (err) {
                    console.error('Error loading businesses:', err);
                    showToast('Error al cargar la lista de negocios', 'error');
                } finally {
                    setLoadingBusinesses(false);
                }
            } else {
                if (permissions?.business_id) {
                    setSelectedBusinessId(permissions.business_id);
                }
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    useEffect(() => {
        let cancelled = false;
        getActiveIntegrationTypesAction()
            .then((res: any) => {
                if (cancelled) return;
                const types = res?.data || [];
                const woo = types.find((t: any) => t.id === 4 || /wooc/i.test(t.code || ''));
                if (woo?.image_url) setLogoUrl(woo.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

    const connectionReady = !!formData.store_url && !!formData.consumer_key && !!formData.consumer_secret;

    const handleTestConnection = async () => {
        if (!connectionReady) {
            showToast('Debes ingresar la URL de la tienda, Consumer Key y Consumer Secret para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials: WooCommerceCredentials = {
                consumer_key: formData.consumer_key,
                consumer_secret: formData.consumer_secret,
            };

            const config: WooCommerceConfig = {
                store_url: formData.store_url,
            };

            const result = await testConnectionRawAction('woocommerce', config, credentials);

            if (result.success) {
                showToast('Conexion exitosa con WooCommerce', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con WooCommerce');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const config: WooCommerceConfig = {
                store_url: formData.store_url,
            };

            if (isEdit && integrationId) {
                const credentials: any = {};
                if (formData.consumer_key) credentials.consumer_key = formData.consumer_key;
                if (formData.consumer_secret) credentials.consumer_secret = formData.consumer_secret;

                const response: any = await updateIntegrationAction(integrationId, {
                    name: formData.name,
                    store_id: formData.store_url,
                    config: config as any,
                    credentials: Object.keys(credentials).length > 0 ? credentials : undefined,
                    is_testing: isSuperAdmin ? isTesting : undefined,
                });

                if (!response || response.success === false) {
                    throw new Error(response?.message || 'Error al actualizar integracion');
                }
                showToast('Integracion WooCommerce actualizada', 'success');
                onSuccess?.();
                return;
            }

            if (isSuperAdmin && !selectedBusinessId) {
                setErrorModal('Debes seleccionar un negocio antes de crear la integracion.');
                setLoading(false);
                return;
            }

            const credentials: WooCommerceCredentials = {
                consumer_key: formData.consumer_key,
                consumer_secret: formData.consumer_secret,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `woocommerce_${Date.now()}`,
                integration_type_id: 4,
                category: 'ecommerce',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
                is_testing: isSuperAdmin ? isTesting : false,
            });

            if (response.success) {
                showToast('Integracion WooCommerce creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al guardar la integracion de WooCommerce');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-3" autoComplete="off">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-center gap-3">
                    <span
                        className="flex h-11 w-11 items-center justify-center rounded-xl overflow-hidden shrink-0"
                        style={{ backgroundColor: logoUrl && !logoFailed ? GREEN_SOFT : GREEN, border: `1px solid ${GREEN_BORDER}` }}
                    >
                        {logoUrl && !logoFailed ? (
                            <img
                                src={logoUrl}
                                alt="WooCommerce"
                                className="h-8 w-8 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <ShoppingBagIcon className="h-6 w-6 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-lg font-bold text-gray-900 dark:text-white leading-tight">WooCommerce</h2>
                        <p className="text-xs text-gray-500 dark:text-gray-300">
                            Conecta tu tienda para sincronizar ordenes y productos a Probability.
                        </p>
                    </div>
                </div>
                <span
                    className="inline-flex items-center gap-2 self-start rounded-full px-3 py-1 text-[11px] font-semibold"
                    style={connectionReady
                        ? { backgroundColor: GREEN_SOFT, border: `1px solid ${GREEN_BORDER}`, color: GREEN_DARK }
                        : { backgroundColor: '#f3f4f6', border: '1px solid #e5e7eb', color: '#6b7280' }}
                >
                    <span className={connectionReady ? 'h-2 w-2 rounded-full animate-pulse' : 'h-2 w-2 rounded-full'} style={{ backgroundColor: connectionReady ? 'var(--color-primary)' : '#9ca3af' }} />
                    {connectionReady ? 'Listo para probar' : 'Datos incompletos'}
                </span>
            </div>

            <div
                className="rounded-xl p-4 dark:bg-gray-800/60"
                style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex items-center gap-2 mb-3">
                    <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                        <Cog6ToothIcon className="w-4.5 h-4.5" style={{ color: GREEN, width: 16, height: 16 }} />
                    </span>
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Configuracion general</h3>
                </div>

                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>
                            Nombre de la Integracion <span style={{ color: GREEN }}>*</span>
                        </label>
                        <input
                            type="text"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            placeholder="Ej: WooCommerce Tienda Principal"
                            required
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Nombre descriptivo para identificar esta integracion</span>
                        </p>
                    </div>

                    <div>
                        <label className={fieldLabel}>
                            URL de la tienda <span style={{ color: GREEN }}>*</span>
                        </label>
                        <input
                            type="url"
                            value={formData.store_url}
                            onChange={(e) => setFormData({ ...formData, store_url: e.target.value })}
                            placeholder="https://mitienda.com"
                            required
                            autoComplete="off"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>URL completa de tu tienda WooCommerce (sin barra final)</span>
                        </p>
                    </div>

                    {isSuperAdmin && !isEdit && (
                        <div className="md:col-span-2">
                            <label className={fieldLabel}>
                                Negocio <span style={{ color: GREEN }}>*</span>
                            </label>
                            {loadingBusinesses ? (
                                <div className="flex items-center gap-2 p-3 bg-white dark:bg-gray-800 rounded-xl" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                                    <svg className="animate-spin h-5 w-5" style={{ color: GREEN }} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                    </svg>
                                    <span className="text-sm text-gray-600 dark:text-gray-300">Cargando negocios...</span>
                                </div>
                            ) : (
                                <Select
                                    value={selectedBusinessId?.toString() || ''}
                                    onChange={(e) => setSelectedBusinessId(Number(e.target.value))}
                                    options={[
                                        { value: '', label: '-- Selecciona un negocio --' },
                                        ...businesses.map((business) => ({
                                            value: business.id.toString(),
                                            label: business.name,
                                        })),
                                    ]}
                                    required
                                    className="bg-white dark:bg-gray-800"
                                />
                            )}
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Negocio al que pertenecera esta integracion</span>
                            </p>
                        </div>
                    )}
                </div>
            </div>

            <div
                className="rounded-xl p-4 dark:bg-gray-800/60"
                style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex items-center gap-2 mb-3">
                    <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                        <KeyIcon style={{ color: GREEN, width: 16, height: 16 }} />
                    </span>
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Credenciales de acceso</h3>
                </div>

                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>
                            Consumer Key <span style={{ color: GREEN }}>*</span>
                        </label>
                        <SecretInput
                            value={formData.consumer_key}
                            onChange={(e) => setFormData({ ...formData, consumer_key: e.target.value })}
                            placeholder="ck_xxxxxxxxxxxxxxxxxxxx"
                            required={!isEdit}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>

                    <div>
                        <label className={fieldLabel}>
                            Consumer Secret <span style={{ color: GREEN }}>*</span>
                        </label>
                        <SecretInput
                            value={formData.consumer_secret}
                            onChange={(e) => setFormData({ ...formData, consumer_secret: e.target.value })}
                            placeholder="cs_xxxxxxxxxxxxxxxxxxxx"
                            required={!isEdit}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>
                </div>

                <button
                    type="button"
                    onClick={handleTestConnection}
                    disabled={testingConnection || loading || !connectionReady}
                    className="mt-3.5 w-full flex items-center justify-center gap-2 rounded-lg py-2 text-[13px] font-semibold transition-colors disabled:opacity-50"
                    style={{
                        border: `2px dashed ${GREEN_BORDER}`,
                        backgroundColor: 'rgba(234, 250, 240, 0.5)',
                        color: GREEN,
                    }}
                >
                    {testingConnection ? (
                        <>
                            <svg className="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Probando...
                        </>
                    ) : (
                        <>
                            <BoltIcon className="w-4 h-4" />
                            Probar conexion
                        </>
                    )}
                </button>
            </div>

            {(isSuperAdmin || isTesting) && (
                <div
                    className="rounded-xl p-4 dark:bg-gray-800/60"
                    style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
                >
                    <div className="flex items-center justify-between gap-3">
                        <div className="flex items-start gap-2">
                            <BeakerIcon className="w-4 h-4 mt-0.5" style={{ color: '#d97706' }} />
                            <div>
                                <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                                    Modo de pruebas
                                    {isTesting && (
                                        <span className="inline-flex items-center rounded px-1.5 py-0.5 text-[10px] font-semibold uppercase bg-amber-50 text-amber-700 border border-amber-200 dark:bg-amber-900/30 dark:text-amber-400 dark:border-amber-800">
                                            Activo
                                        </span>
                                    )}
                                </h4>
                                <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                                    {isSuperAdmin
                                        ? 'Redirige las peticiones a la tienda de pruebas (mock) configurada en el tipo de integracion, sin tocar tu tienda real.'
                                        : 'Esta integracion esta en modo de pruebas. Solo un super admin puede cambiarlo.'}
                                </p>
                            </div>
                        </div>
                        <button
                            type="button"
                            role="switch"
                            aria-checked={isTesting}
                            disabled={!isSuperAdmin}
                            onClick={() => { if (isSuperAdmin) setIsTesting((v) => !v); }}
                            title={isSuperAdmin ? undefined : 'Solo un super admin puede cambiar el modo de pruebas'}
                            className={`relative inline-flex h-6 w-11 flex-shrink-0 items-center rounded-full transition-colors ${isTesting ? 'bg-amber-500' : 'bg-gray-300 dark:bg-gray-600'} ${!isSuperAdmin ? 'opacity-60 cursor-not-allowed' : ''}`}
                        >
                            <span className={`inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform ${isTesting ? 'translate-x-5' : 'translate-x-0.5'}`} />
                        </button>
                    </div>
                </div>
            )}

            <div
                className="rounded-xl p-4 dark:bg-gray-800/60"
                style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex flex-col gap-1 mb-3 sm:flex-row sm:items-center sm:justify-between">
                    <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100 flex items-center gap-1.5">
                        <InformationCircleIcon className="w-4 h-4 text-gray-400" />
                        Como obtener tus credenciales
                    </h4>
                    <p className="text-[11px] text-gray-500 dark:text-gray-400">
                        En WordPress:{' '}
                        <strong style={{ color: GREEN_DARK }}>WooCommerce</strong>
                        <span className="mx-1 text-gray-400">&rarr;</span>
                        <strong style={{ color: GREEN_DARK }}>Ajustes</strong>
                        <span className="mx-1 text-gray-400">&rarr;</span>
                        <strong style={{ color: GREEN_DARK }}>Avanzado</strong>
                        <span className="mx-1 text-gray-400">&rarr;</span>
                        <strong style={{ color: GREEN_DARK }}>REST API</strong>
                    </p>
                </div>
                <ol className="grid grid-cols-2 gap-y-4 sm:grid-cols-5 sm:gap-y-0">
                    {GUIDE_STEPS.map((step, i) => (
                        <li key={i} className="flex flex-col">
                            <div className="flex items-center">
                                <span
                                    className="flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full text-[10px] font-bold text-white"
                                    style={{ backgroundColor: GREEN }}
                                >
                                    {i + 1}
                                </span>
                                {i < GUIDE_STEPS.length - 1 && (
                                    <span className="hidden sm:block flex-1 h-px mx-2" style={{ backgroundColor: INPUT_BORDER }} />
                                )}
                            </div>
                            <span className="mt-1.5 pr-2 text-[11px] text-gray-500 dark:text-gray-400 leading-snug">{step}</span>
                        </li>
                    ))}
                </ol>

                <div className="mt-4 pt-3" style={{ borderTop: `1px solid ${INPUT_BORDER}` }}>
                    <button
                        type="button"
                        onClick={() => setShowHelpImages((v) => !v)}
                        className="flex w-full items-center justify-between rounded-lg px-3 py-2 text-[12px] font-semibold text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-800/60 transition-colors"
                    >
                        <span className="flex items-center gap-1.5">
                            <PhotoIcon className="w-4 h-4" style={{ color: GREEN_DARK }} />
                            Ver imagenes de ayuda paso a paso
                        </span>
                        <ChevronDownIcon
                            className={`w-4 h-4 text-gray-400 transition-transform ${showHelpImages ? 'rotate-180' : ''}`}
                        />
                    </button>

                    {showHelpImages && (
                        <div className="mt-3 grid grid-cols-1 gap-4 sm:grid-cols-2">
                            {HELP_IMAGES.map((img, i) => (
                                <figure key={i} className="flex flex-col">
                                    <a href={img.src} target="_blank" rel="noopener noreferrer" className="block">
                                        <img
                                            src={img.src}
                                            alt={img.caption}
                                            loading="lazy"
                                            className="w-full rounded-lg border object-contain hover:opacity-95 transition-opacity"
                                            style={{ borderColor: INPUT_BORDER, backgroundColor: '#fff' }}
                                        />
                                    </a>
                                    <figcaption className="mt-1.5 text-[11px] text-gray-500 dark:text-gray-400 leading-snug">
                                        {img.caption}
                                    </figcaption>
                                </figure>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            {isEdit && integrationId && (
                <div
                    className="rounded-xl p-4 dark:bg-gray-800/60"
                    style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
                >
                    <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                        <div>
                            <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100 flex items-center gap-1.5">
                                <ArrowPathIcon className="w-4 h-4" style={{ color: GREEN_DARK }} />
                                Sincronizar productos a WooCommerce
                            </h4>
                            <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                                Publica los productos de Probability en tu tienda: crea los que no existen y actualiza el stock de los ya mapeados.
                            </p>
                        </div>
                        <button
                            type="button"
                            onClick={() => setProductSyncOpen(true)}
                            className="inline-flex items-center justify-center gap-1.5 self-start rounded-lg px-3 py-1.5 text-[12px] font-semibold text-white transition-colors"
                            style={{ backgroundColor: GREEN }}
                            onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                            onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                        >
                            <ArrowPathIcon className="w-3.5 h-3.5" />
                            Sincronizar productos
                        </button>
                    </div>
                </div>
            )}

            {isEdit && integrationId && (
                <div
                    className="rounded-xl p-4 dark:bg-gray-800/60"
                    style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
                >
                    <WooWebhookManager integrationId={integrationId} />
                </div>
            )}

            {isEdit && integrationId && (
                <div
                    className="rounded-xl p-4 dark:bg-gray-800/60"
                    style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
                >
                    <div className="flex flex-col gap-1 mb-2 sm:flex-row sm:items-center sm:justify-between">
                        <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100 flex items-center gap-1.5">
                            <TruckIcon className="w-4 h-4" style={{ color: GREEN_DARK }} />
                            Plugin de cotizacion de envios en el checkout
                        </h4>
                        <span
                            className="text-[10px] font-semibold px-2 py-0.5 rounded-full self-start"
                            style={{ backgroundColor: GREEN_SOFT, color: GREEN_DARK, border: `1px solid ${GREEN_BORDER}` }}
                        >
                            Opcional
                        </span>
                    </div>

                    <p className="text-[12px] text-gray-500 dark:text-gray-400 leading-relaxed mb-3">
                        Instala este plugin en tu WordPress para que tus clientes vean las tarifas reales de las
                        transportadoras (EnvioClick y otras) directamente en el checkout, calculadas por Probability
                        segun la direccion de envio. Requiere que tengas una transportadora y una direccion de origen
                        configuradas en Probability.
                    </p>

                    <ol className="grid grid-cols-1 gap-y-2 mb-4 sm:grid-cols-2 sm:gap-x-4">
                        {PLUGIN_STEPS.map((step, i) => (
                            <li key={i} className="flex items-start gap-2">
                                <span
                                    className="flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full text-[10px] font-bold text-white"
                                    style={{ backgroundColor: GREEN }}
                                >
                                    {i + 1}
                                </span>
                                <span className="text-[11px] text-gray-600 dark:text-gray-300 leading-snug">{step}</span>
                            </li>
                        ))}
                    </ol>

                    <div className="flex flex-col gap-3">
                        <button
                            type="button"
                            onClick={handleDownloadPlugin}
                            disabled={downloadingPlugin}
                            className="inline-flex items-center justify-center gap-2 px-4 py-2 text-[13px] font-semibold rounded-lg text-white transition-colors self-start disabled:opacity-60"
                            style={{ backgroundColor: GREEN }}
                        >
                            <ArrowDownTrayIcon className="w-4 h-4" />
                            {downloadingPlugin ? 'Descargando...' : 'Descargar plugin (.zip)'}
                        </button>

                        <div>
                            <label className={fieldLabel}>Clave de conexion</label>
                            <p className="text-[11px] text-gray-400 dark:text-gray-500 mb-1.5">
                                Copiala y pegala en los ajustes del metodo de envio "Probability" en tu WordPress.
                            </p>
                            {loadingConn ? (
                                <div className="text-[12px] text-gray-400 py-2">Cargando clave de conexion...</div>
                            ) : connInfo && connInfo.revoked ? (
                                <div className="flex flex-col gap-2">
                                    <div
                                        className="text-[12px] rounded-lg px-3 py-2"
                                        style={{ backgroundColor: '#fff4f4', color: '#b42318', border: '1px solid #f3c9c9' }}
                                    >
                                        La clave esta revocada. La tienda no cotizara hasta que generes una clave nueva.
                                    </div>
                                    <button
                                        type="button"
                                        onClick={handleRotateKey}
                                        disabled={rotating}
                                        className="px-3 py-2 text-[13px] font-semibold rounded-lg text-white flex items-center justify-center gap-1.5 self-start disabled:opacity-60"
                                        style={{ backgroundColor: GREEN }}
                                    >
                                        <ArrowPathIcon className="w-4 h-4" />
                                        {rotating ? 'Generando...' : 'Generar clave nueva'}
                                    </button>
                                </div>
                            ) : connInfo ? (
                                <div className="flex flex-col gap-2">
                                    <div className="flex flex-col gap-2 sm:flex-row sm:items-stretch">
                                        <textarea
                                            readOnly
                                            value={connInfo.connection_key}
                                            onFocus={(e) => e.currentTarget.select()}
                                            className={`${inputCls} font-mono text-[11px] resize-none flex-1`}
                                            rows={2}
                                            style={{ borderColor: INPUT_BORDER }}
                                        />
                                        <button
                                            type="button"
                                            onClick={handleCopyKey}
                                            className="px-3 py-2 text-[13px] font-semibold rounded-lg flex items-center justify-center gap-1.5 shrink-0 transition-colors"
                                            style={{ backgroundColor: GREEN_SOFT, color: GREEN_DARK, border: `1px solid ${GREEN_BORDER}` }}
                                        >
                                            {copiedKey ? (
                                                <>
                                                    <ClipboardDocumentCheckIcon className="w-4 h-4" />
                                                    Copiado
                                                </>
                                            ) : (
                                                <>
                                                    <ClipboardDocumentIcon className="w-4 h-4" />
                                                    Copiar
                                                </>
                                            )}
                                        </button>
                                    </div>
                                    <div className="flex items-center gap-2 flex-wrap">
                                        <button
                                            type="button"
                                            onClick={handleRotateKey}
                                            disabled={rotating}
                                            className="px-3 py-1.5 text-[12px] font-semibold rounded-lg flex items-center gap-1.5 disabled:opacity-60"
                                            style={{ backgroundColor: '#ffffff', color: '#374151', border: `1px solid ${INPUT_BORDER}` }}
                                        >
                                            <ArrowPathIcon className="w-3.5 h-3.5" />
                                            {rotating ? 'Rotando...' : 'Rotar clave'}
                                        </button>
                                        <button
                                            type="button"
                                            onClick={() => setShowRevokeConfirm(true)}
                                            disabled={revoking}
                                            className="px-3 py-1.5 text-[12px] font-semibold rounded-lg flex items-center gap-1.5 disabled:opacity-60"
                                            style={{ backgroundColor: '#ffffff', color: '#b42318', border: '1px solid #f3c9c9' }}
                                        >
                                            <NoSymbolIcon className="w-3.5 h-3.5" />
                                            {revoking ? 'Revocando...' : 'Revocar'}
                                        </button>
                                    </div>
                                    <p className="text-[10px] text-gray-400">
                                        Rotar genera una clave nueva y desactiva la anterior. Revocar detiene la cotizacion hasta generar una nueva.
                                    </p>
                                </div>
                            ) : (
                                <div className="text-[12px] text-gray-400 py-2">
                                    No se pudo cargar la clave de conexion. Guarda la integracion e intenta de nuevo.
                                </div>
                            )}
                        </div>
                    </div>
                </div>
            )}

            <ConfirmModal
                isOpen={showRevokeConfirm}
                onClose={() => setShowRevokeConfirm(false)}
                onConfirm={doRevokeKey}
                title="Revocar clave de conexion"
                message="Al revocar, la tienda dejara de cotizar envios hasta que generes una clave nueva. Esta accion invalida la clave actual. Deseas continuar?"
                confirmText="Revocar"
                type="danger"
            />

            {isEdit && integrationId && (
                <WooProductSyncModal
                    isOpen={productSyncOpen}
                    onClose={() => setProductSyncOpen(false)}
                    integrationId={integrationId}
                    businessId={selectedBusinessId}
                />
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
                    disabled={loading}
                    className="px-5 py-2 text-[13px] font-semibold rounded-lg text-white flex items-center justify-center gap-2 transition-colors disabled:opacity-60"
                    style={{ backgroundColor: GREEN }}
                    onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                    onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                >
                    {loading ? (
                        <>
                            <svg className="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            {isEdit ? 'Guardando...' : 'Conectando...'}
                        </>
                    ) : (
                        isEdit ? 'Guardar integracion' : 'Crear integracion'
                    )}
                </button>
            </div>

            {errorModal && (
                <Modal
                    isOpen={!!errorModal}
                    onClose={() => setErrorModal(null)}
                    title="Error"
                    size="sm"
                >
                    <div className="p-4">
                        <Alert type="error">{errorModal}</Alert>
                    </div>
                </Modal>
            )}
        </form>
    );
}
