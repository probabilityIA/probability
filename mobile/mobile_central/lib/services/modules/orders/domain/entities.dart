class OrderStatusInfo {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String? category;
  final String? color;

  OrderStatusInfo({
    required this.id,
    required this.code,
    required this.name,
    this.description,
    this.category,
    this.color,
  });

  factory OrderStatusInfo.fromJson(Map<String, dynamic> json) {
    return OrderStatusInfo(
      id: json['id'] ?? 0,
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      category: json['category'],
      color: json['color'],
    );
  }
}

class PaymentStatusInfo {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String? category;
  final String? color;

  PaymentStatusInfo({
    required this.id,
    required this.code,
    required this.name,
    this.description,
    this.category,
    this.color,
  });

  factory PaymentStatusInfo.fromJson(Map<String, dynamic> json) {
    return PaymentStatusInfo(
      id: json['id'] ?? 0,
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      category: json['category'],
      color: json['color'],
    );
  }
}

class FulfillmentStatusInfo {
  final int id;
  final String code;
  final String name;
  final String? description;
  final String? category;
  final String? color;

  FulfillmentStatusInfo({
    required this.id,
    required this.code,
    required this.name,
    this.description,
    this.category,
    this.color,
  });

  factory FulfillmentStatusInfo.fromJson(Map<String, dynamic> json) {
    return FulfillmentStatusInfo(
      id: json['id'] ?? 0,
      code: json['code'] ?? '',
      name: json['name'] ?? '',
      description: json['description'],
      category: json['category'],
      color: json['color'],
    );
  }
}

class Order {
  final String id;
  final String createdAt;
  final String updatedAt;
  final String? deletedAt;
  final int? businessId;
  final int integrationId;
  final String integrationType;
  final String? integrationLogoUrl;
  final String? integrationName;
  final String platform;
  final String externalId;
  final String orderNumber;
  final String internalNumber;
  final double subtotal;
  final double tax;
  final double discount;
  final double shippingCost;
  final double? shippingDiscount;
  final double? shippingDiscountPresentment;
  final double totalAmount;
  final String currency;
  final double? codTotal;
  final double? subtotalPresentment;
  final double? taxPresentment;
  final double? discountPresentment;
  final double? shippingCostPresentment;
  final double? totalAmountPresentment;
  final String? currencyPresentment;
  final int? customerId;
  final String customerName;
  final String? customerFirstName;
  final String? customerLastName;
  final String customerEmail;
  final String customerPhone;
  final String customerDni;
  final String shippingStreet;
  final String shippingCity;
  final String shippingState;
  final String shippingCountry;
  final String shippingPostalCode;
  final String? shippingHouse;
  final String? shippingBarrio;
  final double? shippingLat;
  final double? shippingLng;
  final int paymentMethodId;
  final bool isPaid;
  final String? paidAt;
  final String? trackingNumber;
  final String? trackingLink;
  final String? guideId;
  final String? guideLink;
  final String? deliveryDate;
  final String? deliveredAt;
  final double? deliveryProbability;
  final int? warehouseId;
  final String warehouseName;
  final int? driverId;
  final String driverName;
  final bool isLastMile;
  final double? weight;
  final double? height;
  final double? width;
  final double? length;
  final String? boxes;
  final int? orderTypeId;
  final String orderTypeName;
  final String status;
  final String originalStatus;
  final int? statusId;
  final OrderStatusInfo? orderStatus;
  final int? paymentStatusId;
  final int? fulfillmentStatusId;
  final PaymentStatusInfo? paymentStatus;
  final FulfillmentStatusInfo? fulfillmentStatus;
  final String? notes;
  final String? coupon;
  final bool? approved;
  final int? userId;
  final String userName;
  final bool? isConfirmed;
  final String? novelty;
  final bool invoiceable;
  final String? invoiceUrl;
  final String? invoiceId;
  final String? invoiceProvider;
  final String? invoiceStatus;
  final String? orderStatusUrl;
  final dynamic items;
  final dynamic orderItems;
  final dynamic metadata;
  final dynamic financialDetails;
  final dynamic shippingDetails;
  final dynamic paymentDetails;
  final dynamic fulfillmentDetails;
  final String occurredAt;
  final String importedAt;
  final List<String>? negativeFactors;

