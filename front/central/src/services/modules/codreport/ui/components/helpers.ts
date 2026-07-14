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
    return d.toLocaleDateString('es-CO', { day: '2-digit', month: 'short', year: 'numeric' });
}

export function formatDateTime(s?: string | null): string {
    if (!s) return '-';
    const d = new Date(s);
    if (isNaN(d.getTime())) return '-';
    return d.toLocaleString('es-CO', {
        day: '2-digit',
        month: 'short',
        year: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: true,
    });
}

export function browserTimeZone(): string {
    try {
        return Intl.DateTimeFormat().resolvedOptions().timeZone || '';
    } catch {
        return '';
    }
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

const RANGE_DAYS: Record<string, number> = { today: 1, week: 7, month: 30, '3months': 90 };

export function resolveRangeDates(
    range: RangeKey,
    customStart?: string,
    customEnd?: string,
): { start_date?: string; end_date?: string } {
    if (range === 'custom') {
        if (!customStart || !customEnd) return {};
        const [ys, ms, ds] = customStart.split('-').map(Number);
        const [ye, me, de] = customEnd.split('-').map(Number);
        if (!ys || !ms || !ds || !ye || !me || !de) return {};
        return {
            start_date: new Date(ys, ms - 1, ds, 0, 0, 0, 0).toISOString(),
            end_date: new Date(ye, me - 1, de, 23, 59, 59, 999).toISOString(),
        };
    }
    const days = RANGE_DAYS[range];
    if (!days) return {};
    const now = new Date();
    const start = new Date(now.getFullYear(), now.getMonth(), now.getDate() - (days - 1), 0, 0, 0, 0);
    const end = new Date(now.getFullYear(), now.getMonth(), now.getDate(), 23, 59, 59, 999);
    return { start_date: start.toISOString(), end_date: end.toISOString() };
}

export const CARRIER_COLORS =['#7c3aed', '#0ea5e9', '#f59e0b', '#10b981', '#ef4444', '#ec4899', '#6366f1', '#14b8a6'];

export function carrierColor(index: number): string {
    return CARRIER_COLORS[index % CARRIER_COLORS.length];
}
