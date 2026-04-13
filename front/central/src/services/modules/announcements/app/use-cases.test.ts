import { describe, it, expect, vi, beforeEach } from 'vitest';
import { AnnouncementUseCases } from './use-cases';
import { IAnnouncementRepository } from '../domain/ports';
import {
    AnnouncementInfo,
    AnnouncementStats,
    AnnouncementCategory,
    AnnouncementsListResponse,
    DeleteAnnouncementResponse,
    CreateAnnouncementDTO,
    UpdateAnnouncementDTO,
    RegisterViewDTO,
    ChangeStatusDTO,
} from '../domain/types';

const makeAnnouncement = (overrides: Partial<AnnouncementInfo> = {}): AnnouncementInfo => ({
    id: 1,
    business_id: null,
    category_id: 1,
    category: { id: 1, code: 'informative', name: 'Informativo', icon: 'info', color: '#3b82f6' },
    title: 'Test Announcement',
    message: 'Test message',
    display_type: 'modal_text',
    frequency_type: 'once',
    priority: 0,
    is_global: true,
    status: 'active',
    starts_at: null,
    ends_at: null,
    force_redisplay: false,
    created_by_id: 1,
    created_at: '2026-04-01T00:00:00Z',
    updated_at: '2026-04-01T00:00:00Z',
    images: [],
    links: [],
    targets: [],
    ...overrides,
});

const makeStats = (overrides: Partial<AnnouncementStats> = {}): AnnouncementStats => ({
    total_views: 100,
    unique_users: 50,
    total_clicks: 25,
    total_acceptances: 10,
    total_closed: 40,
    ...overrides,
});

const makeCategory = (overrides: Partial<AnnouncementCategory> = {}): AnnouncementCategory => ({
    id: 1,
    code: 'promotion',
    name: 'Promocion',
    icon: 'tag',
    color: '#10b981',
    ...overrides,
});

