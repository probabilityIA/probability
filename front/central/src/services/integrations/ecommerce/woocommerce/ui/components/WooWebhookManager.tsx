'use client';

import { useState, useEffect, useCallback } from 'react';
import { Alert, Spinner } from '@/shared/ui';
import { BoltIcon, TrashIcon, PlusIcon } from '@heroicons/react/24/outline';
import { useToast } from '@/shared/providers/toast-provider';
import {
    listWebhooksAction,
    createWebhookAction,
    deleteWebhookAction,
} from '@/services/integrations/core/infra/actions';

interface WooWebhookManagerProps {
    integrationId: number;
}

interface WooWebhookInfo {
    id: string;
    address: string;
    topic: string;
    format?: string;
    created_at?: string;
}

const GREEN = 'var(--color-primary)';
const GREEN_DARK = 'color-mix(in srgb, var(--color-primary) 85%, black)';
const GREEN_SOFT = 'color-mix(in srgb, var(--color-primary) 10%, white)';
const GREEN_BORDER = 'color-mix(in srgb, var(--color-primary) 25%, white)';
const INPUT_BORDER = '#e9e9f0';

const TOPIC_LABELS: Record<string, string> = {
    'order.created': 'Orden creada',
    'order.updated': 'Orden actualizada',
    'order.deleted': 'Orden eliminada',
    'order.restored': 'Orden restaurada',
};

function topicLabel(topic: string): string {
    return TOPIC_LABELS[topic] || topic;
}

