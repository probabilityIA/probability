'use client';

import { useState, useEffect, useCallback } from 'react';
import { getOrderStatusesAction, deleteOrderStatusAction } from '../../infra/actions';
import { OrderStatusInfo } from '../../domain/types';
import { getStatusBadgeStyle } from '@/shared/utils/color-utils';
import { Spinner } from '@/shared/ui';
import OrderStatusForm from './OrderStatusForm';

interface OrderStatusCatalogModalProps {
    isOpen: boolean;
    onClose: () => void;
}

type View = 'list' | 'create' | 'edit';

export default function OrderStatusCatalogModal({ isOpen, onClose }: OrderStatusCatalogModalProps) {
    const [statuses, setStatuses] = useState<OrderStatusInfo[]>([]);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [view, setView] = useState<View>('list');
    const [editingStatus, setEditingStatus] = useState<OrderStatusInfo | null>(null);
    const [deleteConfirm, setDeleteConfirm] = useState<OrderStatusInfo | null>(null);
    const [deleteLoading, setDeleteLoading] = useState(false);
    const [deleteError, setDeleteError] = useState<string | null>(null);

    const loadStatuses = useCallback(() => {
        setLoading(true);
        setError(null);
        getOrderStatusesAction()
            .then((res) => {
                const sorted = [...(res.data || [])].sort((a, b) => (a.priority ?? 0) - (b.priority ?? 0));
                setStatuses(sorted);
            })
            .catch((err) => setError(err.message || 'Error al cargar estados'))
            .finally(() => setLoading(false));
    }, []);

    useEffect(() => {
        if (!isOpen) {
            setView('list');
            setEditingStatus(null);
            setDeleteConfirm(null);
            setDeleteError(null);
            return;
        }
        loadStatuses();
    }, [isOpen, loadStatuses]);

    const handleFormSuccess = (saved: OrderStatusInfo) => {
        setView('list');
        setEditingStatus(null);
        loadStatuses();
    };

    const handleEdit = (status: OrderStatusInfo) => {
        setEditingStatus(status);
        setView('edit');
    };

    const handleDelete = async () => {
        if (!deleteConfirm) return;
        setDeleteLoading(true);
        setDeleteError(null);
        try {
            const res = await deleteOrderStatusAction(deleteConfirm.id);
            if (!res.success) {
                setDeleteError(res.message || 'Error al eliminar');
                return;
            }
            setDeleteConfirm(null);
            loadStatuses();
        } catch (err: any) {
            setDeleteError(err.message || 'Error al eliminar');
        } finally {
            setDeleteLoading(false);
        }
    };

    if (!isOpen) return null;

    const title = view === 'create'
        ? 'Nuevo Estado de Probability'
        : view === 'edit'
            ? 'Editar Estado de Probability'
            : 'Estados de Probability';

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center">
            {/* Backdrop */}
            <div className="absolute inset-0 bg-black/50" onClick={onClose} />

            {/* Modal */}
            <div className="relative bg-white rounded-xl shadow-2xl w-full max-w-4xl mx-4 max-h-[90vh] flex flex-col">
                {/* Header */}
                <div className="flex items-center justify-between px-6 py-4 border-b border-gray-200">
                    <div className="flex items-center gap-3">
                        {view !== 'list' && (
                            <button
                                onClick={() => { setView('list'); setEditingStatus(null); }}
                                className="p-1 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                                aria-label="Volver"
                            >
                                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                                </svg>
                            </button>
                        )}
                        <div>
                            <h2 className="text-xl font-semibold text-gray-900">{title}</h2>
                            {view === 'list' && (
                                <p className="text-sm text-gray-500 mt-0.5">
                                    Catálogo de estados internos — ordenados por prioridad de ciclo de vida
                                </p>
                            )}
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        {view === 'list' && (
                            <button
                                onClick={() => setView('create')}
                                className="flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 transition-colors"
                            >
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                                </svg>
                                Nuevo Estado
                            </button>
                        )}
                        <button
                            onClick={onClose}
                            className="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors"
                            aria-label="Cerrar"
                        >
                            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                        </button>
                    </div>
                </div>

                {/* Body */}
                <div className="flex-1 overflow-y-auto px-6 py-4">
                    {/* Create / Edit form */}
                    {(view === 'create' || view === 'edit') && (
                        <OrderStatusForm
                            status={view === 'edit' ? editingStatus ?? undefined : undefined}
                            onSuccess={handleFormSuccess}
                            onCancel={() => { setView('list'); setEditingStatus(null); }}
                        />
                    )}

                    {/* List */}
                    {view === 'list' && (
                        <>
                            {loading && (
                                <div className="flex justify-center items-center py-12">
                                    <Spinner size="lg" />
                                </div>
                            )}

                            {error && (
                                <div className="bg-red-50 text-red-700 px-4 py-3 rounded-lg text-sm">
                                    {error}
                                </div>
                            )}

                            {!loading && !error && statuses.length > 0 && (
                                <table className="w-full">
                                    <thead>
                                        <tr className="text-left border-b border-gray-100">
                                            <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider w-10">ID</th>
                                            <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider w-8">P.</th>
                                            <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Estado</th>
                                            <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Código</th>
                                            <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Categoría</th>
                                            <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Descripción</th>
                                            <th className="pb-3 text-xs font-semibold text-gray-400 uppercase tracking-wider w-20">Acciones</th>
                                        </tr>
                                    </thead>
                                    <tbody className="divide-y divide-gray-50">
                                        {statuses.map((status) => (
                                            <StatusRow
                                                key={status.id}
                                                status={status}
                                                onEdit={() => handleEdit(status)}
                                                onDelete={() => { setDeleteConfirm(status); setDeleteError(null); }}
                                            />
                                        ))}
                                    </tbody>
                                </table>
                            )}
                        </>
                    )}
                </div>

                {/* Footer */}
                {view === 'list' && (
                    <div className="px-6 py-3 border-t border-gray-100 bg-gray-50 rounded-b-xl">
                        <p className="text-xs text-gray-400">
                            {statuses.length} estados registrados · P. = Prioridad en el ciclo de vida
                        </p>
                    </div>
                )}
            </div>

            {/* Delete Confirm Modal */}
            {deleteConfirm && (
                <div className="absolute inset-0 z-60 flex items-center justify-center">
                    <div className="relative bg-white rounded-xl shadow-2xl w-full max-w-md mx-4 p-6">
                        <h3 className="text-lg font-semibold text-gray-900 mb-2">Eliminar estado</h3>
                        <p className="text-sm text-gray-600 mb-1">
                            ¿Estás seguro de que deseas eliminar el estado{' '}
                            <span className="font-semibold">{deleteConfirm.name}</span>?
                        </p>
                        <p className="text-xs text-gray-500 mb-4">
                            Esta acción es permanente. Si existen mapeos que dependan de este estado, no podrá eliminarse.
                        </p>

                        {deleteError && (
                            <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm mb-4">
                                {deleteError}
                            </div>
                        )}

                        <div className="flex justify-end gap-3">
                            <button
                                onClick={() => { setDeleteConfirm(null); setDeleteError(null); }}
                                disabled={deleteLoading}
                                className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50"
                            >
                                Cancelar
                            </button>
                            <button
                                onClick={handleDelete}
                                disabled={deleteLoading}
                                className="px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 disabled:opacity-50"
                            >
                                {deleteLoading ? 'Eliminando...' : 'Eliminar'}
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}

function StatusRow({
    status,
    onEdit,
    onDelete,
}: {
    status: OrderStatusInfo;
    onEdit: () => void;
    onDelete: () => void;
}) {
    const badgeStyle = getStatusBadgeStyle(status.color);

    return (
        <tr className="hover:bg-gray-50 transition-colors">
            {/* ID */}
            <td className="py-3 pr-4 text-sm font-mono text-gray-400">{status.id}</td>

            {/* Prioridad */}
            <td className="py-3 pr-4">
                <span className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-gray-100 text-xs font-bold text-gray-600">
                    {status.priority ?? 0}
                </span>
            </td>

            {/* Estado (badge con color) */}
            <td className="py-3 pr-4">
                <span
                    className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium"
                    style={badgeStyle}
                >
                    <span
                        className="w-1.5 h-1.5 rounded-full flex-shrink-0"
                        style={{ backgroundColor: status.color ?? '#9CA3AF' }}
                    />
                    {status.name}
                </span>
            </td>

            {/* Código */}
            <td className="py-3 pr-4 text-sm font-mono text-gray-600">{status.code}</td>

            {/* Categoría */}
            <td className="py-3 pr-4">
                <span className="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded font-medium">
                    {status.category || '—'}
                </span>
            </td>

            {/* Descripción */}
            <td className="py-3 pr-4 text-sm text-gray-500 max-w-xs">
                {status.description || <span className="text-gray-300">—</span>}
            </td>

            {/* Acciones */}
            <td className="py-3">
                <div className="flex items-center gap-1">
                    <button
                        onClick={onEdit}
                        className="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-colors"
                        title="Editar"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                    </button>
                    <button
                        onClick={onDelete}
                        className="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                        title="Eliminar"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2}
                                d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                    </button>
                </div>
            </td>
        </tr>
    );
}
