import { describe, it, expect } from 'vitest';
import { findDaneCode, resolveCityState } from './dane-lookup';
import danes from '@/app/(auth)/shipments/generate/resources/municipios_dane_extendido.json';

// -----------------------------------------------------------------
// Simula el payload real que envia WooCommerce: shipping.state llega
// como codigo ISO 3166-2 ("CO-MET"), nunca como nombre completo.
// -----------------------------------------------------------------

const wooCommerceOrder = (overrides: Partial<{ shipping_city: string; shipping_state: string; destination_dane_code: string | null }> = {}) => ({
    shipping_city: 'GRANADA',
    shipping_state: 'CO-MET',
    destination_dane_code: null as string | null,
    ...overrides,
});

describe('findDaneCode - caso real orden #145802 (WooCommerce, ciudad duplicada)', () => {
    it('resuelve GRANADA + CO-MET al municipio de Meta, no Cundinamarca ni Antioquia', () => {
        const code = findDaneCode('GRANADA', 'CO-MET');
        expect(code).toBe('50313000');
    });

    it('reproduce el bug original si se ignora el mapeo ISO (regresion)', () => {
        // Antes del fix, "CO-MET" nunca calzaba con "META" en el match exacto,
        // y el fallback por ciudad devolvia Cundinamarca (25312000) por el
        // orden de iteracion de claves numericas en JS, no Antioquia (05313000)
        // como sugeriria el orden textual del JSON.
        const withoutIsoFix = (city: string, state: string) => {
            const norm = (s: string) => s.normalize('NFD').replace(/[̀-ͯ]/g, '').toUpperCase().trim();
            const targetCity = norm(city);
            const entries = Object.entries(danes);
            const exact = entries.find(([_, d]: [string, any]) => norm(d.ciudad) === targetCity && norm(d.departamento) === norm(state));
            if (exact) return exact[0];
            const cityOnly = entries.find(([_, d]: [string, any]) => norm(d.ciudad) === targetCity);
            return cityOnly ? cityOnly[0] : null;
        };

        expect(withoutIsoFix('GRANADA', 'CO-MET')).toBe('25312000');
        expect(findDaneCode('GRANADA', 'CO-MET')).not.toBe('25312000');
    });

    it.each([
        ['GRANADA', 'CO-ANT', '05313000', 'ANTIOQUIA'],
        ['GRANADA', 'CO-CUN', '25312000', 'CUNDINAMARCA'],
        ['GRANADA', 'CO-MET', '50313000', 'META'],
    ])('desambigua %s + %s -> %s (%s)', (city, isoState, expectedCode) => {
        expect(findDaneCode(city, isoState)).toBe(expectedCode);
    });

    it('sigue funcionando si el estado ya llega como nombre completo (compatibilidad)', () => {
        expect(findDaneCode('GRANADA', 'META')).toBe('50313000');
        expect(findDaneCode('GRANADA', 'Meta')).toBe('50313000');
    });

    it('normaliza otros codigos ISO comunes de ordenes reales (Antioquia, Santander)', () => {
        expect(findDaneCode('MEDELLIN', 'CO-ANT')).toBe(findDaneCode('MEDELLIN', 'ANTIOQUIA'));
        expect(findDaneCode('BUCARAMANGA', 'CO-SAN')).toBe(findDaneCode('BUCARAMANGA', 'SANTANDER'));
    });

    it('cae a busqueda solo por ciudad si el departamento no matchea ningun ISO ni nombre conocido', () => {
        const code = findDaneCode('GRANADA', 'CO-XXX');
        expect(code).not.toBeNull();
    });

    it('retorna null si la ciudad esta vacia', () => {
        expect(findDaneCode('', 'CO-MET')).toBeNull();
    });
});

describe('resolveCityState - validacion en dos niveles (geozona verificada + fallback ISO)', () => {
    it('Nivel 1: prioriza el destination_dane_code verificado aunque el texto sea ambiguo', () => {
        const order = wooCommerceOrder({ destination_dane_code: '50313' });
        const resolved = resolveCityState(order.shipping_city, order.shipping_state, order.destination_dane_code);
        expect(resolved).toEqual({ ciudad: 'GRANADA', departamento: 'META' });
    });

    it('Nivel 1 ignora un destination_dane_code invalido y cae al Nivel 2', () => {
        const order = wooCommerceOrder({ destination_dane_code: '99999' });
        const resolved = resolveCityState(order.shipping_city, order.shipping_state, order.destination_dane_code);
        expect(resolved?.departamento).toBe('META'); // resuelto por findDaneCode con fix ISO
    });

    it('Nivel 2: sin destination_dane_code, resuelve por texto con el mapeo ISO', () => {
        const order = wooCommerceOrder({ destination_dane_code: null });
        const resolved = resolveCityState(order.shipping_city, order.shipping_state, order.destination_dane_code);
        expect(resolved).toEqual({ ciudad: 'GRANADA', departamento: 'META' });
    });

    it('replica el flujo completo de la orden 145802 y confirma que NO cae a Cundinamarca', () => {
        // Antes del fix este payload real generaba la guia con destino Cundinamarca
        const order = {
            shipping_city: 'GRANADA',
            shipping_state: 'CO-MET',
            destination_dane_code: '50313', // resuelto por geocodificacion (lat/lng reales de la orden)
        };
        const resolved = resolveCityState(order.shipping_city, order.shipping_state, order.destination_dane_code);
        expect(resolved?.departamento).toBe('META');
        expect(resolved?.departamento).not.toBe('CUNDINAMARCA');
    });

    it('retorna null cuando no hay match posible', () => {
        expect(resolveCityState('CIUDAD-INEXISTENTE-XYZ', '', null)).toBeNull();
    });
});
