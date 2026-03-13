class IntegrationConfig {
  final Map<String, dynamic> data;

  IntegrationConfig({this.data = const {}});

  factory IntegrationConfig.fromJson(Map<String, dynamic>? json) {
    return IntegrationConfig(data: json ?? {});
  }

  Map<String, dynamic> toJson() => data;
}

class IntegrationCredentials {
  final Map<String, dynamic> data;

  IntegrationCredentials({this.data = const {}});

  factory IntegrationCredentials.fromJson(Map<String, dynamic>? json) {
    return IntegrationCredentials(data: json ?? {});
  }

  Map<String, dynamic> toJson() => data;
}

class IntegrationTypeInfo {
  final int id;
  final String name;
  final String code;
  final String? imageUrl;

  IntegrationTypeInfo({
    required this.id,
    required this.name,
    required this.code,
    this.imageUrl,
  });

  factory IntegrationTypeInfo.fromJson(Map<String, dynamic> json) {
    return IntegrationTypeInfo(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      imageUrl: json['image_url'],
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'id': id,
      'name': name,
      'code': code,
    };
    if (imageUrl != null) json['image_url'] = imageUrl;
    return json;
  }
}

class Integration {
  final int id;
  final String name;
  final String code;
  final int integrationTypeId;
  final String type;
  final String category;
  final String? categoryName;
  final String? categoryColor;
  final int? businessId;
  final String? businessName;
  final String? storeId;
  final bool isActive;
  final bool isDefault;
  final bool isTesting;
  final IntegrationConfig config;
  final IntegrationCredentials? credentials;
  final String? description;
  final int createdById;
  final int? updatedById;
  final String createdAt;
  final String updatedAt;
  final IntegrationTypeInfo? integrationType;

  Integration({
    required this.id,
    required this.name,
    required this.code,
    required this.integrationTypeId,
    required this.type,
    required this.category,
    this.categoryName,
    this.categoryColor,
    this.businessId,
    this.businessName,
    this.storeId,
    required this.isActive,
    required this.isDefault,
    required this.isTesting,
    required this.config,
    this.credentials,
    this.description,
    required this.createdById,
    this.updatedById,
    required this.createdAt,
    required this.updatedAt,
    this.integrationType,
  });

  factory Integration.fromJson(Map<String, dynamic> json) {
    return Integration(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      integrationTypeId: json['integration_type_id'] ?? 0,
      type: json['type'] ?? '',
      category: json['category'] ?? '',
      categoryName: json['category_name'],
      categoryColor: json['category_color'],
      businessId: json['business_id'],
      businessName: json['business_name'],
      storeId: json['store_id'],
      isActive: json['is_active'] ?? false,
      isDefault: json['is_default'] ?? false,
      isTesting: json['is_testing'] ?? false,
      config: IntegrationConfig.fromJson(
          json['config'] as Map<String, dynamic>?),
      credentials: json['credentials'] != null
          ? IntegrationCredentials.fromJson(
              json['credentials'] as Map<String, dynamic>?)
          : null,
      description: json['description'],
      createdById: json['created_by_id'] ?? 0,
      updatedById: json['updated_by_id'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      integrationType: json['integration_type'] != null
          ? IntegrationTypeInfo.fromJson(json['integration_type'])
          : null,
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'id': id,
      'name': name,
      'code': code,
      'integration_type_id': integrationTypeId,
      'type': type,
      'category': category,
      'is_active': isActive,
      'is_default': isDefault,
      'is_testing': isTesting,
      'config': config.toJson(),
      'created_by_id': createdById,
      'created_at': createdAt,
      'updated_at': updatedAt,
    };
    if (categoryName != null) json['category_name'] = categoryName;
    if (categoryColor != null) json['category_color'] = categoryColor;
    if (businessId != null) json['business_id'] = businessId;
    if (businessName != null) json['business_name'] = businessName;
    if (storeId != null) json['store_id'] = storeId;
    if (credentials != null) json['credentials'] = credentials!.toJson();
    if (description != null) json['description'] = description;
    if (updatedById != null) json['updated_by_id'] = updatedById;
    if (integrationType != null) {
      json['integration_type'] = integrationType!.toJson();
    }
    return json;
  }
}

class CreateIntegrationDTO {
  final String name;
  final String code;
  final int integrationTypeId;
  final String? type;
  final String category;
  final int? businessId;
  final String? storeId;
  final bool? isActive;
  final bool? isDefault;
  final bool? isTesting;
  final IntegrationConfig? config;
  final IntegrationCredentials? credentials;
  final String? description;

