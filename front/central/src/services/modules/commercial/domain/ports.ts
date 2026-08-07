import { GetProspectsParams, PaginatedProspectsResponse, ProspectStatsResponse } from './types';

export interface ICommercialRepository {
  getProspects(params?: GetProspectsParams): Promise<PaginatedProspectsResponse>;
  getProspectsStats(): Promise<ProspectStatsResponse>;
}
