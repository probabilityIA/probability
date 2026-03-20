import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/orders/domain/entities.dart';

void main() {
  group('OrderStatusInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'code': 'pending',
        'name': 'Pending',
        'description': 'Order is pending',
        'category': 'open',
        'color': '#FFA500',
      };

      final status = OrderStatusInfo.fromJson(json);

      expect(status.id, 1);
      expect(status.code, 'pending');
      expect(status.name, 'Pending');
      expect(status.description, 'Order is pending');
      expect(status.category, 'open');
      expect(status.color, '#FFA500');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 2,
        'code': 'shipped',
        'name': 'Shipped',
      };

      final status = OrderStatusInfo.fromJson(json);

      expect(status.id, 2);
      expect(status.code, 'shipped');
      expect(status.name, 'Shipped');
      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final status = OrderStatusInfo.fromJson(json);

      expect(status.id, 0);
      expect(status.code, '');
      expect(status.name, '');
    });
  });

  group('PaymentStatusInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 5,
        'code': 'paid',
        'name': 'Paid',
        'description': 'Payment received',
        'category': 'completed',
        'color': '#00FF00',
      };

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.id, 5);
      expect(status.code, 'paid');
      expect(status.name, 'Paid');
      expect(status.description, 'Payment received');
      expect(status.category, 'completed');
      expect(status.color, '#00FF00');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'code': 'pending',
        'name': 'Pending',
      };

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final status = PaymentStatusInfo.fromJson(json);

      expect(status.id, 0);
      expect(status.code, '');
      expect(status.name, '');
    });
  });

  group('FulfillmentStatusInfo', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 3,
        'code': 'fulfilled',
        'name': 'Fulfilled',
        'description': 'Order fulfilled',
        'category': 'done',
        'color': '#0000FF',
      };

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.id, 3);
      expect(status.code, 'fulfilled');
      expect(status.name, 'Fulfilled');
      expect(status.description, 'Order fulfilled');
      expect(status.category, 'done');
      expect(status.color, '#0000FF');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'code': 'unfulfilled',
        'name': 'Unfulfilled',
      };

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.description, isNull);
      expect(status.category, isNull);
      expect(status.color, isNull);
    });

    test('fromJson defaults to safe values when fields are missing', () {
      final json = <String, dynamic>{};

      final status = FulfillmentStatusInfo.fromJson(json);

      expect(status.id, 0);
      expect(status.code, '');
      expect(status.name, '');
    });
  });

  group('Order', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 123,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
        'deleted_at': '2026-01-03T00:00:00Z',
        'business_id': 5,
        'integration_id': 10,
        'integration_type': 'shopify',
        'integration_logo_url': 'https://example.com/logo.png',
        'integration_name': 'My Shopify',
        'platform': 'shopify',
        'external_id': 'ext-123',
        'order_number': 'ORD-001',
        'internal_number': 'INT-001',
        'subtotal': 100.50,
        'tax': 19.0,
        'discount': 5.0,
        'shipping_cost': 10.0,
        'shipping_discount': 2.0,
        'shipping_discount_presentment': 2.5,
        'total_amount': 122.50,
        'currency': 'COP',
        'cod_total': 50.0,
        'subtotal_presentment': 100.0,
        'tax_presentment': 19.0,
        'discount_presentment': 5.0,
        'shipping_cost_presentment': 10.0,
        'total_amount_presentment': 124.0,
        'currency_presentment': 'USD',
        'customer_id': 42,
        'customer_name': 'John Doe',
        'customer_first_name': 'John',
        'customer_last_name': 'Doe',
        'customer_email': 'john@example.com',
        'customer_phone': '+573001234567',
        'customer_dni': '123456789',
        'shipping_street': 'Calle 10 #5-20',
        'shipping_city': 'Bogota',
        'shipping_state': 'Cundinamarca',
        'shipping_country': 'CO',
        'shipping_postal_code': '110111',
        'shipping_house': 'Apt 301',
        'shipping_barrio': 'Chapinero',
        'shipping_lat': 4.6097,
        'shipping_lng': -74.0817,
        'payment_method_id': 1,
        'is_paid': true,
        'paid_at': '2026-01-01T12:00:00Z',
        'tracking_number': 'TRK-001',
        'tracking_link': 'https://tracking.com/TRK-001',
        'guide_id': 'G-001',
        'guide_link': 'https://guide.com/G-001',
        'delivery_date': '2026-01-05',
        'delivered_at': '2026-01-05T10:00:00Z',
        'delivery_probability': 0.85,
        'warehouse_id': 3,
        'warehouse_name': 'Main Warehouse',
        'driver_id': 7,
        'driver_name': 'Carlos',
        'is_last_mile': true,
        'weight': 2.5,
        'height': 10.0,
        'width': 20.0,
        'length': 30.0,
        'boxes': '2',
        'order_type_id': 1,
        'order_type_name': 'Standard',
        'status': 'confirmed',
        'original_status': 'paid',
        'status_id': 2,
        'order_status': {
          'id': 2,
          'code': 'confirmed',
          'name': 'Confirmed',
        },
        'payment_status_id': 1,
        'fulfillment_status_id': 3,
        'payment_status': {
          'id': 1,
          'code': 'paid',
          'name': 'Paid',
        },
        'fulfillment_status': {
          'id': 3,
          'code': 'shipped',
          'name': 'Shipped',
        },
        'notes': 'Handle with care',
        'coupon': 'SAVE10',
        'approved': true,
        'user_id': 99,
        'user_name': 'admin',
        'is_confirmed': true,
        'novelty': 'None',
        'invoiceable': true,
        'invoice_url': 'https://invoice.com/1',
        'invoice_id': 'INV-001',
        'invoice_provider': 'softpymes',
        'invoice_status': 'issued',
        'order_status_url': 'https://status.com/1',
        'items': [
          {'name': 'Item 1'}
        ],
        'order_items': [
          {'sku': 'SKU-1'}
        ],
        'metadata': {'source': 'api'},
        'financial_details': {'fee': 1.5},
        'shipping_details': {'carrier': 'DHL'},
        'payment_details': {'method': 'credit_card'},
        'fulfillment_details': {'status': 'shipped'},
        'occurred_at': '2026-01-01T00:00:00Z',
        'imported_at': '2026-01-01T01:00:00Z',
        'negative_factors': ['late_payment', 'address_issue'],
      };

      final order = Order.fromJson(json);

      expect(order.id, '123');
      expect(order.createdAt, '2026-01-01T00:00:00Z');
      expect(order.updatedAt, '2026-01-02T00:00:00Z');
      expect(order.deletedAt, '2026-01-03T00:00:00Z');
      expect(order.businessId, 5);
      expect(order.integrationId, 10);
      expect(order.integrationType, 'shopify');
      expect(order.integrationLogoUrl, 'https://example.com/logo.png');
      expect(order.integrationName, 'My Shopify');
      expect(order.platform, 'shopify');
      expect(order.externalId, 'ext-123');
      expect(order.orderNumber, 'ORD-001');
      expect(order.internalNumber, 'INT-001');
      expect(order.subtotal, 100.50);
      expect(order.tax, 19.0);
      expect(order.discount, 5.0);
      expect(order.shippingCost, 10.0);
      expect(order.shippingDiscount, 2.0);
      expect(order.shippingDiscountPresentment, 2.5);
      expect(order.totalAmount, 122.50);
      expect(order.currency, 'COP');
      expect(order.codTotal, 50.0);
      expect(order.subtotalPresentment, 100.0);
      expect(order.taxPresentment, 19.0);
      expect(order.discountPresentment, 5.0);
      expect(order.shippingCostPresentment, 10.0);
      expect(order.totalAmountPresentment, 124.0);
      expect(order.currencyPresentment, 'USD');
      expect(order.customerId, 42);
      expect(order.customerName, 'John Doe');
      expect(order.customerFirstName, 'John');
      expect(order.customerLastName, 'Doe');
      expect(order.customerEmail, 'john@example.com');
      expect(order.customerPhone, '+573001234567');
      expect(order.customerDni, '123456789');
      expect(order.shippingStreet, 'Calle 10 #5-20');
      expect(order.shippingCity, 'Bogota');
      expect(order.shippingState, 'Cundinamarca');
      expect(order.shippingCountry, 'CO');
      expect(order.shippingPostalCode, '110111');
      expect(order.shippingHouse, 'Apt 301');
      expect(order.shippingBarrio, 'Chapinero');
      expect(order.shippingLat, 4.6097);
      expect(order.shippingLng, -74.0817);
      expect(order.paymentMethodId, 1);
      expect(order.isPaid, true);
      expect(order.paidAt, '2026-01-01T12:00:00Z');
      expect(order.trackingNumber, 'TRK-001');
      expect(order.trackingLink, 'https://tracking.com/TRK-001');
      expect(order.guideId, 'G-001');
      expect(order.guideLink, 'https://guide.com/G-001');
      expect(order.deliveryDate, '2026-01-05');
      expect(order.deliveredAt, '2026-01-05T10:00:00Z');
      expect(order.deliveryProbability, 0.85);
      expect(order.warehouseId, 3);
      expect(order.warehouseName, 'Main Warehouse');
      expect(order.driverId, 7);
      expect(order.driverName, 'Carlos');
      expect(order.isLastMile, true);
      expect(order.weight, 2.5);
      expect(order.height, 10.0);
      expect(order.width, 20.0);
      expect(order.length, 30.0);
      expect(order.boxes, '2');
      expect(order.orderTypeId, 1);
      expect(order.orderTypeName, 'Standard');
      expect(order.status, 'confirmed');
      expect(order.originalStatus, 'paid');
      expect(order.statusId, 2);
      expect(order.orderStatus, isNotNull);
      expect(order.orderStatus!.code, 'confirmed');
      expect(order.paymentStatusId, 1);
      expect(order.fulfillmentStatusId, 3);
      expect(order.paymentStatus, isNotNull);
      expect(order.paymentStatus!.code, 'paid');
      expect(order.fulfillmentStatus, isNotNull);
      expect(order.fulfillmentStatus!.code, 'shipped');
      expect(order.notes, 'Handle with care');
      expect(order.coupon, 'SAVE10');
      expect(order.approved, true);
      expect(order.userId, 99);
      expect(order.userName, 'admin');
      expect(order.isConfirmed, true);
      expect(order.novelty, 'None');
      expect(order.invoiceable, true);
      expect(order.invoiceUrl, 'https://invoice.com/1');
      expect(order.invoiceId, 'INV-001');
      expect(order.invoiceProvider, 'softpymes');
      expect(order.invoiceStatus, 'issued');
      expect(order.orderStatusUrl, 'https://status.com/1');
      expect(order.items, isNotNull);
      expect(order.orderItems, isNotNull);
      expect(order.metadata, isNotNull);
      expect(order.financialDetails, isNotNull);
      expect(order.shippingDetails, isNotNull);
      expect(order.paymentDetails, isNotNull);
      expect(order.fulfillmentDetails, isNotNull);
      expect(order.occurredAt, '2026-01-01T00:00:00Z');
      expect(order.importedAt, '2026-01-01T01:00:00Z');
      expect(order.negativeFactors, ['late_payment', 'address_issue']);
    });

    test('fromJson handles null optional fields and defaults', () {
      final json = <String, dynamic>{};

      final order = Order.fromJson(json);

      expect(order.id, '');
      expect(order.createdAt, '');
      expect(order.updatedAt, '');
      expect(order.deletedAt, isNull);
      expect(order.businessId, isNull);
      expect(order.integrationId, 0);
      expect(order.integrationType, '');
      expect(order.integrationLogoUrl, isNull);
      expect(order.integrationName, isNull);
      expect(order.platform, '');
      expect(order.externalId, '');
      expect(order.orderNumber, '');
      expect(order.internalNumber, '');
      expect(order.subtotal, 0.0);
      expect(order.tax, 0.0);
      expect(order.discount, 0.0);
      expect(order.shippingCost, 0.0);
      expect(order.shippingDiscount, isNull);
      expect(order.totalAmount, 0.0);
      expect(order.currency, '');
      expect(order.codTotal, isNull);
      expect(order.customerId, isNull);
      expect(order.customerName, '');
      expect(order.customerEmail, '');
      expect(order.customerPhone, '');
      expect(order.customerDni, '');
      expect(order.shippingStreet, '');
      expect(order.shippingCity, '');
      expect(order.shippingState, '');
      expect(order.shippingCountry, '');
      expect(order.shippingPostalCode, '');
      expect(order.paymentMethodId, 0);
      expect(order.isPaid, false);
      expect(order.warehouseName, '');
      expect(order.driverName, '');
      expect(order.isLastMile, false);
      expect(order.orderTypeName, '');
      expect(order.status, '');
      expect(order.originalStatus, '');
      expect(order.orderStatus, isNull);
      expect(order.paymentStatus, isNull);
      expect(order.fulfillmentStatus, isNull);
      expect(order.userName, '');
      expect(order.invoiceable, false);
      expect(order.occurredAt, '');
      expect(order.importedAt, '');
      expect(order.negativeFactors, isNull);
    });

    test('fromJson converts id to string', () {
      final json = {'id': 456};

      final order = Order.fromJson(json);

      expect(order.id, '456');
    });

    test('fromJson handles null id', () {
      final json = {'id': null};

      final order = Order.fromJson(json);

      expect(order.id, '');
    });

    test('fromJson parses nested order_status', () {
      final json = {
        'order_status': {
          'id': 5,
          'code': 'processing',
          'name': 'Processing',
          'description': 'Being processed',
          'category': 'active',
          'color': '#FFFF00',
        },
      };

      final order = Order.fromJson(json);

      expect(order.orderStatus, isNotNull);
      expect(order.orderStatus!.id, 5);
      expect(order.orderStatus!.code, 'processing');
      expect(order.orderStatus!.name, 'Processing');
    });

    test('fromJson handles null nested statuses', () {
      final json = {
        'order_status': null,
        'payment_status': null,
        'fulfillment_status': null,
      };

      final order = Order.fromJson(json);

      expect(order.orderStatus, isNull);
      expect(order.paymentStatus, isNull);
      expect(order.fulfillmentStatus, isNull);
    });

    test('fromJson handles null negative_factors', () {
      final json = {'negative_factors': null};

      final order = Order.fromJson(json);

      expect(order.negativeFactors, isNull);
    });

    test('fromJson handles empty negative_factors list', () {
      final json = {'negative_factors': []};

      final order = Order.fromJson(json);

      expect(order.negativeFactors, isEmpty);
    });
  });

  group('GetOrdersParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetOrdersParams(
        page: 2,
        pageSize: 25,
        businessId: 5,
        integrationId: 10,
        integrationType: 'shopify',
        status: 'confirmed',
        customerEmail: 'john@example.com',
        customerPhone: '+573001234567',
        orderNumber: 'ORD-001',
        internalNumber: 'INT-001',
        platform: 'shopify',
        currencyPresentment: 'USD',
        isPaid: true,
        isCod: false,
        paymentStatusId: 1,
        fulfillmentStatusId: 2,
        warehouseId: 3,
        driverId: 7,
        startDate: '2026-01-01',
        endDate: '2026-01-31',
        invoiceStatus: 'issued',
        sortBy: 'created_at',
        sortOrder: 'desc',
      );

      final query = params.toQueryParams();

      expect(query['page'], 2);
      expect(query['page_size'], 25);
      expect(query['business_id'], 5);
      expect(query['integration_id'], 10);
      expect(query['integration_type'], 'shopify');
      expect(query['status'], 'confirmed');
      expect(query['customer_email'], 'john@example.com');
      expect(query['customer_phone'], '+573001234567');
      expect(query['order_number'], 'ORD-001');
      expect(query['internal_number'], 'INT-001');
      expect(query['platform'], 'shopify');
      expect(query['currency_presentment'], 'USD');
      expect(query['is_paid'], true);
      expect(query['is_cod'], false);
      expect(query['payment_status_id'], 1);
      expect(query['fulfillment_status_id'], 2);
      expect(query['warehouse_id'], 3);
      expect(query['driver_id'], 7);
      expect(query['start_date'], '2026-01-01');
      expect(query['end_date'], '2026-01-31');
      expect(query['invoice_status'], 'issued');
      expect(query['sort_by'], 'created_at');
      expect(query['sort_order'], 'desc');
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetOrdersParams();

      final query = params.toQueryParams();

      expect(query, isEmpty);
    });

    test('toQueryParams includes only provided fields', () {
      final params = GetOrdersParams(page: 1, status: 'pending');

      final query = params.toQueryParams();

      expect(query.length, 2);
      expect(query['page'], 1);
      expect(query['status'], 'pending');
    });
  });

  group('CreateOrderDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      final json = dto.toJson();

      expect(json['integration_id'], 10);
      expect(json['integration_type'], 'shopify');
      expect(json['platform'], 'shopify');
      expect(json['external_id'], 'ext-1');
      expect(json['subtotal'], 100.0);
      expect(json['total_amount'], 119.0);
      expect(json['payment_method_id'], 1);
    });

    test('toJson includes all optional fields when provided', () {
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
        businessId: 5,
        orderNumber: 'ORD-001',
        customerName: 'John',
        customerEmail: 'john@test.com',
        customerPhone: '+573001234567',
        customerDni: '123456789',
        tax: 19.0,
        discount: 5.0,
        shippingCost: 10.0,
        currency: 'COP',
        codTotal: 50.0,
        shippingStreet: 'Calle 10',
        shippingCity: 'Bogota',
        shippingState: 'Cundinamarca',
        shippingCountry: 'CO',
        shippingPostalCode: '110111',
        isPaid: true,
        warehouseId: 3,
        status: 'pending',
        notes: 'Test note',
        invoiceable: true,
        items: [{'name': 'Item 1'}],
        metadata: {'key': 'value'},
      );

      final json = dto.toJson();

      expect(json['business_id'], 5);
      expect(json['order_number'], 'ORD-001');
      expect(json['customer_name'], 'John');
      expect(json['customer_email'], 'john@test.com');
      expect(json['customer_phone'], '+573001234567');
      expect(json['customer_dni'], '123456789');
      expect(json['tax'], 19.0);
      expect(json['discount'], 5.0);
      expect(json['shipping_cost'], 10.0);
      expect(json['currency'], 'COP');
      expect(json['cod_total'], 50.0);
      expect(json['shipping_street'], 'Calle 10');
      expect(json['shipping_city'], 'Bogota');
      expect(json['shipping_state'], 'Cundinamarca');
      expect(json['shipping_country'], 'CO');
      expect(json['shipping_postal_code'], '110111');
      expect(json['is_paid'], true);
      expect(json['warehouse_id'], 3);
      expect(json['status'], 'pending');
      expect(json['notes'], 'Test note');
      expect(json['invoiceable'], true);
      expect(json['items'], isNotNull);
      expect(json['metadata'], isNotNull);
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateOrderDTO(
        integrationId: 10,
        integrationType: 'shopify',
        platform: 'shopify',
        externalId: 'ext-1',
        subtotal: 100.0,
        totalAmount: 119.0,
        paymentMethodId: 1,
      );

      final json = dto.toJson();

      expect(json.containsKey('business_id'), false);
      expect(json.containsKey('order_number'), false);
      expect(json.containsKey('customer_name'), false);
      expect(json.containsKey('tax'), false);
      expect(json.containsKey('notes'), false);
      expect(json.containsKey('items'), false);
      expect(json.containsKey('metadata'), false);
    });
  });

  group('UpdateOrderDTO', () {
    test('toJson returns empty map when all fields are null', () {
      final dto = UpdateOrderDTO();

      final json = dto.toJson();

      expect(json, isEmpty);
    });

    test('toJson includes only provided fields', () {
      final dto = UpdateOrderDTO(
        status: 'shipped',
        trackingNumber: 'TRK-001',
        isPaid: true,
      );

      final json = dto.toJson();

      expect(json.length, 3);
      expect(json['status'], 'shipped');
      expect(json['tracking_number'], 'TRK-001');
      expect(json['is_paid'], true);
    });

    test('toJson includes all fields when all provided', () {
      final dto = UpdateOrderDTO(
        subtotal: 100.0,
        tax: 19.0,
        discount: 5.0,
        shippingCost: 10.0,
        totalAmount: 124.0,
        currency: 'COP',
        codTotal: 50.0,
        customerName: 'John',
        customerEmail: 'john@test.com',
        customerPhone: '+573001234567',
        customerDni: '123456789',
        shippingStreet: 'Calle 10',
        shippingCity: 'Bogota',
        shippingState: 'Cundinamarca',
        shippingCountry: 'CO',
        shippingPostalCode: '110111',
        paymentMethodId: 1,
        isPaid: true,
        trackingNumber: 'TRK-001',
        trackingLink: 'https://tracking.com/TRK-001',
        warehouseId: 3,
        warehouseName: 'Main',
        driverId: 7,
        driverName: 'Carlos',
        status: 'shipped',
        statusId: 2,
        paymentStatusId: 1,
        fulfillmentStatusId: 3,
        notes: 'Updated note',
        isConfirmed: true,
        confirmationStatus: 'confirmed',
        novelty: 'None',
        invoiceable: true,
        items: [{'name': 'Item 1'}],
        metadata: {'key': 'value'},
      );

      final json = dto.toJson();

      expect(json['subtotal'], 100.0);
      expect(json['tax'], 19.0);
      expect(json['discount'], 5.0);
      expect(json['shipping_cost'], 10.0);
      expect(json['total_amount'], 124.0);
      expect(json['currency'], 'COP');
      expect(json['cod_total'], 50.0);
      expect(json['customer_name'], 'John');
      expect(json['customer_email'], 'john@test.com');
      expect(json['customer_phone'], '+573001234567');
      expect(json['customer_dni'], '123456789');
      expect(json['shipping_street'], 'Calle 10');
      expect(json['shipping_city'], 'Bogota');
      expect(json['shipping_state'], 'Cundinamarca');
      expect(json['shipping_country'], 'CO');
      expect(json['shipping_postal_code'], '110111');
      expect(json['payment_method_id'], 1);
      expect(json['is_paid'], true);
      expect(json['tracking_number'], 'TRK-001');
      expect(json['tracking_link'], 'https://tracking.com/TRK-001');
      expect(json['warehouse_id'], 3);
      expect(json['warehouse_name'], 'Main');
      expect(json['driver_id'], 7);
      expect(json['driver_name'], 'Carlos');
      expect(json['status'], 'shipped');
      expect(json['status_id'], 2);
      expect(json['payment_status_id'], 1);
      expect(json['fulfillment_status_id'], 3);
      expect(json['notes'], 'Updated note');
      expect(json['is_confirmed'], true);
      expect(json['confirmation_status'], 'confirmed');
      expect(json['novelty'], 'None');
      expect(json['invoiceable'], true);
      expect(json['items'], isNotNull);
      expect(json['metadata'], isNotNull);
    });
  });
}
