import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useResourceForm } from './useResourceForm';
import { Resource } from '../../domain/types';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    createResourceAction: vi.fn(),
    updateResourceAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import { createResourceAction, updateResourceAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeInitialResource = (overrides: Partial<Resource> = {}): Resource => ({
    id: 1,
    name: 'orders',
    description: 'Gestión de pedidos',
    business_type_id: 1,
    business_type_name: 'Restaurante',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useResourceForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con formData por defecto (name="", business_type_id=null)', () => {
        const { result } = renderHook(() => useResourceForm());

        expect(result.current.formData.name).toBe('');
        expect(result.current.formData.description).toBe('');
        expect(result.current.formData.business_type_id).toBeNull();
        expect(result.current.loading).toBe(false);
        expect(result.current.error).toBeNull();
    });

    it('debería popular formData con initialData', async () => {
        const data = makeInitialResource();

        const { result } = renderHook(() => useResourceForm(data));

        await waitFor(() => {
            expect(result.current.formData.name).toBe('orders');
        });

        expect(result.current.formData.description).toBe('Gestión de pedidos');
        expect(result.current.formData.business_type_id).toBe(1);
    });

    it('debería actualizar campo con handleChange', () => {
        const { result } = renderHook(() => useResourceForm());

        act(() => {
            result.current.handleChange('name', 'products');
        });

        expect(result.current.formData.name).toBe('products');
    });

    it('debería llamar createResourceAction sin initialData', async () => {
        vi.mocked(createResourceAction).mockResolvedValue({ success: true, message: 'OK', data: makeInitialResource() });

        const { result } = renderHook(() => useResourceForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(createResourceAction).toHaveBeenCalledOnce();
        expect(updateResourceAction).not.toHaveBeenCalled();
    });

    it('debería llamar updateResourceAction con initialData', async () => {
        vi.mocked(updateResourceAction).mockResolvedValue({ success: true, message: 'OK', data: makeInitialResource() });
        const data = makeInitialResource();

        const { result } = renderHook(() => useResourceForm(data));
        await waitFor(() => expect(result.current.formData.name).toBe('orders'));

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(updateResourceAction).toHaveBeenCalledWith(1, expect.objectContaining({ name: 'orders' }));
    });

    it('debería retornar true y llamar onSuccess cuando response.success es true', async () => {
        vi.mocked(createResourceAction).mockResolvedValue({ success: true, message: 'OK', data: makeInitialResource() });
        const onSuccess = vi.fn();

        const { result } = renderHook(() => useResourceForm(undefined, onSuccess));

        await act(async () => {
            await result.current.submit();
        });

        expect(onSuccess).toHaveBeenCalledOnce();
    });

    it('debería capturar error de response cuando success es false', async () => {
        vi.mocked(createResourceAction).mockResolvedValue({ success: false, message: 'Nombre duplicado', data: null as any });

        const { result } = renderHook(() => useResourceForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Nombre duplicado');
    });

    it('debería capturar error cuando submit lanza excepción', async () => {
        vi.mocked(createResourceAction).mockRejectedValue(new Error('Error de servidor'));

        const { result } = renderHook(() => useResourceForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Error de servidor');
    });

    it('debería manejar loading durante submit', async () => {
        let resolvePromise: (value: any) => void;
        vi.mocked(createResourceAction).mockReturnValue(
            new Promise(resolve => { resolvePromise = resolve; })
        );

        const { result } = renderHook(() => useResourceForm());

        expect(result.current.loading).toBe(false);

        let submitPromise: Promise<boolean>;
        act(() => {
            submitPromise = result.current.submit();
        });

        expect(result.current.loading).toBe(true);

        await act(async () => {
            resolvePromise!({ success: true, message: 'OK', data: makeInitialResource() });
            await submitPromise!;
        });

        expect(result.current.loading).toBe(false);
    });
});
