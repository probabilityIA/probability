'use client';

import { useState, useEffect, useCallback } from 'react';
import { Button, Alert, Spinner } from '@/shared/ui';
import { useToast } from '@/shared/providers/toast-provider';
import {
    listWebhooksAction,
    createWebhookAction,
    deleteWebhookAction,
} from '@/services/integrations/core/infra/actions';
import type { ShopifyWebhookInfo } from '@/services/integrations/core/domain/types';

interface ShopifyWebhookManagerProps {
    integrationId: number;
}

export default function ShopifyWebhookManager({ integrationId }: ShopifyWebhookManagerProps) {
    const { showToast } = useToast();
    const [webhooks, setWebhooks] = useState<ShopifyWebhookInfo[]>([]);
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
            showToast(result.message || 'Webhooks creados correctamente', 'success');
            await fetchWebhooks();
        }
        setCreating(false);
    };

    const handleDelete = async (webhookId: string, topic: string) => {
        if (!confirm(`Eliminar webhook "${topic}"?`)) return;

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

    const formatDate = (dateStr: string) => {
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

    if (loading) {
        return (
            <div className="flex items-center justify-center py-8">
                <Spinner />
                <span className="ml-2 text-sm text-gray-500 dark:text-gray-400 dark:text-gray-400">Cargando webhooks...</span>
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="flex items-center justify-between">
                <p className="text-sm text-gray-500 dark:text-gray-400 dark:text-gray-400">
                    {webhooks.length === 0
                        ? 'No hay webhooks configurados. Crea los webhooks para recibir eventos de Shopify.'
                        : `${webhooks.length} webhook${webhooks.length !== 1 ? 's' : ''} activo${webhooks.length !== 1 ? 's' : ''}`}
                </p>
                <Button
                    onClick={handleCreate}
                    loading={creating}
                    disabled={creating}
                    variant="primary"
                    size="sm"
                >
                    {creating ? 'Creando...' : 'Crear Webhooks'}
                </Button>
            </div>

            {webhooks.length > 0 && (
                <div className="overflow-x-auto border border-gray-200 dark:border-gray-700 rounded-lg">
                    <table className="min-w-full divide-y divide-gray-200 text-sm">
                        <thead className="bg-gray-50 dark:bg-gray-700">
                            <tr>
                                <th className="px-4 py-2 text-left font-medium text-gray-600 dark:text-gray-300 dark:text-gray-300">Topic</th>
                                <th className="px-4 py-2 text-left font-medium text-gray-600 dark:text-gray-300 dark:text-gray-300">URL</th>
                                <th className="px-4 py-2 text-left font-medium text-gray-600 dark:text-gray-300 dark:text-gray-300">Formato</th>
                                <th className="px-4 py-2 text-left font-medium text-gray-600 dark:text-gray-300 dark:text-gray-300">Fecha</th>
                                <th className="px-4 py-2 text-right font-medium text-gray-600 dark:text-gray-300 dark:text-gray-300">Acciones</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 bg-white dark:bg-gray-800">
                            {webhooks.map((wh) => (
                                <tr key={wh.id} className="hover:bg-gray-50 dark:bg-gray-700">
                                    <td className="px-4 py-2 font-mono text-xs">{wh.topic}</td>
                                    <td className="px-4 py-2 text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400 max-w-[300px] truncate" title={wh.address}>
                                        {wh.address}
                                    </td>
                                    <td className="px-4 py-2 text-xs uppercase">{wh.format}</td>
                                    <td className="px-4 py-2 text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400">{formatDate(wh.created_at)}</td>
                                    <td className="px-4 py-2 text-right">
                                        <Button
                                            onClick={() => handleDelete(wh.id, wh.topic)}
                                            loading={deletingId === wh.id}
                                            disabled={deletingId === wh.id}
                                            variant="danger"
                                            size="sm"
                                        >
                                            Eliminar
                                        </Button>
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
