import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/my_integrations/app/use_cases.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/entities.dart';
import 'package:mobile_central/services/modules/my_integrations/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IMyIntegrationsRepository
// ---------------------------------------------------------------------------
class MockMyIntegrationsRepository implements IMyIntegrationsRepository {
  final List<String> calls = [];

  PaginatedResponse<MyIntegration>? getIntegrationsResult;
  MyIntegration? getIntegrationByIdResult;

  Exception? errorToThrow;

  // Captured args
  GetMyIntegrationsParams? lastParams;
  int? lastId;
  int? lastBusinessId;

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<PaginatedResponse<MyIntegration>> getIntegrations(
      GetMyIntegrationsParams? params) async {
    lastParams = params;
    _trackCall('getIntegrations');
    return getIntegrationsResult!;
  }

  @override
  Future<MyIntegration> getIntegrationById(int id,
      {int? businessId}) async {
    lastId = id;
    lastBusinessId = businessId;
    _trackCall('getIntegrationById');
    return getIntegrationByIdResult!;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
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

MyIntegration _makeIntegration({int id = 1, String name = 'Test'}) {
  return MyIntegration(
    id: id,
    createdAt: '',
    updatedAt: '',
    businessId: 1,
    integrationTypeId: 1,
    name: name,
    isActive: true,
  );
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockMyIntegrationsRepository mockRepo;
  late MyIntegrationsUseCases useCases;

  setUp(() {
    mockRepo = MockMyIntegrationsRepository();
    useCases = MyIntegrationsUseCases(mockRepo);
  });

  group('getIntegrations', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getIntegrationsResult = PaginatedResponse<MyIntegration>(
        data: [_makeIntegration(id: 1, name: 'Shopify')],
        pagination: _makePagination(),
      );

      final params = GetMyIntegrationsParams(page: 1, pageSize: 10);
      final result = await useCases.getIntegrations(params);

      expect(mockRepo.calls, ['getIntegrations']);
      expect(mockRepo.lastParams, params);
      expect(result.data.length, 1);
      expect(result.data.first.name, 'Shopify');
    });

    test('delegates with null params', () async {
      mockRepo.getIntegrationsResult = PaginatedResponse<MyIntegration>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getIntegrations(null);

      expect(mockRepo.calls, ['getIntegrations']);
      expect(mockRepo.lastParams, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.getIntegrations(null),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getIntegrationById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getIntegrationByIdResult =
          _makeIntegration(id: 42, name: 'Amazon');

      final result = await useCases.getIntegrationById(42, businessId: 5);

      expect(mockRepo.calls, ['getIntegrationById']);
      expect(mockRepo.lastId, 42);
      expect(mockRepo.lastBusinessId, 5);
      expect(result.id, 42);
      expect(result.name, 'Amazon');
    });

    test('delegates without businessId', () async {
      mockRepo.getIntegrationByIdResult = _makeIntegration(id: 10);

      await useCases.getIntegrationById(10);

      expect(mockRepo.calls, ['getIntegrationById']);
      expect(mockRepo.lastId, 10);
      expect(mockRepo.lastBusinessId, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('not found');

      expect(
        () => useCases.getIntegrationById(999),
        throwsA(isA<Exception>()),
      );
    });
  });
}
