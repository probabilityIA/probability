import type { SaveTourProgressInput, SkipAllTour, TourProgress } from './types';

export interface ITourRepository {
    listProgress(businessId?: number): Promise<TourProgress[]>;
    saveProgress(input: SaveTourProgressInput, businessId?: number): Promise<TourProgress>;
    resetTour(tourKey: string, businessId?: number): Promise<void>;
    resetAll(businessId?: number): Promise<void>;
    skipAll(tours: SkipAllTour[], businessId?: number): Promise<void>;
}
