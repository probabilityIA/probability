'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select, Modal } from '@/shared/ui';
import { AlegraCredentials } from '../../domain/types';
import { updateIntegrationAction, testConnectionRawAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    CheckBadgeIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    EyeIcon,
    EyeSlashIcon
} from '@heroicons/react/24/outline';

interface AlegraEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: any;
        credentials?: AlegraCredentials;
        business_id?: number | null;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function AlegraEditForm({ integrationId, initialData, onSuccess, onCancel }: AlegraEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [showToken, setShowToken] = useState(false);

    // Business selection for super admins
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId] = useState<number | null>(initialData.business_id || null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData.name,
        email: initialData.credentials?.email || '',
        token: initialData.credentials?.token || '',
        base_url: initialData.credentials?.base_url || '',
    });

    // Check if user is super admin and load businesses
    useEffect(() => {
        const checkUserAndLoadBusinesses = async () => {
            const permissions = TokenStorage.getPermissions();
            const isSuperUser = permissions?.is_super || false;
            setIsSuperAdmin(isSuperUser);

            if (isSuperUser) {
                setLoadingBusinesses(true);
                try {
                    const response = await getBusinessesSimpleAction();
                    if (response.success && response.data) {
                        setBusinesses(response.data);
                    }
                } catch (err) {
                    console.error('Error loading businesses:', err);
                    showToast('Error al cargar la lista de negocios', 'error');
                } finally {
                    setLoadingBusinesses(false);
                }
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    const handleTestConnection = async () => {
        if (!formData.email || !formData.token) {
            showToast('Debes ingresar Email y Token API para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials = {
                email: formData.email,
                token: formData.token,
                base_url: formData.base_url || undefined,
            };

            const result = await testConnectionRawAction('alegra', {}, credentials);

            if (result.success) {
                showToast('Conexion exitosa con Alegra', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con Alegra');
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
                config: {},
            };

            // Only include credentials if they were filled in
            if (formData.email && formData.token) {
                const credentials: AlegraCredentials = {
                    email: formData.email,
                    token: formData.token,
                    base_url: formData.base_url || undefined,
                };
                updateData.credentials = credentials;
            }

            const response = await updateIntegrationAction(integrationId, updateData);

            if (response.success) {
                showToast('Integracion Alegra actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al actualizar integracion de Alegra');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-8" autoComplete="off">
            {/* Header */}
            <div className="border-b border-gray-200 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-purple-50 rounded-lg">
                        <CheckBadgeIcon className="w-6 h-6 text-purple-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900">
                        Editar Alegra Facturacion
                    </h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">
                    Actualiza la configuracion de tu integracion con Alegra.
                </p>
            </div>

            {/* Configuracion General */}
            <div className="bg-gray-50 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-4">
                    <Cog6ToothIcon className="w-5 h-5 text-gray-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Configuracion General
                    </h3>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Nombre de la Integracion <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="Ej: Alegra Facturacion Principal"
                        required
                        className="bg-white"
                    />
                </div>

                {/* Business info - Read only for super admins */}
                {isSuperAdmin && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Negocio
                        </label>
                        {loadingBusinesses ? (
                            <div className="flex items-center gap-2 p-3 bg-gray-50 rounded-lg">
                                <svg className="animate-spin h-5 w-5 text-purple-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
                                <span className="text-sm text-gray-600">Cargando negocios...</span>
                            </div>
                        ) : (
                            <Select
                                value={selectedBusinessId?.toString() || ''}
                                onChange={() => {}}
                                options={[
                                    { value: '', label: '-- Sin negocio asignado --' },
                                    ...businesses.map((business) => ({
                                        value: business.id.toString(),
                                        label: business.name,
                                    })),
                                ]}
                                className="bg-white"
                                disabled={true}
                            />
                        )}
                        <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>El negocio no puede ser modificado despues de la creacion</span>
                        </p>
                    </div>
                )}
            </div>

            {/* Credenciales */}
            <div className="bg-gradient-to-br from-purple-50 to-fuchsia-50 rounded-xl p-6 space-y-4 border border-purple-100">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-purple-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Credenciales de Acceso
                    </h3>
                </div>
                <p className="text-sm text-purple-900 -mt-2 mb-4 flex items-start gap-2">
                    <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                    <span>
                        Aqui se muestran las credenciales actuales. Puedes modificarlas si necesitas actualizarlas.
                        Todos los campos deben estar completos para actualizar las credenciales.
                    </span>
                </p>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Email de la Cuenta
                    </label>
                    <Input
                        type="email"
                        value={formData.email}
                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                        placeholder="email@empresa.com"
                        autoComplete="off"
                        data-1p-ignore
                        className="bg-white text-sm"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Token API
                    </label>
                    <div className="relative">
                        <Input
                            type={showToken ? "text" : "password"}
                            value={formData.token}
                            onChange={(e) => setFormData({ ...formData, token: e.target.value })}
                            placeholder="Token o API Key de Alegra"
                            autoComplete="new-password"
                            data-1p-ignore
                            className="bg-white font-mono text-sm pr-10"
                        />
                        <button
                            type="button"
                            onClick={() => setShowToken(!showToken)}
                            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                            tabIndex={-1}
                        >
                            {showToken ? (
                                <EyeSlashIcon className="w-5 h-5" />
                            ) : (
                                <EyeIcon className="w-5 h-5" />
                            )}
                        </button>
                    </div>
                </div>

                {/* URL de la API */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        URL de la API
                    </label>
                    <Input
                        type="url"
                        value={formData.base_url}
                        onChange={(e) => setFormData({ ...formData, base_url: e.target.value })}
                        placeholder="https://api.alegra.com"
                        autoComplete="off"
                        className="bg-white font-mono text-sm"
                    />
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Dejar vacio para usar la URL de produccion de Alegra.</span>
                    </p>
                </div>

                {/* Test Connection Button */}
                <div className="pt-2">
                    <Button
                        type="button"
                        variant="outline"
                        className="w-full bg-white hover:bg-purple-50 border-purple-200 text-purple-700 font-semibold"
                        onClick={handleTestConnection}
                        disabled={testingConnection || loading || !formData.email || !formData.token}
                    >
                        {testingConnection ? (
                            <>
                                <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-purple-700" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                </svg>
                                Probando...
                            </>
                        ) : (
                            <>
                                <CheckBadgeIcon className="w-4 h-4 mr-2" />
                                Probar Conexion
                            </>
                        )}
                    </Button>
                </div>
            </div>

            {/* Action Buttons */}
            <div className="flex justify-between items-center gap-3 pt-6 border-t border-gray-200">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="min-w-[140px] bg-gray-100 hover:bg-gray-200 text-gray-700 border border-gray-300"
                    >
                        <ArrowLeftIcon className="w-4 h-4 mr-2" />
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading}
                    className="min-w-[200px] bg-purple-600 hover:bg-purple-700 text-white font-semibold"
                >
                    {loading ? (
                        <>
                            <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Actualizando...
                        </>
                    ) : (
                        <>
                            <CheckBadgeIcon className="w-5 h-5 mr-2" />
                            Actualizar Integracion
                        </>
                    )}
                </Button>
            </div>
            {/* Error Modal */}
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
