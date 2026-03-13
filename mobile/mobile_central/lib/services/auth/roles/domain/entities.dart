class Role {
  final int id;
  final String name;
  final String? code;
  final String? description;
  final int? level;
  final bool isSystem;
  final int? scopeId;
  final int? businessTypeId;

  Role({
    required this.id,
    required this.name,
    this.code,
    this.description,
    this.level,
    required this.isSystem,
    this.scopeId,
    this.businessTypeId,
  });

  factory Role.fromJson(Map<String, dynamic> json) {
    return Role(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'],
      description: json['description'],
      level: json['level'],
      isSystem: json['is_system'] ?? false,
      scopeId: json['scope_id'],
      businessTypeId: json['business_type_id'],
    );
  }
}

class RolePermission {
  final int id;
  final String? resource;
  final String? action;
  final String? description;
  final int? scopeId;

  RolePermission({
    required this.id,
    this.resource,
    this.action,
    this.description,
    this.scopeId,
  });

  factory RolePermission.fromJson(Map<String, dynamic> json) {
    return RolePermission(
      id: json['id'] ?? 0,
      resource: json['resource'],
      action: json['action'],
      description: json['description'],
      scopeId: json['scope_id'],
    );
  }
}

class RolePermissionsResponse {
  final int roleId;
  final String roleName;
  final List<RolePermission> permissions;
  final int count;

  RolePermissionsResponse({
    required this.roleId,
    required this.roleName,
    required this.permissions,
    required this.count,
  });

  factory RolePermissionsResponse.fromJson(Map<String, dynamic> json) {
    return RolePermissionsResponse(
      roleId: json['role_id'] ?? 0,
      roleName: json['role_name'] ?? '',
      permissions: (json['permissions'] as List<dynamic>?)
              ?.map((p) => RolePermission.fromJson(p))
              .toList() ??
          [],
      count: json['count'] ?? 0,
    );
  }
}

class GetRolesParams {
  final int? page;
  final int? pageSize;
  final int? businessTypeId;
  final int? scopeId;
  final bool? isSystem;
  final String? name;
  final int? level;

  GetRolesParams({
    this.page,
    this.pageSize,
    this.businessTypeId,
    this.scopeId,
    this.isSystem,
    this.name,
    this.level,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (businessTypeId != null) params['business_type_id'] = businessTypeId;
    if (scopeId != null) params['scope_id'] = scopeId;
    if (isSystem != null) params['is_system'] = isSystem;
    if (name != null && name!.isNotEmpty) params['name'] = name;
    if (level != null) params['level'] = level;
    return params;
  }
}

class CreateRoleDTO {
  final String name;
  final String? code;
  final String? description;
  final int? level;
  final int? scopeId;
  final int? businessTypeId;

  CreateRoleDTO({
    required this.name,
    this.code,
    this.description,
    this.level,
    this.scopeId,
    this.businessTypeId,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'name': name};
    if (code != null) json['code'] = code;
    if (description != null) json['description'] = description;
    if (level != null) json['level'] = level;
    if (scopeId != null) json['scope_id'] = scopeId;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}

class UpdateRoleDTO {
  final String? name;
  final String? code;
  final String? description;
  final int? level;
  final int? scopeId;
  final int? businessTypeId;

  UpdateRoleDTO({
    this.name,
    this.code,
    this.description,
    this.level,
    this.scopeId,
    this.businessTypeId,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (name != null) json['name'] = name;
    if (code != null) json['code'] = code;
    if (description != null) json['description'] = description;
    if (level != null) json['level'] = level;
    if (scopeId != null) json['scope_id'] = scopeId;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}

class AssignPermissionsDTO {
  final List<int> permissionIds;

  AssignPermissionsDTO({required this.permissionIds});

  Map<String, dynamic> toJson() => {'permission_ids': permissionIds};
}
