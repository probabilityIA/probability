class NotificationType {
  final int id;
  final String name;
  final String code;
  final String? description;
  final String? icon;
  final bool isActive;
  final Map<String, dynamic>? configSchema;
  final String createdAt;
  final String updatedAt;

  NotificationType({
    required this.id,
    required this.name,
    required this.code,
    this.description,
    this.icon,
    required this.isActive,
    this.configSchema,
    required this.createdAt,
    required this.updatedAt,
  });

  factory NotificationType.fromJson(Map<String, dynamic> json) {
    return NotificationType(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      description: json['description'],
      icon: json['icon'],
      isActive: json['is_active'] ?? true,
      configSchema: json['config_schema'] != null
          ? Map<String, dynamic>.from(json['config_schema'])
          : null,
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class NotificationEventType {
  final int id;
  final int notificationTypeId;
  final String eventCode;
  final String eventName;
  final String? description;
  final Map<String, dynamic>? templateConfig;
  final bool isActive;
  final List<int>? allowedOrderStatusIds;
  final String createdAt;
  final String updatedAt;
  final NotificationType? notificationType;

  NotificationEventType({
    required this.id,
    required this.notificationTypeId,
    required this.eventCode,
    required this.eventName,
    this.description,
    this.templateConfig,
    required this.isActive,
    this.allowedOrderStatusIds,
    required this.createdAt,
    required this.updatedAt,
    this.notificationType,
  });

  factory NotificationEventType.fromJson(Map<String, dynamic> json) {
    return NotificationEventType(
      id: json['id'] ?? 0,
      notificationTypeId: json['notification_type_id'] ?? 0,
      eventCode: json['event_code'] ?? '',
      eventName: json['event_name'] ?? '',
      description: json['description'],
      templateConfig: json['template_config'] != null
          ? Map<String, dynamic>.from(json['template_config'])
          : null,
      isActive: json['is_active'] ?? true,
      allowedOrderStatusIds: (json['allowed_order_status_ids'] as List<dynamic>?)
          ?.map((e) => e as int)
          .toList(),
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      notificationType: json['notification_type'] != null
          ? NotificationType.fromJson(json['notification_type'])
          : null,
    );
  }
}

class NotificationConfig {
  final int id;
  final int businessId;
  final int integrationId;
  final int notificationTypeId;
  final int notificationEventTypeId;
  final bool enabled;
  final Map<String, dynamic>? filters;
  final String? description;
  final String createdAt;
  final String updatedAt;
  final String? deletedAt;
  final NotificationType? notificationType;
  final NotificationEventType? notificationEventType;
  final List<int>? orderStatusIds;
  final String? notificationTypeName;
  final String? notificationEventName;
  final String? eventType;
  final List<String>? channels;

  NotificationConfig({
    required this.id,
    required this.businessId,
    required this.integrationId,
    required this.notificationTypeId,
    required this.notificationEventTypeId,
    required this.enabled,
    this.filters,
    this.description,
    required this.createdAt,
    required this.updatedAt,
    this.deletedAt,
    this.notificationType,
    this.notificationEventType,
    this.orderStatusIds,
    this.notificationTypeName,
    this.notificationEventName,
    this.eventType,
    this.channels,
  });

  factory NotificationConfig.fromJson(Map<String, dynamic> json) {
    return NotificationConfig(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      integrationId: json['integration_id'] ?? 0,
      notificationTypeId: json['notification_type_id'] ?? 0,
      notificationEventTypeId: json['notification_event_type_id'] ?? 0,
      enabled: json['enabled'] ?? false,
      filters: json['filters'] != null
          ? Map<String, dynamic>.from(json['filters'])
          : null,
      description: json['description'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      deletedAt: json['deleted_at'],
      notificationType: json['notification_type'] != null
          ? NotificationType.fromJson(json['notification_type'])
          : null,
      notificationEventType: json['notification_event_type'] != null
          ? NotificationEventType.fromJson(json['notification_event_type'])
          : null,
      orderStatusIds: (json['order_status_ids'] as List<dynamic>?)
          ?.map((e) => e as int)
          .toList(),
      notificationTypeName: json['notification_type_name'],
      notificationEventName: json['notification_event_name'],
      eventType: json['event_type'],
      channels: (json['channels'] as List<dynamic>?)
          ?.map((e) => e.toString())
          .toList(),
    );
  }
}

class CreateConfigDTO {
  final int businessId;
  final int integrationId;
  final int notificationTypeId;
  final int notificationEventTypeId;
  final bool? enabled;
  final Map<String, dynamic>? filters;
  final String? description;
  final List<int>? orderStatusIds;

  CreateConfigDTO({
    required this.businessId,
    required this.integrationId,
    required this.notificationTypeId,
    required this.notificationEventTypeId,
    this.enabled,
    this.filters,
    this.description,
    this.orderStatusIds,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'business_id': businessId,
      'integration_id': integrationId,
      'notification_type_id': notificationTypeId,
      'notification_event_type_id': notificationEventTypeId,
    };
    if (enabled != null) json['enabled'] = enabled;
    if (filters != null) json['filters'] = filters;
    if (description != null) json['description'] = description;
    if (orderStatusIds != null) json['order_status_ids'] = orderStatusIds;
    return json;
  }
}

class UpdateConfigDTO {
  final int? integrationId;
  final int? notificationTypeId;
  final int? notificationEventTypeId;
  final bool? enabled;
  final Map<String, dynamic>? filters;
  final String? description;
  final List<int>? orderStatusIds;

  UpdateConfigDTO({
    this.integrationId,
    this.notificationTypeId,
    this.notificationEventTypeId,
    this.enabled,
    this.filters,
    this.description,
    this.orderStatusIds,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (integrationId != null) json['integration_id'] = integrationId;
    if (notificationTypeId != null) json['notification_type_id'] = notificationTypeId;
    if (notificationEventTypeId != null) json['notification_event_type_id'] = notificationEventTypeId;
    if (enabled != null) json['enabled'] = enabled;
    if (filters != null) json['filters'] = filters;
    if (description != null) json['description'] = description;
    if (orderStatusIds != null) json['order_status_ids'] = orderStatusIds;
    return json;
  }
}

class ConfigFilter {
  final int? businessId;
  final int? integrationId;
  final int? notificationTypeId;
  final int? notificationEventTypeId;

  ConfigFilter({
    this.businessId,
    this.integrationId,
    this.notificationTypeId,
    this.notificationEventTypeId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (businessId != null) params['business_id'] = businessId;
    if (integrationId != null) params['integration_id'] = integrationId;
    if (notificationTypeId != null) params['notification_type_id'] = notificationTypeId;
    if (notificationEventTypeId != null) params['notification_event_type_id'] = notificationEventTypeId;
    return params;
  }
}

class SyncRule {
  final int? id;
  final int notificationTypeId;
  final int notificationEventTypeId;
  final bool enabled;
  final String description;
  final List<int> orderStatusIds;

  SyncRule({
    this.id,
    required this.notificationTypeId,
    required this.notificationEventTypeId,
    required this.enabled,
    required this.description,
    required this.orderStatusIds,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'notification_type_id': notificationTypeId,
      'notification_event_type_id': notificationEventTypeId,
      'enabled': enabled,
      'description': description,
      'order_status_ids': orderStatusIds,
    };
    if (id != null) json['id'] = id;
    return json;
  }
}

class SyncConfigsDTO {
  final int integrationId;
  final List<SyncRule> rules;

  SyncConfigsDTO({
    required this.integrationId,
    required this.rules,
  });

  Map<String, dynamic> toJson() {
    return {
      'integration_id': integrationId,
      'rules': rules.map((r) => r.toJson()).toList(),
    };
  }
}

class SyncConfigsResponse {
  final int created;
  final int updated;
  final int deleted;
  final List<NotificationConfig> configs;

  SyncConfigsResponse({
    required this.created,
    required this.updated,
    required this.deleted,
    required this.configs,
  });

  factory SyncConfigsResponse.fromJson(Map<String, dynamic> json) {
    return SyncConfigsResponse(
      created: json['created'] ?? 0,
      updated: json['updated'] ?? 0,
      deleted: json['deleted'] ?? 0,
      configs: (json['configs'] as List<dynamic>?)
              ?.map((e) => NotificationConfig.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class OrderStatus {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String category;
  final bool isActive;
  final String? icon;
  final String? color;
  final String createdAt;
  final String updatedAt;

  OrderStatus({
    required this.id,
    required this.code,
    required this.name,
    this.description,
    required this.category,
    required this.isActive,
    this.icon,
    this.color,
    required this.createdAt,
    required this.updatedAt,
  });

  factory OrderStatus.fromJson(Map<String, dynamic> json) {
    return OrderStatus(
      id: json['id'] ?? 0,
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      category: json['category'] ?? '',
      isActive: json['is_active'] ?? true,
      icon: json['icon'],
      color: json['color'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class Integration {
  final int id;
  final String name;
  final String code;
  final String type;
  final int? businessId;
  final bool isActive;
  final int? integrationTypeId;
  final String? integrationTypeName;
  final String? integrationTypeIcon;
  final String createdAt;
  final String updatedAt;

  Integration({
    required this.id,
    required this.name,
    required this.code,
    required this.type,
    this.businessId,
    required this.isActive,
    this.integrationTypeId,
    this.integrationTypeName,
    this.integrationTypeIcon,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Integration.fromJson(Map<String, dynamic> json) {
    return Integration(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      type: json['type'] ?? '',
      businessId: json['business_id'],
      isActive: json['is_active'] ?? true,
      integrationTypeId: json['integration_type_id'],
      integrationTypeName: json['integration_type_name'],
      integrationTypeIcon: json['integration_type_icon'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class MessageAuditLog {
  final String id;
  final String conversationId;
  final String messageId;
  final String direction;
  final String templateName;
  final String content;
  final String status;
  final String? deliveredAt;
  final String? readAt;
  final String createdAt;
  final String phoneNumber;
  final String orderNumber;
  final int businessId;

  MessageAuditLog({
    required this.id,
    required this.conversationId,
    required this.messageId,
    required this.direction,
    required this.templateName,
    required this.content,
    required this.status,
    this.deliveredAt,
    this.readAt,
    required this.createdAt,
    required this.phoneNumber,
    required this.orderNumber,
    required this.businessId,
  });

  factory MessageAuditLog.fromJson(Map<String, dynamic> json) {
    return MessageAuditLog(
      id: json['id']?.toString() ?? '',
      conversationId: json['conversation_id'] ?? '',
      messageId: json['message_id'] ?? '',
      direction: json['direction'] ?? '',
      templateName: json['template_name'] ?? '',
      content: json['content'] ?? '',
      status: json['status'] ?? '',
      deliveredAt: json['delivered_at'],
      readAt: json['read_at'],
      createdAt: json['created_at'] ?? '',
      phoneNumber: json['phone_number'] ?? '',
      orderNumber: json['order_number'] ?? '',
      businessId: json['business_id'] ?? 0,
    );
  }
}

class MessageAuditStats {
  final int totalSent;
  final int totalDelivered;
  final int totalRead;
  final int totalFailed;
  final double successRate;

  MessageAuditStats({
    required this.totalSent,
    required this.totalDelivered,
    required this.totalRead,
    required this.totalFailed,
    required this.successRate,
  });

  factory MessageAuditStats.fromJson(Map<String, dynamic> json) {
    return MessageAuditStats(
      totalSent: json['total_sent'] ?? 0,
      totalDelivered: json['total_delivered'] ?? 0,
      totalRead: json['total_read'] ?? 0,
      totalFailed: json['total_failed'] ?? 0,
      successRate: (json['success_rate'] ?? 0).toDouble(),
    );
  }
}

class MessageAuditFilter {
  final int businessId;
  final String? status;
  final String? direction;
  final String? templateName;
  final String? dateFrom;
  final String? dateTo;
  final int? page;
  final int? pageSize;

  MessageAuditFilter({
    required this.businessId,
    this.status,
    this.direction,
    this.templateName,
    this.dateFrom,
    this.dateTo,
    this.page,
    this.pageSize,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{
      'business_id': businessId,
    };
    if (status != null) params['status'] = status;
    if (direction != null) params['direction'] = direction;
    if (templateName != null) params['template_name'] = templateName;
    if (dateFrom != null) params['date_from'] = dateFrom;
    if (dateTo != null) params['date_to'] = dateTo;
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    return params;
  }
}

class PaginatedMessageAuditResponse {
  final List<MessageAuditLog> data;
  final int total;
  final int page;
  final int pageSize;
  final int totalPages;

  PaginatedMessageAuditResponse({
    required this.data,
    required this.total,
    required this.page,
    required this.pageSize,
    required this.totalPages,
  });

  factory PaginatedMessageAuditResponse.fromJson(Map<String, dynamic> json) {
    return PaginatedMessageAuditResponse(
      data: (json['data'] as List<dynamic>?)
              ?.map((e) => MessageAuditLog.fromJson(e))
              .toList() ??
          [],
      total: json['total'] ?? 0,
      page: json['page'] ?? 1,
      pageSize: json['page_size'] ?? 10,
      totalPages: json['total_pages'] ?? 0,
    );
  }
}
