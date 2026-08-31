'use client';

import { useCallback, useEffect, useRef, useState } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { createIntegrationAction, updateIntegrationAction } from '@/services/integrations/core/infra/actions';
import { getBusinessesSimpleAction } from '@/services/auth/business/infra/actions';
import { TokenStorage } from '@/shared/utils/token-storage';
import { CheckCircleIcon, ExclamationTriangleIcon, ShieldExclamationIcon } from '@heroicons/react/24/outline';

type Status = 'processing' | 'confirm' | 'saving' | 'success' | 'error';

interface PendingConnection {
    tokenData: any;
    sellerId: string;
    nickname: string;
    accountName: string;
    accountEmail: string;
    businessId: number | null;
    businessName: string;
    integrationName: string;
    integrationCode: string;
    isTesting: boolean;
    reconnectId: number | null;
}

function friendlyError(raw: string): string {
    const text = String(raw || '');
    if (/SQLSTATE|duplicate key|violates unique constraint|ux_integrations_store_type/i.test(text)) {
        return 'Esta cuenta de MercadoLibre ya está conectada. Verifica que estas autorizando la cuenta correcta.';
    }
    return text || 'Error al completar la conexión con MercadoLibre';
}

export function MercadoLibreOAuthCallback() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const hasRun = useRef(false);
    const [status, setStatus] = useState<Status>('processing');
    const [message, setMessage] = useState('Conectando con MercadoLibre...');
    const [pending, setPending] = useState<PendingConnection | null>(null);

    useEffect(() => {
        if (hasRun.current) return;
        hasRun.current = true;

        const run = async () => {
            const oauthStatus = searchParams.get('meli_oauth');

            if (oauthStatus === 'error') {
                setStatus('error');
                setMessage(searchParams.get('message') || 'Error al conectar con MercadoLibre');
                return;
            }

            const integrationName = searchParams.get('integration_name') || 'MercadoLibre';
            const integrationCode = searchParams.get('integration_code') || `mercado_libre_${Date.now()}`;
            const state = searchParams.get('state');
            const businessIdParam = searchParams.get('business_id');
            const exchangeToken = searchParams.get('exchange_token');
            const isTesting = searchParams.get('is_testing') === 'true';

            if (!state || !exchangeToken) {
                setStatus('error');
                setMessage('Faltan parámetros del flujo OAuth. Intenta conectar nuevamente.');
                return;
            }

            try {
                const sessionToken = TokenStorage.getSessionToken();
                const apiBaseUrl = process.env.NEXT_PUBLIC_API_BASE_URL || '/api/v1';
                const tokenResponse = await fetch(
                    `${apiBaseUrl}/integrations/meli/oauth/token?state=${encodeURIComponent(state)}&integration_name=${encodeURIComponent(integrationName)}&integration_code=${encodeURIComponent(integrationCode)}&exchange_token=${encodeURIComponent(exchangeToken)}`,
                    {
                        headers: { 'Authorization': `Bearer ${sessionToken}` },
                        credentials: 'include',
                    }
                );

                if (tokenResponse.status === 410) {
                    throw new Error('El token de autorización expiro. Inicia la conexión con MercadoLibre nuevamente.');
                }

                const tokenData = await tokenResponse.json();
                if (!tokenResponse.ok || !tokenData.success) {
                    throw new Error(tokenData.error || 'No se pudieron recuperar las credenciales de MercadoLibre');
                }

                const sellerId = tokenData.seller_id ? String(tokenData.seller_id) : '';
                const businessId = businessIdParam && Number(businessIdParam) > 0 ? Number(businessIdParam) : null;

                let reconnectId: number | null = null;
                try {
                    const raw = localStorage.getItem('meli_reconnect');
                    if (raw) {
                        const parsed = JSON.parse(raw);
                        if (parsed?.integration_id) reconnectId = Number(parsed.integration_id);
                    }
                } catch { }

                let businessName = '';
                if (businessId) {
                    try {
                        const r: any = await getBusinessesSimpleAction();
                        const match = (r?.data || []).find((b: any) => Number(b.id) === businessId);
                        if (match?.name) businessName = String(match.name);
                    } catch { }
                }

                setPending({
                    tokenData,
                    sellerId,
                    nickname: tokenData.nickname ? String(tokenData.nickname) : '',
                    accountName: tokenData.account_name ? String(tokenData.account_name) : '',
                    accountEmail: tokenData.account_email ? String(tokenData.account_email) : '',
                    businessId,
                    businessName,
                    integrationName,
                    integrationCode,
                    isTesting: isTesting || tokenData.is_testing || false,
                    reconnectId,
                });
                setStatus('confirm');
            } catch (err: any) {
                setStatus('error');
                setMessage(friendlyError(err?.message));
            }
        };

        run();
    }, [searchParams]);

    const confirmConnection = useCallback(async () => {
        if (!pending) return;
        setStatus('saving');
        setMessage('Guardando la conexión...');

        const { tokenData, sellerId, businessId, integrationName, integrationCode, isTesting, reconnectId } = pending;
        const sessionToken = TokenStorage.getSessionToken();

        const config = {
            app_id: tokenData.client_id,
            seller_id: tokenData.seller_id,
            token_expires_at: tokenData.expires_at,
        } as any;
        const credentials = {
            access_token: tokenData.access_token,
            refresh_token: tokenData.refresh_token,
            client_id: tokenData.client_id,
            client_secret: tokenData.client_secret,
        } as any;

        try {
            if (reconnectId) {
                localStorage.removeItem('meli_reconnect');
                const updateResponse: any = await updateIntegrationAction(reconnectId, {
                    store_id: sellerId || undefined,
                    config,
                    credentials,
                });
                if (!updateResponse || updateResponse.success === false) {
                    throw new Error(updateResponse?.message || 'Error al reconectar la integración');
                }
                setStatus('success');
                setMessage('MercadoLibre reconectado exitosamente. Redirigiendo...');
                setTimeout(() => router.push('/integrations'), 2000);
                return;
            }

            const response = await createIntegrationAction({
                name: integrationName,
                code: integrationCode,
                integration_type_id: 3,
                category: 'ecommerce',
                store_id: sellerId,
                is_active: true,
                is_default: false,
                config,
                credentials,
                business_id: businessId,
                is_testing: isTesting,
            }, sessionToken || undefined);

            if (!response || response.success === false) {
                throw new Error(response?.message || 'Error al crear la integración');
            }

            setStatus('success');
            setMessage('MercadoLibre conectado exitosamente. Redirigiendo...');
            setTimeout(() => router.push('/integrations'), 2000);
        } catch (err: any) {
            setStatus('error');
            setMessage(friendlyError(err?.message));
        }
    }, [pending, router]);

    const cancelConnection = () => {
        try {
            localStorage.removeItem('meli_reconnect');
        } catch { }
        router.push('/integrations');
    };

    return (
        <div className="flex min-h-[60vh] items-center justify-center p-6">
            <div className="w-full max-w-md rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 p-8 text-center shadow-sm">
                {(status === 'processing' || status === 'saving') && (
                    <>
                        <svg className="mx-auto h-10 w-10 animate-spin text-[var(--color-primary)]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        <h2 className="mt-4 text-lg font-bold text-gray-900 dark:text-white">Conectando</h2>
                        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{message}</p>
                    </>
                )}
                {status === 'confirm' && pending && (
                    <div className="text-left">
                        <ShieldExclamationIcon className="mx-auto h-12 w-12 text-amber-500" />
                        <h2 className="mt-4 text-center text-lg font-bold text-gray-900 dark:text-white">
                            {pending.reconnectId ? 'Confirma la reconexion' : 'Confirma la conexión'}
                        </h2>
                        <p className="mt-1 text-center text-sm text-gray-500 dark:text-gray-400">
                            Revisa que la cuenta y el negocio sean los correctos antes de guardar.
                        </p>

                        <dl className="mt-5 space-y-3 rounded-xl border border-gray-200 dark:border-gray-700 p-4">
                            <div>
                                <dt className="text-[11px] uppercase tracking-wide text-gray-400">Cuenta de MercadoLibre</dt>
                                <dd className="text-sm font-semibold text-gray-900 dark:text-white break-words">
                                    {pending.accountName || pending.nickname || 'Nombre no disponible'}
                                </dd>
                                {pending.accountName && pending.nickname && (
                                    <dd className="text-xs text-gray-500 dark:text-gray-400 break-words">Usuario {pending.nickname}</dd>
                                )}
                                {pending.accountEmail && (
                                    <dd className="text-xs text-gray-500 dark:text-gray-400 break-words">{pending.accountEmail}</dd>
                                )}
                                <dd className="text-[11px] text-gray-400">Vendedor {pending.sellerId || 'sin id'}</dd>
                            </div>
                            <div>
                                <dt className="text-[11px] uppercase tracking-wide text-gray-400">Se conectara al negocio</dt>
                                <dd className="text-sm font-semibold text-gray-900 dark:text-white break-words">
                                    {pending.businessName || (pending.businessId ? `Negocio #${pending.businessId}` : 'Tu negocio')}
                                </dd>
                            </div>
                        </dl>

                        <p className="mt-4 text-[11px] text-gray-400">
                            MercadoLibre autoriza la cuenta que este abierta en este navegador. Si el nombre de arriba no es
                            el de este negocio, cancela y vuelve a intentar desde una ventana de incognito.
                        </p>

                        <div className="mt-5 flex gap-3">
                            <button
                                type="button"
                                onClick={cancelConnection}
                                className="flex-1 rounded-lg border border-gray-300 dark:border-gray-600 px-4 py-2 text-sm font-semibold text-gray-700 dark:text-gray-200"
                            >
                                Cancelar
                            </button>
                            <button
                                type="button"
                                onClick={confirmConnection}
                                className="flex-1 rounded-lg bg-[var(--color-primary)] px-4 py-2 text-sm font-semibold text-white"
                            >
                                Confirmar
                            </button>
                        </div>
                    </div>
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
