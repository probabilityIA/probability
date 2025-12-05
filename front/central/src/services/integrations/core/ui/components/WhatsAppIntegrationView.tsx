'use client';

import { useState } from 'react';
import { Button, Alert } from '@/shared/ui';

interface WhatsAppConfig {
    phone_number_id?: string;
    business_account_id?: string;
    webhook_url?: string;
    [key: string]: any;
}

interface WhatsAppCredentials {
    access_token?: string;
    [key: string]: any;
}

interface WhatsAppIntegrationViewProps {
    integration: {
        id: number;
        name: string;
        code: string;
        config?: WhatsAppConfig;
        credentials?: WhatsAppCredentials;
        is_active: boolean;
        created_at: string;
        updated_at: string;
    };
    onTestConnection?: () => Promise<boolean>;
    onRefresh?: () => void;
}

export default function WhatsAppIntegrationView({
    integration,
    onTestConnection,
    onRefresh
}: WhatsAppIntegrationViewProps) {
    const [testing, setTesting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [testSuccess, setTestSuccess] = useState(false);

    const handleTestConnection = async () => {
        setTesting(true);
        setError(null);
        setTestSuccess(false);

        try {
            if (onTestConnection) {
                const success = await onTestConnection();
                if (success) {
                    setTestSuccess(true);
                } else {
                    setError('No se pudo conectar con WhatsApp');
                }
            }
        } catch (err: any) {
            console.error('Test connection error:', err);
            setError(err.message || 'Error al probar la conexión');
        } finally {
            setTesting(false);
        }
    };

    const maskSensitiveData = (value: string | undefined) => {
        if (!value) return 'No configurado';
        if (value.length <= 8) return '••••••••';
        return value.substring(0, 4) + '••••••••' + value.substring(value.length - 4);
    };

    return (
        <div className="space-y-4">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {testSuccess && (
                <Alert type="success" onClose={() => setTestSuccess(false)}>
                    ✓ Conexión exitosa con WhatsApp
                </Alert>
            )}

            <Alert type="info">
                <div className="flex items-start">
                    <svg className="w-5 h-5 mr-2 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                        <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd" />
                    </svg>
                    <div>
                        <p className="font-medium">Integración Interna de WhatsApp</p>
                        <p className="text-sm mt-1">Esta es la integración global de WhatsApp para toda la plataforma. Solo los administradores pueden modificar su configuración.</p>
                    </div>
                </div>
            </Alert>

            {/* Basic Info */}
            <div className="bg-gray-50 rounded-lg p-4 space-y-3">
                <h3 className="text-sm font-semibold text-gray-700">Información Básica</h3>

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Nombre</label>
                        <p className="text-sm text-gray-900">{integration.name}</p>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Código</label>
                        <p className="text-sm text-gray-900 font-mono">{integration.code}</p>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Estado</label>
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${integration.is_active
                                ? 'bg-green-100 text-green-800'
                                : 'bg-red-100 text-red-800'
                            }`}>
                            {integration.is_active ? '✓ Activa' : '✗ Inactiva'}
                        </span>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Última actualización</label>
                        <p className="text-sm text-gray-900">{new Date(integration.updated_at).toLocaleDateString()}</p>
                    </div>
                </div>
            </div>

            {/* Configuration */}
            <div className="bg-blue-50 rounded-lg p-4 space-y-3">
                <h3 className="text-sm font-semibold text-gray-700">Configuración</h3>

                <div className="space-y-2">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Phone Number ID</label>
                        <p className="text-sm text-gray-900 font-mono">{integration.config?.phone_number_id || 'No configurado'}</p>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Business Account ID</label>
                        <p className="text-sm text-gray-900 font-mono">{integration.config?.business_account_id || 'No configurado'}</p>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 mb-1">Webhook URL</label>
                        <p className="text-sm text-gray-900 break-all">{integration.config?.webhook_url || 'No configurado'}</p>
                    </div>
                </div>
            </div>

            {/* Credentials (Masked) */}
            <div className="bg-yellow-50 rounded-lg p-4 space-y-3">
                <div className="flex items-center justify-between">
                    <h3 className="text-sm font-semibold text-gray-700">Credenciales</h3>
                    <svg className="w-4 h-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                    </svg>
                </div>

                <div>
                    <label className="block text-xs font-medium text-gray-500 mb-1">Access Token</label>
                    <p className="text-sm text-gray-900 font-mono">{maskSensitiveData(integration.credentials?.access_token)}</p>
                    <p className="text-xs text-gray-500 mt-1">Por seguridad, el token está oculto</p>
                </div>
            </div>

            {/* Actions */}
            <div className="flex justify-between items-center pt-4 border-t">
                <Button
                    type="button"
                    onClick={onRefresh}
                    variant="outline"
                >
                    🔄 Actualizar
                </Button>

                <Button
                    type="button"
                    onClick={handleTestConnection}
                    disabled={testing}
                    loading={testing}
                    variant="primary"
                >
                    {testing ? 'Probando...' : '🔌 Probar Conexión'}
                </Button>
            </div>
        </div>
    );
}
