'use client';

import { useState } from 'react';
import {
    EyeIcon,
    EyeSlashIcon,
    KeyIcon,
    BeakerIcon,
    GlobeAltIcon,
    ClipboardDocumentIcon,
    ClipboardDocumentCheckIcon,
    LinkIcon,
} from '@heroicons/react/24/outline';
import { Input } from '@/shared/ui';

export interface BoldPlatformCredentials {
    api_key: string;
    secret_key: string;
    test_api_key: string;
    test_secret_key: string;
    link_api_key: string;
    link_secret_key: string;
    test_link_api_key: string;
    test_link_secret_key: string;
}

interface BoldTypeCredentialsFormProps {
    credentials: BoldPlatformCredentials;
    onChange: (credentials: BoldPlatformCredentials) => void;
    isEditing?: boolean;
    webhookUrlProd?: string;
    webhookUrlTest?: string;
}

interface SecretInputProps {
    value: string;
    onChange: (v: string) => void;
    placeholder: string;
    label: string;
    helper: string;
}

function SecretInput({ value, onChange, placeholder, label, helper }: SecretInputProps) {
    const [show, setShow] = useState(false);
    return (
        <div>
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1.5">{label}</label>
            <div className="relative">
                <Input
                    type={show ? 'text' : 'password'}
                    value={value}
                    onChange={(e) => onChange(e.target.value)}
                    placeholder={placeholder}
                    autoComplete="off"
                    className="bg-white dark:bg-gray-800 font-mono text-sm pr-10"
                />
                <button
                    type="button"
                    onClick={() => setShow(!show)}
                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700"
                    tabIndex={-1}
                >
                    {show ? <EyeSlashIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
                </button>
            </div>
            <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">{helper}</p>
        </div>
    );
}

export default function BoldTypeCredentialsForm({
    credentials,
    onChange,
    isEditing = false,
    webhookUrlProd,
    webhookUrlTest,
}: BoldTypeCredentialsFormProps) {
    const placeholderProd = isEditing ? 'Dejar vacío para mantener actual' : 'Pega aquí tu llave de producción';
    const placeholderTest = isEditing ? 'Dejar vacío para mantener actual' : 'Pega aquí tu llave de pruebas';

    const [copiedKey, setCopiedKey] = useState<'prod' | 'test' | null>(null);
    const hasWebhookUrls = Boolean(webhookUrlProd || webhookUrlTest);

    const handleCopy = async (key: 'prod' | 'test', value: string) => {
        await navigator.clipboard.writeText(value);
        setCopiedKey(key);
        setTimeout(() => setCopiedKey(null), 2000);
    };

    return (
        <div className="space-y-4">
            <div className="flex items-center gap-2">
                <KeyIcon className="w-5 h-5 text-blue-600" />
                <h3 className="text-sm font-semibold text-gray-900 dark:text-white">
                    Credenciales de Plataforma Bold
                </h3>
            </div>

            <p className="text-xs text-gray-600 dark:text-gray-300">
                Configura ambos juegos de credenciales (producción y sandbox). Bold usará uno u otro según el flag <code>is_testing</code> de cada negocio. Se almacenan encriptadas (AES-256-GCM).
            </p>

            {hasWebhookUrls ? (
                <div className="p-4 rounded-xl border border-blue-200 bg-blue-50 dark:bg-blue-950/20 dark:border-blue-800">
                    <div className="flex items-center gap-2 mb-2">
                        <LinkIcon className="w-5 h-5 text-blue-700" />
                        <h4 className="font-semibold text-blue-900 dark:text-blue-200">
                            Webhooks de Bold
                        </h4>
                    </div>
                    <div className="text-xs text-gray-600 dark:text-gray-300 mb-3 space-y-1">
                        <p>
                            En el panel de Bold abre <strong>Integraciones → Webhooks → Crear webhook</strong> y crea <strong>dos</strong> webhooks (uno por ambiente). En cada uno marca los 4 eventos: <strong>Venta aprobada</strong>, <strong>Venta rechazada</strong>, <strong>Anulación aprobada</strong>, <strong>Anulación rechazada</strong>.
                        </p>
                        <p>
                            Bold firma cada llamada con HMAC-SHA256 usando la <strong>Secret Key</strong> del ambiente correspondiente; el backend valida la firma antes de procesar el evento.
                        </p>
                    </div>

                    <div className="grid grid-cols-1 lg:grid-cols-2 gap-3">
                        {webhookUrlProd && (
                            <div className="p-3 rounded-lg border border-emerald-200 bg-emerald-50/60 dark:bg-emerald-950/20 dark:border-emerald-800">
                                <div className="flex items-center gap-2 mb-2">
                                    <GlobeAltIcon className="w-4 h-4 text-emerald-700" />
                                    <span className="text-xs font-semibold text-emerald-900 dark:text-emerald-200 uppercase tracking-wide">
                                        Producción
                                    </span>
                                </div>
                                <p className="text-[11px] text-gray-600 dark:text-gray-300 mb-2">
                                    En Bold deja <strong>desmarcado</strong> "¿Este es un webhook de prueba?". Se valida con la <strong>Secret Key</strong> de producción.
                                </p>
                                <div className="flex gap-2">
                                    <div className="flex-1 px-2 py-1.5 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded font-mono text-xs text-gray-700 dark:text-gray-200 select-all break-all">
                                        {webhookUrlProd}
                                    </div>
                                    <button
                                        type="button"
                                        onClick={() => handleCopy('prod', webhookUrlProd)}
                                        className="px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded hover:bg-gray-100 dark:bg-gray-700 transition-colors flex items-center gap-1 text-xs shrink-0"
                                        title="Copiar URL de producción"
                                    >
                                        {copiedKey === 'prod' ? (
                                            <>
                                                <ClipboardDocumentCheckIcon className="w-4 h-4 text-green-600" />
                                                <span className="text-green-600">Copiado</span>
                                            </>
                                        ) : (
                                            <>
                                                <ClipboardDocumentIcon className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                                                <span>Copiar</span>
                                            </>
                                        )}
                                    </button>
                                </div>
                            </div>
                        )}

                        {webhookUrlTest && (
                            <div className="p-3 rounded-lg border border-amber-200 bg-amber-50/60 dark:bg-amber-950/20 dark:border-amber-800">
                                <div className="flex items-center gap-2 mb-2">
                                    <BeakerIcon className="w-4 h-4 text-amber-700" />
                                    <span className="text-xs font-semibold text-amber-900 dark:text-amber-200 uppercase tracking-wide">
                                        Pruebas (Sandbox)
                                    </span>
                                </div>
                                <p className="text-[11px] text-gray-600 dark:text-gray-300 mb-2">
                                    En Bold <strong>marca</strong> "¿Este es un webhook de prueba?". Se valida con la <strong>Test Secret Key</strong>.
                                </p>
                                <div className="flex gap-2">
                                    <div className="flex-1 px-2 py-1.5 bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded font-mono text-xs text-gray-700 dark:text-gray-200 select-all break-all">
                                        {webhookUrlTest}
                                    </div>
                                    <button
                                        type="button"
                                        onClick={() => handleCopy('test', webhookUrlTest)}
                                        className="px-2 py-1.5 border border-gray-300 dark:border-gray-600 rounded hover:bg-gray-100 dark:bg-gray-700 transition-colors flex items-center gap-1 text-xs shrink-0"
                                        title="Copiar URL de pruebas"
                                    >
                                        {copiedKey === 'test' ? (
                                            <>
                                                <ClipboardDocumentCheckIcon className="w-4 h-4 text-green-600" />
                                                <span className="text-green-600">Copiado</span>
                                            </>
                                        ) : (
                                            <>
                                                <ClipboardDocumentIcon className="w-4 h-4 text-gray-500 dark:text-gray-400" />
                                                <span>Copiar</span>
                                            </>
                                        )}
                                    </button>
                                </div>
                            </div>
                        )}
                    </div>
                </div>
            ) : (
                <div className="p-3 rounded-lg border border-yellow-200 bg-yellow-50 dark:bg-yellow-950/20 dark:border-yellow-800 text-xs text-yellow-900 dark:text-yellow-200">
                    Las URLs de webhook no están disponibles porque <code>WEBHOOK_BASE_URL</code> no está configurado en el backend.
                </div>
            )}

            <div className="space-y-2">
                <div className="flex items-center gap-2">
                    <KeyIcon className="w-4 h-4 text-blue-600" />
                    <h4 className="text-sm font-semibold text-gray-900 dark:text-white">
                        Botón de Pago (checkout embebido)
                    </h4>
                </div>
                <p className="text-xs text-gray-600 dark:text-gray-300">
                    Llaves del producto <strong>Botón de Pago</strong>. Se usan para generar la firma de integridad del checkout y validar los webhooks que Bold envía.
                </p>
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                    <div className="space-y-3 p-4 rounded-xl border border-emerald-200 bg-emerald-50 dark:bg-emerald-950/20 dark:border-emerald-800">
                        <div className="flex items-center gap-2 pb-2 border-b border-emerald-200 dark:border-emerald-800">
                            <GlobeAltIcon className="w-5 h-5 text-emerald-700" />
                            <h4 className="font-semibold text-emerald-900 dark:text-emerald-200">Producción</h4>
                        </div>
                        <SecretInput
                            label="Identity Key (api_key)"
                            value={credentials.api_key}
                            onChange={(v) => onChange({ ...credentials, api_key: v })}
                            placeholder={placeholderProd}
                            helper="Llave pública de identidad del Botón de Pago."
                        />
                        <SecretInput
                            label="Secret Key"
                            value={credentials.secret_key}
                            onChange={(v) => onChange({ ...credentials, secret_key: v })}
                            placeholder={placeholderProd}
                            helper="Firma requests y valida webhooks (HMAC-SHA256)."
                        />
                    </div>

                    <div className="space-y-3 p-4 rounded-xl border border-amber-200 bg-amber-50 dark:bg-amber-950/20 dark:border-amber-800">
                        <div className="flex items-center gap-2 pb-2 border-b border-amber-200 dark:border-amber-800">
                            <BeakerIcon className="w-5 h-5 text-amber-700" />
                            <h4 className="font-semibold text-amber-900 dark:text-amber-200">Sandbox (Pruebas)</h4>
                        </div>
                        <SecretInput
                            label="Identity Key (api_key)"
                            value={credentials.test_api_key}
                            onChange={(v) => onChange({ ...credentials, test_api_key: v })}
                            placeholder={placeholderTest}
                            helper="Llave de identidad para sandbox del Botón de Pago."
                        />
                        <SecretInput
                            label="Secret Key"
                            value={credentials.test_secret_key}
                            onChange={(v) => onChange({ ...credentials, test_secret_key: v })}
                            placeholder={placeholderTest}
                            helper="Secret usado para firmar/validar en sandbox."
                        />
                    </div>
                </div>
            </div>

            <div className="space-y-2">
                <div className="flex items-center gap-2">
                    <KeyIcon className="w-4 h-4 text-purple-600" />
                    <h4 className="text-sm font-semibold text-gray-900 dark:text-white">
                        API Payment Links (consultas server-to-server)
                    </h4>
                </div>
                <p className="text-xs text-gray-600 dark:text-gray-300">
                    Llaves del producto <strong>API de Payment Links</strong> (host <code>integrations.api.bold.co</code>). Se usan para consultar el estado de una transacción cuando el webhook tarde o no llegue (polling sync) y para crear payment links en los flujos de orden.
                </p>
                <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
                    <div className="space-y-3 p-4 rounded-xl border border-emerald-200 bg-emerald-50 dark:bg-emerald-950/20 dark:border-emerald-800">
                        <div className="flex items-center gap-2 pb-2 border-b border-emerald-200 dark:border-emerald-800">
                            <GlobeAltIcon className="w-5 h-5 text-emerald-700" />
                            <h4 className="font-semibold text-emerald-900 dark:text-emerald-200">Producción</h4>
                        </div>
                        <SecretInput
                            label="Identity Key API (link_api_key)"
                            value={credentials.link_api_key}
                            onChange={(v) => onChange({ ...credentials, link_api_key: v })}
                            placeholder={placeholderProd}
                            helper="Llave pública del API de Payment Links de Bold."
                        />
                        <SecretInput
                            label="Secret Key API"
                            value={credentials.link_secret_key}
                            onChange={(v) => onChange({ ...credentials, link_secret_key: v })}
                            placeholder={placeholderProd}
                            helper="Secret del API (firma de requests salientes)."
                        />
                    </div>

                    <div className="space-y-3 p-4 rounded-xl border border-amber-200 bg-amber-50 dark:bg-amber-950/20 dark:border-amber-800">
                        <div className="flex items-center gap-2 pb-2 border-b border-amber-200 dark:border-amber-800">
                            <BeakerIcon className="w-5 h-5 text-amber-700" />
                            <h4 className="font-semibold text-amber-900 dark:text-amber-200">Sandbox (Pruebas)</h4>
                        </div>
                        <SecretInput
                            label="Identity Key API (test_link_api_key)"
                            value={credentials.test_link_api_key}
                            onChange={(v) => onChange({ ...credentials, test_link_api_key: v })}
                            placeholder={placeholderTest}
                            helper="Llave de identidad del API de Payment Links en sandbox."
                        />
                        <SecretInput
                            label="Secret Key API"
                            value={credentials.test_link_secret_key}
                            onChange={(v) => onChange({ ...credentials, test_link_secret_key: v })}
                            placeholder={placeholderTest}
                            helper="Secret del API en sandbox."
                        />
                    </div>
                </div>
            </div>
        </div>
    );
}
