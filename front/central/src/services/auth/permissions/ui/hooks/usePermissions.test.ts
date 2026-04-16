import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { usePermissions } from './usePermissions';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getPermissionsAction: vi.fn(),
    deletePermissionAction: vi.fn(),
}));

import { getPermissionsAction, deletePermissionAction } from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makePermission = (id: number, name: string) => ({
    id,
    name,
    code: name.replace(':', '_'),
    resource: 'orders',
    action: 'read',
    resource_id: 1,
    action_id: 1,
    scope_id: 1,
    scope_name: 'Platform',
    scope_code: 'platform',
    business_type_id: 1,
    business_type_name: 'Restaurante',
});

const defaultResponse = {
    success: true,
    data: [makePermission(1, 'read:orders'), makePermission(2, 'write:orders')],
    total: 2,
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('usePermissions', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y permisos vacíos', () => {
        vi.mocked(getPermissionsAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => usePermissions());

        expect(result.current.loading).toBe(true);
        expect(result.current.permissions).toEqual([]);
        expect(result.current.error).toBeNull();
    });

    it('debería cargar permisos exitosamente', async () => {
        vi.mocked(getPermissionsAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => usePermissions());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.permissions).toHaveLength(2);
        expect(result.current.permissions[0].name).toBe('read:orders');
        expect(result.current.error).toBeNull();
    });

    it('debería capturar error cuando getPermissionsAction falla', async () => {
        vi.mocked(getPermissionsAction).mockRejectedValue(new Error('Error de red'));

        const { result } = renderHook(() => usePermissions());

        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
        expect(result.current.permissions).toEqual([]);
    });

    it('debería pasar filtro de nombre a la action', async () => {
        vi.mocked(getPermissionsAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => usePermissions());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setSearchName('read');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getPermissionsAction).mock.calls.at(-1)![0];
        expect(lastCall?.name).toBe('read');
    });

    it('debería pasar filtro de scope a la action', async () => {
        vi.mocked(getPermissionsAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => usePermissions());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setFilterScope('1');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getPermissionsAction).mock.calls.at(-1)![0];
        expect(lastCall?.scope_id).toBe(1);
    });

    it('debería pasar filtro de business type a la action', async () => {
        vi.mocked(getPermissionsAction).mockResolvedValue(defaultResponse);

        const { result } = renderHook(() => usePermissions());
        await waitFor(() => expect(result.current.loading).toBe(false));

        act(() => {
            result.current.setFilterBusinessType('2');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        const lastCall = vi.mocked(getPermissionsAction).mock.calls.at(-1)![0];
        expect(lastCall?.business_type_id).toBe(2);
    });

    it('debería eliminar permiso y refrescar lista', async () => {
        vi.mocked(getPermissionsAction).mockResolvedValue(defaultResponse);
        vi.mocked(deletePermissionAction).mockResolvedValue({ success: true, message: 'OK' });

        const { result } = renderHook(() => usePermissions());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deletePermission(1);
        });

        expect(deleteResult).toBe(true);
        expect(deletePermissionAction).toHaveBeenCalledWith(1);
    });

    it('debería capturar error cuando deletePermission falla', async () => {
        vi.mocked(getPermissionsAction).mockResolvedValue(defaultResponse);
        vi.mocked(deletePermissionAction).mockRejectedValue(new Error('No se puede eliminar'));

        const { result } = renderHook(() => usePermissions());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deletePermission(1);
        });

        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('No se puede eliminar');
    });
});
