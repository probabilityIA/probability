'use client';

import { useEffect, useMemo, useState } from 'react';
import {
    KeyIcon,
    BeakerIcon,
    InformationCircleIcon,
    LinkIcon,
    GlobeAltIcon,
    ClipboardDocumentIcon,
    ClipboardDocumentCheckIcon,
    BoltIcon,
    ShieldCheckIcon,
    LockClosedIcon,
    ArrowTopRightOnSquareIcon,
} from '@heroicons/react/24/outline';
import { SecretInput } from '@/shared/ui';

export interface TiendanubePlatformCredentials {
    client_id: string;
    client_secret: string;
    redirect_uri: string;
    app_url: string;
    user_agent: string;
    test_client_id: string;
    test_client_secret: string;
    test_redirect_uri: string;
    test_store_id: string;
}

export const EMPTY_TIENDANUBE_PLATFORM_CREDENTIALS: TiendanubePlatformCredentials = {
    client_id: '',
    client_secret: '',
    redirect_uri: '',
    app_url: '',
    user_agent: '',
    test_client_id: '',
    test_client_secret: '',
    test_redirect_uri: '',
    test_store_id: '',
};

interface TiendanubeTypeCredentialsFormProps {
    credentials: TiendanubePlatformCredentials;
    onChange: (credentials: TiendanubePlatformCredentials) => void;
    isEditing?: boolean;
    webhookUrls?: Record<string, string>;
}

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1';
const inputCls =
    'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)] font-mono';
const INPUT_BORDER = '#e9e9f0';

const PARTNERS_URL = 'https://partners.tiendanube.com/applications';
const PRIVACY_DOCS_URL = 'https://tiendanube.github.io/api-documentation/resources/webhook';

const WEBHOOK_TOPICS = [
    'order/created',
    'order/paid',
    'order/updated',
    'order/cancelled',
    'product/updated',
    'app/uninstalled',
];

const REQUIRED_SCOPES = [
    'read_products',
    'write_products',
    'read_orders',
    'write_orders',
    'read_customers',
    'read_shipping',
];

