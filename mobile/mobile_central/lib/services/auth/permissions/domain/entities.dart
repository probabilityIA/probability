class Permission {
  final int id;
  final String? resource;
  final String? action;
  final String? description;
  final int? scopeId;
  final int? businessTypeId;
  final String? businessTypeName;

  Permission({
    required this.id,
    this.resource,
    this.action,
    this.description,
    this.scopeId,
    this.businessTypeId,
    this.businessTypeName,
  });

  factory Permission.fromJson(Map<String, dynamic> json) {
    return Permission(
      id: json['id'] ?? 0,
      resource: json['resource'],
      action: json['action'],
      description: json['description'],
      scopeId: json['scope_id'],
      businessTypeId: json['business_type_id'],
      businessTypeName: json['business_type_name'],
    );
  }
}

class GetPermissionsParams {
  final int? page;
  final int? pageSize;
  final int? businessTypeId;
  final String? name;
  final int? scopeId;
  final String? resource;

  GetPermissionsParams({
    this.page,
    this.pageSize,
    this.businessTypeId,
    this.name,
    this.scopeId,
    this.resource,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (businessTypeId != null) params['business_type_id'] = businessTypeId;
    if (name != null && name!.isNotEmpty) params['name'] = name;
    if (scopeId != null) params['scope_id'] = scopeId;
    if (resource != null && resource!.isNotEmpty) params['resource'] = resource;
    return params;
  }
}

class CreatePermissionDTO {
  final String resource;
  final String action;
  final String? description;
  final int? scopeId;
  final int? businessTypeId;

  CreatePermissionDTO({
    required this.resource,
    required this.action,
    this.description,
    this.scopeId,
    this.businessTypeId,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'resource': resource,
      'action': action,
    };
    if (description != null) json['description'] = description;
    if (scopeId != null) json['scope_id'] = scopeId;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}

class UpdatePermissionDTO {
  final String? resource;
  final String? action;
  final String? description;
  final int? scopeId;
  final int? businessTypeId;

  UpdatePermissionDTO({
    this.resource,
    this.action,
    this.description,
    this.scopeId,
    this.businessTypeId,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (resource != null) json['resource'] = resource;
    if (action != null) json['action'] = action;
    if (description != null) json['description'] = description;
    if (scopeId != null) json['scope_id'] = scopeId;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}

class BulkCreatePermissionsDTO {
  final List<CreatePermissionDTO> permissions;

  BulkCreatePermissionsDTO({required this.permissions});

  Map<String, dynamic> toJson() => {
        'permissions': permissions.map((p) => p.toJson()).toList(),
      };
}
