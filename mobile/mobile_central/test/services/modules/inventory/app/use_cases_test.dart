import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/inventory/app/use_cases.dart';
import 'package:mobile_central/services/modules/inventory/domain/entities.dart';
import 'package:mobile_central/services/modules/inventory/domain/ports.dart';
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

  // Captured args
  String? lastProductId;
  int? lastWarehouseId;
  int? lastBusinessId;
  GetInventoryParams? lastInventoryParams;
  AdjustStockDTO? lastAdjustDTO;
  TransferStockDTO? lastTransferDTO;
  GetMovementsParams? lastMovementsParams;
  GetMovementTypesParams? lastMovementTypesParams;

  void _trackCall(String name) {
    calls.add(name);
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<InventoryLevel>> getProductInventory(String productId,
      {int? businessId}) async {
    lastProductId = productId;
    lastBusinessId = businessId;
    _trackCall('getProductInventory');
    return getProductInventoryResult!;
  }

  @override
  Future<PaginatedResponse<InventoryLevel>> getWarehouseInventory(
      int warehouseId, GetInventoryParams? params) async {
    lastWarehouseId = warehouseId;
    lastInventoryParams = params;
    _trackCall('getWarehouseInventory');
    return getWarehouseInventoryResult!;
  }

  @override
  Future<StockMovement> adjustStock(AdjustStockDTO data,
      {int? businessId}) async {
    lastAdjustDTO = data;
    lastBusinessId = businessId;
    _trackCall('adjustStock');
    return adjustStockResult!;
  }

  @override
  Future<Map<String, dynamic>> transferStock(TransferStockDTO data,
      {int? businessId}) async {
    lastTransferDTO = data;
    lastBusinessId = businessId;
    _trackCall('transferStock');
    return transferStockResult!;
  }

  @override
  Future<PaginatedResponse<StockMovement>> getMovements(
      GetMovementsParams? params) async {
    lastMovementsParams = params;
    _trackCall('getMovements');
    return getMovementsResult!;
  }

  @override
  Future<PaginatedResponse<MovementType>> getMovementTypes(
      GetMovementTypesParams? params) async {
    lastMovementTypesParams = params;
    _trackCall('getMovementTypes');
    return getMovementTypesResult!;
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

InventoryLevel _makeInventoryLevel({int id = 1}) {
  return InventoryLevel(
    id: id,
    productId: 'prod-$id',
    warehouseId: 1,
    businessId: 1,
    quantity: 100,
    reservedQty: 10,
    availableQty: 90,
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

MovementType _makeMovementType({int id = 1}) {
  return MovementType(
    id: id,
    code: 'ADJ',
    name: 'Adjustment',
    description: 'Stock adjustment',
    isActive: true,
    direction: 'in',
    createdAt: '',
    updatedAt: '',
  );
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------
void main() {
  late MockInventoryRepository mockRepo;
  late InventoryUseCases useCases;

  setUp(() {
    mockRepo = MockInventoryRepository();
    useCases = InventoryUseCases(mockRepo);
  });

  group('getProductInventory', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getProductInventoryResult = [_makeInventoryLevel(id: 1)];

      final result =
          await useCases.getProductInventory('prod-1', businessId: 5);

      expect(mockRepo.calls, ['getProductInventory']);
      expect(mockRepo.lastProductId, 'prod-1');
      expect(mockRepo.lastBusinessId, 5);
      expect(result.length, 1);
      expect(result.first.id, 1);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('network error');

      expect(
        () => useCases.getProductInventory('prod-1'),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getWarehouseInventory', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getWarehouseInventoryResult = PaginatedResponse<InventoryLevel>(
        data: [_makeInventoryLevel()],
        pagination: _makePagination(),
      );

      final params = GetInventoryParams(page: 1, pageSize: 10);
      final result = await useCases.getWarehouseInventory(5, params);

      expect(mockRepo.calls, ['getWarehouseInventory']);
      expect(mockRepo.lastWarehouseId, 5);
      expect(mockRepo.lastInventoryParams, params);
      expect(result.data.length, 1);
    });

    test('delegates with null params', () async {
      mockRepo.getWarehouseInventoryResult = PaginatedResponse<InventoryLevel>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getWarehouseInventory(1, null);

      expect(mockRepo.calls, ['getWarehouseInventory']);
      expect(mockRepo.lastInventoryParams, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('fetch error');

      expect(
        () => useCases.getWarehouseInventory(1, null),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('adjustStock', () {
    test('delegates to repository and returns result', () async {
      mockRepo.adjustStockResult = _makeMovement(id: 99);

      final dto = AdjustStockDTO(
        productId: 'prod-1',
        warehouseId: 5,
        quantity: 10,
        reason: 'recount',
      );
      final result = await useCases.adjustStock(dto, businessId: 3);

      expect(mockRepo.calls, ['adjustStock']);
      expect(mockRepo.lastAdjustDTO, dto);
      expect(mockRepo.lastBusinessId, 3);
      expect(result.id, 99);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('adjust error');

      expect(
        () => useCases.adjustStock(
          AdjustStockDTO(
            productId: 'p1',
            warehouseId: 1,
            quantity: 1,
            reason: 'r',
          ),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('transferStock', () {
    test('delegates to repository and returns result', () async {
      mockRepo.transferStockResult = {'success': true};

      final dto = TransferStockDTO(
        productId: 'prod-1',
        fromWarehouseId: 1,
        toWarehouseId: 2,
        quantity: 50,
      );
      final result = await useCases.transferStock(dto, businessId: 7);

      expect(mockRepo.calls, ['transferStock']);
      expect(mockRepo.lastTransferDTO, dto);
      expect(mockRepo.lastBusinessId, 7);
      expect(result['success'], true);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('transfer error');

      expect(
        () => useCases.transferStock(
          TransferStockDTO(
            productId: 'p1',
            fromWarehouseId: 1,
            toWarehouseId: 2,
            quantity: 1,
          ),
        ),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getMovements', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getMovementsResult = PaginatedResponse<StockMovement>(
        data: [_makeMovement()],
        pagination: _makePagination(),
      );

      final params = GetMovementsParams(page: 1, pageSize: 10);
      final result = await useCases.getMovements(params);

      expect(mockRepo.calls, ['getMovements']);
      expect(mockRepo.lastMovementsParams, params);
      expect(result.data.length, 1);
    });

    test('delegates with null params', () async {
      mockRepo.getMovementsResult = PaginatedResponse<StockMovement>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getMovements(null);

      expect(mockRepo.calls, ['getMovements']);
      expect(mockRepo.lastMovementsParams, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('movements error');

      expect(
        () => useCases.getMovements(null),
        throwsA(isA<Exception>()),
      );
    });
  });

  group('getMovementTypes', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getMovementTypesResult = PaginatedResponse<MovementType>(
        data: [_makeMovementType()],
        pagination: _makePagination(),
      );

      final params = GetMovementTypesParams(activeOnly: true);
      final result = await useCases.getMovementTypes(params);

      expect(mockRepo.calls, ['getMovementTypes']);
      expect(mockRepo.lastMovementTypesParams, params);
      expect(result.data.length, 1);
    });

    test('delegates with null params', () async {
      mockRepo.getMovementTypesResult = PaginatedResponse<MovementType>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getMovementTypes(null);

      expect(mockRepo.calls, ['getMovementTypes']);
      expect(mockRepo.lastMovementTypesParams, isNull);
    });

    test('propagates error from repository', () {
      mockRepo.errorToThrow = Exception('types error');

      expect(
        () => useCases.getMovementTypes(null),
        throwsA(isA<Exception>()),
      );
    });
  });
}
