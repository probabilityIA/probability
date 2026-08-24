class MobileOrderSummary {
  final String id;
  final String orderNumber;
  final String internalNumber;
  final String integrationType;
  final String platform;
  final String status;
  final bool isPaid;
  final bool isCod;
  final double? codTotal;
  final double subtotal;
  final double tax;
  final double discount;
  final double shippingCost;
  final double totalAmount;
  final String currency;
  final String customerName;
  final String customerEmail;
  final String customerPhone;
  final String shippingStreet;
  final String shippingCity;
  final String warehouseName;
  final String userName;
  final String createdAt;

  MobileOrderSummary({
    required this.id,
    required this.orderNumber,
    required this.internalNumber,
    required this.integrationType,
    required this.platform,
    required this.status,
    required this.isPaid,
    required this.isCod,
    this.codTotal,
    required this.subtotal,
    required this.tax,
    required this.discount,
    required this.shippingCost,
    required this.totalAmount,
    required this.currency,
    required this.customerName,
    required this.customerEmail,
    required this.customerPhone,
    required this.shippingStreet,
    required this.shippingCity,
    required this.warehouseName,
    required this.userName,
    required this.createdAt,
  });

  factory MobileOrderSummary.fromJson(Map<String, dynamic> json) {
    double num_(dynamic v) => (v as num?)?.toDouble() ?? 0;
    return MobileOrderSummary(
      id: json['id']?.toString() ?? '',
      orderNumber: json['order_number'] ?? '',
      internalNumber: json['internal_number'] ?? '',
      integrationType: json['integration_type'] ?? '',
      platform: json['platform'] ?? '',
      status: json['status'] ?? '',
      isPaid: json['is_paid'] ?? false,
      isCod: json['is_cod'] ?? false,
      codTotal: (json['cod_total'] as num?)?.toDouble(),
      subtotal: num_(json['subtotal']),
      tax: num_(json['tax']),
      discount: num_(json['discount']),
      shippingCost: num_(json['shipping_cost']),
      totalAmount: num_(json['total_amount']),
      currency: json['currency'] ?? 'COP',
      customerName: json['customer_name'] ?? '',
      customerEmail: json['customer_email'] ?? '',
      customerPhone: json['customer_phone'] ?? '',
      shippingStreet: json['shipping_street'] ?? '',
      shippingCity: json['shipping_city'] ?? '',
      warehouseName: json['warehouse_name'] ?? '',
      userName: json['user_name'] ?? '',
      createdAt: json['created_at'] ?? '',
    );
  }
}

class MobileOrderLine {
  final String sku;
  final String name;
  final int quantity;
  final double unitPrice;
  final double totalPrice;

  MobileOrderLine({
    required this.sku,
    required this.name,
    required this.quantity,
    required this.unitPrice,
    required this.totalPrice,
  });

  factory MobileOrderLine.fromJson(Map<String, dynamic> json) {
    return MobileOrderLine(
      sku: json['sku'] ?? '',
      name: json['name'] ?? '',
      quantity: (json['quantity'] as num?)?.toInt() ?? 0,
      unitPrice: (json['unit_price'] as num?)?.toDouble() ?? 0,
      totalPrice: (json['total_price'] as num?)?.toDouble() ?? 0,
    );
  }
}

class MobileShipmentSummary {
  final int id;
  final String? trackingNumber;
  final String? carrier;
  final String? guideUrl;
  final String status;
  final String? carrierStatus;
  final String destinationCity;
  final double? totalCost;
  final double? carrierCost;
  final double? appliedMargin;
  final double? codCarrierFee;
  final double? codProbabilityMargin;

  MobileShipmentSummary({
    required this.id,
    this.trackingNumber,
    this.carrier,
    this.guideUrl,
    required this.status,
    this.carrierStatus,
    required this.destinationCity,
    this.totalCost,
    this.carrierCost,
    this.appliedMargin,
    this.codCarrierFee,
    this.codProbabilityMargin,
  });

  factory MobileShipmentSummary.fromJson(Map<String, dynamic> json) {
    return MobileShipmentSummary(
      id: (json['id'] as num?)?.toInt() ?? 0,
      trackingNumber: json['tracking_number'],
      carrier: json['carrier'],
      guideUrl: json['guide_url'],
      status: json['status'] ?? '',
      carrierStatus: json['carrier_status'],
      destinationCity: json['destination_city'] ?? '',
      totalCost: (json['total_cost'] as num?)?.toDouble(),
      carrierCost: (json['carrier_cost'] as num?)?.toDouble(),
      appliedMargin: (json['applied_margin'] as num?)?.toDouble(),
      codCarrierFee: (json['cod_carrier_fee'] as num?)?.toDouble(),
      codProbabilityMargin: (json['cod_probability_margin'] as num?)?.toDouble(),
    );
  }
}

class MobileInvoiceSummary {
  final int id;
  final String invoiceNumber;
  final String status;
  final double totalAmount;
  final String? invoiceUrl;
  final String? cufe;
  final String? issuedAt;

  MobileInvoiceSummary({
    required this.id,
    required this.invoiceNumber,
    required this.status,
    required this.totalAmount,
    this.invoiceUrl,
    this.cufe,
    this.issuedAt,
  });

  factory MobileInvoiceSummary.fromJson(Map<String, dynamic> json) {
    return MobileInvoiceSummary(
      id: (json['id'] as num?)?.toInt() ?? 0,
      invoiceNumber: json['invoice_number'] ?? '',
      status: json['status'] ?? '',
      totalAmount: (json['total_amount'] as num?)?.toDouble() ?? 0,
      invoiceUrl: json['invoice_url'],
      cufe: json['cufe'],
      issuedAt: json['issued_at'],
    );
  }
}

class MobileOrderFull {
  final MobileOrderSummary order;
  final List<MobileOrderLine> items;
  final MobileShipmentSummary? shipment;
  final MobileInvoiceSummary? invoice;

  MobileOrderFull({
    required this.order,
    required this.items,
    this.shipment,
    this.invoice,
  });

  factory MobileOrderFull.fromJson(Map<String, dynamic> json) {
    return MobileOrderFull(
      order: MobileOrderSummary.fromJson(
        Map<String, dynamic>.from(json['order'] ?? const {}),
      ),
      items: (json['items'] as List<dynamic>?)
              ?.map((e) => MobileOrderLine.fromJson(Map<String, dynamic>.from(e)))
              .toList() ??
          const [],
      shipment: json['shipment'] == null
          ? null
          : MobileShipmentSummary.fromJson(Map<String, dynamic>.from(json['shipment'])),
      invoice: json['invoice'] == null
          ? null
          : MobileInvoiceSummary.fromJson(Map<String, dynamic>.from(json['invoice'])),
    );
  }
}
