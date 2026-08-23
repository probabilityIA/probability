import 'dart:convert';
import 'package:flutter/foundation.dart';
import '../../../../../core/errors/error_parser.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../core/storage/token_storage.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/login_repository.dart';

class LoginProvider extends ChangeNotifier {
  LoginProvider({
    required TokenStorage tokenStorage,
    required ApiClient apiClient,
  })  : _tokenStorage = tokenStorage,
        _apiClient = apiClient;

  final TokenStorage _tokenStorage;
  final ApiClient _apiClient;

  bool _isLoading = false;
  String? _error;
  UserInfo? _user;
  bool _isSuperAdmin = false;
  List<BusinessInfo> _businesses = [];
  UserRolesPermissionsResponse? _rolesPermissions;

  bool get isLoading => _isLoading;
  String? get error => _error;
  UserInfo? get user => _user;
  bool get isSuperAdmin => _isSuperAdmin;
  List<BusinessInfo> get businesses => _businesses;
  UserRolesPermissionsResponse? get rolesPermissions => _rolesPermissions;
  bool get isLoggedIn => _user != null;

  String? get businessName =>
      _rolesPermissions?.businessName ??
      (_businesses.isNotEmpty ? _businesses.first.name : null);

  String? get roleName => _rolesPermissions?.role;

  BusinessInfo? get activeBusiness =>
      _businesses.isNotEmpty ? _businesses.first : null;

  LoginUseCases get _useCases => LoginUseCases(LoginApiRepository(_apiClient));

  void clearError() {
    if (_error == null) return;
    _error = null;
    notifyListeners();
  }

  Future<bool> login(String email, String password) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.login(email, password);
      final data = response.data;

      await _tokenStorage.saveToken(data.token.trim());
      _apiClient.setToken(data.token.trim());

      _user = data.user;
      _isSuperAdmin = data.isSuperAdmin;
      _businesses = data.businesses;

      await _persistSession();
      await _fetchRolesPermissions();

      _isLoading = false;
      notifyListeners();
      return true;
    } catch (e) {
      _error = parseError(e);
      _isLoading = false;
      notifyListeners();
      return false;
    }
  }

  Future<void> _persistSession() async {
    await _tokenStorage.saveUserData(jsonEncode({
      'id': _user?.id,
      'name': _user?.name,
      'email': _user?.email,
      'phone': _user?.phone,
      'avatar_url': _user?.avatarUrl,
    }));
    await _tokenStorage.saveSessionData(jsonEncode({
      'is_super_admin': _isSuperAdmin,
      'businesses': _businesses
          .map((b) => {
                'id': b.id,
                'name': b.name,
                'logo_url': b.logoUrl,
                'primary_color': b.primaryColor,
                'secondary_color': b.secondaryColor,
                'accent_color': b.accentColor,
              })
          .toList(),
    }));
  }

  Future<void> _fetchRolesPermissions() async {
    try {
      _rolesPermissions = await _useCases.getRolesPermissions();
    } catch (_) {}
  }

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
      _error = parseError(e);
      _isLoading = false;
      notifyListeners();
      return null;
    }
  }

  Future<List<RecoveryChannel>> recoveryChannels(String email) =>
      _useCases.getRecoveryChannels(email);

  Future<SimpleAuthResponse> forgotPassword(String email, String channel) =>
      _useCases.forgotPassword(email, channel);

  Future<SimpleAuthResponse> verifyOtp(String email, String code) =>
      _useCases.verifyOtp(email, code);

  Future<SimpleAuthResponse> resetPassword(String token, String newPassword) =>
      _useCases.resetPassword(token, newPassword);

  Future<void> restoreSession() async {
    final token = await _tokenStorage.getToken();
    if (token == null) return;

    _apiClient.setToken(token);

    try {
      _rolesPermissions = await _useCases.getRolesPermissions();
      _isSuperAdmin = _rolesPermissions?.isSuper ?? false;

      final userData = await _tokenStorage.getUserData();
      if (userData != null) {
        final json = jsonDecode(userData) as Map<String, dynamic>;
        _user = UserInfo(
          id: json['id'] ?? 0,
          name: json['name'] ?? '',
          email: json['email'] ?? '',
          phone: json['phone'],
          avatarUrl: json['avatar_url'],
          isActive: true,
        );
      }

      final sessionData = await _tokenStorage.getSessionData();
      if (sessionData != null) {
        final json = jsonDecode(sessionData) as Map<String, dynamic>;
        _isSuperAdmin = json['is_super_admin'] ?? _isSuperAdmin;
        _businesses = (json['businesses'] as List<dynamic>?)
                ?.map((e) => BusinessInfo.fromJson(Map<String, dynamic>.from(e)))
                .toList() ??
            [];
      }

      notifyListeners();
    } catch (_) {
      await logout();
    }
  }

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

  bool hasPermission(String resource, String action) {
    if (_isSuperAdmin) return true;
    if (_rolesPermissions == null) return false;
    return _rolesPermissions!.resources.any(
      (r) => r.resource == resource && r.actions.contains(action),
    );
  }
}
