import { IDriverRepository } from '../domain/ports';
import { GetDriversParams, CreateDriverDTO, UpdateDriverDTO } from '../domain/types';

export class DriverUseCases {
    constructor(private repository: IDriverRepository) {}

    async getDrivers(params?: GetDriversParams) {
        return this.repository.getDrivers(params);
    }

    async getDriverById(id: number, businessId?: number) {
        return this.repository.getDriverById(id, businessId);
    }

    async createDriver(data: CreateDriverDTO, businessId?: number) {
        return this.repository.createDriver(data, businessId);
    }

    async updateDriver(id: number, data: UpdateDriverDTO, businessId?: number) {
        return this.repository.updateDriver(id, data, businessId);
    }

    async deleteDriver(id: number, businessId?: number) {
        return this.repository.deleteDriver(id, businessId);
    }
}
