import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useIntegrationsSimple } from './useIntegrationsSimple';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getIntegrationsSimpleAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import { getIntegrationsSimpleAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const defaultResponse = {
    success: true,
    message: 'OK',
    data: [
        { id: 1, name: 'Shopify A', type: 'shopify', category: 'ecommerce', category_name: 'E-commerce', business_id: 1, is_active: true },
        { id: 2, name: 'WhatsApp B', type: 'whatsapp', category: 'messaging', category_name: 'Mensajería', business_id: 1, is_active: true },
    ],
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useIntegrationsSimple', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true e integraciones vacías', () => {
        vi.mocked(getIntegrationsSimpleAction).mockReturnValue(new Promise(() => {}));
        const { result } = renderHook(() => useIntegrationsSimple());
        expect(result.current.loading).toBe(true);
        expect(result.current.integrations).toEqual([]);
    });

    it('debería cargar integraciones simples exitosamente', async () => {
        vi.mocked(getIntegrationsSimpleAction).mockResolvedValue(defaultResponse);
        const { result } = renderHook(() => useIntegrationsSimple());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.integrations).toHaveLength(2);
        expect(result.current.error).toBeNull();
    });

    it('debería pasar businessId al action', async () => {
        vi.mocked(getIntegrationsSimpleAction).mockResolvedValue(defaultResponse);
        renderHook(() => useIntegrationsSimple({ businessId: 5 }));
        await waitFor(() => {
            expect(getIntegrationsSimpleAction).toHaveBeenCalledWith(5);
        });
    });

    it('debería capturar error de response cuando success es false', async () => {
        vi.mocked(getIntegrationsSimpleAction).mockResolvedValue({ success: false, message: 'No autorizado', data: [] });
        const { result } = renderHook(() => useIntegrationsSimple());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.error).toBe('No autorizado');
    });

    it('debería capturar error cuando lanza excepción', async () => {
        vi.mocked(getIntegrationsSimpleAction).mockRejectedValue(new Error('Error de red'));
        const { result } = renderHook(() => useIntegrationsSimple());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.error).toBe('Error de red');
    });

    it('debería exponer refresh para recargar', async () => {
        vi.mocked(getIntegrationsSimpleAction).mockResolvedValue(defaultResponse);
        const { result } = renderHook(() => useIntegrationsSimple());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getIntegrationsSimpleAction).mock.calls.length;
        await act(async () => { await result.current.refresh(); });
        expect(vi.mocked(getIntegrationsSimpleAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });
});
