'use client';

import { useState } from 'react';
import { Modal, Alert } from '@/shared/ui';
import { IntegrationCategory, IntegrationType } from '../../domain/types';
import { CategorySelector } from './CategorySelector';
import { ProviderSelector } from './ProviderSelector';
import { createIntegrationAction, testConnectionRawAction } from '../../infra/actions';

// Importar formularios específicos por tipo de integración
import { SoftpymesConfigForm } from '@/services/integrations/invoicing/softpymes/ui/components';

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

            // TODO: Agregar más formularios específicos aquí
            // case 'shopify':
            //     return <ShopifyConfigForm onSuccess={onSuccess} onCancel={onBack} />;
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
