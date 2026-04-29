'use client';

import { useState, useEffect, useCallback } from 'react';
import { PencilIcon, TrashIcon, MinusIcon, PlusIcon } from '@heroicons/react/24/outline';
import { listShippingMarginsAction, deleteShippingMarginAction, updateShippingMarginAction } from '../../infra/actions';
import { ShippingMargin, GetShippingMarginsParams } from '../../domain/types';
import { Alert, Table, Spinner } from '@/shared/ui';
import { getActionError } from '@/shared/utils/action-result';

const STEP = 100;

interface InlineEditorProps {
    value: number;
    busy: boolean;
    onChange: (delta: number) => void;
}

function InlineMarginEditor({ value, busy, onChange }: InlineEditorProps) {
    return (
        <div className="inline-flex items-center gap-2 justify-center">
            <button
                type="button"
                disabled={busy || value <= 0}
                onClick={() => onChange(-STEP)}
                className="p-1 rounded-md bg-red-500 text-white hover:bg-red-600 dark:bg-red-600 dark:hover:bg-red-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors shadow-sm"
                title="Restar 100 COP"
            >
                <MinusIcon className="w-3.5 h-3.5" />
            </button>
            <span className={`font-mono text-sm tabular-nums min-w-[5.5rem] text-center ${busy ? 'opacity-50' : ''} text-gray-900 dark:text-white`}>
                {formatCOP(value)}
            </span>
            <button
                type="button"
                disabled={busy}
                onClick={() => onChange(STEP)}
                className="p-1 rounded-md bg-green-500 text-white hover:bg-green-600 dark:bg-green-600 dark:hover:bg-green-700 disabled:opacity-40 disabled:cursor-not-allowed transition-colors shadow-sm"
                title="Sumar 100 COP"
            >
                <PlusIcon className="w-3.5 h-3.5" />
            </button>
        </div>
    );
}

interface Props {
    onEdit?: (m: ShippingMargin) => void;
    onRefreshRef?: (ref: () => void) => void;
    selectedBusinessId?: number;
}

const formatCOP = (n: number) =>
    new Intl.NumberFormat('es-CO', { style: 'currency', currency: 'COP', maximumFractionDigits: 0 }).format(n);

export default function ShippingMarginList({ onEdit, onRefreshRef, selectedBusinessId }: Props) {
    const [items, setItems] = useState<ShippingMargin[]>([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [page, setPage] = useState(1);
    const [pageSize, setPageSize] = useState(20);
    const [total, setTotal] = useState(0);
    const [totalPages, setTotalPages] = useState(1);
    const [savingId, setSavingId] = useState<number | null>(null);

    const handleAdjust = async (m: ShippingMargin, field: 'margin_amount' | 'insurance_margin', delta: number) => {
        const current = m[field];
        const next = Math.max(0, current + delta);
        if (next === current) return;
        const previous = items;
        setItems((prev) => prev.map((it) => (it.id === m.id ? { ...it, [field]: next } : it)));
        setSavingId(m.id);
        try {
            await updateShippingMarginAction(
                m.id,
                {
                    carrier_name: m.carrier_name,
                    margin_amount: field === 'margin_amount' ? next : m.margin_amount,
                    insurance_margin: field === 'insurance_margin' ? next : m.insurance_margin,
                    is_active: m.is_active,
                },
                selectedBusinessId,
            );
        } catch (err: any) {
            setItems(previous);
            setError(getActionError(err, 'Error al actualizar el margen'));
        } finally {
            setSavingId(null);
        }
    };

    const fetchItems = useCallback(async () => {
        setLoading(true);
        setError(null);
        try {
            const params: GetShippingMarginsParams = { page, page_size: pageSize };
            if (selectedBusinessId) params.business_id = selectedBusinessId;
            const response = await listShippingMarginsAction(params);
            setItems(response.data || []);
            setTotal(response.total || 0);
            setTotalPages(response.total_pages || 1);
            setPage(response.page || page);
        } catch (err: any) {
            setError(getActionError(err, 'Error al cargar margenes'));
        } finally {
            setLoading(false);
        }
    }, [page, pageSize, selectedBusinessId]);

    useEffect(() => {
        fetchItems();
    }, [fetchItems]);

    useEffect(() => {
        onRefreshRef?.(fetchItems);
    }, [fetchItems, onRefreshRef]);

    useEffect(() => {
        setPage(1);
    }, [selectedBusinessId]);

    const handleDelete = async (m: ShippingMargin) => {
        if (!confirm(`Eliminar margen para "${m.carrier_name}"? Esta accion no se puede deshacer.`)) return;
        try {
            await deleteShippingMarginAction(m.id, selectedBusinessId);
            fetchItems();
        } catch (err: any) {
            setError(getActionError(err, 'Error al eliminar el margen'));
        }
    };

    const columns = [
        { key: 'carrier', label: 'Transportadora' },
        { key: 'margin', label: 'Margen flete', align: 'center' as const },
        { key: 'insurance', label: 'Margen seguro', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (m: ShippingMargin) => ({
        carrier: (
            <div>
                <span className="font-medium text-gray-900 dark:text-white">{m.carrier_name}</span>
                <span className="block text-xs text-gray-500 dark:text-gray-400">{m.carrier_code}</span>
            </div>
        ),
        margin: (
            <InlineMarginEditor
                value={m.margin_amount}
                busy={savingId === m.id}
                onChange={(delta) => handleAdjust(m, 'margin_amount', delta)}
            />
        ),
        insurance: (
            <InlineMarginEditor
                value={m.insurance_margin}
                busy={savingId === m.id}
                onChange={(delta) => handleAdjust(m, 'insurance_margin', delta)}
            />
        ),
        status: (
            <span
                className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    m.is_active
                        ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                        : 'bg-gray-100 text-gray-800 dark:bg-gray-700 dark:text-gray-200'
                }`}
            >
                {m.is_active ? 'Activo' : 'Inactivo'}
            </span>
        ),
        actions: (
            <div className="flex justify-end gap-2">
                {onEdit && (
                    <button
                        onClick={() => onEdit(m)}
                        className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md transition-colors"
                        title="Editar"
                    >
                        <PencilIcon className="w-4 h-4" />
                    </button>
                )}
                <button
                    onClick={() => handleDelete(m)}
                    className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md transition-colors"
                    title="Eliminar"
                >
                    <TrashIcon className="w-4 h-4" />
                </button>
            </div>
        ),
    });

    if (loading && items.length === 0) {
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

            <div className="bg-blue-50 dark:bg-blue-900/20 border border-blue-200 dark:border-blue-800 rounded-lg p-3 text-xs text-blue-800 dark:text-blue-200">
                Estos margenes se suman automaticamente al precio que cobra cada transportadora al generar guias.
                EnvioClick es el intermediario y no lleva margen propio: el margen aplica al carrier real (Servientrega, Interrapidisimo, etc.).
            </div>

            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                <Table
                    columns={columns}
                    data={items.map(renderRow)}
                    keyExtractor={(_, index) => String(items[index]?.id || index)}
                    emptyMessage="No hay margenes configurados. Agrega uno para empezar."
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
