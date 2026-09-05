'use client';

import { useCallback, useEffect, useState } from 'react';
import { ArrowPathIcon } from '@heroicons/react/24/outline';
import { Alert, Button } from '@/shared/ui';
import {
    getWhatsAppTemplatesStatusAction,
    provisionWhatsAppTemplatesAction,
} from '../../infra/actions';
import { WhatsAppTemplateStatus } from '../../domain/types';

interface WhatsAppTemplatesPanelProps {
    businessId?: number;
    enabled: boolean;
}

const STATUS_STYLES: Record<string, string> = {
    APPROVED: 'bg-green-100 text-green-800',
    PENDING: 'bg-yellow-100 text-yellow-800',
    REJECTED: 'bg-red-100 text-red-700',
    DISABLED: 'bg-gray-200 text-gray-700',
    ERROR: 'bg-red-100 text-red-700',
};

const STATUS_LABELS: Record<string, string> = {
    APPROVED: 'Aprobada',
    PENDING: 'En revisión',
    REJECTED: 'Rechazada',
    DISABLED: 'Deshabilitada',
    ERROR: 'Error',
    UNKNOWN: 'Sin estado',
};

export default function WhatsAppTemplatesPanel({ businessId, enabled }: WhatsAppTemplatesPanelProps) {
    const [templates, setTemplates] = useState<WhatsAppTemplateStatus[]>([]);
    const [wabaId, setWabaId] = useState('');
    const [loading, setLoading] = useState(false);
    const [provisioning, setProvisioning] = useState(false);
    const [message, setMessage] = useState<{ type: 'success' | 'error' | 'info'; text: string } | null>(null);

    const load = useCallback(
        async (refresh: boolean) => {
            setLoading(true);
            setMessage(null);
            try {
                const result = await getWhatsAppTemplatesStatusAction(businessId, refresh);
                if (result.success && result.data) {
                    setTemplates(result.data.templates || []);
                    setWabaId(result.data.waba_id || '');
                } else {
                    setMessage({ type: 'error', text: result.message || 'No se pudo consultar el estado' });
                }
            } catch (err: any) {
                setMessage({ type: 'error', text: err?.message || 'No se pudo consultar el estado' });
            } finally {
                setLoading(false);
            }
        },
        [businessId]
    );

    useEffect(() => {
        if (enabled) {
            load(false);
        }
    }, [enabled, load]);

    const handleProvision = async () => {
        setProvisioning(true);
        setMessage(null);
        try {
            const result = await provisionWhatsAppTemplatesAction(businessId);
            if (result.success && result.data) {
                setTemplates(result.data.templates || []);
                setWabaId(result.data.waba_id || '');
                const failed = Object.keys(result.data.failed || {});
                setMessage({
                    type: failed.length > 0 ? 'error' : 'success',
                    text: `Creadas: ${result.data.created.length}. Ya existían: ${result.data.already_exists.length}.${
                        failed.length > 0 ? ` Fallaron: ${failed.join(', ')}` : ''
                    }`,
                });
            } else {
                setMessage({ type: 'error', text: result.message || 'No se pudieron crear las plantillas' });
            }
        } catch (err: any) {
            setMessage({ type: 'error', text: err?.message || 'No se pudieron crear las plantillas' });
        } finally {
            setProvisioning(false);
        }
    };

    if (!enabled) {
        return (
            <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <h4 className="text-sm font-medium text-gray-700 dark:text-gray-200">Plantillas</h4>
                <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                    Estás usando el número de Probability, que ya tiene todas las plantillas aprobadas.
                    Conecta tu propio número para gestionar las tuyas.
                </p>
            </div>
        );
    }

    const pending = templates.filter((t) => t.status !== 'APPROVED').length;

    return (
        <div className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3">
            <div className="flex items-start justify-between gap-3">
                <div>
                    <h4 className="text-sm font-medium text-gray-700 dark:text-gray-200">Plantillas</h4>
                    <p className="text-xs text-gray-500 dark:text-gray-400 mt-1">
                        {wabaId ? `Cuenta de WhatsApp ${wabaId}. ` : ''}
                        Meta puede tardar horas en aprobarlas. Mientras una plantilla no esté aprobada, los
                        mensajes que la usan no se pueden enviar.
                    </p>
                </div>
                <button
                    type="button"
                    onClick={() => load(true)}
                    disabled={loading}
                    className="p-2 rounded-md border border-gray-300 dark:border-gray-600 hover:bg-gray-100 dark:hover:bg-gray-700 shrink-0"
                    title="Actualizar estado"
                >
                    <ArrowPathIcon className={`w-4 h-4 ${loading ? 'animate-spin' : ''}`} />
                </button>
            </div>

            {templates.length === 0 ? (
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    {loading ? 'Consultando plantillas...' : 'Todavía no hay plantillas en tu cuenta.'}
                </p>
            ) : (
                <>
                    <p className="text-xs text-gray-500 dark:text-gray-400">
                        {templates.length} plantillas, {pending} pendientes de aprobación.
                    </p>
                    <div className="max-h-72 overflow-y-auto divide-y divide-gray-100 dark:divide-gray-700">
                        {templates.map((template) => (
                            <div
                                key={`${template.name}-${template.language}`}
                                className="py-2 flex items-center justify-between gap-3"
                            >
                                <div className="min-w-0">
                                    <p className="text-sm text-gray-900 dark:text-white font-mono truncate">
                                        {template.name}
                                    </p>
                                    {template.reason && (
                                        <p className="text-xs text-red-600 dark:text-red-400 truncate">
                                            {template.reason}
                                        </p>
                                    )}
                                </div>
                                <span
                                    className={`px-2 py-1 rounded-full text-xs font-medium shrink-0 ${
                                        STATUS_STYLES[template.status] || 'bg-gray-100 text-gray-700'
                                    }`}
                                >
                                    {STATUS_LABELS[template.status] || template.status}
                                </span>
                            </div>
                        ))}
                    </div>
                </>
            )}

            <Button
                type="button"
                variant="outline"
                onClick={handleProvision}
                disabled={provisioning}
                loading={provisioning}
                className="w-full"
            >
                Crear las plantillas que faltan
            </Button>

            {message && (
                <Alert type={message.type === 'info' ? 'success' : message.type} onClose={() => setMessage(null)}>
                    {message.text}
                </Alert>
            )}
        </div>
    );
}
