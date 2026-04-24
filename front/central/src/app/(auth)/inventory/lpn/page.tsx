'use client';

import { useState } from 'react';
import { PlusIcon, ArrowsRightLeftIcon, ArchiveBoxXMarkIcon, PencilIcon, TrashIcon } from '@heroicons/react/24/outline';
import { Alert, Table, Spinner, Button } from '@/shared/ui';
import { ConfirmModal } from '@/shared/ui/confirm-modal';
import { FormModal } from '@/shared/ui/form-modal';
import { usePermissions } from '@/shared/contexts/permissions-context';
import { useInventoryBusiness } from '@/shared/contexts/inventory-business-context';
import { useLPNs } from '@/services/modules/inventory/ui/hooks/useCapture';
import { createLPNAction, updateLPNAction, deleteLPNAction, moveLPNAction, dissolveLPNAction, mergeLPNAction } from '@/services/modules/inventory/infra/actions/capture';
import { LicensePlate } from '@/services/modules/inventory/domain/capture-types';

export default function LPNPage() {
    const { isSuperAdmin } = usePermissions();
    const { selectedBusinessId } = useInventoryBusiness();
    const businessId = isSuperAdmin ? selectedBusinessId ?? undefined : undefined;

    const { lpns, loading, error, page, pageSize, total, totalPages, filters, setPage, setPageSize, setFilters, refresh } = useLPNs({ business_id: businessId });

    const [modal, setModal] = useState<null | { mode: 'create' } | { mode: 'edit'; lpn: LicensePlate } | { mode: 'move'; lpn: LicensePlate } | { mode: 'merge'; lpn: LicensePlate }>(null);
    const [deleting, setDeleting] = useState<LicensePlate | null>(null);
    const [dissolving, setDissolving] = useState<LicensePlate | null>(null);

    const [form, setForm] = useState({ code: '', lpn_type: 'pallet', location_id: '' });
    const [mergeTargetId, setMergeTargetId] = useState('');
    const [submitting, setSubmitting] = useState(false);
    const [formError, setFormError] = useState<string | null>(null);

    const requiresBusiness = isSuperAdmin && selectedBusinessId === null;

    const openCreate = () => { setForm({ code: '', lpn_type: 'pallet', location_id: '' }); setFormError(null); setModal({ mode: 'create' }); };
    const openEdit = (lpn: LicensePlate) => { setForm({ code: lpn.code, lpn_type: lpn.lpn_type, location_id: lpn.current_location_id?.toString() || '' }); setFormError(null); setModal({ mode: 'edit', lpn }); };
    const openMove = (lpn: LicensePlate) => { setForm({ code: lpn.code, lpn_type: lpn.lpn_type, location_id: '' }); setFormError(null); setModal({ mode: 'move', lpn }); };
    const openMerge = (lpn: LicensePlate) => { setMergeTargetId(''); setFormError(null); setModal({ mode: 'merge', lpn }); };

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!modal) return;
        setSubmitting(true);
        setFormError(null);
        try {
            let r;
            if (modal.mode === 'create') {
                r = await createLPNAction({ code: form.code, lpn_type: form.lpn_type, location_id: form.location_id ? Number(form.location_id) : null }, businessId);
            } else if (modal.mode === 'edit') {
                r = await updateLPNAction(modal.lpn.id, { code: form.code, lpn_type: form.lpn_type, location_id: form.location_id ? Number(form.location_id) : null }, businessId);
            } else if (modal.mode === 'move') {
                if (!form.location_id) { setFormError('Ubicación requerida'); return; }
                r = await moveLPNAction(modal.lpn.id, { new_location_id: Number(form.location_id) }, businessId);
            } else {
                if (!mergeTargetId) { setFormError('LPN destino requerida'); return; }
                r = await mergeLPNAction(modal.lpn.id, { target_lpn_id: Number(mergeTargetId) }, businessId);
            }
            if (!r.success) { setFormError(r.error || 'Error'); return; }
            setModal(null);
            refresh();
        } finally { setSubmitting(false); }
    };

    const handleDelete = async () => {
        if (!deleting) return;
        await deleteLPNAction(deleting.id, businessId);
        setDeleting(null);
        refresh();
    };

    const handleDissolve = async () => {
        if (!dissolving) return;
        await dissolveLPNAction(dissolving.id, businessId);
        setDissolving(null);
        refresh();
    };

    const typeStyles: Record<string, string> = {
        pallet: 'bg-purple-100 dark:bg-purple-900 text-purple-800 dark:text-purple-200',
        case: 'bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200',
        tote: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
    };
    const statusStyles: Record<string, string> = {
        active: 'bg-green-100 dark:bg-green-900 text-green-800 dark:text-green-200',
        in_transit: 'bg-yellow-100 dark:bg-yellow-900 text-yellow-800 dark:text-yellow-200',
        dissolved: 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300',
    };

    const columns = [
        { key: 'code', label: 'Código' },
        { key: 'type', label: 'Tipo', align: 'center' as const },
        { key: 'location', label: 'Ubicación', align: 'center' as const },
        { key: 'items', label: 'Líneas', align: 'center' as const },
        { key: 'status', label: 'Estado', align: 'center' as const },
        { key: 'actions', label: 'Acciones', align: 'right' as const },
    ];

    const renderRow = (l: LicensePlate) => ({
        code: <span className="font-mono font-medium text-gray-900 dark:text-white">{l.code}</span>,
        type: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${typeStyles[l.lpn_type] || 'bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300'}`}>{l.lpn_type}</span>,
        location: <span className="text-sm text-gray-600 dark:text-gray-300">{l.current_location_id ? `#${l.current_location_id}` : '—'}</span>,
        items: <span className="text-sm text-gray-600 dark:text-gray-300">{l.Lines?.length || 0}</span>,
        status: <span className={`inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium ${statusStyles[l.status] || statusStyles.dissolved}`}>{l.status}</span>,
        actions: (
            <div className="flex justify-end gap-1">
                {l.status === 'active' && (
                    <>
                        <button onClick={() => openMove(l)} className="p-2 bg-blue-500 hover:bg-blue-600 text-white rounded-md" title="Mover"><ArrowsRightLeftIcon className="w-4 h-4" /></button>
                        <button onClick={() => openMerge(l)} className="p-2 bg-purple-500 hover:bg-purple-600 text-white rounded-md" title="Merge">⇉</button>
                        <button onClick={() => setDissolving(l)} className="p-2 bg-orange-500 hover:bg-orange-600 text-white rounded-md" title="Disolver"><ArchiveBoxXMarkIcon className="w-4 h-4" /></button>
                    </>
                )}
                <button onClick={() => openEdit(l)} className="p-2 bg-yellow-500 hover:bg-yellow-600 text-white rounded-md"><PencilIcon className="w-4 h-4" /></button>
                <button onClick={() => setDeleting(l)} className="p-2 bg-red-500 hover:bg-red-600 text-white rounded-md"><TrashIcon className="w-4 h-4" /></button>
            </div>
        ),
    });

    return (
        <div className="min-h-screen bg-gray-50 dark:bg-gray-900 w-full px-4 sm:px-6 lg:px-8 py-4 sm:py-6 lg:py-8">
            <div className="space-y-4">
                <div className="flex items-center justify-between flex-wrap gap-3">                    {!requiresBusiness && (
                        <button onClick={openCreate} className="px-5 py-3 btn-business-primary text-white font-bold rounded-lg shadow-lg hover:shadow-2xl transition-all duration-300 transform hover:scale-110 hover:-translate-y-1 flex items-center gap-2">
                            <PlusIcon className="w-4 h-4" /> Nueva LPN
                        </button>
                    )}
                </div>

                {requiresBusiness ? <Alert type="info">Selecciona un negocio.</Alert> : (
                    <>
                        <div className="flex gap-3 items-end flex-wrap">
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Tipo</label>
                                <select value={filters.lpn_type || ''} onChange={(e) => setFilters({ lpn_type: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value="">Todos</option>
                                    <option value="pallet">Pallet</option>
                                    <option value="case">Caja</option>
                                    <option value="tote">Tote</option>
                                </select>
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                                <select value={filters.status || ''} onChange={(e) => setFilters({ status: e.target.value || undefined })} className="px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                    <option value="">Todos</option>
                                    <option value="active">Activas</option>
                                    <option value="in_transit">En tránsito</option>
                                    <option value="dissolved">Disueltas</option>
                                </select>
                            </div>
                        </div>

                        {error && <Alert type="error">{error}</Alert>}

                        {loading && lpns.length === 0 ? <div className="flex justify-center p-8"><Spinner size="lg" /></div> : (
                            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 overflow-hidden">
                                <Table columns={columns} data={lpns.map(renderRow)} keyExtractor={(_, i) => String(lpns[i]?.id || i)} emptyMessage="Sin LPNs" loading={loading} pagination={{ currentPage: page, totalPages, totalItems: total, itemsPerPage: pageSize, onPageChange: setPage, onItemsPerPageChange: setPageSize }} />
                            </div>
                        )}
                    </>
                )}
            </div>

            {modal && (
                <FormModal isOpen={true} onClose={() => setModal(null)} title={
                    modal.mode === 'create' ? 'Crear LPN' :
                    modal.mode === 'edit' ? `Editar LPN ${modal.lpn.code}` :
                    modal.mode === 'move' ? `Mover LPN ${modal.lpn.code}` :
                    `Merge LPN ${modal.lpn.code} → ...`
                }>
                    <form onSubmit={handleSubmit} className="p-6 space-y-4">
                        {formError && <Alert type="error">{formError}</Alert>}

                        {(modal.mode === 'create' || modal.mode === 'edit') && (
                            <>
                                <div className="grid grid-cols-2 gap-4">
                                    <div>
                                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Código *</label>
                                        <input required value={form.code} onChange={(e) => setForm({ ...form, code: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                                    </div>
                                    <div>
                                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Tipo</label>
                                        <select value={form.lpn_type} onChange={(e) => setForm({ ...form, lpn_type: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                                            <option value="pallet">Pallet</option>
                                            <option value="case">Caja</option>
                                            <option value="tote">Tote</option>
                                        </select>
                                    </div>
                                </div>
                                <div>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Ubicación inicial (ID)</label>
                                    <input type="number" value={form.location_id} onChange={(e) => setForm({ ...form, location_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                                </div>
                            </>
                        )}

                        {modal.mode === 'move' && (
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Nueva ubicación (ID) *</label>
                                <input required type="number" value={form.location_id} onChange={(e) => setForm({ ...form, location_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                            </div>
                        )}

                        {modal.mode === 'merge' && (
                            <div>
                                <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">LPN destino (ID) *</label>
                                <input required type="number" value={mergeTargetId} onChange={(e) => setMergeTargetId(e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                                <p className="text-xs text-gray-500 mt-1">Esta LPN (#{modal.lpn.id}) se disolverá y sus líneas pasarán a la destino.</p>
                            </div>
                        )}

                        <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                            <Button type="button" variant="outline" onClick={() => setModal(null)}>Cancelar</Button>
                            <Button type="submit" variant="primary" disabled={submitting}>{submitting ? 'Guardando...' : 'Guardar'}</Button>
                        </div>
                    </form>
                </FormModal>
            )}

            {deleting && <ConfirmModal isOpen={true} onClose={() => setDeleting(null)} onConfirm={handleDelete} title="Eliminar LPN" message={`Se eliminará la LPN ${deleting.code}. Acción irreversible.`} confirmText="Eliminar" type="danger" />}
            {dissolving && <ConfirmModal isOpen={true} onClose={() => setDissolving(null)} onConfirm={handleDissolve} title="Disolver LPN" message={`Se disolverá la LPN ${dissolving.code}. Los productos quedarán sueltos en su ubicación actual.`} confirmText="Disolver" type="warning" />}
        </div>
    );
}
