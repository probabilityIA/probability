'use client';

import { useState, useEffect, useMemo, useCallback, useRef } from 'react';
import {
    PlayIcon,
    ArrowPathIcon,
    PencilIcon,
    PowerIcon,
    TrashIcon
} from '@heroicons/react/24/outline';
import { useIntegrations } from '../hooks/useIntegrations';
import { useSSE } from '@/shared/hooks/use-sse';
import { getActiveIntegrationTypesAction } from '../../infra/actions';
import {
    Integration,
    SyncOrdersParams
} from '../../domain/types';
import { getOrderStatusMappingsAction } from '@/services/modules/orderstatus/infra/actions';
import { OrderStatusMapping } from '@/services/modules/orderstatus/domain/types';
import { Input, Button, Badge, Spinner, Table, Alert, ConfirmModal, Select, DateRangePicker } from '@/shared/ui';
import { DynamicFilters, FilterOption, ActiveFilter } from '@/shared/ui/dynamic-filters';
import { useToast } from '@/shared/providers/toast-provider';
import { TokenStorage } from '@/shared/utils/token-storage';
import { playNotificationSound } from '@/shared/utils';

interface IntegrationListProps {
    onEdit?: (integration: Integration) => void;
    filterCategory?: string;
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

export default function IntegrationList({ onEdit, filterCategory: propFilterCategory }: IntegrationListProps) {
    const {
        integrations,
        loading,
        error,
        page,
        setPage,
        totalPages,
        search,
        setSearch,
        filterType,
        setFilterType,
        filterCategory,
        setFilterCategory,
        deleteIntegration,
        toggleActive,
        testConnection,
        syncOrders,
        setError
    } = useIntegrations(propFilterCategory || '');

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

    // Estado para tipos de integración para filtros
    const [integrationTypes, setIntegrationTypes] = useState<Array<{ value: string; label: string }>>([]);

    // Obtener tipos de integración para el filtro
    useEffect(() => {
        const fetchIntegrationTypes = async () => {
            try {
                const token = TokenStorage.getSessionToken();
                const response = await getActiveIntegrationTypesAction(token);
                if (response.success && response.data) {
                    const options = response.data.map((type) => ({
                        value: type.code || type.name.toLowerCase().replace(/\s+/g, '_'),
                        label: type.name,
                    }));
                    setIntegrationTypes(options);
                }
            } catch (err) {
                console.error('Error fetching integration types for filter:', err);
            }
        };
        fetchIntegrationTypes();
    }, []);

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
        totalFetched: number | null; // Total de órdenes publicadas a cola (set by integration.sync.completed)
        fetchDuration: string | null; // Duración de la fase de fetch
        orders: Array<{
            orderNumber: string;
            status: 'created' | 'rejected' | 'updated';
            reason?: string;
            createdAt?: string;
            orderStatus?: string;
            timestamp: Date;
        }>;
    } | null>(null);

    // Estado para sincronización por lotes
    const [batchSync, setBatchSync] = useState<BatchSyncState | null>(null);
    const batchSyncRef = useRef<BatchSyncState | null>(null);
    useEffect(() => { batchSyncRef.current = batchSync; }, [batchSync]);

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
                        if (!batch) return prev;
                        const newOrderCount = batch.orderCount + increment;
                        const updatedBatch = { ...batch, orderCount: newOrderCount };

                        // Si ya llegaron todas las órdenes de este lote, completar y avanzar
                        if (updatedBatch.totalFetched !== null && newOrderCount >= updatedBatch.totalFetched) {
                            updatedBatch.status = 'completed';
                            updatedBatch.completedAt = new Date();
                            const newCompleted = prev.completedBatches + 1;
                            const nextIdx = idx + 1;
                            const updatedBatches = prev.batches.map((b, i) => {
                                if (i === idx) return updatedBatch;
                                if (i === nextIdx && b.status === 'pending') return { ...b, status: 'processing' as const };
                                return b;
                            });
                            const allDone = newCompleted + prev.failedBatches >= prev.totalBatches;
                            if (allDone) setSyncing(false);
                            const newState = { ...prev, batches: updatedBatches, completedBatches: newCompleted, currentOrderBatchIndex: nextIdx };
                            batchSyncRef.current = newState;
                            return newState;
                        }

