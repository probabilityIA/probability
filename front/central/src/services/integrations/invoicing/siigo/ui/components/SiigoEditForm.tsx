'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Alert, Modal, SecretInput } from '@/shared/ui';
import { SiigoCredentials } from '../../domain/types';
import { updateIntegrationAction, testConnectionRawAction, getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    CheckBadgeIcon,
    InformationCircleIcon,
    BeakerIcon,
    DocumentTextIcon,
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
} from './SiigoFormKit';
import { SiigoInventorySection, InventorySyncConfig } from './SiigoInventorySection';
import { SiigoInventorySyncModal } from './SiigoInventorySyncModal';
import { SiigoWebhookManager } from './SiigoWebhookManager';

function initialInventoryConfig(config: any): InventorySyncConfig {
    const cfg = config || {};
    const rawMappings = Array.isArray(cfg.inventory_warehouse_mappings) ? cfg.inventory_warehouse_mappings : [];
    return {
        enabled: cfg.inventory_sync_enabled === true,
        mode: cfg.inventory_warehouse_mode === 'mapped' ? 'mapped' : 'single',
        single_warehouse_id: Number(cfg.inventory_single_warehouse_id) || 0,
        mappings: rawMappings.map((m: any) => ({
            velocity_warehouse_id: Number(m.velocity_warehouse_id) || 0,
            siigo_warehouse_id: Number(m.siigo_warehouse_id) || 0,
        })),
        product_sync_enabled: cfg.product_sync_enabled === true,
    };
}

function inventoryConfigToMap(inv: InventorySyncConfig) {
    return {
        inventory_sync_enabled: inv.enabled,
        inventory_warehouse_mode: inv.mode,
        inventory_single_warehouse_id: inv.single_warehouse_id,
        inventory_warehouse_mappings: inv.mappings.filter((m) => m.velocity_warehouse_id > 0),
        product_sync_enabled: inv.product_sync_enabled,
    };
}

