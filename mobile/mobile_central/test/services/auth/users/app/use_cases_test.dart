import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/users/app/use_cases.dart';
import 'package:mobile_central/services/auth/users/domain/entities.dart';
import 'package:mobile_central/services/auth/users/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IUserRepository
// ---------------------------------------------------------------------------
class MockUserRepository implements IUserRepository {
  // Captured arguments
  GetUsersParams? lastGetUsersParams;
  int? lastGetUserByIdArg;
  CreateUserDTO? lastCreateUserDTO;
  int? lastUpdateUserId;
  UpdateUserDTO? lastUpdateUserDTO;
  int? lastDeleteUserId;
  int? lastAssignRolesUserId;
  AssignRolesDTO? lastAssignRolesDTO;

  // Configurable return values / errors
  PaginatedResponse<User>? getUsersResult;
  User? getUserByIdResult;
  User? createUserResult;
  User? updateUserResult;
  Exception? errorToThrow;

  @override
  Future<PaginatedResponse<User>> getUsers(GetUsersParams? params) async {
    lastGetUsersParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getUsersResult ??
        PaginatedResponse(
          data: [],
          pagination: _emptyPagination(),
        );
  }

  @override
  Future<User> getUserById(int id) async {
    lastGetUserByIdArg = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getUserByIdResult ?? _defaultUser();
  }

  @override
  Future<User> createUser(CreateUserDTO data) async {
    lastCreateUserDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createUserResult ?? _defaultUser();
  }

  @override
  Future<User> updateUser(int id, UpdateUserDTO data) async {
    lastUpdateUserId = id;
    lastUpdateUserDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateUserResult ?? _defaultUser();
  }

  @override
  Future<void> deleteUser(int id) async {
    lastDeleteUserId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<void> assignRoles(int userId, AssignRolesDTO data) async {
    lastAssignRolesUserId = userId;
    lastAssignRolesDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
  }

  // Helpers
  static User _defaultUser() => User(
        id: 1,
        name: 'Test',
        email: 'test@test.com',
        isActive: true,
        isSuperUser: false,
        businessRoleAssignments: [],
      );

  static Pagination _emptyPagination() => Pagination(
        currentPage: 1,
        perPage: 10,
        total: 0,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
}

void main() {
  late MockUserRepository mockRepo;
  late UserUseCases useCases;

  setUp(() {
    mockRepo = MockUserRepository();
    useCases = UserUseCases(mockRepo);
  });

  group('getUsers', () {
    test('delegates to repository with params', () async {
      final params = GetUsersParams(page: 2, pageSize: 15, name: 'Alice');

      await useCases.getUsers(params);

      expect(mockRepo.lastGetUsersParams, same(params));
    });

    test('delegates to repository with null params', () async {
      await useCases.getUsers(null);

      expect(mockRepo.lastGetUsersParams, isNull);
    });

    test('returns the response from repository', () async {
      final expectedUsers = [
        User(
          id: 1,
          name: 'A',
          email: 'a@a.com',
          isActive: true,
          isSuperUser: false,
          businessRoleAssignments: [],
        ),
      ];
      final expectedPagination = Pagination(
        currentPage: 1,
        perPage: 10,
        total: 1,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
      mockRepo.getUsersResult = PaginatedResponse(
        data: expectedUsers,
        pagination: expectedPagination,
      );

      final result = await useCases.getUsers(null);

      expect(result.data, hasLength(1));
      expect(result.data[0].name, 'A');
      expect(result.pagination.total, 1);
    });
  });

  group('getUserById', () {
    test('delegates to repository with correct id', () async {
      await useCases.getUserById(42);

      expect(mockRepo.lastGetUserByIdArg, 42);
    });

    test('returns user from repository', () async {
      mockRepo.getUserByIdResult = User(
        id: 42,
        name: 'Found User',
        email: 'found@test.com',
        isActive: true,
        isSuperUser: false,
        businessRoleAssignments: [],
      );

      final user = await useCases.getUserById(42);

      expect(user.id, 42);
      expect(user.name, 'Found User');
    });
  });

  group('createUser', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateUserDTO(name: 'New', email: 'new@test.com');

      await useCases.createUser(dto);

      expect(mockRepo.lastCreateUserDTO, same(dto));
    });

    test('returns created user from repository', () async {
      mockRepo.createUserResult = User(
        id: 99,
        name: 'Created',
        email: 'created@test.com',
        isActive: true,
        isSuperUser: false,
        businessRoleAssignments: [],
      );

      final user = await useCases.createUser(
        CreateUserDTO(name: 'Created', email: 'created@test.com'),
      );

      expect(user.id, 99);
      expect(user.name, 'Created');
    });
  });

  group('updateUser', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = UpdateUserDTO(name: 'Updated');

      await useCases.updateUser(10, dto);

      expect(mockRepo.lastUpdateUserId, 10);
      expect(mockRepo.lastUpdateUserDTO, same(dto));
    });

    test('returns updated user from repository', () async {
      mockRepo.updateUserResult = User(
        id: 10,
        name: 'Updated',
        email: 'u@test.com',
        isActive: true,
        isSuperUser: false,
        businessRoleAssignments: [],
      );

      final user = await useCases.updateUser(
        10,
        UpdateUserDTO(name: 'Updated'),
      );

      expect(user.name, 'Updated');
    });
  });

  group('deleteUser', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteUser(77);

      expect(mockRepo.lastDeleteUserId, 77);
    });
  });

  group('assignRoles', () {
    test('delegates to repository with correct userId and DTO', () async {
      final dto = AssignRolesDTO(
        assignments: [RoleAssignment(businessId: 1, roleId: 2)],
      );

      await useCases.assignRoles(55, dto);

      expect(mockRepo.lastAssignRolesUserId, 55);
      expect(mockRepo.lastAssignRolesDTO, same(dto));
    });
  });

  group('error propagation', () {
    test('getUsers propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.getUsers(null),
        throwsA(isA<Exception>()),
      );
    });

    test('getUserById propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('not found');

      expect(
        () => useCases.getUserById(1),
        throwsA(isA<Exception>()),
      );
    });

    test('createUser propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('validation error');

      expect(
        () => useCases.createUser(
          CreateUserDTO(name: 'X', email: 'x@x.com'),
        ),
        throwsA(isA<Exception>()),
      );
    });

    test('updateUser propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('server error');

      expect(
        () => useCases.updateUser(1, UpdateUserDTO(name: 'Y')),
        throwsA(isA<Exception>()),
      );
    });

    test('deleteUser propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('forbidden');

      expect(
        () => useCases.deleteUser(1),
        throwsA(isA<Exception>()),
      );
    });

    test('assignRoles propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('bad request');

      expect(
        () => useCases.assignRoles(
          1,
          AssignRolesDTO(assignments: []),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });
}
