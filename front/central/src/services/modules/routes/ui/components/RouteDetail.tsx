'use client';

import { useState, useEffect, useCallback } from 'react';
import { ArrowLeftIcon, PlayIcon, CheckCircleIcon, XCircleIcon, CheckIcon } from '@heroicons/react/24/outline';
import { RouteDetail as RouteDetailType, RouteStopInfo } from '../../domain/types';
import {
    getRouteByIdAction,
    startRouteAction,
    completeRouteAction,
    updateStopStatusAction,
} from '../../infra/actions';
import { Alert, Spinner, Button } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface RouteDetailProps {
    routeId: number;
    businessId?: number;
    onBack: () => void;
    onRefreshList?: () => void;
}

const ROUTE_STATUS_LABELS: Record<string, string> = {
    planned: 'Planificada',
    in_progress: 'En progreso',
    completed: 'Completada',
    cancelled: 'Cancelada',
};

const ROUTE_STATUS_COLORS: Record<string, string> = {
    planned: 'bg-gray-100 text-gray-700 dark:text-gray-200',
    in_progress: 'bg-blue-100 text-blue-700',
    completed: 'bg-green-100 text-green-700',
    cancelled: 'bg-red-100 text-red-700',
};

const STOP_STATUS_LABELS: Record<string, string> = {
    pending: 'Pendiente',
    arrived: 'En sitio',
    delivered: 'Entregado',
    failed: 'Fallido',
    skipped: 'Omitido',
};

const STOP_STATUS_COLORS: Record<string, string> = {
    pending: 'bg-gray-100 text-gray-700 dark:text-gray-200',
    arrived: 'bg-yellow-100 text-yellow-700',
    delivered: 'bg-green-100 text-green-700',
    failed: 'bg-red-100 text-red-700',
    skipped: 'bg-gray-100 text-gray-500 dark:text-gray-400',
};

