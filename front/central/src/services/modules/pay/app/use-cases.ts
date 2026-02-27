import { IPayGatewayRepository } from '../domain/ports';
import { PaymentGatewayType } from '../domain/types';

export class PayGatewayUseCases {
    constructor(private repo: IPayGatewayRepository) {}

    async getPaymentGatewayTypes(): Promise<PaymentGatewayType[]> {
        const response = await this.repo.listPaymentGatewayTypes();
        if (!response.success) return [];
        return response.data;
    }
}
