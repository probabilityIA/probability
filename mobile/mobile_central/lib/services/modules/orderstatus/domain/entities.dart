class OrderStatusInfo {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String? category;
  final String? color;
  final int? priority;
  final bool? isActive;

  OrderStatusInfo({required this.id, required this.code, required this.name, this.description, this.category, this.color, this.priority, this.isActive});

  factory OrderStatusInfo.fromJson(Map<String, dynamic> json) {
    return OrderStatusInfo(
      id: json['id'] ?? 0, code: json['code'] ?? '', name: json['name'] ?? '',
      description: json['description'], category: json['category'], color: json['color'],
      priority: json['priority'], isActive: json['is_active'],
    );
  }
}

class IntegrationTypeInfo {
  final int id;
  final String code;
  final String name;
  final String? imageUrl;

  IntegrationTypeInfo({required this.id, required this.code, required this.name, this.imageUrl});

  factory IntegrationTypeInfo.fromJson(Map<String, dynamic> json) {
    return IntegrationTypeInfo(id: json['id'] ?? 0, code: json['code'] ?? '', name: json['name'] ?? '', imageUrl: json['image_url']);
  }
}

class OrderStatusMapping {
  final int id;
  final int integrationTypeId;
  final IntegrationTypeInfo? integrationType;
  final String originalStatus;
  final int orderStatusId;
  final OrderStatusInfo? orderStatus;
  final bool isActive;
  final String description;
  final String createdAt;
  final String updatedAt;

  OrderStatusMapping({
    required this.id, required this.integrationTypeId, this.integrationType,
    required this.originalStatus, required this.orderStatusId, this.orderStatus,
    required this.isActive, required this.description, required this.createdAt, required this.updatedAt,
  });

  factory OrderStatusMapping.fromJson(Map<String, dynamic> json) {
    return OrderStatusMapping(
      id: json['id'] ?? 0, integrationTypeId: json['integration_type_id'] ?? 0,
      integrationType: json['integration_type'] != null ? IntegrationTypeInfo.fromJson(json['integration_type']) : null,
      originalStatus: json['original_status'] ?? '', orderStatusId: json['order_status_id'] ?? 0,
      orderStatus: json['order_status'] != null ? OrderStatusInfo.fromJson(json['order_status']) : null,
      isActive: json['is_active'] ?? true, description: json['description'] ?? '',
      createdAt: json['created_at'] ?? '', updatedAt: json['updated_at'] ?? '',
    );
  }
}

class ChannelStatusInfo {
  final int id;
  final int integrationTypeId;
  final IntegrationTypeInfo? integrationType;
  final String code;
  final String name;
  final String? description;
  final bool isActive;
  final int displayOrder;

  ChannelStatusInfo({required this.id, required this.integrationTypeId, this.integrationType, required this.code, required this.name, this.description, required this.isActive, required this.displayOrder});

  factory ChannelStatusInfo.fromJson(Map<String, dynamic> json) {
    return ChannelStatusInfo(
      id: json['id'] ?? 0, integrationTypeId: json['integration_type_id'] ?? 0,
      integrationType: json['integration_type'] != null ? IntegrationTypeInfo.fromJson(json['integration_type']) : null,
      code: json['code'] ?? '', name: json['name'] ?? '', description: json['description'],
      isActive: json['is_active'] ?? true, displayOrder: json['display_order'] ?? 0,
    );
  }
}

class GetOrderStatusMappingsParams {
  final int? page;
  final int? pageSize;
  final int? integrationTypeId;
  final bool? isActive;

  GetOrderStatusMappingsParams({this.page, this.pageSize, this.integrationTypeId, this.isActive});

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (integrationTypeId != null) params['integration_type_id'] = integrationTypeId;
    if (isActive != null) params['is_active'] = isActive;
    return params;
  }
}

class CreateOrderStatusDTO {
  final String code;
  final String name;
  final String? description;
  final String? category;
  final String? color;
  final int? priority;
  final bool? isActive;

  CreateOrderStatusDTO({required this.code, required this.name, this.description, this.category, this.color, this.priority, this.isActive});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'code': code, 'name': name};
    if (description != null) json['description'] = description;
    if (category != null) json['category'] = category;
    if (color != null) json['color'] = color;
    if (priority != null) json['priority'] = priority;
    if (isActive != null) json['is_active'] = isActive;
    return json;
  }
}

class CreateOrderStatusMappingDTO {
  final int integrationTypeId;
  final String originalStatus;
  final int orderStatusId;
  final String? description;

  CreateOrderStatusMappingDTO({required this.integrationTypeId, required this.originalStatus, required this.orderStatusId, this.description});

  Map<String, dynamic> toJson() => {
    'integration_type_id': integrationTypeId, 'original_status': originalStatus,
    'order_status_id': orderStatusId, if (description != null) 'description': description,
  };
}

class UpdateOrderStatusMappingDTO {
  final String originalStatus;
  final int orderStatusId;
  final String? description;

  UpdateOrderStatusMappingDTO({required this.originalStatus, required this.orderStatusId, this.description});

  Map<String, dynamic> toJson() => {
    'original_status': originalStatus, 'order_status_id': orderStatusId,
    if (description != null) 'description': description,
  };
}

class CreateChannelStatusDTO {
  final int integrationTypeId;
  final String code;
  final String name;
  final String? description;
  final bool isActive;
  final int displayOrder;

  CreateChannelStatusDTO({required this.integrationTypeId, required this.code, required this.name, this.description, required this.isActive, required this.displayOrder});

  Map<String, dynamic> toJson() => {
    'integration_type_id': integrationTypeId, 'code': code, 'name': name,
    'is_active': isActive, 'display_order': displayOrder,
    if (description != null) 'description': description,
  };
}

class UpdateChannelStatusDTO {
  final String code;
  final String name;
  final String? description;
  final bool isActive;
  final int displayOrder;

  UpdateChannelStatusDTO({required this.code, required this.name, this.description, required this.isActive, required this.displayOrder});

  Map<String, dynamic> toJson() => {
    'code': code, 'name': name, 'is_active': isActive, 'display_order': displayOrder,
    if (description != null) 'description': description,
  };
}
