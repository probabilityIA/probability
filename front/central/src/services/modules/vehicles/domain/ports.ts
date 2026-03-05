import {
    VehicleInfo,
    VehiclesListResponse,
    GetVehiclesParams,
    CreateVehicleDTO,
    UpdateVehicleDTO,
    DeleteVehicleResponse,
} from './types';

export interface IVehicleRepository {
    getVehicles(params?: GetVehiclesParams): Promise<VehiclesListResponse>;
    getVehicleById(id: number, businessId?: number): Promise<VehicleInfo>;
    createVehicle(data: CreateVehicleDTO, businessId?: number): Promise<VehicleInfo>;
    updateVehicle(id: number, data: UpdateVehicleDTO, businessId?: number): Promise<VehicleInfo>;
    deleteVehicle(id: number, businessId?: number): Promise<DeleteVehicleResponse>;
}
