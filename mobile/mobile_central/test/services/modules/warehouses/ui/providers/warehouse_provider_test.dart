import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/warehouses/app/use_cases.dart';
import 'package:mobile_central/services/modules/warehouses/domain/entities.dart';
import 'package:mobile_central/services/modules/warehouses/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockWarehouseRepository implements IWarehouseRepository {
  PaginatedResponse<Warehouse>? getWarehousesResult;
  WarehouseDetail? getWarehouseByIdResult;
  Warehouse? createWarehouseResult;
  Warehouse? updateWarehouseResult;
  List<WarehouseLocation>? getLocationsResult;
  WarehouseLocation? createLocationResult;
  WarehouseLocation? updateLocationResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;
  int? capturedDeleteWarehouseId;
  int? capturedDeleteLocationId;

  @override
  Future<PaginatedResponse<Warehouse>> getWarehouses(GetWarehousesParams? params) async {
    calls.add('getWarehouses');
    if (errorToThrow != null) throw errorToThrow!;
    return getWarehousesResult!;
  }

  @override
  Future<WarehouseDetail> getWarehouseById(int id, {int? businessId}) async {
    calls.add('getWarehouseById');
    if (errorToThrow != null) throw errorToThrow!;
    return getWarehouseByIdResult!;
  }

  @override
  Future<Warehouse> createWarehouse(CreateWarehouseDTO data, {int? businessId}) async {
    calls.add('createWarehouse');
    if (errorToThrow != null) throw errorToThrow!;
    return createWarehouseResult!;
  }

  @override
  Future<Warehouse> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId}) async {
    calls.add('updateWarehouse');
    if (errorToThrow != null) throw errorToThrow!;
    return updateWarehouseResult!;
  }

  @override
  Future<void> deleteWarehouse(int id, {int? businessId}) async {
    calls.add('deleteWarehouse');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<WarehouseLocation>> getLocations(int warehouseId, {int? businessId}) async {
    calls.add('getLocations');
    if (errorToThrow != null) throw errorToThrow!;
    return getLocationsResult!;
  }

  @override
  Future<WarehouseLocation> createLocation(int warehouseId, CreateLocationDTO data, {int? businessId}) async {
    calls.add('createLocation');
    if (errorToThrow != null) throw errorToThrow!;
    return createLocationResult!;
  }

  @override
  Future<WarehouseLocation> updateLocation(int warehouseId, int locationId, UpdateLocationDTO data, {int? businessId}) async {
    calls.add('updateLocation');
    if (errorToThrow != null) throw errorToThrow!;
    return updateLocationResult!;
  }

  @override
  Future<void> deleteLocation(int warehouseId, int locationId, {int? businessId}) async {
    calls.add('deleteLocation');
    capturedDeleteWarehouseId = warehouseId;
    capturedDeleteLocationId = locationId;
    if (errorToThrow != null) throw errorToThrow!;
  }
}

// --- Testable Provider ---

class TestableWarehouseProvider {
  final WarehouseUseCases _useCases;

  List<Warehouse> _warehouses = [];
  Pagination? _pagination;
  WarehouseDetail? _selectedWarehouse;
  List<WarehouseLocation> _locations = [];
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableWarehouseProvider(this._useCases);

