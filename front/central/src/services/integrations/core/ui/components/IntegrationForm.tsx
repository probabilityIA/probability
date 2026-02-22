'use client';

import { useState, useEffect } from 'react';
import { createIntegrationAction, updateIntegrationAction, getActiveIntegrationTypesAction, testIntegrationAction, testConnectionRawAction, getWebhookUrlAction } from '../../infra/actions';
import { Integration, IntegrationType, WebhookInfo } from '../../domain/types';
import { Alert } from '@/shared/ui';
import ShopifyIntegrationForm from './shopify/ShopifyIntegrationForm';
import WhatsAppIntegrationView from './whatsapp/WhatsAppIntegrationView';
import { SoftpymesConfigForm, SoftpymesEditForm } from '@/services/integrations/invoicing/softpymes/ui/components';
import { FactusConfigForm, FactusEditForm } from '@/services/integrations/invoicing/factus/ui';

// IDs constantes de tipos de integración (tabla integration_types)
const INTEGRATION_TYPE_IDS = {
    SHOPIFY: 1,
    WHATSAPP: 2,
    MERCADO_LIBRE: 3,
    WOOCOMMERCE: 4,
    SOFTPYMES: 5,
    FACTUS: 7,
} as const;

interface IntegrationFormProps {
    integration?: Integration;
    onSuccess?: () => void;
    onCancel?: () => void;
    onTypeSelected?: (hasTypeSelected: boolean) => void;
}

