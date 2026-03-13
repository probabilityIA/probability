// ============================================
// ENTITIES
// ============================================

class RouteInfo {
  final int id;
  final int businessId;
  final int? driverId;
  final String? driverName;
  final int? vehicleId;
  final String? vehiclePlate;
  final String status;
  final String date;
  final String? startTime;
  final String? endTime;
  final String? originAddress;
  final int totalStops;
  final int completedStops;
  final int failedStops;
  final String? notes;
  final String createdAt;
  final String updatedAt;

  RouteInfo({
    required this.id,
    required this.businessId,
    this.driverId,
    this.driverName,
    this.vehicleId,
    this.vehiclePlate,
    required this.status,
    required this.date,
    this.startTime,
    this.endTime,
    this.originAddress,
    required this.totalStops,
    required this.completedStops,
    required this.failedStops,
    this.notes,
    required this.createdAt,
    required this.updatedAt,
  });

  factory RouteInfo.fromJson(Map<String, dynamic> json) {
    return RouteInfo(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      driverId: json['driver_id'],
      driverName: json['driver_name'],
      vehicleId: json['vehicle_id'],
      vehiclePlate: json['vehicle_plate'],
      status: json['status'] ?? '',
      date: json['date'] ?? '',
      startTime: json['start_time'],
      endTime: json['end_time'],
      originAddress: json['origin_address'],
      totalStops: json['total_stops'] ?? 0,
      completedStops: json['completed_stops'] ?? 0,
      failedStops: json['failed_stops'] ?? 0,
      notes: json['notes'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class RouteStopInfo {
  final int id;
  final int routeId;
  final String? orderId;
  final int sequence;
  final String status;
  final String address;
  final String? city;
  final double? lat;
  final double? lng;
  final String customerName;
  final String? customerPhone;
  final String? estimatedArrival;
  final String? actualArrival;
  final String? actualDeparture;
  final String? deliveryNotes;
  final String? failureReason;
  final String createdAt;
  final String updatedAt;

  RouteStopInfo({
    required this.id,
    required this.routeId,
    this.orderId,
    required this.sequence,
    required this.status,
    required this.address,
    this.city,
    this.lat,
    this.lng,
    required this.customerName,
    this.customerPhone,
    this.estimatedArrival,
    this.actualArrival,
    this.actualDeparture,
    this.deliveryNotes,
    this.failureReason,
    required this.createdAt,
    required this.updatedAt,
  });

  factory RouteStopInfo.fromJson(Map<String, dynamic> json) {
    return RouteStopInfo(
      id: json['id'] ?? 0,
      routeId: json['route_id'] ?? 0,
      orderId: json['order_id'],
      sequence: json['sequence'] ?? 0,
      status: json['status'] ?? '',
      address: json['address'] ?? '',
      city: json['city'],
      lat: json['lat']?.toDouble(),
      lng: json['lng']?.toDouble(),
      customerName: json['customer_name'] ?? '',
      customerPhone: json['customer_phone'],
      estimatedArrival: json['estimated_arrival'],
      actualArrival: json['actual_arrival'],
      actualDeparture: json['actual_departure'],
      deliveryNotes: json['delivery_notes'],
      failureReason: json['failure_reason'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }
}

class RouteDetail extends RouteInfo {
  final String? actualStartTime;
  final String? actualEndTime;
  final int? originWarehouseId;
  final double? originLat;
  final double? originLng;
  final double? totalDistanceKm;
  final double? totalDurationMin;
  final List<RouteStopInfo> stops;

  RouteDetail({
    required super.id,
    required super.businessId,
    super.driverId,
    super.driverName,
    super.vehicleId,
    super.vehiclePlate,
    required super.status,
    required super.date,
    super.startTime,
    super.endTime,
    super.originAddress,
    required super.totalStops,
    required super.completedStops,
    required super.failedStops,
    super.notes,
    required super.createdAt,
    required super.updatedAt,
    this.actualStartTime,
    this.actualEndTime,
    this.originWarehouseId,
    this.originLat,
    this.originLng,
    this.totalDistanceKm,
    this.totalDurationMin,
    required this.stops,
  });

  factory RouteDetail.fromJson(Map<String, dynamic> json) {
    return RouteDetail(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      driverId: json['driver_id'],
      driverName: json['driver_name'],
      vehicleId: json['vehicle_id'],
      vehiclePlate: json['vehicle_plate'],
      status: json['status'] ?? '',
      date: json['date'] ?? '',
      startTime: json['start_time'],
      endTime: json['end_time'],
      originAddress: json['origin_address'],
      totalStops: json['total_stops'] ?? 0,
      completedStops: json['completed_stops'] ?? 0,
      failedStops: json['failed_stops'] ?? 0,
      notes: json['notes'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      actualStartTime: json['actual_start_time'],
      actualEndTime: json['actual_end_time'],
      originWarehouseId: json['origin_warehouse_id'],
      originLat: json['origin_lat']?.toDouble(),
      originLng: json['origin_lng']?.toDouble(),
      totalDistanceKm: json['total_distance_km']?.toDouble(),
      totalDurationMin: json['total_duration_min']?.toDouble(),
      stops: (json['stops'] as List<dynamic>?)
              ?.map((e) => RouteStopInfo.fromJson(e))
              .toList() ??
          [],
    );
  }
}

// ============================================
// DTOs
// ============================================

class CreateRouteStopDTO {
  final String? orderId;
  final String address;
  final String? city;
  final double? lat;
  final double? lng;
  final String customerName;
  final String? customerPhone;
  final String? deliveryNotes;

  CreateRouteStopDTO({
    this.orderId,
    required this.address,
    this.city,
    this.lat,
    this.lng,
    required this.customerName,
    this.customerPhone,
    this.deliveryNotes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'address': address,
      'customer_name': customerName,
    };
    if (orderId != null) json['order_id'] = orderId;
    if (city != null) json['city'] = city;
    if (lat != null) json['lat'] = lat;
    if (lng != null) json['lng'] = lng;
    if (customerPhone != null) json['customer_phone'] = customerPhone;
    if (deliveryNotes != null) json['delivery_notes'] = deliveryNotes;
    return json;
  }
}

class CreateRouteDTO {
  final String date;
  final int? driverId;
  final int? vehicleId;
  final String? originAddress;
  final double? originLat;
  final double? originLng;
  final String? notes;
  final List<CreateRouteStopDTO>? stops;

  CreateRouteDTO({
    required this.date,
    this.driverId,
    this.vehicleId,
    this.originAddress,
    this.originLat,
    this.originLng,
    this.notes,
    this.stops,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'date': date};
    if (driverId != null) json['driver_id'] = driverId;
    if (vehicleId != null) json['vehicle_id'] = vehicleId;
    if (originAddress != null) json['origin_address'] = originAddress;
    if (originLat != null) json['origin_lat'] = originLat;
    if (originLng != null) json['origin_lng'] = originLng;
    if (notes != null) json['notes'] = notes;
    if (stops != null) json['stops'] = stops!.map((s) => s.toJson()).toList();
    return json;
  }
}

class UpdateRouteDTO {
  final int? driverId;
  final int? vehicleId;
  final String? date;
  final String? originAddress;
  final double? originLat;
  final double? originLng;
  final String? notes;

  UpdateRouteDTO({
    this.driverId,
    this.vehicleId,
    this.date,
    this.originAddress,
    this.originLat,
    this.originLng,
    this.notes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (driverId != null) json['driver_id'] = driverId;
    if (vehicleId != null) json['vehicle_id'] = vehicleId;
    if (date != null) json['date'] = date;
    if (originAddress != null) json['origin_address'] = originAddress;
    if (originLat != null) json['origin_lat'] = originLat;
    if (originLng != null) json['origin_lng'] = originLng;
    if (notes != null) json['notes'] = notes;
    return json;
  }
}

class AddStopDTO {
  final String? orderId;
  final String address;
  final String? city;
  final double? lat;
  final double? lng;
  final String customerName;
  final String? customerPhone;
  final String? deliveryNotes;

  AddStopDTO({
    this.orderId,
    required this.address,
    this.city,
    this.lat,
    this.lng,
    required this.customerName,
    this.customerPhone,
    this.deliveryNotes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'address': address,
      'customer_name': customerName,
    };
    if (orderId != null) json['order_id'] = orderId;
    if (city != null) json['city'] = city;
    if (lat != null) json['lat'] = lat;
    if (lng != null) json['lng'] = lng;
    if (customerPhone != null) json['customer_phone'] = customerPhone;
    if (deliveryNotes != null) json['delivery_notes'] = deliveryNotes;
    return json;
  }
}

class UpdateStopDTO {
  final String? address;
  final String? city;
  final double? lat;
  final double? lng;
  final String? customerName;
  final String? customerPhone;
  final String? deliveryNotes;

  UpdateStopDTO({
    this.address,
    this.city,
    this.lat,
    this.lng,
    this.customerName,
    this.customerPhone,
    this.deliveryNotes,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (address != null) json['address'] = address;
    if (city != null) json['city'] = city;
    if (lat != null) json['lat'] = lat;
    if (lng != null) json['lng'] = lng;
    if (customerName != null) json['customer_name'] = customerName;
    if (customerPhone != null) json['customer_phone'] = customerPhone;
    if (deliveryNotes != null) json['delivery_notes'] = deliveryNotes;
    return json;
  }
}

class UpdateStopStatusDTO {
  final String status;
  final String? failureReason;

  UpdateStopStatusDTO({required this.status, this.failureReason});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'status': status};
    if (failureReason != null) json['failure_reason'] = failureReason;
    return json;
  }
}

class ReorderStopsDTO {
  final List<int> stopIds;

  ReorderStopsDTO({required this.stopIds});

  Map<String, dynamic> toJson() => {'stop_ids': stopIds};
}

class GetRoutesParams {
  final int? page;
  final int? pageSize;
  final String? status;
  final int? driverId;
  final String? dateFrom;
  final String? dateTo;
  final String? search;
  final int? businessId;

  GetRoutesParams({
    this.page,
    this.pageSize,
    this.status,
    this.driverId,
    this.dateFrom,
    this.dateTo,
    this.search,
    this.businessId,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (status != null) params['status'] = status;
    if (driverId != null) params['driver_id'] = driverId;
    if (dateFrom != null) params['date_from'] = dateFrom;
    if (dateTo != null) params['date_to'] = dateTo;
    if (search != null) params['search'] = search;
    if (businessId != null) params['business_id'] = businessId;
    return params;
  }
}

// ============================================
// FORM OPTIONS
// ============================================

class DriverOption {
  final int id;
  final String firstName;
  final String lastName;
  final String phone;
  final String identification;
  final String status;
  final String licenseType;

  DriverOption({
    required this.id,
    required this.firstName,
    required this.lastName,
    required this.phone,
    required this.identification,
    required this.status,
    required this.licenseType,
  });

  factory DriverOption.fromJson(Map<String, dynamic> json) {
    return DriverOption(
      id: json['id'] ?? 0,
      firstName: json['first_name'] ?? '',
      lastName: json['last_name'] ?? '',
      phone: json['phone'] ?? '',
      identification: json['identification'] ?? '',
      status: json['status'] ?? '',
      licenseType: json['license_type'] ?? '',
    );
  }
}

class VehicleOption {
  final int id;
  final String type;
  final String licensePlate;
  final String brand;
  final String vehicleModel;
  final String status;

  VehicleOption({
    required this.id,
    required this.type,
    required this.licensePlate,
    required this.brand,
    required this.vehicleModel,
    required this.status,
  });

  factory VehicleOption.fromJson(Map<String, dynamic> json) {
    return VehicleOption(
      id: json['id'] ?? 0,
      type: json['type'] ?? '',
      licensePlate: json['license_plate'] ?? '',
      brand: json['brand'] ?? '',
      vehicleModel: json['vehicle_model'] ?? '',
      status: json['status'] ?? '',
    );
  }
}

class AssignableOrder {
  final String id;
  final String orderNumber;
  final String customerName;
  final String customerPhone;
  final String address;
  final String city;
  final double? lat;
  final double? lng;
  final double totalAmount;
  final int itemCount;
  final String createdAt;

  AssignableOrder({
    required this.id,
    required this.orderNumber,
    required this.customerName,
    required this.customerPhone,
    required this.address,
    required this.city,
    this.lat,
    this.lng,
    required this.totalAmount,
    required this.itemCount,
    required this.createdAt,
  });

  factory AssignableOrder.fromJson(Map<String, dynamic> json) {
    return AssignableOrder(
      id: json['id']?.toString() ?? '',
      orderNumber: json['order_number'] ?? '',
      customerName: json['customer_name'] ?? '',
      customerPhone: json['customer_phone'] ?? '',
      address: json['address'] ?? '',
      city: json['city'] ?? '',
      lat: json['lat']?.toDouble(),
      lng: json['lng']?.toDouble(),
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      itemCount: json['item_count'] ?? 0,
      createdAt: json['created_at'] ?? '',
    );
  }
}
