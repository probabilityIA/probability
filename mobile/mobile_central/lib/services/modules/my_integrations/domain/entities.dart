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
      integrationTypeName: json['integration_type_name'],
      integrationTypeCode: json['integration_type_code'],
      categoryCode: json['category_code'],
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
