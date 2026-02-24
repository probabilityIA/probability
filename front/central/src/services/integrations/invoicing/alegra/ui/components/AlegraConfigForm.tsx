'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select, Modal } from '@/shared/ui';
import { AlegraCredentials } from '../../domain/types';
import { createIntegrationAction, testConnectionRawAction } from '@/services/integrations/core/infra/actions';
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

interface AlegraConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function AlegraConfigForm({ onSuccess, onCancel }: AlegraConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [showToken, setShowToken] = useState(false);

    // Business selection for super admins
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: '',
        email: '',
        token: '',
        base_url: '',
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
            } else {
                if (permissions?.business_id) {
                    setSelectedBusinessId(permissions.business_id);
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
            if (isSuperAdmin && !selectedBusinessId) {
                setErrorModal('Debes seleccionar un negocio antes de crear la integracion.');
                setLoading(false);
                return;
            }

            const credentials: AlegraCredentials = {
                email: formData.email,
                token: formData.token,
                base_url: formData.base_url || undefined,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `alegra_${Date.now()}`,
                integration_type_id: 9,
                category: 'invoicing',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: {} as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
            });

            if (response.success) {
                showToast('Integracion Alegra creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al crear la integracion de Alegra');
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
                        Alegra Facturacion Electronica
                    </h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">
                    Conecta tu cuenta de Alegra para gestionar facturacion electronica automaticamente desde Probability.
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
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Nombre descriptivo para identificar esta integracion en el sistema</span>
                    </p>
                </div>

                {/* Business Selector - Only for Super Admins */}
                {isSuperAdmin && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Negocio <span className="text-red-500">*</span>
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
                                onChange={(e) => setSelectedBusinessId(Number(e.target.value))}
                                options={[
                                    { value: '', label: '-- Selecciona un negocio --' },
                                    ...businesses.map((business) => ({
                                        value: business.id.toString(),
                                        label: business.name,
                                    })),
                                ]}
                                required
                                className="bg-white"
                            />
                        )}
                        <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>Selecciona el negocio al que pertenecera esta integracion</span>
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
                        Obten tus credenciales desde el panel de Alegra en
                        <strong> Configuracion &rarr; Integraciones &rarr; API</strong>
                    </span>
                </p>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Email de la Cuenta <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="email"
                        value={formData.email}
                        onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                        placeholder="email@empresa.com"
                        required
                        autoComplete="off"
                        data-1p-ignore
                        className="bg-white text-sm"
                    />
                    <p className="text-xs text-gray-500 mt-1.5">
                        Email con el que accedes a Alegra
                    </p>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Token API <span className="text-red-500">*</span>
                    </label>
                    <div className="relative">
                        <Input
                            type={showToken ? "text" : "password"}
                            value={formData.token}
                            onChange={(e) => setFormData({ ...formData, token: e.target.value })}
                            placeholder="Token o API Key de Alegra"
                            required
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

                <div className="bg-purple-100 border border-purple-200 rounded-lg p-4 mt-4">
                    <h4 className="text-sm font-semibold text-purple-900 mb-3 flex items-center gap-2">
                        <InformationCircleIcon className="w-4 h-4" />
                        Como obtener tus credenciales
                    </h4>
                    <ol className="text-xs text-purple-800 space-y-2 list-decimal list-inside ml-1">
                        <li>Ingresa a tu cuenta en <strong>app.alegra.com</strong></li>
                        <li>Ve a <strong>Configuracion &rarr; Integraciones &rarr; API</strong></li>
                        <li>Copia tu <strong>Email</strong> y <strong>Token API</strong></li>
                    </ol>
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
                            Conectando...
                        </>
                    ) : (
                        <>
                            <CheckBadgeIcon className="w-5 h-5 mr-2" />
                            Crear Integracion
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
