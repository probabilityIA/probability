'use client';

import { useState } from 'react';
import { Modal, SuperAdminBusinessSelector } from '@/shared/ui';
import { Concept } from '../../domain/types';
import { createAccountingEntryAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface EntryFormModalProps {
    isOpen: boolean;
    onClose: () => void;
    onCreated: () => void;
    concepts: Concept[];
}

function today(): string {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
}

export function EntryFormModal({ isOpen, onClose, onCreated, concepts }: EntryFormModalProps) {
    const { showToast } = useToast();
    const [conceptId, setConceptId] = useState<number>(0);
    const [amount, setAmount] = useState('');
    const [entryDate, setEntryDate] = useState(today());
    const [description, setDescription] = useState('');
    const [businessId, setBusinessId] = useState<number | null>(null);
    const [saving, setSaving] = useState(false);

    const activeConcepts = concepts.filter((c) => c.is_active);

    const reset = () => {
        setConceptId(0);
        setAmount('');
        setEntryDate(today());
        setDescription('');
        setBusinessId(null);
    };

    const handleClose = () => {
        reset();
        onClose();
    };

    const handleSubmit = async () => {
        const parsedAmount = Number(amount);
        if (!conceptId) {
            showToast('Selecciona un concepto', 'error');
            return;
        }
        if (!parsedAmount || parsedAmount <= 0) {
            showToast('Ingresa un monto valido', 'error');
            return;
        }
        if (!entryDate) {
            showToast('Selecciona una fecha', 'error');
            return;
        }
        setSaving(true);
        const result = await createAccountingEntryAction({
            concept_id: conceptId,
            related_business_id: businessId || undefined,
            entry_date: entryDate,
            amount: parsedAmount,
            description,
        });
        setSaving(false);
        if (result.success) {
            showToast('Movimiento creado', 'success');
            reset();
            onCreated();
        } else {
            showToast(result.error, 'error');
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={handleClose} title="Nuevo movimiento manual" size="lg">
            <div className="space-y-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Concepto *</label>
                    <select
                        value={conceptId}
                        onChange={(e) => setConceptId(Number(e.target.value))}
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    >
                        <option value={0}>-- Selecciona un concepto --</option>
                        {activeConcepts.map((c) => (
                            <option key={c.id} value={c.id}>
                                {c.name} ({c.kind === 'INCOME' ? 'Ingreso' : 'Gasto'})
                            </option>
                        ))}
                    </select>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Monto (COP) *</label>
                        <input
                            type="number"
                            min="0"
                            step="1"
                            value={amount}
                            onChange={(e) => setAmount(e.target.value)}
                            placeholder="0"
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Fecha *</label>
                        <input
                            type="date"
                            value={entryDate}
                            onChange={(e) => setEntryDate(e.target.value)}
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Descripción</label>
                    <textarea
                        value={description}
                        onChange={(e) => setDescription(e.target.value)}
                        rows={3}
                        placeholder="Ej: pago a colaboradores"
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Negocio relacionado (opcional)</label>
                    <SuperAdminBusinessSelector
                        value={businessId}
                        onChange={setBusinessId}
                        variant="default"
                        placeholder="Sin negocio"
                    />
                </div>

                <div className="flex justify-end gap-3 pt-2">
                    <button
                        onClick={handleClose}
                        className="px-4 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                        Cancelar
                    </button>
                    <button
                        onClick={handleSubmit}
                        disabled={saving}
                        className="px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors disabled:opacity-50"
                    >
                        {saving ? 'Guardando...' : 'Crear movimiento'}
                    </button>
                </div>
            </div>
        </Modal>
    );
}
