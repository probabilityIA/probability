import { PaymentMethodsResponse } from './types';

export interface IPaymentMethodRepository {
    getPaymentMethods(): Promise<PaymentMethodsResponse>;
}