  Order({
    required this.id,
    required this.createdAt,
    required this.updatedAt,
    this.deletedAt,
    this.businessId,
    required this.integrationId,
    required this.integrationType,
    this.integrationLogoUrl,
    this.integrationName,
    required this.platform,
    required this.externalId,
    required this.orderNumber,
    required this.internalNumber,
    required this.subtotal,
    required this.tax,
    required this.discount,
    required this.shippingCost,
    this.shippingDiscount,
    this.shippingDiscountPresentment,
    required this.totalAmount,
    required this.currency,
    this.codTotal,
    this.subtotalPresentment,
    this.taxPresentment,
    this.discountPresentment,
    this.shippingCostPresentment,
    this.totalAmountPresentment,
    this.currencyPresentment,
    this.customerId,
    required this.customerName,
    this.customerFirstName,
    this.customerLastName,
    required this.customerEmail,
    required this.customerPhone,
    required this.customerDni,
    required this.shippingStreet,
    required this.shippingCity,
    required this.shippingState,
    required this.shippingCountry,
    required this.shippingPostalCode,
    this.shippingHouse,
    this.shippingBarrio,
    this.shippingLat,
    this.shippingLng,
    required this.paymentMethodId,
    required this.isPaid,
    this.paidAt,
    this.trackingNumber,
    this.trackingLink,
    this.guideId,
    this.guideLink,
    this.deliveryDate,
    this.deliveredAt,
    this.deliveryProbability,
    this.warehouseId,
    required this.warehouseName,
    this.driverId,
    required this.driverName,
    required this.isLastMile,
    this.weight,
    this.height,
    this.width,
    this.length,
    this.boxes,
    this.orderTypeId,
    required this.orderTypeName,
    required this.status,
    required this.originalStatus,
    this.statusId,
    this.orderStatus,
    this.paymentStatusId,
    this.fulfillmentStatusId,
    this.paymentStatus,
    this.fulfillmentStatus,
    this.notes,
    this.coupon,
    this.approved,
    this.userId,
    required this.userName,
    this.isConfirmed,
    this.novelty,
    required this.invoiceable,
    this.invoiceUrl,
    this.invoiceId,
    this.invoiceProvider,
    this.invoiceStatus,
    this.orderStatusUrl,
    this.items,
    this.orderItems,
    this.metadata,
    this.financialDetails,
    this.shippingDetails,
    this.paymentDetails,
    this.fulfillmentDetails,
    required this.occurredAt,
    required this.importedAt,
    this.negativeFactors,
  });

