class Shipment {
  final int id;
  final String createdAt;
  final String updatedAt;
  final String? orderId;
  final String? clientName;
  final String? destinationAddress;
  final String? trackingNumber;
  final String? trackingUrl;
  final String? carrier;
  final String? carrierCode;
  final String? guideId;
  final String? guideUrl;
  final String status;
  final String? shippedAt;
  final String? deliveredAt;
  final double? shippingCost;
  final double? insuranceCost;
  final double? totalCost;
  final double? weight;
  final double? height;
  final double? width;
  final double? length;
  final String? warehouseName;
  final String? driverName;
  final bool isLastMile;
  final bool isTest;
  final String? estimatedDelivery;
  final String? deliveryNotes;
  final String? customerName;
  final String? customerEmail;
  final String? customerPhone;
  final String? customerDni;
  final String? orderNumber;

  Shipment({
    required this.id, required this.createdAt, required this.updatedAt, this.orderId,
    this.clientName, this.destinationAddress, this.trackingNumber, this.trackingUrl,
    this.carrier, this.carrierCode, this.guideId, this.guideUrl, required this.status,
    this.shippedAt, this.deliveredAt, this.shippingCost, this.insuranceCost, this.totalCost,
    this.weight, this.height, this.width, this.length, this.warehouseName, this.driverName,
    required this.isLastMile, required this.isTest, this.estimatedDelivery, this.deliveryNotes,
    this.customerName, this.customerEmail, this.customerPhone, this.customerDni, this.orderNumber,
  });

  factory Shipment.fromJson(Map<String, dynamic> json) {
    return Shipment(
      id: json['id'] ?? 0, createdAt: json['created_at'] ?? '', updatedAt: json['updated_at'] ?? '',
      orderId: json['order_id'], clientName: json['client_name'], destinationAddress: json['destination_address'],
      trackingNumber: json['tracking_number'], trackingUrl: json['tracking_url'],
      carrier: json['carrier'], carrierCode: json['carrier_code'],
      guideId: json['guide_id'], guideUrl: json['guide_url'], status: json['status'] ?? 'pending',
      shippedAt: json['shipped_at'], deliveredAt: json['delivered_at'],
      shippingCost: json['shipping_cost']?.toDouble(), insuranceCost: json['insurance_cost']?.toDouble(),
      totalCost: json['total_cost']?.toDouble(), weight: json['weight']?.toDouble(),
      height: json['height']?.toDouble(), width: json['width']?.toDouble(), length: json['length']?.toDouble(),
      warehouseName: json['warehouse_name'], driverName: json['driver_name'],
      isLastMile: json['is_last_mile'] ?? false, isTest: json['is_test'] ?? false,
      estimatedDelivery: json['estimated_delivery'], deliveryNotes: json['delivery_notes'],
      customerName: json['customer_name'], customerEmail: json['customer_email'],
      customerPhone: json['customer_phone'], customerDni: json['customer_dni'], orderNumber: json['order_number'],
    );
  }
}

class OriginAddress {
  final int id;
  final int businessId;
  final String alias;
  final String company;
  final String firstName;
  final String lastName;
  final String email;
  final String phone;
  final String street;
  final String? suburb;
  final String cityDaneCode;
  final String city;
  final String state;
  final String? postalCode;
  final bool isDefault;
  final String createdAt;
  final String updatedAt;

  OriginAddress({
    required this.id, required this.businessId, required this.alias, required this.company,
    required this.firstName, required this.lastName, required this.email, required this.phone,
    required this.street, this.suburb, required this.cityDaneCode, required this.city,
    required this.state, this.postalCode, required this.isDefault,
    required this.createdAt, required this.updatedAt,
  });

  factory OriginAddress.fromJson(Map<String, dynamic> json) {
    return OriginAddress(
      id: json['id'] ?? 0, businessId: json['business_id'] ?? 0, alias: json['alias'] ?? '',
      company: json['company'] ?? '', firstName: json['first_name'] ?? '', lastName: json['last_name'] ?? '',
      email: json['email'] ?? '', phone: json['phone'] ?? '', street: json['street'] ?? '',
      suburb: json['suburb'], cityDaneCode: json['city_dane_code'] ?? '', city: json['city'] ?? '',
      state: json['state'] ?? '', postalCode: json['postal_code'], isDefault: json['is_default'] ?? false,
      createdAt: json['created_at'] ?? '', updatedAt: json['updated_at'] ?? '',
    );
  }
}

class EnvioClickRate {
  final int idRate;
  final int idProduct;
  final String product;
  final int idCarrier;
  final String carrier;
  final double flete;
  final int deliveryDays;
  final String quotationType;
  final double? minimumInsurance;
  final double? extraInsurance;
  final bool? cod;

  EnvioClickRate({
    required this.idRate, required this.idProduct, required this.product,
    required this.idCarrier, required this.carrier, required this.flete,
    required this.deliveryDays, required this.quotationType, this.minimumInsurance,
    this.extraInsurance, this.cod,
  });

  factory EnvioClickRate.fromJson(Map<String, dynamic> json) {
    return EnvioClickRate(
      idRate: json['idRate'] ?? 0, idProduct: json['idProduct'] ?? 0, product: json['product'] ?? '',
      idCarrier: json['idCarrier'] ?? 0, carrier: json['carrier'] ?? '', flete: (json['flete'] ?? 0).toDouble(),
      deliveryDays: json['deliveryDays'] ?? 0, quotationType: json['quotationType'] ?? '',
      minimumInsurance: json['minimumInsurance']?.toDouble(), extraInsurance: json['extraInsurance']?.toDouble(),
      cod: json['cod'],
    );
  }
}

