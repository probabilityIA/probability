import { describe, it, expect, vi, beforeEach } from 'vitest';
import { CustomerUseCases } from './use-cases';
import { ICustomerRepository } from '../domain/ports';
import {
    CustomerInfo,
    CustomerDetail,
    CustomerSummary,
    CustomerAddress,
    CustomerProduct,
    CustomerOrderItem,
    CustomersListResponse,
    CustomerAddressListResponse,
    CustomerProductListResponse,
    CustomerOrderItemListResponse,
    CreateCustomerDTO,
    UpdateCustomerDTO,
    DeleteCustomerResponse,
} from '../domain/types';

const makeCustomerInfo = (overrides: Partial<CustomerInfo> = {}): CustomerInfo => ({
    id: 1,
    business_id: 1,
    name: 'Juan Perez',
    email: 'juan@test.com',
    phone: '3001234567',
    dni: '123456789',
    total_orders: 5,
    created_at: '2026-03-01T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
    ...overrides,
});

const makeCustomerDetail = (overrides: Partial<CustomerDetail> = {}): CustomerDetail => ({
    ...makeCustomerInfo(),
    order_count: 5,
    total_spent: 250000,
    last_order_at: '2026-03-15T00:00:00Z',
    ...overrides,
});

const makeSummary = (overrides: Partial<CustomerSummary> = {}): CustomerSummary => ({
    id: 1,
    customer_id: 1,
    business_id: 1,
    total_orders: 10,
    delivered_orders: 7,
    cancelled_orders: 1,
    in_progress_orders: 2,
    total_spent: 500000,
    avg_ticket: 50000,
    total_paid_orders: 8,
    avg_delivery_score: 85.5,
    first_order_at: '2026-01-01T00:00:00Z',
    last_order_at: '2026-03-15T00:00:00Z',
    preferred_platform: 'shopify',
    last_updated_at: '2026-03-15T00:00:00Z',
    ...overrides,
});

const makeAddress = (overrides: Partial<CustomerAddress> = {}): CustomerAddress => ({
    id: 1,
    customer_id: 1,
    business_id: 1,
    street: 'Calle 100 #15-20',
    city: 'Bogota',
    state: 'Cundinamarca',
    country: 'Colombia',
    postal_code: '110111',
    times_used: 3,
    last_used_at: '2026-03-10T00:00:00Z',
    ...overrides,
});

const makeProduct = (overrides: Partial<CustomerProduct> = {}): CustomerProduct => ({
    id: 1,
    customer_id: 1,
    business_id: 1,
    product_id: 'PRD_001',
    product_name: 'Creatina 300g',
    product_sku: 'PT01004',
    product_image: null,
    times_ordered: 3,
    total_quantity: 5,
    total_spent: 150000,
    first_ordered_at: '2026-01-15T00:00:00Z',
    last_ordered_at: '2026-03-10T00:00:00Z',
    ...overrides,
});

const makeOrderItem = (overrides: Partial<CustomerOrderItem> = {}): CustomerOrderItem => ({
    id: 1,
    customer_id: 1,
    business_id: 1,
    order_id: 'abc-123',
    order_number: '#80001',
    product_id: 'PRD_001',
    product_name: 'Creatina 300g',
    product_sku: 'PT01004',
    product_image: null,
    quantity: 2,
    unit_price: 50000,
    total_price: 100000,
    order_status: 'delivered',
    ordered_at: '2026-03-10T00:00:00Z',
    ...overrides,
});

