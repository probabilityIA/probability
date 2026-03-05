import { IVehicleRepository } from '../domain/ports';
import { GetVehiclesParams, CreateVehicleDTO, UpdateVehicleDTO } from '../domain/types';

export class VehicleUseCases {
    constructor(private repository: IVehicleRepository) {}

    async getVehicles(params?: GetVehiclesParams) {
        return this.repository.getVehicles(params);
    }

    async getVehicleById(id: number, businessId?: number) {
        return this.repository.getVehicleById(id, businessId);
    }

    async createVehicle(data: CreateVehicleDTO, businessId?: number) {
        return this.repository.createVehicle(data, businessId);
    }

    async updateVehicle(id: number, data: UpdateVehicleDTO, businessId?: number) {
        return this.repository.updateVehicle(id, data, businessId);
    }

    async deleteVehicle(id: number, businessId?: number) {
        return this.repository.deleteVehicle(id, businessId);
    }
}
