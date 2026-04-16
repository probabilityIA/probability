import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/orderstatus/app/use_cases.dart';
import 'package:mobile_central/services/modules/orderstatus/domain/entities.dart';
import 'package:mobile_central/services/modules/orderstatus/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IOrderStatusRepository
// ---------------------------------------------------------------------------
class MockOrderStatusRepository implements IOrderStatusRepository {
  // Captured arguments
  GetOrderStatusMappingsParams? lastGetMappingsParams;
  int? lastGetMappingByIdArg;
  CreateOrderStatusMappingDTO? lastCreateMappingDTO;
  int? lastUpdateMappingId;
  UpdateOrderStatusMappingDTO? lastUpdateMappingDTO;
  int? lastDeleteMappingId;
  int? lastToggleMappingActiveId;
  bool? lastGetStatusesIsActive;
  CreateOrderStatusDTO? lastCreateStatusDTO;
  int? lastGetStatusByIdArg;
  int? lastUpdateStatusId;
  CreateOrderStatusDTO? lastUpdateStatusDTO;
  int? lastDeleteStatusId;
  int? lastGetChannelStatusesIntegrationTypeId;
  bool? lastGetChannelStatusesIsActive;
  CreateChannelStatusDTO? lastCreateChannelStatusDTO;
  int? lastUpdateChannelStatusId;
  UpdateChannelStatusDTO? lastUpdateChannelStatusDTO;
  int? lastDeleteChannelStatusId;

  // Configurable return values / errors
  PaginatedResponse<OrderStatusMapping>? getMappingsResult;
  OrderStatusMapping? getMappingByIdResult;
  OrderStatusMapping? createMappingResult;
  OrderStatusMapping? updateMappingResult;
  OrderStatusMapping? toggleMappingActiveResult;
  List<OrderStatusInfo>? getStatusesResult;
  OrderStatusInfo? createStatusResult;
  OrderStatusInfo? getStatusByIdResult;
  OrderStatusInfo? updateStatusResult;
  List<IntegrationTypeInfo>? getEcommerceIntegrationTypesResult;
  List<ChannelStatusInfo>? getChannelStatusesResult;
  ChannelStatusInfo? createChannelStatusResult;
  ChannelStatusInfo? updateChannelStatusResult;
  Exception? errorToThrow;

  @override
  Future<PaginatedResponse<OrderStatusMapping>> getMappings(GetOrderStatusMappingsParams? params) async {
    lastGetMappingsParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getMappingsResult ?? PaginatedResponse(data: [], pagination: _emptyPagination());
  }

  @override
  Future<OrderStatusMapping> getMappingById(int id) async {
    lastGetMappingByIdArg = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getMappingByIdResult ?? _defaultMapping();
  }

  @override
  Future<OrderStatusMapping> createMapping(CreateOrderStatusMappingDTO data) async {
    lastCreateMappingDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createMappingResult ?? _defaultMapping();
  }

  @override
  Future<OrderStatusMapping> updateMapping(int id, UpdateOrderStatusMappingDTO data) async {
    lastUpdateMappingId = id;
    lastUpdateMappingDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateMappingResult ?? _defaultMapping();
  }

