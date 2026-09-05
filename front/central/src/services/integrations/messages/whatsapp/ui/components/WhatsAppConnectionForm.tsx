'use client';

import { useState } from 'react';
import { Alert, Button, Input, SecretInput } from '@/shared/ui';
import { WhatsAppConnectionValues } from '../../domain/types';

interface WhatsAppConnectionFormProps {
    config: Record<string, any>;
    onSave: (
        config: Record<string, any>,
        credentials: Record<string, any>
    ) => Promise<{ success: boolean; message?: string }>;
    onSaved?: () => void;
}

function readConfig(config: Record<string, any>): WhatsAppConnectionValues {
    return {
        waba_id: config?.waba_id ? String(config.waba_id) : '',
        phone_number_id: config?.phone_number_id ? String(config.phone_number_id) : '',
        access_token: '',
        use_platform_token: config?.use_platform_token !== false,
    };
}

export default function WhatsAppConnectionForm({ config, onSave, onSaved }: WhatsAppConnectionFormProps) {
    const [values, setValues] = useState<WhatsAppConnectionValues>(() => readConfig(config));
    const [saving, setSaving] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    const hasStoredToken = Boolean(config?.phone_number_id) && config?.use_platform_token === false;

    const handleChange = (field: keyof WhatsAppConnectionValues, value: string | boolean) => {
        setValues((prev) => ({ ...prev, [field]: value }));
    };

    const handleSubmit = async () => {
        setMessage(null);

        if (!values.use_platform_token) {
            if (!values.phone_number_id.trim() || !values.waba_id.trim()) {
                setMessage({
                    type: 'error',
                    text: 'Para usar tu propio número necesitas el WABA ID y el Phone Number ID',
                });
                return;
            }
            if (!values.access_token.trim() && !hasStoredToken) {
                setMessage({ type: 'error', text: 'Ingresa el token de acceso de tu cuenta de WhatsApp' });
                return;
            }
        }

        setSaving(true);
        try {
            const nextConfig: Record<string, any> = {
                ...config,
                use_platform_token: values.use_platform_token,
            };

            if (values.use_platform_token) {
                delete nextConfig.waba_id;
                delete nextConfig.phone_number_id;
            } else {
                nextConfig.waba_id = values.waba_id.trim();
                nextConfig.phone_number_id = values.phone_number_id.trim();
            }

            const credentials: Record<string, any> = {};
            if (!values.use_platform_token && values.access_token.trim()) {
                credentials.access_token = values.access_token.trim();
            }

            const result = await onSave(nextConfig, credentials);
            if (result.success) {
                setMessage({ type: 'success', text: 'Conexión de WhatsApp guardada' });
                setValues((prev) => ({ ...prev, access_token: '' }));
                onSaved?.();
            } else {
                setMessage({ type: 'error', text: result.message || 'Error al guardar la conexión' });
            }
        } catch (err: any) {
            setMessage({ type: 'error', text: err?.message || 'Error al guardar la conexión' });
        } finally {
            setSaving(false);
        }
    };

    return (
        <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-4">
            <div>
                <h4 className="text-sm font-medium text-gray-700 dark:text-gray-200">Número de WhatsApp</h4>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    Puedes enviar desde el número de Probability o conectar el tuyo. Para usar el tuyo,
                    comparte tu cuenta de WhatsApp con Probability desde tu Business Manager de Meta y pega
                    aquí los identificadores.
                </p>
            </div>

            <div className="flex items-center gap-3">
                <button
                    type="button"
                    role="switch"
                    aria-checked={!values.use_platform_token}
                    onClick={() => handleChange('use_platform_token', !values.use_platform_token)}
                    className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none ${
                        !values.use_platform_token ? 'bg-green-500' : 'bg-gray-300 dark:bg-gray-600'
                    }`}
                >
                    <span
                        className={`inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform ${
                            !values.use_platform_token ? 'translate-x-6' : 'translate-x-1'
                        }`}
                    />
                </button>
                <div>
                    <p className="text-sm font-medium text-gray-700 dark:text-gray-200">
                        Usar mi propio número
                    </p>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                        Apagado: los mensajes salen del número de Probability, como hasta ahora.
                    </p>
                </div>
            </div>

            {!values.use_platform_token && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            WABA ID *
                        </label>
                        <Input
                            type="text"
                            value={values.waba_id}
                            onChange={(e) => handleChange('waba_id', e.target.value)}
                            placeholder="123456789012345"
                            className="font-mono"
                        />
                        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                            ID de tu cuenta de WhatsApp Business en Meta
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Phone Number ID *
                        </label>
                        <Input
                            type="text"
                            value={values.phone_number_id}
                            onChange={(e) => handleChange('phone_number_id', e.target.value)}
                            placeholder="123456789012345"
                            className="font-mono"
                        />
                        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                            ID del número desde el que se envían los mensajes
                        </p>
                    </div>

                    <div className="md:col-span-2">
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-1">
                            Token de acceso {hasStoredToken ? '' : '*'}
                        </label>
                        <SecretInput
                            value={values.access_token}
                            onChange={(e) => handleChange('access_token', e.target.value)}
                            placeholder={hasStoredToken ? 'Guárdalo vacío para no cambiarlo' : 'EAAxxxxxxxxx...'}
                            className="font-mono"
                        />
                        <p className="mt-1 text-xs text-gray-500 dark:text-gray-400">
                            Token permanente con permisos de mensajería sobre tu cuenta de WhatsApp
                        </p>
                    </div>
                </div>
            )}

            <Button type="button" variant="primary" onClick={handleSubmit} disabled={saving} loading={saving}>
                Guardar conexión
            </Button>

            {message && (
                <Alert type={message.type} onClose={() => setMessage(null)}>
                    {message.text}
                </Alert>
            )}
        </div>
    );
}
