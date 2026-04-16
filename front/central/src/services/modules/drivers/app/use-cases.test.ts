import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DriverUseCases } from './use-cases';
import { IDriverRepository } from '../domain/ports';
import {
    DriverInfo,
    DriversListResponse,
    DeleteDriverResponse,
    CreateDriverDTO,
    UpdateDriverDTO,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeDriver = (overrides: Partial<DriverInfo> = {}): DriverInfo => ({
    id: 1,
    business_id: 1,
    first_name: 'Juan',
    last_name: 'Perez',
    email: 'juan.perez@example.com',
    phone: '+573001234567',
    identification: '1234567890',
    status: 'active',
    photo_url: '',
    license_type: 'B1',
    license_expiry: '2027-12-31',
    warehouse_id: 1,
    notes: null,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeDriversListResponse = (overrides: Partial<DriversListResponse> = {}): DriversListResponse => ({
    data: [makeDriver()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
    ...overrides,
});

const deleteSuccess: DeleteDriverResponse = { message: 'Conductor eliminado exitosamente' };
const deleteError: DeleteDriverResponse = { error: 'Conductor no encontrado' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IDriverRepository {
    return {
        getDrivers: vi.fn(),
        getDriverById: vi.fn(),
        createDriver: vi.fn(),
        updateDriver: vi.fn(),
        deleteDriver: vi.fn(),
    } as unknown as IDriverRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('DriverUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: DriverUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new DriverUseCases(repo as unknown as IDriverRepository);
    });

    // ---------------------------------------------------------------
    // getDrivers
    // ---------------------------------------------------------------
    describe('getDrivers', () => {
        it('debería retornar la lista paginada de conductores cuando el repositorio tiene éxito', async () => {
            const response = makeDriversListResponse();
            vi.mocked(repo.getDrivers).mockResolvedValue(response);

            const result = await useCases.getDrivers({ page: 1, page_size: 10 });

            expect(result).toEqual(response);
            expect(repo.getDrivers).toHaveBeenCalledOnce();
            expect(repo.getDrivers).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            const response = makeDriversListResponse();
            vi.mocked(repo.getDrivers).mockResolvedValue(response);

            await useCases.getDrivers();

            expect(repo.getDrivers).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getDrivers).mockRejectedValue(expectedError);

            await expect(useCases.getDrivers()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getDriverById
    // ---------------------------------------------------------------
    describe('getDriverById', () => {
        it('debería retornar un conductor por su ID', async () => {
            const driver = makeDriver();
            vi.mocked(repo.getDriverById).mockResolvedValue(driver);

            const result = await useCases.getDriverById(1);

            expect(result).toEqual(driver);
            expect(repo.getDriverById).toHaveBeenCalledOnce();
            expect(repo.getDriverById).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const driver = makeDriver();
            vi.mocked(repo.getDriverById).mockResolvedValue(driver);

            await useCases.getDriverById(1, 5);

            expect(repo.getDriverById).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando el conductor no existe', async () => {
            const expectedError = new Error('Conductor no encontrado');
            vi.mocked(repo.getDriverById).mockRejectedValue(expectedError);

            await expect(useCases.getDriverById(999)).rejects.toThrow('Conductor no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createDriver
    // ---------------------------------------------------------------
    describe('createDriver', () => {
        const dto: CreateDriverDTO = {
            first_name: 'Carlos',
            last_name: 'Lopez',
            phone: '+573009876543',
            identification: '9876543210',
            license_type: 'C1',
        };

        it('debería crear un conductor y retornar la respuesta del repositorio', async () => {
            const createdDriver = makeDriver({ first_name: 'Carlos', last_name: 'Lopez' });
            vi.mocked(repo.createDriver).mockResolvedValue(createdDriver);

            const result = await useCases.createDriver(dto);

            expect(result).toEqual(createdDriver);
            expect(repo.createDriver).toHaveBeenCalledOnce();
            expect(repo.createDriver).toHaveBeenCalledWith(dto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const createdDriver = makeDriver();
            vi.mocked(repo.createDriver).mockResolvedValue(createdDriver);

            await useCases.createDriver(dto, 3);

            expect(repo.createDriver).toHaveBeenCalledWith(dto, 3);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            const expectedError = new Error('Identificación duplicada');
            vi.mocked(repo.createDriver).mockRejectedValue(expectedError);

            await expect(useCases.createDriver(dto)).rejects.toThrow('Identificación duplicada');
        });
    });

    // ---------------------------------------------------------------
    // updateDriver
    // ---------------------------------------------------------------
    describe('updateDriver', () => {
        const updateDto: UpdateDriverDTO = { first_name: 'Carlos Actualizado', phone: '+573001111111' };

        it('debería actualizar un conductor y retornar la respuesta del repositorio', async () => {
            const updatedDriver = makeDriver({ first_name: 'Carlos Actualizado' });
            vi.mocked(repo.updateDriver).mockResolvedValue(updatedDriver);

            const result = await useCases.updateDriver(1, updateDto);

            expect(result).toEqual(updatedDriver);
            expect(repo.updateDriver).toHaveBeenCalledOnce();
            expect(repo.updateDriver).toHaveBeenCalledWith(1, updateDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const updatedDriver = makeDriver();
            vi.mocked(repo.updateDriver).mockResolvedValue(updatedDriver);

            await useCases.updateDriver(1, updateDto, 5);

            expect(repo.updateDriver).toHaveBeenCalledWith(1, updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            const expectedError = new Error('Conductor no encontrado');
            vi.mocked(repo.updateDriver).mockRejectedValue(expectedError);

            await expect(useCases.updateDriver(99, updateDto)).rejects.toThrow('Conductor no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteDriver
    // ---------------------------------------------------------------
    describe('deleteDriver', () => {
        it('debería eliminar un conductor y retornar confirmación', async () => {
            vi.mocked(repo.deleteDriver).mockResolvedValue(deleteSuccess);

            const result = await useCases.deleteDriver(1);

            expect(result).toEqual(deleteSuccess);
            expect(repo.deleteDriver).toHaveBeenCalledWith(1, undefined);
        });

        it('debería retornar respuesta de error cuando el conductor no existe', async () => {
            vi.mocked(repo.deleteDriver).mockResolvedValue(deleteError);

            const result = await useCases.deleteDriver(999);

            expect(result.error).toBeDefined();
            expect(result.message).toBeUndefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteDriver).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteDriver(1)).rejects.toThrow('Network error');
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteDriver).mockResolvedValue(deleteSuccess);

            await useCases.deleteDriver(1, 5);

            expect(repo.deleteDriver).toHaveBeenCalledWith(1, 5);
        });
    });
});
