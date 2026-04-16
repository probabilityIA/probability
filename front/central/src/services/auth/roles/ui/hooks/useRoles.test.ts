import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useRoles } from './useRoles';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getRolesAction: vi.fn(),
    deleteRoleAction: vi.fn(),
}));

import { getRolesAction, deleteRoleAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeRole = (id: number, name: string) => ({
    id,
    name,
    code: name.toLowerCase(),
    description: `Rol ${name}`,
    level: 1,
    is_system: false,
    scope_id: 1,
    scope_name: 'Platform',
    scope_code: 'platform',
    business_type_id: 1,
    business_type_name: 'Restaurante',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
});

const generateRoles = (count: number) =>
    Array.from({ length: count }, (_, i) => makeRole(i + 1, `Role ${i + 1}`));

const makeDefaultResponse = (roles = [makeRole(1, 'Admin'), makeRole(2, 'Editor')]) => ({
    success: true,
    message: 'OK',
    data: roles,
    pagination: { current_page: 1, per_page: 10, total: roles.length, last_page: 1, has_next: false, has_prev: false },
});

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useRoles', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y roles vacíos', () => {
        vi.mocked(getRolesAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => useRoles());

        expect(result.current.loading).toBe(true);
        expect(result.current.roles).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    it('debería cargar roles exitosamente', async () => {
        vi.mocked(getRolesAction).mockResolvedValue(makeDefaultResponse());

        const { result } = renderHook(() => useRoles());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.roles).toHaveLength(2);
        expect(result.current.roles[0].name).toBe('Admin');
        expect(result.current.allRoles).toHaveLength(2);
        expect(result.current.error).toBeNull();
    });

    it('debería capturar error cuando getRolesAction falla', async () => {
        vi.mocked(getRolesAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => useRoles());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
    });

    it('debería paginar del lado del cliente correctamente', async () => {
        const roles25 = generateRoles(25);
        vi.mocked(getRolesAction).mockResolvedValue(makeDefaultResponse(roles25));

        const { result } = renderHook(() => useRoles());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        // Página 1: 20 roles (pageSize default = 20)
        expect(result.current.roles).toHaveLength(20);
        expect(result.current.allRoles).toHaveLength(25);

        // Cambiar a página 2
        act(() => {
            result.current.setPage(2);
        });

        expect(result.current.roles).toHaveLength(5);
    });

    it('debería calcular totalPages correctamente', async () => {
        const roles25 = generateRoles(25);
        vi.mocked(getRolesAction).mockResolvedValue(makeDefaultResponse(roles25));

        const { result } = renderHook(() => useRoles());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.total).toBe(25);
        expect(result.current.totalPages).toBe(2); // 25 / 20 = 2
    });

    it('debería aplicar filtros al refrescar', async () => {
        vi.mocked(getRolesAction).mockResolvedValue(makeDefaultResponse());

        const { result } = renderHook(() => useRoles());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setFilters({ name: 'Admin', page: 1, page_size: 20 });
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getRolesAction).mock.calls.at(-1)![0];
        expect(lastCall?.name).toBe('Admin');
    });

    it('debería eliminar rol y refrescar lista', async () => {
        vi.mocked(getRolesAction).mockResolvedValue(makeDefaultResponse());
        vi.mocked(deleteRoleAction).mockResolvedValue({ success: true, message: 'OK' });

        const { result } = renderHook(() => useRoles());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteRole(1);
        });

        expect(deleteResult).toBe(true);
        expect(deleteRoleAction).toHaveBeenCalledWith(1);
    });

    it('debería capturar error cuando deleteRole falla', async () => {
        vi.mocked(getRolesAction).mockResolvedValue(makeDefaultResponse());
        vi.mocked(deleteRoleAction).mockRejectedValue(new Error('No se puede eliminar'));

        const { result } = renderHook(() => useRoles());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteRole(1);
        });

        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('No se puede eliminar');
    });

    it('debería exponer refresh que recarga roles', async () => {
        vi.mocked(getRolesAction).mockResolvedValue(makeDefaultResponse());

        const { result } = renderHook(() => useRoles());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getRolesAction).mock.calls.length;

        await act(async () => {
            await result.current.refresh();
        });

        expect(vi.mocked(getRolesAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });
});
