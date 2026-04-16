import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/login/app/use_cases.dart';
import 'package:mobile_central/services/auth/login/domain/entities.dart';
import 'package:mobile_central/services/auth/login/domain/ports.dart';

// Manual mock that implements ILoginRepository with controllable behavior.
class MockLoginRepository implements ILoginRepository {
  // Captured call arguments for verification
  String? lastLoginEmail;
  String? lastLoginPassword;
  String? lastCurrentPassword;
  String? lastNewPassword;
  int? lastGeneratePasswordUserId;
  bool getRolesPermissionsCalled = false;

  // Configurable return values
  LoginSuccessResponse? loginResult;
  ChangePasswordResponse? changePasswordResult;
  GeneratePasswordResponse? generatePasswordResult;
  UserRolesPermissionsResponse? rolesPermissionsResult;

  // Configurable errors
  Exception? loginError;
  Exception? changePasswordError;
  Exception? generatePasswordError;
  Exception? rolesPermissionsError;

  @override
  Future<LoginSuccessResponse> login(String email, String password) async {
    lastLoginEmail = email;
    lastLoginPassword = password;
    if (loginError != null) throw loginError!;
    return loginResult!;
  }

  @override
  Future<ChangePasswordResponse> changePassword(
      String currentPassword, String newPassword) async {
    lastCurrentPassword = currentPassword;
    lastNewPassword = newPassword;
    if (changePasswordError != null) throw changePasswordError!;
    return changePasswordResult!;
  }

  @override
  Future<GeneratePasswordResponse> generatePassword({int? userId}) async {
    lastGeneratePasswordUserId = userId;
    if (generatePasswordError != null) throw generatePasswordError!;
    return generatePasswordResult!;
  }

  @override
  Future<UserRolesPermissionsResponse> getRolesPermissions() async {
    getRolesPermissionsCalled = true;
    if (rolesPermissionsError != null) throw rolesPermissionsError!;
    return rolesPermissionsResult!;
  }
}

void main() {
  late MockLoginRepository mockRepo;
  late LoginUseCases useCases;

  setUp(() {
    mockRepo = MockLoginRepository();
    useCases = LoginUseCases(mockRepo);
  });

  group('login', () {
    test('delegates to repository with correct arguments', () async {
      final expectedResponse = LoginSuccessResponse(
        success: true,
        data: LoginResponse(
          user: UserInfo(id: 1, name: 'Test', email: 'test@e.com', isActive: true),
          token: 'tok123',
          requirePasswordChange: false,
          businesses: [],
          isSuperAdmin: false,
        ),
      );
      mockRepo.loginResult = expectedResponse;

      final result = await useCases.login('test@e.com', 'pass123');

      expect(mockRepo.lastLoginEmail, 'test@e.com');
      expect(mockRepo.lastLoginPassword, 'pass123');
      expect(result.success, true);
      expect(result.data.token, 'tok123');
      expect(result.data.user.name, 'Test');
    });

    test('propagates exception from repository', () async {
      mockRepo.loginError = Exception('Invalid credentials');

      expect(
        () => useCases.login('bad@e.com', 'wrong'),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('changePassword', () {
    test('delegates to repository with correct arguments', () async {
      mockRepo.changePasswordResult = ChangePasswordResponse(
        success: true,
        message: 'Password changed',
      );

      final result = await useCases.changePassword('oldPass', 'newPass');

      expect(mockRepo.lastCurrentPassword, 'oldPass');
      expect(mockRepo.lastNewPassword, 'newPass');
      expect(result.success, true);
      expect(result.message, 'Password changed');
    });

    test('propagates exception from repository', () async {
      mockRepo.changePasswordError = Exception('Weak password');

      expect(
        () => useCases.changePassword('old', 'weak'),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('generatePassword', () {
    test('delegates to repository without userId', () async {
      mockRepo.generatePasswordResult = GeneratePasswordResponse(
        success: true,
        message: 'Generated',
        password: 'newPass123',
      );

      final result = await useCases.generatePassword();

      expect(mockRepo.lastGeneratePasswordUserId, isNull);
      expect(result.success, true);
      expect(result.password, 'newPass123');
    });

    test('delegates to repository with userId', () async {
      mockRepo.generatePasswordResult = GeneratePasswordResponse(
        success: true,
        message: 'Generated for user',
        password: 'userPass456',
      );

      final result = await useCases.generatePassword(userId: 42);

      expect(mockRepo.lastGeneratePasswordUserId, 42);
      expect(result.success, true);
      expect(result.password, 'userPass456');
    });

    test('propagates exception from repository', () async {
      mockRepo.generatePasswordError = Exception('User not found');

      expect(
        () => useCases.generatePassword(userId: 999),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getRolesPermissions', () {
    test('delegates to repository', () async {
      mockRepo.rolesPermissionsResult = UserRolesPermissionsResponse(
        isSuper: true,
        businessId: 5,
        businessName: 'Test Biz',
        role: 'admin',
        resources: [
          ResourcePermission(resource: 'orders', actions: ['read', 'write']),
        ],
      );

      final result = await useCases.getRolesPermissions();

      expect(mockRepo.getRolesPermissionsCalled, true);
      expect(result.isSuper, true);
      expect(result.businessId, 5);
      expect(result.resources.length, 1);
      expect(result.resources[0].resource, 'orders');
    });

    test('propagates exception from repository', () async {
      mockRepo.rolesPermissionsError = Exception('Unauthorized');

      expect(
        () => useCases.getRolesPermissions(),
        throwsA(isA<Exception>()),
      );
    });
  });
}
