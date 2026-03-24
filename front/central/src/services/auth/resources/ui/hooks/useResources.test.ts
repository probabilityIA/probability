import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useResources } from './useResources';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getResourcesAction: vi.fn(),
    deleteResourceAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import { getResourcesAction, deleteResourceAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeResource = (id: number, name: string) => ({
    id,
    name,
    description: `Desc ${name}`,
    business_type_id: 1,
    business_type_name: 'Restaurante',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
});

const defaultResponse = {
    success: true,
    message: 'OK',
    data: {
        resources: [makeResource(1, 'orders'), makeResource(2, 'products')],
        total: 2,
        page: 1,
        page_size: 20,
        total_pages: 1,
    },
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useResources', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y recursos vacíos', () => {
        vi.mocked(getResourcesAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useResources());

        expect(result.current.loading).toBe(true);
        expect(result.current.resources).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    it('debería cargar recursos exitosamente', async () => {
        vi.mocked(getResourcesAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useResources());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.resources).toHaveLength(2);
        expect(result.current.total).toBe(2);
        expect(result.current.totalPages).toBe(1);
        expect(result.current.error).toBeNull();
    });

    it('debería capturar error de response cuando success es false', async () => {
        vi.mocked(getResourcesAction).mockResolvedValue({
            success: false,
            message: 'No autorizado',
            data: null as any,
        });

        const { result } = renderHook(() => useResources());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('No autorizado');
    });

    it('debería capturar error cuando getResourcesAction lanza excepción', async () => {
        vi.mocked(getResourcesAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useResources());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
    });

    it('debería cambiar de página y refrescar', async () => {
        vi.mocked(getResourcesAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useResources());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getResourcesAction).mock.calls.length;

        act(() => {
            result.current.setPage(2);
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getResourcesAction).mock.calls.at(-1)![0];
        expect(lastCall?.page).toBe(2);
        expect(vi.mocked(getResourcesAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });

    it('debería aplicar filtros y refrescar', async () => {
        vi.mocked(getResourcesAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useResources());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setFilters({ name: 'orders' });
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getResourcesAction).mock.calls.at(-1)![0];
        expect(lastCall?.name).toBe('orders');
    });

    it('debería eliminar recurso y refrescar lista cuando success es true', async () => {
        vi.mocked(getResourcesAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteResourceAction).mockResolvedValue({ success: true, message: 'OK' });

        const { result } = renderHook(() => useResources());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteResource(1);
        });

        expect(deleteResult).toBe(true);
        expect(deleteResourceAction).toHaveBeenCalledWith(1);
    });

    it('debería capturar error cuando deleteResource retorna success false', async () => {
        vi.mocked(getResourcesAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteResourceAction).mockResolvedValue({ success: false, message: 'No se puede eliminar' });

        const { result } = renderHook(() => useResources());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteResource(1);
        });

        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('No se puede eliminar');
    });

    it('debería capturar error cuando deleteResource lanza excepción', async () => {
        vi.mocked(getResourcesAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteResourceAction).mockRejectedValue(new Error('Error de servidor'));

        const { result } = renderHook(() => useResources());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteResource(1);
        });

        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('Error de servidor');
    });
});
