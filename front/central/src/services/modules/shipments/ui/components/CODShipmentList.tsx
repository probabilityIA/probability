'use client';

import { useCallback, useEffect, useState } from 'react';
import {
    Search, Package, Truck, Calendar, MapPin, X, RefreshCw,
    DollarSign, CheckCircle2, Clock, XCircle, FileText,
    User, Hash, Building2, AlertCircle, ChevronLeft, ChevronRight,
} from 'lucide-react';
import { getCODShipmentsAction, collectCODAction, trackShipmentAction } from '../../infra/actions';
import { Shipment, EnvioClickTrackHistory } from '../../domain/types';
import { MiniAddressMap } from './MiniAddressMap';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';

interface Props {
    selectedBusinessId?: number | null;
}

const STATUS_CONFIG: Record<string, { label: string; cls: string }> = {
    pending: { label: 'Pendiente', cls: 'bg-amber-100 text-amber-700 border-amber-200' },
    in_transit: { label: 'En tránsito', cls: 'bg-blue-100 text-blue-700 border-blue-200' },
    delivered: { label: 'Entregado', cls: 'bg-emerald-100 text-emerald-700 border-emerald-200' },
    failed: { label: 'Fallido', cls: 'bg-red-100 text-red-700 border-red-200' },
    cancelled: { label: 'Cancelado', cls: 'bg-gray-100 text-gray-600 border-gray-200' },
};

function StatusBadge({ status }: { status: string }) {
    const cfg = STATUS_CONFIG[status] || { label: status, cls: 'bg-gray-100 text-gray-600 border-gray-200' };
    return (
        <span className={`inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold border ${cfg.cls}`}>
            {cfg.label}
        </span>
    );
}

function CODBadge({ isPaid }: { isPaid?: boolean }) {
    if (isPaid) {
        return (
            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold border bg-emerald-100 text-emerald-700 border-emerald-200">
                <CheckCircle2 size={12} /> Recaudado
            </span>
        );
    }
    return (
        <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full text-xs font-semibold border bg-orange-100 text-orange-700 border-orange-200">
            <Clock size={12} /> Por cobrar
        </span>
    );
}

function formatMoney(amount?: number, currency: string = 'COP') {
    if (amount == null) return '—';
    try {
        return new Intl.NumberFormat('es-CO', { style: 'currency', currency, maximumFractionDigits: 0 }).format(amount);
    } catch {
        return `${amount}`;
    }
}

function formatDate(s?: string) {
    if (!s) return '—';
    return new Date(s).toLocaleString('es-CO', {
        day: '2-digit', month: 'short', year: 'numeric',
        hour: '2-digit', minute: '2-digit', hour12: false, timeZone: 'America/Bogota',
    });
}

