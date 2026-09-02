import {
    Sprint,
    PaginatedSprints,
    ListSprintsParams,
    CreateSprintDTO,
    UpdateSprintDTO,
    SprintStatus,
} from './types';

export interface ISprintRepository {
    list(params?: ListSprintsParams): Promise<PaginatedSprints>;
    get(id: number): Promise<Sprint>;
    create(data: CreateSprintDTO): Promise<Sprint>;
    update(id: number, data: UpdateSprintDTO): Promise<Sprint>;
    remove(id: number): Promise<void>;
    changeStatus(id: number, status: SprintStatus): Promise<Sprint>;
}
