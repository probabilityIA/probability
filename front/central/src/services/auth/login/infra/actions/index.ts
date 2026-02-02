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

/**
 * Server Action para login en desarrollo local.
 *
 * En producción, el login se hace con fetch directo desde el cliente para
 * que el navegador reciba la cookie Partitioned directamente (necesario para Shopify iframe).
 *
 * En desarrollo local, este Server Action se usa para evitar problemas de proxy con cookies.
 */
export async function loginServerAction(email: string, password: string) {
    try {
        const response = await fetch('http://localhost:3050/api/v1/auth/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });

        if (!response.ok) {
            const errorData = await response.json();
            return {
                success: false,
                error: errorData.error || errorData.message || 'Error al iniciar sesión',
            };
        }

        // Extraer Set-Cookie header del backend
        const setCookieHeader = response.headers.get('set-cookie');

        if (setCookieHeader) {
            // Parsear el cookie manualmente
            const tokenMatch = setCookieHeader.match(/session_token=([^;]+)/);
            const maxAgeMatch = setCookieHeader.match(/Max-Age=(\d+)/);

            if (tokenMatch && tokenMatch[1]) {
                const cookieStore = await cookies();

                // Setear cookie usando Next.js cookies API
                cookieStore.set('session_token', tokenMatch[1], {
                    maxAge: maxAgeMatch ? parseInt(maxAgeMatch[1]) : 7 * 24 * 60 * 60, // 7 días por defecto
                    path: '/',
                    httpOnly: true,
                    secure: false, // En local dev no usamos HTTPS
                    sameSite: 'lax', // En local dev no necesitamos 'none'
                });
            }
        }

        const data = await response.json();
        return {
            success: true,
            data,
        };
    } catch (error: any) {
        return {
            success: false,
            error: error.message || 'Error al conectar con el servidor',
        };
    }
}


