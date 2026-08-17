import { describe, it, expect } from 'vitest';
import { rateGuideCost, rateCarrierFee, rateTotalCost, rateAppliesCOD } from './rate-pricing';

const coordinadora = {
    flete: 7796,
    minimumInsurance: 580,
    extraInsurance: 920,
    cod: true,
    codCarrierFee: 6116,
    codProbabilityMargin: 500,
};

const servientrega = {
    flete: 15390,
    minimumInsurance: 0,
    extraInsurance: 0,
    cod: true,
    codCarrierFee: 8499,
    codProbabilityMargin: 500,
};

const sinCOD = {
    flete: 11424,
    minimumInsurance: 0,
    extraInsurance: 1500,
    cod: false,
};

describe('rate-pricing', () => {
    it('suma flete y seguro minimo siempre', () => {
        expect(rateGuideCost(coordinadora)).toBe(8376);
    });

    it('suma el margen COD solo cuando hay contra entrega', () => {
        expect(rateGuideCost(coordinadora, { cod: true })).toBe(8876);
        expect(rateGuideCost(coordinadora, { cod: false })).toBe(8376);
    });

    it('no aplica margen COD a una tarifa que no soporta contra entrega', () => {
        expect(rateAppliesCOD(sinCOD, { cod: true })).toBe(false);
        expect(rateGuideCost(sinCOD, { cod: true })).toBe(11424);
    });

    it('suma el seguro adicional solo si se asegura', () => {
        expect(rateGuideCost(coordinadora, { insured: true })).toBe(9296);
        expect(rateGuideCost(sinCOD, { insured: true })).toBe(12924);
    });

    it('la comision del carrier solo cuenta con contra entrega', () => {
        expect(rateCarrierFee(coordinadora, { cod: true })).toBe(6116);
        expect(rateCarrierFee(coordinadora, { cod: false })).toBe(0);
    });

    it('el total suma guia mas comision', () => {
        expect(rateTotalCost(coordinadora, { cod: true })).toBe(14992);
        expect(rateTotalCost(servientrega, { cod: true })).toBe(24389);
    });

    it('tolera tarifas incompletas', () => {
        expect(rateGuideCost(null)).toBe(0);
        expect(rateGuideCost({}, { cod: true, insured: true })).toBe(0);
        expect(rateTotalCost({ flete: 1000 } as any, { cod: true })).toBe(1000);
    });
});