interface SiigoEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: any;
        credentials?: SiigoCredentials;
        business_id?: number | null;
        is_testing?: boolean;
        base_url_test?: string;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function SiigoEditForm({ integrationId, initialData, onSuccess, onCancel }: SiigoEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [isTesting, setIsTesting] = useState(initialData.is_testing || false);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId] = useState<number | null>(initialData.business_id || null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);

    const [inventorySync, setInventorySync] = useState<InventorySyncConfig>(initialInventoryConfig(initialData.config));
    const [showSyncModal, setShowSyncModal] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData.name,
        username: initialData.credentials?.username || '',
        access_key: initialData.credentials?.access_key || '',
        account_id: initialData.credentials?.account_id || '',
        partner_id: initialData.credentials?.partner_id || '',
    });

    const connected = !!(formData.username && formData.access_key);
    const canSyncNow = inventorySync.enabled;

    useEffect(() => {
        let cancelled = false;
        getActiveIntegrationTypesAction()
            .then((res: any) => {
                if (cancelled) return;
                const types = res?.data || [];
                const siigo = types.find((t: any) => t.id === 8 || /siigo/i.test(t.code || ''));
                if (siigo?.image_url) setLogoUrl(siigo.image_url);
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
        if (!formData.username || !formData.access_key) {
            showToast('Debes ingresar Usuario y Clave de Acceso para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials = {
                username: formData.username,
                access_key: formData.access_key,
                account_id: formData.account_id || undefined,
                partner_id: formData.partner_id,
            };

            const result = await testConnectionRawAction('siigo', { is_testing: isTesting }, credentials);

            if (result.success) {
                showToast('Conexion exitosa con Siigo', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con Siigo');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const updateData: any = {
                name: formData.name,
                config: {
                    ...(initialData.config || {}),
                    ...inventoryConfigToMap(inventorySync),
                },
                is_testing: isTesting,
            };

            if (formData.username && formData.access_key) {
                const credentials: SiigoCredentials = {
                    username: formData.username,
                    access_key: formData.access_key,
                    account_id: formData.account_id || undefined,
                    partner_id: formData.partner_id,
                };
                updateData.credentials = credentials;
            }

            const response = await updateIntegrationAction(integrationId, updateData);

            if (response.success) {
                showToast('Integracion Siigo actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al actualizar integracion de Siigo');
        } finally {
            setLoading(false);
        }
    };

    return (
        <>
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
                                alt="Siigo"
                                className="h-8 w-8 object-contain"
                                onError={() => setLogoFailed(true)}
                            />
                        ) : (
                            <DocumentTextIcon className="h-6 w-6 text-white" />
                        )}
                    </span>
                    <div>
                        <h2 className="text-base font-bold text-gray-900 dark:text-white leading-tight">Siigo Facturacion Electronica</h2>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-0.5">
                            Datos de tu cuenta Siigo. Puedes modificar las credenciales si es necesario.
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

            <SectionCard icon={<Cog6ToothIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Configuracion General">
                <div className="space-y-3">
                    <div>
                        <label className={fieldLabel}>
                            Nombre de la Integracion <span style={{ color: GREEN }}>*</span>
                        </label>
                        <input
                            type="text"
                            required
                            placeholder="Ej: Siigo Facturacion Principal"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
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
                <p className={`${fieldHint} mb-3`}>
                    <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                    <span>Se muestran las credenciales actuales. Modificalas solo si necesitas actualizarlas; todos los campos deben quedar completos.</span>
                </p>
                <div className="space-y-3">
                    <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                        <div>
                            <label className={fieldLabel}>Usuario API</label>
                            <input
                                type="text"
                                placeholder="usuario@empresa.com"
                                value={formData.username}
                                onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                                autoComplete="off"
                                data-1p-ignore
                                className={inputCls}
                                style={{ borderColor: INPUT_BORDER }}
                            />
                        </div>
                        <div>
                            <label className={fieldLabel}>Clave de Acceso (Access Key)</label>
                            <SecretInput
                                value={formData.access_key}
                                onChange={(e) => setFormData({ ...formData, access_key: e.target.value })}
                                placeholder="Clave de acceso API"
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                        </div>
                    </div>

                    <div className="grid grid-cols-1 gap-x-4 gap-y-3 sm:grid-cols-2">
                        <div>
                            <label className={fieldLabel}>Account ID</label>
                            <input
                                type="text"
                                placeholder="ID de cuenta/suscripcion"
                                value={formData.account_id}
                                onChange={(e) => setFormData({ ...formData, account_id: e.target.value })}
                                autoComplete="off"
                                data-1p-ignore
                                className={`${inputCls} font-mono`}
                                style={{ borderColor: INPUT_BORDER }}
                            />
                        </div>
                        <div>
                            <label className={fieldLabel}>Partner ID</label>
                            <input
                                type="text"
                                placeholder="Partner ID"
                                value={formData.partner_id}
                                onChange={(e) => setFormData({ ...formData, partner_id: e.target.value })}
                                autoComplete="off"
                                data-1p-ignore
                                className={`${inputCls} font-mono`}
                                style={{ borderColor: INPUT_BORDER }}
                            />
                        </div>
                    </div>

                    <button
                        type="button"
                        onClick={handleTestConnection}
                        disabled={testingConnection || loading || !formData.username || !formData.access_key}
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
                </div>
            </SectionCard>

            <SiigoInventorySection
                value={inventorySync}
                onChange={setInventorySync}
                businessId={selectedBusinessId}
                integrationId={integrationId}
                onSyncNow={() => setShowSyncModal(true)}
                canSyncNow={canSyncNow}
            />

            <SectionCard icon={<BoltIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Webhooks de Siigo">
                <SiigoWebhookManager integrationId={integrationId} />
            </SectionCard>

            <SectionCard icon={<BeakerIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Modo de Pruebas">
                <div className="rounded-lg bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                    <ToggleRow
                        icon={<BeakerIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="Activar modo testing"
                        subtitle="Las facturas quedaran marcadas como TEST y usaran la URL de pruebas de Siigo"
                        checked={isTesting}
                        onToggle={() => setIsTesting(!isTesting)}
                    />
                    {isTesting && initialData.base_url_test && (
                        <p className="px-3 pb-2.5 -mt-1 text-[11px] font-mono text-orange-700 dark:text-orange-400 break-all">
                            Sandbox: {initialData.base_url_test}
                        </p>
                    )}
                </div>
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
        <SiigoInventorySyncModal
            isOpen={showSyncModal}
            onClose={() => setShowSyncModal(false)}
            integrationId={integrationId}
            businessId={selectedBusinessId}
        />
        </>
    );
}
