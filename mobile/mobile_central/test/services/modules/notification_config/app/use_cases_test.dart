import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/notification_config/app/use_cases.dart';
import 'package:mobile_central/services/modules/notification_config/domain/entities.dart';
import 'package:mobile_central/services/modules/notification_config/domain/ports.dart';

// ---------------------------------------------------------------------------
// Manual mock for INotificationConfigRepository
// ---------------------------------------------------------------------------
class MockNotificationConfigRepository
    implements INotificationConfigRepository {
  final List<String> calls = [];

  NotificationConfig? createResult;
  NotificationConfig? getByIdResult;
  NotificationConfig? updateResult;
  List<NotificationConfig>? listResult;
  SyncConfigsResponse? syncResult;

  Exception? errorToThrow;

  // Captured args
  int? lastId;
  int? lastBusinessId;
  CreateConfigDTO? lastCreateDTO;
  UpdateConfigDTO? lastUpdateDTO;
  ConfigFilter? lastFilter;
  SyncConfigsDTO? lastSyncDTO;

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<NotificationConfig> create(CreateConfigDTO dto,
      {int? businessId}) async {
    lastCreateDTO = dto;
    lastBusinessId = businessId;
    _trackCall('create');
    return createResult!;
  }

  @override
  Future<NotificationConfig> getById(int id, {int? businessId}) async {
    lastId = id;
    lastBusinessId = businessId;
    _trackCall('getById');
    return getByIdResult!;
  }

  @override
  Future<NotificationConfig> update(int id, UpdateConfigDTO dto,
      {int? businessId}) async {
    lastId = id;
    lastUpdateDTO = dto;
    lastBusinessId = businessId;
    _trackCall('update');
    return updateResult!;
  }

  @override
  Future<void> delete(int id, {int? businessId}) async {
    lastId = id;
    lastBusinessId = businessId;
    _trackCall('delete');
  }

  @override
  Future<List<NotificationConfig>> list({ConfigFilter? filter}) async {
    lastFilter = filter;
    _trackCall('list');
    return listResult!;
  }

  @override
  Future<SyncConfigsResponse> syncByIntegration(SyncConfigsDTO dto,
      {int? businessId}) async {
    lastSyncDTO = dto;
    lastBusinessId = businessId;
    _trackCall('syncByIntegration');
    return syncResult!;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
NotificationConfig _makeConfig({int id = 1}) {
  return NotificationConfig(
    id: id,
    businessId: 1,
    integrationId: 1,
    notificationTypeId: 1,
    notificationEventTypeId: 1,
    enabled: true,
    createdAt: '',
    updatedAt: '',
  );
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockNotificationConfigRepository mockRepo;
  late NotificationConfigUseCases useCases;

  setUp(() {
    mockRepo = MockNotificationConfigRepository();
    useCases = NotificationConfigUseCases(mockRepo);
  });

  group('create', () {
    test('delegates to repository and returns result', () async {
      mockRepo.createResult = _makeConfig(id: 99);

      final dto = CreateConfigDTO(
        businessId: 1,
        integrationId: 2,
        notificationTypeId: 3,
        notificationEventTypeId: 4,
      );
      final result = await useCases.create(dto, businessId: 5);

      expect(mockRepo.calls, ['create']);
      expect(mockRepo.lastCreateDTO, dto);
      expect(mockRepo.lastBusinessId, 5);
      expect(result.id, 99);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('create error');

      expect(
        () => useCases.create(
          CreateConfigDTO(
            businessId: 1,
            integrationId: 1,
            notificationTypeId: 1,
            notificationEventTypeId: 1,
          ),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getByIdResult = _makeConfig(id: 42);

      final result = await useCases.getById(42, businessId: 5);

      expect(mockRepo.calls, ['getById']);
      expect(mockRepo.lastId, 42);
      expect(mockRepo.lastBusinessId, 5);
      expect(result.id, 42);
    });

    test('delegates without businessId', () async {
      mockRepo.getByIdResult = _makeConfig(id: 10);

      await useCases.getById(10);

      expect(mockRepo.lastId, 10);
      expect(mockRepo.lastBusinessId, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('not found');

      expect(
        () => useCases.getById(999),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('update', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = UpdateConfigDTO(enabled: false);
      mockRepo.updateResult = _makeConfig(id: 5);

      final result = await useCases.update(5, dto, businessId: 3);

      expect(mockRepo.calls, ['update']);
      expect(mockRepo.lastId, 5);
      expect(mockRepo.lastUpdateDTO, dto);
      expect(mockRepo.lastBusinessId, 3);
      expect(result.id, 5);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('update error');

      expect(
        () => useCases.update(1, UpdateConfigDTO()),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('delete', () {
    test('delegates to repository with correct id', () async {
      await useCases.delete(8, businessId: 3);

      expect(mockRepo.calls, ['delete']);
      expect(mockRepo.lastId, 8);
      expect(mockRepo.lastBusinessId, 3);
    });

    test('delegates without businessId', () async {
      await useCases.delete(10);

      expect(mockRepo.lastId, 10);
      expect(mockRepo.lastBusinessId, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('delete error');

      expect(
        () => useCases.delete(1),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('list', () {
    test('delegates to repository and returns result', () async {
      mockRepo.listResult = [_makeConfig(id: 1), _makeConfig(id: 2)];

      final filter = ConfigFilter(businessId: 5);
      final result = await useCases.list(filter: filter);

      expect(mockRepo.calls, ['list']);
      expect(mockRepo.lastFilter, filter);
      expect(result.length, 2);
    });

    test('delegates with null filter', () async {
      mockRepo.listResult = [];

      await useCases.list();

      expect(mockRepo.calls, ['list']);
      expect(mockRepo.lastFilter, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('list error');

      expect(
        () => useCases.list(),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('syncByIntegration', () {
    test('delegates to repository and returns result', () async {
      mockRepo.syncResult = SyncConfigsResponse(
        created: 2,
        updated: 1,
        deleted: 0,
        configs: [_makeConfig()],
      );

      final dto = SyncConfigsDTO(
        integrationId: 5,
        rules: [
          SyncRule(
            notificationTypeId: 1,
            notificationEventTypeId: 2,
            enabled: true,
            description: 'Rule',
            orderStatusIds: [1],
          ),
        ],
      );
      final result = await useCases.syncByIntegration(dto, businessId: 3);

      expect(mockRepo.calls, ['syncByIntegration']);
      expect(mockRepo.lastSyncDTO, dto);
      expect(mockRepo.lastBusinessId, 3);
      expect(result.created, 2);
      expect(result.updated, 1);
      expect(result.configs.length, 1);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('sync error');

      expect(
        () => useCases.syncByIntegration(
          SyncConfigsDTO(integrationId: 1, rules: []),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });
}
