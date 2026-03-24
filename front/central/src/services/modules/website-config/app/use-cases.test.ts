import { describe, it, expect, vi, beforeEach } from 'vitest';
import { WebsiteConfigUseCases } from './use-cases';
import { IWebsiteConfigRepository } from '../domain/ports';
import { WebsiteConfigData, UpdateWebsiteConfigDTO } from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeConfig = (overrides: Partial<WebsiteConfigData> = {}): WebsiteConfigData => ({
    id: 1,
    business_id: 5,
    template: 'default',
    show_hero: true,
    show_about: true,
    show_featured_products: true,
    show_full_catalog: false,
    show_testimonials: true,
    show_location: true,
    show_contact: true,
    show_social_media: false,
    show_whatsapp: true,
    hero_content: { title: 'Bienvenido', subtitle: 'Tienda online' },
    about_content: { text: 'Sobre nosotros' },
    testimonials_content: [{ author: 'Juan', text: 'Excelente' }],
    location_content: { address: 'Calle 100' },
    contact_content: { email: 'info@test.com' },
    social_media_content: null,
    whatsapp_content: { phone: '+573001234567' },
    ...overrides,
});

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IWebsiteConfigRepository {
    return {
        getConfig: vi.fn(),
        updateConfig: vi.fn(),
    };
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('WebsiteConfigUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: WebsiteConfigUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new WebsiteConfigUseCases(repo as unknown as IWebsiteConfigRepository);
    });

    // ---------------------------------------------------------------
    // getConfig
    // ---------------------------------------------------------------
    describe('getConfig', () => {
        it('debería retornar la configuración del sitio web cuando el repositorio tiene éxito', async () => {
            const config = makeConfig();
            vi.mocked(repo.getConfig).mockResolvedValue(config);

            const result = await useCases.getConfig();

            expect(result).toEqual(config);
            expect(repo.getConfig).toHaveBeenCalledOnce();
            expect(repo.getConfig).toHaveBeenCalledWith(undefined);
        });

        it('debería pasar el businessId al repositorio cuando se proporciona', async () => {
            const config = makeConfig({ business_id: 10 });
            vi.mocked(repo.getConfig).mockResolvedValue(config);

            const result = await useCases.getConfig(10);

            expect(result).toEqual(config);
            expect(repo.getConfig).toHaveBeenCalledWith(10);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            const expectedError = new Error('Configuración no encontrada');
            vi.mocked(repo.getConfig).mockRejectedValue(expectedError);

            await expect(useCases.getConfig()).rejects.toThrow('Configuración no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // updateConfig
    // ---------------------------------------------------------------
    describe('updateConfig', () => {
        const updateDto: UpdateWebsiteConfigDTO = {
            show_hero: false,
            show_about: true,
            template: 'modern',
        };

        it('debería actualizar la configuración y retornar la respuesta del repositorio', async () => {
            const updatedConfig = makeConfig({ show_hero: false, template: 'modern' });
            vi.mocked(repo.updateConfig).mockResolvedValue(updatedConfig);

            const result = await useCases.updateConfig(updateDto);

            expect(result).toEqual(updatedConfig);
            expect(repo.updateConfig).toHaveBeenCalledOnce();
            expect(repo.updateConfig).toHaveBeenCalledWith(updateDto, undefined);
        });

        it('debería pasar el businessId al repositorio cuando se proporciona', async () => {
            const updatedConfig = makeConfig({ show_hero: false, business_id: 7 });
            vi.mocked(repo.updateConfig).mockResolvedValue(updatedConfig);

            const result = await useCases.updateConfig(updateDto, 7);

            expect(result).toEqual(updatedConfig);
            expect(repo.updateConfig).toHaveBeenCalledWith(updateDto, 7);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            const expectedError = new Error('Error de validación');
            vi.mocked(repo.updateConfig).mockRejectedValue(expectedError);

            await expect(useCases.updateConfig(updateDto)).rejects.toThrow('Error de validación');
        });

        it('debería permitir actualizar secciones individuales', async () => {
            const partialDto: UpdateWebsiteConfigDTO = { show_whatsapp: true };
            const updatedConfig = makeConfig({ show_whatsapp: true });
            vi.mocked(repo.updateConfig).mockResolvedValue(updatedConfig);

            const result = await useCases.updateConfig(partialDto, 5);

            expect(result).toEqual(updatedConfig);
            expect(repo.updateConfig).toHaveBeenCalledWith(partialDto, 5);
        });

        it('debería permitir actualizar contenido de secciones', async () => {
            const contentDto: UpdateWebsiteConfigDTO = {
                hero_content: { title: 'Nuevo título', subtitle: 'Nueva descripción' },
                testimonials_content: [
                    { author: 'María', text: 'Muy bueno' },
                    { author: 'Pedro', text: 'Excelente servicio' },
                ],
            };
            const updatedConfig = makeConfig({
                hero_content: contentDto.hero_content!,
                testimonials_content: contentDto.testimonials_content!,
            });
            vi.mocked(repo.updateConfig).mockResolvedValue(updatedConfig);

            const result = await useCases.updateConfig(contentDto);

            expect(result).toEqual(updatedConfig);
            expect(repo.updateConfig).toHaveBeenCalledWith(contentDto, undefined);
        });
    });
});
