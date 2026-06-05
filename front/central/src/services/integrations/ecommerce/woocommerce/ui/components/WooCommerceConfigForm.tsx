'use client';

import { useState, FormEvent, useEffect } from 'react';
import { Button, Input, Alert, Select, Modal, SecretInput } from '@/shared/ui';
import { WooCommerceCredentials, WooCommerceConfig } from '../../domain/types';
import { createIntegrationAction, updateIntegrationAction, testConnectionRawAction } from '@/services/integrations/core/infra/actions';
import { useToast } from '@/shared/providers/toast-provider';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import {
    KeyIcon,
    Cog6ToothIcon,
    ShoppingBagIcon,
    InformationCircleIcon,
    ArrowLeftIcon,
} from '@heroicons/react/24/outline';

interface WooCommerceConfigFormProps {
    onSuccess?: () => void;
    onCancel?: () => void;
    isEdit?: boolean;
    integrationId?: number;
    initialData?: {
        name?: string;
        store_id?: string;
        config?: any;
        credentials?: any;
        business_id?: number | null;
    };
}

const fieldLabel = 'block text-sm font-medium text-gray-700 dark:text-gray-200 mb-2';
const fieldHint = 'text-xs text-gray-500 dark:text-gray-400 mt-1.5 flex items-start gap-1';

