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
        try {
            setError(null);
            const result = await listWebhooksAction(integrationId);
            setWebhooks(result.data || []);
        } catch (err: any) {
            setError(err.message || 'Error al cargar webhooks');
        } finally {
            setLoading(false);
        }
    }, [integrationId]);

    useEffect(() => {
        fetchWebhooks();
    }, [fetchWebhooks]);

    const handleCreate = async () => {
        setCreating(true);
        setError(null);
        try {
            const result = await createWebhookAction(integrationId);
            showToast(result.message || 'Webhooks creados correctamente', 'success');
            await fetchWebhooks();
        } catch (err: any) {
            const msg = err.message || 'Error al crear webhooks';
            setError(msg);
            showToast(msg, 'error');
        } finally {
            setCreating(false);
        }
    };

    const handleDelete = async (webhookId: string, topic: string) => {
        if (!confirm(`Eliminar webhook "${topic}"?`)) return;

        setDeletingId(webhookId);
        setError(null);
        try {
            await deleteWebhookAction(integrationId, webhookId);
            showToast('Webhook eliminado', 'success');
            setWebhooks((prev) => prev.filter((w) => w.id !== webhookId));
        } catch (err: any) {
            const msg = err.message || 'Error al eliminar webhook';
            setError(msg);
            showToast(msg, 'error');
        } finally {
            setDeletingId(null);
        }
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
                <span className="ml-2 text-sm text-gray-500">Cargando webhooks...</span>
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
                <p className="text-sm text-gray-500">
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
                <div className="overflow-x-auto border border-gray-200 rounded-lg">
                    <table className="min-w-full divide-y divide-gray-200 text-sm">
                        <thead className="bg-gray-50">
                            <tr>
                                <th className="px-4 py-2 text-left font-medium text-gray-600">Topic</th>
                                <th className="px-4 py-2 text-left font-medium text-gray-600">URL</th>
                                <th className="px-4 py-2 text-left font-medium text-gray-600">Formato</th>
                                <th className="px-4 py-2 text-left font-medium text-gray-600">Fecha</th>
                                <th className="px-4 py-2 text-right font-medium text-gray-600">Acciones</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 bg-white">
                            {webhooks.map((wh) => (
                                <tr key={wh.id} className="hover:bg-gray-50">
                                    <td className="px-4 py-2 font-mono text-xs">{wh.topic}</td>
                                    <td className="px-4 py-2 text-xs text-gray-500 max-w-[300px] truncate" title={wh.address}>
                                        {wh.address}
                                    </td>
                                    <td className="px-4 py-2 text-xs uppercase">{wh.format}</td>
                                    <td className="px-4 py-2 text-xs text-gray-500">{formatDate(wh.created_at)}</td>
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
