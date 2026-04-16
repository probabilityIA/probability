import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useClientSearch } from './useClientSearch';

// -----------------------------------------------------------------
// Mock del módulo de Server Actions (customers)
// -----------------------------------------------------------------

vi.mock('../../../customers/infra/actions', () => ({
    getCustomersAction: vi.fn(),
}));

import { getCustomersAction } from '../../../customers/infra/actions';

// -----------------------------------------------------------------
// Helpers: datos de prueba
// -----------------------------------------------------------------

const makeCustomer = (id: number, name: string) => ({
    id,
    business_id: 1,
    name,
    email: `${name.toLowerCase().replace(/\s/g, '.')}@test.com`,
    phone: '+573001234567',
    dni: `100000000${id}`,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
});

const defaultCustomersResponse = {
    success: true,
    message: 'OK',
    data: [
        makeCustomer(1, 'Juan Pérez'),
        makeCustomer(2, 'Juan García'),
    ],
    total: 2,
    page: 1,
    page_size: 5,
    total_pages: 1,
};

const emptyCustomersResponse = {
    success: true,
    message: 'OK',
    data: [],
    total: 0,
    page: 1,
    page_size: 5,
    total_pages: 0,
};

// Use a very short debounce for tests with real timers
const SHORT_DEBOUNCE = 10;

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useClientSearch', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    const defaultOptions = { businessId: 1, debounceMs: SHORT_DEBOUNCE, minChars: 2 };

    // ---------------------------------------------------------------
    // Estado inicial
    // ---------------------------------------------------------------
    it('debería iniciar con estado vacío', () => {
        const { result } = renderHook(() => useClientSearch(defaultOptions));

        expect(result.current.results).toEqual([]);
        expect(result.current.loading).toBe(false);
        expect(result.current.searched).toBe(false);
    });

    // ---------------------------------------------------------------
    // Búsqueda exitosa con debounce
    // ---------------------------------------------------------------
    it('debería buscar clientes después del debounce', async () => {
        vi.mocked(getCustomersAction).mockResolvedValue(defaultCustomersResponse);

        const { result } = renderHook(() => useClientSearch(defaultOptions));

        act(() => {
            result.current.search('Juan');
        });

        // Antes del debounce: loading true
        expect(result.current.loading).toBe(true);

        await waitFor(() => {
            expect(result.current.searched).toBe(true);
        });

        expect(result.current.results).toHaveLength(2);
        expect(result.current.results[0].name).toBe('Juan Pérez');
        expect(result.current.loading).toBe(false);
        expect(getCustomersAction).toHaveBeenCalledWith({
            search: 'Juan',
            business_id: 1,
            page: 1,
            page_size: 5,
        });
    });

    // ---------------------------------------------------------------
    // Búsqueda con menos caracteres del mínimo
    // ---------------------------------------------------------------
    it('debería no buscar cuando el término tiene menos caracteres que minChars', () => {
        const { result } = renderHook(() => useClientSearch(defaultOptions));

        act(() => {
            result.current.search('J');
        });

        expect(result.current.results).toEqual([]);
        expect(result.current.loading).toBe(false);
        expect(result.current.searched).toBe(false);
        expect(getCustomersAction).not.toHaveBeenCalled();
    });

    // ---------------------------------------------------------------
    // Búsqueda con término vacío
    // ---------------------------------------------------------------
    it('debería limpiar resultados cuando el término está vacío', () => {
        const { result } = renderHook(() => useClientSearch(defaultOptions));

        act(() => {
            result.current.search('');
        });

        expect(result.current.results).toEqual([]);
        expect(result.current.loading).toBe(false);
        expect(result.current.searched).toBe(false);
    });

    // ---------------------------------------------------------------
    // Búsqueda sin businessId
    // ---------------------------------------------------------------
    it('debería no buscar cuando businessId es 0', () => {
        const { result } = renderHook(() =>
            useClientSearch({ businessId: 0, debounceMs: SHORT_DEBOUNCE, minChars: 2 })
        );

        act(() => {
            result.current.search('Juan');
        });

        expect(getCustomersAction).not.toHaveBeenCalled();
        expect(result.current.loading).toBe(false);
    });

    // ---------------------------------------------------------------
    // Búsqueda con resultado vacío
    // ---------------------------------------------------------------
    it('debería manejar una búsqueda sin resultados', async () => {
        vi.mocked(getCustomersAction).mockResolvedValue(emptyCustomersResponse);

        const { result } = renderHook(() => useClientSearch(defaultOptions));

        act(() => {
            result.current.search('ZZZZZ');
        });

        await waitFor(() => {
            expect(result.current.searched).toBe(true);
        });

        expect(result.current.results).toEqual([]);
        expect(result.current.loading).toBe(false);
    });

    // ---------------------------------------------------------------
    // Manejo de error en búsqueda
    // ---------------------------------------------------------------
    it('debería manejar error en la búsqueda sin explotar', async () => {
        vi.mocked(getCustomersAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useClientSearch(defaultOptions));

        act(() => {
            result.current.search('Juan');
        });

        await waitFor(() => {
            expect(result.current.searched).toBe(true);
        });

        expect(result.current.results).toEqual([]);
        expect(result.current.loading).toBe(false);
    });

    // ---------------------------------------------------------------
    // Función clear
    // ---------------------------------------------------------------
    it('debería limpiar resultados al llamar clear', async () => {
        vi.mocked(getCustomersAction).mockResolvedValue(defaultCustomersResponse);

        const { result } = renderHook(() => useClientSearch(defaultOptions));

        // Primero buscar
        act(() => {
            result.current.search('Juan');
        });

        await waitFor(() => {
            expect(result.current.results).toHaveLength(2);
        });

        // Luego limpiar
        act(() => {
            result.current.clear();
        });

        expect(result.current.results).toEqual([]);
        expect(result.current.loading).toBe(false);
        expect(result.current.searched).toBe(false);
    });

    // ---------------------------------------------------------------
    // pageSize personalizado
    // ---------------------------------------------------------------
    it('debería respetar el pageSize personalizado', async () => {
        vi.mocked(getCustomersAction).mockResolvedValue(defaultCustomersResponse);

        const { result } = renderHook(() =>
            useClientSearch({ businessId: 1, debounceMs: SHORT_DEBOUNCE, minChars: 2, pageSize: 10 })
        );

        act(() => {
            result.current.search('Juan');
        });

        await waitFor(() => {
            expect(getCustomersAction).toHaveBeenCalledWith(
                expect.objectContaining({ page_size: 10 })
            );
        });
    });
});
