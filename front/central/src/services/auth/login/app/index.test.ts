import { describe, it, expect, vi, beforeEach } from 'vitest';
import { LoginUseCase } from './index';
import { ILoginRepository } from '../domain';
import {
    LoginSuccessResponse,
    ChangePasswordResponse,
    GeneratePasswordResponse,
    UserRolesPermissionsSuccessResponse,
} from '../infra/repository/mapper/response';

// -----------------------------------------------------------------
// Helpers: datos de prueba reutilizables
// -----------------------------------------------------------------

const TOKEN = 'test-jwt-token';

const makeLoginResponse = (): LoginSuccessResponse => ({
    success: true,
    data: {
        user: {
            id: 1,
            name: 'Juan Pérez',
            email: 'juan@test.com',
            phone: '3001234567',
            avatar_url: '',
            is_active: true,
        },
        token: TOKEN,
        require_password_change: false,
        businesses: [],
        scope: 'platform',
        is_super_admin: false,
    },
});

const makeChangePasswordResponse = (): ChangePasswordResponse => ({
    success: true,
    message: 'Contraseña actualizada correctamente',
});

const makeGeneratePasswordResponse = (): GeneratePasswordResponse => ({
    success: true,
    email: 'juan@test.com',
    password: 'newPassword123',
    message: 'Contraseña generada correctamente',
});

const makeRolesPermissionsResponse = (): UserRolesPermissionsSuccessResponse => ({
    success: true,
    data: {
        is_super: false,
        business_id: 1,
        business_name: 'Mi Negocio',
        business_type_id: 1,
        business_type_name: 'Restaurante',
        role: { id: 1, name: 'Admin', description: 'Administrador' },
        resources: [
            { resource: 'orders', actions: ['read', 'write'], active: true },
        ],
    },
});

// -----------------------------------------------------------------
// Mock del repositorio
// -----------------------------------------------------------------

function createMockRepository(): ILoginRepository {
    return {
        login: vi.fn(),
        changePassword: vi.fn(),
        generatePassword: vi.fn(),
        getRolesPermissions: vi.fn(),
    } as unknown as ILoginRepository;
}

// -----------------------------------------------------------------
// Suite principal
// -----------------------------------------------------------------

describe('LoginUseCase', () => {
    let repo: ReturnType<typeof createMockRepository>;
    let useCases: LoginUseCase;

    beforeEach(() => {
        repo = createMockRepository();
        useCases = new LoginUseCase(repo as unknown as ILoginRepository);
    });

    // ---------------------------------------------------------------
    // login
    // ---------------------------------------------------------------
    describe('login', () => {
        const credentials = { email: 'juan@test.com', password: 'password123' };

        it('debería retornar la respuesta de login exitoso', async () => {
            const loginResponse = makeLoginResponse();
            vi.mocked(repo.login).mockResolvedValue(loginResponse);

            const result = await useCases.login(credentials);

            expect(result).toEqual(loginResponse);
            expect(repo.login).toHaveBeenCalledOnce();
            expect(repo.login).toHaveBeenCalledWith(credentials);
        });

        it('debería propagar el error cuando las credenciales son inválidas', async () => {
            vi.mocked(repo.login).mockRejectedValue(new Error('Credenciales inválidas'));

            await expect(useCases.login(credentials)).rejects.toThrow('Credenciales inválidas');
        });
    });

    // ---------------------------------------------------------------
    // changePassword
    // ---------------------------------------------------------------
    describe('changePassword', () => {
        const data = { current_password: 'oldPass123', new_password: 'newPass456' };

        it('debería cambiar la contraseña y retornar confirmación', async () => {
            const response = makeChangePasswordResponse();
            vi.mocked(repo.changePassword).mockResolvedValue(response);

            const result = await useCases.changePassword(data, TOKEN);

            expect(result).toEqual(response);
            expect(repo.changePassword).toHaveBeenCalledOnce();
            expect(repo.changePassword).toHaveBeenCalledWith(data, TOKEN);
        });

        it('debería pasar el token al repositorio', async () => {
            vi.mocked(repo.changePassword).mockResolvedValue(makeChangePasswordResponse());

            await useCases.changePassword(data, 'custom-token');

            expect(repo.changePassword).toHaveBeenCalledWith(data, 'custom-token');
        });

        it('debería propagar el error cuando el cambio falla', async () => {
            vi.mocked(repo.changePassword).mockRejectedValue(new Error('Contraseña actual incorrecta'));

            await expect(useCases.changePassword(data, TOKEN)).rejects.toThrow('Contraseña actual incorrecta');
        });
    });

    // ---------------------------------------------------------------
    // generatePassword
    // ---------------------------------------------------------------
    describe('generatePassword', () => {
        const data = { user_id: 1 };

        it('debería generar una nueva contraseña y retornar la respuesta', async () => {
            const response = makeGeneratePasswordResponse();
            vi.mocked(repo.generatePassword).mockResolvedValue(response);

            const result = await useCases.generatePassword(data, TOKEN);

            expect(result).toEqual(response);
            expect(repo.generatePassword).toHaveBeenCalledOnce();
            expect(repo.generatePassword).toHaveBeenCalledWith(data, TOKEN);
        });

        it('debería pasar el token al repositorio', async () => {
            vi.mocked(repo.generatePassword).mockResolvedValue(makeGeneratePasswordResponse());

            await useCases.generatePassword(data, 'another-token');

            expect(repo.generatePassword).toHaveBeenCalledWith(data, 'another-token');
        });

        it('debería propagar el error cuando la generación falla', async () => {
            vi.mocked(repo.generatePassword).mockRejectedValue(new Error('Usuario no encontrado'));

            await expect(useCases.generatePassword(data, TOKEN)).rejects.toThrow('Usuario no encontrado');
        });
    });

    // ---------------------------------------------------------------
    // getRolesPermissions
    // ---------------------------------------------------------------
    describe('getRolesPermissions', () => {
        it('debería retornar los roles y permisos del usuario', async () => {
            const response = makeRolesPermissionsResponse();
            vi.mocked(repo.getRolesPermissions).mockResolvedValue(response);

            const result = await useCases.getRolesPermissions(TOKEN);

            expect(result).toEqual(response);
            expect(repo.getRolesPermissions).toHaveBeenCalledOnce();
            expect(repo.getRolesPermissions).toHaveBeenCalledWith(TOKEN);
        });

        it('debería propagar el error cuando la consulta falla', async () => {
            vi.mocked(repo.getRolesPermissions).mockRejectedValue(new Error('Token expirado'));

            await expect(useCases.getRolesPermissions(TOKEN)).rejects.toThrow('Token expirado');
        });
    });
});
