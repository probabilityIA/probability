import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ActionUseCases } from './use-cases';
import { IActionRepository } from '../domain/ports';
import {
    Action,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeAction = (overrides: Partial<Action> = {}): Action => ({
    id: 1,
    name: 'read',
    description: 'Acción de lectura',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const paginatedActions: PaginatedResponse<Action> = {
    success: true,
    message: 'OK',
    data: {
        actions: [makeAction(), makeAction({ id: 2, name: 'write' })],
        total: 2,
        page: 1,
        page_size: 10,
        total_pages: 1,
    },
};

const singleAction: SingleResponse<Action> = {
    success: true,
    message: 'OK',
    data: makeAction(),
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IActionRepository {
    return {
        getActions: vi.fn(),
        getActionById: vi.fn(),
        createAction: vi.fn(),
        updateAction: vi.fn(),
        deleteAction: vi.fn(),
    } as unknown as IActionRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('ActionUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: ActionUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new ActionUseCases(repo as unknown as IActionRepository);
    });

    // ---------------------------------------------------------------
    // getActions
    // ---------------------------------------------------------------
    describe('getActions', () => {
        it('debería retornar la lista paginada de acciones cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getActions).mockResolvedValue(paginatedActions);

            const result = await useCases.getActions({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedActions);
            expect(repo.getActions).toHaveBeenCalledOnce();
            expect(repo.getActions).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getActions).mockResolvedValue(paginatedActions);

            await useCases.getActions();

            expect(repo.getActions).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getActions).mockRejectedValue(new Error('Fallo de base de datos'));

            await expect(useCases.getActions()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getActionById
    // ---------------------------------------------------------------
    describe('getActionById', () => {
        it('debería retornar una acción por ID', async () => {
            vi.mocked(repo.getActionById).mockResolvedValue(singleAction);

            const result = await useCases.getActionById(1);

            expect(result).toEqual(singleAction);
            expect(repo.getActionById).toHaveBeenCalledOnce();
            expect(repo.getActionById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando la acción no existe', async () => {
            vi.mocked(repo.getActionById).mockRejectedValue(new Error('Acción no encontrada'));

            await expect(useCases.getActionById(999)).rejects.toThrow('Acción no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // createAction
    // ---------------------------------------------------------------
    describe('createAction', () => {
        const dto = { name: 'delete', description: 'Acción de eliminación' };

        it('debería crear una acción y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createAction).mockResolvedValue(singleAction);

            const result = await useCases.createAction(dto);

            expect(result).toEqual(singleAction);
            expect(repo.createAction).toHaveBeenCalledOnce();
            expect(repo.createAction).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createAction).mockRejectedValue(new Error('Nombre duplicado'));

            await expect(useCases.createAction(dto)).rejects.toThrow('Nombre duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateAction
    // ---------------------------------------------------------------
    describe('updateAction', () => {
        const updateDto = { name: 'read_updated', description: 'Lectura actualizada' };

        it('debería actualizar una acción y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<Action> = {
                ...singleAction,
                data: makeAction({ name: 'read_updated' }),
            };
            vi.mocked(repo.updateAction).mockResolvedValue(updatedResponse);

            const result = await useCases.updateAction(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateAction).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateAction).mockRejectedValue(new Error('Acción no encontrada'));

            await expect(useCases.updateAction(99, updateDto)).rejects.toThrow('Acción no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteAction
    // ---------------------------------------------------------------
    describe('deleteAction', () => {
        it('debería eliminar una acción y retornar confirmación', async () => {
            vi.mocked(repo.deleteAction).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteAction(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteAction).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando la acción no existe', async () => {
            vi.mocked(repo.deleteAction).mockResolvedValue(actionError);

            const result = await useCases.deleteAction(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteAction).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteAction(1)).rejects.toThrow('Network error');
        });
    });
});
