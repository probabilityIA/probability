import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/warehouses/app/use_cases.dart';
import 'package:mobile_central/services/modules/warehouses/domain/entities.dart';
import 'package:mobile_central/services/modules/warehouses/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockWarehouseRepository implements IWarehouseRepository {
  final List<String> calls = [];

  PaginatedResponse<Warehouse>? getWarehousesResult;
  WarehouseDetail? getWarehouseByIdResult;
  Warehouse? createWarehouseResult;
  Warehouse? updateWarehouseResult;
  List<WarehouseLocation>? getLocationsResult;
  WarehouseLocation? createLocationResult;
  WarehouseLocation? updateLocationResult;

  Exception? errorToThrow;

  GetWarehousesParams? capturedParams;
  int? capturedId;
  int? capturedBusinessId;
  CreateWarehouseDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateWarehouseDTO? capturedUpdateData;
  int? capturedDeleteId;
  int? capturedWarehouseId;
  CreateLocationDTO? capturedCreateLocationData;
  int? capturedLocationId;
  UpdateLocationDTO? capturedUpdateLocationData;
  int? capturedDeleteWarehouseId;
  int? capturedDeleteLocationId;

  @override
  Future<PaginatedResponse<Warehouse>> getWarehouses(GetWarehousesParams? params) async {
    calls.add('getWarehouses');
    capturedParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getWarehousesResult!;
  }

  @override
  Future<WarehouseDetail> getWarehouseById(int id, {int? businessId}) async {
    calls.add('getWarehouseById');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getWarehouseByIdResult!;
  }

  @override
  Future<Warehouse> createWarehouse(CreateWarehouseDTO data, {int? businessId}) async {
    calls.add('createWarehouse');
    capturedCreateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createWarehouseResult!;
  }

  @override
  Future<Warehouse> updateWarehouse(int id, UpdateWarehouseDTO data, {int? businessId}) async {
    calls.add('updateWarehouse');
    capturedUpdateId = id;
    capturedUpdateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateWarehouseResult!;
  }

  @override
  Future<void> deleteWarehouse(int id, {int? businessId}) async {
    calls.add('deleteWarehouse');
    capturedDeleteId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
  }

  @override
  Future<List<WarehouseLocation>> getLocations(int warehouseId, {int? businessId}) async {
    calls.add('getLocations');
    capturedWarehouseId = warehouseId;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getLocationsResult!;
  }

  @override
  Future<WarehouseLocation> createLocation(int warehouseId, CreateLocationDTO data, {int? businessId}) async {
    calls.add('createLocation');
    capturedWarehouseId = warehouseId;
    capturedCreateLocationData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createLocationResult!;
  }

  @override
  Future<WarehouseLocation> updateLocation(int warehouseId, int locationId, UpdateLocationDTO data, {int? businessId}) async {
    calls.add('updateLocation');
    capturedWarehouseId = warehouseId;
    capturedLocationId = locationId;
    capturedUpdateLocationData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateLocationResult!;
  }

  @override
  Future<void> deleteLocation(int warehouseId, int locationId, {int? businessId}) async {
    calls.add('deleteLocation');
    capturedDeleteWarehouseId = warehouseId;
    capturedDeleteLocationId = locationId;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
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

Pagination _makePagination() {
  return Pagination(
    currentPage: 1,
    perPage: 10,
    total: 1,
    lastPage: 1,
    hasNext: false,
    hasPrev: false,
  );
}

// --- Tests ---

void main() {
  late MockWarehouseRepository mockRepo;
  late WarehouseUseCases useCases;

  setUp(() {
    mockRepo = MockWarehouseRepository();
    useCases = WarehouseUseCases(mockRepo);
  });

  group('getWarehouses', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<Warehouse>(
        data: [_makeWarehouse()],
        pagination: _makePagination(),
      );
      mockRepo.getWarehousesResult = expected;
      final params = GetWarehousesParams(page: 1, pageSize: 10);

      final result = await useCases.getWarehouses(params);

      expect(result.data.length, 1);
      expect(result.data[0].name, 'TestWarehouse');
      expect(mockRepo.calls, ['getWarehouses']);
      expect(mockRepo.capturedParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getWarehousesResult = PaginatedResponse<Warehouse>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getWarehouses(null);

      expect(mockRepo.capturedParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getWarehouses(null), throwsException);
    });
  });

  group('getWarehouseById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getWarehouseByIdResult = _makeWarehouseDetail(id: 42);

      final result = await useCases.getWarehouseById(42);

      expect(result.id, 42);
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getWarehouseById']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getWarehouseByIdResult = _makeWarehouseDetail();

      await useCases.getWarehouseById(1, businessId: 5);

      expect(mockRepo.capturedBusinessId, 5);
    });
  });

  group('createWarehouse', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateWarehouseDTO(name: 'New WH', code: 'NWH');
      mockRepo.createWarehouseResult = _makeWarehouse(id: 10, name: 'New WH');

      final result = await useCases.createWarehouse(dto);

      expect(result.id, 10);
      expect(result.name, 'New WH');
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createWarehouse']);
    });

    test('passes businessId to repository', () async {
      final dto = CreateWarehouseDTO(name: 'N', code: 'C');
      mockRepo.createWarehouseResult = _makeWarehouse();

      await useCases.createWarehouse(dto, businessId: 7);

      expect(mockRepo.capturedBusinessId, 7);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateWarehouseDTO(name: 'x', code: 'y');

      expect(() => useCases.createWarehouse(dto), throwsException);
    });
  });

  group('updateWarehouse', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateWarehouseDTO(name: 'Updated', code: 'UPD');
      mockRepo.updateWarehouseResult = _makeWarehouse(id: 5, name: 'Updated');

      final result = await useCases.updateWarehouse(5, dto);

      expect(result.id, 5);
      expect(result.name, 'Updated');
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateWarehouse']);
    });
  });

  group('deleteWarehouse', () {
    test('delegates to repository with correct id', () async {
      await useCases.deleteWarehouse(7);

      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deleteWarehouse']);
    });

    test('passes businessId to repository', () async {
      await useCases.deleteWarehouse(1, businessId: 2);

      expect(mockRepo.capturedBusinessId, 2);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      expect(() => useCases.deleteWarehouse(7), throwsException);
    });
  });

  group('getLocations', () {
    test('delegates to repository with correct warehouseId', () async {
      mockRepo.getLocationsResult = [_makeLocation(), _makeLocation(id: 2)];

      final result = await useCases.getLocations(1);

      expect(result.length, 2);
      expect(mockRepo.capturedWarehouseId, 1);
      expect(mockRepo.calls, ['getLocations']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getLocationsResult = [];

      await useCases.getLocations(1, businessId: 3);

      expect(mockRepo.capturedBusinessId, 3);
    });
  });

  group('createLocation', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateLocationDTO(name: 'New Loc', code: 'NL');
      mockRepo.createLocationResult = _makeLocation(id: 20, name: 'New Loc');

      final result = await useCases.createLocation(1, dto);

      expect(result.id, 20);
      expect(result.name, 'New Loc');
      expect(mockRepo.capturedWarehouseId, 1);
      expect(mockRepo.capturedCreateLocationData, dto);
      expect(mockRepo.calls, ['createLocation']);
    });
  });

  group('updateLocation', () {
    test('delegates to repository with correct ids and data', () async {
      final dto = UpdateLocationDTO(name: 'Updated', code: 'UL');
      mockRepo.updateLocationResult = _makeLocation(id: 5, name: 'Updated');

      final result = await useCases.updateLocation(1, 5, dto);

      expect(result.id, 5);
      expect(mockRepo.capturedWarehouseId, 1);
      expect(mockRepo.capturedLocationId, 5);
      expect(mockRepo.capturedUpdateLocationData, dto);
      expect(mockRepo.calls, ['updateLocation']);
    });
  });

  group('deleteLocation', () {
    test('delegates to repository with correct ids', () async {
      await useCases.deleteLocation(1, 5);

      expect(mockRepo.capturedDeleteWarehouseId, 1);
      expect(mockRepo.capturedDeleteLocationId, 5);
      expect(mockRepo.calls, ['deleteLocation']);
    });

    test('passes businessId to repository', () async {
      await useCases.deleteLocation(1, 5, businessId: 3);

      expect(mockRepo.capturedBusinessId, 3);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      expect(() => useCases.deleteLocation(1, 5), throwsException);
    });
  });
}
