'use server';

import { cookies } from 'next/headers';
import { env } from '@/shared/config/env';
import { PublicSiteApiRepository } from '../repository/api-repository';
import { PublicSiteUseCases } from '../../app/use-cases';
import { ContactFormDTO, CreateCheckoutInput, TiendaRegisterDTO, TiendaAddressInput, TiendaProfileUpdate } from '../../domain/types';

function getUseCases() {
    const repository = new PublicSiteApiRepository();
    return new PublicSiteUseCases(repository);
}

export const getPublicBusinessAction = async (slug: string) => {
    try {
        return await getUseCases().getBusinessPage(slug);
    } catch (error: any) {
        console.error('Get Public Business Action Error:', error.message);
        return null;
    }
};

export const getPublicCatalogAction = async (slug: string, params?: { page?: number; page_size?: number; search?: string; category?: string }) => {
    try {
        return await getUseCases().getCatalog(slug, params);
    } catch (error: any) {
        console.error('Get Public Catalog Action Error:', error.message);
        return { data: [], total: 0, page: 1, page_size: 12, total_pages: 0 };
    }
};

export const getPublicCategoriesAction = async (slug: string) => {
    try {
        return await getUseCases().getCategories(slug);
    } catch {
        return [];
    }
};

export const getPublicProductAction = async (slug: string, productId: string) => {
    try {
        return await getUseCases().getProduct(slug, productId);
    } catch (error: any) {
        console.error('Get Public Product Action Error:', error.message);
        return null;
    }
};

export const submitContactAction = async (slug: string, data: ContactFormDTO) => {
    try {
        return await getUseCases().submitContact(slug, data);
    } catch (error: any) {
        console.error('Submit Contact Action Error:', error.message);
        return { success: false, message: error.message || 'Error al enviar mensaje' };
    }
};

export const createCheckoutSessionAction = async (slug: string, data: CreateCheckoutInput) => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value;
        const session = await getUseCases().createCheckoutSession(slug, data, token);
        return { success: true as const, data: session };
    } catch (error: any) {
        console.error('Create Checkout Session Action Error:', error.message);
        return { success: false as const, error: error.message || 'Error al iniciar el pago' };
    }
};

export const getCheckoutStatusAction = async (slug: string, reference: string) => {
    try {
        const result = await getUseCases().getCheckoutStatus(slug, reference);
        return { success: true as const, data: result };
    } catch (error: any) {
        console.error('Get Checkout Status Action Error:', error.message);
        return { success: false as const, error: error.message || 'Error al consultar el estado del pago' };
    }
};

export const tiendaLoginAction = async (email: string, password: string) => {
    try {
        const response = await fetch(`${env.API_BASE_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password }),
        });
        const data = await response.json();
        if (!response.ok) {
            return { success: false as const, error: data.error || data.message || 'Credenciales invalidas' };
        }

        const setCookieHeader = response.headers.get('set-cookie');
        const tokenMatch = setCookieHeader?.match(/session_token=([^;]+)/);
        if (!tokenMatch?.[1]) {
            return { success: false as const, error: 'No se recibio la sesión del servidor' };
        }

        const maxAgeMatch = setCookieHeader?.match(/Max-Age=(\d+)/);
        const cookieStore = await cookies();
        cookieStore.set('session_token', tokenMatch[1], {
            maxAge: maxAgeMatch ? parseInt(maxAgeMatch[1]) : 7 * 24 * 60 * 60,
            path: '/',
            httpOnly: true,
            secure: process.env.NODE_ENV === 'production',
            sameSite: 'lax',
        });

        return { success: true as const, name: data?.data?.user?.name || '' };
    } catch (error: any) {
        console.error('Tienda Login Action Error:', error.message);
        return { success: false as const, error: 'Error al conectar con el servidor' };
    }
};

export const tiendaRegisterAction = async (slug: string, data: TiendaRegisterDTO) => {
    try {
        await getUseCases().registerCustomer(slug, data);
        return await tiendaLoginAction(data.email, data.password);
    } catch (error: any) {
        console.error('Tienda Register Action Error:', error.message);
        return { success: false as const, error: error.message || 'Error al crear la cuenta' };
    }
};

export const tiendaLogoutAction = async () => {
    const cookieStore = await cookies();
    cookieStore.delete('session_token');
    return { success: true as const };
};

export const getTiendaPricesAction = async (slug: string) => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value;
        if (!token) return null;
        return await getUseCases().getMyPrices(slug, token);
    } catch {
        return null;
    }
};

export const getTiendaOrdersAction = async (slug: string, page = 1) => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value;
        if (!token) return null;
        return await getUseCases().getMyOrders(slug, token, page);
    } catch {
        return null;
    }
};

async function withToken<T>(fn: (token: string) => Promise<T>): Promise<T | null> {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value;
        if (!token) return null;
        return await fn(token);
    } catch (error: any) {
        console.error('Tienda authed action error:', error?.message);
        return null;
    }
}

export const updateTiendaProfileAction = async (slug: string, data: TiendaProfileUpdate) =>
    withToken(token => getUseCases().updateProfile(slug, token, data));

export const uploadTiendaAvatarAction = async (slug: string, formData: FormData) =>
    withToken(token => getUseCases().uploadAvatar(slug, token, formData));

export const listTiendaAddressesAction = async (slug: string) =>
    withToken(token => getUseCases().listAddresses(slug, token));

export const addTiendaAddressAction = async (slug: string, data: TiendaAddressInput) =>
    withToken(token => getUseCases().addAddress(slug, token, data));

export const setTiendaPrimaryAddressAction = async (slug: string, id: number) =>
    withToken(token => getUseCases().setPrimaryAddress(slug, token, id));

export const deleteTiendaAddressAction = async (slug: string, id: number) =>
    withToken(token => getUseCases().deleteAddress(slug, token, id));

export const getTiendaSessionAction = async (slug: string) => {
    try {
        const cookieStore = await cookies();
        const token = cookieStore.get('session_token')?.value;
        if (!token) return null;
        return await getUseCases().getSession(slug, token);
    } catch {
        return null;
    }
};
