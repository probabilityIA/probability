'use client';

import { Alert, Table, Spinner } from '@/shared/ui';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { useSyncLogs } from '@/services/modules/inventory/ui/hooks/useCapture';
import { InventorySyncLog } from '@/services/modules/inventory/domain/capture-types';

export default function SyncLogsPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const { logs, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters } = useSyncLogs({ business_id: businessId });

    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    const statusStyles: Record<string, string> = {
        success: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        failed: 'bg-red-100 dark:bg-red-900 text-red-800 dark:text-red-200',
        pending: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        duplicate: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300',
    };

    const directionStyles: Record<string, string> = {
        in: 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200',
        out: 'bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200',
    };

    const columns = [
        { key: 'id', label: '#', align: 'center' as const },
        { key: 'integration', label: 'Integración', align: 'center' as const },
        { key: 'direction', label: 'Dirección', align: 'center' as const },
        { key: 'hash', label: 'Hash' },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'error', label: 'Error' },
        { key: 'synced_at', label: 'Sincronizado', align: 'center' as const },
    ];

    const renderRow = (l: InventorySyncLog) => ({
        id: <span className="text-xs text-gray-500 font-mono">#{l.id}</span>,
        integration: <span className="text-sm text-gray-600 dark:text-gray-300">{l.integration_id ? `#${l.integration_id}` : '—'}</span>,
        direction: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${directionStyles[l.direction] || 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{l.direction === 'in' ? '⬇️ entrada' : '⬆️ salida'}</span>,
        hash: <span className="font-mono text-xs text-gray-500 dark:text-gray-400">{l.payload_hash.slice(0, 16)}…</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[l.status] || statusStyles.duplicate}`}>{l.status}</span>,
        error: l.error ? <span className="text-xs text-red-600 dark:text-red-400 truncate block max-w-sm" title={l.error}>{l.error}</span> : <span className="text-xs text-gray-400">—</span>,
        synced_at: <span className="text-xs text-gray-500">{l.synced_at ? new Date(l.synced_at).toLocaleString() : new Date(l.created_at).toLocaleString()}</span>,
    });

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">                {requiresBusiness ? <Alert type="info">Selecciona un negocio.</Alert> : (
                    <>
                        <div className="flex gap-3 items-end flex-wrap">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Integración (ID)</label>
                                <input type="number" value={filters.integration_id ?? ''} onChange={(e) => setFilters({ integration_id: e.target.value ? Number(e.target.value) : undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm w-32" />
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Dirección</label>
                                <select value={filters.direction || ''} onChange={(e) => setFilters({ direction: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value="">Ambas</option>
                                    <option value="in">Entrada</option>
                                    <option value="out">Salida</option>
                                </select>
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                                <select value={filters.status || ''} onChange={(e) => setFilters({ status: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value="">Todos</option>
                                    <option value="success">Exitoso</option>
                                    <option value="failed">Fallido</option>
                                    <option value="pending">Pendiente</option>
                                    <option value="duplicate">Duplicado</option>
                                </select>
                            </div>
                        </div>

                        {error && <Alert type="error">{error}</Alert>}

                        {loading && logs.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                                <Table columns={columns} data={logs.map(renderRow)} keyExtractor={(_, i) => String(logs[i]?.id || i)} emptyMessage="Sin logs" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
}
