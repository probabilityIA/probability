class UserInfo {
  final int id;
  final String name;
  final String email;
  final String? phone;
  final String? avatarUrl;
  final bool isActive;
  final String? lastLoginAt;

  UserInfo({
    required this.id,
    required this.name,
    required this.email,
    this.phone,
    this.avatarUrl,
    required this.isActive,
    this.lastLoginAt,
  });

  factory UserInfo.fromJson(Map<String, dynamic> json) {
    return UserInfo(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      email: json['email'] ?? '',
      phone: json['phone'],
      avatarUrl: json['avatar_url'],
      isActive: json['is_active'] ?? false,
      lastLoginAt: json['last_login_at'],
    );
  }
}

class BusinessInfo {
  final int id;
  final String name;
  final String? logoUrl;
  final String? primaryColor;
  final String? secondaryColor;
  final String? accentColor;

  BusinessInfo({
    required this.id,
    required this.name,
    this.logoUrl,
    this.primaryColor,
    this.secondaryColor,
    this.accentColor,
  });

  factory BusinessInfo.fromJson(Map<String, dynamic> json) {
    return BusinessInfo(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      logoUrl: json['logo_url'],
      primaryColor: json['primary_color'],
      secondaryColor: json['secondary_color'],
      accentColor: json['accent_color'],
    );
  }
}

class LoginResponse {
  final UserInfo user;
  final String token;
  final bool requirePasswordChange;
  final List<BusinessInfo> businesses;
  final String? scope;
  final bool isSuperAdmin;

  LoginResponse({
    required this.user,
    required this.token,
    required this.requirePasswordChange,
    required this.businesses,
    this.scope,
    required this.isSuperAdmin,
  });

  factory LoginResponse.fromJson(Map<String, dynamic> json) {
    return LoginResponse(
      user: UserInfo.fromJson(json['user'] ?? {}),
      token: json['token'] ?? '',
      requirePasswordChange: json['require_password_change'] ?? false,
      businesses: (json['businesses'] as List<dynamic>?)
              ?.map((b) => BusinessInfo.fromJson(b))
              .toList() ??
          [],
      scope: json['scope'],
      isSuperAdmin: json['is_super_admin'] ?? false,
    );
  }
}

class LoginSuccessResponse {
  final bool success;
  final LoginResponse data;

  LoginSuccessResponse({required this.success, required this.data});

  factory LoginSuccessResponse.fromJson(Map<String, dynamic> json) {
    return LoginSuccessResponse(
      success: json['success'] ?? false,
      data: LoginResponse.fromJson(json['data'] ?? json),
    );
  }
}

class UserRolesPermissionsResponse {
  final bool isSuper;
  final int businessId;
  final String? businessName;
  final int? businessTypeId;
  final String? businessTypeName;
  final String? role;
  final List<ResourcePermission> resources;
  final String? subscriptionStatus;

  UserRolesPermissionsResponse({
    required this.isSuper,
    required this.businessId,
    this.businessName,
    this.businessTypeId,
    this.businessTypeName,
    this.role,
    required this.resources,
    this.subscriptionStatus,
  });

  factory UserRolesPermissionsResponse.fromJson(Map<String, dynamic> json) {
    return UserRolesPermissionsResponse(
      isSuper: json['is_super'] ?? false,
      businessId: json['business_id'] ?? 0,
      businessName: json['business_name'],
      businessTypeId: json['business_type_id'],
      businessTypeName: json['business_type_name'],
      role: json['role'],
      resources: (json['resources'] as List<dynamic>?)
              ?.map((r) => ResourcePermission.fromJson(r))
              .toList() ??
          [],
      subscriptionStatus: json['subscription_status'],
    );
  }
}

class ResourcePermission {
  final String resource;
  final List<String> actions;

  ResourcePermission({required this.resource, required this.actions});

  factory ResourcePermission.fromJson(Map<String, dynamic> json) {
    return ResourcePermission(
      resource: json['resource'] ?? '',
      actions:
          (json['actions'] as List<dynamic>?)?.cast<String>() ?? [],
    );
  }
}

class ChangePasswordResponse {
  final bool success;
  final String message;

  ChangePasswordResponse({required this.success, required this.message});

  factory ChangePasswordResponse.fromJson(Map<String, dynamic> json) {
    return ChangePasswordResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
    );
  }
}

class GeneratePasswordResponse {
  final bool success;
  final String message;
  final String? password;

  GeneratePasswordResponse({
    required this.success,
    required this.message,
    this.password,
  });

  factory GeneratePasswordResponse.fromJson(Map<String, dynamic> json) {
    return GeneratePasswordResponse(
      success: json['success'] ?? false,
      message: json['message'] ?? '',
      password: json['password'],
    );
  }
}
