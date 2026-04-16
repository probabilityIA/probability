import { describe, it, expect, vi, beforeEach } from 'vitest';
import { DashboardUseCases } from './use-cases';
import { IDashboardRepository } from '../domain/ports';
import { DashboardStatsResponse, DashboardStats } from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeDashboardStats = (overrides: Partial<DashboardStats> = {}): DashboardStats => ({
    total_orders: 150,
    orders_today: 12,
    orders_by_integration_type: [
        { integration_type: 'shopify', count: 80 },
        { integration_type: 'whatsapp', count: 70 },
    ],
    top_customers: [
        { customer_name: 'Juan Perez', customer_email: 'juan@test.com', order_count: 25 },
    ],
    orders_by_location: [
        { city: 'Bogota', state: 'Cundinamarca', order_count: 50 },
    ],
    top_drivers: [
        { driver_name: 'Carlos Lopez', order_count: 30 },
    ],
    drivers_by_location: [
        { driver_name: 'Carlos Lopez', city: 'Bogota', state: 'Cundinamarca', order_count: 30 },
    ],
    top_products: [
        { product_name: 'Camiseta', product_id: 'P001', sku: 'SKU001', order_count: 40, total_sold: 200 },
    ],
    products_by_category: [
        { category: 'Ropa', count: 60 },
    ],
    products_by_brand: [
        { brand: 'MarcaX', count: 45 },
    ],
    shipments_by_status: [
        { status: 'delivered', count: 100 },
    ],
    shipments_by_carrier: [
        { carrier: 'Servientrega', count: 80 },
    ],
    shipments_by_carrier_today: [
        { carrier: 'Servientrega', count: 5 },
    ],
    shipments_by_warehouse: [
        { warehouse_name: 'Bodega Principal', count: 90 },
    ],
    shipments_by_day_of_week: [
        { date: '2026-03-23', day_name: 'Lunes', count: 20 },
    ],
    orders_by_department: [
        { department: 'Cundinamarca', count: 50 },
    ],
    ...overrides,
});

const makeStatsResponse = (overrides: Partial<DashboardStatsResponse> = {}): DashboardStatsResponse => ({
    success: true,
    message: 'OK',
    data: makeDashboardStats(),
    ...overrides,
});

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IDashboardRepository {
    return {
        getStats: vi.fn(),
    } as unknown as IDashboardRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('DashboardUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: DashboardUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new DashboardUseCases(repo as unknown as IDashboardRepository);
    });

    // ---------------------------------------------------------------
    // getStats
    // ---------------------------------------------------------------
    describe('getStats', () => {
        it('debería retornar las estadísticas del dashboard cuando el repositorio tiene éxito', async () => {
            const response = makeStatsResponse();
            vi.mocked(repo.getStats).mockResolvedValue(response);

            const result = await useCases.getStats(1, 2);

            expect(result).toEqual(response);
            expect(repo.getStats).toHaveBeenCalledOnce();
            expect(repo.getStats).toHaveBeenCalledWith(1, 2, undefined, undefined, undefined);
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            const response = makeStatsResponse();
            vi.mocked(repo.getStats).mockResolvedValue(response);

            await useCases.getStats();

            expect(repo.getStats).toHaveBeenCalledWith(undefined, undefined, undefined, undefined, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const response = makeStatsResponse();
            vi.mocked(repo.getStats).mockResolvedValue(response);

            await useCases.getStats(5);

            expect(repo.getStats).toHaveBeenCalledWith(5, undefined, undefined, undefined, undefined);
        });

        it('debería pasar el integrationId cuando se proporciona', async () => {
            const response = makeStatsResponse();
            vi.mocked(repo.getStats).mockResolvedValue(response);

            await useCases.getStats(undefined, 3);

            expect(repo.getStats).toHaveBeenCalledWith(undefined, 3, undefined, undefined, undefined);
        });

        it('debería pasar el weekStartDate cuando se proporciona', async () => {
            const response = makeStatsResponse();
            vi.mocked(repo.getStats).mockResolvedValue(response);
            const weekStart = new Date('2026-03-16');

            await useCases.getStats(1, 2, weekStart);

            expect(repo.getStats).toHaveBeenCalledWith(1, 2, weekStart, undefined, undefined);
        });

        it('debería retornar estadísticas con orders_by_business para super admin', async () => {
            const statsWithBusiness = makeDashboardStats({
                orders_by_business: [
                    { business_id: 1, business_name: 'Negocio 1', order_count: 100 },
                    { business_id: 2, business_name: 'Negocio 2', order_count: 50 },
                ],
            });
            const response = makeStatsResponse({ data: statsWithBusiness });
            vi.mocked(repo.getStats).mockResolvedValue(response);

            const result = await useCases.getStats();

            expect(result.data.orders_by_business).toHaveLength(2);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Servicio no disponible');
            vi.mocked(repo.getStats).mockRejectedValue(expectedError);

            await expect(useCases.getStats()).rejects.toThrow('Servicio no disponible');
        });

        it('debería propagar errores de red', async () => {
            vi.mocked(repo.getStats).mockRejectedValue(new Error('Network error'));

            await expect(useCases.getStats(1, 2)).rejects.toThrow('Network error');
        });
    });
});
