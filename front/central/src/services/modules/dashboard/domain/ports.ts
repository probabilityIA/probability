import { DashboardStatsResponse } from './types';

export interface IDashboardRepository {
    getStats(businessId?: number): Promise<DashboardStatsResponse>;
}
