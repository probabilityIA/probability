'use client';

import { useEffect, useState, useCallback } from 'react';
import { Modal, Spinner, Alert } from '@/shared/ui';
import { ProfitReportDetailResponse } from '../../domain/types';
import { shippingProfitReportDetailAction } from '../../infra/actions';
import { getActionError } from '@/shared/utils/action-result';

interface Props {
    isOpen: boolean;
    onClose: () => void;
    carrier: string;
    carrierLabel: string;
    from: string;
    to: string;
    selectedBusinessId?: number;
}

const fmt = (n: number) => '$ ' + Math.round(n).toLocaleString('es-CO');
const fmtDate = (s: string) => {
    try {
        return new Date(s).toLocaleString('es-CO', { dateStyle: 'short', timeStyle: 'short' });
    } catch {
        return s;
    }
};

export default function ShippingProfitDetailModal({ isOpen, onClose, carrier, carrierLabel, from, to, selectedBusinessId }: Props) {
    const [page, setPage] = useState(1);
    const [pageSize] = useState(20);
    const [data, setData] = useState<ProfitReportDetailResponse | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const fetchDetail = useCallback(async () => {
        if (!isOpen || !carrier) return;
        setLoading(true);
        setError(null);
        try {
            const r = await shippingProfitReportDetailAction({
                business_id: selectedBusinessId,
                carrier,
                from,
                to,
                page,
                page_size: pageSize,
            });
            setData(r);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar el detalle'));
        } finally {
            setLoading(false);
        }
    }, [isOpen, carrier, from, to, page, pageSize, selectedBusinessId]);

    useEffect(() => {
        fetchDetail();
    }, [fetchDetail]);

    useEffect(() => {
        if (isOpen) setPage(1);
    }, [isOpen, carrier, from, to]);

    const totalPages = data?.total_pages ?? 0;

    return (
        <Modal isOpen={isOpen} onClose={onClose} size="6xl" title={`Detalle de guias - ${carrierLabel}`}>
            <div className="space-y-4">
                <div className="text-sm text-gray-600 dark:text-gray-300">
                    Rango: <span className="font-medium">{from}</span> a <span className="font-medium">{to}</span>
                    {data && (
                        <span className="ml-3">
                            Total: <span className="font-semibold">{data.total}</span> guias
                        </span>
                    )}
                </div>

                {error && <Alert type="error">{error}</Alert>}

                <div className="bg-white dark:bg-gray-800 rounded-lg border border-gray-200 dark:border-gray-700 overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead className="bg-gray-100 dark:bg-gray-700">
                            <tr>
                                <th className="text-left px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Fecha</th>
                                <th className="text-left px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Orden</th>
                                <th className="text-left px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Guia</th>
                                <th className="text-left px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Estado</th>
                                <th className="text-right px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Cobrado</th>
                                <th className="text-right px-4 py-3 font-semibold text-gray-700 dark:text-gray-200">Costo carrier</th>
                                <th className="text-right px-4 py-3 font-semibold text-emerald-700 dark:text-emerald-300">Ganancia</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {loading ? (
                                <tr>
                                    <td colSpan={7} className="text-center py-10"><Spinner size="lg" /></td>
                                </tr>
                            ) : data?.data && data.data.length > 0 ? (
                                data.data.map((r) => (
                                    <tr key={r.shipment_id} className="hover:bg-gray-50 dark:hover:bg-gray-700/50">
                                        <td className="px-4 py-3 text-gray-700 dark:text-gray-200 whitespace-nowrap">{fmtDate(r.created_at)}</td>
                                        <td className="px-4 py-3 text-gray-900 dark:text-white font-medium">{r.order_number || '—'}</td>
                                        <td className="px-4 py-3 text-gray-700 dark:text-gray-200 font-mono text-xs">{r.tracking_number || '—'}</td>
                                        <td className="px-4 py-3 text-gray-700 dark:text-gray-200">{r.status}</td>
                                        <td className="px-4 py-3 text-right text-blue-700 dark:text-blue-300">{fmt(r.customer_charge)}</td>
                                        <td className="px-4 py-3 text-right text-orange-700 dark:text-orange-300">{fmt(r.carrier_cost)}</td>
                                        <td className="px-4 py-3 text-right font-semibold text-emerald-700 dark:text-emerald-300">{fmt(r.profit)}</td>
                                    </tr>
                                ))
                            ) : (
                                <tr>
                                    <td colSpan={7} className="text-center py-10 text-gray-400">Sin guias en el rango</td>
                                </tr>
                            )}
                        </tbody>
                    </table>
                </div>

                {data && totalPages > 1 && (
                    <div className="flex items-center justify-between text-sm">
                        <div className="text-gray-600 dark:text-gray-300">
                            Pagina <span className="font-semibold">{data.page}</span> de <span className="font-semibold">{totalPages}</span>
                        </div>
                        <div className="flex items-center gap-2">
                            <button
                                onClick={() => setPage((p) => Math.max(1, p - 1))}
                                disabled={page <= 1 || loading}
                                className="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md disabled:opacity-50 text-gray-700 dark:text-gray-200"
                            >
                                Anterior
                            </button>
                            <button
                                onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
                                disabled={page >= totalPages || loading}
                                className="px-3 py-1.5 border border-gray-300 dark:border-gray-600 rounded-md disabled:opacity-50 text-gray-700 dark:text-gray-200"
                            >
                                Siguiente
                            </button>
                        </div>
                    </div>
                )}
            </div>
        </Modal>
    );
}
