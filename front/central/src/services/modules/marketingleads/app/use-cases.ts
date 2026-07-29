import { IMarketingLeadsRepository } from '../domain/ports';
import { GetMarketingLeadsParams } from '../domain/types';

export class MarketingLeadsUseCases {
    constructor(private repository: IMarketingLeadsRepository) { }

    async getLeads(params?: GetMarketingLeadsParams) {
        return this.repository.getLeads(params);
    }
}
