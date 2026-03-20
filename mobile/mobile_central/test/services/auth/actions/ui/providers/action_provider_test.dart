import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/actions/app/use_cases.dart';
import 'package:mobile_central/services/auth/actions/domain/entities.dart';
import 'package:mobile_central/services/auth/actions/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockActionRepository implements IActionRepository {
  PaginatedResponse<ActionEntity>? getActionsResult;
  ActionEntity? createActionResult;
  ActionEntity? updateActionResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<ActionEntity>> getActions(
      GetActionsParams? params) async {
    calls.add('getActions');
    if (errorToThrow != null) throw errorToThrow!;
    return getActionsResult!;
  }

  @override
  Future<ActionEntity> getActionById(int id) async {
    calls.add('getActionById');
    if (errorToThrow != null) throw errorToThrow!;
    return ActionEntity(id: id, name: 'test');
  }

  @override
  Future<ActionEntity> createAction(CreateActionDTO data) async {
    calls.add('createAction');
    if (errorToThrow != null) throw errorToThrow!;
    return createActionResult!;
  }

  @override
  Future<ActionEntity> updateAction(int id, UpdateActionDTO data) async {
    calls.add('updateAction');
    if (errorToThrow != null) throw errorToThrow!;
    return updateActionResult!;
  }

  @override
  Future<void> deleteAction(int id) async {
    calls.add('deleteAction');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---

class TestableActionProvider {
  final ActionUseCases _useCases;

  List<ActionEntity> _actions = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableActionProvider(this._useCases);

  List<ActionEntity> get actions => _actions;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchActions({GetActionsParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.getActions(params);
      _actions = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<ActionEntity?> createAction(CreateActionDTO data) async {
    try {
      return await _useCases.createAction(data);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateAction(int id, UpdateActionDTO data) async {
    try {
      await _useCases.updateAction(id, data);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteAction(int id) async {
    try {
      await _useCases.deleteAction(id);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }
}

// --- Helpers ---

Pagination _makePagination({int total = 5}) {
  return Pagination(
    currentPage: 1,
    perPage: 20,
    total: total,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

ActionEntity _makeAction({int id = 1, String name = 'read'}) {
  return ActionEntity(id: id, name: name);
}

// --- Tests ---

void main() {
  late MockActionRepository mockRepo;
  late ActionUseCases useCases;
  late TestableActionProvider provider;

  setUp(() {
    mockRepo = MockActionRepository();
    useCases = ActionUseCases(mockRepo);
    provider = TestableActionProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty actions list', () {
      expect(provider.actions, isEmpty);
    });

    test('starts with null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchActions', () {
    test('notifies listeners twice (loading start and end)', () async {
      mockRepo.getActionsResult = PaginatedResponse<ActionEntity>(
        data: [_makeAction()],
        pagination: _makePagination(),
      );

      await provider.fetchActions();

      expect(provider.notifications.length, 2);
    });

    test('populates actions and pagination on success', () async {
      final pagination = _makePagination(total: 3);
      mockRepo.getActionsResult = PaginatedResponse<ActionEntity>(
        data: [
          _makeAction(id: 1, name: 'read'),
          _makeAction(id: 2, name: 'write'),
          _makeAction(id: 3, name: 'delete'),
        ],
        pagination: pagination,
      );

      await provider.fetchActions();

      expect(provider.actions.length, 3);
      expect(provider.actions[0].name, 'read');
      expect(provider.actions[1].name, 'write');
      expect(provider.actions[2].name, 'delete');
      expect(provider.pagination?.total, 3);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchActions();

      expect(provider.error, contains('Server error'));
      expect(provider.actions, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchActions();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getActionsResult = PaginatedResponse<ActionEntity>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchActions();

      expect(provider.error, isNull);
    });

    test('passes custom params to use cases', () async {
      mockRepo.getActionsResult = PaginatedResponse<ActionEntity>(
        data: [],
        pagination: _makePagination(),
      );
      final params = GetActionsParams(page: 2, name: 'write');

      await provider.fetchActions(params: params);

      expect(mockRepo.calls, contains('getActions'));
    });
  });

  group('createAction', () {
    test('returns created action on success', () async {
      final dto = CreateActionDTO(name: 'export');
      mockRepo.createActionResult = _makeAction(id: 10, name: 'export');

      final result = await provider.createAction(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.name, 'export');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateActionDTO(name: 'fail');

      final result = await provider.createAction(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateAction', () {
    test('returns true on success', () async {
      final dto = UpdateActionDTO(name: 'updated');
      mockRepo.updateActionResult = _makeAction(id: 5, name: 'updated');

      final result = await provider.updateAction(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateActionDTO(name: 'fail');

      final result = await provider.updateAction(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteAction', () {
    test('returns true on success', () async {
      final result = await provider.deleteAction(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteAction(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });
}
