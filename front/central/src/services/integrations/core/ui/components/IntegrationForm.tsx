'use client';

import { useState, useEffect } from 'react';
import { createIntegrationAction, updateIntegrationAction, getActiveIntegrationTypesAction, testIntegrationAction, testConnectionRawAction, getWebhookUrlAction, activateIntegrationAction, deactivateIntegrationAction } from '../../infra/actions';
import { Integration, IntegrationType, WebhookInfo } from '../../domain/types';
import { Alert, Button } from '@/shared/ui';
import { ShopifyIntegrationForm } from '@/services/integrations/ecommerce/shopify/ui';
import { WooCommerceConfigForm } from '@/services/integrations/ecommerce/woocommerce/ui';
import { MercadoLibreConfigForm } from '@/services/integrations/ecommerce/mercadolibre/ui';
import { WhatsAppIntegrationView } from '@/services/integrations/messages/whatsapp/ui';
import { SoftpymesConfigForm, SoftpymesEditForm } from '@/services/integrations/invoicing/softpymes/ui/components';
import { FactusConfigForm, FactusEditForm } from '@/services/integrations/invoicing/factus/ui';
import { SiigoConfigForm, SiigoEditForm } from '@/services/integrations/invoicing/siigo/ui';
import { AlegraConfigForm, AlegraEditForm } from '@/services/integrations/invoicing/alegra/ui';
import { WorldOfficeConfigForm, WorldOfficeEditForm } from '@/services/integrations/invoicing/world_office/ui';
import { HelisaConfigForm, HelisaEditForm } from '@/services/integrations/invoicing/helisa/ui';
import { EnvioClickConfigForm, EnvioClickEditForm } from '@/services/integrations/transport/envioclick/ui';
import { EnviameConfigForm, EnviameEditForm } from '@/services/integrations/transport/enviame/ui';
import { TuConfigForm, TuEditForm } from '@/services/integrations/transport/tu/ui';
import { MiPaqueteConfigForm, MiPaqueteEditForm } from '@/services/integrations/transport/mipaquete/ui';
import { ShipitConfigForm, ShipitEditForm } from '@/services/integrations/transport/shipit/ui';
import { VTEXConfigForm, VTEXEditForm } from '@/services/integrations/ecommerce/vtex/ui';
import { TiendanubeOAuthForm, TiendanubeEditForm } from '@/services/integrations/ecommerce/tiendanube/ui';
import { MagentoConfigForm, MagentoEditForm } from '@/services/integrations/ecommerce/magento/ui';
import { AmazonConfigForm, AmazonEditForm } from '@/services/integrations/ecommerce/amazon/ui';
import { FalabellaConfigForm, FalabellaEditForm } from '@/services/integrations/ecommerce/falabella/ui';
import { ExitoConfigForm, ExitoEditForm } from '@/services/integrations/ecommerce/exito/ui';
import { TikTokConfigForm, TikTokEditForm } from '@/services/integrations/ecommerce/tiktok/ui';
import { JumpsellerConfigForm } from '@/services/integrations/ecommerce/jumpseller/ui';
import { BoldConfigForm, BoldEditForm } from '@/services/integrations/pay/bold/ui/components';
import { ChatBubbleLeftRightIcon } from '@heroicons/react/24/outline';
import { getActionError } from '@/shared/utils/action-result';

const INTEGRATION_TYPE_IDS = {
    SHOPIFY: 1,
    WHATSAPP: 2,
    MERCADO_LIBRE: 3,
    WOOCOMMERCE: 4,
    SOFTPYMES: 5,
    FACTUS: 7,
    SIIGO: 8,
    ALEGRA: 9,
    WORLD_OFFICE: 10,
    HELISA: 11,
    ENVIOCLICK: 12,
    ENVIAME: 13,
    TU: 14,
    MIPAQUETE: 15,
    SHIPIT: 34,
    VTEX: 16,
    TIENDANUBE: 17,
    MAGENTO: 18,
    AMAZON: 19,
    FALABELLA: 20,
    EXITO: 21,
    BOLD: 23,
    JUMPSELLER: 33,
    TIKTOK: 35,
} as const;

