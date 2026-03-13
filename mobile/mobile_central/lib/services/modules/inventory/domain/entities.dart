// ============================================
// ENTITIES
// ============================================

class InventoryLevel {
  final int id;
  final String productId;
  final int warehouseId;
  final int? locationId;
  final int businessId;
  final int quantity;
  final int reservedQty;
  final int availableQty;
  final int? minStock;
  final int? maxStock;
  final int? reorderPoint;
  final String? productName;
  final String? productSku;
  final String? warehouseName;
  final String? warehouseCode;
  final String createdAt;
  final String updatedAt;

  InventoryLevel({
    required this.id,
    required this.productId,
    required this.warehouseId,
    this.locationId,
    required this.businessId,
    required this.quantity,
    required this.reservedQty,
    required this.availableQty,
    this.minStock,
    this.maxStock,
    this.reorderPoint,
    this.productName,
    this.productSku,
    this.warehouseName,
    this.warehouseCode,
    required this.createdAt,
    required this.updatedAt,
  });

  factory InventoryLevel.fromJson(Map<String, dynamic> json) {
    return InventoryLevel(
      id: json['id'] ?? 0,
      productId: json['product_id']?.toString() ?? '',
      warehouseId: json['warehouse_id'] ?? 0,
      locationId: json['location_id'],
      businessId: json['business_id'] ?? 0,
      quantity: json['quantity'] ?? 0,
      reservedQty: json['reserved_qty'] ?? 0,
      availableQty: json['available_qty'] ?? 0,
      minStock: json['min_stock'],
      maxStock: json['max_stock'],
      reorderPoint: json['reorder_point'],
      productName: json['product_name'],
      productSku: json['product_sku'],
      warehouseName: json['warehouse_name'],
      warehouseCode: json['warehouse_code'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class StockMovement {
  final int id;
  final String productId;
  final int warehouseId;
  final int? locationId;
  final int businessId;
  final int movementTypeId;
  final String movementTypeCode;
  final String movementTypeName;
  final String reason;
  final int quantity;
  final int previousQty;
  final int newQty;
  final String? referenceType;
  final String? referenceId;
  final int? integrationId;
  final String notes;
  final int? createdById;
  final String? productName;
  final String? productSku;
  final String? warehouseName;
  final String createdAt;

  StockMovement({
    required this.id,
    required this.productId,
    required this.warehouseId,
    this.locationId,
    required this.businessId,
    required this.movementTypeId,
    required this.movementTypeCode,
    required this.movementTypeName,
    required this.reason,
    required this.quantity,
    required this.previousQty,
    required this.newQty,
    this.referenceType,
    this.referenceId,
    this.integrationId,
    required this.notes,
    this.createdById,
    this.productName,
    this.productSku,
    this.warehouseName,
    required this.createdAt,
  });

  factory StockMovement.fromJson(Map<String, dynamic> json) {
    return StockMovement(
      id: json['id'] ?? 0,
      productId: json['product_id']?.toString() ?? '',
      warehouseId: json['warehouse_id'] ?? 0,
      locationId: json['location_id'],
      businessId: json['business_id'] ?? 0,
      movementTypeId: json['movement_type_id'] ?? 0,
      movementTypeCode: json['movement_type_code'] ?? '',
      movementTypeName: json['movement_type_name'] ?? '',
      reason: json['reason'] ?? '',
      quantity: json['quantity'] ?? 0,
      previousQty: json['previous_qty'] ?? 0,
      newQty: json['new_qty'] ?? 0,
      referenceType: json['reference_type'],
      referenceId: json['reference_id'],
      integrationId: json['integration_id'],
      notes: json['notes'] ?? '',
      createdById: json['created_by_id'],
      productName: json['product_name'],
      productSku: json['product_sku'],
      warehouseName: json['warehouse_name'],
      createdAt: json['created_at'] ?? '',
    );
  }
}

class MovementType {
  final int id;
  final String code;
  final String name;
  final String description;
  final bool isActive;
  final String direction;
  final String createdAt;
  final String updatedAt;

  MovementType({
    required this.id,
    required this.code,
    required this.name,
    required this.description,
    required this.isActive,
    required this.direction,
    required this.createdAt,
    required this.updatedAt,
  });

  factory MovementType.fromJson(Map<String, dynamic> json) {
    return MovementType(
      id: json['id'] ?? 0,
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      isActive: json['is_active'] ?? false,
      direction: json['direction'] ?? '',
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

// ============================================
// DTOs
// ============================================

class GetInventoryParams {
  final int? page;
  final int? pageSize;
  final String? search;
  final bool? lowStock;
  final int? businessId;

  GetInventoryParams({
    this.page,
    this.pageSize,
    this.search,
    this.lowStock,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (search != null) params['search'] = search;
    if (lowStock != null) params['low_stock'] = lowStock;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

class GetMovementsParams {
  final int? page;
  final int? pageSize;
  final String? productId;
  final int? warehouseId;
  final String? type;
  final int? businessId;

  GetMovementsParams({
    this.page,
    this.pageSize,
    this.productId,
    this.warehouseId,
    this.type,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (productId != null) params['product_id'] = productId;
    if (warehouseId != null) params['warehouse_id'] = warehouseId;
    if (type != null) params['type'] = type;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

class GetMovementTypesParams {
  final int? page;
  final int? pageSize;
  final bool? activeOnly;
  final int? businessId;

  GetMovementTypesParams({
    this.page,
    this.pageSize,
    this.activeOnly,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (activeOnly != null) params['active_only'] = activeOnly;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

class AdjustStockDTO {
  final String productId;
  final int warehouseId;
  final int? locationId;
  final int quantity;
  final String reason;
  final String? notes;

  AdjustStockDTO({
    required this.productId,
    required this.warehouseId,
    this.locationId,
    required this.quantity,
    required this.reason,
    this.notes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'product_id': productId,
      'warehouse_id': warehouseId,
      'quantity': quantity,
      'reason': reason,
    };
    if (locationId != null) json['location_id'] = locationId;
    if (notes != null) json['notes'] = notes;
    return json;
  }
}

class TransferStockDTO {
  final String productId;
  final int fromWarehouseId;
  final int toWarehouseId;
  final int? fromLocationId;
  final int? toLocationId;
  final int quantity;
  final String? reason;
  final String? notes;

  TransferStockDTO({
    required this.productId,
    required this.fromWarehouseId,
    required this.toWarehouseId,
    this.fromLocationId,
    this.toLocationId,
    required this.quantity,
    this.reason,
    this.notes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'product_id': productId,
      'from_warehouse_id': fromWarehouseId,
      'to_warehouse_id': toWarehouseId,
      'quantity': quantity,
    };
    if (fromLocationId != null) json['from_location_id'] = fromLocationId;
    if (toLocationId != null) json['to_location_id'] = toLocationId;
    if (reason != null) json['reason'] = reason;
    if (notes != null) json['notes'] = notes;
    return json;
  }
}
