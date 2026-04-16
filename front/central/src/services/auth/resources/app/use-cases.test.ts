import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ResourceUseCases } from './use-cases';
import { IResourceRepository } from '../domain/ports';
import {
    Resource,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeResource = (overrides: Partial<Resource> = {}): Resource => ({
    id: 1,
    name: 'orders',
    description: 'Gestión de pedidos',
    business_type_id: 1,
    business_type_name: 'Restaurante',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const paginatedResources: PaginatedResponse<Resource> = {
    success: true,
    message: 'OK',
    data: {
        resources: [makeResource(), makeResource({ id: 2, name: 'products' })],
        total: 2,
        page: 1,
        page_size: 10,
        total_pages: 1,
    },
};

const singleResource: SingleResponse<Resource> = {
    success: true,
    message: 'OK',
    data: makeResource(),
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IResourceRepository {
    return {
        getResources: vi.fn(),
        getResourceById: vi.fn(),
        createResource: vi.fn(),
        updateResource: vi.fn(),
        deleteResource: vi.fn(),
    } as unknown as IResourceRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('ResourceUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: ResourceUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new ResourceUseCases(repo as unknown as IResourceRepository);
    });

    // ---------------------------------------------------------------
    // getResources
    // ---------------------------------------------------------------
    describe('getResources', () => {
        it('debería retornar la lista paginada de recursos cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getResources).mockResolvedValue(paginatedResources);

            const result = await useCases.getResources({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedResources);
            expect(repo.getResources).toHaveBeenCalledOnce();
            expect(repo.getResources).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getResources).mockResolvedValue(paginatedResources);

            await useCases.getResources();

            expect(repo.getResources).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getResources).mockRejectedValue(new Error('Fallo de base de datos'));

            await expect(useCases.getResources()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getResourceById
    // ---------------------------------------------------------------
    describe('getResourceById', () => {
        it('debería retornar un recurso por ID', async () => {
            vi.mocked(repo.getResourceById).mockResolvedValue(singleResource);

            const result = await useCases.getResourceById(1);

            expect(result).toEqual(singleResource);
            expect(repo.getResourceById).toHaveBeenCalledOnce();
            expect(repo.getResourceById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando el recurso no existe', async () => {
            vi.mocked(repo.getResourceById).mockRejectedValue(new Error('Recurso no encontrado'));

            await expect(useCases.getResourceById(999)).rejects.toThrow('Recurso no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createResource
    // ---------------------------------------------------------------
    describe('createResource', () => {
        const dto = { name: 'shipments', description: 'Gestión de envíos', business_type_id: 1 };

        it('debería crear un recurso y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createResource).mockResolvedValue(singleResource);

            const result = await useCases.createResource(dto);

            expect(result).toEqual(singleResource);
            expect(repo.createResource).toHaveBeenCalledOnce();
            expect(repo.createResource).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createResource).mockRejectedValue(new Error('Nombre duplicado'));

            await expect(useCases.createResource(dto)).rejects.toThrow('Nombre duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateResource
    // ---------------------------------------------------------------
    describe('updateResource', () => {
        const updateDto = { name: 'orders_updated', description: 'Pedidos actualizado' };

        it('debería actualizar un recurso y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<Resource> = {
                ...singleResource,
                data: makeResource({ name: 'orders_updated' }),
            };
            vi.mocked(repo.updateResource).mockResolvedValue(updatedResponse);

            const result = await useCases.updateResource(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateResource).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateResource).mockRejectedValue(new Error('Recurso no encontrado'));

            await expect(useCases.updateResource(99, updateDto)).rejects.toThrow('Recurso no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteResource
    // ---------------------------------------------------------------
    describe('deleteResource', () => {
        it('debería eliminar un recurso y retornar confirmación', async () => {
            vi.mocked(repo.deleteResource).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteResource(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteResource).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el recurso no existe', async () => {
            vi.mocked(repo.deleteResource).mockResolvedValue(actionError);

            const result = await useCases.deleteResource(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteResource).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteResource(1)).rejects.toThrow('Network error');
        });
    });
});