interface IntegrationFormProps {
    integration?: Integration;
    onSuccess?: () => void;
    onCancel?: () => void;
    onRefresh?: () => void;
    onTypeSelected?: (hasTypeSelected: boolean) => void;
}

export default function IntegrationForm({ integration, onSuccess, onCancel, onRefresh, onTypeSelected }: IntegrationFormProps) {
    const [integrationTypes, setIntegrationTypes] = useState<IntegrationType[]>([]);
    const [selectedType, setSelectedType] = useState<IntegrationType | null>(null);
    const [loadingTypes, setLoadingTypes] = useState(true);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (onTypeSelected) {
            onTypeSelected(!!selectedType);
        }
    }, [selectedType, onTypeSelected]);

    useEffect(() => {
        const fetchIntegrationTypes = async () => {
            console.log('🔍 Fetching integration types...');
            try {
                const response = await getActiveIntegrationTypesAction();
                console.log('📦 Integration types response:', response);

                if (response.success && response.data) {
                    console.log('✅ Integration types loaded:', response.data);
                    setIntegrationTypes(response.data);

                    if (integration) {
                        const type = response.data.find(t => t.id === integration.integration_type_id);
                        setSelectedType(type || null);
                    }

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
    }, [integration?.id, integration?.integration_type_id]);

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

    const [whatsappName, setWhatsappName] = useState('');
    const [creatingWhatsapp, setCreatingWhatsapp] = useState(false);

    const handleWhatsAppCreate = async () => {
        if (!selectedType || !whatsappName.trim()) return;
        setCreatingWhatsapp(true);
        setError(null);
        try {
            const code = whatsappName.trim().toLowerCase().replace(/\s+/g, '_');
            const result = await createIntegrationAction({
                name: whatsappName.trim(),
                code,
                integration_type_id: selectedType.id,
                category: selectedType.category?.code || selectedType.integration_category?.code || 'messaging',
                business_id: null,
                is_active: true,
                is_default: false,
            });
            if (result.success) {
                onSuccess?.();
            } else {
                setError(result.message || 'Error al crear la integración de WhatsApp');
            }
        } catch (err: any) {
            setError(getActionError(err, 'Error al crear la integración de WhatsApp'));
        } finally {
            setCreatingWhatsapp(false);
        }
    };

    const handleWhatsAppToggleActive = async (id: number, currentlyActive: boolean) => {
        try {
            if (currentlyActive) {
                const result = await deactivateIntegrationAction(id);
                return result.success;
            } else {
                const result = await activateIntegrationAction(id);
                return result.success;
            }
        } catch {
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
        code?: string;
        store_id: string;
        config: any;
        credentials: any;
        business_id?: number | null;
        is_testing?: boolean;
    }) => {
        if (!integration) return;

        try {
            const updateData: any = {
                name: data.name,
                store_id: data.store_id,
                config: data.config,
                is_testing: data.is_testing,
            };
            if (data.code) {
                updateData.code = data.code;
            }

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
            setError(getActionError(err, 'Error al actualizar la integración'));
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

    if (integration) {
        console.log('📋 Integration recibida para editar:', integration);

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
                        is_testing: integration.is_testing,
                        base_url_test: selectedType?.base_url_test,
                    }}
                    isEdit={true}
                    integrationId={integration.id}
                />
            );
        }

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.WOOCOMMERCE) {
            return (
                <WooCommerceConfigForm
                    isEdit={true}
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        store_id: integration.store_id,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                        is_testing: integration.is_testing,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.JUMPSELLER) {
            return (
                <JumpsellerConfigForm
                    isEdit={true}
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        store_id: integration.store_id,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                        is_testing: integration.is_testing,
                    }}
                    integrationTypeBaseURLTest={selectedType?.base_url_test}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.MERCADO_LIBRE) {
            return (
                <MercadoLibreConfigForm
                    isEdit={true}
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        store_id: integration.store_id,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        console.log('🔍 Verificando tipo de integración:', {
            hasSelectedType: !!selectedType,
            selectedTypeId: selectedType?.id,
            isWhatsApp: selectedType?.id === INTEGRATION_TYPE_IDS.WHATSAPP,
        });

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.WHATSAPP) {
            return (
                <WhatsAppIntegrationView
                    integration={{
                        id: integration.id,
                        name: integration.name,
                        code: integration.code,
                        business_id: integration.business_id,
                        config: parsedConfig,
                        credentials: integration.credentials || {},
                        is_active: integration.is_active,
                        created_at: integration.created_at,
                        updated_at: integration.updated_at,
                    }}
                    imageUrl={selectedType.image_url}
                    onToggleActive={handleWhatsAppToggleActive}
                    onUpdateConfig={async (id, config) => {
                        const result = await updateIntegrationAction(id, { config });
                        return { success: result.success, message: result.message };
                    }}
                    onTestConnection={async (id) => {
                        const result = await testIntegrationAction(id);
                        return { success: result.success, message: result.message };
                    }}
                    onRefresh={onRefresh}
                />
            );
        }

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
                        is_testing: integration.is_testing,
                        base_url_test: selectedType.base_url_test,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.SIIGO) {
            return (
                <SiigoEditForm
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                        is_testing: integration.is_testing,
                        base_url_test: selectedType.base_url_test,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.ALEGRA) {
            return (
                <AlegraEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.WORLD_OFFICE) {
            return (
                <WorldOfficeEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.HELISA) {
            return (
                <HelisaEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.BOLD) {
            return (
                <BoldEditForm
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                        is_testing: integration.is_testing,
                        base_url_test: selectedType.base_url_test,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.ENVIOCLICK) {
            console.log('✅ Usando EnvioClickEditForm');
            return (
                <EnvioClickEditForm
                    integrationId={integration.id}
                    initialData={{
                        name: integration.name,
                        config: parsedConfig as any,
                        credentials: integration.credentials as any,
                        business_id: integration.business_id,
                        is_testing: integration.is_testing,
                        base_url_test: selectedType.base_url_test,
                    }}
                    onSuccess={onSuccess}
                    onCancel={onCancel}
                />
            );
        }

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.ENVIAME) {
            return (
                <EnviameEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.TU) {
            return (
                <TuEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.MIPAQUETE) {
            return (
                <MiPaqueteEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.SHIPIT) {
            return (
                <ShipitEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.VTEX) {
            return (
                <VTEXEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.TIENDANUBE) {
            return (
                <TiendanubeEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.MAGENTO) {
            return (
                <MagentoEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.AMAZON) {
            return (
                <AmazonEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.FALABELLA) {
            return (
                <FalabellaEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.EXITO) {
            return (
                <ExitoEditForm
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

        if (selectedType && selectedType.id === INTEGRATION_TYPE_IDS.TIKTOK) {
            return (
                <TikTokEditForm
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

    return (
        <div className="space-y-6 w-full max-w-full overflow-x-hidden">

            {!selectedType && integrationTypes.length > 0 && (
                <div className="bg-white dark:bg-gray-800 p-4 rounded-lg w-full">
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-2">
                        Selecciona el tipo de integración *
                    </label>
                    <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-4 w-full max-w-full overflow-x-hidden">
                        {integrationTypes.map((type) => (
                            <button
                                key={type.id}
                                type="button"
                                onClick={() => handleTypeChange(type.id)}
                                className="p-4 border-2 rounded-lg text-center transition-all hover:border-blue-300 hover:shadow-lg border-gray-200 dark:border-gray-700 w-full h-full flex flex-col justify-center items-center min-h-[140px] shadow-md"
                            >

                                <div className="flex items-center justify-center mb-4">
                                    <div className="flex-shrink-0">
                                        {type.image_url ? (
                                            <img
                                                src={type.image_url}
                                                alt={type.name}
                                                className="w-14 h-14 object-contain rounded-lg shadow-md"
                                                onError={(e) => {

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
                                                    <div className="w-14 h-14 flex items-center justify-center bg-gray-100 dark:bg-gray-700 rounded-lg text-gray-400 text-base font-semibold shadow-md">
                                                        {type.name.charAt(0).toUpperCase()}
                                                    </div>
                                                )}
                                            </>
                                        )}
                                    </div>
                                </div>

                                <div className="flex-1 flex flex-col justify-center items-center">
                                    <h4 className="font-semibold text-gray-900 dark:text-white text-base break-words mb-1">{type.name}</h4>
                                    <p className="text-sm text-gray-500 dark:text-gray-400 dark:text-gray-400 break-words">{type.code}</p>
                                </div>
                            </button>
                        ))}
                    </div>
                </div>
            )}

            {!selectedType && integrationTypes.length === 0 && (
                <div className="text-center py-8">
                    <p className="text-gray-600 dark:text-gray-300 dark:text-gray-300">No hay tipos de integración disponibles.</p>
                </div>
            )}

            {selectedType && (
                <div>
                    {selectedType.id === INTEGRATION_TYPE_IDS.SHOPIFY && (
                        <ShopifyIntegrationForm
                            onSubmit={handleShopifySubmit}
                            onCancel={onCancel}
                            onTestConnection={handleTestConnection}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.WHATSAPP && (
                        <div className="space-y-3 w-full">
                            <div
                                className="flex flex-col gap-3 rounded-xl p-4 sm:flex-row sm:items-center dark:bg-gray-800/60"
                                style={{
                                    backgroundColor: 'color-mix(in srgb, var(--color-primary) 10%, white)',
                                    border: '1px solid color-mix(in srgb, var(--color-primary) 25%, white)',
                                }}
                            >
                                <span
                                    className="flex h-11 w-11 items-center justify-center rounded-xl overflow-hidden shrink-0 bg-white dark:bg-gray-900"
                                    style={{ border: '1px solid color-mix(in srgb, var(--color-primary) 25%, white)' }}
                                >
                                    <img
                                        src={selectedType.image_url || '/integrations/whatsapp.png'}
                                        alt="WhatsApp"
                                        className="h-8 w-8 object-contain"
                                    />
                                </span>
                                <div>
                                    <h2 className="text-base font-bold text-gray-900 dark:text-white leading-tight">
                                        Conectar WhatsApp
                                    </h2>
                                    <p className="text-xs text-gray-600 dark:text-gray-300 mt-0.5">
                                        {'Ponle nombre a la integraci\u00f3n y cr\u00e9ala. Despu\u00e9s, al editarla, conectas tu propio n\u00famero.'}
                                    </p>
                                </div>
                            </div>

                            <div
                                className="rounded-xl p-4 dark:bg-gray-800/60"
                                style={{ backgroundColor: '#fafafd', border: '1px solid #eceaf3' }}
                            >
                                <div className="flex items-center gap-2 mb-3">
                                    <span
                                        className="flex h-7 w-7 items-center justify-center rounded-md"
                                        style={{ backgroundColor: 'color-mix(in srgb, var(--color-primary) 10%, white)' }}
                                    >
                                        <ChatBubbleLeftRightIcon
                                            style={{ color: 'var(--color-primary)', width: 16, height: 16 }}
                                        />
                                    </span>
                                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">
                                        Datos de la {'integraci\u00f3n'}
                                    </h3>
                                </div>

                                <label className="block text-[13px] font-semibold text-gray-900 dark:text-gray-100 mb-1">
                                    Nombre de la {'integraci\u00f3n'}{' '}
                                    <span style={{ color: 'var(--color-primary)' }}>*</span>
                                </label>
                                <input
                                    type="text"
                                    value={whatsappName}
                                    onChange={(e) => setWhatsappName(e.target.value)}
                                    placeholder="Ej: WhatsApp - Mi Negocio"
                                    className="w-full px-3 py-2 text-sm rounded-lg border bg-white dark:bg-gray-800 text-gray-900 dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[var(--color-primary)]/30 focus:border-[var(--color-primary)]"
                                    style={{ borderColor: '#e9e9f0' }}
                                    disabled={creatingWhatsapp}
                                />
                                <p className="text-[11px] text-gray-400 dark:text-gray-500 mt-1">
                                    {'Las notificaciones a tus clientes se configuran despu\u00e9s, en el m\u00f3dulo de Notificaciones.'}
                                </p>
                            </div>

                            <div className="flex items-center justify-end gap-2 pt-1">
                                {onCancel && (
                                    <button
                                        type="button"
                                        onClick={onCancel}
                                        disabled={creatingWhatsapp}
                                        className="px-5 py-2 text-[13px] font-semibold rounded-lg bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                                        style={{ border: '1px solid #e9e9f0' }}
                                    >
                                        Cancelar
                                    </button>
                                )}
                                <button
                                    type="button"
                                    onClick={handleWhatsAppCreate}
                                    disabled={creatingWhatsapp || !whatsappName.trim()}
                                    className="px-5 py-2 text-[13px] font-semibold rounded-lg text-white flex items-center justify-center gap-2 transition-colors disabled:opacity-60"
                                    style={{ backgroundColor: 'var(--color-primary)' }}
                                >
                                    {creatingWhatsapp && (
                                        <span className="h-3.5 w-3.5 rounded-full border-2 border-white/40 border-t-white animate-spin" />
                                    )}
                                    Crear {'integraci\u00f3n'}
                                </button>
                            </div>
                        </div>
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.SOFTPYMES && (
                        <SoftpymesConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                            integrationTypeBaseURLTest={selectedType.base_url_test}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.FACTUS && (
                        <FactusConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.SIIGO && (
                        <SiigoConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                            integrationTypeBaseURLTest={selectedType.base_url_test}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.ALEGRA && (
                        <AlegraConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.WORLD_OFFICE && (
                        <WorldOfficeConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.HELISA && (
                        <HelisaConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.BOLD && (
                        <BoldConfigForm
                            integrationTypeId={selectedType.id}
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                            integrationTypeBaseURLTest={selectedType.base_url_test}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.ENVIOCLICK && (
                        <EnvioClickConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                            integrationTypeBaseURL={selectedType.base_url}
                            integrationTypeBaseURLTest={selectedType.base_url_test}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.ENVIAME && (
                        <EnviameConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.TU && (
                        <TuConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.MIPAQUETE && (
                        <MiPaqueteConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.SHIPIT && (
                        <ShipitConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.VTEX && (
                        <VTEXConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.TIENDANUBE && (
                        <TiendanubeOAuthForm
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.MAGENTO && (
                        <MagentoConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.AMAZON && (
                        <AmazonConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.FALABELLA && (
                        <FalabellaConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.EXITO && (
                        <ExitoConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.id === INTEGRATION_TYPE_IDS.JUMPSELLER && (
                        <JumpsellerConfigForm
                            onSuccess={onSuccess}
                            onCancel={onCancel}
                            integrationTypeBaseURLTest={selectedType.base_url_test}
                        />
                    )}

                    {selectedType.id !== INTEGRATION_TYPE_IDS.SHOPIFY &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.WHATSAPP &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.SOFTPYMES &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.FACTUS &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.SIIGO &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.ALEGRA &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.WORLD_OFFICE &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.HELISA &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.ENVIOCLICK &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.ENVIAME &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.TU &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.MIPAQUETE &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.SHIPIT &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.VTEX &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.TIENDANUBE &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.MAGENTO &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.AMAZON &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.FALABELLA &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.EXITO &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.JUMPSELLER &&
                     selectedType.id !== INTEGRATION_TYPE_IDS.BOLD && (
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
