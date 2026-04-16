import { describe, it, expect, vi, beforeEach } from 'vitest';
import { OrderStatusMappingUseCases } from './use-cases';
import { IOrderStatusMappingRepository } from '../domain/ports';
import {
    OrderStatusMapping,
    OrderStatusInfo,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    GetOrderStatusMappingsParams,
    CreateOrderStatusMappingDTO,
    UpdateOrderStatusMappingDTO,
    CreateOrderStatusDTO,
    UpdateOrderStatusDTO,
    EcommerceIntegrationType,
    ChannelStatusInfo,
    CreateChannelStatusDTO,
    UpdateChannelStatusDTO,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeOrderStatusMapping = (overrides: Partial<OrderStatusMapping> = {}): OrderStatusMapping => ({
    id: 1,
    integration_type_id: 1,
    integration_type: { id: 1, code: 'shopify', name: 'Shopify' },
    original_status: 'fulfilled',
    order_status_id: 5,
    order_status: { id: 5, code: 'completed', name: 'Completada' },
    is_active: true,
    description: 'Mapeo de prueba',
    created_at: '2026-03-01T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
    ...overrides,
});

const makeOrderStatusInfo = (overrides: Partial<OrderStatusInfo> = {}): OrderStatusInfo => ({
    id: 1,
    code: 'pending',
    name: 'Pendiente',
    description: 'Orden pendiente',
    category: 'active',
    color: '#FFA500',
    priority: 1,
    is_active: true,
    ...overrides,
});

const makeEcommerceIntegrationType = (overrides: Partial<EcommerceIntegrationType> = {}): EcommerceIntegrationType => ({
    id: 1,
    code: 'shopify',
    name: 'Shopify',
    ...overrides,
});

const makeChannelStatusInfo = (overrides: Partial<ChannelStatusInfo> = {}): ChannelStatusInfo => ({
    id: 1,
    integration_type_id: 1,
    code: 'fulfilled',
    name: 'Fulfilled',
    is_active: true,
    display_order: 1,
    ...overrides,
});

const paginatedMappings: PaginatedResponse<OrderStatusMapping> = {
    success: true,
    message: 'OK',
    data: [makeOrderStatusMapping()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const singleMapping: SingleResponse<OrderStatusMapping> = {
    success: true,
    message: 'OK',
    data: makeOrderStatusMapping(),
};

const singleStatus: SingleResponse<OrderStatusInfo> = {
    success: true,
    message: 'OK',
    data: makeOrderStatusInfo(),
};

const singleChannelStatus: SingleResponse<ChannelStatusInfo> = {
    success: true,
    message: 'OK',
    data: makeChannelStatusInfo(),
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IOrderStatusMappingRepository {
    return {
        getOrderStatusMappings: vi.fn(),
        getOrderStatusMappingById: vi.fn(),
        createOrderStatusMapping: vi.fn(),
        updateOrderStatusMapping: vi.fn(),
        deleteOrderStatusMapping: vi.fn(),
        toggleOrderStatusMappingActive: vi.fn(),
        getOrderStatuses: vi.fn(),
        createOrderStatus: vi.fn(),
        getOrderStatusById: vi.fn(),
        updateOrderStatus: vi.fn(),
        deleteOrderStatus: vi.fn(),
        getEcommerceIntegrationTypes: vi.fn(),
        getChannelStatuses: vi.fn(),
        createChannelStatus: vi.fn(),
        updateChannelStatus: vi.fn(),
        deleteChannelStatus: vi.fn(),
    } as unknown as IOrderStatusMappingRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('OrderStatusMappingUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: OrderStatusMappingUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new OrderStatusMappingUseCases(repo as unknown as IOrderStatusMappingRepository);
    });

    // ---------------------------------------------------------------
    // getOrderStatusMappings
    // ---------------------------------------------------------------
    describe('getOrderStatusMappings', () => {
        it('debería retornar la lista paginada de mapeos cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getOrderStatusMappings).mockResolvedValue(paginatedMappings);

            const result = await useCases.getOrderStatusMappings({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedMappings);
            expect(repo.getOrderStatusMappings).toHaveBeenCalledOnce();
            expect(repo.getOrderStatusMappings).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getOrderStatusMappings).mockResolvedValue(paginatedMappings);

            await useCases.getOrderStatusMappings();

            expect(repo.getOrderStatusMappings).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getOrderStatusMappings).mockRejectedValue(new Error('Fallo de conexión'));

            await expect(useCases.getOrderStatusMappings()).rejects.toThrow('Fallo de conexión');
        });
    });

    // ---------------------------------------------------------------
    // getOrderStatusMappingById
    // ---------------------------------------------------------------
    describe('getOrderStatusMappingById', () => {
        it('debería retornar un mapeo por ID', async () => {
            vi.mocked(repo.getOrderStatusMappingById).mockResolvedValue(singleMapping);

            const result = await useCases.getOrderStatusMappingById(1);

            expect(result).toEqual(singleMapping);
            expect(repo.getOrderStatusMappingById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando el mapeo no existe', async () => {
            vi.mocked(repo.getOrderStatusMappingById).mockRejectedValue(new Error('No encontrado'));

            await expect(useCases.getOrderStatusMappingById(999)).rejects.toThrow('No encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createOrderStatusMapping
    // ---------------------------------------------------------------
    describe('createOrderStatusMapping', () => {
        const dto: CreateOrderStatusMappingDTO = {
            integration_type_id: 1,
            original_status: 'paid',
            order_status_id: 3,
            description: 'Nuevo mapeo',
        };

        it('debería crear un mapeo y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createOrderStatusMapping).mockResolvedValue(singleMapping);

            const result = await useCases.createOrderStatusMapping(dto);

            expect(result).toEqual(singleMapping);
            expect(repo.createOrderStatusMapping).toHaveBeenCalledOnce();
            expect(repo.createOrderStatusMapping).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createOrderStatusMapping).mockRejectedValue(new Error('Mapeo duplicado'));

            await expect(useCases.createOrderStatusMapping(dto)).rejects.toThrow('Mapeo duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateOrderStatusMapping
    // ---------------------------------------------------------------
    describe('updateOrderStatusMapping', () => {
        const updateDto: UpdateOrderStatusMappingDTO = {
            original_status: 'shipped',
            order_status_id: 4,
            description: 'Mapeo actualizado',
        };

        it('debería actualizar un mapeo y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<OrderStatusMapping> = {
                ...singleMapping,
                data: makeOrderStatusMapping({ original_status: 'shipped', description: 'Mapeo actualizado' }),
            };
            vi.mocked(repo.updateOrderStatusMapping).mockResolvedValue(updatedResponse);

            const result = await useCases.updateOrderStatusMapping(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateOrderStatusMapping).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateOrderStatusMapping).mockRejectedValue(new Error('Mapeo no encontrado'));

            await expect(useCases.updateOrderStatusMapping(99, updateDto)).rejects.toThrow('Mapeo no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteOrderStatusMapping
    // ---------------------------------------------------------------
    describe('deleteOrderStatusMapping', () => {
        it('debería eliminar un mapeo y retornar confirmación', async () => {
            vi.mocked(repo.deleteOrderStatusMapping).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteOrderStatusMapping(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteOrderStatusMapping).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el mapeo no existe', async () => {
            vi.mocked(repo.deleteOrderStatusMapping).mockResolvedValue(actionError);

            const result = await useCases.deleteOrderStatusMapping(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteOrderStatusMapping).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteOrderStatusMapping(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // toggleOrderStatusMappingActive
    // ---------------------------------------------------------------
    describe('toggleOrderStatusMappingActive', () => {
        it('debería alternar el estado activo de un mapeo', async () => {
            const toggledResponse: SingleResponse<OrderStatusMapping> = {
                ...singleMapping,
                data: makeOrderStatusMapping({ is_active: false }),
            };
            vi.mocked(repo.toggleOrderStatusMappingActive).mockResolvedValue(toggledResponse);

            const result = await useCases.toggleOrderStatusMappingActive(1);

            expect(result).toEqual(toggledResponse);
            expect(repo.toggleOrderStatusMappingActive).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando el toggle falla', async () => {
            vi.mocked(repo.toggleOrderStatusMappingActive).mockRejectedValue(new Error('No autorizado'));

            await expect(useCases.toggleOrderStatusMappingActive(1)).rejects.toThrow('No autorizado');
        });
    });

    // ---------------------------------------------------------------
    // getOrderStatuses
    // ---------------------------------------------------------------
    describe('getOrderStatuses', () => {
        const statusesResponse = {
            success: true,
            data: [makeOrderStatusInfo()],
            message: 'OK',
        };

        it('debería retornar la lista de estados de orden', async () => {
            vi.mocked(repo.getOrderStatuses).mockResolvedValue(statusesResponse);

            const result = await useCases.getOrderStatuses();

            expect(result).toEqual(statusesResponse);
            expect(repo.getOrderStatuses).toHaveBeenCalledWith(undefined);
        });

        it('debería filtrar por isActive cuando se proporciona', async () => {
            vi.mocked(repo.getOrderStatuses).mockResolvedValue(statusesResponse);

            await useCases.getOrderStatuses(true);

            expect(repo.getOrderStatuses).toHaveBeenCalledWith(true);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getOrderStatuses).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getOrderStatuses()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // createOrderStatus
    // ---------------------------------------------------------------
    describe('createOrderStatus', () => {
        const dto: CreateOrderStatusDTO = {
            code: 'in_transit',
            name: 'En Tránsito',
            description: 'Orden en camino',
            category: 'active',
            color: '#0000FF',
        };

        it('debería crear un estado de orden y retornar la respuesta', async () => {
            vi.mocked(repo.createOrderStatus).mockResolvedValue(singleStatus);

            const result = await useCases.createOrderStatus(dto);

            expect(result).toEqual(singleStatus);
            expect(repo.createOrderStatus).toHaveBeenCalledOnce();
            expect(repo.createOrderStatus).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createOrderStatus).mockRejectedValue(new Error('Código duplicado'));

            await expect(useCases.createOrderStatus(dto)).rejects.toThrow('Código duplicado');
        });
    });

    // ---------------------------------------------------------------
    // getOrderStatusById
    // ---------------------------------------------------------------
    describe('getOrderStatusById', () => {
        it('debería retornar un estado de orden por ID', async () => {
            vi.mocked(repo.getOrderStatusById).mockResolvedValue(singleStatus);

            const result = await useCases.getOrderStatusById(1);

            expect(result).toEqual(singleStatus);
            expect(repo.getOrderStatusById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando el estado no existe', async () => {
            vi.mocked(repo.getOrderStatusById).mockRejectedValue(new Error('Estado no encontrado'));

            await expect(useCases.getOrderStatusById(999)).rejects.toThrow('Estado no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // updateOrderStatus
    // ---------------------------------------------------------------
    describe('updateOrderStatus', () => {
        const updateDto: UpdateOrderStatusDTO = {
            code: 'pending',
            name: 'Pendiente Actualizado',
            color: '#FF0000',
        };

        it('debería actualizar un estado de orden y retornar la respuesta', async () => {
            const updatedResponse: SingleResponse<OrderStatusInfo> = {
                ...singleStatus,
                data: makeOrderStatusInfo({ name: 'Pendiente Actualizado', color: '#FF0000' }),
            };
            vi.mocked(repo.updateOrderStatus).mockResolvedValue(updatedResponse);

            const result = await useCases.updateOrderStatus(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateOrderStatus).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateOrderStatus).mockRejectedValue(new Error('Estado no encontrado'));

            await expect(useCases.updateOrderStatus(99, updateDto)).rejects.toThrow('Estado no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteOrderStatus
    // ---------------------------------------------------------------
    describe('deleteOrderStatus', () => {
        it('debería eliminar un estado de orden y retornar confirmación', async () => {
            vi.mocked(repo.deleteOrderStatus).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteOrderStatus(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteOrderStatus).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el estado no existe', async () => {
            vi.mocked(repo.deleteOrderStatus).mockResolvedValue(actionError);

            const result = await useCases.deleteOrderStatus(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteOrderStatus).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteOrderStatus(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // getEcommerceIntegrationTypes
    // ---------------------------------------------------------------
    describe('getEcommerceIntegrationTypes', () => {
        const integrationTypesResponse = {
            success: true,
            data: [makeEcommerceIntegrationType()],
            message: 'OK',
        };

        it('debería retornar la lista de tipos de integración ecommerce', async () => {
            vi.mocked(repo.getEcommerceIntegrationTypes).mockResolvedValue(integrationTypesResponse);

            const result = await useCases.getEcommerceIntegrationTypes();

            expect(result).toEqual(integrationTypesResponse);
            expect(repo.getEcommerceIntegrationTypes).toHaveBeenCalledOnce();
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getEcommerceIntegrationTypes).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getEcommerceIntegrationTypes()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // getChannelStatuses
    // ---------------------------------------------------------------
    describe('getChannelStatuses', () => {
        const channelStatusesResponse = {
            success: true,
            data: [makeChannelStatusInfo()],
            message: 'OK',
        };

        it('debería retornar los estados de canal por tipo de integración', async () => {
            vi.mocked(repo.getChannelStatuses).mockResolvedValue(channelStatusesResponse);

            const result = await useCases.getChannelStatuses(1);

            expect(result).toEqual(channelStatusesResponse);
            expect(repo.getChannelStatuses).toHaveBeenCalledWith(1, undefined);
        });

        it('debería filtrar por isActive cuando se proporciona', async () => {
            vi.mocked(repo.getChannelStatuses).mockResolvedValue(channelStatusesResponse);

            await useCases.getChannelStatuses(1, true);

            expect(repo.getChannelStatuses).toHaveBeenCalledWith(1, true);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getChannelStatuses).mockRejectedValue(new Error('Error interno'));

            await expect(useCases.getChannelStatuses(1)).rejects.toThrow('Error interno');
        });
    });

    // ---------------------------------------------------------------
    // createChannelStatus
    // ---------------------------------------------------------------
    describe('createChannelStatus', () => {
        const dto: CreateChannelStatusDTO = {
            integration_type_id: 1,
            code: 'processing',
            name: 'Processing',
            is_active: true,
            display_order: 2,
        };

        it('debería crear un estado de canal y retornar la respuesta', async () => {
            vi.mocked(repo.createChannelStatus).mockResolvedValue(singleChannelStatus);

            const result = await useCases.createChannelStatus(dto);

            expect(result).toEqual(singleChannelStatus);
            expect(repo.createChannelStatus).toHaveBeenCalledOnce();
            expect(repo.createChannelStatus).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createChannelStatus).mockRejectedValue(new Error('Código duplicado'));

            await expect(useCases.createChannelStatus(dto)).rejects.toThrow('Código duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateChannelStatus
    // ---------------------------------------------------------------
    describe('updateChannelStatus', () => {
        const updateDto: UpdateChannelStatusDTO = {
            code: 'fulfilled',
            name: 'Fulfilled Updated',
            is_active: true,
            display_order: 1,
        };

        it('debería actualizar un estado de canal y retornar la respuesta', async () => {
            const updatedResponse: SingleResponse<ChannelStatusInfo> = {
                ...singleChannelStatus,
                data: makeChannelStatusInfo({ name: 'Fulfilled Updated' }),
            };
            vi.mocked(repo.updateChannelStatus).mockResolvedValue(updatedResponse);

            const result = await useCases.updateChannelStatus(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateChannelStatus).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateChannelStatus).mockRejectedValue(new Error('Estado no encontrado'));

            await expect(useCases.updateChannelStatus(99, updateDto)).rejects.toThrow('Estado no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteChannelStatus
    // ---------------------------------------------------------------
    describe('deleteChannelStatus', () => {
        it('debería eliminar un estado de canal y retornar confirmación', async () => {
            vi.mocked(repo.deleteChannelStatus).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteChannelStatus(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteChannelStatus).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el estado no existe', async () => {
            vi.mocked(repo.deleteChannelStatus).mockResolvedValue(actionError);

            const result = await useCases.deleteChannelStatus(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteChannelStatus).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteChannelStatus(1)).rejects.toThrow('Network error');
        });
    });
});
