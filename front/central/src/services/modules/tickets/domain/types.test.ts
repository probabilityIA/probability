import { describe, it, expect } from 'vitest';
import {
    STATUS_META,
    PRIORITY_META,
    TYPE_META,
    SEVERITY_META,
    SOURCE_META,
    AREA_META,
    TICKET_STATUSES,
    TICKET_PRIORITIES,
    TICKET_TYPES,
    TICKET_AREAS,
    TICKET_SEVERITIES,
    TICKET_SOURCES,
} from './types';

const sorted = (values: readonly string[]) => [...values].sort();

describe('TICKET_STATUSES y STATUS_META', () => {
    it('declara los ocho estados en el orden del flujo', () => {
        expect(TICKET_STATUSES).toEqual([
            'open',
            'in_review',
            'in_development',
            'testing',
            'blocked',
            'resolved',
            'closed',
            'wont_fix',
        ]);
    });

    it('cada estado declarado tiene metadatos completos', () => {
        TICKET_STATUSES.forEach((status) => {
            const meta = STATUS_META[status];
            expect(meta, `falta metadato para ${status}`).toBeTruthy();
            expect(meta.label.length).toBeGreaterThan(0);
            expect(meta.color).toMatch(/^text-/);
            expect(meta.bg).toMatch(/^bg-/);
            expect(meta.ring).toMatch(/^ring-/);
        });
    });

    it('no hay metadatos huerfanos sin estado declarado', () => {
        expect(sorted(Object.keys(STATUS_META))).toEqual(sorted(TICKET_STATUSES));
    });

    it('no repite etiquetas entre estados', () => {
        const labels = TICKET_STATUSES.map((s) => STATUS_META[s].label);
        expect(new Set(labels).size).toBe(labels.length);
    });

    it('usa las etiquetas en espanol esperadas', () => {
        expect(STATUS_META.open.label).toBe('Abierto');
        expect(STATUS_META.in_review.label).toBe('En revisi\u00f3n');
        expect(STATUS_META.wont_fix.label).toBe('No se har\u00e1');
    });
});

describe('TICKET_PRIORITIES y PRIORITY_META', () => {
    it('declara las cuatro prioridades de menor a mayor', () => {
        expect(TICKET_PRIORITIES).toEqual(['low', 'medium', 'high', 'critical']);
    });

    it('cada prioridad declarada tiene metadatos completos', () => {
        TICKET_PRIORITIES.forEach((priority) => {
            const meta = PRIORITY_META[priority];
            expect(meta, `falta metadato para ${priority}`).toBeTruthy();
            expect(meta.label.length).toBeGreaterThan(0);
            expect(meta.color).toMatch(/^text-/);
            expect(meta.bg).toMatch(/^bg-/);
        });
    });

    it('no hay metadatos huerfanos sin prioridad declarada', () => {
        expect(sorted(Object.keys(PRIORITY_META))).toEqual(sorted(TICKET_PRIORITIES));
    });

    it('usa las etiquetas en espanol esperadas', () => {
        expect(PRIORITY_META.low.label).toBe('Baja');
        expect(PRIORITY_META.medium.label).toBe('Media');
        expect(PRIORITY_META.high.label).toBe('Alta');
        expect(PRIORITY_META.critical.label).toBe('Critica');
    });
});

describe('TICKET_TYPES y TYPE_META', () => {
    it('declara los nueve tipos de ticket', () => {
        expect(TICKET_TYPES).toEqual([
            'bug',
            'improvement',
            'feature',
            'data',
            'integration',
            'support',
            'complaint',
            'claim',
            'question',
        ]);
    });

    it('cada tipo declarado tiene etiqueta e icono', () => {
        TICKET_TYPES.forEach((type) => {
            const meta = TYPE_META[type];
            expect(meta, `falta metadato para ${type}`).toBeTruthy();
            expect(meta.label.length).toBeGreaterThan(0);
            expect(meta.icon).toMatch(/^[A-Z]{3}$/);
        });
    });

    it('no hay metadatos huerfanos sin tipo declarado', () => {
        expect(sorted(Object.keys(TYPE_META))).toEqual(sorted(TICKET_TYPES));
    });

    it('no repite iconos entre tipos', () => {
        const icons = TICKET_TYPES.map((t) => TYPE_META[t].icon);
        expect(new Set(icons).size).toBe(icons.length);
    });

    it('usa las etiquetas en espanol esperadas', () => {
        expect(TYPE_META.bug).toEqual({ label: 'Bug', icon: 'BUG' });
        expect(TYPE_META.integration.label).toBe('Integraci\u00f3n');
        expect(TYPE_META.feature.label).toBe('Nueva funcionalidad');
    });
});

describe('TICKET_AREAS y AREA_META', () => {
    it('declara las tres areas', () => {
        expect(TICKET_AREAS).toEqual(['comercial', 'soporte', 'desarrollo']);
    });

    it('cada area declarada tiene metadatos completos', () => {
        TICKET_AREAS.forEach((area) => {
            const meta = AREA_META[area];
            expect(meta, `falta metadato para ${area}`).toBeTruthy();
            expect(meta.label.length).toBeGreaterThan(0);
            expect(meta.color).toMatch(/^text-/);
            expect(meta.bg).toMatch(/^bg-/);
        });
    });

    it('no hay metadatos huerfanos sin area declarada', () => {
        expect(sorted(Object.keys(AREA_META))).toEqual(sorted(TICKET_AREAS));
    });

    it('usa las etiquetas en espanol esperadas', () => {
        expect(AREA_META.comercial.label).toBe('Comercial');
        expect(AREA_META.soporte.label).toBe('Soporte');
        expect(AREA_META.desarrollo.label).toBe('Desarrollo');
    });
});

describe('TICKET_SEVERITIES y SEVERITY_META', () => {
    it('incluye la severidad vacia como primera opcion', () => {
        expect(TICKET_SEVERITIES).toEqual(['', 'low', 'medium', 'high']);
    });

    it('cada severidad declarada tiene etiqueta', () => {
        TICKET_SEVERITIES.forEach((severity) => {
            expect(SEVERITY_META[severity], `falta metadato para "${severity}"`).toBeTruthy();
            expect(SEVERITY_META[severity].label.length).toBeGreaterThan(0);
        });
    });

    it('no hay metadatos huerfanos sin severidad declarada', () => {
        expect(sorted(Object.keys(SEVERITY_META))).toEqual(sorted(TICKET_SEVERITIES));
    });

    it('etiqueta la severidad vacia como sin severidad', () => {
        expect(SEVERITY_META[''].label).toBe('Sin severidad');
    });
});

describe('TICKET_SOURCES y SOURCE_META', () => {
    it('declara los dos origenes', () => {
        expect(TICKET_SOURCES).toEqual(['internal', 'business']);
    });

    it('cada origen declarado tiene etiqueta', () => {
        TICKET_SOURCES.forEach((source) => {
            expect(SOURCE_META[source], `falta metadato para ${source}`).toBeTruthy();
            expect(SOURCE_META[source].label.length).toBeGreaterThan(0);
        });
    });

    it('no hay metadatos huerfanos sin origen declarado', () => {
        expect(sorted(Object.keys(SOURCE_META))).toEqual(sorted(TICKET_SOURCES));
    });

    it('usa las etiquetas en espanol esperadas', () => {
        expect(SOURCE_META.internal.label).toBe('Interno');
        expect(SOURCE_META.business.label).toBe('Negocio');
    });
});
