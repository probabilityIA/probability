import { describe, it, expect } from 'vitest';
import { NextRequest } from 'next/server';
import { debeIrALogin, esRutaPublica } from './proxy';

function pedir(url: string, opciones: { cookie?: boolean; header?: string } = {}) {
    const req = new NextRequest(new URL(url, 'https://www.probabilityia.com.co'));
    if (opciones.cookie) req.cookies.set('session_token', 'abc123');
    if (opciones.header) req.headers.set('x-shopify-shop-domain', opciones.header);
    return req;
}

describe('rutas publicas', () => {
    it('deja pasar login, recuperacion, verificaciones y la tienda', () => {
        ['/', '/login', '/forgot-password', '/reset-password', '/verify-code',
         '/verify-demo', '/verify-email', '/storefront/registro',
         '/tienda', '/tienda/demo/producto/5'].forEach((r) => {
            expect(esRutaPublica(r)).toBe(true);
        });
    });

    it('no considera publicas las rutas privadas', () => {
        ['/orders', '/invoicing/invoices', '/home', '/wallet',
         '/loginfalso', '/tiendafalsa'].forEach((r) => {
            expect(esRutaPublica(r)).toBe(false);
        });
    });
});

describe('corte de sesion', () => {
    it('sin cookie manda al login una ruta privada', () => {
        expect(debeIrALogin(pedir('/orders?order_id=abc'))).toBe(true);
    });

    it('con cookie deja pasar', () => {
        expect(debeIrALogin(pedir('/orders?order_id=abc', { cookie: true }))).toBe(false);
    });

    it('no toca las rutas publicas aunque no haya sesion', () => {
        expect(debeIrALogin(pedir('/login'))).toBe(false);
        expect(debeIrALogin(pedir('/tienda/demo'))).toBe(false);
    });

    it('deja pasar el flujo embebido de Shopify, que aun no tiene cookie', () => {
        expect(debeIrALogin(pedir('/home?shop=demo.myshopify.com'))).toBe(false);
        expect(debeIrALogin(pedir('/home?embedded=1'))).toBe(false);
        expect(debeIrALogin(pedir('/home', { header: 'demo.myshopify.com' }))).toBe(false);
    });
});
