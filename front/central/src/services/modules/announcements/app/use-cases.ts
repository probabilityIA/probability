import { IAnnouncementRepository } from '../domain/ports';
import {
    GetAnnouncementsParams,
    CreateAnnouncementDTO,
    UpdateAnnouncementDTO,
    RegisterViewDTO,
    ChangeStatusDTO,
} from '../domain/types';

export class AnnouncementUseCases {
    constructor(private repository: IAnnouncementRepository) {}

    async getAnnouncements(params?: GetAnnouncementsParams) {
        return this.repository.getAnnouncements(params);
    }

    async getAnnouncementById(id: number) {
        return this.repository.getAnnouncementById(id);
    }

    async createAnnouncement(data: CreateAnnouncementDTO) {
        return this.repository.createAnnouncement(data);
    }

    async updateAnnouncement(id: number, data: UpdateAnnouncementDTO) {
        return this.repository.updateAnnouncement(id, data);
    }

    async deleteAnnouncement(id: number) {
        return this.repository.deleteAnnouncement(id);
    }

    async getActiveAnnouncements(businessId?: number) {
        return this.repository.getActiveAnnouncements(businessId);
    }

    async registerView(announcementId: number, data: RegisterViewDTO) {
        return this.repository.registerView(announcementId, data);
    }

    async getStats(announcementId: number) {
        return this.repository.getStats(announcementId);
    }

    async listCategories() {
        return this.repository.listCategories();
    }

    async changeStatus(id: number, data: ChangeStatusDTO) {
        return this.repository.changeStatus(id, data);
    }

    async forceRedisplay(id: number) {
        return this.repository.forceRedisplay(id);
    }
}
