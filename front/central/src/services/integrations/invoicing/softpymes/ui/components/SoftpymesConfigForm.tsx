'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select } from '@/shared/ui';
import { SoftpymesConfig, SoftpymesCredentials } from '../../domain/types';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { testConnectionRawAction } from '@/services/integrations/core/infra/actions';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    BuildingOfficeIcon,
    KeyIcon,
    Cog6ToothIcon,
    CheckBadgeIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
    EyeIcon,
    EyeSlashIcon,
    BeakerIcon,
} from '@heroicons/react/24/outline';

interface SoftpymesConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    integrationTypeBaseURLTest?: string;
}

export function SoftpymesConfigForm({ onSuccess, onCancel, integrationTypeBaseURLTest }: SoftpymesConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showApiSecret, setShowApiSecret] = useState(false);
    const [isTesting, setIsTesting] = useState(false);

    // Business selection for super admins
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: '',
        company_nit: '',
        company_name: '',
        referer: '', // Identificación de la instancia del cliente (requerido por API)
        default_customer_nit: '', // NIT por defecto para clientes sin DNI
        resolution_id: '' as string | number, // ID de resolución de facturación
        branch_code: '001', // Código de sucursal
        customer_branch_code: '000', // Código de sucursal del cliente (Softpymes default)
        seller_nit: '', // NIT del vendedor
        api_key: '',
        api_secret: '',
    });

    // Check if user is super admin and load businesses
    useEffect(() => {
        const checkUserAndLoadBusinesses = async () => {
            const permissions = TokenStorage.getPermissions();
            const isSuperUser = permissions?.is_super || false;
            setIsSuperAdmin(isSuperUser);

            // If super admin, load businesses for selection
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
                // If not super admin, auto-select current business
                if (permissions?.business_id) {
                    setSelectedBusinessId(permissions.business_id);
                }
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    const handleTestConnection = async () => {
        // Validar que se hayan ingresado las credenciales y referer
        if (!formData.api_key || !formData.api_secret) {
            showToast('Debes ingresar API Key y API Secret para probar la conexión', 'warning');
            return;
        }

        if (!formData.referer) {
            showToast('Debes ingresar el Referer (identificación de instancia) para probar la conexión', 'warning');
            return;
        }

        setTestingConnection(true);
        setError(null);

        try {
            const config = {
                company_nit: formData.company_nit,
                company_name: formData.company_name,
                referer: formData.referer,
            };

            const credentials = {
                api_key: formData.api_key,
                api_secret: formData.api_secret,
            };

            const result = await testConnectionRawAction('softpymes', config, credentials);

            if (result.success) {
                showToast('✅ Conexión exitosa con Softpymes', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexión');
            }
        } catch (err: any) {
            setError(err.message || 'Error al probar conexión');
            showToast('❌ Error al conectar con Softpymes: ' + err.message, 'error');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try {
            // Validate business selection for super admins
            if (isSuperAdmin && !selectedBusinessId) {
                showToast('Debes seleccionar un negocio', 'warning');
                setLoading(false);
                return;
            }

            const config: SoftpymesConfig = {
                company_nit: formData.company_nit,
                company_name: formData.company_name,
                referer: formData.referer,
                default_customer_nit: formData.default_customer_nit || undefined,
                resolution_id: formData.resolution_id ? Number(formData.resolution_id) : undefined,
                branch_code: formData.branch_code || undefined,
                customer_branch_code: formData.customer_branch_code || undefined,
                seller_nit: formData.seller_nit || undefined,
            };

            const credentials: SoftpymesCredentials = {
                api_key: formData.api_key,
                api_secret: formData.api_secret,
            };

            // Create integration
            // business_id is sent explicitly if super admin selected one,
            // otherwise backend will use the one from JWT
            const response = await createIntegrationAction({
                name: formData.name,
                code: `softpymes_${Date.now()}`,
                integration_type_id: 5, // Softpymes integration type ID (from database)
                category: 'invoicing', // Softpymes es una integración de facturación
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
                is_testing: isTesting,
            });

            if (response.success) {
                showToast('Integración Softpymes creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al crear integración');
            showToast('Error al crear integración Softpymes', 'error');
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
                        <CheckBadgeIcon className="w-6 h-6 text-blue-600" />
                    </div>
                    <h2 className="text-2xl font-bold text-gray-900">
                        Softpymes Facturación
                    </h2>
                </div>
                <p className="text-sm text-gray-600 ml-14">
                    Conecta tu cuenta de Softpymes para gestionar facturación electrónica automáticamente desde Probability.
                </p>
            </div>

            {error && (
                <Alert type="error">
                    {error}
                </Alert>
            )}

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
                        placeholder="Ej: Softpymes Facturación Principal"
                        required
                        className="bg-white"
                    />
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>Nombre descriptivo para identificar esta integración en el sistema</span>
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
                            <span>Selecciona el negocio al que pertenecerá esta integración</span>
                        </p>
                    </div>
                )}
            </div>

            {/* Información de la Empresa */}
            <div className="bg-gray-50 rounded-xl p-6 space-y-4">
                <div className="flex items-center gap-2 mb-4">
                    <BuildingOfficeIcon className="w-5 h-5 text-gray-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Información de la Empresa
                    </h3>
                </div>
                <p className="text-sm text-gray-600 -mt-2 mb-4">
                    Datos de tu empresa registrados en Softpymes
                </p>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            NIT de la Empresa <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={formData.company_nit}
                            onChange={(e) => setFormData({ ...formData, company_nit: e.target.value })}
                            placeholder="900123456-7"
                            required
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5">
                            Incluye el dígito de verificación
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Razón Social <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={formData.company_name}
                            onChange={(e) => setFormData({ ...formData, company_name: e.target.value })}
                            placeholder="Mi Empresa SAS"
                            required
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5">
                            Nombre registrado en Softpymes
                        </p>
                    </div>
                </div>

                {/* Referer - Identificación de instancia */}
                <div className="mt-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Referer / Identificación de Instancia <span className="text-red-500">*</span>
                    </label>
                    <Input
                        type="text"
                        value={formData.referer}
                        onChange={(e) => setFormData({ ...formData, referer: e.target.value })}
                        placeholder="probability-empresa-123"
                        required
                        className="bg-white"
                    />
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>
                            Identificador único de tu instancia en Softpymes (requerido por la API para el header Referer)
                        </span>
                    </p>
                </div>

                {/* NIT por defecto para clientes sin DNI */}
                <div className="mt-4">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        NIT por Defecto para Clientes sin DNI
                    </label>
                    <Input
                        type="text"
                        value={formData.default_customer_nit}
                        onChange={(e) => setFormData({ ...formData, default_customer_nit: e.target.value })}
                        placeholder="222222222222"
                        className="bg-white"
                    />
                    <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                        <span>
                            NIT que se usará cuando un cliente no tenga DNI. En Colombia, el consumidor final es <strong>222222222222</strong>
                        </span>
                    </p>
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
                    Datos requeridos por Softpymes para generar facturas electrónicas
                </p>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Resolution ID <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="number"
                            value={formData.resolution_id}
                            onChange={(e) => setFormData({ ...formData, resolution_id: e.target.value })}
                            placeholder="1"
                            min="1"
                            required
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>
                                ID de la resolución de facturación DIAN. Se obtiene desde el panel de Softpymes o consultando el endpoint <strong>/resolutions</strong>
                            </span>
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Código de Sucursal
                        </label>
                        <Input
                            type="text"
                            value={formData.branch_code}
                            onChange={(e) => setFormData({ ...formData, branch_code: e.target.value })}
                            placeholder="001"
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5">
                            Código de sucursal en Softpymes (por defecto: 001)
                        </p>
                    </div>
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            Código Sucursal del Cliente
                        </label>
                        <Input
                            type="text"
                            value={formData.customer_branch_code}
                            onChange={(e) => setFormData({ ...formData, customer_branch_code: e.target.value })}
                            placeholder="001"
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5">
                            Softpymes asigna &quot;000&quot; por defecto al crear clientes. Solo cambiar si tus clientes tienen otra sucursal.
                        </p>
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            NIT del Vendedor
                        </label>
                        <Input
                            type="text"
                            value={formData.seller_nit}
                            onChange={(e) => setFormData({ ...formData, seller_nit: e.target.value })}
                            placeholder="Dejar vac&iacute;o si no hay vendedores"
                            className="bg-white"
                        />
                        <p className="text-xs text-gray-500 mt-1.5 flex items-start gap-1">
                            <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                            <span>
                                Solo si tienes vendedores registrados en Softpymes. Si no tienes, <strong>d&eacute;jalo vac&iacute;o</strong> para evitar errores.
                            </span>
                        </p>
                    </div>
                </div>
            </div>

            {/* Credenciales de API */}
            <div className="bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl p-6 space-y-4 border border-blue-100">
                <div className="flex items-center gap-2 mb-4">
                    <KeyIcon className="w-5 h-5 text-blue-700" />
                    <h3 className="text-lg font-semibold text-gray-900">
                        Credenciales de API
                    </h3>
                </div>
                <p className="text-sm text-blue-900 -mt-2 mb-4 flex items-start gap-2">
                    <InformationCircleIcon className="w-5 h-5 flex-shrink-0 mt-0.5" />
                    <span>
                        Obtén tus credenciales desde el panel de administración de Softpymes en la sección
                        <strong> Configuración → API Keys</strong>
                    </span>
                </p>

                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            API Key <span className="text-red-500">*</span>
                        </label>
                        <Input
                            type="text"
                            value={formData.api_key}
                            onChange={(e) => setFormData({ ...formData, api_key: e.target.value })}
                            placeholder="Ingresa tu API Key de Softpymes"
                            required
                            autoComplete="off"
                            data-1p-ignore
                            className="bg-white font-mono text-sm"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                            API Secret <span className="text-red-500">*</span>
                        </label>
                        <div className="relative">
                            <Input
                                type={showApiSecret ? "text" : "password"}
                                value={formData.api_secret}
                                onChange={(e) => setFormData({ ...formData, api_secret: e.target.value })}
                                placeholder="Ingresa tu API Secret de Softpymes"
                                required
                                autoComplete="new-password"
                                data-1p-ignore
                                className="bg-white font-mono text-sm pr-10"
                            />
                            <button
                                type="button"
                                onClick={() => setShowApiSecret(!showApiSecret)}
                                className="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500 hover:text-gray-700 focus:outline-none"
                                tabIndex={-1}
                            >
                                {showApiSecret ? (
                                    <EyeSlashIcon className="w-5 h-5" />
                                ) : (
                                    <EyeIcon className="w-5 h-5" />
                                )}
                            </button>
                        </div>
                    </div>

                    {/* Test Connection Button */}
                    <div className="pt-2">
                        <Button
                            type="button"
                            variant="outline"
                            className="w-full bg-white hover:bg-blue-50 border-blue-200 text-blue-700 font-semibold"
                            onClick={handleTestConnection}
                            disabled={testingConnection || loading || !formData.api_key || !formData.api_secret}
                        >
                            {testingConnection ? (
                                <>
                                    <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-blue-700" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                        <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                        <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                    </svg>
                                    Probando...
                                </>
                            ) : (
                                <>
                                    <CheckBadgeIcon className="w-4 h-4 mr-2" />
                                    Probar Conexión
                                </>
                            )}
                        </Button>
                    </div>
                </div>

                <div className="bg-blue-100 border border-blue-200 rounded-lg p-4 mt-4">
                    <h4 className="text-sm font-semibold text-blue-900 mb-3 flex items-center gap-2">
                        <InformationCircleIcon className="w-4 h-4" />
                        Requisitos para facturar correctamente
                    </h4>
                    <ol className="text-xs text-blue-800 space-y-2 list-decimal list-inside ml-1">
                        <li>En Softpymes ve a <strong>Configuraci&oacute;n &rarr; API Keys</strong> y genera tu API Key y API Secret.</li>
                        <li>El <strong>Referer</strong> es el NIT de tu empresa <strong>sin d&iacute;gito de verificaci&oacute;n</strong> (ej: 901497840).</li>
                        <li>Debes tener una <strong>resoluci&oacute;n DIAN activa</strong>. Obt&eacute;n el Resolution ID desde Softpymes en <em>Facturaci&oacute;n &rarr; Resoluciones</em>.</li>
                        <li>Los <strong>SKU de tus productos</strong> en Probability deben coincidir con los c&oacute;digos de art&iacute;culos en Softpymes.</li>
                        <li>Configura un <strong>NIT por defecto</strong> para clientes sin DNI. En Colombia el consumidor final es <code className="bg-blue-200 px-1 rounded">222222222222</code>.</li>
                        <li>El <strong>NIT del vendedor</strong> solo si tienes vendedores registrados en Softpymes. Si no, <strong>d&eacute;jalo vac&iacute;o</strong>.</li>
                    </ol>
                </div>
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
                            Las facturas generadas quedarán marcadas como TEST y usarán la URL de pruebas de Softpymes.
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
                        Modo de pruebas activado. Las facturas generadas con esta integración quedarán marcadas como <strong>TEST</strong> y no serán enviadas a la DIAN.
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
                            <CheckBadgeIcon className="w-5 h-5 mr-2" />
                            Crear Integración
                        </>
                    )}
                </Button>
            </div>
        </form>
    );
}