                        const updatedBatches = prev.batches.map((b, i) => i === idx ? updatedBatch : b);
                        const newState = { ...prev, batches: updatedBatches };
                        batchSyncRef.current = newState;
                        return newState;
                    });
                };

                switch (eventType) {
                    case 'integration.sync.order.created': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const createdAt = eventData.created_at || eventData.synced_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;

                        setSyncProgress(prev => {
                            const base = prev || { total: 0, created: 0, rejected: 0, updated: 0, totalFetched: null, fetchDuration: null, orders: [] };
                            const updated = { ...base, created: base.created + 1, total: base.total + 1, orders: [{ orderNumber, status: 'created' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }, ...base.orders] };
                            if (!batchSyncRef.current && updated.totalFetched !== null && updated.created + updated.updated + updated.rejected >= updated.totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });
                        if (batchSyncRef.current) {
                            advanceBatchOrderCount(1);
                        }
                        playNotificationSound();
                        showToast(`Orden creada: #${orderNumber}`, 'success');
                        break;
                    }
                    case 'integration.sync.order.updated': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const createdAt = eventData.created_at || eventData.updated_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;

                        setSyncProgress(prev => {
                            const base = prev || { total: 0, created: 0, rejected: 0, updated: 0, totalFetched: null, fetchDuration: null, orders: [] };
                            const updated = { ...base, updated: base.updated + 1, total: base.total + 1, orders: [{ orderNumber, status: 'updated' as const, createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }, ...base.orders] };
                            if (!batchSyncRef.current && updated.totalFetched !== null && updated.created + updated.updated + updated.rejected >= updated.totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });
                        if (batchSyncRef.current) {
                            advanceBatchOrderCount(1);
                        }
                        playNotificationSound();
                        showToast(`Orden actualizada: #${orderNumber}`, 'info');
                        break;
                    }
                    case 'integration.sync.order.rejected': {
                        const orderNumber = eventData.order_number || event.metadata?.order_number || 'Desconocida';
                        const reason = eventData.reason || event.metadata?.reason || 'Error desconocido';
                        const error = eventData.error || event.metadata?.error || '';
                        const createdAt = eventData.created_at || eventData.rejected_at || event.metadata?.created_at || null;
                        const orderStatus = eventData.status || event.metadata?.status || null;

                        setSyncProgress(prev => {
                            const base = prev || { total: 0, created: 0, rejected: 0, updated: 0, totalFetched: null, fetchDuration: null, orders: [] };
                            const updated = { ...base, rejected: base.rejected + 1, total: base.total + 1, orders: [{ orderNumber, status: 'rejected' as const, reason: reason + (error ? `: ${error}` : ''), createdAt: createdAt || undefined, orderStatus: orderStatus || undefined, timestamp: new Date() }, ...base.orders] };
                            if (!batchSyncRef.current && updated.totalFetched !== null && updated.created + updated.updated + updated.rejected >= updated.totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });
                        if (batchSyncRef.current) {
                            advanceBatchOrderCount(1);
                        }
                        playNotificationSound();
                        showToast(`Orden rechazada: #${orderNumber} - ${reason}${error ? `: ${error}` : ''}`, 'error');
                        break;
                    }
                    case 'integration.sync.started': {
                        // En modo batch, ignorar: el provider emite esto por CADA lote y resetea contadores
                        if (batchSyncRef.current) break;

                        const integrationId = event.integration_id;
                        const integration = integrations.find(i => i.id === integrationId);
                        const integrationName = integration?.name || `Integración ${integrationId}`;

                        setSyncProgress({ total: 0, created: 0, rejected: 0, updated: 0, totalFetched: null, fetchDuration: null, orders: [] });
                        showToast(`Sincronización iniciada: ${integrationName}`, 'info');
                        break;
                    }
                    case 'integration.sync.completed': {
                        // En modo batch, capturar totalFetched en el lote actual (no en el global)
                        if (batchSyncRef.current) {
                            const batchTotalFetched = Number(eventData.total_fetched) || 0;
                            setBatchSync(prev => {
                                if (!prev) return prev;
                                const idx = prev.currentOrderBatchIndex;
                                const batch = prev.batches[idx];
                                if (!batch) return prev;
                                const updatedBatch = { ...batch, totalFetched: batchTotalFetched };

                                // Si totalFetched=0 o ya llegaron todas las órdenes, completar y avanzar
                                if (batchTotalFetched === 0 || updatedBatch.orderCount >= batchTotalFetched) {
                                    updatedBatch.status = 'completed';
                                    updatedBatch.completedAt = new Date();
                                    const newCompleted = prev.completedBatches + 1;
                                    const nextIdx = idx + 1;
                                    const updatedBatches = prev.batches.map((b, i) => {
                                        if (i === idx) return updatedBatch;
                                        if (i === nextIdx && b.status === 'pending') return { ...b, status: 'processing' as const };
                                        return b;
                                    });
                                    const allDone = newCompleted + prev.failedBatches >= prev.totalBatches;
                                    if (allDone) setSyncing(false);
                                    const newState = { ...prev, batches: updatedBatches, completedBatches: newCompleted, currentOrderBatchIndex: nextIdx };
                                    batchSyncRef.current = newState;
                                    return newState;
                                }

                                const updatedBatches = prev.batches.map((b, i) => i === idx ? updatedBatch : b);
                                const newState = { ...prev, batches: updatedBatches };
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
                                ...(prev || { total: 0, created: 0, rejected: 0, updated: 0, orders: [] }),
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
                                ...(prev || { total: 0, created: 0, rejected: 0, updated: 0, orders: [] }),
                                totalFetched,
                                fetchDuration: duration,
                            };
                            // Si ya se procesaron todas las órdenes, marcar como terminado
                            if (updated.created + updated.updated + updated.rejected >= totalFetched) {
                                setSyncing(false);
                            }
                            return updated;
                        });

                        playNotificationSound();
                        showToast(`Fetch completado: ${integrationName} - ${totalFetched} órdenes obtenidas, procesando...`, 'info');
                        break;
                    }
                    case 'integration.sync.failed': {
                        // En modo batch, los fallos se manejan con batch.failed
                        if (batchSyncRef.current) break;

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
                        };
                        // Set ref synchronously so subsequent events in the same tick see batch mode
                        batchSyncRef.current = newBatchState;
                        setBatchSync(newBatchState);

                        // Initialize order progress
                        setSyncProgress({ total: 0, created: 0, rejected: 0, updated: 0, totalFetched: null, fetchDuration: null, orders: [] });
                        showToast(`Sincronización por lotes iniciada: ${totalBatches} lotes`, 'info');
                        break;
                    }
                    case 'integration.sync.batch.completed': {
                        const batchIndex = Number(eventData.batch_index);
                        const duration = eventData.duration || '';
                        const batchDateFrom = eventData.created_at_min || '';
                        const batchDateTo = eventData.created_at_max || '';

                        // Respaldo: guardar duración y fechas reales.
                        // La completitud real la maneja advanceBatchOrderCount cuando orderCount >= totalFetched.
                        setBatchSync(prev => {
                            if (!prev) return prev;
                            const batch = prev.batches.find(b => b.batchIndex === batchIndex);
                            if (!batch) return prev;
                            // Solo actualizar metadata (duración, fechas). No cambiar status.
                            const updatedBatches = prev.batches.map(b => {
                                if (b.batchIndex === batchIndex) {
                                    return { ...b, duration, dateFrom: batchDateFrom || b.dateFrom, dateTo: batchDateTo || b.dateTo };
                                }
                                return b;
                            });
                            return { ...prev, batches: updatedBatches };
                        });
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

    // Filtros dinámicos - sincronizar con el hook
    const [filters, setFilters] = useState<{
        search?: string;
        type?: string;
    }>({});

    // Definir filtros disponibles
    const availableFilters: FilterOption[] = useMemo(() => {
        return [
            {
                key: 'search',
                label: 'Nombre',
                type: 'text',
                placeholder: 'Buscar por nombre...',
            },
            {
                key: 'type',
                label: 'Tipo',
                type: 'select',
                options: integrationTypes,
            },
        ];
    }, [integrationTypes]);

    // Convertir filtros a ActiveFilter[]
    const activeFilters: ActiveFilter[] = useMemo(() => {
        const active: ActiveFilter[] = [];

        if (filters.search) {
            active.push({
                key: 'search',
                label: 'Nombre',
                value: filters.search,
                type: 'text',
            });
        }

        if (filters.type) {
            active.push({
                key: 'type',
                label: 'Tipo',
                value: filters.type,
                type: 'select',
            });
        }

        return active;
    }, [filters]);

    // Manejar agregar filtro
    const handleAddFilter = useCallback((filterKey: string, value: any) => {
        setFilters((prev) => {
            const newFilters = { ...prev, [filterKey]: value };
            // Sincronizar con el estado del hook
            if (filterKey === 'search') {
                setSearch(value);
            } else if (filterKey === 'type') {
                setFilterType(value);
            }
            return newFilters;
        });
        setPage(1);
    }, [setSearch, setFilterType, setPage]);

    // Manejar eliminar filtro
    const handleRemoveFilter = useCallback((filterKey: string) => {
        setFilters((prev) => {
            const newFilters = { ...prev };
            delete (newFilters as any)[filterKey];
            // Sincronizar con el estado del hook
            if (filterKey === 'search') {
                setSearch('');
            } else if (filterKey === 'type') {
                setFilterType('');
            }
            return newFilters;
        });
        setPage(1);
    }, [setSearch, setFilterType, setPage]);

    // Manejar cambio de ordenamiento
    const handleSortChange = useCallback((sortBy: string, sortOrder: 'asc' | 'desc') => {
        // Por ahora no tenemos ordenamiento en integraciones, pero podemos agregarlo después
        console.log('Sort changed:', sortBy, sortOrder);
    }, []);

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

    const columns = [
        { key: 'id', label: 'ID' },
        { key: 'logo', label: 'Logo' },
        { key: 'name', label: 'Nombre' },
        { key: 'type', label: 'Tipo' },
        { key: 'category', label: 'Categoría' },
        { key: 'business', label: 'Empresa' },
        { key: 'status', label: 'Estado' },
        { key: 'actions', label: 'Acciones' }
    ];

    const renderRow = (integration: Integration) => ({
        id: integration.id,
        logo: (
            <div className="flex items-center justify-center">
                {integration.integration_type?.image_url ? (
                    <img
                        src={integration.integration_type.image_url}
                        alt={integration.integration_type.name || integration.name}
                        className="w-12 h-12 object-contain border border-gray-200 rounded-lg p-1 bg-white"
                        onError={(e) => {
                            // Si la imagen falla al cargar, mostrar un placeholder
                            (e.target as HTMLImageElement).style.display = 'none';
                            const parent = (e.target as HTMLImageElement).parentElement;
                            if (parent) {
                                parent.innerHTML = '<div class="w-12 h-12 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-xs">Sin logo</div>';
                            }
                        }}
                    />
                ) : (
                    <div className="w-12 h-12 flex items-center justify-center bg-gray-100 rounded-lg text-gray-400 text-xs">
                        {integration.name.charAt(0).toUpperCase()}
                    </div>
                )}
            </div>
        ),
        name: (
            <div>
                <div className="text-sm font-medium text-gray-900">{integration.name}</div>
                {integration.description && (
                    <div className="text-sm text-gray-500">{integration.description}</div>
                )}
            </div>
        ),
        type: (
            <div className="text-sm text-gray-700">
                {integration.integration_type?.name || (
                    <span className="text-gray-400 text-xs">—</span>
                )}
            </div>
        ),
        category: (
            <span
                className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium text-white"
                style={{ backgroundColor: integration.category_color || '#6B7280' }}
            >
                {integration.category_name || integration.category || 'Sin categoría'}
            </span>
        ),
        business: (
            <div className="text-sm text-gray-700">
                {integration.business_name || (
                    <span className="text-gray-400 text-xs">Sin empresa</span>
                )}
            </div>
        ),
        status: (
            <div className="flex items-center gap-2">
                <Badge type={integration.is_active ? 'success' : 'error'}>
                    {integration.is_active ? 'Activo' : 'Inactivo'}
                </Badge>
                {integration.is_default && (
                    <Badge type="primary">Por defecto</Badge>
                )}
            </div>
        ),
        actions: (
            <div className="flex gap-2 items-center">
                <button
                    onClick={() => handleTest(integration.id)}
                    className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
                    title="Probar conexión"
                    aria-label="Probar conexión"
                >
                    <PlayIcon className="w-4 h-4" />
                </button>
                {integration.category === 'ecommerce' && (
                    <button
                        onClick={() => handleSyncClick(integration.id, integration.name)}
                        className="p-2 bg-green-500 hover:bg-green-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-green-500 focus:ring-offset-2"
                        title="Sincronizar órdenes"
                        aria-label="Sincronizar órdenes"
                    >
                        <ArrowPathIcon className="w-4 h-4" />
                    </button>
                )}
                {onEdit && (
                    <button
                        onClick={() => onEdit(integration)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-yellow-500 focus:ring-offset-2"
                        title="Editar integración"
                        aria-label="Editar integración"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => toggleActive(integration.id, integration.is_active)}
                    className={`p-2 rounded-md transition-colors duration-200 focus:ring-2 focus:ring-offset-2 ${integration.is_active
                        ? 'bg-orange-500 hover:bg-orange-600 text-white focus:ring-orange-500'
                        : 'bg-gray-500 hover:bg-gray-600 text-white focus:ring-gray-500'
                        }`}
                    title={integration.is_active ? 'Desactivar integración' : 'Activar integración'}
                    aria-label={integration.is_active ? 'Desactivar integración' : 'Activar integración'}
                >
                    <PowerIcon className="w-4 h-4" />
                </button>
                <button
                    onClick={() => handleDeleteClick(integration.id)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors duration-200 focus:ring-2 focus:ring-red-500 focus:ring-offset-2"
                    title="Eliminar integración"
                    aria-label="Eliminar integración"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        )
    });

    return (
        <div className="space-y-4">
            {/* Dynamic Filters */}
            <DynamicFilters
                availableFilters={availableFilters}
                activeFilters={activeFilters}
                onAddFilter={handleAddFilter}
                onRemoveFilter={handleRemoveFilter}
                sortBy="created_at"
                sortOrder="desc"
                onSortChange={handleSortChange}
                sortOptions={[
                    { value: 'created_at', label: 'Ordenar por fecha de creación' },
                    { value: 'updated_at', label: 'Ordenar por fecha de actualización' },
                    { value: 'name', label: 'Ordenar por nombre' },
                ]}
            />

            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white rounded-b-lg rounded-t-none shadow-sm border border-gray-200 border-t-0 overflow-hidden">
                <Table
                    columns={columns}
                    data={integrations.map(renderRow)}
                    emptyMessage="No hay integraciones disponibles"
                />
            </div>

            {/* Pagination */}
            {totalPages > 1 && (
                <div className="flex justify-center items-center gap-2">
                    <Button
                        onClick={() => setPage(page - 1)}
                        disabled={page === 1}
                        variant="primary"
                    >
                        Anterior
                    </Button>
                    <span className="text-sm text-gray-700">
                        Página {page} de {totalPages}
                    </span>
                    <Button
                        onClick={() => setPage(page + 1)}
                        disabled={page === totalPages}
                        variant="primary"
                    >
                        Siguiente
                    </Button>
                </div>
            )}

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
                    <div className={`bg-white rounded-lg shadow-xl flex flex-col max-h-[90vh] overflow-hidden transition-all duration-300 ${isShowingProgress ? 'w-full max-w-5xl' : 'w-full max-w-lg'}`}>
                        {/* Header */}
                        <div className="px-6 py-4 border-b border-gray-200 flex-shrink-0">
                            <div className="flex items-center justify-between">
                                <div className="flex items-center gap-3">
                                    <div>
                                        <h3 className="text-lg font-semibold text-gray-900">
                                            Sincronizar Órdenes
                                        </h3>
                                        <p className="text-sm text-gray-500 mt-0.5">
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
                                        className="text-gray-400 hover:text-gray-600 transition-colors"
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
                        <div className="flex-1 min-h-0 overflow-y-auto">
                            {/* Phase 1: Form (before sync starts) */}
                            {!isShowingProgress && (
                                <div className="px-6 py-4 space-y-4">
                                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                                        <div className="sm:col-span-2">
                                            <label className="block text-sm font-medium text-gray-700 mb-1">Rango de Fechas</label>
                                            <DateRangePicker
                                                startDate={syncFilters.created_at_min}
                                                endDate={syncFilters.created_at_max}
                                                onChange={(startDate, endDate) => {
                                                    setSyncFilters(prev => ({ ...prev, created_at_min: startDate || '', created_at_max: endDate || '' }));
                                                }}
                                                placeholder="Seleccionar rango de fechas (opcional)"
                                            />
                                        </div>
                                        <div>
                                            <label className="block text-sm font-medium text-gray-700 mb-1">Estado de Orden</label>
                                            {loadingMappings ? (
                                                <div className="flex items-center gap-2 text-sm text-gray-500"><Spinner size="sm" /> Cargando...</div>
                                            ) : (
                                                <Select
                                                    value={syncFilters.status || 'any'}
                                                    onChange={(e) => setSyncFilters(prev => ({ ...prev, status: e.target.value }))}
                                                    options={orderStatusOptions.length > 0 ? orderStatusOptions : [{ value: 'any', label: 'Todos' }]}
                                                    disabled={orderStatusOptions.length === 0}
                                                />
                                            )}
                                        </div>
                                        <div>
                                            <label className="block text-sm font-medium text-gray-700 mb-1">Estado Financiero</label>
                                            {loadingMappings ? (
                                                <div className="flex items-center gap-2 text-sm text-gray-500"><Spinner size="sm" /> Cargando...</div>
                                            ) : (
                                                <Select
                                                    value={syncFilters.financial_status || 'any'}
                                                    onChange={(e) => setSyncFilters(prev => ({ ...prev, financial_status: e.target.value }))}
                                                    options={financialStatusOptions.length > 0 ? financialStatusOptions : [{ value: 'any', label: 'Todos' }]}
                                                    disabled={financialStatusOptions.length === 0}
                                                />
                                            )}
                                        </div>
                                        <div className="sm:col-span-2">
                                            <label className="block text-sm font-medium text-gray-700 mb-1">Estado de Envío</label>
                                            {loadingMappings ? (
                                                <div className="flex items-center gap-2 text-sm text-gray-500"><Spinner size="sm" /> Cargando...</div>
                                            ) : (
                                                <Select
                                                    value={syncFilters.fulfillment_status || 'any'}
                                                    onChange={(e) => setSyncFilters(prev => ({ ...prev, fulfillment_status: e.target.value }))}
                                                    options={fulfillmentStatusOptions.length > 0 ? fulfillmentStatusOptions : [{ value: 'any', label: 'Todos' }]}
                                                    disabled={fulfillmentStatusOptions.length === 0}
                                                />
                                            )}
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
                                                        <h4 className="text-sm font-semibold text-gray-700">Progreso por Lotes</h4>
                                                        <span className="text-xs text-gray-500">
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
                                                        <span className="inline-flex items-center gap-1 px-2 py-0.5 bg-gray-100 text-gray-600 rounded-full">
                                                            <div className="w-2 h-2 rounded-full bg-gray-400" /> {batchSync.totalBatches - batchSync.completedBatches - batchSync.failedBatches} pendientes
                                                        </span>
                                                    </div>
                                                </div>

                                                {/* Batch history */}
                                                <div className="border border-gray-200 rounded-lg overflow-hidden">
                                                    <div className="p-2 bg-gray-50 border-b border-gray-200 sticky top-0 z-10">
                                                        <p className="text-xs font-semibold text-gray-700">Historial de Lotes</p>
                                                    </div>
                                                    <div className="max-h-[50vh] overflow-y-auto divide-y divide-gray-100">
                                                        {batchSync.batches.map((batch) => (
                                                            <div
                                                                key={batch.batchIndex}
                                                                className={`p-3 text-xs transition-colors ${
                                                                    batch.status === 'processing' ? 'bg-blue-50 border-l-2 border-l-blue-500' :
                                                                    batch.status === 'completed' ? 'bg-white hover:bg-gray-50' :
                                                                    batch.status === 'failed' ? 'bg-red-50 hover:bg-red-100' :
                                                                    'bg-gray-50'
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
                                                                            <div className="w-4 h-4 rounded-full border-2 border-gray-300" />
                                                                        )}
                                                                        <span className="font-medium text-gray-800">Lote {batch.batchIndex + 1}</span>
                                                                    </div>
                                                                    <div className="flex items-center gap-2">
                                                                        {batch.duration && (
                                                                            <span className="text-gray-400">{batch.duration}</span>
                                                                        )}
                                                                        <span className={`px-1.5 py-0.5 rounded text-[10px] font-medium ${
                                                                            batch.status === 'completed' ? 'bg-green-100 text-green-700' :
                                                                            batch.status === 'failed' ? 'bg-red-100 text-red-700' :
                                                                            batch.status === 'processing' ? 'bg-blue-100 text-blue-700' :
                                                                            'bg-gray-100 text-gray-500'
                                                                        }`}>
                                                                            {batch.status === 'completed' ? 'Completado' :
                                                                             batch.status === 'failed' ? 'Fallido' :
                                                                             batch.status === 'processing' ? 'Procesando' :
                                                                             'Pendiente'}
                                                                        </span>
                                                                    </div>
                                                                </div>
                                                                <div className="mt-1 ml-6 text-gray-500">
                                                                    {formatShortDate(batch.dateFrom)} → {formatShortDate(batch.dateTo)}
                                                                    {batch.orderCount > 0 && <span> · {batch.orderCount}{batch.totalFetched !== null ? `/${batch.totalFetched}` : ''} {batch.orderCount === 1 ? 'orden' : 'órdenes'}</span>}
                                                                    {batch.totalFetched === 0 && batch.status === 'completed' && <span> · Sin órdenes</span>}
                                                                </div>
                                                                {/* Mini progress bar per batch */}
                                                                {batch.status === 'processing' && batch.totalFetched !== null && batch.totalFetched > 0 && (
                                                                    <div className="mt-1.5 ml-6">
                                                                        <div className="w-full bg-gray-200 rounded-full h-1.5 overflow-hidden">
                                                                            <div
                                                                                className="h-full bg-blue-500 transition-all duration-300 ease-out"
                                                                                style={{ width: `${Math.min(100, Math.round((batch.orderCount / batch.totalFetched) * 100))}%` }}
                                                                            />
                                                                        </div>
                                                                    </div>
                                                                )}
                                                                {batch.status === 'completed' && batch.totalFetched !== null && batch.totalFetched > 0 && (
                                                                    <div className="mt-1.5 ml-6">
                                                                        <div className="w-full bg-gray-200 rounded-full h-1.5 overflow-hidden">
                                                                            <div className="h-full bg-green-500 w-full" />
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
                                                            <span className="text-gray-600">Creadas: <strong>{syncProgress.created}</strong></span>
                                                        </div>
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-2.5 h-2.5 rounded-full bg-blue-500" />
                                                            <span className="text-gray-600">Actualizadas: <strong>{syncProgress.updated}</strong></span>
                                                        </div>
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-2.5 h-2.5 rounded-full bg-red-500" />
                                                            <span className="text-gray-600">Rechazadas: <strong>{syncProgress.rejected}</strong></span>
                                                        </div>
                                                    </div>
                                                )}
                                                {/* Orders feed */}
                                                <div className="border border-gray-200 rounded-lg overflow-hidden">
                                                    <div className="p-2 bg-gray-50 border-b border-gray-200 sticky top-0 z-10">
                                                        <p className="text-xs font-semibold text-gray-700">Órdenes Procesadas</p>
                                                    </div>
                                                    <div className="max-h-[60vh] overflow-y-auto divide-y divide-gray-100">
                                                        {syncProgress && syncProgress.orders.length > 0 ? (
                                                            syncProgress.orders.slice(0, 100).map((order, index) => (
                                                                <div
                                                                    key={index}
                                                                    className={`px-3 py-2 text-xs ${
                                                                        order.status === 'created' ? 'bg-green-50' :
                                                                        order.status === 'updated' ? 'bg-blue-50' :
                                                                        'bg-red-50'
                                                                    }`}
                                                                >
                                                                    <div className="flex items-center justify-between gap-1">
                                                                        <div className="flex items-center gap-1.5 min-w-0">
                                                                            <div className={`w-1.5 h-1.5 rounded-full flex-shrink-0 ${
                                                                                order.status === 'created' ? 'bg-green-500' :
                                                                                order.status === 'updated' ? 'bg-blue-500' :
                                                                                'bg-red-500'
                                                                            }`} />
                                                                            <span className="font-medium text-gray-800 truncate">#{order.orderNumber}</span>
                                                                            {order.orderStatus && (
                                                                                <span className="px-1.5 py-0.5 bg-gray-200 text-gray-600 rounded text-[10px] flex-shrink-0">
                                                                                    {order.orderStatus}
                                                                                </span>
                                                                            )}
                                                                        </div>
                                                                        <span className="text-gray-400 text-[10px] flex-shrink-0">
                                                                            {order.timestamp.toLocaleTimeString('es-ES', { hour: '2-digit', minute: '2-digit', second: '2-digit' })}
                                                                        </span>
                                                                    </div>
                                                                    {order.status === 'rejected' && order.reason && (
                                                                        <div className="text-red-600 text-[10px] mt-0.5 ml-3 truncate" title={order.reason}>
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
                                                            <div className="p-2 bg-gray-50 text-xs text-gray-500 text-center">
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
                                            <div>
                                                <div className="flex justify-between items-center mb-2">
                                                    <h4 className="text-sm font-semibold text-gray-700">Progreso de Sincronización</h4>
                                                    {syncProgress && (
                                                        <span className="text-xs text-gray-500">
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
                                                            <span className="text-gray-600">Creadas: <strong>{syncProgress.created}</strong></span>
                                                        </div>
                                                        <div className="flex items-center gap-1">
                                                            <div className="w-3 h-3 rounded-full bg-red-500" />
                                                            <span className="text-gray-600">Rechazadas: <strong>{syncProgress.rejected}</strong></span>
                                                        </div>
                                                        {syncProgress.updated > 0 && (
                                                            <div className="flex items-center gap-1">
                                                                <div className="w-3 h-3 rounded-full bg-yellow-500" />
                                                                <span className="text-gray-600">Actualizadas: <strong>{syncProgress.updated}</strong></span>
                                                            </div>
                                                        )}
                                                    </div>
                                                )}
                                            </div>
                                            {/* Orders list */}
                                            {syncProgress && syncProgress.orders.length > 0 && (
                                                <div className="border border-gray-200 rounded-lg overflow-hidden">
                                                    <div className="p-2 bg-gray-50 border-b border-gray-200 sticky top-0">
                                                        <p className="text-xs font-semibold text-gray-700">Órdenes Procesadas</p>
                                                    </div>
                                                    <div className="max-h-64 overflow-y-auto divide-y divide-gray-100">
                                                        {syncProgress.orders.slice(0, 50).map((order, index) => (
                                                            <div
                                                                key={index}
                                                                className={`p-3 text-xs ${order.status === 'created' ? 'bg-green-50 hover:bg-green-100' : order.status === 'updated' ? 'bg-blue-50 hover:bg-blue-100' : 'bg-red-50 hover:bg-red-100'} transition-colors`}
                                                            >
                                                                <div className="flex items-start justify-between gap-2">
                                                                    <div className="flex-1 min-w-0">
                                                                        <div className="flex items-center gap-2 mb-1">
                                                                            <div className={`w-2 h-2 rounded-full flex-shrink-0 ${order.status === 'created' ? 'bg-green-500' : order.status === 'updated' ? 'bg-blue-500' : 'bg-red-500'}`} />
                                                                            <span className="font-medium text-gray-800">#{order.orderNumber}</span>
                                                                            {order.orderStatus && (
                                                                                <span className="px-2 py-0.5 bg-gray-200 text-gray-700 rounded text-xs font-medium">{order.orderStatus}</span>
                                                                            )}
                                                                        </div>
                                                                        {order.createdAt && (
                                                                            <div className="text-gray-600 text-xs ml-4">
                                                                                Creada: {(() => {
                                                                                    try {
                                                                                        const date = new Date(order.createdAt);
                                                                                        return isNaN(date.getTime()) ? order.createdAt : date.toLocaleString('es-ES', { year: 'numeric', month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit' });
                                                                                    } catch { return order.createdAt; }
                                                                                })()}
                                                                            </div>
                                                                        )}
                                                                        {order.status === 'rejected' && order.reason && (
                                                                            <div className="text-red-600 text-xs ml-4 mt-1">{order.reason}</div>
                                                                        )}
                                                                    </div>
                                                                    <span className="text-gray-400 text-xs flex-shrink-0">{order.timestamp.toLocaleTimeString()}</span>
                                                                </div>
                                                            </div>
                                                        ))}
                                                    </div>
                                                    {syncProgress.orders.length > 50 && (
                                                        <div className="p-2 bg-gray-50 text-xs text-gray-500 text-center">
                                                            Y {syncProgress.orders.length - 50} órdenes más...
                                                        </div>
                                                    )}
                                                </div>
                                            )}
                                        </div>
                                    )}
                                </div>
                            )}
                        </div>

                        {/* Footer */}
                        <div className="px-6 py-4 border-t border-gray-200 flex items-center justify-between flex-shrink-0 bg-white">
                            <div className="text-xs text-gray-500">
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
