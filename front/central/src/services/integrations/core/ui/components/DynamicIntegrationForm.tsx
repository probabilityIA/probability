'use client';

import { useState, useEffect } from 'react';
import DynamicField, { FieldSchema } from './DynamicField';
import { Button, Alert } from '@/shared/ui';
import { IntegrationType } from '../../domain/types';

interface DynamicIntegrationFormProps {
    integrationType: IntegrationType;
    onSubmit: (data: {
        name: string;
        code: string;
        config: Record<string, any>;
        credentials: Record<string, any>;
        business_id?: number | null;
    }) => Promise<void>;
    onCancel?: () => void;
    initialData?: {
        name?: string;
        code?: string;
        config?: Record<string, any>;
        credentials?: Record<string, any>;
        business_id?: number | null;
    };
    isEdit?: boolean;
    onTest?: (config: Record<string, any>, credentials: Record<string, any>) => Promise<{ success: boolean; message?: string }>;
}

export default function DynamicIntegrationForm({
    integrationType,
    onSubmit,
    onCancel,
    onTest,
    initialData,
    isEdit = false
}: DynamicIntegrationFormProps) {
    const [formData, setFormData] = useState({
        name: initialData?.name || '',
        code: initialData?.code || '',
        business_id: initialData?.business_id || null,
    });

    const [configData, setConfigData] = useState<Record<string, any>>(initialData?.config || {});
    const [credentialsData, setCredentialsData] = useState<Record<string, any>>(initialData?.credentials || {});
    
    // Inicializar testPhone desde config si existe (para edición)
    const [testPhone, setTestPhone] = useState<string>(initialData?.config?.test_phone_number || '');

    // Actualizar estados cuando initialData cambie (importante para edición)
    useEffect(() => {
        if (initialData) {
            console.log('🔄 Actualizando DynamicIntegrationForm con initialData:', initialData);
            
            setFormData({
                name: initialData.name || '',
                code: initialData.code || '',
                business_id: initialData.business_id || null,
            });
            
            // Parse config si es string
            let parsedConfig = initialData.config || {};
            if (typeof initialData.config === 'string') {
                try {
                    parsedConfig = JSON.parse(initialData.config);
                    console.log('✅ Config parseado desde string:', parsedConfig);
                } catch (e) {
                    console.error('❌ Error parsing config:', e);
                    parsedConfig = {};
                }
            } else if (initialData.config) {
                // Si ya es un objeto, usarlo directamente
                parsedConfig = initialData.config;
                console.log('✅ Config ya es objeto:', parsedConfig);
            }
            
            console.log('📦 ConfigData que se va a establecer:', parsedConfig);
            setConfigData(parsedConfig);
            
            // Parse credentials si es string
            let parsedCredentials = initialData.credentials || {};
            if (typeof initialData.credentials === 'string') {
                try {
                    parsedCredentials = JSON.parse(initialData.credentials);
                    console.log('✅ Credentials parseado desde string:', parsedCredentials);
                } catch (e) {
                    console.error('❌ Error parsing credentials:', e);
                    parsedCredentials = {};
                }
            } else if (initialData.credentials) {
                parsedCredentials = initialData.credentials;
                console.log('✅ Credentials ya es objeto:', parsedCredentials);
                console.log('🔍 Keys en credentials:', Object.keys(parsedCredentials));
                console.log('🔍 Valores en credentials:', Object.entries(parsedCredentials).map(([k, v]) => ({ key: k, hasValue: !!v, valueLength: typeof v === 'string' ? v.length : 'N/A' })));
            }
            
            console.log('🔐 CredentialsData que se va a establecer:', parsedCredentials);
            setCredentialsData(parsedCredentials);
            setTestPhone(parsedConfig?.test_phone_number || '');
        }
    }, [initialData]);

    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [errors, setErrors] = useState<Record<string, string>>({});

    // Parse schemas
    const configSchema = integrationType.config_schema as any;
    const credentialsSchema = integrationType.credentials_schema as any;

    const configFields = configSchema?.properties ? Object.entries(configSchema.properties) as [string, FieldSchema][] : [];
    const credentialsFields = credentialsSchema?.properties ? Object.entries(credentialsSchema.properties) as [string, FieldSchema][] : [];

    // Filtrar test_phone_number de los campos dinámicos de configuración
    // porque se maneja manualmente más abajo
    const filteredConfigFields = configFields.filter(([name]) => name !== 'test_phone_number');

    // Sort fields by order
    const sortedConfigFields = filteredConfigFields.sort((a, b) => (a[1].order || 999) - (b[1].order || 999));
    const sortedCredentialsFields = credentialsFields.sort((a, b) => (a[1].order || 999) - (b[1].order || 999));

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        setErrors({});

        // Validate required fields
        const newErrors: Record<string, string> = {};

        sortedConfigFields.forEach(([name, schema]) => {
            if (schema.required && !configData[name]) {
                newErrors[`config_${name}`] = schema.error_message || 'Este campo es requerido';
            }
        });

        sortedCredentialsFields.forEach(([name, schema]) => {
            if (schema.required && !credentialsData[name]) {
                newErrors[`credentials_${name}`] = schema.error_message || 'Este campo es requerido';
            }
        });

        if (Object.keys(newErrors).length > 0) {
            setErrors(newErrors);
            setLoading(false);
            return;
        }

        try {
            // Asegurar que test_phone_number esté en el config si se ingresó
            const finalConfig = {
                ...configData,
                ...(testPhone && { test_phone_number: testPhone })
            };
            
            await onSubmit({
                name: formData.name,
                code: formData.code,
                config: finalConfig,
                credentials: credentialsData,
                business_id: formData.business_id,
            });
        } catch (err: any) {
            console.error('Error saving integration:', err);
            setError(err.message || 'Error al guardar la integración');
        } finally {
            setLoading(false);
        }
    };

    const handleTestConnection = async () => {
        if (!testPhone) {
            setError('Por favor ingresa un número de teléfono de prueba');
            return;
        }

        if (!onTest) {
            setError('La función de prueba no está disponible');
            return;
        }

        setLoading(true);
        setError(null);

        try {
            // Include test phone in config for the backend
            const configWithTestPhone = {
                ...configData,
                test_phone_number: testPhone
            };

            const result = await onTest(configWithTestPhone, credentialsData);

            if (result.success) {
                alert(result.message || `Mensaje de prueba enviado a ${testPhone}`);
            } else {
                setError(result.message || 'Error al enviar mensaje de prueba');
            }
        } catch (err: any) {
            setError(err.message || 'Error al enviar mensaje de prueba');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {/* Configuration Section - All fields together */}
            <div>
                <div className="flex items-center gap-3 mb-4 pb-2 border-b-2 border-gray-900">
                    {/* WhatsApp Logo */}
                    {(integrationType.code.toLowerCase() === 'whatsapp') && (
                        <div className="flex-shrink-0">
                            <div className="w-12 h-12 bg-green-500 rounded-full flex items-center justify-center">
                                <svg className="w-7 h-7 text-white" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M17.472 14.382c-.297-.149-1.758-.867-2.03-.967-.273-.099-.471-.148-.67.15-.197.297-.767.966-.94 1.164-.173.199-.347.223-.644.075-.297-.15-1.255-.463-2.39-1.475-.883-.788-1.48-1.761-1.653-2.059-.173-.297-.018-.458.13-.606.134-.133.298-.347.446-.52.149-.174.198-.298.298-.497.099-.198.05-.371-.025-.52-.075-.149-.669-1.612-.916-2.207-.242-.579-.487-.5-.669-.51-.173-.008-.371-.01-.57-.01-.198 0-.52.074-.792.372-.272.297-1.04 1.016-1.04 2.479 0 1.462 1.065 2.875 1.213 3.074.149.198 2.096 3.2 5.077 4.487.709.306 1.262.489 1.694.625.712.227 1.36.195 1.871.118.571-.085 1.758-.719 2.006-1.413.248-.694.248-1.289.173-1.413-.074-.124-.272-.198-.57-.347m-5.421 7.403h-.004a9.87 9.87 0 01-5.031-1.378l-.361-.214-3.741.982.998-3.648-.235-.374a9.86 9.86 0 01-1.51-5.26c.001-5.45 4.436-9.884 9.888-9.884 2.64 0 5.122 1.03 6.988 2.898a9.825 9.825 0 012.893 6.994c-.003 5.45-4.437 9.884-9.885 9.884m8.413-18.297A11.815 11.815 0 0012.05 0C5.495 0 .16 5.335.157 11.892c0 2.096.547 4.142 1.588 5.945L.057 24l6.305-1.654a11.882 11.882 0 005.683 1.448h.005c6.554 0 11.89-5.335 11.893-11.893a11.821 11.821 0 00-3.48-8.413Z" />
                                </svg>
                            </div>
                        </div>
                    )}
                    <h3 className="text-sm font-semibold text-gray-900">
                        Configuración
                    </h3>
                </div>

                {/* Name and Config Fields in 2 columns */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3 mb-3">
                    {/* Name Input */}
                    <div>
                        <label className="block text-xs font-bold text-gray-700 mb-1.5 flex items-center gap-1">
                            Nombre de la Integración <span className="text-red-500">*</span>
                            <span className="group relative">
                                <span className="text-gray-400 cursor-help">ⓘ</span>
                                <span className="invisible group-hover:visible absolute left-0 top-6 w-48 bg-gray-900 text-white text-xs rounded py-1 px-2 z-10">
                                    Nombre descriptivo para identificar esta integración
                                </span>
                            </span>
                        </label>
                        <input
                            type="text"
                            required
                            placeholder="Ej: WhatsApp Principal"
                            value={formData.name}
                            onChange={(e) => {
                                const name = e.target.value;
                                setFormData({
                                    ...formData,
                                    name,
                                    code: name.toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_]/g, '')
                                });
                            }}
                            className="w-full px-3 py-2 text-sm text-gray-900 bg-gray-50 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-white"
                        />
                    </div>

                    {/* Dynamic Config Fields */}
                    {sortedConfigFields.map(([name, schema]) => {
                        const fieldValue = configData[name];
                        console.log(`🔍 Campo config "${name}":`, fieldValue, 'de configData:', configData);
                        return (
                            <DynamicField
                                key={name}
                                name={name}
                                schema={schema}
                                value={fieldValue}
                                onChange={(value) => setConfigData({ ...configData, [name]: value })}
                                error={errors[`config_${name}`]}
                            />
                        );
                    })}
                </div>

                {/* Credentials Fields */}
                {sortedCredentialsFields.length > 0 && (
                    <div className="grid grid-cols-1 md:grid-cols-3 gap-3 mt-3">
                        {sortedCredentialsFields.map(([name, schema]) => {
                            const fieldValue = credentialsData[name];
                            console.log(`🔍 Campo credential "${name}":`, {
                                value: fieldValue,
                                type: typeof fieldValue,
                                isUndefined: fieldValue === undefined,
                                isNull: fieldValue === null,
                                isEmpty: fieldValue === '',
                                credentialsDataKeys: Object.keys(credentialsData),
                                allCredentialsData: credentialsData
                            });
                            return (
                                <DynamicField
                                    key={name}
                                    name={name}
                                    schema={{ ...schema, boldLabel: true }}
                                    value={fieldValue ?? ''}
                                    onChange={(value) => setCredentialsData({ ...credentialsData, [name]: value })}
                                    error={errors[`credentials_${name}`]}
                                />
                            );
                        })}

                        {/* Test Phone Number */}
                        <div>
                            <label className="block text-xs font-bold text-gray-700 mb-1.5 flex items-center gap-1">
                                Número de Prueba
                                <span className="group relative">
                                    <span className="text-gray-400 cursor-help">ⓘ</span>
                                    <span className="invisible group-hover:visible absolute left-0 top-6 w-48 bg-gray-900 text-white text-xs rounded py-1 px-2 z-10">
                                        Número para mensaje de prueba
                                    </span>
                                </span>
                            </label>
                            <input
                                type="tel"
                                placeholder="+57 300 123 4567"
                                value={testPhone}
                                onChange={(e) => {
                                    const value = e.target.value;
                                    setTestPhone(value);
                                    // Guardar también en configData para que se persista
                                    setConfigData({ ...configData, test_phone_number: value });
                                }}
                                className="w-full px-3 py-2 text-sm text-gray-900 bg-gray-50 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent focus:bg-white"
                            />
                        </div>

                        {/* Test Button */}
                        <div>
                            <label className="block text-xs font-bold text-transparent mb-1.5 select-none">
                                Action
                            </label>
                            <Button
                                type="button"
                                onClick={handleTestConnection}
                                disabled={!testPhone || loading}
                                variant="outline"
                                className="w-full h-[38px] flex items-center justify-center gap-1 text-xs"
                            >
                                <span className="inline-block">🧪</span>
                                <span>Probar</span>
                            </Button>
                        </div>
                    </div>
                )}
            </div>

            {/* Setup Instructions - Blue Info Box at Bottom */}
            {integrationType.setup_instructions && (
                <div className="bg-blue-100 border border-blue-300 rounded-lg p-4">
                    <h3 className="text-sm font-bold text-blue-900 mb-2">
                        📋 Instrucciones de Configuración
                    </h3>
                    <pre className="whitespace-pre-wrap text-xs text-blue-800 leading-relaxed font-sans">
                        {integrationType.setup_instructions}
                    </pre>
                </div>
            )}

            {/* Action Buttons - Clean */}
            <div className="flex justify-end space-x-3 pt-4 border-t border-gray-200">
                {onCancel && (
                    <Button
                        type="button"
                        onClick={onCancel}
                        variant="outline"
                    >
                        Cancelar
                    </Button>
                )}
                <Button
                    type="submit"
                    disabled={loading}
                    loading={loading}
                    variant="primary"
                >
                    {isEdit ? 'Actualizar Integración' : 'Crear Integración'}
                </Button>
            </div>
        </form>
    );
}
