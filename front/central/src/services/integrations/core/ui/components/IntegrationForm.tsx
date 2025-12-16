'use client';

import { useState, useEffect } from 'react';
import { createIntegrationAction, updateIntegrationAction, getActiveIntegrationTypesAction, testIntegrationAction, testConnectionRawAction } from '../../infra/actions';
import { Integration, IntegrationType } from '../../domain/types';
import { Alert } from '@/shared/ui';
import ShopifyIntegrationForm from './shopify/ShopifyIntegrationForm';
import WhatsAppIntegrationView from './whatsapp/WhatsAppIntegrationView';
import DynamicIntegrationForm from './DynamicIntegrationForm';

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

        // Show edit form for WhatsApp and other dynamic types
        if (selectedType && (selectedType.code.toLowerCase() === 'whatsapp' || selectedType.code.toLowerCase() === 'whatsap')) {
            if (selectedType.config_schema && selectedType.credentials_schema) {
                return (
                    <DynamicIntegrationForm
                        integrationType={selectedType}
                        isEdit={true}
                        initialData={{
                            name: integration.name,
                            code: integration.code,
                            config: parsedConfig,
                            credentials: integration.credentials || {}, // Credenciales desencriptadas (si están disponibles)
                            business_id: integration.business_id,
                        }}
                        onSubmit={async (data) => {
                            try {
                                if (!integration.id) {
                                    throw new Error('ID de integración no encontrado');
                                }
                                // Solo enviar credenciales si hay valores (no vacío)
                                const updateData: any = {
                                    name: data.name,
                                    code: data.code,
                                    config: data.config,
                                };
                                // Solo incluir credenciales si hay valores ingresados
                                if (data.credentials && Object.keys(data.credentials).length > 0) {
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
                        }}
                        onTest={async (config, credentials) => {
                            try {
                                const result = await testConnectionRawAction(selectedType.code, config, credentials);
                                return {
                                    success: result.success,
                                    message: result.message
                                };
                            } catch (error: any) {
                                return {
                                    success: false,
                                    message: error.message
                                };
                            }
                        }}
                        onCancel={onCancel}
                    />
                );
            }
        }

        // For other types with schemas, show dynamic form
        if (selectedType && selectedType.config_schema && selectedType.credentials_schema) {
            return (
                <DynamicIntegrationForm
                    integrationType={selectedType}
                    isEdit={true}
                    initialData={{
                        name: integration.name,
                        code: integration.code,
                        config: parsedConfig,
                        credentials: integration.credentials || {}, // Credenciales desencriptadas (si están disponibles)
                        business_id: integration.business_id,
                    }}
                    onSubmit={async (data) => {
                        try {
                            if (!integration.id) {
                                throw new Error('ID de integración no encontrado');
                            }
                            // Solo enviar credenciales si hay valores (no vacío)
                            const updateData: any = {
                                name: data.name,
                                code: data.code,
                                config: data.config,
                            };
                            // Solo incluir credenciales si hay valores ingresados
                            if (data.credentials && Object.keys(data.credentials).length > 0) {
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
                    }}
                    onTest={async (config, credentials) => {
                        try {
                            const result = await testConnectionRawAction(selectedType.code, config, credentials);
                            return {
                                success: result.success,
                                message: result.message
                            };
                        } catch (error: any) {
                            return {
                                success: false,
                                message: error.message
                            };
                        }
                    }}
                    onCancel={onCancel}
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
                                                    if (type.code.toLowerCase() === 'shopify') {
                                                        target.src = '/integrations/shopify.png';
                                                    } else if (type.code.toLowerCase() === 'whatsapp' || type.code.toLowerCase() === 'whatsap') {
                                                        target.src = '/integrations/whatsapp.png';
                                                    } else {
                                                        target.style.display = 'none';
                                                    }
                                                }}
                                            />
                                        ) : (
                                            // Fallback a imágenes hardcodeadas si no hay imagen_url
                                            <>
                                                {type.code.toLowerCase() === 'shopify' && (
                                                    <img
                                                        src="/integrations/shopify.png"
                                                        alt="Shopify"
                                                        className="w-14 h-14 object-contain rounded-lg shadow-md"
                                                    />
                                                )}
                                                {(type.code.toLowerCase() === 'whatsapp' || type.code.toLowerCase() === 'whatsap') && (
                                                    <img
                                                        src="/integrations/whatsapp.png"
                                                        alt="WhatsApp"
                                                        className="w-14 h-14 object-contain rounded-lg shadow-md"
                                                    />
                                                )}
                                                {type.code.toLowerCase() !== 'shopify' && type.code.toLowerCase() !== 'whatsapp' && type.code.toLowerCase() !== 'whatsap' && (
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
                    {selectedType.code.toLowerCase() === 'shopify' && (
                        <ShopifyIntegrationForm
                            onSubmit={handleShopifySubmit}
                            onCancel={onCancel}
                            onTestConnection={handleTestConnection}
                        />
                    )}

                    {selectedType.code.toLowerCase() !== 'shopify' && selectedType.config_schema && selectedType.credentials_schema && (
                        <DynamicIntegrationForm
                            integrationType={selectedType}
                            onSubmit={async (data) => {
                                try {
                                    const result = await createIntegrationAction({
                                        name: data.name,
                                        code: data.code,
                                        integration_type_id: selectedType.id,
                                        category: selectedType.category,
                                        business_id: data.business_id || null,
                                        config: data.config,
                                        credentials: data.credentials,
                                        is_active: true,
                                    });

                                    if (result.success) {
                                        onSuccess?.();
                                    } else {
                                        setError(result.message || 'Error al crear la integración');
                                    }
                                } catch (err: any) {
                                    setError(err.message || 'Error al crear la integración');
                                }
                            }}
                            onTest={async (config, credentials) => {
                                try {
                                    const result = await testConnectionRawAction(selectedType.code, config, credentials);
                                    return {
                                        success: result.success,
                                        message: result.message
                                    };
                                } catch (error: any) {
                                    return {
                                        success: false,
                                        message: error.message
                                    };
                                }
                            }}
                            onCancel={onCancel}
                        />
                    )}

                    {selectedType.code.toLowerCase() !== 'shopify' && (!selectedType.config_schema || !selectedType.credentials_schema) && (
                        <Alert type="warning">
                            <p className="font-medium">Esquema no configurado</p>
                            <p className="text-sm mt-1">Este tipo de integración aún no tiene un esquema configurado. Por favor, configura los schemas en el módulo de administración.</p>
                        </Alert>
                    )}
                </div>
            )}
        </div>
    );
}
