import { IDashboardRepository } from '../domain/ports';
import { DashboardStatsResponse } from '../domain/types';

export class DashboardUseCases {
    constructor(private repository: IDashboardRepository) { }

    async getStats(businessId?: number): Promise<DashboardStatsResponse> {
        return this.repository.getStats(businessId);
    }
}
