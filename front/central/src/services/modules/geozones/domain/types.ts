export type GeozoneType = 'country' | 'state' | 'city' | 'admin_district' | 'locality' | 'neighborhood' | 'barrio' | 'custom';

export interface Geozone {
    id: number;
    business_id: number;
    parent_id: number | null;
    type: GeozoneType;
    code: string | null;
    name: string;
    geometry?: GeoJSONGeometry | null;
    centroid?: GeoJSONPoint | null;
    properties: Record<string, unknown>;
    is_active: boolean;
}

export interface GeoJSONPoint {
    type: 'Point';
    coordinates: [number, number];
}

export interface GeoJSONPolygon {
    type: 'Polygon';
    coordinates: number[][][];
}

export interface GeoJSONMultiPolygon {
    type: 'MultiPolygon';
    coordinates: number[][][][];
}

export type GeoJSONGeometry = GeoJSONPoint | GeoJSONPolygon | GeoJSONMultiPolygon;

export interface CreateGeozoneDTO {
    parent_id?: number | null;
    type: GeozoneType;
    code?: string | null;
    name: string;
    geometry: GeoJSONGeometry;
    properties?: Record<string, unknown>;
}

export interface GetGeozonesParams {
    page?: number;
    page_size?: number;
    type?: GeozoneType | '';
    parent_id?: number;
    code?: string;
    search?: string;
    include_geometry?: boolean;
    business_id?: number;
}

export interface GeozonesListResponse {
    data: Geozone[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export interface LookupParams {
    lat: number;
    lng: number;
    type?: GeozoneType;
    business_id?: number;
}

export interface LookupResponse {
    data: Geozone[];
}

export interface BulkImportFeature {
    type: 'Feature';
    geometry: GeoJSONGeometry;
    properties: {
        type: GeozoneType;
        code?: string | null;
        name: string;
        parent_code?: string | null;
    };
}

export interface BulkImportRequest {
    type: 'FeatureCollection';
    features: BulkImportFeature[];
}

export interface BulkImportResponse {
    created: number;
    skipped: number;
    errors?: string[];
}

export interface DisplayFeature {
    type: 'Feature';
    geometry: GeoJSONGeometry;
    properties: {
        id: number;
        type: GeozoneType;
        code?: string | null;
        name: string;
    };
}

export interface DisplayFeatureCollection {
    type: 'FeatureCollection';
    features: DisplayFeature[];
}

export interface DrillState {
    level: 'country' | 'state' | 'city' | 'admin_district' | 'neighborhood';
    state?: { id: number; name: string };
    city?: { id: number; name: string };
    adminDistrict?: { id: number; name: string };
    neighborhood?: { id: number; name: string };
}
