import { ICustomerRepository } from '../domain/ports';
import { GetCustomersParams, CreateCustomerDTO, UpdateCustomerDTO } from '../domain/types';

export class CustomerUseCases {
    constructor(private repository: ICustomerRepository) {}

    async getCustomers(params?: GetCustomersParams) {
        return this.repository.getCustomers(params);
    }

    async getCustomerById(id: number) {
        return this.repository.getCustomerById(id);
    }

    async createCustomer(data: CreateCustomerDTO) {
        return this.repository.createCustomer(data);
    }

    async updateCustomer(id: number, data: UpdateCustomerDTO) {
        return this.repository.updateCustomer(id, data);
    }

    async deleteCustomer(id: number) {
        return this.repository.deleteCustomer(id);
    }
}
