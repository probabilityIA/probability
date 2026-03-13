import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class MyIntegrationsUseCases {
  final IMyIntegrationsRepository _repository;

  MyIntegrationsUseCases(this._repository);

  Future<PaginatedResponse<MyIntegration>> getIntegrations(GetMyIntegrationsParams? params) {
    return _repository.getIntegrations(params);
  }

  Future<MyIntegration> getIntegrationById(int id, {int? businessId}) {
    return _repository.getIntegrationById(id, businessId: businessId);
  }
}
