import { describe, it, expect, vi, beforeEach } from 'vitest';
import { IntegrationUseCases } from './use-cases';
import { IIntegrationRepository } from '../domain/ports';
import {
    Integration,
    IntegrationType,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeIntegration = (overrides: Partial<Integration> = {}): Integration => ({
    id: 1,
    name: 'Mi Shopify',
    code: 'shopify-1',
    integration_type_id: 1,
    type: 'shopify',
    category: 'ecommerce',
    business_id: 1,
    is_active: true,
    is_default: false,
    is_testing: false,
    config: {},
    created_by_id: 1,
    updated_by_id: null,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeIntegrationType = (overrides: Partial<IntegrationType> = {}): IntegrationType => ({
    id: 1,
    name: 'Shopify',
    code: 'shopify',
    is_active: true,
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const paginatedIntegrations: PaginatedResponse<Integration> = {
    success: true,
    message: 'OK',
    data: [makeIntegration(), makeIntegration({ id: 2, name: 'Mi WhatsApp', type: 'whatsapp' })],
    total: 2,
    page: 1,
    page_size: 10,
    total_pages: 1,
};

const singleIntegration: SingleResponse<Integration> = {
    success: true,
    message: 'OK',
    data: makeIntegration(),
};

const singleIntegrationType: SingleResponse<IntegrationType> = {
    success: true,
    message: 'OK',
    data: makeIntegrationType(),
};

const integrationTypesList: SingleResponse<IntegrationType[]> = {
    success: true,
    message: 'OK',
    data: [makeIntegrationType(), makeIntegrationType({ id: 2, name: 'WhatsApp', code: 'whatsapp' })],
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IIntegrationRepository {
    return {
        getIntegrations: vi.fn(),
        getIntegrationById: vi.fn(),
        getIntegrationByType: vi.fn(),
        createIntegration: vi.fn(),
        updateIntegration: vi.fn(),
        deleteIntegration: vi.fn(),
        testConnection: vi.fn(),
        activateIntegration: vi.fn(),
        deactivateIntegration: vi.fn(),
        setAsDefault: vi.fn(),
        syncOrders: vi.fn(),
        getSyncStatus: vi.fn(),
        testIntegration: vi.fn(),
        testConnectionRaw: vi.fn(),
        getWebhookUrl: vi.fn(),
        listWebhooks: vi.fn(),
        deleteWebhook: vi.fn(),
        verifyWebhooks: vi.fn(),
        createWebhook: vi.fn(),
        getIntegrationTypes: vi.fn(),
        getActiveIntegrationTypes: vi.fn(),
        getIntegrationTypeById: vi.fn(),
        getIntegrationTypeByCode: vi.fn(),
        createIntegrationType: vi.fn(),
        updateIntegrationType: vi.fn(),
        deleteIntegrationType: vi.fn(),
        getIntegrationTypePlatformCredentials: vi.fn(),
        getIntegrationCategories: vi.fn(),
    } as unknown as IIntegrationRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('IntegrationUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: IntegrationUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new IntegrationUseCases(repo as unknown as IIntegrationRepository);
    });

    // ---------------------------------------------------------------
    // Integrations CRUD
    // ---------------------------------------------------------------
    describe('getIntegrations', () => {
        it('debería retornar la lista paginada de integraciones', async () => {
            vi.mocked(repo.getIntegrations).mockResolvedValue(paginatedIntegrations);

            const result = await useCases.getIntegrations({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedIntegrations);
            expect(repo.getIntegrations).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getIntegrations).mockResolvedValue(paginatedIntegrations);
            await useCases.getIntegrations();
            expect(repo.getIntegrations).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getIntegrations).mockRejectedValue(new Error('DB error'));
            await expect(useCases.getIntegrations()).rejects.toThrow('DB error');
        });
    });

    describe('getIntegrationById', () => {
        it('debería retornar una integración por ID', async () => {
            vi.mocked(repo.getIntegrationById).mockResolvedValue(singleIntegration);
            const result = await useCases.getIntegrationById(1);
            expect(result).toEqual(singleIntegration);
            expect(repo.getIntegrationById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error', async () => {
            vi.mocked(repo.getIntegrationById).mockRejectedValue(new Error('Not found'));
            await expect(useCases.getIntegrationById(999)).rejects.toThrow('Not found');
        });
    });

    describe('getIntegrationByType', () => {
        it('debería retornar una integración por tipo y businessId', async () => {
            vi.mocked(repo.getIntegrationByType).mockResolvedValue(singleIntegration);
            const result = await useCases.getIntegrationByType('shopify', 1);
            expect(result).toEqual(singleIntegration);
            expect(repo.getIntegrationByType).toHaveBeenCalledWith('shopify', 1);
        });

        it('debería funcionar sin businessId', async () => {
            vi.mocked(repo.getIntegrationByType).mockResolvedValue(singleIntegration);
            await useCases.getIntegrationByType('shopify');
            expect(repo.getIntegrationByType).toHaveBeenCalledWith('shopify', undefined);
        });
    });

    describe('createIntegration', () => {
        const dto = { name: 'New', code: 'new', integration_type_id: 1, category: 'ecommerce', business_id: 1 };

        it('debería crear una integración', async () => {
            vi.mocked(repo.createIntegration).mockResolvedValue(singleIntegration);
            const result = await useCases.createIntegration(dto);
            expect(result).toEqual(singleIntegration);
            expect(repo.createIntegration).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error', async () => {
            vi.mocked(repo.createIntegration).mockRejectedValue(new Error('Duplicate'));
            await expect(useCases.createIntegration(dto)).rejects.toThrow('Duplicate');
        });
    });

    describe('updateIntegration', () => {
        it('debería actualizar una integración', async () => {
            vi.mocked(repo.updateIntegration).mockResolvedValue(singleIntegration);
            const result = await useCases.updateIntegration(1, { name: 'Updated' });
            expect(result).toEqual(singleIntegration);
            expect(repo.updateIntegration).toHaveBeenCalledWith(1, { name: 'Updated' });
        });

        it('debería propagar el error', async () => {
            vi.mocked(repo.updateIntegration).mockRejectedValue(new Error('Not found'));
            await expect(useCases.updateIntegration(99, {})).rejects.toThrow('Not found');
        });
    });

    describe('deleteIntegration', () => {
        it('debería eliminar una integración', async () => {
            vi.mocked(repo.deleteIntegration).mockResolvedValue(actionSuccess);
            const result = await useCases.deleteIntegration(1);
            expect(result).toEqual(actionSuccess);
        });

        it('debería retornar error response', async () => {
            vi.mocked(repo.deleteIntegration).mockResolvedValue(actionError);
            const result = await useCases.deleteIntegration(999);
            expect(result.success).toBe(false);
        });

        it('debería propagar excepción', async () => {
            vi.mocked(repo.deleteIntegration).mockRejectedValue(new Error('Network'));
            await expect(useCases.deleteIntegration(1)).rejects.toThrow('Network');
        });
    });

    // ---------------------------------------------------------------
    // Operations
    // ---------------------------------------------------------------
    describe('testConnection', () => {
        it('debería probar conexión', async () => {
            vi.mocked(repo.testConnection).mockResolvedValue(actionSuccess);
            const result = await useCases.testConnection(1);
            expect(result).toEqual(actionSuccess);
            expect(repo.testConnection).toHaveBeenCalledWith(1);
        });
    });

    describe('activateIntegration / deactivateIntegration', () => {
        it('debería activar', async () => {
            vi.mocked(repo.activateIntegration).mockResolvedValue(actionSuccess);
            const result = await useCases.activateIntegration(1);
            expect(result).toEqual(actionSuccess);
        });

        it('debería desactivar', async () => {
            vi.mocked(repo.deactivateIntegration).mockResolvedValue(actionSuccess);
            const result = await useCases.deactivateIntegration(1);
            expect(result).toEqual(actionSuccess);
        });
    });

    describe('setAsDefault', () => {
        it('debería establecer como default', async () => {
            vi.mocked(repo.setAsDefault).mockResolvedValue(singleIntegration);
            const result = await useCases.setAsDefault(1);
            expect(result).toEqual(singleIntegration);
        });
    });

    describe('syncOrders', () => {
        it('debería sincronizar órdenes con params', async () => {
            vi.mocked(repo.syncOrders).mockResolvedValue(actionSuccess);
            const result = await useCases.syncOrders(1, { status: 'open' });
            expect(result).toEqual(actionSuccess);
            expect(repo.syncOrders).toHaveBeenCalledWith(1, { status: 'open' });
        });
    });

    describe('getSyncStatus', () => {
        it('debería obtener estado de sincronización', async () => {
            const status = { success: true, in_progress: false };
            vi.mocked(repo.getSyncStatus).mockResolvedValue(status);
            const result = await useCases.getSyncStatus(1, 5);
            expect(result).toEqual(status);
            expect(repo.getSyncStatus).toHaveBeenCalledWith(1, 5);
        });
    });

    describe('testConnectionRaw', () => {
        it('debería probar conexión con datos crudos', async () => {
            vi.mocked(repo.testConnectionRaw).mockResolvedValue(actionSuccess);
            const result = await useCases.testConnectionRaw('shopify', { url: 'x' }, { key: 'y' });
            expect(result).toEqual(actionSuccess);
            expect(repo.testConnectionRaw).toHaveBeenCalledWith('shopify', { url: 'x' }, { key: 'y' });
        });
    });

    // ---------------------------------------------------------------
    // Integration Types
    // ---------------------------------------------------------------
    describe('getIntegrationTypes', () => {
        it('debería retornar tipos de integración', async () => {
            vi.mocked(repo.getIntegrationTypes).mockResolvedValue(integrationTypesList);
            const result = await useCases.getIntegrationTypes(1);
            expect(result).toEqual(integrationTypesList);
            expect(repo.getIntegrationTypes).toHaveBeenCalledWith(1);
        });

        it('debería funcionar sin categoryId', async () => {
            vi.mocked(repo.getIntegrationTypes).mockResolvedValue(integrationTypesList);
            await useCases.getIntegrationTypes();
            expect(repo.getIntegrationTypes).toHaveBeenCalledWith(undefined);
        });
    });

    describe('getActiveIntegrationTypes', () => {
        it('debería retornar tipos activos', async () => {
            vi.mocked(repo.getActiveIntegrationTypes).mockResolvedValue(integrationTypesList);
            const result = await useCases.getActiveIntegrationTypes();
            expect(result).toEqual(integrationTypesList);
        });
    });

    describe('getIntegrationTypeById', () => {
        it('debería retornar un tipo por ID', async () => {
            vi.mocked(repo.getIntegrationTypeById).mockResolvedValue(singleIntegrationType);
            const result = await useCases.getIntegrationTypeById(1);
            expect(result).toEqual(singleIntegrationType);
        });
    });

    describe('getIntegrationTypeByCode', () => {
        it('debería retornar un tipo por código', async () => {
            vi.mocked(repo.getIntegrationTypeByCode).mockResolvedValue(singleIntegrationType);
            const result = await useCases.getIntegrationTypeByCode('shopify');
            expect(result).toEqual(singleIntegrationType);
            expect(repo.getIntegrationTypeByCode).toHaveBeenCalledWith('shopify');
        });
    });

    describe('createIntegrationType', () => {
        it('debería crear un tipo', async () => {
            vi.mocked(repo.createIntegrationType).mockResolvedValue(singleIntegrationType);
            const result = await useCases.createIntegrationType({ name: 'New', category_id: 1 });
            expect(result).toEqual(singleIntegrationType);
        });
    });

    describe('updateIntegrationType', () => {
        it('debería actualizar un tipo', async () => {
            vi.mocked(repo.updateIntegrationType).mockResolvedValue(singleIntegrationType);
            const result = await useCases.updateIntegrationType(1, { name: 'Updated' });
            expect(result).toEqual(singleIntegrationType);
            expect(repo.updateIntegrationType).toHaveBeenCalledWith(1, { name: 'Updated' });
        });
    });

    describe('deleteIntegrationType', () => {
        it('debería eliminar un tipo', async () => {
            vi.mocked(repo.deleteIntegrationType).mockResolvedValue(actionSuccess);
            const result = await useCases.deleteIntegrationType(1);
            expect(result).toEqual(actionSuccess);
        });
    });

    describe('getIntegrationTypePlatformCredentials', () => {
        it('debería retornar credenciales de plataforma', async () => {
            const creds = { success: true, message: 'OK', data: { api_key: '***' } };
            vi.mocked(repo.getIntegrationTypePlatformCredentials).mockResolvedValue(creds);
            const result = await useCases.getIntegrationTypePlatformCredentials(1);
            expect(result).toEqual(creds);
        });
    });

    // ---------------------------------------------------------------
    // Webhooks
    // ---------------------------------------------------------------
    describe('getWebhookUrl', () => {
        it('debería retornar URL del webhook', async () => {
            const resp = { success: true, data: { url: 'https://hook.example.com', method: 'POST', description: 'Webhook' } };
            vi.mocked(repo.getWebhookUrl).mockResolvedValue(resp);
            const result = await useCases.getWebhookUrl(1);
            expect(result).toEqual(resp);
        });
    });

    describe('listWebhooks', () => {
        it('debería listar webhooks', async () => {
            const resp = { success: true, data: [] };
            vi.mocked(repo.listWebhooks).mockResolvedValue(resp);
            const result = await useCases.listWebhooks(1);
            expect(result).toEqual(resp);
        });
    });

    describe('deleteWebhook', () => {
        it('debería eliminar un webhook', async () => {
            const resp = { success: true, message: 'Deleted' };
            vi.mocked(repo.deleteWebhook).mockResolvedValue(resp);
            const result = await useCases.deleteWebhook(1, 'wh-123');
            expect(result).toEqual(resp);
            expect(repo.deleteWebhook).toHaveBeenCalledWith(1, 'wh-123');
        });
    });

    describe('verifyWebhooks', () => {
        it('debería verificar webhooks', async () => {
            const resp = { success: true, data: [], message: 'OK' };
            vi.mocked(repo.verifyWebhooks).mockResolvedValue(resp);
            const result = await useCases.verifyWebhooks(1);
            expect(result).toEqual(resp);
        });
    });

    describe('createWebhook', () => {
        it('debería crear webhook', async () => {
            const resp = { success: true, data: { existing_webhooks: [], deleted_webhooks: [], created_webhooks: ['orders'], webhook_url: 'https://x' }, message: 'OK' };
            vi.mocked(repo.createWebhook).mockResolvedValue(resp);
            const result = await useCases.createWebhook(1);
            expect(result).toEqual(resp);
        });
    });

    // ---------------------------------------------------------------
    // Categories
    // ---------------------------------------------------------------
    describe('getIntegrationCategories', () => {
        it('debería retornar categorías', async () => {
            const resp = { success: true, message: 'OK', data: [{ id: 1, code: 'ecommerce', name: 'E-commerce', display_order: 1, is_active: true, is_visible: true, created_at: '', updated_at: '' }] };
            vi.mocked(repo.getIntegrationCategories).mockResolvedValue(resp);
            const result = await useCases.getIntegrationCategories();
            expect(result).toEqual(resp);
        });
    });
});
