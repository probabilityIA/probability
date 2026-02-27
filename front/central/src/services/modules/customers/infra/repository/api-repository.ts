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

            if (!res.ok) {
                throw new Error(data.error || data.message || 'An error occurred');
            }

            return data;
        } catch (error) {
            throw error;
        }
    }

    /** Agrega ?business_id=X a la url si se provee (para super admin) */
    private withBusinessId(path: string, businessId?: number): string {
        if (!businessId) return path;
        const sep = path.includes('?') ? '&' : '?';
        return `${path}${sep}business_id=${businessId}`;
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

    async getCustomerById(id: number, businessId?: number): Promise<CustomerDetail> {
        return this.fetch<CustomerDetail>(this.withBusinessId(`/customers/${id}`, businessId));
    }

    async createCustomer(data: CreateCustomerDTO, businessId?: number): Promise<CustomerInfo> {
        return this.fetch<CustomerInfo>(this.withBusinessId('/customers', businessId), {
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    async updateCustomer(id: number, data: UpdateCustomerDTO, businessId?: number): Promise<CustomerInfo> {
        return this.fetch<CustomerInfo>(this.withBusinessId(`/customers/${id}`, businessId), {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    async deleteCustomer(id: number, businessId?: number): Promise<DeleteCustomerResponse> {
        return this.fetch<DeleteCustomerResponse>(this.withBusinessId(`/customers/${id}`, businessId), {
            method: 'DELETE',
        });
    }
}
