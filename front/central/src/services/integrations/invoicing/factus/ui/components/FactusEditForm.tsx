'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select, Modal } from '@/shared/ui';
import { FactusConfig, FactusCredentials } from '../../domain/types';
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

interface FactusEditFormProps {
    integrationId: number;
    initialData: {
        name: string;
        config: FactusConfig;
        credentials?: FactusCredentials;
        business_id?: number | null;
    };
    onSuccess?: () => void;
    onCancel?: () => void;
}

export function FactusEditForm({ integrationId, initialData, onSuccess, onCancel }: FactusEditFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [showClientSecret, setShowClientSecret] = useState(false);
    const [showPassword, setShowPassword] = useState(false);

    // Business selection for super admins
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId] = useState<number | null>(initialData.business_id || null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData.name,
        numbering_range_id: initialData.config.numbering_range_id ?? '' as string | number,
        default_tax_rate: initialData.config.default_tax_rate || '19.00',
        payment_form: initialData.config.payment_form || '1',
        payment_method_code: initialData.config.payment_method_code || '10',
        legal_organization_id: initialData.config.legal_organization_id || '2',
        tribute_id: initialData.config.tribute_id || '21',
        identification_document_id: initialData.config.identification_document_id || '3',
        municipality_id: initialData.config.municipality_id || '',
        client_id: initialData.credentials?.client_id || '',
        client_secret: initialData.credentials?.client_secret || '',
        username: initialData.credentials?.username || '',
        password: initialData.credentials?.password || '',
        api_url: initialData.credentials?.api_url || '',
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
        if (!formData.client_id || !formData.client_secret || !formData.username || !formData.password) {
            showToast('Debes ingresar Client ID, Client Secret, Usuario y Contrasena para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials = {
                client_id: formData.client_id,
                client_secret: formData.client_secret,
                username: formData.username,
                password: formData.password,
                api_url: formData.api_url || undefined,
            };

            const result = await testConnectionRawAction('factus', {}, credentials);

            if (result.success) {
                showToast('Conexion exitosa con Factus', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con Factus');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const config: FactusConfig = {
                numbering_range_id: Number(formData.numbering_range_id),
                default_tax_rate: formData.default_tax_rate || undefined,
                payment_form: formData.payment_form || undefined,
                payment_method_code: formData.payment_method_code || undefined,
                legal_organization_id: formData.legal_organization_id || undefined,
                tribute_id: formData.tribute_id || undefined,
                identification_document_id: formData.identification_document_id || undefined,
                municipality_id: formData.municipality_id || undefined,
            };

            const updateData: any = {
                name: formData.name,
                config: config,
            };

            // Only include credentials if they were filled in
            if (formData.client_id && formData.client_secret && formData.username && formData.password) {
                const credentials: FactusCredentials = {
                    client_id: formData.client_id,
                    client_secret: formData.client_secret,
                    username: formData.username,
                    password: formData.password,
                    api_url: formData.api_url || undefined,
                };
                updateData.credentials = credentials;
            }

            const response = await updateIntegrationAction(integrationId, updateData);

            if (response.success) {
                showToast('Integración Factus actualizada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al actualizar integración');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al actualizar integracion de Factus');
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
                        <CheckBadgeIcon className="w-6 h-6 text-green-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900">
                        Editar Factus Facturación
                    </h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">
                    Actualiza la configuración de tu integración con Factus.
                </p>
            </div>

            {/* Configuración General */}
            <div className="bg-gray-50 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-4">
                    <Cog6ToothIcon className="w-5 h-5 text-gray-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Configuración General
                    </h3>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Nombre de la Integración <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.name}
                        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                        placeholder="Ej: Factus Facturación Principal"
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
                                <svg className="animate-spin h-5 w-5 text-green-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
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
                            <span>El negocio no puede ser modificado después de la creación</span>
                        </p>
                    </div>
                )}
            </div>

            {/* Credenciales OAuth2 */}
            <div className="bg-gradient-to-br from-green-50 to-emerald-50 rounded-xl p-6 space-y-4 border border-green-100">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-green-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Credenciales de Acceso
                    </h3>
                </div>
                <p className="text-sm text-green-900 -mt-2 mb-4 flex items-start gap-2">
                    <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                    <span>
                        Aquí se muestran las credenciales actuales. Puedes modificarlas si necesitas actualizarlas.
                        Todos los campos deben estar completos para actualizar las credenciales.
                    </span>
                </p>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Client ID
                        </label>
                        <Input
                            type="text"
                            value={formData.client_id}
                            onChange={(e) => setFormData({ ...formData, client_id: e.target.value })}
                            placeholder="Client ID de OAuth2"
                            autoComplete="off"
                            data-1p-ignore
                            className="bg-white font-mono text-sm"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Client Secret
                        </label>
                        <div className="relative">
                            <Input
                                type={showClientSecret ? "text" : "password"}
                                value={formData.client_secret}
                                onChange={(e) => setFormData({ ...formData, client_secret: e.target.value })}
                                placeholder="Client Secret de OAuth2"
                                autoComplete="new-password"
                                data-1p-ignore
                                className="bg-white font-mono text-sm pr-10"
                            />
                            <button
                                type="button"
                                onClick={() => setShowClientSecret(!showClientSecret)}
                                className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                                tabIndex={-1}
                            >
                                {showClientSecret ? (
                                    <EyeSlashIcon className="w-5 h-5" />
                                ) : (
                                    <EyeIcon className="w-5 h-5" />
                                )}
                            </button>
                        </div>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Usuario (Email)
                        </label>
                        <Input
                            type="email"
                            value={formData.username}
                            onChange={(e) => setFormData({ ...formData, username: e.target.value })}
                            placeholder="email@empresa.com"
                            autoComplete="off"
                            data-1p-ignore
                            className="bg-white text-sm"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Contraseña
                        </label>
                        <div className="relative">
                            <Input
                                type={showPassword ? "text" : "password"}
                                value={formData.password}
                                onChange={(e) => setFormData({ ...formData, password: e.target.value })}
                                placeholder="Contraseña de tu cuenta Factus"
                                autoComplete="new-password"
                                data-1p-ignore
                                className="bg-white text-sm pr-10"
                            />
                            <button
                                type="button"
                                onClick={() => setShowPassword(!showPassword)}
                                className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                                tabIndex={-1}
                            >
                                {showPassword ? (
                                    <EyeSlashIcon className="w-5 h-5" />
                                ) : (
                                    <EyeIcon className="w-5 h-5" />
                                )}
                            </button>
                        </div>
                    </div>
                </div>

                {/* URL de la API */}
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        URL de la API
                    </label>
                    <Input
                        type="url"
                        value={formData.api_url}
                        onChange={(e) => setFormData({ ...formData, api_url: e.target.value })}
                        placeholder="https://api.factus.com.co"
                        autoComplete="off"
                        className="bg-white font-mono text-sm"
                    />
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Dejar vacío para usar la URL de producción de Factus. Útil para entornos sandbox.</span>
                    </p>
                </div>

                {/* Test Connection Button */}
                <div className="pt-2">
                    <Button
                        type="button"
                        variant="outline"
                        className="w-full bg-white hover:bg-green-50 border-green-200 text-green-700 font-semibold"
                        onClick={handleTestConnection}
                        disabled={testingConnection || loading || !formData.client_id || !formData.client_secret || !formData.username || !formData.password}
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
            </div>

            {/* Configuración de Facturación */}
            <div className="bg-blue-50 rounded-xl p-6 space-y-4 border border-blue-100">
                <div className="flex items-center gap-2 mb-4">
                    <Cog6ToothIcon className="w-5 h-5 text-blue-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Configuración de Facturación
                    </h3>
                </div>
                <p className="text-sm text-blue-900 -mt-2 mb-4">
                    Datos requeridos por Factus para generar facturas electrónicas
                </p>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        ID Rango de Numeración <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="number"
                        value={formData.numbering_range_id}
                        onChange={(e) => setFormData({ ...formData, numbering_range_id: Number(e.target.value) })}
                        placeholder="Ej: 1"
                        min="1"
                        required
                        className="bg-white"
                    />
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>
                            ID del rango de numeración habilitado en Factus. Lo encuentras en
                            <strong> Configuración → Numeración</strong>
                        </span>
                    </p>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Tasa de IVA por Defecto
                        </label>
                        <Input
                            type="text"
                            value={formData.default_tax_rate}
                            onChange={(e) => setFormData({ ...formData, default_tax_rate: e.target.value })}
                            placeholder="19.00"
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5">
                            Porcentaje de IVA (por defecto: 19.00)
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Forma de Pago
                        </label>
                        <Select
                            value={formData.payment_form}
                            onChange={(e) => setFormData({ ...formData, payment_form: e.target.value })}
                            options={[
                                { value: '1', label: '1 - Contado' },
                                { value: '2', label: '2 - Crédito' },
                            ]}
                            className="bg-white"
                        />
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Método de Pago
                        </label>
                        <Select
                            value={formData.payment_method_code}
                            onChange={(e) => setFormData({ ...formData, payment_method_code: e.target.value })}
                            options={[
                                { value: '10', label: '10 - Efectivo' },
                                { value: '42', label: '42 - Consignación bancaria' },
                                { value: '47', label: '47 - Transferencia' },
                                { value: '48', label: '48 - Tarjeta crédito' },
                                { value: '49', label: '49 - Tarjeta débito' },
                            ]}
                            className="bg-white"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Organización Legal del Cliente
                        </label>
                        <Select
                            value={formData.legal_organization_id}
                            onChange={(e) => setFormData({ ...formData, legal_organization_id: e.target.value })}
                            options={[
                                { value: '1', label: '1 - Jurídica' },
                                { value: '2', label: '2 - Natural' },
                            ]}
                            className="bg-white"
                        />
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Régimen Tributario (Tribute ID)
                        </label>
                        <Input
                            type="text"
                            value={formData.tribute_id}
                            onChange={(e) => setFormData({ ...formData, tribute_id: e.target.value })}
                            placeholder="21"
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5">
                            Régimen tributario DIAN del cliente (21 = No responsable de IVA)
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Tipo de Documento del Cliente
                        </label>
                        <Select
                            value={formData.identification_document_id}
                            onChange={(e) => setFormData({ ...formData, identification_document_id: e.target.value })}
                            options={[
                                { value: '3', label: '3 - Cédula de ciudadanía' },
                                { value: '13', label: '13 - NIT' },
                                { value: '2', label: '2 - Cédula de extranjería' },
                                { value: '11', label: '11 - Registro civil' },
                                { value: '12', label: '12 - Tarjeta de identidad' },
                                { value: '31', label: '31 - NIT extranjero' },
                            ]}
                            className="bg-white"
                        />
                    </div>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        ID Municipio del Cliente
                    </label>
                    <Input
                        type="text"
                        value={formData.municipality_id}
                        onChange={(e) => setFormData({ ...formData, municipality_id: e.target.value })}
                        placeholder="Dejar vacío para usar el municipio del cliente"
                        className="bg-white"
                    />
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>
                            ID del municipio en Factus. Solo si quieres forzar un municipio fijo.
                        </span>
                    </p>
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
                    className="min-w-[200px] bg-green-600 hover:bg-green-700 text-white font-semibold"
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
                            Actualizar Integración
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
