'use client';

import { useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { ConfirmModal, Modal } from '@/shared/ui';
import { Invoice } from '../../domain/types';
import {
    cancelAccountingInvoiceAction,
    deleteAccountingInvoiceAction,
    payAccountingInvoiceAction,
    sendAccountingInvoiceAction,
} from '../../infra/actions';
import { formatCOP, formatEntryDate, invoiceStatusBadgeClass, invoiceStatusLabel } from '../format';
import { useToast } from '@/shared/providers/toast-provider';

interface InvoiceDetailViewProps {
    invoice: Invoice;
}

const printStyles = `
@media print {
    body * {
        visibility: hidden;
    }
    .invoice-print-area, .invoice-print-area * {
        visibility: visible;
    }
    .invoice-print-area {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
        border: none !important;
        box-shadow: none !important;
    }
}
`;

export function InvoiceDetailView({ invoice }: InvoiceDetailViewProps) {
    const router = useRouter();
    const { showToast } = useToast();
    const [showSend, setShowSend] = useState(false);
    const [sendEmail, setSendEmail] = useState(invoice.email_to || '');
    const [sending, setSending] = useState(false);
    const [showPay, setShowPay] = useState(false);
    const [showCancel, setShowCancel] = useState(false);
    const [showDelete, setShowDelete] = useState(false);
    const [busy, setBusy] = useState(false);

    const isDraft = invoice.status === 'DRAFT';
    const isPaid = invoice.status === 'PAID';
    const canSend = invoice.status === 'DRAFT' || invoice.status === 'SENT';
    const canPay = invoice.status === 'DRAFT' || invoice.status === 'SENT';
    const canCancel = invoice.status !== 'CANCELLED';
    const canDelete = invoice.status === 'DRAFT' || invoice.status === 'CANCELLED';

    const handleSend = async () => {
        setSending(true);
        const result = await sendAccountingInvoiceAction(invoice.id, sendEmail.trim() || undefined);
        setSending(false);
        if (result.success) {
            showToast('Factura enviada por correo', 'success');
            setShowSend(false);
            router.refresh();
        } else {
            showToast(result.error, 'error');
        }
    };

    const runAction = async (
        action: () => Promise<{ success: true; data: unknown } | { success: false; error: string }>,
        successMessage: string,
        redirect?: string,
    ) => {
        setBusy(true);
        const result = await action();
        setBusy(false);
        if (result.success) {
            showToast(successMessage, 'success');
            if (redirect) {
                router.push(redirect);
            } else {
                router.refresh();
            }
        } else {
            showToast(result.error, 'error');
        }
    };

    return (
        <div className="space-y-4 max-w-4xl">
            <style>{printStyles}</style>

            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 print:hidden">
                <div className="flex items-center gap-3">
                    <Link
                        href="/accounting/facturas"
                        className="p-2 rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                        title="Volver a facturas"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
                        </svg>
                    </Link>
                    <div>
                        <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Factura {invoice.number}</h1>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                            {invoice.business_name || `Negocio #${invoice.business_id}`}
                        </p>
                    </div>
                </div>
                <div className="flex flex-wrap items-center gap-2">
                    <button
                        onClick={() => window.print()}
                        className="inline-flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                    >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4H7v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                        </svg>
                        Imprimir
                    </button>
                    {canSend && (
                        <button
                            onClick={() => {
                                setSendEmail(invoice.email_to || '');
                                setShowSend(true);
                            }}
                            disabled={busy}
                            className="inline-flex items-center gap-1.5 px-3 py-2 text-sm font-semibold rounded-lg bg-blue-600 hover:bg-blue-700 text-white transition-colors disabled:opacity-50"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                            </svg>
                            Enviar por correo
                        </button>
                    )}
                    {isDraft && (
                        <Link
                            href={`/accounting/facturas/${invoice.id}/editar`}
                            className="inline-flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                            Editar
                        </Link>
                    )}
                    {canPay && (
                        <button
                            onClick={() => setShowPay(true)}
                            disabled={busy}
                            className="inline-flex items-center gap-1.5 px-3 py-2 text-sm font-semibold rounded-lg bg-emerald-600 hover:bg-emerald-700 text-white transition-colors disabled:opacity-50"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                            </svg>
                            Marcar pagada
                        </button>
                    )}
                    {canCancel && (
                        <button
                            onClick={() => setShowCancel(true)}
                            disabled={busy}
                            className="inline-flex items-center gap-1.5 px-3 py-2 text-sm font-medium rounded-lg border border-red-300 dark:border-red-800 text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors disabled:opacity-50"
                        >
                            Cancelar
                        </button>
                    )}
                    {canDelete && (
                        <button
                            onClick={() => setShowDelete(true)}
                            disabled={busy}
                            className="p-2 rounded-lg text-red-600 dark:text-red-400 hover:bg-red-50 dark:hover:bg-red-900/30 transition-colors disabled:opacity-50"
                            title="Eliminar factura"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                        </button>
                    )}
                </div>
            </div>

            <div className="invoice-print-area bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm p-6 sm:p-8">
                <div className="flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 pb-6 border-b border-gray-200 dark:border-gray-700">
                    <div>
                        <p className="text-2xl font-bold text-purple-700 dark:text-purple-300">Probability</p>
                        <p className="text-sm text-gray-500 dark:text-gray-400 mt-1">Factura de servicios</p>
                    </div>
                    <div className="sm:text-right">
                        <p className="text-lg font-semibold text-gray-900 dark:text-white">{invoice.number}</p>
                        <span className={`inline-block mt-1 px-2.5 py-0.5 rounded-full text-xs font-semibold ${invoiceStatusBadgeClass(invoice.status)}`}>
                            {invoiceStatusLabel(invoice.status)}
                        </span>
                    </div>
                </div>

                <div className="grid grid-cols-1 sm:grid-cols-2 gap-6 py-6">
                    <div>
                        <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase mb-1">Cliente</p>
                        <p className="text-sm font-medium text-gray-900 dark:text-white">
                            {invoice.business_name || `Negocio #${invoice.business_id}`}
                        </p>
                        {invoice.email_to && (
                            <p className="text-sm text-gray-500 dark:text-gray-400">{invoice.email_to}</p>
                        )}
                        <p className="text-xs text-gray-400 dark:text-gray-500 mt-1">Concepto: {invoice.concept_name}</p>
                    </div>
                    <div className="sm:text-right space-y-1 text-sm">
                        <p className="text-gray-700 dark:text-gray-300">
                            <span className="text-gray-500 dark:text-gray-400">Emision:</span>{' '}
                            <span className="font-medium">{formatEntryDate(invoice.issue_date)}</span>
                        </p>
                        {invoice.due_date && (
                            <p className="text-gray-700 dark:text-gray-300">
                                <span className="text-gray-500 dark:text-gray-400">Vencimiento:</span>{' '}
                                <span className="font-medium">{formatEntryDate(invoice.due_date)}</span>
                            </p>
                        )}
                        {invoice.sent_at && (
                            <p className="text-gray-700 dark:text-gray-300">
                                <span className="text-gray-500 dark:text-gray-400">Enviada:</span>{' '}
                                <span className="font-medium">{formatEntryDate(invoice.sent_at)}</span>
                            </p>
                        )}
                        {invoice.paid_at && (
                            <p className="text-gray-700 dark:text-gray-300">
                                <span className="text-gray-500 dark:text-gray-400">Pagada:</span>{' '}
                                <span className="font-medium">{formatEntryDate(invoice.paid_at)}</span>
                            </p>
                        )}
                    </div>
                </div>

                <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-900/40">
                                <th className="px-4 py-2 font-medium">Descripcion</th>
                                <th className="px-4 py-2 font-medium text-right">Cantidad</th>
                                <th className="px-4 py-2 font-medium text-right">Precio unitario</th>
                                <th className="px-4 py-2 font-medium text-right">Subtotal</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {(invoice.items || []).map((item) => (
                                <tr key={item.id} className="text-gray-800 dark:text-gray-200">
                                    <td className="px-4 py-2">{item.description}</td>
                                    <td className="px-4 py-2 text-right">{item.quantity}</td>
                                    <td className="px-4 py-2 text-right whitespace-nowrap">{formatCOP(item.unit_price)}</td>
                                    <td className="px-4 py-2 text-right font-medium whitespace-nowrap">{formatCOP(item.amount)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>

                <div className="flex justify-end pt-4">
                    <div className="w-full sm:w-72 space-y-1.5 text-sm">
                        <div className="flex justify-between text-gray-700 dark:text-gray-300">
                            <span>Subtotal</span>
                            <span className="font-medium">{formatCOP(invoice.subtotal)}</span>
                        </div>
                        {(invoice.tax_detail || []).map((tax) => (
                            <div key={tax.tax_code} className="flex justify-between text-gray-500 dark:text-gray-400">
                                <span>{tax.tax_name} ({tax.rate_percent}%)</span>
                                <span>{formatCOP(tax.amount)}</span>
                            </div>
                        ))}
                        <div className="flex justify-between pt-1.5 border-t border-gray-200 dark:border-gray-700 text-base font-semibold text-gray-900 dark:text-white">
                            <span>Total</span>
                            <span>{formatCOP(invoice.total)}</span>
                        </div>
                    </div>
                </div>

                {invoice.notes && (
                    <div className="mt-6 pt-4 border-t border-gray-200 dark:border-gray-700">
                        <p className="text-xs font-semibold text-gray-500 dark:text-gray-400 uppercase mb-1">Notas</p>
                        <p className="text-sm text-gray-700 dark:text-gray-300 whitespace-pre-line">{invoice.notes}</p>
                    </div>
                )}
            </div>

            <Modal isOpen={showSend} onClose={() => setShowSend(false)} title="Enviar factura por correo" size="md">
                <div className="space-y-4">
                    <p className="text-sm text-gray-600 dark:text-gray-300">
                        Se enviara la factura {invoice.number} al correo indicado. Si lo dejas vacio se usara el correo de la factura o el del primer usuario del negocio.
                    </p>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Correo destino</label>
                        <input
                            type="email"
                            value={sendEmail}
                            onChange={(e) => setSendEmail(e.target.value)}
                            placeholder="cliente@correo.com"
                            className="w-full px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        />
                    </div>
                    <div className="flex justify-end gap-3 pt-2">
                        <button
                            onClick={() => setShowSend(false)}
                            className="px-4 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors"
                        >
                            Cancelar
                        </button>
                        <button
                            onClick={handleSend}
                            disabled={sending}
                            className="px-4 py-2 text-sm font-semibold rounded-lg bg-blue-600 hover:bg-blue-700 text-white transition-colors disabled:opacity-50"
                        >
                            {sending ? 'Enviando...' : 'Enviar'}
                        </button>
                    </div>
                </div>
            </Modal>

            <ConfirmModal
                isOpen={showPay}
                onClose={() => setShowPay(false)}
                onConfirm={() => runAction(() => payAccountingInvoiceAction(invoice.id), 'Factura marcada como pagada')}
                title="Marcar factura como pagada"
                message={`Se marcara la factura ${invoice.number} como pagada por ${formatCOP(invoice.total)} y se creara el movimiento contable correspondiente.`}
                confirmText="Marcar pagada"
                type="info"
            />

            <ConfirmModal
                isOpen={showCancel}
                onClose={() => setShowCancel(false)}
                onConfirm={() => runAction(() => cancelAccountingInvoiceAction(invoice.id), 'Factura cancelada')}
                title="Cancelar factura"
                message={
                    isPaid
                        ? `La factura ${invoice.number} esta PAGADA. Al cancelarla se revertira el movimiento contable del pago. Esta accion no se puede deshacer.`
                        : `Se cancelara la factura ${invoice.number}. Esta accion no se puede deshacer.`
                }
                confirmText="Cancelar factura"
                type="danger"
            />

            <ConfirmModal
                isOpen={showDelete}
                onClose={() => setShowDelete(false)}
                onConfirm={() => runAction(() => deleteAccountingInvoiceAction(invoice.id), 'Factura eliminada', '/accounting/facturas')}
                title="Eliminar factura"
                message={`Se eliminara definitivamente la factura ${invoice.number}. Esta accion no se puede deshacer.`}
                confirmText="Eliminar"
                type="danger"
            />
        </div>
    );
}
