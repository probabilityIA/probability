import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../domain/entities.dart';
import '../../domain/ports.dart';

class MyIntegrationsApiRepository implements IMyIntegrationsRepository {
  final ApiClient _client;

  MyIntegrationsApiRepository(this._client);

  @override
  Future<PaginatedResponse<MyIntegration>> getIntegrations(GetMyIntegrationsParams? params) async {
    final response = await _client.get(
      '/integrations',
      queryParameters: params?.toQueryParams(),
    );
    final data = response.data;
    final items = (data['data'] as List<dynamic>?)
            ?.map((e) => MyIntegration.fromJson(e))
            .toList() ??
        [];
    final pagination = Pagination.fromJson(data['pagination'] ?? {});
    return PaginatedResponse(data: items, pagination: pagination);
  }

  @override
  Future<MyIntegration> getIntegrationById(int id, {int? businessId}) async {
    final qp = businessId != null ? {'business_id': businessId} : null;
    final response = await _client.get('/integrations/$id', queryParameters: qp);
    return MyIntegration.fromJson(response.data['data'] ?? response.data);
  }
}