  CreateIntegrationDTO({
    required this.name,
    required this.code,
    required this.integrationTypeId,
    this.type,
    required this.category,
    this.businessId,
    this.storeId,
    this.isActive,
    this.isDefault,
    this.isTesting,
    this.config,
    this.credentials,
    this.description,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'code': code,
      'integration_type_id': integrationTypeId,
      'category': category,
    };
    if (type != null) json['type'] = type;
    if (businessId != null) json['business_id'] = businessId;
    if (storeId != null) json['store_id'] = storeId;
    if (isActive != null) json['is_active'] = isActive;
    if (isDefault != null) json['is_default'] = isDefault;
    if (isTesting != null) json['is_testing'] = isTesting;
    if (config != null) json['config'] = config!.toJson();
    if (credentials != null) json['credentials'] = credentials!.toJson();
    if (description != null) json['description'] = description;
    return json;
  }
}

class UpdateIntegrationDTO {
  final String? name;
  final String? code;
  final String? storeId;
  final bool? isActive;
  final bool? isDefault;
  final bool? isTesting;
  final IntegrationConfig? config;
  final IntegrationCredentials? credentials;
  final String? description;

  UpdateIntegrationDTO({
    this.name,
    this.code,
    this.storeId,
    this.isActive,
    this.isDefault,
    this.isTesting,
    this.config,
    this.credentials,
    this.description,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (name != null) json['name'] = name;
    if (code != null) json['code'] = code;
    if (storeId != null) json['store_id'] = storeId;
    if (isActive != null) json['is_active'] = isActive;
    if (isDefault != null) json['is_default'] = isDefault;
    if (isTesting != null) json['is_testing'] = isTesting;
    if (config != null) json['config'] = config!.toJson();
    if (credentials != null) json['credentials'] = credentials!.toJson();
    if (description != null) json['description'] = description;
    return json;
  }
}

class GetIntegrationsParams {
  final int? page;
  final int? pageSize;
  final String? type;
  final String? category;
  final int? categoryId;
  final int? businessId;
  final bool? isActive;
  final String? search;

  GetIntegrationsParams({
    this.page,
    this.pageSize,
    this.type,
    this.category,
    this.categoryId,
    this.businessId,
    this.isActive,
    this.search,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (type != null) params['type'] = type;
    if (category != null) params['category'] = category;
    if (categoryId != null) params['category_id'] = categoryId;
    if (businessId != null) params['business_id'] = businessId;
    if (isActive != null) params['is_active'] = isActive;
    if (search != null) params['search'] = search;
    return params;
  }
}

class ActionResponse {
  final bool success;
  final String message;
  final String? error;

  ActionResponse({
    required this.success,
    required this.message,
    this.error,
  });

  factory ActionResponse.fromJson(Map<String, dynamic> json) {
    return ActionResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      error: json['error'],
    );
  }
}

class IntegrationType {
  final int id;
  final String name;
  final String code;
  final String? description;
  final String? icon;
  final String? imageUrl;
  final IntegrationCategory? category;
  final int? categoryId;
  final IntegrationCategory? integrationCategory;
  final bool isActive;
  final bool? inDevelopment;
  final dynamic configSchema;
  final dynamic credentialsSchema;
  final String? setupInstructions;
  final String? baseUrl;
  final String? baseUrlTest;
  final bool? hasPlatformCredentials;
  final List<String>? platformCredentialKeys;
  final String createdAt;
  final String updatedAt;

  IntegrationType({
    required this.id,
    required this.name,
    required this.code,
    this.description,
    this.icon,
    this.imageUrl,
    this.category,
    this.categoryId,
    this.integrationCategory,
    required this.isActive,
    this.inDevelopment,
    this.configSchema,
    this.credentialsSchema,
    this.setupInstructions,
    this.baseUrl,
    this.baseUrlTest,
    this.hasPlatformCredentials,
    this.platformCredentialKeys,
    required this.createdAt,
    required this.updatedAt,
  });

