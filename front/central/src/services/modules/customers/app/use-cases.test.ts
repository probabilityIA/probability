import { describe, it, expect, vi, beforeEach } from 'vitest';
import { CustomerUseCases } from './use-cases';
import { ICustomerRepository } from '../domain/ports';
import {
    CustomerInfo,
    CustomerDetail,
    CustomersListResponse,
    CreateCustomerDTO,
    UpdateCustomerDTO,
    DeleteCustomerResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeCustomerInfo = (overrides: Partial<CustomerInfo> = {}): CustomerInfo => ({
    id: 1,
    business_id: 1,
    name: 'Juan Perez',
    email: 'juan@test.com',
    phone: '3001234567',
    dni: '123456789',
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

const customersListResponse: CustomersListResponse = {
    data: [makeCustomerInfo()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const deleteSuccess: DeleteCustomerResponse = { message: 'Cliente eliminado' };
const deleteError: DeleteCustomerResponse = { error: 'Cliente no encontrado' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): ICustomerRepository {
    return {
        getCustomers: vi.fn(),
        getCustomerById: vi.fn(),
        createCustomer: vi.fn(),
        updateCustomer: vi.fn(),
        deleteCustomer: vi.fn(),
    } as unknown as ICustomerRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('CustomerUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: CustomerUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new CustomerUseCases(repo as unknown as ICustomerRepository);
    });

    // ---------------------------------------------------------------
    // getCustomers
    // ---------------------------------------------------------------
    describe('getCustomers', () => {
        it('debería retornar la lista paginada de clientes cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getCustomers).mockResolvedValue(customersListResponse);

            const result = await useCases.getCustomers({ page: 1, page_size: 10 });

            expect(result).toEqual(customersListResponse);
            expect(repo.getCustomers).toHaveBeenCalledOnce();
            expect(repo.getCustomers).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getCustomers).mockResolvedValue(customersListResponse);

            await useCases.getCustomers();

            expect(repo.getCustomers).toHaveBeenCalledWith(undefined);
        });

        it('debería pasar search y business_id como filtros', async () => {
            vi.mocked(repo.getCustomers).mockResolvedValue(customersListResponse);

            await useCases.getCustomers({ search: 'Juan', business_id: 5 });

            expect(repo.getCustomers).toHaveBeenCalledWith({ search: 'Juan', business_id: 5 });
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getCustomers).mockRejectedValue(expectedError);

            await expect(useCases.getCustomers()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getCustomerById
    // ---------------------------------------------------------------
    describe('getCustomerById', () => {
        it('debería retornar el detalle de un cliente por su ID', async () => {
            const detail = makeCustomerDetail();
            vi.mocked(repo.getCustomerById).mockResolvedValue(detail);

            const result = await useCases.getCustomerById(1);

            expect(result).toEqual(detail);
            expect(repo.getCustomerById).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.getCustomerById).mockResolvedValue(makeCustomerDetail());

            await useCases.getCustomerById(1, 5);

            expect(repo.getCustomerById).toHaveBeenCalledWith(1, 5);
        });

        it('debería propagar el error cuando el cliente no existe', async () => {
            vi.mocked(repo.getCustomerById).mockRejectedValue(new Error('Cliente no encontrado'));

            await expect(useCases.getCustomerById(999)).rejects.toThrow('Cliente no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createCustomer
    // ---------------------------------------------------------------
    describe('createCustomer', () => {
        const dto: CreateCustomerDTO = {
            name: 'Nuevo Cliente',
            email: 'nuevo@test.com',
            phone: '3009876543',
            dni: '987654321',
        };

        it('debería crear un cliente y retornar la respuesta del repositorio', async () => {
            const newCustomer = makeCustomerInfo({ id: 2, name: 'Nuevo Cliente' });
            vi.mocked(repo.createCustomer).mockResolvedValue(newCustomer);

            const result = await useCases.createCustomer(dto);

            expect(result).toEqual(newCustomer);
            expect(repo.createCustomer).toHaveBeenCalledOnce();
            expect(repo.createCustomer).toHaveBeenCalledWith(dto, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.createCustomer).mockResolvedValue(makeCustomerInfo());

            await useCases.createCustomer(dto, 5);

            expect(repo.createCustomer).toHaveBeenCalledWith(dto, 5);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createCustomer).mockRejectedValue(new Error('Email duplicado'));

            await expect(useCases.createCustomer(dto)).rejects.toThrow('Email duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateCustomer
    // ---------------------------------------------------------------
    describe('updateCustomer', () => {
        const updateDto: UpdateCustomerDTO = { name: 'Cliente Actualizado', email: 'actualizado@test.com' };

        it('debería actualizar un cliente y retornar la respuesta del repositorio', async () => {
            const updatedCustomer = makeCustomerInfo({ name: 'Cliente Actualizado' });
            vi.mocked(repo.updateCustomer).mockResolvedValue(updatedCustomer);

            const result = await useCases.updateCustomer(1, updateDto);

            expect(result).toEqual(updatedCustomer);
            expect(repo.updateCustomer).toHaveBeenCalledWith(1, updateDto, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.updateCustomer).mockResolvedValue(makeCustomerInfo());

            await useCases.updateCustomer(1, updateDto, 5);

            expect(repo.updateCustomer).toHaveBeenCalledWith(1, updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateCustomer).mockRejectedValue(new Error('Cliente no encontrado'));

            await expect(useCases.updateCustomer(99, updateDto)).rejects.toThrow('Cliente no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteCustomer
    // ---------------------------------------------------------------
    describe('deleteCustomer', () => {
        it('debería eliminar un cliente y retornar confirmación', async () => {
            vi.mocked(repo.deleteCustomer).mockResolvedValue(deleteSuccess);

            const result = await useCases.deleteCustomer(1);

            expect(result).toEqual(deleteSuccess);
            expect(repo.deleteCustomer).toHaveBeenCalledWith(1, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteCustomer).mockResolvedValue(deleteSuccess);

            await useCases.deleteCustomer(1, 5);

            expect(repo.deleteCustomer).toHaveBeenCalledWith(1, 5);
        });

        it('debería retornar respuesta de error cuando el cliente no existe', async () => {
            vi.mocked(repo.deleteCustomer).mockResolvedValue(deleteError);

            const result = await useCases.deleteCustomer(999);

            expect(result.error).toBeDefined();
            expect(result.message).toBeUndefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteCustomer).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteCustomer(1)).rejects.toThrow('Network error');
        });
    });
});
