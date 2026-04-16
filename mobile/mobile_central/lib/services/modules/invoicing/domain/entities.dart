class InvoiceItem {
  final int id;
  final int invoiceId;
  final String? productSku;
  final String productName;
  final String? description;
  final int quantity;
  final double unitPrice;
  final double totalPrice;
  final double tax;
  final double? taxRate;
  final double discount;

  InvoiceItem({
    required this.id,
    required this.invoiceId,
    this.productSku,
    required this.productName,
    this.description,
    required this.quantity,
    required this.unitPrice,
    required this.totalPrice,
    required this.tax,
    this.taxRate,
    required this.discount,
  });

  factory InvoiceItem.fromJson(Map<String, dynamic> json) {
    return InvoiceItem(
      id: json['id'] ?? 0,
      invoiceId: json['invoice_id'] ?? 0,
      productSku: json['product_sku'],
      productName: json['product_name'] ?? '',
      description: json['description'],
      quantity: json['quantity'] ?? 0,
      unitPrice: (json['unit_price'] ?? 0).toDouble(),
      totalPrice: (json['total_price'] ?? 0).toDouble(),
      tax: (json['tax'] ?? 0).toDouble(),
      taxRate: json['tax_rate']?.toDouble(),
      discount: (json['discount'] ?? 0).toDouble(),
    );
  }
}

class Invoice {
  final int id;
  final String orderId;
  final String? orderNumber;
  final int businessId;
  final int integrationId;
  final int invoicingProviderId;
  final String invoiceNumber;
  final String? externalId;
  final String status;
  final double totalAmount;
  final double subtotal;
  final double tax;
  final double discount;
  final String currency;
  final String customerName;
  final String? customerEmail;
  final String? customerDni;
  final String? invoiceUrl;
  final String? pdfUrl;
  final String? xmlUrl;
  final String? cufe;
  final String? issuedAt;
  final String? cancelledAt;
  final String? errorMessage;
  final Map<String, dynamic>? metadata;
  final Map<String, dynamic>? providerResponse;
  final String? providerLogoUrl;
  final String? providerName;
  final String createdAt;
  final String updatedAt;
  final List<InvoiceItem>? items;

  Invoice({
    required this.id,
    required this.orderId,
    this.orderNumber,
    required this.businessId,
    required this.integrationId,
    required this.invoicingProviderId,
    required this.invoiceNumber,
    this.externalId,
    required this.status,
    required this.totalAmount,
    required this.subtotal,
    required this.tax,
    required this.discount,
    required this.currency,
    required this.customerName,
    this.customerEmail,
    this.customerDni,
    this.invoiceUrl,
    this.pdfUrl,
    this.xmlUrl,
    this.cufe,
    this.issuedAt,
    this.cancelledAt,
    this.errorMessage,
    this.metadata,
    this.providerResponse,
    this.providerLogoUrl,
    this.providerName,
    required this.createdAt,
    required this.updatedAt,
    this.items,
  });

  factory Invoice.fromJson(Map<String, dynamic> json) {
    return Invoice(
      id: json['id'] ?? 0,
      orderId: json['order_id']?.toString() ?? '',
      orderNumber: json['order_number'],
      businessId: json['business_id'] ?? 0,
      integrationId: json['integration_id'] ?? 0,
      invoicingProviderId: json['invoicing_provider_id'] ?? 0,
      invoiceNumber: json['invoice_number'] ?? '',
      externalId: json['external_id'],
      status: json['status'] ?? '',
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      subtotal: (json['subtotal'] ?? 0).toDouble(),
      tax: (json['tax'] ?? 0).toDouble(),
      discount: (json['discount'] ?? 0).toDouble(),
      currency: json['currency'] ?? '',
      customerName: json['customer_name'] ?? '',
      customerEmail: json['customer_email'],
      customerDni: json['customer_dni'],
      invoiceUrl: json['invoice_url'],
      pdfUrl: json['pdf_url'],
      xmlUrl: json['xml_url'],
      cufe: json['cufe'],
      issuedAt: json['issued_at'],
      cancelledAt: json['cancelled_at'],
      errorMessage: json['error_message'],
      metadata: json['metadata'] != null ? Map<String, dynamic>.from(json['metadata']) : null,
      providerResponse: json['provider_response'] != null ? Map<String, dynamic>.from(json['provider_response']) : null,
      providerLogoUrl: json['provider_logo_url'],
      providerName: json['provider_name'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      items: (json['items'] as List<dynamic>?)?.map((e) => InvoiceItem.fromJson(e)).toList(),
    );
  }
}

class InvoicingConfig {
  final int id;
  final int businessId;
  final List<int> integrationIds;
  final int? invoicingIntegrationId;
  final int? invoicingProviderId;
  final bool enabled;
  final bool autoInvoice;
  final Map<String, dynamic>? filters;
  final Map<String, dynamic>? config;
  final String? description;
  final String createdAt;
  final String updatedAt;
  final List<String>? integrationNames;
  final String? providerName;
  final String? providerImageUrl;