  factory Order.fromJson(Map<String, dynamic> json) {
    return Order(
      id: json['id']?.toString() ?? '',
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      deletedAt: json['deleted_at'],
      businessId: json['business_id'],
      integrationId: json['integration_id'] ?? 0,
      integrationType: json['integration_type'] ?? '',
      integrationLogoUrl: json['integration_logo_url'],
      integrationName: json['integration_name'],
      platform: json['platform'] ?? '',
      externalId: json['external_id'] ?? '',
      orderNumber: json['order_number'] ?? '',
      internalNumber: json['internal_number'] ?? '',
      subtotal: (json['subtotal'] ?? 0).toDouble(),
      tax: (json['tax'] ?? 0).toDouble(),
      discount: (json['discount'] ?? 0).toDouble(),
      shippingCost: (json['shipping_cost'] ?? 0).toDouble(),
      shippingDiscount: json['shipping_discount']?.toDouble(),
      shippingDiscountPresentment: json['shipping_discount_presentment']?.toDouble(),
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      currency: json['currency'] ?? '',
      codTotal: json['cod_total']?.toDouble(),
      subtotalPresentment: json['subtotal_presentment']?.toDouble(),
      taxPresentment: json['tax_presentment']?.toDouble(),
      discountPresentment: json['discount_presentment']?.toDouble(),
      shippingCostPresentment: json['shipping_cost_presentment']?.toDouble(),
      totalAmountPresentment: json['total_amount_presentment']?.toDouble(),
      currencyPresentment: json['currency_presentment'],
      customerId: json['customer_id'],
      customerName: json['customer_name'] ?? '',
      customerFirstName: json['customer_first_name'],
      customerLastName: json['customer_last_name'],
      customerEmail: json['customer_email'] ?? '',
      customerPhone: json['customer_phone'] ?? '',
      customerDni: json['customer_dni'] ?? '',
      shippingStreet: json['shipping_street'] ?? '',
      shippingCity: json['shipping_city'] ?? '',
      shippingState: json['shipping_state'] ?? '',
      shippingCountry: json['shipping_country'] ?? '',
      shippingPostalCode: json['shipping_postal_code'] ?? '',
      shippingHouse: json['shipping_house'],
      shippingBarrio: json['shipping_barrio'],
      shippingLat: json['shipping_lat']?.toDouble(),
      shippingLng: json['shipping_lng']?.toDouble(),
      paymentMethodId: json['payment_method_id'] ?? 0,
      isPaid: json['is_paid'] ?? false,
      paidAt: json['paid_at'],
      trackingNumber: json['tracking_number'],
      trackingLink: json['tracking_link'],
      guideId: json['guide_id'],
      guideLink: json['guide_link'],
      deliveryDate: json['delivery_date'],
      deliveredAt: json['delivered_at'],
      deliveryProbability: json['delivery_probability']?.toDouble(),
      warehouseId: json['warehouse_id'],
      warehouseName: json['warehouse_name'] ?? '',
      driverId: json['driver_id'],
      driverName: json['driver_name'] ?? '',
      isLastMile: json['is_last_mile'] ?? false,
      weight: json['weight']?.toDouble(),
      height: json['height']?.toDouble(),
      width: json['width']?.toDouble(),
      length: json['length']?.toDouble(),
      boxes: json['boxes'],
      orderTypeId: json['order_type_id'],
      orderTypeName: json['order_type_name'] ?? '',
      status: json['status'] ?? '',
      originalStatus: json['original_status'] ?? '',
      statusId: json['status_id'],
      orderStatus: json['order_status'] != null
          ? OrderStatusInfo.fromJson(json['order_status'])
          : null,
      paymentStatusId: json['payment_status_id'],
      fulfillmentStatusId: json['fulfillment_status_id'],
      paymentStatus: json['payment_status'] != null
          ? PaymentStatusInfo.fromJson(json['payment_status'])
          : null,
      fulfillmentStatus: json['fulfillment_status'] != null
          ? FulfillmentStatusInfo.fromJson(json['fulfillment_status'])
          : null,
      notes: json['notes'],
      coupon: json['coupon'],
      approved: json['approved'],
      userId: json['user_id'],
      userName: json['user_name'] ?? '',
      isConfirmed: json['is_confirmed'],
      novelty: json['novelty'],
      invoiceable: json['invoiceable'] ?? false,
      invoiceUrl: json['invoice_url'],
      invoiceId: json['invoice_id'],
      invoiceProvider: json['invoice_provider'],
      invoiceStatus: json['invoice_status'],
      orderStatusUrl: json['order_status_url'],
      items: json['items'],
      orderItems: json['order_items'],
      metadata: json['metadata'],
      financialDetails: json['financial_details'],
      shippingDetails: json['shipping_details'],
      paymentDetails: json['payment_details'],
      fulfillmentDetails: json['fulfillment_details'],
      occurredAt: json['occurred_at'] ?? '',
      importedAt: json['imported_at'] ?? '',
      negativeFactors: (json['negative_factors'] as List<dynamic>?)
          ?.map((e) => e.toString())
          .toList(),
    );
  }
}

class GetOrdersParams {
  final int? page;
  final int? pageSize;
  final int? businessId;
  final int? integrationId;
  final String? integrationType;
  final String? status;
  final String? customerEmail;
  final String? customerPhone;
  final String? orderNumber;
  final String? internalNumber;
  final String? platform;
  final String? currencyPresentment;
  final bool? isPaid;
  final bool? isCod;
  final int? paymentStatusId;
  final int? fulfillmentStatusId;
  final int? warehouseId;
  final int? driverId;
  final String? startDate;
  final String? endDate;
  final String? invoiceStatus;
  final String? sortBy;
  final String? sortOrder;

