import { env } from '@/shared/config/env';
import { ICustomerRepository } from '../../domain/ports';
import {
    CustomerInfo,
    CustomerDetail,
    CustomersListResponse,
    GetCustomersParams,
    CreateCustomerDTO,
    UpdateCustomerDTO,
    DeleteCustomerResponse,
} from '../../domain/types';

export class CustomerApiRepository implements ICustomerRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        console.log(`[API Request] ${options.method || 'GET'} ${url}`, {
            headers: options.headers,
            body: options.body,
        });

        const headers: Record<string, string> = {
            Accept: 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        try {
            const res = await fetch(url, { ...options, headers });
            const data = await res.json();

            console.log(`[API Response] ${res.status} ${url}`, data);

            if (!res.ok) {
                console.error(`[API Error] ${res.status} ${url}`, data);
                throw new Error(data.error || data.message || 'An error occurred');
            }

            return data;
        } catch (error) {
            console.error(`[API Network Error] ${url}`, error);
            throw error;
        }
    }

    async getCustomers(params?: GetCustomersParams): Promise<CustomersListResponse> {
        const searchParams = new URLSearchParams();
        if (params) {
            Object.entries(params).forEach(([key, value]) => {
                if (value !== undefined && value !== null && value !== '') {
                    searchParams.append(key, String(value));
                }
            });
        }
        const query = searchParams.toString();
        return this.fetch<CustomersListResponse>(`/customers${query ? `?${query}` : ''}`);
    }

    async getCustomerById(id: number): Promise<CustomerDetail> {
        return this.fetch<CustomerDetail>(`/customers/${id}`);
    }

    async createCustomer(data: CreateCustomerDTO): Promise<CustomerInfo> {
        return this.fetch<CustomerInfo>('/customers', {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateCustomer(id: number, data: UpdateCustomerDTO): Promise<CustomerInfo> {
        return this.fetch<CustomerInfo>(`/customers/${id}`, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteCustomer(id: number): Promise<DeleteCustomerResponse> {
        return this.fetch<DeleteCustomerResponse>(`/customers/${id}`, {
            method: 'DELETE',
        });
    }
}
