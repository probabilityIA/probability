import {
    Action,
    PaginatedResponse,
    SingleResponse,
    ActionResponse,
    GetActionsParams,
    CreateActionDTO,
    UpdateActionDTO
} from './types';

export interface IActionRepository {
    getActions(params?: GetActionsParams): Promise<PaginatedResponse<Action>>;
    getActionById(id: number): Promise<SingleResponse<Action>>;
    createAction(data: CreateActionDTO): Promise<SingleResponse<Action>>;
    updateAction(id: number, data: UpdateActionDTO): Promise<SingleResponse<Action>>;
    deleteAction(id: number): Promise<ActionResponse>;
}
