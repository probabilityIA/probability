'use client';

import { useState, useEffect, useCallback } from 'react';
import { getMovementsAction } from '../../infra/actions';
import { StockMovement, GetMovementsParams } from '../../domain/types';
import { Spinner } from '@/shared/ui';

interface Props {
    productId?: string;
    warehouseId?: number;
    businessId?: number;
}

const PAGE_SIZE = 15;

const DIRECTION_STYLES: Record<string, { bg: string; text: string; prefix: string }> = {
    in: { bg: 'bg-green-100', text: 'text-green-800', prefix: '+' },
    out: { bg: 'bg-red-100', text: 'text-red-800', prefix: '-' },
    reserve: { bg: 'bg-blue-100', text: 'text-blue-800', prefix: '' },
    release: { bg: 'bg-amber-100', text: 'text-amber-800', prefix: '' },
    confirm: { bg: 'bg-red-100', text: 'text-red-800', prefix: '-' },
    neutral: { bg: 'bg-gray-100', text: 'text-gray-800', prefix: '' },
};

function getDirectionFromCode(code: string): string {
    if (['inbound', 'return_stock', 'adjustment_in'].includes(code)) return 'in';
    if (['outbound', 'sale', 'adjustment_out'].includes(code)) return 'out';
    if (code === 'reserve') return 'reserve';
    if (code === 'release') return 'release';
    if (code === 'confirm_sale') return 'confirm';
    return 'neutral';
}

function parseNotes(notes: string) {
    const rsvMatch = notes?.match(/Reservado:\s*(-?\d+)/i);
    const liberadoMatch = notes?.match(/Liberado:\s*(-?\d+)/i);
    const dispMatch = notes?.match(/Disponible:\s*(\d+)\s*->\s*(\d+)/i);
    const confirmedMatch = notes?.match(/Confirmado:\s*(\d+)/i);
    const rsvReleasedMatch = notes?.match(/Reserva liberada:\s*(\d+)/i);
    return {
        reserved: rsvMatch ? parseInt(rsvMatch[1], 10) : null,
        availPrev: dispMatch ? parseInt(dispMatch[1], 10) : null,
        availNew: dispMatch ? parseInt(dispMatch[2], 10) : null,
        rsvReleased: rsvReleasedMatch ? parseInt(rsvReleasedMatch[1], 10) : (liberadoMatch ? parseInt(liberadoMatch[1], 10) : null),
        confirmed: confirmedMatch ? parseInt(confirmedMatch[1], 10) : null,
    };
}

function MovementRow({ m, showProduct, showWarehouse }: { m: StockMovement; showProduct: boolean; showWarehouse: boolean }) {
    const direction = getDirectionFromCode(m.movement_type_code);
    const style = DIRECTION_STYLES[direction] || DIRECTION_STYLES.neutral;
    const isReservation = m.movement_type_code === 'reserve';
    const isRelease = m.movement_type_code === 'release';
    const isConfirmSale = m.movement_type_code === 'confirm_sale';
    const parsed = (isReservation || isRelease || isConfirmSale) ? parseNotes(m.notes || '') : null;

    return (
        <tr className="hover:bg-gray-50 transition-colors">
            <td className="px-3 py-2 text-xs text-gray-500 whitespace-nowrap">
                {new Date(m.created_at).toLocaleDateString('es-CO', {
                    day: '2-digit', month: 'short', year: 'numeric',
                    hour: '2-digit', minute: '2-digit',
                })}
            </td>
            {showProduct && (
                <td className="px-3 py-2">
                    <span className="text-sm font-medium text-gray-900">{m.product_name || m.product_id}</span>
                    {m.product_sku && (
                        <span className="block text-xs text-gray-400 font-mono">{m.product_sku}</span>
                    )}
                </td>
            )}
            {showWarehouse && (
                <td className="px-3 py-2 text-sm text-gray-600">{m.warehouse_name || '-'}</td>
            )}
            <td className="px-3 py-2">
                <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
                    {m.movement_type_name || m.movement_type_code}
                </span>
            </td>
            <td className="px-3 py-2 text-center">
                <span className={`text-sm font-semibold ${direction === 'in' ? 'text-green-700' : direction === 'out' || direction === 'confirm' ? 'text-red-700' : 'text-gray-600'}`}>
                    {isReservation || isRelease ? (
                        <span className="text-gray-300 text-xs">&mdash;</span>
                    ) : (
                        <>{style.prefix}{Math.abs(m.quantity)}</>
                    )}
                </span>
            </td>
            <td className="px-3 py-2 text-center text-xs text-gray-500">
                {m.previous_qty} &rarr; {m.new_qty}
            </td>
            <td className="px-3 py-2 text-center text-xs">
                {!parsed ? (
                    <span className="text-gray-300">&mdash;</span>
                ) : isConfirmSale && parsed.rsvReleased !== null ? (
                    <span className="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-semibold bg-red-100 text-red-800">
                        -{parsed.rsvReleased} rsv
                    </span>
                ) : (isReservation || isRelease) && parsed.reserved !== null ? (
                    <div className="flex flex-col items-center gap-0.5">
                        <span className={`inline-flex items-center px-1.5 py-0.5 rounded text-xs font-semibold ${isReservation ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'}`}>
                            {isReservation ? '+' : '-'}{Math.abs(parsed.reserved)} rsv
                        </span>
                        {parsed.availPrev !== null && (
                            <span className="text-gray-400">disp: {parsed.availPrev} &rarr; {parsed.availNew}</span>
                        )}
                    </div>
                ) : (
                    <span className="text-gray-300">&mdash;</span>
                )}
            </td>
            <td className="px-3 py-2 text-sm text-gray-600 max-w-[160px] truncate">{m.reason}</td>
        </tr>
    );
}

