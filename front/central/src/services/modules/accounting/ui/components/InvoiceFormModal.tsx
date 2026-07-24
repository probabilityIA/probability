'use client';

import { useMemo, useState } from 'react';
import { Modal, SuperAdminBusinessSelector } from '@/shared/ui';
import { Concept, Invoice, InvoiceItemInputDTO, Tax } from '../../domain/types';
import { createAccountingInvoiceAction, updateAccountingInvoiceAction } from '../../infra/actions';
import { formatCOP, todayISO } from '../format';
import { useToast } from '@/shared/providers/toast-provider';

interface ItemRow {
    description: string;
    quantity: string;
    unit_price: string;
}

interface InvoiceFormModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSaved: (invoice: Invoice) => void;
    concepts: Concept[];
    taxes: Tax[];
    invoice?: Invoice;
}

const inputClass = 'w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500';

function emptyRow(): ItemRow {
    return { description: '', quantity: '1', unit_price: '' };
}

function rowAmount(row: ItemRow): number {
    const quantity = Number(row.quantity) || 0;
    const unitPrice = Number(row.unit_price) || 0;
    return quantity * unitPrice;
}

export function InvoiceFormModal({ isOpen, onClose, onSaved, concepts, taxes, invoice }: InvoiceFormModalProps) {
    const { showToast } = useToast();
    const isEdit = Boolean(invoice);

    const incomeConcepts = useMemo(
        () => concepts.filter((c) => c.kind === 'INCOME' && c.is_active),
        [concepts],
    );
    const activeTaxes = useMemo(() => taxes.filter((t) => t.is_active), [taxes]);

    const defaultConceptId = useMemo(() => {
        if (invoice) return invoice.concept_id;
        const serviceInvoice = incomeConcepts.find((c) => c.code === 'SERVICE_INVOICE');
        return serviceInvoice?.id ?? 0;
    }, [invoice, incomeConcepts]);

    const initialTaxIds = useMemo(() => {
        if (!invoice) return [];
        const codes = (invoice.tax_detail || []).map((d) => d.tax_code);
        return activeTaxes.filter((t) => codes.includes(t.code)).map((t) => t.id);
    }, [invoice, activeTaxes]);

    const [businessId, setBusinessId] = useState<number | null>(invoice?.business_id ?? null);
    const [conceptId, setConceptId] = useState<number>(defaultConceptId);
    const [issueDate, setIssueDate] = useState(invoice?.issue_date?.slice(0, 10) || todayISO());
    const [dueDate, setDueDate] = useState(invoice?.due_date?.slice(0, 10) || '');
    const [emailTo, setEmailTo] = useState(invoice?.email_to || '');
    const [notes, setNotes] = useState(invoice?.notes || '');
    const [taxIds, setTaxIds] = useState<number[]>(initialTaxIds);
    const [rows, setRows] = useState<ItemRow[]>(() => {
        if (invoice?.items?.length) {
            return invoice.items.map((item) => ({
                description: item.description,
                quantity: String(item.quantity),
                unit_price: String(item.unit_price),
            }));
        }
        return [emptyRow()];
    });
    const [saving, setSaving] = useState(false);

    const subtotal = rows.reduce((acc, row) => acc + rowAmount(row), 0);
    const selectedTaxes = activeTaxes.filter((t) => taxIds.includes(t.id));
    const taxLines = selectedTaxes.map((tax) => ({
        tax,
        amount: Math.round((subtotal * tax.rate_percent) / 100),
    }));
    const taxTotal = taxLines.reduce((acc, line) => acc + line.amount, 0);
    const total = subtotal + taxTotal;

    const updateRow = (index: number, patch: Partial<ItemRow>) => {
        setRows((prev) => prev.map((row, i) => (i === index ? { ...row, ...patch } : row)));
    };

    const removeRow = (index: number) => {
        setRows((prev) => (prev.length > 1 ? prev.filter((_, i) => i !== index) : prev));
    };

    const toggleTax = (taxId: number) => {
        setTaxIds((prev) => (prev.includes(taxId) ? prev.filter((id) => id !== taxId) : [...prev, taxId]));
    };

    const handleSubmit = async () => {
        if (!isEdit && !businessId) {
            showToast('Selecciona el negocio cliente', 'error');
            return;
        }
        if (!conceptId) {
            showToast('Selecciona un concepto', 'error');
            return;
        }
        if (!issueDate) {
            showToast('Selecciona la fecha de emision', 'error');
            return;
        }
        const items: InvoiceItemInputDTO[] = [];
        for (const row of rows) {
            const description = row.description.trim();
            const quantity = Number(row.quantity);
            const unitPrice = Number(row.unit_price);
            if (!description) {
                showToast('Cada item debe tener descripcion', 'error');
                return;
            }
            if (!quantity || quantity <= 0) {
                showToast('Cada item debe tener cantidad mayor a cero', 'error');
                return;
            }
            if (!unitPrice || unitPrice <= 0) {
                showToast('Cada item debe tener precio unitario mayor a cero', 'error');
                return;
            }
            items.push({ description, quantity, unit_price: unitPrice });
        }
        if (items.length === 0) {
            showToast('Agrega al menos un item', 'error');
            return;
        }

        setSaving(true);
        const payload = {
            concept_id: conceptId,
            issue_date: issueDate,
            due_date: dueDate || undefined,
            notes: notes || undefined,
            email_to: emailTo || undefined,
            tax_ids: taxIds.length > 0 ? taxIds : undefined,
            items,
        };
        const result = isEdit && invoice
            ? await updateAccountingInvoiceAction(invoice.id, payload)
            : await createAccountingInvoiceAction({ ...payload, business_id: businessId as number });
        setSaving(false);

        if (result.success) {
            showToast(isEdit ? 'Factura actualizada' : 'Factura creada', 'success');
            onSaved(result.data);
        } else {
            showToast(result.error, 'error');
        }
    };

    return (
        <Modal
            isOpen={isOpen}
            onClose={onClose}
            title={isEdit ? `Editar factura ${invoice?.number}` : 'Nueva factura'}
            size="4xl"
        >
            <div className="space-y-4">
                <p className="text-sm text-gray-500 dark:text-gray-400">
                    {isEdit ? 'Solo las facturas en borrador se pueden editar' : 'Emitir una factura a un negocio cliente'}
                </p>

                {!isEdit && (
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Negocio cliente *</label>
                        <SuperAdminBusinessSelector
                            value={businessId}
                            onChange={setBusinessId}
                            variant="default"
                            placeholder="Selecciona un negocio"
                        />
                    </div>
                )}
                {isEdit && invoice && (
                    <div className="text-sm text-gray-700 dark:text-gray-300">
                        <span className="font-medium">Cliente:</span> {invoice.business_name || `#${invoice.business_id}`}
                    </div>
                )}

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Concepto *</label>
                        <select
                            value={conceptId}
                            onChange={(e) => setConceptId(Number(e.target.value))}
                            className={inputClass}
                        >
                            <option value={0}>-- Selecciona un concepto --</option>
                            {incomeConcepts.map((c) => (
                                <option key={c.id} value={c.id}>{c.name}</option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Correo destino (opcional)</label>
                        <input
                            type="email"
                            value={emailTo}
                            onChange={(e) => setEmailTo(e.target.value)}
                            placeholder="cliente@correo.com"
                            className={inputClass}
                        />
                    </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Fecha de emision *</label>
                        <input
                            type="date"
                            value={issueDate}
                            onChange={(e) => setIssueDate(e.target.value)}
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Vencimiento (opcional)</label>
                        <input
                            type="date"
                            value={dueDate}
                            onChange={(e) => setDueDate(e.target.value)}
                            className={inputClass}
                        />
                    </div>
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Notas</label>
                    <textarea
                        value={notes}
                        onChange={(e) => setNotes(e.target.value)}
                        rows={2}
                        placeholder="Ej: factura por servicios del mes"
                        className={inputClass}
                    />
                </div>

                <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Impuestos aplicables</label>
                    {activeTaxes.length === 0 ? (
                        <p className="text-sm text-gray-500 dark:text-gray-400">No hay impuestos activos configurados</p>
                    ) : (
                        <div className="flex flex-wrap gap-3">
                            {activeTaxes.map((tax) => (
                                <label
                                    key={tax.id}
                                    className="inline-flex items-center gap-2 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700"
                                >
                                    <input
                                        type="checkbox"
                                        checked={taxIds.includes(tax.id)}
                                        onChange={() => toggleTax(tax.id)}
                                        className="rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                                    />
                                    <span>{tax.name} ({tax.rate_percent}%)</span>
                                </label>
                            ))}
                        </div>
                    )}
                </div>

                <div className="pt-4 border-t border-gray-200 dark:border-gray-700 space-y-4">
                    <div className="flex items-center justify-between">
                        <h4 className="text-sm font-semibold text-gray-900 dark:text-white uppercase">Items</h4>
                        <button
                            onClick={() => setRows((prev) => [...prev, emptyRow()])}
                            className="inline-flex items-center gap-1.5 px-3 py-1.5 text-sm font-medium rounded-lg border border-purple-300 dark:border-purple-700 text-purple-700 dark:text-purple-300 hover:bg-purple-50 dark:hover:bg-purple-900/30 transition-colors"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                            Agregar item
                        </button>
                    </div>

                    <div className="overflow-x-auto">
                        <table className="w-full text-sm">
                            <thead>
                                <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase">
                                    <th className="pb-2 font-medium">Descripcion</th>
                                    <th className="pb-2 font-medium w-24">Cantidad</th>
                                    <th className="pb-2 font-medium w-40">Precio unitario</th>
                                    <th className="pb-2 font-medium w-32 text-right">Subtotal</th>
                                    <th className="pb-2 w-10"></th>
                                </tr>
                            </thead>
                            <tbody>
                                {rows.map((row, index) => (
                                    <tr key={index} className="align-top">
                                        <td className="py-1 pr-2">
                                            <input
                                                type="text"
                                                value={row.description}
                                                onChange={(e) => updateRow(index, { description: e.target.value })}
                                                placeholder="Descripcion del servicio"
                                                className={inputClass}
                                            />
                                        </td>
                                        <td className="py-1 pr-2">
                                            <input
                                                type="number"
                                                min="1"
                                                step="1"
                                                value={row.quantity}
                                                onChange={(e) => updateRow(index, { quantity: e.target.value })}
                                                className={inputClass}
                                            />
                                        </td>
                                        <td className="py-1 pr-2">
                                            <input
                                                type="number"
                                                min="0"
                                                step="1"
                                                value={row.unit_price}
                                                onChange={(e) => updateRow(index, { unit_price: e.target.value })}
                                                placeholder="0"
                                                className={inputClass}
                                            />
                                        </td>
                                        <td className="py-1 text-right whitespace-nowrap pt-3 font-medium text-gray-800 dark:text-gray-200">
                                            {formatCOP(rowAmount(row))}
                                        </td>
                                        <td className="py-1 text-right pt-2">
                                            <button
                                                onClick={() => removeRow(index)}
                                                disabled={rows.length <= 1}
                                                className="p-1.5 rounded-lg text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 disabled:opacity-30 transition-colors"
                                                title="Quitar item"
                                            >
                                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                                                </svg>
                                            </button>
                                        </td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    </div>

                    <div className="flex justify-end">
                        <div className="w-full sm:w-72 space-y-1.5 text-sm">
                            <div className="flex justify-between text-gray-700 dark:text-gray-300">
                                <span>Subtotal</span>
                                <span className="font-medium">{formatCOP(subtotal)}</span>
                            </div>
                            {taxLines.map((line) => (
                                <div key={line.tax.id} className="flex justify-between text-gray-500 dark:text-gray-400">
                                    <span>{line.tax.name} ({line.tax.rate_percent}%)</span>
                                    <span>{formatCOP(line.amount)}</span>
                                </div>
                            ))}
                            <div className="flex justify-between pt-1.5 border-t border-gray-200 dark:border-gray-700 text-base font-semibold text-gray-900 dark:text-white">
                                <span>Total</span>
                                <span>{formatCOP(total)}</span>
                            </div>
                        </div>
                    </div>
                </div>

                <div className="flex justify-end gap-3 pt-2 border-t border-gray-200 dark:border-gray-700">
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
                        {saving ? 'Guardando...' : isEdit ? 'Guardar cambios' : 'Crear factura'}
                    </button>
                </div>
            </div>
        </Modal>
    );
}
