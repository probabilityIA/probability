import { env } from '@/shared/config/env';
import {
    LoginRequest,
    ChangePasswordRequest,
    GeneratePasswordRequest,
    mapLoginRequest,
    mapChangePasswordRequest,
    mapGeneratePasswordRequest
} from './mapper/request';
import {
    LoginSuccessResponse,
    UserRolesPermissionsSuccessResponse,
    ChangePasswordResponse,
    GeneratePasswordResponse,
    mapLoginResponse,
    mapUserRolesPermissionsResponse,
    mapChangePasswordResponse,
    mapGeneratePasswordResponse
} from './mapper/response';

import { ILoginRepository } from '../../domain';

export class LoginRepository implements ILoginRepository {
    private baseUrl: string;

    constructor() {
        this.baseUrl = env.API_BASE_URL;
    }

    private async fetch<T>(url: string, options: RequestInit = {}): Promise<T> {
        console.log(`[API Request] ${options.method || 'GET'} ${url}`, {
            headers: options.headers,
            body: options.body
        });

        try {
            const res = await fetch(url, {
                ...options,
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers,
                },
            });

            const data = await res.json();

            console.log(`[API Response] ${res.status} ${url}`, data);

            if (!res.ok) {
                console.error(`[API Error] ${res.status} ${url}`, data);
                throw new Error(data.error || data.message || res.statusText || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    // Nueva función que retorna Response completo (para leer headers)
    private async fetchWithResponse<T>(url: string, options: RequestInit = {}): Promise<{ data: T; response: Response }> {
        console.log(`[API Request] ${options.method || 'GET'} ${url}`, {
            headers: options.headers,
            body: options.body
        });

        try {
            const res = await fetch(url, {
                ...options,
                headers: {
                    'Content-Type': 'application/json',
                    ...options.headers,
                },
            });

            const data = await res.json();

            console.log(`[API Response] ${res.status} ${url}`, data);

            if (!res.ok) {
                console.error(`[API Error] ${res.status} ${url}`, data);
                throw new Error(data.error || data.message || res.statusText || 'An error occurred');
            }

            return { data, response: res };
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    /**
     * Autentica un usuario
     * POST /auth/login
     */
    async login(credentials: LoginRequest): Promise<LoginSuccessResponse> {
        const payload = mapLoginRequest(credentials);
        const data = await this.fetch<any>(`${this.baseUrl}/auth/login`, {
            method: 'POST',
            body: JSON.stringify(payload),
            cache: 'no-store',
        });
        return mapLoginResponse(data);
    }

    /**
     * Cambia la contraseña del usuario
     * POST /auth/change-password
     */
    async changePassword(data: ChangePasswordRequest, token: string): Promise<ChangePasswordResponse> {
        const payload = mapChangePasswordRequest(data);
        const resData = await this.fetch<any>(`${this.baseUrl}/auth/change-password`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });
        return mapChangePasswordResponse(resData);
    }

    /**
     * Genera una nueva contraseña (admin o propio usuario)
     * POST /auth/generate-password
     */
    async generatePassword(data: GeneratePasswordRequest, token: string): Promise<GeneratePasswordResponse> {
        const payload = mapGeneratePasswordRequest(data);
        const resData = await this.fetch<any>(`${this.baseUrl}/auth/generate-password`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(payload),
        });
        return mapGeneratePasswordResponse(resData);
    }

    /**
     * Obtiene roles y permisos del usuario
     * GET /auth/roles-permissions
     */
    async getRolesPermissions(token: string): Promise<UserRolesPermissionsSuccessResponse> {
        const data = await this.fetch<any>(`${this.baseUrl}/auth/roles-permissions`, {
            method: 'GET',
            headers: {
                'Authorization': `Bearer ${token}`,
            },
            cache: 'no-store',
        });
        return mapUserRolesPermissionsResponse(data);
    }


}
