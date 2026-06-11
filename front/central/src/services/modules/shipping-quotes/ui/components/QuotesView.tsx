'use client';

import { useCallback, useEffect, useState } from 'react';
import { useToast } from '@/shared/providers/toast-provider';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { getIntegrationsAction, setShopifyAutoGuideAction } from '@/services/integrations/core/infra/actions';
import { getOrderByIdAction } from '@/services/modules/orders/infra/actions';
import { Order } from '@/services/modules/orders/domain/types';
import OrderDetails from '@/services/modules/orders/ui/components/OrderDetails';
import { Modal } from '@/shared/ui';
import { getSavedQuotesAction, SavedQuote } from '../../infra/actions';
import CreateOrderFromQuoteModal from './CreateOrderFromQuoteModal';

interface QuotesViewProps {
    selectedBusinessId: number | null;
}

interface AutoGuideIntegration {
    id: number;
    name: string;
    type: string;
    imageUrl: string;
    enabled: boolean;
    loading: boolean;
}

const statusStyles: Record<string, string> = {
    created: 'bg-gray-100 text-gray-700',
    associated: 'bg-blue-100 text-blue-700',
    guide_generated: 'bg-emerald-100 text-emerald-700',
    failed: 'bg-red-100 text-red-700',
    expired: 'bg-amber-100 text-amber-700',
};

const statusLabels: Record<string, string> = {
    created: 'Creada',
    associated: 'Asociada',
    guide_generated: 'Guía generada',
    failed: 'Fallida',
    expired: 'Expirada',
};

