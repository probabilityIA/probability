// ============================================
// ENTITIES
// ============================================

class VehicleInfo {
  final int id;
  final int businessId;
  final String type;
  final String licensePlate;
  final String brand;
  final String model;
  final int? year;
  final String color;
  final String status;
  final double? weightCapacityKg;
  final double? volumeCapacityM3;
  final String photoUrl;
  final String? insuranceExpiry;
  final String? registrationExpiry;
  final String createdAt;
  final String updatedAt;

  VehicleInfo({
    required this.id,
    required this.businessId,
    required this.type,
    required this.licensePlate,
    required this.brand,
    required this.model,
    this.year,
    required this.color,
    required this.status,
    this.weightCapacityKg,
    this.volumeCapacityM3,
    required this.photoUrl,
    this.insuranceExpiry,
    this.registrationExpiry,
    required this.createdAt,
    required this.updatedAt,
  });

  factory VehicleInfo.fromJson(Map<String, dynamic> json) {
    return VehicleInfo(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      type: json['type'] ?? '',
      licensePlate: json['license_plate'] ?? '',
      brand: json['brand'] ?? '',
      model: json['model'] ?? '',
      year: json['year'],
      color: json['color'] ?? '',
      status: json['status'] ?? '',
      weightCapacityKg: json['weight_capacity_kg']?.toDouble(),
      volumeCapacityM3: json['volume_capacity_m3']?.toDouble(),
      photoUrl: json['photo_url'] ?? '',
      insuranceExpiry: json['insurance_expiry'],
      registrationExpiry: json['registration_expiry'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

// ============================================
// DTOs
// ============================================

class CreateVehicleDTO {
  final String type;
  final String licensePlate;
  final String? brand;
  final String? model;
  final int? year;
  final String? color;
  final double? weightCapacityKg;
  final double? volumeCapacityM3;
  final String? insuranceExpiry;
  final String? registrationExpiry;

  CreateVehicleDTO({
    required this.type,
    required this.licensePlate,
    this.brand,
    this.model,
    this.year,
    this.color,
    this.weightCapacityKg,
    this.volumeCapacityM3,
    this.insuranceExpiry,
    this.registrationExpiry,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'type': type,
      'license_plate': licensePlate,
    };
    if (brand != null) json['brand'] = brand;
    if (model != null) json['model'] = model;
    if (year != null) json['year'] = year;
    if (color != null) json['color'] = color;
    if (weightCapacityKg != null) json['weight_capacity_kg'] = weightCapacityKg;
    if (volumeCapacityM3 != null) json['volume_capacity_m3'] = volumeCapacityM3;
    if (insuranceExpiry != null) json['insurance_expiry'] = insuranceExpiry;
    if (registrationExpiry != null) json['registration_expiry'] = registrationExpiry;
    return json;
  }
}

class UpdateVehicleDTO {
  final String? type;
  final String? licensePlate;
  final String? brand;
  final String? model;
  final int? year;
  final String? color;
  final String? status;
  final double? weightCapacityKg;
  final double? volumeCapacityM3;
  final String? insuranceExpiry;
  final String? registrationExpiry;

  UpdateVehicleDTO({
    this.type,
    this.licensePlate,
    this.brand,
    this.model,
    this.year,
    this.color,
    this.status,
    this.weightCapacityKg,
    this.volumeCapacityM3,
    this.insuranceExpiry,
    this.registrationExpiry,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (type != null) json['type'] = type;
    if (licensePlate != null) json['license_plate'] = licensePlate;
    if (brand != null) json['brand'] = brand;
    if (model != null) json['model'] = model;
    if (year != null) json['year'] = year;
    if (color != null) json['color'] = color;
    if (status != null) json['status'] = status;
    if (weightCapacityKg != null) json['weight_capacity_kg'] = weightCapacityKg;
    if (volumeCapacityM3 != null) json['volume_capacity_m3'] = volumeCapacityM3;
    if (insuranceExpiry != null) json['insurance_expiry'] = insuranceExpiry;
    if (registrationExpiry != null) json['registration_expiry'] = registrationExpiry;
    return json;
  }
}

class GetVehiclesParams {
  final int? page;
  final int? pageSize;
  final String? search;
  final String? type;
  final String? status;
  final int? businessId;

  GetVehiclesParams({
    this.page,
    this.pageSize,
    this.search,
    this.type,
    this.status,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (search != null) params['search'] = search;
    if (type != null) params['type'] = type;
    if (status != null) params['status'] = status;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}
