import type { RangeKey } from '../../domain/types';

export function formatMoney(amount?: number, currency: string = 'COP'): string {
    if (amount == null || isNaN(amount)) return '$0';
    try {
        return new Intl.NumberFormat('es-CO', { style: 'currency', currency, maximumFractionDigits: 0 }).format(amount);
    } catch {
        return `$${Math.round(amount).toLocaleString('es-CO')}`;
    }
}

export function formatMoneyShort(amount?: number): string {
    if (amount == null || isNaN(amount)) return '0';
    const abs = Math.abs(amount);
    if (abs >= 1_000_000) return `${(amount / 1_000_000).toFixed(1)}M`;
    if (abs >= 1_000) return `${(amount / 1_000).toFixed(0)}k`;
    return `${Math.round(amount)}`;
}

export function formatDate(s?: string | null): string {
    if (!s) return '-';
    const d = new Date(s);
    if (isNaN(d.getTime())) return '-';
    return d.toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric', timeZone: 'America/Bogota' });
}

const MONTHS_ES = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic'];

export function formatDateOnly(s?: string | null): string {
    if (!s) return '-';
    const [y, m, d] = s.slice(0, 10).split('-').map(Number);
    if (!y || !m || !d) return formatDate(s);
    return `${String(d).padStart(2, '0')} ${MONTHS_ES[m - 1]} ${y}`;
}

export function carrierLabel(c?: string): string {
    if (!c) return 'Sin transportadora';
    return c.charAt(0).toUpperCase() + c.slice(1).toLowerCase();
}

export const RANGE_OPTIONS: { key: RangeKey; label: string }[] = [
    { key: 'today', label: 'Hoy' },
    { key: 'week', label: 'Semana' },
    { key: 'month', label: 'Mes' },
    { key: '3months', label: '3 meses' },
];

export const CARRIER_COLORS = ['#7c3aed', '#0ea5e9', '#f59e0b', '#10b981', '#ef4444', '#ec4899', '#6366f1', '#14b8a6'];

export function carrierColor(index: number): string {
    return CARRIER_COLORS[index % CARRIER_COLORS.length];
}
