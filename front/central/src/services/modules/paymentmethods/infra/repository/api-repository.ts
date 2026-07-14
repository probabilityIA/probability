import { env } from '@/shared/config/env';
import { IPaymentMethodRepository } from '../../domain/ports';
import { PaginatedPaymentMethodsResponse, PaymentMethodsResponse } from '../../domain/types';

export class PaymentMethodApiRepository implements IPaymentMethodRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    private async fetch<T>(path: string, options: RequestInit = {}): Promise<T> {
        const url = `${this.baseUrl}${path}`;

        const headers: Record<string, string> = {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            ...(options.headers as Record<string, string> || {}),
        };

        if (this.token) {
            headers['Authorization'] = `Bearer ${this.token}`;
        }

        const res = await fetch(url, { ...options, headers, cache: 'no-store' });
        const data = await res.json();

        if (!res.ok) {
            console.error(`[API Error] ${res.status} ${url}`, data);
            throw new Error(data.message || data.error || 'An error occurred');
        }

        return data;
    }

    async getPaymentMethods(): Promise<PaymentMethodsResponse> {
        const res = await this.fetch<PaginatedPaymentMethodsResponse>(
            '/payments/methods?is_active=true&pageSize=100'
        );
        return {
            success: res.success ?? true,
            message: res.message,
            data: (res.data || []).map((m) => ({ id: m.id, code: m.code, name: m.name })),
        };
    }
}
