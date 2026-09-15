import { describe, expect, it } from 'vitest';
import { MAX_HISTORY, buildHistory, hubEnvironmentFor, isCurrentRoute, isHubDestination, pickSuggestions, rateLimitText } from './use-cases';

describe('destinos de Tus Integraciones', () => {
    it('reconoce el hub y traduce cada accion a su ambiente', () => {
        expect(isHubDestination('integrations.hub')).toBe(true);
        expect(isHubDestination('integrations')).toBe(false);
        expect(hubEnvironmentFor('integrations.hub')).toBeNull();
        expect(hubEnvironmentFor('integrations.hub.inventory')).toBe('inventory');
        expect(hubEnvironmentFor('integrations.hub.orders')).toBe('orders_compare');
    });
});
import type { ChatEntry } from '../domain/types';

describe('buildHistory', () => {
    it('solo manda turnos de usuario y asistente, sin avisos ni lineas de sistema', () => {
        const entries: ChatEntry[] = [
            { id: '1', kind: 'assistant', text: 'Hola', destination: null },
            { id: '2', kind: 'user', text: 'ordenes' },
            { id: '3', kind: 'system', text: 'Te llev\u00e9 a \u00d3rdenes' },
            { id: '4', kind: 'notice', text: 'fallo', tone: 'error', retryable: true },
        ];
        expect(buildHistory(entries)).toEqual([
            { role: 'assistant', text: 'Hola' },
            { role: 'user', text: 'ordenes' },
        ]);
    });

    it('recorta a la historia reciente', () => {
        const entries: ChatEntry[] = Array.from({ length: 30 }, (_, i) => ({ id: String(i), kind: 'user', text: `m${i}` }));
        const history = buildHistory(entries);
        expect(history).toHaveLength(MAX_HISTORY);
        expect(history[history.length - 1].text).toBe('m29');
    });
});

describe('isCurrentRoute', () => {
    it('reconoce la ruta y sus subrutas sin confundir prefijos', () => {
        expect(isCurrentRoute('/orders', '/orders')).toBe(true);
        expect(isCurrentRoute('/shipments/cod', '/shipments')).toBe(true);
        expect(isCurrentRoute('/ordersx', '/orders')).toBe(false);
    });
});

describe('pickSuggestions', () => {
    it('solo sugiere preguntas de modulos que el usuario puede ver', () => {
        const suggestions = pickSuggestions((key) => key === 'wallet' || key === 'home');
        expect(suggestions).toHaveLength(2);
        expect(suggestions[0]).toContain('billetera');
    });
});

describe('rateLimitText', () => {
    it('sin hora de reinicio no inventa una', () => {
        expect(rateLimitText(30, null)).toBe('Llegaste a 30 mensajes en una hora. Intenta de nuevo m\u00e1s tarde.');
    });
});
