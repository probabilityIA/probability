import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/role_repository.dart';

class RoleProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<Role> _roles = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  final int _page = 1;
  final int _pageSize = 20;

  RoleProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<Role> get roles => _roles;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  RoleUseCases get _useCases =>
      RoleUseCases(RoleApiRepository(_apiClient));

  Future<void> fetchRoles({GetRolesParams? params}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final queryParams = params ??
          GetRolesParams(page: _page, pageSize: _pageSize);
      final response = await _useCases.getRoles(queryParams);
      _roles = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<Role?> createRole(CreateRoleDTO data) async {
    try {
      return await _useCases.createRole(data);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateRole(int id, UpdateRoleDTO data) async {
    try {
      await _useCases.updateRole(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteRole(int id) async {
    try {
      await _useCases.deleteRole(id);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }

  Future<RolePermissionsResponse?> getRolePermissions(int roleId) async {
    try {
      return await _useCases.getRolePermissions(roleId);
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return null;
    }
  }

  Future<bool> assignPermissions(
      int roleId, AssignPermissionsDTO data) async {
    try {
      await _useCases.assignPermissions(roleId, data);
      return true;
    } catch (e) {
      _error = e.toString();
      notifyListeners();
      return false;
    }
  }
}
