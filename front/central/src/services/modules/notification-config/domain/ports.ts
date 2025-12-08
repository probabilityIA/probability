import { NotificationConfig, CreateConfigDTO, UpdateConfigDTO, ConfigFilter } from "./types";

export interface INotificationConfigRepository {
    create(dto: CreateConfigDTO): Promise<NotificationConfig>;
    getById(id: number): Promise<NotificationConfig>;
    update(id: number, dto: UpdateConfigDTO): Promise<NotificationConfig>;
    delete(id: number): Promise<void>;
    list(filter?: ConfigFilter): Promise<NotificationConfig[]>;
}
