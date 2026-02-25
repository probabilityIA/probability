import { describe, it, expect, vi, beforeEach } from 'vitest';
import { BusinessUseCases } from './use-cases';
import { IBusinessRepository } from '../domain/ports';
import {
    Business,
    BusinessType,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    BusinessConfiguredResources,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makePagination = () => ({
    current_page: 1,
    per_page: 10,
    total: 1,
    last_page: 1,
    has_next: false,
    has_prev: false,
});

const makeBusiness = (overrides: Partial<Business> = {}): Business => ({
    id: 1,
    name: 'Negocio de prueba',
    business_type_id: 1,
    is_active: true,
    ...overrides,
});

const makeBusinessType = (overrides: Partial<BusinessType> = {}): BusinessType => ({
    id: 1,
    name: 'Restaurante',
    code: 'restaurant',
    is_active: true,
    ...overrides,
});

const makeConfiguredResources = (overrides: Partial<BusinessConfiguredResources> = {}): BusinessConfiguredResources => ({
    business_id: 1,
    resources: [{ resource_id: 10, resource_name: 'Delivery', is_active: true }],
    total: 1,
    active: 1,
    inactive: 0,
    ...overrides,
});

const paginatedBusinesses: PaginatedResponse<Business> = {
    success: true,
    message: 'OK',
    data: [makeBusiness()],
    pagination: makePagination(),
};

const singleBusiness: SingleResponse<Business> = {
    success: true,
    message: 'OK',
    data: makeBusiness(),
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

/**
 * Crea un mock completo de IBusinessRepository usando vi.fn().
 * Incluye también los métodos activateBusiness / deactivateBusiness
 * que el use case expone aunque no estén aún en el port formal.
 */
function createMockRepository(): IBusinessRepository {
    return {
        // Business
        getBusinesses: vi.fn(),
        getBusinessById: vi.fn(),
        createBusiness: vi.fn(),
        updateBusiness: vi.fn(),
        deleteBusiness: vi.fn(),
        activateBusiness: vi.fn(),
        deactivateBusiness: vi.fn(),
        // Configured Resources
        getConfiguredResources: vi.fn(),
        getBusinessConfiguredResources: vi.fn(),
        activateResource: vi.fn(),
        deactivateResource: vi.fn(),
        // Business Types
        getBusinessTypes: vi.fn(),
        getBusinessTypeById: vi.fn(),
        createBusinessType: vi.fn(),
        updateBusinessType: vi.fn(),
        deleteBusinessType: vi.fn(),
    } as unknown as IBusinessRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('BusinessUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: BusinessUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new BusinessUseCases(repo as unknown as IBusinessRepository);
    });

    // ---------------------------------------------------------------
    // getBusinesses
    // ---------------------------------------------------------------
    describe('getBusinesses', () => {
        it('debería retornar la lista paginada de negocios cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getBusinesses).mockResolvedValue(paginatedBusinesses);

            const result = await useCases.getBusinesses({ page: 1, per_page: 10 });

            expect(result).toEqual(paginatedBusinesses);
            expect(repo.getBusinesses).toHaveBeenCalledOnce();
            expect(repo.getBusinesses).toHaveBeenCalledWith({ page: 1, per_page: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getBusinesses).mockResolvedValue(paginatedBusinesses);

            await useCases.getBusinesses();

            expect(repo.getBusinesses).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getBusinesses).mockRejectedValue(expectedError);

            await expect(useCases.getBusinesses()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // createBusiness
    // ---------------------------------------------------------------
    describe('createBusiness', () => {
        const dto = {
            name: 'Nuevo negocio',
            code: 'nuevo-neg',
            business_type_id: 1,
        };

        it('debería crear un negocio y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createBusiness).mockResolvedValue(singleBusiness);

            const result = await useCases.createBusiness(dto);

            expect(result).toEqual(singleBusiness);
            expect(repo.createBusiness).toHaveBeenCalledOnce();
            expect(repo.createBusiness).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            const expectedError = new Error('Nombre duplicado');
            vi.mocked(repo.createBusiness).mockRejectedValue(expectedError);

            await expect(useCases.createBusiness(dto)).rejects.toThrow('Nombre duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateBusiness
    // ---------------------------------------------------------------
    describe('updateBusiness', () => {
        const updateDto = { name: 'Negocio actualizado' };

        it('debería actualizar un negocio y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<Business> = {
                ...singleBusiness,
                data: makeBusiness({ name: 'Negocio actualizado' }),
            };
            vi.mocked(repo.updateBusiness).mockResolvedValue(updatedResponse);

            const result = await useCases.updateBusiness(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateBusiness).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            const expectedError = new Error('Negocio no encontrado');
            vi.mocked(repo.updateBusiness).mockRejectedValue(expectedError);

            await expect(useCases.updateBusiness(99, updateDto)).rejects.toThrow('Negocio no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteBusiness
    // ---------------------------------------------------------------
    describe('deleteBusiness', () => {
        it('debería eliminar un negocio y retornar confirmación', async () => {
            vi.mocked(repo.deleteBusiness).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteBusiness(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteBusiness).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el negocio no existe', async () => {
            vi.mocked(repo.deleteBusiness).mockResolvedValue(actionError);

            const result = await useCases.deleteBusiness(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteBusiness).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteBusiness(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // activateBusiness
    // ---------------------------------------------------------------
    describe('activateBusiness', () => {
        it('debería activar un negocio y retornar confirmación', async () => {
            vi.mocked(repo.activateBusiness).mockResolvedValue(actionSuccess);

            const result = await useCases.activateBusiness(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.activateBusiness).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando la activación falla', async () => {
            vi.mocked(repo.activateBusiness).mockRejectedValue(new Error('No autorizado'));

            await expect(useCases.activateBusiness(1)).rejects.toThrow('No autorizado');
        });
    });

    // ---------------------------------------------------------------
    // deactivateBusiness
    // ---------------------------------------------------------------
    describe('deactivateBusiness', () => {
        it('debería desactivar un negocio y retornar confirmación', async () => {
            vi.mocked(repo.deactivateBusiness).mockResolvedValue(actionSuccess);

            const result = await useCases.deactivateBusiness(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deactivateBusiness).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando la desactivación falla', async () => {
            vi.mocked(repo.deactivateBusiness).mockRejectedValue(new Error('Negocio ya inactivo'));

            await expect(useCases.deactivateBusiness(1)).rejects.toThrow('Negocio ya inactivo');
        });
    });

    // ---------------------------------------------------------------
    // getBusinessTypes
    // ---------------------------------------------------------------
    describe('getBusinessTypes', () => {
        it('debería retornar la lista de tipos de negocio', async () => {
            const paginatedTypes: PaginatedResponse<BusinessType> = {
                success: true,
                message: 'OK',
                data: [makeBusinessType()],
                pagination: makePagination(),
            };
            vi.mocked(repo.getBusinessTypes).mockResolvedValue(paginatedTypes);

            const result = await useCases.getBusinessTypes();

            expect(result).toEqual(paginatedTypes);
            expect(repo.getBusinessTypes).toHaveBeenCalledOnce();
        });

        it('debería propagar el error cuando la consulta de tipos falla', async () => {
            vi.mocked(repo.getBusinessTypes).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getBusinessTypes()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // activateResource / deactivateResource
    // ---------------------------------------------------------------
    describe('activateResource', () => {
        it('debería activar un recurso y retornar confirmación', async () => {
            vi.mocked(repo.activateResource).mockResolvedValue(actionSuccess);

            const result = await useCases.activateResource(10, 1);

            expect(result).toEqual(actionSuccess);
            expect(repo.activateResource).toHaveBeenCalledWith(10, 1);
        });

        it('debería activar un recurso sin businessId opcional', async () => {
            vi.mocked(repo.activateResource).mockResolvedValue(actionSuccess);

            await useCases.activateResource(10);

            expect(repo.activateResource).toHaveBeenCalledWith(10, undefined);
        });
    });

    describe('deactivateResource', () => {
        it('debería desactivar un recurso y retornar confirmación', async () => {
            vi.mocked(repo.deactivateResource).mockResolvedValue(actionSuccess);

            const result = await useCases.deactivateResource(10, 1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deactivateResource).toHaveBeenCalledWith(10, 1);
        });

        it('debería propagar el error cuando la desactivación del recurso falla', async () => {
            vi.mocked(repo.deactivateResource).mockRejectedValue(new Error('Recurso no encontrado'));

            await expect(useCases.deactivateResource(99)).rejects.toThrow('Recurso no encontrado');
        });
    });
});
