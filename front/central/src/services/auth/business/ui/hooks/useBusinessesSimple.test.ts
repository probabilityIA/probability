import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useBusinessesSimple } from './useBusinessesSimple';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getBusinessesSimpleAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import { getBusinessesSimpleAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const defaultResponse = {
    success: true,
    message: 'OK',
    data: [
        { id: 1, name: 'Negocio A' },
        { id: 2, name: 'Negocio B' },
    ],
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useBusinessesSimple', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y businesses vacíos', () => {
        vi.mocked(getBusinessesSimpleAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useBusinessesSimple());

        expect(result.current.loading).toBe(true);
        expect(result.current.businesses).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    it('debería cargar negocios simples exitosamente cuando success es true', async () => {
        vi.mocked(getBusinessesSimpleAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useBusinessesSimple());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.businesses).toHaveLength(2);
        expect(result.current.businesses[0].name).toBe('Negocio A');
        expect(result.current.error).toBeNull();
    });

    it('debería capturar error del response cuando success es false', async () => {
        vi.mocked(getBusinessesSimpleAction).mockResolvedValue({
            success: false,
            message: 'No autorizado',
            data: [],
        });

        const { result } = renderHook(() => useBusinessesSimple());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('No autorizado');
        expect(result.current.businesses).toEqual([]);
    });

    it('debería capturar error cuando getBusinessesSimpleAction lanza excepción', async () => {
        vi.mocked(getBusinessesSimpleAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useBusinessesSimple());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
    });

    it('debería exponer refresh para recargar negocios', async () => {
        vi.mocked(getBusinessesSimpleAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useBusinessesSimple());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getBusinessesSimpleAction).mock.calls.length;

        await act(async () => {
            await result.current.refresh();
        });

        expect(vi.mocked(getBusinessesSimpleAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });
});
