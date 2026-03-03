'use client';

import { useState, useEffect, useCallback } from 'react';
import { getMovementsAction } from '../../infra/actions';
import { StockMovement, GetMovementsParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';

interface StockMovementListProps {
    warehouseId?: number;
    selectedBusinessId?: number;
    onRefreshRef?: (ref: () => void) => void;
}

const DIRECTION_STYLES: Record<string, { bg: string; text: string; prefix: string }> = {
    in: { bg: 'bg-green-100', text: 'text-green-800', prefix: '+' },
    out: { bg: 'bg-red-100', text: 'text-red-800', prefix: '-' },
    neutral: { bg: 'bg-gray-100', text: 'text-gray-800', prefix: '' },
};

function getDirectionFromCode(code: string): string {
    const inCodes = ['inbound', 'return', 'adjustment_in'];
    const outCodes = ['outbound', 'sale', 'adjustment_out'];
    if (inCodes.includes(code)) return 'in';
    if (outCodes.includes(code)) return 'out';
    return 'neutral';
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
            setError(err.message || 'Error al cargar movimientos');
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
        { key: 'before_after', label: 'Antes / Después', align: 'center' as const },
        { key: 'reason', label: 'Razón' },
    ];

    const renderRow = (movement: StockMovement) => {
        const direction = getDirectionFromCode(movement.movement_type_code);
        const style = DIRECTION_STYLES[direction] || DIRECTION_STYLES.neutral;

        return {
            date: (
                <span className="text-sm text-gray-500">
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
                    <span className="font-medium text-gray-900">{movement.product_name || movement.product_id}</span>
                    {movement.product_sku && (
                        <span className="block text-xs text-gray-500 font-mono">{movement.product_sku}</span>
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
                <span className={`text-sm font-semibold ${direction === 'in' ? 'text-green-700' : direction === 'out' ? 'text-red-700' : 'text-gray-700'}`}>
                    {style.prefix}{Math.abs(movement.quantity)}
                </span>
            ),
            before_after: (
                <span className="text-xs text-gray-500">
                    {movement.previous_qty} &rarr; {movement.new_qty}
                </span>
            ),
            reason: (
                <div>
                    <span className="text-sm text-gray-700">{movement.reason}</span>
                    {movement.notes && (
                        <span className="block text-xs text-gray-400 truncate max-w-[200px]">{movement.notes}</span>
                    )}
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

            <div className="bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
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
