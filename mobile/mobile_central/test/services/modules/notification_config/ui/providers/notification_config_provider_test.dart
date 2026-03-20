import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/core/network/api_client.dart';
import 'package:mobile_central/services/modules/notification_config/app/use_cases.dart';
import 'package:mobile_central/services/modules/notification_config/domain/entities.dart';
import 'package:mobile_central/services/modules/notification_config/domain/ports.dart';
import 'package:mobile_central/services/modules/notification_config/ui/providers/notification_config_provider.dart';

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

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<NotificationConfig> create(CreateConfigDTO dto,
      {int? businessId}) async {
    _trackCall('create');
    return createResult!;
  }

  @override
  Future<NotificationConfig> getById(int id, {int? businessId}) async {
    _trackCall('getById');
    return getByIdResult!;
  }

  @override
  Future<NotificationConfig> update(int id, UpdateConfigDTO dto,
      {int? businessId}) async {
    _trackCall('update');
    return updateResult!;
  }

  @override
  Future<void> delete(int id, {int? businessId}) async {
    _trackCall('delete');
  }

  @override
  Future<List<NotificationConfig>> list({ConfigFilter? filter}) async {
    _trackCall('list');
    return listResult!;
  }

  @override
  Future<SyncConfigsResponse> syncByIntegration(SyncConfigsDTO dto,
      {int? businessId}) async {
    _trackCall('syncByIntegration');
    return syncResult!;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
NotificationConfig _makeConfig({int id = 1, bool enabled = true}) {
  return NotificationConfig(
    id: id,
    businessId: 1,
    integrationId: 1,
    notificationTypeId: 1,
    notificationEventTypeId: 1,
    enabled: enabled,
    createdAt: '',
    updatedAt: '',
  );
}

NotificationConfigProvider _createProvider(
    MockNotificationConfigRepository mockRepo) {
  final apiClient = ApiClient();
  final useCases = NotificationConfigUseCases(mockRepo);
  return NotificationConfigProvider(apiClient: apiClient, useCases: useCases);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockNotificationConfigRepository mockRepo;
  late NotificationConfigProvider provider;

  setUp(() {
    mockRepo = MockNotificationConfigRepository();
    provider = _createProvider(mockRepo);
  });

  group('Initial state', () {
    test('has empty configs list', () {
      expect(provider.configs, isEmpty);
    });

    test('is not loading', () {
      expect(provider.isLoading, false);
    });

    test('has no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchConfigs', () {
    test('updates configs on success', () async {
      mockRepo.listResult = [
        _makeConfig(id: 1),
        _makeConfig(id: 2),
      ];

      await provider.fetchConfigs();

      expect(provider.configs.length, 2);
      expect(provider.configs[0].id, 1);
      expect(provider.configs[1].id, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets isLoading during fetch and clears after', () async {
      final loadingStates = <bool>[];
      provider.addListener(() {
        loadingStates.add(provider.isLoading);
      });

      mockRepo.listResult = [];

      await provider.fetchConfigs();

      expect(loadingStates, [true, false]);
    });

    test('clears previous error before fetching', () async {
      mockRepo.errorToThrow = Exception('first error');
      await provider.fetchConfigs();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.listResult = [];

      await provider.fetchConfigs();

      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('fetch failed');

      await provider.fetchConfigs();

      expect(provider.error, contains('fetch failed'));
      expect(provider.isLoading, false);
    });
  });

  group('getById', () {
    test('returns config on success', () async {
      mockRepo.getByIdResult = _makeConfig(id: 42);

      final result = await provider.getById(42);

      expect(result, isNotNull);
      expect(result!.id, 42);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('not found');

      final result = await provider.getById(999);

      expect(result, isNull);
      expect(provider.error, contains('not found'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('error');

      await provider.getById(1);

      expect(notified, true);
    });
  });

  group('createConfig', () {
    test('returns config on success', () async {
      mockRepo.createResult = _makeConfig(id: 99);

      final dto = CreateConfigDTO(
        businessId: 1,
        integrationId: 2,
        notificationTypeId: 3,
        notificationEventTypeId: 4,
      );
      final result = await provider.createConfig(dto);

      expect(result, isNotNull);
      expect(result!.id, 99);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('create failed');

      final dto = CreateConfigDTO(
        businessId: 1,
        integrationId: 2,
        notificationTypeId: 3,
        notificationEventTypeId: 4,
      );
      final result = await provider.createConfig(dto);

      expect(result, isNull);
      expect(provider.error, contains('create failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('create error');

      await provider.createConfig(
        CreateConfigDTO(
          businessId: 1,
          integrationId: 1,
          notificationTypeId: 1,
          notificationEventTypeId: 1,
        ),
      );

      expect(notified, true);
    });
  });

  group('updateConfig', () {
    test('returns config on success', () async {
      mockRepo.updateResult = _makeConfig(id: 5, enabled: false);

      final dto = UpdateConfigDTO(enabled: false);
      final result = await provider.updateConfig(5, dto);

      expect(result, isNotNull);
      expect(result!.id, 5);
      expect(result.enabled, false);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('update failed');

      final result = await provider.updateConfig(1, UpdateConfigDTO());

      expect(result, isNull);
      expect(provider.error, contains('update failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('update error');

      await provider.updateConfig(1, UpdateConfigDTO());

      expect(notified, true);
    });
  });

  group('deleteConfig', () {
    test('returns true on success', () async {
      final result = await provider.deleteConfig(5);

      expect(result, true);
      expect(mockRepo.calls, ['delete']);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('delete failed');

      final result = await provider.deleteConfig(5);

      expect(result, false);
      expect(provider.error, contains('delete failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('delete error');

      await provider.deleteConfig(1);

      expect(notified, true);
    });
  });

  group('syncByIntegration', () {
    test('returns response on success', () async {
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
      final result = await provider.syncByIntegration(dto);

      expect(result, isNotNull);
      expect(result!.created, 2);
      expect(result.updated, 1);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('sync failed');

      final result = await provider.syncByIntegration(
        SyncConfigsDTO(integrationId: 1, rules: []),
      );

      expect(result, isNull);
      expect(provider.error, contains('sync failed'));
    });

    test('notifies listeners on failure', () async {
      var notified = false;
      provider.addListener(() => notified = true);

      mockRepo.errorToThrow = Exception('sync error');

      await provider.syncByIntegration(
        SyncConfigsDTO(integrationId: 1, rules: []),
      );

      expect(notified, true);
    });
  });
}
