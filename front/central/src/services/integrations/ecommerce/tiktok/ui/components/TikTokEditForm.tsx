'use client';

import { useState, FormEvent } from 'react';
import { Button, Input, Alert, Modal, SecretInput } from '@/shared/ui';
import { TikTokCredentials } from '../../domain/types';
import { updateIntegrationAction, testConnectionRawAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import {
    KeyIcon,
    Cog6ToothIcon,
    CheckBadgeIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    MusicalNoteIcon,
} from '@heroicons/react/24/outline';

interface TikTokEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: any;
        credentials?: TikTokCredentials;
        business_id?: number | null;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function TikTokEditForm({ integrationId, initialData, onSuccess, onCancel }: TikTokEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);

    const [formData, setFormData] = useState({
        name: initialData.name,
        shop_name: initialData.config?.shop_name || '',
        client_key: initialData.credentials?.client_key || '',
        client_secret: initialData.credentials?.client_secret || '',
    });

    const handleTestConnection = async () => {
        if (!formData.client_key || !formData.client_secret) {
            showToast('Debes ingresar Client Key y Client Secret para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials = {
                client_key: formData.client_key,
                client_secret: formData.client_secret,
            };

            const config = formData.shop_name ? { shop_name: formData.shop_name } : {};

            const result = await testConnectionRawAction('tiktok', config, credentials);

            if (result.success) {
                showToast('Conexion exitosa con TikTok', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexión');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con TikTok');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const updateData: any = {
                name: formData.name,
                config: { shop_name: formData.shop_name },
            };

            if (formData.client_key && formData.client_secret) {
                const credentials: TikTokCredentials = {
                    client_key: formData.client_key,
                    client_secret: formData.client_secret,
                };
                updateData.credentials = credentials;
            }

            const response = await updateIntegrationAction(integrationId, updateData);

            if (response.success) {
                showToast('Integracion TikTok actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integración');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al actualizar la integración de TikTok');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-8" autoComplete="off">
            <div className="border-b border-gray-200 dark:border-gray-700 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-gray-100 dark:bg-gray-700 rounded-lg">
                        <MusicalNoteIcon className="w-6 h-6 text-gray-900 dark:text-white" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900 dark:text-white">
                        Editar TikTok
                    </h2>
                </div>
                <p className="text-sm text-gray-600 dark:text-gray-300 ml-14">
                    Actualiza la configuracion y credenciales de tu integracion con TikTok.
                </p>
            </div>

            <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-4">
                    <Cog6ToothIcon className="w-5 h-5 text-gray-700 dark:text-gray-200" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                        Configuracion General
                    </h3>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Nombre de la Integracion <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="Ej: TikTok Shop Principal"
                        required
                        className="bg-white dark:bg-gray-800"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Nombre de la tienda (TikTok Shop) <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.shop_name}
                        onChange={(e) => setFormData({ ...formData, shop_name: e.target.value })}
                        placeholder="mitienda"
                        required
                        className="bg-white dark:bg-gray-800"
                    />
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Nombre exacto de tu TikTok Shop tal como aparece en la plataforma</span>
                    </p>
                </div>
            </div>

            <div className="bg-gradient-to-br from-cyan-50 to-rose-50 dark:from-gray-800 dark:to-gray-800 rounded-xl p-6 space-y-4 border border-cyan-100 dark:border-gray-700">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-cyan-700 dark:text-cyan-400" />
                    <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                        Credenciales de Acceso
                    </h3>
                </div>
                <p className="text-sm text-cyan-900 dark:text-cyan-200 -mt-2 mb-4 flex items-start gap-2">
                    <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                    <span>Deja los campos vacíos para conservar las credenciales actuales</span>
                </p>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Client Key
                    </label>
                    <Input
                        type="text"
                        value={formData.client_key}
                        onChange={(e) => setFormData({ ...formData, client_key: e.target.value })}
                        placeholder="awxxxxxxxxxxxxxxxx"
                        autoComplete="off"
                        data-1p-ignore
                        className="bg-white dark:bg-gray-800 font-mono text-sm"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2">
                        Client Secret
                    </label>
                    <SecretInput
                        value={formData.client_secret}
                        onChange={(e) => setFormData({ ...formData, client_secret: e.target.value })}
                        placeholder="Client Secret de TikTok"
                        className="bg-white dark:bg-gray-800 font-mono text-sm"
                    />
                </div>

                <div className="pt-2">
                    <Button
                        type="button"
                        variant="outline"
                        className="w-full bg-white dark:bg-gray-800 hover:bg-cyan-50 border-cyan-200 text-cyan-700 font-semibold"
                        onClick={handleTestConnection}
                        disabled={testingConnection || loading || !formData.client_key || !formData.client_secret}
                    >
                        {testingConnection ? 'Probando...' : (
                            <>
                                <CheckBadgeIcon className="w-4 h-4 mr-2" />
                                Probar Conexion
                            </>
                        )}
                    </Button>
                </div>
            </div>

            <div className="flex justify-between items-center gap-3 pt-6 border-t border-gray-200 dark:border-gray-700">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="min-w-[140px] bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600"
                    >
                        <ArrowLeftIcon className="w-4 h-4 mr-2" />
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading}
                    className="min-w-[200px] bg-gray-900 hover:bg-black text-white font-semibold"
                >
                    {loading ? 'Guardando...' : (
                        <>
                            <CheckBadgeIcon className="w-5 h-5 mr-2" />
                            Guardar Cambios
                        </>
                    )}
                </Button>
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
        </form>
    );
}
