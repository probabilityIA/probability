import 'dart:convert';
import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../core/storage/token_storage.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/login_repository.dart';

class LoginProvider extends ChangeNotifier {
  final TokenStorage _tokenStorage;
  final ApiClient _apiClient;

  bool _isLoading = false;
  String? _error;
  UserInfo? _user;
  bool _isSuperAdmin = false;
  List<BusinessInfo> _businesses = [];
  UserRolesPermissionsResponse? _rolesPermissions;

  LoginProvider({
    required TokenStorage tokenStorage,
    required ApiClient apiClient,
  })  : _tokenStorage = tokenStorage,
        _apiClient = apiClient;

  bool get isLoading => _isLoading;
  String? get error => _error;
  UserInfo? get user => _user;
  bool get isSuperAdmin => _isSuperAdmin;
  List<BusinessInfo> get businesses => _businesses;
  UserRolesPermissionsResponse? get rolesPermissions => _rolesPermissions;
  bool get isLoggedIn => _user != null;

  LoginUseCases get _useCases =>
      LoginUseCases(LoginApiRepository(_apiClient));

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