function CopyField({ value, disabled }: { value: string; disabled?: boolean }) {
    const [copied, setCopied] = useState(false);

    const handleCopy = async () => {
        if (!value) return;
        await navigator.clipboard.writeText(value);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="flex items-stretch gap-2">
            <input
                type="text"
                readOnly
                value={value}
                placeholder={disabled ? 'Configura WEBHOOK_BASE_URL en el backend' : ''}
                onFocus={(e) => e.currentTarget.select()}
                className={`${inputCls} flex-1 text-gray-700 dark:text-gray-200`}
                style={{ borderColor: INPUT_BORDER }}
            />
            <button
                type="button"
                onClick={handleCopy}
                disabled={!value}
                className="px-3 py-2 text-[13px] font-semibold rounded-lg flex items-center gap-1.5 shrink-0 text-blue-800 dark:text-blue-200 bg-blue-100 dark:bg-blue-900/40 border border-blue-200 dark:border-blue-800 hover:bg-blue-200 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
            >
                {copied ? (
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
    );
}

export default function TiendanubeTypeCredentialsForm({
    credentials,
    onChange,
    isEditing = false,
    webhookUrls,
}: TiendanubeTypeCredentialsFormProps) {
    const set = (patch: Partial<TiendanubePlatformCredentials>) => onChange({ ...credentials, ...patch });
    const placeholderSecret = isEditing ? 'Dejar vacio para mantener actual' : 'Client Secret de la aplicacion';

    const [apiOrigin, setApiOrigin] = useState('');

    useEffect(() => {
        const apiBase = (process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1').replace(/\/$/, '');
        setApiOrigin(apiBase.startsWith('http') ? apiBase : `${window.location.origin}${apiBase}`);
    }, []);

    const url = (key: string, path: string) => webhookUrls?.[key] || (apiOrigin ? `${apiOrigin}${path}` : '');

    const prodWebhook = url('production', '/tiendanube/webhook');
    const testWebhook = webhookUrls?.sandbox || '';
    const callbackURL = apiOrigin ? `${apiOrigin}/tiendanube/callback` : '';
    const storeRedactURL = url('store_redact', '/tiendanube/webhook/store-redact');
    const customersRedactURL = url('customers_redact', '/tiendanube/webhook/customers-redact');
    const customersDataURL = url('customers_data_request', '/tiendanube/webhook/customers-data-request');

    const authorizeURL = useMemo(() => {
        const appID = credentials.client_id.trim();
        if (!appID) return '';
        return `https://www.tiendanube.com/apps/${appID}/authorize`;
    }, [credentials.client_id]);

    return (
        <div className="space-y-4">
            <div className="flex flex-wrap items-center justify-between gap-2">
                <div className="flex items-center gap-2">
                    <KeyIcon className="w-5 h-5 text-[var(--color-primary)]" />
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Credenciales de la Aplicacion Tiendanube</h3>
                </div>
                <a
                    href={PARTNERS_URL}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-[12px] font-semibold text-[var(--color-primary)] hover:underline"
                >
                    Panel de Partners
                    <ArrowTopRightOnSquareIcon className="w-3.5 h-3.5" />
                </a>
            </div>

            <p className="text-xs text-gray-600 dark:text-gray-300 flex items-start gap-1.5">
                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                <span>
                    Crea una aplicacion en <strong>partners.tiendanube.com</strong> y pega aqui su <strong>App ID</strong> y{' '}
                    <strong>Client Secret</strong> (seccion &quot;Llaves de acceso&quot;). Todos los negocios instalan esta misma
                    aplicacion via OAuth y cada tienda queda con su propio access_token. Se guardan encriptadas (AES-256-GCM).
                </span>
            </p>

            <div className="rounded-xl border border-emerald-200 bg-emerald-50/50 dark:bg-emerald-950/20 dark:border-emerald-800 p-4">
                <div className="flex items-center gap-2 mb-3">
                    <GlobeAltIcon className="w-4 h-4 text-emerald-700" />
                    <span className="text-xs font-bold uppercase tracking-wide text-emerald-900 dark:text-emerald-200">Produccion</span>
                </div>
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>App ID (Client ID)</label>
                        <input
                            type="text"
                            value={credentials.client_id}
                            onChange={(e) => set({ client_id: e.target.value })}
                            placeholder="39928"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>Numero de la aplicacion en el panel de Partners.</p>
                    </div>
                    <div>
                        <label className={fieldLabel}>Client Secret</label>
                        <SecretInput
                            value={credentials.client_secret}
                            onChange={(e) => set({ client_secret: e.target.value })}
                            placeholder={placeholderSecret}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                        <p className={fieldHint}>Tambien firma los webhooks (HMAC SHA-256).</p>
                    </div>
                    <div>
                        <label className={fieldLabel}>URL de redireccion despues de la instalacion</label>
                        <input
                            type="text"
                            value={credentials.redirect_uri}
                            onChange={(e) => set({ redirect_uri: e.target.value })}
                            placeholder={callbackURL || 'https://tu-dominio/api/v1/tiendanube/callback'}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            Debe apuntar al callback de Probability y ser identica a la registrada en la pestana
                            Configuracion de tu App. Si en el Partner Portal queda la URL de partners.tiendanube.com,
                            al comerciante le aparece una pantalla con un comando curl en vez de volver a Probability.
                        </p>
                    </div>
                    <div>
                        <label className={fieldLabel}>Pagina de la aplicacion</label>
                        <input
                            type="text"
                            value={credentials.app_url}
                            onChange={(e) => set({ app_url: e.target.value })}
                            placeholder="https://probabilityia.com.co/integrations"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>Opcional. URL publica que Tiendanube muestra al comerciante.</p>
                    </div>
                    <div className="md:col-span-2">
                        <label className={fieldLabel}>User-Agent</label>
                        <input
                            type="text"
                            value={credentials.user_agent}
                            onChange={(e) => set({ user_agent: e.target.value })}
                            placeholder="Probability (soporte@probabilityia.com.co)"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>
                            Tiendanube exige identificar la app con nombre y correo de contacto en cada request. Sin esto el API
                            responde 400.
                        </p>
                    </div>
                </div>

                <div className="mt-4 pt-3 border-t border-emerald-200 dark:border-emerald-800">
                    <label className={fieldLabel}>URL de redireccion a registrar en el Partner Portal</label>
                    <CopyField value={callbackURL} disabled={!callbackURL} />
                    <p className={fieldHint}>
                        Copiala en Configuracion &gt; URL de redireccion de tu App en partners.tiendanube.com. Sin
                        esto el flujo OAuth no puede devolver el codigo a Probability.
                    </p>
                </div>

                {authorizeURL && (
                    <div className="mt-4 pt-3 border-t border-emerald-200 dark:border-emerald-800">
                        <label className={fieldLabel}>URL de autorizacion generada</label>
                        <CopyField value={authorizeURL} />
                        <p className={fieldHint}>A esta URL se envia al comerciante para instalar la app en su tienda.</p>
                    </div>
                )}
            </div>

            <div className="rounded-xl border border-blue-200 bg-blue-50/60 dark:bg-blue-950/20 dark:border-blue-800 p-4 space-y-3">
                <div className="flex items-center gap-2">
                    <BoltIcon className="w-4 h-4 text-blue-700" />
                    <span className="text-xs font-bold uppercase tracking-wide text-blue-900 dark:text-blue-200">Webhooks</span>
                </div>
                <p className="text-[11px] text-gray-600 dark:text-gray-300">
                    Registra esta URL como destino de los webhooks de tu App. Tiendanube firma cada envio con HMAC SHA-256 usando
                    el Client Secret y lo manda en el header <code>x-linkedstore-hmac-sha256</code>.
                </p>
                <div>
                    <label className={fieldLabel}>URL de webhook (produccion)</label>
                    <CopyField value={prodWebhook} disabled={!prodWebhook} />
                </div>
                {testWebhook && (
                    <div>
                        <label className={fieldLabel}>URL de webhook (pruebas)</label>
                        <CopyField value={testWebhook} />
                    </div>
                )}
                <div>
                    <span className="block text-[11px] font-semibold text-gray-600 dark:text-gray-300 mb-1.5">
                        Eventos que consume la plataforma
                    </span>
                    <div className="flex flex-wrap gap-1.5">
                        {WEBHOOK_TOPICS.map((topic) => (
                            <span
                                key={topic}
                                className="px-2 py-0.5 rounded-md text-[11px] font-mono bg-white dark:bg-gray-800 text-blue-800 dark:text-blue-200 border border-blue-200 dark:border-blue-800"
                            >
                                {topic}
                            </span>
                        ))}
                    </div>
                </div>
            </div>

            <div className="rounded-xl border border-purple-200 bg-purple-50/50 dark:bg-purple-950/20 dark:border-purple-800 p-4 space-y-3">
                <div className="flex flex-wrap items-center justify-between gap-2">
                    <div className="flex items-center gap-2">
                        <LockClosedIcon className="w-4 h-4 text-purple-700" />
                        <span className="text-xs font-bold uppercase tracking-wide text-purple-900 dark:text-purple-200">
                            Privacidad (obligatorio)
                        </span>
                    </div>
                    <a
                        href={PRIVACY_DOCS_URL}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 text-[12px] font-semibold text-purple-800 dark:text-purple-200 hover:underline"
                    >
                        Documentacion
                        <ArrowTopRightOnSquareIcon className="w-3.5 h-3.5" />
                    </a>
                </div>
                <p className="text-[11px] text-gray-600 dark:text-gray-300">
                    Tiendanube exige estas tres URLs para aprobar la aplicacion. Copialas en la seccion{' '}
                    <strong>Privacidad</strong> del panel de Partners.
                </p>
                <div className="space-y-3">
                    <div>
                        <label className={fieldLabel}>URL webhook store redact</label>
                        <CopyField value={storeRedactURL} disabled={!storeRedactURL} />
                        <p className={fieldHint}>La tienda desinstalo la app: se purgan sus datos.</p>
                    </div>
                    <div>
                        <label className={fieldLabel}>URL webhook customers redact</label>
                        <CopyField value={customersRedactURL} disabled={!customersRedactURL} />
                        <p className={fieldHint}>Un comprador pidio borrar sus datos personales.</p>
                    </div>
                    <div>
                        <label className={fieldLabel}>URL webhook customers data request</label>
                        <CopyField value={customersDataURL} disabled={!customersDataURL} />
                        <p className={fieldHint}>Un comprador pidio una copia de sus datos personales.</p>
                    </div>
                </div>
            </div>

            <div className="rounded-xl border border-gray-200 bg-gray-50/70 dark:bg-gray-800/40 dark:border-gray-700 p-4 space-y-2">
                <div className="flex items-center gap-2">
                    <ShieldCheckIcon className="w-4 h-4 text-gray-600 dark:text-gray-300" />
                    <span className="text-xs font-bold uppercase tracking-wide text-gray-700 dark:text-gray-200">
                        Permisos a habilitar en el panel
                    </span>
                </div>
                <p className="text-[11px] text-gray-500 dark:text-gray-400">
                    Tiendanube no recibe scopes por OAuth: se configuran en la seccion <strong>Permisos</strong> de la App. Estos
                    son los minimos para sincronizar productos, inventario y ordenes.
                </p>
                <div className="flex flex-wrap gap-1.5">
                    {REQUIRED_SCOPES.map((scope) => (
                        <span
                            key={scope}
                            className="px-2 py-0.5 rounded-md text-[11px] font-mono bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-gray-700"
                        >
                            {scope}
                        </span>
                    ))}
                </div>
            </div>

            <div className="rounded-xl border border-amber-200 bg-amber-50/50 dark:bg-amber-950/20 dark:border-amber-800 p-4">
                <div className="flex items-center gap-2 mb-3">
                    <BeakerIcon className="w-4 h-4 text-amber-600" />
                    <span className="text-xs font-bold uppercase tracking-wide text-amber-900 dark:text-amber-200">
                        Pruebas (opcional)
                    </span>
                </div>
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>App ID de prueba</label>
                        <input
                            type="text"
                            value={credentials.test_client_id}
                            onChange={(e) => set({ test_client_id: e.target.value })}
                            placeholder="App ID de la aplicacion en desarrollo"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>Client Secret de prueba</label>
                        <SecretInput
                            value={credentials.test_client_secret}
                            onChange={(e) => set({ test_client_secret: e.target.value })}
                            placeholder={placeholderSecret}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>URL de redireccion de prueba</label>
                        <input
                            type="text"
                            value={credentials.test_redirect_uri}
                            onChange={(e) => set({ test_redirect_uri: e.target.value })}
                            placeholder="http://localhost:3050/api/v1/tiendanube/callback"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>ID de la tienda demo</label>
                        <input
                            type="text"
                            value={credentials.test_store_id}
                            onChange={(e) => set({ test_store_id: e.target.value })}
                            placeholder="1234567"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>La tienda de prueba que creaste desde el panel de Partners.</p>
                    </div>
                </div>
                <p className="text-[11px] text-amber-800 dark:text-amber-300 mt-3 flex items-start gap-1.5">
                    <LinkIcon className="w-3.5 h-3.5 mt-0.5 flex-shrink-0" />
                    <span>
                        El modo pruebas usa la <strong>URL de Pruebas (Sandbox)</strong> definida arriba en URLs del API, no la de
                        produccion.
                    </span>
                </p>
            </div>
        </div>
    );
}
