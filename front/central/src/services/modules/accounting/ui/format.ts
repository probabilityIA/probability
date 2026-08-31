const copFormatter = new Intl.NumberFormat('es-CO', {
    style: 'currency',
    currency: 'COP',
    maximumFractionDigits: 0,
});

export function formatCOP(value: number): string {
    return copFormatter.format(value || 0);
}

export function formatEntryDate(value: string): string {
    if (!value) return '-';
    return value.slice(0, 10);
}

export function firstDayOfCurrentMonth(): string {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`;
}

export function lastDayOfCurrentMonth(): string {
    const now = new Date();
    const last = new Date(now.getFullYear(), now.getMonth() + 1, 0);
    return `${last.getFullYear()}-${String(last.getMonth() + 1).padStart(2, '0')}-${String(last.getDate()).padStart(2, '0')}`;
}

export function kindLabel(kind: string): string {
    if (kind === 'INCOME') return 'Ingreso';
    if (kind === 'EXPENSE') return 'Gasto';
    return kind;
}

export function taxKindLabel(kind: string): string {
    if (kind === 'CHARGE') return 'Cargo';
    if (kind === 'WITHHOLDING') return 'Retención';
    if (kind === 'OTHER') return 'Otro';
    return kind;
}

export function taxKindBadgeClass(kind: string): string {
    if (kind === 'CHARGE') return 'bg-purple-100 text-purple-700 dark:bg-purple-900/40 dark:text-purple-300';
    if (kind === 'WITHHOLDING') return 'bg-amber-100 text-amber-700 dark:bg-amber-900/40 dark:text-amber-300';
    return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300';
}

export function invoiceStatusLabel(status: string): string {
    if (status === 'DRAFT') return 'Borrador';
    if (status === 'SENT') return 'Enviada';
    if (status === 'PAID') return 'Pagada';
    if (status === 'CANCELLED') return 'Cancelada';
    return status;
}

export function invoiceStatusBadgeClass(status: string): string {
    if (status === 'SENT') return 'bg-blue-100 text-blue-700 dark:bg-blue-900/40 dark:text-blue-300';
    if (status === 'PAID') return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-900/40 dark:text-emerald-300';
    if (status === 'CANCELLED') return 'bg-red-100 text-red-700 dark:bg-red-900/40 dark:text-red-300';
    return 'bg-gray-100 text-gray-700 dark:bg-gray-700 dark:text-gray-300';
}

export function todayISO(): string {
    const now = new Date();
    return `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`;
}
