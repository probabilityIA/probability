import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useUsers } from './useUsers';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getUsersAction: vi.fn(),
    deleteUserAction: vi.fn(),
}));

import { getUsersAction, deleteUserAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeUser = (id: number, name: string) => ({
    id,
    name,
    email: `${name.toLowerCase().replace(' ', '.')}@test.com`,
    phone: '3001234567',
    is_active: true,
    is_super_user: false,
    scope_id: 2,
    scope_code: 'business',
    scope_name: 'Business',
    business_role_assignments: [],
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
});

const makePagination = (overrides = {}) => ({
    current_page: 1,
    per_page: 10,
    total: 2,
    last_page: 1,
    has_next: false,
    has_prev: false,
    ...overrides,
});

const defaultResponse = {
    success: true,
    data: [makeUser(1, 'Juan Pérez'), makeUser(2, 'María López')],
    pagination: makePagination(),
    message: 'OK',
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useUsers', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y usuarios vacíos', () => {
        vi.mocked(getUsersAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useUsers());

        expect(result.current.loading).toBe(true);
        expect(result.current.users).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    it('debería cargar usuarios exitosamente', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useUsers());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.users).toHaveLength(2);
        expect(result.current.pagination).toEqual(makePagination());
        expect(result.current.error).toBeNull();
    });

    it('debería capturar error cuando getUsersAction falla', async () => {
        vi.mocked(getUsersAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useUsers());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
    });

    it('debería cambiar de página y refrescar', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useUsers());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setPage(2);
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getUsersAction).mock.calls.at(-1)![0];
        expect(lastCall?.page).toBe(2);
    });

    it('debería pasar filtro de nombre', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useUsers());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setSearchName('Juan');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getUsersAction).mock.calls.at(-1)![0];
        expect(lastCall?.name).toBe('Juan');
    });

    it('debería pasar filtro de email', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useUsers());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setSearchEmail('juan@test.com');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getUsersAction).mock.calls.at(-1)![0];
        expect(lastCall?.email).toBe('juan@test.com');
    });

    it('debería pasar filtro de is_active como boolean', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useUsers());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setFilterIsActive('true');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getUsersAction).mock.calls.at(-1)![0];
        expect(lastCall?.is_active).toBe(true);

        act(() => {
            result.current.setFilterIsActive('false');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall2 = vi.mocked(getUsersAction).mock.calls.at(-1)![0];
        expect(lastCall2?.is_active).toBe(false);
    });

    it('debería pasar filtro de role_id como número', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => useUsers());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setFilterRoleId('3');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getUsersAction).mock.calls.at(-1)![0];
        expect(lastCall?.role_id).toBe(3);
    });

    it('debería eliminar usuario y refrescar lista', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteUserAction).mockResolvedValue({ success: true, message: 'OK' });

        const { result } = renderHook(() => useUsers());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteUser(1);
        });

        expect(deleteResult).toBe(true);
        expect(deleteUserAction).toHaveBeenCalledWith(1);
    });

    it('debería capturar error cuando deleteUser falla', async () => {
        vi.mocked(getUsersAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteUserAction).mockRejectedValue(new Error('No se puede eliminar'));

        const { result } = renderHook(() => useUsers());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteUser(1);
        });

        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('No se puede eliminar');
    });
});
