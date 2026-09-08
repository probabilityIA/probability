import { Ticket, TicketStatus } from '../../domain/types';

export type MetaTone = 'danger' | 'warn' | 'muted' | 'ok';

export interface TicketMeta {
    label: string;
    tone: MetaTone;
}

export const META_TONE: Record<MetaTone, string> = {
    danger: 'text-red-600 dark:text-red-400 font-semibold',
    warn: 'text-amber-600 dark:text-amber-400 font-semibold',
    ok: 'text-emerald-600 dark:text-emerald-400 font-semibold',
    muted: 'text-gray-500 dark:text-gray-400 font-medium',
};

const CLOSED_STATUSES: TicketStatus[] = ['resolved', 'closed', 'wont_fix'];

const HOUR = 3600000;
const DAY = 86400000;

export const isClosedStatus = (status: TicketStatus) => CLOSED_STATUSES.includes(status);

export const initials = (name?: string) => {
    if (!name) return '?';
    const parts = name.trim().split(/\s+/);
    const first = parts[0]?.[0] || '';
    const last = parts.length > 1 ? parts[parts.length - 1][0] : '';
    return (first + last).toUpperCase();
};

const elapsed = (iso?: string | null, now = Date.now()) => {
    if (!iso) return null;
    const ms = new Date(iso).getTime();
    if (Number.isNaN(ms)) return null;
    return now - ms;
};

export const ageLabel = (iso?: string | null, now = Date.now()) => {
    const ms = elapsed(iso, now);
    if (ms === null) return '';
    if (ms < HOUR) return 'ahora';
    if (ms < DAY) return 'hace ' + Math.floor(ms / HOUR) + 'h';
    return 'hace ' + Math.floor(ms / DAY) + 'd';
};

export const slaMeta = (ticket: Ticket, now = Date.now()): TicketMeta | null => {
    if (!ticket.due_date) return null;
    const due = new Date(ticket.due_date).getTime();
    if (Number.isNaN(due)) return null;
    const left = due - now;
    if (left < 0) {
        const over = Math.abs(left);
        const label = over < DAY ? '-' + Math.max(1, Math.floor(over / HOUR)) + 'h' : '-' + Math.floor(over / DAY) + 'd';
        return { label: label, tone: 'danger' };
    }
    if (left < DAY) return { label: 'SLA ' + Math.max(1, Math.floor(left / HOUR)) + 'h', tone: 'danger' };
    const days = Math.ceil(left / DAY);
    return { label: 'SLA ' + days + 'd', tone: days <= 3 ? 'warn' : 'muted' };
};

export const isDueSoon = (ticket: Ticket, now = Date.now()) => {
    const meta = slaMeta(ticket, now);
    return !!meta && (meta.tone === 'danger' || meta.tone === 'warn');
};

const stalledMeta = (ticket: Ticket, now = Date.now()): TicketMeta | null => {
    const ms = elapsed(ticket.updated_at, now);
    if (ms === null) return null;
    const days = Math.floor(ms / DAY);
    if (days >= 7) return { label: days + 'd en col.', tone: 'danger' };
    if (days >= 3) return { label: days + 'd en col.', tone: 'warn' };
    return null;
};

export const cardMeta = (ticket: Ticket, trackStale: boolean, now = Date.now()): TicketMeta => {
    if (isClosedStatus(ticket.status)) return { label: '\u2713 Resuelto', tone: 'ok' };
    const sla = slaMeta(ticket, now);
    if (sla) return sla;
    if (trackStale) {
        const stalled = stalledMeta(ticket, now);
        if (stalled) return stalled;
    }
    return { label: ageLabel(ticket.created_at, now), tone: 'muted' };
};

export const slaRank = (ticket: Ticket, now = Date.now()) => {
    if (isClosedStatus(ticket.status)) return Number.MAX_SAFE_INTEGER;
    if (!ticket.due_date) return Number.MAX_SAFE_INTEGER - 1;
    const due = new Date(ticket.due_date).getTime();
    return Number.isNaN(due) ? Number.MAX_SAFE_INTEGER - 1 : due - now;
};

export const bySla = (now = Date.now()) => (a: Ticket, b: Ticket) => {
    const diff = slaRank(a, now) - slaRank(b, now);
    if (diff !== 0) return diff;
    const aAt = new Date(a.created_at).getTime() || 0;
    const bAt = new Date(b.created_at).getTime() || 0;
    return aAt - bAt;
};

export const channelLabel = (ticket: Ticket) =>
    ticket.business_name || (ticket.business_id ? '#' + ticket.business_id : 'Interno');
