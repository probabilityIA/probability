'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Alert, Modal, SecretInput, CodIncludesShippingToggle } from '@/shared/ui';
import { TiendanubeCredentials, TiendanubeConfig } from '../../domain/types';
import { createIntegrationAction, testConnectionRawAction, getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    ShoppingBagIcon,
    InformationCircleIcon,
    CheckBadgeIcon,
    BeakerIcon,
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

const TIENDANUBE_TYPE_ID = 17;

interface TiendanubeConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    integrationTypeBaseURLTest?: string;
}

export function TiendanubeConfigForm({ onSuccess, onCancel, integrationTypeBaseURLTest }: TiendanubeConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [codIncludesShipping, setCodIncludesShipping] = useState<boolean>(true);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [isTesting, setIsTesting] = useState(false);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);

    const [formData, setFormData] = useState({
        name: '',
        store_id: '',
        access_token: '',
    });

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
            } else {
                if (permissions?.business_id) {
                    setSelectedBusinessId(permissions.business_id);
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
                is_testing: isTesting,
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
            if (isSuperAdmin && !selectedBusinessId) {
                setErrorModal('Debes seleccionar un negocio antes de crear la integracion.');
                setLoading(false);
                return;
            }

            const credentials: TiendanubeCredentials = {
                access_token: formData.access_token,
            };

            const config: TiendanubeConfig = {
                cod_includes_shipping: codIncludesShipping,
                store_id: formData.store_id || undefined,
                is_testing: isTesting,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `tiendanube_${Date.now()}`,
                integration_type_id: TIENDANUBE_TYPE_ID,
                category: 'ecommerce',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
                is_testing: isTesting,
            });

            if (response.success) {
                showToast('Integracion Tiendanube creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al crear la integracion de Tiendanube');
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
                        <h2 className="text-base font-bold text-gray-900 dark:text-white leading-tight">Tiendanube</h2>
                        <p className="text-xs text-gray-600 dark:text-gray-300 mt-0.5">
                            Conecta tu tienda Tiendanube para sincronizar ordenes y productos automaticamente con Probability.
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
                            <label className={fieldLabel}>
                                Access Token <span style={{ color: GREEN }}>*</span>
                            </label>
                            <SecretInput
                                value={formData.access_token}
                                onChange={(e) => setFormData({ ...formData, access_token: e.target.value })}
                                placeholder="Access Token de Tiendanube"
                                required
                                autoComplete="off"
                                data-1p-ignore
                                className="w-full bg-white dark:bg-gray-800 font-mono text-sm rounded-lg"
                            />
                            <p className={fieldHint}>
                                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                <span>Lo encuentras en Tiendanube: Aplicaciones &gt; Credenciales de API</span>
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
                            Como obtener tus credenciales
                        </h4>
                        <ol className="text-[11px] text-gray-600 dark:text-gray-300 space-y-1 list-decimal list-inside ml-1">
                            <li>Ingresa al panel de administracion de tu tienda en <strong>Tiendanube</strong></li>
                            <li>Ve a <strong>Aplicaciones</strong> y luego a <strong>Credenciales de API</strong></li>
                            <li>Genera un nuevo <strong>Access Token</strong> y copialo</li>
                            <li>Copia el <strong>Store ID</strong> de tu tienda (opcional)</li>
                        </ol>
                    </div>
                </div>
            </SectionCard>

            <SectionCard icon={<BeakerIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Modo de Pruebas">
                <div className="rounded-lg bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                    <ToggleRow
                        icon={<BeakerIcon className="w-4 h-4" style={{ color: GREEN }} />}
                        title="Activar modo pruebas"
                        subtitle="Apunta al simulador interno de Tiendanube: puedes usar credenciales ficticias"
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

            <CodIncludesShippingToggle checked={codIncludesShipping} onToggle={() => setCodIncludesShipping((v) => !v)} />

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
                            Conectando...
                        </>
                    ) : (
                        <>
                            <CheckBadgeIcon className="w-4 h-4" />
                            Crear Integracion
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
