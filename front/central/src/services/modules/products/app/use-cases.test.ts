import { describe, it, expect, vi, beforeEach } from 'vitest';
import { ProductUseCases } from './use-cases';
import { IProductRepository } from '../domain/ports';
import {
    Product,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    CreateProductDTO,
    UpdateProductDTO,
    AddProductIntegrationDTO,
    ProductIntegrationsResponse,
    UploadImageResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeProduct = (overrides: Partial<Product> = {}): Product => ({
    id: '1',
    created_at: '2026-03-01T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
    business_id: 1,
    sku: 'SKU-001',
    name: 'Producto de prueba',
    price: 50000,
    currency: 'COP',
    stock: 100,
    manage_stock: true,
    track_inventory: true,
    status: 'active',
    is_active: true,
    ...overrides,
});

const paginatedProducts: PaginatedResponse<Product> = {
    success: true,
    message: 'OK',
    data: [makeProduct()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const singleProduct: SingleResponse<Product> = {
    success: true,
    message: 'OK',
    data: makeProduct(),
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IProductRepository {
    return {
        getProducts: vi.fn(),
        getProductById: vi.fn(),
        createProduct: vi.fn(),
        updateProduct: vi.fn(),
        deleteProduct: vi.fn(),
        uploadProductImage: vi.fn(),
        addProductIntegration: vi.fn(),
        removeProductIntegration: vi.fn(),
        getProductIntegrations: vi.fn(),
    } as unknown as IProductRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('ProductUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: ProductUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new ProductUseCases(repo as unknown as IProductRepository);
    });

    // ---------------------------------------------------------------
    // getProducts
    // ---------------------------------------------------------------
    describe('getProducts', () => {
        it('debería retornar la lista paginada de productos cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getProducts).mockResolvedValue(paginatedProducts);

            const result = await useCases.getProducts({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedProducts);
            expect(repo.getProducts).toHaveBeenCalledOnce();
            expect(repo.getProducts).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getProducts).mockResolvedValue(paginatedProducts);

            await useCases.getProducts();

            expect(repo.getProducts).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Fallo de base de datos');
            vi.mocked(repo.getProducts).mockRejectedValue(expectedError);

            await expect(useCases.getProducts()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getProductById
    // ---------------------------------------------------------------
    describe('getProductById', () => {
        it('debería retornar un producto por su ID', async () => {
            vi.mocked(repo.getProductById).mockResolvedValue(singleProduct);

            const result = await useCases.getProductById('1');

            expect(result).toEqual(singleProduct);
            expect(repo.getProductById).toHaveBeenCalledWith('1', undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.getProductById).mockResolvedValue(singleProduct);

            await useCases.getProductById('1', 5);

            expect(repo.getProductById).toHaveBeenCalledWith('1', 5);
        });

        it('debería propagar el error cuando el producto no existe', async () => {
            vi.mocked(repo.getProductById).mockRejectedValue(new Error('Producto no encontrado'));

            await expect(useCases.getProductById('999')).rejects.toThrow('Producto no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createProduct
    // ---------------------------------------------------------------
    describe('createProduct', () => {
        const dto: CreateProductDTO = {
            business_id: 1,
            sku: 'SKU-002',
            name: 'Nuevo Producto',
            price: 30000,
            stock: 50,
        };

        it('debería crear un producto y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createProduct).mockResolvedValue(singleProduct);

            const result = await useCases.createProduct(dto);

            expect(result).toEqual(singleProduct);
            expect(repo.createProduct).toHaveBeenCalledOnce();
            expect(repo.createProduct).toHaveBeenCalledWith(dto, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.createProduct).mockResolvedValue(singleProduct);

            await useCases.createProduct(dto, 5);

            expect(repo.createProduct).toHaveBeenCalledWith(dto, 5);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createProduct).mockRejectedValue(new Error('SKU duplicado'));

            await expect(useCases.createProduct(dto)).rejects.toThrow('SKU duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateProduct
    // ---------------------------------------------------------------
    describe('updateProduct', () => {
        const updateDto: UpdateProductDTO = { name: 'Producto Actualizado', price: 55000 };

        it('debería actualizar un producto y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<Product> = {
                ...singleProduct,
                data: makeProduct({ name: 'Producto Actualizado', price: 55000 }),
            };
            vi.mocked(repo.updateProduct).mockResolvedValue(updatedResponse);

            const result = await useCases.updateProduct('1', updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateProduct).toHaveBeenCalledWith('1', updateDto, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.updateProduct).mockResolvedValue(singleProduct);

            await useCases.updateProduct('1', updateDto, 5);

            expect(repo.updateProduct).toHaveBeenCalledWith('1', updateDto, 5);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateProduct).mockRejectedValue(new Error('Producto no encontrado'));

            await expect(useCases.updateProduct('999', updateDto)).rejects.toThrow('Producto no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteProduct
    // ---------------------------------------------------------------
    describe('deleteProduct', () => {
        it('debería eliminar un producto y retornar confirmación', async () => {
            vi.mocked(repo.deleteProduct).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteProduct('1');

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteProduct).toHaveBeenCalledWith('1', undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.deleteProduct).mockResolvedValue(actionSuccess);

            await useCases.deleteProduct('1', 5);

            expect(repo.deleteProduct).toHaveBeenCalledWith('1', 5);
        });

        it('debería retornar respuesta de error cuando el producto no existe', async () => {
            vi.mocked(repo.deleteProduct).mockResolvedValue(actionError);

            const result = await useCases.deleteProduct('999');

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteProduct).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteProduct('1')).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // uploadProductImage
    // ---------------------------------------------------------------
    describe('uploadProductImage', () => {
        it('debería subir una imagen y retornar la URL', async () => {
            const uploadResponse: UploadImageResponse = {
                success: true,
                message: 'Imagen subida',
                image_url: 'https://cdn.example.com/img.jpg',
            };
            const formData = new FormData();
            vi.mocked(repo.uploadProductImage).mockResolvedValue(uploadResponse);

            const result = await useCases.uploadProductImage('1', formData);

            expect(result).toEqual(uploadResponse);
            expect(repo.uploadProductImage).toHaveBeenCalledWith('1', formData, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            const uploadResponse: UploadImageResponse = {
                success: true,
                message: 'OK',
                image_url: 'https://cdn.example.com/img.jpg',
            };
            const formData = new FormData();
            vi.mocked(repo.uploadProductImage).mockResolvedValue(uploadResponse);

            await useCases.uploadProductImage('1', formData, 5);

            expect(repo.uploadProductImage).toHaveBeenCalledWith('1', formData, 5);
        });

        it('debería propagar el error cuando la subida falla', async () => {
            const formData = new FormData();
            vi.mocked(repo.uploadProductImage).mockRejectedValue(new Error('Archivo muy grande'));

            await expect(useCases.uploadProductImage('1', formData)).rejects.toThrow('Archivo muy grande');
        });
    });

    // ---------------------------------------------------------------
    // addProductIntegration
    // ---------------------------------------------------------------
    describe('addProductIntegration', () => {
        const integrationDto: AddProductIntegrationDTO = {
            integration_id: 10,
            external_product_id: 'ext-prod-001',
        };

        it('debería agregar una integración al producto y retornar la respuesta', async () => {
            const response: SingleResponse<any> = { success: true, message: 'OK', data: { id: 1 } };
            vi.mocked(repo.addProductIntegration).mockResolvedValue(response);

            const result = await useCases.addProductIntegration('1', integrationDto);

            expect(result).toEqual(response);
            expect(repo.addProductIntegration).toHaveBeenCalledWith('1', integrationDto, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            const response: SingleResponse<any> = { success: true, message: 'OK', data: { id: 1 } };
            vi.mocked(repo.addProductIntegration).mockResolvedValue(response);

            await useCases.addProductIntegration('1', integrationDto, 5);

            expect(repo.addProductIntegration).toHaveBeenCalledWith('1', integrationDto, 5);
        });

        it('debería propagar el error cuando falla agregar la integración', async () => {
            vi.mocked(repo.addProductIntegration).mockRejectedValue(new Error('Integración duplicada'));

            await expect(useCases.addProductIntegration('1', integrationDto)).rejects.toThrow('Integración duplicada');
        });
    });

    // ---------------------------------------------------------------
    // removeProductIntegration
    // ---------------------------------------------------------------
    describe('removeProductIntegration', () => {
        it('debería remover una integración del producto y retornar confirmación', async () => {
            vi.mocked(repo.removeProductIntegration).mockResolvedValue(actionSuccess);

            const result = await useCases.removeProductIntegration('1', 10);

            expect(result).toEqual(actionSuccess);
            expect(repo.removeProductIntegration).toHaveBeenCalledWith('1', 10, undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            vi.mocked(repo.removeProductIntegration).mockResolvedValue(actionSuccess);

            await useCases.removeProductIntegration('1', 10, 5);

            expect(repo.removeProductIntegration).toHaveBeenCalledWith('1', 10, 5);
        });

        it('debería propagar el error cuando falla la remoción', async () => {
            vi.mocked(repo.removeProductIntegration).mockRejectedValue(new Error('Integración no encontrada'));

            await expect(useCases.removeProductIntegration('1', 99)).rejects.toThrow('Integración no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // getProductIntegrations
    // ---------------------------------------------------------------
    describe('getProductIntegrations', () => {
        it('debería retornar las integraciones de un producto', async () => {
            const integrationsResponse: ProductIntegrationsResponse = {
                success: true,
                message: 'OK',
                data: [
                    {
                        id: 1,
                        product_id: '1',
                        integration_id: 10,
                        external_product_id: 'ext-001',
                        created_at: '2026-03-01T00:00:00Z',
                        updated_at: '2026-03-01T00:00:00Z',
                    },
                ],
                total: 1,
            };
            vi.mocked(repo.getProductIntegrations).mockResolvedValue(integrationsResponse);

            const result = await useCases.getProductIntegrations('1');

            expect(result).toEqual(integrationsResponse);
            expect(repo.getProductIntegrations).toHaveBeenCalledWith('1', undefined);
        });

        it('debería pasar businessId cuando se proporciona', async () => {
            const integrationsResponse: ProductIntegrationsResponse = {
                success: true,
                message: 'OK',
                data: [],
                total: 0,
            };
            vi.mocked(repo.getProductIntegrations).mockResolvedValue(integrationsResponse);

            await useCases.getProductIntegrations('1', 5);

            expect(repo.getProductIntegrations).toHaveBeenCalledWith('1', 5);
        });

        it('debería propagar el error cuando la consulta de integraciones falla', async () => {
            vi.mocked(repo.getProductIntegrations).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getProductIntegrations('1')).rejects.toThrow('Servicio no disponible');
        });
    });
});
