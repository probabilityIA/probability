/**
 * API Client universal que detecta el contexto (iframe vs navegador normal)
 * y usa la estrategia de autenticación correcta automáticamente
 *
 * - Iframe de Shopify: Fetch directo con token de sessionStorage
 * - Navegador normal: Server Actions con cookies HttpOnly
 */

import { env } from '@/shared/config/env';
import { TokenStorage } from './token-storage';

/**
 * Detecta si estamos en un iframe
 */
function isInIframe(): boolean {
    if (typeof window === 'undefined') return false;
    try {
        return window.self !== window.top;
    } catch (e) {
        return true;
    }
}

/**
 * Detecta si estamos en un iframe de Shopify
 */
function isShopifyIframe(): boolean {
    if (typeof window === 'undefined') return false;
    try {
        const referrer = document.referrer.toLowerCase();
        return (
            isInIframe() &&
            (referrer.includes('shopify.com') ||
             referrer.includes('myshopify.com'))
        );
    } catch (e) {
        return false;
    }
}

/**
 * Cliente API universal
 * Automáticamente usa la estrategia correcta según el contexto
 */
export class UniversalApiClient {
    private baseUrl: string;

    constructor() {
        this.baseUrl = env.API_BASE_URL;
    }

    /**
     * Fetch universal que funciona en cualquier contexto
     */
    async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        // Preparar headers
        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        // En iframe de Shopify, usar token de sessionStorage
        if (isShopifyIframe()) {
            const token = TokenStorage.getSessionToken();
            if (token) {
                headers['Authorization'] = `Bearer ${token}`;
                console.log('🛍️ Shopify iframe: Usando token de sessionStorage');
            }
        } else {
            // En navegador normal, las cookies se envían automáticamente
            console.log('🌐 Navegador normal: Usando cookies HttpOnly');
        }

        try {
            const response = await fetch(url, {
                ...options,
                headers,
                credentials: 'include', // Incluir cookies
            });

            const data = await response.json();

            if (!response.ok) {
                console.error(`[API Error] ${response.status} ${url}`, data);
                throw new Error(data.message || data.error || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    /**
     * GET request
     */
    async get<T>(path: string, params?: Record<string, any>): Promise<T> {
        const queryParams = params ? `?${new URLSearchParams(params)}` : '';
        return this.fetch<T>(`${path}${queryParams}`, { method: 'GET' });
    }

    /**
     * POST request
     */
    async post<T>(path: string, body?: any): Promise<T> {
        return this.fetch<T>(path, {
            method: 'POST',
            body: JSON.stringify(body),
        });
    }

    /**
     * PUT request
     */
    async put<T>(path: string, body?: any): Promise<T> {
        return this.fetch<T>(path, {
            method: 'PUT',
            body: JSON.stringify(body),
        });
    }

    /**
     * DELETE request
     */
    async delete<T>(path: string): Promise<T> {
        return this.fetch<T>(path, { method: 'DELETE' });
    }

    /**
     * PATCH request
     */
    async patch<T>(path: string, body?: any): Promise<T> {
        return this.fetch<T>(path, {
            method: 'PATCH',
            body: JSON.stringify(body),
        });
    }
}

/**
 * Instancia singleton del cliente API
 */
export const apiClient = new UniversalApiClient();

/**
 * Hook para usar el cliente API en componentes React
 */
export function useApiClient() {
    return apiClient;
}
