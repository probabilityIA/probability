import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useResourceConfig } from './useResourceConfig';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getBusinessConfiguredResourcesAction: vi.fn(),
    activateResourceAction: vi.fn(),
    deactivateResourceAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any) =>
        err instanceof Error ? err.message : 'Error'
    ),
}));

import {
    getBusinessConfiguredResourcesAction,
    activateResourceAction,
    deactivateResourceAction,
} from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeConfigResponse = () => ({
    success: true,
    message: 'OK',
    data: {
        business_id: 1,
        resources: [
            { resource_id: 10, resource_name: 'Delivery', is_active: true },
            { resource_id: 20, resource_name: 'Pickup', is_active: false },
        ],
        total: 2,
        active: 1,
        inactive: 1,
    },
});

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useResourceConfig', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y config null', () => {
        vi.mocked(getBusinessConfiguredResourcesAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useResourceConfig(1));

        expect(result.current.loading).toBe(true);
        expect(result.current.config).toBeNull();
        expect(result.current.error).toBeNull();
    });

    it('debería cargar configuración de recursos por businessId', async () => {
        vi.mocked(getBusinessConfiguredResourcesAction).mockResolvedValue(makeConfigResponse());

        const { result } = renderHook(() => useResourceConfig(1));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.config).toBeDefined();
        expect(result.current.config!.resources).toHaveLength(2);
        expect(getBusinessConfiguredResourcesAction).toHaveBeenCalledWith(1);
    });

    it('debería capturar error cuando getBusinessConfiguredResourcesAction falla', async () => {
        vi.mocked(getBusinessConfiguredResourcesAction).mockRejectedValue(new Error('No encontrado'));

        const { result } = renderHook(() => useResourceConfig(1));

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('No encontrado');
    });

    it('debería no cargar cuando businessId es 0 (falsy)', () => {
        const { result } = renderHook(() => useResourceConfig(0));

        expect(getBusinessConfiguredResourcesAction).not.toHaveBeenCalled();
    });

    it('debería desactivar recurso activo al llamar toggleResource', async () => {
        vi.mocked(getBusinessConfiguredResourcesAction).mockResolvedValue(makeConfigResponse());
        vi.mocked(deactivateResourceAction).mockResolvedValue({ success: true, message: 'OK' });

        const { result } = renderHook(() => useResourceConfig(1));
        await waitFor(() => expect(result.current.loading).toBe(false));

        await act(async () => {
            await result.current.toggleResource({ resource_id: 10, resource_name: 'Delivery', is_active: true });
        });

        expect(deactivateResourceAction).toHaveBeenCalledWith(10, 1);
        expect(activateResourceAction).not.toHaveBeenCalled();
    });

    it('debería activar recurso inactivo al llamar toggleResource', async () => {
        vi.mocked(getBusinessConfiguredResourcesAction).mockResolvedValue(makeConfigResponse());
        vi.mocked(activateResourceAction).mockResolvedValue({ success: true, message: 'OK' });

        const { result } = renderHook(() => useResourceConfig(1));
        await waitFor(() => expect(result.current.loading).toBe(false));

        await act(async () => {
            await result.current.toggleResource({ resource_id: 20, resource_name: 'Pickup', is_active: false });
        });

        expect(activateResourceAction).toHaveBeenCalledWith(20, 1);
        expect(deactivateResourceAction).not.toHaveBeenCalled();
    });

    it('debería capturar error y limpiar actionLoading cuando toggleResource falla', async () => {
        vi.mocked(getBusinessConfiguredResourcesAction).mockResolvedValue(makeConfigResponse());
        vi.mocked(deactivateResourceAction).mockRejectedValue(new Error('Error de servidor'));

        const { result } = renderHook(() => useResourceConfig(1));
        await waitFor(() => expect(result.current.loading).toBe(false));

        await act(async () => {
            await result.current.toggleResource({ resource_id: 10, resource_name: 'Delivery', is_active: true });
        });

        expect(result.current.error).toBe('Error de servidor');
        expect(result.current.actionLoading).toBeNull();
    });
});
