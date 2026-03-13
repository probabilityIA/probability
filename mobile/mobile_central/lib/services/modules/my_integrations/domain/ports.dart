import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IMyIntegrationsRepository {
  Future<PaginatedResponse<MyIntegration>> getIntegrations(GetMyIntegrationsParams? params);
  Future<MyIntegration> getIntegrationById(int id, {int? businessId});
}
