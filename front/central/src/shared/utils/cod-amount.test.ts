import { describe, it, expect } from 'vitest';
import { codCustomerCharge, codBusinessNet, codEffectiveCarrierFee } from './cod-amount';

describe('cod-amount', () => {
    it('orden del plugin de WooCommerce con guia: cod_total neto + comision del checkout (orden 15789)', () => {
        const order = { cod_total: 68002, cod_carrier_fee: 5365, cod_includes_shipping: true, cod_checkout_carrier_fee: 5365 };
        expect(codCustomerCharge(order)).toBe(73367);
        expect(codBusinessNet(order)).toBe(68002);
    });

    it('orden del plugin antes de la guia: usa la comision del checkout', () => {
        const order = { cod_total: 68002, cod_carrier_fee: 0, cod_includes_shipping: true, cod_checkout_carrier_fee: 5365 };
        expect(codCustomerCharge(order)).toBe(73367);
        expect(codEffectiveCarrierFee(order)).toBe(5365);
        expect(codBusinessNet(order)).toBe(68002);
    });

    it('orden del plugin con comision real distinta: se cobra lo de la tienda y el negocio asume (orden 15788)', () => {
        const order = { cod_total: 65990, cod_carrier_fee: 6116, cod_includes_shipping: true, cod_checkout_carrier_fee: 5238 };
        expect(codCustomerCharge(order)).toBe(71228);
        expect(codBusinessNet(order)).toBe(65112);
    });

    it('orden manual: la comision se suma al cod_total (VIG-0069)', () => {
        const order = { cod_total: 135929, cod_carrier_fee: 9451, cod_includes_shipping: false };
        expect(codCustomerCharge(order)).toBe(145380);
        expect(codBusinessNet(order)).toBe(135929);
    });

    it('canal que cobra el total en su checkout sin plugin: no se suma la comision', () => {
        const order = { cod_total: 191322, cod_carrier_fee: 6116, cod_includes_shipping: true };
        expect(codCustomerCharge(order)).toBe(191322);
        expect(codBusinessNet(order)).toBe(185206);
    });

    it('sin cod_total no inventa un monto', () => {
        expect(codCustomerCharge({ cod_total: 0, cod_carrier_fee: 5365, cod_checkout_carrier_fee: 5365 })).toBe(0);
        expect(codBusinessNet({ cod_carrier_fee: 5365 })).toBe(0);
    });
});
