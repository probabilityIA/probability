'use client';

import { useState } from 'react';
import { Modal, Alert, Button } from '@/shared/ui';
import { ChatBubbleLeftRightIcon } from '@heroicons/react/24/outline';
import { IntegrationCategory, IntegrationType } from '../../domain/types';
import { CategorySelector } from './CategorySelector';
import { ProviderSelector } from './ProviderSelector';
import { createIntegrationAction, testConnectionRawAction } from '../../infra/actions';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useBusinessesSimple } from '@/services/auth/business/ui/hooks/useBusinessesSimple';

// Importar formularios específicos por tipo de integración
import { SoftpymesConfigForm } from '@/services/integrations/invoicing/softpymes/ui/components';
import { FactusConfigForm } from '@/services/integrations/invoicing/factus/ui';
import { SiigoConfigForm } from '@/services/integrations/invoicing/siigo/ui';
import { AlegraConfigForm } from '@/services/integrations/invoicing/alegra/ui';
import { WorldOfficeConfigForm } from '@/services/integrations/invoicing/world_office/ui';
import { HelisaConfigForm } from '@/services/integrations/invoicing/helisa/ui';
import { EnvioClickConfigForm } from '@/services/integrations/transport/envioclick/ui';
import { EnviameConfigForm } from '@/services/integrations/transport/enviame/ui';
import { TuConfigForm } from '@/services/integrations/transport/tu/ui';
import { MiPaqueteConfigForm } from '@/services/integrations/transport/mipaquete/ui';
import { ShopifyIntegrationForm } from '@/services/integrations/ecommerce/shopify/ui';
import { VTEXConfigForm } from '@/services/integrations/ecommerce/vtex/ui';
import { TiendanubeConfigForm } from '@/services/integrations/ecommerce/tiendanube/ui';
import { MagentoConfigForm } from '@/services/integrations/ecommerce/magento/ui';
import { AmazonConfigForm } from '@/services/integrations/ecommerce/amazon/ui';
import { FalabellaConfigForm } from '@/services/integrations/ecommerce/falabella/ui';
import { ExitoConfigForm } from '@/services/integrations/ecommerce/exito/ui';
import { WooCommerceConfigForm } from '@/services/integrations/ecommerce/woocommerce/ui';

// IDs constantes de tipos de integración (tabla integration_types)
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
    VTEX: 16,
    TIENDANUBE: 17,
    MAGENTO: 18,
    AMAZON: 19,
    FALABELLA: 20,
    EXITO: 21,
} as const;

interface CreateIntegrationModalProps {
    isOpen: boolean;
    onClose: () => void;
    categories: IntegrationCategory[];
    onSuccess: () => void;
}

type Step = 1 | 2 | 3;

