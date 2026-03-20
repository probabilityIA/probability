import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/auth/actions/app/use_cases.dart';
import 'package:mobile_central/services/auth/actions/domain/entities.dart';
import 'package:mobile_central/services/auth/actions/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockActionRepository implements IActionRepository {
  final List<String> calls = [];

  PaginatedResponse<ActionEntity>? getActionsResult;
  ActionEntity? getActionByIdResult;
  ActionEntity? createActionResult;
  ActionEntity? updateActionResult;

  Exception? errorToThrow;

  GetActionsParams? capturedGetParams;
  int? capturedId;
  CreateActionDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateActionDTO? capturedUpdateData;
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<ActionEntity>> getActions(
      GetActionsParams? params) async {
    calls.add('getActions');
    capturedGetParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getActionsResult!;
  }

  @override
  Future<ActionEntity> getActionById(int id) async {
    calls.add('getActionById');
    capturedId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getActionByIdResult!;
  }

  @override
  Future<ActionEntity> createAction(CreateActionDTO data) async {
    calls.add('createAction');
    capturedCreateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createActionResult!;
  }

  @override
  Future<ActionEntity> updateAction(int id, UpdateActionDTO data) async {
    calls.add('updateAction');
    capturedUpdateId = id;
    capturedUpdateData = data;
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

// --- Helpers ---

ActionEntity _makeAction({int id = 1, String name = 'read'}) {
  return ActionEntity(id: id, name: name);
}

Pagination _makePagination() {
  return Pagination(
    currentPage: 1,
    perPage: 10,
    total: 1,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

// --- Tests ---

void main() {
  late MockActionRepository mockRepo;
  late ActionUseCases useCases;

  setUp(() {
    mockRepo = MockActionRepository();
    useCases = ActionUseCases(mockRepo);
  });

  group('getActions', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<ActionEntity>(
        data: [_makeAction()],
        pagination: _makePagination(),
      );
      mockRepo.getActionsResult = expected;
      final params = GetActionsParams(page: 1, pageSize: 10);

      final result = await useCases.getActions(params);

      expect(result.data.length, 1);
      expect(result.data[0].name, 'read');
      expect(mockRepo.calls, ['getActions']);
      expect(mockRepo.capturedGetParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getActionsResult = PaginatedResponse<ActionEntity>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getActions(null);

      expect(mockRepo.capturedGetParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getActions(null), throwsException);
    });
  });

  group('getActionById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getActionByIdResult = _makeAction(id: 42, name: 'write');

      final result = await useCases.getActionById(42);

      expect(result.id, 42);
      expect(result.name, 'write');
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getActionById']);
    });
  });

  group('createAction', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateActionDTO(name: 'delete');
      mockRepo.createActionResult = _makeAction(id: 99, name: 'delete');

      final result = await useCases.createAction(dto);

      expect(result.id, 99);
      expect(result.name, 'delete');
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createAction']);
    });
  });

  group('updateAction', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateActionDTO(name: 'updated');
      mockRepo.updateActionResult = _makeAction(id: 5, name: 'updated');

      final result = await useCases.updateAction(5, dto);

      expect(result.id, 5);
      expect(result.name, 'updated');
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateAction']);
    });
  });

  group('deleteAction', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteAction(7);

      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deleteAction']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.deleteAction(7), throwsException);
    });
  });
}
