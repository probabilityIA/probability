'use client';

import { useEffect, useMemo, useState } from 'react';
import { X, FileText, RefreshCw, Truck, Download, Search, ChevronLeft, ChevronRight } from 'lucide-react';
import {
    getManifestCarriersAction,
    getManifestPendingAction,
    generateManifestPdfAction,
    ManifestCarrierOption,
    ManifestPendingShipment,
} from '../../infra/actions';
import { getCarrierLogo } from '@/shared/utils/carrier-logos';

function formatDate(s?: string): string {
    if (!s) return '—';
    try {
        const d = new Date(s);
        return d.toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: '2-digit' });
    } catch { return '—'; }
}

interface ManifestModalProps {
    isOpen: boolean;
    onClose: () => void;
    businessId: number | null;
}

const PAGE_SIZE_OPTIONS = [10, 25, 50, 100];

export function ManifestModal({ isOpen, onClose, businessId }: ManifestModalProps) {
    const [loadingCarriers, setLoadingCarriers] = useState(false);
    const [loadingRows, setLoadingRows] = useState(false);
    const [carriers, setCarriers] = useState<ManifestCarrierOption[]>([]);
    const [carrier, setCarrier] = useState<string>('');
    const [rows, setRows] = useState<ManifestPendingShipment[]>([]);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(25);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [selected, setSelected] = useState<Set<number>>(new Set());
    const [search, setSearch] = useState('');
    const [generating, setGenerating] = useState(false);

    useEffect(() => {
        if (!isOpen || !businessId) return;
        setLoadingCarriers(true);
        setSelected(new Set());
        setRows([]);
        setCarrier('');
        getManifestCarriersAction(businessId)
            .then((res) => {
                const list = res.success ? res.data : [];
                setCarriers(list);
                if (list.length > 0) setCarrier(list[0].carrier);
            })
            .finally(() => setLoadingCarriers(false));
    }, [isOpen, businessId]);

    useEffect(() => {
        if (!isOpen || !businessId || !carrier) return;
        setLoadingRows(true);
        getManifestPendingAction(businessId, carrier, page, pageSize)
            .then((res) => {
                setRows(res.success ? res.data : []);
                setTotal(res.total || 0);
                setTotalPages(res.total_pages || 1);
            })
            .finally(() => setLoadingRows(false));
    }, [isOpen, businessId, carrier, page, pageSize]);

    useEffect(() => { setPage(1); }, [carrier, pageSize]);

    const filteredRows = useMemo(() => {
        if (!search.trim()) return rows;
        const q = search.trim().toLowerCase();
        return rows.filter((s) =>
            s.order_number?.toLowerCase().includes(q) ||
            s.tracking_number?.toLowerCase().includes(q) ||
            s.customer_name?.toLowerCase().includes(q) ||
            s.destination_city?.toLowerCase().includes(q),
        );
    }, [rows, search]);

    const allPageSelected = filteredRows.length > 0 && filteredRows.every((s) => selected.has(s.shipment_id));

    const toggle = (id: number) => {
        const next = new Set(selected);
        if (next.has(id)) next.delete(id);
        else next.add(id);
        setSelected(next);
    };

    const togglePage = () => {
        const next = new Set(selected);
        if (allPageSelected) filteredRows.forEach((s) => next.delete(s.shipment_id));
        else filteredRows.forEach((s) => next.add(s.shipment_id));
        setSelected(next);
    };

    const handleGenerate = async () => {
        if (!businessId || selected.size === 0) return;
        setGenerating(true);
        try {
            const res = await generateManifestPdfAction(businessId, Array.from(selected), carrier);
            if (res.success && res.blob) {
                const a = document.createElement('a');
                a.href = res.blob;
                a.download = res.filename || 'manifiesto.pdf';
                document.body.appendChild(a);
                a.click();
                a.remove();
            } else {
                alert(res.message || 'Error al generar PDF');
            }
        } finally {
            setGenerating(false);
        }
    };

    if (!isOpen) return null;

    const carrierLogo = getCarrierLogo(carrier);

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4">
            <div className="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl w-full max-w-6xl max-h-[92vh] flex flex-col">
                <div className="flex items-center justify-between px-5 py-4 border-b border-gray-100 dark:border-gray-700 bg-gradient-to-r from-indigo-50 to-purple-50 dark:from-indigo-950/40 dark:to-purple-950/40 rounded-t-2xl">
                    <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-xl bg-indigo-100 dark:bg-indigo-900/40 flex items-center justify-center">
                            <FileText size={20} className="text-indigo-600 dark:text-indigo-300" />
                        </div>
                        <div>
                            <h3 className="text-lg font-bold text-gray-900 dark:text-white">Manifiesto de Recoleccion</h3>
                            <p className="text-xs text-gray-500 dark:text-gray-400">Selecciona los envios pendientes a entregar al transportador</p>
                        </div>
                    </div>
                    <button onClick={onClose} className="p-2 rounded-lg hover:bg-white/60 dark:hover:bg-gray-700 text-gray-500">
                        <X size={18} />
                    </button>
                </div>

                <div className="px-5 py-3 border-b border-gray-100 dark:border-gray-700 flex items-center gap-3 flex-wrap">
                    <div className="flex items-center gap-2">
                        {carrierLogo ? (
                            <img src={carrierLogo} alt={carrier} className="h-7 w-auto" />
                        ) : (
                            <Truck size={16} className="text-gray-400" />
                        )}
                        <select
                            value={carrier}
                            onChange={(e) => setCarrier(e.target.value)}
                            disabled={loadingCarriers || carriers.length === 0}
                            className="px-3 py-2 border border-gray-200 dark:border-gray-600 rounded-lg text-sm bg-gray-50 dark:bg-gray-700 text-gray-700 dark:text-gray-200 min-w-[200px] font-semibold"
                        >
                            {carriers.length === 0 && <option value="">Sin transportadoras</option>}
                            {carriers.map((c) => (
                                <option key={c.carrier} value={c.carrier}>{c.carrier} ({c.count})</option>
                            ))}
                        </select>
                    </div>
                    <div className="relative flex-1 min-w-[200px]">
                        <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-gray-400" />
                        <input
                            type="text"
                            value={search}
                            onChange={(e) => setSearch(e.target.value)}
                            placeholder="Buscar en pagina actual..."
                            className="w-full pl-9 pr-3 py-2 text-sm border border-gray-200 dark:border-gray-600 rounded-lg bg-gray-50 dark:bg-gray-700 text-gray-900 dark:text-white"
                        />
                    </div>
                    <div className="text-xs font-semibold text-gray-600 dark:text-gray-300">
                        Seleccionados: <span className="text-indigo-600 dark:text-indigo-300">{selected.size}</span>
                    </div>
                </div>

                <div className="flex-1 overflow-auto">
                    {loadingRows || loadingCarriers ? (
                        <div className="flex items-center justify-center py-16 text-gray-400 gap-2">
                            <RefreshCw size={18} className="animate-spin" /> Cargando...
                        </div>
                    ) : !businessId ? (
                        <div className="text-center py-16 text-gray-400 text-sm">Selecciona un negocio para ver los envios pendientes.</div>
                    ) : carriers.length === 0 ? (
                        <div className="text-center py-16 text-gray-400 text-sm">No hay envios pendientes de recoleccion.</div>
                    ) : filteredRows.length === 0 ? (
                        <div className="text-center py-16 text-gray-400 text-sm">Sin resultados para esta pagina.</div>
                    ) : (
                        <table className="w-full text-sm">
                            <thead className="bg-gray-50 dark:bg-gray-700/40 sticky top-0">
                                <tr className="text-[10px] uppercase tracking-wider text-gray-500 dark:text-gray-400">
                                    <th className="px-3 py-2 w-10 text-left">
                                        <input type="checkbox" checked={allPageSelected} onChange={togglePage} className="rounded" />
                                    </th>
                                    <th className="px-2 py-2 text-left font-bold">Guia</th>
                                    <th className="px-2 py-2 text-left font-bold">Orden</th>
                                    <th className="px-2 py-2 text-left font-bold">Cliente</th>
                                    <th className="px-2 py-2 text-left font-bold">Ciudad</th>
                                    <th className="px-2 py-2 text-left font-bold">F. Orden</th>
                                    <th className="px-2 py-2 text-left font-bold">F. Guia</th>
                                    <th className="px-2 py-2 text-left font-bold">Estados</th>
                                </tr>
                            </thead>
                            <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                                {filteredRows.map((s) => (
                                    <tr key={s.shipment_id} className="hover:bg-gray-50 dark:hover:bg-gray-700/30">
                                        <td className="px-3 py-2">
                                            <input type="checkbox" checked={selected.has(s.shipment_id)} onChange={() => toggle(s.shipment_id)} className="rounded" />
                                        </td>
                                        <td className="px-2 py-2 font-mono text-xs text-gray-500 dark:text-gray-400 truncate max-w-[140px]">{s.tracking_number || '—'}</td>
                                        <td className="px-2 py-2 text-purple-700 dark:text-purple-300 font-semibold text-xs truncate max-w-[110px]">{s.order_number}</td>
                                        <td className="px-2 py-2 text-gray-700 dark:text-gray-200 truncate max-w-[200px]">{s.customer_name || '—'}</td>
                                        <td className="px-2 py-2 text-gray-500 dark:text-gray-400 text-xs truncate max-w-[130px]">{s.destination_city}</td>
                                        <td className="px-2 py-2 text-[11px] text-gray-500 dark:text-gray-400 truncate">{formatDate(s.order_created_at)}</td>
                                        <td className="px-2 py-2 text-[11px] text-gray-500 dark:text-gray-400 truncate">{formatDate(s.shipment_created_at)}</td>
                                        <td className="px-2 py-2">
                                            <div className="flex flex-col gap-0.5">
                                                <span className="inline-block px-1.5 py-0.5 rounded text-[9px] font-semibold bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300 truncate">Guia: {s.shipment_status || '—'}</span>
                                                <span className="inline-block px-1.5 py-0.5 rounded text-[9px] font-semibold bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300 truncate">Ord: {s.order_status || '—'}</span>
                                            </div>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    )}
                </div>

                <div className="px-5 py-3 border-t border-gray-100 dark:border-gray-700 flex items-center justify-between gap-3 flex-wrap bg-gray-50 dark:bg-gray-800/60 rounded-b-2xl">
                    <div className="flex items-center gap-3 text-xs text-gray-600 dark:text-gray-300">
                        <label className="flex items-center gap-2">
                            Filas:
                            <select
                                value={pageSize}
                                onChange={(e) => setPageSize(Number(e.target.value))}
                                className="px-2 py-1 border border-gray-200 dark:border-gray-600 rounded bg-white dark:bg-gray-700 text-gray-700 dark:text-gray-200"
                            >
                                {PAGE_SIZE_OPTIONS.map((n) => <option key={n} value={n}>{n}</option>)}
                            </select>
                        </label>
                        <span>
                            Pagina <span className="font-semibold text-gray-800 dark:text-gray-100">{page}</span> de <span className="font-semibold text-gray-800 dark:text-gray-100">{totalPages}</span> · {total} envios
                        </span>
                        <div className="flex items-center gap-1">
                            <button
                                onClick={() => setPage((p) => Math.max(1, p - 1))}
                                disabled={page <= 1 || loadingRows}
                                className="p-1.5 rounded border border-gray-200 dark:border-gray-600 disabled:opacity-40 bg-white dark:bg-gray-700"
                            >
                                <ChevronLeft size={14} />
                            </button>
                            <button
                                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                                disabled={page >= totalPages || loadingRows}
                                className="p-1.5 rounded border border-gray-200 dark:border-gray-600 disabled:opacity-40 bg-white dark:bg-gray-700"
                            >
                                <ChevronRight size={14} />
                            </button>
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        <button onClick={onClose} className="px-4 py-2 text-sm font-semibold text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-lg">
                            Cancelar
                        </button>
                        <button
                            onClick={handleGenerate}
                            disabled={selected.size === 0 || generating}
                            className="px-4 py-2 text-sm font-semibold text-white bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 rounded-lg flex items-center gap-2"
                        >
                            {generating ? <RefreshCw size={14} className="animate-spin" /> : <Download size={14} />}
                            Generar PDF
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