export default function IntegrationForm({ integration, onSuccess, onCancel, onTypeSelected }: IntegrationFormProps) {
    const [integrationTypes, setIntegrationTypes] = useState<IntegrationType[]>([]);
    const [selectedType, setSelectedType] = useState<IntegrationType | null>(null);
    const [loadingTypes, setLoadingTypes] = useState(true);
    const [error, setError] = useState<string | null>(null);

    // Notify parent when type is selected
    useEffect(() => {
        if (onTypeSelected) {
            onTypeSelected(!!selectedType);
        }
    }, [selectedType, onTypeSelected]);

    // Fetch integration types on mount
    useEffect(() => {
        const fetchIntegrationTypes = async () => {
            console.log('🔍 Fetching integration types...');
            try {
                const response = await getActiveIntegrationTypesAction();
                console.log('📦 Integration types response:', response);

                if (response.success && response.data) {
                    console.log('✅ Integration types loaded:', response.data);
                    setIntegrationTypes(response.data);

                    // Set selected type ONLY if editing an existing integration
                    if (integration) {
                        const type = response.data.find(t => t.id === integration.integration_type_id);
                        setSelectedType(type || null);
                    }
                    // Don't auto-select first type when creating new
                } else {
                    console.warn('⚠️ No integration types in response:', response);
                    setError('No se encontraron tipos de integración');
                }
            } catch (err) {
                console.error('❌ Error fetching integration types:', err);
                setError('Error al cargar los tipos de integración');
            } finally {
                setLoadingTypes(false);
            }
        };

        fetchIntegrationTypes();
    }, [integration]);

    const handleTypeChange = (typeId: number) => {
        const type = integrationTypes.find(t => t.id === typeId);
        setSelectedType(type || null);
    };

    const handleShopifySubmit = async (data: {
        name: string;
        code: string;
        store_id: string;
        config: any;
        credentials: any;
        business_id?: number | null;
    }) => {
        if (!selectedType) return;

        const integrationData = {
            name: data.name,
            code: data.code,
            store_id: data.store_id,
            integration_type_id: selectedType.id,
            category: selectedType.category?.code || selectedType.integration_category?.code || 'external',
            business_id: data.business_id || null,
            config: data.config,
            credentials: data.credentials,
            is_active: true,
            is_default: false,
        };

        await createIntegrationAction(integrationData);
        onSuccess?.();
    };

    const handleTestConnection = async (config: any, credentials: any) => {
        if (!selectedType) {
            console.error('No hay tipo de integración seleccionado');
            return false;
        }

        try {
            const result = await testConnectionRawAction(selectedType.code, config, credentials);
            if (result.success) {
                console.log('✅ Conexión probada exitosamente:', result.message);
                return true;
            } else {
                console.error('❌ Error al probar conexión:', result.message);
                return false;
            }
        } catch (error: any) {
            console.error('❌ Error al probar conexión:', error);
            return false;
        }
    };

    const handleWhatsAppTest = async () => {
        if (!integration) return false;

        try {
            const result = await testIntegrationAction(integration.id);
            return result.success;
        } catch (error) {
            return false;
        }
    };

    const handleGetWebhook = async (): Promise<WebhookInfo | null> => {
        if (!integration) return null;

        try {
            const result = await getWebhookUrlAction(integration.id);
            if (result.success && result.data) {
                return result.data;
            }
            return null;
        } catch (error) {
            console.error('Error getting webhook:', error);
            return null;
        }
    };

    const handleShopifyUpdate = async (data: {
        name: string;
        code: string;
        store_id: string;
        config: any;
        credentials: any;
        business_id?: number | null;
    }) => {
        if (!integration) return;

        try {
            const updateData: any = {
                name: data.name,
                code: data.code,
                store_id: data.store_id,
                config: data.config,
            };
            // Solo incluir credenciales si hay valores ingresados
            if (data.credentials && Object.keys(data.credentials).some(k => data.credentials[k])) {
                updateData.credentials = data.credentials;
            }
            const result = await updateIntegrationAction(integration.id, updateData);

            if (result.success) {
                onSuccess?.();
            } else {
                setError(result.message || 'Error al actualizar la integración');
            }
        } catch (err: any) {
            setError(err.message || 'Error al actualizar la integración');
        }
    };

    if (loadingTypes) {
        return <div className="text-center py-8">Cargando tipos de integración...</div>;
    }

    if (error) {
        return (
            <Alert type="error" onClose={() => setError(null)}>
                {error}
            </Alert>
        );
    }

    // If editing an existing integration
    if (integration) {
        console.log('📋 Integration recibida para editar:', integration);

        // Parse config if it's a string
        let parsedConfig = integration.config || {};
        if (typeof integration.config === 'string') {
            try {
                parsedConfig = JSON.parse(integration.config);
                console.log('✅ Config parseado en IntegrationForm:', parsedConfig);
            } catch (e) {
                console.error('❌ Error parsing config:', e);
                parsedConfig = {};
            }
        } else if (integration.config) {
            parsedConfig = integration.config;
            console.log('✅ Config ya es objeto en IntegrationForm:', parsedConfig);
        }

        // Show edit form for Shopify with webhook support
        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.SHOPIFY) {
            console.log('🛒 Editando Shopify - store_id:', integration.store_id);
            console.log('🛒 Editando Shopify - config:', parsedConfig);
            console.log('🛒 Editando Shopify - credentials:', integration.credentials);
            return (
                <ShopifyIntegrationForm
                    onSubmit={handleShopifyUpdate}
                    onCancel={onCancel}
                    onTestConnection={handleTestConnection}
                    onGetWebhook={handleGetWebhook}
                    initialData={{
                        name: integration.name,
                        code: integration.code,
                        store_id: integration.store_id,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                    }}
                    isEdit={true}
                    integrationId={integration.id}
                />
            );
        }

        // Show WhatsApp view (read-only with webhook info)
        console.log('🔍 Verificando tipo de integración:', {
            hasSelectedType: !!selectedType,
            selectedTypeId: selectedType?.id,
            isWhatsApp: selectedType?.id === INTEGRATION_TYPE_IDS.WHATSAPP,
        });

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.WHATSAPP) {
            console.log('✅ Usando WhatsAppIntegrationView');
            return (
                <WhatsAppIntegrationView
                    integration={{
                        id: integration.id,
                        name: integration.name,
                        code: integration.code,
                        config: parsedConfig,
                        credentials: integration.credentials || {},
                        is_active: integration.is_active,
                        created_at: integration.created_at,
                        updated_at: integration.updated_at,
                    }}
                    onTestConnection={handleWhatsAppTest}
                    onRefresh={onSuccess}
                />
            );
        }

        // Show edit form for Softpymes
        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.SOFTPYMES) {
            console.log('✅ Usando SoftpymesEditForm');
            return (
                <SoftpymesEditForm
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        // Show edit form for Factus
        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.FACTUS) {
            console.log('✅ Usando FactusEditForm');
            return (
                <FactusEditForm
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        // For other types that don't have a specific form yet
        return (
            <Alert type="info">
                <div className="space-y-3">
                    <p className="font-semibold">Formulario de Edición No Disponible</p>
                    <p>
                        El formulario de edición para <strong>{selectedType?.name}</strong> aún no está implementado.
                    </p>
                    <p className="text-sm">
                        Cada tipo de integración requiere su propio formulario personalizado.
                        Por favor, contacta al equipo de desarrollo para implementar este formulario.
                    </p>
                </div>
            </Alert>
        );
    }

    // Creating new integration - show type selector first if no type selected
    return (
        <div className="space-y-6 w-full max-w-full overflow-x-hidden">
            {/* Type Selector - Show when no type is selected */}
            {!selectedType && integrationTypes.length > 0 && (
                <div className="bg-white p-4 rounded-lg w-full">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Selecciona el tipo de integración *
                    </label>
                    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 w-full max-w-full overflow-x-hidden">
                        {integrationTypes.map((type) => (
                            <button
                                key={type.id}
                                type="button"
                                onClick={() => handleTypeChange(type.id)}
                                className="p-4 border-2 rounded-lg text-center transition-all hover:border-blue-300 hover:shadow-lg border-gray-200 w-full h-full flex flex-col justify-center items-center min-h-[140px] shadow-md"
                            >
                                {/* Logo centrado */}
                                <div className="flex items-center justify-center mb-4">
                                    <div className="flex-shrink-0">
                                        {type.image_url ? (
                                            <img
                                                src={type.image_url}
                                                alt={type.name}
                                                className="w-14 h-14 object-contain rounded-lg shadow-md"
                                                onError={(e) => {
                                                    // Fallback a imágenes hardcodeadas si la imagen falla
                                                    const target = e.target as HTMLImageElement;
                                                    if (type.id === INTEGRATION_TYPE_IDS.SHOPIFY) {
                                                        target.src = '/integrations/shopify.png';
                                                    } else if (type.id === INTEGRATION_TYPE_IDS.WHATSAPP) {
                                                        target.src = '/integrations/whatsapp.png';
                                                    } else {
                                                        target.style.display = 'none';
                                                    }
                                                }}
                                            />
                                        ) : (
                                            // Fallback a imágenes hardcodeadas si no hay imagen_url
                                            <>
                                                {type.id === INTEGRATION_TYPE_IDS.SHOPIFY && (
                                                    <img
                                                        src="/integrations/shopify.png"
                                                        alt="Shopify"
                                                        className="w-14 h-14 object-contain rounded-lg shadow-md"
                                                    />
                                                )}
                                                {type.id === INTEGRATION_TYPE_IDS.WHATSAPP && (
                                                    <img
                                                        src="/integrations/whatsapp.png"
                                                        alt="WhatsApp"
                                                        className="w-14 h-14 object-contain rounded-lg shadow-md"
                                                    />
                                                )}
                                                {type.id !== INTEGRATION_TYPE_IDS.SHOPIFY && type.id !== INTEGRATION_TYPE_IDS.WHATSAPP && (
                                                    <div className="w-14 h-14 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-base font-semibold shadow-md">
                                                        {type.name.charAt(0).toUpperCase()}
                                                    </div>
                                                )}
                                            </>
                                        )}
                                    </div>
                                </div>

                                {/* Contenido de texto - Nombre y código centrados */}
                                <div className="flex-1 flex flex-col justify-center items-center">
                                    <h4 className="font-semibold text-gray-900 text-base break-words mb-1">{type.name}</h4>
                                    <p className="text-sm text-gray-500 break-words">{type.code}</p>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {/* Show message if no types available */}
            {!selectedType && integrationTypes.length === 0 && (
                <div className="text-center py-8">
                    <p className="text-gray-600">No hay tipos de integración disponibles.</p>
                </div>
            )}

            {/* Render specific form based on selected type */}
            {selectedType && (
                <div>
                    {selectedType.id === INTEGRATION_TYPE_IDS.SHOPIFY && (
                        <ShopifyIntegrationForm
                            onSubmit={handleShopifySubmit}
                            onCancel={onCancel}
                            onTestConnection={handleTestConnection}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.SOFTPYMES && (
                        <SoftpymesConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.FACTUS && (
                        <FactusConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id !== INTEGRATION_TYPE_IDS.SHOPIFY &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.SOFTPYMES &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.FACTUS && (
                        <Alert type="warning">
                            <div className="space-y-3">
                                <p className="font-semibold">Formulario No Disponible</p>
                                <p>
                                    El formulario de configuración para <strong>{selectedType.name}</strong> aún no está implementado.
                                </p>
                                <p className="text-sm">
                                    Cada tipo de integración requiere su propio formulario personalizado.
                                    Por favor, selecciona una integración con formulario disponible (ej: Shopify) o contacta al equipo de desarrollo.
                                </p>
                            </div>
                        </Alert>
                    )}
                </div>
            )}
        </div>
    );
}
