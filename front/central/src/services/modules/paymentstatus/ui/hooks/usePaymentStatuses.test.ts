import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { usePaymentStatuses } from './usePaymentStatuses';

// -----------------------------------------------------------------
// Mock del módulo de Server Actions
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getPaymentStatusesAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback: string) =>
        err instanceof Error ? err.message : fallback
    ),
}));

import { getPaymentStatusesAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers: datos de prueba
// -----------------------------------------------------------------

const makePaymentStatus = (id: number, code: string, name: string) => ({
    id,
    code,
    name,
    description: `Descripción de ${name}`,
    category: 'payment',
    color: '#10B981',
    icon: 'check',
    is_active: true,
});

const defaultSuccessResponse = {
    success: true,
    message: 'OK',
    data: [
        makePaymentStatus(1, 'pending', 'Pendiente'),
        makePaymentStatus(2, 'paid', 'Pagado'),
        makePaymentStatus(3, 'failed', 'Fallido'),
    ],
};

const emptyResponse = {
    success: true,
    message: 'OK',
    data: [],
};

const errorResponse = {
    success: false,
    message: 'Error al cargar estados de pago',
    data: [],
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('usePaymentStatuses', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    // ---------------------------------------------------------------
    // Estado inicial
    // ---------------------------------------------------------------
    it('debería iniciar con loading en true y sin estados de pago', () => {
        vi.mocked(getPaymentStatusesAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => usePaymentStatuses());

        expect(result.current.loading).toBe(true);
        expect(result.current.paymentStatuses).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    // ---------------------------------------------------------------
    // Carga exitosa
    // ---------------------------------------------------------------
    it('debería cargar estados de pago exitosamente', async () => {
        vi.mocked(getPaymentStatusesAction).mockResolvedValue(defaultSuccessResponse);

        const { result } = renderHook(() => usePaymentStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.paymentStatuses).toHaveLength(3);
        expect(result.current.paymentStatuses[0].code).toBe('pending');
        expect(result.current.paymentStatuses[1].name).toBe('Pagado');
        expect(result.current.paymentStatuses[2].code).toBe('failed');
        expect(result.current.error).toBeNull();
    });

    // ---------------------------------------------------------------
    // Respuesta vacía exitosa
    // ---------------------------------------------------------------
    it('debería manejar una respuesta vacía exitosa', async () => {
        vi.mocked(getPaymentStatusesAction).mockResolvedValue(emptyResponse);

        const { result } = renderHook(() => usePaymentStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.paymentStatuses).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    // ---------------------------------------------------------------
    // Respuesta con success: false
    // ---------------------------------------------------------------
    it('debería establecer error cuando la respuesta tiene success: false', async () => {
        vi.mocked(getPaymentStatusesAction).mockResolvedValue(errorResponse);

        const { result } = renderHook(() => usePaymentStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error al cargar estados de pago');
        expect(result.current.paymentStatuses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Manejo de error (excepción)
    // ---------------------------------------------------------------
    it('debería capturar el error cuando la action lanza una excepción', async () => {
        vi.mocked(getPaymentStatusesAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => usePaymentStatuses());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
        expect(result.current.paymentStatuses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Parámetro isActive
    // ---------------------------------------------------------------
    it('debería llamar a la action con is_active=true por defecto', async () => {
        vi.mocked(getPaymentStatusesAction).mockResolvedValue(defaultSuccessResponse);

        renderHook(() => usePaymentStatuses());

        await waitFor(() => {
            expect(getPaymentStatusesAction).toHaveBeenCalledWith({ is_active: true });
        });
    });

    it('debería llamar a la action con is_active=false cuando se pasa como parámetro', async () => {
        vi.mocked(getPaymentStatusesAction).mockResolvedValue(defaultSuccessResponse);

        renderHook(() => usePaymentStatuses(false));

        await waitFor(() => {
            expect(getPaymentStatusesAction).toHaveBeenCalledWith({ is_active: false });
        });
    });

    it('debería volver a cargar cuando cambia el parámetro isActive', async () => {
        vi.mocked(getPaymentStatusesAction).mockResolvedValue(defaultSuccessResponse);

        const { rerender } = renderHook(
            ({ isActive }) => usePaymentStatuses(isActive),
            { initialProps: { isActive: true } }
        );

        await waitFor(() => {
            expect(getPaymentStatusesAction).toHaveBeenCalledTimes(1);
        });

        rerender({ isActive: false });

        await waitFor(() => {
            expect(getPaymentStatusesAction).toHaveBeenCalledTimes(2);
        });

        expect(vi.mocked(getPaymentStatusesAction).mock.calls[1][0]).toEqual({ is_active: false });
    });

    // ---------------------------------------------------------------
    // Función refresh
    // ---------------------------------------------------------------
    it('debería exponer la función refresh que vuelve a cargar los estados', async () => {
        vi.mocked(getPaymentStatusesAction).mockResolvedValue(defaultSuccessResponse);

        const { result } = renderHook(() => usePaymentStatuses());

        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getPaymentStatusesAction).mock.calls.length;

        await act(async () => {
            await result.current.refresh();
        });

        expect(vi.mocked(getPaymentStatusesAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });
});
