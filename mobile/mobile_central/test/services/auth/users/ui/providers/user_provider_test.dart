import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/users/app/use_cases.dart';
import 'package:mobile_central/services/auth/users/domain/entities.dart';
import 'package:mobile_central/services/auth/users/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// We cannot directly inject a mock repository into UserProvider because it
// creates its own use cases internally via ApiClient. To test the provider's
// state management logic in isolation, we build a TestableUserProvider that
// accepts an IUserRepository and exercises the same logic.
// ---------------------------------------------------------------------------

// ---------------------------------------------------------------------------
// Manual mock for IUserRepository
// ---------------------------------------------------------------------------
class MockUserRepository implements IUserRepository {
  PaginatedResponse<User>? getUsersResult;
  User? getUserByIdResult;
  User? createUserResult;
  User? updateUserResult;
  Exception? errorToThrow;

  GetUsersParams? lastGetUsersParams;
  int? lastDeleteId;
  int? lastAssignRolesUserId;
  AssignRolesDTO? lastAssignRolesDTO;

  @override
  Future<PaginatedResponse<User>> getUsers(GetUsersParams? params) async {
    lastGetUsersParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getUsersResult ??
        PaginatedResponse(
          data: [],
          pagination: _defaultPagination(),
        );
  }

  @override
  Future<User> getUserById(int id) async {
    if (errorToThrow != null) throw errorToThrow!;
    return getUserByIdResult ?? _defaultUser();
  }

  @override
  Future<User> createUser(CreateUserDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return createUserResult ?? _defaultUser();
  }

  @override
  Future<User> updateUser(int id, UpdateUserDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return updateUserResult ?? _defaultUser();
  }

