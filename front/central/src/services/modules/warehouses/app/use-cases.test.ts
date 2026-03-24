import { describe, it, expect, vi, beforeEach } from 'vitest';
import { WarehouseUseCases } from './use-cases';
import { IWarehouseRepository } from '../domain/ports';
import {
    Warehouse,
    WarehouseDetail,
    WarehouseLocation,
    WarehousesListResponse,
    CreateWarehouseDTO,
    UpdateWarehouseDTO,
    CreateLocationDTO,
    UpdateLocationDTO,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeWarehouse = (overrides: Partial<Warehouse> = {}): Warehouse => ({
    id: 1,
    business_id: 1,
    name: 'Bodega Principal',
    code: 'BOD-001',
    address: 'Calle 100 #15-20',
    city: 'Bogota',
    state: 'Cundinamarca',
    country: 'CO',
    zip_code: '110111',
    phone: '3001234567',
    contact_name: 'Juan Perez',
    contact_email: 'juan@test.com',
    is_active: true,
    is_default: true,
    is_fulfillment: false,
    company: 'Test Company',
    first_name: 'Juan',
    last_name: 'Perez',
    email: 'juan@test.com',
    suburb: 'Centro',
    city_dane_code: '11001',
    postal_code: '110111',
    street: 'Calle 100',
    latitude: 4.6097,
    longitude: -74.0817,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeLocation = (overrides: Partial<WarehouseLocation> = {}): WarehouseLocation => ({
    id: 1,
    warehouse_id: 1,
    name: 'Pasillo A',
    code: 'PA-001',
    type: 'shelf',
    is_active: true,
    is_fulfillment: false,
    capacity: 100,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeWarehouseDetail = (overrides: Partial<WarehouseDetail> = {}): WarehouseDetail => ({
    ...makeWarehouse(),
    locations: [makeLocation()],
    ...overrides,
});

const warehousesListResponse: WarehousesListResponse = {
    data: [makeWarehouse()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IWarehouseRepository {
    return {
        getWarehouses: vi.fn(),
        getWarehouseById: vi.fn(),
        createWarehouse: vi.fn(),
        updateWarehouse: vi.fn(),
        deleteWarehouse: vi.fn(),
        getLocations: vi.fn(),
        createLocation: vi.fn(),
        updateLocation: vi.fn(),
        deleteLocation: vi.fn(),
    } as unknown as IWarehouseRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('WarehouseUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: WarehouseUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new WarehouseUseCases(repo as unknown as IWarehouseRepository);
    });

    // ---------------------------------------------------------------
    // getWarehouses
    // ---------------------------------------------------------------
    describe('getWarehouses', () => {
        it('debería retornar la lista paginada de bodegas cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getWarehouses).mockResolvedValue(warehousesListResponse);

            const result = await useCases.getWarehouses({ page: 1, page_size: 10 });

            expect(result).toEqual(warehousesListResponse);
            expect(repo.getWarehouses).toHaveBeenCalledOnce();
            expect(repo.getWarehouses).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getWarehouses).mockResolvedValue(warehousesListResponse);

            await useCases.getWarehouses();

            expect(repo.getWarehouses).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getWarehouses).mockRejectedValue(expectedError);

            await expect(useCases.getWarehouses()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getWarehouseById
    // ---------------------------------------------------------------
    describe('getWarehouseById', () => {
        it('debería retornar el detalle de una bodega por su ID', async () => {
            const detail = makeWarehouseDetail();
            vi.mocked(repo.getWarehouseById).mockResolvedValue(detail);

            const result = await useCases.getWarehouseById(1);

            expect(result).toEqual(detail);
            expect(repo.getWarehouseById).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const detail = makeWarehouseDetail();
            vi.mocked(repo.getWarehouseById).mockResolvedValue(detail);

            await useCases.getWarehouseById(1, 5);

            expect(repo.getWarehouseById).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando la bodega no existe', async () => {
            vi.mocked(repo.getWarehouseById).mockRejectedValue(new Error('Bodega no encontrada'));

            await expect(useCases.getWarehouseById(999)).rejects.toThrow('Bodega no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // createWarehouse
    // ---------------------------------------------------------------
    describe('createWarehouse', () => {
        const dto: CreateWarehouseDTO = {
            name: 'Bodega Nueva',
            code: 'BOD-002',
        };

        it('debería crear una bodega y retornar la respuesta del repositorio', async () => {
            const created = makeWarehouse({ id: 2, name: 'Bodega Nueva', code: 'BOD-002' });
            vi.mocked(repo.createWarehouse).mockResolvedValue(created);

            const result = await useCases.createWarehouse(dto);

            expect(result).toEqual(created);
            expect(repo.createWarehouse).toHaveBeenCalledOnce();
            expect(repo.createWarehouse).toHaveBeenCalledWith(dto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const created = makeWarehouse({ id: 2 });
            vi.mocked(repo.createWarehouse).mockResolvedValue(created);

            await useCases.createWarehouse(dto, 5);

            expect(repo.createWarehouse).toHaveBeenCalledWith(dto, 5);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createWarehouse).mockRejectedValue(new Error('Código duplicado'));

            await expect(useCases.createWarehouse(dto)).rejects.toThrow('Código duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateWarehouse
    // ---------------------------------------------------------------
    describe('updateWarehouse', () => {
        const updateDto: UpdateWarehouseDTO = { name: 'Bodega Actualizada', code: 'BOD-001' };

        it('debería actualizar una bodega y retornar la respuesta del repositorio', async () => {
            const updated = makeWarehouse({ name: 'Bodega Actualizada' });
            vi.mocked(repo.updateWarehouse).mockResolvedValue(updated);

            const result = await useCases.updateWarehouse(1, updateDto);

            expect(result).toEqual(updated);
            expect(repo.updateWarehouse).toHaveBeenCalledWith(1, updateDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const updated = makeWarehouse();
            vi.mocked(repo.updateWarehouse).mockResolvedValue(updated);

            await useCases.updateWarehouse(1, updateDto, 5);

            expect(repo.updateWarehouse).toHaveBeenCalledWith(1, updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateWarehouse).mockRejectedValue(new Error('Bodega no encontrada'));

            await expect(useCases.updateWarehouse(99, updateDto)).rejects.toThrow('Bodega no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteWarehouse
    // ---------------------------------------------------------------
    describe('deleteWarehouse', () => {
        it('debería eliminar una bodega correctamente', async () => {
            vi.mocked(repo.deleteWarehouse).mockResolvedValue(undefined);

            await useCases.deleteWarehouse(1);

            expect(repo.deleteWarehouse).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteWarehouse).mockResolvedValue(undefined);

            await useCases.deleteWarehouse(1, 5);

            expect(repo.deleteWarehouse).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteWarehouse).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteWarehouse(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // getLocations
    // ---------------------------------------------------------------
    describe('getLocations', () => {
        it('debería retornar la lista de ubicaciones de una bodega', async () => {
            const locations = [makeLocation()];
            vi.mocked(repo.getLocations).mockResolvedValue(locations);

            const result = await useCases.getLocations(1);

            expect(result).toEqual(locations);
            expect(repo.getLocations).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.getLocations).mockResolvedValue([]);

            await useCases.getLocations(1, 5);

            expect(repo.getLocations).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getLocations).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getLocations(1)).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // createLocation
    // ---------------------------------------------------------------
    describe('createLocation', () => {
        const dto: CreateLocationDTO = { name: 'Pasillo B', code: 'PB-001' };

        it('debería crear una ubicación y retornar la respuesta del repositorio', async () => {
            const created = makeLocation({ id: 2, name: 'Pasillo B', code: 'PB-001' });
            vi.mocked(repo.createLocation).mockResolvedValue(created);

            const result = await useCases.createLocation(1, dto);

            expect(result).toEqual(created);
            expect(repo.createLocation).toHaveBeenCalledOnce();
            expect(repo.createLocation).toHaveBeenCalledWith(1, dto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const created = makeLocation({ id: 2 });
            vi.mocked(repo.createLocation).mockResolvedValue(created);

            await useCases.createLocation(1, dto, 5);

            expect(repo.createLocation).toHaveBeenCalledWith(1, dto, 5);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createLocation).mockRejectedValue(new Error('Código duplicado'));

            await expect(useCases.createLocation(1, dto)).rejects.toThrow('Código duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateLocation
    // ---------------------------------------------------------------
    describe('updateLocation', () => {
        const updateDto: UpdateLocationDTO = { name: 'Pasillo A Actualizado', code: 'PA-001' };

        it('debería actualizar una ubicación y retornar la respuesta del repositorio', async () => {
            const updated = makeLocation({ name: 'Pasillo A Actualizado' });
            vi.mocked(repo.updateLocation).mockResolvedValue(updated);

            const result = await useCases.updateLocation(1, 1, updateDto);

            expect(result).toEqual(updated);
            expect(repo.updateLocation).toHaveBeenCalledWith(1, 1, updateDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const updated = makeLocation();
            vi.mocked(repo.updateLocation).mockResolvedValue(updated);

            await useCases.updateLocation(1, 1, updateDto, 5);

            expect(repo.updateLocation).toHaveBeenCalledWith(1, 1, updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateLocation).mockRejectedValue(new Error('Ubicación no encontrada'));

            await expect(useCases.updateLocation(1, 99, updateDto)).rejects.toThrow('Ubicación no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteLocation
    // ---------------------------------------------------------------
    describe('deleteLocation', () => {
        it('debería eliminar una ubicación correctamente', async () => {
            vi.mocked(repo.deleteLocation).mockResolvedValue(undefined);

            await useCases.deleteLocation(1, 1);

            expect(repo.deleteLocation).toHaveBeenCalledWith(1, 1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteLocation).mockResolvedValue(undefined);

            await useCases.deleteLocation(1, 1, 5);

            expect(repo.deleteLocation).toHaveBeenCalledWith(1, 1, 5);
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteLocation).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteLocation(1, 1)).rejects.toThrow('Network error');
        });
    });
});
