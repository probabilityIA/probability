class Business {
  final int id;
  final String name;
  final String? logoUrl;
  final String? primaryColor;
  final String? secondaryColor;
  final String? accentColor;
  final String? navbarColor;
  final String? navbarImageUrl;
  final String? domain;
  final bool? hasDelivery;
  final bool? hasPickup;
  final int? businessTypeId;
  final String? businessTypeName;
  final bool isActive;

  Business({
    required this.id,
    required this.name,
    this.logoUrl,
    this.primaryColor,
    this.secondaryColor,
    this.accentColor,
    this.navbarColor,
    this.navbarImageUrl,
    this.domain,
    this.hasDelivery,
    this.hasPickup,
    this.businessTypeId,
    this.businessTypeName,
    required this.isActive,
  });

  factory Business.fromJson(Map<String, dynamic> json) {
    return Business(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      logoUrl: json['logo_url'],
      primaryColor: json['primary_color'],
      secondaryColor: json['secondary_color'],
      accentColor: json['accent_color'],
      navbarColor: json['navbar_color'],
      navbarImageUrl: json['navbar_image_url'],
      domain: json['domain'],
      hasDelivery: json['has_delivery'],
      hasPickup: json['has_pickup'],
      businessTypeId: json['business_type_id'],
      businessTypeName: json['business_type_name'],
      isActive: json['is_active'] ?? true,
    );
  }
}

class BusinessSimple {
  final int id;
  final String name;
  final String? logoUrl;
  final String? primaryColor;
  final String? secondaryColor;

  BusinessSimple({
    required this.id,
    required this.name,
    this.logoUrl,
    this.primaryColor,
    this.secondaryColor,
  });

  factory BusinessSimple.fromJson(Map<String, dynamic> json) {
    return BusinessSimple(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      logoUrl: json['logo_url'],
      primaryColor: json['primary_color'],
      secondaryColor: json['secondary_color'],
    );
  }
}

class BusinessType {
  final int id;
  final String name;
  final String? code;
  final String? description;
  final String? icon;

  BusinessType({
    required this.id,
    required this.name,
    this.code,
    this.description,
    this.icon,
  });

  factory BusinessType.fromJson(Map<String, dynamic> json) {
    return BusinessType(
      id: json['id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'],
      description: json['description'],
      icon: json['icon'],
    );
  }
}

class ConfiguredResource {
  final int id;
  final String? name;
  final bool isActive;

  ConfiguredResource({
    required this.id,
    this.name,
    required this.isActive,
  });

  factory ConfiguredResource.fromJson(Map<String, dynamic> json) {
    return ConfiguredResource(
      id: json['id'] ?? 0,
      name: json['name'],
      isActive: json['is_active'] ?? false,
    );
  }
}

class GetBusinessesParams {
  final int? page;
  final int? pageSize;
  final String? name;
  final int? businessTypeId;

  GetBusinessesParams({this.page, this.pageSize, this.name, this.businessTypeId});

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (name != null && name!.isNotEmpty) params['name'] = name;
    if (businessTypeId != null) params['business_type_id'] = businessTypeId;
    return params;
  }
}

class CreateBusinessDTO {
  final String name;
  final String? primaryColor;
  final String? secondaryColor;
  final String? accentColor;
  final String? navbarColor;
  final String? domain;
  final bool? hasDelivery;
  final bool? hasPickup;
  final int? businessTypeId;

  CreateBusinessDTO({
    required this.name,
    this.primaryColor,
    this.secondaryColor,
    this.accentColor,
    this.navbarColor,
    this.domain,
    this.hasDelivery,
    this.hasPickup,
    this.businessTypeId,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'name': name};
    if (primaryColor != null) json['primary_color'] = primaryColor;
    if (secondaryColor != null) json['secondary_color'] = secondaryColor;
    if (accentColor != null) json['accent_color'] = accentColor;
    if (navbarColor != null) json['navbar_color'] = navbarColor;
    if (domain != null) json['domain'] = domain;
    if (hasDelivery != null) json['has_delivery'] = hasDelivery;
    if (hasPickup != null) json['has_pickup'] = hasPickup;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}

class UpdateBusinessDTO {
  final String? name;
  final String? primaryColor;
  final String? secondaryColor;
  final String? accentColor;
  final String? navbarColor;
  final String? domain;
  final bool? hasDelivery;
  final bool? hasPickup;
  final int? businessTypeId;

  UpdateBusinessDTO({
    this.name,
    this.primaryColor,
    this.secondaryColor,
    this.accentColor,
    this.navbarColor,
    this.domain,
    this.hasDelivery,
    this.hasPickup,
    this.businessTypeId,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (name != null) json['name'] = name;
    if (primaryColor != null) json['primary_color'] = primaryColor;
    if (secondaryColor != null) json['secondary_color'] = secondaryColor;
    if (accentColor != null) json['accent_color'] = accentColor;
    if (navbarColor != null) json['navbar_color'] = navbarColor;
    if (domain != null) json['domain'] = domain;
    if (hasDelivery != null) json['has_delivery'] = hasDelivery;
    if (hasPickup != null) json['has_pickup'] = hasPickup;
    if (businessTypeId != null) json['business_type_id'] = businessTypeId;
    return json;
  }
}
