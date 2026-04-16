import 'package:flutter/foundation.dart';
import '../../../../../core/network/api_client.dart';
import '../../../../../shared/types/paginated_response.dart';
import '../../app/use_cases.dart';
import '../../domain/entities.dart';
import '../../infra/repository/action_repository.dart';
import '../../../../../core/errors/error_parser.dart';

class ActionProvider extends ChangeNotifier {
  final ApiClient _apiClient;

  List<ActionEntity> _actions = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  ActionProvider({required ApiClient apiClient}) : _apiClient = apiClient;

  List<ActionEntity> get actions => _actions;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  ActionUseCases get _useCases =>
      ActionUseCases(ActionApiRepository(_apiClient));

  Future<void> fetchActions({GetActionsParams? params}) async {
    _isLoading = true;
    _error = null;
    notifyListeners();

    try {
      final response = await _useCases.getActions(params);
      _actions = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = parseError(e);
    }

    _isLoading = false;
    notifyListeners();
  }

  Future<ActionEntity?> createAction(CreateActionDTO data) async {
    try {
      return await _useCases.createAction(data);
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return null;
    }
  }

  Future<bool> updateAction(int id, UpdateActionDTO data) async {
    try {
      await _useCases.updateAction(id, data);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }

  Future<bool> deleteAction(int id) async {
    try {
      await _useCases.deleteAction(id);
      return true;
    } catch (e) {
      _error = parseError(e);
      notifyListeners();
      return false;
    }
  }
}
