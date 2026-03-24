import { describe, it, expect, vi, beforeEach } from 'vitest';
import { OrderUseCases } from './use-cases';
import { IOrderRepository } from '../domain/ports';
import {
    Order,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    CreateOrderDTO,
    UpdateOrderDTO,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeOrder = (overrides: Partial<Order> = {}): Order => ({
    id: '1',
    created_at: '2026-03-01T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
    integration_id: 1,
    integration_type: 'shopify',
    platform: 'shopify',
    external_id: 'ext-001',
    order_number: 'ORD-001',
    internal_number: 'INT-001',
    subtotal: 50000,
    tax: 9500,
    discount: 0,
    shipping_cost: 5000,
    total_amount: 64500,
    currency: 'COP',
    customer_name: 'Juan Perez',
    customer_email: 'juan@test.com',
    customer_phone: '3001234567',
    customer_dni: '123456789',
    shipping_street: 'Calle 10 #5-20',
    shipping_city: 'Bogota',
    shipping_state: 'Cundinamarca',
    shipping_country: 'CO',
    shipping_postal_code: '110111',
    payment_method_id: 1,
    is_paid: true,
    warehouse_name: 'Bodega Principal',
    driver_name: '',
    is_last_mile: false,
    order_type_name: 'standard',
    status: 'pending',
    original_status: 'pending',
    user_name: 'admin',
    invoiceable: true,
    occurred_at: '2026-03-01T00:00:00Z',
    imported_at: '2026-03-01T00:00:00Z',
    ...overrides,
});

const paginatedOrders: PaginatedResponse<Order> = {
    success: true,
    message: 'OK',
    data: [makeOrder()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const singleOrder: SingleResponse<Order> = {
    success: true,
    message: 'OK',
    data: makeOrder(),
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IOrderRepository {
    return {
        getOrders: vi.fn(),
        getOrderById: vi.fn(),
        createOrder: vi.fn(),
        updateOrder: vi.fn(),
        deleteOrder: vi.fn(),
        getOrderRaw: vi.fn(),
        getAIRecommendation: vi.fn(),
        requestConfirmation: vi.fn(),
    } as unknown as IOrderRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('OrderUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: OrderUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new OrderUseCases(repo as unknown as IOrderRepository);
    });

    // ---------------------------------------------------------------
    // getOrders
    // ---------------------------------------------------------------
    describe('getOrders', () => {
        it('debería retornar la lista paginada de órdenes cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getOrders).mockResolvedValue(paginatedOrders);

            const result = await useCases.getOrders({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedOrders);
            expect(repo.getOrders).toHaveBeenCalledOnce();
            expect(repo.getOrders).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getOrders).mockResolvedValue(paginatedOrders);

            await useCases.getOrders();

            expect(repo.getOrders).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getOrders).mockRejectedValue(expectedError);

            await expect(useCases.getOrders()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getOrderById
    // ---------------------------------------------------------------
    describe('getOrderById', () => {
        it('debería retornar una orden por su ID', async () => {
            vi.mocked(repo.getOrderById).mockResolvedValue(singleOrder);

            const result = await useCases.getOrderById('1');

            expect(result).toEqual(singleOrder);
            expect(repo.getOrderById).toHaveBeenCalledWith('1');
        });

        it('debería propagar el error cuando la orden no existe', async () => {
            vi.mocked(repo.getOrderById).mockRejectedValue(new Error('Orden no encontrada'));

            await expect(useCases.getOrderById('999')).rejects.toThrow('Orden no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // createOrder
    // ---------------------------------------------------------------
    describe('createOrder', () => {
        const dto: CreateOrderDTO = {
            integration_id: 1,
            integration_type: 'shopify',
            platform: 'shopify',
            external_id: 'ext-002',
            subtotal: 30000,
            total_amount: 35000,
            payment_method_id: 1,
        };

        it('debería crear una orden y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createOrder).mockResolvedValue(singleOrder);

            const result = await useCases.createOrder(dto);

            expect(result).toEqual(singleOrder);
            expect(repo.createOrder).toHaveBeenCalledOnce();
            expect(repo.createOrder).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createOrder).mockRejectedValue(new Error('Datos inválidos'));

            await expect(useCases.createOrder(dto)).rejects.toThrow('Datos inválidos');
        });
    });

    // ---------------------------------------------------------------
    // updateOrder
    // ---------------------------------------------------------------
    describe('updateOrder', () => {
        const updateDto: UpdateOrderDTO = { customer_name: 'Nombre Actualizado' };

        it('debería actualizar una orden y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<Order> = {
                ...singleOrder,
                data: makeOrder({ customer_name: 'Nombre Actualizado' }),
            };
            vi.mocked(repo.updateOrder).mockResolvedValue(updatedResponse);

            const result = await useCases.updateOrder('1', updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateOrder).toHaveBeenCalledWith('1', updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateOrder).mockRejectedValue(new Error('Orden no encontrada'));

            await expect(useCases.updateOrder('999', updateDto)).rejects.toThrow('Orden no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteOrder
    // ---------------------------------------------------------------
    describe('deleteOrder', () => {
        it('debería eliminar una orden y retornar confirmación', async () => {
            vi.mocked(repo.deleteOrder).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteOrder('1');

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteOrder).toHaveBeenCalledWith('1');
        });

        it('debería retornar respuesta de error cuando la orden no existe', async () => {
            vi.mocked(repo.deleteOrder).mockResolvedValue(actionError);

            const result = await useCases.deleteOrder('999');

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteOrder).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteOrder('1')).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // getOrderRaw
    // ---------------------------------------------------------------
    describe('getOrderRaw', () => {
        it('debería retornar los datos crudos de una orden', async () => {
            const rawResponse = { success: true, message: 'OK', data: { raw: 'data' } };
            vi.mocked(repo.getOrderRaw).mockResolvedValue(rawResponse);

            const result = await useCases.getOrderRaw('1');

            expect(result).toEqual(rawResponse);
            expect(repo.getOrderRaw).toHaveBeenCalledWith('1');
        });

        it('debería propagar el error cuando falla la obtención de datos crudos', async () => {
            vi.mocked(repo.getOrderRaw).mockRejectedValue(new Error('No encontrado'));

            await expect(useCases.getOrderRaw('999')).rejects.toThrow('No encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getAIRecommendation
    // ---------------------------------------------------------------
    describe('getAIRecommendation', () => {
        it('debería retornar la recomendación de IA para origen y destino', async () => {
            const aiResponse = { recommendation: 'Enviar por Servientrega' };
            vi.mocked(repo.getAIRecommendation).mockResolvedValue(aiResponse);

            const result = await useCases.getAIRecommendation('Bogota', 'Medellin');

            expect(result).toEqual(aiResponse);
            expect(repo.getAIRecommendation).toHaveBeenCalledWith('Bogota', 'Medellin');
        });

        it('debería propagar el error cuando la IA falla', async () => {
            vi.mocked(repo.getAIRecommendation).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getAIRecommendation('Bogota', 'Medellin')).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // requestConfirmation
    // ---------------------------------------------------------------
    describe('requestConfirmation', () => {
        it('debería solicitar confirmación de una orden y retornar respuesta', async () => {
            vi.mocked(repo.requestConfirmation).mockResolvedValue(actionSuccess);

            const result = await useCases.requestConfirmation('1');

            expect(result).toEqual(actionSuccess);
            expect(repo.requestConfirmation).toHaveBeenCalledWith('1');
        });

        it('debería propagar el error cuando la solicitud de confirmación falla', async () => {
            vi.mocked(repo.requestConfirmation).mockRejectedValue(new Error('Orden ya confirmada'));

            await expect(useCases.requestConfirmation('1')).rejects.toThrow('Orden ya confirmada');
        });
    });
});