export function CreateIntegrationModal({
    isOpen,
    onClose,
    categories,
    onSuccess,
}: CreateIntegrationModalProps) {
    const [step, setStep] = useState<Step>(1);
    const [selectedCategory, setSelectedCategory] = useState<IntegrationCategory | null>(null);
    const [selectedProvider, setSelectedProvider] = useState<IntegrationType | null>(null);

    const handleCategorySelect = (category: IntegrationCategory) => {
        setSelectedCategory(category);
        setStep(2);
    };

    const handleProviderSelect = (provider: IntegrationType) => {
        setSelectedProvider(provider);
        setStep(3);
    };

    const handleSuccess = () => {
        // Reset state
        setStep(1);
        setSelectedCategory(null);
        setSelectedProvider(null);
        onSuccess();
    };

    const handleClose = () => {
        // Reset state on close
        setStep(1);
        setSelectedCategory(null);
        setSelectedProvider(null);
        onClose();
    };

    const handleBackToCategories = () => {
        setStep(1);
        setSelectedCategory(null);
        setSelectedProvider(null);
    };

    const handleBackToProviders = () => {
        setStep(2);
        setSelectedProvider(null);
    };

    // Determine modal size based on step
    const getModalSize = (): 'sm' | 'md' | 'lg' | 'xl' | '2xl' | '4xl' | '5xl' | '6xl' | '7xl' | 'full' => {
        switch (step) {
            case 1: // Category selection
                return '4xl';
            case 2: // Provider selection
                return '4xl';
            case 3: // Form
                return 'full';
            default:
                return '4xl';
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={handleClose}
            title={step === 1 ? 'Nueva Integración' : step === 2 ? selectedCategory?.name : selectedProvider?.name}
            size={getModalSize()}
        >
            {/* Step 1: Select Category */}
            {step === 1 && (
                <CategorySelector
                    categories={categories}
                    onSelect={handleCategorySelect}
                />
            )}

            {/* Step 2: Select Provider */}
            {step === 2 && selectedCategory && (
                <ProviderSelector
                    category={selectedCategory}
                    onSelect={handleProviderSelect}
                    onBack={handleBackToCategories}
                />
            )}

            {/* Step 3: Configuration Form */}
            {step === 3 && selectedProvider && (
                <FormWrapper
                    integrationType={selectedProvider}
                    onSuccess={handleSuccess}
                    onCancel={handleClose}
                    onBack={handleBackToProviders}
                />
            )}
        </Modal>
    );
}

// Internal component to wrap the form with submission logic
interface FormWrapperProps {
    integrationType: IntegrationType;
    onSuccess: () => void;
    onCancel: () => void;
    onBack: () => void;
}

function FormWrapper({ integrationType, onSuccess, onCancel, onBack }: FormWrapperProps) {
    // Si la integración está en desarrollo, mostrar mensaje
    if (integrationType.in_development) {
        return (
            <div className="p-6">
                <button
                    onClick={onBack}
                    className="text-blue-600 hover:text-blue-800 mb-6 flex items-center gap-2 font-medium transition-colors"
                >
                    <span>&larr;</span>
                    <span>Volver a proveedores</span>
                </button>
                <Alert type="info">
                    <div className="space-y-3">
                        <p className="font-semibold">Proximamente</p>
                        <p>
                            La integración con <strong>{integrationType.name}</strong> estará disponible próximamente.
                            Estamos trabajando para habilitarla lo antes posible.
                        </p>
                    </div>
                </Alert>
            </div>
        );
    }

    // Renderizar formulario específico según el ID del tipo de integración
    const renderSpecificForm = () => {
        switch (integrationType.id) {
            case INTEGRATION_TYPE_IDS.SHOPIFY:
                return (
                    <ShopifyIntegrationForm
                        onSubmit={async () => onSuccess()}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.WHATSAPP:
                return <WhatsAppActivateForm integrationType={integrationType} onSuccess={onSuccess} onBack={onBack} />;

            case INTEGRATION_TYPE_IDS.SOFTPYMES:
                return (
                    <SoftpymesConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                        integrationTypeBaseURLTest={integrationType.base_url_test}
                    />
                );

            case INTEGRATION_TYPE_IDS.FACTUS:
                return (
                    <FactusConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.SIIGO:
                return (
                    <SiigoConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.ALEGRA:
                return (
                    <AlegraConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.WORLD_OFFICE:
                return (
                    <WorldOfficeConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.HELISA:
                return (
                    <HelisaConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.ENVIOCLICK:
                return (
                    <EnvioClickConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                        integrationTypeBaseURL={integrationType.base_url}
                        integrationTypeBaseURLTest={integrationType.base_url_test}
                    />
                );

            case INTEGRATION_TYPE_IDS.ENVIAME:
                return (
                    <EnviameConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.TU:
                return (
                    <TuConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.MIPAQUETE:
                return (
                    <MiPaqueteConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.VTEX:
                return (
                    <VTEXConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.TIENDANUBE:
                return (
                    <TiendanubeConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.MAGENTO:
                return (
                    <MagentoConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.AMAZON:
                return (
                    <AmazonConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.FALABELLA:
                return (
                    <FalabellaConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.EXITO:
                return (
                    <ExitoConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case INTEGRATION_TYPE_IDS.WOOCOMMERCE:
                return (
                    <WooCommerceConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            default:
                return (
                    <Alert type="warning">
                        <div className="space-y-3">
                            <p className="font-semibold">Formulario No Disponible</p>
                            <p>
                                El formulario de configuración para <strong>{integrationType.name}</strong> aún no está implementado.
                            </p>
                            <p className="text-sm">
                                Cada tipo de integración requiere su propio formulario personalizado.
                                Por favor, contacta al equipo de desarrollo para implementar este formulario.
                            </p>
                        </div>
                    </Alert>
                );
        }
    };

    return (
        <div className="p-6">
            <button
                onClick={onBack}
                className="text-blue-600 hover:text-blue-800 mb-6 flex items-center gap-2 font-medium transition-colors"
            >
                <span>←</span>
                <span>Volver a proveedores</span>
            </button>

            {renderSpecificForm()}
        </div>
    );
}

// Formulario simple para activar WhatsApp (sin configuración, usa defaults del tipo)
// Super admin debe seleccionar negocio antes de crear
function WhatsAppActivateForm({
    integrationType,
    onSuccess,
    onBack,
}: {
    integrationType: IntegrationType;
    onSuccess: () => void;
    onBack: () => void;
}) {
    const { isSuperAdmin } = usePermissions();
    const { businesses, loading: loadingBusinesses } = useBusinessesSimple();
    const [selectedBusinessId, setSelectedBusinessId] = useState<number | null>(null);
    const [usePlatformToken, setUsePlatformToken] = useState(true);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const requiresBusinessSelection = isSuperAdmin && selectedBusinessId === null;

    const handleActivate = async () => {
        if (isSuperAdmin && !selectedBusinessId) return;

        setLoading(true);
        setError(null);

        try {
            const result = await createIntegrationAction({
                name: 'WhatsApp',
                code: 'whatsapp',
                integration_type_id: integrationType.id,
                category: integrationType.category?.code || integrationType.integration_category?.code || 'messaging',
                business_id: isSuperAdmin ? selectedBusinessId : null,
                is_active: true,
                is_default: true,
                config: { use_platform_token: usePlatformToken },
                credentials: {},
            });

            if (result.success) {
                onSuccess();
            } else {
                setError(result.message || 'Error al activar WhatsApp');
            }
        } catch (err: any) {
            setError(err.message || 'Error al activar WhatsApp');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="max-w-md mx-auto py-8 space-y-6">
            <div className="flex flex-col items-center text-center">
                <div className="w-16 h-16 bg-green-100 rounded-full flex items-center justify-center mb-4">
                    <ChatBubbleLeftRightIcon className="w-8 h-8 text-green-600" />
                </div>
                <h3 className="text-lg font-semibold text-gray-900">Activar WhatsApp</h3>
                <p className="text-sm text-gray-500 mt-2">
                    La integración de WhatsApp usa la configuración global del tipo de integración.
                    Las notificaciones se configuran desde el módulo de Notificaciones.
                </p>
            </div>

            {/* Selector de negocio - solo para super admin */}
            {isSuperAdmin && (
                <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <label className="block text-sm font-medium text-blue-800 mb-2">
                        Selecciona un negocio *
                    </label>
                    {loadingBusinesses ? (
                        <p className="text-sm text-blue-600">Cargando negocios...</p>
                    ) : (
                        <select
                            value={selectedBusinessId?.toString() ?? ''}
                            onChange={(e) => setSelectedBusinessId(e.target.value ? Number(e.target.value) : null)}
                            className="w-full px-3 py-2 border border-blue-300 rounded-lg text-sm bg-white focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                        >
                            <option value="">— Selecciona un negocio —</option>
                            {businesses.map(b => (
                                <option key={b.id} value={b.id}>{b.name} (ID: {b.id})</option>
                            ))}
                        </select>
                    )}
                </div>
            )}

            {/* Usar credenciales del tipo de integración */}
            <label className="flex items-start gap-3 cursor-pointer p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                <input
                    type="checkbox"
                    checked={usePlatformToken}
                    onChange={(e) => setUsePlatformToken(e.target.checked)}
                    className="mt-0.5 h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                />
                <div>
                    <span className="text-sm font-medium text-gray-900">Usar credenciales del tipo de integración</span>
                    <p className="text-xs text-gray-500 mt-1">
                        Usa las credenciales globales configuradas en el tipo de integración WhatsApp (access_token, phone_number_id).
                        Desactiva esta opción si este negocio tiene sus propias credenciales de Meta.
                    </p>
                </div>
            </label>

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="flex gap-3">
                <Button
                    type="button"
                    variant="outline"
                    onClick={onBack}
                    className="flex-1"
                >
                    Cancelar
                </Button>
                <Button
                    type="button"
                    variant="primary"
                    onClick={handleActivate}
                    disabled={loading || requiresBusinessSelection}
                    loading={loading}
                    className="flex-1"
                >
                    Activar WhatsApp
                </Button>
            </div>
        </div>
    );
}
