import { IUserRepository } from '../domain/ports';
import {
    GetUsersParams,
    CreateUserDTO,
    UpdateUserDTO,
    AssignRolesDTO
} from '../domain/types';

export class UserUseCases {
    constructor(private repository: IUserRepository) { }

    async getUsers(params?: GetUsersParams) {
        return this.repository.getUsers(params);
    }

    async getUserById(id: number) {
        return this.repository.getUserById(id);
    }

    async createUser(data: CreateUserDTO) {
        return this.repository.createUser(data);
    }

    async updateUser(id: number, data: UpdateUserDTO) {
        return this.repository.updateUser(id, data);
    }

    async deleteUser(id: number) {
        return this.repository.deleteUser(id);
    }

    async assignRoles(id: number, data: AssignRolesDTO) {
        return this.repository.assignRoles(id, data);
    }
}
