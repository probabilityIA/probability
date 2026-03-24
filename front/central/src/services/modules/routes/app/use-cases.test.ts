import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RouteUseCases } from './use-cases';
import { IRouteRepository } from '../domain/ports';
import {
    RouteInfo,
    RouteDetail,
    RouteStopInfo,
    RoutesListResponse,
    DeleteRouteResponse,
    DriverOption,
    VehicleOption,
    AssignableOrder,
    CreateRouteDTO,
    UpdateRouteDTO,
    AddStopDTO,
    UpdateStopDTO,
    UpdateStopStatusDTO,
    ReorderStopsDTO,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeRouteInfo = (overrides: Partial<RouteInfo> = {}): RouteInfo => ({
    id: 1,
    business_id: 1,
    driver_id: 10,
    driver_name: 'Carlos Lopez',
    vehicle_id: 5,
    vehicle_plate: 'ABC123',
    status: 'pending',
    date: '2026-03-24',
    start_time: null,
    end_time: null,
    origin_address: 'Calle 100 #15-20',
    total_stops: 3,
    completed_stops: 0,
    failed_stops: 0,
    notes: null,
    created_at: '2026-03-24T00:00:00Z',
    updated_at: '2026-03-24T00:00:00Z',
    ...overrides,
});

const makeStopInfo = (overrides: Partial<RouteStopInfo> = {}): RouteStopInfo => ({
    id: 1,
    route_id: 1,
    order_id: 'ORD-001',
    sequence: 1,
    status: 'pending',
    address: 'Calle 50 #10-30',
    city: 'Bogota',
    lat: 4.65,
    lng: -74.08,
    customer_name: 'Maria Garcia',
    customer_phone: '3009876543',
    estimated_arrival: null,
    actual_arrival: null,
    actual_departure: null,
    delivery_notes: null,
    failure_reason: null,
    created_at: '2026-03-24T00:00:00Z',
    updated_at: '2026-03-24T00:00:00Z',
    ...overrides,
});

const makeRouteDetail = (overrides: Partial<RouteDetail> = {}): RouteDetail => ({
    ...makeRouteInfo(),
    actual_start_time: null,
    actual_end_time: null,
    origin_warehouse_id: null,
    origin_lat: null,
    origin_lng: null,
    total_distance_km: null,
    total_duration_min: null,
    stops: [makeStopInfo()],
    ...overrides,
});

const makeDriverOption = (overrides: Partial<DriverOption> = {}): DriverOption => ({
    id: 10,
    first_name: 'Carlos',
    last_name: 'Lopez',
    phone: '3001112233',
    identification: '1234567890',
    status: 'available',
    license_type: 'B1',
    ...overrides,
});

const makeVehicleOption = (overrides: Partial<VehicleOption> = {}): VehicleOption => ({
    id: 5,
    type: 'van',
    license_plate: 'ABC123',
    brand: 'Chevrolet',
    vehicle_model: 'N300',
    status: 'available',
    ...overrides,
});

const makeAssignableOrder = (overrides: Partial<AssignableOrder> = {}): AssignableOrder => ({
    id: 'ORD-001',
    order_number: '10001',
    customer_name: 'Maria Garcia',
    customer_phone: '3009876543',
    address: 'Calle 50 #10-30',
    city: 'Bogota',
    lat: null,
    lng: null,
    total_amount: 150000,
    item_count: 3,
    created_at: '2026-03-24T00:00:00Z',
    ...overrides,
});

const routesListResponse: RoutesListResponse = {
    data: [makeRouteInfo()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const deleteSuccess: DeleteRouteResponse = { message: 'Eliminado exitosamente' };
const deleteError: DeleteRouteResponse = { error: 'Ruta no encontrada' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IRouteRepository {
    return {
        // Route CRUD
        getRoutes: vi.fn(),
        getRouteById: vi.fn(),
        createRoute: vi.fn(),
        updateRoute: vi.fn(),
        deleteRoute: vi.fn(),
        // Route lifecycle
        startRoute: vi.fn(),
        completeRoute: vi.fn(),
        // Stop management
        addStop: vi.fn(),
        updateStop: vi.fn(),
        deleteStop: vi.fn(),
        updateStopStatus: vi.fn(),
        reorderStops: vi.fn(),
        // Form options
        getAvailableDrivers: vi.fn(),
        getAvailableVehicles: vi.fn(),
        getAssignableOrders: vi.fn(),
    } as unknown as IRouteRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('RouteUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: RouteUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new RouteUseCases(repo as unknown as IRouteRepository);
    });

    // ---------------------------------------------------------------
    // getRoutes
    // ---------------------------------------------------------------
    describe('getRoutes', () => {
        it('debería retornar la lista paginada de rutas cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getRoutes).mockResolvedValue(routesListResponse);

            const result = await useCases.getRoutes({ page: 1, page_size: 10 });

            expect(result).toEqual(routesListResponse);
            expect(repo.getRoutes).toHaveBeenCalledOnce();
            expect(repo.getRoutes).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getRoutes).mockResolvedValue(routesListResponse);

            await useCases.getRoutes();

            expect(repo.getRoutes).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getRoutes).mockRejectedValue(expectedError);

            await expect(useCases.getRoutes()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getRouteById
    // ---------------------------------------------------------------
    describe('getRouteById', () => {
        it('debería retornar el detalle de una ruta por su ID', async () => {
            const detail = makeRouteDetail();
            vi.mocked(repo.getRouteById).mockResolvedValue(detail);

            const result = await useCases.getRouteById(1);

            expect(result).toEqual(detail);
            expect(repo.getRouteById).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const detail = makeRouteDetail();
            vi.mocked(repo.getRouteById).mockResolvedValue(detail);

            await useCases.getRouteById(1, 5);

            expect(repo.getRouteById).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando la ruta no existe', async () => {
            vi.mocked(repo.getRouteById).mockRejectedValue(new Error('Ruta no encontrada'));

            await expect(useCases.getRouteById(999)).rejects.toThrow('Ruta no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // createRoute
    // ---------------------------------------------------------------
    describe('createRoute', () => {
        const dto: CreateRouteDTO = {
            date: '2026-03-25',
            driver_id: 10,
            stops: [{ address: 'Calle 50 #10-30', customer_name: 'Maria Garcia' }],
        };

        it('debería crear una ruta y retornar la respuesta del repositorio', async () => {
            const created = makeRouteInfo({ id: 2 });
            vi.mocked(repo.createRoute).mockResolvedValue(created);

            const result = await useCases.createRoute(dto);

            expect(result).toEqual(created);
            expect(repo.createRoute).toHaveBeenCalledOnce();
            expect(repo.createRoute).toHaveBeenCalledWith(dto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const created = makeRouteInfo({ id: 2 });
            vi.mocked(repo.createRoute).mockResolvedValue(created);

            await useCases.createRoute(dto, 5);

            expect(repo.createRoute).toHaveBeenCalledWith(dto, 5);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createRoute).mockRejectedValue(new Error('Fecha inválida'));

            await expect(useCases.createRoute(dto)).rejects.toThrow('Fecha inválida');
        });
    });

    // ---------------------------------------------------------------
    // updateRoute
    // ---------------------------------------------------------------
    describe('updateRoute', () => {
        const updateDto: UpdateRouteDTO = { notes: 'Ruta actualizada' };

        it('debería actualizar una ruta y retornar la respuesta del repositorio', async () => {
            const updated = makeRouteInfo({ notes: 'Ruta actualizada' });
            vi.mocked(repo.updateRoute).mockResolvedValue(updated);

            const result = await useCases.updateRoute(1, updateDto);

            expect(result).toEqual(updated);
            expect(repo.updateRoute).toHaveBeenCalledWith(1, updateDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const updated = makeRouteInfo();
            vi.mocked(repo.updateRoute).mockResolvedValue(updated);

            await useCases.updateRoute(1, updateDto, 5);

            expect(repo.updateRoute).toHaveBeenCalledWith(1, updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateRoute).mockRejectedValue(new Error('Ruta no encontrada'));

            await expect(useCases.updateRoute(99, updateDto)).rejects.toThrow('Ruta no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteRoute
    // ---------------------------------------------------------------
    describe('deleteRoute', () => {
        it('debería eliminar una ruta y retornar confirmación', async () => {
            vi.mocked(repo.deleteRoute).mockResolvedValue(deleteSuccess);

            const result = await useCases.deleteRoute(1);

            expect(result).toEqual(deleteSuccess);
            expect(repo.deleteRoute).toHaveBeenCalledWith(1, undefined);
        });

        it('debería retornar respuesta de error cuando la ruta no existe', async () => {
            vi.mocked(repo.deleteRoute).mockResolvedValue(deleteError);

            const result = await useCases.deleteRoute(999);

            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteRoute).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteRoute(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // startRoute
    // ---------------------------------------------------------------
    describe('startRoute', () => {
        it('debería iniciar una ruta y retornar el detalle actualizado', async () => {
            const started = makeRouteDetail({ status: 'in_progress' });
            vi.mocked(repo.startRoute).mockResolvedValue(started);

            const result = await useCases.startRoute(1);

            expect(result).toEqual(started);
            expect(repo.startRoute).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const started = makeRouteDetail({ status: 'in_progress' });
            vi.mocked(repo.startRoute).mockResolvedValue(started);

            await useCases.startRoute(1, 5);

            expect(repo.startRoute).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando la ruta no puede iniciarse', async () => {
            vi.mocked(repo.startRoute).mockRejectedValue(new Error('Ruta ya iniciada'));

            await expect(useCases.startRoute(1)).rejects.toThrow('Ruta ya iniciada');
        });
    });

    // ---------------------------------------------------------------
    // completeRoute
    // ---------------------------------------------------------------
    describe('completeRoute', () => {
        it('debería completar una ruta y retornar el detalle actualizado', async () => {
            const completed = makeRouteDetail({ status: 'completed' });
            vi.mocked(repo.completeRoute).mockResolvedValue(completed);

            const result = await useCases.completeRoute(1);

            expect(result).toEqual(completed);
            expect(repo.completeRoute).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const completed = makeRouteDetail({ status: 'completed' });
            vi.mocked(repo.completeRoute).mockResolvedValue(completed);

            await useCases.completeRoute(1, 5);

            expect(repo.completeRoute).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando la ruta no puede completarse', async () => {
            vi.mocked(repo.completeRoute).mockRejectedValue(new Error('Paradas pendientes'));

            await expect(useCases.completeRoute(1)).rejects.toThrow('Paradas pendientes');
        });
    });

    // ---------------------------------------------------------------
    // addStop
    // ---------------------------------------------------------------
    describe('addStop', () => {
        const dto: AddStopDTO = { address: 'Calle 70 #20-10', customer_name: 'Pedro Martinez' };

        it('debería agregar una parada y retornar la respuesta del repositorio', async () => {
            const created = makeStopInfo({ id: 2, address: 'Calle 70 #20-10', customer_name: 'Pedro Martinez' });
            vi.mocked(repo.addStop).mockResolvedValue(created);

            const result = await useCases.addStop(1, dto);

            expect(result).toEqual(created);
            expect(repo.addStop).toHaveBeenCalledOnce();
            expect(repo.addStop).toHaveBeenCalledWith(1, dto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const created = makeStopInfo({ id: 2 });
            vi.mocked(repo.addStop).mockResolvedValue(created);

            await useCases.addStop(1, dto, 5);

            expect(repo.addStop).toHaveBeenCalledWith(1, dto, 5);
        });

        it('debería propagar el error cuando la adición falla', async () => {
            vi.mocked(repo.addStop).mockRejectedValue(new Error('Ruta no encontrada'));

            await expect(useCases.addStop(1, dto)).rejects.toThrow('Ruta no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // updateStop
    // ---------------------------------------------------------------
    describe('updateStop', () => {
        const updateDto: UpdateStopDTO = { address: 'Calle 80 #25-15' };

        it('debería actualizar una parada y retornar la respuesta del repositorio', async () => {
            const updated = makeStopInfo({ address: 'Calle 80 #25-15' });
            vi.mocked(repo.updateStop).mockResolvedValue(updated);

            const result = await useCases.updateStop(1, 1, updateDto);

            expect(result).toEqual(updated);
            expect(repo.updateStop).toHaveBeenCalledWith(1, 1, updateDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const updated = makeStopInfo();
            vi.mocked(repo.updateStop).mockResolvedValue(updated);

            await useCases.updateStop(1, 1, updateDto, 5);

            expect(repo.updateStop).toHaveBeenCalledWith(1, 1, updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateStop).mockRejectedValue(new Error('Parada no encontrada'));

            await expect(useCases.updateStop(1, 99, updateDto)).rejects.toThrow('Parada no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteStop
    // ---------------------------------------------------------------
    describe('deleteStop', () => {
        it('debería eliminar una parada y retornar confirmación', async () => {
            vi.mocked(repo.deleteStop).mockResolvedValue(deleteSuccess);

            const result = await useCases.deleteStop(1, 1);

            expect(result).toEqual(deleteSuccess);
            expect(repo.deleteStop).toHaveBeenCalledWith(1, 1, undefined);
        });

        it('debería retornar respuesta de error cuando la parada no existe', async () => {
            vi.mocked(repo.deleteStop).mockResolvedValue(deleteError);

            const result = await useCases.deleteStop(1, 999);

            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteStop).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteStop(1, 1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // updateStopStatus
    // ---------------------------------------------------------------
    describe('updateStopStatus', () => {
        const statusDto: UpdateStopStatusDTO = { status: 'completed' };

        it('debería actualizar el estado de una parada y retornar la respuesta', async () => {
            const updated = makeStopInfo({ status: 'completed' });
            vi.mocked(repo.updateStopStatus).mockResolvedValue(updated);

            const result = await useCases.updateStopStatus(1, 1, statusDto);

            expect(result).toEqual(updated);
            expect(repo.updateStopStatus).toHaveBeenCalledWith(1, 1, statusDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const updated = makeStopInfo({ status: 'completed' });
            vi.mocked(repo.updateStopStatus).mockResolvedValue(updated);

            await useCases.updateStopStatus(1, 1, statusDto, 5);

            expect(repo.updateStopStatus).toHaveBeenCalledWith(1, 1, statusDto, 5);
        });

        it('debería propagar el error cuando la actualización de estado falla', async () => {
            vi.mocked(repo.updateStopStatus).mockRejectedValue(new Error('Estado inválido'));

            await expect(useCases.updateStopStatus(1, 1, statusDto)).rejects.toThrow('Estado inválido');
        });
    });

    // ---------------------------------------------------------------
    // reorderStops
    // ---------------------------------------------------------------
    describe('reorderStops', () => {
        const reorderDto: ReorderStopsDTO = { stop_ids: [3, 1, 2] };

        it('debería reordenar las paradas y retornar el detalle actualizado', async () => {
            const detail = makeRouteDetail();
            vi.mocked(repo.reorderStops).mockResolvedValue(detail);

            const result = await useCases.reorderStops(1, reorderDto);

            expect(result).toEqual(detail);
            expect(repo.reorderStops).toHaveBeenCalledWith(1, reorderDto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const detail = makeRouteDetail();
            vi.mocked(repo.reorderStops).mockResolvedValue(detail);

            await useCases.reorderStops(1, reorderDto, 5);

            expect(repo.reorderStops).toHaveBeenCalledWith(1, reorderDto, 5);
        });

        it('debería propagar el error cuando el reordenamiento falla', async () => {
            vi.mocked(repo.reorderStops).mockRejectedValue(new Error('IDs inválidos'));

            await expect(useCases.reorderStops(1, reorderDto)).rejects.toThrow('IDs inválidos');
        });
    });

    // ---------------------------------------------------------------
    // getAvailableDrivers
    // ---------------------------------------------------------------
    describe('getAvailableDrivers', () => {
        it('debería retornar la lista de conductores disponibles', async () => {
            const drivers = [makeDriverOption()];
            vi.mocked(repo.getAvailableDrivers).mockResolvedValue(drivers);

            const result = await useCases.getAvailableDrivers();

            expect(result).toEqual(drivers);
            expect(repo.getAvailableDrivers).toHaveBeenCalledWith(undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.getAvailableDrivers).mockResolvedValue([]);

            await useCases.getAvailableDrivers(5);

            expect(repo.getAvailableDrivers).toHaveBeenCalledWith(5);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getAvailableDrivers).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getAvailableDrivers()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // getAvailableVehicles
    // ---------------------------------------------------------------
    describe('getAvailableVehicles', () => {
        it('debería retornar la lista de vehículos disponibles', async () => {
            const vehicles = [makeVehicleOption()];
            vi.mocked(repo.getAvailableVehicles).mockResolvedValue(vehicles);

            const result = await useCases.getAvailableVehicles();

            expect(result).toEqual(vehicles);
            expect(repo.getAvailableVehicles).toHaveBeenCalledWith(undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.getAvailableVehicles).mockResolvedValue([]);

            await useCases.getAvailableVehicles(5);

            expect(repo.getAvailableVehicles).toHaveBeenCalledWith(5);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getAvailableVehicles).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getAvailableVehicles()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // getAssignableOrders
    // ---------------------------------------------------------------
    describe('getAssignableOrders', () => {
        it('debería retornar la lista de órdenes asignables', async () => {
            const orders = [makeAssignableOrder()];
            vi.mocked(repo.getAssignableOrders).mockResolvedValue(orders);

            const result = await useCases.getAssignableOrders();

            expect(result).toEqual(orders);
            expect(repo.getAssignableOrders).toHaveBeenCalledWith(undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            vi.mocked(repo.getAssignableOrders).mockResolvedValue([]);

            await useCases.getAssignableOrders(5);

            expect(repo.getAssignableOrders).toHaveBeenCalledWith(5);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getAssignableOrders).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getAssignableOrders()).rejects.toThrow('Servicio no disponible');
        });
    });
});