export function WooWebhookManager({ integrationId }: WooWebhookManagerProps) {
    const { showToast } = useToast();
    const [webhooks, setWebhooks] = useState<WooWebhookInfo[]>([]);
    const [loading, setLoading] = useState(true);
    const [creating, setCreating] = useState(false);
    const [deletingId, setDeletingId] = useState<string | null>(null);
    const [error, setError] = useState<string | null>(null);

    const fetchWebhooks = useCallback(async () => {
        setError(null);
        const result: any = await listWebhooksAction(integrationId);
        if (result.success === false) {
            setError(result.message);
        } else {
            setWebhooks(result.data || []);
        }
        setLoading(false);
    }, [integrationId]);

    useEffect(() => {
        fetchWebhooks();
    }, [fetchWebhooks]);

    const handleCreate = async () => {
        setCreating(true);
        setError(null);
        const result: any = await createWebhookAction(integrationId);
        if (result.success === false) {
            setError(result.message);
            showToast(result.message, 'error');
        } else {
            showToast(result.message || 'Webhooks de ordenes creados', 'success');
            await fetchWebhooks();
        }
        setCreating(false);
    };

    const handleDelete = async (webhookId: string, topic: string) => {
        if (!confirm(`Eliminar webhook "${topicLabel(topic)}"?`)) return;

        setDeletingId(webhookId);
        setError(null);
        const result: any = await deleteWebhookAction(integrationId, webhookId);
        if (result.success === false) {
            setError(result.message);
            showToast(result.message, 'error');
        } else {
            showToast('Webhook eliminado', 'success');
            setWebhooks((prev) => prev.filter((w) => w.id !== webhookId));
        }
        setDeletingId(null);
    };

    const formatDate = (dateStr?: string) => {
        if (!dateStr) return '-';
        try {
            return new Date(dateStr).toLocaleDateString('es-CO', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
            });
        } catch {
            return dateStr;
        }
    };

    const thCls = 'px-3 py-2 text-left text-[10px] font-semibold uppercase tracking-wider text-gray-400 dark:text-gray-500';

    if (loading) {
        return (
            <div className="flex items-center justify-center py-6">
                <Spinner />
                <span className="ml-2 text-[13px] text-gray-500 dark:text-gray-400">Cargando webhooks...</span>
            </div>
        );
    }

    return (
        <div className="space-y-3">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="flex flex-col gap-2 sm:flex-row sm:items-center sm:justify-between">
                <div className="flex items-center gap-2">
                    <span className="flex h-7 w-7 items-center justify-center rounded-md" style={{ backgroundColor: GREEN_SOFT }}>
                        <BoltIcon style={{ color: GREEN, width: 16, height: 16 }} />
                    </span>
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Webhooks de ordenes</h3>
                    {webhooks.length > 0 && (
                        <span
                            className="inline-flex items-center rounded-full px-2 py-0.5 text-[11px] font-semibold"
                            style={{ backgroundColor: GREEN_SOFT, border: `1px solid ${GREEN_BORDER}`, color: GREEN_DARK }}
                        >
                            {webhooks.length} activo{webhooks.length !== 1 ? 's' : ''}
                        </span>
                    )}
                </div>
                <button
                    type="button"
                    onClick={handleCreate}
                    disabled={creating}
                    className="inline-flex items-center justify-center gap-1.5 self-start rounded-lg px-3 py-1.5 text-[12px] font-semibold text-white transition-colors disabled:opacity-60"
                    style={{ backgroundColor: GREEN }}
                    onMouseEnter={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN_DARK; }}
                    onMouseLeave={(e) => { (e.currentTarget as HTMLButtonElement).style.backgroundColor = GREEN; }}
                >
                    {creating ? (
                        <>
                            <svg className="animate-spin h-3.5 w-3.5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                            </svg>
                            Creando...
                        </>
                    ) : (
                        <>
                            <PlusIcon className="w-3.5 h-3.5" />
                            Crear webhooks
                        </>
                    )}
                </button>
            </div>

            <p className="text-[11px] text-gray-400 dark:text-gray-500">
                Probability creara automaticamente los webhooks en tu tienda WooCommerce. Cuando se cree o actualice una orden,
                la recibiras aqui en tiempo real, sin tener que configurarlos a mano.
            </p>

            {webhooks.length === 0 ? (
                <p className="text-[11px] text-gray-400 dark:text-gray-500">
                    No hay webhooks configurados. Haz clic en &quot;Crear webhooks&quot; para recibir las ordenes de WooCommerce.
                </p>
            ) : (
                <div
                    className="overflow-x-auto rounded-lg bg-white dark:bg-gray-800"
                    style={{ border: `1px solid ${INPUT_BORDER}` }}
                >
                    <table className="min-w-full text-sm">
                        <thead className="bg-[#f6f6fa] dark:bg-gray-700/60">
                            <tr>
                                <th className={thCls}>Evento</th>
                                <th className={thCls}>URL destino</th>
                                <th className={thCls}>Creado</th>
                                <th className={thCls}></th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {webhooks.map((wh) => (
                                <tr key={wh.id} className="hover:bg-gray-50/60 dark:hover:bg-gray-700/40">
                                    <td className="px-3 py-2 text-[11px] text-gray-800 dark:text-gray-200 whitespace-nowrap">{topicLabel(wh.topic)}</td>
                                    <td className="px-3 py-2 text-[11px] text-gray-500 dark:text-gray-400 max-w-[260px] truncate" title={wh.address}>
                                        {wh.address}
                                    </td>
                                    <td className="px-3 py-2 text-[11px] text-gray-500 dark:text-gray-400 whitespace-nowrap">{formatDate(wh.created_at)}</td>
                                    <td className="px-3 py-2 text-right">
                                        <button
                                            type="button"
                                            onClick={() => handleDelete(wh.id, wh.topic)}
                                            disabled={deletingId === wh.id}
                                            title="Eliminar webhook"
                                            className="inline-flex h-7 w-7 items-center justify-center rounded-md text-gray-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors disabled:opacity-50"
                                        >
                                            {deletingId === wh.id ? (
                                                <svg className="animate-spin h-4 w-4" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                                                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                                                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                                                </svg>
                                            ) : (
                                                <TrashIcon className="h-4 w-4" />
                                            )}
                                        </button>
                                    </td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>
            )}
        </div>
    );
}
