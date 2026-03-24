import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PermissionUseCases } from './use-cases';
import { IPermissionRepository } from '../domain/ports';
import {
    Permission,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    BulkCreatePermissionsResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makePermission = (overrides: Partial<Permission> = {}): Permission => ({
    id: 1,
    name: 'read:orders',
    code: 'read_orders',
    description: 'Leer pedidos',
    resource: 'orders',
    action: 'read',
    resource_id: 1,
    action_id: 1,
    scope_id: 1,
    scope_name: 'Platform',
    scope_code: 'platform',
    business_type_id: 1,
    business_type_name: 'Restaurante',
    ...overrides,
});

const paginatedPermissions: PaginatedResponse<Permission> = {
    success: true,
    data: [makePermission(), makePermission({ id: 2, name: 'write:orders' })],
    total: 2,
};

const singlePermission: SingleResponse<Permission> = {
    success: true,
    data: makePermission(),
    message: 'OK',
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

const bulkSuccess: BulkCreatePermissionsResponse = {
    success: true,
    message: 'OK',
    results: [
        { name: 'read:orders', success: true },
        { name: 'write:orders', success: true },
    ],
};

const bulkPartial: BulkCreatePermissionsResponse = {
    success: true,
    message: 'Parcialmente completado',
    results: [
        { name: 'read:orders', success: true },
        { name: 'write:orders', success: false, error: 'Ya existe' },
    ],
};

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IPermissionRepository {
    return {
        getPermissions: vi.fn(),
        getPermissionById: vi.fn(),
        getPermissionsByScope: vi.fn(),
        getPermissionsByResource: vi.fn(),
        createPermission: vi.fn(),
        updatePermission: vi.fn(),
        deletePermission: vi.fn(),
        createPermissionsBulk: vi.fn(),
    } as unknown as IPermissionRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('PermissionUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: PermissionUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new PermissionUseCases(repo as unknown as IPermissionRepository);
    });

    // ---------------------------------------------------------------
    // getPermissions
    // ---------------------------------------------------------------
    describe('getPermissions', () => {
        it('debería retornar la lista de permisos cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getPermissions).mockResolvedValue(paginatedPermissions);

            const result = await useCases.getPermissions({ scope_id: 1 });

            expect(result).toEqual(paginatedPermissions);
            expect(repo.getPermissions).toHaveBeenCalledOnce();
            expect(repo.getPermissions).toHaveBeenCalledWith({ scope_id: 1 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getPermissions).mockResolvedValue(paginatedPermissions);

            await useCases.getPermissions();

            expect(repo.getPermissions).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getPermissions).mockRejectedValue(new Error('Fallo de base de datos'));

            await expect(useCases.getPermissions()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getPermissionById
    // ---------------------------------------------------------------
    describe('getPermissionById', () => {
        it('debería retornar un permiso por ID', async () => {
            vi.mocked(repo.getPermissionById).mockResolvedValue(singlePermission);

            const result = await useCases.getPermissionById(1);

            expect(result).toEqual(singlePermission);
            expect(repo.getPermissionById).toHaveBeenCalledOnce();
            expect(repo.getPermissionById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando el permiso no existe', async () => {
            vi.mocked(repo.getPermissionById).mockRejectedValue(new Error('Permiso no encontrado'));

            await expect(useCases.getPermissionById(999)).rejects.toThrow('Permiso no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getPermissionsByScope
    // ---------------------------------------------------------------
    describe('getPermissionsByScope', () => {
        it('debería retornar permisos filtrados por scope', async () => {
            vi.mocked(repo.getPermissionsByScope).mockResolvedValue(paginatedPermissions);

            const result = await useCases.getPermissionsByScope(1);

            expect(result).toEqual(paginatedPermissions);
            expect(repo.getPermissionsByScope).toHaveBeenCalledOnce();
            expect(repo.getPermissionsByScope).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getPermissionsByScope).mockRejectedValue(new Error('Scope no encontrado'));

            await expect(useCases.getPermissionsByScope(999)).rejects.toThrow('Scope no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getPermissionsByResource
    // ---------------------------------------------------------------
    describe('getPermissionsByResource', () => {
        it('debería retornar permisos filtrados por recurso', async () => {
            vi.mocked(repo.getPermissionsByResource).mockResolvedValue(paginatedPermissions);

            const result = await useCases.getPermissionsByResource('orders');

            expect(result).toEqual(paginatedPermissions);
            expect(repo.getPermissionsByResource).toHaveBeenCalledOnce();
            expect(repo.getPermissionsByResource).toHaveBeenCalledWith('orders');
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getPermissionsByResource).mockRejectedValue(new Error('Recurso no encontrado'));

            await expect(useCases.getPermissionsByResource('invalid')).rejects.toThrow('Recurso no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createPermission
    // ---------------------------------------------------------------
    describe('createPermission', () => {
        const dto = { name: 'delete:orders', resource_id: 1, action_id: 3, scope_id: 1 };

        it('debería crear un permiso y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createPermission).mockResolvedValue(singlePermission);

            const result = await useCases.createPermission(dto);

            expect(result).toEqual(singlePermission);
            expect(repo.createPermission).toHaveBeenCalledOnce();
            expect(repo.createPermission).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createPermission).mockRejectedValue(new Error('Nombre duplicado'));

            await expect(useCases.createPermission(dto)).rejects.toThrow('Nombre duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updatePermission
    // ---------------------------------------------------------------
    describe('updatePermission', () => {
        const updateDto = {
            name: 'read:orders_updated',
            code: 'read_orders_updated',
            description: 'Lectura actualizada',
            resource_id: 1,
            action_id: 1,
            scope_id: 1,
        };

        it('debería actualizar un permiso y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<Permission> = {
                ...singlePermission,
                data: makePermission({ name: 'read:orders_updated' }),
            };
            vi.mocked(repo.updatePermission).mockResolvedValue(updatedResponse);

            const result = await useCases.updatePermission(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updatePermission).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updatePermission).mockRejectedValue(new Error('Permiso no encontrado'));

            await expect(useCases.updatePermission(99, updateDto)).rejects.toThrow('Permiso no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deletePermission
    // ---------------------------------------------------------------
    describe('deletePermission', () => {
        it('debería eliminar un permiso y retornar confirmación', async () => {
            vi.mocked(repo.deletePermission).mockResolvedValue(actionSuccess);

            const result = await useCases.deletePermission(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deletePermission).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el permiso no existe', async () => {
            vi.mocked(repo.deletePermission).mockResolvedValue(actionError);

            const result = await useCases.deletePermission(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deletePermission).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deletePermission(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // createPermissionsBulk
    // ---------------------------------------------------------------
    describe('createPermissionsBulk', () => {
        const bulkData = [
            { name: 'read:orders', resource_id: 1, action_id: 1, scope_id: 1 },
            { name: 'write:orders', resource_id: 1, action_id: 2, scope_id: 1 },
        ];

        it('debería crear permisos en lote y retornar los resultados', async () => {
            vi.mocked(repo.createPermissionsBulk).mockResolvedValue(bulkSuccess);

            const result = await useCases.createPermissionsBulk(bulkData);

            expect(result).toEqual(bulkSuccess);
            expect(repo.createPermissionsBulk).toHaveBeenCalledOnce();
            expect(repo.createPermissionsBulk).toHaveBeenCalledWith(bulkData);
        });

        it('debería manejar resultados parciales cuando algunos permisos fallan', async () => {
            vi.mocked(repo.createPermissionsBulk).mockResolvedValue(bulkPartial);

            const result = await useCases.createPermissionsBulk(bulkData);

            expect(result.success).toBe(true);
            expect(result.results[0].success).toBe(true);
            expect(result.results[1].success).toBe(false);
            expect(result.results[1].error).toBe('Ya existe');
        });

        it('debería propagar el error cuando la creación en lote falla', async () => {
            vi.mocked(repo.createPermissionsBulk).mockRejectedValue(new Error('Error de servidor'));

            await expect(useCases.createPermissionsBulk(bulkData)).rejects.toThrow('Error de servidor');
        });
    });
});
