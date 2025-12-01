import { env } from '@/shared/config/env';

export class UserRepository {
    private baseUrl: string;

    constructor() {
        this.baseUrl = env.API_BASE_URL;
    }

    /**
     * Obtiene lista de usuarios
     * GET /users
     */
    async getUsers(params: any = {}, token: string): Promise<any> {
        const searchParams = new URLSearchParams();
        Object.keys(params).forEach(key => {
            if (params[key] !== undefined && params[key] !== null) {
                searchParams.append(key, params[key].toString());
            }
        });

        const response = await fetch(`${this.baseUrl}/users?${searchParams.toString()}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            cache: 'no-store', // SSR: Datos frescos
        });

        if (!response.ok) {
            throw new Error(`Error al obtener usuarios: ${response.statusText}`);
        }

        return response.json();
    }

    /**
     * Crea un nuevo usuario
     * POST /users
     */
    async createUser(userData: any, token: string): Promise<any> {
        const response = await fetch(`${this.baseUrl}/users`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(userData),
        });

        if (!response.ok) {
            throw new Error(`Error al crear usuario: ${response.statusText}`);
        }

        return response.json();
    }

    /**
     * Obtiene un usuario por ID
     * GET /users/{id}
     */
    async getUserById(id: number | string, token: string): Promise<any> {
        const response = await fetch(`${this.baseUrl}/users/${id}`, {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            cache: 'no-store',
        });

        if (!response.ok) {
            throw new Error(`Error al obtener usuario ${id}: ${response.statusText}`);
        }

        return response.json();
    }

    /**
     * Actualiza un usuario
     * PUT /users/{id}
     */
    async updateUser(id: number | string, userData: any, token: string): Promise<any> {
        const response = await fetch(`${this.baseUrl}/users/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(userData),
        });

        if (!response.ok) {
            throw new Error(`Error al actualizar usuario ${id}: ${response.statusText}`);
        }

        return response.json();
    }

    /**
     * Elimina un usuario
     * DELETE /users/{id}
     */
    async deleteUser(id: number | string, token: string): Promise<any> {
        const response = await fetch(`${this.baseUrl}/users/${id}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
        });

        if (!response.ok) {
            throw new Error(`Error al eliminar usuario ${id}: ${response.statusText}`);
        }

        return response.json();
    }

    /**
     * Asigna roles a un usuario
     * POST /users/{id}/assign-role
     */
    async assignRole(id: number | string, assignments: any, token: string): Promise<any> {
        const response = await fetch(`${this.baseUrl}/users/${id}/assign-role`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`,
            },
            body: JSON.stringify(assignments),
        });

        if (!response.ok) {
            throw new Error(`Error al asignar roles al usuario ${id}: ${response.statusText}`);
        }

        return response.json();
    }
}