export default function RouteDetail({ routeId, businessId, onBack, onRefreshList }: RouteDetailProps) {
    const [route, setRoute] = useState<RouteDetailType | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [actionLoading, setActionLoading] = useState<string | null>(null);
    const [failureStopId, setFailureStopId] = useState<number | null>(null);
    const [failureReason, setFailureReason] = useState('');

    const fetchRoute = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const data = await getRouteByIdAction(routeId, businessId);
            setRoute(data);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar la ruta'));
        } finally {
            setLoading(false);
        }
    }, [routeId, businessId]);

    useEffect(() => {
        fetchRoute();
    }, [fetchRoute]);

    const handleStartRoute = async () => {
        if (!confirm('Iniciar esta ruta? El estado cambiara a "En progreso".')) return;
        setActionLoading('start');
        try {
            const updated = await startRouteAction(routeId, businessId);
            setRoute(updated);
            onRefreshList?.();
        } catch (err: any) {
            setError(getActionError(err, 'Error al iniciar la ruta'));
        } finally {
            setActionLoading(null);
        }
    };

    const handleCompleteRoute = async () => {
        if (!confirm('Completar esta ruta? El estado cambiara a "Completada".')) return;
        setActionLoading('complete');
        try {
            const updated = await completeRouteAction(routeId, businessId);
            setRoute(updated);
            onRefreshList?.();
        } catch (err: any) {
            setError(getActionError(err, 'Error al completar la ruta'));
        } finally {
            setActionLoading(null);
        }
    };

    const handleStopDelivered = async (stop: RouteStopInfo) => {
        setActionLoading(`stop-${stop.id}`);
        try {
            await updateStopStatusAction(routeId, stop.id, { status: 'delivered' }, businessId);
            await fetchRoute();
            onRefreshList?.();
        } catch (err: any) {
            setError(getActionError(err, 'Error al actualizar el estado de la parada'));
        } finally {
            setActionLoading(null);
        }
    };

    const handleStopFailed = async (stop: RouteStopInfo) => {
        setFailureStopId(stop.id);
        setFailureReason('');
    };

    const confirmStopFailed = async () => {
        if (failureStopId === null) return;
        setActionLoading(`stop-${failureStopId}`);
        try {
            await updateStopStatusAction(
                routeId,
                failureStopId,
                { status: 'failed', failure_reason: failureReason || undefined },
                businessId
            );
            setFailureStopId(null);
            setFailureReason('');
            await fetchRoute();
            onRefreshList?.();
        } catch (err: any) {
            setError(getActionError(err, 'Error al marcar la parada como fallida'));
        } finally {
            setActionLoading(null);
        }
    };

    const formatDate = (dateStr: string) => {
        const date = new Date(dateStr);
        return date.toLocaleDateString('es-CO', { day: '2-digit', month: '2-digit', year: 'numeric' });
    };

    const formatTime = (dateStr: string | null) => {
        if (!dateStr) return null;
        const date = new Date(dateStr);
        return date.toLocaleTimeString('es-CO', { hour: '2-digit', minute: '2-digit' });
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center p-12">
                <Spinner size="lg" />
            </div>
        );
    }

    if (error && !route) {
        return (
            <div className="space-y-4">
                <Alert type="error">{error}</Alert>
                <Button variant="outline" onClick={onBack}>
                    <ArrowLeftIcon className="w-4 h-4 mr-2" />
                    Volver
                </Button>
            </div>
        );
    }

    if (!route) return null;

    const progressPercent = route.total_stops > 0
        ? Math.round((route.completed_stops / route.total_stops) * 100)
        : 0;

    const sortedStops = [...(route.stops || [])].sort((a, b) => a.sequence - b.sequence);

    return (
        <div className="space-y-6">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            {/* Route header card */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6">
                <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4">
                    <div className="space-y-2">
                        <div className="flex items-center gap-3">
                            <h2 className="text-lg font-semibold text-gray-900 dark:text-white">
                                Ruta del {formatDate(route.date)}
                            </h2>
                            <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${ROUTE_STATUS_COLORS[route.status] || 'bg-gray-100 text-gray-700 dark:text-gray-200'}`}>
                                {ROUTE_STATUS_LABELS[route.status] || route.status}
                            </span>
                        </div>
                        <div className="flex flex-wrap gap-x-6 gap-y-1 text-sm text-gray-600 dark:text-gray-300">
                            {route.driver_name && (
                                <span>Conductor: <strong>{route.driver_name}</strong></span>
                            )}
                            {route.vehicle_plate && (
                                <span>Vehiculo: <strong>{route.vehicle_plate}</strong></span>
                            )}
                            {route.origin_address && (
                                <span>Origen: {route.origin_address}</span>
                            )}
                        </div>
                        {route.notes && (
                            <p className="text-sm text-gray-500 dark:text-gray-400 italic">{route.notes}</p>
                        )}
                    </div>

                    <div className="flex items-center gap-2">
                        {route.status === 'planned' && (
                            <Button
                                variant="primary"
                                onClick={handleStartRoute}
                                disabled={actionLoading === 'start'}
                            >
                                <PlayIcon className="w-4 h-4 mr-1" />
                                {actionLoading === 'start' ? 'Iniciando...' : 'Iniciar Ruta'}
                            </Button>
                        )}
                        {route.status === 'in_progress' && (
                            <Button
                                variant="primary"
                                onClick={handleCompleteRoute}
                                disabled={actionLoading === 'complete'}
                            >
                                <CheckCircleIcon className="w-4 h-4 mr-1" />
                                {actionLoading === 'complete' ? 'Completando...' : 'Completar Ruta'}
                            </Button>
                        )}
                    </div>
                </div>

                {/* Progress bar */}
                <div className="mt-4">
                    <div className="flex items-center justify-between text-sm text-gray-600 dark:text-gray-300 mb-1">
                        <span>Progreso</span>
                        <span>{route.completed_stops}/{route.total_stops} paradas ({progressPercent}%)</span>
                    </div>
                    <div className="w-full bg-gray-200 rounded-full h-2.5">
                        <div
                            className="bg-green-500 h-2.5 rounded-full transition-all"
                            style={{ width: `${progressPercent}%` }}
                        />
                    </div>
                    {route.failed_stops > 0 && (
                        <p className="text-xs text-red-500 mt-1">{route.failed_stops} parada(s) fallida(s)</p>
                    )}
                </div>

                {/* Time info */}
                <div className="flex flex-wrap gap-x-6 gap-y-1 text-xs text-gray-500 dark:text-gray-400 mt-3">
                    {route.actual_start_time && (
                        <span>Inicio real: {formatTime(route.actual_start_time)}</span>
                    )}
                    {route.actual_end_time && (
                        <span>Fin real: {formatTime(route.actual_end_time)}</span>
                    )}
                    {route.total_distance_km != null && (
                        <span>Distancia: {route.total_distance_km} km</span>
                    )}
                    {route.total_duration_min != null && (
                        <span>Duracion: {route.total_duration_min} min</span>
                    )}
                </div>
            </div>

            {/* Stops list */}
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
                <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700">
                    <h3 className="text-base font-semibold text-gray-900 dark:text-white">
                        Paradas ({sortedStops.length})
                    </h3>
                </div>

                {sortedStops.length === 0 ? (
                    <div className="px-6 py-8 text-center text-sm text-gray-400">
                        No hay paradas en esta ruta
                    </div>
                ) : (
                    <div className="divide-y divide-gray-100">
                        {sortedStops.map((stop) => (
                            <div key={stop.id} className="px-6 py-4 flex flex-col sm:flex-row sm:items-center gap-3">
                                {/* Sequence number */}
                                <div className="flex-shrink-0 w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center text-sm font-semibold text-gray-600 dark:text-gray-300">
                                    {stop.sequence}
                                </div>

                                {/* Stop info */}
                                <div className="flex-1 min-w-0">
                                    <div className="flex items-center gap-2 mb-0.5">
                                        <span className="font-medium text-gray-900 dark:text-white text-sm truncate">
                                            {stop.customer_name}
                                        </span>
                                        <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${STOP_STATUS_COLORS[stop.status] || 'bg-gray-100 text-gray-700 dark:text-gray-200'}`}>
                                            {STOP_STATUS_LABELS[stop.status] || stop.status}
                                        </span>
                                    </div>
                                    <p className="text-sm text-gray-500 dark:text-gray-400 truncate">{stop.address}</p>
                                    <div className="flex flex-wrap gap-x-4 gap-y-0.5 text-xs text-gray-400 mt-0.5">
                                        {stop.order_id && (
                                            <span>Orden: {stop.order_id.substring(0, 8)}...</span>
                                        )}
                                        {stop.customer_phone && (
                                            <span>Tel: {stop.customer_phone}</span>
                                        )}
                                        {stop.actual_arrival && (
                                            <span>Llegada: {formatTime(stop.actual_arrival)}</span>
                                        )}
                                        {stop.actual_departure && (
                                            <span>Salida: {formatTime(stop.actual_departure)}</span>
                                        )}
                                        {stop.failure_reason && (
                                            <span className="text-red-400">Motivo: {stop.failure_reason}</span>
                                        )}
                                    </div>
                                </div>

                                {/* Action buttons for in_progress routes with pending stops */}
                                {route.status === 'in_progress' && stop.status === 'pending' && (
                                    <div className="flex gap-2 flex-shrink-0">
                                        <button
                                            onClick={() => handleStopDelivered(stop)}
                                            disabled={actionLoading === `stop-${stop.id}`}
                                            className="inline-flex items-center gap-1 px-3 py-1.5 bg-green-500 hover:bg-green-600 text-white text-xs font-medium rounded-md transition-colors disabled:opacity-50"
                                            title="Marcar como entregado"
                                        >
                                            <CheckIcon className="w-3.5 h-3.5" />
                                            Entregado
                                        </button>
                                        <button
                                            onClick={() => handleStopFailed(stop)}
                                            disabled={actionLoading === `stop-${stop.id}`}
                                            className="inline-flex items-center gap-1 px-3 py-1.5 bg-red-500 hover:bg-red-600 text-white text-xs font-medium rounded-md transition-colors disabled:opacity-50"
                                            title="Marcar como fallido"
                                        >
                                            <XCircleIcon className="w-3.5 h-3.5" />
                                            Fallido
                                        </button>
                                    </div>
                                )}
                            </div>
                        ))}
                    </div>
                )}
            </div>

            {/* Failure reason modal */}
            {failureStopId !== null && (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-xl shadow-xl w-full max-w-md">
                        <div className="px-6 py-4 border-b">
                            <h3 className="text-lg font-semibold text-gray-900 dark:text-white">Motivo del fallo</h3>
                        </div>
                        <div className="p-6 space-y-4">
                            <textarea
                                value={failureReason}
                                onChange={(e) => setFailureReason(e.target.value)}
                                placeholder="Describe el motivo por el que no se pudo entregar (opcional)..."
                                rows={3}
                                className="w-full px-3 py-2 border border-gray-300 rounded-lg text-sm focus:ring-2 focus:ring-blue-500 focus:border-transparent resize-none"
                            />
                            <div className="flex justify-end gap-3">
                                <Button
                                    variant="outline"
                                    onClick={() => {
                                        setFailureStopId(null);
                                        setFailureReason('');
                                    }}
                                >
                                    Cancelar
                                </Button>
                                <Button
                                    variant="primary"
                                    onClick={confirmStopFailed}
                                    disabled={actionLoading !== null}
                                >
                                    {actionLoading ? 'Guardando...' : 'Confirmar fallo'}
                                </Button>
                            </div>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
