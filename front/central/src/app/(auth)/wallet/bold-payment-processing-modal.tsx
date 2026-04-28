'use client';

import { useEffect, useRef, useState } from 'react';
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/outline';
import { Spinner } from '@/shared/ui';
import { useSSE } from '@/shared/hooks/use-sse';
import { syncBoldRechargeAction } from '@/services/modules/pay/infra/actions';

type Status = 'waiting' | 'success' | 'failed' | 'timeout';

interface BoldPaymentProcessingModalProps {
    open: boolean;
    orderId: string;
    amount: number;
    businessId?: number;
    onClose: () => void;
    onResolved: (status: 'success' | 'failed' | 'timeout', payload?: { newBalance?: number; reason?: string }) => void;
}

const TIMEOUT_MS = 90_000;
const POLLING_INTERVAL_MS = 5_000;
const POLLING_FIRST_DELAY_MS = 8_000;
const EVENT_OK = 'wallet.recharge.completed';
const EVENT_FAIL = 'wallet.recharge.failed';

const formatCOP = (n: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(n);

export function BoldPaymentProcessingModal({
    open,
    orderId,
    amount,
    businessId,
    onClose,
    onResolved,
}: BoldPaymentProcessingModalProps) {
    const [status, setStatus] = useState<Status>('waiting');
    const [reason, setReason] = useState<string | null>(null);
    const [newBalance, setNewBalance] = useState<number | null>(null);
    const [secondsLeft, setSecondsLeft] = useState(Math.floor(TIMEOUT_MS / 1000));
    const resolvedRef = useRef(false);

    useEffect(() => {
        if (!open) {
            resolvedRef.current = false;
            setStatus('waiting');
            setReason(null);
            setNewBalance(null);
            setSecondsLeft(Math.floor(TIMEOUT_MS / 1000));
        }
    }, [open]);

    useEffect(() => {
        if (!open || status !== 'waiting') return;
        const startedAt = Date.now();
        const interval = setInterval(() => {
            const elapsed = Date.now() - startedAt;
            const left = Math.max(0, Math.floor((TIMEOUT_MS - elapsed) / 1000));
            setSecondsLeft(left);
            if (left <= 0 && !resolvedRef.current) {
                resolvedRef.current = true;
                setStatus('timeout');
                onResolved('timeout');
                clearInterval(interval);
            }
        }, 1000);
        return () => clearInterval(interval);
    }, [open, status, onResolved]);

    useEffect(() => {
        if (!open || status !== 'waiting' || !orderId) return;

        let cancelled = false;
        let timer: ReturnType<typeof setTimeout> | undefined;

        const tick = async () => {
            if (cancelled || resolvedRef.current) return;
            try {
                await syncBoldRechargeAction(orderId, businessId);
            } catch {
                // Ignore — el SSE seguirá escuchando y al próximo intento volvemos a probar
            }
            if (!cancelled && !resolvedRef.current) {
                timer = setTimeout(tick, POLLING_INTERVAL_MS);
            }
        };

        timer = setTimeout(tick, POLLING_FIRST_DELAY_MS);
        return () => {
            cancelled = true;
            if (timer) clearTimeout(timer);
        };
    }, [open, status, orderId, businessId]);

    useSSE({
        enabled: open && status === 'waiting' && !!orderId,
        eventTypes: [EVENT_OK, EVENT_FAIL],
        orderIds: orderId ? [orderId] : undefined,
        businessId,
        onMessage: (event) => {
            if (resolvedRef.current) return;
            try {
                const payload = JSON.parse(event.data);
                const evType = payload.type || event.type;
                const data = payload.data || {};
                const dataOrderId = data.order_id || payload.metadata?.order_id;
                if (dataOrderId && dataOrderId !== orderId) return;
                if (evType === EVENT_OK) {
                    resolvedRef.current = true;
                    setNewBalance(typeof data.new_balance === 'number' ? data.new_balance : null);
                    setStatus('success');
                    onResolved('success', { newBalance: data.new_balance });
                } else if (evType === EVENT_FAIL) {
                    resolvedRef.current = true;
                    const r = String(data.reason || 'Pago rechazado por la pasarela');
                    setReason(r);
                    setStatus('failed');
                    onResolved('failed', { reason: r });
                }
            } catch {
                // ignore non-JSON keep-alive frames
            }
        },
    });

    if (!open) return null;

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm">
            <div className="w-full max-w-md mx-4 rounded-2xl bg-white dark:bg-gray-800 shadow-2xl p-6">
                {status === 'waiting' && (
                    <div className="flex flex-col items-center text-center space-y-4">
                        <Spinner />
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                            Procesando pago...
                        </h2>
                        <p className="text-sm text-gray-600 dark:text-gray-300">
                            Esperando confirmación de Bold.
                        </p>
                        <div className="w-full p-3 rounded-lg bg-gray-50 dark:bg-gray-700 text-xs space-y-1">
                            <div className="flex justify-between">
                                <span className="text-gray-500 dark:text-gray-400">Monto</span>
                                <span className="font-mono font-semibold text-gray-900 dark:text-white">{formatCOP(amount)}</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-gray-500 dark:text-gray-400">Orden</span>
                                <span className="font-mono text-gray-700 dark:text-gray-200 truncate ml-2">{orderId}</span>
                            </div>
                            <div className="flex justify-between">
                                <span className="text-gray-500 dark:text-gray-400">Tiempo restante</span>
                                <span className="font-mono text-gray-700 dark:text-gray-200">{secondsLeft}s</span>
                            </div>
                        </div>
                        <button
                            type="button"
                            onClick={onClose}
                            className="text-xs text-gray-500 dark:text-gray-400 hover:underline"
                        >
                            Seguir esperando en segundo plano
                        </button>
                    </div>
                )}

                {status === 'success' && (
                    <div className="flex flex-col items-center text-center space-y-4">
                        <CheckCircleIcon className="w-14 h-14 text-emerald-600" />
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">¡Pago confirmado!</h2>
                        <p className="text-sm text-gray-600 dark:text-gray-300">
                            Tu billetera fue recargada con {formatCOP(amount)}.
                        </p>
                        {newBalance !== null && (
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                                Nuevo saldo: <span className="font-semibold">{formatCOP(newBalance)}</span>
                            </p>
                        )}
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 rounded-lg bg-emerald-600 text-white text-sm hover:bg-emerald-700"
                        >
                            Cerrar
                        </button>
                    </div>
                )}

                {status === 'failed' && (
                    <div className="flex flex-col items-center text-center space-y-4">
                        <XCircleIcon className="w-14 h-14 text-red-600" />
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Pago rechazado</h2>
                        <p className="text-sm text-gray-600 dark:text-gray-300">{reason || 'Bold rechazó la transacción.'}</p>
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 rounded-lg bg-red-600 text-white text-sm hover:bg-red-700"
                        >
                            Cerrar
                        </button>
                    </div>
                )}

                {status === 'timeout' && (
                    <div className="flex flex-col items-center text-center space-y-4">
                        <div className="w-14 h-14 rounded-full bg-yellow-100 dark:bg-yellow-950 flex items-center justify-center">
                            <span className="text-2xl">⏳</span>
                        </div>
                        <h2 className="text-lg font-semibold text-gray-900 dark:text-white">Pago en proceso</h2>
                        <p className="text-sm text-gray-600 dark:text-gray-300">
                            No recibimos confirmación de Bold en 90s. Tu transacción aparecerá como completada en cuanto Bold nos avise.
                        </p>
                        <p className="text-xs text-gray-500 dark:text-gray-400 break-all">Orden: {orderId}</p>
                        <button
                            type="button"
                            onClick={onClose}
                            className="px-4 py-2 rounded-lg bg-gray-700 text-white text-sm hover:bg-gray-800"
                        >
                            Entendido
                        </button>
                    </div>
                )}
            </div>
        </div>
    );
}
