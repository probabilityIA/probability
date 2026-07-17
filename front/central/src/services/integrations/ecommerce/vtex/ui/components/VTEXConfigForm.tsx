'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Alert, Modal, SecretInput } from '@/shared/ui';
import { VTEXCredentials, VTEXConfig } from '../../domain/types';
import { createIntegrationAction, updateIntegrationAction, testConnectionRawAction, getActiveIntegrationTypesAction, syncOrdersAction } from '@/services/integrations/core/infra/actions';
import { VTEXWarehouseSection, VTEXWarehouseMapping, VTEXWarehousesInfo } from './VTEXWarehouseSection';
import { VTEXWebhookManager } from './VTEXWebhookManager';
import { getVTEXWarehousesAction, syncVTEXProductsAction, syncVTEXInventoryAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    CheckBadgeIcon,
    InformationCircleIcon,
    ArrowPathIcon,
    ArrowsRightLeftIcon,
    ShoppingBagIcon,
    BuildingStorefrontIcon,
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

const VTEX_TYPE_ID = 16;

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

const VTEX_SUFFIXES = [
    '.myvtex.com',
    '.vtexcommercestable.com.br',
    '.vtexcommercebeta.com.br',
    '.vtexlocal.com.br',
    '.vtexcommerce.com.br',
    '.vtex.com',
];

export const cleanAccountName = (raw: string): string => {
    let account = raw.trim().toLowerCase();
    account = account.replace(/^https?:\/\//, '');
    account = account.split('/')[0];
    account = account.split('?')[0];
    account = account.split(':')[0];
    for (const suffix of VTEX_SUFFIXES) {
        if (account.endsWith(suffix)) {
            account = account.slice(0, -suffix.length);
            break;
        }
    }
    return account.split('.')[0];
};

interface VTEXConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    isEdit?: boolean;
    integrationId?: number;
    initialData?: {
        name?: string;
        config?: any;
        credentials?: any;
        business_id?: number | null;
    };
}

export function VTEXConfigForm({ onSuccess, onCancel, isEdit, integrationId, initialData }: VTEXConfigFormProps) {
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

    const [inventorySyncEnabled, setInventorySyncEnabled] = useState<boolean>(!!initialData?.config?.inventory_sync_enabled);
    const [statusSyncEnabled, setStatusSyncEnabled] = useState<boolean>(!!initialData?.config?.status_sync_enabled);
    const [isSeller, setIsSeller] = useState<boolean>(!!initialData?.config?.is_seller);

    const [warehouseMappings, setWarehouseMappings] = useState<VTEXWarehouseMapping[]>(() => {
        const raw = initialData?.config?.vtex_warehouse_mappings;
        if (!Array.isArray(raw)) return [];
        return raw.map((m: any) => ({
            internal_warehouse_id: Number(m.internal_warehouse_id) || 0,
            vtex_warehouse_id: m.vtex_warehouse_id ? String(m.vtex_warehouse_id) : '',
        }));
    });
    const [warehousesInfo, setWarehousesInfo] = useState<VTEXWarehousesInfo | null>(null);
    const [syncingProducts, setSyncingProducts] = useState(false);
    const [syncingInventory, setSyncingInventory] = useState(false);
    const [syncingOrders, setSyncingOrders] = useState(false);

    const inventorySyncSaved = !!initialData?.config?.inventory_sync_enabled;

    const [ordersFrom, setOrdersFrom] = useState<string>(defaultOrdersFrom());
    const [ordersTo, setOrdersTo] = useState<string>(defaultOrdersTo());

    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        account_name: initialData?.config?.account_name || '',
        app_key: initialData?.credentials?.app_key || '',
        app_token: initialData?.credentials?.app_token || '',
    });

    useEffect(() => {
        let cancelled = false;
        getActiveIntegrationTypesAction()
            .then((res: any) => {
                if (cancelled) return;
                const types = res?.data || [];
                const vtex = types.find((t: any) => t.id === VTEX_TYPE_ID || /vtex/i.test(t.code || ''));
                if (vtex?.image_url) setLogoUrl(vtex.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

    useEffect(() => {
        if (!inventorySyncEnabled || !isEdit || !integrationId) return;
        let cancelled = false;
        getVTEXWarehousesAction(integrationId, selectedBusinessId ?? undefined)
            .then((res: any) => {
                if (cancelled || !res?.success) return;
                setWarehousesInfo({ warehouses: Array.isArray(res.warehouses) ? res.warehouses : [] });
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
        if (!formData.app_key || !formData.app_token) {
            showToast('Debes ingresar App Key y App Token para probar la conexion', 'warning');
            return;
        }

        if (!formData.account_name) {
            showToast('Debes ingresar el nombre de la cuenta VTEX', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials: VTEXCredentials = {
                app_key: formData.app_key,
                app_token: formData.app_token,
            };

            const config: VTEXConfig = {
                account_name: cleanAccountName(formData.account_name),
            };

            const result = await testConnectionRawAction('vtex', config as any, credentials as any);

            if (result.success) {
                showToast('Conexion exitosa con VTEX', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con VTEX');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSyncProducts = async () => {
        if (!integrationId) return;
        setSyncingProducts(true);
        try {
            const res: any = await syncVTEXProductsAction(integrationId, selectedBusinessId ?? undefined);
            if (res?.success) {
                showToast('Sincronizacion de productos iniciada', 'success');
            } else {
                setErrorModal(res?.message || 'Error al sincronizar productos');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al sincronizar productos');
        } finally {
            setSyncingProducts(false);
        }
    };

    const handleSyncInventory = async () => {
        if (!integrationId) return;
        setSyncingInventory(true);
        try {
            const res: any = await syncVTEXInventoryAction(integrationId, selectedBusinessId ?? undefined);
            if (res?.success) {
                showToast('Sincronizacion de inventario iniciada', 'success');
            } else {
                setErrorModal(res?.message || 'Error al sincronizar inventario');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al sincronizar inventario');
        } finally {
            setSyncingInventory(false);
        }
    };

    const handleSyncOrders = async () => {
        if (!integrationId) return;

        if (ordersFrom && ordersTo && ordersFrom > ordersTo) {
            showToast('La fecha inicial no puede ser mayor que la fecha final', 'warning');
            return;
        }

        setSyncingOrders(true);
        try {
            const res: any = await syncOrdersAction(integrationId, {
                created_at_min: `${ordersFrom}T00:00:00.000Z`,
                created_at_max: `${ordersTo}T23:59:59.999Z`,
            });
            if (res?.success !== false) {
                showToast('Sincronizacion de ordenes iniciada', 'success');
            } else {
                setErrorModal(res?.message || 'Error al sincronizar ordenes');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al sincronizar ordenes');
        } finally {
            setSyncingOrders(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const cleanMappings = warehouseMappings
                .filter((m) => m.internal_warehouse_id > 0 && m.vtex_warehouse_id.trim() !== '')
                .map((m) => ({
                    internal_warehouse_id: m.internal_warehouse_id,
                    vtex_warehouse_id: m.vtex_warehouse_id.trim(),
                }));

            const config = {
                account_name: cleanAccountName(formData.account_name),
                is_seller: isSeller,
                inventory_sync_enabled: inventorySyncEnabled,
                status_sync_enabled: statusSyncEnabled,
                vtex_warehouse_mappings: cleanMappings,
            };

            if (isEdit && integrationId) {
                const editCredentials: any = {};
                if (formData.app_key) editCredentials.app_key = formData.app_key;
                if (formData.app_token) editCredentials.app_token = formData.app_token;

                const response: any = await updateIntegrationAction(integrationId, {
                    name: formData.name,
                    config: config as any,
                    credentials: Object.keys(editCredentials).length > 0 ? editCredentials : undefined,
                });

                if (!response || response.success === false) {
                    throw new Error(response?.message || 'Error al actualizar integracion');
                }
                showToast('Integracion VTEX actualizada', 'success');
                onSuccess?.();
                return;
            }

            if (isSuperAdmin && !selectedBusinessId) {
                setErrorModal('Debes seleccionar un negocio antes de crear la integracion.');
                setLoading(false);
                return;
            }

            const credentials: VTEXCredentials = {
                app_key: formData.app_key,
                app_token: formData.app_token,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `vtex_${Date.now()}`,
                integration_type_id: VTEX_TYPE_ID,
                category: 'ecommerce',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
            });

            if (response.success) {
                showToast('Integracion VTEX creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al guardar la integracion de VTEX');
        } finally {
            setLoading(false);
        }
    };

    const accountPreview = formData.account_name ? cleanAccountName(formData.account_name) : '';

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
                                alt="VTEX"
                                className="h-8 w-8 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <ShoppingBagIcon className="h-6 w-6 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-base font-bold text-gray-900 dark:text-white leading-tight">VTEX</h2>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-0.5">
                            Conecta tu tienda VTEX para sincronizar ordenes, productos e inventario con Probability.
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
                            placeholder="Ej: VTEX Tienda Principal"
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

                    <div>
                        <label className={fieldLabel}>
                            Nombre de la Cuenta VTEX <span style={{ color: GREEN }}>*</span>
                        </label>
                        <input
                            type="text"
                            required
                            placeholder="Ej: mitienda"
                            value={formData.account_name}
                            onChange={(e) => setFormData({ ...formData, account_name: e.target.value })}
                            autoComplete="off"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>
                                {accountPreview
                                    ? `Se conectara a ${accountPreview}.vtexcommercestable.com.br`
                                    : 'El nombre de tu cuenta VTEX, tal como aparece en la URL del Admin'}
                            </span>
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

                    {isEdit && (
                        <p className={fieldHint}>
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>El negocio no puede ser modificado despues de la creacion</span>
                        </p>
                    )}
                </div>
            </SectionCard>

            <SectionCard icon={<KeyIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Credenciales de Acceso">
                <div className="space-y-3">
                    <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                        <div>
                            <label className={fieldLabel}>
                                App Key <span style={{ color: GREEN }}>*</span>
                            </label>
                            <SecretInput
                                value={formData.app_key}
                                onChange={(e) => setFormData({ ...formData, app_key: e.target.value })}
                                placeholder="X-VTEX-API-AppKey"
                                required={!isEdit}
                                autoComplete="off"
                                data-1p-ignore
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Lo encuentras en VTEX: Configuracion de la cuenta &gt; Claves de API</span>
                            </p>
                        </div>
                        <div>
                            <label className={fieldLabel}>
                                App Token <span style={{ color: GREEN }}>*</span>
                            </label>
                            <SecretInput
                                value={formData.app_token}
                                onChange={(e) => setFormData({ ...formData, app_token: e.target.value })}
                                placeholder="X-VTEX-API-AppToken"
                                required={!isEdit}
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Solo se muestra una vez al crear la clave en VTEX</span>
                            </p>
                        </div>
                    </div>

                    {isEdit && (
                        <p className="text-[11px] text-gray-500 dark:text-gray-400">
                            Deja los campos vacios para conservar las credenciales actuales.
                        </p>
                    )}

                    <button
                        type="button"
                        onClick={handleTestConnection}
                        disabled={testingConnection || loading || !formData.app_key || !formData.app_token}
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
                            <li>Ingresa al <strong>Admin de VTEX</strong> de tu tienda</li>
                            <li>Ve a <strong>Configuracion de la cuenta</strong> y luego a <strong>Claves de API</strong></li>
                            <li>Crea una nueva clave y copia el valor de <strong>App Key</strong></li>
                            <li>Copia el valor de <strong>App Token</strong> (solo se muestra una vez)</li>
                            <li>Asigna a la clave los roles de <strong>OMS</strong>, <strong>Catalogo</strong> y <strong>Logistica</strong></li>
                        </ol>
                    </div>
                </div>
            </SectionCard>

            <SectionCard icon={<ArrowsRightLeftIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Sincronizacion">
                <div className="rounded-lg bg-white dark:bg-gray-800 divide-y divide-gray-100 dark:divide-gray-700" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                    <ToggleRow
                        icon={<ArrowPathIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="Sincronizar inventario hacia VTEX"
                        subtitle="Envia el stock de Probability a los SKUs de tu tienda VTEX"
                        checked={inventorySyncEnabled}
                        onToggle={() => setInventorySyncEnabled(!inventorySyncEnabled)}
                    />
                    <ToggleRow
                        icon={<ArrowsRightLeftIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="Sincronizar estados hacia VTEX"
                        subtitle="Actualiza el estado de las ordenes en VTEX cuando cambian en Probability"
                        checked={statusSyncEnabled}
                        onToggle={() => setStatusSyncEnabled(!statusSyncEnabled)}
                    />
                    <ToggleRow
                        icon={<BuildingStorefrontIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="La cuenta es un seller"
                        subtitle="Activalo si vendes como seller en un marketplace VTEX y no como tienda propia"
                        checked={isSeller}
                        onToggle={() => setIsSeller(!isSeller)}
                    />
                </div>

                <VTEXWarehouseSection
                    enabled={inventorySyncEnabled}
                    value={warehouseMappings}
                    onChange={setWarehouseMappings}
                    businessId={selectedBusinessId}
                    warehousesInfo={warehousesInfo}
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
                            disabled={!inventorySyncSaved || syncingInventory}
                            className="mt-3 w-full inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors disabled:opacity-60"
                            style={{ backgroundColor: GREEN }}
                            onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                            onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                        >
                            {syncingInventory ? <Spinner className="animate-spin h-3.5 w-3.5" /> : <ArrowPathIcon className="w-3.5 h-3.5" />}
                            Sincronizar inventario ahora
                        </button>
                    </>
                )}
            </SectionCard>

            {isEdit && integrationId && (
                <SectionCard icon={<ArrowPathIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Sincronizar productos">
                    <p className="text-[11px] text-gray-500 dark:text-gray-400">
                        Cruza los productos por SKU contra el RefId de VTEX: crea en Probability los que falten y asocia los que coinciden.
                    </p>
                    <button
                        type="button"
                        onClick={handleSyncProducts}
                        disabled={syncingProducts}
                        className="mt-3 w-full inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors disabled:opacity-60"
                        style={{ backgroundColor: GREEN }}
                        onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                        onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                    >
                        {syncingProducts ? <Spinner className="animate-spin h-3.5 w-3.5" /> : <ArrowPathIcon className="w-3.5 h-3.5" />}
                        Sincronizar productos
                    </button>
                </SectionCard>
            )}

            {isEdit && integrationId && (
                <SectionCard icon={<ShoppingBagIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Sincronizar ordenes">
                    <p className="text-[11px] text-gray-500 dark:text-gray-400">
                        Trae las ordenes de VTEX del periodo elegido. Los estados de VTEX (ready-for-handling, payment-pending, canceled) se traducen automaticamente a los de Probability.
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
                        disabled={syncingOrders}
                        className="mt-3 w-full inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors disabled:opacity-60"
                        style={{ backgroundColor: GREEN }}
                        onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                        onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                    >
                        {syncingOrders ? <Spinner className="animate-spin h-3.5 w-3.5" /> : <ArrowPathIcon className="w-3.5 h-3.5" />}
                        Sincronizar ordenes
                    </button>
                </SectionCard>
            )}

            {isEdit && integrationId && (
                <SectionCard icon={<BoltIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Webhooks">
                    <VTEXWebhookManager integrationId={integrationId} businessId={selectedBusinessId} />
                </SectionCard>
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
