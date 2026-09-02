import { ISprintRepository } from '../domain/ports';
import { CreateSprintDTO, UpdateSprintDTO, ListSprintsParams, SprintStatus } from '../domain/types';

export class SprintUseCases {
    constructor(private repo: ISprintRepository) {}

    list(params?: ListSprintsParams) { return this.repo.list(params); }
    get(id: number) { return this.repo.get(id); }
    create(data: CreateSprintDTO) { return this.repo.create(data); }
    update(id: number, data: UpdateSprintDTO) { return this.repo.update(id, data); }
    remove(id: number) { return this.repo.remove(id); }
    changeStatus(id: number, status: SprintStatus) { return this.repo.changeStatus(id, status); }
}
