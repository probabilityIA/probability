'use client';

import { useState } from 'react';
import { Modal } from '@/shared/ui';
import { IntegrationCategory, IntegrationType } from '../../domain/types';
import { CategorySelector } from './CategorySelector';
import { ProviderSelector } from './ProviderSelector';
import DynamicIntegrationForm from './DynamicIntegrationForm';
import { createIntegrationAction, testConnectionRawAction } from '../../infra/actions';

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
    const handleSubmit = async (data: {
        name: string;
        code: string;
        config: Record<string, any>;
        credentials: Record<string, any>;
        business_id?: number | null;
    }) => {
        const integrationData = {
            name: data.name,
            code: data.code,
            integration_type_id: integrationType.id,
            category: integrationType.category,
            business_id: data.business_id || null,
            config: data.config,
            credentials: data.credentials,
            is_active: true,
            is_default: false,
        };

        await createIntegrationAction(integrationData);
        onSuccess();
    };

    const handleTest = async (config: Record<string, any>, credentials: Record<string, any>) => {
        try {
            const result = await testConnectionRawAction(integrationType.code, config, credentials);
            return result;
        } catch (error: any) {
            return { success: false, message: error.message || 'Error al probar conexión' };
        }
    };

    return (
        <div className="p-6">
            <button
                onClick={onBack}
                className="text-blue-600 hover:text-blue-800 mb-4 flex items-center gap-2"
            >
                <span>←</span>
                <span>Volver a proveedores</span>
            </button>

            <DynamicIntegrationForm
                integrationType={integrationType}
                onSubmit={handleSubmit}
                onCancel={onCancel}
                onTest={handleTest}
            />
        </div>
    );
}
