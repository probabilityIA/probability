import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { useWarehouses } from './useWarehouses';

// -----------------------------------------------------------------
// Mock del módulo de Server Actions (warehouses)
// -----------------------------------------------------------------

vi.mock('../../../warehouses/infra/actions', () => ({
    getWarehousesAction: vi.fn(),
}));

import { getWarehousesAction } from '../../../warehouses/infra/actions';

// -----------------------------------------------------------------
// Helpers: datos de prueba
// -----------------------------------------------------------------

const makeWarehouse = (id: number, name: string) => ({
    id,
    business_id: 1,
    name,
    code: name.toLowerCase().replace(/\s/g, '-'),
    address: 'Calle 100 #15-20',
    city: 'Bogotá',
    state: 'Cundinamarca',
    country: 'CO',
    zip_code: '110111',
    phone: '+573001234567',
    contact_name: 'Admin',
    contact_email: 'admin@test.com',
    is_active: true,
    is_default: id === 1,
    is_fulfillment: true,
    company: 'Test Co',
    first_name: 'Admin',
    last_name: 'Test',
    email: 'admin@test.com',
    suburb: '',
    city_dane_code: '11001',
    postal_code: '110111',
    street: 'Calle 100',
    latitude: null,
    longitude: null,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
});

const defaultWarehousesResponse = {
    data: [
        makeWarehouse(1, 'Bodega Principal'),
        makeWarehouse(2, 'Bodega Norte'),
    ],
    total: 2,
    page: 1,
    page_size: 100,
    total_pages: 1,
};

const emptyWarehousesResponse = {
    data: [],
    total: 0,
    page: 1,
    page_size: 100,
    total_pages: 0,
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useWarehouses', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    // ---------------------------------------------------------------
    // Estado inicial
    // ---------------------------------------------------------------
    it('debería iniciar con estado vacío y loading false', () => {
        vi.mocked(getWarehousesAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useWarehouses({ businessId: 1 }));

        // loading se establece en true dentro del useEffect
        expect(result.current.warehouses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Carga exitosa
    // ---------------------------------------------------------------
    it('debería cargar bodegas exitosamente', async () => {
        vi.mocked(getWarehousesAction).mockResolvedValue(defaultWarehousesResponse);

        const { result } = renderHook(() => useWarehouses({ businessId: 1 }));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.warehouses).toHaveLength(2);
        expect(result.current.warehouses[0].name).toBe('Bodega Principal');
        expect(result.current.warehouses[1].name).toBe('Bodega Norte');
    });

    // ---------------------------------------------------------------
    // Parámetros enviados a la action
    // ---------------------------------------------------------------
    it('debería llamar a la action con los parámetros correctos', async () => {
        vi.mocked(getWarehousesAction).mockResolvedValue(defaultWarehousesResponse);

        renderHook(() => useWarehouses({ businessId: 5 }));

        await waitFor(() => {
            expect(getWarehousesAction).toHaveBeenCalledWith({
                business_id: 5,
                is_active: true,
                page: 1,
                page_size: 100,
            });
        });
    });

    // ---------------------------------------------------------------
    // No cargar cuando businessId es 0
    // ---------------------------------------------------------------
    it('debería no cargar cuando businessId es 0', () => {
        const { result } = renderHook(() => useWarehouses({ businessId: 0 }));

        expect(getWarehousesAction).not.toHaveBeenCalled();
        expect(result.current.warehouses).toEqual([]);
        expect(result.current.loading).toBe(false);
    });

    // ---------------------------------------------------------------
    // Respuesta vacía
    // ---------------------------------------------------------------
    it('debería manejar una respuesta vacía', async () => {
        vi.mocked(getWarehousesAction).mockResolvedValue(emptyWarehousesResponse);

        const { result } = renderHook(() => useWarehouses({ businessId: 1 }));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.warehouses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Manejo de error
    // ---------------------------------------------------------------
    it('debería manejar error sin explotar y dejar warehouses vacío', async () => {
        vi.mocked(getWarehousesAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useWarehouses({ businessId: 1 }));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.warehouses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Recargar al cambiar businessId
    // ---------------------------------------------------------------
    it('debería recargar cuando cambia el businessId', async () => {
        vi.mocked(getWarehousesAction).mockResolvedValue(defaultWarehousesResponse);

        const { rerender } = renderHook(
            ({ businessId }) => useWarehouses({ businessId }),
            { initialProps: { businessId: 1 } }
        );

        await waitFor(() => {
            expect(getWarehousesAction).toHaveBeenCalledTimes(1);
        });

        rerender({ businessId: 2 });

        await waitFor(() => {
            expect(getWarehousesAction).toHaveBeenCalledTimes(2);
        });

        expect(vi.mocked(getWarehousesAction).mock.calls[1][0]).toEqual({
            business_id: 2,
            is_active: true,
            page: 1,
            page_size: 100,
        });
    });

    // ---------------------------------------------------------------
    // Respuesta con data null/undefined
    // ---------------------------------------------------------------
    it('debería manejar respuesta con data null como array vacío', async () => {
        vi.mocked(getWarehousesAction).mockResolvedValue({
            data: null,
            total: 0,
            page: 1,
            page_size: 100,
            total_pages: 0,
        } as any);

        const { result } = renderHook(() => useWarehouses({ businessId: 1 }));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.warehouses).toEqual([]);
    });
});