class EnvioClickTrackHistory {
  final String date;
  final String status;
  final String description;
  final String location;

  EnvioClickTrackHistory({required this.date, required this.status, required this.description, required this.location});

  factory EnvioClickTrackHistory.fromJson(Map<String, dynamic> json) {
    return EnvioClickTrackHistory(
      date: json['date'] ?? '', status: json['status'] ?? '',
      description: json['description'] ?? '', location: json['location'] ?? '',
    );
  }
}

class GetShipmentsParams {
  final int? page;
  final int? pageSize;
  final String? orderId;
  final String? trackingNumber;
  final String? carrier;
  final String? status;
  final String? customerName;
  final int? businessId;
  final String? startDate;
  final String? endDate;
  final String? sortBy;
  final String? sortOrder;
  final bool? isTest;

  GetShipmentsParams({this.page, this.pageSize, this.orderId, this.trackingNumber, this.carrier, this.status, this.customerName, this.businessId, this.startDate, this.endDate, this.sortBy, this.sortOrder, this.isTest});

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (orderId != null) params['order_id'] = orderId;
    if (trackingNumber != null) params['tracking_number'] = trackingNumber;
    if (carrier != null) params['carrier'] = carrier;
    if (status != null) params['status'] = status;
    if (customerName != null) params['customer_name'] = customerName;
    if (businessId != null) params['business_id'] = businessId;
    if (startDate != null) params['start_date'] = startDate;
    if (endDate != null) params['end_date'] = endDate;
    if (sortBy != null) params['sort_by'] = sortBy;
    if (sortOrder != null) params['sort_order'] = sortOrder;
    if (isTest != null) params['is_test'] = isTest;
    return params;
  }
}

class EnvioClickAddress {
  final String? company;
  final String? firstName;
  final String? lastName;
  final String? email;
  final String? phone;
  final String address;
  final String? suburb;
  final String? crossStreet;
  final String? reference;
  final String daneCode;

  EnvioClickAddress({this.company, this.firstName, this.lastName, this.email, this.phone, required this.address, this.suburb, this.crossStreet, this.reference, required this.daneCode});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'address': address, 'daneCode': daneCode};
    if (company != null) json['company'] = company;
    if (firstName != null) json['firstName'] = firstName;
    if (lastName != null) json['lastName'] = lastName;
    if (email != null) json['email'] = email;
    if (phone != null) json['phone'] = phone;
    if (suburb != null) json['suburb'] = suburb;
    if (crossStreet != null) json['crossStreet'] = crossStreet;
    if (reference != null) json['reference'] = reference;
    return json;
  }
}

class EnvioClickPackage {
  final double weight;
  final double height;
  final double width;
  final double length;

  EnvioClickPackage({required this.weight, required this.height, required this.width, required this.length});

  Map<String, dynamic> toJson() => {'weight': weight, 'height': height, 'width': width, 'length': length};
}

class EnvioClickQuoteRequest {
  final int? businessId;
  final int? idRate;
  final String? myShipmentReference;
  final String? externalOrderId;
  final String? orderUuid;
  final bool? requestPickup;
  final String? pickupDate;
  final bool? insurance;
  final String description;
  final double contentValue;
  final double? codValue;
  final bool includeGuideCost;
  final String codPaymentMethod;
  final double? totalCost;
  final List<EnvioClickPackage> packages;
  final EnvioClickAddress origin;
  final EnvioClickAddress destination;

  EnvioClickQuoteRequest({
    this.businessId, this.idRate, this.myShipmentReference, this.externalOrderId,
    this.orderUuid, this.requestPickup, this.pickupDate, this.insurance,
    required this.description, required this.contentValue, this.codValue,
    required this.includeGuideCost, required this.codPaymentMethod, this.totalCost,
    required this.packages, required this.origin, required this.destination,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'description': description, 'contentValue': contentValue,
      'includeGuideCost': includeGuideCost, 'codPaymentMethod': codPaymentMethod,
      'packages': packages.map((p) => p.toJson()).toList(),
      'origin': origin.toJson(), 'destination': destination.toJson(),
    };
    if (businessId != null) json['business_id'] = businessId;
    if (idRate != null) json['idRate'] = idRate;
    if (myShipmentReference != null) json['myShipmentReference'] = myShipmentReference;
    if (externalOrderId != null) json['external_order_id'] = externalOrderId;
    if (orderUuid != null) json['order_uuid'] = orderUuid;
    if (requestPickup != null) json['requestPickup'] = requestPickup;
    if (pickupDate != null) json['pickupDate'] = pickupDate;
    if (insurance != null) json['insurance'] = insurance;
    if (codValue != null) json['codValue'] = codValue;
    if (totalCost != null) json['totalCost'] = totalCost;
    return json;
  }
}

class CreateOriginAddressDTO {
  final String alias;
  final String company;
  final String firstName;
  final String lastName;
  final String email;
  final String phone;
  final String street;
  final String? suburb;
  final String cityDaneCode;
  final String city;
  final String state;
  final String? postalCode;
  final bool? isDefault;

  CreateOriginAddressDTO({required this.alias, required this.company, required this.firstName, required this.lastName, required this.email, required this.phone, required this.street, this.suburb, required this.cityDaneCode, required this.city, required this.state, this.postalCode, this.isDefault});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'alias': alias, 'company': company, 'first_name': firstName, 'last_name': lastName,
      'email': email, 'phone': phone, 'street': street, 'city_dane_code': cityDaneCode,
      'city': city, 'state': state,
    };
    if (suburb != null) json['suburb'] = suburb;
    if (postalCode != null) json['postal_code'] = postalCode;
    if (isDefault != null) json['is_default'] = isDefault;
    return json;
  }
}
