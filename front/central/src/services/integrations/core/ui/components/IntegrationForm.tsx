'use client';

import { useState, useEffect } from 'react';
import { createIntegrationAction, updateIntegrationAction, getActiveIntegrationTypesAction, testIntegrationAction } from '../../infra/actions';
import { Integration, IntegrationType } from '../../domain/types';
import { Alert } from '@/shared/ui';
import ShopifyIntegrationForm from './ShopifyIntegrationForm';
import WhatsAppIntegrationView from './WhatsAppIntegrationView';

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
        config: any;
        credentials: any;
        business_id?: number | null;
    }) => {
        if (!selectedType) return;

        const integrationData = {
            name: data.name,
            code: data.code,
            integration_type_id: selectedType.id,
            category: selectedType.category,
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
        // For now, we'll need to create the integration first, then test
        // In a real scenario, you might want a separate test endpoint
        try {
            // This would require a backend endpoint that tests without saving
            // For now, return true as placeholder
            console.log('Testing connection with:', { config, credentials });
            return true;
        } catch (error) {
            console.error('Test connection error:', error);
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
        // WhatsApp is read-only
        if (selectedType?.code.toLowerCase() === 'whatsapp' || selectedType?.code.toLowerCase() === 'whatsap') {
            return (
                <WhatsAppIntegrationView
                    integration={integration}
                    onTestConnection={handleWhatsAppTest}
                    onRefresh={() => window.location.reload()}
                />
            );
        }

        // For other types, show a generic message for now
        return (
            <div className="text-center py-8">
                <p className="text-gray-600">La edición de integraciones de tipo {selectedType?.name} aún no está implementada.</p>
            </div>
        );
    }

    // Creating new integration - show type selector first if no type selected
    return (
        <div className="space-y-6">
            {/* Type Selector - Show when no type is selected */}
            {!selectedType && integrationTypes.length > 0 && (
                <div className="bg-gray-50 p-4 rounded-lg">
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                        Selecciona el tipo de integración *
                    </label>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                        {integrationTypes.map((type) => (
                            <button
                                key={type.id}
                                type="button"
                                onClick={() => handleTypeChange(type.id)}
                                className="p-4 border-2 rounded-lg text-left transition-all hover:border-blue-300 hover:shadow-md border-gray-200"
                            >
                                <div className="flex items-start gap-3">
                                    {/* Logo */}
                                    {type.code.toLowerCase() === 'shopify' && (
                                        <img
                                            src="/integrations/shopify.png"
                                            alt="Shopify"
                                            className="w-10 h-10 object-contain"
                                        />
                                    )}
                                    {(type.code.toLowerCase() === 'whatsapp' || type.code.toLowerCase() === 'whatsap') && (
                                        <img
                                            src="/integrations/whatsapp.png"
                                            alt="WhatsApp"
                                            className="w-10 h-10 object-contain"
                                        />
                                    )}

                                    {/* Content */}
                                    <div className="flex-1">
                                        <div className="flex items-center justify-between">
                                            <h4 className="font-semibold text-gray-900">{type.name}</h4>
                                            <span className={`px-2 py-1 text-xs rounded-full ${type.category === 'external'
                                                ? 'bg-blue-100 text-blue-700'
                                                : 'bg-purple-100 text-purple-700'
                                                }`}>
                                                {type.category === 'external' ? 'Externa' : 'Interna'}
                                            </span>
                                        </div>
                                        <p className="text-sm text-gray-500 mt-1">{type.code}</p>
                                    </div>
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
                    {selectedType.code.toLowerCase() === 'shopify' && (
                        <ShopifyIntegrationForm
                            onSubmit={handleShopifySubmit}
                            onCancel={onCancel}
                            onTestConnection={handleTestConnection}
                        />
                    )}

                    {(selectedType.code.toLowerCase() === 'whatsapp' || selectedType.code.toLowerCase() === 'whatsap') && (
                        <Alert type="info">
                            <p className="font-medium">WhatsApp es una integración interna</p>
                            <p className="text-sm mt-1">No se pueden crear nuevas integraciones de WhatsApp. Solo existe una integración global para toda la plataforma.</p>
                        </Alert>
                    )}

                    {selectedType.code.toLowerCase() !== 'shopify' && selectedType.code.toLowerCase() !== 'whatsapp' && selectedType.code.toLowerCase() !== 'whatsap' && (
                        <Alert type="warning">
                            <p className="font-medium">Tipo de integración no soportado</p>
                            <p className="text-sm mt-1">El formulario para {selectedType.name} aún no está implementado.</p>
                        </Alert>
                    )}
                </div>
            )}
        </div>
    );
}
