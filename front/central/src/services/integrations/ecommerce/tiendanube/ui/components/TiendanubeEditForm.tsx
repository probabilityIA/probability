'use client';

import { useState, FormEvent, useEffect, useMemo } from 'react';
import dynamic from 'next/dynamic';
import { Alert, Modal, SecretInput } from '@/shared/ui';
import type { Integration } from '@/services/integrations/core/domain/types';
import { TiendanubeCredentials, TiendanubeConfig } from '../../domain/types';
import { updateIntegrationAction, testConnectionRawAction, getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import { TiendanubeWebhookManager } from './TiendanubeWebhookManager';
import { TiendanubeInventorySection, TiendanubeInventoryConfig } from './TiendanubeInventorySection';
import { TiendanubeProductSyncModal } from './TiendanubeProductSyncModal';
import { TiendanubeOrderSyncModal } from './TiendanubeOrderSyncModal';
import { TiendanubeInventorySyncModal } from './TiendanubeInventorySyncModal';
import {
    ChannelStatusSyncSection,
    readChannelStatusSyncConfig,
    writeChannelStatusSyncConfig,
    type ChannelStatusSyncConfig,
} from '@/services/integrations/core/ui/components/ChannelStatusSyncSection';
import {
    ArchiveBoxIcon,
    ArrowPathIcon,
    BoltIcon,
    KeyIcon,
    Cog6ToothIcon,
    ShoppingBagIcon,
    InformationCircleIcon,
    CheckBadgeIcon,
    ScaleIcon,
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
    Spinner,
} from '@/services/integrations/invoicing/siigo/ui/components/SiigoFormKit';

const TIENDANUBE_TYPE_ID = 17;

const InventoryCompareStandalone = dynamic(
    () => import('@/services/modules/my-integrations/ui/components/InventoryCompareStandalone')
        .then((m) => m.InventoryCompareStandalone),
    { ssr: false },
);

interface TiendanubeEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: any;
        credentials?: TiendanubeCredentials;
        business_id?: number | null;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function TiendanubeEditForm({ integrationId, initialData, onSuccess, onCancel }: TiendanubeEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId] = useState<number | null>(initialData.business_id || null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData.name,
        store_id: initialData.config?.store_id || '',
        access_token: initialData.credentials?.access_token || '',
    });

    const [productSyncOpen, setProductSyncOpen] = useState(false);
    const [orderSyncOpen, setOrderSyncOpen] = useState(false);
    const [ordersFrom, setOrdersFrom] = useState('');
    const [ordersTo, setOrdersTo] = useState('');

    const handleSyncOrders = () => {
        if (!integrationId) return;
        if (ordersFrom && ordersTo && ordersFrom > ordersTo) {
            showToast('La fecha inicial no puede ser mayor que la fecha final', 'warning');
            return;
        }
        setOrderSyncOpen(true);
    };

    const [inventory, setInventory] = useState<TiendanubeInventoryConfig>({
        enabled: initialData.config?.inventory_sync_enabled === true,
        single_warehouse_id: Number(initialData.config?.inventory_single_warehouse_id) || 0,
    });

    const [statusSync, setStatusSync] = useState<ChannelStatusSyncConfig>(readChannelStatusSyncConfig(initialData.config));
    const [inventorySyncOpen, setInventorySyncOpen] = useState(false);
    const [inventoryCompareOpen, setInventoryCompareOpen] = useState(false);

    const compareIntegration = useMemo(() => ({
        id: integrationId,
        name: formData.name || 'Tiendanube',
        code: 'tiendanube',
        integration_type_id: TIENDANUBE_TYPE_ID,
        type: 'tiendanube',
        category: 'ecommerce',
        business_id: selectedBusinessId,
        store_id: formData.store_id,
        is_active: true,
        is_default: false,
        is_testing: false,
        config: initialData.config || {},
        created_by_id: 0,
        updated_by_id: null,
        created_at: '',
        updated_at: '',
        integration_type: {
            id: TIENDANUBE_TYPE_ID,
            name: 'Tiendanube',
            code: 'tiendanube',
            image_url: logoUrl || undefined,
        },
    }) as Integration, [integrationId, formData.name, formData.store_id, selectedBusinessId, initialData.config, logoUrl]);

    useEffect(() => {
        let cancelled = false;
        getActiveIntegrationTypesAction()
            .then((res: any) => {
                if (cancelled) return;
                const types = res?.data || [];
                const tiendanube = types.find((t: any) => t.id === TIENDANUBE_TYPE_ID || /tiendanube/i.test(t.code || ''));
                if (tiendanube?.image_url) setLogoUrl(tiendanube.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

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
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    const handleTestConnection = async () => {
        if (!formData.access_token) {
            showToast('Debes ingresar el Access Token para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials = {
                access_token: formData.access_token,
            };

            const config: TiendanubeConfig = {
                store_id: formData.store_id || undefined,
            };

            const result = await testConnectionRawAction('tiendanube', config, credentials);

            if (result.success) {
                showToast('Conexion exitosa con Tiendanube', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con Tiendanube');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const config: TiendanubeConfig = {
                ...(initialData.config || {}),
                store_id: formData.store_id || undefined,
                inventory_sync_enabled: inventory.enabled,
                inventory_single_warehouse_id: inventory.single_warehouse_id || undefined,
                ...writeChannelStatusSyncConfig(statusSync),
            };

            const updateData: any = {
                name: formData.name,
                store_id: formData.store_id || undefined,
                config: config,
            };

            if (formData.access_token) {
                const credentials: TiendanubeCredentials = {
                    access_token: formData.access_token,
                };
                updateData.credentials = credentials;
            }

            const response = await updateIntegrationAction(integrationId, updateData);

            if (response.success) {
                showToast('Integracion Tiendanube actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al actualizar integracion de Tiendanube');
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
                                alt="Tiendanube"
                                className="h-8 w-8 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <ShoppingBagIcon className="h-6 w-6 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-base font-bold text-gray-900 dark:text-white leading-tight">Editar Tiendanube</h2>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-0.5">
                            Actualiza la configuracion de tu integracion con Tiendanube.
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
                            placeholder="Ej: Tiendanube Tienda Principal"
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

                    {isSuperAdmin && (
                        <div>
                            <label className={fieldLabel}>Negocio</label>
                            {loadingBusinesses ? (
                                <div className="flex items-center gap-2 rounded-lg px-3 py-2 bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                                    <Spinner className="animate-spin h-4 w-4 text-gray-400" />
                                    <span className="text-sm text-gray-600 dark:text-gray-300">Cargando negocios...</span>
                                </div>
                            ) : (
                                <select
                                    value={selectedBusinessId?.toString() || ''}
                                    onChange={() => { }}
                                    disabled
                                    className={`${inputCls} opacity-70 cursor-not-allowed`}
                                    style={{ borderColor: INPUT_BORDER }}
                                >
                                    <option value="">-- Sin negocio asignado --</option>
                                    {businesses.map((business) => (
                                        <option key={business.id} value={business.id.toString()}>{business.name}</option>
                                    ))}
                                </select>
                            )}
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>El negocio no puede ser modificado despues de la creacion</span>
                            </p>
                        </div>
                    )}
                </div>
            </SectionCard>

            <SectionCard icon={<KeyIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Credenciales de Acceso">
                <div className="space-y-3">
                    <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                        <div>
                            <label className={fieldLabel}>Store ID</label>
                            <input
                                type="text"
                                value={formData.store_id}
                                onChange={(e) => setFormData({ ...formData, store_id: e.target.value })}
                                placeholder="Ej: 1234567"
                                autoComplete="off"
                                data-1p-ignore
                                className={inputCls}
                                style={{ borderColor: INPUT_BORDER }}
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>ID de tu tienda en Tiendanube (opcional)</span>
                            </p>
                        </div>
                        <div>
                            <label className={fieldLabel}>Access Token</label>
                            <SecretInput
                                value={formData.access_token}
                                onChange={(e) => setFormData({ ...formData, access_token: e.target.value })}
                                placeholder="Access Token de Tiendanube"
                                autoComplete="off"
                                data-1p-ignore
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Dejalo como esta si no necesitas cambiarlo</span>
                            </p>
                        </div>
                    </div>

                    <button
                        type="button"
                        onClick={handleTestConnection}
                        disabled={testingConnection || loading || !formData.access_token}
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
                            Sobre estas credenciales
                        </h4>
                        <p className="text-[11px] text-gray-600 dark:text-gray-300">
                            El Access Token y el Store ID los entrega Tiendanube al autorizar la app, no se generan a
                            mano. Solo editalos si necesitas reemplazarlos por unos obtenidos de otra forma; para
                            reconectar la tienda es preferible volver a ejecutar el flujo de autorizacion.
                        </p>
                    </div>
                </div>
            </SectionCard>

            <SectionCard icon={<ArrowPathIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Sincronizacion">
                <div>
                    <h4 className="text-[12px] font-bold text-gray-900 dark:text-gray-100">Productos</h4>
                    <p className="mt-1 text-[11px] text-gray-500 dark:text-gray-400">
                        Compara el catalogo de Probability con el de tu tienda y decide que crear o actualizar en cada lado.
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
                    <h4 className="text-[12px] font-bold text-gray-900 dark:text-gray-100">Estados de las ordenes</h4>
                    <p className="mt-1 mb-3 text-[11px] text-gray-500 dark:text-gray-400">
                        Define en que direccion viajan los cambios de estado entre Probability y tu tienda.
                    </p>
                    <ChannelStatusSyncSection
                        channelName="Tiendanube"
                        value={statusSync}
                        onChange={setStatusSync}
                        accentColor={GREEN}
                    />
                </div>

                <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
                    <h4 className="text-[12px] font-bold text-gray-900 dark:text-gray-100">Ordenes</h4>
                    <p className="mt-1 text-[11px] text-gray-500 dark:text-gray-400">
                        Trae las ordenes de Tiendanube del periodo elegido. Sus estados (open, closed, cancelled) y el de
                        pago se traducen automaticamente a los de Probability.
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
            </SectionCard>

            <TiendanubeProductSyncModal
                isOpen={productSyncOpen}
                onClose={() => setProductSyncOpen(false)}
                integrationId={integrationId}
                businessId={selectedBusinessId}
            />
            <TiendanubeOrderSyncModal
                isOpen={orderSyncOpen}
                onClose={() => setOrderSyncOpen(false)}
                integrationId={integrationId}
                businessId={selectedBusinessId}
                createdAtMin={ordersFrom}
                createdAtMax={ordersTo}
            />

            <SectionCard icon={<ArchiveBoxIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Inventario">
                <TiendanubeInventorySection
                    value={inventory}
                    onChange={setInventory}
                    businessId={selectedBusinessId}
                />

                <div className="mt-4 pt-4 border-t border-gray-100 dark:border-gray-700">
                    <h4 className="text-[12px] font-bold text-gray-900 dark:text-gray-100">Sincronizacion masiva</h4>
                    <p className="mt-1 text-[11px] text-gray-500 dark:text-gray-400">
                        Envia el stock de Probability a todas las variantes que compartan SKU con tu tienda. Compara
                        primero si quieres revisar que cambiaria antes de enviarlo.
                    </p>

                    <div className="mt-3 grid grid-cols-1 gap-2 sm:grid-cols-2">
                        <button
                            type="button"
                            onClick={() => setInventoryCompareOpen(true)}
                            disabled={!inventory.enabled}
                            title={inventory.enabled
                                ? 'Ver que cambiaria antes de enviar el stock'
                                : 'Activa la sincronizacion de inventario para poder comparar'}
                            className="inline-flex items-center justify-center gap-1.5 rounded-lg border py-2 text-[12px] font-semibold transition-colors disabled:cursor-not-allowed disabled:opacity-40"
                            style={{ borderColor: GREEN_BORDER, backgroundColor: GREEN_SOFT, color: GREEN_DARK }}
                        >
                            <ScaleIcon className="w-3.5 h-3.5" />
                            Comparar inventario
                        </button>
                        <button
                            type="button"
                            onClick={() => setInventorySyncOpen(true)}
                            disabled={!inventory.enabled}
                            title={inventory.enabled
                                ? 'Enviar el stock de todos los productos asociados'
                                : 'Activa la sincronizacion de inventario para poder enviarlo'}
                            className="inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors disabled:cursor-not-allowed disabled:opacity-40"
                            style={{ backgroundColor: GREEN }}
                            onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                            onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                        >
                            <ArrowPathIcon className="w-3.5 h-3.5" />
                            Sincronizar inventario
                        </button>
                    </div>

                    {!inventory.enabled && (
                        <p className="mt-2 flex items-start gap-1 text-[11px] text-amber-600 dark:text-amber-400">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Activa el interruptor de arriba y guarda para habilitar el envio masivo de stock.</span>
                        </p>
                    )}
                </div>
            </SectionCard>

            <TiendanubeInventorySyncModal
                isOpen={inventorySyncOpen}
                onClose={() => setInventorySyncOpen(false)}
                integrationId={integrationId}
                businessId={selectedBusinessId}
            />

            {inventoryCompareOpen && (
                <InventoryCompareStandalone
                    integration={compareIntegration}
                    businessId={selectedBusinessId}
                    onClose={() => setInventoryCompareOpen(false)}
                />
            )}

            <SectionCard icon={<BoltIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Webhooks">
                <TiendanubeWebhookManager integrationId={integrationId} />
            </SectionCard>

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
                            Actualizando...
                        </>
                    ) : (
                        <>
                            <CheckBadgeIcon className="w-4 h-4" />
                            Actualizar Integracion
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
