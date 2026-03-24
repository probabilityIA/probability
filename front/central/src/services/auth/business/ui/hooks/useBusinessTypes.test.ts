import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useBusinessTypes, useBusinessTypeForm } from './useBusinessTypes';
import { BusinessType } from '../../domain/types';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getBusinessTypesAction: vi.fn(),
    deleteBusinessTypeAction: vi.fn(),
    createBusinessTypeAction: vi.fn(),
    updateBusinessTypeAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any) =>
        err instanceof Error ? err.message : 'Error'
    ),
}));

import {
    getBusinessTypesAction,
    deleteBusinessTypeAction,
    createBusinessTypeAction,
    updateBusinessTypeAction,
} from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeBusinessType = (id: number, name: string): BusinessType => ({
    id,
    name,
    code: name.toLowerCase(),
    description: `Tipo ${name}`,
    icon: 'icon',
    is_active: true,
});

const defaultTypesResponse = {
    success: true,
    message: 'OK',
    data: [makeBusinessType(1, 'Restaurante'), makeBusinessType(2, 'Tienda')],
    pagination: { current_page: 1, per_page: 10, total: 2, last_page: 1, has_next: false, has_prev: false },
};

// =================================================================
// useBusinessTypes
// =================================================================

describe('useBusinessTypes', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y tipos vacíos', () => {
        vi.mocked(getBusinessTypesAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useBusinessTypes());

        expect(result.current.loading).toBe(true);
        expect(result.current.types).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    it('debería cargar tipos exitosamente', async () => {
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);

        const { result } = renderHook(() => useBusinessTypes());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.types).toHaveLength(2);
        expect(result.current.types[0].name).toBe('Restaurante');
        expect(result.current.error).toBeNull();
    });

    it('debería capturar error cuando getBusinessTypesAction falla', async () => {
        vi.mocked(getBusinessTypesAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useBusinessTypes());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
        expect(result.current.types).toEqual([]);
    });

    it('debería eliminar tipo y refrescar lista cuando deleteType tiene éxito', async () => {
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);
        vi.mocked(deleteBusinessTypeAction).mockResolvedValue({ success: true, message: 'OK' });

        const { result } = renderHook(() => useBusinessTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteType(1);
        });

        expect(deleteResult).toBe(true);
        expect(deleteBusinessTypeAction).toHaveBeenCalledWith(1);
    });

    it('debería capturar error cuando deleteType falla y retornar false', async () => {
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);
        vi.mocked(deleteBusinessTypeAction).mockRejectedValue(new Error('No se puede eliminar'));

        const { result } = renderHook(() => useBusinessTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteType(1);
        });

        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('No se puede eliminar');
    });

    it('debería exponer refresh que recarga tipos', async () => {
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);

        const { result } = renderHook(() => useBusinessTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getBusinessTypesAction).mock.calls.length;

        await act(async () => {
            await result.current.refresh();
        });

        expect(vi.mocked(getBusinessTypesAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });
});

// =================================================================
// useBusinessTypeForm
// =================================================================

describe('useBusinessTypeForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con formData por defecto', () => {
        const { result } = renderHook(() => useBusinessTypeForm());

        expect(result.current.formData.name).toBe('');
        expect(result.current.formData.code).toBe('');
        expect(result.current.formData.is_active).toBe(true);
        expect(result.current.loading).toBe(false);
        expect(result.current.error).toBeNull();
    });

    it('debería popular formData con initialData', async () => {
        const data = makeBusinessType(1, 'Restaurante');

        const { result } = renderHook(() => useBusinessTypeForm(data));

        await waitFor(() => {
            expect(result.current.formData.name).toBe('Restaurante');
        });

        expect(result.current.formData.code).toBe('restaurante');
    });

    it('debería llamar createBusinessTypeAction sin initialData', async () => {
        vi.mocked(createBusinessTypeAction).mockResolvedValue({ success: true, message: 'OK', data: {} });

        const { result } = renderHook(() => useBusinessTypeForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(createBusinessTypeAction).toHaveBeenCalledOnce();
        expect(updateBusinessTypeAction).not.toHaveBeenCalled();
    });

    it('debería llamar updateBusinessTypeAction con initialData', async () => {
        vi.mocked(updateBusinessTypeAction).mockResolvedValue({ success: true, message: 'OK', data: {} });
        const data = makeBusinessType(5, 'Tienda');

        const { result } = renderHook(() => useBusinessTypeForm(data));
        await waitFor(() => expect(result.current.formData.name).toBe('Tienda'));

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(updateBusinessTypeAction).toHaveBeenCalledWith(5, expect.objectContaining({ name: 'Tienda' }));
    });

    it('debería retornar true y llamar onSuccess al crear exitosamente', async () => {
        vi.mocked(createBusinessTypeAction).mockResolvedValue({ success: true, message: 'OK', data: {} });
        const onSuccess = vi.fn();

        const { result } = renderHook(() => useBusinessTypeForm(undefined, onSuccess));

        await act(async () => {
            await result.current.submit();
        });

        expect(onSuccess).toHaveBeenCalledOnce();
    });

    it('debería capturar error y retornar false cuando submit falla', async () => {
        vi.mocked(createBusinessTypeAction).mockRejectedValue(new Error('Nombre duplicado'));

        const { result } = renderHook(() => useBusinessTypeForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Nombre duplicado');
    });

    it('debería manejar loading durante submit', async () => {
        let resolvePromise: (value: any) => void;
        vi.mocked(createBusinessTypeAction).mockReturnValue(
            new Promise(resolve => { resolvePromise = resolve; })
        );

        const { result } = renderHook(() => useBusinessTypeForm());

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
