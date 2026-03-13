// ============================================
// ENTITIES
// ============================================

class Warehouse {
  final int id;
  final int businessId;
  final String name;
  final String code;
  final String address;
  final String city;
  final String state;
  final String country;
  final String zipCode;
  final String phone;
  final String contactName;
  final String contactEmail;
  final bool isActive;
  final bool isDefault;
  final bool isFulfillment;
  final String company;
  final String firstName;
  final String lastName;
  final String email;
  final String suburb;
  final String cityDaneCode;
  final String postalCode;
  final String street;
  final double? latitude;
  final double? longitude;
  final String createdAt;
  final String updatedAt;

  Warehouse({
    required this.id,
    required this.businessId,
    required this.name,
    required this.code,
    required this.address,
    required this.city,
    required this.state,
    required this.country,
    required this.zipCode,
    required this.phone,
    required this.contactName,
    required this.contactEmail,
    required this.isActive,
    required this.isDefault,
    required this.isFulfillment,
    required this.company,
    required this.firstName,
    required this.lastName,
    required this.email,
    required this.suburb,
    required this.cityDaneCode,
    required this.postalCode,
    required this.street,
    this.latitude,
    this.longitude,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Warehouse.fromJson(Map<String, dynamic> json) {
    return Warehouse(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      address: json['address'] ?? '',
      city: json['city'] ?? '',
      state: json['state'] ?? '',
      country: json['country'] ?? '',
      zipCode: json['zip_code'] ?? '',
      phone: json['phone'] ?? '',
      contactName: json['contact_name'] ?? '',
      contactEmail: json['contact_email'] ?? '',
      isActive: json['is_active'] ?? false,
      isDefault: json['is_default'] ?? false,
      isFulfillment: json['is_fulfillment'] ?? false,
      company: json['company'] ?? '',
      firstName: json['first_name'] ?? '',
      lastName: json['last_name'] ?? '',
      email: json['email'] ?? '',
      suburb: json['suburb'] ?? '',
      cityDaneCode: json['city_dane_code'] ?? '',
      postalCode: json['postal_code'] ?? '',
      street: json['street'] ?? '',
      latitude: json['latitude']?.toDouble(),
      longitude: json['longitude']?.toDouble(),
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class WarehouseDetail extends Warehouse {
  final List<WarehouseLocation> locations;

  WarehouseDetail({
    required super.id,
    required super.businessId,
    required super.name,
    required super.code,
    required super.address,
    required super.city,
    required super.state,
    required super.country,
    required super.zipCode,
    required super.phone,
    required super.contactName,
    required super.contactEmail,
    required super.isActive,
    required super.isDefault,
    required super.isFulfillment,
    required super.company,
    required super.firstName,
    required super.lastName,
    required super.email,
    required super.suburb,
    required super.cityDaneCode,
    required super.postalCode,
    required super.street,
    super.latitude,
    super.longitude,
    required super.createdAt,
    required super.updatedAt,
    required this.locations,
  });

  factory WarehouseDetail.fromJson(Map<String, dynamic> json) {
    return WarehouseDetail(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      address: json['address'] ?? '',
      city: json['city'] ?? '',
      state: json['state'] ?? '',
      country: json['country'] ?? '',
      zipCode: json['zip_code'] ?? '',
      phone: json['phone'] ?? '',
      contactName: json['contact_name'] ?? '',
      contactEmail: json['contact_email'] ?? '',
      isActive: json['is_active'] ?? false,
      isDefault: json['is_default'] ?? false,
      isFulfillment: json['is_fulfillment'] ?? false,
      company: json['company'] ?? '',
      firstName: json['first_name'] ?? '',
      lastName: json['last_name'] ?? '',
      email: json['email'] ?? '',
      suburb: json['suburb'] ?? '',
      cityDaneCode: json['city_dane_code'] ?? '',
      postalCode: json['postal_code'] ?? '',
      street: json['street'] ?? '',
      latitude: json['latitude']?.toDouble(),
      longitude: json['longitude']?.toDouble(),
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      locations: (json['locations'] as List<dynamic>?)
              ?.map((e) => WarehouseLocation.fromJson(e))
              .toList() ??
          [],
    );
  }
}

class WarehouseLocation {
  final int id;
  final int warehouseId;
  final String name;
  final String code;
  final String type;
  final bool isActive;
  final bool isFulfillment;
  final int? capacity;
  final String createdAt;
  final String updatedAt;

  WarehouseLocation({
    required this.id,
    required this.warehouseId,
    required this.name,
    required this.code,
    required this.type,
    required this.isActive,
    required this.isFulfillment,
    this.capacity,
    required this.createdAt,
    required this.updatedAt,
  });

  factory WarehouseLocation.fromJson(Map<String, dynamic> json) {
    return WarehouseLocation(
      id: json['id'] ?? 0,
      warehouseId: json['warehouse_id'] ?? 0,
      name: json['name'] ?? '',
      code: json['code'] ?? '',
      type: json['type'] ?? '',
      isActive: json['is_active'] ?? false,
      isFulfillment: json['is_fulfillment'] ?? false,
      capacity: json['capacity'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

// ============================================
// DTOs
// ============================================

class CreateWarehouseDTO {
  final String name;
  final String code;
  final String? address;
  final String? city;
  final String? state;
  final String? country;
  final String? zipCode;
  final String? phone;
  final String? contactName;
  final String? contactEmail;
  final bool? isDefault;
  final bool? isFulfillment;
  final String? company;
  final String? firstName;
  final String? lastName;
  final String? email;
  final String? suburb;
  final String? cityDaneCode;
  final String? postalCode;
  final String? street;
  final double? latitude;
  final double? longitude;

  CreateWarehouseDTO({
    required this.name,
    required this.code,
    this.address,
    this.city,
    this.state,
    this.country,
    this.zipCode,
    this.phone,
    this.contactName,
    this.contactEmail,
    this.isDefault,
    this.isFulfillment,
    this.company,
    this.firstName,
    this.lastName,
    this.email,
    this.suburb,
    this.cityDaneCode,
    this.postalCode,
    this.street,
    this.latitude,
    this.longitude,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'code': code,
    };
    if (address != null) json['address'] = address;
    if (city != null) json['city'] = city;
    if (state != null) json['state'] = state;
    if (country != null) json['country'] = country;
    if (zipCode != null) json['zip_code'] = zipCode;
    if (phone != null) json['phone'] = phone;
    if (contactName != null) json['contact_name'] = contactName;
    if (contactEmail != null) json['contact_email'] = contactEmail;
    if (isDefault != null) json['is_default'] = isDefault;
    if (isFulfillment != null) json['is_fulfillment'] = isFulfillment;
    if (company != null) json['company'] = company;
    if (firstName != null) json['first_name'] = firstName;
    if (lastName != null) json['last_name'] = lastName;
    if (email != null) json['email'] = email;
    if (suburb != null) json['suburb'] = suburb;
    if (cityDaneCode != null) json['city_dane_code'] = cityDaneCode;
    if (postalCode != null) json['postal_code'] = postalCode;
    if (street != null) json['street'] = street;
    if (latitude != null) json['latitude'] = latitude;
    if (longitude != null) json['longitude'] = longitude;
    return json;
  }
}

class UpdateWarehouseDTO {
  final String name;
  final String code;
  final String? address;
  final String? city;
  final String? state;
  final String? country;
  final String? zipCode;
  final String? phone;
  final String? contactName;
  final String? contactEmail;
  final bool? isActive;
  final bool? isDefault;
  final bool? isFulfillment;
  final String? company;
  final String? firstName;
  final String? lastName;
  final String? email;
  final String? suburb;
  final String? cityDaneCode;
  final String? postalCode;
  final String? street;
  final double? latitude;
  final double? longitude;

  UpdateWarehouseDTO({
    required this.name,
    required this.code,
    this.address,
    this.city,
    this.state,
    this.country,
    this.zipCode,
    this.phone,
    this.contactName,
    this.contactEmail,
    this.isActive,
    this.isDefault,
    this.isFulfillment,
    this.company,
    this.firstName,
    this.lastName,
    this.email,
    this.suburb,
    this.cityDaneCode,
    this.postalCode,
    this.street,
    this.latitude,
    this.longitude,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'code': code,
    };
    if (address != null) json['address'] = address;
    if (city != null) json['city'] = city;
    if (state != null) json['state'] = state;
    if (country != null) json['country'] = country;
    if (zipCode != null) json['zip_code'] = zipCode;
    if (phone != null) json['phone'] = phone;
    if (contactName != null) json['contact_name'] = contactName;
    if (contactEmail != null) json['contact_email'] = contactEmail;
    if (isActive != null) json['is_active'] = isActive;
    if (isDefault != null) json['is_default'] = isDefault;
    if (isFulfillment != null) json['is_fulfillment'] = isFulfillment;
    if (company != null) json['company'] = company;
    if (firstName != null) json['first_name'] = firstName;
    if (lastName != null) json['last_name'] = lastName;
    if (email != null) json['email'] = email;
    if (suburb != null) json['suburb'] = suburb;
    if (cityDaneCode != null) json['city_dane_code'] = cityDaneCode;
    if (postalCode != null) json['postal_code'] = postalCode;
    if (street != null) json['street'] = street;
    if (latitude != null) json['latitude'] = latitude;
    if (longitude != null) json['longitude'] = longitude;
    return json;
  }
}

class GetWarehousesParams {
  final int? page;
  final int? pageSize;
  final String? search;
  final bool? isActive;
  final bool? isFulfillment;
  final int? businessId;

  GetWarehousesParams({
    this.page,
    this.pageSize,
    this.search,
    this.isActive,
    this.isFulfillment,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (search != null) params['search'] = search;
    if (isActive != null) params['is_active'] = isActive;
    if (isFulfillment != null) params['is_fulfillment'] = isFulfillment;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

class CreateLocationDTO {
  final String name;
  final String code;
  final String? type;
  final bool? isFulfillment;
  final int? capacity;

  CreateLocationDTO({
    required this.name,
    required this.code,
    this.type,
    this.isFulfillment,
    this.capacity,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'code': code,
    };
    if (type != null) json['type'] = type;
    if (isFulfillment != null) json['is_fulfillment'] = isFulfillment;
    if (capacity != null) json['capacity'] = capacity;
    return json;
  }
}

class UpdateLocationDTO {
  final String name;
  final String code;
  final String? type;
  final bool? isActive;
  final bool? isFulfillment;
  final int? capacity;

  UpdateLocationDTO({
    required this.name,
    required this.code,
    this.type,
    this.isActive,
    this.isFulfillment,
    this.capacity,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'name': name,
      'code': code,
    };
    if (type != null) json['type'] = type;
    if (isActive != null) json['is_active'] = isActive;
    if (isFulfillment != null) json['is_fulfillment'] = isFulfillment;
    if (capacity != null) json['capacity'] = capacity;
    return json;
  }
}
