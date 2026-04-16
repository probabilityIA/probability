import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useBusinessForm } from './useBusinessForm';
import { Business } from '../../domain/types';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    createBusinessAction: vi.fn(),
    updateBusinessAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import { createBusinessAction, updateBusinessAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeInitialData = (overrides: Partial<Business> = {}): Business => ({
    id: 1,
    name: 'Mi Negocio',
    code: 'mi-negocio',
    business_type_id: 1,
    is_active: true,
    timezone: 'America/Bogota',
    address: 'Calle 123',
    description: 'Descripción',
    primary_color: '#ff0000',
    secondary_color: '#00ff00',
    tertiary_color: '#0000ff',
    quaternary_color: '#ffffff',
    enable_delivery: true,
    enable_pickup: false,
    enable_reservations: false,
    ...overrides,
} as Business);

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useBusinessForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    // ---------------------------------------------------------------
    // Estado inicial
    // ---------------------------------------------------------------
    it('debería iniciar con formData por defecto cuando no hay initialData', () => {
        const { result } = renderHook(() => useBusinessForm());

        expect(result.current.formData.name).toBe('');
        expect(result.current.formData.code).toBe('');
        expect(result.current.formData.business_type_id).toBe(1);
        expect(result.current.formData.is_active).toBe(true);
        expect(result.current.loading).toBe(false);
        expect(result.current.error).toBeNull();
    });

    it('debería popular formData con initialData cuando se proporciona', async () => {
        const data = makeInitialData();
        const { result } = renderHook(() => useBusinessForm(data));

        await waitFor(() => {
            expect(result.current.formData.name).toBe('Mi Negocio');
        });

        expect(result.current.formData.code).toBe('mi-negocio');
        expect(result.current.formData.primary_color).toBe('#ff0000');
        expect(result.current.formData.enable_delivery).toBe(true);
    });

    // ---------------------------------------------------------------
    // handleChange
    // ---------------------------------------------------------------
    it('debería actualizar un campo con handleChange', () => {
        const { result } = renderHook(() => useBusinessForm());

        act(() => {
            result.current.handleChange('name', 'Nuevo Nombre');
        });

        expect(result.current.formData.name).toBe('Nuevo Nombre');
    });

    // ---------------------------------------------------------------
    // Submit - crear
    // ---------------------------------------------------------------
    it('debería llamar createBusinessAction al hacer submit sin initialData', async () => {
        vi.mocked(createBusinessAction).mockResolvedValue({ success: true, message: 'OK', data: {} });

        const { result } = renderHook(() => useBusinessForm());

        act(() => {
            result.current.handleChange('name', 'Negocio Nuevo');
        });

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(createBusinessAction).toHaveBeenCalledOnce();
        expect(updateBusinessAction).not.toHaveBeenCalled();
    });

    // ---------------------------------------------------------------
    // Submit - actualizar
    // ---------------------------------------------------------------
    it('debería llamar updateBusinessAction al hacer submit con initialData', async () => {
        vi.mocked(updateBusinessAction).mockResolvedValue({ success: true, message: 'OK', data: {} });
        const data = makeInitialData();

        const { result } = renderHook(() => useBusinessForm(data));

        await waitFor(() => {
            expect(result.current.formData.name).toBe('Mi Negocio');
        });

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(updateBusinessAction).toHaveBeenCalledOnce();
        expect(updateBusinessAction).toHaveBeenCalledWith(1, expect.objectContaining({ name: 'Mi Negocio' }));
        expect(createBusinessAction).not.toHaveBeenCalled();
    });

    // ---------------------------------------------------------------
    // Submit - onSuccess callback
    // ---------------------------------------------------------------
    it('debería llamar onSuccess cuando submit tiene éxito', async () => {
        vi.mocked(createBusinessAction).mockResolvedValue({ success: true, message: 'OK', data: {} });
        const onSuccess = vi.fn();

        const { result } = renderHook(() => useBusinessForm(undefined, onSuccess));

        await act(async () => {
            await result.current.submit();
        });

        expect(onSuccess).toHaveBeenCalledOnce();
    });

    // ---------------------------------------------------------------
    // Submit - error
    // ---------------------------------------------------------------
    it('debería capturar error y retornar false cuando submit falla', async () => {
        vi.mocked(createBusinessAction).mockRejectedValue(new Error('Nombre duplicado'));

        const { result } = renderHook(() => useBusinessForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Nombre duplicado');
        expect(result.current.loading).toBe(false);
    });

    // ---------------------------------------------------------------
    // Loading state
    // ---------------------------------------------------------------
    it('debería manejar loading correctamente durante submit', async () => {
        let resolvePromise: (value: any) => void;
        vi.mocked(createBusinessAction).mockReturnValue(
            new Promise(resolve => { resolvePromise = resolve; })
        );

        const { result } = renderHook(() => useBusinessForm());

        expect(result.current.loading).toBe(false);

        let submitPromise: Promise<boolean>;
        act(() => {
            submitPromise = result.current.submit();
        });

        expect(result.current.loading).toBe(true);

        await act(async () => {
            resolvePromise!({ success: true, message: 'OK', data: {} });
            await submitPromise!;
        });

        expect(result.current.loading).toBe(false);
    });
});
