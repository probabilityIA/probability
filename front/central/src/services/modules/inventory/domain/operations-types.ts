import { PaginatedResponse, PaginationParams } from './types';

export interface PutawayRule {
    id: number;
    business_id: number;
    product_id: string | null;
    category_id: number | null;
    target_zone_id: number;
    priority: number;
    strategy: string;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface PutawaySuggestion {
    id: number;
    business_id: number;
    product_id: string;
    recommended_location_id: number;
    quantity: number;
    status: string;
    rule_id: number | null;
    reason: string;
    actual_location_id: number | null;
    confirmed_at: string | null;
    confirmed_by_id: number | null;
    created_at: string;
    updated_at: string;
}

export interface ReplenishmentTask {
    id: number;
    business_id: number;
    product_id: string;
    warehouse_id: number;
    from_location_id: number | null;
    to_location_id: number | null;
    quantity: number;
    status: string;
    triggered_by: string;
    assigned_to_id: number | null;
    assigned_at: string | null;
    completed_at: string | null;
    notes: string;
    created_at: string;
    updated_at: string;
}

export interface CrossDockLink {
    id: number;
    business_id: number;
    inbound_shipment_id: number | null;
    outbound_order_id: string;
    product_id: string;
    quantity: number;
    status: string;
    executed_at: string | null;
    created_at: string;
    updated_at: string;
}

export interface ProductVelocity {
    id: number;
    business_id: number;
    product_id: string;
    warehouse_id: number;
    period: string;
    units_moved: number;
    rank: string;
    computed_at: string;
}

export interface SlottingRunResult {
    business_id: number;
    warehouse_id: number;
    period: string;
    total_scanned: number;
    velocities: ProductVelocity[];
}

export interface PutawaySuggestResult {
    suggestions: PutawaySuggestion[];
    unresolved_items: string[];
}

export interface ReplenishmentDetectResult {
    created: number;
    tasks: ReplenishmentTask[];
}

export interface CreatePutawayRuleDTO {
    product_id?: string | null;
    category_id?: number | null;
    target_zone_id: number;
    priority?: number;
    strategy?: string;
    is_active?: boolean;
}

export interface UpdatePutawayRuleDTO {
    product_id?: string | null;
    category_id?: number | null;
    target_zone_id?: number;
    priority?: number;
    strategy?: string;
    is_active?: boolean;
}

export interface SuggestPutawayInput {
    items: { product_id: string; quantity: number }[];
}

export interface ConfirmPutawayInput {
    actual_location_id: number;
}

export interface CreateReplenishmentTaskDTO {
    product_id: string;
    warehouse_id: number;
    from_location_id?: number | null;
    to_location_id?: number | null;
    quantity: number;
    triggered_by?: string;
    notes?: string;
}

export interface AssignReplenishmentInput {
    user_id: number;
}

export interface CompleteReplenishmentInput {
    notes?: string;
}

export interface CreateCrossDockLinkDTO {
    inbound_shipment_id?: number | null;
    outbound_order_id: string;
    product_id: string;
    quantity: number;
}

export interface RunSlottingInput {
    warehouse_id: number;
    period?: string;
}

export interface GetPutawayRulesParams extends PaginationParams {
    active_only?: boolean;
    business_id?: number;
}

export interface GetPutawaySuggestionsParams extends PaginationParams {
    status?: string;
    business_id?: number;
}

export interface GetReplenishmentTasksParams extends PaginationParams {
    warehouse_id?: number;
    status?: string;
    assigned_to?: number;
    business_id?: number;
}

export interface GetCrossDockLinksParams extends PaginationParams {
    outbound_order_id?: string;
    status?: string;
    business_id?: number;
}

export interface GetVelocitiesParams {
    warehouse_id: number;
    period?: string;
    rank?: string;
    limit?: number;
    business_id?: number;
}

export type PutawayRuleListResponse = PaginatedResponse<PutawayRule>;
export type PutawaySuggestionListResponse = PaginatedResponse<PutawaySuggestion>;
export type ReplenishmentTaskListResponse = PaginatedResponse<ReplenishmentTask>;
export type CrossDockLinkListResponse = PaginatedResponse<CrossDockLink>;