  List<Warehouse> get warehouses => _warehouses;
  Pagination? get pagination => _pagination;
  WarehouseDetail? get selectedWarehouse => _selectedWarehouse;
  List<WarehouseLocation> get locations => _locations;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchWarehouses({GetWarehousesParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.getWarehouses(params);
      _warehouses = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchWarehouseDetail(int id, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _selectedWarehouse = await _useCases.getWarehouseById(id, businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<Warehouse?> createWarehouse(CreateWarehouseDTO data, {int? businessId}) async {
    try {
      return await _useCases.createWarehouse(data, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId}) async {
    try {
      await _useCases.updateWarehouse(id, data, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteWarehouse(int id, {int? businessId}) async {
    try {
      await _useCases.deleteWarehouse(id, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<void> fetchLocations(int warehouseId, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      _locations = await _useCases.getLocations(warehouseId, businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<bool> deleteLocation(int warehouseId, int locationId, {int? businessId}) async {
    try {
      await _useCases.deleteLocation(warehouseId, locationId, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }
}

// --- Helpers ---

Warehouse _makeWarehouse({int id = 1, String name = 'TestWarehouse'}) {
  return Warehouse(
    id: id,
    businessId: 1,
    name: name,
    code: 'WH-001',
    address: 'Addr',
    city: 'City',
    state: 'State',
    country: 'CO',
    zipCode: '110111',
    phone: '123',
    contactName: 'CN',
    contactEmail: 'ce@t.com',
    isActive: true,
    isDefault: false,
    isFulfillment: false,
    company: 'Co',
    firstName: 'F',
    lastName: 'L',
    email: 'e@t.com',
    suburb: 'S',
    cityDaneCode: '11001',
    postalCode: '110111',
    street: 'St',
    createdAt: '2026-01-01',
    updatedAt: '2026-01-01',
  );
}

WarehouseDetail _makeWarehouseDetail({int id = 1}) {
  return WarehouseDetail(
    id: id,
    businessId: 1,
    name: 'Detail WH',
    code: 'DWH',
    address: 'A',
    city: 'C',
    state: 'S',
    country: 'CO',
    zipCode: '110111',
    phone: '1',
    contactName: 'CN',
    contactEmail: 'ce@t.com',
    isActive: true,
    isDefault: false,
    isFulfillment: false,
    company: 'Co',
    firstName: 'F',
    lastName: 'L',
    email: 'e@t.com',
    suburb: 'S',
    cityDaneCode: '11001',
    postalCode: '110111',
    street: 'St',
    createdAt: '2026-01-01',
    updatedAt: '2026-01-01',
    locations: [],
  );
}

WarehouseLocation _makeLocation({int id = 1, String name = 'Shelf A'}) {
  return WarehouseLocation(
    id: id,
    warehouseId: 1,
    name: name,
    code: 'SA',
    type: 'shelf',
    isActive: true,
    isFulfillment: false,
    createdAt: '2026-01-01',
    updatedAt: '2026-01-01',
  );
}

Pagination _makePagination({int total = 1, int currentPage = 1, int lastPage = 1}) {
  return Pagination(
    currentPage: currentPage,
    perPage: 10,
    total: total,
    lastPage: lastPage,
    hasNext: currentPage < lastPage,
    hasPrev: currentPage > 1,
  );
}

// --- Tests ---

void main() {
  late MockWarehouseRepository mockRepo;
  late WarehouseUseCases useCases;
  late TestableWarehouseProvider provider;

  setUp(() {
    mockRepo = MockWarehouseRepository();
    useCases = WarehouseUseCases(mockRepo);
    provider = TestableWarehouseProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty warehouses list', () {
      expect(provider.warehouses, isEmpty);
    });

    test('starts with null pagination', () {
      expect(provider.pagination, isNull);
    });

    test('starts with null selected warehouse', () {
      expect(provider.selectedWarehouse, isNull);
    });

    test('starts with empty locations', () {
      expect(provider.locations, isEmpty);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchWarehouses', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getWarehousesResult = PaginatedResponse<Warehouse>(
        data: [_makeWarehouse()],
        pagination: _makePagination(),
      );

      await provider.fetchWarehouses();

      expect(provider.notifications.length, 2);
    });

    test('populates warehouses and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getWarehousesResult = PaginatedResponse<Warehouse>(
        data: [_makeWarehouse(id: 1), _makeWarehouse(id: 2, name: 'Second')],
        pagination: pagination,
      );

      await provider.fetchWarehouses();

      expect(provider.warehouses.length, 2);
      expect(provider.warehouses[0].id, 1);
      expect(provider.warehouses[1].name, 'Second');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchWarehouses();

      expect(provider.error, contains('Server error'));
      expect(provider.warehouses, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchWarehouses();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getWarehousesResult = PaginatedResponse<Warehouse>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchWarehouses();

      expect(provider.error, isNull);
    });
  });

  group('fetchWarehouseDetail', () {
    test('populates selected warehouse on success', () async {
      mockRepo.getWarehouseByIdResult = _makeWarehouseDetail(id: 42);

      await provider.fetchWarehouseDetail(42);

      expect(provider.selectedWarehouse, isNotNull);
      expect(provider.selectedWarehouse!.id, 42);
      expect(provider.isLoading, false);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Not found');

      await provider.fetchWarehouseDetail(99);

      expect(provider.error, contains('Not found'));
      expect(provider.selectedWarehouse, isNull);
    });
  });

  group('createWarehouse', () {
    test('returns created warehouse on success', () async {
      final dto = CreateWarehouseDTO(name: 'New WH', code: 'NWH');
      mockRepo.createWarehouseResult = _makeWarehouse(id: 10, name: 'New WH');

      final result = await provider.createWarehouse(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.name, 'New WH');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateWarehouseDTO(name: 'x', code: 'y');

      final result = await provider.createWarehouse(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateWarehouse', () {
    test('returns true on success', () async {
      final dto = UpdateWarehouseDTO(name: 'Updated', code: 'UPD');
      mockRepo.updateWarehouseResult = _makeWarehouse(id: 5, name: 'Updated');

      final result = await provider.updateWarehouse(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateWarehouseDTO(name: 'x', code: 'y');

      final result = await provider.updateWarehouse(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteWarehouse', () {
    test('returns true on success', () async {
      final result = await provider.deleteWarehouse(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteWarehouse(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });

  group('fetchLocations', () {
    test('populates locations on success', () async {
      mockRepo.getLocationsResult = [_makeLocation(), _makeLocation(id: 2, name: 'Shelf B')];

      await provider.fetchLocations(1);

      expect(provider.locations.length, 2);
      expect(provider.locations[0].name, 'Shelf A');
      expect(provider.locations[1].name, 'Shelf B');
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Location error');

      await provider.fetchLocations(1);

      expect(provider.error, contains('Location error'));
    });
  });

  group('deleteLocation', () {
    test('returns true on success', () async {
      final result = await provider.deleteLocation(1, 5);

      expect(result, true);
      expect(mockRepo.capturedDeleteWarehouseId, 1);
      expect(mockRepo.capturedDeleteLocationId, 5);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete location failed');

      final result = await provider.deleteLocation(1, 5);

      expect(result, false);
      expect(provider.error, contains('Delete location failed'));
    });
  });
}
