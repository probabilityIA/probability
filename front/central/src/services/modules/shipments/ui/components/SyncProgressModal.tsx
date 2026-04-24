'use client';

import { useEffect, useMemo, useRef, useState } from 'react';
import { X, CheckCircle2, Loader2, AlertCircle, Package, RefreshCw } from 'lucide-react';
import { syncShipmentStatusAction } from '../../infra/actions';
import { useShipmentSSE } from '../hooks/useShipmentSSE';
import type { ShipmentSSEEventData } from '../../domain/types';

interface SyncProgressModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId: number | null;
    onCompleted?: () => void;
}

interface SyncUpdate {
    shipmentId?: number;
    orderNumber?: string;
    customerName?: string;
    trackingNumber: string;
    previousStatus?: string;
    newStatus: string;
    rawStatus?: string;
    timestamp: number;
}

const STATUS_COLORS: Record<string, string> = {
    pending: 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300',
    picked_up: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/40 dark:text-indigo-300',
    in_transit: 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300',
    out_for_delivery: 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-300',
    delivered: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300',
    on_hold: 'bg-orange-100 text-orange-700 dark:bg-orange-900/40 dark:text-orange-300',
    failed: 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300',
    returned: 'bg-rose-100 text-rose-700 dark:bg-rose-900/40 dark:text-rose-300',
    cancelled: 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300',
};

const STATUS_LABELS: Record<string, string> = {
    pending: 'Pendiente',
    picked_up: 'Recolectado',
    in_transit: 'En Tránsito',
    out_for_delivery: 'En Reparto',
    delivered: 'Entregado',
    on_hold: 'Novedad',
    failed: 'Fallido',
    returned: 'Devuelto',
    cancelled: 'Cancelado',
};

