enum DashboardPeriod { today, week, month, all }

extension DashboardPeriodX on DashboardPeriod {
  String get label {
    switch (this) {
      case DashboardPeriod.today:
        return 'Hoy';
      case DashboardPeriod.week:
        return '7 dias';
      case DashboardPeriod.month:
        return '30 dias';
      case DashboardPeriod.all:
        return 'Todo';
    }
  }

  int? get days {
    switch (this) {
      case DashboardPeriod.today:
        return 0;
      case DashboardPeriod.week:
        return 6;
      case DashboardPeriod.month:
        return 29;
      case DashboardPeriod.all:
        return null;
    }
  }

  ({String start, String end})? rangeFrom(DateTime now) {
    final span = days;
    if (span == null) return null;
    final end = DateTime(now.year, now.month, now.day);
    final start = end.subtract(Duration(days: span));
    return (start: _ymd(start), end: _ymd(end));
  }
}

String _ymd(DateTime d) =>
    '${d.year.toString().padLeft(4, '0')}-'
    '${d.month.toString().padLeft(2, '0')}-'
    '${d.day.toString().padLeft(2, '0')}';

int _asInt(dynamic value) {
  if (value is int) return value;
  if (value is num) return value.round();
  if (value is String) return int.tryParse(value) ?? 0;
  return 0;
}

double _asDouble(dynamic value) {
  if (value is double) return value;
  if (value is num) return value.toDouble();
  if (value is String) return double.tryParse(value) ?? 0;
  return 0;
}

int? _asIntOrNull(dynamic value) {
  if (value == null) return null;
  return _asInt(value);
}

class OrderCountByIntegrationType {
  final String integrationType;
  final int count;

  OrderCountByIntegrationType({
    required this.integrationType,
    required this.count,
  });

  factory OrderCountByIntegrationType.fromJson(Map<String, dynamic> json) {
    return OrderCountByIntegrationType(
      integrationType: json['integration_type'] ?? '',
      count: _asInt(json['count']),
    );
  }
}

class TopCustomer {
  final String customerName;
  final String customerEmail;
  final int orderCount;

  TopCustomer({
    required this.customerName,
    required this.customerEmail,
    required this.orderCount,
  });

  factory TopCustomer.fromJson(Map<String, dynamic> json) {
    return TopCustomer(
      customerName: json['customer_name'] ?? '',
      customerEmail: json['customer_email'] ?? '',
      orderCount: _asInt(json['order_count']),
    );
  }
}

class OrderCountByLocation {
  final String city;
  final String state;
  final int orderCount;

  OrderCountByLocation({
    required this.city,
    required this.state,
    required this.orderCount,
  });

  factory OrderCountByLocation.fromJson(Map<String, dynamic> json) {
    return OrderCountByLocation(
      city: json['city'] ?? '',
      state: json['state'] ?? '',
      orderCount: _asInt(json['order_count']),
    );
  }
}

class TopDriver {
  final String driverName;
  final int? driverId;
  final int orderCount;

  TopDriver({
    required this.driverName,
    this.driverId,
    required this.orderCount,
  });

  factory TopDriver.fromJson(Map<String, dynamic> json) {
    return TopDriver(
      driverName: json['driver_name'] ?? '',
      driverId: _asIntOrNull(json['driver_id']),
      orderCount: _asInt(json['order_count']),
    );
  }
}

class DriverByLocation {
  final String driverName;
  final String city;
  final String state;
  final int orderCount;

  DriverByLocation({
    required this.driverName,
    required this.city,
    required this.state,
    required this.orderCount,
  });

  factory DriverByLocation.fromJson(Map<String, dynamic> json) {
    return DriverByLocation(
      driverName: json['driver_name'] ?? '',
      city: json['city'] ?? '',
      state: json['state'] ?? '',
      orderCount: _asInt(json['order_count']),
    );
  }
}

class TopProduct {
  final String productName;
  final String productId;
  final String sku;
  final int orderCount;
  final double totalSold;

  TopProduct({
    required this.productName,
    required this.productId,
    required this.sku,
    required this.orderCount,
    required this.totalSold,
  });

  factory TopProduct.fromJson(Map<String, dynamic> json) {
    return TopProduct(
      productName: json['product_name'] ?? '',
      productId: json['product_id'] ?? '',
      sku: json['sku'] ?? '',
      orderCount: _asInt(json['order_count']),
      totalSold: _asDouble(json['total_sold']),
    );
  }
}

class ProductByCategory {
  final String category;
  final int count;

  ProductByCategory({
    required this.category,
    required this.count,
  });

  factory ProductByCategory.fromJson(Map<String, dynamic> json) {
    return ProductByCategory(
      category: json['category'] ?? '',
      count: _asInt(json['count']),
    );
  }
}

class ProductByBrand {
  final String brand;
  final int count;

  ProductByBrand({
    required this.brand,
    required this.count,
  });

  factory ProductByBrand.fromJson(Map<String, dynamic> json) {
    return ProductByBrand(
      brand: json['brand'] ?? '',
      count: _asInt(json['count']),
    );
  }
}

class ShipmentsByStatus {
  final String status;
  final int count;

  ShipmentsByStatus({
    required this.status,
    required this.count,
  });

  factory ShipmentsByStatus.fromJson(Map<String, dynamic> json) {
    return ShipmentsByStatus(
      status: json['status'] ?? '',
      count: _asInt(json['count']),
    );
  }
}

