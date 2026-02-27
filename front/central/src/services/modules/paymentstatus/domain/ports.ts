import { GetPaymentStatusesParams, PaymentStatusesResponse } from './types';

export interface IPaymentStatusRepository {
    getPaymentStatuses(params?: GetPaymentStatusesParams): Promise<PaymentStatusesResponse>;
}
