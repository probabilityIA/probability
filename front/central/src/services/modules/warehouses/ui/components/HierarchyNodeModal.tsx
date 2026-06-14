'use client';

import { useEffect, useState } from 'react';
import { FormModal, Button, Alert } from '@/shared/ui';
import {
    createZoneAction,
    updateZoneAction,
    createAisleAction,
    updateAisleAction,
    createRackAction,
    updateRackAction,
    createRackLevelAction,
    updateRackLevelAction,
} from '../../infra/actions/hierarchy';

export type NodeType = 'zone' | 'aisle' | 'rack' | 'level';

function generateCode(name: string, prefix: string): string {
    const slug = name.trim().toUpperCase().replace(/[^A-Z0-9]+/g, '-').replace(/^-|-$/g, '').slice(0, 15);
    return slug || prefix + '-' + Math.random().toString(36).slice(2, 5).toUpperCase();
}

interface Props {
    warehouseId: number;
    businessId?: number;
    mode: 'create' | 'edit';
    type: NodeType;
    parentId: number | null;
    initial?: Record<string, any>;
    onClose: () => void;
    onSuccess: () => void;
}

const labels: Record<NodeType, string> = {
    zone: 'Zona',
    aisle: 'Pasillo',
    rack: 'Rack',
    level: 'Nivel',
};

