'use client';

import { useState, useEffect } from 'react';
import { Select, Modal, Alert } from '@/shared/ui';
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

interface JumpsellerOAuthFormProps {
    onCancel?: () => void;
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
    'Haz clic en Conectar con Jumpseller',
    'Inicia sesion en tu tienda Jumpseller',
    'Autoriza el acceso a Probability',
    'Regresas automaticamente y queda conectada',
];

export function JumpsellerOAuthForm({ onCancel }: JumpsellerOAuthFormProps) {
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
                const js = types.find((t: any) => t.id === 33 || /jumpseller/i.test(t.code || ''));
                if (js?.image_url) setLogoUrl(js.image_url);
            })
            .catch(() => { });
        return () => { cancelled = true; };
    }, []);

    const handleVerify = async () => {
        setVerifying(true);
        setVerifyResult(null);
        try {
            const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
            const response = await fetch(`${apiBaseUrl}/integrations/jumpseller/verify-app?is_testing=${isTesting}`, {
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
            const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
            const response = await fetch(`${apiBaseUrl}/integrations/jumpseller/connect`, {
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
            setErrorModal(err.message || 'Error al conectar con Jumpseller');
            setConnecting(false);
        }
    };

    return (
        <div className="space-y-3 w-full">
            <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-center gap-3">
                    <span
                        className="flex h-11 w-11 items-center justify-center rounded-xl overflow-hidden shrink-0"
                        style={{ backgroundColor: logoUrl && !logoFailed ? GREEN_SOFT : GREEN, border: `1px solid ${GREEN_BORDER}` }}
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
                        <h2 className="text-lg font-bold text-gray-900 dark:text-white leading-tight">Jumpseller</h2>
                        <p className="text-xs text-gray-500 dark:text-gray-300">
                            Conecta tu tienda con OAuth para sincronizar ordenes, stock y estados.
                        </p>
                    </div>
                </div>
                <span
                    className="inline-flex items-center gap-2 self-start rounded-full px-3 py-1 text-[11px] font-semibold"
                    style={{ backgroundColor: '#f3f4f6', border: '1px solid #e5e7eb', color: '#6b7280' }}
                >
                    <span className="h-2 w-2 rounded-full" style={{ backgroundColor: '#9ca3af' }} />
                    Sin conectar
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
                    <div className={isSuperAdmin ? '' : 'md:col-span-2'}>
                        <label className={fieldLabel}>
                            Nombre de la Integracion <span style={{ color: GREEN }}>*</span>
                        </label>
                        <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            placeholder="Ej: Jumpseller Principal"
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

            {isSuperAdmin && (
                <div
                    className="rounded-xl p-4 dark:bg-gray-800/60"
                    style={{ backgroundColor: '#ffffff', border: `1px solid ${CARD_BORDER}` }}
                >
                    <div className="flex items-center justify-between gap-3">
                        <div className="flex items-start gap-2">
                            <BeakerIcon className="w-4 h-4 mt-0.5" style={{ color: '#d97706' }} />
                            <div>
                                <h4 className="text-[13px] font-bold text-gray-900 dark:text-gray-100">Modo de pruebas</h4>
                                <p className="mt-0.5 text-[11px] text-gray-400 dark:text-gray-500">
                                    Conecta usando la app de pruebas (credenciales de test del tipo de integracion) en vez de la de produccion.
                                </p>
                            </div>
                        </div>
                        <button
                            type="button"
                            role="switch"
                            aria-checked={isTesting}
                            onClick={() => setIsTesting((v) => !v)}
                            className={`relative inline-flex h-6 w-11 flex-shrink-0 items-center rounded-full transition-colors ${isTesting ? 'bg-amber-500' : 'bg-gray-300 dark:bg-gray-600'}`}
                        >
                            <span className={`inline-block h-5 w-5 transform rounded-full bg-white shadow transition-transform ${isTesting ? 'translate-x-5' : 'translate-x-0.5'}`} />
                        </button>
                    </div>
                </div>
            )}

            <div
                className="rounded-xl p-4 dark:bg-gray-800/60"
                style={{ backgroundColor: CARD_BG, border: `1px solid ${CARD_BORDER}` }}
            >
                <div className="flex items-center gap-2 mb-3">
                    <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                        <LinkIcon style={{ color: GREEN, width: 16, height: 16 }} />
                    </span>
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Conexion OAuth</h3>
                </div>

                <p className="text-[12px] text-gray-500 dark:text-gray-400 leading-relaxed mb-3">
                    Te redirigimos a Jumpseller para que autorices el acceso. No necesitas pegar tokens: al
                    volver, la integracion queda conectada automaticamente y los tokens se guardan de forma segura.
                </p>

                <ol className="grid grid-cols-2 gap-y-4 sm:grid-cols-4 sm:gap-y-0 mb-4">
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
                        className="mb-3 flex items-start gap-2 rounded-lg px-3 py-2 text-[12px] font-medium"
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
                                <svg className="animate-spin h-4 w-4" style={{ color: GREEN }} xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
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
                                <svg className="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
                                Redirigiendo a Jumpseller...
                            </>
                        ) : (
                            <>
                                <LinkIcon className="w-4 h-4" />
                                Conectar con Jumpseller
                            </>
                        )}
                    </button>
                </div>
            </div>

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
        </div>
    );
}
