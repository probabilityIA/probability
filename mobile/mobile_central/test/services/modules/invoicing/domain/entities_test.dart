import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/invoicing/domain/entities.dart';

void main() {
  // =========================================================================
  // InvoiceItem
  // =========================================================================
  group('InvoiceItem', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'invoice_id': 10,
        'product_sku': 'SKU-001',
        'product_name': 'Widget',
        'description': 'A widget item',
        'quantity': 5,
        'unit_price': 100.50,
        'total_price': 502.50,
        'tax': 95.47,
        'tax_rate': 19.0,
        'discount': 10.0,
      };

      final item = InvoiceItem.fromJson(json);

      expect(item.id, 1);
      expect(item.invoiceId, 10);
      expect(item.productSku, 'SKU-001');
      expect(item.productName, 'Widget');
      expect(item.description, 'A widget item');
      expect(item.quantity, 5);
      expect(item.unitPrice, 100.50);
      expect(item.totalPrice, 502.50);
      expect(item.tax, 95.47);
      expect(item.taxRate, 19.0);
      expect(item.discount, 10.0);
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final item = InvoiceItem.fromJson(json);

      expect(item.id, 0);
      expect(item.invoiceId, 0);
      expect(item.productName, '');
      expect(item.quantity, 0);
      expect(item.unitPrice, 0.0);
      expect(item.totalPrice, 0.0);
      expect(item.tax, 0.0);
      expect(item.discount, 0.0);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'invoice_id': 1,
        'product_name': 'Test',
        'quantity': 1,
        'unit_price': 10,
        'total_price': 10,
        'tax': 0,
        'discount': 0,
      };

      final item = InvoiceItem.fromJson(json);

      expect(item.productSku, isNull);
      expect(item.description, isNull);
      expect(item.taxRate, isNull);
    });

    test('fromJson converts numeric values to double', () {
      final json = {
        'id': 1,
        'invoice_id': 1,
        'product_name': 'Test',
        'quantity': 2,
        'unit_price': 100,
        'total_price': 200,
        'tax': 38,
        'discount': 0,
      };

      final item = InvoiceItem.fromJson(json);

      expect(item.unitPrice, isA<double>());
      expect(item.totalPrice, isA<double>());
      expect(item.tax, isA<double>());
      expect(item.discount, isA<double>());
    });
  });

  // =========================================================================
  // Invoice
  // =========================================================================
  group('Invoice', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 42,
        'order_id': 'ord-123',
        'order_number': 'ORD-001',
        'business_id': 1,
        'integration_id': 2,
        'invoicing_provider_id': 3,
        'invoice_number': 'INV-001',
        'external_id': 'EXT-001',
        'status': 'completed',
        'total_amount': 1000.0,
        'subtotal': 840.34,
        'tax': 159.66,
        'discount': 50.0,
        'currency': 'COP',
        'customer_name': 'John Doe',
        'customer_email': 'john@example.com',
        'customer_dni': '123456789',
        'invoice_url': 'https://example.com/inv',
        'pdf_url': 'https://example.com/inv.pdf',
        'xml_url': 'https://example.com/inv.xml',
        'cufe': 'CUFE-123',
        'issued_at': '2026-01-01',
        'cancelled_at': null,
        'error_message': null,
        'metadata': {'key': 'value'},
        'provider_response': {'code': 200},
        'provider_logo_url': 'https://example.com/logo.png',
        'provider_name': 'Softpymes',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
        'items': [
          {
            'id': 1,
            'invoice_id': 42,
            'product_name': 'Item 1',
            'quantity': 2,
            'unit_price': 500,
            'total_price': 1000,
            'tax': 190,
            'discount': 0,
          }
        ],
      };

      final invoice = Invoice.fromJson(json);

      expect(invoice.id, 42);
      expect(invoice.orderId, 'ord-123');
      expect(invoice.orderNumber, 'ORD-001');
      expect(invoice.businessId, 1);
      expect(invoice.integrationId, 2);
      expect(invoice.invoicingProviderId, 3);
      expect(invoice.invoiceNumber, 'INV-001');
      expect(invoice.externalId, 'EXT-001');
      expect(invoice.status, 'completed');
      expect(invoice.totalAmount, 1000.0);
      expect(invoice.subtotal, 840.34);
      expect(invoice.tax, 159.66);
      expect(invoice.discount, 50.0);
      expect(invoice.currency, 'COP');
      expect(invoice.customerName, 'John Doe');
      expect(invoice.customerEmail, 'john@example.com');
      expect(invoice.customerDni, '123456789');
      expect(invoice.invoiceUrl, 'https://example.com/inv');
      expect(invoice.pdfUrl, 'https://example.com/inv.pdf');
      expect(invoice.xmlUrl, 'https://example.com/inv.xml');
      expect(invoice.cufe, 'CUFE-123');
      expect(invoice.issuedAt, '2026-01-01');
      expect(invoice.cancelledAt, isNull);
      expect(invoice.errorMessage, isNull);
      expect(invoice.metadata, {'key': 'value'});
      expect(invoice.providerResponse, {'code': 200});
      expect(invoice.providerLogoUrl, 'https://example.com/logo.png');
      expect(invoice.providerName, 'Softpymes');
      expect(invoice.createdAt, '2026-01-01');
      expect(invoice.updatedAt, '2026-01-02');
      expect(invoice.items, isNotNull);
      expect(invoice.items!.length, 1);
      expect(invoice.items!.first.productName, 'Item 1');
    });

    test('fromJson handles defaults for missing required fields', () {
      final json = <String, dynamic>{};

      final invoice = Invoice.fromJson(json);

      expect(invoice.id, 0);
      expect(invoice.orderId, '');
      expect(invoice.businessId, 0);
      expect(invoice.integrationId, 0);
      expect(invoice.invoicingProviderId, 0);
      expect(invoice.invoiceNumber, '');
      expect(invoice.status, '');
      expect(invoice.totalAmount, 0.0);
      expect(invoice.subtotal, 0.0);
      expect(invoice.tax, 0.0);
      expect(invoice.discount, 0.0);
      expect(invoice.currency, '');
      expect(invoice.customerName, '');
      expect(invoice.createdAt, '');
      expect(invoice.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'order_id': 'o1',
        'business_id': 1,
        'integration_id': 1,
        'invoicing_provider_id': 1,
        'invoice_number': 'INV',
        'status': 'pending',
        'total_amount': 100,
        'subtotal': 84,
        'tax': 16,
        'discount': 0,
        'currency': 'COP',
        'customer_name': 'Test',
        'created_at': '',
        'updated_at': '',
      };

      final invoice = Invoice.fromJson(json);

      expect(invoice.orderNumber, isNull);
      expect(invoice.externalId, isNull);
      expect(invoice.customerEmail, isNull);
      expect(invoice.customerDni, isNull);
      expect(invoice.invoiceUrl, isNull);
      expect(invoice.pdfUrl, isNull);
      expect(invoice.xmlUrl, isNull);
      expect(invoice.cufe, isNull);
      expect(invoice.issuedAt, isNull);
      expect(invoice.cancelledAt, isNull);
      expect(invoice.errorMessage, isNull);
      expect(invoice.metadata, isNull);
      expect(invoice.providerResponse, isNull);
      expect(invoice.providerLogoUrl, isNull);
      expect(invoice.providerName, isNull);
      expect(invoice.items, isNull);
    });

    test('fromJson converts order_id to string', () {
      final json = {
        'id': 1,
        'order_id': 12345,
        'business_id': 1,
        'integration_id': 1,
        'invoicing_provider_id': 1,
        'invoice_number': 'INV',
        'status': 'pending',
        'total_amount': 0,
        'subtotal': 0,
        'tax': 0,
        'discount': 0,
        'currency': 'COP',
        'customer_name': 'Test',
        'created_at': '',
        'updated_at': '',
      };

      final invoice = Invoice.fromJson(json);
      expect(invoice.orderId, '12345');
    });
  });

  // =========================================================================
  // InvoicingConfig
  // =========================================================================
  group('InvoicingConfig', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'integration_ids': [1, 2, 3],
        'invoicing_integration_id': 10,
        'invoicing_provider_id': 20,
        'enabled': true,
        'auto_invoice': true,
        'filters': {'min_amount': 1000},
        'config': {'key': 'value'},
        'description': 'Test config',
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
        'integration_names': ['Shopify', 'Amazon'],
        'provider_name': 'Softpymes',
        'provider_image_url': 'https://example.com/logo.png',
      };

      final config = InvoicingConfig.fromJson(json);

      expect(config.id, 1);
      expect(config.businessId, 5);
      expect(config.integrationIds, [1, 2, 3]);
      expect(config.invoicingIntegrationId, 10);
      expect(config.invoicingProviderId, 20);
      expect(config.enabled, true);
      expect(config.autoInvoice, true);
      expect(config.filters, {'min_amount': 1000});
      expect(config.config, {'key': 'value'});
      expect(config.description, 'Test config');
      expect(config.createdAt, '2026-01-01');
      expect(config.updatedAt, '2026-01-02');
      expect(config.integrationNames, ['Shopify', 'Amazon']);
      expect(config.providerName, 'Softpymes');
      expect(config.providerImageUrl, 'https://example.com/logo.png');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final config = InvoicingConfig.fromJson(json);

      expect(config.id, 0);
      expect(config.businessId, 0);
      expect(config.integrationIds, isEmpty);
      expect(config.enabled, false);
      expect(config.autoInvoice, false);
      expect(config.createdAt, '');
      expect(config.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'integration_ids': <int>[],
        'enabled': true,
        'auto_invoice': false,
        'created_at': '',
        'updated_at': '',
      };

      final config = InvoicingConfig.fromJson(json);

      expect(config.invoicingIntegrationId, isNull);
      expect(config.invoicingProviderId, isNull);
      expect(config.filters, isNull);
      expect(config.config, isNull);
      expect(config.description, isNull);
      expect(config.integrationNames, isNull);
      expect(config.providerName, isNull);
      expect(config.providerImageUrl, isNull);
    });
  });

  // =========================================================================
  // CreditNote
  // =========================================================================
  group('CreditNote', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'invoice_id': 42,
        'credit_note_number': 'CN-001',
        'external_id': 'EXT-CN-001',
        'amount': 500.0,
        'reason': 'Product return',
        'note_type': 'full',
        'status': 'completed',
        'note_url': 'https://example.com/cn',
        'pdf_url': 'https://example.com/cn.pdf',
        'xml_url': 'https://example.com/cn.xml',
        'cufe': 'CUFE-CN-001',
        'issued_at': '2026-02-01',
        'created_at': '2026-02-01',
      };

      final note = CreditNote.fromJson(json);

      expect(note.id, 1);
      expect(note.invoiceId, 42);
      expect(note.creditNoteNumber, 'CN-001');
      expect(note.externalId, 'EXT-CN-001');
      expect(note.amount, 500.0);
      expect(note.reason, 'Product return');
      expect(note.noteType, 'full');
      expect(note.status, 'completed');
      expect(note.noteUrl, 'https://example.com/cn');
      expect(note.pdfUrl, 'https://example.com/cn.pdf');
      expect(note.xmlUrl, 'https://example.com/cn.xml');
      expect(note.cufe, 'CUFE-CN-001');
      expect(note.issuedAt, '2026-02-01');
      expect(note.createdAt, '2026-02-01');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final note = CreditNote.fromJson(json);

      expect(note.id, 0);
      expect(note.invoiceId, 0);
      expect(note.creditNoteNumber, '');
      expect(note.amount, 0.0);
      expect(note.reason, '');
      expect(note.noteType, '');
      expect(note.status, '');
      expect(note.createdAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'invoice_id': 1,
        'credit_note_number': 'CN',
        'amount': 100,
        'reason': 'r',
        'note_type': 'full',
        'status': 'pending',
        'created_at': '',
      };

      final note = CreditNote.fromJson(json);

      expect(note.externalId, isNull);
      expect(note.noteUrl, isNull);
      expect(note.pdfUrl, isNull);
      expect(note.xmlUrl, isNull);
      expect(note.cufe, isNull);
      expect(note.issuedAt, isNull);
    });
  });

  // =========================================================================
  // SyncLog
  // =========================================================================
  group('SyncLog', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'invoice_id': 10,
        'operation_type': 'create',
        'status': 'completed',
        'error_message': null,
        'error_code': null,
        'retry_count': 0,
        'max_retries': 3,
        'next_retry_at': null,
        'triggered_by': 'auto',
        'duration_ms': 250,
        'started_at': '2026-01-01T10:00:00Z',
        'completed_at': '2026-01-01T10:00:01Z',
        'created_at': '2026-01-01T10:00:00Z',
      };

      final log = SyncLog.fromJson(json);

      expect(log.id, 1);
      expect(log.invoiceId, 10);
      expect(log.operationType, 'create');
      expect(log.status, 'completed');
      expect(log.errorMessage, isNull);
      expect(log.errorCode, isNull);
      expect(log.retryCount, 0);
      expect(log.maxRetries, 3);
      expect(log.nextRetryAt, isNull);
      expect(log.triggeredBy, 'auto');
      expect(log.durationMs, 250);
      expect(log.startedAt, '2026-01-01T10:00:00Z');
      expect(log.completedAt, '2026-01-01T10:00:01Z');
      expect(log.createdAt, '2026-01-01T10:00:00Z');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final log = SyncLog.fromJson(json);

      expect(log.id, 0);
      expect(log.invoiceId, 0);
      expect(log.operationType, '');
      expect(log.status, '');
      expect(log.retryCount, 0);
      expect(log.maxRetries, 0);
      expect(log.triggeredBy, '');
      expect(log.startedAt, '');
      expect(log.createdAt, '');
    });
  });

  // =========================================================================
  // InvoicingStats
  // =========================================================================
  group('InvoicingStats', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'total_invoices': 100,
        'total_amount': 5000000.0,
        'pending_invoices': 5,
        'failed_invoices': 2,
        'success_rate': 93.0,
        'last_invoice_date': '2026-03-01',
      };

      final stats = InvoicingStats.fromJson(json);

      expect(stats.totalInvoices, 100);
      expect(stats.totalAmount, 5000000.0);
      expect(stats.pendingInvoices, 5);
      expect(stats.failedInvoices, 2);
      expect(stats.successRate, 93.0);
      expect(stats.lastInvoiceDate, '2026-03-01');
    });

    test('fromJson handles defaults for missing fields', () {
      final json = <String, dynamic>{};

      final stats = InvoicingStats.fromJson(json);

      expect(stats.totalInvoices, 0);
      expect(stats.totalAmount, 0.0);
      expect(stats.pendingInvoices, 0);
      expect(stats.failedInvoices, 0);
      expect(stats.successRate, 0.0);
      expect(stats.lastInvoiceDate, isNull);
    });
  });

  // =========================================================================
  // InvoiceFilters
  // =========================================================================
  group('InvoiceFilters', () {
    test('toQueryParams includes all set fields', () {
      final filters = InvoiceFilters(
        businessId: 1,
        orderId: 'ord-1',
        integrationId: 2,
        status: 'completed',
        currency: 'COP',
        providerId: 3,
        invoiceNumber: 'INV-001',
        orderNumber: 'ORD-001',
        customerName: 'John',
        startDate: '2026-01-01',
        endDate: '2026-12-31',
        page: 1,
        pageSize: 20,
      );

      final query = filters.toQueryParams();

      expect(query['business_id'], 1);
      expect(query['order_id'], 'ord-1');
      expect(query['integration_id'], 2);
      expect(query['status'], 'completed');
      expect(query['currency'], 'COP');
      expect(query['provider_id'], 3);
      expect(query['invoice_number'], 'INV-001');
      expect(query['order_number'], 'ORD-001');
      expect(query['customer_name'], 'John');
      expect(query['start_date'], '2026-01-01');
      expect(query['end_date'], '2026-12-31');
      expect(query['page'], 1);
      expect(query['page_size'], 20);
    });

    test('toQueryParams omits null fields', () {
      final filters = InvoiceFilters();
      final query = filters.toQueryParams();
      expect(query, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final filters = InvoiceFilters(status: 'pending', page: 1);
      final query = filters.toQueryParams();

      expect(query.length, 2);
      expect(query['status'], 'pending');
      expect(query['page'], 1);
    });
  });

  // =========================================================================
  // ConfigFilters
  // =========================================================================
  group('ConfigFilters', () {
    test('toQueryParams includes all set fields', () {
      final filters = ConfigFilters(
        businessId: 1,
        integrationId: 2,
        providerId: 3,
        enabled: true,
        page: 1,
        pageSize: 10,
      );

      final query = filters.toQueryParams();

      expect(query['business_id'], 1);
      expect(query['integration_id'], 2);
      expect(query['provider_id'], 3);
      expect(query['enabled'], true);
      expect(query['page'], 1);
      expect(query['page_size'], 10);
    });

    test('toQueryParams omits null fields', () {
      final filters = ConfigFilters();
      final query = filters.toQueryParams();
      expect(query, isEmpty);
    });
  });

  // =========================================================================
  // CreateInvoiceDTO
  // =========================================================================
  group('CreateInvoiceDTO', () {
    test('toJson includes all fields', () {
      final dto = CreateInvoiceDTO(
        orderId: 'ord-1',
        businessId: 5,
        integrationId: 10,
      );

      final json = dto.toJson();

      expect(json['order_id'], 'ord-1');
      expect(json['business_id'], 5);
      expect(json['integration_id'], 10);
    });
  });

  // =========================================================================
  // CreateConfigDTO
  // =========================================================================
  group('CreateConfigDTO', () {
    test('toJson includes all required fields', () {
      final dto = CreateConfigDTO(
        businessId: 1,
        integrationIds: [1, 2],
        invoicingIntegrationId: 5,
      );

      final json = dto.toJson();

      expect(json['business_id'], 1);
      expect(json['integration_ids'], [1, 2]);
      expect(json['invoicing_integration_id'], 5);
    });

    test('toJson includes optional fields when set', () {
      final dto = CreateConfigDTO(
        businessId: 1,
        integrationIds: [1],
        invoicingIntegrationId: 5,
        enabled: true,
        autoInvoice: false,
        filters: {'min_amount': 1000},
        config: {'key': 'val'},
        description: 'Test',
      );

      final json = dto.toJson();

      expect(json['enabled'], true);
      expect(json['auto_invoice'], false);
      expect(json['filters'], {'min_amount': 1000});
      expect(json['config'], {'key': 'val'});
      expect(json['description'], 'Test');
    });

    test('toJson omits null optional fields', () {
      final dto = CreateConfigDTO(
        businessId: 1,
        integrationIds: [1],
        invoicingIntegrationId: 5,
      );

      final json = dto.toJson();

      expect(json.containsKey('enabled'), false);
      expect(json.containsKey('auto_invoice'), false);
      expect(json.containsKey('filters'), false);
      expect(json.containsKey('config'), false);
      expect(json.containsKey('description'), false);
    });
  });

  // =========================================================================
  // UpdateConfigDTO
  // =========================================================================
  group('UpdateConfigDTO', () {
    test('toJson includes all set fields', () {
      final dto = UpdateConfigDTO(
        enabled: true,
        autoInvoice: false,
        filters: {'min': 100},
        config: {'k': 'v'},
        invoicingIntegrationId: 5,
        integrationIds: [1, 2],
      );

      final json = dto.toJson();

      expect(json['enabled'], true);
      expect(json['auto_invoice'], false);
      expect(json['filters'], {'min': 100});
      expect(json['config'], {'k': 'v'});
      expect(json['invoicing_integration_id'], 5);
      expect(json['integration_ids'], [1, 2]);
    });

    test('toJson returns empty map when all fields null', () {
      final dto = UpdateConfigDTO();
      final json = dto.toJson();
      expect(json, isEmpty);
    });
  });

  // =========================================================================
  // CreateCreditNoteDTO
  // =========================================================================
  group('CreateCreditNoteDTO', () {
    test('toJson includes all fields', () {
      final dto = CreateCreditNoteDTO(
        invoiceId: 42,
        amount: 500.0,
        reason: 'return',
        noteType: 'full',
      );

      final json = dto.toJson();

      expect(json['invoice_id'], 42);
      expect(json['amount'], 500.0);
      expect(json['reason'], 'return');
      expect(json['note_type'], 'full');
    });
  });

  // =========================================================================
  // BulkCreateInvoicesDTO
  // =========================================================================
  group('BulkCreateInvoicesDTO', () {
    test('toJson includes order_ids always', () {
      final dto = BulkCreateInvoicesDTO(orderIds: ['o1', 'o2']);

      final json = dto.toJson();

      expect(json['order_ids'], ['o1', 'o2']);
      expect(json.containsKey('business_id'), false);
    });

    test('toJson includes businessId when set', () {
      final dto = BulkCreateInvoicesDTO(orderIds: ['o1'], businessId: 5);

      final json = dto.toJson();

      expect(json['order_ids'], ['o1']);
      expect(json['business_id'], 5);
    });

    test('toJson omits null businessId', () {
      final dto = BulkCreateInvoicesDTO(orderIds: ['o1']);

      final json = dto.toJson();

      expect(json.containsKey('business_id'), false);
    });
  });
}