export function SyncProgressModal({ isOpen, onClose, businessId, onCompleted }: SyncProgressModalProps) {
    const [phase, setPhase] = useState<'idle' | 'starting' | 'running' | 'done' | 'error'>('idle');
    const [correlationPrefix, setCorrelationPrefix] = useState<string | null>(null);
    const [total, setTotal] = useState(0);
    const [batches, setBatches] = useState(0);
    const [updates, setUpdates] = useState<SyncUpdate[]>([]);
    const [errorMessage, setErrorMessage] = useState<string | null>(null);

    const effectiveBusinessId = useMemo(() => businessId ?? 0, [businessId]);
    const businessIdRef = useRef(businessId);
    const correlationPrefixRef = useRef<string | null>(null);
    const phaseRef = useRef(phase);
    useEffect(() => {
        businessIdRef.current = businessId;
    }, [businessId]);
    useEffect(() => {
        correlationPrefixRef.current = correlationPrefix;
    }, [correlationPrefix]);
    useEffect(() => {
        phaseRef.current = phase;
    }, [phase]);

    useShipmentSSE({
        businessId: effectiveBusinessId,
        onTrackingUpdated: (data: ShipmentSSEEventData) => {
            const currentPhase = phaseRef.current;
            if (currentPhase === 'idle' || currentPhase === 'done' || currentPhase === 'error') return;

            const corr = (data.correlation_id as string | undefined) || '';
            const prefix = correlationPrefixRef.current;

            const matches = prefix ? corr.startsWith(prefix) : corr.startsWith('sync-');
            if (!matches) return;

            const payload = (data.tracking as Record<string, any> | undefined) || (data as Record<string, any>);

            setUpdates((prev) => {
                const next: SyncUpdate = {
                    shipmentId: payload.shipment_id as number | undefined,
                    orderNumber: payload.order_number as string | undefined,
                    customerName: payload.customer_name as string | undefined,
                    trackingNumber: (payload.tracking_number as string) || '',
                    previousStatus: payload.previous_status as string | undefined,
                    newStatus: (payload.new_status as string) || (payload.probability_status as string) || '',
                    rawStatus: payload.raw_status as string | undefined,
                    timestamp: Date.now(),
                };
                return [next, ...prev].slice(0, 200);
            });
        },
    });

    useEffect(() => {
        if (!isOpen) return;

        let cancelled = false;
        const run = async () => {
            setPhase('starting');
            setUpdates([]);
            setErrorMessage(null);

            const currentBusinessId = businessIdRef.current;
            if (!currentBusinessId) {
                if (!cancelled) {
                    setErrorMessage('Selecciona un negocio antes de sincronizar');
                    setPhase('error');
                }
                return;
            }

            const result: any = await syncShipmentStatusAction({
                provider: 'envioclick',
                business_id: currentBusinessId,
            });

            if (cancelled) return;

            if (!result.success) {
                setErrorMessage(result.message || result.error || 'No se pudo iniciar la sincronización');
                setPhase('error');
                return;
            }

            const corr = result.correlation_id as string | undefined;
            const totalItems = (result.total_shipments as number | undefined) ?? 0;

            if (totalItems === 0) {
                setErrorMessage(result.message || 'No hay envíos para sincronizar en el rango indicado');
                setPhase('error');
                return;
            }

            setCorrelationPrefix(corr ?? null);
            setTotal(totalItems);
            setBatches((result.batches as number | undefined) ?? 0);
            setPhase('running');
        };
        run();

        return () => {
            cancelled = true;
        };
    }, [isOpen]);

    useEffect(() => {
        if (phase === 'running' && total > 0 && updates.length >= total) {
            setPhase('done');
            onCompleted?.();
        }
    }, [phase, total, updates.length, onCompleted]);

    useEffect(() => {
        if (!isOpen) {
            setPhase('idle');
            setCorrelationPrefix(null);
            setTotal(0);
            setBatches(0);
            setUpdates([]);
            setErrorMessage(null);
        }
    }, [isOpen]);

    if (!isOpen) return null;

    const processed = updates.length;
    const progressPct = total > 0 ? Math.min(100, Math.round((processed / total) * 100)) : 0;

    const stats = updates.reduce<Record<string, number>>((acc, u) => {
        acc[u.newStatus] = (acc[u.newStatus] || 0) + 1;
        return acc;
    }, {});

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm p-4">
            <div className="bg-white dark:bg-gray-900 rounded-2xl shadow-2xl w-full max-w-3xl max-h-[90vh] flex flex-col overflow-hidden border border-gray-200 dark:border-gray-700">
                <div className="flex items-center justify-between px-6 py-5 border-b border-gray-100 dark:border-gray-800 bg-gradient-to-r from-emerald-50 to-teal-50 dark:from-emerald-950/40 dark:to-teal-950/40">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-emerald-500/10 dark:bg-emerald-400/10 flex items-center justify-center">
                            <RefreshCw size={18} className={`text-emerald-600 dark:text-emerald-400 ${phase === 'running' || phase === 'starting' ? 'animate-spin' : ''}`} />
                        </div>
                        <div>
                            <h2 className="text-lg font-bold text-gray-900 dark:text-white">Sincronización de Estados de Guías</h2>
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                                {phase === 'starting' && 'Iniciando...'}
                                {phase === 'running' && `${processed} de ${total} envíos · ${batches} batch${batches !== 1 ? 'es' : ''}`}
                                {phase === 'done' && `Completado · ${processed} envíos actualizados`}
                                {phase === 'error' && 'Error'}
                            </p>
                        </div>
                    </div>
                    <button
                        onClick={onClose}
                        disabled={phase === 'starting' || phase === 'running'}
                        className="p-2 rounded-lg hover:bg-white/50 dark:hover:bg-white/10 text-gray-500 dark:text-gray-400 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
                    >
                        <X size={18} />
                    </button>
                </div>

                {phase === 'error' && (
                    <div className="p-6 flex flex-col items-center text-center gap-3">
                        <AlertCircle size={48} className="text-red-500" />
                        <p className="text-gray-700 dark:text-gray-300 font-medium">{errorMessage}</p>
                        <button
                            onClick={onClose}
                            className="px-4 py-2 bg-gray-200 dark:bg-gray-700 hover:bg-gray-300 dark:hover:bg-gray-600 rounded-lg text-sm font-semibold text-gray-800 dark:text-gray-200 transition-colors"
                        >
                            Cerrar
                        </button>
                    </div>
                )}

                {(phase === 'starting' || phase === 'running' || phase === 'done') && (
                    <>
                        <div className="px-6 py-4 border-b border-gray-100 dark:border-gray-800">
                            <div className="flex items-center justify-between mb-2">
                                <span className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wider">Progreso</span>
                                <span className="text-sm font-bold text-gray-900 dark:text-white tabular-nums">{progressPct}%</span>
                            </div>
                            <div className="h-2 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
                                <div
                                    className={`h-full rounded-full transition-all duration-300 ease-out ${
                                        phase === 'done'
                                            ? 'bg-gradient-to-r from-emerald-500 to-teal-500'
                                            : 'bg-gradient-to-r from-emerald-400 to-teal-400'
                                    }`}
                                    style={{ width: `${progressPct}%` }}
                                />
                            </div>
                            {Object.keys(stats).length > 0 && (
                                <div className="flex flex-wrap gap-2 mt-3">
                                    {Object.entries(stats).map(([status, count]) => (
                                        <span
                                            key={status}
                                            className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold ${
                                                STATUS_COLORS[status] || 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300'
                                            }`}
                                        >
                                            {STATUS_LABELS[status] || status}
                                            <span className="opacity-70 tabular-nums">{count}</span>
                                        </span>
                                    ))}
                                </div>
                            )}
                        </div>

                        <div className="flex-1 overflow-y-auto px-4 py-3 space-y-2 bg-gray-50/50 dark:bg-gray-900/50">
                            {phase === 'starting' && (
                                <div className="flex items-center justify-center py-12 gap-3 text-gray-500 dark:text-gray-400">
                                    <Loader2 size={20} className="animate-spin" />
                                    <span className="text-sm">Preparando sincronización...</span>
                                </div>
                            )}

                            {phase !== 'starting' && updates.length === 0 && (
                                <div className="flex flex-col items-center justify-center py-12 gap-3 text-gray-400 dark:text-gray-500">
                                    <Package size={32} />
                                    <span className="text-sm">Esperando actualizaciones...</span>
                                </div>
                            )}

                            {updates.map((u, idx) => {
                                const changed = u.previousStatus && u.previousStatus !== u.newStatus;
                                return (
                                    <div
                                        key={`${u.shipmentId}-${u.timestamp}-${idx}`}
                                        className={`bg-white dark:bg-gray-800 border rounded-xl p-3 shadow-sm transition-all ${
                                            changed
                                                ? 'border-emerald-200 dark:border-emerald-800'
                                                : 'border-gray-100 dark:border-gray-700'
                                        }`}
                                    >
                                        <div className="flex items-center justify-between gap-3">
                                            <div className="flex-1 min-w-0">
                                                <div className="flex items-center gap-2 mb-1">
                                                    <span className="text-xs font-mono text-gray-500 dark:text-gray-400">
                                                        {u.orderNumber ? `Orden ${u.orderNumber}` : 'Orden -'}
                                                    </span>
                                                    <span className="text-gray-300 dark:text-gray-600">·</span>
                                                    <span className="text-xs font-mono text-gray-500 dark:text-gray-400 truncate">
                                                        {u.trackingNumber}
                                                    </span>
                                                </div>
                                                {u.customerName && (
                                                    <p className="text-sm font-semibold text-gray-800 dark:text-gray-200 truncate">
                                                        {u.customerName}
                                                    </p>
                                                )}
                                            </div>
                                            <div className="flex items-center gap-2 flex-shrink-0">
                                                {changed && u.previousStatus && (
                                                    <>
                                                        <span
                                                            className={`px-2 py-0.5 rounded-full text-[10px] font-semibold ${
                                                                STATUS_COLORS[u.previousStatus] || 'bg-gray-100 text-gray-600'
                                                            } opacity-60`}
                                                        >
                                                            {STATUS_LABELS[u.previousStatus] || u.previousStatus}
                                                        </span>
                                                        <span className="text-gray-400 dark:text-gray-600 text-xs">→</span>
                                                    </>
                                                )}
                                                <span
                                                    className={`inline-flex items-center gap-1 px-2.5 py-1 rounded-full text-xs font-semibold ${
                                                        STATUS_COLORS[u.newStatus] || 'bg-gray-100 text-gray-700'
                                                    }`}
                                                >
                                                    {phase === 'running' && idx === 0 && (
                                                        <CheckCircle2 size={12} className="text-current" />
                                                    )}
                                                    {STATUS_LABELS[u.newStatus] || u.newStatus}
                                                </span>
                                            </div>
                                        </div>
                                        {u.rawStatus && u.rawStatus !== STATUS_LABELS[u.newStatus] && (
                                            <p className="text-[10px] text-gray-400 dark:text-gray-500 mt-1.5 italic">
                                                Reporte carrier: &quot;{u.rawStatus}&quot;
                                            </p>
                                        )}
                                    </div>
                                );
                            })}
                        </div>

                        <div className="px-6 py-4 border-t border-gray-100 dark:border-gray-800 bg-gray-50 dark:bg-gray-900/50 flex items-center justify-end gap-3">
                            {phase === 'done' ? (
                                <button
                                    onClick={onClose}
                                    className="px-5 py-2 bg-emerald-600 hover:bg-emerald-700 text-white rounded-lg text-sm font-semibold transition-colors flex items-center gap-2"
                                >
                                    <CheckCircle2 size={16} />
                                    Listo
                                </button>
                            ) : (
                                <span className="text-xs text-gray-500 dark:text-gray-400 italic">
                                    {phase === 'running' ? 'Actualizando en tiempo real...' : 'Procesando...'}
                                </span>
                            )}
                        </div>
                    </>
                )}
            </div>
        </div>
    );
}
