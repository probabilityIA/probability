'use client';

import { useState } from 'react';
import { EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline';
import { Input, Button, Modal } from '@/shared/ui';
import { testConnectionRawAction } from '@/services/integrations/core/infra/actions';

export interface WhatsAppPlatformCredentials {
    whatsapp_url: string;
    webhook_callback_url: string;
    phone_number_id: string;
    access_token: string;
    verify_token: string;
    test_phone_number: string;
}

interface WhatsAppTypeCredentialsFormProps {
    credentials: WhatsAppPlatformCredentials;
    onChange: (credentials: WhatsAppPlatformCredentials) => void;
    isEditing?: boolean;
}

export default function WhatsAppTypeCredentialsForm({
    credentials,
    onChange,
    isEditing = false,
}: WhatsAppTypeCredentialsFormProps) {
    const [showAccessToken, setShowAccessToken] = useState(false);
    const [showVerifyToken, setShowVerifyToken] = useState(false);
    const [testing, setTesting] = useState(false);
    const [testResult, setTestResult] = useState<{ type: 'success' | 'error'; message: string } | null>(null);

    const handleChange = (field: keyof WhatsAppPlatformCredentials, value: string) => {
        onChange({ ...credentials, [field]: value });
    };

    const handleTestConnection = async () => {
        setTesting(true);
        setTestResult(null);

        if (!credentials.whatsapp_url.trim() || !credentials.phone_number_id.trim() || !credentials.access_token.trim()) {
            setTestResult({ type: 'error', message: 'WhatsApp URL, Phone Number ID y Access Token son requeridos para probar' });
            setTesting(false);
            return;
        }

        if (!credentials.test_phone_number.trim()) {
            setTestResult({ type: 'error', message: 'Ingresa un numero de pruebas para enviar el mensaje' });
            setTesting(false);
            return;
        }

        try {
            const config = {
                whatsapp_url: credentials.whatsapp_url.trim(),
                phone_number_id: credentials.phone_number_id.trim(),
                test_phone_number: credentials.test_phone_number.trim(),
            };
            const creds = {
                access_token: credentials.access_token.trim(),
            };

            const result = await testConnectionRawAction('whatsapp', config, creds);

            if (result.success) {
                setTestResult({ type: 'success', message: 'Mensaje de prueba (hello_world) enviado correctamente' });
            } else {
                setTestResult({ type: 'error', message: result.message || 'Error al enviar mensaje de prueba' });
            }
        } catch (err: any) {
            setTestResult({ type: 'error', message: err.message || 'Error al probar la conexion' });
        } finally {
            setTesting(false);
        }
    };

    return (
        <div className="space-y-4">
            <h3 className="text-sm font-medium text-gray-700">
                Credenciales de WhatsApp (se encriptan)
            </h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 p-4 border border-gray-200 rounded-lg bg-gray-50">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        WhatsApp API URL *
                    </label>
                    <Input
                        type="url"
                        value={credentials.whatsapp_url}
                        onChange={(e) => handleChange('whatsapp_url', e.target.value)}
                        placeholder="https://graph.facebook.com/v22.0"
                        className="font-mono"
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        URL base de la API de WhatsApp Cloud (Meta Graph API)
                    </p>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Webhook Callback URL
                    </label>
                    <Input
                        type="url"
                        value={credentials.webhook_callback_url}
                        onChange={(e) => handleChange('webhook_callback_url', e.target.value)}
                        placeholder="https://api.tudominio.com/api/integrations/whatsapp/webhook"
                        className="font-mono"
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        URL que se configura en Meta &rarr; WhatsApp &rarr; Configuration &rarr; Webhook
                    </p>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Phone Number ID *
                    </label>
                    <Input
                        type="text"
                        name="wa_phone_number_id"
                        autoComplete="new-password"
                        value={credentials.phone_number_id}
                        onChange={(e) => handleChange('phone_number_id', e.target.value)}
                        placeholder="123456789012345"
                        className="font-mono"
                    />
                    <p className="mt-1 text-xs text-gray-500">
                        ID del numero de telefono en Meta Business Manager
                    </p>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Access Token *
                    </label>
                    <div className="relative">
                        <Input
                            type={showAccessToken ? 'text' : 'password'}
                            name="wa_access_token"
                            autoComplete="new-password"
                            value={credentials.access_token}
                            onChange={(e) => handleChange('access_token', e.target.value)}
                            placeholder="EAAxxxxxxxxx..."
                            className="font-mono pr-10"
                        />
                        <button
                            type="button"
                            onClick={() => setShowAccessToken((v) => !v)}
                            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
                        >
                            {showAccessToken ? <EyeSlashIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
                        </button>
                    </div>
                    <p className="mt-1 text-xs text-gray-500">
                        Token de acceso permanente de la app en Meta
                    </p>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Verify Token (Webhook)
                    </label>
                    <div className="relative">
                        <Input
                            type={showVerifyToken ? 'text' : 'password'}
                            name="wa_verify_token"
                            autoComplete="new-password"
                            value={credentials.verify_token}
                            onChange={(e) => handleChange('verify_token', e.target.value)}
                            placeholder="mi_token_secreto"
                            className="font-mono pr-10"
                        />
                        <button
                            type="button"
                            onClick={() => setShowVerifyToken((v) => !v)}
                            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition-colors"
                        >
                            {showVerifyToken ? <EyeSlashIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
                        </button>
                    </div>
                    <p className="mt-1 text-xs text-gray-500">
                        Token de verificacion para el webhook (Meta &rarr; Configuration)
                    </p>
                </div>

                {/* Test connection - full width */}
                <div className="md:col-span-2 border-t border-gray-300 pt-4">
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Numero de pruebas
                    </label>
                    <div className="flex gap-2">
                        <Input
                            type="text"
                            value={credentials.test_phone_number}
                            onChange={(e) => handleChange('test_phone_number', e.target.value)}
                            placeholder="+573001234567"
                            className="font-mono flex-1"
                        />
                        <Button
                            type="button"
                            variant="outline"
                            onClick={handleTestConnection}
                            disabled={testing}
                            loading={testing}
                        >
                            {testing ? 'Enviando...' : 'Probar conexion'}
                        </Button>
                    </div>
                    <p className="mt-1 text-xs text-gray-500">
                        Se enviara un mensaje de prueba (hello_world) a este numero para verificar las credenciales
                    </p>
                </div>

            </div>
            {isEditing && (
                <p className="text-xs text-gray-500">
                    Deja los campos vacios para no modificar los valores actuales.
                </p>
            )}

            {/* Modal de resultado del test */}
            <Modal
                isOpen={testResult !== null}
                onClose={() => setTestResult(null)}
                size="sm"
                showCloseButton={false}
                zIndex={70}
            >
                {testResult && (
                    <div className="flex flex-col items-center text-center py-4">
                        {testResult.type === 'success' ? (
                            <CheckCircleIcon className="w-16 h-16 text-green-500 mb-4" />
                        ) : (
                            <XCircleIcon className="w-16 h-16 text-red-500 mb-4" />
                        )}
                        <h3 className={`text-lg font-semibold mb-2 ${testResult.type === 'success' ? 'text-green-700' : 'text-red-700'}`}>
                            {testResult.type === 'success' ? 'Conexion exitosa' : 'Error de conexion'}
                        </h3>
                        <p className="text-sm text-gray-600 mb-6">
                            {testResult.message}
                        </p>
                        <Button
                            type="button"
                            variant={testResult.type === 'success' ? 'primary' : 'outline'}
                            onClick={() => setTestResult(null)}
                        >
                            Cerrar
                        </Button>
                    </div>
                )}
            </Modal>
        </div>
    );
}
