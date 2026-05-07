import { IGeozoneRepository } from '../domain/ports';
import {
    GetGeozonesParams,
    CreateGeozoneDTO,
    LookupParams,
    BulkImportRequest,
    GeozoneType,
    ProbabilityRequest,
} from '../domain/types';

export class GeozoneUseCases {
    constructor(private repo: IGeozoneRepository) {}

    list(params?: GetGeozonesParams) { return this.repo.list(params); }
    getById(id: number, includeGeom = true, businessId?: number) { return this.repo.getById(id, includeGeom, businessId); }
    create(data: CreateGeozoneDTO, businessId?: number) { return this.repo.create(data, businessId); }
    bulkImport(data: BulkImportRequest, businessId?: number) { return this.repo.bulkImport(data, businessId); }
    lookup(params: LookupParams) { return this.repo.lookup(params); }
    remove(id: number, businessId?: number) { return this.repo.remove(id, businessId); }
    getForDisplay(type: GeozoneType | '', zoom: number, bbox?: string) { return this.repo.getForDisplay(type, zoom, bbox); }
    probability(req: ProbabilityRequest) { return this.repo.probability(req); }
}
