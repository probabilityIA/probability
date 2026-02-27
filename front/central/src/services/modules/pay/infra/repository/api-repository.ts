import { env } from '@/shared/config/env';
import { IPayGatewayRepository } from '../../domain/ports';
import { PaymentGatewayType, PaymentGatewayTypesResponse } from '../../domain/types';

export class PayGatewayApiRepository implements IPayGatewayRepository {
    private baseUrl: string;
    private token: string | null;

    constructor(token?: string | null) {
        this.baseUrl = env.API_BASE_URL;
        this.token = token || null;
    }

    async listPaymentGatewayTypes(): Promise<PaymentGatewayTypesResponse> {
        const url = `${this.baseUrl}/integration-types/active`;

        const res = await fetch(url, {
            headers: {
                'Authorization': `Bearer ${this.token}`,
                'Content-Type': 'application/json',
            },
            cache: 'no-store',
        });

        if (!res.ok) {
            return { success: false, data: [], message: 'Error al obtener métodos de pago' };
        }

        const json = await res.json();
        const allTypes = json.data || [];

        const paymentTypes: PaymentGatewayType[] = allTypes
            .filter((it: any) =>
                it.category?.name === 'Pagos' ||
                it.integration_category?.name === 'Pagos'
            )
            .map((it: any) => ({
                id: it.id,
                name: it.name,
                code: it.code,
                image_url: it.image_url,
                is_active: it.is_active ?? true,
                in_development: it.in_development ?? false,
            }));

        return { success: true, data: paymentTypes };
    }
}
