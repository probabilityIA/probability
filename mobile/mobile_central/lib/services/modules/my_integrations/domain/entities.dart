class IntegrationStats {
  const IntegrationStats({
    required this.integrationId,
    this.ordersCount = 0,
    this.ordersInProgress = 0,
    this.ordersDelivered = 0,
    this.ordersCancelled = 0,
    this.ordersReturned = 0,
    this.productsCount = 0,
    this.lastOrderAt,
  });

  final int integrationId;
  final int ordersCount;
  final int ordersInProgress;
  final int ordersDelivered;
  final int ordersCancelled;
  final int ordersReturned;
  final int productsCount;
  final String? lastOrderAt;

  bool get hasOrders => ordersCount > 0;

  static int _int(dynamic v) {
    if (v is int) return v;
    if (v is num) return v.round();
    return 0;
  }

  factory IntegrationStats.fromJson(Map<String, dynamic> json) {
    return IntegrationStats(
      integrationId: _int(json['integration_id']),
      ordersCount: _int(json['orders_count']),
      ordersInProgress: _int(json['orders_in_progress']),
      ordersDelivered: _int(json['orders_delivered']),
      ordersCancelled: _int(json['orders_cancelled']),
      ordersReturned: _int(json['orders_returned']),
      productsCount: _int(json['products_count']),
      lastOrderAt: json['last_order_at'],
    );
  }

  IntegrationStats operator +(IntegrationStats other) => IntegrationStats(
        integrationId: 0,
        ordersCount: ordersCount + other.ordersCount,
        ordersInProgress: ordersInProgress + other.ordersInProgress,
        ordersDelivered: ordersDelivered + other.ordersDelivered,
        ordersCancelled: ordersCancelled + other.ordersCancelled,
        ordersReturned: ordersReturned + other.ordersReturned,
        productsCount: productsCount + other.productsCount,
      );
}

class MyIntegration {
  final int id;
  final String createdAt;
  final String updatedAt;
  final String? deletedAt;
  final int businessId;
  final int integrationTypeId;
  final String? integrationTypeName;
  final String? integrationTypeCode;
  final String? categoryCode;
  final String? imageUrl;
  final String name;
  final bool isActive;
  final Map<String, dynamic>? credentials;
  final Map<String, dynamic>? config;

  MyIntegration({
    required this.id,
    required this.createdAt,
    required this.updatedAt,
    this.deletedAt,
    required this.businessId,
    required this.integrationTypeId,
    this.integrationTypeName,
    this.integrationTypeCode,
    this.categoryCode,
    this.imageUrl,
    required this.name,
    required this.isActive,
    this.credentials,
    this.config,
  });

  factory MyIntegration.fromJson(Map<String, dynamic> json) {
    return MyIntegration(
      id: json['id'] ?? 0,
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      deletedAt: json['deleted_at'],
      businessId: json['business_id'] ?? 0,
      integrationTypeId: json['integration_type_id'] ?? 0,
      integrationTypeName: json['integration_type_name'] ?? json['integration_type']?['name'],
      integrationTypeCode: json['integration_type_code'] ?? json['integration_type']?['code'],
      categoryCode: json['category_code'] ?? json['category'],
      imageUrl: json['integration_type']?['image_url'] ?? json['image_url'],
      name: json['name'] ?? '',
      isActive: json['is_active'] ?? false,
      credentials: json['credentials'] != null
          ? Map<String, dynamic>.from(json['credentials'])
          : null,
      config: json['config'] != null
          ? Map<String, dynamic>.from(json['config'])
          : null,
    );
  }
}

class IntegrationCategory {
  final String code;
  final String icon;

  IntegrationCategory({
    required this.code,
    required this.icon,
  });
}

/// Channel codes: where orders originate (parallel)
const List<String> channelCodes = ['platform', 'ecommerce'];

/// Service codes: where orders are processed (independent from hub)
const List<String> serviceCodes = ['messaging', 'invoicing', 'shipping', 'payment'];

/// Category icon mapping
const Map<String, String> categoryIcons = {
  'platform': 'puzzle',
  'ecommerce': 'cart',
  'invoicing': 'receipt',
  'messaging': 'chat',
  'payment': 'credit_card',
  'shipping': 'local_shipping',
};

class GetMyIntegrationsParams {
  final int? page;
  final int? pageSize;
  final int? businessId;
  final String? categoryCode;
  final bool? isActive;

  GetMyIntegrationsParams({
    this.page,
    this.pageSize,
    this.businessId,
    this.categoryCode,
    this.isActive,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (businessId != null) params['business_id'] = businessId;
    if (categoryCode != null) params['category_code'] = categoryCode;
    if (isActive != null) params['is_active'] = isActive;
    return params;
  }
}
