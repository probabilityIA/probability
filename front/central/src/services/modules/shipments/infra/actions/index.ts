'use server';

import { ShipmentApiRepository } from '../repository/api-repository';
import { ShipmentUseCases } from '../../app/use-cases';
import { GetShipmentsParams } from '../../domain/types';

const getUseCases = async () => {
    const repository = new ShipmentApiRepository();
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