export default function QuotesView({ selectedBusinessId }: QuotesViewProps) {
    const { showToast } = useToast();
    const { isSuperAdmin } = usePermissions();

    const [integrations, setIntegrations] = useState<AutoGuideIntegration[]>([]);
    const [quotes, setQuotes] = useState<SavedQuote[]>([]);
    const [loading, setLoading] = useState(false);
    const [page, setPage] = useState(1);
    const [totalPages, setTotalPages] = useState(1);
    const [total, setTotal] = useState(0);
    const [sourceFilter, setSourceFilter] = useState('');
    const [statusFilter, setStatusFilter] = useState('');
    const [createFromQuote, setCreateFromQuote] = useState<SavedQuote | null>(null);
    const [viewOrder, setViewOrder] = useState<Order | null>(null);
    const [loadingOrderId, setLoadingOrderId] = useState<string | null>(null);

    const needsBusiness = isSuperAdmin && !selectedBusinessId;

    const loadIntegrations = useCallback(async () => {
        const res: any = await getIntegrationsAction({
            business_id: selectedBusinessId || undefined,
            is_active: true,
            page: 1,
            page_size: 200,
        });
        const data = (res?.data || []) as any[];
        const ecommerce = data.filter((i) => i.type === 'shopify' || i.integration_type_id === 1);
        setIntegrations(
            ecommerce.map((i) => ({
                id: i.id,
                name: i.name,
                type: i.type,
                imageUrl: i.integration_type?.image_url || (i.type ? `/integrations/${i.type}.png` : ''),
                enabled: i.config?.auto_generate_guide_enabled === true,
                loading: false,
            }))
        );
    }, [selectedBusinessId]);

    const loadQuotes = useCallback(async () => {
        setLoading(true);
        const res = await getSavedQuotesAction({
            businessId: selectedBusinessId,
            page,
            pageSize: 10,
            source: sourceFilter,
            status: statusFilter,
        });
        setQuotes(res.data);
        setTotalPages(res.total_pages || 1);
        setTotal(res.total || 0);
        setLoading(false);
    }, [selectedBusinessId, page, sourceFilter, statusFilter]);

    useEffect(() => {
        if (needsBusiness) return;
        loadIntegrations();
    }, [needsBusiness, loadIntegrations]);

    useEffect(() => {
        if (needsBusiness) return;
        loadQuotes();
    }, [needsBusiness, loadQuotes]);

    const openOrderDetails = async (orderUuid: string) => {
        setLoadingOrderId(orderUuid);
        try {
            const res: any = await getOrderByIdAction(orderUuid);
            const order = res?.data;
            if (res?.success && order) {
                setViewOrder(order);
            } else {
                showToast(res?.message || 'No se pudo cargar la orden', 'error');
            }
        } catch {
            showToast('No se pudo cargar la orden', 'error');
        } finally {
            setLoadingOrderId(null);
        }
    };

    const toggleAutoGuide = async (integ: AutoGuideIntegration) => {
        setIntegrations((prev) => prev.map((i) => (i.id === integ.id ? { ...i, loading: true } : i)));
        const next = !integ.enabled;
        const result: any = await setShopifyAutoGuideAction(integ.id, next);
        if (!result || result.success === false) {
            showToast(result?.message || 'No se pudo actualizar la generación automática', 'error');
            setIntegrations((prev) => prev.map((i) => (i.id === integ.id ? { ...i, loading: false } : i)));
        } else {
            showToast(next ? 'Generación automática de guía activada' : 'Generación automática desactivada', 'success');
            setIntegrations((prev) => prev.map((i) => (i.id === integ.id ? { ...i, enabled: next, loading: false } : i)));
        }
    };

    if (needsBusiness) {
        return (
            <div className="flex items-center justify-center h-64 text-gray-500 dark:text-gray-400">
                Selecciona un negocio para ver sus cotizaciones.
            </div>
        );
    }

    return (
        <div className="space-y-6">
            <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm overflow-hidden">
                <div className="flex items-center gap-2 px-5 py-3 border-b border-gray-100 dark:border-gray-700 bg-gray-50/70 dark:bg-gray-800/60">
                    <span>🧾</span>
                    <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Generación automática de guía</h3>
                </div>
                <div className="p-5">
                    <p className="text-xs text-gray-500 dark:text-gray-400 mb-4">
                        Cuando una orden llega con una transportadora elegida en el checkout, la guía se genera sola con esa transportadora.
                        Configúralo por integración.
                    </p>
                    {integrations.length === 0 ? (
                        <p className="text-sm text-gray-400">No hay integraciones de e-commerce activas en este negocio.</p>
                    ) : (
                        <div className="divide-y divide-gray-100 dark:divide-gray-700">
                            {integrations.map((integ) => (
                                <div key={integ.id} className="flex items-center justify-between gap-3 py-2.5">
                                    <div className="flex items-center gap-3 min-w-0">
                                        <span className="w-9 h-9 rounded-lg border border-gray-200 dark:border-gray-700 bg-white flex items-center justify-center overflow-hidden shrink-0">
                                            {integ.imageUrl ? (
                                                // eslint-disable-next-line @next/next/no-img-element
                                                <img
                                                    src={integ.imageUrl}
                                                    alt={integ.type}
                                                    className="w-7 h-7 object-contain"
                                                    onError={(e) => {
                                                        const el = e.currentTarget as HTMLImageElement;
                                                        el.style.display = 'none';
                                                        if (el.parentElement) el.parentElement.textContent = (integ.type || '?').charAt(0).toUpperCase();
                                                    }}
                                                />
                                            ) : (
                                                <span className="text-xs text-gray-400">{(integ.type || '?').charAt(0).toUpperCase()}</span>
                                            )}
                                        </span>
                                        <div className="min-w-0">
                                            <p className="text-sm font-medium text-gray-800 dark:text-gray-100 truncate">{integ.name}</p>
                                            <p className="text-xs text-gray-400 capitalize">{integ.type}</p>
                                        </div>
                                    </div>
                                    <button
                                        type="button"
                                        onClick={() => toggleAutoGuide(integ)}
                                        disabled={integ.loading}
                                        className={`relative inline-flex h-5 w-9 items-center rounded-full transition-colors focus:outline-none shrink-0 disabled:opacity-50 ${integ.enabled ? 'bg-emerald-500' : 'bg-gray-300 dark:bg-gray-600'}`}
                                    >
                                        <span className={`inline-block h-3.5 w-3.5 transform rounded-full bg-white shadow-sm transition-transform ${integ.enabled ? 'translate-x-5' : 'translate-x-0.5'}`} />
                                    </button>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            </div>

            <div className="rounded-xl border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 shadow-sm overflow-hidden">
                <div className="flex flex-wrap items-center justify-between gap-3 px-5 py-3 border-b border-gray-100 dark:border-gray-700 bg-gray-50/70 dark:bg-gray-800/60">
                    <h3 className="text-sm font-semibold text-gray-800 dark:text-gray-100">Cotizaciones guardadas ({total})</h3>
                    <div className="flex items-center gap-2">
                        <select
                            value={sourceFilter}
                            onChange={(e) => { setPage(1); setSourceFilter(e.target.value); }}
                            className="text-xs rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 px-2 py-1.5"
                        >
                            <option value="">Todos los orígenes</option>
                            <option value="shopify">Shopify</option>
                            <option value="woocommerce">WooCommerce</option>
                            <option value="panel">Panel</option>
                        </select>
                        <select
                            value={statusFilter}
                            onChange={(e) => { setPage(1); setStatusFilter(e.target.value); }}
                            className="text-xs rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 px-2 py-1.5"
                        >
                            <option value="">Todos los estados</option>
                            <option value="created">Creada</option>
                            <option value="associated">Asociada</option>
                            <option value="guide_generated">Guía generada</option>
                            <option value="failed">Fallida</option>
                        </select>
                    </div>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="text-left text-xs text-gray-500 dark:text-gray-400 border-b border-gray-100 dark:border-gray-700">
                                <th className="px-4 py-2 font-medium">Fecha</th>
                                <th className="px-4 py-2 font-medium">Origen</th>
                                <th className="px-4 py-2 font-medium">Transportadora elegida</th>
                                <th className="px-4 py-2 font-medium">Tarifas</th>
                                <th className="px-4 py-2 font-medium">Orden</th>
                                <th className="px-4 py-2 font-medium">Estado</th>
                                <th className="px-4 py-2 font-medium text-right">Acciones</th>
                            </tr>
                        </thead>
                        <tbody>
                            {loading ? (
                                <tr><td colSpan={7} className="px-4 py-8 text-center text-gray-400">Cargando...</td></tr>
                            ) : quotes.length === 0 ? (
                                <tr><td colSpan={7} className="px-4 py-8 text-center text-gray-400">No hay cotizaciones.</td></tr>
                            ) : (
                                quotes.map((q) => (
                                    <tr key={q.id} className="border-b border-gray-50 dark:border-gray-700/50 hover:bg-gray-50/60 dark:hover:bg-gray-700/30">
                                        <td className="px-4 py-2.5 text-gray-600 dark:text-gray-300 whitespace-nowrap">
                                            {new Date(q.created_at).toLocaleString('es-CO', { dateStyle: 'short', timeStyle: 'short' })}
                                        </td>
                                        <td className="px-4 py-2.5 capitalize">{q.source}</td>
                                        <td className="px-4 py-2.5">{q.selected_carrier || <span className="text-gray-400">—</span>}</td>
                                        <td className="px-4 py-2.5 text-gray-500">{q.rates?.length || 0}</td>
                                        <td className="px-4 py-2.5 font-mono text-xs">
                                            {q.order_uuid ? (
                                                <button
                                                    onClick={() => openOrderDetails(q.order_uuid!)}
                                                    disabled={loadingOrderId === q.order_uuid}
                                                    className="text-purple-600 dark:text-purple-400 hover:underline font-semibold disabled:opacity-50"
                                                    title="Ver detalle de la orden"
                                                >
                                                    {loadingOrderId === q.order_uuid ? 'Cargando...' : `#${q.order_number || q.order_uuid.slice(0, 8)}`}
                                                </button>
                                            ) : <span className="text-gray-400">—</span>}
                                        </td>
                                        <td className="px-4 py-2.5">
                                            <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[q.status] || 'bg-gray-100 text-gray-700'}`}>
                                                {statusLabels[q.status] || q.status}
                                            </span>
                                        </td>
                                        <td className="px-4 py-2.5 text-right">
                                            {!q.order_uuid && (q.rates?.length || 0) > 0 && (
                                                <button
                                                    onClick={() => setCreateFromQuote(q)}
                                                    className="text-xs font-semibold px-3 py-1.5 rounded-md bg-purple-600 hover:bg-purple-700 text-white"
                                                >
                                                    Crear orden
                                                </button>
                                            )}
                                        </td>
                                    </tr>
                                ))
                            )}
                        </tbody>
                    </table>
                </div>

                {totalPages > 1 && (
                    <div className="flex items-center justify-between px-5 py-3 border-t border-gray-100 dark:border-gray-700">
                        <span className="text-xs text-gray-500">Página {page} de {totalPages}</span>
                        <div className="flex gap-2">
                            <button
                                onClick={() => setPage((p) => Math.max(1, p - 1))}
                                disabled={page <= 1}
                                className="text-xs px-3 py-1.5 rounded-md border border-gray-200 dark:border-gray-600 disabled:opacity-40"
                            >Anterior</button>
                            <button
                                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                                disabled={page >= totalPages}
                                className="text-xs px-3 py-1.5 rounded-md border border-gray-200 dark:border-gray-600 disabled:opacity-40"
                            >Siguiente</button>
                        </div>
                    </div>
                )}
            </div>

            {createFromQuote && (
                <CreateOrderFromQuoteModal
                    quote={createFromQuote}
                    businessId={selectedBusinessId}
                    onClose={() => setCreateFromQuote(null)}
                    onSuccess={loadQuotes}
                />
            )}

            <Modal isOpen={!!viewOrder} onClose={() => setViewOrder(null)} title={undefined} size="full">
                {viewOrder && (
                    <OrderDetails initialOrder={viewOrder} onClose={() => setViewOrder(null)} />
                )}
            </Modal>
        </div>
    );
}
