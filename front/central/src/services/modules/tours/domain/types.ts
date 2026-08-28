export type TourStatus = 'pending' | 'in_progress' | 'completed' | 'skipped';

export interface TourProgress {
    tour_key: string;
    version: number;
    status: TourStatus;
    step_index: number;
    completed_at?: string;
    updated_at?: string;
}

export interface SaveTourProgressInput {
    tour_key: string;
    version: number;
    status: TourStatus;
    step_index: number;
}

export type TourPlacement = 'top' | 'bottom' | 'left' | 'right' | 'center';

export interface SkipAllTour {
    tour_key: string;
    version: number;
}

export interface TourStep {
    id: string;
    title: string;
    body: string;
    target?: string;
    placement?: TourPlacement;
    route?: string;
    optional?: boolean;
}

export interface TourDefinition {
    key: string;
    version: number;
    title: string;
    routes: string[];
    resource?: string;
    autoStart: boolean;
    legacyStorageKey?: string;
    steps: TourStep[];
}
