import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/pagination/paged_list_controller.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/user_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class UserProvider extends ChangeNotifier {
  UserProvider({required ApiClient apiClient, UserUseCases? useCases})
      : _apiClient = apiClient,
        _injectedUseCases = useCases {
    list = PagedListController<User>(fetcher: _fetchPage);
    list.addListener(notifyListeners);
  }

  final ApiClient _apiClient;
  final UserUseCases? _injectedUseCases;
  late final PagedListController<User> list;

  int? _businessId;
  String? _error;
  String _nameFilter = '';
  String _emailFilter = '';

  List<User> get users => list.loadedItems;
  bool get isLoading => list.isLoading;
  String? get error => _error ?? list.error;

  UserUseCases get _useCases =>
      _injectedUseCases ?? UserUseCases(UserApiRepository(_apiClient));

  Future<PaginatedResponse<User>> _fetchPage(int page, int pageSize) {
    return _useCases.getUsers(GetUsersParams(
      page: page,
      pageSize: pageSize,
      name: _nameFilter.isNotEmpty ? _nameFilter : null,
      email: _emailFilter.isNotEmpty ? _emailFilter : null,
      businessId: _businessId,
    ));
  }

  Future<void> fetchUsers({int? businessId}) {
    _businessId = businessId;
    _error = null;
    return list.refresh();
  }

  void setFilters({String? name, String? email}) {
    _nameFilter = name ?? _nameFilter;
    _emailFilter = email ?? _emailFilter;
  }

  void resetFilters() {
    _nameFilter = '';
    _emailFilter = '';
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

  @override
  void dispose() {
    list.removeListener(notifyListeners);
    list.dispose();
    super.dispose();
  }
}
