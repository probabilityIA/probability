import {
    createConfigAction,
    updateConfigAction,
    deleteConfigAction,
    listConfigsAction,
    getConfigAction
} from "../infra/actions";
import { CreateConfigDTO, UpdateConfigDTO, ConfigFilter } from "../domain/types";

export class NotificationConfigUseCases {
    async createConfig(dto: CreateConfigDTO) {
        return createConfigAction(dto);
    }

    async updateConfig(id: number, dto: UpdateConfigDTO) {
        return updateConfigAction(id, dto);
    }

    async deleteConfig(id: number) {
        return deleteConfigAction(id);
    }

    async listConfigs(filter?: ConfigFilter) {
        return listConfigsAction(filter);
    }

    async getConfig(id: number) {
        return getConfigAction(id);
    }
}

export const notificationConfigUseCases = new NotificationConfigUseCases();
