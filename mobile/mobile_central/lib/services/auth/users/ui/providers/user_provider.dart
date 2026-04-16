import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/user_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class UserProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<User> _users = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _nameFilter = '';
  String _emailFilter = '';

  UserProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<User> get users => _users;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;

  UserUseCases get _useCases =>
      UserUseCases(UserApiRepository(_apiClient));

  Future<void> fetchUsers({int? businessId}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final params = GetUsersParams(
        page: _page,
        pageSize: _pageSize,
        name: _nameFilter.isNotEmpty ? _nameFilter : null,
        email: _emailFilter.isNotEmpty ? _emailFilter : null,
        businessId: businessId,
      );
      final response = await _useCases.getUsers(params);
      _users = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  void setPage(int page) {
    _page = page;
  }

  void setFilters({String? name, String? email}) {
    _nameFilter = name ?? _nameFilter;
    _emailFilter = email ?? _emailFilter;
    _page = 1;
  }

  void resetFilters() {
    _nameFilter = '';
    _emailFilter = '';
    _page = 1;
  }

  Future<User?> createUser(CreateUserDTO data) async {
    try {
      final user = await _useCases.createUser(data);
      return user;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateUser(int id, UpdateUserDTO data) async {
    try {
      await _useCases.updateUser(id, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteUser(int id) async {
    try {
      await _useCases.deleteUser(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> assignRoles(int userId, AssignRolesDTO data) async {
    try {
      await _useCases.assignRoles(userId, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }
}