  @override
  Future<void> deleteMapping(int id) async {
    lastDeleteMappingId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<OrderStatusMapping> toggleMappingActive(int id) async {
    lastToggleMappingActiveId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return toggleMappingActiveResult ?? _defaultMapping();
  }

  @override
  Future<List<OrderStatusInfo>> getStatuses({bool? isActive}) async {
    lastGetStatusesIsActive = isActive;
    if (errorToThrow != null) throw errorToThrow!;
    return getStatusesResult ?? [];
  }

  @override
  Future<OrderStatusInfo> createStatus(CreateOrderStatusDTO data) async {
    lastCreateStatusDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createStatusResult ?? _defaultStatus();
  }

  @override
  Future<OrderStatusInfo> getStatusById(int id) async {
    lastGetStatusByIdArg = id;
    if (errorToThrow != null) throw errorToThrow!;
    return getStatusByIdResult ?? _defaultStatus();
  }

  @override
  Future<OrderStatusInfo> updateStatus(int id, CreateOrderStatusDTO data) async {
    lastUpdateStatusId = id;
    lastUpdateStatusDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateStatusResult ?? _defaultStatus();
  }

  @override
  Future<void> deleteStatus(int id) async {
    lastDeleteStatusId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<IntegrationTypeInfo>> getEcommerceIntegrationTypes() async {
    if (errorToThrow != null) throw errorToThrow!;
    return getEcommerceIntegrationTypesResult ?? [];
  }

  @override
  Future<List<ChannelStatusInfo>> getChannelStatuses(int integrationTypeId, {bool? isActive}) async {
    lastGetChannelStatusesIntegrationTypeId = integrationTypeId;
    lastGetChannelStatusesIsActive = isActive;
    if (errorToThrow != null) throw errorToThrow!;
    return getChannelStatusesResult ?? [];
  }

  @override
  Future<ChannelStatusInfo> createChannelStatus(CreateChannelStatusDTO data) async {
    lastCreateChannelStatusDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createChannelStatusResult ?? _defaultChannelStatus();
  }

  @override
  Future<ChannelStatusInfo> updateChannelStatus(int id, UpdateChannelStatusDTO data) async {
    lastUpdateChannelStatusId = id;
    lastUpdateChannelStatusDTO = data;
    if (errorToThrow != null) throw errorToThrow!;
    return updateChannelStatusResult ?? _defaultChannelStatus();
  }

  @override
  Future<void> deleteChannelStatus(int id) async {
    lastDeleteChannelStatusId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  // Helpers
  static OrderStatusMapping _defaultMapping() => OrderStatusMapping(
        id: 1,
        integrationTypeId: 5,
        originalStatus: 'paid',
        orderStatusId: 2,
        isActive: true,
        description: 'Test mapping',
        createdAt: '2026-01-01T00:00:00Z',
        updatedAt: '2026-01-01T00:00:00Z',
      );

  static OrderStatusInfo _defaultStatus() => OrderStatusInfo(
        id: 1,
        code: 'pending',
        name: 'Pending',
      );

  static ChannelStatusInfo _defaultChannelStatus() => ChannelStatusInfo(
        id: 1,
        integrationTypeId: 5,
        code: 'fulfilled',
        name: 'Fulfilled',
        isActive: true,
        displayOrder: 1,
      );

  static Pagination _emptyPagination() => Pagination(
        currentPage: 1,
        perPage: 10,
        total: 0,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
}

void main() {
  late MockOrderStatusRepository mockRepo;
  late OrderStatusUseCases useCases;

  setUp(() {
    mockRepo = MockOrderStatusRepository();
    useCases = OrderStatusUseCases(mockRepo);
  });

  group('getMappings', () {
    test('delegates to repository with params', () async {
      final params = GetOrderStatusMappingsParams(page: 2, pageSize: 15);

      await useCases.getMappings(params);

      expect(mockRepo.lastGetMappingsParams, same(params));
    });

    test('delegates to repository with null params', () async {
      await useCases.getMappings(null);

      expect(mockRepo.lastGetMappingsParams, isNull);
    });

    test('returns the response from repository', () async {
      final expectedMappings = [MockOrderStatusRepository._defaultMapping()];
      mockRepo.getMappingsResult = PaginatedResponse(
        data: expectedMappings,
        pagination: Pagination(currentPage: 1, perPage: 10, total: 1, lastPage: 1, hasNext: false, hasPrev: false),
      );

      final result = await useCases.getMappings(null);

      expect(result.data, hasLength(1));
      expect(result.data[0].id, 1);
    });
  });

  group('getMappingById', () {
    test('delegates to repository with correct id', () async {
      await useCases.getMappingById(42);

      expect(mockRepo.lastGetMappingByIdArg, 42);
    });
  });

  group('createMapping', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateOrderStatusMappingDTO(integrationTypeId: 5, originalStatus: 'paid', orderStatusId: 2);

      await useCases.createMapping(dto);

      expect(mockRepo.lastCreateMappingDTO, same(dto));
    });
  });

  group('updateMapping', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = UpdateOrderStatusMappingDTO(originalStatus: 'shipped', orderStatusId: 3);

      await useCases.updateMapping(10, dto);

      expect(mockRepo.lastUpdateMappingId, 10);
      expect(mockRepo.lastUpdateMappingDTO, same(dto));
    });
  });

  group('deleteMapping', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteMapping(77);

      expect(mockRepo.lastDeleteMappingId, 77);
    });
  });

  group('toggleMappingActive', () {
    test('delegates to repository with correct id', () async {
      await useCases.toggleMappingActive(5);

      expect(mockRepo.lastToggleMappingActiveId, 5);
    });
  });

  group('getStatuses', () {
    test('delegates to repository with isActive param', () async {
      await useCases.getStatuses(isActive: true);

      expect(mockRepo.lastGetStatusesIsActive, true);
    });

    test('delegates to repository without isActive param', () async {
      await useCases.getStatuses();

      expect(mockRepo.lastGetStatusesIsActive, isNull);
    });

    test('returns list from repository', () async {
      mockRepo.getStatusesResult = [
        OrderStatusInfo(id: 1, code: 'a', name: 'A'),
        OrderStatusInfo(id: 2, code: 'b', name: 'B'),
      ];

      final result = await useCases.getStatuses();

      expect(result, hasLength(2));
    });
  });

  group('createStatus', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateOrderStatusDTO(code: 'new', name: 'New');

      await useCases.createStatus(dto);

      expect(mockRepo.lastCreateStatusDTO, same(dto));
    });
  });

  group('getStatusById', () {
    test('delegates to repository with correct id', () async {
      await useCases.getStatusById(3);

      expect(mockRepo.lastGetStatusByIdArg, 3);
    });
  });

  group('updateStatus', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = CreateOrderStatusDTO(code: 'updated', name: 'Updated');

      await useCases.updateStatus(5, dto);

      expect(mockRepo.lastUpdateStatusId, 5);
      expect(mockRepo.lastUpdateStatusDTO, same(dto));
    });
  });

  group('deleteStatus', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteStatus(8);

      expect(mockRepo.lastDeleteStatusId, 8);
    });
  });

  group('getEcommerceIntegrationTypes', () {
    test('returns list from repository', () async {
      mockRepo.getEcommerceIntegrationTypesResult = [
        IntegrationTypeInfo(id: 1, code: 'shopify', name: 'Shopify'),
      ];

      final result = await useCases.getEcommerceIntegrationTypes();

      expect(result, hasLength(1));
      expect(result[0].code, 'shopify');
    });
  });

  group('getChannelStatuses', () {
    test('delegates to repository with correct integrationTypeId', () async {
      await useCases.getChannelStatuses(5);

      expect(mockRepo.lastGetChannelStatusesIntegrationTypeId, 5);
    });

    test('delegates to repository with isActive param', () async {
      await useCases.getChannelStatuses(5, isActive: true);

      expect(mockRepo.lastGetChannelStatusesIsActive, true);
    });
  });

  group('createChannelStatus', () {
    test('delegates to repository with correct DTO', () async {
      final dto = CreateChannelStatusDTO(integrationTypeId: 5, code: 'new', name: 'New', isActive: true, displayOrder: 1);

      await useCases.createChannelStatus(dto);

      expect(mockRepo.lastCreateChannelStatusDTO, same(dto));
    });
  });

  group('updateChannelStatus', () {
    test('delegates to repository with correct id and DTO', () async {
      final dto = UpdateChannelStatusDTO(code: 'updated', name: 'Updated', isActive: true, displayOrder: 2);

      await useCases.updateChannelStatus(10, dto);

      expect(mockRepo.lastUpdateChannelStatusId, 10);
      expect(mockRepo.lastUpdateChannelStatusDTO, same(dto));
    });
  });

  group('deleteChannelStatus', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteChannelStatus(15);

      expect(mockRepo.lastDeleteChannelStatusId, 15);
    });
  });

  group('error propagation', () {
    test('getMappings propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('network error');
      expect(() => useCases.getMappings(null), throwsA(isA<Exception>()));
    });

    test('getMappingById propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('not found');
      expect(() => useCases.getMappingById(1), throwsA(isA<Exception>()));
    });

    test('createMapping propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('validation error');
      expect(
        () => useCases.createMapping(CreateOrderStatusMappingDTO(integrationTypeId: 1, originalStatus: 'x', orderStatusId: 1)),
        throwsA(isA<Exception>()),
      );
    });

    test('updateMapping propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('server error');
      expect(
        () => useCases.updateMapping(1, UpdateOrderStatusMappingDTO(originalStatus: 'x', orderStatusId: 1)),
        throwsA(isA<Exception>()),
      );
    });

    test('deleteMapping propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('forbidden');
      expect(() => useCases.deleteMapping(1), throwsA(isA<Exception>()));
    });

    test('toggleMappingActive propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('error');
      expect(() => useCases.toggleMappingActive(1), throwsA(isA<Exception>()));
    });

    test('getStatuses propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('error');
      expect(() => useCases.getStatuses(), throwsA(isA<Exception>()));
    });

    test('createStatus propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('error');
      expect(() => useCases.createStatus(CreateOrderStatusDTO(code: 'x', name: 'X')), throwsA(isA<Exception>()));
    });

    test('deleteStatus propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('error');
      expect(() => useCases.deleteStatus(1), throwsA(isA<Exception>()));
    });

    test('getEcommerceIntegrationTypes propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('error');
      expect(() => useCases.getEcommerceIntegrationTypes(), throwsA(isA<Exception>()));
    });

    test('createChannelStatus propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('error');
      expect(
        () => useCases.createChannelStatus(CreateChannelStatusDTO(integrationTypeId: 1, code: 'x', name: 'X', isActive: true, displayOrder: 1)),
        throwsA(isA<Exception>()),
      );
    });

    test('deleteChannelStatus propagates repository exceptions', () {
      mockRepo.errorToThrow = Exception('error');
      expect(() => useCases.deleteChannelStatus(1), throwsA(isA<Exception>()));
    });
  });
}
