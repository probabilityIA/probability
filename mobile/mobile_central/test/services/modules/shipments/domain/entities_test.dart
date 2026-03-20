import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/shipments/domain/entities.dart';

void main() {
  group('Shipment', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'created_at': '2026-01-01T00:00:00Z',
        'updated_at': '2026-01-02T00:00:00Z',
        'order_id': 'ord-123',
        'client_name': 'Maria',
        'destination_address': 'Calle 100 #10-30',
        'tracking_number': 'TRK-001',
        'tracking_url': 'https://track.example.com/TRK-001',
        'carrier': 'Servientrega',
        'carrier_code': 'SER',
        'guide_id': 'G-001',
        'guide_url': 'https://guide.example.com/G-001',
        'status': 'in_transit',
        'shipped_at': '2026-01-01T10:00:00Z',
        'delivered_at': null,
        'shipping_cost': 15000.0,
        'insurance_cost': 2000.0,
        'total_cost': 17000.0,
        'weight': 2.5,
        'height': 30.0,
        'width': 20.0,
        'length': 40.0,
        'warehouse_name': 'Main Warehouse',
        'driver_name': 'Carlos',
        'is_last_mile': true,
        'is_test': false,
        'estimated_delivery': '2026-01-03',
        'delivery_notes': 'Ring bell',
        'customer_name': 'Maria Garcia',
        'customer_email': 'maria@example.com',
        'customer_phone': '+57123456789',
        'customer_dni': '123456789',
        'order_number': 'ORD-001',
      };

      final s = Shipment.fromJson(json);

      expect(s.id, 1);
      expect(s.createdAt, '2026-01-01T00:00:00Z');
      expect(s.updatedAt, '2026-01-02T00:00:00Z');
      expect(s.orderId, 'ord-123');
      expect(s.clientName, 'Maria');
      expect(s.destinationAddress, 'Calle 100 #10-30');
      expect(s.trackingNumber, 'TRK-001');
      expect(s.trackingUrl, 'https://track.example.com/TRK-001');
      expect(s.carrier, 'Servientrega');
      expect(s.carrierCode, 'SER');
      expect(s.guideId, 'G-001');
      expect(s.guideUrl, 'https://guide.example.com/G-001');
      expect(s.status, 'in_transit');
      expect(s.shippedAt, '2026-01-01T10:00:00Z');
      expect(s.deliveredAt, isNull);
      expect(s.shippingCost, 15000.0);
      expect(s.insuranceCost, 2000.0);
      expect(s.totalCost, 17000.0);
      expect(s.weight, 2.5);
      expect(s.height, 30.0);
      expect(s.width, 20.0);
      expect(s.length, 40.0);
      expect(s.warehouseName, 'Main Warehouse');
      expect(s.driverName, 'Carlos');
      expect(s.isLastMile, true);
      expect(s.isTest, false);
      expect(s.estimatedDelivery, '2026-01-03');
      expect(s.deliveryNotes, 'Ring bell');
      expect(s.customerName, 'Maria Garcia');
      expect(s.customerEmail, 'maria@example.com');
      expect(s.customerPhone, '+57123456789');
      expect(s.customerDni, '123456789');
      expect(s.orderNumber, 'ORD-001');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final s = Shipment.fromJson(json);

      expect(s.id, 0);
      expect(s.createdAt, '');
      expect(s.updatedAt, '');
      expect(s.status, 'pending');
      expect(s.isLastMile, false);
      expect(s.isTest, false);
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'created_at': 'c',
        'updated_at': 'u',
        'status': 'pending',
        'is_last_mile': false,
        'is_test': false,
      };

      final s = Shipment.fromJson(json);

      expect(s.orderId, isNull);
      expect(s.clientName, isNull);
      expect(s.destinationAddress, isNull);
      expect(s.trackingNumber, isNull);
      expect(s.trackingUrl, isNull);
      expect(s.carrier, isNull);
      expect(s.carrierCode, isNull);
      expect(s.guideId, isNull);
      expect(s.guideUrl, isNull);
      expect(s.shippedAt, isNull);
      expect(s.deliveredAt, isNull);
      expect(s.shippingCost, isNull);
      expect(s.insuranceCost, isNull);
      expect(s.totalCost, isNull);
      expect(s.weight, isNull);
      expect(s.height, isNull);
      expect(s.width, isNull);
      expect(s.length, isNull);
      expect(s.warehouseName, isNull);
      expect(s.driverName, isNull);
      expect(s.estimatedDelivery, isNull);
      expect(s.deliveryNotes, isNull);
      expect(s.customerName, isNull);
      expect(s.customerEmail, isNull);
      expect(s.customerPhone, isNull);
      expect(s.customerDni, isNull);
      expect(s.orderNumber, isNull);
    });
  });

  group('OriginAddress', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'id': 1,
        'business_id': 5,
        'alias': 'Main Office',
        'company': 'My Company',
        'first_name': 'John',
        'last_name': 'Doe',
        'email': 'john@example.com',
        'phone': '+57123',
        'street': 'Calle 100 #10-30',
        'suburb': 'Chapinero',
        'city_dane_code': '11001',
        'city': 'Bogota',
        'state': 'Cundinamarca',
        'postal_code': '110111',
        'is_default': true,
        'created_at': '2026-01-01',
        'updated_at': '2026-01-02',
      };

      final addr = OriginAddress.fromJson(json);

      expect(addr.id, 1);
      expect(addr.businessId, 5);
      expect(addr.alias, 'Main Office');
      expect(addr.company, 'My Company');
      expect(addr.firstName, 'John');
      expect(addr.lastName, 'Doe');
      expect(addr.email, 'john@example.com');
      expect(addr.phone, '+57123');
      expect(addr.street, 'Calle 100 #10-30');
      expect(addr.suburb, 'Chapinero');
      expect(addr.cityDaneCode, '11001');
      expect(addr.city, 'Bogota');
      expect(addr.state, 'Cundinamarca');
      expect(addr.postalCode, '110111');
      expect(addr.isDefault, true);
      expect(addr.createdAt, '2026-01-01');
      expect(addr.updatedAt, '2026-01-02');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final addr = OriginAddress.fromJson(json);

      expect(addr.id, 0);
      expect(addr.businessId, 0);
      expect(addr.alias, '');
      expect(addr.company, '');
      expect(addr.firstName, '');
      expect(addr.lastName, '');
      expect(addr.email, '');
      expect(addr.phone, '');
      expect(addr.street, '');
      expect(addr.cityDaneCode, '');
      expect(addr.city, '');
      expect(addr.state, '');
      expect(addr.isDefault, false);
      expect(addr.createdAt, '');
      expect(addr.updatedAt, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'id': 1,
        'business_id': 1,
        'alias': 'a',
        'company': 'c',
        'first_name': 'f',
        'last_name': 'l',
        'email': 'e',
        'phone': 'p',
        'street': 's',
        'city_dane_code': 'd',
        'city': 'c',
        'state': 's',
        'is_default': false,
        'created_at': 'c',
        'updated_at': 'u',
      };

      final addr = OriginAddress.fromJson(json);

      expect(addr.suburb, isNull);
      expect(addr.postalCode, isNull);
    });
  });

  group('EnvioClickRate', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'idRate': 1,
        'idProduct': 10,
        'product': 'Express',
        'idCarrier': 20,
        'carrier': 'Servientrega',
        'flete': 15000.0,
        'deliveryDays': 2,
        'quotationType': 'standard',
        'minimumInsurance': 1000.0,
        'extraInsurance': 500.0,
        'cod': true,
      };

      final rate = EnvioClickRate.fromJson(json);

      expect(rate.idRate, 1);
      expect(rate.idProduct, 10);
      expect(rate.product, 'Express');
      expect(rate.idCarrier, 20);
      expect(rate.carrier, 'Servientrega');
      expect(rate.flete, 15000.0);
      expect(rate.deliveryDays, 2);
      expect(rate.quotationType, 'standard');
      expect(rate.minimumInsurance, 1000.0);
      expect(rate.extraInsurance, 500.0);
      expect(rate.cod, true);
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final rate = EnvioClickRate.fromJson(json);

      expect(rate.idRate, 0);
      expect(rate.idProduct, 0);
      expect(rate.product, '');
      expect(rate.idCarrier, 0);
      expect(rate.carrier, '');
      expect(rate.flete, 0.0);
      expect(rate.deliveryDays, 0);
      expect(rate.quotationType, '');
    });

    test('fromJson handles null optional fields', () {
      final json = {
        'idRate': 1,
        'idProduct': 1,
        'product': 'p',
        'idCarrier': 1,
        'carrier': 'c',
        'flete': 100,
        'deliveryDays': 1,
        'quotationType': 'q',
      };

      final rate = EnvioClickRate.fromJson(json);

      expect(rate.minimumInsurance, isNull);
      expect(rate.extraInsurance, isNull);
      expect(rate.cod, isNull);
    });
  });

  group('EnvioClickTrackHistory', () {
    test('fromJson parses all fields correctly', () {
      final json = {
        'date': '2026-01-01',
        'status': 'in_transit',
        'description': 'Package picked up',
        'location': 'Bogota',
      };

      final history = EnvioClickTrackHistory.fromJson(json);

      expect(history.date, '2026-01-01');
      expect(history.status, 'in_transit');
      expect(history.description, 'Package picked up');
      expect(history.location, 'Bogota');
    });

    test('fromJson uses defaults for missing fields', () {
      final json = <String, dynamic>{};

      final history = EnvioClickTrackHistory.fromJson(json);

      expect(history.date, '');
      expect(history.status, '');
      expect(history.description, '');
      expect(history.location, '');
    });
  });

  group('GetShipmentsParams', () {
    test('toQueryParams includes all non-null fields', () {
      final params = GetShipmentsParams(
        page: 1,
        pageSize: 20,
        orderId: 'ord-1',
        trackingNumber: 'TRK-001',
        carrier: 'Servientrega',
        status: 'in_transit',
        customerName: 'Maria',
        businessId: 5,
        startDate: '2026-01-01',
        endDate: '2026-12-31',
        sortBy: 'created_at',
        sortOrder: 'desc',
        isTest: false,
      );

      final qp = params.toQueryParams();

      expect(qp['page'], 1);
      expect(qp['page_size'], 20);
      expect(qp['order_id'], 'ord-1');
      expect(qp['tracking_number'], 'TRK-001');
      expect(qp['carrier'], 'Servientrega');
      expect(qp['status'], 'in_transit');
      expect(qp['customer_name'], 'Maria');
      expect(qp['business_id'], 5);
      expect(qp['start_date'], '2026-01-01');
      expect(qp['end_date'], '2026-12-31');
      expect(qp['sort_by'], 'created_at');
      expect(qp['sort_order'], 'desc');
      expect(qp['is_test'], false);
    });

    test('toQueryParams excludes null fields', () {
      final params = GetShipmentsParams(page: 1);

      final qp = params.toQueryParams();

      expect(qp.length, 1);
      expect(qp.containsKey('page'), true);
      expect(qp.containsKey('carrier'), false);
    });

    test('toQueryParams returns empty map when all fields are null', () {
      final params = GetShipmentsParams();

      final qp = params.toQueryParams();

      expect(qp, isEmpty);
    });
  });

  group('EnvioClickAddress', () {
    test('toJson includes required fields', () {
      final addr = EnvioClickAddress(address: 'Calle 100', daneCode: '11001');

      final json = addr.toJson();

      expect(json['address'], 'Calle 100');
      expect(json['daneCode'], '11001');
    });

    test('toJson includes all non-null optional fields', () {
      final addr = EnvioClickAddress(
        company: 'Company',
        firstName: 'John',
        lastName: 'Doe',
        email: 'john@ex.com',
        phone: '+57123',
        address: 'Calle 100',
        suburb: 'Chapinero',
        crossStreet: 'Carrera 10',
        reference: 'Near park',
        daneCode: '11001',
      );

      final json = addr.toJson();

      expect(json['company'], 'Company');
      expect(json['firstName'], 'John');
      expect(json['lastName'], 'Doe');
      expect(json['email'], 'john@ex.com');
      expect(json['phone'], '+57123');
      expect(json['suburb'], 'Chapinero');
      expect(json['crossStreet'], 'Carrera 10');
      expect(json['reference'], 'Near park');
    });

    test('toJson excludes null optional fields', () {
      final addr = EnvioClickAddress(address: 'A', daneCode: 'D');

      final json = addr.toJson();

      expect(json.length, 2);
      expect(json.containsKey('company'), false);
      expect(json.containsKey('firstName'), false);
    });
  });

  group('EnvioClickPackage', () {
    test('toJson produces correct structure', () {
      final pkg = EnvioClickPackage(weight: 2.5, height: 30, width: 20, length: 40);

      final json = pkg.toJson();

      expect(json['weight'], 2.5);
      expect(json['height'], 30.0);
      expect(json['width'], 20.0);
      expect(json['length'], 40.0);
    });
  });

  group('EnvioClickQuoteRequest', () {
    test('toJson includes required fields', () {
      final req = EnvioClickQuoteRequest(
        description: 'Package',
        contentValue: 50000,
        includeGuideCost: true,
        codPaymentMethod: 'cash',
        packages: [EnvioClickPackage(weight: 1, height: 10, width: 10, length: 10)],
        origin: EnvioClickAddress(address: 'Origin', daneCode: '11001'),
        destination: EnvioClickAddress(address: 'Dest', daneCode: '76001'),
      );

      final json = req.toJson();

      expect(json['description'], 'Package');
      expect(json['contentValue'], 50000);
      expect(json['includeGuideCost'], true);
      expect(json['codPaymentMethod'], 'cash');
      expect((json['packages'] as List).length, 1);
      expect(json['origin']['address'], 'Origin');
      expect(json['destination']['address'], 'Dest');
    });

    test('toJson includes all non-null optional fields', () {
      final req = EnvioClickQuoteRequest(
        businessId: 5,
        idRate: 10,
        myShipmentReference: 'REF-001',
        externalOrderId: 'ext-ord-1',
        orderUuid: 'uuid-123',
        requestPickup: true,
        pickupDate: '2026-03-01',
        insurance: true,
        description: 'Desc',
        contentValue: 50000,
        codValue: 50000,
        includeGuideCost: true,
        codPaymentMethod: 'cash',
        totalCost: 17000,
        packages: [EnvioClickPackage(weight: 1, height: 10, width: 10, length: 10)],
        origin: EnvioClickAddress(address: 'O', daneCode: 'D'),
        destination: EnvioClickAddress(address: 'D', daneCode: 'D'),
      );

      final json = req.toJson();

      expect(json['business_id'], 5);
      expect(json['idRate'], 10);
      expect(json['myShipmentReference'], 'REF-001');
      expect(json['external_order_id'], 'ext-ord-1');
      expect(json['order_uuid'], 'uuid-123');
      expect(json['requestPickup'], true);
      expect(json['pickupDate'], '2026-03-01');
      expect(json['insurance'], true);
      expect(json['codValue'], 50000);
      expect(json['totalCost'], 17000);
    });

    test('toJson excludes null optional fields', () {
      final req = EnvioClickQuoteRequest(
        description: 'D',
        contentValue: 100,
        includeGuideCost: false,
        codPaymentMethod: 'none',
        packages: [],
        origin: EnvioClickAddress(address: 'O', daneCode: 'D'),
        destination: EnvioClickAddress(address: 'D', daneCode: 'D'),
      );

      final json = req.toJson();

      expect(json.containsKey('business_id'), false);
      expect(json.containsKey('idRate'), false);
      expect(json.containsKey('myShipmentReference'), false);
      expect(json.containsKey('insurance'), false);
      expect(json.containsKey('codValue'), false);
      expect(json.containsKey('totalCost'), false);
    });
  });

  group('CreateOriginAddressDTO', () {
    test('toJson includes required fields', () {
      final dto = CreateOriginAddressDTO(
        alias: 'Office',
        company: 'Co',
        firstName: 'John',
        lastName: 'Doe',
        email: 'j@e.com',
        phone: '+57123',
        street: 'Calle 1',
        cityDaneCode: '11001',
        city: 'Bogota',
        state: 'Cundinamarca',
      );

      final json = dto.toJson();

      expect(json['alias'], 'Office');
      expect(json['company'], 'Co');
      expect(json['first_name'], 'John');
      expect(json['last_name'], 'Doe');
      expect(json['email'], 'j@e.com');
      expect(json['phone'], '+57123');
      expect(json['street'], 'Calle 1');
      expect(json['city_dane_code'], '11001');
      expect(json['city'], 'Bogota');
      expect(json['state'], 'Cundinamarca');
    });

    test('toJson includes optional fields when present', () {
      final dto = CreateOriginAddressDTO(
        alias: 'a',
        company: 'c',
        firstName: 'f',
        lastName: 'l',
        email: 'e',
        phone: 'p',
        street: 's',
        suburb: 'Chapinero',
        cityDaneCode: 'd',
        city: 'c',
        state: 's',
        postalCode: '110111',
        isDefault: true,
      );

      final json = dto.toJson();

      expect(json['suburb'], 'Chapinero');
      expect(json['postal_code'], '110111');
      expect(json['is_default'], true);
    });

    test('toJson excludes null optional fields', () {
      final dto = CreateOriginAddressDTO(
        alias: 'a',
        company: 'c',
        firstName: 'f',
        lastName: 'l',
        email: 'e',
        phone: 'p',
        street: 's',
        cityDaneCode: 'd',
        city: 'c',
        state: 's',
      );

      final json = dto.toJson();

      expect(json.containsKey('suburb'), false);
      expect(json.containsKey('postal_code'), false);
      expect(json.containsKey('is_default'), false);
    });
  });
}
