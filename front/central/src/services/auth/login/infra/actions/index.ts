'use server';

import { LoginRepository } from '../repository';
import { LoginUseCase } from '../../app';
import {
    LoginRequest,
    ChangePasswordRequest,
    GeneratePasswordRequest,
    GenerateBusinessTokenRequest
} from '../repository/mapper/request';
import {
    LoginSuccessResponse,
    UserRolesPermissionsSuccessResponse,
    ChangePasswordResponse,
    GeneratePasswordResponse,
    GenerateBusinessTokenSuccessResponse
} from '../repository/mapper/response';

// Instancia del repositorio y caso de uso (Singleton implícito por módulo)
const repository = new LoginRepository();
const useCase = new LoginUseCase(repository);

/**
 * Server Action para autenticar usuario
 */
export const loginAction = async (credentials: LoginRequest): Promise<LoginSuccessResponse> => {
    try {
        return await useCase.login(credentials);
    } catch (error: any) {
        console.error('Login Action Error:', error.message);
        throw new Error(error.message); // Re-throw to be caught by client
    }
};

/**
 * Server Action para cambiar contraseña
 */
export const changePasswordAction = async (data: ChangePasswordRequest, token: string): Promise<ChangePasswordResponse> => {
    return useCase.changePassword(data, token);
};

/**
 * Server Action para generar contraseña
 */
export const generatePasswordAction = async (data: GeneratePasswordRequest, token: string): Promise<GeneratePasswordResponse> => {
    return useCase.generatePassword(data, token);
};

/**
 * Server Action para obtener roles y permisos
 */
export const getRolesPermissionsAction = async (token: string): Promise<UserRolesPermissionsSuccessResponse> => {
    return useCase.getRolesPermissions(token);
};

/**
 * Server Action para generar token de negocio
 */
export const generateBusinessTokenAction = async (data: GenerateBusinessTokenRequest, token: string): Promise<GenerateBusinessTokenSuccessResponse> => {
    return useCase.generateBusinessToken(data, token);
};
