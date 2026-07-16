'use client';

import { useState, useEffect, useCallback, useRef } from 'react';
import {
    PlayIcon,
    ArrowPathIcon,
    PencilIcon,
    TrashIcon
} from '@heroicons/react/24/outline';
import { useIntegrations } from '../hooks/useIntegrations';
import { useSSE } from '@/shared/hooks/use-sse';
import {
    Integration,
    SyncOrdersParams
} from '../../domain/types';
import { getOrderStatusMappingsAction } from '@/services/modules/orderstatus/infra/actions';
import { OrderStatusMapping } from '@/services/modules/orderstatus/domain/types';
import { Input, Button, Badge, Spinner, Alert, ConfirmModal, Select, DateRangePicker } from '@/shared/ui';
import { useToast } from '@/shared/providers/toast-provider';
import { TokenStorage } from '@/shared/utils/token-storage';
import { playNotificationSound } from '@/shared/utils';

interface IntegrationListProps {
    onEdit?: (integration: Integration) => void;
    filterCategory?: string;
    businessId?: number | null;
    onCreate?: () => void;
}

// Estado de un lote individual
interface BatchInfo {
    batchIndex: number;
    status: 'pending' | 'processing' | 'completed' | 'failed';
    dateFrom: string;
    dateTo: string;
    duration?: string;
    error?: string;
    completedAt?: Date;
    orderCount: number;
    totalFetched: number | null;
}

// Estado completo de sincronización por lotes
interface BatchSyncState {
    jobId: string;
    totalBatches: number;
    completedBatches: number;
    failedBatches: number;
    dateFrom: string;
    dateTo: string;
    chunkDays: number;
    batches: BatchInfo[];
    currentOrderBatchIndex: number; // Índice del lote que está recibiendo órdenes
    currentFetchBatchIndex: number; // Índice del lote que el provider está fetcheando
}

// Helper para formatear fecha corta
function formatShortDate(dateStr: string): string {
    try {
        const d = new Date(dateStr);
        if (isNaN(d.getTime())) return dateStr;
        return d.toLocaleDateString('es-ES', { day: '2-digit', month: 'short', year: 'numeric' });
    } catch {
        return dateStr;
    }
}

// Estado inicial para los filtros de sincronización
const initialSyncFilters: SyncOrdersParams = {
    created_at_min: '',
    created_at_max: '',
    status: 'any',
    financial_status: 'any',
    fulfillment_status: 'any'
};

const ECOMMERCE_TYPE_IDS = new Set<number>([1, 3, 4, 16, 17, 18, 19, 20, 21]);

function isEcommerceIntegration(integration: { category?: string; integration_type_id?: number }): boolean {
    return integration.category === 'ecommerce' || ECOMMERCE_TYPE_IDS.has(integration.integration_type_id ?? -1);
}

