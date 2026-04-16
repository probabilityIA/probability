import {
    AnnouncementInfo,
    AnnouncementStats,
    AnnouncementCategory,
    AnnouncementsListResponse,
    GetAnnouncementsParams,
    CreateAnnouncementDTO,
    UpdateAnnouncementDTO,
    RegisterViewDTO,
    ChangeStatusDTO,
    DeleteAnnouncementResponse,
    UploadImageResponse,
    DeleteImageResponse,
} from './types';

export interface IAnnouncementRepository {
    getAnnouncements(params?: GetAnnouncementsParams): Promise<AnnouncementsListResponse>;
    getAnnouncementById(id: number): Promise<AnnouncementInfo>;
    createAnnouncement(data: CreateAnnouncementDTO): Promise<AnnouncementInfo>;
    updateAnnouncement(id: number, data: UpdateAnnouncementDTO): Promise<AnnouncementInfo>;
    deleteAnnouncement(id: number): Promise<DeleteAnnouncementResponse>;
    getActiveAnnouncements(businessId?: number): Promise<AnnouncementInfo[]>;
    registerView(announcementId: number, data: RegisterViewDTO): Promise<void>;
    getStats(announcementId: number): Promise<AnnouncementStats>;
    listCategories(): Promise<AnnouncementCategory[]>;
    changeStatus(id: number, data: ChangeStatusDTO): Promise<void>;
    forceRedisplay(id: number): Promise<void>;
    uploadImage(announcementId: number, formData: FormData): Promise<UploadImageResponse>;
    deleteImage(announcementId: number, imageId: number): Promise<DeleteImageResponse>;
}
