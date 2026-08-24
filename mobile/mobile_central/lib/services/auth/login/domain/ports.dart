import 'entities.dart';

abstract class ILoginRepository {
  Future<LoginSuccessResponse> login(String email, String password);
  Future<ChangePasswordResponse> changePassword(
      String currentPassword, String newPassword);
  Future<GeneratePasswordResponse> generatePassword({int? userId});
  Future<UserRolesPermissionsResponse> getRolesPermissions();
  Future<List<RecoveryChannel>> getRecoveryChannels(String email);
  Future<SimpleAuthResponse> forgotPassword(String email, String channel);
  Future<SimpleAuthResponse> verifyOtp(String email, String code);
  Future<SimpleAuthResponse> resetPassword(String token, String newPassword);
}