  InvoicingConfig({
    required this.id,
    required this.businessId,
    required this.integrationIds,
    this.invoicingIntegrationId,
    this.invoicingProviderId,
    required this.enabled,
    required this.autoInvoice,
    this.filters,
    this.config,
    this.description,
    required this.createdAt,
    required this.updatedAt,
    this.integrationNames,
    this.providerName,
    this.providerImageUrl,
  });

  factory InvoicingConfig.fromJson(Map<String, dynamic> json) {
    return InvoicingConfig(
      id: json['id'] ?? 0,
      businessId: json['business_id'] ?? 0,
      integrationIds: (json['integration_ids'] as List<dynamic>?)?.map((e) => e as int).toList() ?? [],
      invoicingIntegrationId: json['invoicing_integration_id'],
      invoicingProviderId: json['invoicing_provider_id'],
      enabled: json['enabled'] ?? false,
      autoInvoice: json['auto_invoice'] ?? false,
      filters: json['filters'] != null ? Map<String, dynamic>.from(json['filters']) : null,
      config: json['config'] != null ? Map<String, dynamic>.from(json['config']) : null,
      description: json['description'],
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
      integrationNames: (json['integration_names'] as List<dynamic>?)?.map((e) => e.toString()).toList(),
      providerName: json['provider_name'],
      providerImageUrl: json['provider_image_url'],
    );
  }
}

class CreditNote {
  final int id;
  final int invoiceId;
  final String creditNoteNumber;
  final String? externalId;
  final double amount;
  final String reason;
  final String noteType;
  final String status;
  final String? noteUrl;
  final String? pdfUrl;
  final String? xmlUrl;
  final String? cufe;
  final String? issuedAt;
  final String createdAt;

  CreditNote({
    required this.id,
    required this.invoiceId,
    required this.creditNoteNumber,
    this.externalId,
    required this.amount,
    required this.reason,
    required this.noteType,
    required this.status,
    this.noteUrl,
    this.pdfUrl,
    this.xmlUrl,
    this.cufe,
    this.issuedAt,
    required this.createdAt,
  });

  factory CreditNote.fromJson(Map<String, dynamic> json) {
    return CreditNote(
      id: json['id'] ?? 0,
      invoiceId: json['invoice_id'] ?? 0,
      creditNoteNumber: json['credit_note_number'] ?? '',
      externalId: json['external_id'],
      amount: (json['amount'] ?? 0).toDouble(),
      reason: json['reason'] ?? '',
      noteType: json['note_type'] ?? '',
      status: json['status'] ?? '',
      noteUrl: json['note_url'],
      pdfUrl: json['pdf_url'],
      xmlUrl: json['xml_url'],
      cufe: json['cufe'],
      issuedAt: json['issued_at'],
      createdAt: json['created_at'] ?? '',
    );
  }
}

class SyncLog {
  final int id;
  final int invoiceId;
  final String operationType;
  final String status;
  final String? errorMessage;
  final String? errorCode;
  final int retryCount;
  final int maxRetries;
  final String? nextRetryAt;
  final String triggeredBy;
  final int? durationMs;
  final String startedAt;
  final String? completedAt;
  final String createdAt;

  SyncLog({
    required this.id,
    required this.invoiceId,
    required this.operationType,
    required this.status,
    this.errorMessage,
    this.errorCode,
    required this.retryCount,
    required this.maxRetries,
    this.nextRetryAt,
    required this.triggeredBy,
    this.durationMs,
    required this.startedAt,
    this.completedAt,
    required this.createdAt,
  });

  factory SyncLog.fromJson(Map<String, dynamic> json) {
    return SyncLog(
      id: json['id'] ?? 0,
      invoiceId: json['invoice_id'] ?? 0,
      operationType: json['operation_type'] ?? '',
      status: json['status'] ?? '',
      errorMessage: json['error_message'],
      errorCode: json['error_code'],
      retryCount: json['retry_count'] ?? 0,
      maxRetries: json['max_retries'] ?? 0,
      nextRetryAt: json['next_retry_at'],
      triggeredBy: json['triggered_by'] ?? '',
      durationMs: json['duration_ms'],
      startedAt: json['started_at'] ?? '',
      completedAt: json['completed_at'],
      createdAt: json['created_at'] ?? '',
    );
  }
}

class InvoicingStats {
  final int totalInvoices;
  final double totalAmount;
  final int pendingInvoices;
  final int failedInvoices;
  final double successRate;
  final String? lastInvoiceDate;

  InvoicingStats({
    required this.totalInvoices,
    required this.totalAmount,
    required this.pendingInvoices,
    required this.failedInvoices,
    required this.successRate,
    this.lastInvoiceDate,
  });

