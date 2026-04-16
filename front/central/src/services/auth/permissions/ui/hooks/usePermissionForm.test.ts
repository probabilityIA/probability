import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { usePermissionForm } from './usePermissionForm';
import { Permission } from '../../domain/types';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    createPermissionAction: vi.fn(),
    updatePermissionAction: vi.fn(),
}));

vi.mock('@/services/auth/resources/infra/actions', () => ({
    getResourcesAction: vi.fn(),
}));

vi.mock('@/services/auth/actions/infra/actions', () => ({
    getActionsAction: vi.fn(),
}));

import { createPermissionAction, updatePermissionAction } from '../../infra/actions';
import { getResourcesAction } from '@/services/auth/resources/infra/actions';
import { getActionsAction } from '@/services/auth/actions/infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeResourcesResponse = () => ({
    success: true,
    message: 'OK',
    data: {
        resources: [
            { id: 1, name: 'orders', description: 'Pedidos', business_type_id: 1, business_type_name: 'Restaurante', created_at: '', updated_at: '' },
            { id: 2, name: 'products', description: 'Productos', business_type_id: 1, business_type_name: 'Restaurante', created_at: '', updated_at: '' },
        ],
        total: 2,
        page: 1,
        page_size: 100,
        total_pages: 1,
    },
});

const makeActionsResponse = () => ({
    success: true,
    message: 'OK',
    data: {
        actions: [
            { id: 1, name: 'read', description: 'Lectura', created_at: '', updated_at: '' },
            { id: 2, name: 'write', description: 'Escritura', created_at: '', updated_at: '' },
        ],
        total: 2,
        page: 1,
        page_size: 100,
        total_pages: 1,
    },
});

const makeInitialPermission = (overrides: Partial<Permission> = {}): Permission => ({
    id: 1,
    name: 'read:orders',
    code: 'read_orders',
    description: 'Leer pedidos',
    resource: 'orders',
    action: 'read',
    resource_id: 1,
    action_id: 1,
    scope_id: 1,
    scope_name: 'Platform',
    scope_code: 'platform',
    business_type_id: 1,
    business_type_name: 'Restaurante',
    ...overrides,
});

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('usePermissionForm', () => {
    beforeEach(() => {
        vi.clearAllMocks();
        vi.mocked(getResourcesAction).mockResolvedValue(makeResourcesResponse());
        vi.mocked(getActionsAction).mockResolvedValue(makeActionsResponse());
    });

    it('debería iniciar con formData por defecto y loadingData en true', () => {
        vi.mocked(getResourcesAction).mockReturnValue(new Promise(() => {}));
        vi.mocked(getActionsAction).mockReturnValue(new Promise(() => {}));

        const { result } = renderHook(() => usePermissionForm());

        expect(result.current.formData.name).toBe('');
        expect(result.current.formData.business_type_id).toBe(1);
        expect(result.current.loadingData).toBe(true);
        expect(result.current.loading).toBe(false);
    });

    it('debería cargar recursos y acciones al montar', async () => {
        const { result } = renderHook(() => usePermissionForm());

        await waitFor(() => {
            expect(result.current.loadingData).toBe(false);
        });

        expect(getResourcesAction).toHaveBeenCalledWith({ page_size: 100 });
        expect(getActionsAction).toHaveBeenCalledWith({ page_size: 100 });
        expect(result.current.resources).toHaveLength(2);
        expect(result.current.actions).toHaveLength(2);
    });

    it('debería manejar error en carga de datos sin romper el hook', async () => {
        vi.mocked(getResourcesAction).mockRejectedValue(new Error('Error'));
        vi.mocked(getActionsAction).mockRejectedValue(new Error('Error'));

        const { result } = renderHook(() => usePermissionForm());

        await waitFor(() => {
            expect(result.current.loadingData).toBe(false);
        });

        expect(result.current.resources).toEqual([]);
        expect(result.current.actions).toEqual([]);
    });

    it('debería popular formData con initialData cuando se proporciona', async () => {
        const initialData = makeInitialPermission();

        const { result } = renderHook(() => usePermissionForm(initialData));

        await waitFor(() => {
            expect(result.current.formData.name).toBe('read:orders');
        });

        expect(result.current.formData.resource_id).toBe(1);
        expect(result.current.formData.action_id).toBe(1);
        expect(result.current.formData.scope_id).toBe(1);
        expect(result.current.formData.business_type_id).toBe(1);
    });

    it('debería actualizar campo con handleChange', async () => {
        const { result } = renderHook(() => usePermissionForm());
        await waitFor(() => expect(result.current.loadingData).toBe(false));

        act(() => {
            result.current.handleChange('name', 'write:orders');
        });

        expect(result.current.formData.name).toBe('write:orders');
    });

    it('debería llamar createPermissionAction sin initialData con business_type_id=1', async () => {
        vi.mocked(createPermissionAction).mockResolvedValue({ success: true, data: makeInitialPermission(), message: 'OK' });

        const { result } = renderHook(() => usePermissionForm());
        await waitFor(() => expect(result.current.loadingData).toBe(false));

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(createPermissionAction).toHaveBeenCalledOnce();
        expect(createPermissionAction).toHaveBeenCalledWith(expect.objectContaining({ business_type_id: 1 }));
        expect(updatePermissionAction).not.toHaveBeenCalled();
    });

    it('debería llamar updatePermissionAction con initialData', async () => {
        vi.mocked(updatePermissionAction).mockResolvedValue({ success: true, data: makeInitialPermission(), message: 'OK' });
        const initialData = makeInitialPermission();

        const { result } = renderHook(() => usePermissionForm(initialData));
        await waitFor(() => expect(result.current.formData.name).toBe('read:orders'));

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(true);
        expect(updatePermissionAction).toHaveBeenCalledWith(1, expect.objectContaining({ name: 'read:orders', business_type_id: 1 }));
    });

    it('debería retornar true y llamar onSuccess cuando submit tiene éxito', async () => {
        vi.mocked(createPermissionAction).mockResolvedValue({ success: true, data: makeInitialPermission(), message: 'OK' });
        const onSuccess = vi.fn();

        const { result } = renderHook(() => usePermissionForm(undefined, onSuccess));
        await waitFor(() => expect(result.current.loadingData).toBe(false));

        await act(async () => {
            await result.current.submit();
        });

        expect(onSuccess).toHaveBeenCalledOnce();
    });

    it('debería capturar error y retornar false cuando submit falla', async () => {
        vi.mocked(createPermissionAction).mockRejectedValue(new Error('Nombre duplicado'));

        const { result } = renderHook(() => usePermissionForm());
        await waitFor(() => expect(result.current.loadingData).toBe(false));

        let submitResult: boolean | undefined;
        await act(async () => {
            submitResult = await result.current.submit();
        });

        expect(submitResult).toBe(false);
        expect(result.current.error).toBe('Nombre duplicado');
    });
});