class ShipmentsByCarrier {
  final String carrier;
  final int count;

  ShipmentsByCarrier({
    required this.carrier,
    required this.count,
  });

  factory ShipmentsByCarrier.fromJson(Map<String, dynamic> json) {
    return ShipmentsByCarrier(
      carrier: json['carrier'] ?? '',
      count: _asInt(json['count']),
    );
  }
}

class ShipmentsByWarehouse {
  final String warehouseName;
  final int? warehouseId;
  final int count;

  ShipmentsByWarehouse({
    required this.warehouseName,
    this.warehouseId,
    required this.count,
  });

  factory ShipmentsByWarehouse.fromJson(Map<String, dynamic> json) {
    return ShipmentsByWarehouse(
      warehouseName: json['warehouse_name'] ?? '',
      warehouseId: _asIntOrNull(json['warehouse_id']),
      count: _asInt(json['count']),
    );
  }
}

class OrdersByBusiness {
  final int businessId;
  final String businessName;
  final int orderCount;

  OrdersByBusiness({
    required this.businessId,
    required this.businessName,
    required this.orderCount,
  });

  factory OrdersByBusiness.fromJson(Map<String, dynamic> json) {
    return OrdersByBusiness(
      businessId: _asInt(json['business_id']),
      businessName: json['business_name'] ?? '',
      orderCount: _asInt(json['order_count']),
    );
  }
}

class DashboardStats {
  final int totalOrders;
  final List<OrderCountByIntegrationType> ordersByIntegrationType;
  final List<TopCustomer> topCustomers;
  final List<OrderCountByLocation> ordersByLocation;
  final List<TopDriver> topDrivers;
  final List<DriverByLocation> driversByLocation;
  final List<TopProduct> topProducts;
  final List<ProductByCategory> productsByCategory;
  final List<ProductByBrand> productsByBrand;
  final List<ShipmentsByStatus> shipmentsByStatus;
  final List<ShipmentsByCarrier> shipmentsByCarrier;
  final List<ShipmentsByWarehouse> shipmentsByWarehouse;
  final List<OrdersByBusiness>? ordersByBusiness;

  DashboardStats({
    required this.totalOrders,
    required this.ordersByIntegrationType,
    required this.topCustomers,
    required this.ordersByLocation,
    required this.topDrivers,
    required this.driversByLocation,
    required this.topProducts,
    required this.productsByCategory,
    required this.productsByBrand,
    required this.shipmentsByStatus,
    required this.shipmentsByCarrier,
    required this.shipmentsByWarehouse,
    this.ordersByBusiness,
  });

  factory DashboardStats.fromJson(Map<String, dynamic> json) {
    return DashboardStats(
      totalOrders: _asInt(json['total_orders']),
      ordersByIntegrationType:
          (json['orders_by_integration_type'] as List<dynamic>?)
                  ?.map((e) => OrderCountByIntegrationType.fromJson(e))
                  .toList() ??
              [],
      topCustomers: (json['top_customers'] as List<dynamic>?)
              ?.map((e) => TopCustomer.fromJson(e))
              .toList() ??
          [],
      ordersByLocation: (json['orders_by_location'] as List<dynamic>?)
              ?.map((e) => OrderCountByLocation.fromJson(e))
              .toList() ??
          [],
      topDrivers: (json['top_drivers'] as List<dynamic>?)
              ?.map((e) => TopDriver.fromJson(e))
              .toList() ??
          [],
      driversByLocation: (json['drivers_by_location'] as List<dynamic>?)
              ?.map((e) => DriverByLocation.fromJson(e))
              .toList() ??
          [],
      topProducts: (json['top_products'] as List<dynamic>?)
              ?.map((e) => TopProduct.fromJson(e))
              .toList() ??
          [],
      productsByCategory: (json['products_by_category'] as List<dynamic>?)
              ?.map((e) => ProductByCategory.fromJson(e))
              .toList() ??
          [],
      productsByBrand: (json['products_by_brand'] as List<dynamic>?)
              ?.map((e) => ProductByBrand.fromJson(e))
              .toList() ??
          [],
      shipmentsByStatus: (json['shipments_by_status'] as List<dynamic>?)
              ?.map((e) => ShipmentsByStatus.fromJson(e))
              .toList() ??
          [],
      shipmentsByCarrier: (json['shipments_by_carrier'] as List<dynamic>?)
              ?.map((e) => ShipmentsByCarrier.fromJson(e))
              .toList() ??
          [],
      shipmentsByWarehouse: (json['shipments_by_warehouse'] as List<dynamic>?)
              ?.map((e) => ShipmentsByWarehouse.fromJson(e))
              .toList() ??
          [],
      ordersByBusiness: (json['orders_by_business'] as List<dynamic>?)
              ?.map((e) => OrdersByBusiness.fromJson(e))
              .toList(),
    );
  }
}

class DashboardStatsResponse {
  final bool success;
  final String message;
  final DashboardStats data;

  DashboardStatsResponse({
    required this.success,
    required this.message,
    required this.data,
  });

  factory DashboardStatsResponse.fromJson(Map<String, dynamic> json) {
    return DashboardStatsResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      data: DashboardStats.fromJson(json['data'] ?? {}),
    );
  }
}
