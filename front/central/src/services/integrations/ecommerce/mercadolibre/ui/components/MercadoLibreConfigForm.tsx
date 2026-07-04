'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Select, Modal, Alert, SecretInput } from '@/shared/ui';
import { MercadoLibreCredentials, MercadoLibreConfig } from '../../domain/types';
import { createIntegrationAction, updateIntegrationAction, testConnectionRawAction, getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    ShoppingBagIcon,
    InformationCircleIcon,
    BoltIcon,
    ArrowPathIcon,
} from '@heroicons/react/24/outline';
import { MercadoLibreProductSyncModal } from './MercadoLibreProductSyncModal';
import { MercadoLibreInventorySection, MeliInventoryConfig } from './MercadoLibreInventorySection';
import { MercadoLibreInventorySyncModal } from './MercadoLibreInventorySyncModal';

interface MercadoLibreConfigFormProps {
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
    'Crea una aplicacion en developers.mercadolibre.com',
    'Copia el App ID y la Secret Key',
    'Autoriza la app con tu cuenta vendedor',
    'Genera el Access Token y Refresh Token',
    'Configura la URL de notificaciones (webhook)',
];

export function MercadoLibreConfigForm({ onSuccess, onCancel, isEdit, integrationId, initialData }: MercadoLibreConfigFormProps) {
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
    const [productSyncOpen, setProductSyncOpen] = useState(false);
    const [inventorySyncOpen, setInventorySyncOpen] = useState(false);
    const [inventorySync, setInventorySync] = useState<MeliInventoryConfig>(() => {
        const c: any = initialData?.config || {};
        return {
            enabled: !!c.inventory_sync_enabled,
            mode: c.inventory_warehouse_mode === 'single' ? 'single' : 'sum',
            single_warehouse_id: Number(c.inventory_single_warehouse_id) || 0,
            warehouse_ids: Array.isArray(c.inventory_warehouse_ids) ? c.inventory_warehouse_ids.map(Number) : [],
        };
    });

    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        client_id: initialData?.credentials?.client_id || '',
        client_secret: initialData?.credentials?.client_secret || '',
        access_token: initialData?.credentials?.access_token || '',
        refresh_token: initialData?.credentials?.refresh_token || '',
        seller_id: String(initialData?.config?.seller_id ?? initialData?.store_id ?? ''),
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
                const meli = types.find((t: any) => t.id === 3 || /mercado/i.test(t.code || ''));
                if (meli?.image_url) setLogoUrl(meli.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

    const connectionReady = !!formData.client_id && !!formData.client_secret && !!formData.access_token && !!formData.refresh_token;

    const handleTestConnection = async () => {
        if (!formData.access_token) {
            showToast('Debes ingresar el Access Token para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials: MercadoLibreCredentials = {
                client_id: formData.client_id,
                client_secret: formData.client_secret,
                access_token: formData.access_token,
                refresh_token: formData.refresh_token,
                seller_id: formData.seller_id,
            };

            const result = await testConnectionRawAction('mercado_libre', {}, credentials as any);

            if (result.success) {
                showToast('Conexion exitosa con MercadoLibre', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con MercadoLibre');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const config: any = {
                ...(initialData?.config || {}),
                seller_id: formData.seller_id || undefined,
                inventory_sync_enabled: inventorySync.enabled,
                inventory_warehouse_mode: inventorySync.mode,
                inventory_single_warehouse_id: inventorySync.single_warehouse_id,
                inventory_warehouse_ids: inventorySync.warehouse_ids,
            };

            if (isEdit && integrationId) {
                const credentials: any = {};
                if (formData.client_id) credentials.client_id = formData.client_id;
                if (formData.client_secret) credentials.client_secret = formData.client_secret;
                if (formData.access_token) credentials.access_token = formData.access_token;
                if (formData.refresh_token) credentials.refresh_token = formData.refresh_token;
                if (formData.seller_id) credentials.seller_id = formData.seller_id;

                const response: any = await updateIntegrationAction(integrationId, {
                    name: formData.name,
                    store_id: formData.seller_id || undefined,
                    config: config as any,
                    credentials: Object.keys(credentials).length > 0 ? credentials : undefined,
                });

                if (!response || response.success === false) {
                    throw new Error(response?.message || 'Error al actualizar integracion');
                }
                showToast('Integracion MercadoLibre actualizada', 'success');
                onSuccess?.();
                return;
            }

            if (isSuperAdmin && !selectedBusinessId) {
                setErrorModal('Debes seleccionar un negocio antes de crear la integracion.');
                setLoading(false);
                return;
            }

            const credentials: MercadoLibreCredentials = {
                client_id: formData.client_id,
                client_secret: formData.client_secret,
                access_token: formData.access_token,
                refresh_token: formData.refresh_token,
                seller_id: formData.seller_id,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `mercado_libre_${Date.now()}`,
                integration_type_id: 3,
                category: 'ecommerce',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                store_id: formData.seller_id || undefined,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
            });

            if (response.success) {
                showToast('Integracion MercadoLibre creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al guardar la integracion de MercadoLibre');
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
                                alt="MercadoLibre"
                                className="h-8 w-8 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <ShoppingBagIcon className="h-6 w-6 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-lg font-bold text-gray-900 dark:text-white leading-tight">MercadoLibre</h2>
                        <p className="text-xs text-gray-500 dark:text-gray-300">
                            Conecta tu cuenta para sincronizar ordenes y recibir notificaciones en tiempo real.
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
                        <Cog6ToothIcon style={{ color: GREEN, width: 16, height: 16 }} />
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
                            placeholder="Ej: MercadoLibre Principal"
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
                        <label className={fieldLabel}>Seller ID</label>
                        <input
                            type="text"
                            value={formData.seller_id}
                            onChange={(e) => setFormData({ ...formData, seller_id: e.target.value })}
                            placeholder="ID del vendedor (user_id de MercadoLibre)"
                            autoComplete="off"
                            className={`${inputCls} font-mono`}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Se usa para asociar las notificaciones (webhooks) a esta cuenta</span>
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
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Credenciales OAuth</h3>
                </div>

                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>
                            App ID (Client ID) <span style={{ color: GREEN }}>*</span>
                        </label>
                        <input
                            type="text"
                            value={formData.client_id}
                            onChange={(e) => setFormData({ ...formData, client_id: e.target.value })}
                            placeholder="Ej: 1234567890123456"
                            required={!isEdit}
                            autoComplete="off"
                            className={`${inputCls} font-mono`}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>

                    <div>
                        <label className={fieldLabel}>
                            Secret Key (Client Secret) <span style={{ color: GREEN }}>*</span>
                        </label>
                        <SecretInput
                            value={formData.client_secret}
                            onChange={(e) => setFormData({ ...formData, client_secret: e.target.value })}
                            placeholder="Secret Key de la aplicacion"
                            required={!isEdit}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>

                    <div>
                        <label className={fieldLabel}>
                            Access Token <span style={{ color: GREEN }}>*</span>
                        </label>
                        <SecretInput
                            value={formData.access_token}
                            onChange={(e) => setFormData({ ...formData, access_token: e.target.value })}
                            placeholder="APP_USR-..."
                            required={!isEdit}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>

                    <div>
                        <label className={fieldLabel}>
                            Refresh Token <span style={{ color: GREEN }}>*</span>
                        </label>
                        <SecretInput
                            value={formData.refresh_token}
                            onChange={(e) => setFormData({ ...formData, refresh_token: e.target.value })}
                            placeholder="TG-..."
                            required={!isEdit}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>
                </div>

                <button
                    type="button"
                    onClick={handleTestConnection}
                    disabled={testingConnection || loading || !formData.access_token}
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
                        En <strong style={{ color: GREEN_DARK }}>developers.mercadolibre.com</strong>
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
                                Sincronizar productos
                            </h4>
                            <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                                Cruza los productos por SKU y te muestra que falta en cada lado; eliges crear en MercadoLibre lo que solo esta en Probability, o al reves.
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

            <MercadoLibreInventorySection
                value={inventorySync}
                onChange={setInventorySync}
                businessId={selectedBusinessId}
                integrationId={isEdit ? integrationId : undefined}
                onSyncNow={isEdit && integrationId ? () => setInventorySyncOpen(true) : undefined}
                canSyncNow={inventorySync.enabled}
            />

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

            {isEdit && integrationId && (
                <MercadoLibreProductSyncModal
                    isOpen={productSyncOpen}
                    onClose={() => setProductSyncOpen(false)}
                    integrationId={integrationId}
                    businessId={selectedBusinessId}
                />
            )}

            {isEdit && integrationId && (
                <MercadoLibreInventorySyncModal
                    isOpen={inventorySyncOpen}
                    onClose={() => setInventorySyncOpen(false)}
                    integrationId={integrationId}
                    businessId={selectedBusinessId}
                />
            )}

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
