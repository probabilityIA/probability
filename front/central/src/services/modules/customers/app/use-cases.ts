import { ICustomerRepository } from '../domain/ports';
import { GetCustomersParams, CreateCustomerDTO, UpdateCustomerDTO } from '../domain/types';

export class CustomerUseCases {
    constructor(private repository: ICustomerRepository) {}

    async getCustomers(params?: GetCustomersParams) {
        return this.repository.getCustomers(params);
    }

    async getCustomerById(id: number, businessId?: number) {
        return this.repository.getCustomerById(id, businessId);
    }

    async createCustomer(data: CreateCustomerDTO, businessId?: number) {
        return this.repository.createCustomer(data, businessId);
    }

    async updateCustomer(id: number, data: UpdateCustomerDTO, businessId?: number) {
        return this.repository.updateCustomer(id, data, businessId);
    }

    async deleteCustomer(id: number, businessId?: number) {
        return this.repository.deleteCustomer(id, businessId);
    }
}
