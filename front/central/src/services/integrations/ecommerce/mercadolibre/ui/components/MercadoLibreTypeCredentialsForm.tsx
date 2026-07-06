'use client';

import { useState, useEffect } from 'react';
import {
    KeyIcon,
    GlobeAltIcon,
    BeakerIcon,
    InformationCircleIcon,
    LinkIcon,
    ClipboardDocumentIcon,
    ClipboardDocumentCheckIcon,
} from '@heroicons/react/24/outline';
import { SecretInput } from '@/shared/ui';

export interface MercadoLibrePlatformCredentials {
    client_id: string;
    client_secret: string;
    auth_domain: string;
    test_client_id: string;
    test_client_secret: string;
    test_auth_domain: string;
}

interface MercadoLibreTypeCredentialsFormProps {
    credentials: MercadoLibrePlatformCredentials;
    onChange: (credentials: MercadoLibrePlatformCredentials) => void;
    isEditing?: boolean;
}

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)] font-mono';
const INPUT_BORDER = '#e9e9f0';

export default function MercadoLibreTypeCredentialsForm({
    credentials,
    onChange,
    isEditing = false,
}: MercadoLibreTypeCredentialsFormProps) {
    const set = (patch: Partial<MercadoLibrePlatformCredentials>) => onChange({ ...credentials, ...patch });
    const placeholderSecret = isEditing ? 'Dejar vacio para mantener actual' : 'Secret Key de la aplicacion';

    const [urls, setUrls] = useState({ redirect: '', notifications: '' });
    const [copiedKey, setCopiedKey] = useState<'redirect' | 'notifications' | null>(null);

    useEffect(() => {
        const apiBase = (process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1').replace(/\/$/, '');
        const build = (path: string) => apiBase.startsWith('http')
            ? `${apiBase}${path}`
            : `${window.location.origin}${apiBase}${path}`;
        setUrls({
            redirect: build('/meli/callback'),
            notifications: build('/meli/notifications'),
        });
    }, []);

    const handleCopy = async (key: 'redirect' | 'notifications', value: string) => {
        if (!value) return;
        await navigator.clipboard.writeText(value);
        setCopiedKey(key);
        setTimeout(() => setCopiedKey(null), 2000);
    };

    return (
        <div className="space-y-4">
            <div className="flex items-center gap-2">
                <KeyIcon className="w-5 h-5 text-[var(--color-primary)]" />
                <h3 className="text-sm font-bold text-gray-900 dark:text-white">Credenciales de Plataforma MercadoLibre</h3>
            </div>

            <p className="text-xs text-gray-600 dark:text-gray-300 flex items-start gap-1.5">
                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                <span>
                    Configura los dos juegos de credenciales OAuth (produccion y sandbox). El backend usa uno u otro
                    segun el flag <code>is_testing</code> de cada integracion. Se guardan encriptadas (AES-256-GCM).
                </span>
            </p>

            <div className="rounded-xl border border-blue-200 bg-blue-50/60 dark:bg-blue-950/20 dark:border-blue-800 p-4 space-y-3">
                <div className="flex items-center gap-2">
                    <LinkIcon className="w-4 h-4 text-blue-700" />
                    <span className="text-xs font-bold text-blue-900 dark:text-blue-200">URLs para registrar en tu app de MercadoLibre</span>
                </div>
                <p className="text-[11px] text-gray-600 dark:text-gray-300">
                    El sistema genera estas URLs automaticamente. Copialas y pegalas en <strong>developers.mercadolibre.com</strong>. Son las mismas para produccion y sandbox.
                </p>

                <div>
                    <label className="text-[12px] font-semibold text-blue-900 dark:text-blue-200">Redirect URI</label>
                    <p className="text-[11px] text-gray-500 dark:text-gray-400 mb-1">Configuracion y scopes &rarr; Redirect URIs</p>
                    <div className="flex items-stretch gap-2">
                        <input
                            type="text"
                            readOnly
                            value={urls.redirect}
                            onFocus={(e) => e.currentTarget.select()}
                            className={`${inputCls} flex-1`}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <button
                            type="button"
                            onClick={() => handleCopy('redirect', urls.redirect)}
                            className="px-3 py-2 text-[13px] font-semibold rounded-lg flex items-center gap-1.5 shrink-0 text-blue-800 dark:text-blue-200 bg-blue-100 dark:bg-blue-900/40 border border-blue-200 dark:border-blue-800 hover:bg-blue-200 transition-colors"
                        >
                            {copiedKey === 'redirect' ? (
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
                </div>

                <div>
                    <label className="text-[12px] font-semibold text-blue-900 dark:text-blue-200">Notificaciones callback URL</label>
                    <p className="text-[11px] text-gray-500 dark:text-gray-400 mb-1">Configuracion de notificaciones &rarr; Notificaciones callbacks URL</p>
                    <div className="flex items-stretch gap-2">
                        <input
                            type="text"
                            readOnly
                            value={urls.notifications}
                            onFocus={(e) => e.currentTarget.select()}
                            className={`${inputCls} flex-1`}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <button
                            type="button"
                            onClick={() => handleCopy('notifications', urls.notifications)}
                            className="px-3 py-2 text-[13px] font-semibold rounded-lg flex items-center gap-1.5 shrink-0 text-blue-800 dark:text-blue-200 bg-blue-100 dark:bg-blue-900/40 border border-blue-200 dark:border-blue-800 hover:bg-blue-200 transition-colors"
                        >
                            {copiedKey === 'notifications' ? (
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
                </div>

                <p className="text-[11px] text-amber-700 dark:text-amber-400">
                    MercadoLibre exige HTTPS. En local usa un tunel (ngrok) y registra esas URLs.
                </p>
            </div>

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
                            placeholder="Ej: 1234567890123456"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>Secret Key (Client Secret)</label>
                        <SecretInput
                            value={credentials.client_secret}
                            onChange={(e) => set({ client_secret: e.target.value })}
                            placeholder={placeholderSecret}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>Auth Domain</label>
                        <input
                            type="text"
                            value={credentials.auth_domain}
                            onChange={(e) => set({ auth_domain: e.target.value })}
                            placeholder="auth.mercadolibre.com.co"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>Opcional. Default: auth.mercadolibre.com.co</p>
                    </div>
                </div>
            </div>

            <div className="rounded-xl border border-amber-200 bg-amber-50/50 dark:bg-amber-950/20 dark:border-amber-800 p-4">
                <div className="flex items-center gap-2 mb-3">
                    <BeakerIcon className="w-4 h-4 text-amber-600" />
                    <span className="text-xs font-bold uppercase tracking-wide text-amber-900 dark:text-amber-200">Sandbox / Pruebas</span>
                </div>
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>App ID (Client ID) de prueba</label>
                        <input
                            type="text"
                            value={credentials.test_client_id}
                            onChange={(e) => set({ test_client_id: e.target.value })}
                            placeholder="Ej: 9876543210987654"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>Secret Key de prueba</label>
                        <SecretInput
                            value={credentials.test_client_secret}
                            onChange={(e) => set({ test_client_secret: e.target.value })}
                            placeholder={placeholderSecret}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>Auth Domain de prueba</label>
                        <input
                            type="text"
                            value={credentials.test_auth_domain}
                            onChange={(e) => set({ test_auth_domain: e.target.value })}
                            placeholder="auth.mercadolibre.com.co"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>Opcional. Si vacio, usa el de produccion.</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
