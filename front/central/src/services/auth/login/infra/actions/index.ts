'use server';

import { cookies } from 'next/headers';
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
        const response = await useCase.login(credentials);

        // Set cookie for server-side access
        const cookieStore = await cookies();
        cookieStore.set('session_token', response.data.token, {
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            path: '/',
            maxAge: 60 * 60 * 24 * 7 // 1 week
        });

        return response;
    } catch (error: any) {
        console.error('Login Action Error:', error.message);
        throw new Error(error.message); // Re-throw to be caught by client
    }
};

/**
 * Server Action para cambiar contraseña
 */
/**
 * Server Action para cambiar contraseña
 */
export const changePasswordAction = async (data: ChangePasswordRequest, token: string): Promise<ChangePasswordResponse> => {
    try {
        return await useCase.changePassword(data, token);
    } catch (error: any) {
        console.error('Change Password Action Error:', error.message);
        throw new Error(error.message);
    }
};

/**
 * Server Action para generar contraseña
 */
export const generatePasswordAction = async (data: GeneratePasswordRequest, token: string): Promise<GeneratePasswordResponse> => {
    try {
        return await useCase.generatePassword(data, token);
    } catch (error: any) {
        console.error('Generate Password Action Error:', error.message);
        throw new Error(error.message);
    }
};

/**
 * Server Action para obtener roles y permisos
 */
export const getRolesPermissionsAction = async (token: string): Promise<UserRolesPermissionsSuccessResponse> => {
    try {
        return await useCase.getRolesPermissions(token);
    } catch (error: any) {
        console.error('Get Roles Permissions Action Error:', error.message);
        throw new Error(error.message);
    }
};


