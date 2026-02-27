import {
    CustomerInfo,
    CustomerDetail,
    CustomersListResponse,
    GetCustomersParams,
    CreateCustomerDTO,
    UpdateCustomerDTO,
    DeleteCustomerResponse,
} from './types';

export interface ICustomerRepository {
    getCustomers(params?: GetCustomersParams): Promise<CustomersListResponse>;
    getCustomerById(id: number): Promise<CustomerDetail>;
    createCustomer(data: CreateCustomerDTO): Promise<CustomerInfo>;
    updateCustomer(id: number, data: UpdateCustomerDTO): Promise<CustomerInfo>;
    deleteCustomer(id: number): Promise<DeleteCustomerResponse>;
}
