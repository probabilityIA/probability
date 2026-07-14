import { IPaymentMethodRepository } from '../domain/ports';

export class PaymentMethodUseCases {
    constructor(private repository: IPaymentMethodRepository) {}

    async getPaymentMethods() {
        return this.repository.getPaymentMethods();
    }
}
