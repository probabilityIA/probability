'use client';

import { useCallback, useEffect, useState } from 'react';
import { ChevronLeft, ChevronRight, Loader2 } from 'lucide-react';
import { DateRangePicker } from '@/shared/ui/date-range-picker';
import { getWalletSpendSummaryAction, listWalletSpendTransactionsAction, ConceptTotal, GuideStatusTotal } from '../../infra/actions';
import { CONCEPT_LABELS } from '../../domain/concept';

interface Props {
    businessId?: number;
}

interface SpendTransaction {
    ID: string;
    Amount: number;
    Type: string;
    Status: string;
    Concept: string;
    Reference: string;
    CreatedAt: string;
}

const CARD = 'bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl';

const CONCEPT_COLOR: Record<string, { bg: string; c: string }> = {
    GUIDE: { bg: '#eef2ff', c: '#4338ca' },
    SUBSCRIPTION: { bg: '#f3e8ff', c: '#7e22ce' },
    EXTRA_USAGE: { bg: '#fef3c7', c: '#b45309' },
    OTHER: { bg: '#f1f5f9', c: '#475569' },
};

const GUIDE_STATUS_LABELS: Record<string, string> = {
    pending: 'Pendiente',
    picked_up: 'Recolectada',
    in_transit: 'En tránsito',
    out_for_delivery: 'En reparto',
    delivered: 'Entregada',
    on_hold: 'Novedad',
    returned: 'Devuelta',
    failed: 'Fallida',
    cancelled: 'Cancelada',
    desconocido: 'Sin guía asociada',
};

const GUIDE_STATUS_BADGE: Record<string, string> = {
    pending: 'bg-amber-100 text-amber-700 dark:bg-amber-900/30 dark:text-amber-300',
    picked_up: 'bg-indigo-100 text-indigo-700 dark:bg-indigo-900/30 dark:text-indigo-300',
    in_transit: 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-300',
    out_for_delivery: 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-300',
    delivered: 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/30 dark:text-emerald-300',
    on_hold: 'bg-orange-100 text-orange-700 dark:bg-orange-900/30 dark:text-orange-300',
    returned: 'bg-rose-100 text-rose-700 dark:bg-rose-900/30 dark:text-rose-300',
    failed: 'bg-red-100 text-red-700 dark:bg-red-900/30 dark:text-red-300',
    cancelled: 'bg-gray-200 text-gray-700 dark:bg-gray-700 dark:text-gray-300',
    desconocido: 'bg-gray-100 text-gray-500 dark:bg-gray-700 dark:text-gray-400',
};

const guideStatusLabel = (status: string) => GUIDE_STATUS_LABELS[status] || status || 'Sin estado';
const guideStatusBadgeClass = (status: string) => GUIDE_STATUS_BADGE[status] || 'bg-gray-100 text-gray-600 dark:bg-gray-700 dark:text-gray-300';

function ConceptPill({ concept }: { concept: string }) {
    const m = CONCEPT_COLOR[concept] || CONCEPT_COLOR.OTHER;
    return (
        <span
            className="inline-flex items-center px-2 py-0.5 rounded-full text-[11px] font-semibold"
            style={{ backgroundColor: m.bg, color: m.c }}
        >
            {CONCEPT_LABELS[concept] || concept}
        </span>
    );
}

