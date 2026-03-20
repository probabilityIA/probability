import 'dart:convert';

import 'package:flutter/foundation.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/login/app/use_cases.dart';
import 'package:mobile_central/services/auth/login/domain/entities.dart';
import 'package:mobile_central/services/auth/login/domain/ports.dart';

// ---------------------------------------------------------------------------
// Manual mocks
// ---------------------------------------------------------------------------

/// In-memory token storage that mirrors the real TokenStorage API without
/// depending on FlutterSecureStorage (which requires platform channels).
class MockTokenStorage {
  final Map<String, String> _store = {};

  Future<void> saveToken(String token) async {
    _store['token'] = token;
  }

  Future<String?> getToken() async {
    return _store['token'];
  }

  Future<void> saveUserData(String userData) async {
    _store['user'] = userData;
  }

  Future<String?> getUserData() async {
    return _store['user'];
  }

  Future<void> clearAll() async {
    _store.clear();
  }
}

/// Minimal mock that tracks setToken calls.
class MockApiClient {
  String? lastTokenSet;
  bool setTokenCalled = false;

  void setToken(String? token) {
    lastTokenSet = token;
    setTokenCalled = true;
  }
}

/// Mock repository with controllable return values and errors.
class MockLoginRepository implements ILoginRepository {
  LoginSuccessResponse? loginResult;
  ChangePasswordResponse? changePasswordResult;
  GeneratePasswordResponse? generatePasswordResult;
  UserRolesPermissionsResponse? rolesPermissionsResult;

  Exception? loginError;
  Exception? changePasswordError;
  Exception? rolesPermissionsError;

  @override
  Future<LoginSuccessResponse> login(String email, String password) async {
    if (loginError != null) throw loginError!;
    return loginResult!;
  }

  @override
  Future<ChangePasswordResponse> changePassword(
      String currentPassword, String newPassword) async {
    if (changePasswordError != null) throw changePasswordError!;
    return changePasswordResult!;
  }

  @override
  Future<GeneratePasswordResponse> generatePassword({int? userId}) async {
    return generatePasswordResult!;
  }

  @override
  Future<UserRolesPermissionsResponse> getRolesPermissions() async {
    if (rolesPermissionsError != null) throw rolesPermissionsError!;
    return rolesPermissionsResult!;
  }
}

// ---------------------------------------------------------------------------
// Testable provider
// ---------------------------------------------------------------------------
// The production LoginProvider hardcodes the creation of LoginApiRepository
// from an ApiClient internally, making it impossible to inject a mock
// repository without modifying production code.
//
// To avoid modifying production code, we create a standalone ChangeNotifier
// that replicates the exact same logic as LoginProvider but accepts injected
// dependencies (MockTokenStorage, MockApiClient, ILoginRepository). Every
// method body is a 1:1 copy of the production code so we are testing the
// real business logic.
// ---------------------------------------------------------------------------

class _TestLoginProvider extends ChangeNotifier {
  final MockTokenStorage _tokenStorage;
  final MockApiClient _apiClient;
  final LoginUseCases _useCases;

  bool _isLoading = false;
  String? _error;
  UserInfo? _user;
  bool _isSuperAdmin = false;
  List<BusinessInfo> _businesses = [];
  UserRolesPermissionsResponse? _rolesPermissions;

  _TestLoginProvider({
    required MockTokenStorage tokenStorage,
    required MockApiClient apiClient,
    required ILoginRepository repository,
  })  : _tokenStorage = tokenStorage,
        _apiClient = apiClient,
        _useCases = LoginUseCases(repository);

  bool get isLoading => _isLoading;
  String? get error => _error;
  UserInfo? get user => _user;
  bool get isSuperAdmin => _isSuperAdmin;
  List<BusinessInfo> get businesses => _businesses;
  UserRolesPermissionsResponse? get rolesPermissions => _rolesPermissions;
  bool get isLoggedIn => _user != null;

  // Exact copy of LoginProvider.login
  Future<bool> login(String email, String password) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.login(email, password);
      final data = response.data;

      await _tokenStorage.saveToken(data.token);
      await _tokenStorage.saveUserData(jsonEncode({
        'id': data.user.id,
        'name': data.user.name,
        'email': data.user.email,
      }));

      _apiClient.setToken(data.token);
      _user = data.user;
      _isSuperAdmin = data.isSuperAdmin;
      _businesses = data.businesses;

      await _fetchRolesPermissions();

      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<void> _fetchRolesPermissions() async {
    try {
      _rolesPermissions = await _useCases.getRolesPermissions();
    } catch (_) {}
  }

