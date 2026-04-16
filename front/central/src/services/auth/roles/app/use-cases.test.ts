import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RoleUseCases } from './use-cases';
import { IRoleRepository } from '../domain/ports';
import {
    Role,
    Permission,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    RolePermissionsResponse,
    AssignPermissionsResponse,
} from '../domain/types';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const makePagination = () => ({
    current_page: 1,
    per_page: 10,
    total: 1,
    last_page: 1,
    has_next: false,
    has_prev: false,
});

const makeRole = (overrides: Partial<Role> = {}): Role => ({
    id: 1,
    name: 'Admin',
    code: 'admin',
    description: 'Administrador',
    level: 1,
    is_system: false,
    scope_id: 1,
    scope_name: 'Platform',
    scope_code: 'platform',
    business_type_id: 1,
    business_type_name: 'Restaurante',
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const makeRolePermission = (overrides: Partial<Permission> = {}): Permission => ({
    id: 1,
    resource: 'orders',
    action: 'read',
    description: 'Leer pedidos',
    scope_id: 1,
    scope_name: 'Platform',
    scope_code: 'platform',
    ...overrides,
});

const paginatedRoles: PaginatedResponse<Role> = {
    success: true,
    message: 'OK',
    data: [makeRole(), makeRole({ id: 2, name: 'Editor' })],
    pagination: makePagination(),
};

const singleRole: SingleResponse<Role> = {
    success: true,
    data: makeRole(),
    message: 'OK',
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

const rolePermissionsResponse: RolePermissionsResponse = {
    success: true,
    message: 'OK',
    role_id: 1,
    role_name: 'Admin',
    permissions: [makeRolePermission(), makeRolePermission({ id: 2, action: 'write' })],
    count: 2,
};

const assignPermissionsResponse: AssignPermissionsResponse = {
    success: true,
    message: 'OK',
    role_id: 1,
    permission_ids: [1, 2, 3],
};

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IRoleRepository {
    return {
        getRoles: vi.fn(),
        getRoleById: vi.fn(),
        getRolesByScope: vi.fn(),
        getRolesByLevel: vi.fn(),
        getSystemRoles: vi.fn(),
        createRole: vi.fn(),
        updateRole: vi.fn(),
        deleteRole: vi.fn(),
        assignPermissions: vi.fn(),
        getRolePermissions: vi.fn(),
        removePermissionFromRole: vi.fn(),
    } as unknown as IRoleRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('RoleUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: RoleUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new RoleUseCases(repo as unknown as IRoleRepository);
    });

    // ---------------------------------------------------------------
    // getRoles
    // ---------------------------------------------------------------
    describe('getRoles', () => {
        it('debería retornar la lista paginada de roles cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getRoles).mockResolvedValue(paginatedRoles);

            const result = await useCases.getRoles({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedRoles);
            expect(repo.getRoles).toHaveBeenCalledOnce();
            expect(repo.getRoles).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getRoles).mockResolvedValue(paginatedRoles);

            await useCases.getRoles();

            expect(repo.getRoles).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getRoles).mockRejectedValue(new Error('Fallo de base de datos'));

            await expect(useCases.getRoles()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getRoleById
    // ---------------------------------------------------------------
    describe('getRoleById', () => {
        it('debería retornar un rol por ID', async () => {
            vi.mocked(repo.getRoleById).mockResolvedValue(singleRole);

            const result = await useCases.getRoleById(1);

            expect(result).toEqual(singleRole);
            expect(repo.getRoleById).toHaveBeenCalledOnce();
            expect(repo.getRoleById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando el rol no existe', async () => {
            vi.mocked(repo.getRoleById).mockRejectedValue(new Error('Rol no encontrado'));

            await expect(useCases.getRoleById(999)).rejects.toThrow('Rol no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getRolesByScope
    // ---------------------------------------------------------------
    describe('getRolesByScope', () => {
        it('debería retornar roles filtrados por scope', async () => {
            vi.mocked(repo.getRolesByScope).mockResolvedValue(paginatedRoles);

            const result = await useCases.getRolesByScope(1);

            expect(result).toEqual(paginatedRoles);
            expect(repo.getRolesByScope).toHaveBeenCalledOnce();
            expect(repo.getRolesByScope).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getRolesByScope).mockRejectedValue(new Error('Scope no encontrado'));

            await expect(useCases.getRolesByScope(999)).rejects.toThrow('Scope no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getRolesByLevel
    // ---------------------------------------------------------------
    describe('getRolesByLevel', () => {
        it('debería retornar roles filtrados por nivel', async () => {
            vi.mocked(repo.getRolesByLevel).mockResolvedValue(paginatedRoles);

            const result = await useCases.getRolesByLevel(1);

            expect(result).toEqual(paginatedRoles);
            expect(repo.getRolesByLevel).toHaveBeenCalledOnce();
            expect(repo.getRolesByLevel).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getRolesByLevel).mockRejectedValue(new Error('Nivel inválido'));

            await expect(useCases.getRolesByLevel(-1)).rejects.toThrow('Nivel inválido');
        });
    });

    // ---------------------------------------------------------------
    // getSystemRoles
    // ---------------------------------------------------------------
    describe('getSystemRoles', () => {
        it('debería retornar los roles del sistema', async () => {
            const systemRoles = {
                ...paginatedRoles,
                data: [makeRole({ is_system: true })],
            };
            vi.mocked(repo.getSystemRoles).mockResolvedValue(systemRoles);

            const result = await useCases.getSystemRoles();

            expect(result).toEqual(systemRoles);
            expect(repo.getSystemRoles).toHaveBeenCalledOnce();
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getSystemRoles).mockRejectedValue(new Error('Servicio no disponible'));

            await expect(useCases.getSystemRoles()).rejects.toThrow('Servicio no disponible');
        });
    });

    // ---------------------------------------------------------------
    // createRole
    // ---------------------------------------------------------------
    describe('createRole', () => {
        const dto = {
            name: 'Viewer',
            description: 'Solo lectura',
            level: 3,
            is_system: false,
            scope_id: 1,
            business_type_id: 1,
        };

        it('debería crear un rol y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createRole).mockResolvedValue(singleRole);

            const result = await useCases.createRole(dto);

            expect(result).toEqual(singleRole);
            expect(repo.createRole).toHaveBeenCalledOnce();
            expect(repo.createRole).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createRole).mockRejectedValue(new Error('Nombre duplicado'));

            await expect(useCases.createRole(dto)).rejects.toThrow('Nombre duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateRole
    // ---------------------------------------------------------------
    describe('updateRole', () => {
        const updateDto = { name: 'Admin Actualizado' };

        it('debería actualizar un rol y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<Role> = {
                ...singleRole,
                data: makeRole({ name: 'Admin Actualizado' }),
            };
            vi.mocked(repo.updateRole).mockResolvedValue(updatedResponse);

            const result = await useCases.updateRole(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateRole).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateRole).mockRejectedValue(new Error('Rol no encontrado'));

            await expect(useCases.updateRole(99, updateDto)).rejects.toThrow('Rol no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteRole
    // ---------------------------------------------------------------
    describe('deleteRole', () => {
        it('debería eliminar un rol y retornar confirmación', async () => {
            vi.mocked(repo.deleteRole).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteRole(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteRole).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el rol no existe', async () => {
            vi.mocked(repo.deleteRole).mockResolvedValue(actionError);

            const result = await useCases.deleteRole(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteRole).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteRole(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // assignPermissions
    // ---------------------------------------------------------------
    describe('assignPermissions', () => {
        const assignDto = { permission_ids: [1, 2, 3] };

        it('debería asignar permisos a un rol y retornar la respuesta', async () => {
            vi.mocked(repo.assignPermissions).mockResolvedValue(assignPermissionsResponse);

            const result = await useCases.assignPermissions(1, assignDto);

            expect(result).toEqual(assignPermissionsResponse);
            expect(repo.assignPermissions).toHaveBeenCalledOnce();
            expect(repo.assignPermissions).toHaveBeenCalledWith(1, assignDto);
        });

        it('debería propagar el error cuando la asignación falla', async () => {
            vi.mocked(repo.assignPermissions).mockRejectedValue(new Error('Permiso no encontrado'));

            await expect(useCases.assignPermissions(1, assignDto)).rejects.toThrow('Permiso no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getRolePermissions
    // ---------------------------------------------------------------
    describe('getRolePermissions', () => {
        it('debería retornar los permisos de un rol', async () => {
            vi.mocked(repo.getRolePermissions).mockResolvedValue(rolePermissionsResponse);

            const result = await useCases.getRolePermissions(1);

            expect(result).toEqual(rolePermissionsResponse);
            expect(repo.getRolePermissions).toHaveBeenCalledOnce();
            expect(repo.getRolePermissions).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getRolePermissions).mockRejectedValue(new Error('Rol no encontrado'));

            await expect(useCases.getRolePermissions(999)).rejects.toThrow('Rol no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // removePermissionFromRole
    // ---------------------------------------------------------------
    describe('removePermissionFromRole', () => {
        it('debería eliminar un permiso del rol y retornar confirmación', async () => {
            vi.mocked(repo.removePermissionFromRole).mockResolvedValue(actionSuccess);

            const result = await useCases.removePermissionFromRole(1, 5);

            expect(result).toEqual(actionSuccess);
            expect(repo.removePermissionFromRole).toHaveBeenCalledOnce();
            expect(repo.removePermissionFromRole).toHaveBeenCalledWith(1, 5);
        });

        it('debería retornar respuesta de error cuando no se encuentra', async () => {
            vi.mocked(repo.removePermissionFromRole).mockResolvedValue(actionError);

            const result = await useCases.removePermissionFromRole(1, 999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.removePermissionFromRole).mockRejectedValue(new Error('Network error'));

            await expect(useCases.removePermissionFromRole(1, 5)).rejects.toThrow('Network error');
        });
    });
});
