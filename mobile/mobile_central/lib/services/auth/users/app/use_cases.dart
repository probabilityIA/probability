import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class UserUseCases {
  final IUserRepository _repository;

  UserUseCases(this._repository);

  Future<PaginatedResponse<User>> getUsers(GetUsersParams? params) {
    return _repository.getUsers(params);
  }

  Future<User> getUserById(int id) {
    return _repository.getUserById(id);
  }

  Future<User> createUser(CreateUserDTO data) {
    return _repository.createUser(data);
  }

  Future<User> updateUser(int id, UpdateUserDTO data) {
    return _repository.updateUser(id, data);
  }

  Future<void> deleteUser(int id) {
    return _repository.deleteUser(id);
  }

  Future<void> assignRoles(int userId, AssignRolesDTO data) {
    return _repository.assignRoles(userId, data);
  }
}
