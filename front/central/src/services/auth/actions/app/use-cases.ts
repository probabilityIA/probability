import { IActionRepository } from '../domain/ports';
import {
    GetActionsParams,
    CreateActionDTO,
    UpdateActionDTO
} from '../domain/types';

export class ActionUseCases {
    constructor(private repository: IActionRepository) {}

    async getActions(params?: GetActionsParams) {
        return this.repository.getActions(params);
    }

    async getActionById(id: number) {
        return this.repository.getActionById(id);
    }

    async createAction(data: CreateActionDTO) {
        return this.repository.createAction(data);
    }

    async updateAction(id: number, data: UpdateActionDTO) {
        return this.repository.updateAction(id, data);
    }

    async deleteAction(id: number) {
        return this.repository.deleteAction(id);
    }
}
