import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class ActionApiRepository implements IActionRepository {
  final ApiClient _client;

  ActionApiRepository(this._client);

  @override
  Future<PaginatedResponse<ActionEntity>> getActions(
      GetActionsParams? params) async {
    final response = await _client.get(
      '/actions',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final actions = (data['data'] as List<dynamic>?)
            ?.map((a) => ActionEntity.fromJson(a))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: actions, pagination: pagination);
  }

  @override
  Future<ActionEntity> getActionById(int id) async {
    final response = await _client.get('/actions/$id');
    return ActionEntity.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<ActionEntity> createAction(CreateActionDTO data) async {
    final response = await _client.post('/actions', data: data.toJson());
    return ActionEntity.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<ActionEntity> updateAction(int id, UpdateActionDTO data) async {
    final response = await _client.put('/actions/$id', data: data.toJson());
    return ActionEntity.fromJson(response.data['data'] ?? response.data);
  }

  @override
  Future<void> deleteAction(int id) async {
    await _client.delete('/actions/$id');
  }
}
