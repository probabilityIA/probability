import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/drivers/app/use_cases.dart';
import 'package:mobile_central/services/modules/drivers/domain/entities.dart';
import 'package:mobile_central/services/modules/drivers/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockDriverRepository implements IDriverRepository {
  PaginatedResponse<DriverInfo>? getDriversResult;
  DriverInfo? getDriverByIdResult;
  DriverInfo? createDriverResult;
  DriverInfo? updateDriverResult;
  Map<String, dynamic>? deleteDriverResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<DriverInfo>> getDrivers(GetDriversParams? params) async {
    calls.add('getDrivers');
    if (errorToThrow != null) throw errorToThrow!;
    return getDriversResult!;
  }

  @override
  Future<DriverInfo> getDriverById(int id, {int? businessId}) async {
    calls.add('getDriverById');
    if (errorToThrow != null) throw errorToThrow!;
    return getDriverByIdResult!;
  }

  @override
  Future<DriverInfo> createDriver(CreateDriverDTO data, {int? businessId}) async {
    calls.add('createDriver');
    if (errorToThrow != null) throw errorToThrow!;
    return createDriverResult!;
  }

  @override
  Future<DriverInfo> updateDriver(int id, UpdateDriverDTO data, {int? businessId}) async {
    calls.add('updateDriver');
    if (errorToThrow != null) throw errorToThrow!;
    return updateDriverResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteDriver(int id, {int? businessId}) async {
    calls.add('deleteDriver');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteDriverResult ?? {'message': 'deleted'};
  }
}

// --- Testable Provider ---

class TestableDriverProvider {
  final DriverUseCases _useCases;

  List<DriverInfo> _drivers = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;
  final int _pageSize = 20;
  String _searchFilter = '';
  String _statusFilter = '';

  final List<String> notifications = [];

  TestableDriverProvider(this._useCases);

  List<DriverInfo> get drivers => _drivers;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;
  int get page => _page;
  int get pageSize => _pageSize;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchDrivers({int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final params = GetDriversParams(
        page: _page,
        pageSize: _pageSize,
        search: _searchFilter.isNotEmpty ? _searchFilter : null,
        status: _statusFilter.isNotEmpty ? _statusFilter : null,
        businessId: businessId,
      );
      final response = await _useCases.getDrivers(params);
      _drivers = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }

  Future<DriverInfo?> getDriverById(int id, {int? businessId}) async {
    try {
      return await _useCases.getDriverById(id, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<DriverInfo?> createDriver(CreateDriverDTO data, {int? businessId}) async {
    try {
      final driver = await _useCases.createDriver(data, businessId: businessId);
      return driver;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> updateDriver(int id, UpdateDriverDTO data, {int? businessId}) async {
    try {
      await _useCases.updateDriver(id, data, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<bool> deleteDriver(int id, {int? businessId}) async {
    try {
      await _useCases.deleteDriver(id, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  void setPage(int page) {
    _page = page;
  }

  void setFilters({String? search, String? status}) {
    _searchFilter = search ?? _searchFilter;
    _statusFilter = status ?? _statusFilter;
    _page = 1;
  }

  void resetFilters() {
    _searchFilter = '';
    _statusFilter = '';
    _page = 1;
  }
}

// --- Helpers ---

Pagination _makePagination({
  int currentPage = 1,
  int total = 5,
  int lastPage = 1,
}) {
  return Pagination(
    currentPage: currentPage,
    perPage: 20,
    total: total,
    lastPage: lastPage,
    hasNext: currentPage < lastPage,
    hasPrev: currentPage > 1,
  );
}

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

// --- Tests ---

void main() {
  late MockDriverRepository mockRepo;
  late DriverUseCases useCases;
  late TestableDriverProvider provider;

  setUp(() {
    mockRepo = MockDriverRepository();
    useCases = DriverUseCases(mockRepo);
    provider = TestableDriverProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty drivers list', () {
      expect(provider.drivers, isEmpty);
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

    test('starts at page 1', () {
      expect(provider.page, 1);
    });

    test('has pageSize of 20', () {
      expect(provider.pageSize, 20);
    });
  });

  group('fetchDrivers', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getDriversResult = PaginatedResponse<DriverInfo>(
        data: [_makeDriver()],
        pagination: _makePagination(),
      );

      await provider.fetchDrivers();

      expect(provider.notifications.length, 2);
    });

    test('populates drivers and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getDriversResult = PaginatedResponse<DriverInfo>(
        data: [_makeDriver(id: 1), _makeDriver(id: 2, firstName: 'Carlos')],
        pagination: pagination,
      );

      await provider.fetchDrivers();

      expect(provider.drivers.length, 2);
      expect(provider.drivers[0].id, 1);
      expect(provider.drivers[1].firstName, 'Carlos');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchDrivers();

      expect(provider.error, contains('Server error'));
      expect(provider.drivers, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchDrivers();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getDriversResult = PaginatedResponse<DriverInfo>(
        data: [],
        pagination: _makePagination(),
      );
      await provider.fetchDrivers();

      expect(provider.error, isNull);
    });
  });

  group('getDriverById', () {
    test('returns driver on success', () async {
      mockRepo.getDriverByIdResult = _makeDriver(id: 42, firstName: 'Found');

      final result = await provider.getDriverById(42);

      expect(result, isNotNull);
      expect(result!.id, 42);
      expect(result.firstName, 'Found');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Not found');

      final result = await provider.getDriverById(42);

      expect(result, isNull);
      expect(provider.error, contains('Not found'));
    });
  });

  group('createDriver', () {
    test('returns created driver on success', () async {
      final dto = CreateDriverDTO(
        firstName: 'New',
        lastName: 'Driver',
        phone: '555',
        identification: 'CC-999',
      );
      mockRepo.createDriverResult = _makeDriver(id: 10, firstName: 'New');

      final result = await provider.createDriver(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
      expect(result.firstName, 'New');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateDriverDTO(
        firstName: 'Fail',
        lastName: 'Test',
        phone: '555',
        identification: 'CC-0',
      );

      final result = await provider.createDriver(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateDriver', () {
    test('returns true on success', () async {
      final dto = UpdateDriverDTO(firstName: 'Updated');
      mockRepo.updateDriverResult = _makeDriver(id: 5, firstName: 'Updated');

      final result = await provider.updateDriver(5, dto);

      expect(result, true);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateDriverDTO(firstName: 'Fail');

      final result = await provider.updateDriver(5, dto);

      expect(result, false);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteDriver', () {
    test('returns true on success', () async {
      final result = await provider.deleteDriver(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteDriver(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });

  group('setPage', () {
    test('updates page value', () {
      provider.setPage(3);

      expect(provider.page, 3);
    });
  });

  group('setFilters', () {
    test('updates search filter and resets page to 1', () {
      provider.setPage(5);
      provider.setFilters(search: 'Carlos');

      expect(provider.page, 1);
    });

    test('updates status filter and resets page to 1', () {
      provider.setPage(3);
      provider.setFilters(status: 'active');

      expect(provider.page, 1);
    });

    test('updates both filters simultaneously', () {
      provider.setPage(5);
      provider.setFilters(search: 'Carlos', status: 'active');

      expect(provider.page, 1);
    });
  });

  group('resetFilters', () {
    test('resets all filters and page', () {
      provider.setPage(3);
      provider.setFilters(search: 'test', status: 'active');
      provider.resetFilters();

      expect(provider.page, 1);
    });
  });
}
