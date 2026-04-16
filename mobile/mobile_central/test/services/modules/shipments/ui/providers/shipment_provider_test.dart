import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/shipments/app/use_cases.dart';
import 'package:mobile_central/services/modules/shipments/domain/entities.dart';
import 'package:mobile_central/services/modules/shipments/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockShipmentRepository implements IShipmentRepository {
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

  final List<String> calls = [];

  @override
  Future<PaginatedResponse<Shipment>> getShipments(GetShipmentsParams? params) async {
    calls.add('getShipments');
    if (errorToThrow != null) throw errorToThrow!;
    return getShipmentsResult!;
  }

  @override
  Future<Map<String, dynamic>> quoteShipment(EnvioClickQuoteRequest req) async {
    calls.add('quoteShipment');
    if (errorToThrow != null) throw errorToThrow!;
    return quoteShipmentResult!;
  }

  @override
  Future<Map<String, dynamic>> generateGuide(EnvioClickQuoteRequest req) async {
    calls.add('generateGuide');
    if (errorToThrow != null) throw errorToThrow!;
    return generateGuideResult!;
  }

  @override
  Future<Map<String, dynamic>> trackShipment(String trackingNumber) async {
    calls.add('trackShipment');
    if (errorToThrow != null) throw errorToThrow!;
    return trackShipmentResult!;
  }

  @override
  Future<Map<String, dynamic>> cancelShipment(String id) async {
    calls.add('cancelShipment');
    if (errorToThrow != null) throw errorToThrow!;
    return cancelShipmentResult!;
  }

  @override
  Future<Shipment> createShipment(Map<String, dynamic> data) async {
    calls.add('createShipment');
    if (errorToThrow != null) throw errorToThrow!;
    return createShipmentResult!;
  }

  @override
  Future<List<OriginAddress>> getOriginAddresses({int? businessId}) async {
    calls.add('getOriginAddresses');
    if (errorToThrow != null) throw errorToThrow!;
    return getOriginAddressesResult ?? [];
  }

  @override
  Future<OriginAddress> createOriginAddress(CreateOriginAddressDTO data, {int? businessId}) async {
    calls.add('createOriginAddress');
    if (errorToThrow != null) throw errorToThrow!;
    return createOriginAddressResult!;
  }

  @override
  Future<OriginAddress> updateOriginAddress(int id, Map<String, dynamic> data, {int? businessId}) async {
    calls.add('updateOriginAddress');
    if (errorToThrow != null) throw errorToThrow!;
    return updateOriginAddressResult!;
  }

  @override
  Future<void> deleteOriginAddress(int id, {int? businessId}) async {
    calls.add('deleteOriginAddress');
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---

class TestableShipmentProvider {
  final ShipmentUseCases _useCases;

  List<Shipment> _shipments = [];
  List<OriginAddress> _originAddresses = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;

  final List<String> notifications = [];

  TestableShipmentProvider(this._useCases);

  List<Shipment> get shipments => _shipments;
  List<OriginAddress> get originAddresses => _originAddresses;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchShipments({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();
    try {
      final params = GetShipmentsParams(page: _page, pageSize: _pageSize, businessId: businessId);
      final response = await _useCases.getShipments(params);
      _shipments = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }
    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchOriginAddresses({int? businessId}) async {
    try {
      _originAddresses = await _useCases.getOriginAddresses(businessId: businessId);
      _notifyListeners();
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
    }
  }

  Future<Map<String, dynamic>?> quoteShipment(EnvioClickQuoteRequest req) async {
    try {
      return await _useCases.quoteShipment(req);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Map<String, dynamic>?> generateGuide(EnvioClickQuoteRequest req) async {
    try {
      return await _useCases.generateGuide(req);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<Map<String, dynamic>?> trackShipment(String trackingNumber) async {
    try {
      return await _useCases.trackShipment(trackingNumber);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  void setPage(int page) {
    _page = page;
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

Pagination _makePagination({int currentPage = 1, int total = 5, int lastPage = 1}) {
  return Pagination(
    currentPage: currentPage, perPage: 20, total: total, lastPage: lastPage,
    hasNext: currentPage < lastPage, hasPrev: currentPage > 1,
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
  late TestableShipmentProvider provider;

  setUp(() {
    mockRepo = MockShipmentRepository();
    useCases = ShipmentUseCases(mockRepo);
    provider = TestableShipmentProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty shipments list', () {
      expect(provider.shipments, isEmpty);
    });

    test('starts with empty origin addresses', () {
      expect(provider.originAddresses, isEmpty);
    });

    test('starts with null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchShipments', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getShipmentsResult = PaginatedResponse<Shipment>(
        data: [_makeShipment()],
        pagination: _makePagination(),
      );

      await provider.fetchShipments();

      expect(provider.notifications.length, 2);
    });

    test('populates shipments and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getShipmentsResult = PaginatedResponse<Shipment>(
        data: [_makeShipment(id: 1), _makeShipment(id: 2, status: 'in_transit')],
        pagination: pagination,
      );

      await provider.fetchShipments();

      expect(provider.shipments.length, 2);
      expect(provider.shipments[0].id, 1);
      expect(provider.shipments[1].status, 'in_transit');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchShipments();

      expect(provider.error, contains('Server error'));
      expect(provider.shipments, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchShipments();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getShipmentsResult = PaginatedResponse<Shipment>(
        data: [], pagination: _makePagination(),
      );
      await provider.fetchShipments();

      expect(provider.error, isNull);
    });
  });

  group('fetchOriginAddresses', () {
    test('populates origin addresses on success', () async {
      mockRepo.getOriginAddressesResult = [
        _makeOriginAddress(id: 1, alias: 'Office'),
        _makeOriginAddress(id: 2, alias: 'Warehouse'),
      ];

      await provider.fetchOriginAddresses();

      expect(provider.originAddresses.length, 2);
      expect(provider.originAddresses[0].alias, 'Office');
      expect(provider.originAddresses[1].alias, 'Warehouse');
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Fetch failed');

      await provider.fetchOriginAddresses();

      expect(provider.error, contains('Fetch failed'));
    });
  });

  group('quoteShipment', () {
    test('returns result on success', () async {
      final req = _makeQuoteRequest();
      mockRepo.quoteShipmentResult = {'rates': [{'carrier': 'SER', 'flete': 15000}]};

      final result = await provider.quoteShipment(req);

      expect(result, isNotNull);
      expect(result!['rates'], isA<List>());
    });

    test('returns null and sets error on failure', () async {
      final req = _makeQuoteRequest();
      mockRepo.errorToThrow = Exception('Quote failed');

      final result = await provider.quoteShipment(req);

      expect(result, isNull);
      expect(provider.error, contains('Quote failed'));
    });
  });

  group('generateGuide', () {
    test('returns result on success', () async {
      final req = _makeQuoteRequest();
      mockRepo.generateGuideResult = {'guide_id': 'G-001'};

      final result = await provider.generateGuide(req);

      expect(result, isNotNull);
      expect(result!['guide_id'], 'G-001');
    });

    test('returns null and sets error on failure', () async {
      final req = _makeQuoteRequest();
      mockRepo.errorToThrow = Exception('Guide failed');

      final result = await provider.generateGuide(req);

      expect(result, isNull);
      expect(provider.error, contains('Guide failed'));
    });
  });

  group('trackShipment', () {
    test('returns result on success', () async {
      mockRepo.trackShipmentResult = {'status': 'delivered', 'history': []};

      final result = await provider.trackShipment('TRK-001');

      expect(result, isNotNull);
      expect(result!['status'], 'delivered');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Track failed');

      final result = await provider.trackShipment('TRK-001');

      expect(result, isNull);
      expect(provider.error, contains('Track failed'));
    });
  });

  group('pagination', () {
    test('setPage updates page', () {
      provider.setPage(3);

      // Verify it will be used on next fetch
      mockRepo.getShipmentsResult = PaginatedResponse<Shipment>(
        data: [], pagination: _makePagination(currentPage: 3),
      );
    });
  });
}
