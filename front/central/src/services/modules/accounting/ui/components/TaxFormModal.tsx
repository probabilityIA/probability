'use client';

import { useEffect, useState } from 'react';
import { Modal } from '@/shared/ui';
import { Tax, TaxKind } from '../../domain/types';
import { createAccountingTaxAction, updateAccountingTaxAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface TaxFormModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSaved: () => void;
    tax: Tax | null;
}

export function TaxFormModal({ isOpen, onClose, onSaved, tax }: TaxFormModalProps) {
    const { showToast } = useToast();
    const isEdit = tax !== null;
    const [code, setCode] = useState('');
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [ratePercent, setRatePercent] = useState('');
    const [kind, setKind] = useState<TaxKind>('CHARGE');
    const [isActive, setIsActive] = useState(true);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        if (!isOpen) return;
        setCode(tax?.code || '');
        setName(tax?.name || '');
        setDescription(tax?.description || '');
        setRatePercent(tax ? String(tax.rate_percent) : '');
        setKind(tax?.kind || 'CHARGE');
        setIsActive(tax?.is_active ?? true);
    }, [isOpen, tax]);

    const handleSubmit = async () => {
        const rate = Number(ratePercent);
        if (!name.trim() || (!isEdit && !code.trim())) {
            showToast('Codigo y nombre son requeridos', 'error');
            return;
        }
        if (isNaN(rate) || rate < 0 || rate > 100) {
            showToast('La tarifa debe estar entre 0 y 100', 'error');
            return;
        }
        setSaving(true);
        const result = isEdit
            ? await updateAccountingTaxAction(tax.id, {
                name: name.trim(),
                description: description.trim(),
                rate_percent: rate,
                kind,
                is_active: isActive,
            })
            : await createAccountingTaxAction({
                code: code.trim(),
                name: name.trim(),
                description: description.trim(),
                rate_percent: rate,
                kind,
            });
        setSaving(false);
        if (result.success) {
            showToast(isEdit ? 'Impuesto actualizado' : 'Impuesto creado', 'success');
            onSaved();
        } else {
            showToast(result.error, 'error');
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title={isEdit ? 'Editar impuesto' : 'Nuevo impuesto'} size="md">
            <div className="space-y-4">
                {!isEdit && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Codigo *</label>
                        <input
                            type="text"
                            value={code}
                            onChange={(e) => setCode(e.target.value.toUpperCase())}
                            placeholder="Ej: IVA"
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                )}
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Nombre *</label>
                    <input
                        type="text"
                        value={name}
                        onChange={(e) => setName(e.target.value)}
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Descripcion</label>
                    <textarea
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        rows={2}
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tarifa (%) *</label>
                    <input
                        type="number"
                        min="0"
                        max="100"
                        step="0.01"
                        value={ratePercent}
                        onChange={(e) => setRatePercent(e.target.value)}
                        placeholder="19"
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tipo *</label>
                    <select
                        value={kind}
                        onChange={(e) => setKind(e.target.value as TaxKind)}
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    >
                        <option value="CHARGE">Cargo (suma al total)</option>
                        <option value="WITHHOLDING">Retencion (informativa, no suma)</option>
                        <option value="OTHER">Otro (no aplica a facturas)</option>
                    </select>
                </div>
                {isEdit && (
                    <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                        <input
                            type="checkbox"
                            checked={isActive}
                            onChange={(e) => setIsActive(e.target.checked)}
                            className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                        />
                        Activo
                    </label>
                )}
                <div className="flex justify-end gap-3 pt-2">
                    <button
                        onClick={onClose}
                        className="px-4 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                        Cancelar
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={saving}
                        className="px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors disabled:opacity-50"
                    >
                        {saving ? 'Guardando...' : isEdit ? 'Guardar cambios' : 'Crear impuesto'}
                    </button>
                </div>
            </div>
        </Modal>
    );
}
