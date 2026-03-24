import { describe, it, expect, vi, beforeEach } from 'vitest';
import { NotificationConfigUseCases } from './use-cases';
import {
    NotificationConfig,
    CreateConfigDTO,
    UpdateConfigDTO,
    ConfigFilter,
} from '../domain/types';

// -----------------------------------------------------------------
// Mock de las server actions
// -----------------------------------------------------------------

vi.mock('../infra/actions', () => ({
    createConfigAction: vi.fn(),
    updateConfigAction: vi.fn(),
    deleteConfigAction: vi.fn(),
    listConfigsAction: vi.fn(),
    getConfigAction: vi.fn(),
}));

import {
    createConfigAction,
    updateConfigAction,
    deleteConfigAction,
    listConfigsAction,
    getConfigAction,
} from '../infra/actions';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makeConfig = (overrides: Partial<NotificationConfig> = {}): NotificationConfig => ({
    id: 1,
    business_id: 1,
    integration_id: 10,
    notification_type_id: 2,
    notification_event_type_id: 3,
    enabled: true,
    description: 'Notificacion de prueba',
    created_at: '2026-03-01T00:00:00Z',
    updated_at: '2026-03-01T00:00:00Z',
    ...overrides,
});

const makeSuccessResponse = (data: any = makeConfig()) => ({
    success: true,
    data,
});

const makeErrorResponse = (error: string = 'Algo salio mal') => ({
    success: false,
    error,
});

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('NotificationConfigUseCases', () => {
    let useCases: NotificationConfigUseCases;

    beforeEach(() => {
        vi.clearAllMocks();
        useCases = new NotificationConfigUseCases();
    });

    // ---------------------------------------------------------------
    // createConfig
    // ---------------------------------------------------------------
    describe('createConfig', () => {
        const dto: CreateConfigDTO = {
            business_id: 1,
            integration_id: 10,
            notification_type_id: 2,
            notification_event_type_id: 3,
            enabled: true,
            description: 'Nueva config',
        };

        it('debería crear una configuración y retornar la respuesta exitosa', async () => {
            const response = makeSuccessResponse(makeConfig({ description: 'Nueva config' }));
            vi.mocked(createConfigAction).mockResolvedValue(response);

            const result = await useCases.createConfig(dto);

            expect(result).toEqual(response);
            expect(createConfigAction).toHaveBeenCalledOnce();
            expect(createConfigAction).toHaveBeenCalledWith(dto);
        });

        it('debería retornar respuesta de error cuando la creación falla', async () => {
            const response = makeErrorResponse('Configuración duplicada');
            vi.mocked(createConfigAction).mockResolvedValue(response);

            const result = await useCases.createConfig(dto);

            expect(result.success).toBe(false);
            expect(result.error).toBe('Configuración duplicada');
        });

        it('debería propagar el error cuando la action lanza una excepción', async () => {
            vi.mocked(createConfigAction).mockRejectedValue(new Error('Error de red'));

            await expect(useCases.createConfig(dto)).rejects.toThrow('Error de red');
        });
    });

    // ---------------------------------------------------------------
    // updateConfig
    // ---------------------------------------------------------------
    describe('updateConfig', () => {
        const updateDto: UpdateConfigDTO = {
            enabled: false,
            description: 'Config actualizada',
        };

        it('debería actualizar una configuración y retornar la respuesta', async () => {
            const response = makeSuccessResponse(makeConfig({ enabled: false, description: 'Config actualizada' }));
            vi.mocked(updateConfigAction).mockResolvedValue(response);

            const result = await useCases.updateConfig(1, updateDto);

            expect(result).toEqual(response);
            expect(updateConfigAction).toHaveBeenCalledOnce();
            expect(updateConfigAction).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(updateConfigAction).mockRejectedValue(new Error('Config no encontrada'));

            await expect(useCases.updateConfig(99, updateDto)).rejects.toThrow('Config no encontrada');
        });
    });

    // ---------------------------------------------------------------
    // deleteConfig
    // ---------------------------------------------------------------
    describe('deleteConfig', () => {
        it('debería eliminar una configuración y retornar confirmación', async () => {
            const response = { success: true };
            vi.mocked(deleteConfigAction).mockResolvedValue(response);

            const result = await useCases.deleteConfig(1);

            expect(result).toEqual(response);
            expect(deleteConfigAction).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando la config no existe', async () => {
            const response = makeErrorResponse('Config no encontrada');
            vi.mocked(deleteConfigAction).mockResolvedValue(response);

            const result = await useCases.deleteConfig(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando la action lanza un error de red', async () => {
            vi.mocked(deleteConfigAction).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteConfig(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // listConfigs
    // ---------------------------------------------------------------
    describe('listConfigs', () => {
        it('debería retornar la lista de configuraciones con filtro', async () => {
            const configs = [makeConfig(), makeConfig({ id: 2, integration_id: 20 })];
            const response = makeSuccessResponse(configs);
            const filter: ConfigFilter = { business_id: 1, integration_id: 10 };
            vi.mocked(listConfigsAction).mockResolvedValue(response);

            const result = await useCases.listConfigs(filter);

            expect(result).toEqual(response);
            expect(listConfigsAction).toHaveBeenCalledOnce();
            expect(listConfigsAction).toHaveBeenCalledWith(filter);
        });

        it('debería llamar a la action sin filtro cuando no se pasa parámetro', async () => {
            const response = makeSuccessResponse([]);
            vi.mocked(listConfigsAction).mockResolvedValue(response);

            await useCases.listConfigs();

            expect(listConfigsAction).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(listConfigsAction).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.listConfigs()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // getConfig
    // ---------------------------------------------------------------
    describe('getConfig', () => {
        it('debería retornar una configuración por ID', async () => {
            const response = makeSuccessResponse(makeConfig());
            vi.mocked(getConfigAction).mockResolvedValue(response);

            const result = await useCases.getConfig(1);

            expect(result).toEqual(response);
            expect(getConfigAction).toHaveBeenCalledOnce();
            expect(getConfigAction).toHaveBeenCalledWith(1);
        });

        it('debería retornar error cuando la config no existe', async () => {
            const response = makeErrorResponse('No encontrada');
            vi.mocked(getConfigAction).mockResolvedValue(response);

            const result = await useCases.getConfig(999);

            expect(result.success).toBe(false);
            expect(result.error).toBe('No encontrada');
        });

        it('debería propagar el error cuando la action lanza una excepción', async () => {
            vi.mocked(getConfigAction).mockRejectedValue(new Error('Error interno'));

            await expect(useCases.getConfig(1)).rejects.toThrow('Error interno');
        });
    });
});