export default function IntegrationList({ onEdit, filterCategory: propFilterCategory, businessId = null, onCreate }: IntegrationListProps) {
    const {
        integrations,
        loading,
        loadingMore,
        error,
        setPage,
        hasMore,
        loadMore,
        total,
        filterCategory,
        setFilterCategory,
        deleteIntegration,
        toggleActive,
        testConnection,
        syncOrders,
        setError
    } = useIntegrations(propFilterCategory || '', businessId);

    const loadMoreRef = useRef<HTMLDivElement | null>(null);

    useEffect(() => {
        const node = loadMoreRef.current;
        if (!node || !hasMore) return;
        const observer = new IntersectionObserver(
            (entries) => {
                if (entries[0].isIntersecting && hasMore && !loading && !loadingMore) {
                    loadMore();
                }
            },
            { rootMargin: '200px' }
        );
        observer.observe(node);
        return () => observer.disconnect();
    }, [hasMore, loading, loadingMore, loadMore]);


    // Sincronizar cambios de categoría después del montaje inicial (cambio de pestaña)
    // No causa doble fetch en el montaje porque el hook ya inicia con propFilterCategory
    useEffect(() => {
        if (propFilterCategory !== undefined && propFilterCategory !== filterCategory) {
            setPage(1);
            setFilterCategory(propFilterCategory);
        }
    }, [propFilterCategory]);

    const [deleteModal, setDeleteModal] = useState<{ show: boolean; id: number | null }>({
        show: false,
        id: null
    });

    // Estado para el modal de sincronización
    const [syncModal, setSyncModal] = useState<{ show: boolean; id: number | null; name: string; integrationTypeId?: number }>({
        show: false,
        id: null,
        name: ''
    });
    const [syncFilters, setSyncFilters] = useState<SyncOrdersParams>(initialSyncFilters);
    const [syncing, setSyncing] = useState(false);
    const [syncError, setSyncError] = useState<string | null>(null);

    // Estados para el progreso de sincronización en tiempo real
    const [syncProgress, setSyncProgress] = useState<{
        total: number;
        created: number;
        rejected: number;
        updated: number;
        skipped: number;
        totalFetched: number | null; // Total de órdenes publicadas a cola (set by integration.sync.completed)
        fetchDuration: string | null; // Duración de la fase de fetch
        orders: Array<{
            orderNumber: string;
            status: 'created' | 'rejected' | 'updated' | 'skipped';
            reason?: string;
            createdAt?: string;
            orderStatus?: string;
            timestamp: Date;
            batchIndex?: number;
        }>;
    } | null>(null);

    // Estado para sincronización por lotes
    const [batchSync, setBatchSync] = useState<BatchSyncState | null>(null);
    const batchSyncRef = useRef<BatchSyncState | null>(null);
    useEffect(() => { batchSyncRef.current = batchSync; }, [batchSync]);

    // Timers para completar lotes: se inician con batch.completed y se reinician con cada orden
    const batchCompletionTimers = useRef<Map<number, NodeJS.Timeout>>(new Map());
    const batchCompletedFlags = useRef<Set<number>>(new Set());
    const syncCompletionTimerRef = useRef<NodeJS.Timeout | null>(null);

    // Completar un lote y avanzar al siguiente
    const completeBatch = useCallback((batchIndex: number) => {
        setBatchSync(prev => {
            if (!prev) return prev;
            const batch = prev.batches[batchIndex];
            if (!batch || batch.status !== 'processing') return prev;
            console.log(`[Lote ${batchIndex + 1}] Completado: ${batch.orderCount} órdenes`);
            const updatedBatch: BatchInfo = { ...batch, status: 'completed', completedAt: new Date() };
            const newCompleted = prev.completedBatches + 1;
            const nextIdx = batchIndex + 1;
            const updatedBatches = prev.batches.map((b, i) => {
                if (b.batchIndex === batchIndex) return updatedBatch;
                if (i === nextIdx && b.status === 'pending') return { ...b, status: 'processing' as const };
                return b;
            });
            const allDone = newCompleted + prev.failedBatches >= prev.totalBatches;
            if (allDone) setSyncing(false);
            const newState = { ...prev, batches: updatedBatches, completedBatches: newCompleted, currentOrderBatchIndex: nextIdx };
            batchSyncRef.current = newState;
            return newState;
        });
    }, []);

    // Iniciar/reiniciar timer de completado para un lote
    const startBatchCompletionTimer = useCallback((batchIndex: number, delayMs: number = 3000) => {
        const existing = batchCompletionTimers.current.get(batchIndex);
        if (existing) clearTimeout(existing);
        const timer = setTimeout(() => {
            console.log(`[Lote ${batchIndex + 1}] Timer: completando tras ${delayMs / 1000}s sin órdenes nuevas`);
            completeBatch(batchIndex);
            batchCompletionTimers.current.delete(batchIndex);
        }, delayMs);
        batchCompletionTimers.current.set(batchIndex, timer);
    }, [completeBatch]);

    // Estados para las opciones de filtros dinámicos desde la BD
    const [orderStatusOptions, setOrderStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [financialStatusOptions, setFinancialStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [fulfillmentStatusOptions, setFulfillmentStatusOptions] = useState<Array<{ value: string; label: string; mappedStatus?: string }>>([]);
    const [loadingMappings, setLoadingMappings] = useState(false);

    // Toast para notificaciones
    const { showToast } = useToast();

    // Obtener businessId del usuario actual
    const [currentBusinessId, setCurrentBusinessId] = useState<number | undefined>(undefined);

    useEffect(() => {
        const businessesData = TokenStorage.getBusinessesData();
        if (businessesData && businessesData.length > 0) {
            // Usar el primer business o el activo si existe
            const activeBusinessId = localStorage.getItem('active_business_id');
            const businessId = activeBusinessId
                ? parseInt(activeBusinessId, 10)
                : businessesData[0].id;
            setCurrentBusinessId(businessId);
        }
    }, []);

    // Hook para escuchar eventos de sincronización en tiempo real via SSE centralizado (modules/events)
    // Solo escuchar cuando el modal está abierto y hay una integración seleccionada
    const integrationEventTypes = syncModal.show ? [
        'integration.sync.order.created',
        'integration.sync.order.updated',
        'integration.sync.order.rejected',
        'integration.sync.started',
        'integration.sync.completed',
        'integration.sync.failed',
        'integration.sync.batched.started',
        'integration.sync.batch.processing',
        'integration.sync.batch.completed',
        'integration.sync.batch.failed'
    ] : [];

    const { isConnected } = useSSE({
        businessId: currentBusinessId,
        integrationId: syncModal.show && syncModal.id ? syncModal.id : undefined,
        eventTypes: integrationEventTypes,
        onMessage: (messageEvent: MessageEvent) => {
            try {
                const event = JSON.parse(messageEvent.data);
                const eventType = event.event_type || event.type || messageEvent.type;

                const eventData = event.data?.data || event.data || {};

                // Extraer integration_id del evento (puede estar a nivel raíz, en data, o en metadata)
                const eventIntegrationId = Number(event.integration_id || eventData.integration_id || event.metadata?.integration_id || 0);

                // Solo procesar si es de la integración que está sincronizando
                if (syncModal.id && eventIntegrationId && eventIntegrationId !== syncModal.id) return;

                // Helper: incrementa orderCount en el lote actual y completa+avanza si llegaron todas las órdenes
                const advanceBatchOrderCount = (increment: number) => {
                    setBatchSync(prev => {
                        if (!prev) return prev;
                        const idx = prev.currentOrderBatchIndex;
                        const batch = prev.batches[idx];
                        if (!batch || batch.status === 'completed' || batch.status === 'failed') return prev;
                        const newOrderCount = batch.orderCount + increment;
                        const updatedBatch = { ...batch, orderCount: newOrderCount };
                        console.log(`[Lote ${idx + 1}] orden recibida: ${newOrderCount}`);

                        const updatedBatches = prev.batches.map((b, i) => i === idx ? updatedBatch : b);
                        const newState = { ...prev, batches: updatedBatches };
                        batchSyncRef.current = newState;

                        // Si batch.completed ya llegó para este lote, reiniciar timer (3s sin órdenes → completar)
                        if (batchCompletedFlags.current.has(idx)) {
                            startBatchCompletionTimer(idx);
                        }

                        return newState;
                    });
                };

                switch (eventType) {
                    case 'integration.sync.order.created': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const createdAt = eventData.created_at || eventData.synced_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;
                        const orderBatchIdx = batchSyncRef.current?.currentOrderBatchIndex;

                        setSyncProgress(prev => {
                            const base = prev || { total: 0, created: 0, rejected: 0, updated: 0, skipped: 0, totalFetched: null, fetchDuration: null, orders: [] };
                            const updated = { ...base, created: base.created + 1, total: base.total + 1, orders: [{ orderNumber, status: 'created' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date(), batchIndex: orderBatchIdx }, ...base.orders] };
                            if (!batchSyncRef.current && updated.totalFetched !== null && updated.created + updated.updated + updated.rejected + updated.skipped >= updated.totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });
                        if (batchSyncRef.current) {
                            advanceBatchOrderCount(1);
                        }
                        break;
                    }
                    case 'integration.sync.order.updated': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const createdAt = eventData.created_at || eventData.updated_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;
                        const orderBatchIdx = batchSyncRef.current?.currentOrderBatchIndex;

                        setSyncProgress(prev => {
                            const base = prev || { total: 0, created: 0, rejected: 0, updated: 0, skipped: 0, totalFetched: null, fetchDuration: null, orders: [] };
                            const updated = { ...base, updated: base.updated + 1, total: base.total + 1, orders: [{ orderNumber, status: 'updated' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date(), batchIndex: orderBatchIdx }, ...base.orders] };
                            if (!batchSyncRef.current && updated.totalFetched !== null && updated.created + updated.updated + updated.rejected + updated.skipped >= updated.totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });
                        if (batchSyncRef.current) {
                            advanceBatchOrderCount(1);
                        }
                        break;
                    }
                    case 'integration.sync.order.rejected': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const reason = eventData.reason || event.metadata?.reason || 'Error desconocido';
                        const error = eventData.error || event.metadata?.error || '';
                        const isSkipped = eventData.skipped === true || event.metadata?.skipped === true;
                        const createdAt = eventData.created_at || eventData.rejected_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;
                        const orderBatchIdx = batchSyncRef.current?.currentOrderBatchIndex;

                        setSyncProgress(prev => {
                            const base = prev || { total: 0, created: 0, rejected: 0, updated: 0, skipped: 0, totalFetched: null, fetchDuration: null, orders: [] };
                            const updated = isSkipped
                                ? { ...base, skipped: base.skipped + 1, total: base.total + 1, orders: [{ orderNumber, status: 'skipped' as const, reason, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date(), batchIndex: orderBatchIdx }, ...base.orders] }
                                : { ...base, rejected: base.rejected + 1, total: base.total + 1, orders: [{ orderNumber, status: 'rejected' as const, reason: reason + (error ? `: ${error}` : ''), createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date(), batchIndex: orderBatchIdx }, ...base.orders] };
                            if (!batchSyncRef.current && updated.totalFetched !== null && updated.created + updated.updated + updated.rejected + updated.skipped >= updated.totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });
                        if (batchSyncRef.current) {
                            advanceBatchOrderCount(1);
                        }
                        break;
                    }
                    case 'integration.sync.started': {
                        // En modo batch, ignorar: el provider emite esto por CADA lote y resetea contadores
                        if (batchSyncRef.current) break;

                        const integrationId = event.integration_id;
                        const integration = integrations.find(i => i.id === integrationId);
                        const integrationName = integration?.name || `Integración ${integrationId}`;

                        setSyncProgress({ total: 0, created: 0, rejected: 0, updated: 0, skipped: 0, totalFetched: null, fetchDuration: null, orders: [] });
                        showToast(`Sincronización iniciada: ${integrationName}`, 'info');
                        break;
                    }
                    case 'integration.sync.completed': {
                        // En modo batch, solo guardar metadata (totalFetched). La completitud la maneja batch.completed + timer.
                        if (batchSyncRef.current) {
                            const batchTotalFetched = Number(eventData.total_fetched) || 0;
                            const fetchIdx = batchSyncRef.current.currentFetchBatchIndex;
                            console.log(`[Lote ${fetchIdx + 1}] sync.completed: totalFetched=${batchTotalFetched}`);
                            setBatchSync(prev => {
                                if (!prev) return prev;
                                const idx = prev.currentFetchBatchIndex;
                                const batch = prev.batches[idx];
                                if (!batch) return prev;
                                const updatedBatches = prev.batches.map((b, i) =>
                                    i === idx ? { ...b, totalFetched: batchTotalFetched } : b
                                );
                                const newState = { ...prev, batches: updatedBatches, currentFetchBatchIndex: idx + 1 };
                                batchSyncRef.current = newState;
                                return newState;
                            });
                            break;
                        }

                        const integrationId = event.integration_id;
                        const integration = integrations.find(i => i.id === integrationId);
                        const integrationName = integration?.name || `Integración ${integrationId}`;
                        const totalFetched = Number(eventData.total_fetched) || 0;
                        const duration = eventData.duration || '';

                        // Si no se obtuvieron órdenes, la sincronización terminó inmediatamente
                        if (totalFetched === 0) {
                            setSyncProgress(prev => ({
                                ...(prev || { total: 0, created: 0, rejected: 0, updated: 0, skipped: 0, orders: [] }),
                                totalFetched: 0,
                                fetchDuration: duration,
                            }));
                            setSyncing(false);
                            playNotificationSound();
                            showToast(`Sincronización completada: ${integrationName} - No se encontraron órdenes`, 'info');
                            break;
                        }

                        // Guardar totalFetched pero NO sobreescribir contadores acumulados
                        // Los contadores se actualizan con eventos individuales (order.created/updated/rejected)
                        setSyncProgress(prev => {
                            const updated = {
                                ...(prev || { total: 0, created: 0, rejected: 0, updated: 0, skipped: 0, orders: [] }),
                                totalFetched,
                                fetchDuration: duration,
                            };
                            // Si ya se procesaron todas las órdenes, marcar como terminado
                            if (updated.created + updated.updated + updated.rejected + updated.skipped >= totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });

                        playNotificationSound();
                        showToast(`Fetch completado: ${integrationName} - ${totalFetched} órdenes obtenidas, procesando...`, 'info');

                        // Fallback: si tras el fetch no llegan todos los eventos por-orden
                        // (evento perdido, timing SSE), no dejar el modal colgado.
                        if (syncCompletionTimerRef.current) clearTimeout(syncCompletionTimerRef.current);
                        syncCompletionTimerRef.current = setTimeout(() => {
                            setSyncing(false);
                        }, 15000);
                        break;
                    }
                    case 'integration.sync.failed': {
                        // En modo batch, avanzar fetchIndex (no habrá sync.completed para este lote)
                        if (batchSyncRef.current) {
                            setBatchSync(prev => {
                                if (!prev) return prev;
                                const newState = { ...prev, currentFetchBatchIndex: prev.currentFetchBatchIndex + 1 };
                                batchSyncRef.current = newState;
                                return newState;
                            });
                            break;
                        }

                        const integrationId = event.integration_id;
                        const integration = integrations.find(i => i.id === integrationId);
                        const integrationName = integration?.name || `Integración ${integrationId}`;
                        const failError = eventData.error || event.metadata?.error || 'Error desconocido';

                        setSyncing(false);
                        setSyncError(failError);
                        playNotificationSound();
                        showToast(`Sincronización fallida: ${integrationName} - ${failError}`, 'error');
                        break;
                    }
                    case 'integration.sync.batched.started': {
                        const jobId = eventData.job_id || '';
                        const totalBatches = Number(eventData.total_batches) || 0;
                        const dateFrom = eventData.date_from || '';
                        const dateTo = eventData.date_to || '';
                        const chunkDays = Number(eventData.chunk_days) || 7;

                        // Pre-calculate batches from date range
                        const batches: BatchInfo[] = [];
                        const startDate = new Date(dateFrom);
                        for (let i = 0; i < totalBatches; i++) {
                            const batchStart = new Date(startDate.getTime() + i * chunkDays * 24 * 60 * 60 * 1000);
                            const batchEnd = new Date(batchStart.getTime() + chunkDays * 24 * 60 * 60 * 1000);
                            // Cap the last batch at dateTo
                            const endCap = new Date(dateTo);
                            batches.push({
                                batchIndex: i,
                                status: i === 0 ? 'processing' : 'pending',
                                dateFrom: batchStart.toISOString(),
                                dateTo: (batchEnd > endCap ? endCap : batchEnd).toISOString(),
                                orderCount: 0,
                                totalFetched: null,
                            });
                        }

                        const newBatchState: BatchSyncState = {
                            jobId,
                            totalBatches,
                            completedBatches: 0,
                            failedBatches: 0,
                            dateFrom,
                            dateTo,
                            chunkDays,
                            batches,
                            currentOrderBatchIndex: 0,
                            currentFetchBatchIndex: 0,
                        };
                        // Set ref synchronously so subsequent events in the same tick see batch mode
                        batchSyncRef.current = newBatchState;
                        setBatchSync(newBatchState);

                        // Initialize order progress
                        setSyncProgress({ total: 0, created: 0, rejected: 0, updated: 0, skipped: 0, totalFetched: null, fetchDuration: null, orders: [] });
                        showToast(`Sincronización por lotes iniciada: ${totalBatches} lotes`, 'info');
                        break;
                    }
                    case 'integration.sync.batch.processing': {
                        // El backend emite esto ANTES de procesar cada lote.
                        // Usamos esto para completar el lote anterior y avanzar currentOrderBatchIndex.
                        const batchIndex = Number(eventData.batch_index);
                        console.log(`[Lote ${batchIndex + 1}] batch.processing del backend`);

                        setBatchSync(prev => {
                            if (!prev) return prev;
                            const prevBatchIdx = batchIndex - 1;
                            let newCompleted = prev.completedBatches;

                            const updatedBatches = prev.batches.map((b, i) => {
                                // Completar el lote anterior si estaba processing
                                if (i === prevBatchIdx && b.status === 'processing') {
                                    newCompleted++;
                                    console.log(`[Lote ${prevBatchIdx + 1}] Completado por batch.processing(${batchIndex}): ${b.orderCount} órdenes`);
                                    return { ...b, status: 'completed' as const, completedAt: new Date() };
                                }
                                // Marcar el lote actual como processing
                                if (i === batchIndex) {
                                    return { ...b, status: 'processing' as const };
                                }
                                return b;
                            });

                            // Limpiar timer del lote anterior si había
                            if (batchCompletionTimers.current.has(prevBatchIdx)) {
                                clearTimeout(batchCompletionTimers.current.get(prevBatchIdx)!);
                                batchCompletionTimers.current.delete(prevBatchIdx);
                            }
                            batchCompletedFlags.current.delete(prevBatchIdx);

                            const allDone = newCompleted + prev.failedBatches >= prev.totalBatches;
                            if (allDone) setSyncing(false);

                            const newState = { ...prev, batches: updatedBatches, completedBatches: newCompleted, currentOrderBatchIndex: batchIndex };
                            batchSyncRef.current = newState;
                            return newState;
                        });
                        break;
                    }
                    case 'integration.sync.batch.completed': {
                        // Para el ÚLTIMO lote, batch.processing no llega después.
                        // Usamos timer de 3s: si no llegan más órdenes, completar.
                        const batchIndex = Number(eventData.batch_index);
                        const duration = eventData.duration || '';
                        const batchDateFrom = eventData.created_at_min || '';
                        const batchDateTo = eventData.created_at_max || '';
                        console.log(`[Lote ${batchIndex + 1}] batch.completed del backend → iniciando timer 3s`);

                        // Guardar metadata (duración, fechas reales)
                        setBatchSync(prev => {
                            if (!prev) return prev;
                            const updatedBatches = prev.batches.map(b => {
                                if (b.batchIndex === batchIndex) {
                                    return { ...b, duration, dateFrom: batchDateFrom || b.dateFrom, dateTo: batchDateTo || b.dateTo };
                                }
                                return b;
                            });
                            return { ...prev, batches: updatedBatches };
                        });

                        // Solo iniciar timer si es el último lote (o si batch.processing no lo completó ya)
                        batchCompletedFlags.current.add(batchIndex);
                        startBatchCompletionTimer(batchIndex);
                        break;
                    }
                    case 'integration.sync.batch.failed': {
                        const batchIndex = Number(eventData.batch_index);
                        const duration = eventData.duration || '';
                        const batchError = eventData.error || 'Error desconocido';
                        const batchDateFrom = eventData.created_at_min || '';
                        const batchDateTo = eventData.created_at_max || '';

                        setBatchSync(prev => {
                            if (!prev) return prev;
                            const updatedBatches = prev.batches.map(b => {
                                if (b.batchIndex === batchIndex) {
                                    return { ...b, status: 'failed' as const, duration, error: batchError, dateFrom: batchDateFrom || b.dateFrom, dateTo: batchDateTo || b.dateTo };
                                }
                                return b;
                            });
                            // Avanzar al siguiente lote
                            const nextIdx = prev.currentOrderBatchIndex + 1;
                            if (nextIdx < prev.totalBatches) {
                                const nextBatch = updatedBatches.find(b => b.batchIndex === nextIdx);
                                if (nextBatch && nextBatch.status === 'pending') {
                                    updatedBatches[updatedBatches.indexOf(nextBatch)] = { ...nextBatch, status: 'processing' as const };
                                }
                            }
                            const newFailed = prev.failedBatches + 1;
                            const allDone = prev.completedBatches + newFailed >= prev.totalBatches;
                            if (allDone) setSyncing(false);
                            const newState = { ...prev, batches: updatedBatches, failedBatches: newFailed, currentOrderBatchIndex: nextIdx };
                            batchSyncRef.current = newState;
                            return newState;
                        });
                        playNotificationSound();
                        showToast(`Lote ${batchIndex + 1} falló: ${batchError}`, 'error');
                        break;
                    }
                }
            } catch (err) {
                console.error('Error parsing integration SSE event:', err);
            }
        },
        onError: () => {
            console.error('Error en conexión SSE de eventos de integraciones');
        },
        onOpen: () => {
            console.log('Conexión SSE de eventos de integraciones establecida');
        },
    });

    // Safety timeout (fallback): normalmente created+updated+rejected === totalFetched
    // gracias a los eventos order.rejected del backend. Este timeout solo aplica si se pierde algún evento SSE.
    // En modo batch, la finalización la maneja batch.completed/batch.failed exclusivamente.
    useEffect(() => {
        if (!syncing || !syncProgress?.totalFetched) return;
        if (batchSyncRef.current) return;
        const processed = syncProgress.created + syncProgress.updated + syncProgress.rejected;
        if (processed >= syncProgress.totalFetched) return; // Already done

        const timer = setTimeout(() => {
            setSyncing(false);
        }, 10000);

        return () => clearTimeout(timer);
    }, [syncing, syncProgress?.totalFetched, syncProgress?.created, syncProgress?.updated, syncProgress?.rejected]);

    const handleDeleteClick = (id: number) => {
        setDeleteModal({ show: true, id });
    };

    const handleDeleteConfirm = async () => {
        if (deleteModal.id) {
            const success = await deleteIntegration(deleteModal.id);
            if (success) {
                setDeleteModal({ show: false, id: null });
            }
        }
    };

    const handleTest = async (id: number) => {
        const result = await testConnection(id);
        if (result.success) {
            alert('✅ Conexión exitosa');
        } else {
            alert(`❌ Error: ${result.message}`);
        }
    };

    // Estados de Shopify por categoría (para agrupar los mapeos)
    const SHOPIFY_ORDER_STATUSES = ['any', 'open', 'closed', 'cancelled'];
    const SHOPIFY_FINANCIAL_STATUSES = ['any', 'authorized', 'pending', 'paid', 'partially_paid', 'refunded', 'voided', 'partially_refunded', 'unpaid'];
    const SHOPIFY_FULFILLMENT_STATUSES = ['any', 'shipped', 'partial', 'unfulfilled', 'unshipped'];

    // Cargar mapeos de estados para la integración
    const loadStatusMappings = async (integrationTypeId: number) => {
        setLoadingMappings(true);
        try {
            const response = await getOrderStatusMappingsAction({
                integration_type_id: integrationTypeId,
                is_active: true
            });

            // El backend devuelve { data: [...], total: number }
            // Puede venir con o sin el campo success dependiendo del endpoint
            const mappings: OrderStatusMapping[] = (response as any).data || response.data || [];

            console.log('Status mappings response:', response);
            console.log('Status mappings loaded:', mappings);

            if (mappings && mappings.length > 0) {
                // Agrupar por tipo de estado
                const orderStatusMap = mappings
                    .filter(m => SHOPIFY_ORDER_STATUSES.includes(m.original_status))
                    .map(m => ({
                        value: m.original_status,
                        label: `${m.original_status}${m.order_status ? ` → ${m.order_status.name}` : ''}`,
                        mappedStatus: m.order_status?.code
                    }));

                const financialStatusMap = mappings
                    .filter(m => SHOPIFY_FINANCIAL_STATUSES.includes(m.original_status))
                    .map(m => ({
                        value: m.original_status,
                        label: `${m.original_status}${m.order_status ? ` → ${m.order_status.name}` : ''}`,
                        mappedStatus: m.order_status?.code
                    }));

                const fulfillmentStatusMap = mappings
                    .filter(m => SHOPIFY_FULFILLMENT_STATUSES.includes(m.original_status))
                    .map(m => ({
                        value: m.original_status,
                        label: `${m.original_status}${m.order_status ? ` → ${m.order_status.name}` : ''}`,
                        mappedStatus: m.order_status?.code
                    }));

                console.log('Order status options:', orderStatusMap);
                console.log('Financial status options:', financialStatusMap);
                console.log('Fulfillment status options:', fulfillmentStatusMap);

                // Agregar opción "Todos" solo si no existe ya
                const hasAnyInOrder = orderStatusMap.some(opt => opt.value === 'any');
                const hasAnyInFinancial = financialStatusMap.some(opt => opt.value === 'any');
                const hasAnyInFulfillment = fulfillmentStatusMap.some(opt => opt.value === 'any');

                setOrderStatusOptions(
                    orderStatusMap.length > 0
                        ? (hasAnyInOrder ? orderStatusMap : [{ value: 'any', label: 'Todos' }, ...orderStatusMap])
                        : [{ value: 'any', label: 'Todos' }]
                );
                setFinancialStatusOptions(
                    financialStatusMap.length > 0
                        ? (hasAnyInFinancial ? financialStatusMap : [{ value: 'any', label: 'Todos' }, ...financialStatusMap])
                        : [{ value: 'any', label: 'Todos' }]
                );
                setFulfillmentStatusOptions(
                    fulfillmentStatusMap.length > 0
                        ? (hasAnyInFulfillment ? fulfillmentStatusMap : [{ value: 'any', label: 'Todos' }, ...fulfillmentStatusMap])
                        : [{ value: 'any', label: 'Todos' }]
                );
            } else {
                // Si no hay mapeos, usar opciones por defecto
                console.warn('No mappings found, using default options');
                setOrderStatusOptions([{ value: 'any', label: 'Todos' }]);
                setFinancialStatusOptions([{ value: 'any', label: 'Todos' }]);
                setFulfillmentStatusOptions([{ value: 'any', label: 'Todos' }]);
            }
        } catch (err) {
            console.error('Error loading status mappings:', err);
            // En caso de error, usar opciones por defecto
            setOrderStatusOptions([{ value: 'any', label: 'Todos' }]);
            setFinancialStatusOptions([{ value: 'any', label: 'Todos' }]);
            setFulfillmentStatusOptions([{ value: 'any', label: 'Todos' }]);
        } finally {
            setLoadingMappings(false);
        }
    };

    // Abrir modal de sincronización
    const handleSyncClick = async (id: number, name: string) => {
        // Buscar la integración para obtener el integration_type_id
        const integration = integrations.find(i => i.id === id);
        if (!integration) {
            alert('Error: No se encontró la integración');
            return;
        }

        const integrationTypeId = integration.integration_type_id;

        setSyncModal({ show: true, id, name, integrationTypeId });

        // Cargar mapeos de estados
        await loadStatusMappings(integrationTypeId);

        // Consultar si hay una sincronización en curso
        try {
            const { getSyncStatusAction } = await import('../../infra/actions');
            const token = TokenStorage.getSessionToken();
            const syncStatus = await getSyncStatusAction(id, currentBusinessId, token);
            if (syncStatus.success && syncStatus.in_progress && syncStatus.sync_state) {
                // Hay una sincronización en curso, mostrar el estado actual
                setSyncing(true);
                setSyncProgress({
                    total: 0,
                    created: 0,
                    rejected: 0,
                    updated: 0,
                    skipped: 0,
                    totalFetched: null,
                    fetchDuration: null,
                    orders: []
                });
                showToast('🔄 Hay una sincronización en curso. Mostrando progreso actual...', 'info');
            }
        } catch (error: any) {
            console.error('Error al consultar estado de sincronización:', error);
            // Continuar normalmente si hay error
        }

        // Establecer fecha mínima por defecto a 30 días atrás
        const thirtyDaysAgo = new Date();
        thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
        setSyncFilters({
            ...initialSyncFilters,
            created_at_min: thirtyDaysAgo.toISOString().split('T')[0]
        });
    };

    // Ejecutar sincronización con filtros
    const handleSyncConfirm = async () => {
        if (syncModal.id) {
            setSyncing(true);
            setSyncError(null);
            // Inicializar progreso
            setSyncProgress({
                total: 0,
                created: 0,
                rejected: 0,
                updated: 0,
                skipped: 0,
                totalFetched: null,
                fetchDuration: null,
                orders: []
            });
            try {
                // Preparar parámetros (solo enviar los que tienen valor)
                const params: SyncOrdersParams = {};
                if (syncFilters.created_at_min) params.created_at_min = syncFilters.created_at_min;
                if (syncFilters.created_at_max) params.created_at_max = syncFilters.created_at_max;
                if (syncFilters.status && syncFilters.status !== 'any') params.status = syncFilters.status;
                if (syncFilters.financial_status && syncFilters.financial_status !== 'any') params.financial_status = syncFilters.financial_status;
                if (syncFilters.fulfillment_status && syncFilters.fulfillment_status !== 'any') params.fulfillment_status = syncFilters.fulfillment_status;

                const result = await syncOrders(syncModal.id, Object.keys(params).length > 0 ? params : undefined);
                if (result.success) {
                    // No mostrar alert, el evento SSE mostrará la notificación en tiempo real
                    showToast('🔄 Sincronización iniciada. Recibirás notificaciones en tiempo real.', 'info');
                } else {
                    showToast(`❌ Error al iniciar sincronización: ${result.message}`, 'error');
                    setSyncProgress(null);
                    setSyncing(false);
                }
            } catch (error: any) {
                showToast(`❌ Error al iniciar sincronización: ${error.message}`, 'error');
                setSyncProgress(null);
                setSyncing(false);
            }
        }
    };

    // Cerrar modal de sincronización
    const handleSyncCancel = () => {
        setSyncModal({ show: false, id: null, name: '' });
        setSyncFilters(initialSyncFilters);
        setSyncProgress(null);
        setBatchSync(null);
        setSyncing(false);
        setSyncError(null);
        // Limpiar timers y flags de batch
        batchCompletionTimers.current.forEach(t => clearTimeout(t));
        batchCompletionTimers.current.clear();
        batchCompletedFlags.current.clear();
        // Limpiar opciones
        setOrderStatusOptions([]);
        setFinancialStatusOptions([]);
        setFulfillmentStatusOptions([]);
    };

    if (loading) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    const renderCard = (integration: Integration) => {
        const initial = integration.name.charAt(0).toUpperCase();
        const logo = integration.integration_type?.image_url ? (
            <img
                src={integration.integration_type.image_url}
                alt={integration.integration_type.name || integration.name}
                className="w-11 h-11 object-contain border border-gray-200 dark:border-gray-700 rounded-lg p-1 bg-white dark:bg-gray-800 flex-shrink-0"
                onError={(e) => {
                    (e.target as HTMLImageElement).style.display = 'none';
                    const parent = (e.target as HTMLImageElement).parentElement;
                    if (parent) {
                        parent.innerHTML = '<div class="w-11 h-11 flex items-center justify-center bg-gray-100 dark:bg-gray-700 rounded-lg text-gray-400 text-sm font-semibold flex-shrink-0">' + initial + '</div>';
                    }
                }}
            />
        ) : (
            <div className="w-11 h-11 flex items-center justify-center bg-gray-100 dark:bg-gray-700 rounded-lg text-gray-400 text-sm font-semibold flex-shrink-0">
                {initial}
            </div>
        );

        const subtitleParts = [integration.integration_type?.name, integration.business_name].filter(Boolean);
        const categoryLabel = integration.category_name || integration.category;
        const categoryColor = integration.category_color || '#6B7280';

        return (
            <div
                key={integration.id}
                className="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm hover:shadow-md transition-shadow p-4 flex flex-col gap-3"
            >
                <div className="flex items-start justify-between gap-3">
                    <div className="flex items-center gap-3 min-w-0">
                        {logo}
                        <div className="min-w-0">
                            <div className="flex items-center gap-2">
                                <span className="text-sm font-semibold text-gray-900 dark:text-white truncate">{integration.name}</span>
                                <span className="text-[10px] text-gray-400 flex-shrink-0">#{integration.id}</span>
                            </div>
                            <div className="text-xs text-gray-500 dark:text-gray-400 truncate">
                                {subtitleParts.length > 0 ? subtitleParts.join(' · ') : 'Sin tipo'}
                            </div>
                            {categoryLabel && (
                                <span
                                    className="inline-flex items-center mt-1.5 px-2 py-0.5 rounded-full text-[10px] font-medium text-white whitespace-nowrap"
                                    style={{ backgroundColor: categoryColor }}
                                >
                                    {categoryLabel}
                                </span>
                            )}
                        </div>
                    </div>
                    <div className="flex flex-col items-end gap-1 flex-shrink-0">
                        <Badge type={integration.is_active ? 'success' : 'error'}>
                            {integration.is_active ? 'Activo' : 'Inactivo'}
                        </Badge>
                        {integration.is_default && (
                            <Badge type="primary">Por defecto</Badge>
                        )}
                    </div>
                </div>

                <div className="flex flex-wrap items-center justify-between gap-2 border-t border-gray-100 dark:border-gray-700 pt-3">
                    <div className="flex items-center gap-2">
                        <button
                            onClick={() => handleTest(integration.id)}
                            className="inline-flex items-center gap-1.5 px-2.5 py-1.5 text-xs font-medium text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                            title="Probar conexion"
                            aria-label="Probar conexion"
                        >
                            <PlayIcon className="w-3.5 h-3.5" />
                            Probar
                        </button>
                        {isEcommerceIntegration(integration) && (
                            <button
                                onClick={() => handleSyncClick(integration.id, integration.name)}
                                className="inline-flex items-center gap-1.5 px-2.5 py-1.5 text-xs font-medium text-gray-700 dark:text-gray-200 border border-gray-200 dark:border-gray-600 rounded-md hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                                title="Sincronizar ordenes"
                                aria-label="Sincronizar ordenes"
                            >
                                <ArrowPathIcon className="w-3.5 h-3.5" />
                                Sincronizar
                            </button>
                        )}
                    </div>
                    <div className="flex items-center gap-1.5">
                        {onEdit && (
                            <button
                                onClick={() => onEdit(integration)}
                                className="p-1.5 text-gray-500 hover:text-yellow-600 hover:bg-yellow-50 dark:hover:bg-gray-700 rounded-md transition-colors"
                                title="Editar integracion"
                                aria-label="Editar integracion"
                            >
                                <PencilIcon className="w-4 h-4" />
                            </button>
                        )}
                        <button
                            onClick={() => handleDeleteClick(integration.id)}
                            className="p-1.5 text-gray-500 hover:text-red-600 hover:bg-red-50 dark:hover:bg-gray-700 rounded-md transition-colors"
                            title="Eliminar integracion"
                            aria-label="Eliminar integracion"
                        >
                            <TrashIcon className="w-4 h-4" />
                        </button>
                        <button
                            onClick={() => toggleActive(integration.id, integration.is_active)}
                            className={`relative inline-flex h-6 w-11 items-center rounded-full transition-colors focus:outline-none ${integration.is_active ? 'bg-orange-500' : 'bg-gray-300 dark:bg-gray-600'}`}
                            title={integration.is_active ? 'Desactivar integracion' : 'Activar integracion'}
                            aria-label={integration.is_active ? 'Desactivar integracion' : 'Activar integracion'}
                        >
                            <span className={`inline-block h-4 w-4 transform rounded-full bg-white shadow transition-transform ${integration.is_active ? 'translate-x-6' : 'translate-x-1'}`} />
                        </button>
                    </div>
                </div>
            </div>
        );
    };

    return (
        <div className="mx-auto w-full max-w-7xl">
            <div className="rounded-2xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm overflow-hidden">
                <div className="px-6 pt-6 pb-4 border-b border-gray-100 dark:border-gray-700">
                    <h2 className="text-lg font-bold text-gray-900 dark:text-white">Integraciones conectadas</h2>
                    <p className="text-sm text-gray-500 dark:text-gray-400">Las integraciones conectadas a tu cuenta.</p>
                    <span className="inline-flex items-center gap-1.5 mt-3 px-2.5 py-1 rounded-full bg-gray-100 dark:bg-gray-700 text-xs font-medium text-gray-600 dark:text-gray-300">
                        <span className="w-1.5 h-1.5 rounded-full bg-green-500" />
                        {integrations.filter(i => i.is_active).length} activas · {total} en total
                    </span>
                </div>

                <div className="p-6 space-y-4">
                    {error && (
                        <Alert type="error" onClose={() => setError(null)}>
                            {error}
                        </Alert>
                    )}

                    {integrations.length === 0 ? (
                        <div className="rounded-lg border border-dashed border-gray-200 dark:border-gray-700 p-12 text-center text-sm text-gray-500 dark:text-gray-400">
                            No hay integraciones disponibles
                        </div>
                    ) : (
                        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
                            {integrations.map(renderCard)}
                        </div>
                    )}

                    {integrations.length > 0 && (
                        <div ref={loadMoreRef} className="flex justify-center items-center py-2">
                            {loadingMore ? (
                                <div className="flex items-center gap-2 text-sm text-gray-500 dark:text-gray-400">
                                    <Spinner size="sm" />
                                    Cargando mas...
                                </div>
                            ) : hasMore ? (
                                <button
                                    onClick={loadMore}
                                    className="text-sm text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 transition-colors"
                                >
                                    Cargar mas
                                </button>
                            ) : (
                                <span className="text-xs text-gray-400">No hay mas integraciones</span>
                            )}
                        </div>
                    )}
                </div>

                {onCreate && (
                    <div className="border-t border-gray-100 dark:border-gray-700 p-4">
                        <Button variant="primary" onClick={onCreate} className="w-full justify-center gap-2">
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                            Crear Integración
                        </Button>
                    </div>
                )}
            </div>

            <ConfirmModal
                isOpen={deleteModal.show}
                onClose={() => setDeleteModal({ show: false, id: null })}
                onConfirm={handleDeleteConfirm}
                title="Eliminar Integración"
                message="¿Estás seguro de que deseas eliminar esta integración? Esta acción no se puede deshacer."
            />

            {/* Modal de Sincronización con Filtros */}
            {syncModal.show && (() => {
                const hasProgress = syncProgress && (syncProgress.total > 0 || syncProgress.totalFetched !== null);
                const isShowingProgress = syncing || batchSync || syncError || hasProgress;
                const isBatchMode = !!batchSync;
                const batchFinished = batchSync && (batchSync.completedBatches + batchSync.failedBatches >= batchSync.totalBatches);
                const directFinished = !isBatchMode && !syncing && (syncError || hasProgress);
                const allFinished = batchFinished || directFinished;

                return (
                <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
                    <div className={`bg-white dark:bg-gray-800 rounded-lg shadow-xl flex flex-col max-h-[90vh] transition-all duration-300 ${isShowingProgress ? 'w-full max-w-5xl overflow-hidden' : 'w-full max-w-3xl overflow-visible'}`}>
                        {/* Header */}
                        <div className="px-6 py-4 border-b border-gray-200 dark:border-gray-700 flex-shrink-0">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                    <div>
                                        <h3 className="text-lg font-semibold text-gray-900 dark:text-white">
                                            Sincronizar Órdenes
                                        </h3>
                                        <p className="text-sm text-gray-500 dark:text-gray-400 dark:text-gray-400 mt-0.5">
                                            {syncModal.name}
                                        </p>
                                    </div>
                                    {isShowingProgress && (
                                        <div className="flex items-center gap-1.5 ml-2">
                                            <div className={`w-2 h-2 rounded-full ${isConnected ? 'bg-green-500' : 'bg-red-500'} ${isConnected && syncing ? 'animate-pulse' : ''}`} />
                                            <span className="text-xs text-gray-400">{isConnected ? 'SSE conectado' : 'SSE desconectado'}</span>
                                        </div>
                                    )}
                                </div>
                                {!syncing && (
                                    <button
                                        onClick={handleSyncCancel}
                                        className="text-gray-400 hover:text-gray-600 dark:text-gray-300 dark:text-gray-300 transition-colors"
                                        aria-label="Cerrar modal"
                                    >
                                        <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                        </svg>
                                    </button>
                                )}
                            </div>
                        </div>

                        {/* Content */}
                        <div className={`flex-1 min-h-0 ${isShowingProgress ? 'overflow-y-auto' : 'overflow-visible'}`}>
                            {/* Phase 1: Form (before sync starts) */}
                            {!isShowingProgress && (
                                <div className="px-6 py-4 space-y-4">
                                    <div className="space-y-4">
                                        <div>
                                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-200 dark:text-gray-200 mb-1">Rango de Fechas</label>
                                            <DateRangePicker
                                                startDate={syncFilters.created_at_min}
                                                endDate={syncFilters.created_at_max}
                                                onChange={(startDate, endDate) => {
                                                    setSyncFilters(prev => ({ ...prev, created_at_min: startDate || '', created_at_max: endDate || '' }));
                                                }}
                                                placeholder="Seleccionar rango de fechas (opcional)"
                                            />
                                        </div>
                                    </div>
                                    <div className="bg-blue-50 border border-blue-200 rounded-lg p-3">
                                        <p className="text-sm text-blue-700">
                                            Los estados mostrados están mapeados a estados de Probability. Si no especificas filtros, se sincronizarán las órdenes de los últimos 30 días. Los rangos mayores a 14 días se procesarán por lotes automáticamente.
                                        </p>
                                    </div>
                                </div>
                            )}

                            {/* Phase 2: Progress */}
                            {isShowingProgress && (
                                <div className="px-6 py-4">
                                    {isBatchMode ? (
                                        /* Batch mode: two-column layout */
                                        <div className="grid grid-cols-5 gap-6">
                                            {/* Left column: Batches (3/5) */}
                                            <div className="col-span-3 space-y-4">
                                                {/* Batch progress bar */}
                                                <div>
                                                    <div className="flex justify-between items-center mb-2">
                                                        <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-200 dark:text-gray-200">Progreso por Lotes</h4>
                                                        <span className="text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400">
                                                            {batchSync.completedBatches + batchSync.failedBatches} / {batchSync.totalBatches}
                                                            {' '}({batchSync.totalBatches > 0 ? Math.round(((batchSync.completedBatches + batchSync.failedBatches) / batchSync.totalBatches) * 100) : 0}%)
                                                        </span>
                                                    </div>
                                                    {/* Segmented progress bar */}
                                                    <div className="w-full bg-gray-200 rounded-full h-5 overflow-hidden flex">
                                                        {batchSync.batches.map((batch) => {
                                                            const widthPct = 100 / batchSync.totalBatches;
                                                            let bgClass = 'bg-gray-300';
                                                            if (batch.status === 'completed') bgClass = 'bg-green-500';
                                                            else if (batch.status === 'failed') bgClass = 'bg-red-500';
                                                            else if (batch.status === 'processing') bgClass = 'bg-blue-500 animate-pulse';
                                                            return (
                                                                <div
                                                                    key={batch.batchIndex}
                                                                    className={`h-full ${bgClass} transition-all duration-500 ease-out border-r border-white/30 last:border-r-0`}
                                                                    style={{ width: `${widthPct}%` }}
                                                                    title={`Lote ${batch.batchIndex + 1}: ${batch.status}`}
                                                                />
                                                            );
                                                        })}
                                                    </div>
                                                    {/* Summary badges */}
                                                    <div className="flex gap-3 mt-2 text-xs">
                                                        <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-green-100 text-green-700 rounded-full">
                                                            <div className="w-2 h-2 rounded-full bg-green-500" /> {batchSync.completedBatches} completados
                                                        </span>
                                                        {batchSync.failedBatches > 0 && (
                                                            <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-red-100 text-red-700 rounded-full">
                                                                <div className="w-2 h-2 rounded-full bg-red-500" /> {batchSync.failedBatches} fallidos
                                                            </span>
                                                        )}
                                                        <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 dark:text-gray-300 rounded-full">
                                                            <div className="w-2 h-2 rounded-full bg-gray-400" /> {batchSync.totalBatches - batchSync.completedBatches - batchSync.failedBatches} pendientes
                                                        </span>
                                                    </div>
                                                </div>

                                                {/* Batch history */}
                                                <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                                                    <div className="p-2 bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-10">
                                                        <p className="text-xs font-semibold text-gray-700 dark:text-gray-200 dark:text-gray-200">Historial de Lotes</p>
                                                    </div>
                                                    <div className="max-h-[50vh] overflow-y-auto divide-y divide-gray-100">
                                                        {batchSync.batches.map((batch) => (
                                                            <div
                                                                key={batch.batchIndex}
                                                                className={`p-3 text-xs transition-colors ${
                                                                    batch.status === 'processing' ? 'bg-blue-50 border-l-2 border-l-blue-500' :
                                                                    batch.status === 'completed' ? 'bg-white dark:bg-gray-800 hover:bg-gray-50 dark:bg-gray-700' :
                                                                    batch.status === 'failed' ? 'bg-red-50 hover:bg-red-100' :
                                                                    'bg-gray-50 dark:bg-gray-700'
                                                                }`}
                                                            >
                                                                <div className="flex items-center justify-between">
                                                                    <div className="flex items-center gap-2">
                                                                        {/* Status icon */}
                                                                        {batch.status === 'completed' && (
                                                                            <svg className="w-4 h-4 text-green-500" fill="currentColor" viewBox="0 0 20 20">
                                                                                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
                                                                            </svg>
                                                                        )}
                                                                        {batch.status === 'failed' && (
                                                                            <svg className="w-4 h-4 text-red-500" fill="currentColor" viewBox="0 0 20 20">
                                                                                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
                                                                            </svg>
                                                                        )}
                                                                        {batch.status === 'processing' && (
                                                                            <Spinner size="sm" />
                                                                        )}
                                                                        {batch.status === 'pending' && (
                                                                            <div className="w-4 h-4 rounded-full border-2 border-gray-300 dark:border-gray-600" />
                                                                        )}
                                                                        <span className="font-medium text-gray-800 dark:text-gray-100">Lote {batch.batchIndex + 1}</span>
                                                                    </div>
                                                                    <div className="flex items-center gap-2">
                                                                        {batch.duration && (
                                                                            <span className="text-gray-400">{batch.duration}</span>
                                                                        )}
                                                                        <span className={`px-1.5 py-0.5 rounded text-[10px] font-medium ${
                                                                            batch.status === 'completed' ? 'bg-green-100 text-green-700' :
                                                                            batch.status === 'failed' ? 'bg-red-100 text-red-700' :
                                                                            batch.status === 'processing' ? 'bg-blue-100 text-blue-700' :
                                                                            'bg-gray-100 dark:bg-gray-700 text-gray-500 dark:text-gray-400 dark:text-gray-400'
                                                                        }`}>
                                                                            {batch.status === 'completed' ? 'Completado' :
                                                                             batch.status === 'failed' ? 'Fallido' :
                                                                             batch.status === 'processing' ? 'Procesando' :
                                                                             'Pendiente'}
                                                                        </span>
                                                                    </div>
                                                                </div>
                                                                <div className="mt-1 ml-6 text-gray-500 dark:text-gray-400 dark:text-gray-400">
                                                                    {formatShortDate(batch.dateFrom)} → {formatShortDate(batch.dateTo)}
                                                                    {batch.orderCount > 0 && <span> · {batch.orderCount} {batch.orderCount === 1 ? 'orden' : 'órdenes'}</span>}
                                                                    {batch.orderCount === 0 && batch.status === 'completed' && <span> · Sin órdenes</span>}
                                                                </div>
                                                                {/* Animated bar for processing batch */}
                                                                {batch.status === 'processing' && (
                                                                    <div className="mt-1.5 ml-6">
                                                                        <div className="w-full bg-gray-200 rounded-full h-1.5 overflow-hidden">
                                                                            <div className="h-full bg-blue-500 animate-pulse" style={{ width: '100%' }} />
                                                                        </div>
                                                                    </div>
                                                                )}
                                                                {batch.status === 'failed' && batch.error && (
                                                                    <div className="mt-1 ml-6 text-red-600 text-[11px]">
                                                                        {batch.error}
                                                                    </div>
                                                                )}
                                                            </div>
                                                        ))}
                                                    </div>
                                                </div>
                                            </div>

                                            {/* Right column: Orders (2/5) */}
                                            <div className="col-span-2 space-y-3">
                                                {/* Order stats mini */}
                                                {syncProgress && (
                                                    <div className="flex gap-3 text-xs flex-wrap">
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-2.5 h-2.5 rounded-full bg-green-500" />
                                                            <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Creadas: <strong>{syncProgress.created}</strong></span>
                                                        </div>
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-2.5 h-2.5 rounded-full bg-blue-500" />
                                                            <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Actualizadas: <strong>{syncProgress.updated}</strong></span>
                                                        </div>
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-2.5 h-2.5 rounded-full bg-red-500" />
                                                            <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Rechazadas: <strong>{syncProgress.rejected}</strong></span>
                                                        </div>
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-2.5 h-2.5 rounded-full bg-amber-500" />
                                                            <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Omitidas: <strong>{syncProgress.skipped}</strong></span>
                                                        </div>
                                                    </div>
                                                )}
                                                {/* Orders feed */}
                                                <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                                                    <div className="p-2 bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-700 sticky top-0 z-10">
                                                        <p className="text-xs font-semibold text-gray-700 dark:text-gray-200 dark:text-gray-200">Órdenes Procesadas</p>
                                                    </div>
                                                    <div className="max-h-[60vh] overflow-y-auto divide-y divide-gray-100">
                                                        {syncProgress && syncProgress.orders.length > 0 ? (
                                                            syncProgress.orders.slice(0, 100).map((order, index) => (
                                                                <div
                                                                    key={index}
                                                                    className={`px-3 py-2 text-xs ${
                                                                        order.status === 'created' ? 'bg-green-50' :
                                                                        order.status === 'updated' ? 'bg-blue-50' :
                                                                        order.status === 'skipped' ? 'bg-amber-50' :
                                                                        'bg-red-50'
                                                                    }`}
                                                                >
                                                                    <div className="flex items-center justify-between gap-1">
                                                                        <div className="flex items-center gap-1.5 min-w-0">
                                                                            <div className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
                                                                                order.status === 'created' ? 'bg-green-500' :
                                                                                order.status === 'updated' ? 'bg-blue-500' :
                                                                                order.status === 'skipped' ? 'bg-amber-500' :
                                                                                'bg-red-500'
                                                                            }`} />
                                                                            {order.status === 'skipped' && (
                                                                                <span className="px-1.5 py-0.5 bg-amber-100 text-amber-700 rounded text-[10px] flex-shrink-0 font-medium">
                                                                                    Omitida
                                                                                </span>
                                                                            )}
                                                                            <span className="font-medium text-gray-800 dark:text-gray-100 truncate">#{order.orderNumber}</span>
                                                                            {order.batchIndex !== undefined && (
                                                                                <span className="px-1 py-0.5 bg-indigo-100 text-indigo-600 rounded text-[10px] flex-shrink-0 font-medium">
                                                                                    L{order.batchIndex + 1}
                                                                                </span>
                                                                            )}
                                                                            {order.orderStatus && (
                                                                                <span className="px-1.5 py-0.5 bg-gray-200 text-gray-600 dark:text-gray-300 dark:text-gray-300 rounded text-[10px] flex-shrink-0">
                                                                                    {order.orderStatus}
                                                                                </span>
                                                                            )}
                                                                        </div>
                                                                        <span className="text-gray-400 text-[10px] flex-shrink-0">
                                                                            {order.timestamp.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
                                                                        </span>
                                                                    </div>
                                                                    {(order.status === 'rejected' || order.status === 'skipped') && order.reason && (
                                                                        <div className={`text-[10px] mt-0.5 ml-3 truncate ${order.status === 'skipped' ? 'text-amber-600' : 'text-red-600'}`} title={order.reason}>
                                                                            {order.reason}
                                                                        </div>
                                                                    )}
                                                                </div>
                                                            ))
                                                        ) : (
                                                            <div className="p-4 text-center text-xs text-gray-400">
                                                                Esperando órdenes...
                                                            </div>
                                                        )}
                                                        {syncProgress && syncProgress.orders.length > 100 && (
                                                            <div className="p-2 bg-gray-50 dark:bg-gray-700 text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400 text-center">
                                                                Y {syncProgress.orders.length - 100} órdenes más...
                                                            </div>
                                                        )}
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    ) : (
                                        /* Direct mode (no batches): single column with orders */
                                        <div className="space-y-4">
                                            {syncError && (
                                                <div className="bg-red-50 border border-red-200 rounded-lg p-3">
                                                    <p className="text-sm text-red-700 font-medium">Sincronización fallida</p>
                                                    <p className="text-sm text-red-600 mt-1">{syncError}</p>
                                                </div>
                                            )}
                                            {/* Empty result message */}
                                            {syncProgress && syncProgress.totalFetched === 0 && !syncing && (
                                                <div className="flex flex-col items-center justify-center py-12 text-gray-500 dark:text-gray-400 dark:text-gray-400">
                                                    <svg className="w-16 h-16 mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                                                    </svg>
                                                    <p className="text-lg font-medium">No se encontraron órdenes</p>
                                                    <p className="text-sm mt-1">No hay órdenes en el rango de fechas seleccionado</p>
                                                    {syncProgress.fetchDuration && (
                                                        <p className="text-xs mt-2 text-gray-400">Consulta completada en {syncProgress.fetchDuration}</p>
                                                    )}
                                                </div>
                                            )}
                                            {/* Progress section (only when there are orders or still syncing) */}
                                            {!(syncProgress && syncProgress.totalFetched === 0 && !syncing) && (
                                            <>
                                            <div>
                                                <div className="flex justify-between items-center mb-2">
                                                    <h4 className="text-sm font-semibold text-gray-700 dark:text-gray-200 dark:text-gray-200">Progreso de Sincronización</h4>
                                                    {syncProgress && (
                                                        <span className="text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400">
                                                            {(() => {
                                                                const processed = syncProgress.created + syncProgress.updated + syncProgress.rejected;
                                                                if (allFinished && processed > 0) return `${processed} procesadas`;
                                                                const denominator = syncProgress.totalFetched ?? syncProgress.total;
                                                                if (denominator > 0) return `${processed} / ${denominator} procesadas`;
                                                                if (syncProgress.totalFetched === null && syncing) return 'Obteniendo órdenes...';
                                                                return 'Iniciando...';
                                                            })()}
                                                        </span>
                                                    )}
                                                </div>
                                                {/* Progress bar */}
                                                <div className="w-full bg-gray-200 rounded-full h-4 overflow-hidden flex">
                                                    {(() => {
                                                        const processed = (syncProgress?.created ?? 0) + (syncProgress?.updated ?? 0) + (syncProgress?.rejected ?? 0);
                                                        // When finished, use processed as denominator so bar reaches 100%
                                                        const denominator = allFinished && processed > 0
                                                            ? processed
                                                            : (syncProgress?.totalFetched ?? syncProgress?.total ?? 0);
                                                        if (syncProgress && denominator > 0) {
                                                            return (
                                                                <>
                                                                    {syncProgress.created > 0 && (
                                                                        <div className="h-full bg-green-500 transition-all duration-500 ease-out" style={{ width: `${(syncProgress.created / denominator) * 100}%` }} title={`${syncProgress.created} creadas`} />
                                                                    )}
                                                                    {syncProgress.rejected > 0 && (
                                                                        <div className="h-full bg-red-500 transition-all duration-500 ease-out" style={{ width: `${(syncProgress.rejected / denominator) * 100}%` }} title={`${syncProgress.rejected} rechazadas`} />
                                                                    )}
                                                                    {syncProgress.updated > 0 && (
                                                                        <div className="h-full bg-yellow-500 transition-all duration-500 ease-out" style={{ width: `${(syncProgress.updated / denominator) * 100}%` }} title={`${syncProgress.updated} actualizadas`} />
                                                                    )}
                                                                </>
                                                            );
                                                        }
                                                        return <div className="h-full bg-blue-500 animate-pulse" style={{ width: '10%' }} />;
                                                    })()}
                                                </div>
                                                {/* Stats */}
                                                {syncProgress && (
                                                    <div className="flex gap-4 mt-3 text-xs">
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-3 h-3 rounded-full bg-green-500" />
                                                            <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Creadas: <strong>{syncProgress.created}</strong></span>
                                                        </div>
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-3 h-3 rounded-full bg-red-500" />
                                                            <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Rechazadas: <strong>{syncProgress.rejected}</strong></span>
                                                        </div>
                                                        {syncProgress.skipped > 0 && (
                                                            <div className="flex items-center gap-1">
                                                                <div className="w-3 h-3 rounded-full bg-amber-500" />
                                                                <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Omitidas: <strong>{syncProgress.skipped}</strong></span>
                                                            </div>
                                                        )}
                                                        {syncProgress.updated > 0 && (
                                                            <div className="flex items-center gap-1">
                                                                <div className="w-3 h-3 rounded-full bg-yellow-500" />
                                                                <span className="text-gray-600 dark:text-gray-300 dark:text-gray-300">Actualizadas: <strong>{syncProgress.updated}</strong></span>
                                                            </div>
                                                        )}
                                                    </div>
                                                )}
                                            </div>
                                            {/* Orders list */}
                                            {syncProgress && syncProgress.orders.length > 0 && (
                                                <div className="border border-gray-200 dark:border-gray-700 rounded-lg overflow-hidden">
                                                    <div className="p-2 bg-gray-50 dark:bg-gray-700 border-b border-gray-200 dark:border-gray-700 sticky top-0">
                                                        <p className="text-xs font-semibold text-gray-700 dark:text-gray-200 dark:text-gray-200">Órdenes Procesadas</p>
                                                    </div>
                                                    <div className="max-h-64 overflow-y-auto divide-y divide-gray-100">
                                                        {syncProgress.orders.slice(0, 50).map((order, index) => (
                                                            <div
                                                                key={index}
                                                                className={`p-3 text-xs ${order.status === 'created' ? 'bg-green-50 hover:bg-green-100' : order.status === 'updated' ? 'bg-blue-50 hover:bg-blue-100' : order.status === 'skipped' ? 'bg-amber-50 hover:bg-amber-100' : 'bg-red-50 hover:bg-red-100'} transition-colors`}
                                                            >
                                                                <div className="flex items-start justify-between gap-2">
                                                                    <div className="flex-1 min-w-0">
                                                                        <div className="flex items-center gap-2 mb-1">
                                                                            <div className={`w-2 h-2 rounded-full flex-shrink-0 ${order.status === 'created' ? 'bg-green-500' : order.status === 'updated' ? 'bg-blue-500' : order.status === 'skipped' ? 'bg-amber-500' : 'bg-red-500'}`} />
                                                                            <span className="font-medium text-gray-800 dark:text-gray-100">#{order.orderNumber}</span>
                                                                            {order.status === 'skipped' && (
                                                                                <span className="px-2 py-0.5 bg-amber-100 text-amber-700 rounded text-xs font-medium">Omitida</span>
                                                                            )}
                                                                            {order.orderStatus && (
                                                                                <span className="px-2 py-0.5 bg-gray-200 text-gray-700 dark:text-gray-200 dark:text-gray-200 rounded text-xs font-medium">{order.orderStatus}</span>
                                                                            )}
                                                                        </div>
                                                                        {order.createdAt && (
                                                                            <div className="text-gray-600 dark:text-gray-300 dark:text-gray-300 text-xs ml-4">
                                                                                Creada: {(() => {
                                                                                    try {
                                                                                        const date = new Date(order.createdAt);
                                                                                        return isNaN(date.getTime()) ? order.createdAt : date.toLocaleString('es-ES', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' });
                                                                                    } catch { return order.createdAt; }
                                                                                })()}
                                                                            </div>
                                                                        )}
                                                                        {(order.status === 'rejected' || order.status === 'skipped') && order.reason && (
                                                                            <div className={`text-xs ml-4 mt-1 ${order.status === 'skipped' ? 'text-amber-600' : 'text-red-600'}`}>{order.reason}</div>
                                                                        )}
                                                                    </div>
                                                                    <span className="text-gray-400 text-xs flex-shrink-0">{order.timestamp.toLocaleTimeString()}</span>
                                                                </div>
                                                            </div>
                                                        ))}
                                                    </div>
                                                    {syncProgress.orders.length > 50 && (
                                                        <div className="p-2 bg-gray-50 dark:bg-gray-700 text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400 text-center">
                                                            Y {syncProgress.orders.length - 50} órdenes más...
                                                        </div>
                                                    )}
                                                </div>
                                            )}
                                            </>
                                            )}
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>

                        {/* Footer */}
                        <div className="px-6 py-4 border-t border-gray-200 dark:border-gray-700 flex items-center justify-between flex-shrink-0 bg-white dark:bg-gray-800">
                            <div className="text-xs text-gray-500 dark:text-gray-400 dark:text-gray-400">
                                {syncing && isBatchMode && batchSync && (
                                    <span className="flex items-center gap-2">
                                        <Spinner size="sm" />
                                        Procesando lote {Math.min(batchSync.completedBatches + batchSync.failedBatches + 1, batchSync.totalBatches)} de {batchSync.totalBatches}... No cerrar esta ventana.
                                    </span>
                                )}
                                {syncing && !isBatchMode && (
                                    <span className="flex items-center gap-2">
                                        <Spinner size="sm" />
                                        {syncProgress?.totalFetched !== null && syncProgress?.totalFetched !== undefined
                                            ? `Procesando órdenes (${syncProgress.created + syncProgress.updated + syncProgress.rejected}/${syncProgress.totalFetched})... No cerrar.`
                                            : 'Obteniendo órdenes... No cerrar esta ventana.'}
                                    </span>
                                )}
                                {allFinished && isBatchMode && batchSync && batchSync.failedBatches > 0 && (
                                    <span className="text-amber-600 font-medium">
                                        {batchSync.failedBatches} lote(s) fallaron de {batchSync.totalBatches} total.
                                    </span>
                                )}
                                {allFinished && isBatchMode && batchSync && batchSync.failedBatches === 0 && (
                                    <span className="text-green-600 font-medium">
                                        Todos los lotes completados exitosamente.
                                    </span>
                                )}
                                {allFinished && !isBatchMode && syncError && (
                                    <span className="text-red-600 font-medium">
                                        Sincronización fallida: {syncError}
                                    </span>
                                )}
                                {allFinished && !isBatchMode && !syncError && syncProgress && (
                                    <span className="text-green-600 font-medium">
                                        Sincronización completada. {syncProgress.created} creadas, {syncProgress.updated} actualizadas, {syncProgress.rejected} rechazadas.
                                    </span>
                                )}
                            </div>
                            <div className="flex gap-3">
                                {!isShowingProgress && (
                                    <>
                                        <Button variant="outline" onClick={handleSyncCancel}>Cancelar</Button>
                                        <Button variant="primary" onClick={handleSyncConfirm}>Iniciar Sincronización</Button>
                                    </>
                                )}
                                {allFinished && (
                                    <Button variant="primary" onClick={handleSyncCancel}>Cerrar</Button>
                                )}
                                {syncing && (
                                    <Button variant="outline" disabled>Sincronizando...</Button>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
                );
            })()}
        </div>
    );
}
