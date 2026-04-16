import { IResourceRepository } from '../domain/ports';
import {
    GetResourcesParams,
    CreateResourceDTO,
    UpdateResourceDTO
} from '../domain/types';

export class ResourceUseCases {
    constructor(private repository: IResourceRepository) { }

    async getResources(params?: GetResourcesParams) {
        return this.repository.getResources(params);
    }

    async getResourceById(id: number) {
        return this.repository.getResourceById(id);
    }

    async createResource(data: CreateResourceDTO) {
        return this.repository.createResource(data);
    }

    async updateResource(id: number, data: UpdateResourceDTO) {
        return this.repository.updateResource(id, data);
    }

    async deleteResource(id: number) {
        return this.repository.deleteResource(id);
    }
}
