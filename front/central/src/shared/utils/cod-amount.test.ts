import { describe, it, expect } from 'vitest';
import { codCustomerCharge, codBusinessNet } from './cod-amount';

describe('cod-amount', () => {
    it('orden del checkout de WooCommerce: el cod_total ya trae la comision (orden 15789)', () => {
        const order = { cod_total: 73367, cod_carrier_fee: 5365, cod_includes_shipping: true };
        expect(codCustomerCharge(order)).toBe(73367);
        expect(codBusinessNet(order)).toBe(68002);
    });

    it('orden manual: la comision se suma al cod_total (VIG-0069)', () => {
        const order = { cod_total: 135929, cod_carrier_fee: 9451, cod_includes_shipping: false };
        expect(codCustomerCharge(order)).toBe(145380);
        expect(codBusinessNet(order)).toBe(135929);
    });

    it('sin comision conocida todavia muestra el cod_total', () => {
        expect(codCustomerCharge({ cod_total: 73367, cod_carrier_fee: 0, cod_includes_shipping: false })).toBe(73367);
    });

    it('sin cod_total no inventa un monto', () => {
        expect(codCustomerCharge({ cod_total: 0, cod_carrier_fee: 5365 })).toBe(0);
        expect(codBusinessNet({ cod_carrier_fee: 5365 })).toBe(0);
    });
});
