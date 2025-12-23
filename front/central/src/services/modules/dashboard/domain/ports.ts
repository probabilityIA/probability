import { DashboardStatsResponse } from './types';

export interface IDashboardRepository {
    getStats(businessId?: number, integrationId?: number): Promise<DashboardStatsResponse>;
}
