import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IUserRepository {
  Future<PaginatedResponse<User>> getUsers(GetUsersParams? params);
  Future<User> getUserById(int id);
  Future<User> createUser(CreateUserDTO data);
  Future<User> updateUser(int id, UpdateUserDTO data);
  Future<void> deleteUser(int id);
  Future<void> assignRoles(int userId, AssignRolesDTO data);
}
