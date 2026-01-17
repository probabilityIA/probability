'use client';

import { useEffect, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { Alert } from '@/shared/ui';

export default function ShopifyOAuthCallback() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [status, setStatus] = useState<'processing' | 'success' | 'error'>('processing');
    const [message, setMessage] = useState('Procesando autorización de Shopify...');

    useEffect(() => {
        const handleOAuthCallback = async () => {
            const oauthStatus = searchParams.get('shopify_oauth');

            if (oauthStatus === 'success') {
                // Extraer datos de la URL
                const shop = searchParams.get('shop');
                const integrationName = searchParams.get('integration_name');
                const integrationCode = searchParams.get('integration_code');
                const accessToken = searchParams.get('access_token');
                const businessId = searchParams.get('business_id');

                if (!shop || !integrationName || !integrationCode || !accessToken) {
                    setStatus('error');
                    setMessage('Datos de OAuth incompletos');
                    return;
                }

                try {
                    setMessage('Creando integración...');

                    // Crear la integración
                    const response = await createIntegrationAction({
                        name: integrationName,
                        code: integrationCode,
                        integration_type_id: 1, // Shopify
                        category: 'ecommerce',
                        store_id: shop,
                        is_active: true,
                        is_default: false,
                        config: {
                            store_name: shop,
                            api_version: '2024-10',
                            webhook_configured: false,
                        },
                        credentials: {
                            access_token: accessToken,
                        },
                        business_id: businessId ? parseInt(businessId) : null,
                    });

                    if (response.success) {
                        setStatus('success');
                        setMessage('¡Integración creada exitosamente!');

                        // Redirigir a la lista de integraciones después de 2 segundos
                        setTimeout(() => {
                            router.push('/integrations');
                        }, 2000);
                    } else {
                        throw new Error(response.message || 'Error al crear integración');
                    }
                } catch (error: any) {
                    console.error('Error creating integration:', error);
                    setStatus('error');
                    setMessage(error.message || 'Error al crear la integración');
                }
            } else if (oauthStatus === 'error') {
                setStatus('error');
                setMessage(searchParams.get('message') || 'Error en el proceso de autorización');
            }
        };

        handleOAuthCallback();
    }, [searchParams, router]);

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-50">
            <div className="max-w-md w-full p-6">
                {status === 'processing' && (
                    <div className="text-center">
                        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                        <p className="text-gray-700">{message}</p>
                    </div>
                )}

                {status === 'success' && (
                    <Alert type="success">
                        <div className="text-center">
                            <p className="font-medium">{message}</p>
                            <p className="text-sm mt-2">Redirigiendo...</p>
                        </div>
                    </Alert>
                )}

                {status === 'error' && (
                    <Alert type="error">
                        <div className="text-center">
                            <p className="font-medium">{message}</p>
                            <button
                                onClick={() => router.push('/integrations')}
                                className="mt-4 px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700"
                            >
                                Volver a Integraciones
                            </button>
                        </div>
                    </Alert>
                )}
            </div>
        </div>
    );
}
