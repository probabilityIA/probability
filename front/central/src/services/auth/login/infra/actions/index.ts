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

        // ✅ NO setear cookie aquí - el backend ya la setea como HttpOnly
        // El backend Go setea: c.SetCookie("session_token", token, ...)
        // Next.js recibirá esa cookie automáticamente en el navegador

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
 * Lee el token de la cookie HttpOnly automáticamente
 */
export const getRolesPermissionsAction = async (): Promise<UserRolesPermissionsSuccessResponse> => {
    try {
        // Leer token de cookie HttpOnly (seteada por el backend)
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value;

        if (!token) {
            throw new Error('No session token found');
        }

        return await useCase.getRolesPermissions(token);
    } catch (error: any) {
        console.error('Get Roles Permissions Action Error:', error.message);
        throw new Error(error.message);
    }
};