  factory IntegrationType.fromJson(Map<String, dynamic> json) {
    return IntegrationType(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      description: json['description'],
      icon: json['icon'],
      imageUrl: json['image_url'],
      category: json['category'] != null
          ? IntegrationCategory.fromJson(json['category'])
          : null,
      categoryId: json['category_id'],
      integrationCategory: json['integration_category'] != null
          ? IntegrationCategory.fromJson(json['integration_category'])
          : null,
      isActive: json['is_active'] ?? false,
      inDevelopment: json['in_development'],
      configSchema: json['config_schema'],
      credentialsSchema: json['credentials_schema'],
      setupInstructions: json['setup_instructions'],
      baseUrl: json['base_url'],
      baseUrlTest: json['base_url_test'],
      hasPlatformCredentials: json['has_platform_credentials'],
      platformCredentialKeys: (json['platform_credential_keys'] as List<dynamic>?)
          ?.map((e) => e.toString())
          .toList(),
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'id': id,
      'name': name,
      'code': code,
      'is_active': isActive,
      'created_at': createdAt,
      'updated_at': updatedAt,
    };
    if (description != null) json['description'] = description;
    if (icon != null) json['icon'] = icon;
    if (imageUrl != null) json['image_url'] = imageUrl;
    if (category != null) json['category'] = category!.toJson();
    if (categoryId != null) json['category_id'] = categoryId;
    if (inDevelopment != null) json['in_development'] = inDevelopment;
    if (configSchema != null) json['config_schema'] = configSchema;
    if (credentialsSchema != null) {
      json['credentials_schema'] = credentialsSchema;
    }
    if (setupInstructions != null) {
      json['setup_instructions'] = setupInstructions;
    }
    if (baseUrl != null) json['base_url'] = baseUrl;
    if (baseUrlTest != null) json['base_url_test'] = baseUrlTest;
    if (hasPlatformCredentials != null) {
      json['has_platform_credentials'] = hasPlatformCredentials;
    }
    if (platformCredentialKeys != null) {
      json['platform_credential_keys'] = platformCredentialKeys;
    }
    return json;
  }
}

class CreateIntegrationTypeDTO {
  final String name;
  final String? code;
  final String? description;
  final String? icon;
  final int categoryId;
  final bool? isActive;
  final dynamic configSchema;
  final dynamic credentialsSchema;
  final String? setupInstructions;
  final String? baseUrl;
  final String? baseUrlTest;
  final Map<String, String>? platformCredentials;

  CreateIntegrationTypeDTO({
    required this.name,
    this.code,
    this.description,
    this.icon,
    required this.categoryId,
    this.isActive,
    this.configSchema,
    this.credentialsSchema,
    this.setupInstructions,
    this.baseUrl,
    this.baseUrlTest,
    this.platformCredentials,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'category_id': categoryId,
    };
    if (code != null) json['code'] = code;
    if (description != null) json['description'] = description;
    if (icon != null) json['icon'] = icon;
    if (isActive != null) json['is_active'] = isActive;
    if (configSchema != null) json['config_schema'] = configSchema;
    if (credentialsSchema != null) {
      json['credentials_schema'] = credentialsSchema;
    }
    if (setupInstructions != null) {
      json['setup_instructions'] = setupInstructions;
    }
    if (baseUrl != null) json['base_url'] = baseUrl;
    if (baseUrlTest != null) json['base_url_test'] = baseUrlTest;
    if (platformCredentials != null) {
      json['platform_credentials'] = platformCredentials;
    }
    return json;
  }
}

class UpdateIntegrationTypeDTO {
  final String? name;
  final String? code;
  final String? description;
  final String? icon;
  final int? categoryId;
  final bool? isActive;
  final bool? inDevelopment;
  final dynamic configSchema;
  final dynamic credentialsSchema;
  final String? setupInstructions;
  final bool? removeImage;
  final String? baseUrl;
  final String? baseUrlTest;
  final Map<String, String>? platformCredentials;

  UpdateIntegrationTypeDTO({
    this.name,
    this.code,
    this.description,
    this.icon,
    this.categoryId,
    this.isActive,
    this.inDevelopment,
    this.configSchema,
    this.credentialsSchema,
    this.setupInstructions,
    this.removeImage,
    this.baseUrl,
    this.baseUrlTest,
    this.platformCredentials,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (name != null) json['name'] = name;
    if (code != null) json['code'] = code;
    if (description != null) json['description'] = description;
    if (icon != null) json['icon'] = icon;
    if (categoryId != null) json['category_id'] = categoryId;
    if (isActive != null) json['is_active'] = isActive;
    if (inDevelopment != null) json['in_development'] = inDevelopment;
    if (configSchema != null) json['config_schema'] = configSchema;
    if (credentialsSchema != null) {
      json['credentials_schema'] = credentialsSchema;
    }
    if (setupInstructions != null) {
      json['setup_instructions'] = setupInstructions;
    }
    if (removeImage != null) json['remove_image'] = removeImage;
    if (baseUrl != null) json['base_url'] = baseUrl;
    if (baseUrlTest != null) json['base_url_test'] = baseUrlTest;
    if (platformCredentials != null) {
      json['platform_credentials'] = platformCredentials;
    }
    return json;
  }
}

class WebhookInfo {
  final String url;
  final String method;
  final String description;
  final List<String>? events;

