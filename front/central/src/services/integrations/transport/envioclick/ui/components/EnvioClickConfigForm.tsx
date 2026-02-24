'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select } from '@/shared/ui';
import { EnvioClickConfig, EnvioClickCredentials } from '../../domain/types';
import { createIntegrationAction, testConnectionRawAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    TruckIcon,
    CheckBadgeIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    EyeIcon,
    EyeSlashIcon,
    BeakerIcon,
} from '@heroicons/react/24/outline';

interface EnvioClickConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    integrationTypeBaseURL?: string;      // URL de producción global (read-only)
    integrationTypeBaseURLTest?: string;  // URL de pruebas por defecto del tipo
}

export function EnvioClickConfigForm({ onSuccess, onCancel, integrationTypeBaseURL, integrationTypeBaseURLTest }: EnvioClickConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showApiKey, setShowApiKey] = useState(false);
    const [usePlatformToken, setUsePlatformToken] = useState(false);
    const [isTesting, setIsTesting] = useState(false);

    // Business selection for super admins
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: '',
        api_key: '',
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
        if (!usePlatformToken && !formData.api_key) {
            showToast('Debes ingresar la API Key para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);
        setError(null);

        try {
            const config: EnvioClickConfig = {
                use_platform_token: usePlatformToken,
            };

            const credentials: EnvioClickCredentials = usePlatformToken
                ? {}
                : { api_key: formData.api_key };

            const result = await testConnectionRawAction('envioclick', config, credentials);

            if (result.success) {
                showToast('Conexion exitosa con EnvioClick', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setError(err.message || 'Error al probar conexion');
            showToast('Error al conectar con EnvioClick: ' + err.message, 'error');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            if (isSuperAdmin && !selectedBusinessId) {
                showToast('Debes seleccionar un negocio', 'warning');
                setLoading(false);
                return;
            }

            const config: EnvioClickConfig = {
                use_platform_token: usePlatformToken,
            };

            const credentials: EnvioClickCredentials = usePlatformToken
                ? {}
                : { api_key: formData.api_key };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `envioclick_${Date.now()}`,
                integration_type_id: 12,
                category: 'shipping',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
                is_testing: isTesting,
            });

            if (response.success) {
                showToast('Integracion EnvioClick creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setError(err.message || 'Error al crear integracion');
            showToast('Error al crear integracion EnvioClick', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-8" autoComplete="off">
            {/* Header */}
            <div className="border-b border-gray-200 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-green-50 rounded-lg">
                        <TruckIcon className="w-6 h-6 text-green-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900">
                        EnvioClick - Logistica
                    </h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">
                    Conecta tu cuenta de EnvioClick para gestionar envios y cotizaciones de paqueteria desde Probability.
                </p>
            </div>

            {error && (
                <Alert type="error">
                    {error}
                </Alert>
            )}

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
                        placeholder="Ej: EnvioClick Principal"
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
                                <svg className="animate-spin h-5 w-5 text-green-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
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

            {/* Modo de Token */}
            <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl p-6 space-y-4 border border-green-100">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-green-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Credenciales de API
                    </h3>
                </div>

                {/* Toggle: usar token de la plataforma */}
                <div className="flex items-center justify-between p-3 bg-white rounded-lg border border-green-200">
                    <div className="flex-1">
                        <p className="text-sm font-medium text-gray-800">Usar token de la plataforma</p>
                        <p className="text-xs text-gray-500 mt-0.5">
                            Activa esto si no tienes una cuenta propia de EnvioClick. Se usará la cuenta compartida de la plataforma.
                        </p>
                    </div>
                    <button
                        type="button"
                        onClick={() => setUsePlatformToken(!usePlatformToken)}
                        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none ml-4 flex-shrink-0 ${usePlatformToken ? 'bg-green-600' : 'bg-gray-200'}`}
                    >
                        <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${usePlatformToken ? 'translate-x-6' : 'translate-x-1'}`} />
                    </button>
                </div>

                {/* API Key — solo si no usa token de plataforma */}
                {!usePlatformToken && (
                    <div className="space-y-4">
                        <p className="text-sm text-green-900 flex items-start gap-2">
                            <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                            <span>
                                Obtiene tu API Key desde el panel de EnvioClick en la seccion de
                                <strong> Integraciones / API</strong>
                            </span>
                        </p>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-2">
                                API Key <span className="text-red-500">*</span>
                            </label>
                            <div className="relative">
                                <Input
                                    type={showApiKey ? "text" : "password"}
                                    value={formData.api_key}
                                    onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                                    placeholder="Ingresa tu API Key de EnvioClick"
                                    required={!usePlatformToken}
                                    autoComplete="off"
                                    data-1p-ignore
                                    className="bg-white font-mono text-sm pr-10"
                                />
                                <button
                                    type="button"
                                    onClick={() => setShowApiKey(!showApiKey)}
                                    className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                                    tabIndex={-1}
                                >
                                    {showApiKey ? (
                                        <EyeSlashIcon className="w-5 h-5" />
                                    ) : (
                                        <EyeIcon className="w-5 h-5" />
                                    )}
                                </button>
                            </div>
                        </div>
                    </div>
                )}

                {/* Test Connection Button — solo cuando usa token propio */}
                {!usePlatformToken && (
                    <div className="pt-2">
                        <Button
                            type="button"
                            variant="outline"
                            className="w-full bg-white hover:bg-green-50 border-green-200 text-green-700 font-semibold"
                            onClick={handleTestConnection}
                            disabled={testingConnection || loading || !formData.api_key}
                        >
                            {testingConnection ? (
                                <>
                                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-green-700" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
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
                )}
            </div>

            {/* Modo de Pruebas */}
            <div className="bg-orange-50 rounded-xl p-6 space-y-4 border border-orange-200">
                <div className="flex items-center gap-2 mb-2">
                    <BeakerIcon className="w-5 h-5 text-orange-600" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Modo de Pruebas
                    </h3>
                </div>
                <div className="flex items-center justify-between p-3 bg-white rounded-lg border border-orange-200">
                    <div className="flex-1">
                        <p className="text-sm font-medium text-gray-800">Activar modo testing</p>
                        <p className="text-xs text-gray-500 mt-0.5">
                            Los envíos generados quedarán marcados como TEST y usarán la URL de pruebas.
                        </p>
                    </div>
                    <button
                        type="button"
                        onClick={() => setIsTesting(!isTesting)}
                        className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none ml-4 flex-shrink-0 ${isTesting ? 'bg-orange-500' : 'bg-gray-200'}`}
                    >
                        <span className={`inline-block h-4 w-4 transform rounded-full bg-white transition-transform ${isTesting ? 'translate-x-6' : 'translate-x-1'}`} />
                    </button>
                </div>
                {isTesting && (
                    <Alert type="warning">
                        Modo de pruebas activado. Los envíos generados con esta integración quedarán marcados como <strong>TEST</strong> en el sistema.
                        {integrationTypeBaseURLTest && (
                            <p className="mt-2 text-xs font-mono text-orange-800 break-all">
                                URL sandbox: {integrationTypeBaseURLTest}
                            </p>
                        )}
                    </Alert>
                )}
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
                    className="min-w-[200px] bg-green-600 hover:bg-green-700 text-white font-semibold"
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
                            <TruckIcon className="w-5 h-5 mr-2" />
                            Crear Integracion
                        </>
                    )}
                </Button>
            </div>
        </form>
    );
}