  factory InvoicingStats.fromJson(Map<String, dynamic> json) {
    return InvoicingStats(
      totalInvoices: json['total_invoices'] ?? 0,
      totalAmount: (json['total_amount'] ?? 0).toDouble(),
      pendingInvoices: json['pending_invoices'] ?? 0,
      failedInvoices: json['failed_invoices'] ?? 0,
      successRate: (json['success_rate'] ?? 0).toDouble(),
      lastInvoiceDate: json['last_invoice_date'],
    );
  }
}

// Params & DTOs

class InvoiceFilters {
  final int? businessId;
  final String? orderId;
  final int? integrationId;
  final String? status;
  final String? currency;
  final int? providerId;
  final String? invoiceNumber;
  final String? orderNumber;
  final String? customerName;
  final String? startDate;
  final String? endDate;
  final int? page;
  final int? pageSize;

  InvoiceFilters({
    this.businessId,
    this.orderId,
    this.integrationId,
    this.status,
    this.currency,
    this.providerId,
    this.invoiceNumber,
    this.orderNumber,
    this.customerName,
    this.startDate,
    this.endDate,
    this.page,
    this.pageSize,
  });

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (businessId != null) params['business_id'] = businessId;
    if (orderId != null) params['order_id'] = orderId;
    if (integrationId != null) params['integration_id'] = integrationId;
    if (status != null) params['status'] = status;
    if (currency != null) params['currency'] = currency;
    if (providerId != null) params['provider_id'] = providerId;
    if (invoiceNumber != null) params['invoice_number'] = invoiceNumber;
    if (orderNumber != null) params['order_number'] = orderNumber;
    if (customerName != null) params['customer_name'] = customerName;
    if (startDate != null) params['start_date'] = startDate;
    if (endDate != null) params['end_date'] = endDate;
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    return params;
  }
}

class ConfigFilters {
  final int? businessId;
  final int? integrationId;
  final int? providerId;
  final bool? enabled;
  final int? page;
  final int? pageSize;

  ConfigFilters({this.businessId, this.integrationId, this.providerId, this.enabled, this.page, this.pageSize});

  Map<String, dynamic> toQueryParams() {
    final params = <String, dynamic>{};
    if (businessId != null) params['business_id'] = businessId;
    if (integrationId != null) params['integration_id'] = integrationId;
    if (providerId != null) params['provider_id'] = providerId;
    if (enabled != null) params['enabled'] = enabled;
    if (page != null) params['page'] = page;
    if (pageSize != null) params['page_size'] = pageSize;
    return params;
  }
}

class CreateInvoiceDTO {
  final String orderId;
  final int businessId;
  final int integrationId;

  CreateInvoiceDTO({required this.orderId, required this.businessId, required this.integrationId});

  Map<String, dynamic> toJson() => {
        'order_id': orderId,
        'business_id': businessId,
        'integration_id': integrationId,
      };
}

class CreateConfigDTO {
  final int businessId;
  final List<int> integrationIds;
  final int invoicingIntegrationId;
  final bool? enabled;
  final bool? autoInvoice;
  final Map<String, dynamic>? filters;
  final Map<String, dynamic>? config;
  final String? description;

  CreateConfigDTO({
    required this.businessId,
    required this.integrationIds,
    required this.invoicingIntegrationId,
    this.enabled,
    this.autoInvoice,
    this.filters,
    this.config,
    this.description,
  });

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{
      'business_id': businessId,
      'integration_ids': integrationIds,
      'invoicing_integration_id': invoicingIntegrationId,
    };
    if (enabled != null) json['enabled'] = enabled;
    if (autoInvoice != null) json['auto_invoice'] = autoInvoice;
    if (filters != null) json['filters'] = filters;
    if (config != null) json['config'] = config;
    if (description != null) json['description'] = description;
    return json;
  }
}

class UpdateConfigDTO {
  final bool? enabled;
  final bool? autoInvoice;
  final Map<String, dynamic>? filters;
  final Map<String, dynamic>? config;
  final int? invoicingIntegrationId;
  final List<int>? integrationIds;

  UpdateConfigDTO({this.enabled, this.autoInvoice, this.filters, this.config, this.invoicingIntegrationId, this.integrationIds});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{};
    if (enabled != null) json['enabled'] = enabled;
    if (autoInvoice != null) json['auto_invoice'] = autoInvoice;
    if (filters != null) json['filters'] = filters;
    if (config != null) json['config'] = config;
    if (invoicingIntegrationId != null) json['invoicing_integration_id'] = invoicingIntegrationId;
    if (integrationIds != null) json['integration_ids'] = integrationIds;
    return json;
  }
}

class CreateCreditNoteDTO {
  final int invoiceId;
  final double amount;
  final String reason;
  final String noteType;

  CreateCreditNoteDTO({required this.invoiceId, required this.amount, required this.reason, required this.noteType});

  Map<String, dynamic> toJson() => {
        'invoice_id': invoiceId,
        'amount': amount,
        'reason': reason,
        'note_type': noteType,
      };
}

class BulkCreateInvoicesDTO {
  final List<String> orderIds;
  final int? businessId;

  BulkCreateInvoicesDTO({required this.orderIds, this.businessId});

  Map<String, dynamic> toJson() {
    final json = <String, dynamic>{'order_ids': orderIds};
    if (businessId != null) json['business_id'] = businessId;
    return json;
  }
}
