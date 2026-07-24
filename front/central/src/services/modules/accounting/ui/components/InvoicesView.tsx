'use client';

import { useTransition } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/navigation';
import { SuperAdminBusinessSelector } from '@/shared/ui';
import { InvoicesListResponse, InvoiceStatus } from '../../domain/types';
import { formatCOP, formatEntryDate, invoiceStatusBadgeClass, invoiceStatusLabel } from '../format';

interface InvoicesFilters {
    page: number;
    page_size: number;
    status: InvoiceStatus | null;
    business_id: number | null;
}

interface InvoicesViewProps {
    invoices: InvoicesListResponse;
    filters: InvoicesFilters;
    error: string | null;
}

const STATUSES: InvoiceStatus[] = ['DRAFT', 'SENT', 'PAID', 'CANCELLED'];

export function InvoicesView({ invoices, filters, error }: InvoicesViewProps) {
    const router = useRouter();
    const [isPending, startTransition] = useTransition();

    const navigate = (next: Partial<InvoicesFilters>) => {
        const merged = { ...filters, ...next };
        const params = new URLSearchParams();
        params.set('page', String(merged.page));
        params.set('page_size', String(merged.page_size));
        if (merged.status) params.set('status', merged.status);
        if (merged.business_id) params.set('business_id', String(merged.business_id));
        startTransition(() => {
            router.replace(`/accounting/facturas?${params.toString()}`);
        });
    };

    const totalPages = invoices.total_pages || 0;

    return (
        <div className="space-y-4">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Facturas</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        {invoices.total} facturas emitidas
                    </p>
                </div>
                <Link
                    href="/accounting/facturas/nueva"
                    className="inline-flex items-center gap-2 px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    Nueva factura
                </Link>
            </div>

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm p-4">
                <div className="flex flex-wrap items-end gap-3">
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Estado</label>
                        <select
                            value={filters.status ?? ''}
                            onChange={(e) => navigate({ status: (e.target.value || null) as InvoiceStatus | null, page: 1 })}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        >
                            <option value="">Todos</option>
                            {STATUSES.map((status) => (
                                <option key={status} value={status}>{invoiceStatusLabel(status)}</option>
                            ))}
                        </select>
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Negocio</label>
                        <SuperAdminBusinessSelector
                            value={filters.business_id}
                            onChange={(businessId) => navigate({ business_id: businessId, page: 1 })}
                            variant="default"
                            placeholder="Todos los negocios"
                        />
                    </div>
                </div>
            </div>

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 text-sm rounded-lg px-4 py-3">
                    {error}
                </div>
            )}

            <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm overflow-hidden">
                <div className="overflow-x-auto">
                    <table className="w-full text-sm">
                        <thead>
                            <tr className="text-left text-xs text-gray-500 dark:text-gray-400 uppercase bg-gray-50 dark:bg-gray-900/40">
                                <th className="px-4 py-2 font-medium">Numero</th>
                                <th className="px-4 py-2 font-medium">Negocio</th>
                                <th className="px-4 py-2 font-medium">Concepto</th>
                                <th className="px-4 py-2 font-medium">Emision</th>
                                <th className="px-4 py-2 font-medium">Vencimiento</th>
                                <th className="px-4 py-2 font-medium">Estado</th>
                                <th className="px-4 py-2 font-medium text-right">Total</th>
                            </tr>
                        </thead>
                        <tbody className="divide-y divide-gray-100 dark:divide-gray-700">
                            {invoices.data.length === 0 && (
                                <tr>
                                    <td colSpan={7} className="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                                        No hay facturas con los filtros seleccionados
                                    </td>
                                </tr>
                            )}
                            {invoices.data.map((invoice) => (
                                <tr
                                    key={invoice.id}
                                    onClick={() => router.push(`/accounting/facturas/${invoice.id}`)}
                                    className="text-gray-800 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700/40 cursor-pointer transition-colors"
                                >
                                    <td className="px-4 py-2 whitespace-nowrap">
                                        <span className="font-medium text-purple-700 dark:text-purple-300">{invoice.number}</span>
                                    </td>
                                    <td className="px-4 py-2">{invoice.business_name || `#${invoice.business_id}`}</td>
                                    <td className="px-4 py-2">{invoice.concept_name}</td>
                                    <td className="px-4 py-2 whitespace-nowrap">{formatEntryDate(invoice.issue_date)}</td>
                                    <td className="px-4 py-2 whitespace-nowrap">{invoice.due_date ? formatEntryDate(invoice.due_date) : '-'}</td>
                                    <td className="px-4 py-2">
                                        <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${invoiceStatusBadgeClass(invoice.status)}`}>
                                            {invoiceStatusLabel(invoice.status)}
                                        </span>
                                    </td>
                                    <td className="px-4 py-2 text-right font-medium whitespace-nowrap">{formatCOP(invoice.total)}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </div>

                <div className="flex flex-col sm:flex-row items-center justify-between gap-3 px-4 py-3 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center gap-2 text-sm text-gray-600 dark:text-gray-300">
                        <span>Por pagina:</span>
                        <select
                            value={filters.page_size}
                            onChange={(e) => navigate({ page_size: Number(e.target.value), page: 1 })}
                            className="px-2 py-1 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                        >
                            {[10, 25, 50, 100].map((size) => (
                                <option key={size} value={size}>{size}</option>
                            ))}
                        </select>
                    </div>
                    <div className="flex items-center gap-2">
                        <button
                            onClick={() => navigate({ page: filters.page - 1 })}
                            disabled={filters.page <= 1 || isPending}
                            className="px-3 py-1.5 text-sm rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 transition-colors"
                        >
                            Anterior
                        </button>
                        <span className="text-sm text-gray-600 dark:text-gray-300">
                            Pagina {filters.page} de {Math.max(totalPages, 1)}
                        </span>
                        <button
                            onClick={() => navigate({ page: filters.page + 1 })}
                            disabled={filters.page >= totalPages || isPending}
                            className="px-3 py-1.5 text-sm rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 disabled:opacity-40 transition-colors"
                        >
                            Siguiente
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
}
