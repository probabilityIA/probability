'use client';

import { useState, useEffect, useCallback } from 'react';
import { getMovementsAction } from '../../infra/actions';
import { StockMovement, GetMovementsParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

interface StockMovementListProps {
    warehouseId?: number;
    selectedBusinessId?: number;
    onRefreshRef?: (ref: () => void) => void;
}

const DIRECTION_STYLES: Record<string, { bg: string; text: string; prefix: string }> = {
    in: { bg: 'bg-green-100', text: 'text-green-800', prefix: '+' },
    out: { bg: 'bg-red-100', text: 'text-red-800', prefix: '-' },
    reserve: { bg: 'bg-blue-100', text: 'text-blue-800', prefix: '' },
    release: { bg: 'bg-amber-100', text: 'text-amber-800', prefix: '' },
    confirm: { bg: 'bg-red-100', text: 'text-red-800', prefix: '-' },
    neutral: { bg: 'bg-gray-100 dark:bg-gray-700', text: 'text-gray-800 dark:text-gray-100', prefix: '' },
};

function getDirectionFromCode(code: string): string {
    const inCodes = ['inbound', 'return_stock', 'adjustment_in'];
    const outCodes = ['outbound', 'sale', 'adjustment_out'];
    if (inCodes.includes(code)) return 'in';
    if (outCodes.includes(code)) return 'out';
    if (code === 'reserve') return 'reserve';
    if (code === 'release') return 'release';
    if (code === 'confirm_sale') return 'confirm';
    return 'neutral';
}

interface ParsedNotes {
    reserved: number | null;
    availPrev: number | null;
    availNew: number | null;
    confirmed: number | null;
    rsvReleased: number | null;
}

function parseMovementNotes(notes: string): ParsedNotes {
    const rsvMatch = notes?.match(/Reservado:\s*(-?\d+)/i);
    const liberadoMatch = notes?.match(/Liberado:\s*(-?\d+)/i);
    const dispMatch = notes?.match(/Disponible:\s*(\d+)\s*->\s*(\d+)/i);
    const confirmedMatch = notes?.match(/Confirmado:\s*(\d+)/i);
    const rsvReleasedMatch = notes?.match(/Reserva liberada:\s*(\d+)/i);
    return {
        reserved: rsvMatch ? parseInt(rsvMatch[1], 10) : null,
        availPrev: dispMatch ? parseInt(dispMatch[1], 10) : null,
        availNew: dispMatch ? parseInt(dispMatch[2], 10) : null,
        confirmed: confirmedMatch ? parseInt(confirmedMatch[1], 10) : null,
        rsvReleased: rsvReleasedMatch ? parseInt(rsvReleasedMatch[1], 10) : (liberadoMatch ? parseInt(liberadoMatch[1], 10) : null),
    };
}

export default function StockMovementList({ warehouseId, selectedBusinessId, onRefreshRef }: StockMovementListProps) {
    const [movements, setMovements] = useState<StockMovement[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);

    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);

    const fetchMovements = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetMovementsParams = { page, page_size: pageSize };
            if (warehouseId) params.warehouse_id = warehouseId;
            if (selectedBusinessId) params.business_id = selectedBusinessId;

            const response = await getMovementsAction(params);
            setMovements(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar movimientos'));
        } finally {
            setLoading(false);
        }
    }, [warehouseId, page, pageSize, selectedBusinessId]);

    useEffect(() => {
        fetchMovements();
    }, [fetchMovements]);

    useEffect(() => {
        onRefreshRef?.(fetchMovements);
    }, [fetchMovements, onRefreshRef]);

    useEffect(() => {
        setPage(1);
    }, [warehouseId]);

    const columns = [
        { key: 'date', label: 'Fecha' },
        { key: 'product', label: 'Producto' },
        { key: 'type', label: 'Tipo', align: 'center' as const },
        { key: 'quantity', label: 'Cantidad', align: 'center' as const },
        { key: 'before_after', label: 'Stock Total', align: 'center' as const },
        { key: 'reserved_info', label: 'Reservado / Disponible', align: 'center' as const },
        { key: 'reason', label: 'Razón' },
    ];

    const renderRow = (movement: StockMovement) => {
        const direction = getDirectionFromCode(movement.movement_type_code);
        const style = DIRECTION_STYLES[direction] || DIRECTION_STYLES.neutral;
        const isReservation = movement.movement_type_code === 'reserve';
        const isRelease = movement.movement_type_code === 'release';
        const isConfirmSale = movement.movement_type_code === 'confirm_sale';
        const parsed = (isReservation || isRelease || isConfirmSale) ? parseMovementNotes(movement.notes || '') : null;

        return {
            date: (
                <span className="text-sm text-gray-500 dark:text-gray-400">
                    {new Date(movement.created_at).toLocaleDateString('es-CO', {
                        day: '2-digit',
                        month: 'short',
                        year: 'numeric',
                        hour: '2-digit',
                        minute: '2-digit',
                    })}
                </span>
            ),
            product: (
                <div>
                    <span className="font-medium text-gray-900 dark:text-white">{movement.product_name || movement.product_id}</span>
                    {movement.variant_label && (
                        <span className="block text-xs text-purple-600 dark:text-purple-400 font-medium">Variante: {movement.variant_label}</span>
                    )}
                    {movement.product_sku && (
                        <span className="block text-xs text-gray-500 dark:text-gray-400 font-mono">{movement.product_sku}</span>
                    )}
                </div>
            ),
            type: (
                <div className="flex justify-center">
                    <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${style.bg} ${style.text}`}>
                        {movement.movement_type_name || movement.movement_type_code}
                    </span>
                </div>
            ),
            quantity: (
                <span className={`text-sm font-semibold ${direction === 'in' ? 'text-green-700' : direction === 'out' || direction === 'confirm' ? 'text-red-700' : 'text-gray-700 dark:text-gray-200'}`}>
                    {isReservation || isRelease ? (
                        <span className="text-gray-400 dark:text-gray-500 text-xs">—</span>
                    ) : (
                        <>{style.prefix}{Math.abs(movement.quantity)}</>
                    )}
                </span>
            ),
            before_after: (
                <span className="text-xs text-gray-500 dark:text-gray-400">
                    {movement.previous_qty} &rarr; {movement.new_qty}
                </span>
            ),
            reserved_info: (() => {
                if (!parsed) return <span className="text-gray-300 dark:text-gray-600 text-xs">—</span>;

                if (isConfirmSale && parsed.rsvReleased !== null) {
                    return (
                        <div className="flex flex-col items-center gap-0.5">
                            <span className="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-semibold bg-red-100 text-red-800">
                                -{parsed.rsvReleased} rsv
                            </span>
                            <span className="text-xs text-gray-400 dark:text-gray-500">liberado</span>
                        </div>
                    );
                }

                if ((isReservation || isRelease) && (parsed.reserved !== null || parsed.availPrev !== null)) {
                    return (
                        <div className="flex flex-col items-center gap-0.5">
                            {parsed.reserved !== null && (
                                <span className={`inline-flex items-center px-1.5 py-0.5 rounded text-xs font-semibold ${isReservation ? 'bg-blue-100 text-blue-800' : 'bg-amber-100 text-amber-800'}`}>
                                    {isReservation ? '+' : '-'}{Math.abs(parsed.reserved)} rsv
                                </span>
                            )}
                            {parsed.availPrev !== null && parsed.availNew !== null && (
                                <span className="text-xs text-gray-500 dark:text-gray-400">
                                    disp: {parsed.availPrev} &rarr; {parsed.availNew}
                                </span>
                            )}
                        </div>
                    );
                }

                return <span className="text-gray-300 dark:text-gray-600 text-xs">—</span>;
            })(),
            reason: (
                <div>
                    <span className="text-sm text-gray-700 dark:text-gray-200">{movement.reason}</span>
                </div>
            ),
        };
    };

    if (loading && movements.length === 0) {
        return (
            <div className="flex justify-center items-center p-8">
                <Spinner size="lg" />
            </div>
        );
    }

    return (
        <div className="space-y-4">
            {error && (
                <Alert type="error" onClose={() => setError(null)}>
                    {error}
                </Alert>
            )}

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                <Table
                    columns={columns}
                    data={movements.map(renderRow)}
                    keyExtractor={(_, index) => String(movements[index]?.id || index)}
                    emptyMessage="No hay movimientos registrados"
                    loading={loading}
                    pagination={{
                        currentPage: page,
                        totalPages: totalPages,
                        totalItems: total,
                        itemsPerPage: pageSize,
                        onPageChange: (newPage) => setPage(newPage),
                        onItemsPerPageChange: (newSize) => {
                            setPageSize(newSize);
                            setPage(1);
                        },
                    }}
                />
            </div>
        </div>
    );
}
