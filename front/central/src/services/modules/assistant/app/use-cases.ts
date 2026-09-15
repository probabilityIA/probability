import type { AssistantHistoryMessage, ChatEntry } from '../domain/types';

export const MAX_HISTORY = 12;
export const INTRO_DELAY_MS = 20000;
export const ASSISTANT_NAME = 'V\u00eda';

export const GREETING = `Hola, soy ${ASSISTANT_NAME}. Te ayudo a encontrar cualquier cosa en Probability. \u00bfQu\u00e9 necesitas?`;

const SUGGESTIONS: { navKey: string; text: string }[] = [
    { navKey: 'orders', text: '\u00bfD\u00f3nde veo las \u00f3rdenes?' },
    { navKey: 'orders', text: '\u00bfC\u00f3mo genero una gu\u00eda de env\u00edo?' },
    { navKey: 'shipments', text: '\u00bfD\u00f3nde veo el recaudo contra entrega?' },
    { navKey: 'wallet', text: '\u00bfC\u00f3mo recargo la billetera?' },
    { navKey: 'integrations', text: '\u00bfD\u00f3nde conecto mi tienda?' },
    { navKey: 'products', text: '\u00bfD\u00f3nde creo un producto?' },
    { navKey: 'home', text: '\u00bfQu\u00e9 puedo hacer en Probability?' },
];

export function pickSuggestions(hasNav: (key: string) => boolean, max = 4): string[] {
    return SUGGESTIONS.filter((s) => hasNav(s.navKey)).slice(0, max).map((s) => s.text);
}

export function buildHistory(entries: ChatEntry[]): AssistantHistoryMessage[] {
    const history: AssistantHistoryMessage[] = [];
    for (const entry of entries) {
        if (entry.kind === 'user') history.push({ role: 'user', text: entry.text });
        if (entry.kind === 'assistant') history.push({ role: 'assistant', text: entry.text });
    }
    return history.slice(-MAX_HISTORY);
}

export function isCurrentRoute(pathname: string, route: string): boolean {
    return pathname === route || pathname.startsWith(`${route}/`);
}

export type HubEnvironment = 'products' | 'data' | 'inventory' | 'orders_compare';

const HUB_DESTINATIONS: Record<string, HubEnvironment | null> = {
    'integrations.hub': null,
    'integrations.hub.products': 'products',
    'integrations.hub.data': 'data',
    'integrations.hub.inventory': 'inventory',
    'integrations.hub.orders': 'orders_compare',
};

export function isHubDestination(key: string): boolean {
    return key in HUB_DESTINATIONS;
}

export function hubEnvironmentFor(key: string): HubEnvironment | null {
    return HUB_DESTINATIONS[key] ?? null;
}

export function formatResetTime(iso: string | null | undefined): string | null {
    if (!iso) return null;
    const date = new Date(iso);
    if (Number.isNaN(date.getTime())) return null;
    return date.toLocaleTimeString('es-CO', { hour: 'numeric', minute: '2-digit' });
}

export function rateLimitText(limit: number, resetAt: string | null | undefined): string {
    const time = formatResetTime(resetAt);
    const base = `Llegaste a ${limit} mensajes en una hora.`;
    return time ? `${base} Puedes volver a escribir a las ${time}.` : `${base} Intenta de nuevo m\u00e1s tarde.`;
}

export function newConversationId(): string {
    if (typeof crypto !== 'undefined' && typeof crypto.randomUUID === 'function') return crypto.randomUUID();
    return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
        const r = Math.floor(Math.random() * 16);
        return (c === 'x' ? r : (r % 4) + 8).toString(16);
    });
}

export function nextFeedback(current: number | undefined, pressed: 1 | -1): -1 | 0 | 1 {
    return current === pressed ? 0 : pressed;
}

export function percent(part: number, total: number): string {
    if (!total) return '0 %';
    return `${Math.round((part / total) * 100)} %`;
}

export function reviewRange(days: number, now: Date = new Date()): { from: string; to: string } {
    const fmt = (d: Date) => `${d.getFullYear()}-${String(d.getMonth() + 1).padStart(2, '0')}-${String(d.getDate()).padStart(2, '0')}`;
    const start = new Date(now);
    start.setDate(start.getDate() - (days - 1));
    return { from: fmt(start), to: fmt(now) };
}

let sequence = 0;
export function entryId(): string {
    sequence += 1;
    return `${Date.now().toString(36)}-${sequence}`;
}
