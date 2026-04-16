import {
    DriverInfo,
    DriversListResponse,
    GetDriversParams,
    CreateDriverDTO,
    UpdateDriverDTO,
    DeleteDriverResponse,
} from './types';

export interface IDriverRepository {
    getDrivers(params?: GetDriversParams): Promise<DriversListResponse>;
    getDriverById(id: number, businessId?: number): Promise<DriverInfo>;
    createDriver(data: CreateDriverDTO, businessId?: number): Promise<DriverInfo>;
    updateDriver(id: number, data: UpdateDriverDTO, businessId?: number): Promise<DriverInfo>;
    deleteDriver(id: number, businessId?: number): Promise<DeleteDriverResponse>;
}
