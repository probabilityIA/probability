import {
    Geozone,
    GeozonesListResponse,
    GetGeozonesParams,
    CreateGeozoneDTO,
    LookupParams,
    LookupResponse,
    BulkImportRequest,
    BulkImportResponse,
    DisplayFeatureCollection,
    GeozoneType,
    ProbabilityRequest,
    ProbabilityResult,
} from './types';

export interface IGeozoneRepository {
    list(params?: GetGeozonesParams): Promise<GeozonesListResponse>;
    getById(id: number, includeGeom: boolean, businessId?: number): Promise<Geozone>;
    create(data: CreateGeozoneDTO, businessId?: number): Promise<Geozone>;
    bulkImport(data: BulkImportRequest, businessId?: number): Promise<BulkImportResponse>;
    lookup(params: LookupParams): Promise<LookupResponse>;
    remove(id: number, businessId?: number): Promise<void>;
    getForDisplay(geozoneType: GeozoneType | '', zoom: number, bbox?: string): Promise<DisplayFeatureCollection>;
    probability(req: ProbabilityRequest): Promise<ProbabilityResult>;
    getOrderZone(orderId: string, businessId: number): Promise<Geozone | null>;
    probabilityByCarrier(orderId: string, businessId: number): Promise<ProbabilityResult[]>;
}
