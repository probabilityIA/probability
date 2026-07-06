'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';
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
export const changePasswordAction = async (data: ChangePasswordRequest, token?: string): Promise<ChangePasswordResponse> => {
    try {
        if (!token) {
            const cookieStore = await cookies();
            token = cookieStore.get('session_token')?.value;
        }

        if (!token) {
            throw new Error('No se encontró el token de sesión. Por favor, inicia sesión nuevamente.');
        }

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



const DEMO_API_BASE = process.env.API_BASE_URL || process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:3050/api/v1';

export async function demoRegisterAction(payload: { full_name: string; business_name: string; email: string; password: string; phone?: string; channel?: 'email' | 'whatsapp'; }): Promise<{ success: boolean; message?: string; error?: string }> {
    try {
        const res = await fetch(`${DEMO_API_BASE}/auth/demo-register`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload),
            cache: 'no-store',
        });
        const data = await res.json();
        if (!res.ok) {
            return { success: false, error: data.error || data.message || 'No se pudo crear la demo' };
        }
        return { success: true, message: data.message };
    } catch (error: any) {
        return { success: false, error: error.message || 'Error al conectar con el servidor' };
    }
}

export async function verifyEmailAction(token: string): Promise<{ success: boolean; message?: string; error?: string }> {
    try {
        const res = await fetch(`${DEMO_API_BASE}/auth/verify-email`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ token }),
            cache: 'no-store',
        });
        const data = await res.json();
        if (!res.ok) {
            return { success: false, error: data.error || data.message || 'No se pudo verificar la cuenta' };
        }
        return { success: true, message: data.message };
    } catch (error: any) {
        return { success: false, error: error.message || 'Error al conectar con el servidor' };
    }
}

export async function demoVerifyOtpAction(email: string, code: string): Promise<{ success: boolean; message?: string; error?: string }> {
    try {
        const res = await fetch(`${DEMO_API_BASE}/auth/demo-verify-otp`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, code }),
            cache: 'no-store',
        });
        const data = await res.json();
        if (!res.ok || !data.success) {
            return { success: false, error: data.message || data.error || 'Codigo invalido o expirado' };
        }
        return { success: true, message: data.message };
    } catch (error: any) {
        return { success: false, error: error.message || 'Error al conectar con el servidor' };
    }
}

export async function recoveryChannelsAction(email: string): Promise<{ email: boolean; whatsapp: { available: boolean; masked_phone: string }; error?: string }> {
    try {
        const res = await fetch(`${env.API_BASE_URL}/auth/recovery-channels`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email }),
            cache: 'no-store',
        });
        const data = await res.json();
        if (!res.ok) {
            return { email: true, whatsapp: { available: false, masked_phone: '' }, error: data.error || data.message };
        }
        return { email: data.email, whatsapp: data.whatsapp };
    } catch (error: any) {
        return { email: true, whatsapp: { available: false, masked_phone: '' }, error: error.message || 'Error al conectar con el servidor' };
    }
}

export async function forgotPasswordAction(email: string, channel: 'email' | 'whatsapp' = 'email'): Promise<{ success: boolean; message?: string; error?: string }> {
    try {
        const res = await fetch(`${env.API_BASE_URL}/auth/forgot-password`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, channel }),
            cache: 'no-store',
        });
        const data = await res.json();
        if (!res.ok) {
            return { success: false, error: data.error || data.message || 'No se pudo procesar la solicitud' };
        }
        return { success: true, message: data.message };
    } catch (error: any) {
        return { success: false, error: error.message || 'Error al conectar con el servidor' };
    }
}

export async function verifyOtpAction(email: string, code: string): Promise<{ success: boolean; token?: string; message?: string; error?: string }> {
    try {
        const res = await fetch(`${env.API_BASE_URL}/auth/verify-otp`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, code }),
            cache: 'no-store',
        });
        const data = await res.json();
        if (!res.ok || !data.success) {
            return { success: false, error: data.message || data.error || 'Codigo invalido o expirado' };
        }
        return { success: true, token: data.token, message: data.message };
    } catch (error: any) {
        return { success: false, error: error.message || 'Error al conectar con el servidor' };
    }
}

export async function resetPasswordAction(token: string, newPassword: string): Promise<{ success: boolean; message?: string; error?: string }> {
    try {
        const res = await fetch(`${env.API_BASE_URL}/auth/reset-password`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ token, new_password: newPassword }),
            cache: 'no-store',
        });
        const data = await res.json();
        if (!res.ok) {
            return { success: false, error: data.error || data.message || 'No se pudo restablecer la contrasena' };
        }
        return { success: true, message: data.message };
    } catch (error: any) {
        return { success: false, error: error.message || 'Error al conectar con el servidor' };
    }
}
