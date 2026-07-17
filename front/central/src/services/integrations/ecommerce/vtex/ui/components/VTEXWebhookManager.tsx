'use client';

import { useState, useEffect, useCallback } from 'react';
import { useToast } from '@/shared/providers/toast-provider';
import {
    getVTEXWebhookStatusAction,
    registerVTEXWebhookAction,
    unregisterVTEXWebhookAction,
} from '../../infra/actions';
import {
    CheckCircleIcon,
    ExclamationTriangleIcon,
    ArrowPathIcon,
    TrashIcon,
} from '@heroicons/react/24/outline';
import { GREEN, GREEN_DARK, GREEN_BORDER, INPUT_BORDER, Spinner } from '@/services/integrations/invoicing/siigo/ui/components/SiigoFormKit';

interface VTEXWebhook {
    id: string;
    address: string;
    statuses: string[];
    is_ours: boolean;
}

interface VTEXWebhookManagerProps {
    integrationId: number;
    businessId?: number | null;
}

export function VTEXWebhookManager({ integrationId, businessId }: VTEXWebhookManagerProps) {
    const { showToast } = useToast();
    const [loading, setLoading] = useState(true);
    const [working, setWorking] = useState(false);
    const [webhook, setWebhook] = useState<VTEXWebhook | null>(null);
    const [webhookUrl, setWebhookUrl] = useState<string>('');
    const [foreignWarning, setForeignWarning] = useState<string | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const res: any = await getVTEXWebhookStatusAction(integrationId, businessId ?? undefined);
            if (res?.success) {
                setWebhook(res.webhook || null);
                setWebhookUrl(res.webhook_url || '');
            }
        } catch {
        } finally {
            setLoading(false);
        }
    }, [integrationId, businessId]);

    useEffect(() => {
        load();
    }, [load]);

    const doRegister = async (force: boolean) => {
        setWorking(true);
        try {
            const res: any = await registerVTEXWebhookAction(integrationId, businessId ?? undefined, force);
            if (res?.success) {
                showToast('Webhook registrado en VTEX', 'success');
                setForeignWarning(null);
                await load();
                return;
            }
            if (res?.foreign_hook) {
                setForeignWarning(res.message || 'La cuenta ya tiene un webhook de otra herramienta');
                return;
            }
            showToast(res?.message || 'Error al registrar el webhook', 'error');
        } catch (err: any) {
            showToast(err.message || 'Error al registrar el webhook', 'error');
        } finally {
            setWorking(false);
        }
    };

    const doUnregister = async () => {
        setWorking(true);
        try {
            const res: any = await unregisterVTEXWebhookAction(integrationId, businessId ?? undefined);
            if (res?.success) {
                showToast('Webhook eliminado en VTEX', 'success');
                setForeignWarning(null);
                await load();
            } else {
                showToast(res?.message || 'Error al eliminar el webhook', 'error');
            }
        } catch (err: any) {
            showToast(err.message || 'Error al eliminar el webhook', 'error');
        } finally {
            setWorking(false);
        }
    };

    if (loading) {
        return (
            <div className="flex items-center gap-2 text-[12px] text-gray-400">
                <Spinner className="animate-spin h-4 w-4" />
                Consultando el webhook en VTEX...
            </div>
        );
    }

    const registeredByUs = webhook?.is_ours === true;
    const registeredByOther = webhook != null && !webhook.is_ours;

    return (
        <div className="space-y-3">
            <p className="text-[11px] text-gray-500 dark:text-gray-400 leading-snug">
                VTEX admite un solo webhook por cuenta. Al registrarlo se reemplaza cualquier hook anterior, incluido el
                de otra herramienta.
            </p>

            <div className="rounded-lg bg-white dark:bg-gray-800 p-3" style={{ border: `1px solid ${INPUT_BORDER}` }}>
                <p className="text-[11px] font-semibold text-gray-500 dark:text-gray-400 mb-1">URL de Probability</p>
                <p className="text-[11px] font-mono text-gray-700 dark:text-gray-200 break-all">{webhookUrl || '-'}</p>
            </div>

            {registeredByUs && (
                <div className="flex items-start gap-2 rounded-lg p-3" style={{ backgroundColor: 'color-mix(in srgb, var(--color-primary) 8%, white)', border: `1px solid ${GREEN_BORDER}` }}>
                    <CheckCircleIcon className="w-4 h-4 mt-0.5 flex-shrink-0" style={{ color: GREEN }} />
                    <div className="min-w-0">
                        <p className="text-[12px] font-semibold text-gray-900 dark:text-white">Webhook activo</p>
                        <p className="text-[11px] text-gray-600 dark:text-gray-300 mt-0.5">
                            VTEX esta enviando las ordenes a Probability en tiempo real.
                        </p>
                        {webhook.statuses?.length > 0 && (
                            <p className="text-[11px] text-gray-500 dark:text-gray-400 mt-1 font-mono break-all">
                                Estados: {webhook.statuses.join(', ')}
                            </p>
                        )}
                    </div>
                </div>
            )}

            {registeredByOther && (
                <div className="flex items-start gap-2 rounded-lg p-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800">
                    <ExclamationTriangleIcon className="w-4 h-4 mt-0.5 flex-shrink-0 text-amber-600 dark:text-amber-500" />
                    <div className="min-w-0">
                        <p className="text-[12px] font-semibold text-amber-900 dark:text-amber-300">
                            La cuenta ya tiene un webhook de otra herramienta
                        </p>
                        <p className="text-[11px] text-amber-800 dark:text-amber-400 mt-0.5 font-mono break-all">
                            {webhook.address}
                        </p>
                        <p className="text-[11px] text-amber-800 dark:text-amber-400 mt-1">
                            Si lo reemplazas, esa herramienta dejara de recibir las ordenes de VTEX.
                        </p>
                    </div>
                </div>
            )}

            {!webhook && (
                <div className="flex items-start gap-2 rounded-lg p-3 bg-gray-50 dark:bg-gray-800 border border-gray-200 dark:border-gray-700">
                    <ExclamationTriangleIcon className="w-4 h-4 mt-0.5 flex-shrink-0 text-gray-400" />
                    <div className="min-w-0">
                        <p className="text-[12px] font-semibold text-gray-900 dark:text-white">Sin webhook registrado</p>
                        <p className="text-[11px] text-gray-600 dark:text-gray-300 mt-0.5">
                            Las ordenes solo entraran cuando sincronices manualmente.
                        </p>
                    </div>
                </div>
            )}

            {foreignWarning && (
                <div className="rounded-lg p-3 bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800">
                    <p className="text-[12px] font-semibold text-amber-900 dark:text-amber-300 mb-1">
                        Confirma el reemplazo
                    </p>
                    <p className="text-[11px] text-amber-800 dark:text-amber-400 break-all">{foreignWarning}</p>
                    <div className="mt-2 flex gap-2">
                        <button
                            type="button"
                            onClick={() => doRegister(true)}
                            disabled={working}
                            className="rounded-lg px-3 py-1.5 text-[12px] font-semibold text-white bg-amber-600 hover:bg-amber-700 disabled:opacity-60"
                        >
                            Reemplazar de todos modos
                        </button>
                        <button
                            type="button"
                            onClick={() => setForeignWarning(null)}
                            disabled={working}
                            className="rounded-lg px-3 py-1.5 text-[12px] font-semibold bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200 disabled:opacity-60"
                            style={{ border: `1px solid ${INPUT_BORDER}` }}
                        >
                            Cancelar
                        </button>
                    </div>
                </div>
            )}

            <div className="flex flex-col gap-2 sm:flex-row">
                <button
                    type="button"
                    onClick={() => doRegister(false)}
                    disabled={working}
                    className="flex-1 inline-flex items-center justify-center gap-1.5 rounded-lg py-2 text-[12px] font-semibold text-white transition-colors disabled:opacity-60"
                    style={{ backgroundColor: GREEN }}
                    onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                    onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                >
                    {working ? <Spinner className="animate-spin h-3.5 w-3.5" /> : <ArrowPathIcon className="w-3.5 h-3.5" />}
                    {registeredByUs ? 'Volver a registrar' : 'Registrar webhook'}
                </button>

                {registeredByUs && (
                    <button
                        type="button"
                        onClick={doUnregister}
                        disabled={working}
                        className="inline-flex items-center justify-center gap-1.5 rounded-lg px-3 py-2 text-[12px] font-semibold text-red-600 bg-white dark:bg-gray-800 hover:bg-red-50 dark:hover:bg-red-900/20 disabled:opacity-60"
                        style={{ border: `1px solid ${INPUT_BORDER}` }}
                    >
                        <TrashIcon className="w-3.5 h-3.5" />
                        Eliminar
                    </button>
                )}
            </div>
        </div>
    );
}
