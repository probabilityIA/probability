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

  @override
  Future<List<RecoveryChannel>> getRecoveryChannels(String email) async {
    final response = await _client.post(
      '/auth/recovery-channels',
      data: {'email': email},
    );
    final raw = response.data['data'];
    final list = raw is List ? raw : (raw is Map ? raw['channels'] as List? : null);
    return (list ?? [])
        .map((e) => RecoveryChannel.fromJson(Map<String, dynamic>.from(e)))
        .toList();
  }

  @override
  Future<SimpleAuthResponse> forgotPassword(String email, String channel) async {
    final response = await _client.post(
      '/auth/forgot-password',
      data: {'email': email, 'channel': channel},
    );
    return SimpleAuthResponse.fromJson(response.data);
  }

  @override
  Future<SimpleAuthResponse> verifyOtp(String email, String code) async {
    final response = await _client.post(
      '/auth/verify-otp',
      data: {'email': email, 'code': code},
    );
    return SimpleAuthResponse.fromJson(response.data);
  }

  @override
  Future<SimpleAuthResponse> resetPassword(String token, String newPassword) async {
    final response = await _client.post(
      '/auth/reset-password',
      data: {'token': token, 'new_password': newPassword},
    );
    return SimpleAuthResponse.fromJson(response.data);
  }
}
