import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useOrderStatuses } from './useOrderStatuses';

// -----------------------------------------------------------------
// Mock del módulo de Server Actions
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getOrderStatusesSimpleAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback: string) =>
        err instanceof Error ? err.message : fallback
    ),
}));

import { getOrderStatusesSimpleAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers: datos de prueba
// -----------------------------------------------------------------

const makeOrderStatus = (id: number, name: string, code: string) => ({
    id,
    name,
    code,
    color: '#3B82F6',
    is_active: true,
});

const defaultSuccessResponse = {
    success: true,
    message: 'OK',
    data: [
        makeOrderStatus(1, 'Pendiente', 'pending'),
        makeOrderStatus(2, 'Procesando', 'processing'),
        makeOrderStatus(3, 'Completado', 'completed'),
    ],
};

const emptyResponse = {
    success: true,
    message: 'OK',
    data: [],
};

const errorResponse = {
    success: false,
    message: 'Error al cargar estados de orden',
    data: [],
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useOrderStatuses', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    // ---------------------------------------------------------------
    // Estado inicial
    // ---------------------------------------------------------------
    it('debería iniciar con loading en true y sin estados', () => {
        vi.mocked(getOrderStatusesSimpleAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useOrderStatuses());

        expect(result.current.loading).toBe(true);
        expect(result.current.orderStatuses).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    // ---------------------------------------------------------------
    // Carga exitosa
    // ---------------------------------------------------------------
    it('debería cargar estados de orden exitosamente', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockResolvedValue(defaultSuccessResponse);

        const { result } = renderHook(() => useOrderStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.orderStatuses).toHaveLength(3);
        expect(result.current.orderStatuses[0].name).toBe('Pendiente');
        expect(result.current.orderStatuses[1].code).toBe('processing');
        expect(result.current.error).toBeNull();
    });

    // ---------------------------------------------------------------
    // Respuesta vacía exitosa
    // ---------------------------------------------------------------
    it('debería manejar una respuesta vacía exitosa', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockResolvedValue(emptyResponse);

        const { result } = renderHook(() => useOrderStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.orderStatuses).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    // ---------------------------------------------------------------
    // Respuesta con success: false
    // ---------------------------------------------------------------
    it('debería establecer error cuando la respuesta tiene success: false', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockResolvedValue(errorResponse);

        const { result } = renderHook(() => useOrderStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error al cargar estados de orden');
        expect(result.current.orderStatuses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Manejo de error (excepción)
    // ---------------------------------------------------------------
    it('debería capturar el error cuando la action lanza una excepción', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useOrderStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
        expect(result.current.orderStatuses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Parámetro isActive
    // ---------------------------------------------------------------
    it('debería llamar a la action con isActive=true por defecto', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockResolvedValue(defaultSuccessResponse);

        renderHook(() => useOrderStatuses());

        await waitFor(() => {
            expect(getOrderStatusesSimpleAction).toHaveBeenCalledWith(true);
        });
    });

    it('debería llamar a la action con isActive=false cuando se pasa como parámetro', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockResolvedValue(defaultSuccessResponse);

        renderHook(() => useOrderStatuses(false));

        await waitFor(() => {
            expect(getOrderStatusesSimpleAction).toHaveBeenCalledWith(false);
        });
    });

    it('debería volver a cargar cuando cambia el parámetro isActive', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockResolvedValue(defaultSuccessResponse);

        const { rerender } = renderHook(
            ({ isActive }) => useOrderStatuses(isActive),
            { initialProps: { isActive: true } }
        );

        await waitFor(() => {
            expect(getOrderStatusesSimpleAction).toHaveBeenCalledTimes(1);
        });

        rerender({ isActive: false });

        await waitFor(() => {
            expect(getOrderStatusesSimpleAction).toHaveBeenCalledTimes(2);
        });

        expect(vi.mocked(getOrderStatusesSimpleAction).mock.calls[1][0]).toBe(false);
    });

    // ---------------------------------------------------------------
    // Función refresh
    // ---------------------------------------------------------------
    it('debería exponer la función refresh que vuelve a cargar los estados', async () => {
        vi.mocked(getOrderStatusesSimpleAction).mockResolvedValue(defaultSuccessResponse);

        const { result } = renderHook(() => useOrderStatuses());

        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getOrderStatusesSimpleAction).mock.calls.length;

        await act(async () => {
            await result.current.refresh();
        });

        expect(vi.mocked(getOrderStatusesSimpleAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });
});
