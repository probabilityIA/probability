'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { ShipmentApiRepository } from '../repository/api-repository';
import { ShipmentUseCases } from '../../app/use-cases';
import { GetShipmentsParams, CreateShipmentRequest } from '../../domain/types';

const getUseCases = async () => {
    const token = await getAuthToken();
    const repository = new ShipmentApiRepository(token);
    return new ShipmentUseCases(repository);
};

export const getShipmentsAction = async (params?: GetShipmentsParams) => {
    try {
        return await (await getUseCases()).getShipments(params);
    } catch (error: any) {
        console.error('Get Shipments Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al obtener envíos',
            data: [],
            total: 0,
            page: params?.page || 1,
            page_size: params?.page_size || 20,
            total_pages: 0
        };
    }
};

export const trackShipmentAction = async (trackingNumber: string) => {
    try {
        return await (await getUseCases()).trackShipment(trackingNumber);
    } catch (error: any) {
        console.error('Track Shipment Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al rastrear envío',
        };
    }
};

export const cancelShipmentAction = async (id: string) => {
    try {
        return await (await getUseCases()).cancelShipment(id);
    } catch (error: any) {
        console.error('Cancel Shipment Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al cancelar envío',
        };
    }
};

export const createShipmentAction = async (req: CreateShipmentRequest) => {
    try {
        return await (await getUseCases()).createShipment(req);
    } catch (error: any) {
        console.error('Create Shipment Action Error:', error.message);
        return {
            success: false,
            message: error.message || 'Error al crear envío',
        };
    }
};

export const quoteShipmentAction = async (req: any) => {
    try {
        const data = await (await getUseCases()).quoteShipment(req);
        return { success: true, data };
    } catch (error: any) {
        console.error('Quote Shipment Action Error:', error.message);
        return { success: false, message: error.message || 'Error al cotizar envío' };
    }
};

export const generateGuideAction = async (req: any) => {
    try {
        const data = await (await getUseCases()).generateGuide(req);
        return { success: true, data };
    } catch (error: any) {
        console.error('Generate Guide Action Error:', error.message);
        return { success: false, message: error.message || 'Error al generar guía' };
    }
};

// Origin Addresses Actions
export const getOriginAddressesAction = async () => {
    try {
        const data = await (await getUseCases()).getOriginAddresses();
        return { success: true, data };
    } catch (error: any) {
        console.error('Get Origin Addresses Action Error:', error.message);
        return { success: false, message: error.message || 'Error al obtener direcciones de origen' };
    }
};

export const createOriginAddressAction = async (req: any) => {
    try {
        const data = await (await getUseCases()).createOriginAddress(req);
        return { success: true, data };
    } catch (error: any) {
        console.error('Create Origin Address Action Error:', error.message);
        return { success: false, message: error.message || 'Error al crear dirección de origen' };
    }
};

export const updateOriginAddressAction = async (id: number, req: any) => {
    try {
        const data = await (await getUseCases()).updateOriginAddress(id, req);
        return { success: true, data };
    } catch (error: any) {
        console.error('Update Origin Address Action Error:', error.message);
        return { success: false, message: error.message || 'Error al actualizar dirección de origen' };
    }
};

export const deleteOriginAddressAction = async (id: number) => {
    try {
        const data = await (await getUseCases()).deleteOriginAddress(id);
        return { success: true, message: data.message };
    } catch (error: any) {
        console.error('Delete Origin Address Action Error:', error.message);
        return { success: false, message: error.message || 'Error al eliminar dirección de origen' };
    }
};

