'use client';

import { useEffect, useState } from 'react';
import { Modal } from '@/shared/ui';
import { AccountingService } from '../../domain/types';
import { createAccountingServiceAction, updateAccountingServiceAction } from '../../infra/actions';
import { useToast } from '@/shared/providers/toast-provider';

interface ServiceFormModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSaved: () => void;
    service: AccountingService | null;
}

export function ServiceFormModal({ isOpen, onClose, onSaved, service }: ServiceFormModalProps) {
    const { showToast } = useToast();
    const isEdit = service !== null;
    const [code, setCode] = useState('');
    const [name, setName] = useState('');
    const [unspscCode, setUnspscCode] = useState('');
    const [unitPrice, setUnitPrice] = useState('');
    const [isActive, setIsActive] = useState(true);
    const [saving, setSaving] = useState(false);

    useEffect(() => {
        if (!isOpen) return;
        setCode(service?.code || '');
        setName(service?.name || '');
        setUnspscCode(service?.unspsc_code || '');
        setUnitPrice(service && service.unit_price > 0 ? String(service.unit_price) : '');
        setIsActive(service?.is_active ?? true);
    }, [isOpen, service]);

    const handleSubmit = async () => {
        if (!code.trim() || !name.trim()) {
            showToast('Codigo y nombre son requeridos', 'error');
            return;
        }
        const price = Number(unitPrice);
        if (unitPrice !== '' && (isNaN(price) || price < 0)) {
            showToast('El precio base debe ser un numero mayor o igual a cero', 'error');
            return;
        }
        setSaving(true);
        const payload = {
            code: code.trim(),
            name: name.trim(),
            unspsc_code: unspscCode.trim() || undefined,
            unit_price: unitPrice !== '' ? price : undefined,
        };
        const result = isEdit
            ? await updateAccountingServiceAction(service.id, { ...payload, is_active: isActive })
            : await createAccountingServiceAction(payload);
        setSaving(false);
        if (result.success) {
            showToast(isEdit ? 'Servicio actualizado' : 'Servicio creado', 'success');
            onSaved();
        } else {
            showToast(result.error, 'error');
        }
    };

    return (
        <Modal isOpen={isOpen} onClose={onClose} title={isEdit ? 'Editar servicio' : 'Nuevo servicio'} size="md">
            <div className="space-y-4">
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Codigo *</label>
                    <input
                        type="text"
                        value={code}
                        onChange={(e) => setCode(e.target.value.toUpperCase())}
                        placeholder="Ej: PLAN-PRO"
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                </div>
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
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Codigo UNSPSC</label>
                    <input
                        type="text"
                        value={unspscCode}
                        onChange={(e) => setUnspscCode(e.target.value)}
                        placeholder="81112105"
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
                    <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">Clasificacion DIAN/UNSPSC de 8 digitos, opcional</p>
                </div>
                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Precio base</label>
                    <input
                        type="number"
                        min="0"
                        step="1"
                        value={unitPrice}
                        onChange={(e) => setUnitPrice(e.target.value)}
                        placeholder="0"
                        className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                    />
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
                        {saving ? 'Guardando...' : isEdit ? 'Guardar cambios' : 'Crear servicio'}
                    </button>
                </div>
            </div>
        </Modal>
    );
}
