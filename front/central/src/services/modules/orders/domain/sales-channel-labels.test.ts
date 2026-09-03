import { describe, it, expect } from 'vitest';
import { esCanalDeVenta } from '@/services/modules/orders/domain/sales-channel-labels';

describe('esCanalDeVenta', () => {
    it('reconoce los canales externos', () => {
        ['woocommerce', 'shopify', 'mercadolibre', 'MercadoLibre', 'tiendanube'].forEach((p) => {
            expect(esCanalDeVenta(p)).toBe(true);
        });
    });

    it('no marca las ordenes propias', () => {
        ['manual', 'Manual', 'plataforma', 'platform', 'api', '', null, undefined].forEach((p) => {
            expect(esCanalDeVenta(p)).toBe(false);
        });
    });
});
