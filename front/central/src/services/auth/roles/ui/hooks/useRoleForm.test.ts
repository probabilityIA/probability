import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useRoleForm } from './useRoleForm';
import { Role } from '../../domain/types';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    createRoleAction: vi.fn(),
    updateRoleAction: vi.fn(),
}));

import { createRoleAction, updateRoleAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeInitialRole = (overrides: Partial<Role> = {}): Role => ({
    id: 1,
    name: 'Admin',
    code: 'admin',
    description: 'Administrador',
    level: 1,
    is_system: false,
    scope_id: 1,
    scope_name: 'Platform',
    scope_code: 'platform',
    business_type_id: 1,
    business_type_name: 'Restaurante',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useRoleForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con formData por defecto (business_type_id=1, is_system=false)', () => {
        const { result } = renderHook(() => useRoleForm());

        expect(result.current.formData.name).toBe('');
        expect(result.current.formData.description).toBe('');
        expect(result.current.formData.business_type_id).toBe(1);
        expect(result.current.formData.is_system).toBe(false);
        expect(result.current.loading).toBe(false);
        expect(result.current.error).toBeNull();
    });

    it('debería popular formData con initialData', async () => {
        const data = makeInitialRole();

        const { result } = renderHook(() => useRoleForm(data));

        await waitFor(() => {
            expect(result.current.formData.name).toBe('Admin');
        });

        expect(result.current.formData.level).toBe(1);
        expect(result.current.formData.scope_id).toBe(1);
    });

    it('debería actualizar campo con handleChange', () => {
        const { result } = renderHook(() => useRoleForm());

        act(() => {
            result.current.handleChange('name', 'Editor');
        });

        expect(result.current.formData.name).toBe('Editor');
    });

    it('debería llamar createRoleAction sin initialData con business_type_id=1 forzado', async () => {
        vi.mocked(createRoleAction).mockResolvedValue({ success: true, data: makeInitialRole(), message: 'OK' });

        const { result } = renderHook(() => useRoleForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(createRoleAction).toHaveBeenCalledOnce();
        expect(createRoleAction).toHaveBeenCalledWith(expect.objectContaining({ business_type_id: 1 }));
        expect(updateRoleAction).not.toHaveBeenCalled();
    });

    it('debería llamar updateRoleAction con initialData con business_type_id=1 forzado', async () => {
        vi.mocked(updateRoleAction).mockResolvedValue({ success: true, data: makeInitialRole(), message: 'OK' });
        const data = makeInitialRole({ business_type_id: 5 });

        const { result } = renderHook(() => useRoleForm(data));
        await waitFor(() => expect(result.current.formData.name).toBe('Admin'));

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(updateRoleAction).toHaveBeenCalledWith(1, expect.objectContaining({ business_type_id: 1 }));
    });

    it('debería retornar true y llamar onSuccess cuando submit tiene éxito', async () => {
        vi.mocked(createRoleAction).mockResolvedValue({ success: true, data: makeInitialRole(), message: 'OK' });
        const onSuccess = vi.fn();

        const { result } = renderHook(() => useRoleForm(undefined, onSuccess));

        await act(async () => {
            await result.current.submit();
        });

        expect(onSuccess).toHaveBeenCalledOnce();
    });

    it('debería capturar error y retornar false cuando submit falla', async () => {
        vi.mocked(createRoleAction).mockRejectedValue(new Error('Nombre duplicado'));

        const { result } = renderHook(() => useRoleForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Nombre duplicado');
        expect(result.current.loading).toBe(false);
    });
});