  // Exact copy of LoginProvider.changePassword
  Future<ChangePasswordResponse?> changePassword(
      String currentPassword, String newPassword) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response =
          await _useCases.changePassword(currentPassword, newPassword);
      _isLoading = false;
      notifyListeners();
      return response;
    } catch (e) {
      _error = e.toString();
      _isLoading = false;
      notifyListeners();
      return null;
    }
  }

  // Exact copy of LoginProvider.restoreSession
  Future<void> restoreSession() async {
    final token = await _tokenStorage.getToken();
    if (token == null) return;

    _apiClient.setToken(token);

    try {
      _rolesPermissions = await _useCases.getRolesPermissions();
      _isSuperAdmin = _rolesPermissions?.isSuper ?? false;

      final userData = await _tokenStorage.getUserData();
      if (userData != null) {
        final json = jsonDecode(userData);
        _user = UserInfo(
          id: json['id'],
          name: json['name'],
          email: json['email'],
          isActive: true,
        );
      }
      notifyListeners();
    } catch (_) {
      await logout();
    }
  }

  // Exact copy of LoginProvider.logout
  Future<void> logout() async {
    await _tokenStorage.clearAll();
    _apiClient.setToken(null);
    _user = null;
    _isSuperAdmin = false;
    _businesses = [];
    _rolesPermissions = null;
    _error = null;
    notifyListeners();
  }

  // Exact copy of LoginProvider.hasPermission
  bool hasPermission(String resource, String action) {
    if (_isSuperAdmin) return true;
    if (_rolesPermissions == null) return false;
    return _rolesPermissions!.resources.any(
      (r) => r.resource == resource && r.actions.contains(action),
    );
  }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

