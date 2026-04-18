import { PaginatedResponse, PaginationParams } from './types';

export interface CycleCountPlan {
    id: number;
    business_id: number;
    warehouse_id: number;
    name: string;
    strategy: string;
    frequency_days: number;
    next_run_at: string | null;
    is_active: boolean;
    created_at: string;
    updated_at: string;
}

export interface CycleCountTask {
    id: number;
    plan_id: number;
    business_id: number;
    warehouse_id: number;
    scope_type: string;
    scope_id: number | null;
    status: string;
    assigned_to_id: number | null;
    started_at: string | null;
    finished_at: string | null;
    created_at: string;
    updated_at: string;
}

export interface CycleCountLine {
    id: number;
    task_id: number;
    business_id: number;
    product_id: string;
    location_id: number | null;
    lot_id: number | null;
    expected_qty: number;
    counted_qty: number | null;
    variance: number;
    status: string;
    created_at: string;
    updated_at: string;
}

export interface InventoryDiscrepancy {
    id: number;
    task_id: number;
    line_id: number;
    business_id: number;
    status: string;
    resolution_movement_id: number | null;
    reviewed_by_id: number | null;
    reviewed_at: string | null;
    notes: string;
    created_at: string;
    updated_at: string;
}

export interface KardexEntry {
    movement_id: number;
    movement_type_code: string;
    movement_type_name: string;
    quantity: number;
    previous_qty: number;
    new_qty: number;
    running_balance: number;
    reason: string;
    reference_type: string | null;
    reference_id: string | null;
    location_id: number | null;
    lot_id: number | null;
}

export interface KardexExportResult {
    business_id: number;
    product_id: string;
    warehouse_id: number;
    entries: KardexEntry[];
    total_in: number;
    total_out: number;
    final_balance: number;
}

export interface GenerateCountTaskResult {
    task: CycleCountTask;
    lines: CycleCountLine[];
}

export interface SubmitCountLineResult {
    line: CycleCountLine;
    discrepancy?: InventoryDiscrepancy;
}

export interface CreateCycleCountPlanDTO {
    warehouse_id: number;
    name: string;
    strategy?: string;
    frequency_days?: number;
    next_run_at?: string | null;
    is_active?: boolean;
}

export interface UpdateCycleCountPlanDTO {
    warehouse_id?: number;
    name?: string;
    strategy?: string;
    frequency_days?: number;
    next_run_at?: string | null;
    is_active?: boolean;
}

export interface GenerateCountTaskInput {
    plan_id: number;
    scope_type?: string;
    scope_id?: number | null;
}

export interface StartCountTaskInput {
    user_id: number;
}

export interface SubmitCountLineInput {
    counted_qty: number;
}

export interface ApproveDiscrepancyInput {
    notes?: string;
}

export interface RejectDiscrepancyInput {
    reason?: string;
}

export interface KardexQueryInput {
    product_id: string;
    warehouse_id: number;
    from?: string;
    to?: string;
    business_id?: number;
}

export interface GetCountPlansParams extends PaginationParams {
    warehouse_id?: number;
    active_only?: boolean;
    business_id?: number;
}

export interface GetCountTasksParams extends PaginationParams {
    warehouse_id?: number;
    plan_id?: number;
    status?: string;
    business_id?: number;
}

export interface GetCountLinesParams extends PaginationParams {
    task_id: number;
    status?: string;
    business_id?: number;
}

export interface GetDiscrepanciesParams extends PaginationParams {
    task_id?: number;
    status?: string;
    business_id?: number;
}

export type CycleCountPlanListResponse = PaginatedResponse<CycleCountPlan>;
export type CycleCountTaskListResponse = PaginatedResponse<CycleCountTask>;
export type CycleCountLineListResponse = PaginatedResponse<CycleCountLine>;
export type InventoryDiscrepancyListResponse = PaginatedResponse<InventoryDiscrepancy>;
