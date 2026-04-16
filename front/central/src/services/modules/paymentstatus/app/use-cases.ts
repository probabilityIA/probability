import { IPaymentStatusRepository } from '../domain/ports';
import { GetPaymentStatusesParams } from '../domain/types';

export class PaymentStatusUseCases {
    constructor(private repository: IPaymentStatusRepository) {}

    async getPaymentStatuses(params?: GetPaymentStatusesParams) {
        return this.repository.getPaymentStatuses(params);
    }
}
