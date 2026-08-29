'use client';

import { useMemo, useState } from 'react';
import { Modal, SuperAdminBusinessSelector } from '@/shared/ui';
import { AccountingService, ClientProfile, Concept, Invoice, InvoiceItemInputDTO, Tax } from '../../domain/types';
import { createAccountingInvoiceAction, getAccountingClientProfileAction, updateAccountingInvoiceAction } from '../../infra/actions';
import { formatCOP, todayISO } from '../format';
import { useToast } from '@/shared/providers/toast-provider';

interface ItemRow {
    description: string;
    quantity: string;
    unit_price: string;
    service_id: number;
    unspsc_code: string;
}

interface InvoiceFormModalProps {
    isOpen: boolean;
    onClose: () => void;
    onSaved: (invoice: Invoice) => void;
    concepts: Concept[];
    taxes: Tax[];
    services: AccountingService[];
    invoice?: Invoice;
}

const inputClass = 'w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500';

function emptyRow(): ItemRow {
    return { description: '', quantity: '1', unit_price: '', service_id: 0, unspsc_code: '' };
}

function rowAmount(row: ItemRow): number {
    const quantity = Number(row.quantity) || 0;
    const unitPrice = Number(row.unit_price) || 0;
    return quantity * unitPrice;
}

export function InvoiceFormModal({ isOpen, onClose, onSaved, concepts, taxes, services, invoice }: InvoiceFormModalProps) {
    const { showToast } = useToast();
    const isEdit = Boolean(invoice);
    const activeServices = useMemo(() => services.filter((s) => s.is_active), [services]);

    const incomeConcepts = useMemo(
        () => concepts.filter((c) => c.kind === 'INCOME' && c.is_active),
        [concepts],
    );
    const activeTaxes = useMemo(() => taxes.filter((t) => t.is_active), [taxes]);
    const chargeTaxes = useMemo(() => activeTaxes.filter((t) => t.kind === 'CHARGE'), [activeTaxes]);
    const withholdingTaxes = useMemo(() => activeTaxes.filter((t) => t.kind === 'WITHHOLDING'), [activeTaxes]);

    const defaultConceptId = useMemo(() => {
        if (invoice) return invoice.concept_id;
        const serviceInvoice = incomeConcepts.find((c) => c.code === 'SERVICE_INVOICE');
        return serviceInvoice?.id ?? 0;
    }, [invoice, incomeConcepts]);

    const initialTaxIds = useMemo(() => {
        if (!invoice) return [];
        const codes = [
            ...(invoice.tax_detail || []).map((d) => d.tax_code),
            ...(invoice.withholding_detail || []).map((d) => d.tax_code),
        ];
        return activeTaxes.filter((t) => codes.includes(t.code)).map((t) => t.id);
    }, [invoice, activeTaxes]);

    const [businessId, setBusinessId] = useState<number | null>(invoice?.business_id ?? null);
    const [conceptId, setConceptId] = useState<number>(defaultConceptId);
    const [issueDate, setIssueDate] = useState(invoice?.issue_date?.slice(0, 10) || todayISO());
    const [dueDate, setDueDate] = useState(invoice?.due_date?.slice(0, 10) || '');
    const [emailTo, setEmailTo] = useState(invoice?.email_to || '');
    const [customerDocument, setCustomerDocument] = useState(invoice?.customer_document || '');
    const [customerPhone, setCustomerPhone] = useState(invoice?.customer_phone || '');
    const [customerAddress, setCustomerAddress] = useState(invoice?.customer_address || '');
    const [notes, setNotes] = useState(invoice?.notes || '');
    const [isTest, setIsTest] = useState(invoice?.is_test ?? false);
    const [taxIds, setTaxIds] = useState<number[]>(initialTaxIds);
    const [rows, setRows] = useState<ItemRow[]>(() => {
        if (invoice?.items?.length) {
            return invoice.items.map((item) => ({
                description: item.description,
                quantity: String(item.quantity),
                unit_price: String(item.unit_price),
                service_id: item.service_id || 0,
                unspsc_code: item.unspsc_code || '',
            }));
        }
        return [emptyRow()];
    });
    const [saving, setSaving] = useState(false);
    const [profile, setProfile] = useState<ClientProfile | null>(null);

    const subtotal = rows.reduce((acc, row) => acc + rowAmount(row), 0);
    const taxLines = chargeTaxes
        .filter((t) => taxIds.includes(t.id))
        .map((tax) => ({
            tax,
            amount: Math.round((subtotal * tax.rate_percent) / 100),
        }));
    const taxTotal = taxLines.reduce((acc, line) => acc + line.amount, 0);
    const total = subtotal + taxTotal;
    const withholdingLines = withholdingTaxes
        .filter((t) => taxIds.includes(t.id))
        .map((tax) => ({
            tax,
            amount: Math.round((subtotal * tax.rate_percent) / 100),
        }));
    const withholdingTotal = withholdingLines.reduce((acc, line) => acc + line.amount, 0);
    const netReceivable = total - withholdingTotal;

    const handleBusinessChange = async (id: number | null) => {
        setBusinessId(id);
        setProfile(null);
        if (!id) return;
        const result = await getAccountingClientProfileAction(id);
        if (!result.success) return;
        const p = result.data;
        setProfile(p);
        if (!p.configured) return;
        if (p.document_number) setCustomerDocument(p.document_number);
        if (p.phone) setCustomerPhone(p.phone);
        if (p.address) setCustomerAddress(p.address);
        if (p.billing_email) setEmailTo((prev) => prev || p.billing_email);
        const suggested = p.suggested_tax_ids || [];
        if (suggested.length > 0) {
            const visibleIds = [...chargeTaxes, ...withholdingTaxes].map((t) => t.id);
            setTaxIds(suggested.filter((taxId) => visibleIds.includes(taxId)));
        }
    };

    const updateRow = (index: number, patch: Partial<ItemRow>) => {
        setRows((prev) => prev.map((row, i) => (i === index ? { ...row, ...patch } : row)));
    };

    const selectService = (index: number, serviceId: number) => {
        if (!serviceId) {
            updateRow(index, { service_id: 0, unspsc_code: '' });
            return;
        }
        const service = activeServices.find((s) => s.id === serviceId);
        if (!service) return;
        setRows((prev) => prev.map((row, i) => {
            if (i !== index) return row;
            const currentPrice = Number(row.unit_price) || 0;
            return {
                ...row,
                service_id: service.id,
                unspsc_code: service.unspsc_code,
                description: service.name,
                unit_price: currentPrice === 0 && service.unit_price > 0 ? String(service.unit_price) : row.unit_price,
            };
        }));
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
            items.push({
                description,
                quantity,
                unit_price: unitPrice,
                service_id: row.service_id || undefined,
                unspsc_code: row.service_id ? row.unspsc_code || undefined : undefined,
            });
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
            customer_document: customerDocument.trim() || undefined,
            customer_phone: customerPhone.trim() || undefined,
            customer_address: customerAddress.trim() || undefined,
            is_test: isTest,
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
                            onChange={handleBusinessChange}
                            variant="default"
                            placeholder="Selecciona un negocio"
                        />
                        {profile && !profile.configured && (
                            <p className="mt-1 text-xs text-amber-600 dark:text-amber-400">
                                Este negocio no tiene ficha fiscal configurada (se edita en el negocio)
                            </p>
                        )}
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

                <div className="grid grid-cols-1 sm:grid-cols-3 gap-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Documento (NIT/CC)</label>
                        <input
                            type="text"
                            value={customerDocument}
                            onChange={(e) => setCustomerDocument(e.target.value)}
                            placeholder="900123456"
                            className={inputClass}
                        />
                        {profile?.configured && profile.dv && (
                            <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">
                                {profile.document_type || 'NIT'} {profile.document_number}-{profile.dv}
                            </p>
                        )}
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{'Tel\u00e9fono'}</label>
                        <input
                            type="text"
                            value={customerPhone}
                            onChange={(e) => setCustomerPhone(e.target.value)}
                            placeholder="3001234567"
                            className={inputClass}
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">{'Direcci\u00f3n'}</label>
                        <input
                            type="text"
                            value={customerAddress}
                            onChange={(e) => setCustomerAddress(e.target.value)}
                            placeholder="Calle 1 # 2-3"
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
                    <label className="inline-flex items-center gap-2 text-sm text-gray-700 dark:text-gray-300">
                        <input
                            type="checkbox"
                            checked={isTest}
                            onChange={(e) => setIsTest(e.target.checked)}
                            className="w-4 h-4 rounded border-gray-300 text-purple-600 focus:ring-purple-500"
                        />
                        Factura de prueba (modo test)
                    </label>
                    <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">No genera movimientos contables al marcarse pagada</p>
                </div>

                <div className="space-y-3">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Impuestos (suman al total)</label>
                        {chargeTaxes.length === 0 ? (
                            <p className="text-sm text-gray-500 dark:text-gray-400">No hay impuestos activos configurados</p>
                        ) : (
                            <div className="flex flex-wrap gap-3">
                                {chargeTaxes.map((tax) => (
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
                    {withholdingTaxes.length > 0 && (
                        <div>
                            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Retenciones (informativas)</label>
                            <div className="flex flex-wrap gap-3">
                                {withholdingTaxes.map((tax) => (
                                    <label
                                        key={tax.id}
                                        className="inline-flex items-center gap-2 px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg cursor-pointer text-gray-500 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-700"
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
                            <p className="mt-1 text-xs text-gray-400 dark:text-gray-500">No suman al total; se muestran como neto a recibir</p>
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
                                    <th className="pb-2 font-medium w-44">Servicio</th>
                                    <th className="pb-2 font-medium">{'Descripci\u00f3n'}</th>
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
                                            <select
                                                value={row.service_id || 0}
                                                onChange={(e) => selectService(index, Number(e.target.value))}
                                                className={inputClass}
                                            >
                                                <option value={0}>-- Manual --</option>
                                                {activeServices.map((service) => (
                                                    <option key={service.id} value={service.id}>{service.name}</option>
                                                ))}
                                            </select>
                                        </td>
                                        <td className="py-1 pr-2">
                                            <input
                                                type="text"
                                                value={row.description}
                                                onChange={(e) => updateRow(index, { description: e.target.value })}
                                                placeholder="Descripci\u00f3n del servicio"
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
                            {withholdingLines.length > 0 && (
                                <>
                                    {withholdingLines.map((line) => (
                                        <div key={line.tax.id} className="flex justify-between text-gray-400 dark:text-gray-500">
                                            <span>{line.tax.name} ({line.tax.rate_percent}%)</span>
                                            <span>-{formatCOP(line.amount)}</span>
                                        </div>
                                    ))}
                                    <div className="flex justify-between pt-1.5 border-t border-gray-200 dark:border-gray-700 font-semibold text-gray-900 dark:text-white">
                                        <span>Neto a recibir</span>
                                        <span>{formatCOP(netReceivable)}</span>
                                    </div>
                                </>
                            )}
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
