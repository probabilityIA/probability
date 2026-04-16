import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/orderstatus/app/use_cases.dart';
import 'package:mobile_central/services/modules/orderstatus/domain/entities.dart';
import 'package:mobile_central/services/modules/orderstatus/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IOrderStatusRepository
// ---------------------------------------------------------------------------
class MockOrderStatusRepository implements IOrderStatusRepository {
  PaginatedResponse<OrderStatusMapping>? getMappingsResult;
  List<OrderStatusInfo>? getStatusesResult;
  Exception? errorToThrow;

  GetOrderStatusMappingsParams? lastGetMappingsParams;
  bool? lastGetStatusesIsActive;

  @override
  Future<PaginatedResponse<OrderStatusMapping>> getMappings(GetOrderStatusMappingsParams? params) async {
    lastGetMappingsParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getMappingsResult ?? PaginatedResponse(data: [], pagination: _defaultPagination());
  }

  @override
  Future<OrderStatusMapping> getMappingById(int id) async {
    if (errorToThrow != null) throw errorToThrow!;
    return _defaultMapping();
  }

  @override
  Future<OrderStatusMapping> createMapping(CreateOrderStatusMappingDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return _defaultMapping();
  }

  @override
  Future<OrderStatusMapping> updateMapping(int id, UpdateOrderStatusMappingDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return _defaultMapping();
  }

