import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/permission_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class PermissionProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<Permission> _permissions = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  final int _page = 1;
  final int _pageSize = 20;

  PermissionProvider({required ApiClient apiClient})
      : _apiClient = apiClient;

  List<Permission> get permissions => _permissions;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  PermissionUseCases get _useCases =>
      PermissionUseCases(PermissionApiRepository(_apiClient));

  Future<void> fetchPermissions({GetPermissionsParams? params}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final queryParams = params ??
          GetPermissionsParams(page: _page, pageSize: _pageSize);
      final response = await _useCases.getPermissions(queryParams);
      _permissions = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Permission?> createPermission(CreatePermissionDTO data) async {
    try {
      return await _useCases.createPermission(data);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updatePermission(int id, UpdatePermissionDTO data) async {
    try {
      await _useCases.updatePermission(id, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deletePermission(int id) async {
    try {
      await _useCases.deletePermission(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }
}
