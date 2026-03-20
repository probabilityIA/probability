import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/drivers/app/use_cases.dart';
import 'package:mobile_central/services/modules/drivers/domain/entities.dart';
import 'package:mobile_central/services/modules/drivers/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockDriverRepository implements IDriverRepository {
  final List<String> calls = [];

  PaginatedResponse<DriverInfo>? getDriversResult;
  DriverInfo? getDriverByIdResult;
  DriverInfo? createDriverResult;
  DriverInfo? updateDriverResult;
  Map<String, dynamic>? deleteDriverResult;

  Exception? errorToThrow;

  GetDriversParams? capturedGetDriversParams;
  int? capturedId;
  int? capturedBusinessId;
  CreateDriverDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateDriverDTO? capturedUpdateData;
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<DriverInfo>> getDrivers(GetDriversParams? params) async {
    calls.add('getDrivers');
    capturedGetDriversParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getDriversResult!;
  }

  @override
  Future<DriverInfo> getDriverById(int id, {int? businessId}) async {
    calls.add('getDriverById');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getDriverByIdResult!;
  }

  @override
  Future<DriverInfo> createDriver(CreateDriverDTO data, {int? businessId}) async {
    calls.add('createDriver');
    capturedCreateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createDriverResult!;
  }

  @override
  Future<DriverInfo> updateDriver(int id, UpdateDriverDTO data, {int? businessId}) async {
    calls.add('updateDriver');
    capturedUpdateId = id;
    capturedUpdateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateDriverResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteDriver(int id, {int? businessId}) async {
    calls.add('deleteDriver');
    capturedDeleteId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteDriverResult ?? {'message': 'deleted'};
  }
}

// --- Helpers ---

DriverInfo _makeDriver({int id = 1, String firstName = 'TestDriver'}) {
  return DriverInfo(
    id: id,
    businessId: 1,
    firstName: firstName,
    lastName: 'Last',
    email: 'test@test.com',
    phone: '555-1234',
    identification: 'CC-123',
    status: 'active',
    photoUrl: '',
    licenseType: 'B1',
    createdAt: '2026-01-01',
    updatedAt: '2026-01-02',
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
  late MockDriverRepository mockRepo;
  late DriverUseCases useCases;

  setUp(() {
    mockRepo = MockDriverRepository();
    useCases = DriverUseCases(mockRepo);
  });

  group('getDrivers', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<DriverInfo>(
        data: [_makeDriver()],
        pagination: _makePagination(),
      );
      mockRepo.getDriversResult = expected;
      final params = GetDriversParams(page: 1, pageSize: 10);

      final result = await useCases.getDrivers(params);

      expect(result.data.length, 1);
      expect(result.data[0].firstName, 'TestDriver');
      expect(mockRepo.calls, ['getDrivers']);
      expect(mockRepo.capturedGetDriversParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getDriversResult = PaginatedResponse<DriverInfo>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getDrivers(null);

      expect(mockRepo.capturedGetDriversParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getDrivers(null), throwsException);
    });
  });

  group('getDriverById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getDriverByIdResult = _makeDriver(id: 42, firstName: 'Found');

      final result = await useCases.getDriverById(42);

      expect(result.id, 42);
      expect(result.firstName, 'Found');
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getDriverById']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getDriverByIdResult = _makeDriver();

      await useCases.getDriverById(1, businessId: 5);

      expect(mockRepo.capturedBusinessId, 5);
    });
  });

  group('createDriver', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateDriverDTO(
        firstName: 'New',
        lastName: 'Driver',
        phone: '555',
        identification: 'CC-999',
      );
      mockRepo.createDriverResult = _makeDriver(id: 99, firstName: 'New');

      final result = await useCases.createDriver(dto);

      expect(result.id, 99);
      expect(result.firstName, 'New');
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createDriver']);
    });

    test('passes businessId to repository', () async {
      final dto = CreateDriverDTO(
        firstName: 'Test',
        lastName: 'Test',
        phone: '555',
        identification: 'CC-1',
      );
      mockRepo.createDriverResult = _makeDriver();

      await useCases.createDriver(dto, businessId: 7);

      expect(mockRepo.capturedBusinessId, 7);
    });
  });

  group('updateDriver', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateDriverDTO(firstName: 'Updated');
      mockRepo.updateDriverResult = _makeDriver(id: 5, firstName: 'Updated');

      final result = await useCases.updateDriver(5, dto);

      expect(result.id, 5);
      expect(result.firstName, 'Updated');
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateDriver']);
    });

    test('passes businessId to repository', () async {
      final dto = UpdateDriverDTO(firstName: 'Test');
      mockRepo.updateDriverResult = _makeDriver();

      await useCases.updateDriver(1, dto, businessId: 3);

      expect(mockRepo.capturedBusinessId, 3);
    });
  });

  group('deleteDriver', () {
    test('delegates to repository with correct id', () async {
      mockRepo.deleteDriverResult = {'message': 'deleted'};

      final result = await useCases.deleteDriver(7);

      expect(result['message'], 'deleted');
      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deleteDriver']);
    });

    test('passes businessId to repository', () async {
      await useCases.deleteDriver(7, businessId: 2);

      expect(mockRepo.capturedBusinessId, 2);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.deleteDriver(7), throwsException);
    });
  });
}
