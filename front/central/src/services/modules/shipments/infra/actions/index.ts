'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { ShipmentApiRepository } from '../repository/api-repository';
import { ShipmentUseCases } from '../../app/use-cases';
import { GetShipmentsParams } from '../../domain/types';

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
