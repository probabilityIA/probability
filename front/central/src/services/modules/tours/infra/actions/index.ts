'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { TourApiRepository } from '../repository/api-repository';
import type { SaveTourProgressInput, SkipAllTour, TourProgress } from '../../domain/types';

async function getRepository() {
    const token = await getAuthToken();
    return new TourApiRepository(token);
}

export const listTourProgressAction = async (businessId?: number): Promise<TourProgress[]> => {
    const repository = await getRepository();
    return repository.listProgress(businessId);
};

export const saveTourProgressAction = async (
    input: SaveTourProgressInput,
    businessId?: number,
): Promise<TourProgress> => {
    const repository = await getRepository();
    return repository.saveProgress(input, businessId);
};

export const resetTourAction = async (tourKey: string, businessId?: number): Promise<void> => {
    const repository = await getRepository();
    return repository.resetTour(tourKey, businessId);
};

export const resetAllToursAction = async (businessId?: number): Promise<void> => {
    const repository = await getRepository();
    return repository.resetAll(businessId);
};

export const skipAllToursAction = async (tours: SkipAllTour[], businessId?: number): Promise<void> => {
    const repository = await getRepository();
    return repository.skipAll(tours, businessId);
};
