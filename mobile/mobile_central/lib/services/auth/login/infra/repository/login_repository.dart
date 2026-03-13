import '../../../../../core/network/api_client.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class LoginApiRepository implements ILoginRepository {
  final ApiClient _client;

  LoginApiRepository(this._client);

  @override
  Future<LoginSuccessResponse> login(String email, String password) async {
    final response = await _client.post(
      '/auth/login',
      data: {'email': email, 'password': password},
    );
    return LoginSuccessResponse.fromJson(response.data);
  }

  @override
  Future<ChangePasswordResponse> changePassword(
      String currentPassword, String newPassword) async {
    final response = await _client.post(
      '/auth/change-password',
      data: {
        'current_password': currentPassword,
        'new_password': newPassword,
      },
    );
    return ChangePasswordResponse.fromJson(response.data);
  }

  @override
  Future<GeneratePasswordResponse> generatePassword({int? userId}) async {
    final data = <String, dynamic>{};
    if (userId != null) data['user_id'] = userId;
    final response = await _client.post('/auth/generate-password', data: data);
    return GeneratePasswordResponse.fromJson(response.data);
  }

  @override
  Future<UserRolesPermissionsResponse> getRolesPermissions() async {
    final response = await _client.get('/auth/roles-permissions');
    return UserRolesPermissionsResponse.fromJson(response.data);
  }
}