export default function CODShipmentList({ selectedBusinessId }: Props) {
    const [shipments, setShipments] = useState<Shipment[]>([]);
    const [loading, setLoading] = useState(false);
    const [page, setPage] = useState(1);
    const [pageSize] = useState(15);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(0);
    const [statusFilter, setStatusFilter] = useState('');
    const [paidFilter, setPaidFilter] = useState<'' | 'true' | 'false'>('');
    const [search, setSearch] = useState('');
    const [selected, setSelected] = useState<Shipment | null>(null);
    const [error, setError] = useState<string | null>(null);

    const fetchData = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: any = { page, page_size: pageSize };
            if (statusFilter) params.status = statusFilter;
            if (paidFilter) params.is_paid = paidFilter === 'true';
            if (selectedBusinessId) params.business_id = selectedBusinessId;
            const resp = await getCODShipmentsAction(params);
            if ((resp as any).success === false) {
                setError((resp as any).message || 'Error al cargar envios');
                setShipments([]);
                setTotal(0);
                setTotalPages(0);
            } else {
                setShipments(resp.data || []);
                setTotal(resp.total || 0);
                setTotalPages(resp.total_pages || 0);
            }
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, statusFilter, paidFilter, selectedBusinessId]);

    useEffect(() => {
        fetchData();
    }, [fetchData]);

    const filtered = shipments.filter(s => {
        if (!search.trim()) return true;
        const q = search.toLowerCase();
        return (
            s.customer_name?.toLowerCase().includes(q) ||
            s.client_name?.toLowerCase().includes(q) ||
            s.order_number?.toLowerCase().includes(q) ||
            s.tracking_number?.toLowerCase().includes(q)
        );
    });

    return (
        <div className="flex flex-col h-full bg-white dark:bg-gray-800 rounded-xl shadow-sm border border-gray-200 dark:border-gray-700">
            <div className="px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex items-center gap-3 flex-wrap">
                <div className="flex items-center gap-2">
                    <DollarSign className="text-emerald-600" size={20} />
                    <h2 className="text-lg font-bold text-gray-900 dark:text-white">Envíos contra entrega</h2>
                    <span className="text-xs text-gray-500 dark:text-gray-400">({total})</span>
                </div>
                <div className="flex-1" />
                <div className="relative">
                    <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                    <input
                        value={search}
                        onChange={e => setSearch(e.target.value)}
                        placeholder="Buscar cliente, orden, tracking..."
                        className="pl-9 pr-3 py-1.5 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white w-64"
                    />
                </div>
                <select
                    value={statusFilter}
                    onChange={e => { setStatusFilter(e.target.value); setPage(1); }}
                    className="px-2 py-1.5 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                >
                    <option value="">Todos los estados</option>
                    <option value="pending">Pendiente</option>
                    <option value="in_transit">En tránsito</option>
                    <option value="delivered">Entregado</option>
                    <option value="failed">Fallido</option>
                </select>
                <select
                    value={paidFilter}
                    onChange={e => { setPaidFilter(e.target.value as any); setPage(1); }}
                    className="px-2 py-1.5 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                >
                    <option value="">Cobro: todos</option>
                    <option value="false">Por cobrar</option>
                    <option value="true">Recaudado</option>
                </select>
                <button
                    onClick={fetchData}
                    disabled={loading}
                    className="p-1.5 rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-50"
                    title="Refrescar"
                >
                    <RefreshCw size={14} className={loading ? 'animate-spin' : ''} />
                </button>
            </div>

            <div className="flex-1 min-h-0 grid grid-cols-1 lg:grid-cols-3 gap-0 overflow-hidden">
                <div className="lg:col-span-1 border-r border-gray-200 dark:border-gray-700 overflow-y-auto">
                    {loading && (
                        <div className="flex items-center justify-center p-8 text-gray-400">
                            <RefreshCw size={20} className="animate-spin mr-2" /> Cargando...
                        </div>
                    )}
                    {error && (
                        <div className="m-3 p-3 bg-red-50 border border-red-200 rounded-md text-red-700 text-sm flex items-start gap-2">
                            <AlertCircle size={16} className="mt-0.5" /> {error}
                        </div>
                    )}
                    {!loading && filtered.length === 0 && !error && (
                        <div className="text-center text-gray-400 dark:text-gray-500 p-10 text-sm">
                            <Package size={32} className="mx-auto mb-2 opacity-50" />
                            No hay envíos contra entrega.
                        </div>
                    )}
                    <div className="divide-y divide-gray-100 dark:divide-gray-700">
                        {filtered.map(s => {
                            const isSelected = selected?.id === s.id;
                            const carrierLogo = s.carrier ? getCarrierLogo(s.carrier) : null;
                            return (
                                <button
                                    key={s.id}
                                    onClick={() => setSelected(s)}
                                    className={`w-full text-left px-4 py-3 hover:bg-gray-50 dark:hover:bg-gray-700/40 transition-colors ${isSelected ? 'bg-purple-50 dark:bg-purple-900/20 border-l-4 border-purple-500' : ''}`}
                                >
                                    <div className="flex items-start justify-between gap-2 mb-1">
                                        <div className="font-semibold text-sm text-gray-900 dark:text-white truncate">
                                            {s.customer_name || s.client_name || 'Sin cliente'}
                                        </div>
                                        <CODBadge isPaid={s.is_paid} />
                                    </div>
                                    <div className="flex items-center justify-between gap-2 text-xs text-gray-600 dark:text-gray-300">
                                        <span className="font-mono truncate">#{s.order_number || s.tracking_number || s.id}</span>
                                        <span className="font-bold text-emerald-600 dark:text-emerald-400">{formatMoney(s.cod_total, s.order_currency)}</span>
                                    </div>
                                    <div className="flex items-center justify-between gap-2 mt-1">
                                        <StatusBadge status={s.status} />
                                        <div className="flex items-center gap-1.5 text-[10px] text-gray-500 dark:text-gray-400">
                                            {carrierLogo && <img src={carrierLogo} alt={s.carrier} className="h-4 max-w-[60px] object-contain" />}
                                            <span>{formatDate(s.created_at).split(',')[0]}</span>
                                        </div>
                                    </div>
                                </button>
                            );
                        })}
                    </div>

                    {totalPages > 1 && (
                        <div className="flex items-center justify-between px-4 py-2 border-t border-gray-200 dark:border-gray-700 text-xs">
                            <button
                                onClick={() => setPage(p => Math.max(1, p - 1))}
                                disabled={page <= 1}
                                className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30"
                            >
                                <ChevronLeft size={14} />
                            </button>
                            <span className="text-gray-500 dark:text-gray-400">Pág {page} de {totalPages}</span>
                            <button
                                onClick={() => setPage(p => Math.min(totalPages, p + 1))}
                                disabled={page >= totalPages}
                                className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30"
                            >
                                <ChevronRight size={14} />
                            </button>
                        </div>
                    )}
                </div>

                <div className="lg:col-span-2 overflow-y-auto bg-gray-50 dark:bg-gray-900">
                    {selected ? (
                        <CODDetailPanel
                            shipment={selected}
                            businessId={selectedBusinessId || undefined}
                            onClose={() => setSelected(null)}
                            onCollected={() => { fetchData(); }}
                        />
                    ) : (
                        <div className="flex items-center justify-center h-full text-gray-400 dark:text-gray-500 text-sm">
                            <div className="text-center">
                                <Package size={48} className="mx-auto mb-3 opacity-40" />
                                Selecciona un envío para ver el detalle
                            </div>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}

interface DetailProps {
    shipment: Shipment;
    businessId?: number;
    onClose: () => void;
    onCollected: () => void;
}

function CODDetailPanel({ shipment, businessId, onClose, onCollected }: DetailProps) {
    const [collecting, setCollecting] = useState(false);
    const [showConfirm, setShowConfirm] = useState(false);
    const [notes, setNotes] = useState('');
    const [feedback, setFeedback] = useState<{ kind: 'ok' | 'err'; msg: string } | null>(null);
    const [tracking, setTracking] = useState<{ loading: boolean; history?: EnvioClickTrackHistory[]; error?: string }>({ loading: false });

    const canCollect = shipment.status === 'delivered' && !shipment.is_paid;
    const carrierLogo = shipment.carrier ? getCarrierLogo(shipment.carrier) : null;

    useEffect(() => {
        if (!shipment.tracking_number) return;
        setTracking({ loading: true });
        trackShipmentAction(shipment.tracking_number)
            .then((r: any) => {
                if (r?.success && r.data?.history) {
                    setTracking({ loading: false, history: r.data.history });
                } else {
                    setTracking({ loading: false, history: [] });
                }
            })
            .catch(() => setTracking({ loading: false, error: 'No se pudo cargar el historial' }));
    }, [shipment.tracking_number]);

    const handleCollect = async () => {
        setCollecting(true);
        setFeedback(null);
        try {
            const res = await collectCODAction(shipment.id, { notes }, businessId);
            if ((res as any).success === false) {
                setFeedback({ kind: 'err', msg: (res as any).message || 'Error al registrar el cobro' });
            } else {
                setFeedback({ kind: 'ok', msg: 'Cobro registrado exitosamente' });
                setShowConfirm(false);
                setNotes('');
                onCollected();
            }
        } finally {
            setCollecting(false);
        }
    };

    return (
        <div className="p-4 space-y-4">
            <div className="flex items-start justify-between">
                <div>
                    <div className="flex items-center gap-2 mb-1">
                        <h3 className="text-base font-bold text-gray-900 dark:text-white">
                            {shipment.customer_name || shipment.client_name || 'Sin cliente'}
                        </h3>
                        <CODBadge isPaid={shipment.is_paid} />
                        <StatusBadge status={shipment.status} />
                    </div>
                    <p className="text-xs text-gray-500 dark:text-gray-400 font-mono">
                        Orden #{shipment.order_number || shipment.id} · Tracking {shipment.tracking_number || '—'}
                    </p>
                </div>
                <button onClick={onClose} className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700">
                    <X size={16} />
                </button>
            </div>

            <div className="grid grid-cols-2 gap-3">
                <div className="bg-blue-50 dark:bg-blue-900/20 rounded-lg overflow-hidden">
                    <div className="p-3">
                        <div className="flex items-center gap-1.5 mb-1">
                            <Building2 size={12} className="text-blue-600" />
                            <span className="text-[10px] uppercase font-bold text-blue-700 dark:text-blue-300">Origen</span>
                        </div>
                        <p className="text-xs text-blue-900 dark:text-blue-100 truncate">{shipment.warehouse_name || 'Bodega principal'}</p>
                    </div>
                    <MiniAddressMap address={shipment.warehouse_name || 'Medellín'} city="Medellín" color="blue" />
                </div>
                <div className="bg-emerald-50 dark:bg-emerald-900/20 rounded-lg overflow-hidden">
                    <div className="p-3">
                        <div className="flex items-center gap-1.5 mb-1">
                            <MapPin size={12} className="text-emerald-600" />
                            <span className="text-[10px] uppercase font-bold text-emerald-700 dark:text-emerald-300">Destino</span>
                        </div>
                        <p className="text-xs text-emerald-900 dark:text-emerald-100 truncate" title={shipment.destination_address}>
                            {shipment.destination_address || 'Sin destino'}
                        </p>
                    </div>
                    <MiniAddressMap address={shipment.destination_address || 'Colombia'} color="emerald" />
                </div>
            </div>

            <div className="grid grid-cols-2 md:grid-cols-4 gap-2">
                <InfoCard icon={<DollarSign size={12} />} label="Monto a cobrar" value={formatMoney(shipment.cod_total, shipment.order_currency)} highlight />
                <InfoCard icon={<DollarSign size={12} />} label="Total orden" value={formatMoney(shipment.order_total_amount, shipment.order_currency)} />
                <InfoCard icon={<Truck size={12} />} label="Transportadora" value={shipment.carrier || '—'} />
                <InfoCard icon={<Calendar size={12} />} label="Entregado" value={shipment.delivered_at ? formatDate(shipment.delivered_at) : 'Pendiente'} />
            </div>

            {carrierLogo && (
                <div className="flex items-center gap-2 text-xs text-gray-500 dark:text-gray-400">
                    <img src={carrierLogo} alt={shipment.carrier} className="h-5 max-w-[80px] object-contain" />
                    <span>Guía: {shipment.guide_id || '—'}</span>
                    {shipment.guide_url && (
                        <a href={shipment.guide_url} target="_blank" rel="noopener noreferrer" className="ml-2 text-blue-600 hover:underline inline-flex items-center gap-1">
                            <FileText size={11} /> Ver guía
                        </a>
                    )}
                </div>
            )}

            {shipment.is_paid && shipment.paid_at && (
                <div className="bg-emerald-50 border border-emerald-200 rounded-md p-3 flex items-center gap-2 text-sm text-emerald-800">
                    <CheckCircle2 size={16} />
                    <span>Cobrado el {formatDate(shipment.paid_at)}</span>
                </div>
            )}

            {!shipment.is_paid && (
                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                    <div className="flex items-start justify-between gap-3">
                        <div>
                            <h4 className="font-semibold text-sm text-gray-900 dark:text-white mb-1">Registrar cobro contra entrega</h4>
                            <p className="text-xs text-gray-500 dark:text-gray-400">
                                {canCollect
                                    ? `Marca como pagado el monto de ${formatMoney(shipment.cod_total, shipment.order_currency)}.`
                                    : 'Solo se puede registrar el cobro cuando el envío esté entregado.'}
                            </p>
                        </div>
                        <button
                            disabled={!canCollect || collecting}
                            onClick={() => setShowConfirm(true)}
                            className="px-3 py-2 bg-emerald-600 hover:bg-emerald-700 disabled:opacity-50 disabled:cursor-not-allowed text-white text-sm font-semibold rounded-md inline-flex items-center gap-1.5"
                        >
                            <CheckCircle2 size={14} /> Marcar recaudado
                        </button>
                    </div>
                </div>
            )}

            {feedback && (
                <div className={`p-3 rounded-md text-sm ${feedback.kind === 'ok' ? 'bg-emerald-50 border border-emerald-200 text-emerald-800' : 'bg-red-50 border border-red-200 text-red-700'}`}>
                    {feedback.msg}
                </div>
            )}

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-lg p-4">
                <h4 className="font-semibold text-sm text-gray-900 dark:text-white mb-3 flex items-center gap-1.5">
                    <Truck size={14} /> Historial de tracking
                </h4>
                {tracking.loading && (
                    <div className="text-xs text-gray-400 flex items-center gap-2">
                        <RefreshCw size={12} className="animate-spin" /> Cargando historial...
                    </div>
                )}
                {tracking.error && <div className="text-xs text-red-600">{tracking.error}</div>}
                {!tracking.loading && tracking.history && tracking.history.length === 0 && (
                    <div className="text-xs text-gray-400">Sin eventos registrados.</div>
                )}
                {tracking.history && tracking.history.length > 0 && (
                    <ol className="space-y-2">
                        {tracking.history.map((ev, i) => (
                            <li key={i} className="flex gap-3 text-xs">
                                <div className="flex-shrink-0 w-2 h-2 rounded-full bg-blue-500 mt-1.5" />
                                <div className="flex-1">
                                    <div className="font-semibold text-gray-900 dark:text-white">{ev.status}</div>
                                    <div className="text-gray-600 dark:text-gray-400">{ev.description}</div>
                                    <div className="text-gray-400 dark:text-gray-500 mt-0.5">{ev.location} · {ev.date}</div>
                                </div>
                            </li>
                        ))}
                    </ol>
                )}
            </div>

            {showConfirm && (
                <div className="fixed inset-0 bg-black/50 z-50 flex items-center justify-center p-4">
                    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-5">
                        <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">Confirmar cobro</h3>
                        <p className="text-sm text-gray-600 dark:text-gray-300 mb-3">
                            ¿Marcar la orden <span className="font-mono">#{shipment.order_number}</span> como pagada por <strong>{formatMoney(shipment.cod_total, shipment.order_currency)}</strong>?
                        </p>
                        <label className="block text-xs font-semibold text-gray-700 dark:text-gray-200 mb-1">Notas / referencia (opcional)</label>
                        <textarea
                            value={notes}
                            onChange={e => setNotes(e.target.value)}
                            rows={3}
                            maxLength={500}
                            placeholder="Ej: Recibido en efectivo, comprobante #123"
                            className="w-full px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 bg-white dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                        <div className="flex justify-end gap-2 mt-4">
                            <button
                                onClick={() => setShowConfirm(false)}
                                disabled={collecting}
                                className="px-3 py-2 text-sm rounded-md border border-gray-200 dark:border-gray-600 hover:bg-gray-50 dark:hover:bg-gray-700"
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={handleCollect}
                                disabled={collecting}
                                className="px-3 py-2 text-sm rounded-md bg-emerald-600 hover:bg-emerald-700 text-white font-semibold disabled:opacity-50 inline-flex items-center gap-1.5"
                            >
                                {collecting ? <RefreshCw size={14} className="animate-spin" /> : <CheckCircle2 size={14} />}
                                Confirmar cobro
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

function InfoCard({ icon, label, value, highlight }: { icon: React.ReactNode; label: string; value: string; highlight?: boolean }) {
    return (
        <div className={`rounded-md p-2 border ${highlight ? 'bg-emerald-50 border-emerald-200 dark:bg-emerald-900/20 dark:border-emerald-700' : 'bg-white dark:bg-gray-800 border-gray-200 dark:border-gray-700'}`}>
            <div className="flex items-center gap-1 text-[10px] uppercase font-bold text-gray-500 dark:text-gray-400">
                {icon} {label}
            </div>
            <div className={`text-sm font-semibold mt-0.5 ${highlight ? 'text-emerald-700 dark:text-emerald-300' : 'text-gray-900 dark:text-white'}`}>
                {value}
            </div>
        </div>
    );
}
