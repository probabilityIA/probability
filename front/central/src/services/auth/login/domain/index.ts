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

export interface ILoginRepository {
    login(credentials: LoginRequest): Promise<LoginSuccessResponse>;
    changePassword(data: ChangePasswordRequest, token: string): Promise<ChangePasswordResponse>;
    generatePassword(data: GeneratePasswordRequest, token: string): Promise<GeneratePasswordResponse>;
    getRolesPermissions(token: string): Promise<UserRolesPermissionsSuccessResponse>;

}
