'use client';

import { useState } from 'react';
import { FormModal, Button, Alert } from '@/shared/ui';
import { createCountPlanAction, updateCountPlanAction } from '../../infra/actions/audit';
import { CycleCountPlan } from '../../domain/audit-types';

interface Props {
    businessId?: number;
    plan?: CycleCountPlan | null;
    onClose: () => void;
    onSuccess: () => void;
}

export default function CountPlanFormModal({ businessId, plan, onClose, onSuccess }: Props) {
    const isEdit = !!plan;
    const [form, setForm] = useState({
        warehouse_id: plan?.warehouse_id || 0,
        name: plan?.name || '',
        strategy: plan?.strategy || 'abc',
        frequency_days: plan?.frequency_days ?? 30,
        next_run_at: plan?.next_run_at?.slice(0, 10) || '',
        is_active: plan?.is_active ?? true,
    });
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        setError(null);
        try {
            const data: any = {
                warehouse_id: Number(form.warehouse_id),
                name: form.name,
                strategy: form.strategy,
                frequency_days: Number(form.frequency_days),
                next_run_at: form.next_run_at ? `${form.next_run_at}T00:00:00Z` : null,
                is_active: form.is_active,
            };
            const r = isEdit ? await updateCountPlanAction(plan!.id, data, businessId) : await createCountPlanAction(data, businessId);
            if (!r.success) { setError(r.error || 'Error al guardar'); return; }
            onSuccess();
        } catch (err: any) { setError(err.message); } finally { setSubmitting(false); }
    };

    return (
        <FormModal isOpen={true} onClose={onClose} title={isEdit ? 'Editar plan de conteo' : 'Crear plan de conteo'}>
            <form onSubmit={submit} className="p-6 space-y-4">
                {error && <Alert type="error">{error}</Alert>}

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Bodega (ID) *</label>
                        <input required type="number" min={1} value={form.warehouse_id} onChange={(e) => setForm({ ...form, warehouse_id: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Nombre *</label>
                        <input required value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                </div>

                <div className="grid grid-cols-3 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estrategia</label>
                        <select value={form.strategy} onChange={(e) => setForm({ ...form, strategy: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="abc">ABC (alta rotación)</option>
                            <option value="zone">Por zona</option>
                            <option value="random">Aleatoria</option>
                            <option value="full">Inventario completo</option>
                        </select>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Frecuencia (días)</label>
                        <input type="number" min={1} value={form.frequency_days} onChange={(e) => setForm({ ...form, frequency_days: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Próxima ejecución</label>
                        <input type="date" value={form.next_run_at} onChange={(e) => setForm({ ...form, next_run_at: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                </div>

                <div className="flex items-center gap-2">
                    <input id="plan_active" type="checkbox" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} className="h-4 w-4 rounded border-gray-300 text-purple-600" />
                    <label htmlFor="plan_active" className="text-sm text-gray-700 dark:text-gray-200">Plan activo</label>
                </div>

                <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Button type="button" variant="outline" onClick={onClose} disabled={submitting}>Cancelar</Button>
                    <Button type="submit" variant="primary" disabled={submitting}>{submitting ? 'Guardando...' : 'Guardar'}</Button>
                </div>
            </form>
        </FormModal>
    );
}