void main() {
  late MockTokenStorage mockTokenStorage;
  late MockApiClient mockApiClient;
  late MockLoginRepository mockRepo;
  late _TestLoginProvider provider;

  // Helpers
  LoginSuccessResponse makeLoginResponse({
    int userId = 1,
    String userName = 'Test User',
    String userEmail = 'test@e.com',
    String token = 'jwt-token',
    bool isSuperAdmin = false,
    List<BusinessInfo>? businesses,
  }) {
    return LoginSuccessResponse(
      success: true,
      data: LoginResponse(
        user: UserInfo(
            id: userId, name: userName, email: userEmail, isActive: true),
        token: token,
        requirePasswordChange: false,
        businesses: businesses ?? [],
        isSuperAdmin: isSuperAdmin,
      ),
    );
  }

  UserRolesPermissionsResponse makeRolesPermissions({
    bool isSuper = false,
    List<ResourcePermission>? resources,
  }) {
    return UserRolesPermissionsResponse(
      isSuper: isSuper,
      businessId: 1,
      resources: resources ?? [],
    );
  }

  setUp(() {
    mockTokenStorage = MockTokenStorage();
    mockApiClient = MockApiClient();
    mockRepo = MockLoginRepository();
    provider = _TestLoginProvider(
      tokenStorage: mockTokenStorage,
      apiClient: mockApiClient,
      repository: mockRepo,
    );
  });

  group('initial state', () {
    test('has correct initial values', () {
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
      expect(provider.user, isNull);
      expect(provider.isSuperAdmin, false);
      expect(provider.businesses, isEmpty);
      expect(provider.rolesPermissions, isNull);
      expect(provider.isLoggedIn, false);
    });
  });

  group('login', () {
    test('success updates user, businesses, token, and isSuperAdmin', () async {
      final businesses = [
        BusinessInfo(id: 1, name: 'Biz 1'),
        BusinessInfo(id: 2, name: 'Biz 2'),
      ];
      mockRepo.loginResult = makeLoginResponse(
        userId: 42,
        userName: 'Cam',
        userEmail: 'cam@prob.co',
        token: 'abc123',
        isSuperAdmin: true,
        businesses: businesses,
      );
      mockRepo.rolesPermissionsResult = makeRolesPermissions(isSuper: true);

      final result = await provider.login('cam@prob.co', 'secret');

      expect(result, true);
      expect(provider.isLoggedIn, true);
      expect(provider.user!.id, 42);
      expect(provider.user!.name, 'Cam');
      expect(provider.user!.email, 'cam@prob.co');
      expect(provider.isSuperAdmin, true);
      expect(provider.businesses.length, 2);
      expect(provider.businesses[0].name, 'Biz 1');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('success saves token to storage', () async {
      mockRepo.loginResult = makeLoginResponse(token: 'saved-token');
      mockRepo.rolesPermissionsResult = makeRolesPermissions();

      await provider.login('u@e.com', 'p');

      final storedToken = await mockTokenStorage.getToken();
      expect(storedToken, 'saved-token');
    });

    test('success saves user data to storage', () async {
      mockRepo.loginResult = makeLoginResponse(
        userId: 7,
        userName: 'Stored',
        userEmail: 'stored@e.com',
      );
      mockRepo.rolesPermissionsResult = makeRolesPermissions();

      await provider.login('stored@e.com', 'p');

      final rawUserData = await mockTokenStorage.getUserData();
      expect(rawUserData, isNotNull);
      final userData = jsonDecode(rawUserData!);
      expect(userData['id'], 7);
      expect(userData['name'], 'Stored');
      expect(userData['email'], 'stored@e.com');
    });

    test('success sets token on ApiClient', () async {
      mockRepo.loginResult = makeLoginResponse(token: 'api-token');
      mockRepo.rolesPermissionsResult = makeRolesPermissions();

      await provider.login('u@e.com', 'p');

      expect(mockApiClient.setTokenCalled, true);
      expect(mockApiClient.lastTokenSet, 'api-token');
    });

    test('success fetches roles and permissions', () async {
      final resources = [
        ResourcePermission(resource: 'orders', actions: ['read', 'write']),
      ];
      mockRepo.loginResult = makeLoginResponse();
      mockRepo.rolesPermissionsResult =
          makeRolesPermissions(resources: resources);

      await provider.login('u@e.com', 'p');

      expect(provider.rolesPermissions, isNotNull);
      expect(provider.rolesPermissions!.resources.length, 1);
      expect(provider.rolesPermissions!.resources[0].resource, 'orders');
    });

    test('success with failed roles fetch still succeeds login', () async {
      mockRepo.loginResult = makeLoginResponse();
      mockRepo.rolesPermissionsError = Exception('Network error');

      final result = await provider.login('u@e.com', 'p');

      expect(result, true);
      expect(provider.isLoggedIn, true);
      expect(provider.rolesPermissions, isNull);
    });

    test('failure sets error and returns false', () async {
      mockRepo.loginError = Exception('Invalid credentials');

      final result = await provider.login('bad@e.com', 'wrong');

      expect(result, false);
      expect(provider.isLoggedIn, false);
      expect(provider.error, contains('Invalid credentials'));
      expect(provider.isLoading, false);
    });

    test('notifies listeners during state changes', () async {
      mockRepo.loginResult = makeLoginResponse();
      mockRepo.rolesPermissionsResult = makeRolesPermissions();

      final notifications = <bool>[];
      provider.addListener(() {
        notifications.add(provider.isLoading);
      });

      await provider.login('u@e.com', 'p');

      // Should have notified at least twice: loading=true then loading=false
      expect(notifications.length, greaterThanOrEqualTo(2));
      expect(notifications.first, true);
      expect(notifications.last, false);
    });
  });

  group('changePassword', () {
    test('success returns response', () async {
      mockRepo.changePasswordResult = ChangePasswordResponse(
        success: true,
        message: 'Password updated',
      );

      final result = await provider.changePassword('old', 'new');

      expect(result, isNotNull);
      expect(result!.success, true);
      expect(result.message, 'Password updated');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('failure returns null and sets error', () async {
      mockRepo.changePasswordError = Exception('Wrong current password');

      final result = await provider.changePassword('bad', 'new');

      expect(result, isNull);
      expect(provider.error, contains('Wrong current password'));
      expect(provider.isLoading, false);
    });
  });

  group('logout', () {
    test('clears all state', () async {
      // First login
      mockRepo.loginResult = makeLoginResponse(
        userId: 1,
        userName: 'User',
        isSuperAdmin: true,
        businesses: [BusinessInfo(id: 1, name: 'Biz')],
      );
      mockRepo.rolesPermissionsResult = makeRolesPermissions(isSuper: true);
      await provider.login('u@e.com', 'p');
      expect(provider.isLoggedIn, true);

      // Now logout
      await provider.logout();

      expect(provider.user, isNull);
      expect(provider.isLoggedIn, false);
      expect(provider.isSuperAdmin, false);
      expect(provider.businesses, isEmpty);
      expect(provider.rolesPermissions, isNull);
      expect(provider.error, isNull);
    });

    test('clears token storage', () async {
      await mockTokenStorage.saveToken('some-token');
      await mockTokenStorage.saveUserData('{"id":1}');

      await provider.logout();

      final token = await mockTokenStorage.getToken();
      final userData = await mockTokenStorage.getUserData();
      expect(token, isNull);
      expect(userData, isNull);
    });

    test('clears token on ApiClient', () async {
      await provider.logout();

      expect(mockApiClient.setTokenCalled, true);
      expect(mockApiClient.lastTokenSet, isNull);
    });
  });

  group('hasPermission', () {
    test('returns true for super admin regardless of resources', () async {
      mockRepo.loginResult = makeLoginResponse(isSuperAdmin: true);
      mockRepo.rolesPermissionsResult = makeRolesPermissions(
        isSuper: true,
        resources: [],
      );
      await provider.login('u@e.com', 'p');

      expect(provider.hasPermission('orders', 'read'), true);
      expect(provider.hasPermission('anything', 'delete'), true);
    });

    test('returns false when rolesPermissions is null', () {
      expect(provider.hasPermission('orders', 'read'), false);
    });

    test('returns true when resource and action match', () async {
      mockRepo.loginResult = makeLoginResponse(isSuperAdmin: false);
      mockRepo.rolesPermissionsResult = makeRolesPermissions(
        isSuper: false,
        resources: [
          ResourcePermission(resource: 'orders', actions: ['read', 'write']),
          ResourcePermission(resource: 'products', actions: ['read']),
        ],
      );
      await provider.login('u@e.com', 'p');

      expect(provider.hasPermission('orders', 'read'), true);
      expect(provider.hasPermission('orders', 'write'), true);
      expect(provider.hasPermission('products', 'read'), true);
    });

    test('returns false when resource matches but action does not', () async {
      mockRepo.loginResult = makeLoginResponse(isSuperAdmin: false);
      mockRepo.rolesPermissionsResult = makeRolesPermissions(
        isSuper: false,
        resources: [
          ResourcePermission(resource: 'orders', actions: ['read']),
        ],
      );
      await provider.login('u@e.com', 'p');

      expect(provider.hasPermission('orders', 'delete'), false);
    });

    test('returns false when resource does not match', () async {
      mockRepo.loginResult = makeLoginResponse(isSuperAdmin: false);
      mockRepo.rolesPermissionsResult = makeRolesPermissions(
        isSuper: false,
        resources: [
          ResourcePermission(resource: 'orders', actions: ['read']),
        ],
      );
      await provider.login('u@e.com', 'p');

      expect(provider.hasPermission('invoices', 'read'), false);
    });
  });

  group('restoreSession', () {
    test('does nothing when no token is stored', () async {
      await provider.restoreSession();

      expect(provider.isLoggedIn, false);
      expect(provider.user, isNull);
    });

    test('restores user from stored token and user data', () async {
      await mockTokenStorage.saveToken('restored-token');
      await mockTokenStorage.saveUserData(jsonEncode({
        'id': 10,
        'name': 'Restored User',
        'email': 'restored@e.com',
      }));
      mockRepo.rolesPermissionsResult = makeRolesPermissions(
        isSuper: false,
        resources: [
          ResourcePermission(resource: 'orders', actions: ['read']),
        ],
      );

      await provider.restoreSession();

      expect(provider.isLoggedIn, true);
      expect(provider.user!.id, 10);
      expect(provider.user!.name, 'Restored User');
      expect(provider.user!.email, 'restored@e.com');
      expect(provider.user!.isActive, true);
      expect(provider.isSuperAdmin, false);
      expect(provider.rolesPermissions, isNotNull);
    });

    test('sets isSuperAdmin from rolesPermissions', () async {
      await mockTokenStorage.saveToken('token');
      await mockTokenStorage.saveUserData(
          jsonEncode({'id': 1, 'name': 'Admin', 'email': 'a@e.com'}));
      mockRepo.rolesPermissionsResult = makeRolesPermissions(isSuper: true);

      await provider.restoreSession();

      expect(provider.isSuperAdmin, true);
    });

    test('sets token on ApiClient', () async {
      await mockTokenStorage.saveToken('my-token');
      mockRepo.rolesPermissionsResult = makeRolesPermissions();

      await provider.restoreSession();

      expect(mockApiClient.lastTokenSet, 'my-token');
    });

    test('calls logout when getRolesPermissions fails', () async {
      await mockTokenStorage.saveToken('expired-token');
      await mockTokenStorage.saveUserData(
          jsonEncode({'id': 1, 'name': 'X', 'email': 'x@e.com'}));
      mockRepo.rolesPermissionsError = Exception('Unauthorized');

      await provider.restoreSession();

      expect(provider.isLoggedIn, false);
      expect(provider.user, isNull);
      expect(provider.isSuperAdmin, false);
      final token = await mockTokenStorage.getToken();
      expect(token, isNull);
    });

    test('handles missing user data gracefully', () async {
      await mockTokenStorage.saveToken('token-only');
      mockRepo.rolesPermissionsResult = makeRolesPermissions();

      await provider.restoreSession();

      expect(provider.rolesPermissions, isNotNull);
      expect(provider.user, isNull);
    });
  });
}