const customersListResponse: CustomersListResponse = {
    data: [makeCustomerInfo()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const addressListResponse: CustomerAddressListResponse = {
    data: [makeAddress()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const productListResponse: CustomerProductListResponse = {
    data: [makeProduct()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const orderItemListResponse: CustomerOrderItemListResponse = {
    data: [makeOrderItem()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const deleteSuccess: DeleteCustomerResponse = { message: 'Cliente eliminado' };
const deleteError: DeleteCustomerResponse = { error: 'Cliente no encontrado' };

function createMockRepository(): ICustomerRepository {
    return {
        getCustomers: vi.fn(),
        getCustomerById: vi.fn(),
        createCustomer: vi.fn(),
        updateCustomer: vi.fn(),
        deleteCustomer: vi.fn(),
        getCustomerSummary: vi.fn(),
        getCustomerAddresses: vi.fn(),
        getCustomerProducts: vi.fn(),
        getCustomerOrderItems: vi.fn(),
    };
}

describe('CustomerUseCases', () => {
    let repo: ICustomerRepository;
    let useCases: CustomerUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new CustomerUseCases(repo);
    });

    describe('getCustomers', () => {
        it('retorna lista paginada de clientes', async () => {
            vi.mocked(repo.getCustomers).mockResolvedValue(customersListResponse);

            const result = await useCases.getCustomers({ page: 1, page_size: 10 });

            expect(result).toEqual(customersListResponse);
            expect(repo.getCustomers).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('llama al repositorio sin parametros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getCustomers).mockResolvedValue(customersListResponse);

            await useCases.getCustomers();

            expect(repo.getCustomers).toHaveBeenCalledWith(undefined);
        });

        it('pasa search y business_id como filtros', async () => {
            vi.mocked(repo.getCustomers).mockResolvedValue(customersListResponse);

            await useCases.getCustomers({ search: 'Juan', business_id: 5 });

            expect(repo.getCustomers).toHaveBeenCalledWith({ search: 'Juan', business_id: 5 });
        });

        it('propaga el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getCustomers).mockRejectedValue(new Error('DB error'));

            await expect(useCases.getCustomers()).rejects.toThrow('DB error');
        });
    });

    describe('getCustomerById', () => {
        it('retorna el detalle de un cliente por ID', async () => {
            const detail = makeCustomerDetail();
            vi.mocked(repo.getCustomerById).mockResolvedValue(detail);

            const result = await useCases.getCustomerById(1);

            expect(result).toEqual(detail);
            expect(repo.getCustomerById).toHaveBeenCalledWith(1, undefined);
        });

        it('pasa businessId cuando se proporciona', async () => {
            vi.mocked(repo.getCustomerById).mockResolvedValue(makeCustomerDetail());

            await useCases.getCustomerById(1, 5);

            expect(repo.getCustomerById).toHaveBeenCalledWith(1, 5);
        });

        it('propaga error cuando el cliente no existe', async () => {
            vi.mocked(repo.getCustomerById).mockRejectedValue(new Error('Not found'));

            await expect(useCases.getCustomerById(999)).rejects.toThrow('Not found');
        });
    });

    describe('createCustomer', () => {
        const dto: CreateCustomerDTO = {
            name: 'Nuevo Cliente',
            email: 'nuevo@test.com',
            phone: '3009876543',
            dni: '987654321',
        };

        it('crea un cliente y retorna la respuesta', async () => {
            const newCustomer = makeCustomerInfo({ id: 2, name: 'Nuevo Cliente' });
            vi.mocked(repo.createCustomer).mockResolvedValue(newCustomer);

            const result = await useCases.createCustomer(dto);

            expect(result).toEqual(newCustomer);
            expect(repo.createCustomer).toHaveBeenCalledWith(dto, undefined);
        });

        it('pasa businessId cuando se proporciona', async () => {
            vi.mocked(repo.createCustomer).mockResolvedValue(makeCustomerInfo());

            await useCases.createCustomer(dto, 5);

            expect(repo.createCustomer).toHaveBeenCalledWith(dto, 5);
        });

        it('propaga error cuando la creacion falla', async () => {
            vi.mocked(repo.createCustomer).mockRejectedValue(new Error('Email duplicado'));

            await expect(useCases.createCustomer(dto)).rejects.toThrow('Email duplicado');
        });
    });

    describe('updateCustomer', () => {
        const updateDto: UpdateCustomerDTO = { name: 'Cliente Actualizado', email: 'act@test.com' };

        it('actualiza un cliente y retorna la respuesta', async () => {
            const updated = makeCustomerInfo({ name: 'Cliente Actualizado' });
            vi.mocked(repo.updateCustomer).mockResolvedValue(updated);

            const result = await useCases.updateCustomer(1, updateDto);

            expect(result).toEqual(updated);
            expect(repo.updateCustomer).toHaveBeenCalledWith(1, updateDto, undefined);
        });

        it('pasa businessId cuando se proporciona', async () => {
            vi.mocked(repo.updateCustomer).mockResolvedValue(makeCustomerInfo());

            await useCases.updateCustomer(1, updateDto, 5);

            expect(repo.updateCustomer).toHaveBeenCalledWith(1, updateDto, 5);
        });

        it('propaga error cuando la actualizacion falla', async () => {
            vi.mocked(repo.updateCustomer).mockRejectedValue(new Error('Not found'));

            await expect(useCases.updateCustomer(99, updateDto)).rejects.toThrow('Not found');
        });
    });

    describe('deleteCustomer', () => {
        it('elimina un cliente y retorna confirmacion', async () => {
            vi.mocked(repo.deleteCustomer).mockResolvedValue(deleteSuccess);

            const result = await useCases.deleteCustomer(1);

            expect(result).toEqual(deleteSuccess);
            expect(repo.deleteCustomer).toHaveBeenCalledWith(1, undefined);
        });

        it('pasa businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteCustomer).mockResolvedValue(deleteSuccess);

            await useCases.deleteCustomer(1, 5);

            expect(repo.deleteCustomer).toHaveBeenCalledWith(1, 5);
        });

        it('retorna respuesta de error cuando el cliente no existe', async () => {
            vi.mocked(repo.deleteCustomer).mockResolvedValue(deleteError);

            const result = await useCases.deleteCustomer(999);

            expect(result.error).toBeDefined();
            expect(result.message).toBeUndefined();
        });

        it('propaga excepcion en error de red', async () => {
            vi.mocked(repo.deleteCustomer).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteCustomer(1)).rejects.toThrow('Network error');
        });
    });

    describe('getCustomerSummary', () => {
        it('retorna el resumen de un cliente', async () => {
            const summary = makeSummary();
            vi.mocked(repo.getCustomerSummary).mockResolvedValue(summary);

            const result = await useCases.getCustomerSummary(1);

            expect(result).toEqual(summary);
            expect(repo.getCustomerSummary).toHaveBeenCalledWith(1, undefined);
        });

        it('pasa businessId cuando se proporciona', async () => {
            vi.mocked(repo.getCustomerSummary).mockResolvedValue(makeSummary());

            await useCases.getCustomerSummary(1, 5);

            expect(repo.getCustomerSummary).toHaveBeenCalledWith(1, 5);
        });

        it('propaga error cuando no hay resumen', async () => {
            vi.mocked(repo.getCustomerSummary).mockRejectedValue(new Error('Not found'));

            await expect(useCases.getCustomerSummary(999)).rejects.toThrow('Not found');
        });
    });

    describe('getCustomerAddresses', () => {
        it('retorna lista paginada de direcciones', async () => {
            vi.mocked(repo.getCustomerAddresses).mockResolvedValue(addressListResponse);

            const result = await useCases.getCustomerAddresses(1, { page: 1, page_size: 10 });

            expect(result).toEqual(addressListResponse);
            expect(repo.getCustomerAddresses).toHaveBeenCalledWith(1, { page: 1, page_size: 10 });
        });

        it('llama sin params de paginacion', async () => {
            vi.mocked(repo.getCustomerAddresses).mockResolvedValue(addressListResponse);

            await useCases.getCustomerAddresses(1);

            expect(repo.getCustomerAddresses).toHaveBeenCalledWith(1, undefined);
        });

        it('pasa business_id en params', async () => {
            vi.mocked(repo.getCustomerAddresses).mockResolvedValue(addressListResponse);

            await useCases.getCustomerAddresses(1, { business_id: 5 });

            expect(repo.getCustomerAddresses).toHaveBeenCalledWith(1, { business_id: 5 });
        });

        it('propaga error', async () => {
            vi.mocked(repo.getCustomerAddresses).mockRejectedValue(new Error('DB error'));

            await expect(useCases.getCustomerAddresses(1)).rejects.toThrow('DB error');
        });
    });

    describe('getCustomerProducts', () => {
        it('retorna lista paginada de productos', async () => {
            vi.mocked(repo.getCustomerProducts).mockResolvedValue(productListResponse);

            const result = await useCases.getCustomerProducts(1, { page: 1, page_size: 10 });

            expect(result).toEqual(productListResponse);
            expect(repo.getCustomerProducts).toHaveBeenCalledWith(1, { page: 1, page_size: 10 });
        });

        it('llama sin params de paginacion', async () => {
            vi.mocked(repo.getCustomerProducts).mockResolvedValue(productListResponse);

            await useCases.getCustomerProducts(1);

            expect(repo.getCustomerProducts).toHaveBeenCalledWith(1, undefined);
        });

        it('propaga error', async () => {
            vi.mocked(repo.getCustomerProducts).mockRejectedValue(new Error('DB error'));

            await expect(useCases.getCustomerProducts(1)).rejects.toThrow('DB error');
        });
    });

    describe('getCustomerOrderItems', () => {
        it('retorna lista paginada de order items', async () => {
            vi.mocked(repo.getCustomerOrderItems).mockResolvedValue(orderItemListResponse);

            const result = await useCases.getCustomerOrderItems(1, { page: 1, page_size: 10 });

            expect(result).toEqual(orderItemListResponse);
            expect(repo.getCustomerOrderItems).toHaveBeenCalledWith(1, { page: 1, page_size: 10 });
        });

        it('llama sin params de paginacion', async () => {
            vi.mocked(repo.getCustomerOrderItems).mockResolvedValue(orderItemListResponse);

            await useCases.getCustomerOrderItems(1);

            expect(repo.getCustomerOrderItems).toHaveBeenCalledWith(1, undefined);
        });

        it('pasa business_id en params', async () => {
            vi.mocked(repo.getCustomerOrderItems).mockResolvedValue(orderItemListResponse);

            await useCases.getCustomerOrderItems(1, { business_id: 34 });

            expect(repo.getCustomerOrderItems).toHaveBeenCalledWith(1, { business_id: 34 });
        });

        it('propaga error', async () => {
            vi.mocked(repo.getCustomerOrderItems).mockRejectedValue(new Error('Network error'));

            await expect(useCases.getCustomerOrderItems(1)).rejects.toThrow('Network error');
        });
    });
});
