import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/core/network/api_client.dart';
import 'package:mobile_central/services/modules/inventory/app/use_cases.dart';
import 'package:mobile_central/services/modules/inventory/domain/entities.dart';
import 'package:mobile_central/services/modules/inventory/domain/ports.dart';
import 'package:mobile_central/services/modules/inventory/ui/providers/inventory_provider.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// ---------------------------------------------------------------------------
// Manual mock for IInventoryRepository
// ---------------------------------------------------------------------------
class MockInventoryRepository implements IInventoryRepository {
  final List<String> calls = [];

  List<InventoryLevel>? getProductInventoryResult;
  PaginatedResponse<InventoryLevel>? getWarehouseInventoryResult;
  StockMovement? adjustStockResult;
  Map<String, dynamic>? transferStockResult;
  PaginatedResponse<StockMovement>? getMovementsResult;
  PaginatedResponse<MovementType>? getMovementTypesResult;

  Exception? errorToThrow;

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<InventoryLevel>> getProductInventory(String productId,
      {int? businessId}) async {
    _trackCall('getProductInventory');
    return getProductInventoryResult!;
  }

  @override
  Future<PaginatedResponse<InventoryLevel>> getWarehouseInventory(
      int warehouseId, GetInventoryParams? params) async {
    _trackCall('getWarehouseInventory');
    return getWarehouseInventoryResult!;
  }

  @override
  Future<StockMovement> adjustStock(AdjustStockDTO data,
      {int? businessId}) async {
    _trackCall('adjustStock');
    return adjustStockResult!;
  }

  @override
  Future<Map<String, dynamic>> transferStock(TransferStockDTO data,
      {int? businessId}) async {
    _trackCall('transferStock');
    return transferStockResult!;
  }

  @override
  Future<PaginatedResponse<StockMovement>> getMovements(
      GetMovementsParams? params) async {
    _trackCall('getMovements');
    return getMovementsResult!;
  }

  @override
  Future<PaginatedResponse<MovementType>> getMovementTypes(
      GetMovementTypesParams? params) async {
    _trackCall('getMovementTypes');
    return getMovementTypesResult!;
  }
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------
Pagination _makePagination({
  int currentPage = 1,
  int total = 50,
  bool hasNext = true,
}) {
  return Pagination(
    currentPage: currentPage,
    perPage: 20,
    total: total,
    lastPage: 3,
    hasNext: hasNext,
    hasPrev: currentPage > 1,
  );
}

InventoryLevel _makeInventoryLevel({int id = 1, String name = 'Widget'}) {
  return InventoryLevel(
    id: id,
    productId: 'prod-$id',
    warehouseId: 1,
    businessId: 1,
    quantity: 100,
    reservedQty: 10,
    availableQty: 90,
    productName: name,
    createdAt: '',
    updatedAt: '',
  );
}

StockMovement _makeMovement({int id = 1}) {
  return StockMovement(
    id: id,
    productId: 'prod-$id',
    warehouseId: 1,
    businessId: 1,
    movementTypeId: 1,
    movementTypeCode: 'ADJ',
    movementTypeName: 'Adjustment',
    reason: 'test',
    quantity: 10,
    previousQty: 90,
    newQty: 100,
    notes: '',
    createdAt: '',
  );
}

MovementType _makeMovementType({int id = 1, String name = 'Adjustment'}) {
  return MovementType(
    id: id,
    code: 'ADJ',
    name: name,
    description: 'desc',
    isActive: true,
    direction: 'in',
    createdAt: '',
    updatedAt: '',
  );
}

InventoryProvider _createProvider(MockInventoryRepository mockRepo) {
  final apiClient = ApiClient();
  final useCases = InventoryUseCases(mockRepo);
  return InventoryProvider(apiClient: apiClient, useCases: useCases);
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockInventoryRepository mockRepo;
  late InventoryProvider provider;

  setUp(() {
    mockRepo = MockInventoryRepository();
    provider = _createProvider(mockRepo);
  });

  group('Initial state', () {
    test('has empty inventory levels list', () {
      expect(provider.inventoryLevels, isEmpty);
    });

    test('has empty movements list', () {
      expect(provider.movements, isEmpty);
    });

    test('has empty movement types list', () {
      expect(provider.movementTypes, isEmpty);
    });

    test('has null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('has null movements pagination', () {
      expect(provider.movementsPagination, isNull);
    });

    test('is not loading', () {
      expect(provider.isLoading, false);
    });

    test('has no error', () {
      expect(provider.error, isNull);
    });

    test('has default page 1', () {
      expect(provider.page, 1);
    });

    test('has default pageSize 20', () {
      expect(provider.pageSize, 20);
    });
  });

