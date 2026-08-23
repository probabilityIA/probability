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

  Future<List<RecoveryChannel>> getRecoveryChannels(String email) {
    return _repository.getRecoveryChannels(email);
  }

  Future<SimpleAuthResponse> forgotPassword(String email, String channel) {
    return _repository.forgotPassword(email, channel);
  }

  Future<SimpleAuthResponse> verifyOtp(String email, String code) {
    return _repository.verifyOtp(email, code);
  }

  Future<SimpleAuthResponse> resetPassword(String token, String newPassword) {
    return _repository.resetPassword(token, newPassword);
  }
}
