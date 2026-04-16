import { describe, it, expect, vi, beforeEach } from 'vitest';
import { UserUseCases } from './use-cases';
import { IUserRepository } from '../domain/ports';
import {
    User,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
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

const makeUser = (overrides: Partial<User> = {}): User => ({
    id: 1,
    name: 'Juan Pérez',
    email: 'juan@test.com',
    phone: '3001234567',
    is_active: true,
    is_super_user: false,
    scope_id: 2,
    scope_code: 'business',
    scope_name: 'Business',
    business_role_assignments: [],
    created_at: '2026-01-01T00:00:00Z',
    updated_at: '2026-01-01T00:00:00Z',
    ...overrides,
});

const paginatedUsers: PaginatedResponse<User> = {
    success: true,
    data: [makeUser(), makeUser({ id: 2, name: 'María López' })],
    pagination: makePagination(),
    message: 'OK',
};

const singleUser: SingleResponse<User> = {
    success: true,
    data: makeUser(),
    message: 'OK',
};

const actionSuccess: ActionResponse = { success: true, message: 'OK' };
const actionError: ActionResponse = { success: false, message: 'Error', error: 'Something went wrong' };

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): IUserRepository {
    return {
        getUsers: vi.fn(),
        getUserById: vi.fn(),
        createUser: vi.fn(),
        updateUser: vi.fn(),
        deleteUser: vi.fn(),
        assignRoles: vi.fn(),
    } as unknown as IUserRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('UserUseCases', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: UserUseCases;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new UserUseCases(repo as unknown as IUserRepository);
    });

    // ---------------------------------------------------------------
    // getUsers
    // ---------------------------------------------------------------
    describe('getUsers', () => {
        it('debería retornar la lista paginada de usuarios cuando el repositorio tiene éxito', async () => {
            vi.mocked(repo.getUsers).mockResolvedValue(paginatedUsers);

            const result = await useCases.getUsers({ page: 1, page_size: 10 });

            expect(result).toEqual(paginatedUsers);
            expect(repo.getUsers).toHaveBeenCalledOnce();
            expect(repo.getUsers).toHaveBeenCalledWith({ page: 1, page_size: 10 });
        });

        it('debería llamar al repositorio sin parámetros cuando no se pasan filtros', async () => {
            vi.mocked(repo.getUsers).mockResolvedValue(paginatedUsers);

            await useCases.getUsers();

            expect(repo.getUsers).toHaveBeenCalledWith(undefined);
        });

        it('debería propagar el error cuando el repositorio falla', async () => {
            vi.mocked(repo.getUsers).mockRejectedValue(new Error('Fallo de base de datos'));

            await expect(useCases.getUsers()).rejects.toThrow('Fallo de base de datos');
        });
    });

    // ---------------------------------------------------------------
    // getUserById
    // ---------------------------------------------------------------
    describe('getUserById', () => {
        it('debería retornar un usuario por ID', async () => {
            vi.mocked(repo.getUserById).mockResolvedValue(singleUser);

            const result = await useCases.getUserById(1);

            expect(result).toEqual(singleUser);
            expect(repo.getUserById).toHaveBeenCalledOnce();
            expect(repo.getUserById).toHaveBeenCalledWith(1);
        });

        it('debería propagar el error cuando el usuario no existe', async () => {
            vi.mocked(repo.getUserById).mockRejectedValue(new Error('Usuario no encontrado'));

            await expect(useCases.getUserById(999)).rejects.toThrow('Usuario no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // createUser
    // ---------------------------------------------------------------
    describe('createUser', () => {
        const dto = { name: 'Nuevo Usuario', email: 'nuevo@test.com', scope_id: 2 };

        it('debería crear un usuario y retornar la respuesta del repositorio', async () => {
            vi.mocked(repo.createUser).mockResolvedValue(singleUser);

            const result = await useCases.createUser(dto);

            expect(result).toEqual(singleUser);
            expect(repo.createUser).toHaveBeenCalledOnce();
            expect(repo.createUser).toHaveBeenCalledWith(dto);
        });

        it('debería propagar el error cuando la creación falla', async () => {
            vi.mocked(repo.createUser).mockRejectedValue(new Error('Email duplicado'));

            await expect(useCases.createUser(dto)).rejects.toThrow('Email duplicado');
        });
    });

    // ---------------------------------------------------------------
    // updateUser
    // ---------------------------------------------------------------
    describe('updateUser', () => {
        const updateDto = { name: 'Juan Actualizado' };

        it('debería actualizar un usuario y retornar la respuesta del repositorio', async () => {
            const updatedResponse: SingleResponse<User> = {
                ...singleUser,
                data: makeUser({ name: 'Juan Actualizado' }),
            };
            vi.mocked(repo.updateUser).mockResolvedValue(updatedResponse);

            const result = await useCases.updateUser(1, updateDto);

            expect(result).toEqual(updatedResponse);
            expect(repo.updateUser).toHaveBeenCalledWith(1, updateDto);
        });

        it('debería propagar el error cuando la actualización falla', async () => {
            vi.mocked(repo.updateUser).mockRejectedValue(new Error('Usuario no encontrado'));

            await expect(useCases.updateUser(99, updateDto)).rejects.toThrow('Usuario no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // deleteUser
    // ---------------------------------------------------------------
    describe('deleteUser', () => {
        it('debería eliminar un usuario y retornar confirmación', async () => {
            vi.mocked(repo.deleteUser).mockResolvedValue(actionSuccess);

            const result = await useCases.deleteUser(1);

            expect(result).toEqual(actionSuccess);
            expect(repo.deleteUser).toHaveBeenCalledWith(1);
        });

        it('debería retornar respuesta de error cuando el usuario no existe', async () => {
            vi.mocked(repo.deleteUser).mockResolvedValue(actionError);

            const result = await useCases.deleteUser(999);

            expect(result.success).toBe(false);
            expect(result.error).toBeDefined();
        });

        it('debería propagar la excepción cuando el repositorio lanza un error de red', async () => {
            vi.mocked(repo.deleteUser).mockRejectedValue(new Error('Network error'));

            await expect(useCases.deleteUser(1)).rejects.toThrow('Network error');
        });
    });

    // ---------------------------------------------------------------
    // assignRoles
    // ---------------------------------------------------------------
    describe('assignRoles', () => {
        const assignDto = {
            assignments: [
                { business_id: 1, role_id: 2 },
                { business_id: 3, role_id: 4 },
            ],
        };

        it('debería asignar roles a un usuario y retornar confirmación', async () => {
            vi.mocked(repo.assignRoles).mockResolvedValue(actionSuccess);

            const result = await useCases.assignRoles(1, assignDto);

            expect(result).toEqual(actionSuccess);
            expect(repo.assignRoles).toHaveBeenCalledOnce();
            expect(repo.assignRoles).toHaveBeenCalledWith(1, assignDto);
        });

        it('debería pasar los assignments correctos al repositorio', async () => {
            vi.mocked(repo.assignRoles).mockResolvedValue(actionSuccess);

            await useCases.assignRoles(5, assignDto);

            const [userId, data] = vi.mocked(repo.assignRoles).mock.calls[0];
            expect(userId).toBe(5);
            expect(data.assignments).toHaveLength(2);
            expect(data.assignments[0]).toEqual({ business_id: 1, role_id: 2 });
        });

        it('debería propagar el error cuando la asignación falla', async () => {
            vi.mocked(repo.assignRoles).mockRejectedValue(new Error('Rol no encontrado'));

            await expect(useCases.assignRoles(1, assignDto)).rejects.toThrow('Rol no encontrado');
        });
    });
});
