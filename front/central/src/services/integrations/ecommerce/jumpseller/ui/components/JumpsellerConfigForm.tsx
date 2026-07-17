'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Alert, Modal, SecretInput } from '@/shared/ui';
import { JumpsellerCredentials } from '../../domain/types';
import { createIntegrationAction, updateIntegrationAction, testConnectionRawAction, getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { JumpsellerProductSyncModal } from './JumpsellerProductSyncModal';
import { JumpsellerInventorySyncModal } from './JumpsellerInventorySyncModal';
import { JumpsellerOrderSyncModal } from './JumpsellerOrderSyncModal';
import { JumpsellerInventorySection, JumpsellerInventoryConfig, JumpsellerLocationsInfo } from './JumpsellerInventorySection';
import { JumpsellerWebhookManager } from './JumpsellerWebhookManager';
import { getJumpsellerLocationsAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    CheckBadgeIcon,
    InformationCircleIcon,
    BeakerIcon,
    ArrowPathIcon,
    ArrowsRightLeftIcon,
    ShoppingBagIcon,
    PhotoIcon,
    ChevronDownIcon,
    BoltIcon,
} from '@heroicons/react/24/outline';
import {
    GREEN,
    GREEN_DARK,
    GREEN_SOFT,
    GREEN_BORDER,
    INPUT_BORDER,
    fieldLabel,
    fieldHint,
    inputCls,
    SectionCard,
    ToggleRow,
    Spinner,
} from '@/services/integrations/invoicing/siigo/ui/components/SiigoFormKit';

const JUMPSELLER_TYPE_ID = 33;

const toDateInput = (d: Date) => {
    const month = `${d.getMonth() + 1}`.padStart(2, '0');
    const day = `${d.getDate()}`.padStart(2, '0');
    return `${d.getFullYear()}-${month}-${day}`;
};

const defaultOrdersFrom = () => {
    const d = new Date();
    d.setDate(d.getDate() - 30);
    return toDateInput(d);
};

const defaultOrdersTo = () => toDateInput(new Date());

const HELP_IMAGES = [
    {
        src: 'https://probability-media-assets.s3.us-east-1.amazonaws.com/manuals/jumpseller/step-1-credenciales-api.png',
        caption: 'En el panel de Jumpseller entra a Cuenta > Preferencias. En el recuadro "API y MCP" (arriba a la derecha) estan el Login y el Auth Token: copialos con el boton de la derecha de cada campo.',
    },
];

interface JumpsellerConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    integrationTypeBaseURLTest?: string;
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

export function JumpsellerConfigForm({ onSuccess, onCancel, integrationTypeBaseURLTest, isEdit, integrationId, initialData }: JumpsellerConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [isTesting, setIsTesting] = useState<boolean>(!!initialData?.is_testing);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(initialData?.business_id ?? null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);

    const [inventorySyncEnabled, setInventorySyncEnabled] = useState<boolean>(!!initialData?.config?.inventory_sync_enabled);
    const [inventorySync, setInventorySync] = useState<JumpsellerInventoryConfig>(() => {
        const c: any = initialData?.config || {};
        return {
            enabled: !!c.inventory_sync_enabled,
            mode: 'single',
            single_warehouse_id: Number(c.inventory_single_warehouse_id) || 0,
            default_location_id: c.jumpseller_default_location_id ? String(c.jumpseller_default_location_id) : '',
            mappings: Array.isArray(c.jumpseller_location_mappings)
                ? c.jumpseller_location_mappings.map((m: any) => ({
                    internal_warehouse_id: Number(m.internal_warehouse_id) || 0,
                    jumpseller_location_id: Number(m.jumpseller_location_id) || 0,
                }))
                : [],
        };
    });
    const [locationsInfo, setLocationsInfo] = useState<JumpsellerLocationsInfo | null>(null);
    const [statusSyncEnabled, setStatusSyncEnabled] = useState<boolean>(!!initialData?.config?.status_sync_enabled);
    const [productSyncOpen, setProductSyncOpen] = useState(false);
    const [showHelpImages, setShowHelpImages] = useState(false);

    const inventorySyncSaved = !!initialData?.config?.inventory_sync_enabled;

    const [ordersFrom, setOrdersFrom] = useState<string>(defaultOrdersFrom());
    const [ordersTo, setOrdersTo] = useState<string>(defaultOrdersTo());
    const [orderSyncOpen, setOrderSyncOpen] = useState(false);
    const [inventorySyncOpen, setInventorySyncOpen] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        api_key: initialData?.credentials?.api_key || '',
        api_secret: initialData?.credentials?.api_secret || '',
    });

    useEffect(() => {
        let cancelled = false;
        getActiveIntegrationTypesAction()
            .then((res: any) => {
                if (cancelled) return;
                const types = res?.data || [];
                const jumpseller = types.find((t: any) => t.id === JUMPSELLER_TYPE_ID || /jumpseller/i.test(t.code || ''));
                if (jumpseller?.image_url) setLogoUrl(jumpseller.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

    useEffect(() => {
        if (!inventorySyncEnabled || !isEdit || !integrationId) return;
        let cancelled = false;
        getJumpsellerLocationsAction(integrationId, selectedBusinessId ?? undefined)
            .then((res: any) => {
                if (cancelled || !res?.success) return;
                setLocationsInfo({
                    locations: Array.isArray(res.locations) ? res.locations : [],
                    multi_location: !!res.multi_location,
                    subscription_plan: res.subscription_plan || '',
                    stock_origin_name: res.stock_origin_name || '',
                });
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, [inventorySyncEnabled, isEdit, integrationId, selectedBusinessId]);

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

    const handleTestConnection = async () => {
        if (!formData.api_key || !formData.api_secret) {
            showToast('Debes ingresar Login y Auth Token para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials: JumpsellerCredentials = {
                api_key: formData.api_key,
                api_secret: formData.api_secret,
            };

            const result = await testConnectionRawAction('jumpseller', { is_testing: isTesting }, credentials);

            if (result.success) {
                showToast('Conexion exitosa con Jumpseller', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con Jumpseller');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSyncOrders = () => {
        if (!integrationId) return;

        if (ordersFrom && ordersTo && ordersFrom > ordersTo) {
            showToast('La fecha inicial no puede ser mayor que la fecha final', 'warning');
            return;
        }

        setOrderSyncOpen(true);
    };

    const handleSyncInventory = () => {
        if (!integrationId) return;
        setInventorySyncOpen(true);
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const cleanMappings = inventorySync.mappings
                .filter((m) => m.internal_warehouse_id > 0 && m.jumpseller_location_id > 0)
                .map((m) => ({
                    internal_warehouse_id: m.internal_warehouse_id,
                    jumpseller_location_id: m.jumpseller_location_id,
                }));

            const config = {
                inventory_sync_enabled: inventorySyncEnabled,
                status_sync_enabled: statusSyncEnabled,
                inventory_warehouse_mode: inventorySync.mode,
                inventory_single_warehouse_id: inventorySync.single_warehouse_id,
                jumpseller_default_location_id: inventorySync.default_location_id.trim()
                    ? Number(inventorySync.default_location_id.trim()) || 0
                    : 0,
                jumpseller_location_mappings: cleanMappings,
            };

            if (isEdit && integrationId) {
                const editCredentials: any = {};
                if (formData.api_key) editCredentials.api_key = formData.api_key;
                if (formData.api_secret) editCredentials.api_secret = formData.api_secret;

                const response: any = await updateIntegrationAction(integrationId, {
                    name: formData.name,
                    config: config as any,
                    credentials: Object.keys(editCredentials).length > 0 ? editCredentials : undefined,
                    is_testing: isSuperAdmin ? isTesting : undefined,
                });

                if (!response || response.success === false) {
                    throw new Error(response?.message || 'Error al actualizar integracion');
                }
                showToast('Integracion Jumpseller actualizada', 'success');
                onSuccess?.();
                return;
            }

            if (isSuperAdmin && !selectedBusinessId) {
                setErrorModal('Debes seleccionar un negocio antes de crear la integracion.');
                setLoading(false);
                return;
            }

            const credentials: JumpsellerCredentials = {
                api_key: formData.api_key,
                api_secret: formData.api_secret,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `jumpseller_${Date.now()}`,
                integration_type_id: JUMPSELLER_TYPE_ID,
                category: 'ecommerce',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
                is_testing: isTesting,
            });

            if (response.success) {
                showToast('Integracion Jumpseller creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al guardar la integracion de Jumpseller');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-3 w-full" autoComplete="off">
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
                                alt="Jumpseller"
                                className="h-8 w-8 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <ShoppingBagIcon className="h-6 w-6 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-base font-bold text-gray-900 dark:text-white leading-tight">Jumpseller</h2>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-0.5">
                            Conecta tu tienda Jumpseller para sincronizar ordenes, productos e inventario con Probability.
                        </p>
                    </div>
                </div>
            </div>

            <SectionCard icon={<Cog6ToothIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Configuracion General">
                <div className="space-y-3">
                    <div>
                        <label className={fieldLabel}>
                            Nombre de la Integracion <span style={{ color: GREEN }}>*</span>
                        </label>
                        <input
                            type="text"
                            required
                            placeholder="Ej: Jumpseller Tienda Principal"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Nombre descriptivo para identificar esta integracion</span>
                        </p>
                    </div>

                    {isSuperAdmin && !isEdit && (
                        <div>
                            <label className={fieldLabel}>
                                Negocio <span style={{ color: GREEN }}>*</span>
                            </label>
                            {loadingBusinesses ? (
                                <div className="flex items-center gap-2 rounded-lg px-3 py-2 bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                                    <Spinner className="animate-spin h-4 w-4 text-gray-400" />
                                    <span className="text-sm text-gray-600 dark:text-gray-300">Cargando negocios...</span>
                                </div>
                            ) : (
                                <select
                                    value={selectedBusinessId?.toString() || ''}
                                    onChange={(e) => setSelectedBusinessId(Number(e.target.value))}
                                    required
                                    className={inputCls}
                                    style={{ borderColor: INPUT_BORDER }}
                                >
                                    <option value="">-- Selecciona un negocio --</option>
                                    {businesses.map((business) => (
                                        <option key={business.id} value={business.id.toString()}>{business.name}</option>
                                    ))}
                                </select>
                            )}
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Selecciona el negocio al que pertenecera esta integracion</span>
                            </p>
                        </div>
                    )}
                </div>
            </SectionCard>

            <SectionCard icon={<KeyIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Credenciales de Acceso">
                <div className="space-y-3">
                    <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                        <div>
                            <label className={fieldLabel}>
                                Login <span style={{ color: GREEN }}>*</span>
                            </label>
                            <SecretInput
                                value={formData.api_key}
                                onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                                placeholder="Login de la API"
                                required={!isEdit}
                                autoComplete="off"
                                data-1p-ignore
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Lo encuentras en Jumpseller: Cuenta &gt; API</span>
                            </p>
                        </div>
                        <div>
                            <label className={fieldLabel}>
                                Auth Token <span style={{ color: GREEN }}>*</span>
                            </label>
                            <SecretInput
                                value={formData.api_secret}
                                onChange={(e) => setFormData({ ...formData, api_secret: e.target.value })}
                                placeholder="Auth Token de la API"
                                required={!isEdit}
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Lo encuentras en Jumpseller: Cuenta &gt; API</span>
                            </p>
                        </div>
                    </div>

                    <button
                        type="button"
                        onClick={handleTestConnection}
                        disabled={testingConnection || loading || !formData.api_key || !formData.api_secret}
                        className="w-full flex items-center justify-center gap-2 rounded-lg px-4 py-2.5 text-[13px] font-semibold bg-white dark:bg-gray-800 disabled:opacity-50"
                        style={{ border: `1px solid ${GREEN_BORDER}`, color: GREEN_DARK }}
                    >
                        {testingConnection ? (
                            <>
                                <Spinner className="animate-spin h-4 w-4" />
                                Probando...
                            </>
                        ) : (
                            <>
                                <CheckBadgeIcon className="w-4 h-4" />
                                Probar Conexion
                            </>
                        )}
                    </button>

                    <div className="rounded-lg p-3" style={{ backgroundColor: GREEN_SOFT, border: `1px solid ${GREEN_BORDER}` }}>
                        <h4 className="text-[13px] font-semibold text-gray-900 dark:text-white mb-2 flex items-center gap-2">
                            <InformationCircleIcon className="w-4 h-4" style={{ color: GREEN }} />
                            Como obtener tus credenciales
                        </h4>
                        <ol className="text-[11px] text-gray-600 dark:text-gray-300 space-y-1 list-decimal list-inside ml-1">
                            <li>Ingresa al panel de administracion de tu tienda en <strong>jumpseller.com</strong></li>
                            <li>En el menu lateral ve a <strong>Cuenta</strong> y luego a <strong>Preferencias</strong></li>
                            <li>Busca el recuadro <strong>API y MCP</strong> (arriba a la derecha)</li>
                            <li>Copia el valor de <strong>Login</strong></li>
                            <li>Copia el valor de <strong>Auth Token</strong></li>
                        </ol>

                        <div className="mt-3 pt-2" style={{ borderTop: `1px solid ${GREEN_BORDER}` }}>
                            <button
                                type="button"
                                onClick={() => setShowHelpImages((v) => !v)}
                                className="flex w-full items-center justify-between rounded-lg px-2 py-1.5 text-[12px] font-semibold text-gray-700 dark:text-gray-200 hover:bg-white/60 dark:hover:bg-gray-800/60 transition-colors"
                            >
                                <span className="flex items-center gap-1.5">
                                    <PhotoIcon className="w-4 h-4" style={{ color: GREEN_DARK }} />
                                    Ver imagen de ayuda paso a paso
                                </span>
                                <ChevronDownIcon
                                    className={`w-4 h-4 text-gray-400 transition-transform ${showHelpImages ? 'rotate-180' : ''}`}
                                />
                            </button>

                            {showHelpImages && (
                                <div className="mt-3 grid grid-cols-1 gap-4">
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
                </div>
            </SectionCard>

            <SectionCard icon={<ArrowsRightLeftIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Sincronizacion">
                <div className="rounded-lg bg-white dark:bg-gray-800 divide-y divide-gray-100 dark:divide-gray-700" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                    <ToggleRow
                        icon={<ArrowPathIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="Sincronizar inventario hacia Jumpseller"
                        subtitle="Envia el stock de Probability a los productos de tu tienda Jumpseller"
                        checked={inventorySyncEnabled}
                        onToggle={() => setInventorySyncEnabled(!inventorySyncEnabled)}
                    />
                    <ToggleRow
                        icon={<ArrowsRightLeftIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="Sincronizar estados hacia Jumpseller"
                        subtitle="Actualiza el estado de las ordenes en Jumpseller cuando cambian en Probability"
                        checked={statusSyncEnabled}
                        onToggle={() => setStatusSyncEnabled(!statusSyncEnabled)}
                    />
                </div>

                <JumpsellerInventorySection
                    value={{ ...inventorySync, enabled: inventorySyncEnabled }}
                    onChange={setInventorySync}
                    businessId={selectedBusinessId}
                    locationsInfo={locationsInfo}
                />

                {isEdit && integrationId && (
                    <>
                        {!inventorySyncSaved && (
                            <p className="mt-3 text-[11px] text-amber-600 dark:text-amber-500">
                                {inventorySyncEnabled
                                    ? 'Activaste el toggle de inventario pero aun no guardaste. Guarda la integracion para poder sincronizar ahora.'
                                    : 'Activa el toggle de inventario y guarda la integracion para poder sincronizar ahora.'}
                            </p>
                        )}
                        <button
                            type="button"
                            onClick={handleSyncInventory}
                            disabled={!inventorySyncSaved}
                            className="mt-3 w-full inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors disabled:opacity-60"
                            style={{ backgroundColor: GREEN }}
                            onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                            onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                        >
                            <ArrowPathIcon className="w-3.5 h-3.5" />
                            Sincronizar inventario ahora
                        </button>

                        <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
                            <h4 className="text-[12px] font-bold text-gray-900 dark:text-gray-100">Productos</h4>
                            <p className="mt-1 text-[11px] text-gray-500 dark:text-gray-400">
                                Cruza los productos por SKU; crea en Jumpseller o en Probability los que falten y asocia los que coinciden.
                            </p>
                            <button
                                type="button"
                                onClick={() => setProductSyncOpen(true)}
                                className="mt-3 w-full inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors"
                                style={{ backgroundColor: GREEN }}
                                onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                                onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                            >
                                <ArrowPathIcon className="w-3.5 h-3.5" />
                                Sincronizar productos
                            </button>
                        </div>

                        <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
                            <h4 className="text-[12px] font-bold text-gray-900 dark:text-gray-100">Ordenes</h4>
                            <p className="mt-1 text-[11px] text-gray-500 dark:text-gray-400">
                                Trae las ordenes de Jumpseller del periodo elegido. Los estados de Jumpseller (Paid, Canceled, Pending Payment) se traducen automaticamente a los de Probability (En Procesamiento, Cancelada, Pendiente).
                            </p>

                            <div className="mt-3 grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                                <div>
                                    <label className={fieldLabel}>Desde</label>
                                    <input
                                        type="date"
                                        value={ordersFrom}
                                        onChange={(e) => setOrdersFrom(e.target.value)}
                                        className={inputCls}
                                        style={{ borderColor: INPUT_BORDER }}
                                    />
                                </div>
                                <div>
                                    <label className={fieldLabel}>Hasta</label>
                                    <input
                                        type="date"
                                        value={ordersTo}
                                        onChange={(e) => setOrdersTo(e.target.value)}
                                        className={inputCls}
                                        style={{ borderColor: INPUT_BORDER }}
                                    />
                                </div>
                            </div>

                            <button
                                type="button"
                                onClick={handleSyncOrders}
                                className="mt-3 w-full inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors"
                                style={{ backgroundColor: GREEN }}
                                onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                                onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                            >
                                <ArrowPathIcon className="w-3.5 h-3.5" />
                                Sincronizar ordenes
                            </button>
                        </div>
                    </>
                )}
            </SectionCard>

            <SectionCard icon={<BeakerIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Modo de Pruebas">
                <div className="rounded-lg bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                    <ToggleRow
                        icon={<BeakerIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="Activar modo pruebas"
                        subtitle="Apunta al simulador interno de Jumpseller: puedes usar credenciales ficticias"
                        checked={isTesting}
                        onToggle={() => setIsTesting(!isTesting)}
                    />
                    {isTesting && integrationTypeBaseURLTest && (
                        <p className="px-3 pb-2.5 -mt-1 text-[11px] font-mono text-orange-700 dark:text-orange-400 break-all">
                            Sandbox: {integrationTypeBaseURLTest}
                        </p>
                    )}
                </div>
            </SectionCard>

            {isEdit && integrationId && (
                <SectionCard icon={<BoltIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Webhooks">
                    <JumpsellerWebhookManager integrationId={integrationId} />
                </SectionCard>
            )}

            {isEdit && integrationId && (
                <>
                    <JumpsellerProductSyncModal
                        isOpen={productSyncOpen}
                        onClose={() => setProductSyncOpen(false)}
                        integrationId={integrationId}
                        businessId={initialData?.business_id ?? null}
                    />
                    <JumpsellerOrderSyncModal
                        isOpen={orderSyncOpen}
                        onClose={() => setOrderSyncOpen(false)}
                        integrationId={integrationId}
                        businessId={initialData?.business_id ?? null}
                        createdAtMin={ordersFrom}
                        createdAtMax={ordersTo}
                    />
                    <JumpsellerInventorySyncModal
                        isOpen={inventorySyncOpen}
                        onClose={() => setInventorySyncOpen(false)}
                        integrationId={integrationId}
                        businessId={initialData?.business_id ?? null}
                    />
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
                    disabled={loading}
                    className="px-5 py-2 text-[13px] font-semibold rounded-lg text-white flex items-center justify-center gap-2 transition-colors disabled:opacity-60"
                    style={{ backgroundColor: GREEN }}
                    onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                    onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                >
                    {loading ? (
                        <>
                            <Spinner className="animate-spin h-4 w-4 text-white" />
                            {isEdit ? 'Guardando...' : 'Conectando...'}
                        </>
                    ) : (
                        <>
                            <CheckBadgeIcon className="w-4 h-4" />
                            {isEdit ? 'Guardar Integracion' : 'Crear Integracion'}
                        </>
                    )}
                </button>
            </div>

            {errorModal && (
                <Modal isOpen={!!errorModal} onClose={() => setErrorModal(null)} title="Error" size="sm">
                    <div className="p-4">
                        <Alert type="error">{errorModal}</Alert>
                    </div>
                </Modal>
            )}
        </form>
    );
}
