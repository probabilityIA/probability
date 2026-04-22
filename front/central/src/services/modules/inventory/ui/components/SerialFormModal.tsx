'use client';

import { useState } from 'react';
import { FormModal, Button, Alert } from '@/shared/ui';
import { createSerialAction, updateSerialAction } from '../../infra/actions/traceability';
import { InventorySerial } from '../../domain/traceability-types';

interface Props {
    businessId?: number;
    serial?: InventorySerial | null;
    onClose: () => void;
    onSuccess: () => void;
}

export default function SerialFormModal({ businessId, serial, onClose, onSuccess }: Props) {
    const isEdit = !!serial;
    const [form, setForm] = useState({
        product_id: serial?.product_id || '',
        serial_number: serial?.serial_number || '',
        lot_id: serial?.lot_id || '',
        location_id: serial?.current_location_id || '',
        state_code: 'available',
    });
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const change = (k: string, v: any) => setForm((f) => ({ ...f, [k]: v }));

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        setError(null);
        try {
            const data: any = {
                lot_id: form.lot_id ? Number(form.lot_id) : null,
                location_id: form.location_id ? Number(form.location_id) : null,
                state_code: form.state_code,
            };
            let result;
            if (isEdit) {
                result = await updateSerialAction(serial!.id, data, businessId);
            } else {
                result = await createSerialAction({ ...data, product_id: form.product_id, serial_number: form.serial_number }, businessId);
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
        <FormModal isOpen={true} onClose={onClose} title={isEdit ? 'Editar serie' : 'Crear serie'}>
            <form onSubmit={submit} className="p-6 space-y-4">
                {error && <Alert type="error">{error}</Alert>}

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU / ID) *</label>
                        <input required disabled={isEdit} value={form.product_id} onChange={(e) => change('product_id', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm disabled:opacity-60" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Número de serie *</label>
                        <input required disabled={isEdit} value={form.serial_number} onChange={(e) => change('serial_number', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm disabled:opacity-60" />
                    </div>
                </div>

                <div className="grid grid-cols-3 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Lote (ID)</label>
                        <input type="number" value={form.lot_id} onChange={(e) => change('lot_id', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Ubicación (ID)</label>
                        <input type="number" value={form.location_id} onChange={(e) => change('location_id', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estado</label>
                        <select value={form.state_code} onChange={(e) => change('state_code', e.target.value)} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="available">Disponible</option>
                            <option value="reserved">Reservado</option>
                            <option value="on_hold">En espera</option>
                            <option value="damaged">Averiado</option>
                            <option value="quarantine">Cuarentena</option>
                            <option value="expired">Vencido</option>
                            <option value="in_transit">En tránsito</option>
                            <option value="returned">Devuelto</option>
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