  GetOrdersParams({
    this.page,
    this.pageSize,
    this.businessId,
    this.integrationId,
    this.integrationType,
    this.status,
    this.customerEmail,
    this.customerPhone,
    this.orderNumber,
    this.internalNumber,
    this.platform,
    this.currencyPresentment,
    this.isPaid,
    this.isCod,
    this.paymentStatusId,
    this.fulfillmentStatusId,
    this.warehouseId,
    this.driverId,
    this.startDate,
    this.endDate,
    this.invoiceStatus,
    this.sortBy,
    this.sortOrder,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    if (businessId != null) params['business_id'] = businessId;
    if (integrationId != null) params['integration_id'] = integrationId;
    if (integrationType != null) params['integration_type'] = integrationType;
    if (status != null) params['status'] = status;
    if (customerEmail != null) params['customer_email'] = customerEmail;
    if (customerPhone != null) params['customer_phone'] = customerPhone;
    if (orderNumber != null) params['order_number'] = orderNumber;
    if (internalNumber != null) params['internal_number'] = internalNumber;
    if (platform != null) params['platform'] = platform;
    if (currencyPresentment != null) params['currency_presentment'] = currencyPresentment;
    if (isPaid != null) params['is_paid'] = isPaid;
    if (isCod != null) params['is_cod'] = isCod;
    if (paymentStatusId != null) params['payment_status_id'] = paymentStatusId;
    if (fulfillmentStatusId != null) params['fulfillment_status_id'] = fulfillmentStatusId;
    if (warehouseId != null) params['warehouse_id'] = warehouseId;
    if (driverId != null) params['driver_id'] = driverId;
    if (startDate != null) params['start_date'] = startDate;
    if (endDate != null) params['end_date'] = endDate;
    if (invoiceStatus != null) params['invoice_status'] = invoiceStatus;
    if (sortBy != null) params['sort_by'] = sortBy;
    if (sortOrder != null) params['sort_order'] = sortOrder;
    return params;
  }
}

class CreateOrderDTO {
  final int? businessId;
  final int integrationId;
  final String integrationType;
  final String platform;
  final String externalId;
  final String? orderNumber;
  final double subtotal;
  final double totalAmount;
  final int paymentMethodId;
  final String? customerName;
  final String? customerEmail;
  final String? customerPhone;
  final String? customerDni;
  final double? tax;
  final double? discount;
  final double? shippingCost;
  final String? currency;
  final double? codTotal;
  final String? shippingStreet;
  final String? shippingCity;
  final String? shippingState;
  final String? shippingCountry;
  final String? shippingPostalCode;
  final bool? isPaid;
  final int? warehouseId;
  final String? status;
  final String? notes;
  final bool? invoiceable;
  final dynamic items;
  final dynamic metadata;

  CreateOrderDTO({
    this.businessId,
    required this.integrationId,
    required this.integrationType,
    required this.platform,
    required this.externalId,
    this.orderNumber,
    required this.subtotal,
    required this.totalAmount,
    required this.paymentMethodId,
    this.customerName,
    this.customerEmail,
    this.customerPhone,
    this.customerDni,
    this.tax,
    this.discount,
    this.shippingCost,
    this.currency,
    this.codTotal,
    this.shippingStreet,
    this.shippingCity,
    this.shippingState,
    this.shippingCountry,
    this.shippingPostalCode,
    this.isPaid,
    this.warehouseId,
    this.status,
    this.notes,
    this.invoiceable,
    this.items,
    this.metadata,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'integration_id': integrationId,
      'integration_type': integrationType,
      'platform': platform,
      'external_id': externalId,
      'subtotal': subtotal,
      'total_amount': totalAmount,
      'payment_method_id': paymentMethodId,
    };
    if (businessId != null) json['business_id'] = businessId;
    if (orderNumber != null) json['order_number'] = orderNumber;
    if (customerName != null) json['customer_name'] = customerName;
    if (customerEmail != null) json['customer_email'] = customerEmail;
    if (customerPhone != null) json['customer_phone'] = customerPhone;
    if (customerDni != null) json['customer_dni'] = customerDni;
    if (tax != null) json['tax'] = tax;
    if (discount != null) json['discount'] = discount;
    if (shippingCost != null) json['shipping_cost'] = shippingCost;
    if (currency != null) json['currency'] = currency;
    if (codTotal != null) json['cod_total'] = codTotal;
    if (shippingStreet != null) json['shipping_street'] = shippingStreet;
    if (shippingCity != null) json['shipping_city'] = shippingCity;
    if (shippingState != null) json['shipping_state'] = shippingState;
    if (shippingCountry != null) json['shipping_country'] = shippingCountry;
    if (shippingPostalCode != null) json['shipping_postal_code'] = shippingPostalCode;
    if (isPaid != null) json['is_paid'] = isPaid;
    if (warehouseId != null) json['warehouse_id'] = warehouseId;
    if (status != null) json['status'] = status;
    if (notes != null) json['notes'] = notes;
    if (invoiceable != null) json['invoiceable'] = invoiceable;
    if (items != null) json['items'] = items;
    if (metadata != null) json['metadata'] = metadata;
    return json;
  }
}

class UpdateOrderDTO {
  final double? subtotal;
  final double? tax;
  final double? discount;
  final double? shippingCost;
  final double? totalAmount;
  final String? currency;
  final double? codTotal;
  final String? customerName;
  final String? customerEmail;
  final String? customerPhone;
  final String? customerDni;
  final String? shippingStreet;
  final String? shippingCity;
  final String? shippingState;
  final String? shippingCountry;
  final String? shippingPostalCode;
  final int? paymentMethodId;
  final bool? isPaid;
  final String? trackingNumber;
  final String? trackingLink;
  final int? warehouseId;
  final String? warehouseName;
  final int? driverId;
  final String? driverName;
  final String? status;
  final int? statusId;
  final int? paymentStatusId;
  final int? fulfillmentStatusId;
  final String? notes;
  final bool? isConfirmed;
  final String? confirmationStatus;
  final String? novelty;
  final bool? invoiceable;
  final dynamic items;
  final dynamic metadata;

