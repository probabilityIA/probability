import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useCategories } from './useCategories';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getIntegrationCategoriesAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import { getIntegrationCategoriesAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const defaultResponse = {
    success: true,
    message: 'OK',
    data: [
        { id: 1, code: 'ecommerce', name: 'E-commerce', display_order: 1, is_active: true, is_visible: true, created_at: '', updated_at: '' },
        { id: 2, code: 'invoicing', name: 'Facturación', display_order: 2, is_active: true, is_visible: true, created_at: '', updated_at: '' },
    ],
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useCategories', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y categorías vacías', () => {
        vi.mocked(getIntegrationCategoriesAction).mockReturnValue(new Promise(() => {}));
        const { result } = renderHook(() => useCategories());
        expect(result.current.loading).toBe(true);
        expect(result.current.categories).toEqual([]);
    });

    it('debería cargar categorías exitosamente', async () => {
        vi.mocked(getIntegrationCategoriesAction).mockResolvedValue(defaultResponse);
        const { result } = renderHook(() => useCategories());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.categories).toHaveLength(2);
        expect(result.current.categories[0].code).toBe('ecommerce');
        expect(result.current.error).toBeNull();
    });

    it('debería capturar error de response cuando success es false', async () => {
        vi.mocked(getIntegrationCategoriesAction).mockResolvedValue({ success: false, message: 'No autorizado', data: null as any });
        const { result } = renderHook(() => useCategories());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.error).toBe('No autorizado');
    });

    it('debería capturar error cuando lanza excepción', async () => {
        vi.mocked(getIntegrationCategoriesAction).mockRejectedValue(new Error('Error de red'));
        const { result } = renderHook(() => useCategories());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.error).toBe('Error de red');
    });

    it('debería exponer refresh para recargar categorías', async () => {
        vi.mocked(getIntegrationCategoriesAction).mockResolvedValue(defaultResponse);
        const { result } = renderHook(() => useCategories());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getIntegrationCategoriesAction).mock.calls.length;
        await act(async () => { await result.current.refresh(); });
        expect(vi.mocked(getIntegrationCategoriesAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });
});
