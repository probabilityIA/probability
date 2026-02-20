import { describe, it, expect, vi, beforeEach } from 'vitest';
import { renderHook, waitFor, act } from '@testing-library/react';
import { useBusinesses } from './useBusinesses';

// -----------------------------------------------------------------
// Mock del módulo de Server Actions
// El archivo de actions usa 'use server' y depende de headers de Next.js,
// por lo que debe ser mockeado completamente para correr en jsdom.
// -----------------------------------------------------------------

vi.mock('../../infra/actions', () => ({
    getBusinessesAction: vi.fn(),
    deleteBusinessAction: vi.fn(),
    getBusinessTypesAction: vi.fn(),
}));

// Importar los mocks tipados DESPUÉS del vi.mock
import {
    getBusinessesAction,
    deleteBusinessAction,
    getBusinessTypesAction,
} from '../../infra/actions';

// -----------------------------------------------------------------
// Helpers: datos de prueba
// -----------------------------------------------------------------

const makePagination = (overrides = {}) => ({
    current_page: 1,
    per_page: 10,
    total: 2,
    last_page: 1,
    has_next: false,
    has_prev: false,
    ...overrides,
});

const makeBusiness = (id: number, name: string) => ({
    id,
    name,
    business_type_id: 1,
    is_active: true,
});

const makeBusinessType = (id: number, name: string) => ({
    id,
    name,
    code: name.toLowerCase(),
    is_active: true,
});

const defaultBusinessesResponse = {
    success: true,
    message: 'OK',
    data: [makeBusiness(1, 'Negocio A'), makeBusiness(2, 'Negocio B')],
    pagination: makePagination(),
};

