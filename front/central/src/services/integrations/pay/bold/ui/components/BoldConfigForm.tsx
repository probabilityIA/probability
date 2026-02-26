'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select } from '@/shared/ui';
import { BoldCredentials } from '../../domain/types';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    CreditCardIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    EyeIcon,
    EyeSlashIcon,
    BeakerIcon,
} from '@heroicons/react/24/outline';

interface BoldConfigFormProps {
    integrationTypeId: number;
    onSuccess?: () => void;
    onCancel?: () => void;
    integrationTypeBaseURLTest?: string;
}

export function BoldConfigForm({ integrationTypeId, onSuccess, onCancel, integrationTypeBaseURLTest }: BoldConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showApiKey, setShowApiKey] = useState(false);

    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: '',
        api_key: '',
        environment: 'sandbox' as 'sandbox' | 'production',
    });

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
                const permissions = TokenStorage.getPermissions();
                if (permissions?.business_id) {
                    setSelectedBusinessId(permissions.business_id);
                }
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

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

            const credentials: BoldCredentials = {
                api_key: formData.api_key,
                environment: formData.environment,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `bold_${Date.now()}`,
                integration_type_id: integrationTypeId,
                category: 'pay',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: {} as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
                is_testing: formData.environment === 'sandbox',
            });

            if (response.success) {
                showToast('Integración Bold creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al crear integración');
            showToast('Error al crear integración Bold', 'error');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-8" autoComplete="off">
            {/* Header */}
            <div className="border-b border-gray-200 pb-6">
                <div className="flex items-center gap-3 mb-2">
                    <div className="p-2 bg-blue-50 rounded-lg">
                        <CreditCardIcon className="w-6 h-6 text-blue-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900">Bold Pagos</h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">
                    Conecta tu cuenta de Bold para crear links de pago y recibir pagos con tarjeta, PSE y más.
                </p>
            </div>

            {error && <Alert type="error">{error}</Alert>}

            {/* Configuración General */}
            <div className="bg-gray-50 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-4">
                    <Cog6ToothIcon className="w-5 h-5 text-gray-700" />
                    <h3 className="text-lg font-semibold text-gray-900">Configuración General</h3>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Nombre de la Integración <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="Ej: Bold Principal"
                        required
                        className="bg-white"
                    />
                </div>

                {isSuperAdmin && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Negocio <span className="text-red-500">*</span>
                        </label>
                        {loadingBusinesses ? (
                            <div className="flex items-center gap-2 p-3 bg-gray-100 rounded-lg">
                                <svg className="animate-spin h-5 w-5 text-blue-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
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
                                    ...businesses.map((b) => ({ value: b.id.toString(), label: b.name })),
                                ]}
                                required
                                className="bg-white"
                            />
                        )}
                    </div>
                )}
            </div>

            {/* Credenciales API */}
            <div className="bg-gradient-to-br from-blue-50 to-cyan-50 rounded-xl p-6 space-y-4 border border-blue-100">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-blue-700" />
                    <h3 className="text-lg font-semibold text-gray-900">Credenciales API</h3>
                </div>
                <p className="text-sm text-blue-900 -mt-2 mb-4 flex items-start gap-2">
                    <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                    <span>
                        Obtén tu API Key desde el panel de Bold en <strong>Configuración → Integraciones → API</strong>.
                    </span>
                </p>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Ambiente <span className="text-red-500">*</span>
                    </label>
                    <Select
                        value={formData.environment}
                        onChange={(e) => setFormData({ ...formData, environment: e.target.value as 'sandbox' | 'production' })}
                        options={[
                            { value: 'sandbox', label: 'Sandbox (Pruebas)' },
                            { value: 'production', label: 'Producción' },
                        ]}
                        className="bg-white"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        API Key <span className="text-red-500">*</span>
                    </label>
                    <div className="relative">
                        <Input
                            type={showApiKey ? 'text' : 'password'}
                            value={formData.api_key}
                            onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                            placeholder="Ingresa tu API Key de Bold"
                            required
                            autoComplete="new-password"
                            data-1p-ignore
                            className="bg-white font-mono text-sm pr-10"
                        />
                        <button
                            type="button"
                            onClick={() => setShowApiKey(!showApiKey)}
                            className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                            tabIndex={-1}
                        >
                            {showApiKey ? <EyeSlashIcon className="w-5 h-5" /> : <EyeIcon className="w-5 h-5" />}
                        </button>
                    </div>
                </div>

                <div className="bg-blue-100 border border-blue-200 rounded-lg p-4 mt-4">
                    <h4 className="text-sm font-semibold text-blue-900 mb-2 flex items-center gap-2">
                        <InformationCircleIcon className="w-4 h-4" />
                        Métodos de pago disponibles
                    </h4>
                    <ul className="text-xs text-blue-800 space-y-1 list-disc list-inside">
                        <li>Tarjeta de crédito y débito</li>
                        <li>PSE (Pagos Seguros en Línea)</li>
                        <li>Nequi</li>
                        <li>Botón Bancolombia</li>
                    </ul>
                </div>
            </div>

            {/* Modo Sandbox */}
            {formData.environment === 'sandbox' && (
                <div className="bg-orange-50 rounded-xl p-4 border border-orange-200">
                    <div className="flex items-start gap-2">
                        <BeakerIcon className="w-5 h-5 text-orange-600 flex-shrink-0 mt-0.5" />
                        <div>
                            <p className="text-sm font-semibold text-orange-900">Modo Sandbox activado</p>
                            <p className="text-xs text-orange-700 mt-0.5">
                                Los links de pago serán creados en el ambiente de pruebas. No se realizarán cobros reales.
                                {integrationTypeBaseURLTest && (
                                    <span className="block mt-1 font-mono break-all">{integrationTypeBaseURLTest}</span>
                                )}
                            </p>
                        </div>
                    </div>
                </div>
            )}

            {/* Buttons */}
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
                    className="min-w-[200px] bg-blue-600 hover:bg-blue-700 text-white font-semibold"
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
                            <CreditCardIcon className="w-5 h-5 mr-2" />
                            Crear Integración
                        </>
                    )}
                </Button>
            </div>
        </form>
    );
}
