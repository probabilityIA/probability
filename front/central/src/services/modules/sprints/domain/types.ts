export type SprintStatus = 'planned' | 'active' | 'closed';

export interface Sprint {
    id: number;
    name: string;
    goal: string;
    start_date: string;
    end_date: string;
    status: SprintStatus;
    created_by_id: number;
    created_by_name: string;
    ticket_count: number;
    done_count: number;
    created_at: string;
    updated_at: string;
}

export interface CreateSprintDTO {
    name: string;
    goal?: string;
    start_date: string;
    end_date: string;
    status?: SprintStatus;
}

export interface UpdateSprintDTO {
    name: string;
    goal: string;
    start_date: string;
    end_date: string;
    status: SprintStatus;
}

export interface ListSprintsParams {
    status?: string;
    page?: number;
    page_size?: number;
}

export interface PaginatedSprints {
    data: Sprint[];
    total: number;
    page: number;
    page_size: number;
    total_pages: number;
}

export const SPRINT_STATUSES: SprintStatus[] = ['planned', 'active', 'closed'];

export const SPRINT_STATUS_META: Record<SprintStatus, { label: string; color: string; bg: string; ring: string }> = {
    planned: { label: 'Planeado', color: 'text-blue-700 dark:text-blue-200', bg: 'bg-blue-100 dark:bg-blue-900/40', ring: 'ring-blue-300' },
    active: { label: 'Activo', color: 'text-emerald-700 dark:text-emerald-200', bg: 'bg-emerald-100 dark:bg-emerald-900/40', ring: 'ring-emerald-300' },
    closed: { label: 'Cerrado', color: 'text-gray-700 dark:text-gray-200', bg: 'bg-gray-200 dark:bg-gray-700', ring: 'ring-gray-300' },
};

export const SPRINT_STATUS_RANK: Record<SprintStatus, number> = {
    active: 0,
    planned: 1,
    closed: 2,
};