const listResponse: AnnouncementsListResponse = {
    data: [makeAnnouncement()],
    total: 1,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const deleteResponse: DeleteAnnouncementResponse = { success: true, message: 'Eliminado' };

function createMockRepository(): IAnnouncementRepository {
    return {
        getAnnouncements: vi.fn(),
        getAnnouncementById: vi.fn(),
        createAnnouncement: vi.fn(),
        updateAnnouncement: vi.fn(),
        deleteAnnouncement: vi.fn(),
        getActiveAnnouncements: vi.fn(),
        registerView: vi.fn(),
        getStats: vi.fn(),
        listCategories: vi.fn(),
        changeStatus: vi.fn(),
        forceRedisplay: vi.fn(),
        uploadImage: vi.fn(),
        deleteImage: vi.fn(),
    };
}

describe('AnnouncementUseCases', () => {
    let repo: IAnnouncementRepository;
    let useCases: AnnouncementUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new AnnouncementUseCases(repo);
    });

    describe('getAnnouncements', () => {
        it('retorna lista paginada de anuncios', async () => {
            vi.mocked(repo.getAnnouncements).mockResolvedValue(listResponse);

            const result = await useCases.getAnnouncements({ page: 1, page_size: 10 });

            expect(result).toEqual(listResponse);
            expect(repo.getAnnouncements).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('llama sin parametros', async () => {
            vi.mocked(repo.getAnnouncements).mockResolvedValue(listResponse);

            await useCases.getAnnouncements();

            expect(repo.getAnnouncements).toHaveBeenCalledWith(undefined);
        });

        it('pasa filtros de status y busqueda', async () => {
            vi.mocked(repo.getAnnouncements).mockResolvedValue(listResponse);

            await useCases.getAnnouncements({ status: 'active', search: 'promo' });

            expect(repo.getAnnouncements).toHaveBeenCalledWith({ status: 'active', search: 'promo' });
        });

        it('propaga error del repositorio', async () => {
            vi.mocked(repo.getAnnouncements).mockRejectedValue(new Error('DB error'));

            await expect(useCases.getAnnouncements()).rejects.toThrow('DB error');
        });
    });

    describe('getAnnouncementById', () => {
        it('retorna un anuncio por ID', async () => {
            const announcement = makeAnnouncement();
            vi.mocked(repo.getAnnouncementById).mockResolvedValue(announcement);

            const result = await useCases.getAnnouncementById(1);

            expect(result).toEqual(announcement);
            expect(repo.getAnnouncementById).toHaveBeenCalledWith(1);
        });

        it('propaga error cuando no existe', async () => {
            vi.mocked(repo.getAnnouncementById).mockRejectedValue(new Error('Not found'));

            await expect(useCases.getAnnouncementById(999)).rejects.toThrow('Not found');
        });
    });

    describe('createAnnouncement', () => {
        const dto: CreateAnnouncementDTO = {
            category_id: 1,
            title: 'Nuevo',
            message: 'Contenido',
            display_type: 'modal_text',
            frequency_type: 'once',
            priority: 0,
            is_global: true,
            links: [],
            target_ids: [],
        };

        it('crea un anuncio', async () => {
            const created = makeAnnouncement({ id: 2, title: 'Nuevo' });
            vi.mocked(repo.createAnnouncement).mockResolvedValue(created);

            const result = await useCases.createAnnouncement(dto);

            expect(result).toEqual(created);
            expect(repo.createAnnouncement).toHaveBeenCalledWith(dto);
        });

        it('propaga error de creacion', async () => {
            vi.mocked(repo.createAnnouncement).mockRejectedValue(new Error('Validation error'));

            await expect(useCases.createAnnouncement(dto)).rejects.toThrow('Validation error');
        });
    });

    describe('updateAnnouncement', () => {
        const dto: UpdateAnnouncementDTO = {
            category_id: 1,
            title: 'Actualizado',
            message: 'Nuevo contenido',
            display_type: 'modal_text',
            frequency_type: 'daily',
            priority: 5,
            is_global: true,
            links: [],
            target_ids: [],
        };

        it('actualiza un anuncio', async () => {
            const updated = makeAnnouncement({ title: 'Actualizado' });
            vi.mocked(repo.updateAnnouncement).mockResolvedValue(updated);

            const result = await useCases.updateAnnouncement(1, dto);

            expect(result).toEqual(updated);
            expect(repo.updateAnnouncement).toHaveBeenCalledWith(1, dto);
        });

        it('propaga error', async () => {
            vi.mocked(repo.updateAnnouncement).mockRejectedValue(new Error('Not found'));

            await expect(useCases.updateAnnouncement(99, dto)).rejects.toThrow('Not found');
        });
    });

    describe('deleteAnnouncement', () => {
        it('elimina un anuncio', async () => {
            vi.mocked(repo.deleteAnnouncement).mockResolvedValue(deleteResponse);

            const result = await useCases.deleteAnnouncement(1);

            expect(result).toEqual(deleteResponse);
            expect(repo.deleteAnnouncement).toHaveBeenCalledWith(1);
        });

        it('propaga error', async () => {
            vi.mocked(repo.deleteAnnouncement).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteAnnouncement(1)).rejects.toThrow('Network error');
        });
    });

    describe('getActiveAnnouncements', () => {
        it('retorna anuncios activos sin businessId', async () => {
            const actives = [makeAnnouncement({ status: 'active' })];
            vi.mocked(repo.getActiveAnnouncements).mockResolvedValue(actives);

            const result = await useCases.getActiveAnnouncements();

            expect(result).toEqual(actives);
            expect(repo.getActiveAnnouncements).toHaveBeenCalledWith(undefined);
        });

        it('pasa businessId', async () => {
            vi.mocked(repo.getActiveAnnouncements).mockResolvedValue([]);

            await useCases.getActiveAnnouncements(5);

            expect(repo.getActiveAnnouncements).toHaveBeenCalledWith(5);
        });
    });

    describe('registerView', () => {
        const viewDto: RegisterViewDTO = { action: 'viewed' };

        it('registra una vista', async () => {
            vi.mocked(repo.registerView).mockResolvedValue(undefined);

            await useCases.registerView(1, viewDto);

            expect(repo.registerView).toHaveBeenCalledWith(1, viewDto);
        });

        it('registra click en link', async () => {
            const clickDto: RegisterViewDTO = { action: 'clicked_link', link_id: 5 };
            vi.mocked(repo.registerView).mockResolvedValue(undefined);

            await useCases.registerView(1, clickDto);

            expect(repo.registerView).toHaveBeenCalledWith(1, clickDto);
        });
    });

    describe('getStats', () => {
        it('retorna estadisticas de un anuncio', async () => {
            const stats = makeStats();
            vi.mocked(repo.getStats).mockResolvedValue(stats);

            const result = await useCases.getStats(1);

            expect(result).toEqual(stats);
            expect(repo.getStats).toHaveBeenCalledWith(1);
        });
    });

    describe('listCategories', () => {
        it('retorna lista de categorias', async () => {
            const cats = [makeCategory(), makeCategory({ id: 2, code: 'alert', name: 'Alerta' })];
            vi.mocked(repo.listCategories).mockResolvedValue(cats);

            const result = await useCases.listCategories();

            expect(result).toEqual(cats);
            expect(result).toHaveLength(2);
        });
    });

    describe('changeStatus', () => {
        const dto: ChangeStatusDTO = { status: 'inactive' };

        it('cambia el estado de un anuncio', async () => {
            vi.mocked(repo.changeStatus).mockResolvedValue(undefined);

            await useCases.changeStatus(1, dto);

            expect(repo.changeStatus).toHaveBeenCalledWith(1, dto);
        });

        it('propaga error', async () => {
            vi.mocked(repo.changeStatus).mockRejectedValue(new Error('Invalid status'));

            await expect(useCases.changeStatus(1, dto)).rejects.toThrow('Invalid status');
        });
    });

    describe('forceRedisplay', () => {
        it('fuerza re-visualizacion', async () => {
            vi.mocked(repo.forceRedisplay).mockResolvedValue(undefined);

            await useCases.forceRedisplay(1);

            expect(repo.forceRedisplay).toHaveBeenCalledWith(1);
        });

        it('propaga error', async () => {
            vi.mocked(repo.forceRedisplay).mockRejectedValue(new Error('Not found'));

            await expect(useCases.forceRedisplay(999)).rejects.toThrow('Not found');
        });
    });

    describe('uploadImage', () => {
        it('sube imagen al anuncio', async () => {
            const uploadResponse = { success: true, data: { id: 1, image_url: 'https://s3/img.png', sort_order: 0 } };
            const formData = new FormData();
            vi.mocked(repo.uploadImage).mockResolvedValue(uploadResponse);

            const result = await useCases.uploadImage(1, formData);

            expect(repo.uploadImage).toHaveBeenCalledWith(1, formData);
            expect(result).toEqual(uploadResponse);
        });

        it('propaga error', async () => {
            vi.mocked(repo.uploadImage).mockRejectedValue(new Error('Upload failed'));

            await expect(useCases.uploadImage(1, new FormData())).rejects.toThrow('Upload failed');
        });
    });

    describe('deleteImage', () => {
        it('elimina imagen del anuncio', async () => {
            const deleteImgResponse = { success: true, message: 'image deleted' };
            vi.mocked(repo.deleteImage).mockResolvedValue(deleteImgResponse);

            const result = await useCases.deleteImage(1, 5);

            expect(repo.deleteImage).toHaveBeenCalledWith(1, 5);
            expect(result).toEqual(deleteImgResponse);
        });

        it('propaga error', async () => {
            vi.mocked(repo.deleteImage).mockRejectedValue(new Error('Not found'));

            await expect(useCases.deleteImage(1, 999)).rejects.toThrow('Not found');
        });
    });
});
