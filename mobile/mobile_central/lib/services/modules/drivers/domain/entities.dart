// ============================================
// ENTITIES
// ============================================

class DriverInfo {
  final int id;
  final int businessId;
  final String firstName;
  final String lastName;
  final String email;
  final String phone;
  final String identification;
  final String status;
  final String photoUrl;
  final String licenseType;
  final String? licenseExpiry;
  final int? warehouseId;
  final String? notes;
  final String createdAt;
  final String updatedAt;

  DriverInfo({
    required this.id,
    required this.businessId,
    required this.firstName,
    required this.lastName,
    required this.email,
    required this.phone,
    required this.identification,
    required this.status,
    required this.photoUrl,
    required this.licenseType,
    this.licenseExpiry,
    this.warehouseId,
    this.notes,
    required this.createdAt,
    required this.updatedAt,
  });

  factory DriverInfo.fromJson(Map<String, dynamic> json) {
    return DriverInfo(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      firstName: json['first_name'] ?? '',
      lastName: json['last_name'] ?? '',
      email: json['email'] ?? '',
      phone: json['phone'] ?? '',
      identification: json['identification'] ?? '',
      status: json['status'] ?? '',
      photoUrl: json['photo_url'] ?? '',
      licenseType: json['license_type'] ?? '',
      licenseExpiry: json['license_expiry'],
      warehouseId: json['warehouse_id'],
      notes: json['notes'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

// ============================================
// DTOs
// ============================================

class CreateDriverDTO {
  final String firstName;
  final String lastName;
  final String? email;
  final String phone;
  final String identification;
  final String? licenseType;
  final String? licenseExpiry;
  final int? warehouseId;
  final String? notes;

  CreateDriverDTO({
    required this.firstName,
    required this.lastName,
    this.email,
    required this.phone,
    required this.identification,
    this.licenseType,
    this.licenseExpiry,
    this.warehouseId,
    this.notes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'first_name': firstName,
      'last_name': lastName,
      'phone': phone,
      'identification': identification,
    };
    if (email != null) json['email'] = email;
    if (licenseType != null) json['license_type'] = licenseType;
    if (licenseExpiry != null) json['license_expiry'] = licenseExpiry;
    if (warehouseId != null) json['warehouse_id'] = warehouseId;
    if (notes != null) json['notes'] = notes;
    return json;
  }
}

class UpdateDriverDTO {
  final String? firstName;
  final String? lastName;
  final String? email;
  final String? phone;
  final String? identification;
  final String? status;
  final String? licenseType;
  final String? licenseExpiry;
  final int? warehouseId;
  final String? notes;

  UpdateDriverDTO({
    this.firstName,
    this.lastName,
    this.email,
    this.phone,
    this.identification,
    this.status,
    this.licenseType,
    this.licenseExpiry,
    this.warehouseId,
    this.notes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (firstName != null) json['first_name'] = firstName;
    if (lastName != null) json['last_name'] = lastName;
    if (email != null) json['email'] = email;
    if (phone != null) json['phone'] = phone;
    if (identification != null) json['identification'] = identification;
    if (status != null) json['status'] = status;
    if (licenseType != null) json['license_type'] = licenseType;
    if (licenseExpiry != null) json['license_expiry'] = licenseExpiry;
    if (warehouseId != null) json['warehouse_id'] = warehouseId;
    if (notes != null) json['notes'] = notes;
    return json;
  }
}

class GetDriversParams {
  final int? page;
  final int? pageSize;
  final String? search;
  final String? status;
  final int? businessId;

  GetDriversParams({
    this.page,
    this.pageSize,
    this.search,
    this.status,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (search != null) params['search'] = search;
    if (status != null) params['status'] = status;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}
