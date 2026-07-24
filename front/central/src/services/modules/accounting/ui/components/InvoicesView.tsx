'use client';

import { useState, useTransition } from 'react';
import { useRouter } from 'next/navigation';
import { SuperAdminBusinessSelector, Table } from '@/shared/ui';
import { Concept, InvoicesListResponse, InvoiceStatus, Tax } from '../../domain/types';
import { formatCOP, formatEntryDate, invoiceStatusBadgeClass, invoiceStatusLabel } from '../format';
import { InvoiceFormModal } from './InvoiceFormModal';

interface InvoicesFilters {
    page: number;
    page_size: number;
    status: InvoiceStatus | null;
    business_id: number | null;
    is_test: 'true' | 'false' | null;
}

interface InvoicesViewProps {
    invoices: InvoicesListResponse;
    concepts: Concept[];
    taxes: Tax[];
    filters: InvoicesFilters;
    error: string | null;
    openCreate?: boolean;
}

const STATUSES: InvoiceStatus[] = ['DRAFT', 'SENT', 'PAID', 'CANCELLED'];

function buildListUrl(filters: InvoicesFilters): string {
    const params = new URLSearchParams();
    params.set('page', String(filters.page));
    params.set('page_size', String(filters.page_size));
    if (filters.status) params.set('status', filters.status);
    if (filters.business_id) params.set('business_id', String(filters.business_id));
    if (filters.is_test) params.set('is_test', filters.is_test);
    return `/accounting/facturas?${params.toString()}`;
}

export function InvoicesView({ invoices, concepts, taxes, filters, error, openCreate = false }: InvoicesViewProps) {
    const router = useRouter();
    const [isPending, startTransition] = useTransition();
    const [showForm, setShowForm] = useState(openCreate);
    const [fromQuery] = useState(openCreate);

    const navigate = (next: Partial<InvoicesFilters>) => {
        startTransition(() => {
            router.replace(buildListUrl({ ...filters, ...next }));
        });
    };

    const closeForm = () => {
        setShowForm(false);
        if (fromQuery) {
            router.replace(buildListUrl(filters));
        }
    };

    const hasActiveFilters = Boolean(filters.status || filters.business_id || filters.is_test);
    const totalPages = invoices.total_pages || 0;

    const columns = [
        { key: 'number', label: 'Numero' },
        { key: 'business', label: 'Negocio' },
        { key: 'concept', label: 'Concepto' },
        { key: 'issue_date', label: 'Emision' },
        { key: 'due_date', label: 'Vencimiento' },
        { key: 'status', label: 'Estado' },
        { key: 'total', label: 'Total', align: 'right' as const },
    ];

    const rows = invoices.data.map((invoice) => ({
        number: (
            <span className="inline-flex items-center gap-1.5 whitespace-nowrap">
                <span className="font-medium text-purple-700 dark:text-purple-300">{invoice.number}</span>
                {invoice.is_test && (
                    <span className="inline-block px-1.5 py-0.5 rounded text-[10px] font-bold bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300">
                        TEST
                    </span>
                )}
            </span>
        ),
        business: (
            <span className="text-gray-800 dark:text-gray-200">{invoice.business_name || `#${invoice.business_id}`}</span>
        ),
        concept: (
            <span className="text-gray-800 dark:text-gray-200">{invoice.concept_name}</span>
        ),
        issue_date: (
            <span className="text-gray-800 dark:text-gray-200 whitespace-nowrap">{formatEntryDate(invoice.issue_date)}</span>
        ),
        due_date: (
            <span className="text-gray-800 dark:text-gray-200 whitespace-nowrap">
                {invoice.due_date ? formatEntryDate(invoice.due_date) : '-'}
            </span>
        ),
        status: (
            <span className={`inline-block px-2 py-0.5 rounded-full text-xs font-semibold ${invoiceStatusBadgeClass(invoice.status)}`}>
                {invoiceStatusLabel(invoice.status)}
            </span>
        ),
        total: (
            <span className="font-medium text-gray-900 dark:text-white whitespace-nowrap">{formatCOP(invoice.total)}</span>
        ),
    }));

    return (
        <div className="space-y-4">
            <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3">
                <div>
                    <h1 className="text-xl font-semibold text-gray-900 dark:text-white">Facturas</h1>
                    <p className="text-sm text-gray-500 dark:text-gray-400 mt-0.5">
                        {invoices.total} facturas emitidas
                    </p>
                </div>
                <button
                    onClick={() => setShowForm(true)}
                    className="inline-flex items-center gap-2 px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors"
                >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    Nueva factura
                </button>
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
                        <label className="block text-xs font-medium text-gray-500 dark:text-gray-400 mb-1">Tipo</label>
                        <select
                            value={filters.is_test ?? ''}
                            onChange={(e) => navigate({ is_test: (e.target.value || null) as 'true' | 'false' | null, page: 1 })}
                            className="px-3 py-2 text-sm border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 focus:outline-none focus:ring-2 focus:ring-purple-500"
                        >
                            <option value="">Todas</option>
                            <option value="false">Reales</option>
                            <option value="true">De prueba</option>
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
                    {hasActiveFilters && (
                        <button
                            onClick={() => navigate({ status: null, business_id: null, is_test: null, page: 1 })}
                            disabled={isPending}
                            className="px-4 py-2 text-sm font-medium rounded-lg border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 hover:bg-gray-50 dark:hover:bg-gray-700 transition-colors disabled:opacity-50"
                        >
                            Limpiar filtros
                        </button>
                    )}
                </div>
            </div>

            {error && (
                <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 text-red-700 dark:text-red-300 text-sm rounded-lg px-4 py-3">
                    {error}
                </div>
            )}

            {invoices.data.length === 0 && !isPending ? (
                <div className="bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-sm">
                    <div className="flex flex-col items-center justify-center py-16 text-center px-4">
                        <svg className="w-12 h-12 text-gray-300 dark:text-gray-600 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                        </svg>
                        <p className="text-gray-500 dark:text-gray-400 text-sm mb-4">
                            {hasActiveFilters
                                ? 'No hay facturas con los filtros seleccionados'
                                : 'Aun no has emitido facturas'}
                        </p>
                        <button
                            onClick={() => setShowForm(true)}
                            className="inline-flex items-center gap-2 px-4 py-2 text-sm font-semibold rounded-lg bg-purple-600 hover:bg-purple-700 text-white transition-colors"
                        >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                            Nueva factura
                        </button>
                    </div>
                </div>
            ) : (
                <Table
                    columns={columns}
                    data={rows}
                    keyExtractor={(_, index) => invoices.data[index]?.id ?? index}
                    loading={isPending}
                    emptyMessage="No hay facturas con los filtros seleccionados"
                    onRowClick={(_, index) => {
                        const invoice = invoices.data[index];
                        if (invoice) router.push(`/accounting/facturas/${invoice.id}`);
                    }}
                    pagination={{
                        currentPage: filters.page,
                        totalPages: Math.max(totalPages, 1),
                        totalItems: invoices.total,
                        itemsPerPage: filters.page_size,
                        onPageChange: (page) => navigate({ page }),
                        onItemsPerPageChange: (pageSize) => navigate({ page_size: pageSize, page: 1 }),
                        itemsPerPageOptions: [10, 25, 50, 100],
                    }}
                />
            )}

            {showForm && (
                <InvoiceFormModal
                    isOpen
                    onClose={closeForm}
                    concepts={concepts}
                    taxes={taxes}
                    onSaved={(invoice) => router.push(`/accounting/facturas/${invoice.id}`)}
                />
            )}
        </div>
    );
}
