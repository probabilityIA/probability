'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { ShippingConfigApiRepository } from '../repository/api-repository';
import { ShippingConfigUseCases } from '../../app/use-cases';
import { SaveShippingConfigRequest } from '../../domain/types';

const getUseCases = async () => {
    const token = await getAuthToken();
    const repository = new ShippingConfigApiRepository(token);
    return new ShippingConfigUseCases(repository);
};

export const getShippingConfigAction = async (businessId?: number) => {
    try {
        const data = await (await getUseCases()).getOverview(businessId);
        return { success: true, data };
    } catch (error: any) {
        console.error('Get Shipping Config Action Error:', error.message);
        return { success: false, message: error.message || 'Error al cargar la configuración de envíos' };
    }
};

export const saveShippingConfigAction = async (req: SaveShippingConfigRequest, businessId?: number) => {
    try {
        const data = await (await getUseCases()).saveBusinessConfig(req, businessId);
        return { success: true, data };
    } catch (error: any) {
        console.error('Save Shipping Config Action Error:', error.message);
        return { success: false, message: error.message || 'Error al guardar la configuración de envíos' };
    }
};

export const saveWarehouseShippingConfigAction = async (warehouseId: number, req: SaveShippingConfigRequest, businessId?: number) => {
    try {
        const data = await (await getUseCases()).saveWarehouseConfig(warehouseId, req, businessId);
        return { success: true, data };
    } catch (error: any) {
        console.error('Save Warehouse Shipping Config Action Error:', error.message);
        return { success: false, message: error.message || 'Error al guardar la configuración de la bodega' };
    }
};

export const deleteWarehouseShippingConfigAction = async (warehouseId: number, businessId?: number) => {
    try {
        await (await getUseCases()).deleteWarehouseConfig(warehouseId, businessId);
        return { success: true };
    } catch (error: any) {
        console.error('Delete Warehouse Shipping Config Action Error:', error.message);
        return { success: false, message: error.message || 'Error al eliminar la configuración de la bodega' };
    }
};

export const setDefaultWarehouseAction = async (warehouseId: number, businessId?: number) => {
    try {
        await (await getUseCases()).setDefaultWarehouse(warehouseId, businessId);
        return { success: true };
    } catch (error: any) {
        console.error('Set Default Warehouse Action Error:', error.message);
        return { success: false, message: error.message || 'Error al cambiar la bodega predeterminada' };
    }
};
