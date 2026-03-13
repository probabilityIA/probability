import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class UserApiRepository implements IUserRepository {
  final ApiClient _client;

  UserApiRepository(this._client);

  @override
  Future<PaginatedResponse<User>> getUsers(GetUsersParams? params) async {
    final response = await _client.get(
      '/users',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final users = (data['data'] as List<dynamic>?)
            ?.map((u) => User.fromJson(u))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: users, pagination: pagination);
  }

  @override
  Future<User> getUserById(int id) async {
    final response = await _client.get('/users/$id');
    return User.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<User> createUser(CreateUserDTO data) async {
    final response = await _client.post('/users', data: data.toJson());
    return User.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<User> updateUser(int id, UpdateUserDTO data) async {
    final response = await _client.put('/users/$id', data: data.toJson());
    return User.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteUser(int id) async {
    await _client.delete('/users/$id');
  }

  @override
  Future<void> assignRoles(int userId, AssignRolesDTO data) async {
    await _client.post('/users/$userId/assign-role', data: data.toJson());
  }
}
