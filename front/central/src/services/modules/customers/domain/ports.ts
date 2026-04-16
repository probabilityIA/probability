import {
    CustomerInfo,
    CustomerDetail,
    CustomerSummary,
    CustomersListResponse,
    CustomerAddressListResponse,
    CustomerProductListResponse,
    CustomerOrderItemListResponse,
    GetCustomersParams,
    PaginationParams,
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
    getCustomerSummary(customerId: number, businessId?: number): Promise<CustomerSummary>;
    getCustomerAddresses(customerId: number, params?: PaginationParams): Promise<CustomerAddressListResponse>;
    getCustomerProducts(customerId: number, params?: PaginationParams): Promise<CustomerProductListResponse>;
    getCustomerOrderItems(customerId: number, params?: PaginationParams): Promise<CustomerOrderItemListResponse>;
}
