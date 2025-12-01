import { env } from '@/shared/config/env';
import {
    LoginRequest,
    ChangePasswordRequest,
    GeneratePasswordRequest,
    GenerateBusinessTokenRequest,
    mapLoginRequest,
    mapChangePasswordRequest,
    mapGeneratePasswordRequest,
    mapGenerateBusinessTokenRequest
} from './mapper/request';
import {
    LoginSuccessResponse,
    UserRolesPermissionsSuccessResponse,
    ChangePasswordResponse,
    GeneratePasswordResponse,
    GenerateBusinessTokenSuccessResponse,
    mapLoginResponse,
    mapUserRolesPermissionsResponse,
    mapChangePasswordResponse,
    mapGeneratePasswordResponse,
    mapGenerateBusinessTokenResponse
} from './mapper/response';

import { ILoginRepository } from '../../domain';

export class LoginRepository implements ILoginRepository {
    private baseUrl: string;

    constructor() {
        this.baseUrl = env.API_BASE_URL;
    }

    /**
     * Autentica un usuario
     * POST /auth/login
     */
    async login(credentials: LoginRequest): Promise<LoginSuccessResponse> {
        const payload = mapLoginRequest(credentials);
        const response = await fetch(`${this.baseUrl}/auth/login`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(payload),
            cache: 'no-store', // SSR: No cachear login
        });

        if (!response.ok) {
            let errorMessage = `Error en login: ${response.statusText}`;
            try {
                const errorData = await response.json();
                if (errorData && errorData.error) {
                    errorMessage = errorData.error;
                }
            } catch (e) {
                // Si no se puede parsear el JSON, usar el statusText por defecto
            }
            throw new Error(errorMessage);
        }

        const data = await response.json();
        return mapLoginResponse(data);
    }

    /**
     * Cambia la contraseña del usuario
     * POST /auth/change-password
     */
    async changePassword(data: ChangePasswordRequest, token: string): Promise<ChangePasswordResponse> {
        const payload = mapChangePasswordRequest(data);
        const response = await fetch(`${this.baseUrl}/auth/change-password`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            throw new Error(`Error al cambiar contraseña: ${response.statusText}`);
        }

        const resData = await response.json();
        return mapChangePasswordResponse(resData);
    }

    /**
     * Genera una nueva contraseña (admin o propio usuario)
     * POST /auth/generate-password
     */
    async generatePassword(data: GeneratePasswordRequest, token: string): Promise<GeneratePasswordResponse> {
        const payload = mapGeneratePasswordRequest(data);
        const response = await fetch(`${this.baseUrl}/auth/generate-password`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            throw new Error(`Error al generar contraseña: ${response.statusText}`);
        }

        const resData = await response.json();
        return mapGeneratePasswordResponse(resData);
    }

    /**
     * Obtiene roles y permisos del usuario
     * GET /auth/roles-permissions
     */
    async getRolesPermissions(token: string): Promise<UserRolesPermissionsSuccessResponse> {
        const response = await fetch(`${this.baseUrl}/auth/roles-permissions`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            cache: 'no-store', // SSR: Datos frescos siempre
        });

        if (!response.ok) {
            throw new Error(`Error al obtener roles y permisos: ${response.statusText}`);
        }

        const data = await response.json();
        return mapUserRolesPermissionsResponse(data);
    }

    /**
     * Genera un token de negocio
     * POST /auth/business-token
     */
    async generateBusinessToken(data: GenerateBusinessTokenRequest, token: string): Promise<GenerateBusinessTokenSuccessResponse> {
        const payload = mapGenerateBusinessTokenRequest(data);
        const response = await fetch(`${this.baseUrl}/auth/business-token`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });

        if (!response.ok) {
            throw new Error(`Error al generar token de negocio: ${response.statusText}`);
        }

        const resData = await response.json();
        return mapGenerateBusinessTokenResponse(resData);
    }
}
