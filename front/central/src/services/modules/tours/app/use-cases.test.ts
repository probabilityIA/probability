import { afterEach, describe, expect, it } from 'vitest';
import { findTourForRoute, resolveVisibleSteps, shouldAutoStart } from './use-cases';
import { TOUR_LIST } from '../content';
import type { TourDefinition, TourProgress, TourStatus } from '../domain/types';

function tour(parcial: Partial<TourDefinition>): TourDefinition {
    return {
        key: 'demo',
        version: 1,
        title: 'Demo',
        routes: ['/demo'],
        autoStart: true,
        steps: [{ id: 'a', title: 'A', body: 'b' }],
        ...parcial,
    };
}

function progreso(status: TourStatus, version = 1): TourProgress {
    return { tour_key: 'demo', version, status, step_index: 0 };
}

describe('findTourForRoute', () => {
    const definiciones = [
        tour({ key: 'shipments', routes: ['/shipments'] }),
        tour({ key: 'shipments-cod', routes: ['/shipments/cod'] }),
        tour({ key: 'orders', routes: ['/orders'] }),
    ];

    it('devuelve el tour de la ruta exacta', () => {
        expect(findTourForRoute(definiciones, '/orders')?.key).toBe('orders');
    });

    it('prefiere la ruta mas especifica cuando dos hacen match', () => {
        expect(findTourForRoute(definiciones, '/shipments/cod')?.key).toBe('shipments-cod');
    });

    it('hace match con subrutas del mismo modulo', () => {
        expect(findTourForRoute(definiciones, '/shipments/123')?.key).toBe('shipments');
    });

    it('no hace match parcial de un segmento', () => {
        expect(findTourForRoute(definiciones, '/orders-compare')).toBeUndefined();
    });

    it('devuelve undefined cuando ninguna ruta aplica', () => {
        expect(findTourForRoute(definiciones, '/profile')).toBeUndefined();
    });
});

describe('shouldAutoStart', () => {
    it('arranca para un usuario nuevo sin progreso', () => {
        expect(shouldAutoStart(tour({}), undefined)).toBe(true);
    });

    it('nunca vuelve a arrancar si el usuario lo omitio, aunque suba la version', () => {
        expect(shouldAutoStart(tour({ version: 5 }), progreso('skipped', 1))).toBe(false);
    });

    it('no repite un tour completado con la misma version', () => {
        expect(shouldAutoStart(tour({ version: 2 }), progreso('completed', 2))).toBe(false);
    });

    it('vuelve a mostrar un tour completado cuando sube la version', () => {
        expect(shouldAutoStart(tour({ version: 3 }), progreso('completed', 2))).toBe(true);
    });

    it('retoma un tour que quedo a medias', () => {
        expect(shouldAutoStart(tour({}), progreso('in_progress'))).toBe(true);
    });

    it('respeta autoStart en false', () => {
        expect(shouldAutoStart(tour({ autoStart: false }), undefined)).toBe(false);
    });
});

describe('registro de tours', () => {
    it('no tiene claves duplicadas', () => {
        const claves = TOUR_LIST.map((t) => t.key);
        expect(new Set(claves).size).toBe(claves.length);
    });

    it('cada tour tiene al menos un paso y version valida', () => {
        for (const t of TOUR_LIST) {
            expect(t.steps.length, `${t.key} sin pasos`).toBeGreaterThan(0);
            expect(t.version, `${t.key} con version invalida`).toBeGreaterThan(0);
            expect(t.routes.length, `${t.key} sin rutas`).toBeGreaterThan(0);
        }
    });

    it('los ids de paso son unicos dentro de cada tour', () => {
        for (const t of TOUR_LIST) {
            const ids = t.steps.map((s) => s.id);
            expect(new Set(ids).size, `${t.key} tiene ids repetidos`).toBe(ids.length);
        }
    });

    it('el primer paso de cada tour es de bienvenida sin ancla', () => {
        for (const t of TOUR_LIST) {
            expect(t.steps[0].target, `${t.key} arranca anclado`).toBeUndefined();
        }
    });

    it('toda ruta declarada empieza con barra', () => {
        for (const t of TOUR_LIST) {
            for (const ruta of t.routes) {
                expect(ruta.startsWith('/'), `${t.key}: ruta invalida ${ruta}`).toBe(true);
            }
        }
    });

    it('cada ruta resuelve al tour que la declara', () => {
        for (const t of TOUR_LIST) {
            for (const ruta of t.routes) {
                const encontrado = findTourForRoute(TOUR_LIST, ruta);
                expect(encontrado, `${ruta} no resuelve a ningun tour`).toBeDefined();
                expect(encontrado?.routes).toContain(ruta);
            }
        }
    });
});

describe('resolveVisibleSteps', () => {
    afterEach(() => {
        document.body.innerHTML = '';
    });

    it('descarta los pasos opcionales cuyo elemento no existe en la pagina', () => {
        document.body.innerHTML = '<a href="/orders">Ordenes</a>';

        const definicion = tour({
            steps: [
                { id: 'welcome', title: 'Hola', body: 'b' },
                { id: 'orders', title: 'Ordenes', body: 'b', target: 'a[href="/orders"]', optional: true },
                { id: 'invoicing', title: 'Facturacion', body: 'b', target: 'a[href="/invoicing/invoices"]', optional: true },
            ],
        });

        const resuelto = resolveVisibleSteps(definicion);

        expect(resuelto.steps.map((s) => s.id)).toEqual(['welcome', 'orders']);
    });

    it('conserva los pasos obligatorios aunque el elemento no exista todavia', () => {
        const definicion = tour({
            steps: [
                { id: 'welcome', title: 'Hola', body: 'b' },
                { id: 'create', title: 'Crear', body: 'b', target: '[data-tour="x"]' },
            ],
        });

        expect(resolveVisibleSteps(definicion).steps).toHaveLength(2);
    });

    it('conserva los pasos que navegan antes de anclarse', () => {
        const definicion = tour({
            steps: [
                { id: 'welcome', title: 'Hola', body: 'b' },
                { id: 'otro', title: 'Otro', body: 'b', target: '#x', route: '/otro', optional: true },
            ],
        });

        expect(resolveVisibleSteps(definicion).steps).toHaveLength(2);
    });

    it('devuelve la misma definicion cuando todos los pasos aplican', () => {
        const definicion = tour({ steps: [{ id: 'welcome', title: 'Hola', body: 'b' }] });

        expect(resolveVisibleSteps(definicion)).toBe(definicion);
    });
});
