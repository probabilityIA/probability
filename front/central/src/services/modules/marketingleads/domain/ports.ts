import { GetMarketingLeadsParams, PaginatedLeadsResponse } from './types';

export interface IMarketingLeadsRepository {
    getLeads(params?: GetMarketingLeadsParams): Promise<PaginatedLeadsResponse>;
}