  @override
  Future<void> deleteMapping(int id) async {
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<OrderStatusMapping> toggleMappingActive(int id) async {
    if (errorToThrow != null) throw errorToThrow!;
    return _defaultMapping();
  }

  @override
  Future<List<OrderStatusInfo>> getStatuses({bool? isActive}) async {
    lastGetStatusesIsActive = isActive;
    if (errorToThrow != null) throw errorToThrow!;
    return getStatusesResult ?? [];
  }

  @override
  Future<OrderStatusInfo> createStatus(CreateOrderStatusDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return OrderStatusInfo(id: 1, code: 'new', name: 'New');
  }

  @override
  Future<OrderStatusInfo> getStatusById(int id) async {
    if (errorToThrow != null) throw errorToThrow!;
    return OrderStatusInfo(id: id, code: 'test', name: 'Test');
  }

  @override
  Future<OrderStatusInfo> updateStatus(int id, CreateOrderStatusDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return OrderStatusInfo(id: id, code: data.code, name: data.name);
  }

  @override
  Future<void> deleteStatus(int id) async {
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<IntegrationTypeInfo>> getEcommerceIntegrationTypes() async {
    if (errorToThrow != null) throw errorToThrow!;
    return [];
  }

  @override
  Future<List<ChannelStatusInfo>> getChannelStatuses(int integrationTypeId, {bool? isActive}) async {
    if (errorToThrow != null) throw errorToThrow!;
    return [];
  }

  @override
  Future<ChannelStatusInfo> createChannelStatus(CreateChannelStatusDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return ChannelStatusInfo(id: 1, integrationTypeId: data.integrationTypeId, code: data.code, name: data.name, isActive: data.isActive, displayOrder: data.displayOrder);
  }

  @override
  Future<ChannelStatusInfo> updateChannelStatus(int id, UpdateChannelStatusDTO data) async {
    if (errorToThrow != null) throw errorToThrow!;
    return ChannelStatusInfo(id: id, integrationTypeId: 1, code: data.code, name: data.name, isActive: data.isActive, displayOrder: data.displayOrder);
  }

  @override
  Future<void> deleteChannelStatus(int id) async {
    if (errorToThrow != null) throw errorToThrow!;
  }

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

  static Pagination _defaultPagination() => Pagination(
        currentPage: 1,
        perPage: 20,
        total: 0,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
}

// ---------------------------------------------------------------------------
// Testable provider that mirrors OrderStatusProvider logic
// ---------------------------------------------------------------------------
class TestableOrderStatusProvider {
  final IOrderStatusRepository _repository;

  List<OrderStatusMapping> _mappings = [];
  List<OrderStatusInfo> _statuses = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;

  int _notifyCount = 0;

  TestableOrderStatusProvider(this._repository);

  List<OrderStatusMapping> get mappings => _mappings;
  List<OrderStatusInfo> get statuses => _statuses;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;
  int get notifyCount => _notifyCount;

  OrderStatusUseCases get _useCases => OrderStatusUseCases(_repository);

  void _notifyListeners() {
    _notifyCount++;
  }

  Future<void> fetchMappings({int? integrationTypeId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();
    try {
      final params = GetOrderStatusMappingsParams(page: _page, pageSize: _pageSize, integrationTypeId: integrationTypeId);
      final response = await _useCases.getMappings(params);
      _mappings = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }
    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchStatuses({bool? isActive}) async {
    try {
      _statuses = await _useCases.getStatuses(isActive: isActive);
      _notifyListeners();
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
    }
  }

  void setPage(int page) {
    _page = page;
  }
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockOrderStatusRepository mockRepo;
  late TestableOrderStatusProvider provider;

  setUp(() {
    mockRepo = MockOrderStatusRepository();
    provider = TestableOrderStatusProvider(mockRepo);
  });

  group('initial state', () {
    test('has empty mappings list', () {
      expect(provider.mappings, isEmpty);
    });

    test('has empty statuses list', () {
      expect(provider.statuses, isEmpty);
    });

    test('pagination is null', () {
      expect(provider.pagination, isNull);
    });

    test('isLoading is false', () {
      expect(provider.isLoading, false);
    });

    test('error is null', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchMappings', () {
    test('updates mappings list and pagination on success', () async {
      final testMappings = [MockOrderStatusRepository._defaultMapping()];
      final testPagination = Pagination(
        currentPage: 1,
        perPage: 20,
        total: 1,
        lastPage: 1,
        hasNext: false,
        hasPrev: false,
      );
      mockRepo.getMappingsResult = PaginatedResponse(
        data: testMappings,
        pagination: testPagination,
      );

      await provider.fetchMappings();

      expect(provider.mappings, hasLength(1));
      expect(provider.mappings[0].id, 1);
      expect(provider.pagination, isNotNull);
      expect(provider.pagination!.total, 1);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('server down');

      await provider.fetchMappings();

      expect(provider.error, contains('server down'));
      expect(provider.mappings, isEmpty);
      expect(provider.isLoading, false);
    });

    test('notifies listeners twice (loading start and end)', () async {
      await provider.fetchMappings();

      expect(provider.notifyCount, 2);
    });

    test('passes correct params including page and pageSize', () async {
      await provider.fetchMappings();

      expect(mockRepo.lastGetMappingsParams, isNotNull);
      expect(mockRepo.lastGetMappingsParams!.page, 1);
      expect(mockRepo.lastGetMappingsParams!.pageSize, 20);
    });

    test('passes integrationTypeId when provided', () async {
      await provider.fetchMappings(integrationTypeId: 5);

      expect(mockRepo.lastGetMappingsParams!.integrationTypeId, 5);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('fail');
      await provider.fetchMappings();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      await provider.fetchMappings();

      expect(provider.error, isNull);
    });
  });

  group('fetchStatuses', () {
    test('updates statuses list on success', () async {
      mockRepo.getStatusesResult = [
        OrderStatusInfo(id: 1, code: 'pending', name: 'Pending'),
        OrderStatusInfo(id: 2, code: 'confirmed', name: 'Confirmed'),
      ];

      await provider.fetchStatuses();

      expect(provider.statuses, hasLength(2));
      expect(provider.statuses[0].code, 'pending');
      expect(provider.statuses[1].code, 'confirmed');
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('fetch statuses failed');

      await provider.fetchStatuses();

      expect(provider.error, contains('fetch statuses failed'));
    });

    test('passes isActive when provided', () async {
      await provider.fetchStatuses(isActive: true);

      expect(mockRepo.lastGetStatusesIsActive, true);
    });

    test('does not pass isActive when not provided', () async {
      await provider.fetchStatuses();

      expect(mockRepo.lastGetStatusesIsActive, isNull);
    });

    test('notifies listeners on success', () async {
      await provider.fetchStatuses();

      expect(provider.notifyCount, 1);
    });

    test('notifies listeners on failure', () async {
      mockRepo.errorToThrow = Exception('err');

      await provider.fetchStatuses();

      expect(provider.notifyCount, 1);
    });
  });

  group('setPage', () {
    test('updates internal page value', () {
      provider.setPage(3);

      expect(provider.page, 3);
    });

    test('does not notify listeners', () {
      provider.setPage(5);

      expect(provider.notifyCount, 0);
    });

    test('fetchMappings uses the updated page', () async {
      provider.setPage(4);

      await provider.fetchMappings();

      expect(mockRepo.lastGetMappingsParams!.page, 4);
    });
  });

  group('loading states', () {
    test('isLoading is false before fetch', () {
      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch completes', () async {
      await provider.fetchMappings();

      expect(provider.isLoading, false);
    });

    test('isLoading is false after fetch fails', () async {
      mockRepo.errorToThrow = Exception('fail');

      await provider.fetchMappings();

      expect(provider.isLoading, false);
    });
  });
}
