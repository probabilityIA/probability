import { PaymentGatewayTypesResponse } from './types';

export interface IPayGatewayRepository {
    listPaymentGatewayTypes(): Promise<PaymentGatewayTypesResponse>;
}