const formatCurrency = (amount: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(amount);

const formatDate = (dateStr?: string | null) => {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' });
};

const todayISO = () => new Date().toISOString().slice(0, 10);
const daysAgoISO = (days: number) => {
    const d = new Date();
    d.setDate(d.getDate() - days);
    return d.toISOString().slice(0, 10);
};

export default function BusinessWalletFinancialSummary({ businessId }: Props) {
    const [startDate, setStartDate] = useState<string | undefined>(daysAgoISO(29));
    const [endDate, setEndDate] = useState<string | undefined>(todayISO());
    const [byConcept, setByConcept] = useState<ConceptTotal[]>([]);
    const [guideStatusBreakdown, setGuideStatusBreakdown] = useState<GuideStatusTotal[]>([]);
    const [loadingSummary, setLoadingSummary] = useState(true);
    const [transactions, setTransactions] = useState<SpendTransaction[]>([]);
    const [loadingTx, setLoadingTx] = useState(true);
    const [conceptFilter, setConceptFilter] = useState('');
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(0);
    const pageSize = 15;

    const handleRangeChange = (s: string | undefined, e: string | undefined) => {
        setStartDate(s);
        setEndDate(e);
    };

    useEffect(() => {
        if (!startDate || !endDate) return;
        let cancelled = false;
        setLoadingSummary(true);
        getWalletSpendSummaryAction(startDate, endDate, businessId).then((res) => {
            if (cancelled) return;
            if (res.success && res.data) {
                setByConcept(res.data.by_concept || []);
                setGuideStatusBreakdown(res.data.guide_status_breakdown || []);
            }
            setLoadingSummary(false);
        });
        return () => { cancelled = true; };
    }, [businessId, startDate, endDate]);

    useEffect(() => { setPage(1); }, [startDate, endDate, conceptFilter]);

    const loadTransactions = useCallback(async () => {
        if (!startDate || !endDate) return;
        setLoadingTx(true);
        const res = await listWalletSpendTransactionsAction({
            startDate,
            endDate,
            concept: conceptFilter || undefined,
            page,
            pageSize,
            businessId,
        });
        if (res.success) {
            setTransactions(res.data || []);
            setTotal(res.total || 0);
            setTotalPages(res.total_pages || 0);
        }
        setLoadingTx(false);
    }, [businessId, startDate, endDate, conceptFilter, page]);

    useEffect(() => { loadTransactions(); }, [loadTransactions]);

    const guideTotal = byConcept.find((c) => c.concept === 'GUIDE');
    const subscriptionTotal = byConcept.find((c) => c.concept === 'SUBSCRIPTION');
    const grandTotal = byConcept.reduce((acc, c) => acc + c.amount, 0);

    return (
        <div className="space-y-4">
            <div className="flex items-center gap-3 flex-wrap">
                <label className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase tracking-wide">Rango de fechas</label>
                <DateRangePicker startDate={startDate} endDate={endDate} onChange={handleRangeChange} />
            </div>

            {loadingSummary ? (
                <div className="flex justify-center py-8"><Loader2 className="animate-spin text-gray-400" size={22} /></div>
            ) : (
                <>
                    <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                        <div className={`${CARD} p-4`}>
                            <p className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Gastado en guías</p>
                            <p className="text-lg font-bold text-gray-900 dark:text-white">{formatCurrency(guideTotal?.amount || 0)}</p>
                            <p className="text-[10px] text-gray-400 mt-0.5">{guideTotal?.count || 0} guías</p>
                        </div>
                        <div className={`${CARD} p-4`}>
                            <p className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Gastado en membresía</p>
                            <p className="text-lg font-bold text-gray-900 dark:text-white">{formatCurrency(subscriptionTotal?.amount || 0)}</p>
                            <p className="text-[10px] text-gray-400 mt-0.5">{subscriptionTotal?.count || 0} pagos</p>
                        </div>
                        <div className={`${CARD} p-4`}>
                            <p className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Total de guías</p>
                            <p className="text-lg font-bold text-gray-900 dark:text-white">{guideTotal?.count || 0}</p>
                        </div>
                        <div className={`${CARD} p-4`}>
                            <p className="text-[11px] text-gray-400 uppercase tracking-wide mb-1">Total gastado</p>
                            <p className="text-lg font-bold text-violet-600 dark:text-violet-400">{formatCurrency(grandTotal)}</p>
                        </div>
                    </div>

                    {guideStatusBreakdown.length > 0 && (
                        <div className={`${CARD} overflow-hidden`}>
                            <div className="px-4 py-3 border-b border-gray-100 dark:border-gray-700">
                                <h3 className="text-sm font-bold text-gray-900 dark:text-white">Guías por estado</h3>
                                <p className="text-[11px] text-gray-400 mt-0.5">De qué se compone el total de {guideTotal?.count || 0} guías pagadas</p>
                            </div>
                            <div className="overflow-x-auto">
                                <table className="w-full text-sm">
                                    <thead>
                                        <tr className="text-left text-[11px] text-gray-400 uppercase tracking-wide border-b border-gray-100 dark:border-gray-700">
                                            <th className="px-4 py-2 font-medium">Estado</th>
                                            <th className="px-4 py-2 font-medium text-right">Guías</th>
                                            <th className="px-4 py-2 font-medium text-right">Gastado</th>
                                        </tr>
                                    </thead>
                                    <tbody>
                                        {guideStatusBreakdown.map((row) => (
                                            <tr key={row.status} className="border-b border-gray-50 dark:border-gray-700/50">
                                                <td className="px-4 py-2">
                                                    <span className={`inline-block px-2 py-0.5 rounded-full text-[11px] font-semibold ${guideStatusBadgeClass(row.status)}`}>
                                                        {guideStatusLabel(row.status)}
                                                    </span>
                                                </td>
                                                <td className="px-4 py-2 text-right text-gray-600 dark:text-gray-300">{row.count}</td>
                                                <td className="px-4 py-2 text-right font-medium text-gray-900 dark:text-white">{formatCurrency(row.amount)}</td>
                                            </tr>
                                        ))}
                                    </tbody>
                                </table>
                            </div>
                        </div>
                    )}
                </>
            )}

            <div className={`${CARD} overflow-hidden`}>
                <div className="flex items-center justify-between px-4 py-3 border-b border-gray-100 dark:border-gray-700">
                    <h3 className="text-sm font-bold text-gray-900 dark:text-white">Movimientos</h3>
                    <select
                        value={conceptFilter}
                        onChange={(e) => setConceptFilter(e.target.value)}
                        className="text-xs border border-gray-300 dark:border-gray-600 rounded-lg px-2 py-1 bg-white dark:bg-gray-800 text-gray-700 dark:text-gray-200"
                    >
                        <option value="">Todos los conceptos</option>
                        {Object.entries(CONCEPT_LABELS).map(([value, label]) => (
                            <option key={value} value={value}>{label}</option>
                        ))}
                    </select>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="text-left text-[11px] text-gray-400 uppercase tracking-wide border-b border-gray-100 dark:border-gray-700">
                                <th className="px-4 py-2 font-medium">Fecha</th>
                                <th className="px-4 py-2 font-medium">Concepto</th>
                                <th className="px-4 py-2 font-medium">Referencia</th>
                                <th className="px-4 py-2 font-medium text-right">Monto</th>
                            </tr>
                        </thead>
                        <tbody>
                            {loadingTx ? (
                                <tr><td colSpan={4} className="text-center py-8"><Loader2 className="animate-spin text-gray-400 inline" size={18} /></td></tr>
                            ) : transactions.length === 0 ? (
                                <tr><td colSpan={4} className="text-center py-8 text-gray-400 text-sm">Sin movimientos para este filtro</td></tr>
                            ) : transactions.map((tx) => (
                                <tr key={tx.ID} className="border-b border-gray-50 dark:border-gray-700/50">
                                    <td className="px-4 py-2 text-gray-500 dark:text-gray-400">{formatDate(tx.CreatedAt)}</td>
                                    <td className="px-4 py-2"><ConceptPill concept={tx.Concept} /></td>
                                    <td className="px-4 py-2 text-gray-600 dark:text-gray-300 text-xs">{tx.Reference || '—'}</td>
                                    <td className="px-4 py-2 text-right font-medium text-gray-900 dark:text-white">{formatCurrency(tx.Amount)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>

                <div className="flex items-center justify-between px-4 py-3 border-t border-gray-100 dark:border-gray-700 bg-gray-50 dark:bg-gray-900/40">
                    <div className="text-xs text-gray-400 dark:text-gray-500 font-medium">
                        Mostrando <b className="text-gray-700 dark:text-gray-200">{transactions.length}</b> de {total} movimientos
                    </div>
                    {totalPages > 1 && (
                        <div className="flex items-center gap-1">
                            <span className="text-[12px] text-gray-500 mr-1">Pág {page} / {totalPages}</span>
                            <button
                                onClick={() => setPage((p) => Math.max(1, p - 1))}
                                disabled={page <= 1}
                                className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30"
                            >
                                <ChevronLeft size={15} />
                            </button>
                            <button
                                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                                disabled={page >= totalPages}
                                className="p-1 rounded hover:bg-gray-100 dark:hover:bg-gray-700 disabled:opacity-30"
                            >
                                <ChevronRight size={15} />
                            </button>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
}
