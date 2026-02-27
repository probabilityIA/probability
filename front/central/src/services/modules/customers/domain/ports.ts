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
    getCustomerById(id: number, businessId?: number): Promise<CustomerDetail>;
    createCustomer(data: CreateCustomerDTO, businessId?: number): Promise<CustomerInfo>;
    updateCustomer(id: number, data: UpdateCustomerDTO, businessId?: number): Promise<CustomerInfo>;
    deleteCustomer(id: number, businessId?: number): Promise<DeleteCustomerResponse>;
}
