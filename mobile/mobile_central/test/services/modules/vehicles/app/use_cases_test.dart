import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/vehicles/app/use_cases.dart';
import 'package:mobile_central/services/modules/vehicles/domain/entities.dart';
import 'package:mobile_central/services/modules/vehicles/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockVehicleRepository implements IVehicleRepository {
  final List<String> calls = [];

  PaginatedResponse<VehicleInfo>? getVehiclesResult;
  VehicleInfo? getVehicleByIdResult;
  VehicleInfo? createVehicleResult;
  VehicleInfo? updateVehicleResult;
  Map<String, dynamic>? deleteVehicleResult;

  Exception? errorToThrow;

  GetVehiclesParams? capturedParams;
  int? capturedId;
  int? capturedBusinessId;
  CreateVehicleDTO? capturedCreateData;
  int? capturedUpdateId;
  UpdateVehicleDTO? capturedUpdateData;
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<VehicleInfo>> getVehicles(GetVehiclesParams? params) async {
    calls.add('getVehicles');
    capturedParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getVehiclesResult!;
  }

  @override
  Future<VehicleInfo> getVehicleById(int id, {int? businessId}) async {
    calls.add('getVehicleById');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getVehicleByIdResult!;
  }

  @override
  Future<VehicleInfo> createVehicle(CreateVehicleDTO data, {int? businessId}) async {
    calls.add('createVehicle');
    capturedCreateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createVehicleResult!;
  }

  @override
  Future<VehicleInfo> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId}) async {
    calls.add('updateVehicle');
    capturedUpdateId = id;
    capturedUpdateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateVehicleResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteVehicle(int id, {int? businessId}) async {
    calls.add('deleteVehicle');
    capturedDeleteId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteVehicleResult ?? {};
  }
}

// --- Helpers ---

VehicleInfo _makeVehicle({int id = 1, String type = 'truck'}) {
  return VehicleInfo(
    id: id,
    businessId: 1,
    type: type,
    licensePlate: 'ABC-123',
    brand: 'Toyota',
    model: 'Hilux',
    color: 'White',
    status: 'active',
    photoUrl: '',
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
  late MockVehicleRepository mockRepo;
  late VehicleUseCases useCases;

  setUp(() {
    mockRepo = MockVehicleRepository();
    useCases = VehicleUseCases(mockRepo);
  });

  group('getVehicles', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<VehicleInfo>(
        data: [_makeVehicle()],
        pagination: _makePagination(),
      );
      mockRepo.getVehiclesResult = expected;
      final params = GetVehiclesParams(page: 1, pageSize: 10);

      final result = await useCases.getVehicles(params);

      expect(result.data.length, 1);
      expect(result.data[0].type, 'truck');
      expect(mockRepo.calls, ['getVehicles']);
      expect(mockRepo.capturedParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getVehiclesResult = PaginatedResponse<VehicleInfo>(
        data: [],
        pagination: _makePagination(),
      );

      await useCases.getVehicles(null);

      expect(mockRepo.capturedParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getVehicles(null), throwsException);
    });
  });

  group('getVehicleById', () {
    test('delegates to repository with correct id', () async {
      mockRepo.getVehicleByIdResult = _makeVehicle(id: 42, type: 'van');

      final result = await useCases.getVehicleById(42);

      expect(result.id, 42);
      expect(result.type, 'van');
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.calls, ['getVehicleById']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getVehicleByIdResult = _makeVehicle();

      await useCases.getVehicleById(1, businessId: 5);

      expect(mockRepo.capturedBusinessId, 5);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.getVehicleById(99), throwsException);
    });
  });

  group('createVehicle', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateVehicleDTO(type: 'truck', licensePlate: 'NEW-001');
      mockRepo.createVehicleResult = _makeVehicle(id: 10);

      final result = await useCases.createVehicle(dto);

      expect(result.id, 10);
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.calls, ['createVehicle']);
    });

    test('passes businessId to repository', () async {
      final dto = CreateVehicleDTO(type: 'van', licensePlate: 'V-001');
      mockRepo.createVehicleResult = _makeVehicle();

      await useCases.createVehicle(dto, businessId: 7);

      expect(mockRepo.capturedBusinessId, 7);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateVehicleDTO(type: 'x', licensePlate: 'y');

      expect(() => useCases.createVehicle(dto), throwsException);
    });
  });

  group('updateVehicle', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateVehicleDTO(status: 'inactive');
      mockRepo.updateVehicleResult = _makeVehicle(id: 5);

      final result = await useCases.updateVehicle(5, dto);

      expect(result.id, 5);
      expect(mockRepo.capturedUpdateId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.calls, ['updateVehicle']);
    });

    test('passes businessId to repository', () async {
      final dto = UpdateVehicleDTO(color: 'Red');
      mockRepo.updateVehicleResult = _makeVehicle();

      await useCases.updateVehicle(1, dto, businessId: 3);

      expect(mockRepo.capturedBusinessId, 3);
    });
  });

  group('deleteVehicle', () {
    test('delegates to repository with correct id', () async {
      mockRepo.deleteVehicleResult = {'message': 'deleted'};

      final result = await useCases.deleteVehicle(7);

      expect(result['message'], 'deleted');
      expect(mockRepo.capturedDeleteId, 7);
      expect(mockRepo.calls, ['deleteVehicle']);
    });

    test('passes businessId to repository', () async {
      await useCases.deleteVehicle(1, businessId: 2);

      expect(mockRepo.capturedBusinessId, 2);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      expect(() => useCases.deleteVehicle(7), throwsException);
    });
  });
}
