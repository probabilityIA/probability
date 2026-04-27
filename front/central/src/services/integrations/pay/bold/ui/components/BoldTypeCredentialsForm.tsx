'use client';

import { useState } from 'react';
import { EyeIcon, EyeSlashIcon, KeyIcon, BeakerIcon, GlobeAltIcon } from '@heroicons/react/24/outline';
import { Input } from '@/shared/ui';

export interface BoldPlatformCredentials {
    api_key: string;
    secret_key: string;
    test_api_key: string;
    test_secret_key: string;
}

interface BoldTypeCredentialsFormProps {
    credentials: BoldPlatformCredentials;
    onChange: (credentials: BoldPlatformCredentials) => void;
    isEditing?: boolean;
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
}: BoldTypeCredentialsFormProps) {
    const placeholderProd = isEditing ? 'Dejar vacío para mantener actual' : 'Pega aquí tu llave de producción';
    const placeholderTest = isEditing ? 'Dejar vacío para mantener actual' : 'Pega aquí tu llave de pruebas';

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
                        helper="Llave pública de identidad de Bold (panel Comercios → Integraciones)."
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
                        helper="Llave de identidad para el ambiente de pruebas / mock interno."
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
    );
}
