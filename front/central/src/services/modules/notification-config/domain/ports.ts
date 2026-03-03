import { NotificationConfig, CreateConfigDTO, UpdateConfigDTO, ConfigFilter, SyncConfigsDTO, SyncConfigsResponse } from "./types";

export interface INotificationConfigRepository {
    create(dto: CreateConfigDTO, businessId?: number): Promise<NotificationConfig>;
    getById(id: number, businessId?: number): Promise<NotificationConfig>;
    update(id: number, dto: UpdateConfigDTO, businessId?: number): Promise<NotificationConfig>;
    delete(id: number, businessId?: number): Promise<void>;
    list(filter?: ConfigFilter): Promise<NotificationConfig[]>;
    syncByIntegration(dto: SyncConfigsDTO, businessId?: number): Promise<SyncConfigsResponse>;
}