  WebhookInfo({
    required this.url,
    required this.method,
    required this.description,
    this.events,
  });

  factory WebhookInfo.fromJson(Map<String, dynamic> json) {
    return WebhookInfo(
      url: json['url'] ?? '',
      method: json['method'] ?? '',
      description: json['description'] ?? '',
      events: (json['events'] as List<dynamic>?)
          ?.map((e) => e.toString())
          .toList(),
    );
  }
}

class SyncOrdersParams {
  final String? createdAtMin;
  final String? createdAtMax;
  final String? status;
  final String? financialStatus;
  final String? fulfillmentStatus;

  SyncOrdersParams({
    this.createdAtMin,
    this.createdAtMax,
    this.status,
    this.financialStatus,
    this.fulfillmentStatus,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (createdAtMin != null) json['created_at_min'] = createdAtMin;
    if (createdAtMax != null) json['created_at_max'] = createdAtMax;
    if (status != null) json['status'] = status;
    if (financialStatus != null) json['financial_status'] = financialStatus;
    if (fulfillmentStatus != null) {
      json['fulfillment_status'] = fulfillmentStatus;
    }
    return json;
  }
}

class IntegrationSimple {
  final int id;
  final String name;
  final String type;
  final String category;
  final String categoryName;
  final String? categoryColor;
  final String? imageUrl;
  final int? businessId;
  final bool isActive;

  IntegrationSimple({
    required this.id,
    required this.name,
    required this.type,
    required this.category,
    required this.categoryName,
    this.categoryColor,
    this.imageUrl,
    this.businessId,
    required this.isActive,
  });

  factory IntegrationSimple.fromJson(Map<String, dynamic> json) {
    return IntegrationSimple(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      type: json['type'] ?? '',
      category: json['category'] ?? '',
      categoryName: json['category_name'] ?? '',
      categoryColor: json['category_color'],
      imageUrl: json['image_url'],
      businessId: json['business_id'],
      isActive: json['is_active'] ?? false,
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'id': id,
      'name': name,
      'type': type,
      'category': category,
      'category_name': categoryName,
      'is_active': isActive,
    };
    if (categoryColor != null) json['category_color'] = categoryColor;
    if (imageUrl != null) json['image_url'] = imageUrl;
    if (businessId != null) json['business_id'] = businessId;
    return json;
  }
}

class IntegrationCategory {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String? icon;
  final String? color;
  final int displayOrder;
  final int? parentCategoryId;
  final bool isActive;
  final bool isVisible;
  final String createdAt;
  final String updatedAt;

  IntegrationCategory({
    required this.id,
    required this.code,
    required this.name,
    this.description,
    this.icon,
    this.color,
    required this.displayOrder,
    this.parentCategoryId,
    required this.isActive,
    required this.isVisible,
    required this.createdAt,
    required this.updatedAt,
  });

  factory IntegrationCategory.fromJson(Map<String, dynamic> json) {
    return IntegrationCategory(
      id: json['id'] ?? 0,
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      icon: json['icon'],
      color: json['color'],
      displayOrder: json['display_order'] ?? 0,
      parentCategoryId: json['parent_category_id'],
      isActive: json['is_active'] ?? false,
      isVisible: json['is_visible'] ?? false,
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'id': id,
      'code': code,
      'name': name,
      'display_order': displayOrder,
      'is_active': isActive,
      'is_visible': isVisible,
      'created_at': createdAt,
      'updated_at': updatedAt,
    };
    if (description != null) json['description'] = description;
    if (icon != null) json['icon'] = icon;
    if (color != null) json['color'] = color;
    if (parentCategoryId != null) {
      json['parent_category_id'] = parentCategoryId;
    }
    return json;
  }
}
