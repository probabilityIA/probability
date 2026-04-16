import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useIntegrations } from './useIntegrations';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getIntegrationsAction: vi.fn(),
    deleteIntegrationAction: vi.fn(),
    activateIntegrationAction: vi.fn(),
    deactivateIntegrationAction: vi.fn(),
    setAsDefaultAction: vi.fn(),
    testConnectionAction: vi.fn(),
    syncOrdersAction: vi.fn(),
}));

vi.mock('@/shared/utils/token-storage', () => ({
    TokenStorage: { getSessionToken: vi.fn(() => 'test-token') },
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import {
    getIntegrationsAction,
    deleteIntegrationAction,
    activateIntegrationAction,
    deactivateIntegrationAction,
    setAsDefaultAction,
    testConnectionAction,
    syncOrdersAction,
} from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeIntegration = (id: number, name: string) => ({
    id, name, code: name.toLowerCase(), integration_type_id: 1, type: 'shopify',
    category: 'ecommerce', business_id: 1, is_active: true, is_default: false,
    is_testing: false, config: {}, created_by_id: 1, updated_by_id: null,
    created_at: '', updated_at: '',
});

const defaultResponse = {
    success: true, message: 'OK',
    data: [makeIntegration(1, 'Shopify A'), makeIntegration(2, 'WhatsApp B')],
    total: 2, page: 1, page_size: 10, total_pages: 1,
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useIntegrations', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true e integraciones vacías', () => {
        vi.mocked(getIntegrationsAction).mockReturnValue(new Promise(() => {}));
        const { result } = renderHook(() => useIntegrations());
        expect(result.current.loading).toBe(true);
        expect(result.current.integrations).toEqual([]);
    });

    it('debería cargar integraciones exitosamente', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.integrations).toHaveLength(2);
        expect(result.current.totalPages).toBe(1);
    });

    it('debería capturar error cuando getIntegrationsAction falla', async () => {
        vi.mocked(getIntegrationsAction).mockRejectedValue(new Error('Error de red'));
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.error).toBe('Error de red');
    });

    it('debería pasar initialCategory como filtro', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        renderHook(() => useIntegrations('invoicing'));
        await waitFor(() => {
            const lastCall = vi.mocked(getIntegrationsAction).mock.calls.at(-1)![0];
            expect(lastCall?.category).toBe('invoicing');
        });
    });

    it('debería cambiar de página y refrescar', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => { result.current.setPage(2); });
        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getIntegrationsAction).mock.calls.at(-1)![0];
        expect(lastCall?.page).toBe(2);
    });

    it('debería eliminar integración y refrescar', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteIntegrationAction).mockResolvedValue({ success: true, message: 'OK' });
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => { deleteResult = await result.current.deleteIntegration(1); });
        expect(deleteResult).toBe(true);
        expect(deleteIntegrationAction).toHaveBeenCalledWith(1, 'test-token');
    });

    it('debería desactivar integración activa con toggleActive', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        vi.mocked(deactivateIntegrationAction).mockResolvedValue({ success: true, message: 'OK' });
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));

        await act(async () => { await result.current.toggleActive(1, true); });
        expect(deactivateIntegrationAction).toHaveBeenCalledWith(1, 'test-token');
        expect(activateIntegrationAction).not.toHaveBeenCalled();
    });

    it('debería activar integración inactiva con toggleActive', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        vi.mocked(activateIntegrationAction).mockResolvedValue({ success: true, message: 'OK' });
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));

        await act(async () => { await result.current.toggleActive(1, false); });
        expect(activateIntegrationAction).toHaveBeenCalledWith(1, 'test-token');
    });

    it('debería establecer como default', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        vi.mocked(setAsDefaultAction).mockResolvedValue({ success: true, message: 'OK', data: makeIntegration(1, 'x') });
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let setResult: boolean | undefined;
        await act(async () => { setResult = await result.current.setAsDefault(1); });
        expect(setResult).toBe(true);
    });

    it('debería probar conexión', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        vi.mocked(testConnectionAction).mockResolvedValue({ success: true, message: 'Connected' });
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let testResult: any;
        await act(async () => { testResult = await result.current.testConnection(1); });
        expect(testResult.success).toBe(true);
    });

    it('debería sincronizar órdenes', async () => {
        vi.mocked(getIntegrationsAction).mockResolvedValue(defaultResponse);
        vi.mocked(syncOrdersAction).mockResolvedValue({ success: true, message: 'Synced' });
        const { result } = renderHook(() => useIntegrations());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let syncResult: any;
        await act(async () => { syncResult = await result.current.syncOrders(1, { status: 'open' }); });
        expect(syncResult.success).toBe(true);
    });
});
