import { IDashboardRepository } from '../domain/ports';
import { DashboardStatsResponse } from '../domain/types';

export class DashboardUseCases {
    constructor(private repository: IDashboardRepository) { }

    async getStats(businessId?: number, integrationId?: number, weekStartDate?: Date, startDate?: Date, endDate?: Date): Promise<DashboardStatsResponse> {
        return this.repository.getStats(businessId, integrationId, weekStartDate, startDate, endDate);
    }
}