const defaultTypesResponse = {
    success: true,
    message: 'OK',
    data: [makeBusinessType(1, 'Restaurante'), makeBusinessType(2, 'Tienda')],
    pagination: makePagination(),
};

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('useBusinesses', () => {
    beforeEach(() => {
        vi.clearAllMocks();
    });

    // ---------------------------------------------------------------
    // Estado inicial
    // ---------------------------------------------------------------
    it('debería iniciar con loading en true y sin negocios', () => {
        // Arrange: las actions nunca resuelven durante esta verificación
        vi.mocked(getBusinessesAction).mockReturnValue(new Promise(() => {}));
        vi.mocked(getBusinessTypesAction).mockReturnValue(new Promise(() => {}));

        // Act
        const { result } = renderHook(() => useBusinesses());

        // Assert: estado inicial
        expect(result.current.loading).toBe(true);
        expect(result.current.businesses).toEqual([]);
        expect(result.current.error).toBeNull();
        expect(result.current.page).toBe(1);
        expect(result.current.totalPages).toBe(1);
    });

    // ---------------------------------------------------------------
    // Carga exitosa de businesses y tipos
    // ---------------------------------------------------------------
    it('debería cargar negocios y tipos de negocio exitosamente', async () => {
        // Arrange
        vi.mocked(getBusinessesAction).mockResolvedValue(defaultBusinessesResponse);
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);

        // Act
        const { result } = renderHook(() => useBusinesses());

        // Assert: esperar a que termine la carga
        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.businesses).toHaveLength(2);
        expect(result.current.businesses[0].name).toBe('Negocio A');
        expect(result.current.businesses[1].name).toBe('Negocio B');
        expect(result.current.businessTypes).toHaveLength(2);
        expect(result.current.error).toBeNull();
        expect(result.current.totalPages).toBe(1);
    });

    // ---------------------------------------------------------------
    // Manejo de error en la carga de negocios
    // ---------------------------------------------------------------
    it('debería capturar el error cuando getBusinessesAction falla', async () => {
        // Arrange
        vi.mocked(getBusinessesAction).mockRejectedValue(new Error('Error de red'));
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);

        // Act
        const { result } = renderHook(() => useBusinesses());

        // Assert
        await waitFor(() => {
            expect(result.current.loading).toBe(false);
        });

        expect(result.current.error).toBe('Error de red');
        expect(result.current.businesses).toEqual([]);
    });

    // ---------------------------------------------------------------
    // Paginación
    // ---------------------------------------------------------------
    it('debería refrescar negocios cuando la página cambia', async () => {
        // Arrange
        vi.mocked(getBusinessesAction).mockResolvedValue(defaultBusinessesResponse);
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);

        const { result } = renderHook(() => useBusinesses());

        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBeforePageChange = vi.mocked(getBusinessesAction).mock.calls.length;

        // Act: cambiar página
        act(() => {
            result.current.setPage(2);
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        // Assert: getBusinessesAction debe haberse llamado otra vez con page=2
        const callsAfterPageChange = vi.mocked(getBusinessesAction).mock.calls.length;
        expect(callsAfterPageChange).toBeGreaterThan(callsBeforePageChange);

        const lastCall = vi.mocked(getBusinessesAction).mock.calls.at(-1)![0];
        expect(lastCall?.page).toBe(2);
    });

    // ---------------------------------------------------------------
    // Filtro por nombre
    // ---------------------------------------------------------------
    it('debería pasar el filtro de nombre a la action', async () => {
        // Arrange
        vi.mocked(getBusinessesAction).mockResolvedValue(defaultBusinessesResponse);
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);

        const { result } = renderHook(() => useBusinesses());
        await waitFor(() => expect(result.current.loading).toBe(false));

        // Act: cambiar nombre de búsqueda
        act(() => {
            result.current.setSearchName('Negocio A');
        });

        await waitFor(() => expect(result.current.loading).toBe(false));

        // Assert: la última llamada debe incluir el nombre
        const lastCall = vi.mocked(getBusinessesAction).mock.calls.at(-1)![0];
        expect(lastCall?.name).toBe('Negocio A');
    });

    // ---------------------------------------------------------------
    // Eliminar negocio - éxito
    // ---------------------------------------------------------------
    it('debería eliminar un negocio y refrescar la lista cuando deleteBusiness tiene éxito', async () => {
        // Arrange
        vi.mocked(getBusinessesAction).mockResolvedValue(defaultBusinessesResponse);
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);
        vi.mocked(deleteBusinessAction).mockResolvedValue({ success: true, message: 'Eliminado' });

        const { result } = renderHook(() => useBusinesses());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBeforeDelete = vi.mocked(getBusinessesAction).mock.calls.length;

        // Act
        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteBusiness(1);
        });

        // Assert
        expect(deleteResult).toBe(true);
        expect(deleteBusinessAction).toHaveBeenCalledWith(1);
        // Debe haber hecho otra llamada para refrescar la lista
        expect(vi.mocked(getBusinessesAction).mock.calls.length).toBeGreaterThan(callsBeforeDelete);
    });

    // ---------------------------------------------------------------
    // Eliminar negocio - error
    // ---------------------------------------------------------------
    it('debería capturar el error cuando deleteBusinessAction falla', async () => {
        // Arrange
        vi.mocked(getBusinessesAction).mockResolvedValue(defaultBusinessesResponse);
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);
        vi.mocked(deleteBusinessAction).mockRejectedValue(new Error('No se puede eliminar'));

        const { result } = renderHook(() => useBusinesses());
        await waitFor(() => expect(result.current.loading).toBe(false));

        // Act
        let deleteResult: boolean | undefined;
        await act(async () => {
            deleteResult = await result.current.deleteBusiness(99);
        });

        // Assert
        expect(deleteResult).toBe(false);
        expect(result.current.error).toBe('No se puede eliminar');
    });

    // ---------------------------------------------------------------
    // Función refresh expuesta
    // ---------------------------------------------------------------
    it('debería exponer la función refresh que vuelve a cargar los negocios', async () => {
        // Arrange
        vi.mocked(getBusinessesAction).mockResolvedValue(defaultBusinessesResponse);
        vi.mocked(getBusinessTypesAction).mockResolvedValue(defaultTypesResponse);

        const { result } = renderHook(() => useBusinesses());
        await waitFor(() => expect(result.current.loading).toBe(false));

        const callsBefore = vi.mocked(getBusinessesAction).mock.calls.length;

        // Act: llamar refresh manualmente
        await act(async () => {
            await result.current.refresh();
        });

        // Assert
        expect(vi.mocked(getBusinessesAction).mock.calls.length).toBeGreaterThan(callsBefore);
    });

    // ---------------------------------------------------------------
    // Error silencioso en fetchTypes (no bloquea la UI)
    // ---------------------------------------------------------------
    it('debería cargar negocios aunque getBusinessTypesAction falle (error silencioso)', async () => {
        // Arrange: types falla silenciosamente, businesses carga OK
        vi.mocked(getBusinessesAction).mockResolvedValue(defaultBusinessesResponse);
        vi.mocked(getBusinessTypesAction).mockRejectedValue(new Error('Types service down'));

        // Act
        const { result } = renderHook(() => useBusinesses());

        await waitFor(() => expect(result.current.loading).toBe(false));

        // Assert: negocios cargados y sin error en la UI (error de types es silencioso)
        expect(result.current.businesses).toHaveLength(2);
        expect(result.current.error).toBeNull();
        expect(result.current.businessTypes).toEqual([]);
    });
});