export function WooCommerceConfigForm({ onSuccess, onCancel, isEdit, integrationId, initialData }: WooCommerceConfigFormProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(false);
    const [testingConnection, setTestingConnection] = useState(false);
    const [errorModal, setErrorModal] = useState<string | null>(null);
    const [isSuperAdmin, setIsSuperAdmin] = useState(false);
    const [businesses, setBusinesses] = useState<Array<{ id: number; name: string }>>([]);
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [loadingBusinesses, setLoadingBusinesses] = useState(false);

    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        store_url: initialData?.config?.store_url || initialData?.store_id || '',
        consumer_key: initialData?.credentials?.consumer_key || '',
        consumer_secret: initialData?.credentials?.consumer_secret || '',
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
                if (permissions?.business_id) {
                    setSelectedBusinessId(permissions.business_id);
                }
            }
        };

        checkUserAndLoadBusinesses();
    }, []);

    const connectionReady = !!formData.store_url && !!formData.consumer_key && !!formData.consumer_secret;

    const handleTestConnection = async () => {
        if (!connectionReady) {
            showToast('Debes ingresar la URL de la tienda, Consumer Key y Consumer Secret para probar la conexion', 'warning');
            return;
        }

        setTestingConnection(true);

        try {
            const credentials: WooCommerceCredentials = {
                consumer_key: formData.consumer_key,
                consumer_secret: formData.consumer_secret,
            };

            const config: WooCommerceConfig = {
                store_url: formData.store_url,
            };

            const result = await testConnectionRawAction('woocommerce', config, credentials);

            if (result.success) {
                showToast('Conexion exitosa con WooCommerce', 'success');
            } else {
                throw new Error(result.message || 'Error al probar conexion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al conectar con WooCommerce');
        } finally {
            setTestingConnection(false);
        }
    };

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setLoading(true);

        try {
            const config: WooCommerceConfig = {
                store_url: formData.store_url,
            };

            if (isEdit && integrationId) {
                const credentials: any = {};
                if (formData.consumer_key) credentials.consumer_key = formData.consumer_key;
                if (formData.consumer_secret) credentials.consumer_secret = formData.consumer_secret;

                const response: any = await updateIntegrationAction(integrationId, {
                    name: formData.name,
                    store_id: formData.store_url,
                    config: config as any,
                    credentials: Object.keys(credentials).length > 0 ? credentials : undefined,
                });

                if (!response || response.success === false) {
                    throw new Error(response?.message || 'Error al actualizar integracion');
                }
                showToast('Integracion WooCommerce actualizada', 'success');
                onSuccess?.();
                return;
            }

            if (isSuperAdmin && !selectedBusinessId) {
                setErrorModal('Debes seleccionar un negocio antes de crear la integracion.');
                setLoading(false);
                return;
            }

            const credentials: WooCommerceCredentials = {
                consumer_key: formData.consumer_key,
                consumer_secret: formData.consumer_secret,
            };

            const response = await createIntegrationAction({
                name: formData.name,
                code: `woocommerce_${Date.now()}`,
                integration_type_id: 4,
                category: 'ecommerce',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                config: config as any,
                credentials: credentials as any,
                is_active: true,
                is_default: false,
            });

            if (response.success) {
                showToast('Integracion WooCommerce creada exitosamente', 'success');
                onSuccess?.();
            } else {
                throw new Error(response.message || 'Error al crear integracion');
            }
        } catch (err: any) {
            setErrorModal(err.message || 'Error al guardar la integracion de WooCommerce');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-6" autoComplete="off">
            <div className="flex flex-col gap-3 border-b border-gray-200 dark:border-gray-700 pb-5 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-center gap-3">
                    <div className="p-2 bg-purple-50 rounded-lg">
                        <ShoppingBagIcon className="w-6 h-6 text-purple-600" />
                    </div>
                    <div>
                        <h2 className="text-2xl font-bold text-gray-900 dark:text-white">WooCommerce</h2>
                        <p className="text-sm text-gray-600 dark:text-gray-300">
                            Conecta tu tienda para sincronizar ordenes y productos a Probability.
                        </p>
                    </div>
                </div>
                <span
                    className={`inline-flex items-center gap-2 self-start rounded-full px-3 py-1 text-xs font-semibold ${
                        connectionReady
                            ? 'bg-green-50 text-green-700 dark:bg-green-900/30 dark:text-green-300'
                            : 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-300'
                    }`}
                >
                    <span className={`h-2 w-2 rounded-full ${connectionReady ? 'bg-green-500' : 'bg-gray-400'}`} />
                    {connectionReady ? 'Listo para probar' : 'Datos incompletos'}
                </span>
            </div>

            <div className="space-y-6">
                <div className="bg-gray-50 dark:bg-gray-700 rounded-xl p-6">
                        <div className="flex items-center gap-2 mb-5">
                            <Cog6ToothIcon className="w-5 h-5 text-gray-700 dark:text-gray-200" />
                            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Configuracion General</h3>
                        </div>

                        <div className="grid grid-cols-1 gap-x-5 gap-y-4 md:grid-cols-2">
                            <div>
                                <label className={fieldLabel}>
                                    Nombre de la Integracion <span className="text-red-500">*</span>
                                </label>
                                <Input
                                    type="text"
                                    value={formData.name}
                                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                                    placeholder="Ej: WooCommerce Tienda Principal"
                                    required
                                    className="bg-white dark:bg-gray-800"
                                />
                                <p className={fieldHint}>
                                    <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                    <span>Nombre descriptivo para identificar esta integracion</span>
                                </p>
                            </div>

                            <div>
                                <label className={fieldLabel}>
                                    URL de la Tienda <span className="text-red-500">*</span>
                                </label>
                                <Input
                                    type="url"
                                    value={formData.store_url}
                                    onChange={(e) => setFormData({ ...formData, store_url: e.target.value })}
                                    placeholder="https://mitienda.com"
                                    required
                                    autoComplete="off"
                                    className="bg-white dark:bg-gray-800 text-sm"
                                />
                                <p className={fieldHint}>
                                    <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                    <span>URL completa de tu tienda WooCommerce (sin barra final)</span>
                                </p>
                            </div>

                            {isSuperAdmin && !isEdit && (
                                <div className="md:col-span-2">
                                    <label className={fieldLabel}>
                                        Negocio <span className="text-red-500">*</span>
                                    </label>
                                    {loadingBusinesses ? (
                                        <div className="flex items-center gap-2 p-3 bg-white dark:bg-gray-800 rounded-lg">
                                            <svg className="animate-spin h-5 w-5 text-purple-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                            </svg>
                                            <span className="text-sm text-gray-600 dark:text-gray-300">Cargando negocios...</span>
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
                                            className="bg-white dark:bg-gray-800"
                                        />
                                    )}
                                    <p className={fieldHint}>
                                        <InformationCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" />
                                        <span>Negocio al que pertenecera esta integracion</span>
                                    </p>
                                </div>
                            )}
                        </div>
                    </div>

                    <div className="bg-gradient-to-br from-purple-50 to-indigo-50 dark:from-purple-900/20 dark:to-indigo-900/20 rounded-xl p-6 border border-purple-100 dark:border-purple-900/40">
                        <div className="flex items-center gap-2 mb-5">
                            <KeyIcon className="w-5 h-5 text-purple-700 dark:text-purple-300" />
                            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Credenciales de Acceso</h3>
                        </div>

                        <div className="grid grid-cols-1 gap-x-5 gap-y-4 md:grid-cols-2">
                            <div>
                                <label className={fieldLabel}>
                                    Consumer Key <span className="text-red-500">*</span>
                                </label>
                                <SecretInput
                                    value={formData.consumer_key}
                                    onChange={(e) => setFormData({ ...formData, consumer_key: e.target.value })}
                                    placeholder="ck_xxxxxxxxxxxxxxxxxxxx"
                                    required={!isEdit}
                                    className="bg-white dark:bg-gray-800 font-mono text-sm"
                                />
                            </div>

                            <div>
                                <label className={fieldLabel}>
                                    Consumer Secret <span className="text-red-500">*</span>
                                </label>
                                <SecretInput
                                    value={formData.consumer_secret}
                                    onChange={(e) => setFormData({ ...formData, consumer_secret: e.target.value })}
                                    placeholder="cs_xxxxxxxxxxxxxxxxxxxx"
                                    required={!isEdit}
                                    className="bg-white dark:bg-gray-800 font-mono text-sm"
                                />
                            </div>
                        </div>

                        <div className="pt-5">
                            <Button
                                type="button"
                                variant="outline"
                                className="w-full bg-white dark:bg-gray-800 hover:bg-purple-50 border-purple-200 text-purple-700 dark:text-purple-300 font-semibold"
                                onClick={handleTestConnection}
                                disabled={testingConnection || loading || !connectionReady}
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
                                        <ShoppingBagIcon className="w-4 h-4 mr-2" />
                                        Probar Conexion
                                    </>
                                )}
                            </Button>
                        </div>
                    </div>

                    <div className="bg-purple-50 dark:bg-purple-900/20 border border-purple-100 dark:border-purple-900/40 rounded-xl p-5">
                        <div className="flex flex-col gap-1 mb-4 sm:flex-row sm:items-center sm:justify-between">
                            <h4 className="text-sm font-semibold text-purple-900 dark:text-purple-200 flex items-center gap-2">
                                <InformationCircleIcon className="w-5 h-5" />
                                Como obtener tus credenciales
                            </h4>
                            <p className="text-xs text-purple-800 dark:text-purple-200/80">
                                En tu panel de WordPress: <strong>WooCommerce &rarr; Ajustes &rarr; Avanzado &rarr; REST API</strong>
                            </p>
                        </div>
                        <ol className="grid grid-cols-1 gap-3 sm:grid-cols-2 lg:grid-cols-5">
                            {[
                                'Ingresa al panel de WordPress',
                                'Ve a Ajustes / Avanzado / REST API',
                                'Haz clic en "Agregar clave"',
                                'Permisos de Lectura/Escritura',
                                'Copia el Consumer Key y Secret',
                            ].map((step, i) => (
                                <li key={i} className="flex items-start gap-2 rounded-lg bg-white/60 dark:bg-gray-800/40 p-3">
                                    <span className="flex h-5 w-5 flex-shrink-0 items-center justify-center rounded-full bg-purple-600 text-[10px] font-bold text-white">
                                        {i + 1}
                                    </span>
                                    <span className="text-xs text-purple-800 dark:text-purple-200/80">{step}</span>
                                </li>
                            ))}
                        </ol>
                    </div>
                </div>

            <div className="flex flex-col-reverse gap-3 pt-6 border-t border-gray-200 dark:border-gray-700 sm:flex-row sm:justify-between sm:items-center">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        disabled={loading}
                        className="sm:min-w-[140px] bg-gray-100 dark:bg-gray-700 hover:bg-gray-200 text-gray-700 dark:text-gray-200 border border-gray-300 dark:border-gray-600"
                    >
                        <ArrowLeftIcon className="w-4 h-4 mr-2" />
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    variant="primary"
                    disabled={loading}
                    className="sm:min-w-[200px] bg-purple-600 hover:bg-purple-700 text-white font-semibold"
                >
                    {loading ? (
                        <>
                            <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            {isEdit ? 'Guardando...' : 'Conectando...'}
                        </>
                    ) : (
                        <>
                            <ShoppingBagIcon className="w-5 h-5 mr-2" />
                            {isEdit ? 'Guardar Cambios' : 'Crear Integracion'}
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
