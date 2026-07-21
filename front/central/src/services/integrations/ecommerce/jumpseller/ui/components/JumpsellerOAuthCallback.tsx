'use client';

import { useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { createIntegrationAction } from '@/services/integrations/core/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import { CheckCircleIcon, ExclamationTriangleIcon } from '@heroicons/react/24/outline';

type Status = 'processing' | 'success' | 'error';

export function JumpsellerOAuthCallback() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const hasRun = useRef(false);
    const [status, setStatus] = useState<Status>('processing');
    const [message, setMessage] = useState('Conectando con Jumpseller...');

    useEffect(() => {
        if (hasRun.current) return;
        hasRun.current = true;

        const run = async () => {
            const oauthStatus = searchParams.get('jumpseller_oauth');

            if (oauthStatus === 'error') {
                setStatus('error');
                setMessage(searchParams.get('message') || 'Error al conectar con Jumpseller');
                return;
            }

            const integrationName = searchParams.get('integration_name') || 'Jumpseller';
            const integrationCode = `jumpseller_${Date.now()}`;
            const state = searchParams.get('state');
            const businessId = searchParams.get('business_id');
            const exchangeToken = searchParams.get('exchange_token');
            const isTesting = searchParams.get('is_testing') === 'true';

            if (!state || !exchangeToken) {
                setStatus('error');
                setMessage('Faltan parametros del flujo OAuth. Intenta conectar nuevamente.');
                return;
            }

            try {
                const sessionToken = TokenStorage.getSessionToken();
                const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
                const tokenResponse = await fetch(
                    `${apiBaseUrl}/integrations/jumpseller/oauth/token?state=${encodeURIComponent(state)}&exchange_token=${encodeURIComponent(exchangeToken)}`,
                    {
                        headers: { 'Authorization': `Bearer ${sessionToken}` },
                        credentials: 'include',
                    }
                );

                if (tokenResponse.status === 410) {
                    throw new Error('El token de autorizacion expiro. Inicia la conexion con Jumpseller nuevamente.');
                }

                const tokenData = await tokenResponse.json();
                if (!tokenResponse.ok || !tokenData.success) {
                    throw new Error(tokenData.error || 'No se pudieron recuperar las credenciales de Jumpseller');
                }

                const response = await createIntegrationAction({
                    name: integrationName,
                    code: integrationCode,
                    integration_type_id: 33,
                    category: 'ecommerce',
                    store_id: '',
                    is_active: true,
                    is_default: false,
                    config: {
                        auth_method: 'oauth',
                        token_expires_at: tokenData.expires_at,
                    } as any,
                    credentials: {
                        access_token: tokenData.access_token,
                        refresh_token: tokenData.refresh_token,
                    } as any,
                    business_id: businessId && Number(businessId) > 0 ? Number(businessId) : null,
                    is_testing: isTesting || tokenData.is_testing || false,
                }, sessionToken || undefined);

                if (!response || response.success === false) {
                    throw new Error(response?.message || 'Error al crear la integracion');
                }

                setStatus('success');
                setMessage('Jumpseller conectado exitosamente. Redirigiendo...');
                setTimeout(() => router.push('/integrations'), 2000);
            } catch (err: any) {
                setStatus('error');
                setMessage(err.message || 'Error al completar la conexion con Jumpseller');
            }
        };

        run();
    }, [searchParams, router]);

    return (
        <div className="flex min-h-[60vh] items-center justify-center p-6">
            <div className="w-full max-w-md rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-8 text-center shadow-sm">
                {status === 'processing' && (
                    <>
                        <svg className="mx-auto h-10 w-10 animate-spin text-[var(--color-primary)]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <h2 className="mt-4 text-lg font-bold text-gray-900 dark:text-white">Conectando</h2>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{message}</p>
                    </>
                )}
                {status === 'success' && (
                    <>
                        <CheckCircleIcon className="mx-auto h-12 w-12 text-[var(--color-primary)]" />
                        <h2 className="mt-4 text-lg font-bold text-gray-900 dark:text-white">Conectado</h2>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{message}</p>
                    </>
                )}
                {status === 'error' && (
                    <>
                        <ExclamationTriangleIcon className="mx-auto h-12 w-12 text-red-500" />
                        <h2 className="mt-4 text-lg font-bold text-gray-900 dark:text-white">No se pudo conectar</h2>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{message}</p>
                        <button
                            type="button"
                            onClick={() => router.push('/integrations')}
                            className="mt-5 rounded-lg bg-[var(--color-primary)] px-5 py-2 text-sm font-semibold text-white"
                        >
                            Volver a integraciones
                        </button>
                    </>
                )}
            </div>
        </div>
    );
}
