'use client';

import { useState, useEffect } from 'react';
import {
    KeyIcon,
    BeakerIcon,
    InformationCircleIcon,
    LinkIcon,
    GlobeAltIcon,
    ClipboardDocumentIcon,
    ClipboardDocumentCheckIcon,
} from '@heroicons/react/24/outline';
import { SecretInput } from '@/shared/ui';

export interface JumpsellerPlatformCredentials {
    client_id: string;
    client_secret: string;
    scopes: string;
    redirect_uri: string;
    test_client_id: string;
    test_client_secret: string;
    test_redirect_uri: string;
}

interface JumpsellerTypeCredentialsFormProps {
    credentials: JumpsellerPlatformCredentials;
    onChange: (credentials: JumpsellerPlatformCredentials) => void;
    isEditing?: boolean;
}

const fieldLabel = 'block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1';
const fieldHint = 'text-[11px] text-gray-400 dark:text-gray-500 mt-1';
const inputCls = 'w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)] font-mono';
const INPUT_BORDER = '#e9e9f0';

const DEFAULT_SCOPES = 'read_store read_orders write_orders read_products write_products read_customers write_hooks';

export default function JumpsellerTypeCredentialsForm({
    credentials,
    onChange,
    isEditing = false,
}: JumpsellerTypeCredentialsFormProps) {
    const set = (patch: Partial<JumpsellerPlatformCredentials>) => onChange({ ...credentials, ...patch });
    const placeholderSecret = isEditing ? 'Dejar vacio para mantener actual' : 'APP SECRET de la aplicacion';

    const [callbackURL, setCallbackURL] = useState('');
    const [copied, setCopied] = useState(false);

    useEffect(() => {
        if (!credentials.scopes) {
            set({ scopes: DEFAULT_SCOPES });
        }
    }, []);

    useEffect(() => {
        const apiBase = (process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1').replace(/\/$/, '');
        const url = apiBase.startsWith('http')
            ? `${apiBase}/jumpseller/callback`
            : `${window.location.origin}${apiBase}/jumpseller/callback`;
        setCallbackURL(url);
    }, []);

    const handleCopy = async () => {
        if (!callbackURL) return;
        await navigator.clipboard.writeText(callbackURL);
        setCopied(true);
        setTimeout(() => setCopied(false), 2000);
    };

    return (
        <div className="space-y-4">
            <div className="flex items-center gap-2">
                <KeyIcon className="w-5 h-5 text-[var(--color-primary)]" />
                <h3 className="text-sm font-bold text-gray-900 dark:text-white">Credenciales de la Aplicacion Jumpseller</h3>
            </div>

            <p className="text-xs text-gray-600 dark:text-gray-300 flex items-start gap-1.5">
                <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                <span>
                    Registra una App en el panel de Jumpseller (menu Apps) y pega aqui su APP ID y APP SECRET.
                    Todos los negocios conectan con esta misma App via OAuth. Se guardan encriptadas (AES-256-GCM).
                </span>
            </p>

            <div className="rounded-xl border border-blue-200 bg-blue-50/60 dark:bg-blue-950/20 dark:border-blue-800 p-4 space-y-2">
                <div className="flex items-center gap-2">
                    <LinkIcon className="w-4 h-4 text-blue-700" />
                    <span className="text-xs font-bold text-blue-900 dark:text-blue-200">Redirigir URL para registrar en tu App de Jumpseller</span>
                </div>
                <p className="text-[11px] text-gray-600 dark:text-gray-300">
                    Copia esta URL y pegala en el campo <strong>Redirigir URL</strong> de tu App en Jumpseller. Debe coincidir exactamente.
                </p>
                <div className="flex items-stretch gap-2">
                    <input
                        type="text"
                        readOnly
                        value={callbackURL}
                        onFocus={(e) => e.currentTarget.select()}
                        className={`${inputCls} flex-1`}
                        style={{ borderColor: INPUT_BORDER }}
                    />
                    <button
                        type="button"
                        onClick={handleCopy}
                        className="px-3 py-2 text-[13px] font-semibold rounded-lg flex items-center gap-1.5 shrink-0 text-blue-800 dark:text-blue-200 bg-blue-100 dark:bg-blue-900/40 border border-blue-200 dark:border-blue-800 hover:bg-blue-200 transition-colors"
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
            </div>

            <div className="rounded-xl border border-emerald-200 bg-emerald-50/50 dark:bg-emerald-950/20 dark:border-emerald-800 p-4">
                <div className="flex items-center gap-2 mb-3">
                    <GlobeAltIcon className="w-4 h-4 text-emerald-700" />
                    <span className="text-xs font-bold uppercase tracking-wide text-emerald-900 dark:text-emerald-200">Produccion</span>
                </div>
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>APP ID (Client ID)</label>
                        <input
                            type="text"
                            value={credentials.client_id}
                            onChange={(e) => set({ client_id: e.target.value })}
                            placeholder="APP ID de tu aplicacion"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>APP SECRET (Client Secret)</label>
                        <SecretInput
                            value={credentials.client_secret}
                            onChange={(e) => set({ client_secret: e.target.value })}
                            placeholder={placeholderSecret}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>
                    <div className="md:col-span-2">
                        <label className={fieldLabel}>Redirigir URL registrada</label>
                        <input
                            type="text"
                            value={credentials.redirect_uri}
                            onChange={(e) => set({ redirect_uri: e.target.value })}
                            placeholder="Pega aqui la misma Redirigir URL de arriba"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>Opcional. Si vacio, se deriva del dominio actual. Debe ser identica a la de tu App.</p>
                    </div>
                    <div className="md:col-span-2">
                        <label className={fieldLabel}>Scopes</label>
                        <input
                            type="text"
                            value={credentials.scopes}
                            onChange={(e) => set({ scopes: e.target.value })}
                            placeholder={DEFAULT_SCOPES}
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>Opcional. Separados por espacio. Default: {DEFAULT_SCOPES}</p>
                    </div>
                </div>
            </div>

            <div className="rounded-xl border border-amber-200 bg-amber-50/50 dark:bg-amber-950/20 dark:border-amber-800 p-4">
                <div className="flex items-center gap-2 mb-3">
                    <BeakerIcon className="w-4 h-4 text-amber-600" />
                    <span className="text-xs font-bold uppercase tracking-wide text-amber-900 dark:text-amber-200">Pruebas (opcional)</span>
                </div>
                <div className="grid grid-cols-1 gap-x-4 gap-y-3 md:grid-cols-2">
                    <div>
                        <label className={fieldLabel}>APP ID de prueba</label>
                        <input
                            type="text"
                            value={credentials.test_client_id}
                            onChange={(e) => set({ test_client_id: e.target.value })}
                            placeholder="APP ID de la App de desarrollo"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                    </div>
                    <div>
                        <label className={fieldLabel}>APP SECRET de prueba</label>
                        <SecretInput
                            value={credentials.test_client_secret}
                            onChange={(e) => set({ test_client_secret: e.target.value })}
                            placeholder={placeholderSecret}
                            className="bg-white dark:bg-gray-800 font-mono text-sm rounded-xl"
                        />
                    </div>
                    <div className="md:col-span-2">
                        <label className={fieldLabel}>Redirigir URL de prueba</label>
                        <input
                            type="text"
                            value={credentials.test_redirect_uri}
                            onChange={(e) => set({ test_redirect_uri: e.target.value })}
                            placeholder="http://localhost:3050/api/v1/jumpseller/callback"
                            className={inputCls}
                            style={{ borderColor: INPUT_BORDER }}
                        />
                        <p className={fieldHint}>Opcional. Si vacio, se deriva del host actual.</p>
                    </div>
                </div>
            </div>
        </div>
    );
}