export default function MovementsInlineTable({ productId, warehouseId, businessId }: Props) {
    const [movements, setMovements] = useState<StockMovement[]>([]);
    const [loading, setLoading] = useState(true);
    const [page, setPage] = useState(1);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const showProduct = !productId;
    const showWarehouse = !warehouseId;

    const load = useCallback(async () => {
        setLoading(true);
        try {
            const params: GetMovementsParams = { page, page_size: PAGE_SIZE };
            if (productId) params.product_id = productId;
            if (warehouseId) params.warehouse_id = warehouseId;
            if (businessId) params.business_id = businessId;
            const res = await getMovementsAction(params);
            setMovements(res.data ?? []);
            setTotal(res.total ?? 0);
            setTotalPages(res.total_pages ?? 1);
        } finally {
            setLoading(false);
        }
    }, [productId, warehouseId, businessId, page]);

    useEffect(() => { load(); }, [load]);
    useEffect(() => { setPage(1); }, [productId, warehouseId]);

    if (loading && movements.length === 0) {
        return (
            <div className="flex justify-center py-12">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="flex flex-col gap-3">
            <div className="overflow-x-auto">
                <table className="w-full text-left border-collapse">
                    <thead>
                        <tr className="border-b border-gray-200">
                            <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Fecha</th>
                            {showProduct && <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Producto</th>}
                            {showWarehouse && <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Bodega</th>}
                            <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Tipo</th>
                            <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50 text-center">Cantidad</th>
                            <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50 text-center">Stock Total</th>
                            <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50 text-center">Reservado</th>
                            <th className="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50">Razon</th>
                        </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-100">
                        {loading ? (
                            <tr><td colSpan={showProduct && showWarehouse ? 8 : 7} className="py-8 text-center"><Spinner size="sm" /></td></tr>
                        ) : movements.length === 0 ? (
                            <tr><td colSpan={showProduct && showWarehouse ? 8 : 7} className="py-10 text-center text-sm text-gray-400">Sin movimientos</td></tr>
                        ) : (
                            movements.map((m) => (
                                <MovementRow key={m.id} m={m} showProduct={showProduct} showWarehouse={showWarehouse} />
                            ))
                        )}
                    </tbody>
                </table>
            </div>

            {totalPages > 1 && (
                <div className="flex items-center justify-between px-2">
                    <span className="text-xs text-gray-500">{total} movimientos</span>
                    <div className="flex items-center gap-1">
                        <button
                            disabled={page <= 1}
                            onClick={() => setPage((p) => p - 1)}
                            className="px-2 py-1 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50"
                        >
                            Ant
                        </button>
                        <span className="text-xs text-gray-600 px-2">{page} / {totalPages}</span>
                        <button
                            disabled={page >= totalPages}
                            onClick={() => setPage((p) => p + 1)}
                            className="px-2 py-1 text-xs rounded border border-gray-300 disabled:opacity-40 hover:bg-gray-50"
                        >
                            Sig
                        </button>
                    </div>
                </div>
            )}
        </div>
    );
}
