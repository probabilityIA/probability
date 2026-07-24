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