  group('fetchWarehouseInventory', () {
    test('updates inventory levels and pagination on success', () async {
      mockRepo.getWarehouseInventoryResult = PaginatedResponse<InventoryLevel>(
        data: [
          _makeInventoryLevel(id: 1, name: 'A'),
          _makeInventoryLevel(id: 2, name: 'B'),
        ],
        pagination: _makePagination(total: 2),
      );

      await provider.fetchWarehouseInventory(1);

      expect(provider.inventoryLevels.length, 2);
      expect(provider.inventoryLevels[0].productName, 'A');
      expect(provider.inventoryLevels[1].productName, 'B');
      expect(provider.pagination, isNotNull);
      expect(provider.pagination!.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets isLoading during fetch and clears after', () async {
      final loadingStates = <bool>[];
      provider.addListener(() {
        loadingStates.add(provider.isLoading);
      });

      mockRepo.getWarehouseInventoryResult = PaginatedResponse<InventoryLevel>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchWarehouseInventory(1);

      expect(loadingStates, [true, false]);
    });

    test('clears previous error before fetching', () async {
      mockRepo.errorToThrow = Exception('first error');
      await provider.fetchWarehouseInventory(1);
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getWarehouseInventoryResult = PaginatedResponse<InventoryLevel>(
        data: [],
        pagination: _makePagination(),
      );

      await provider.fetchWarehouseInventory(1);

      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('fetch failed');

      await provider.fetchWarehouseInventory(1);

      expect(provider.error, contains('fetch failed'));
      expect(provider.isLoading, false);
    });
  });

  group('getProductInventory', () {
    test('returns inventory levels on success', () async {
      mockRepo.getProductInventoryResult = [_makeInventoryLevel(id: 5)];

      final result = await provider.getProductInventory('prod-5');

      expect(result.length, 1);
      expect(result.first.id, 5);
    });

    test('returns empty list and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('product error');

      final result = await provider.getProductInventory('prod-1');

      expect(result, isEmpty);
      expect(provider.error, contains('product error'));
    });
  });

  group('adjustStock', () {
    test('returns movement on success', () async {
      mockRepo.adjustStockResult = _makeMovement(id: 42);

      final dto = AdjustStockDTO(
        productId: 'prod-1',
        warehouseId: 1,
        quantity: 10,
        reason: 'recount',
      );
      final result = await provider.adjustStock(dto);

      expect(result, isNotNull);
      expect(result!.id, 42);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('adjust failed');

      final dto = AdjustStockDTO(
        productId: 'prod-1',
        warehouseId: 1,
        quantity: 10,
        reason: 'r',
      );
      final result = await provider.adjustStock(dto);

      expect(result, isNull);
      expect(provider.error, contains('adjust failed'));
    });
  });

  group('transferStock', () {
    test('returns true on success', () async {
      mockRepo.transferStockResult = {'success': true};

      final dto = TransferStockDTO(
        productId: 'prod-1',
        fromWarehouseId: 1,
        toWarehouseId: 2,
        quantity: 50,
      );
      final result = await provider.transferStock(dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('transfer failed');

      final dto = TransferStockDTO(
        productId: 'prod-1',
        fromWarehouseId: 1,
        toWarehouseId: 2,
        quantity: 50,
      );
      final result = await provider.transferStock(dto);

      expect(result, false);
      expect(provider.error, contains('transfer failed'));
    });
  });

  group('fetchMovements', () {
    test('updates movements and pagination on success', () async {
      mockRepo.getMovementsResult = PaginatedResponse<StockMovement>(
        data: [_makeMovement(id: 1), _makeMovement(id: 2)],
        pagination: _makePagination(total: 2),
      );

      await provider.fetchMovements();

      expect(provider.movements.length, 2);
      expect(provider.movementsPagination, isNotNull);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('movements failed');

      await provider.fetchMovements();

      expect(provider.error, contains('movements failed'));
      expect(provider.isLoading, false);
    });
  });

  group('fetchMovementTypes', () {
    test('updates movement types on success', () async {
      mockRepo.getMovementTypesResult = PaginatedResponse<MovementType>(
        data: [_makeMovementType(id: 1, name: 'Adjustment')],
        pagination: _makePagination(),
      );

      await provider.fetchMovementTypes();

      expect(provider.movementTypes.length, 1);
      expect(provider.movementTypes[0].name, 'Adjustment');
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('types failed');

      await provider.fetchMovementTypes();

      expect(provider.error, contains('types failed'));
    });
  });

  group('setPage', () {
    test('updates page', () {
      provider.setPage(3);
      expect(provider.page, 3);
    });
  });

  group('setFilters', () {
    test('resets page to 1', () {
      provider.setPage(5);
      provider.setFilters(search: 'test');
      expect(provider.page, 1);
    });
  });

  group('resetFilters', () {
    test('resets page to 1', () {
      provider.setPage(5);
      provider.resetFilters();
      expect(provider.page, 1);
    });
  });
}
