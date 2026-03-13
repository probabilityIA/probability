import '../../../../shared/types/paginated_response.dart';
import 'entities.dart';

abstract class IActionRepository {
  Future<PaginatedResponse<ActionEntity>> getActions(
      GetActionsParams? params);
  Future<ActionEntity> getActionById(int id);
  Future<ActionEntity> createAction(CreateActionDTO data);
  Future<ActionEntity> updateAction(int id, UpdateActionDTO data);
  Future<void> deleteAction(int id);
}
