import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/shipments/app/use_cases.dart';
import 'package:mobile_central/services/modules/shipments/domain/entities.dart';
import 'package:mobile_central/services/modules/shipments/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockShipmentRepository implements IShipmentRepository {
  final List<String> calls = [];

  PaginatedResponse<Shipment>? getShipmentsResult;
  Map<String, dynamic>? quoteShipmentResult;
  Map<String, dynamic>? generateGuideResult;
  Map<String, dynamic>? trackShipmentResult;
  Map<String, dynamic>? cancelShipmentResult;
  Shipment? createShipmentResult;
  List<OriginAddress>? getOriginAddressesResult;
  OriginAddress? createOriginAddressResult;
  OriginAddress? updateOriginAddressResult;

  Exception? errorToThrow;

  GetShipmentsParams? capturedGetShipmentsParams;
  EnvioClickQuoteRequest? capturedQuoteRequest;
  String? capturedTrackingNumber;
  String? capturedCancelId;
  Map<String, dynamic>? capturedCreateData;
  int? capturedBusinessId;
  CreateOriginAddressDTO? capturedCreateOriginData;
  int? capturedUpdateOriginId;
  Map<String, dynamic>? capturedUpdateOriginData;
  int? capturedDeleteOriginId;

  @override
  Future<PaginatedResponse<Shipment>> getShipments(GetShipmentsParams? params) async {
    calls.add('getShipments');
    capturedGetShipmentsParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getShipmentsResult!;
  }

  @override
  Future<Map<String, dynamic>> quoteShipment(EnvioClickQuoteRequest req) async {
    calls.add('quoteShipment');
    capturedQuoteRequest = req;
    if (errorToThrow != null) throw errorToThrow!;
    return quoteShipmentResult!;
  }

  @override
  Future<Map<String, dynamic>> generateGuide(EnvioClickQuoteRequest req) async {
    calls.add('generateGuide');
    capturedQuoteRequest = req;
    if (errorToThrow != null) throw errorToThrow!;
    return generateGuideResult!;
  }

  @override
  Future<Map<String, dynamic>> trackShipment(String trackingNumber) async {
    calls.add('trackShipment');
    capturedTrackingNumber = trackingNumber;
    if (errorToThrow != null) throw errorToThrow!;
    return trackShipmentResult!;
  }

  @override
  Future<Map<String, dynamic>> cancelShipment(String id) async {
    calls.add('cancelShipment');
    capturedCancelId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return cancelShipmentResult!;
  }

  @override
  Future<Shipment> createShipment(Map<String, dynamic> data) async {
    calls.add('createShipment');
    capturedCreateData = data;
    if (errorToThrow != null) throw errorToThrow!;
    return createShipmentResult!;
  }

  @override
  Future<List<OriginAddress>> getOriginAddresses({int? businessId}) async {
    calls.add('getOriginAddresses');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getOriginAddressesResult!;
  }

  @override
  Future<OriginAddress> createOriginAddress(CreateOriginAddressDTO data, {int? businessId}) async {
    calls.add('createOriginAddress');
    capturedCreateOriginData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createOriginAddressResult!;
  }

  @override
  Future<OriginAddress> updateOriginAddress(int id, Map<String, dynamic> data, {int? businessId}) async {
    calls.add('updateOriginAddress');
    capturedUpdateOriginId = id;
    capturedUpdateOriginData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateOriginAddressResult!;
  }

  @override
  Future<void> deleteOriginAddress(int id, {int? businessId}) async {
    calls.add('deleteOriginAddress');
    capturedDeleteOriginId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Helpers ---

Shipment _makeShipment({int id = 1, String status = 'pending'}) {
  return Shipment(
    id: id, createdAt: '2026-01-01', updatedAt: '2026-01-01',
    status: status, isLastMile: false, isTest: false,
  );
}

OriginAddress _makeOriginAddress({int id = 1, String alias = 'Office'}) {
  return OriginAddress(
    id: id, businessId: 1, alias: alias, company: 'Co', firstName: 'F',
    lastName: 'L', email: 'e', phone: 'p', street: 's', cityDaneCode: 'd',
    city: 'c', state: 's', isDefault: false, createdAt: 'c', updatedAt: 'u',
  );
}

Pagination _makePagination() {
  return Pagination(
    currentPage: 1, perPage: 10, total: 1, lastPage: 1,
    hasNext: false, hasPrev: false,
  );
}

EnvioClickQuoteRequest _makeQuoteRequest() {
  return EnvioClickQuoteRequest(
    description: 'Package',
    contentValue: 50000,
    includeGuideCost: true,
    codPaymentMethod: 'cash',
    packages: [EnvioClickPackage(weight: 1, height: 10, width: 10, length: 10)],
    origin: EnvioClickAddress(address: 'Origin', daneCode: '11001'),
    destination: EnvioClickAddress(address: 'Dest', daneCode: '76001'),
  );
}

// --- Tests ---

void main() {
  late MockShipmentRepository mockRepo;
  late ShipmentUseCases useCases;

  setUp(() {
    mockRepo = MockShipmentRepository();
    useCases = ShipmentUseCases(mockRepo);
  });

  group('getShipments', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Shipment>(
        data: [_makeShipment()],
        pagination: _makePagination(),
      );
      mockRepo.getShipmentsResult = expected;
      final params = GetShipmentsParams(page: 1, pageSize: 10);

      final result = await useCases.getShipments(params);

      expect(result.data.length, 1);
      expect(mockRepo.calls, ['getShipments']);
      expect(mockRepo.capturedGetShipmentsParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getShipmentsResult = PaginatedResponse<Shipment>(
        data: [], pagination: _makePagination(),
      );

      await useCases.getShipments(null);

      expect(mockRepo.capturedGetShipmentsParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getShipments(null), throwsException);
    });
  });

  group('quoteShipment', () {
    test('delegates to repository with correct request', () async {
      final req = _makeQuoteRequest();
      mockRepo.quoteShipmentResult = {'rates': []};

      final result = await useCases.quoteShipment(req);

      expect(result['rates'], isA<List>());
      expect(mockRepo.capturedQuoteRequest, req);
      expect(mockRepo.calls, ['quoteShipment']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Quote failed');

      expect(() => useCases.quoteShipment(_makeQuoteRequest()), throwsException);
    });
  });

  group('generateGuide', () {
    test('delegates to repository with correct request', () async {
      final req = _makeQuoteRequest();
      mockRepo.generateGuideResult = {'guide_id': 'G-001'};

      final result = await useCases.generateGuide(req);

      expect(result['guide_id'], 'G-001');
      expect(mockRepo.capturedQuoteRequest, req);
      expect(mockRepo.calls, ['generateGuide']);
    });
  });

  group('trackShipment', () {
    test('delegates to repository with correct tracking number', () async {
      mockRepo.trackShipmentResult = {'status': 'in_transit'};

      final result = await useCases.trackShipment('TRK-001');

      expect(result['status'], 'in_transit');
      expect(mockRepo.capturedTrackingNumber, 'TRK-001');
      expect(mockRepo.calls, ['trackShipment']);
    });
  });

  group('cancelShipment', () {
    test('delegates to repository with correct id', () async {
      mockRepo.cancelShipmentResult = {'message': 'cancelled'};

      final result = await useCases.cancelShipment('1');

      expect(result['message'], 'cancelled');
      expect(mockRepo.capturedCancelId, '1');
      expect(mockRepo.calls, ['cancelShipment']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Cancel failed');

      expect(() => useCases.cancelShipment('1'), throwsException);
    });
  });

  group('createShipment', () {
    test('delegates to repository with correct data', () async {
      final data = {'order_id': 'ord-1', 'carrier': 'SER'};
      mockRepo.createShipmentResult = _makeShipment(id: 99);

      final result = await useCases.createShipment(data);

      expect(result.id, 99);
      expect(mockRepo.capturedCreateData, data);
      expect(mockRepo.calls, ['createShipment']);
    });
  });

  group('getOriginAddresses', () {
    test('delegates to repository with businessId', () async {
      mockRepo.getOriginAddressesResult = [_makeOriginAddress()];

      final result = await useCases.getOriginAddresses(businessId: 5);

      expect(result.length, 1);
      expect(result[0].alias, 'Office');
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['getOriginAddresses']);
    });
  });

  group('createOriginAddress', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateOriginAddressDTO(
        alias: 'New', company: 'Co', firstName: 'F', lastName: 'L',
        email: 'e', phone: 'p', street: 's', cityDaneCode: 'd',
        city: 'c', state: 's',
      );
      mockRepo.createOriginAddressResult = _makeOriginAddress(id: 10, alias: 'New');

      final result = await useCases.createOriginAddress(dto, businessId: 5);

      expect(result.id, 10);
      expect(result.alias, 'New');
      expect(mockRepo.capturedCreateOriginData, dto);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['createOriginAddress']);
    });
  });

  group('updateOriginAddress', () {
    test('delegates to repository with correct id and data', () async {
      final data = {'alias': 'Updated'};
      mockRepo.updateOriginAddressResult = _makeOriginAddress(id: 5, alias: 'Updated');

      final result = await useCases.updateOriginAddress(5, data, businessId: 3);

      expect(result.id, 5);
      expect(result.alias, 'Updated');
      expect(mockRepo.capturedUpdateOriginId, 5);
      expect(mockRepo.capturedUpdateOriginData, data);
      expect(mockRepo.capturedBusinessId, 3);
      expect(mockRepo.calls, ['updateOriginAddress']);
    });
  });

  group('deleteOriginAddress', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteOriginAddress(7, businessId: 3);

      expect(mockRepo.capturedDeleteOriginId, 7);
      expect(mockRepo.capturedBusinessId, 3);
      expect(mockRepo.calls, ['deleteOriginAddress']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      expect(() => useCases.deleteOriginAddress(7), throwsException);
    });
  });
}
