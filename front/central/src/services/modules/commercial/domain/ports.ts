import { GetProspectsParams, PaginatedProspectsResponse, ProspectStatsResponse, UpdateSeenResponse } from './types';

export interface ICommercialRepository {
  getProspects(params?: GetProspectsParams): Promise<PaginatedProspectsResponse>;
  getProspectsStats(): Promise<ProspectStatsResponse>;
  updateProspectSeen(id: number, seen: boolean): Promise<UpdateSeenResponse>;
}
