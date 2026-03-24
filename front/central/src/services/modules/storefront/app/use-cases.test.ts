import { describe, it, expect, vi, beforeEach } from 'vitest';
import { StorefrontUseCases } from './use-cases';
import { IStorefrontRepository } from '../domain/ports';
import {
    StorefrontProduct,
    StorefrontOrder,
    StorefrontOrderItem,
    CreateStorefrontOrderDTO,
    RegisterDTO,
    PaginatedResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeProduct = (overrides: Partial<StorefrontProduct> = {}): StorefrontProduct => ({
    id: 'prod-001',
    name: 'Camiseta Negra',
    description: 'Camiseta de algodón 100%',
    short_description: 'Camiseta algodón',
    price: 45000,
    currency: 'COP',
    image_url: 'https://example.com/camiseta.png',
    sku: 'CAM-001',
    stock_quantity: 50,
    category: 'Ropa',
    brand: 'MarcaX',
    is_featured: true,
    created_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeOrderItem = (overrides: Partial<StorefrontOrderItem> = {}): StorefrontOrderItem => ({
    product_name: 'Camiseta Negra',
    quantity: 2,
    unit_price: 45000,
    total_price: 90000,
    ...overrides,
});

const makeOrder = (overrides: Partial<StorefrontOrder> = {}): StorefrontOrder => ({
    id: 'ord-001',
    order_number: '10001',
    status: 'pending',
    total_amount: 90000,
    currency: 'COP',
    created_at: '2026-03-24T00:00:00Z',
    items: [makeOrderItem()],
    ...overrides,
});

const paginatedProducts: PaginatedResponse<StorefrontProduct> = {
    data: [makeProduct()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const paginatedOrders: PaginatedResponse<StorefrontOrder> = {
    data: [makeOrder()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IStorefrontRepository {
    return {
        getCatalog: vi.fn(),
        getProduct: vi.fn(),
        createOrder: vi.fn(),
        getOrders: vi.fn(),
        getOrder: vi.fn(),
        register: vi.fn(),
    } as unknown as IStorefrontRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('StorefrontUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: StorefrontUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new StorefrontUseCases(repo as unknown as IStorefrontRepository);
    });

    // ---------------------------------------------------------------
    // getCatalog
    // ---------------------------------------------------------------
    describe('getCatalog', () => {
        it('debería retornar el catálogo paginado de productos cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getCatalog).mockResolvedValue(paginatedProducts);

            const result = await useCases.getCatalog({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedProducts);
            expect(repo.getCatalog).toHaveBeenCalledOnce();
            expect(repo.getCatalog).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getCatalog).mockResolvedValue(paginatedProducts);

            await useCases.getCatalog();

            expect(repo.getCatalog).toHaveBeenCalledWith(undefined);
        });

        it('debería pasar filtros de búsqueda y categoría cuando se proporcionan', async () => {
            vi.mocked(repo.getCatalog).mockResolvedValue(paginatedProducts);

            await useCases.getCatalog({ search: 'camiseta', category: 'ropa', business_id: 1 });

            expect(repo.getCatalog).toHaveBeenCalledWith({ search: 'camiseta', category: 'ropa', business_id: 1 });
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getCatalog).mockRejectedValue(new Error('Fallo de base de datos'));

            await expect(useCases.getCatalog()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getProduct
    // ---------------------------------------------------------------
    describe('getProduct', () => {
        it('debería retornar un producto por su ID', async () => {
            const product = makeProduct();
            vi.mocked(repo.getProduct).mockResolvedValue(product);

            const result = await useCases.getProduct('prod-001');

            expect(result).toEqual(product);
            expect(repo.getProduct).toHaveBeenCalledOnce();
            expect(repo.getProduct).toHaveBeenCalledWith('prod-001', undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const product = makeProduct();
            vi.mocked(repo.getProduct).mockResolvedValue(product);

            await useCases.getProduct('prod-001', 5);

            expect(repo.getProduct).toHaveBeenCalledWith('prod-001', 5);
        });

        it('debería propagar el error cuando el producto no existe', async () => {
            vi.mocked(repo.getProduct).mockRejectedValue(new Error('Producto no encontrado'));

            await expect(useCases.getProduct('no-existe')).rejects.toThrow('Producto no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createOrder
    // ---------------------------------------------------------------
    describe('createOrder', () => {
        const dto: CreateStorefrontOrderDTO = {
            items: [{ product_id: 'prod-001', quantity: 2 }],
            notes: 'Entregar en la tarde',
        };

        it('debería crear una orden y retornar la respuesta del repositorio', async () => {
            const response = { message: 'Orden creada exitosamente' };
            vi.mocked(repo.createOrder).mockResolvedValue(response);

            const result = await useCases.createOrder(dto);

            expect(result).toEqual(response);
            expect(repo.createOrder).toHaveBeenCalledOnce();
            expect(repo.createOrder).toHaveBeenCalledWith(dto, undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const response = { message: 'Orden creada exitosamente' };
            vi.mocked(repo.createOrder).mockResolvedValue(response);

            await useCases.createOrder(dto, 5);

            expect(repo.createOrder).toHaveBeenCalledWith(dto, 5);
        });

        it('debería propagar el error cuando la creación de la orden falla', async () => {
            vi.mocked(repo.createOrder).mockRejectedValue(new Error('Stock insuficiente'));

            await expect(useCases.createOrder(dto)).rejects.toThrow('Stock insuficiente');
        });
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

        it('debería propagar el error cuando la consulta de órdenes falla', async () => {
            vi.mocked(repo.getOrders).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getOrders()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // getOrder
    // ---------------------------------------------------------------
    describe('getOrder', () => {
        it('debería retornar una orden por su ID', async () => {
            const order = makeOrder();
            vi.mocked(repo.getOrder).mockResolvedValue(order);

            const result = await useCases.getOrder('ord-001');

            expect(result).toEqual(order);
            expect(repo.getOrder).toHaveBeenCalledOnce();
            expect(repo.getOrder).toHaveBeenCalledWith('ord-001', undefined);
        });

        it('debería pasar el businessId cuando se proporciona', async () => {
            const order = makeOrder();
            vi.mocked(repo.getOrder).mockResolvedValue(order);

            await useCases.getOrder('ord-001', 5);

            expect(repo.getOrder).toHaveBeenCalledWith('ord-001', 5);
        });

        it('debería propagar el error cuando la orden no existe', async () => {
            vi.mocked(repo.getOrder).mockRejectedValue(new Error('Orden no encontrada'));

            await expect(useCases.getOrder('no-existe')).rejects.toThrow('Orden no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // register
    // ---------------------------------------------------------------
    describe('register', () => {
        const dto: RegisterDTO = {
            name: 'Juan Perez',
            email: 'juan@test.com',
            password: 'securePassword123',
            phone: '3001234567',
            business_code: 'tienda-prueba',
        };

        it('debería registrar un usuario y retornar la respuesta del repositorio', async () => {
            const response = { message: 'Usuario registrado exitosamente' };
            vi.mocked(repo.register).mockResolvedValue(response);

            const result = await useCases.register(dto);

            expect(result).toEqual(response);
            expect(repo.register).toHaveBeenCalledOnce();
            expect(repo.register).toHaveBeenCalledWith(dto);
        });

        it('debería registrar un usuario con solo los campos obligatorios', async () => {
            const minimalDto: RegisterDTO = {
                name: 'Ana',
                email: 'ana@test.com',
                password: 'pass123',
                business_code: 'tienda-prueba',
            };
            const response = { message: 'Usuario registrado exitosamente' };
            vi.mocked(repo.register).mockResolvedValue(response);

            await useCases.register(minimalDto);

            expect(repo.register).toHaveBeenCalledWith(minimalDto);
        });

        it('debería propagar el error cuando el registro falla', async () => {
            vi.mocked(repo.register).mockRejectedValue(new Error('Email ya registrado'));

            await expect(useCases.register(dto)).rejects.toThrow('Email ya registrado');
        });
    });
});