  @override
  Future<void> deleteUser(int id) async {
    lastDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> assignRoles(int userId, AssignRolesDTO data) async {
    lastAssignRolesUserId = userId;
    lastAssignRolesDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
  }

  static User _defaultUser() => User(
        id: 1,
        name: 'Test User',
        email: 'test@test.com',
        isActive: true,
        isSuperUser: false,
        businessRoleAssignments: [],
      );

  static Pagination _defaultPagination() => Pagination(
        currentPage: 1,
        perPage: 20,
        total: 0,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
}

// ---------------------------------------------------------------------------
// Testable provider that mirrors UserProvider logic but accepts a repository
// ---------------------------------------------------------------------------
class TestableUserProvider {
  final IUserRepository _repository;

  List<User> _users = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _nameFilter = '';
  String _emailFilter = '';

  int _notifyCount = 0;

  TestableUserProvider(this._repository);

  // Getters matching UserProvider
  List<User> get users => _users;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;
  int get notifyCount => _notifyCount;
  String get nameFilter => _nameFilter;
  String get emailFilter => _emailFilter;

  UserUseCases get _useCases => UserUseCases(_repository);

  void _notifyListeners() {
    _notifyCount++;
  }

  Future<void> fetchUsers({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

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
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
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
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateUser(int id, UpdateUserDTO data) async {
    try {
      await _useCases.updateUser(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteUser(int id) async {
    try {
      await _useCases.deleteUser(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> assignRoles(int userId, AssignRolesDTO data) async {
    try {
      await _useCases.assignRoles(userId, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockUserRepository mockRepo;
  late TestableUserProvider provider;

  setUp(() {
    mockRepo = MockUserRepository();
    provider = TestableUserProvider(mockRepo);
  });

  group('initial state', () {
    test('has empty users list', () {
      expect(provider.users, isEmpty);
    });

    test('pagination is null', () {
      expect(provider.pagination, isNull);
    });

    test('isLoading is false', () {
      expect(provider.isLoading, false);
    });

    test('error is null', () {
      expect(provider.error, isNull);
    });

    test('page defaults to 1', () {
      expect(provider.page, 1);
    });

    test('pageSize defaults to 20', () {
      expect(provider.pageSize, 20);
    });
  });

  group('fetchUsers', () {
    test('updates users list and pagination on success', () async {
      final testUsers = [
        User(
          id: 1,
          name: 'Alice',
          email: 'alice@test.com',
          isActive: true,
          isSuperUser: false,
          businessRoleAssignments: [],
        ),
        User(
          id: 2,
          name: 'Bob',
          email: 'bob@test.com',
          isActive: true,
          isSuperUser: false,
          businessRoleAssignments: [],
        ),
      ];
      final testPagination = Pagination(
        currentPage: 1,
        perPage: 20,
        total: 2,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
      mockRepo.getUsersResult = PaginatedResponse(
        data: testUsers,
        pagination: testPagination,
      );

      await provider.fetchUsers();

      expect(provider.users, hasLength(2));
      expect(provider.users[0].name, 'Alice');
      expect(provider.users[1].name, 'Bob');
      expect(provider.pagination, isNotNull);
      expect(provider.pagination!.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('server down');

      await provider.fetchUsers();

      expect(provider.error, contains('server down'));
      expect(provider.users, isEmpty);
      expect(provider.isLoading, false);
    });

    test('notifies listeners twice (loading start and end)', () async {
      await provider.fetchUsers();

      expect(provider.notifyCount, 2);
    });

    test('passes correct params including page and pageSize', () async {
      await provider.fetchUsers();

      expect(mockRepo.lastGetUsersParams, isNotNull);
      expect(mockRepo.lastGetUsersParams!.page, 1);
      expect(mockRepo.lastGetUsersParams!.pageSize, 20);
    });

    test('passes name and email filters when set', () async {
      provider.setFilters(name: 'John', email: 'john@');

      await provider.fetchUsers();

      expect(mockRepo.lastGetUsersParams!.name, 'John');
      expect(mockRepo.lastGetUsersParams!.email, 'john@');
    });

    test('does not pass empty filters', () async {
      // filters are empty by default
      await provider.fetchUsers();

      expect(mockRepo.lastGetUsersParams!.name, isNull);
      expect(mockRepo.lastGetUsersParams!.email, isNull);
    });

    test('passes businessId when provided', () async {
      await provider.fetchUsers(businessId: 7);

      expect(mockRepo.lastGetUsersParams!.businessId, 7);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('fail');
      await provider.fetchUsers();
      expect(provider.error, isNotNull);

      // Now succeed
      mockRepo.errorToThrow = null;
      await provider.fetchUsers();

      expect(provider.error, isNull);
    });
  });

  group('setPage', () {
    test('updates internal page value', () {
      provider.setPage(3);

      expect(provider.page, 3);
    });

    test('does not notify listeners', () {
      provider.setPage(5);

      expect(provider.notifyCount, 0);
    });

    test('fetchUsers uses the updated page', () async {
      provider.setPage(4);

      await provider.fetchUsers();

      expect(mockRepo.lastGetUsersParams!.page, 4);
    });
  });

  group('setFilters', () {
    test('updates name filter', () {
      provider.setFilters(name: 'Carlos');

      expect(provider.nameFilter, 'Carlos');
    });

    test('updates email filter', () {
      provider.setFilters(email: 'test@');

      expect(provider.emailFilter, 'test@');
    });

    test('resets page to 1', () {
      provider.setPage(5);
      provider.setFilters(name: 'X');

      expect(provider.page, 1);
    });

    test('preserves existing filters when not overridden', () {
      provider.setFilters(name: 'A');
      provider.setFilters(email: 'B');

      expect(provider.nameFilter, 'A');
      expect(provider.emailFilter, 'B');
    });
  });

  group('resetFilters', () {
    test('clears name and email filters', () {
      provider.setFilters(name: 'X', email: 'Y');
      provider.resetFilters();

      expect(provider.nameFilter, '');
      expect(provider.emailFilter, '');
    });

    test('resets page to 1', () {
      provider.setPage(10);
      provider.resetFilters();

      expect(provider.page, 1);
    });
  });

  group('createUser', () {
    test('returns User on success', () async {
      final createdUser = User(
        id: 50,
        name: 'New User',
        email: 'new@test.com',
        isActive: true,
        isSuperUser: false,
        businessRoleAssignments: [],
      );
      mockRepo.createUserResult = createdUser;

      final result = await provider.createUser(
        CreateUserDTO(name: 'New User', email: 'new@test.com'),
      );

      expect(result, isNotNull);
      expect(result!.id, 50);
      expect(result.name, 'New User');
    });

    test('returns null on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('create failed');

      final result = await provider.createUser(
        CreateUserDTO(name: 'Fail', email: 'fail@test.com'),
      );

      expect(result, isNull);
      expect(provider.error, contains('create failed'));
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');

      await provider.createUser(
        CreateUserDTO(name: 'X', email: 'x@x.com'),
      );

      expect(provider.notifyCount, 1);
    });

    test('does not notify listeners on success', () async {
      await provider.createUser(
        CreateUserDTO(name: 'X', email: 'x@x.com'),
      );

      expect(provider.notifyCount, 0);
    });
  });

  group('updateUser', () {
    test('returns true on success', () async {
      final result = await provider.updateUser(
        1,
        UpdateUserDTO(name: 'Updated'),
      );

      expect(result, true);
    });

    test('returns false on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('update failed');

      final result = await provider.updateUser(
        1,
        UpdateUserDTO(name: 'Fail'),
      );

      expect(result, false);
      expect(provider.error, contains('update failed'));
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');

      await provider.updateUser(1, UpdateUserDTO(name: 'X'));

      expect(provider.notifyCount, 1);
    });

    test('does not notify listeners on success', () async {
      await provider.updateUser(1, UpdateUserDTO(name: 'X'));

      expect(provider.notifyCount, 0);
    });
  });

  group('deleteUser', () {
    test('returns true on success', () async {
      final result = await provider.deleteUser(1);

      expect(result, true);
      expect(mockRepo.lastDeleteId, 1);
    });

    test('returns false on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('delete failed');

      final result = await provider.deleteUser(1);

      expect(result, false);
      expect(provider.error, contains('delete failed'));
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');

      await provider.deleteUser(1);

      expect(provider.notifyCount, 1);
    });

    test('does not notify listeners on success', () async {
      await provider.deleteUser(1);

      expect(provider.notifyCount, 0);
    });
  });

  group('assignRoles', () {
    test('returns true on success', () async {
      final dto = AssignRolesDTO(
        assignments: [RoleAssignment(businessId: 1, roleId: 2)],
      );

      final result = await provider.assignRoles(10, dto);

      expect(result, true);
      expect(mockRepo.lastAssignRolesUserId, 10);
      expect(mockRepo.lastAssignRolesDTO, same(dto));
    });

    test('returns false on failure and sets error', () async {
      mockRepo.errorToThrow = Exception('assign failed');

      final result = await provider.assignRoles(
        10,
        AssignRolesDTO(assignments: []),
      );

      expect(result, false);
      expect(provider.error, contains('assign failed'));
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');

      await provider.assignRoles(10, AssignRolesDTO(assignments: []));

      expect(provider.notifyCount, 1);
    });

    test('does not notify listeners on success', () async {
      await provider.assignRoles(10, AssignRolesDTO(assignments: []));

      expect(provider.notifyCount, 0);
    });
  });

  group('loading states', () {
    test('isLoading is false before fetch', () {
      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch completes', () async {
      await provider.fetchUsers();

      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch fails', () async {
      mockRepo.errorToThrow = Exception('fail');

      await provider.fetchUsers();

      expect(provider.isLoading, false);
    });
  });
}
