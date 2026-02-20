'use client';

import { useEffect, useState, useRef } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { Alert } from '@/shared/ui';
import { TokenStorage } from '@/shared/utils';

export default function ShopifyOAuthCallback() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const [status, setStatus] = useState<'processing' | 'success' | 'error'>('processing');
    const [message, setMessage] = useState('Procesando autorización de Shopify...');
    const hasRun = useRef(false);

    useEffect(() => {
        if (hasRun.current) return;

        const handleOAuthCallback = async () => {
            hasRun.current = true;
            const oauthStatus = searchParams.get('shopify_oauth');

            if (oauthStatus === 'success') {
                // Extraer datos de la URL
                const shop = searchParams.get('shop');
                const integrationName = searchParams.get('integration_name');
                const integrationCode = searchParams.get('integration_code');
                const state = searchParams.get('state');
                const businessId = searchParams.get('business_id');

                if (!shop || !integrationName || !integrationCode || !state) {
                    setStatus('error');
                    setMessage('Datos de OAuth incompletos');
                    return;
                }

                try {
                    setMessage('Obteniendo credenciales de forma segura...');

                    const exchangeToken = searchParams.get('exchange_token');

                    // Obtener token desde storage unificado (soporta iframes/cookies/localstorage)
                    const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
                    const sessionToken = TokenStorage.getSessionToken();

                    const tokenResponse = await fetch(
                        `${apiBaseUrl}/integrations/shopify/oauth/token?state=${state}&shop=${shop}&integration_name=${integrationName}&integration_code=${integrationCode}&exchange_token=${exchangeToken || ''}`,
                        {
                            headers: {
                                'Authorization': `Bearer ${sessionToken}`,
                            },
                            credentials: 'include', // Incluir cookies por si acaso
                        }
                    );

                    if (!tokenResponse.ok) {
                        throw new Error('Error al obtener credenciales de Shopify');
                    }

                    const tokenData = await tokenResponse.json();
                    const accessToken = tokenData.access_token;
                    const clientId = tokenData.client_id;
                    const clientSecret = tokenData.client_secret;

                    if (!accessToken) {
                        throw new Error('Token de acceso no recibido');
                    }

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
                            client_id: clientId, // Guardar ID para referencia
                        },
                        credentials: {
                            access_token: accessToken,
                            client_id: clientId,
                            client_secret: clientSecret, // GUARDAR EL SECRETO PARA VALIDAR WEBHOOKS
                        },
                        business_id: businessId ? parseInt(businessId) : null,
                    }, sessionToken || undefined);

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
