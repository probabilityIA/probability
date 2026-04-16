import '../domain/entities.dart';
import '../domain/ports.dart';

class LoginUseCases {
  final ILoginRepository _repository;

  LoginUseCases(this._repository);

  Future<LoginSuccessResponse> login(String email, String password) {
    return _repository.login(email, password);
  }

  Future<ChangePasswordResponse> changePassword(
      String currentPassword, String newPassword) {
    return _repository.changePassword(currentPassword, newPassword);
  }

  Future<GeneratePasswordResponse> generatePassword({int? userId}) {
    return _repository.generatePassword(userId: userId);
  }

  Future<UserRolesPermissionsResponse> getRolesPermissions() {
    return _repository.getRolesPermissions();
  }
}
