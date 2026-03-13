import 'entities.dart';

abstract class ILoginRepository {
  Future<LoginSuccessResponse> login(String email, String password);
  Future<ChangePasswordResponse> changePassword(
      String currentPassword, String newPassword);
  Future<GeneratePasswordResponse> generatePassword({int? userId});
  Future<UserRolesPermissionsResponse> getRolesPermissions();
}
