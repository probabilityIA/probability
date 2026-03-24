import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PublicSiteUseCases } from './use-cases';
import { IPublicSiteRepository } from '../domain/ports';
import {
    PublicBusiness,
    PublicProduct,
    PaginatedResponse,
    ContactFormDTO,
    WebsiteConfig,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeWebsiteConfig = (overrides: Partial<WebsiteConfig> = {}): WebsiteConfig => ({
    template: 'default',
    show_hero: true,
    show_about: true,
    show_featured_products: true,
    show_full_catalog: false,
    show_testimonials: false,
    show_location: true,
    show_contact: true,
    show_social_media: false,
    show_whatsapp: true,
    hero_content: null,
    about_content: null,
    testimonials_content: null,
    location_content: null,
    contact_content: null,
    social_media_content: null,
    whatsapp_content: null,
    ...overrides,
});

const makePublicBusiness = (overrides: Partial<PublicBusiness> = {}): PublicBusiness => ({
    id: 1,
    name: 'Tienda de Prueba',
    code: 'tienda-prueba',
    description: 'Una tienda de ejemplo',
    logo_url: 'https://example.com/logo.png',
    primary_color: '#FF5733',
    secondary_color: '#33FF57',
    tertiary_color: '#3357FF',
    quaternary_color: '#F3F3F3',
    navbar_image_url: 'https://example.com/navbar.png',
    website_config: makeWebsiteConfig(),
    featured_products: [],
    ...overrides,
});

const makePublicProduct = (overrides: Partial<PublicProduct> = {}): PublicProduct => ({
    id: 'prod-001',
    name: 'Producto de Prueba',
    description: 'Descripción del producto',
    short_description: 'Desc corta',
    price: 50000,
    currency: 'COP',
    image_url: 'https://example.com/product.png',
    sku: 'SKU-001',
    stock_quantity: 100,
    category: 'General',
    brand: 'MarcaX',
    is_featured: true,
    created_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const paginatedProducts: PaginatedResponse<PublicProduct> = {
    data: [makePublicProduct()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IPublicSiteRepository {
    return {
        getBusinessPage: vi.fn(),
        getCatalog: vi.fn(),
        getProduct: vi.fn(),
        submitContact: vi.fn(),
    } as unknown as IPublicSiteRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('PublicSiteUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: PublicSiteUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new PublicSiteUseCases(repo as unknown as IPublicSiteRepository);
    });

    // ---------------------------------------------------------------
    // getBusinessPage
    // ---------------------------------------------------------------
    describe('getBusinessPage', () => {
        it('debería retornar la página del negocio por su slug', async () => {
            const business = makePublicBusiness();
            vi.mocked(repo.getBusinessPage).mockResolvedValue(business);

            const result = await useCases.getBusinessPage('tienda-prueba');

            expect(result).toEqual(business);
            expect(repo.getBusinessPage).toHaveBeenCalledOnce();
            expect(repo.getBusinessPage).toHaveBeenCalledWith('tienda-prueba');
        });

        it('debería propagar el error cuando el negocio no existe', async () => {
            vi.mocked(repo.getBusinessPage).mockRejectedValue(new Error('Negocio no encontrado'));

            await expect(useCases.getBusinessPage('no-existe')).rejects.toThrow('Negocio no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getCatalog
    // ---------------------------------------------------------------
    describe('getCatalog', () => {
        it('debería retornar el catálogo paginado de productos', async () => {
            vi.mocked(repo.getCatalog).mockResolvedValue(paginatedProducts);

            const result = await useCases.getCatalog('tienda-prueba', { page: 1, page_size: 10 });

            expect(result).toEqual(paginatedProducts);
            expect(repo.getCatalog).toHaveBeenCalledOnce();
            expect(repo.getCatalog).toHaveBeenCalledWith('tienda-prueba', { page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros opcionales cuando no se pasan', async () => {
            vi.mocked(repo.getCatalog).mockResolvedValue(paginatedProducts);

            await useCases.getCatalog('tienda-prueba');

            expect(repo.getCatalog).toHaveBeenCalledWith('tienda-prueba', undefined);
        });

        it('debería pasar filtros de búsqueda y categoría cuando se proporcionan', async () => {
            vi.mocked(repo.getCatalog).mockResolvedValue(paginatedProducts);

            await useCases.getCatalog('tienda-prueba', { search: 'camiseta', category: 'ropa' });

            expect(repo.getCatalog).toHaveBeenCalledWith('tienda-prueba', { search: 'camiseta', category: 'ropa' });
        });

        it('debería propagar el error cuando la consulta del catálogo falla', async () => {
            vi.mocked(repo.getCatalog).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getCatalog('tienda-prueba')).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // getProduct
    // ---------------------------------------------------------------
    describe('getProduct', () => {
        it('debería retornar un producto por su ID y slug del negocio', async () => {
            const product = makePublicProduct();
            vi.mocked(repo.getProduct).mockResolvedValue(product);

            const result = await useCases.getProduct('tienda-prueba', 'prod-001');

            expect(result).toEqual(product);
            expect(repo.getProduct).toHaveBeenCalledOnce();
            expect(repo.getProduct).toHaveBeenCalledWith('tienda-prueba', 'prod-001');
        });

        it('debería propagar el error cuando el producto no existe', async () => {
            vi.mocked(repo.getProduct).mockRejectedValue(new Error('Producto no encontrado'));

            await expect(useCases.getProduct('tienda-prueba', 'no-existe')).rejects.toThrow('Producto no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // submitContact
    // ---------------------------------------------------------------
    describe('submitContact', () => {
        const contactDto: ContactFormDTO = {
            name: 'Juan Perez',
            email: 'juan@test.com',
            message: 'Me interesa su producto',
        };

        it('debería enviar el formulario de contacto y retornar la respuesta', async () => {
            const response = { message: 'Mensaje enviado exitosamente' };
            vi.mocked(repo.submitContact).mockResolvedValue(response);

            const result = await useCases.submitContact('tienda-prueba', contactDto);

            expect(result).toEqual(response);
            expect(repo.submitContact).toHaveBeenCalledOnce();
            expect(repo.submitContact).toHaveBeenCalledWith('tienda-prueba', contactDto);
        });

        it('debería enviar el formulario con campos opcionales omitidos', async () => {
            const minimalDto: ContactFormDTO = { name: 'Ana', message: 'Hola' };
            const response = { message: 'Mensaje enviado exitosamente' };
            vi.mocked(repo.submitContact).mockResolvedValue(response);

            await useCases.submitContact('tienda-prueba', minimalDto);

            expect(repo.submitContact).toHaveBeenCalledWith('tienda-prueba', minimalDto);
        });

        it('debería propagar el error cuando el envío del formulario falla', async () => {
            vi.mocked(repo.submitContact).mockRejectedValue(new Error('Error de validación'));

            await expect(useCases.submitContact('tienda-prueba', contactDto)).rejects.toThrow('Error de validación');
        });
    });
});
