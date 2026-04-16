import '../../../../shared/types/paginated_response.dart';
import '../domain/entities.dart';
import '../domain/ports.dart';

class ActionUseCases {
  final IActionRepository _repository;

  ActionUseCases(this._repository);

  Future<PaginatedResponse<ActionEntity>> getActions(
      GetActionsParams? params) {
    return _repository.getActions(params);
  }

  Future<ActionEntity> getActionById(int id) {
    return _repository.getActionById(id);
  }

  Future<ActionEntity> createAction(CreateActionDTO data) {
    return _repository.createAction(data);
  }

  Future<ActionEntity> updateAction(int id, UpdateActionDTO data) {
    return _repository.updateAction(id, data);
  }

  Future<void> deleteAction(int id) {
    return _repository.deleteAction(id);
  }
}
