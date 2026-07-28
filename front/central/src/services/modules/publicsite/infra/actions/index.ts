'use server';

import { PublicSiteApiRepository } from '../repository/api-repository';
import { PublicSiteUseCases } from '../../app/use-cases';
import { ContactFormDTO, CreateCheckoutInput } from '../../domain/types';

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
        const session = await getUseCases().createCheckoutSession(slug, data);
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