  UpdateOrderDTO({
    this.subtotal,
    this.tax,
    this.discount,
    this.shippingCost,
    this.totalAmount,
    this.currency,
    this.codTotal,
    this.customerName,
    this.customerEmail,
    this.customerPhone,
    this.customerDni,
    this.shippingStreet,
    this.shippingCity,
    this.shippingState,
    this.shippingCountry,
    this.shippingPostalCode,
    this.paymentMethodId,
    this.isPaid,
    this.trackingNumber,
    this.trackingLink,
    this.warehouseId,
    this.warehouseName,
    this.driverId,
    this.driverName,
    this.status,
    this.statusId,
    this.paymentStatusId,
    this.fulfillmentStatusId,
    this.notes,
    this.isConfirmed,
    this.confirmationStatus,
    this.novelty,
    this.invoiceable,
    this.items,
    this.metadata,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (subtotal != null) json['subtotal'] = subtotal;
    if (tax != null) json['tax'] = tax;
    if (discount != null) json['discount'] = discount;
    if (shippingCost != null) json['shipping_cost'] = shippingCost;
    if (totalAmount != null) json['total_amount'] = totalAmount;
    if (currency != null) json['currency'] = currency;
    if (codTotal != null) json['cod_total'] = codTotal;
    if (customerName != null) json['customer_name'] = customerName;
    if (customerEmail != null) json['customer_email'] = customerEmail;
    if (customerPhone != null) json['customer_phone'] = customerPhone;
    if (customerDni != null) json['customer_dni'] = customerDni;
    if (shippingStreet != null) json['shipping_street'] = shippingStreet;
    if (shippingCity != null) json['shipping_city'] = shippingCity;
    if (shippingState != null) json['shipping_state'] = shippingState;
    if (shippingCountry != null) json['shipping_country'] = shippingCountry;
    if (shippingPostalCode != null) json['shipping_postal_code'] = shippingPostalCode;
    if (paymentMethodId != null) json['payment_method_id'] = paymentMethodId;
    if (isPaid != null) json['is_paid'] = isPaid;
    if (trackingNumber != null) json['tracking_number'] = trackingNumber;
    if (trackingLink != null) json['tracking_link'] = trackingLink;
    if (warehouseId != null) json['warehouse_id'] = warehouseId;
    if (warehouseName != null) json['warehouse_name'] = warehouseName;
    if (driverId != null) json['driver_id'] = driverId;
    if (driverName != null) json['driver_name'] = driverName;
    if (status != null) json['status'] = status;
    if (statusId != null) json['status_id'] = statusId;
    if (paymentStatusId != null) json['payment_status_id'] = paymentStatusId;
    if (fulfillmentStatusId != null) json['fulfillment_status_id'] = fulfillmentStatusId;
    if (notes != null) json['notes'] = notes;
    if (isConfirmed != null) json['is_confirmed'] = isConfirmed;
    if (confirmationStatus != null) json['confirmation_status'] = confirmationStatus;
    if (novelty != null) json['novelty'] = novelty;
    if (invoiceable != null) json['invoiceable'] = invoiceable;
    if (items != null) json['items'] = items;
    if (metadata != null) json['metadata'] = metadata;
    return json;
  }
}
