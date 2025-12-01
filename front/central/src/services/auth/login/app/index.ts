import { ILoginRepository } from '../domain';
import {
    LoginRequest,
    ChangePasswordRequest,
    GeneratePasswordRequest,
    GenerateBusinessTokenRequest
} from '../infra/repository/mapper/request';
import {
    LoginSuccessResponse,
    UserRolesPermissionsSuccessResponse,
    ChangePasswordResponse,
    GeneratePasswordResponse,
    GenerateBusinessTokenSuccessResponse
} from '../infra/repository/mapper/response';

export class LoginUseCase {
    constructor(private readonly repository: ILoginRepository) { }

    async login(credentials: LoginRequest): Promise<LoginSuccessResponse> {
        return this.repository.login(credentials);
    }

    async changePassword(data: ChangePasswordRequest, token: string): Promise<ChangePasswordResponse> {
        return this.repository.changePassword(data, token);
    }

    async generatePassword(data: GeneratePasswordRequest, token: string): Promise<GeneratePasswordResponse> {
        return this.repository.generatePassword(data, token);
    }

    async getRolesPermissions(token: string): Promise<UserRolesPermissionsSuccessResponse> {
        return this.repository.getRolesPermissions(token);
    }

    async generateBusinessToken(data: GenerateBusinessTokenRequest, token: string): Promise<GenerateBusinessTokenSuccessResponse> {
        return this.repository.generateBusinessToken(data, token);
    }
}
