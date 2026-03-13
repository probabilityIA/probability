class BusinessRoleAssignment {
  final int businessId;
  final String? businessName;
  final int roleId;
  final String? roleName;

  BusinessRoleAssignment({
    required this.businessId,
    this.businessName,
    required this.roleId,
    this.roleName,
  });

  factory BusinessRoleAssignment.fromJson(Map<String, dynamic> json) {
    return BusinessRoleAssignment(
      businessId: json['business_id'] ?? 0,
      businessName: json['business_name'],
      roleId: json['role_id'] ?? 0,
      roleName: json['role_name'],
    );
  }
}

class User {
  final int id;
  final String name;
  final String email;
  final String? phone;
  final String? avatarUrl;
  final bool isActive;
  final bool isSuperUser;
  final int? scopeId;
  final List<BusinessRoleAssignment> businessRoleAssignments;

  User({
    required this.id,
    required this.name,
    required this.email,
    this.phone,
    this.avatarUrl,
    required this.isActive,
    required this.isSuperUser,
    this.scopeId,
    required this.businessRoleAssignments,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      email: json['email'] ?? '',
      phone: json['phone'],
      avatarUrl: json['avatar_url'],
      isActive: json['is_active'] ?? false,
      isSuperUser: json['is_super_user'] ?? false,
      scopeId: json['scope_id'],
      businessRoleAssignments: (json['business_role_assignments']
                  as List<dynamic>?)
              ?.map((a) => BusinessRoleAssignment.fromJson(a))
              .toList() ??
          [],
    );
  }
}

class GetUsersParams {
  final int? page;
  final int? pageSize;
  final String? name;
  final String? email;
  final String? phone;
  final bool? isActive;
  final int? roleId;
  final int? businessId;

  GetUsersParams({
    this.page,
    this.pageSize,
    this.name,
    this.email,
    this.phone,
    this.isActive,
    this.roleId,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (name != null && name!.isNotEmpty) params['name'] = name;
    if (email != null && email!.isNotEmpty) params['email'] = email;
    if (phone != null && phone!.isNotEmpty) params['phone'] = phone;
    if (isActive != null) params['is_active'] = isActive;
    if (roleId != null) params['role_id'] = roleId;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

class CreateUserDTO {
  final String name;
  final String email;
  final String? phone;
  final bool isActive;
  final int? scopeId;
  final List<int>? businessIds;

  CreateUserDTO({
    required this.name,
    required this.email,
    this.phone,
    this.isActive = true,
    this.scopeId,
    this.businessIds,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'email': email,
      'is_active': isActive,
    };
    if (phone != null) json['phone'] = phone;
    if (scopeId != null) json['scope_id'] = scopeId;
    if (businessIds != null) json['business_ids'] = businessIds;
    return json;
  }
}

class UpdateUserDTO {
  final String? name;
  final String? email;
  final String? phone;
  final bool? isActive;
  final int? scopeId;
  final bool? removeAvatar;

  UpdateUserDTO({
    this.name,
    this.email,
    this.phone,
    this.isActive,
    this.scopeId,
    this.removeAvatar,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (name != null) json['name'] = name;
    if (email != null) json['email'] = email;
    if (phone != null) json['phone'] = phone;
    if (isActive != null) json['is_active'] = isActive;
    if (scopeId != null) json['scope_id'] = scopeId;
    if (removeAvatar != null) json['remove_avatar'] = removeAvatar;
    return json;
  }
}

class RoleAssignment {
  final int businessId;
  final int roleId;

  RoleAssignment({required this.businessId, required this.roleId});

  Map<String, dynamic> toJson() => {
        'business_id': businessId,
        'role_id': roleId,
      };
}

class AssignRolesDTO {
  final List<RoleAssignment> assignments;

  AssignRolesDTO({required this.assignments});

  Map<String, dynamic> toJson() => {
        'assignments': assignments.map((a) => a.toJson()).toList(),
      };
}
