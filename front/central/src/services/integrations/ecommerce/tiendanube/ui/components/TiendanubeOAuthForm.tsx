'use client';

import { useState, useEffect } from 'react';
import { Modal, Alert, CodIncludesShippingToggle } from '@/shared/ui';
import { getActiveIntegrationTypesAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    Cog6ToothIcon,
    ShoppingBagIcon,
    InformationCircleIcon,
    LinkIcon,
    BeakerIcon,
    CheckCircleIcon,
    ShieldCheckIcon,
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

const GUIDE_STEPS = [
    'Haz clic en Conectar con Tiendanube',
    'Inicia sesion en tu tienda Tiendanube',
    'Autoriza el acceso a Probability',
    'Regresas automaticamente y queda conectada',
];

interface TiendanubeOAuthFormProps {
    onCancel?: () => void;
}

export function TiendanubeOAuthForm({ onCancel }: TiendanubeOAuthFormProps) {
    const { showToast } = useToast();
    const [connecting, setConnecting] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);
    const [logoUrl, setLogoUrl] = useState<string | null>(null);
    const [logoFailed, setLogoFailed] = useState(false);
    const [name, setName] = useState('');
    const [isTesting, setIsTesting] = useState(false);
    const [codIncludesShipping, setCodIncludesShipping] = useState<boolean>(true);
    const [verifying, setVerifying] = useState(false);
    const [verifyResult, setVerifyResult] = useState<{ ok: boolean; message: string } | null>(null);

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
            } else if (permissions?.business_id) {
                setSelectedBusinessId(permissions.business_id);
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
                const tiendanube = types.find((t: any) => t.id === TIENDANUBE_TYPE_ID || /tiendanube/i.test(t.code || ''));
                if (tiendanube?.image_url) setLogoUrl(tiendanube.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

    const handleVerify = async () => {
        setVerifying(true);
        setVerifyResult(null);
        try {
            const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
            const response = await fetch(`${apiBaseUrl}/integrations/tiendanube/verify-app?is_testing=${isTesting}`, {
                headers: { 'Authorization': `Bearer ${TokenStorage.getSessionToken()}` },
                credentials: 'include',
            });
            const data = await response.json();
            setVerifyResult({ ok: !!data.configured, message: data.message || (data.configured ? 'App configurada' : 'App no configurada') });
        } catch (err: any) {
            setVerifyResult({ ok: false, message: err.message || 'No se pudo verificar la App' });
        } finally {
            setVerifying(false);
        }
    };

    const handleConnect = async () => {
        if (!name.trim()) {
            showToast('Ingresa un nombre para la integracion', 'warning');
            return;
        }
        if (isSuperAdmin && !selectedBusinessId) {
            setErrorModal('Debes seleccionar un negocio antes de conectar.');
            return;
        }

        setConnecting(true);
        try {
            sessionStorage.setItem(
                'tiendanube_pending_connection',
                JSON.stringify({ cod_includes_shipping: codIncludesShipping })
            );

            const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
            const response = await fetch(`${apiBaseUrl}/integrations/tiendanube/connect`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${TokenStorage.getSessionToken()}`,
                },
                credentials: 'include',
                body: JSON.stringify({
                    integration_name: name.trim(),
                    business_id: isSuperAdmin ? selectedBusinessId : 0,
                    is_testing: isTesting,
                }),
            });

            const data = await response.json();
            if (!response.ok || !data.success) {
                throw new Error(data.error || data.message || 'Error al iniciar la conexion OAuth');
            }
            if (!data.authorization_url) {
                throw new Error('No se recibio la URL de autorizacion');
            }
            window.location.href = data.authorization_url;
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con Tiendanube');
            setConnecting(false);
        }
    };

    return (
        <div className="space-y-3 w-full">
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
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Ej: Tiendanube Tienda Principal"
                            required
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

            <SectionCard icon={<LinkIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Conexion OAuth">
                <div className="space-y-3">
                    <p className="text-[12px] text-gray-500 dark:text-gray-400 leading-relaxed">
                        Te redirigimos a Tiendanube para que autorices el acceso. No necesitas pegar tokens ni
                        buscar el ID de la tienda: al volver, la integracion queda conectada automaticamente y las
                        credenciales se guardan de forma segura.
                    </p>

                    <ol className="grid grid-cols-2 gap-y-4 sm:grid-cols-4 sm:gap-y-0">
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

                    {verifyResult && (
                        <div
                            className="flex items-start gap-2 rounded-lg px-3 py-2 text-[12px] font-medium"
                            style={verifyResult.ok
                                ? { backgroundColor: '#ecfdf5', border: '1px solid #a7f3d0', color: '#047857' }
                                : { backgroundColor: '#fef2f2', border: '1px solid #fecaca', color: '#b91c1c' }}
                        >
                            {verifyResult.ok
                                ? <CheckCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                : <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />}
                            <span>{verifyResult.message}</span>
                        </div>
                    )}

                    <div className="flex flex-col gap-2 sm:flex-row">
                        <button
                            type="button"
                            onClick={handleVerify}
                            disabled={verifying || connecting}
                            className="flex items-center justify-center gap-2 rounded-lg py-2.5 px-4 text-[13px] font-semibold transition-colors disabled:opacity-60 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 sm:w-auto"
                            style={{ border: `1px solid ${INPUT_BORDER}` }}
                        >
                            {verifying ? (
                                <>
                                    <Spinner className="animate-spin h-4 w-4" />
                                    Verificando...
                                </>
                            ) : (
                                <>
                                    <ShieldCheckIcon className="w-4 h-4" />
                                    Verificar App
                                </>
                            )}
                        </button>

                        <button
                            type="button"
                            onClick={handleConnect}
                            disabled={connecting}
                            className="flex-1 flex items-center justify-center gap-2 rounded-lg py-2.5 text-[13px] font-semibold text-white transition-colors disabled:opacity-60"
                            style={{ backgroundColor: GREEN }}
                            onMouseEnter={(e) => { if (!connecting) (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                            onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                        >
                            {connecting ? (
                                <>
                                    <Spinner className="animate-spin h-4 w-4 text-white" />
                                    Redirigiendo a Tiendanube...
                                </>
                            ) : (
                                <>
                                    <LinkIcon className="w-4 h-4" />
                                    Conectar con Tiendanube
                                </>
                            )}
                        </button>
                    </div>
                </div>
            </SectionCard>

            {isSuperAdmin && (
                <SectionCard icon={<BeakerIcon style={{ color: GREEN, width: 16, height: 16 }} />} title="Modo de Pruebas">
                    <div className="rounded-lg bg-white dark:bg-gray-800" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                        <ToggleRow
                            icon={<BeakerIcon className="w-4 h-4" style={{ color: GREEN }} />}
                            title="Activar modo pruebas"
                            subtitle="Conecta usando la app de pruebas del tipo de integracion en vez de la de produccion"
                            checked={isTesting}
                            onToggle={() => setIsTesting(!isTesting)}
                        />
                    </div>
                </SectionCard>
            )}

            <CodIncludesShippingToggle checked={codIncludesShipping} onToggle={() => setCodIncludesShipping((v) => !v)} />

            <div className="flex flex-col-reverse gap-2.5 pt-3 border-t border-gray-100 dark:border-gray-700 sm:flex-row sm:justify-end sm:items-center">
                {onCancel && (
                    <button
                        type="button"
                        onClick={onCancel}
                        disabled={connecting}
                        className="px-5 py-2 text-[13px] font-semibold rounded-lg bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                        style={{ border: `1px solid ${INPUT_BORDER}` }}
                    >
                        Cancelar
                    </button>
                )}
            </div>

            {errorModal && (
                <Modal isOpen={!!errorModal} onClose={() => setErrorModal(null)} title="Error" size="sm">
                    <div className="p-4">
                        <Alert type="error">{errorModal}</Alert>
                    </div>
                </Modal>
            )}
        </div>
    );
}
