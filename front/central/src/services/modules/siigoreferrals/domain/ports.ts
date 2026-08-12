import { GetSiigoReferralsParams, PaginatedSiigoReferralsResponse } from './types';

export interface ISiigoReferralsRepository {
    getReferrals(params?: GetSiigoReferralsParams): Promise<PaginatedSiigoReferralsResponse>;
}
