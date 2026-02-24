'use client';

import { useState } from 'react';
import { Modal, Alert } from '@/shared/ui';
import { IntegrationCategory, IntegrationType } from '../../domain/types';
import { CategorySelector } from './CategorySelector';
import { ProviderSelector } from './ProviderSelector';
import { createIntegrationAction, testConnectionRawAction } from '../../infra/actions';

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

    const integrationCode = integrationType.code.toLowerCase();

    // Renderizar formulario específico según el código del tipo de integración
    const renderSpecificForm = () => {
        switch (integrationCode) {
            case 'softpymes':
                return (
                    <SoftpymesConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'shopify':
                return (
                    <ShopifyIntegrationForm
                        onSubmit={async () => onSuccess()}
                        onCancel={onBack}
                    />
                );

            case 'factus':
                return (
                    <FactusConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'envioclick':
                return (
                    <EnvioClickConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                        integrationTypeBaseURL={integrationType.base_url}
                        integrationTypeBaseURLTest={integrationType.base_url_test}
                    />
                );

            case 'enviame':
                return (
                    <EnviameConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'tu':
                return (
                    <TuConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'mipaquete':
            case 'mi_paquete':
                return (
                    <MiPaqueteConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'siigo':
                return (
                    <SiigoConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'alegra':
                return (
                    <AlegraConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'world_office':
                return (
                    <WorldOfficeConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'helisa':
                return (
                    <HelisaConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'vtex':
                return (
                    <VTEXConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'tiendanube':
            case 'tienda_nube':
                return (
                    <TiendanubeConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'magento':
            case 'adobe_commerce':
                return (
                    <MagentoConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'amazon':
                return (
                    <AmazonConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'falabella':
                return (
                    <FalabellaConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'exito':
                return (
                    <ExitoConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            case 'woocommerce':
                return (
                    <WooCommerceConfigForm
                        onSuccess={onSuccess}
                        onCancel={onBack}
                    />
                );

            // case 'whatsapp':
            //     return <WhatsAppConfigForm onSuccess={onSuccess} onCancel={onBack} />;
            // case 'mercadolibre':
            //     return <MercadoLibreConfigForm onSuccess={onSuccess} onCancel={onBack} />;

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
