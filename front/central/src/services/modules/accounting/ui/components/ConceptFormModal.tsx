'use client';

import { useEffect, useState } from 'react';
import { Modal } from '@/shared/ui';
import { Concept, ConceptKind } from '../../domain/types';
import { createAccountingConceptAction, updateAccountingConceptAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface ConceptFormModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSaved: () => void;
    concept: Concept | null;
}

export function ConceptFormModal({ isOpen, onClose, onSaved, concept }: ConceptFormModalProps) {
    const { showToast } = useToast();
    const isEdit = concept !== null;
    const [code, setCode] = useState('');
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [kind, setKind] = useState<ConceptKind>('INCOME');
    const [isRealIncome, setIsRealIncome] = useState(false);
    const [isActive, setIsActive] = useState(true);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        if (!isOpen) return;
        setCode(concept?.code || '');
        setName(concept?.name || '');
        setDescription(concept?.description || '');
        setKind(concept?.kind || 'INCOME');
        setIsRealIncome(concept?.is_real_income || false);
        setIsActive(concept?.is_active ?? true);
    }, [isOpen, concept]);

    const handleSubmit = async () => {
        if (!name.trim() || (!isEdit && !code.trim())) {
            showToast('Codigo y nombre son requeridos', 'error');
            return;
        }
        setSaving(true);
        const result = isEdit
            ? await updateAccountingConceptAction(concept.id, {
                name: name.trim(),
                description: description.trim(),
                is_real_income: isRealIncome,
                is_active: isActive,
            })
            : await createAccountingConceptAction({
                code: code.trim(),
                name: name.trim(),
                description: description.trim(),
                kind,
                is_real_income: isRealIncome,
            });
        setSaving(false);
        if (result.success) {
            showToast(isEdit ? 'Concepto actualizado' : 'Concepto creado', 'success');
            onSaved();
        } else {
            showToast(result.error, 'error');
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title={isEdit ? 'Editar concepto' : 'Nuevo concepto'} size="lg">
            <div className="space-y-4">
                {!isEdit && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Codigo *</label>
                        <input
                            type="text"
                            value={code}
                            onChange={(e) => setCode(e.target.value.toUpperCase())}
                            placeholder="Ej: PAGO_COLABORADORES"
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
                {!isEdit && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Tipo *</label>
                        <select
                            value={kind}
                            onChange={(e) => setKind(e.target.value as ConceptKind)}
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        >
                            <option value="INCOME">Ingreso</option>
                            <option value="EXPENSE">Gasto</option>
                        </select>
                    </div>
                )}
                <label className="flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                    <input
                        type="checkbox"
                        checked={isRealIncome}
                        onChange={(e) => setIsRealIncome(e.target.checked)}
                        className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                    />
                    Es ganancia real de la plataforma
                </label>
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
                        {saving ? 'Guardando...' : isEdit ? 'Guardar cambios' : 'Crear concepto'}
                    </button>
                </div>
            </div>
        </Modal>
    );
}
