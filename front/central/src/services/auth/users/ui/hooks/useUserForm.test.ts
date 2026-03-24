import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useUserForm } from './useUserForm';
import { User } from '../../domain/types';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    createUserAction: vi.fn(),
    updateUserAction: vi.fn(),
}));

import { createUserAction, updateUserAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeInitialUser = (overrides: Partial<User> = {}): User => ({
    id: 1,
    name: 'Juan Pérez',
    email: 'juan@test.com',
    phone: '3001234567',
    avatar_url: 'https://example.com/avatar.png',
    is_active: true,
    is_super_user: false,
    scope_id: 2,
    scope_code: 'business',
    scope_name: 'Business',
    business_role_assignments: [
        { business_id: 1, business_name: 'Negocio A', role_id: 2, role_name: 'Admin' },
        { business_id: 3, business_name: 'Negocio B', role_id: 4, role_name: 'Editor' },
    ],
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useUserForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    // ---------------------------------------------------------------
    // Estado inicial
    // ---------------------------------------------------------------
    it('debería iniciar con formData por defecto (scope_id=2, business_ids=[])', () => {
        const { result } = renderHook(() => useUserForm());

        expect(result.current.formData.name).toBe('');
        expect(result.current.formData.email).toBe('');
        expect(result.current.formData.scope_id).toBe(2);
        expect(result.current.formData.business_ids).toEqual([]);
        expect(result.current.loading).toBe(false);
        expect(result.current.error).toBeNull();
        expect(result.current.successMessage).toBeNull();
    });

    it('debería popular formData con initialData incluyendo business_ids de assignments', async () => {
        const data = makeInitialUser();

        const { result } = renderHook(() => useUserForm(data));

        await waitFor(() => {
            expect(result.current.formData.name).toBe('Juan Pérez');
        });

        expect(result.current.formData.email).toBe('juan@test.com');
        expect(result.current.formData.business_ids).toEqual([1, 3]);
    });

    // ---------------------------------------------------------------
    // handleChange
    // ---------------------------------------------------------------
    it('debería actualizar campo con handleChange', () => {
        const { result } = renderHook(() => useUserForm());

        act(() => {
            result.current.handleChange('name', 'María López');
        });

        expect(result.current.formData.name).toBe('María López');
    });

    // ---------------------------------------------------------------
    // handleFileChange
    // ---------------------------------------------------------------
    it('debería manejar handleFileChange con archivo nuevo', () => {
        const { result } = renderHook(() => useUserForm());
        const file = new File(['content'], 'avatar.png', { type: 'image/png' });

        act(() => {
            result.current.handleFileChange(file);
        });

        expect(result.current.avatarFile).toBe(file);
    });

    it('debería marcar removeAvatar cuando se limpia archivo y existía avatar_url', async () => {
        const data = makeInitialUser({ avatar_url: 'https://example.com/avatar.png' });
        const { result } = renderHook(() => useUserForm(data));
        await waitFor(() => expect(result.current.formData.name).toBe('Juan Pérez'));

        act(() => {
            result.current.handleFileChange(null);
        });

        // removeAvatar is internal state but affects submit behavior
        // We test it indirectly via submit
        vi.mocked(updateUserAction).mockResolvedValue({ success: true, data: data, message: 'OK' });

        await act(async () => {
            await result.current.submit();
        });

        expect(updateUserAction).toHaveBeenCalledWith(1, expect.objectContaining({ remove_avatar: true }));
    });

    // ---------------------------------------------------------------
    // Submit - crear
    // ---------------------------------------------------------------
    it('debería llamar createUserAction sin initialData', async () => {
        vi.mocked(createUserAction).mockResolvedValue({ success: true, data: makeInitialUser(), message: 'OK', password: 'gen123' });

        const { result } = renderHook(() => useUserForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(createUserAction).toHaveBeenCalledOnce();
        expect(updateUserAction).not.toHaveBeenCalled();
    });

    it('debería mostrar successMessage con password cuando se crea usuario exitosamente', async () => {
        vi.mocked(createUserAction).mockResolvedValue({ success: true, data: makeInitialUser(), message: 'OK', password: 'generatedPass123' });

        const { result } = renderHook(() => useUserForm());

        await act(async () => {
            await result.current.submit();
        });

        expect(result.current.successMessage).toContain('generatedPass123');
    });

    it('debería NO llamar onSuccess cuando se crea usuario con password', async () => {
        vi.mocked(createUserAction).mockResolvedValue({ success: true, data: makeInitialUser(), message: 'OK', password: 'gen123' });
        const onSuccess = vi.fn();

        const { result } = renderHook(() => useUserForm(undefined, onSuccess));

        await act(async () => {
            await result.current.submit();
        });

        expect(onSuccess).not.toHaveBeenCalled();
    });

    // ---------------------------------------------------------------
    // Submit - actualizar
    // ---------------------------------------------------------------
    it('debería llamar updateUserAction con initialData', async () => {
        vi.mocked(updateUserAction).mockResolvedValue({ success: true, data: makeInitialUser(), message: 'OK' });
        const data = makeInitialUser();
        const onSuccess = vi.fn();

        const { result } = renderHook(() => useUserForm(data, onSuccess));
        await waitFor(() => expect(result.current.formData.name).toBe('Juan Pérez'));

        await act(async () => {
            await result.current.submit();
        });

        expect(updateUserAction).toHaveBeenCalledWith(1, expect.objectContaining({ name: 'Juan Pérez' }));
        expect(onSuccess).toHaveBeenCalledOnce();
    });

    // ---------------------------------------------------------------
    // Submit - errores
    // ---------------------------------------------------------------
    it('debería capturar error de response cuando success es false', async () => {
        vi.mocked(createUserAction).mockResolvedValue({ success: false, data: null as any, message: 'Email duplicado' });

        const { result } = renderHook(() => useUserForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Email duplicado');
    });

    it('debería capturar error cuando submit lanza excepción', async () => {
        vi.mocked(createUserAction).mockRejectedValue(new Error('Error de servidor'));

        const { result } = renderHook(() => useUserForm());

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Error de servidor');
    });
});
