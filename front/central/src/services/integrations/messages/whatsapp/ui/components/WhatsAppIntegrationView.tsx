'use client';

import { useState } from 'react';
import { Alert, Button } from '@/shared/ui';
import { ChatBubbleLeftRightIcon, CheckCircleIcon, PhoneIcon } from '@heroicons/react/24/outline';

interface WhatsAppIntegrationViewProps {
    integration: {
        id: number;
        name: string;
        code: string;
        config?: Record<string, any>;
        credentials?: Record<string, any>;
        is_active: boolean;
        created_at: string;
        updated_at: string;
    };
    imageUrl?: string;
    onToggleActive?: (id: number, currentlyActive: boolean) => Promise<boolean>;
    onUpdateConfig?: (id: number, config: Record<string, any>) => Promise<{ success: boolean; message?: string }>;
    onTestConnection?: (id: number) => Promise<{ success: boolean; message?: string }>;
    onRefresh?: () => void;
}

export default function WhatsAppIntegrationView({
    integration,
    imageUrl,
    onToggleActive,
    onUpdateConfig,
    onTestConnection,
    onRefresh,
}: WhatsAppIntegrationViewProps) {
    const [isActive, setIsActive] = useState(integration.is_active);
    const [toggling, setToggling] = useState(false);
    const [testPhone, setTestPhone] = useState(integration.config?.test_phone_number || '');
    const [saving, setSaving] = useState(false);
    const [testing, setTesting] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error'; text: string } | null>(null);

    const savedPhone = integration.config?.test_phone_number || '';
    const hasUnsavedChanges = testPhone !== savedPhone;

    const handleToggle = async () => {
        if (!onToggleActive) return;
        setToggling(true);
        try {
            const success = await onToggleActive(integration.id, isActive);
            if (success) {
                setIsActive((prev) => !prev);
            }
        } finally {
            setToggling(false);
        }
    };

    const handleSavePhone = async () => {
        if (!onUpdateConfig || !testPhone.trim()) return;
        setSaving(true);
        setMessage(null);
        try {
            const updatedConfig = { ...integration.config, test_phone_number: testPhone.trim() };
            const result = await onUpdateConfig(integration.id, updatedConfig);
            if (result.success) {
                setMessage({ type: 'success', text: 'Numero de prueba guardado' });
                onRefresh?.();
            } else {
                setMessage({ type: 'error', text: result.message || 'Error al guardar' });
            }
        } catch (err: any) {
            setMessage({ type: 'error', text: err.message || 'Error al guardar' });
        } finally {
            setSaving(false);
        }
    };

    const handleTestConnection = async () => {
        if (!onTestConnection) return;

        // Si hay cambios sin guardar, guardar primero
        if (hasUnsavedChanges && onUpdateConfig && testPhone.trim()) {
            setSaving(true);
            try {
                const updatedConfig = { ...integration.config, test_phone_number: testPhone.trim() };
                const saveResult = await onUpdateConfig(integration.id, updatedConfig);
                if (!saveResult.success) {
                    setMessage({ type: 'error', text: saveResult.message || 'Error al guardar antes de probar' });
                    setSaving(false);
                    return;
                }
            } catch (err: any) {
                setMessage({ type: 'error', text: err.message || 'Error al guardar' });
                setSaving(false);
                return;
            }
            setSaving(false);
        }

        setTesting(true);
        setMessage(null);
        try {
            const result = await onTestConnection(integration.id);
            if (result.success) {
                setMessage({ type: 'success', text: 'Mensaje de prueba enviado correctamente' });
            } else {
                setMessage({ type: 'error', text: result.message || 'Error en la prueba de conexion' });
            }
        } catch (err: any) {
            setMessage({ type: 'error', text: err.message || 'Error en la prueba' });
        } finally {
            setTesting(false);
        }
    };

    return (
        <div className="space-y-6 max-w-2xl mx-auto py-4">
            {/* Header con icono */}
            <div className="flex flex-col items-center text-center">
                <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4 overflow-hidden">
                    {imageUrl ? (
                        <img
                            src={imageUrl}
                            alt={integration.name}
                            className="w-12 h-12 object-contain"
                            onError={(e) => {
                                const target = e.target as HTMLImageElement;
                                target.style.display = 'none';
                            }}
                        />
                    ) : (
                        <ChatBubbleLeftRightIcon className="w-8 h-8 text-green-600" />
                    )}
                </div>
                <h3 className="text-lg font-semibold text-gray-900">{integration.name}</h3>
                <p className="text-sm text-gray-500 font-mono">{integration.code}</p>
            </div>

            {/* Estado actual - clickeable para activar/desactivar */}
            <div className="flex items-center justify-center">
                {onToggleActive ? (
                    <button
                        type="button"
                        onClick={handleToggle}
                        disabled={toggling}
                        className={`inline-flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium transition-colors cursor-pointer ${
                            isActive
                                ? 'bg-green-100 text-green-800 hover:bg-red-100 hover:text-red-700'
                                : 'bg-red-100 text-red-700 hover:bg-green-100 hover:text-green-800'
                        } ${toggling ? 'opacity-50 cursor-wait' : ''}`}
                        title={isActive ? 'Clic para desactivar' : 'Clic para activar'}
                    >
                        {toggling ? (
                            <svg className="w-5 h-5 animate-spin" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" />
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" />
                            </svg>
                        ) : isActive ? (
                            <CheckCircleIcon className="w-5 h-5" />
                        ) : (
                            <span className="w-2.5 h-2.5 bg-red-400 rounded-full" />
                        )}
                        {isActive ? 'Activa' : 'Inactiva'}
                    </button>
                ) : (
                    <span className={`inline-flex items-center gap-2 px-4 py-2 rounded-full text-sm font-medium ${
                        isActive
                            ? 'bg-green-100 text-green-800'
                            : 'bg-red-100 text-red-700'
                    }`}>
                        {isActive ? (
                            <>
                                <CheckCircleIcon className="w-5 h-5" />
                                Activa
                            </>
                        ) : (
                            <>
                                <span className="w-2.5 h-2.5 bg-red-400 rounded-full" />
                                Inactiva
                            </>
                        )}
                    </span>
                )}
            </div>

            {/* Fechas */}
            <div className="grid grid-cols-2 gap-4 text-center text-sm">
                <div>
                    <p className="text-gray-500">Creada</p>
                    <p className="text-gray-900 font-medium">{new Date(integration.created_at).toLocaleDateString()}</p>
                </div>
                <div>
                    <p className="text-gray-500">Actualizada</p>
                    <p className="text-gray-900 font-medium">{new Date(integration.updated_at).toLocaleDateString()}</p>
                </div>
            </div>

            {/* Numero de telefono de prueba */}
            {isActive && (onUpdateConfig || onTestConnection) && (
                <div className="border border-gray-200 rounded-lg p-4 space-y-3">
                    <div className="flex items-center gap-2 text-sm font-medium text-gray-700">
                        <PhoneIcon className="w-4 h-4" />
                        Numero de telefono para pruebas
                    </div>
                    <p className="text-xs text-gray-500">
                        Ingresa un numero con codigo de pais (ej: 573001234567) para enviar un mensaje de prueba y verificar que la integracion funciona.
                    </p>
                    <div className="flex gap-2">
                        <input
                            type="text"
                            value={testPhone}
                            onChange={(e) => setTestPhone(e.target.value)}
                            placeholder="573001234567"
                            className="flex-1 px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-green-500 focus:border-green-500"
                        />
                        {onUpdateConfig && (
                            <Button
                                type="button"
                                variant="outline"
                                onClick={handleSavePhone}
                                disabled={saving || !testPhone.trim() || !hasUnsavedChanges}
                                loading={saving}
                                size="sm"
                            >
                                Guardar
                            </Button>
                        )}
                    </div>
                    {onTestConnection && (
                        <Button
                            type="button"
                            variant="primary"
                            onClick={handleTestConnection}
                            disabled={testing || !testPhone.trim()}
                            loading={testing}
                            className="w-full"
                        >
                            Enviar mensaje de prueba
                        </Button>
                    )}
                </div>
            )}

            {/* Mensajes de feedback */}
            {message && (
                <Alert type={message.type} onClose={() => setMessage(null)}>
                    {message.text}
                </Alert>
            )}

        </div>
    );
}
