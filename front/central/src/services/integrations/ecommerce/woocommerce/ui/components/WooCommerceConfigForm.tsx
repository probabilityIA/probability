'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Select, Modal, Alert, SecretInput } from '@/shared/ui';
import { WooCommerceCredentials, WooCommerceConfig } from '../../domain/types';
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
} from '@heroicons/react/24/outline';

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
    };
}

const GREEN = '#1F8A5B';
const GREEN_DARK = '#15803d';
const GREEN_SOFT = '#eafaf0';
const GREEN_BORDER = '#c7eed5';
const CARD_BG = '#fafafd';
const CARD_BORDER = '#eceaf3';
const INPUT_BORDER = '#e9e9f0';

const fieldLabel = 'block text-sm font-semibold text-gray-900 dark:text-gray-100 mb-1.5';
const fieldHint = 'text-xs text-gray-400 dark:text-gray-500 mt-1.5 flex items-start gap-1';
const inputCls = 'w-full px-3.5 py-2.5 text-sm rounded-xl border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-emerald-500/30 focus:border-emerald-500';

const GUIDE_STEPS = [
    'Ingresa al panel de WordPress',
    'Ve a Ajustes / Avanzado / REST API',
    'Haz clic en "Agregar clave"',
    'Permisos de Lectura/Escritura',
    'Copia el Consumer Key y Secret',
];

export function WooCommerceConfigForm({ onSuccess, onCancel, isEdit, integrationId, initialData }: WooCommerceConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);
    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);

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
        <form onSubmit={handleSubmit} className="space-y-4" autoComplete="off">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-center gap-3.5">
                    <span
                        className="flex h-14 w-14 items-center justify-center rounded-2xl overflow-hidden shrink-0"
                        style={{ backgroundColor: logoUrl && !logoFailed ? GREEN_SOFT : GREEN, border: `1px solid ${GREEN_BORDER}` }}
                    >
                        {logoUrl && !logoFailed ? (
                            <img
                                src={logoUrl}
                                alt="WooCommerce"
                                className="h-10 w-10 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <ShoppingBagIcon className="h-7 w-7 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-2xl font-bold text-gray-900 dark:text-white leading-tight">WooCommerce</h2>
                        <p className="text-sm text-gray-500 dark:text-gray-300">
                            Conecta tu tienda para sincronizar ordenes y productos a Probability.
                        </p>
                    </div>
                </div>
                <span
                    className="inline-flex items-center gap-2 self-start rounded-full px-3.5 py-1.5 text-xs font-semibold"
                    style={connectionReady
                        ? { backgroundColor: GREEN_SOFT, border: `1px solid ${GREEN_BORDER}`, color: GREEN_DARK }
                        : { backgroundColor: '#f3f4f6', border: '1px solid #e5e7eb', color: '#6b7280' }}
                >
                    <span className="h-2 w-2 rounded-full" style={{ backgroundColor: connectionReady ? '#22c55e' : '#9ca3af' }} />
                    {connectionReady ? 'Listo para probar' : 'Datos incompletos'}
                </span>
            </div>

            <div
                className="rounded-2xl p-5 dark:bg-gray-800/60"
                style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex items-center gap-2.5 mb-4">
                    <span className="flex h-8 w-8 items-center justify-center rounded-lg" style={{ backgroundColor: GREEN_SOFT }}>
                        <Cog6ToothIcon className="w-4.5 h-4.5" style={{ color: GREEN, width: 18, height: 18 }} />
                    </span>
                    <h3 className="text-base font-bold text-gray-900 dark:text-white">Configuracion general</h3>
                </div>

                <div className="grid grid-cols-1 gap-x-5 gap-y-4 md:grid-cols-2">
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
                className="rounded-2xl p-5 dark:bg-gray-800/60"
                style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex items-center gap-2.5 mb-4">
                    <span className="flex h-8 w-8 items-center justify-center rounded-lg" style={{ backgroundColor: GREEN_SOFT }}>
                        <KeyIcon style={{ color: GREEN, width: 18, height: 18 }} />
                    </span>
                    <h3 className="text-base font-bold text-gray-900 dark:text-white">Credenciales de acceso</h3>
                </div>

                <div className="grid grid-cols-1 gap-x-5 gap-y-4 md:grid-cols-2">
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
                    className="mt-5 w-full flex items-center justify-center gap-2 rounded-xl py-2.5 text-sm font-semibold transition-colors disabled:opacity-50"
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
                className="rounded-2xl p-5 dark:bg-gray-800/60"
                style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex flex-col gap-1.5 mb-4 sm:flex-row sm:items-center sm:justify-between">
                    <h4 className="text-sm font-bold text-gray-900 dark:text-gray-100 flex items-center gap-2">
                        <InformationCircleIcon className="w-5 h-5 text-gray-400" />
                        Como obtener tus credenciales
                    </h4>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
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
                                    className="flex h-6 w-6 flex-shrink-0 items-center justify-center rounded-full text-[11px] font-bold text-white"
                                    style={{ backgroundColor: GREEN }}
                                >
                                    {i + 1}
                                </span>
                                {i < GUIDE_STEPS.length - 1 && (
                                    <span className="hidden sm:block flex-1 h-px mx-2" style={{ backgroundColor: INPUT_BORDER }} />
                                )}
                            </div>
                            <span className="mt-2 pr-3 text-xs text-gray-500 dark:text-gray-400 leading-snug">{step}</span>
                        </li>
                    ))}
                </ol>
            </div>

            <div className="flex flex-col-reverse gap-3 pt-4 border-t border-gray-100 dark:border-gray-700 sm:flex-row sm:justify-end sm:items-center">
                {onCancel && (
                    <button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="px-6 py-2.5 text-sm font-semibold rounded-xl bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                        style={{ border: `1px solid ${INPUT_BORDER}` }}
                    >
                        Cancelar
                    </button>
                )}
                <button
                    type="submit"
                    disabled={loading}
                    className="px-6 py-2.5 text-sm font-semibold rounded-xl text-white flex items-center justify-center gap-2 transition-colors disabled:opacity-60"
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
