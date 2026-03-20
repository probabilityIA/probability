import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/inventory/domain/entities.dart';

void main() {
  // =========================================================================
  // InventoryLevel
  // =========================================================================
  group('InventoryLevel', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'product_id': 'prod-abc',
        'warehouse_id': 10,
        'location_id': 5,
        'business_id': 3,
        'quantity': 100,
        'reserved_qty': 20,
        'available_qty': 80,
        'min_stock': 10,
        'max_stock': 500,
        'reorder_point': 25,
        'product_name': 'Widget',
        'product_sku': 'WDG-001',
        'warehouse_name': 'Main Warehouse',
        'warehouse_code': 'MW-01',
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
      };

      final level = InventoryLevel.fromJson(json);

      expect(level.id, 1);
      expect(level.productId, 'prod-abc');
      expect(level.warehouseId, 10);
      expect(level.locationId, 5);
      expect(level.businessId, 3);
      expect(level.quantity, 100);
      expect(level.reservedQty, 20);
      expect(level.availableQty, 80);
      expect(level.minStock, 10);
      expect(level.maxStock, 500);
      expect(level.reorderPoint, 25);
      expect(level.productName, 'Widget');
      expect(level.productSku, 'WDG-001');
      expect(level.warehouseName, 'Main Warehouse');
      expect(level.warehouseCode, 'MW-01');
      expect(level.createdAt, '2026-01-01T00:00:00Z');
      expect(level.updatedAt, '2026-01-02T00:00:00Z');
    });

    test('fromJson handles defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final level = InventoryLevel.fromJson(json);

      expect(level.id, 0);
      expect(level.productId, '');
      expect(level.warehouseId, 0);
      expect(level.businessId, 0);
      expect(level.quantity, 0);
      expect(level.reservedQty, 0);
      expect(level.availableQty, 0);
      expect(level.createdAt, '');
      expect(level.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'product_id': 'p1',
        'warehouse_id': 2,
        'business_id': 3,
        'quantity': 10,
        'reserved_qty': 0,
        'available_qty': 10,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-01',
      };

      final level = InventoryLevel.fromJson(json);

      expect(level.locationId, isNull);
      expect(level.minStock, isNull);
      expect(level.maxStock, isNull);
      expect(level.reorderPoint, isNull);
      expect(level.productName, isNull);
      expect(level.productSku, isNull);
      expect(level.warehouseName, isNull);
      expect(level.warehouseCode, isNull);
    });

    test('fromJson converts product_id to string', () {
      final json = {
        'id': 1,
        'product_id': 12345,
        'warehouse_id': 1,
        'business_id': 1,
        'quantity': 0,
        'reserved_qty': 0,
        'available_qty': 0,
        'created_at': '',
        'updated_at': '',
      };

      final level = InventoryLevel.fromJson(json);

      expect(level.productId, '12345');
    });
  });

  // =========================================================================
  // StockMovement
  // =========================================================================
  group('StockMovement', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'product_id': 'prod-xyz',
        'warehouse_id': 5,
        'location_id': 3,
        'business_id': 7,
        'movement_type_id': 2,
        'movement_type_code': 'ADJ',
        'movement_type_name': 'Adjustment',
        'reason': 'Inventory count correction',
        'quantity': 15,
        'previous_qty': 100,
        'new_qty': 115,
        'reference_type': 'order',
        'reference_id': 'ORD-001',
        'integration_id': 9,
        'notes': 'Manual adjustment',
        'created_by_id': 1,
        'product_name': 'Widget',
        'product_sku': 'WDG-001',
        'warehouse_name': 'Main',
        'created_at': '2026-03-01T10:00:00Z',
      };

      final movement = StockMovement.fromJson(json);

      expect(movement.id, 42);
      expect(movement.productId, 'prod-xyz');
      expect(movement.warehouseId, 5);
      expect(movement.locationId, 3);
      expect(movement.businessId, 7);
      expect(movement.movementTypeId, 2);
      expect(movement.movementTypeCode, 'ADJ');
      expect(movement.movementTypeName, 'Adjustment');
      expect(movement.reason, 'Inventory count correction');
      expect(movement.quantity, 15);
      expect(movement.previousQty, 100);
      expect(movement.newQty, 115);
      expect(movement.referenceType, 'order');
      expect(movement.referenceId, 'ORD-001');
      expect(movement.integrationId, 9);
      expect(movement.notes, 'Manual adjustment');
      expect(movement.createdById, 1);
      expect(movement.productName, 'Widget');
      expect(movement.productSku, 'WDG-001');
      expect(movement.warehouseName, 'Main');
      expect(movement.createdAt, '2026-03-01T10:00:00Z');
    });

    test('fromJson handles defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final movement = StockMovement.fromJson(json);

      expect(movement.id, 0);
      expect(movement.productId, '');
      expect(movement.warehouseId, 0);
      expect(movement.businessId, 0);
      expect(movement.movementTypeId, 0);
      expect(movement.movementTypeCode, '');
      expect(movement.movementTypeName, '');
      expect(movement.reason, '');
      expect(movement.quantity, 0);
      expect(movement.previousQty, 0);
      expect(movement.newQty, 0);
      expect(movement.notes, '');
      expect(movement.createdAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'product_id': 'p1',
        'warehouse_id': 1,
        'business_id': 1,
        'movement_type_id': 1,
        'movement_type_code': 'ADJ',
        'movement_type_name': 'Adjustment',
        'reason': 'test',
        'quantity': 1,
        'previous_qty': 0,
        'new_qty': 1,
        'notes': '',
        'created_at': '',
      };

      final movement = StockMovement.fromJson(json);

      expect(movement.locationId, isNull);
      expect(movement.referenceType, isNull);
      expect(movement.referenceId, isNull);
      expect(movement.integrationId, isNull);
      expect(movement.createdById, isNull);
      expect(movement.productName, isNull);
      expect(movement.productSku, isNull);
      expect(movement.warehouseName, isNull);
    });
  });

  // =========================================================================
  // MovementType
  // =========================================================================
  group('MovementType', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'ADJ',
        'name': 'Adjustment',
        'description': 'Stock adjustment',
        'is_active': true,
        'direction': 'in',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
      };

      final type = MovementType.fromJson(json);

      expect(type.id, 1);
      expect(type.code, 'ADJ');
      expect(type.name, 'Adjustment');
      expect(type.description, 'Stock adjustment');
      expect(type.isActive, true);
      expect(type.direction, 'in');
      expect(type.createdAt, '2026-01-01');
      expect(type.updatedAt, '2026-01-02');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final type = MovementType.fromJson(json);

      expect(type.id, 0);
      expect(type.code, '');
      expect(type.name, '');
      expect(type.description, '');
      expect(type.isActive, false);
      expect(type.direction, '');
      expect(type.createdAt, '');
      expect(type.updatedAt, '');
    });
  });

  // =========================================================================
  // GetInventoryParams
  // =========================================================================
  group('GetInventoryParams', () {
    test('toQueryParams includes all set fields', () {
      final params = GetInventoryParams(
        page: 2,
        pageSize: 20,
        search: 'widget',
        lowStock: true,
        businessId: 5,
      );

      final query = params.toQueryParams();

      expect(query['page'], 2);
      expect(query['page_size'], 20);
      expect(query['search'], 'widget');
      expect(query['low_stock'], true);
      expect(query['business_id'], 5);
    });

    test('toQueryParams omits null fields', () {
      final params = GetInventoryParams();
      final query = params.toQueryParams();
      expect(query, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final params = GetInventoryParams(page: 1, lowStock: false);
      final query = params.toQueryParams();

      expect(query['page'], 1);
      expect(query['low_stock'], false);
      expect(query.containsKey('page_size'), false);
      expect(query.containsKey('search'), false);
      expect(query.containsKey('business_id'), false);
    });
  });

  // =========================================================================
  // GetMovementsParams
  // =========================================================================
  group('GetMovementsParams', () {
    test('toQueryParams includes all set fields', () {
      final params = GetMovementsParams(
        page: 3,
        pageSize: 50,
        productId: 'prod-1',
        warehouseId: 10,
        type: 'in',
        businessId: 7,
      );

      final query = params.toQueryParams();

      expect(query['page'], 3);
      expect(query['page_size'], 50);
      expect(query['product_id'], 'prod-1');
      expect(query['warehouse_id'], 10);
      expect(query['type'], 'in');
      expect(query['business_id'], 7);
    });

    test('toQueryParams omits null fields', () {
      final params = GetMovementsParams();
      final query = params.toQueryParams();
      expect(query, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final params = GetMovementsParams(productId: 'p1');
      final query = params.toQueryParams();

      expect(query.length, 1);
      expect(query['product_id'], 'p1');
    });
  });

  // =========================================================================
  // GetMovementTypesParams
  // =========================================================================
  group('GetMovementTypesParams', () {
    test('toQueryParams includes all set fields', () {
      final params = GetMovementTypesParams(
        page: 1,
        pageSize: 10,
        activeOnly: true,
        businessId: 2,
      );

      final query = params.toQueryParams();

      expect(query['page'], 1);
      expect(query['page_size'], 10);
      expect(query['active_only'], true);
      expect(query['business_id'], 2);
    });

    test('toQueryParams omits null fields', () {
      final params = GetMovementTypesParams();
      final query = params.toQueryParams();
      expect(query, isEmpty);
    });
  });

  // =========================================================================
  // AdjustStockDTO
  // =========================================================================
  group('AdjustStockDTO', () {
    test('toJson includes all required fields', () {
      final dto = AdjustStockDTO(
        productId: 'prod-1',
        warehouseId: 5,
        quantity: 10,
        reason: 'recount',
      );

      final json = dto.toJson();

      expect(json['product_id'], 'prod-1');
      expect(json['warehouse_id'], 5);
      expect(json['quantity'], 10);
      expect(json['reason'], 'recount');
      expect(json.containsKey('location_id'), false);
      expect(json.containsKey('notes'), false);
    });

    test('toJson includes optional fields when set', () {
      final dto = AdjustStockDTO(
        productId: 'prod-1',
        warehouseId: 5,
        locationId: 3,
        quantity: 10,
        reason: 'recount',
        notes: 'Some notes',
      );

      final json = dto.toJson();

      expect(json['location_id'], 3);
      expect(json['notes'], 'Some notes');
    });

    test('toJson omits null optional fields', () {
      final dto = AdjustStockDTO(
        productId: 'prod-1',
        warehouseId: 5,
        quantity: 10,
        reason: 'recount',
      );

      final json = dto.toJson();

      expect(json.containsKey('location_id'), false);
      expect(json.containsKey('notes'), false);
    });
  });

  // =========================================================================
  // TransferStockDTO
  // =========================================================================
  group('TransferStockDTO', () {
    test('toJson includes all required fields', () {
      final dto = TransferStockDTO(
        productId: 'prod-1',
        fromWarehouseId: 1,
        toWarehouseId: 2,
        quantity: 50,
      );

      final json = dto.toJson();

      expect(json['product_id'], 'prod-1');
      expect(json['from_warehouse_id'], 1);
      expect(json['to_warehouse_id'], 2);
      expect(json['quantity'], 50);
      expect(json.containsKey('from_location_id'), false);
      expect(json.containsKey('to_location_id'), false);
      expect(json.containsKey('reason'), false);
      expect(json.containsKey('notes'), false);
    });

    test('toJson includes optional fields when set', () {
      final dto = TransferStockDTO(
        productId: 'prod-1',
        fromWarehouseId: 1,
        toWarehouseId: 2,
        fromLocationId: 10,
        toLocationId: 20,
        quantity: 50,
        reason: 'relocation',
        notes: 'Moving stock',
      );

      final json = dto.toJson();

      expect(json['from_location_id'], 10);
      expect(json['to_location_id'], 20);
      expect(json['reason'], 'relocation');
      expect(json['notes'], 'Moving stock');
    });

    test('toJson omits null optional fields', () {
      final dto = TransferStockDTO(
        productId: 'prod-1',
        fromWarehouseId: 1,
        toWarehouseId: 2,
        quantity: 50,
      );

      final json = dto.toJson();

      expect(json.containsKey('from_location_id'), false);
      expect(json.containsKey('to_location_id'), false);
      expect(json.containsKey('reason'), false);
      expect(json.containsKey('notes'), false);
    });
  });
}
