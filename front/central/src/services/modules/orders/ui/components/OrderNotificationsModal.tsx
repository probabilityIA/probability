'use client';

import { useCallback, useEffect, useState } from 'react';
import { OrderApiRepository } from '../../infra/repository/api-repository';
import { OrderNotifications, OrderNotificationEvent } from '../../domain/types';

interface Props {
    orderId: string;
    orderNumber: string;
    businessId?: number;
    onClose: () => void;
    onChanged?: () => void;
}

const EVENT_LABELS: Record<string, string> = {
    'order.created': 'Confirmacion de pedido',
    'order.shipped': 'Pedido enviado',
    'order.delivered': 'Pedido entregado',
    'order.canceled': 'Pedido cancelado',
    'invoice.created': 'Factura generada',
    'shipment.guide_generated': 'Guia de envio generada',
};

const STATUS_STYLES: Record<string, string> = {
    sent: 'bg-green-100 text-green-800',
    failed: 'bg-red-100 text-red-800',
    pending: 'bg-gray-100 text-gray-600',
};

const STATUS_LABELS: Record<string, string> = {
    sent: 'Enviado',
    failed: 'Fallido',
    pending: 'Pendiente',
};

function formatTime(value: string): string {
    try {
        return new Date(value).toLocaleString('es-CO', {
            day: '2-digit',
            month: '2-digit',
            hour: '2-digit',
            minute: '2-digit',
        });
    } catch {
        return value;
    }
}

export function OrderNotificationsModal({ orderId, orderNumber, businessId, onClose, onChanged }: Props) {
    const [data, setData] = useState<OrderNotifications | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [resending, setResending] = useState<string | null>(null);

    const load = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const repo = new OrderApiRepository();
            setData(await repo.getOrderNotifications(orderId, businessId));
        } catch (e) {
            setError(e instanceof Error ? e.message : 'Error cargando notificaciones');
        } finally {
            setLoading(false);
        }
    }, [orderId, businessId]);

    useEffect(() => {
        load();
    }, [load]);

    const handleResend = async (event: OrderNotificationEvent) => {
        if (!event.resend_action) return;
        setResending(event.event_code);
        setError(null);
        try {
            const repo = new OrderApiRepository();
            await repo.resendOrderNotification(orderId, event.resend_action, businessId);
            await load();
            onChanged?.();
        } catch (e) {
            setError(e instanceof Error ? e.message : 'Error reenviando la notificacion');
        } finally {
            setResending(null);
        }
    };

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4" onClick={onClose}>
            <div
                className="w-full max-w-lg max-h-[85vh] overflow-hidden rounded-2xl bg-white shadow-2xl flex flex-col"
                onClick={(e) => e.stopPropagation()}
            >
                <div className="flex items-center justify-between border-b px-5 py-4">
                    <div className="flex items-center gap-2">
                        <span className="text-xl" aria-hidden>&#128172;</span>
                        <div>
                            <h3 className="text-sm font-bold text-gray-900">Notificaciones WhatsApp</h3>
                            <p className="text-xs text-gray-500">Orden {orderNumber}</p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        className="rounded-lg px-2 py-1 text-gray-400 hover:bg-gray-100 hover:text-gray-600"
                        aria-label="Cerrar"
                    >
                        &#10005;
                    </button>
                </div>

                {loading && <div className="px-5 py-8 text-center text-sm text-gray-500">Cargando...</div>}

                {error && (
                    <div className="mx-5 mt-4 rounded-lg bg-red-50 px-3 py-2 text-xs text-red-700">{error}</div>
                )}

                {!loading && data && (
                    <div className="flex-1 overflow-y-auto">
                        <div className="border-b bg-gray-50 px-5 py-3">
                            <p className="mb-2 text-xs font-semibold uppercase tracking-wide text-gray-500">
                                Enviadas {data.sent}/{data.expected}
                            </p>
                            <div className="space-y-2">
                                {data.events.length === 0 && (
                                    <p className="text-xs text-gray-500">
                                        Este canal no tiene notificaciones configuradas para la integracion de la orden.
                                    </p>
                                )}
                                {data.events.map((ev) => (
                                    <div key={ev.event_code} className="flex items-center justify-between gap-2">
                                        <div className="flex items-center gap-2">
                                            <span className={`rounded-full px-2 py-0.5 text-[10px] font-bold ${STATUS_STYLES[ev.status]}`}>
                                                {STATUS_LABELS[ev.status]}
                                            </span>
                                            <span className="text-xs text-gray-700">
                                                {EVENT_LABELS[ev.event_code] || ev.event_code}
                                            </span>
                                        </div>
                                        {ev.can_resend && (
                                            <button
                                                onClick={() => handleResend(ev)}
                                                disabled={resending === ev.event_code}
                                                className="rounded-md bg-cyan-600 px-2.5 py-1 text-[11px] font-semibold text-white hover:bg-cyan-700 disabled:opacity-50"
                                            >
                                                {resending === ev.event_code ? 'Enviando...' : 'Reenviar'}
                                            </button>
                                        )}
                                    </div>
                                ))}
                            </div>
                        </div>

                        <div className="space-y-3 px-5 py-4">
                            {data.messages.length === 0 && (
                                <p className="text-center text-xs text-gray-500">
                                    Todavia no hay mensajes para esta orden.
                                </p>
                            )}
                            {data.messages.map((m) => {
                                const outbound = m.direction === 'outbound';
                                return (
                                    <div key={m.id} className={`flex ${outbound ? 'justify-end' : 'justify-start'}`}>
                                        <div
                                            className={`max-w-[80%] rounded-2xl px-3 py-2 text-xs whitespace-pre-wrap break-words ${
                                                outbound ? 'bg-green-50 text-gray-800' : 'bg-gray-100 text-gray-800'
                                            }`}
                                        >
                                            {m.content || <span className="italic text-gray-400">(sin contenido)</span>}
                                            <div className="mt-1 flex items-center justify-end gap-1 text-[10px] text-gray-400">
                                                <span>{formatTime(m.created_at)}</span>
                                                {outbound && (
                                                    <span className={m.status === 'failed' ? 'text-red-500' : ''}>
                                                        {m.status === 'failed' ? 'fallido' : m.status}
                                                    </span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                );
                            })}
                        </div>
                    </div>
                )}
            </div>
        </div>
    );
}
