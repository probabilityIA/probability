'use client';

import { useState, useEffect, useCallback } from 'react';
import {
    getEcommerceIntegrationTypesAction,
    getChannelStatusesAction,
    createChannelStatusAction,
    updateChannelStatusAction,
    deleteChannelStatusAction,
} from '../../infra/actions';
import { EcommerceIntegrationType, ChannelStatusInfo, CreateChannelStatusDTO, UpdateChannelStatusDTO } from '../../domain/types';
import { Spinner } from '@/shared/ui';

interface ChannelStatusManagerProps {
    isOpen: boolean;
    onClose: () => void;
}

type View = 'list' | 'create' | 'edit';

export default function ChannelStatusManager({ isOpen, onClose }: ChannelStatusManagerProps) {
    const [integrationTypes, setIntegrationTypes] = useState<EcommerceIntegrationType[]>([]);
    const [selectedType, setSelectedType] = useState<EcommerceIntegrationType | null>(null);
    const [statuses, setStatuses] = useState<ChannelStatusInfo[]>([]);
    const [loadingTypes, setLoadingTypes] = useState(false);
    const [loadingStatuses, setLoadingStatuses] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [view, setView] = useState<View>('list');
    const [editingStatus, setEditingStatus] = useState<ChannelStatusInfo | null>(null);
    const [deleteConfirm, setDeleteConfirm] = useState<ChannelStatusInfo | null>(null);
    const [deleteLoading, setDeleteLoading] = useState(false);
    const [deleteError, setDeleteError] = useState<string | null>(null);

    const loadIntegrationTypes = useCallback(() => {
        setLoadingTypes(true);
        setError(null);
        getEcommerceIntegrationTypesAction()
            .then((res) => {
                setIntegrationTypes(res.data || []);
                if (res.data && res.data.length > 0 && !selectedType) {
                    setSelectedType(res.data[0]);
                }
            })
            .catch((err) => setError(err.message || 'Error al cargar integraciones'))
            .finally(() => setLoadingTypes(false));
    }, []);

    const loadStatuses = useCallback((integrationTypeId: number) => {
        setLoadingStatuses(true);
        setError(null);
        getChannelStatusesAction(integrationTypeId)
            .then((res) => setStatuses(res.data || []))
            .catch((err) => setError(err.message || 'Error al cargar estados'))
            .finally(() => setLoadingStatuses(false));
    }, []);

    useEffect(() => {
        if (!isOpen) {
            setView('list');
            setEditingStatus(null);
            setDeleteConfirm(null);
            setDeleteError(null);
            return;
        }
        loadIntegrationTypes();
    }, [isOpen, loadIntegrationTypes]);

    useEffect(() => {
        if (selectedType) {
            loadStatuses(selectedType.id);
            setView('list');
            setEditingStatus(null);
        }
    }, [selectedType, loadStatuses]);

    const handleFormSuccess = () => {
        setView('list');
        setEditingStatus(null);
        if (selectedType) loadStatuses(selectedType.id);
    };

    const handleEdit = (status: ChannelStatusInfo) => {
        setEditingStatus(status);
        setView('edit');
    };

    const handleDelete = async () => {
        if (!deleteConfirm) return;
        setDeleteLoading(true);
        setDeleteError(null);
        try {
            const res = await deleteChannelStatusAction(deleteConfirm.id);
            if (!res.success) {
                setDeleteError(res.message || 'Error al eliminar');
                return;
            }
            setDeleteConfirm(null);
            if (selectedType) loadStatuses(selectedType.id);
        } catch (err: any) {
            setDeleteError(err.message || 'Error al eliminar');
        } finally {
            setDeleteLoading(false);
        }
    };

    if (!isOpen) return null;

    const title = view === 'create'
        ? 'Nuevo Estado del Canal'
        : view === 'edit'
            ? 'Editar Estado del Canal'
            : 'Estados por Canal de Integración';

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
                                    Estados nativos reportados por cada plataforma de ecommerce
                                </p>
                            )}
                        </div>
                    </div>
                    <div className="flex items-center gap-2">
                        {view === 'list' && selectedType && (
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
                    {/* Form view */}
                    {(view === 'create' || view === 'edit') && selectedType && (
                        <ChannelStatusForm
                            integrationTypeId={selectedType.id}
                            integrationTypeName={selectedType.name}
                            status={view === 'edit' ? editingStatus ?? undefined : undefined}
                            onSuccess={handleFormSuccess}
                            onCancel={() => { setView('list'); setEditingStatus(null); }}
                        />
                    )}

                    {/* List view */}
                    {view === 'list' && (
                        <>
                            {/* Integration type tabs */}
                            {loadingTypes ? (
                                <div className="flex justify-center py-6"><Spinner size="md" /></div>
                            ) : (
                                <div className="flex flex-wrap gap-2 mb-6">
                                    {integrationTypes.map((it) => (
                                        <button
                                            key={it.id}
                                            onClick={() => setSelectedType(it)}
                                            className={`flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium border transition-colors ${
                                                selectedType?.id === it.id
                                                    ? 'bg-blue-600 text-white border-blue-600'
                                                    : 'bg-white text-gray-600 border-gray-200 hover:border-blue-300 hover:text-blue-600'
                                            }`}
                                        >
                                            {it.image_url && (
                                                <img src={it.image_url} alt={it.name} className="w-4 h-4 object-contain" />
                                            )}
                                            {it.name}
                                        </button>
                                    ))}
                                    {integrationTypes.length === 0 && !loadingTypes && (
                                        <p className="text-sm text-gray-400">No tienes integraciones ecommerce configuradas.</p>
                                    )}
                                </div>
                            )}

                            {error && (
                                <div className="bg-red-50 text-red-700 px-4 py-3 rounded-lg text-sm mb-4">
                                    {error}
                                </div>
                            )}

                            {selectedType && (
                                <>
                                    {loadingStatuses ? (
                                        <div className="flex justify-center py-8"><Spinner size="lg" /></div>
                                    ) : (
                                        <table className="w-full">
                                            <thead>
                                                <tr className="text-left border-b border-gray-100">
                                                    <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider w-10">ID</th>
                                                    <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider w-10">Ord.</th>
                                                    <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Nombre</th>
                                                    <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Código</th>
                                                    <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider">Descripción</th>
                                                    <th className="pb-3 pr-4 text-xs font-semibold text-gray-400 uppercase tracking-wider w-16">Activo</th>
                                                    <th className="pb-3 text-xs font-semibold text-gray-400 uppercase tracking-wider w-20">Acciones</th>
                                                </tr>
                                            </thead>
                                            <tbody className="divide-y divide-gray-50">
                                                {statuses.map((status) => (
                                                    <ChannelStatusRow
                                                        key={status.id}
                                                        status={status}
                                                        onEdit={() => handleEdit(status)}
                                                        onDelete={() => { setDeleteConfirm(status); setDeleteError(null); }}
                                                    />
                                                ))}
                                                {statuses.length === 0 && (
                                                    <tr>
                                                        <td colSpan={7} className="py-8 text-center text-sm text-gray-400">
                                                            No hay estados registrados para {selectedType.name}
                                                        </td>
                                                    </tr>
                                                )}
                                            </tbody>
                                        </table>
                                    )}
                                </>
                            )}
                        </>
                    )}
                </div>

                {/* Footer */}
                {view === 'list' && selectedType && (
                    <div className="px-6 py-3 border-t border-gray-100 bg-gray-50 rounded-b-xl">
                        <p className="text-xs text-gray-400">
                            {statuses.length} estados registrados para {selectedType.name}
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

// ─────────────────────────────────────────────────────────────
// ChannelStatusRow
// ─────────────────────────────────────────────────────────────

function ChannelStatusRow({
    status,
    onEdit,
    onDelete,
}: {
    status: ChannelStatusInfo;
    onEdit: () => void;
    onDelete: () => void;
}) {
    return (
        <tr className="hover:bg-gray-50 transition-colors">
            <td className="py-3 pr-4 text-sm font-mono text-gray-400">{status.id}</td>
            <td className="py-3 pr-4">
                <span className="inline-flex items-center justify-center w-6 h-6 rounded-full bg-gray-100 text-xs font-bold text-gray-600">
                    {status.display_order}
                </span>
            </td>
            <td className="py-3 pr-4 text-sm font-medium text-gray-800">{status.name}</td>
            <td className="py-3 pr-4 text-sm font-mono text-gray-600">{status.code}</td>
            <td className="py-3 pr-4 text-sm text-gray-500 max-w-xs">
                {status.description || <span className="text-gray-300">—</span>}
            </td>
            <td className="py-3 pr-4">
                <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${
                    status.is_active
                        ? 'bg-green-100 text-green-700'
                        : 'bg-gray-100 text-gray-500'
                }`}>
                    {status.is_active ? 'Sí' : 'No'}
                </span>
            </td>
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

// ─────────────────────────────────────────────────────────────
// ChannelStatusForm
// ─────────────────────────────────────────────────────────────

interface ChannelStatusFormProps {
    integrationTypeId: number;
    integrationTypeName: string;
    status?: ChannelStatusInfo;
    onSuccess: () => void;
    onCancel: () => void;
}

function ChannelStatusForm({ integrationTypeId, integrationTypeName, status, onSuccess, onCancel }: ChannelStatusFormProps) {
    const isEdit = !!status;
    const [form, setForm] = useState({
        code: status?.code ?? '',
        name: status?.name ?? '',
        description: status?.description ?? '',
        is_active: status?.is_active ?? true,
        display_order: status?.display_order ?? 0,
    });
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setLoading(true);
        setError(null);
        try {
            if (isEdit && status) {
                const dto: UpdateChannelStatusDTO = {
                    code: form.code,
                    name: form.name,
                    description: form.description || undefined,
                    is_active: form.is_active,
                    display_order: Number(form.display_order),
                };
                const res = await updateChannelStatusAction(status.id, dto);
                if (!res.success) throw new Error(res.message || 'Error al actualizar');
            } else {
                const dto: CreateChannelStatusDTO = {
                    integration_type_id: integrationTypeId,
                    code: form.code,
                    name: form.name,
                    description: form.description || undefined,
                    is_active: form.is_active,
                    display_order: Number(form.display_order),
                };
                const res = await createChannelStatusAction(dto);
                if (!res.success) throw new Error(res.message || 'Error al crear');
            }
            onSuccess();
        } catch (err: any) {
            setError(err.message || 'Error al guardar');
        } finally {
            setLoading(false);
        }
    };

    return (
        <form onSubmit={handleSubmit} className="space-y-4">
            <div className="bg-blue-50 text-blue-700 px-4 py-2 rounded-lg text-sm font-medium">
                Canal: {integrationTypeName}
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Código <span className="text-red-500">*</span>
                    </label>
                    <input
                        type="text"
                        value={form.code}
                        onChange={(e) => setForm(f => ({ ...f, code: e.target.value }))}
                        placeholder="pending, paid, cancelled…"
                        required
                        className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm font-mono focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                        Nombre <span className="text-red-500">*</span>
                    </label>
                    <input
                        type="text"
                        value={form.name}
                        onChange={(e) => setForm(f => ({ ...f, name: e.target.value }))}
                        placeholder="Pendiente, Pagado, Cancelado…"
                        required
                        className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
            </div>

            <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">Descripción</label>
                <textarea
                    value={form.description}
                    onChange={(e) => setForm(f => ({ ...f, description: e.target.value }))}
                    placeholder="Descripción opcional del estado…"
                    rows={2}
                    className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 resize-none"
                />
            </div>

            <div className="grid grid-cols-2 gap-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">Orden de visualización</label>
                    <input
                        type="number"
                        min={0}
                        value={form.display_order}
                        onChange={(e) => setForm(f => ({ ...f, display_order: Number(e.target.value) }))}
                        className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
                    />
                </div>
                <div className="flex items-end pb-2">
                    <label className="flex items-center gap-2 cursor-pointer select-none">
                        <input
                            type="checkbox"
                            checked={form.is_active}
                            onChange={(e) => setForm(f => ({ ...f, is_active: e.target.checked }))}
                            className="w-4 h-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                        />
                        <span className="text-sm font-medium text-gray-700">Activo</span>
                    </label>
                </div>
            </div>

            {error && (
                <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-lg text-sm">
                    {error}
                </div>
            )}

            <div className="flex justify-end gap-3 pt-2">
                <button
                    type="button"
                    onClick={onCancel}
                    disabled={loading}
                    className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50"
                >
                    Cancelar
                </button>
                <button
                    type="submit"
                    disabled={loading}
                    className="px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded-lg hover:bg-blue-700 disabled:opacity-50"
                >
                    {loading ? 'Guardando…' : isEdit ? 'Actualizar' : 'Crear'}
                </button>
            </div>
        </form>
    );
}
