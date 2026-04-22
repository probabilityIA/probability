'use client';

import { useState } from 'react';
import { FormModal, Button, Alert } from '@/shared/ui';
import { createPutawayRuleAction, updatePutawayRuleAction } from '../../infra/actions/operations';
import { PutawayRule } from '../../domain/operations-types';

interface Props {
    businessId?: number;
    rule?: PutawayRule | null;
    onClose: () => void;
    onSuccess: () => void;
}

export default function PutawayRuleFormModal({ businessId, rule, onClose, onSuccess }: Props) {
    const isEdit = !!rule;
    const [form, setForm] = useState({
        product_id: rule?.product_id || '',
        category_id: rule?.category_id || '',
        target_zone_id: rule?.target_zone_id || 0,
        priority: rule?.priority ?? 10,
        strategy: rule?.strategy || 'closest',
        is_active: rule?.is_active ?? true,
    });
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const submit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        setError(null);
        try {
            const data: any = {
                product_id: form.product_id || null,
                category_id: form.category_id ? Number(form.category_id) : null,
                target_zone_id: Number(form.target_zone_id),
                priority: Number(form.priority),
                strategy: form.strategy,
                is_active: form.is_active,
            };
            const r = isEdit ? await updatePutawayRuleAction(rule!.id, data, businessId) : await createPutawayRuleAction(data, businessId);
            if (!r.success) { setError(r.error || 'Error al guardar'); return; }
            onSuccess();
        } catch (err: any) { setError(err.message); } finally { setSubmitting(false); }
    };

    return (
        <FormModal isOpen={true} onClose={onClose} title={isEdit ? 'Editar regla de put-away' : 'Crear regla de put-away'}>
            <form onSubmit={submit} className="p-6 space-y-4">
                {error && <Alert type="error">{error}</Alert>}

                <div className="grid grid-cols-2 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Producto (SKU) [opcional]</label>
                        <input value={form.product_id} onChange={(e) => setForm({ ...form, product_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Categoría (ID) [opcional]</label>
                        <input type="number" value={form.category_id} onChange={(e) => setForm({ ...form, category_id: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                </div>

                <div className="grid grid-cols-3 gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Zona destino *</label>
                        <input required type="number" min={1} value={form.target_zone_id} onChange={(e) => setForm({ ...form, target_zone_id: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Prioridad</label>
                        <input type="number" value={form.priority} onChange={(e) => setForm({ ...form, priority: Number(e.target.value) })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm" />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Estrategia</label>
                        <select value={form.strategy} onChange={(e) => setForm({ ...form, strategy: e.target.value })} className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm">
                            <option value="closest">Más cercano</option>
                            <option value="fewest_items">Menor ocupación</option>
                            <option value="dedicated">Dedicado</option>
                            <option value="random">Aleatorio</option>
                        </select>
                    </div>
                </div>

                <div className="flex items-center gap-2">
                    <input id="pa_active" type="checkbox" checked={form.is_active} onChange={(e) => setForm({ ...form, is_active: e.target.checked })} className="h-4 w-4 rounded border-gray-300 text-purple-600" />
                    <label htmlFor="pa_active" className="text-sm text-gray-700 dark:text-gray-200">Regla activa</label>
                </div>

                <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Button type="button" variant="outline" onClick={onClose} disabled={submitting}>Cancelar</Button>
                    <Button type="submit" variant="primary" disabled={submitting}>{submitting ? 'Guardando...' : 'Guardar'}</Button>
                </div>
            </form>
        </FormModal>
    );
}
