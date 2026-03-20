import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/vehicles/app/use_cases.dart';
import 'package:mobile_central/services/modules/vehicles/domain/entities.dart';
import 'package:mobile_central/services/modules/vehicles/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockVehicleRepository implements IVehicleRepository {
  PaginatedResponse<VehicleInfo>? getVehiclesResult;
  VehicleInfo? getVehicleByIdResult;
  VehicleInfo? createVehicleResult;
  VehicleInfo? updateVehicleResult;
  Map<String, dynamic>? deleteVehicleResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<VehicleInfo>> getVehicles(GetVehiclesParams? params) async {
    calls.add('getVehicles');
    if (errorToThrow != null) throw errorToThrow!;
    return getVehiclesResult!;
  }

  @override
  Future<VehicleInfo> getVehicleById(int id, {int? businessId}) async {
    calls.add('getVehicleById');
    if (errorToThrow != null) throw errorToThrow!;
    return getVehicleByIdResult!;
  }

  @override
  Future<VehicleInfo> createVehicle(CreateVehicleDTO data, {int? businessId}) async {
    calls.add('createVehicle');
    if (errorToThrow != null) throw errorToThrow!;
    return createVehicleResult!;
  }

  @override
  Future<VehicleInfo> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId}) async {
    calls.add('updateVehicle');
    if (errorToThrow != null) throw errorToThrow!;
    return updateVehicleResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteVehicle(int id, {int? businessId}) async {
    calls.add('deleteVehicle');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteVehicleResult ?? {};
  }
}

// --- Testable Provider ---

class TestableVehicleProvider {
  final VehicleUseCases _useCases;

  List<VehicleInfo> _vehicles = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableVehicleProvider(this._useCases);

  List<VehicleInfo> get vehicles => _vehicles;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchVehicles({GetVehiclesParams? params}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.getVehicles(params);
      _vehicles = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<VehicleInfo?> createVehicle(CreateVehicleDTO data, {int? businessId}) async {
    try {
      return await _useCases.createVehicle(data, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateVehicle(int id, UpdateVehicleDTO data, {int? businessId}) async {
    try {
      await _useCases.updateVehicle(id, data, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteVehicle(int id, {int? businessId}) async {
    try {
      await _useCases.deleteVehicle(id, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
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
  late MockVehicleRepository mockRepo;
  late VehicleUseCases useCases;
  late TestableVehicleProvider provider;

  setUp(() {
    mockRepo = MockVehicleRepository();
    useCases = VehicleUseCases(mockRepo);
    provider = TestableVehicleProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty vehicles list', () {
      expect(provider.vehicles, isEmpty);
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

  group('fetchVehicles', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getVehiclesResult = PaginatedResponse<VehicleInfo>(
        data: [_makeVehicle()],
        pagination: _makePagination(),
      );

      await provider.fetchVehicles();

      expect(provider.notifications.length, 2);
    });

    test('populates vehicles and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getVehiclesResult = PaginatedResponse<VehicleInfo>(
        data: [_makeVehicle(id: 1), _makeVehicle(id: 2, type: 'van')],
        pagination: pagination,
      );

      await provider.fetchVehicles();

      expect(provider.vehicles.length, 2);
      expect(provider.vehicles[0].id, 1);
      expect(provider.vehicles[1].type, 'van');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchVehicles();

      expect(provider.error, contains('Server error'));
      expect(provider.vehicles, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchVehicles();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getVehiclesResult = PaginatedResponse<VehicleInfo>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchVehicles();

      expect(provider.error, isNull);
    });
  });

  group('createVehicle', () {
    test('returns created vehicle on success', () async {
      final dto = CreateVehicleDTO(type: 'truck', licensePlate: 'NEW-001');
      mockRepo.createVehicleResult = _makeVehicle(id: 10, type: 'truck');

      final result = await provider.createVehicle(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateVehicleDTO(type: 'x', licensePlate: 'y');

      final result = await provider.createVehicle(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateVehicle', () {
    test('returns true on success', () async {
      final dto = UpdateVehicleDTO(status: 'inactive');
      mockRepo.updateVehicleResult = _makeVehicle(id: 5);

      final result = await provider.updateVehicle(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateVehicleDTO(color: 'Red');

      final result = await provider.updateVehicle(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteVehicle', () {
    test('returns true on success', () async {
      final result = await provider.deleteVehicle(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteVehicle(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });
}
