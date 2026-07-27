import { describe, expect, it } from 'vitest';
import { buildApplyMessage } from './sync-activity-context';

describe('buildApplyMessage', () => {
    it('resume una aplicacion exitosa', () => {
        expect(buildApplyMessage(50, 50, 0, 0)).toBe('50 creados de 50');
    });

    it('marca los que fallaron', () => {
        expect(buildApplyMessage(50, 1, 0, 49)).toBe('1 creados, 49 con error de 50');
    });

    it('agrega el motivo del primer error', () => {
        const message = buildApplyMessage(2, 0, 0, 2, [
            { sku: 'PT01004', error: 'item.attributes: BRAND es requerido' },
            { sku: 'PT01015', error: 'otro' },
        ]);
        expect(message).toBe('2 con error de 2 — PT01004: item.attributes: BRAND es requerido (y 1 mas)');
    });

    it('avisa cuando no habia nada por aplicar', () => {
        expect(buildApplyMessage(0, 0, 0, 0)).toBe('No habia productos por aplicar');
    });
});
