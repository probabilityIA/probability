'use client';

import { useState } from 'react';
import { FormModal, Button, Alert } from '@/shared/ui';
import { createLotAction, updateLotAction } from '../../infra/actions/traceability';
import { InventoryLot } from '../../domain/traceability-types';

interface Props {
    businessId?: number;
    lot?: InventoryLot | null;
    onClose: () => void;
    onSuccess: () => void;
}

export default function LotFormModal({ businessId, lot, onClose, onSuccess }: Props) {
    const isEdit = !!lot;
    const [form, setForm] = useState({
        product_id: lot?.product_id || '',
        lot_code: lot?.lot_code || '',
        manufacture_date: lot?.manufacture_date?.slice(0, 10) || '',
        expiration_date: lot?.expiration_date?.slice(0, 10) || '',
        received_at: lot?.received_at?.slice(0, 10) || '',
        supplier_id: lot?.supplier_id || '',
        status: lot?.status || 'active',
    });
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const change = (k: string, v: any) => setForm((f) => ({ ...f, [k]: v }));

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        setError(null);
        try {
            const toIso = (d: string) => (d ? `${d}T00:00:00Z` : null);
            const data: any = {
                lot_code: form.lot_code,
                manufacture_date: toIso(form.manufacture_date),
                expiration_date: toIso(form.expiration_date),
                received_at: toIso(form.received_at),
                supplier_id: form.supplier_id ? Number(form.supplier_id) : null,
                status: form.status,
            };
            let result;
            if (isEdit) {
                result = await updateLotAction(lot!.id, data, businessId);
            } else {
                result = await createLotAction({ ...data, product_id: form.product_id }, businessId);
            }
            if (!result.success) {
                setError(result.error || 'Error al guardar');
                return;
            }
            onSuccess();
        } catch (err: any) {
            setError(err.message);
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <FormModal isOpen={true} onClose={onClose} title={isEdit ? 'Editar lote' : 'Crear lote'}>
            <form onSubmit={submit} className="p-6 space-y-4">
                {error && <Alert type="error">{error}</Alert>}

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU / ID) *</label>
                        <input
                            required
                            disabled={isEdit}
                            value={form.product_id}
                            onChange={(e) => change('product_id', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm disabled:opacity-60"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Codigo de lote *</label>
                        <input
                            required
                            value={form.lot_code}
                            onChange={(e) => change('lot_code', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                        />
                    </div>
                </div>

                <div className="grid grid-cols-3 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Fabricacion</label>
                        <input type="date" value={form.manufacture_date} onChange={(e) => change('manufacture_date', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Vencimiento</label>
                        <input type="date" value={form.expiration_date} onChange={(e) => change('expiration_date', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Recibido</label>
                        <input type="date" value={form.received_at} onChange={(e) => change('received_at', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                </div>

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Proveedor (ID)</label>
                        <input type="number" value={form.supplier_id} onChange={(e) => change('supplier_id', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                        <select value={form.status} onChange={(e) => change('status', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="active">Activo</option>
                            <option value="expired">Vencido</option>
                            <option value="recalled">Retirado</option>
                            <option value="blocked">Bloqueado</option>
                        </select>
                    </div>
                </div>

                <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Button type="button" variant="outline" onClick={onClose} disabled={submitting}>Cancelar</Button>
                    <Button type="submit" variant="primary" disabled={submitting}>{submitting ? 'Guardando...' : 'Guardar'}</Button>
                </div>
            </form>
        </FormModal>
    );
}
