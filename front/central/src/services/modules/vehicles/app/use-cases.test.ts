import { describe, it, expect, vi, beforeEach } from 'vitest';
import { VehicleUseCases } from './use-cases';
import { IVehicleRepository } from '../domain/ports';
import {
    VehicleInfo,
    VehiclesListResponse,
    DeleteVehicleResponse,
    CreateVehicleDTO,
    UpdateVehicleDTO,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeVehicle = (overrides: Partial<VehicleInfo> = {}): VehicleInfo => ({
    id: 1,
    business_id: 1,
    type: 'van',
    license_plate: 'ABC-123',
    brand: 'Toyota',
    model: 'HiAce',
    year: 2024,
    color: 'Blanco',
    status: 'active',
    weight_capacity_kg: 1500,
    volume_capacity_m3: 8.5,
    photo_url: '',
    insurance_expiry: '2027-06-30',
    registration_expiry: '2027-12-31',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeVehiclesListResponse = (overrides: Partial<VehiclesListResponse> = {}): VehiclesListResponse => ({
    data: [makeVehicle()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
    ...overrides,
});

const deleteSuccess: DeleteVehicleResponse = { message: 'Vehículo eliminado exitosamente' };
const deleteError: DeleteVehicleResponse = { error: 'Vehículo no encontrado' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IVehicleRepository {
    return {
        getVehicles: vi.fn(),
        getVehicleById: vi.fn(),
        createVehicle: vi.fn(),
        updateVehicle: vi.fn(),
        deleteVehicle: vi.fn(),
    } as unknown as IVehicleRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('VehicleUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: VehicleUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new VehicleUseCases(repo as unknown as IVehicleRepository);
    });

    // ---------------------------------------------------------------
    // getVehicles
    // ---------------------------------------------------------------
    describe('getVehicles', () => {
        it('debería retornar la lista paginada de vehículos cuando el repositorio tiene éxito', async () => {
            const response = makeVehiclesListResponse();
            vi.mocked(repo.getVehicles).mockResolvedValue(response);

            const result = await useCases.getVehicles({ page: 1, page_size: 10 });

            expect(result).toEqual(response);
            expect(repo.getVehicles).toHaveBeenCalledOnce();
            expect(repo.getVehicles).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            const response = makeVehiclesListResponse();
            vi.mocked(repo.getVehicles).mockResolvedValue(response);

            await useCases.getVehicles();

            expect(repo.getVehicles).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getVehicles).mockRejectedValue(expectedError);

            await expect(useCases.getVehicles()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getVehicleById
    // ---------------------------------------------------------------
    describe('getVehicleById', () => {
        it('debería retornar un vehículo por su ID', async () => {
            const vehicle = makeVehicle();
            vi.mocked(repo.getVehicleById).mockResolvedValue(vehicle);

            const result = await useCases.getVehicleById(1);

            expect(result).toEqual(vehicle);
            expect(repo.getVehicleById).toHaveBeenCalledOnce();
            expect(repo.getVehicleById).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const vehicle = makeVehicle();
            vi.mocked(repo.getVehicleById).mockResolvedValue(vehicle);

            await useCases.getVehicleById(1, 5);

            expect(repo.getVehicleById).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando el vehículo no existe', async () => {
            const expectedError = new Error('Vehículo no encontrado');
            vi.mocked(repo.getVehicleById).mockRejectedValue(expectedError);

            await expect(useCases.getVehicleById(999)).rejects.toThrow('Vehículo no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createVehicle
    // ---------------------------------------------------------------
    describe('createVehicle', () => {
        const dto: CreateVehicleDTO = {
            type: 'truck',
            license_plate: 'XYZ-789',
            brand: 'Chevrolet',
            model: 'NHR',
            year: 2025,
            color: 'Rojo',
        };

        it('debería crear un vehículo y retornar la respuesta del repositorio', async () => {
            const createdVehicle = makeVehicle({ type: 'truck', license_plate: 'XYZ-789' });
            vi.mocked(repo.createVehicle).mockResolvedValue(createdVehicle);

            const result = await useCases.createVehicle(dto);

            expect(result).toEqual(createdVehicle);
            expect(repo.createVehicle).toHaveBeenCalledOnce();
            expect(repo.createVehicle).toHaveBeenCalledWith(dto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const createdVehicle = makeVehicle();
            vi.mocked(repo.createVehicle).mockResolvedValue(createdVehicle);

            await useCases.createVehicle(dto, 3);

            expect(repo.createVehicle).toHaveBeenCalledWith(dto, 3);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            const expectedError = new Error('Placa duplicada');
            vi.mocked(repo.createVehicle).mockRejectedValue(expectedError);

            await expect(useCases.createVehicle(dto)).rejects.toThrow('Placa duplicada');
        });
    });

    // ---------------------------------------------------------------
    // updateVehicle
    // ---------------------------------------------------------------
    describe('updateVehicle', () => {
        const updateDto: UpdateVehicleDTO = { color: 'Azul', status: 'maintenance' };

        it('debería actualizar un vehículo y retornar la respuesta del repositorio', async () => {
            const updatedVehicle = makeVehicle({ color: 'Azul', status: 'maintenance' });
            vi.mocked(repo.updateVehicle).mockResolvedValue(updatedVehicle);

            const result = await useCases.updateVehicle(1, updateDto);

            expect(result).toEqual(updatedVehicle);
            expect(repo.updateVehicle).toHaveBeenCalledOnce();
            expect(repo.updateVehicle).toHaveBeenCalledWith(1, updateDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const updatedVehicle = makeVehicle();
            vi.mocked(repo.updateVehicle).mockResolvedValue(updatedVehicle);

            await useCases.updateVehicle(1, updateDto, 5);

            expect(repo.updateVehicle).toHaveBeenCalledWith(1, updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            const expectedError = new Error('Vehículo no encontrado');
            vi.mocked(repo.updateVehicle).mockRejectedValue(expectedError);

            await expect(useCases.updateVehicle(99, updateDto)).rejects.toThrow('Vehículo no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteVehicle
    // ---------------------------------------------------------------
    describe('deleteVehicle', () => {
        it('debería eliminar un vehículo y retornar confirmación', async () => {
            vi.mocked(repo.deleteVehicle).mockResolvedValue(deleteSuccess);

            const result = await useCases.deleteVehicle(1);

            expect(result).toEqual(deleteSuccess);
            expect(repo.deleteVehicle).toHaveBeenCalledWith(1, undefined);
        });

        it('debería retornar respuesta de error cuando el vehículo no existe', async () => {
            vi.mocked(repo.deleteVehicle).mockResolvedValue(deleteError);

            const result = await useCases.deleteVehicle(999);

            expect(result.error).toBeDefined();
            expect(result.message).toBeUndefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteVehicle).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteVehicle(1)).rejects.toThrow('Network error');
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteVehicle).mockResolvedValue(deleteSuccess);

            await useCases.deleteVehicle(1, 5);

            expect(repo.deleteVehicle).toHaveBeenCalledWith(1, 5);
        });
    });
});
