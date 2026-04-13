import { DashboardStatsResponse } from './types';

export interface IDashboardRepository {
    getStats(businessId?: number, integrationId?: number, weekStartDate?: Date, startDate?: Date, endDate?: Date): Promise<DashboardStatsResponse>;
}