export default function HierarchyNodeModal({ warehouseId, businessId, mode, type, parentId, initial, onClose, onSuccess }: Props) {
    const [form, setForm] = useState<Record<string, any>>({
        code: '',
        name: '',
        purpose: '',
        color_hex: '#6366f1',
        levels_count: 1,
        ordinal: 1,
        is_active: true,
    });
    const [submitting, setSubmitting] = useState(false);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        if (initial) setForm((f) => ({ ...f, ...initial }));
    }, [initial]);

    const handleChange = (key: string, value: any) => setForm((f) => ({ ...f, [key]: value }));

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setSubmitting(true);
        setError(null);
        try {
            let result: { success: boolean; error?: string };
            const prefixes: Record<NodeType, string> = { zone: 'Z', aisle: 'A', rack: 'R', level: 'L' };
            const resolvedCode = mode === 'edit'
                ? form.code
                : type === 'level'
                    ? `L-${String(form.ordinal).padStart(2, '0')}`
                    : generateCode(form.name, prefixes[type]);
            if (type === 'zone') {
                if (mode === 'create') {
                    result = await createZoneAction({ warehouse_id: warehouseId, code: resolvedCode, name: form.name, purpose: form.purpose, color_hex: form.color_hex, is_active: form.is_active }, businessId);
                } else {
                    result = await updateZoneAction(initial!.id, warehouseId, { code: resolvedCode, name: form.name, purpose: form.purpose, color_hex: form.color_hex, is_active: form.is_active }, businessId);
                }
            } else if (type === 'aisle') {
                const aisleDims = { width_cm: Number(form.width_cm) || 0 };
                if (mode === 'create') {
                    result = await createAisleAction({ zone_id: parentId!, code: resolvedCode, name: form.name, is_active: form.is_active, ...aisleDims }, warehouseId, businessId);
                } else {
                    result = await updateAisleAction(initial!.id, warehouseId, { code: resolvedCode, name: form.name, is_active: form.is_active, ...aisleDims }, businessId);
                }
            } else if (type === 'rack') {
                const rackDims = { width_cm: Number(form.width_cm) || 0, depth_cm: Number(form.depth_cm) || 0, height_cm: Number(form.height_cm) || 0, side: form.side || '' };
                if (mode === 'create') {
                    result = await createRackAction({ aisle_id: parentId!, code: resolvedCode, name: form.name, levels_count: Number(form.levels_count) || 1, is_active: form.is_active, ...rackDims }, warehouseId, businessId);
                } else {
                    result = await updateRackAction(initial!.id, warehouseId, { code: resolvedCode, name: form.name, levels_count: Number(form.levels_count) || 1, is_active: form.is_active, ...rackDims }, businessId);
                }
            } else {
                if (mode === 'create') {
                    result = await createRackLevelAction({ rack_id: parentId!, code: resolvedCode, ordinal: Number(form.ordinal) || 1, is_active: form.is_active }, warehouseId, businessId);
                } else {
                    result = await updateRackLevelAction(initial!.id, warehouseId, { code: resolvedCode, ordinal: Number(form.ordinal) || 1, is_active: form.is_active }, businessId);
                }
            }
            if (!result.success) {
                setError(result.error || 'Error al guardar');
                return;
            }
            onSuccess();
        } catch (err: any) {
            setError(err.message || 'Error inesperado');
        } finally {
            setSubmitting(false);
        }
    };

    return (
        <FormModal isOpen={true} onClose={onClose} title={`${mode === 'create' ? 'Crear' : 'Editar'} ${labels[type]}`}>
            <form onSubmit={handleSubmit} className="p-6 space-y-4">
                {error && <Alert type="error">{error}</Alert>}

                {type !== 'level' && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Nombre *</label>
                        <input
                            type="text"
                            required
                            value={form.name}
                            onChange={(e) => handleChange('name', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm focus:ring-2 focus:ring-indigo-500"
                        />
                    </div>
                )}

                {type === 'zone' && (
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Proposito</label>
                            <select
                                value={form.purpose}
                                onChange={(e) => handleChange('purpose', e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                            >
                                <option value="">(sin especificar)</option>
                                <option value="picking">Picking</option>
                                <option value="bulk">Bulk</option>
                                <option value="receiving">Recibo</option>
                                <option value="shipping">Despacho</option>
                                <option value="cross_dock">Cross-dock</option>
                                <option value="returns">Devoluciones</option>
                                <option value="quarantine">Cuarentena</option>
                                <option value="damaged">Averiado</option>
                            </select>
                        </div>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Color</label>
                            <input
                                type="color"
                                value={form.color_hex}
                                onChange={(e) => handleChange('color_hex', e.target.value)}
                                className="w-full h-10 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 cursor-pointer"
                            />
                        </div>
                    </div>
                )}

                {type === 'rack' && (
                    <>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1"># Niveles</label>
                            <input
                                type="number"
                                min={1}
                                max={20}
                                value={form.levels_count}
                                onChange={(e) => handleChange('levels_count', e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                            />
                        </div>
                        <div className="grid grid-cols-3 gap-3">
                            {([
                                ['width_cm', 'Ancho (m)'],
                                ['depth_cm', 'Fondo (m)'],
                                ['height_cm', 'Alto (m)'],
                            ] as const).map(([key, label]) => (
                                <div key={key}>
                                    <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">{label}</label>
                                    <input
                                        type="number"
                                        min={0}
                                        step={0.1}
                                        value={form[key] ? Number(form[key]) / 100 : ''}
                                        onChange={(e) => handleChange(key, e.target.value ? Math.round(Number(e.target.value) * 100) : 0)}
                                        className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                                    />
                                </div>
                            ))}
                        </div>
                        <div>
                            <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Lado del pasillo</label>
                            <select
                                value={form.side || ''}
                                onChange={(e) => handleChange('side', e.target.value)}
                                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                            >
                                <option value="">(sin lado)</option>
                                <option value="A">Lado A</option>
                                <option value="B">Lado B</option>
                            </select>
                        </div>
                    </>
                )}

                {type === 'aisle' && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Ancho del pasillo (m)</label>
                        <input
                            type="number"
                            min={0}
                            step={0.1}
                            value={form.width_cm ? Number(form.width_cm) / 100 : ''}
                            onChange={(e) => handleChange('width_cm', e.target.value ? Math.round(Number(e.target.value) * 100) : 0)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                        />
                    </div>
                )}

                {type === 'level' && (
                    <div>
                        <label className="block text-xs font-medium text-gray-600 dark:text-gray-300 mb-1">Ordinal</label>
                        <input
                            type="number"
                            min={1}
                            value={form.ordinal}
                            onChange={(e) => handleChange('ordinal', e.target.value)}
                            className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-md bg-white dark:bg-gray-700 text-gray-900 dark:text-white text-sm"
                        />
                    </div>
                )}

                <div className="flex items-center gap-2">
                    <input
                        id="is_active"
                        type="checkbox"
                        checked={!!form.is_active}
                        onChange={(e) => handleChange('is_active', e.target.checked)}
                        className="h-4 w-4 rounded border-gray-300 text-indigo-600 focus:ring-indigo-500"
                    />
                    <label htmlFor="is_active" className="text-sm text-gray-700 dark:text-gray-200">Activo</label>
                </div>

                <div className="flex justify-end gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
                    <Button type="button" variant="outline" onClick={onClose} disabled={submitting}>Cancelar</Button>
                    <Button type="submit" variant="primary" disabled={submitting}>{submitting ? 'Guardando...' : 'Guardar'}</Button>
                </div>
            </form>
        </FormModal>
    );
}
