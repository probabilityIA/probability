import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useIntegrationTypes } from './useIntegrationTypes';

// -----------------------------------------------------------------
// Mocks
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getIntegrationTypesAction: vi.fn(),
    createIntegrationTypeAction: vi.fn(),
    updateIntegrationTypeAction: vi.fn(),
    deleteIntegrationTypeAction: vi.fn(),
}));

vi.mock('@/shared/utils/action-result', () => ({
    getActionError: vi.fn((err: any, fallback?: string) =>
        err instanceof Error ? err.message : fallback || 'Error'
    ),
}));

import {
    getIntegrationTypesAction,
    createIntegrationTypeAction,
    updateIntegrationTypeAction,
    deleteIntegrationTypeAction,
} from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------

const makeType = (id: number, name: string) => ({
    id, name, code: name.toLowerCase(), is_active: true, created_at: '', updated_at: '',
});

const defaultResponse = {
    success: true, message: 'OK',
    data: [makeType(1, 'Shopify'), makeType(2, 'WhatsApp')],
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useIntegrationTypes', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    it('debería iniciar con loading en true y tipos vacíos', () => {
        vi.mocked(getIntegrationTypesAction).mockReturnValue(new Promise(() => {}));
        const { result } = renderHook(() => useIntegrationTypes());
        expect(result.current.loading).toBe(true);
        expect(result.current.integrationTypes).toEqual([]);
    });

    it('debería cargar tipos exitosamente', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue(defaultResponse);
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.integrationTypes).toHaveLength(2);
    });

    it('debería pasar categoryId al action', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue(defaultResponse);
        renderHook(() => useIntegrationTypes(5));
        await waitFor(() => {
            expect(getIntegrationTypesAction).toHaveBeenCalledWith(5);
        });
    });

    it('debería capturar error de response cuando success es false', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue({ success: false, message: 'No autorizado', data: [] });
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.error).toBe('No autorizado');
    });

    it('debería capturar error cuando lanza excepción', async () => {
        vi.mocked(getIntegrationTypesAction).mockRejectedValue(new Error('Error de red'));
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));
        expect(result.current.error).toBe('Error de red');
    });

    it('debería crear tipo y refrescar', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue(defaultResponse);
        vi.mocked(createIntegrationTypeAction).mockResolvedValue({ success: true, message: 'OK', data: makeType(3, 'New') });
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let createResult: boolean | undefined;
        await act(async () => { createResult = await result.current.createIntegrationType({ name: 'New', category_id: 1 }); });
        expect(createResult).toBe(true);
    });

    it('debería capturar error al crear tipo cuando success es false', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue(defaultResponse);
        vi.mocked(createIntegrationTypeAction).mockResolvedValue({ success: false, message: 'Nombre duplicado', data: null as any });
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let createResult: boolean | undefined;
        await act(async () => { createResult = await result.current.createIntegrationType({ name: 'Dup', category_id: 1 }); });
        expect(createResult).toBe(false);
        expect(result.current.error).toBe('Nombre duplicado');
    });

    it('debería actualizar tipo y refrescar', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue(defaultResponse);
        vi.mocked(updateIntegrationTypeAction).mockResolvedValue({ success: true, message: 'OK', data: makeType(1, 'Updated') });
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let updateResult: boolean | undefined;
        await act(async () => { updateResult = await result.current.updateIntegrationType(1, { name: 'Updated' }); });
        expect(updateResult).toBe(true);
    });

    it('debería eliminar tipo y refrescar', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteIntegrationTypeAction).mockResolvedValue({ success: true, message: 'OK' });
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => { deleteResult = await result.current.deleteIntegrationType(1); });
        expect(deleteResult).toBe(true);
    });

    it('debería capturar error al eliminar tipo', async () => {
        vi.mocked(getIntegrationTypesAction).mockResolvedValue(defaultResponse);
        vi.mocked(deleteIntegrationTypeAction).mockRejectedValue(new Error('No se puede eliminar'));
        const { result } = renderHook(() => useIntegrationTypes());
        await waitFor(() => expect(result.current.loading).toBe(false));

        let deleteResult: boolean | undefined;
        await act(async () => { deleteResult = await result.current.deleteIntegrationType(1); });
        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('No se puede eliminar');
    });
});
