'use server';

import { getAuthToken } from '@/shared/utils/server-auth';
import { AnnouncementApiRepository } from '../repository/api-repository';
import { AnnouncementUseCases } from '../../app/use-cases';
import {
    GetAnnouncementsParams,
    CreateAnnouncementDTO,
    UpdateAnnouncementDTO,
    RegisterViewDTO,
    ChangeStatusDTO,
} from '../../domain/types';

async function getUseCases() {
    const token = await getAuthToken();
    const repository = new AnnouncementApiRepository(token);
    return new AnnouncementUseCases(repository);
}

export const getAnnouncementsAction = async (params?: GetAnnouncementsParams) => {
    try {
        return await (await getUseCases()).getAnnouncements(params);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getAnnouncementByIdAction = async (id: number) => {
    try {
        return await (await getUseCases()).getAnnouncementById(id);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const createAnnouncementAction = async (data: CreateAnnouncementDTO) => {
    try {
        return await (await getUseCases()).createAnnouncement(data);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const updateAnnouncementAction = async (id: number, data: UpdateAnnouncementDTO) => {
    try {
        return await (await getUseCases()).updateAnnouncement(id, data);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const deleteAnnouncementAction = async (id: number) => {
    try {
        return await (await getUseCases()).deleteAnnouncement(id);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getActiveAnnouncementsAction = async (businessId?: number) => {
    try {
        return await (await getUseCases()).getActiveAnnouncements(businessId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const registerViewAction = async (announcementId: number, data: RegisterViewDTO) => {
    try {
        return await (await getUseCases()).registerView(announcementId, data);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const getStatsAction = async (announcementId: number) => {
    try {
        return await (await getUseCases()).getStats(announcementId);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const listCategoriesAction = async () => {
    try {
        return await (await getUseCases()).listCategories();
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const changeStatusAction = async (id: number, data: ChangeStatusDTO) => {
    try {
        return await (await getUseCases()).changeStatus(id, data);
    } catch (error: any) {
        throw new Error(error.message);
    }
};

export const forceRedisplayAction = async (id: number) => {
    try {
        return await (await getUseCases()).forceRedisplay(id);
    } catch (error: any) {
        throw new Error(error.message);
    }
};
