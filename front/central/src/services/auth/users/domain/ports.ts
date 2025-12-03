import {
    User,
    PaginatedResponse,
    GetUsersParams,
    SingleResponse,
    CreateUserDTO,
    UpdateUserDTO,
    ActionResponse,
    AssignRolesDTO
} from './types';

export interface IUserRepository {
    getUsers(params?: GetUsersParams): Promise<PaginatedResponse<User>>;
    getUserById(id: number): Promise<SingleResponse<User>>;
    createUser(data: CreateUserDTO): Promise<SingleResponse<User>>;
    updateUser(id: number, data: UpdateUserDTO): Promise<SingleResponse<User>>;
    deleteUser(id: number): Promise<ActionResponse>;
    assignRoles(id: number, data: AssignRolesDTO): Promise<ActionResponse>;
}
