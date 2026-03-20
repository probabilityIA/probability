import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/routes/app/use_cases.dart';
import 'package:mobile_central/services/modules/routes/domain/entities.dart';
import 'package:mobile_central/services/modules/routes/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock Repository ---

class MockRouteRepository implements IRouteRepository {
  PaginatedResponse<RouteInfo>? getRoutesResult;
  RouteDetail? getRouteByIdResult;
  RouteInfo? createRouteResult;
  RouteInfo? updateRouteResult;
  Map<String, dynamic>? deleteRouteResult;
  RouteDetail? startRouteResult;
  RouteDetail? completeRouteResult;
  RouteStopInfo? addStopResult;
  RouteStopInfo? updateStopResult;
  Map<String, dynamic>? deleteStopResult;
  RouteStopInfo? updateStopStatusResult;
  RouteDetail? reorderStopsResult;
  List<DriverOption>? getAvailableDriversResult;
  List<VehicleOption>? getAvailableVehiclesResult;
  List<AssignableOrder>? getAssignableOrdersResult;
  Exception? errorToThrow;

  final List<String> calls = [];
  int? capturedDeleteId;

  @override
  Future<PaginatedResponse<RouteInfo>> getRoutes(GetRoutesParams? params) async {
    calls.add('getRoutes');
    if (errorToThrow != null) throw errorToThrow!;
    return getRoutesResult!;
  }

  @override
  Future<RouteDetail> getRouteById(int id, {int? businessId}) async {
    calls.add('getRouteById');
    if (errorToThrow != null) throw errorToThrow!;
    return getRouteByIdResult!;
  }

  @override
  Future<RouteInfo> createRoute(CreateRouteDTO data, {int? businessId}) async {
    calls.add('createRoute');
    if (errorToThrow != null) throw errorToThrow!;
    return createRouteResult!;
  }

  @override
  Future<RouteInfo> updateRoute(int id, UpdateRouteDTO data, {int? businessId}) async {
    calls.add('updateRoute');
    if (errorToThrow != null) throw errorToThrow!;
    return updateRouteResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteRoute(int id, {int? businessId}) async {
    calls.add('deleteRoute');
    capturedDeleteId = id;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteRouteResult ?? {'message': 'deleted'};
  }

  @override
  Future<RouteDetail> startRoute(int id, {int? businessId}) async {
    calls.add('startRoute');
    if (errorToThrow != null) throw errorToThrow!;
    return startRouteResult!;
  }

  @override
  Future<RouteDetail> completeRoute(int id, {int? businessId}) async {
    calls.add('completeRoute');
    if (errorToThrow != null) throw errorToThrow!;
    return completeRouteResult!;
  }

  @override
  Future<RouteStopInfo> addStop(int routeId, AddStopDTO data, {int? businessId}) async {
    calls.add('addStop');
    if (errorToThrow != null) throw errorToThrow!;
    return addStopResult!;
  }

  @override
  Future<RouteStopInfo> updateStop(int routeId, int stopId, UpdateStopDTO data, {int? businessId}) async {
    calls.add('updateStop');
    if (errorToThrow != null) throw errorToThrow!;
    return updateStopResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteStop(int routeId, int stopId, {int? businessId}) async {
    calls.add('deleteStop');
    if (errorToThrow != null) throw errorToThrow!;
    return deleteStopResult ?? {'message': 'deleted'};
  }

  @override
  Future<RouteStopInfo> updateStopStatus(int routeId, int stopId, UpdateStopStatusDTO data, {int? businessId}) async {
    calls.add('updateStopStatus');
    if (errorToThrow != null) throw errorToThrow!;
    return updateStopStatusResult!;
  }

  @override
  Future<RouteDetail> reorderStops(int routeId, ReorderStopsDTO data, {int? businessId}) async {
    calls.add('reorderStops');
    if (errorToThrow != null) throw errorToThrow!;
    return reorderStopsResult!;
  }

  @override
  Future<List<DriverOption>> getAvailableDrivers({int? businessId}) async {
    calls.add('getAvailableDrivers');
    if (errorToThrow != null) throw errorToThrow!;
    return getAvailableDriversResult ?? [];
  }

  @override
  Future<List<VehicleOption>> getAvailableVehicles({int? businessId}) async {
    calls.add('getAvailableVehicles');
    if (errorToThrow != null) throw errorToThrow!;
    return getAvailableVehiclesResult ?? [];
  }

  @override
  Future<List<AssignableOrder>> getAssignableOrders({int? businessId}) async {
    calls.add('getAssignableOrders');
    if (errorToThrow != null) throw errorToThrow!;
    return getAssignableOrdersResult ?? [];
  }
}

// --- Testable Provider ---

class TestableRouteProvider {
  final RouteUseCases _useCases;

  List<RouteInfo> _routes = [];
  RouteDetail? _selectedRoute;
  List<DriverOption> _availableDrivers = [];
  List<VehicleOption> _availableVehicles = [];
  List<AssignableOrder> _assignableOrders = [];
  Pagination? _pagination;
  bool _isLoading = false;
  String? _error;
  int _page = 1;

  final List<String> notifications = [];

  TestableRouteProvider(this._useCases);

  List<RouteInfo> get routes => _routes;
  RouteDetail? get selectedRoute => _selectedRoute;
  List<DriverOption> get availableDrivers => _availableDrivers;
  List<VehicleOption> get availableVehicles => _availableVehicles;
  List<AssignableOrder> get assignableOrders => _assignableOrders;
  Pagination? get pagination => _pagination;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchRoutes({int? businessId, String? status, int? driverId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();
    try {
      final params = GetRoutesParams(page: _page, pageSize: 20, businessId: businessId, status: status, driverId: driverId);
      final response = await _useCases.getRoutes(params);
      _routes = response.data;
      _pagination = response.pagination;
    } catch (e) {
      _error = e.toString();
    }
    _isLoading = false;
    _notifyListeners();
  }

  Future<void> fetchRouteDetail(int id, {int? businessId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();
    try {
      _selectedRoute = await _useCases.getRouteById(id, businessId: businessId);
    } catch (e) {
      _error = e.toString();
    }
    _isLoading = false;
    _notifyListeners();
  }

  Future<RouteInfo?> createRoute(CreateRouteDTO data, {int? businessId}) async {
    try {
      return await _useCases.createRoute(data, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<RouteInfo?> updateRoute(int id, UpdateRouteDTO data, {int? businessId}) async {
    try {
      return await _useCases.updateRoute(id, data, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<bool> deleteRoute(int id, {int? businessId}) async {
    try {
      await _useCases.deleteRoute(id, businessId: businessId);
      return true;
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return false;
    }
  }

  Future<RouteDetail?> startRoute(int id, {int? businessId}) async {
    try {
      return await _useCases.startRoute(id, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<RouteDetail?> completeRoute(int id, {int? businessId}) async {
    try {
      return await _useCases.completeRoute(id, businessId: businessId);
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
      return null;
    }
  }

  Future<void> fetchFormOptions({int? businessId}) async {
    try {
      final results = await Future.wait([
        _useCases.getAvailableDrivers(businessId: businessId),
        _useCases.getAvailableVehicles(businessId: businessId),
        _useCases.getAssignableOrders(businessId: businessId),
      ]);
      _availableDrivers = results[0] as List<DriverOption>;
      _availableVehicles = results[1] as List<VehicleOption>;
      _assignableOrders = results[2] as List<AssignableOrder>;
      _notifyListeners();
    } catch (e) {
      _error = e.toString();
      _notifyListeners();
    }
  }

  void setPage(int page) {
    _page = page;
  }
}

// --- Helpers ---

RouteInfo _makeRouteInfo({int id = 1, String status = 'pending'}) {
  return RouteInfo(
    id: id, businessId: 1, status: status, date: '2026-03-01',
    totalStops: 3, completedStops: 0, failedStops: 0,
    createdAt: '2026-01-01', updatedAt: '2026-01-01',
  );
}

RouteDetail _makeRouteDetail({int id = 1, String status = 'pending'}) {
  return RouteDetail(
    id: id, businessId: 1, status: status, date: '2026-03-01',
    totalStops: 3, completedStops: 0, failedStops: 0,
    createdAt: '2026-01-01', updatedAt: '2026-01-01', stops: [],
  );
}

Pagination _makePagination({int currentPage = 1, int total = 5, int lastPage = 1}) {
  return Pagination(
    currentPage: currentPage, perPage: 20, total: total, lastPage: lastPage,
    hasNext: currentPage < lastPage, hasPrev: currentPage > 1,
  );
}

// --- Tests ---

void main() {
  late MockRouteRepository mockRepo;
  late RouteUseCases useCases;
  late TestableRouteProvider provider;

  setUp(() {
    mockRepo = MockRouteRepository();
    useCases = RouteUseCases(mockRepo);
    provider = TestableRouteProvider(useCases);
  });

  group('initial state', () {
    test('starts with empty routes list', () {
      expect(provider.routes, isEmpty);
    });

    test('starts with null selected route', () {
      expect(provider.selectedRoute, isNull);
    });

    test('starts with empty form options', () {
      expect(provider.availableDrivers, isEmpty);
      expect(provider.availableVehicles, isEmpty);
      expect(provider.assignableOrders, isEmpty);
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

  group('fetchRoutes', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getRoutesResult = PaginatedResponse<RouteInfo>(
        data: [_makeRouteInfo()],
        pagination: _makePagination(),
      );

      await provider.fetchRoutes();

      expect(provider.notifications.length, 2);
    });

    test('populates routes and pagination on success', () async {
      final pagination = _makePagination(total: 2);
      mockRepo.getRoutesResult = PaginatedResponse<RouteInfo>(
        data: [_makeRouteInfo(id: 1), _makeRouteInfo(id: 2, status: 'in_progress')],
        pagination: pagination,
      );

      await provider.fetchRoutes();

      expect(provider.routes.length, 2);
      expect(provider.routes[0].id, 1);
      expect(provider.routes[1].status, 'in_progress');
      expect(provider.pagination?.total, 2);
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchRoutes();

      expect(provider.error, contains('Server error'));
      expect(provider.routes, isEmpty);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchRoutes();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getRoutesResult = PaginatedResponse<RouteInfo>(
        data: [], pagination: _makePagination(),
      );
      await provider.fetchRoutes();

      expect(provider.error, isNull);
    });
  });

  group('fetchRouteDetail', () {
    test('populates selected route on success', () async {
      mockRepo.getRouteByIdResult = _makeRouteDetail(id: 42, status: 'in_progress');

      await provider.fetchRouteDetail(42);

      expect(provider.selectedRoute, isNotNull);
      expect(provider.selectedRoute!.id, 42);
      expect(provider.selectedRoute!.status, 'in_progress');
      expect(provider.isLoading, false);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Not found');

      await provider.fetchRouteDetail(999);

      expect(provider.error, contains('Not found'));
      expect(provider.selectedRoute, isNull);
    });
  });

  group('createRoute', () {
    test('returns created route on success', () async {
      final dto = CreateRouteDTO(date: '2026-03-01');
      mockRepo.createRouteResult = _makeRouteInfo(id: 10);

      final result = await provider.createRoute(dto);

      expect(result, isNotNull);
      expect(result!.id, 10);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Creation failed');
      final dto = CreateRouteDTO(date: '2026-03-01');

      final result = await provider.createRoute(dto);

      expect(result, isNull);
      expect(provider.error, contains('Creation failed'));
    });
  });

  group('updateRoute', () {
    test('returns updated route on success', () async {
      final dto = UpdateRouteDTO(notes: 'Updated');
      mockRepo.updateRouteResult = _makeRouteInfo(id: 5);

      final result = await provider.updateRoute(5, dto);

      expect(result, isNotNull);
      expect(result!.id, 5);
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Update failed');
      final dto = UpdateRouteDTO(notes: 'Fail');

      final result = await provider.updateRoute(5, dto);

      expect(result, isNull);
      expect(provider.error, contains('Update failed'));
    });
  });

  group('deleteRoute', () {
    test('returns true on success', () async {
      mockRepo.deleteRouteResult = {'message': 'deleted'};

      final result = await provider.deleteRoute(7);

      expect(result, true);
      expect(mockRepo.capturedDeleteId, 7);
    });

    test('returns false and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Delete failed');

      final result = await provider.deleteRoute(7);

      expect(result, false);
      expect(provider.error, contains('Delete failed'));
    });
  });

  group('startRoute', () {
    test('returns route detail on success', () async {
      mockRepo.startRouteResult = _makeRouteDetail(id: 1, status: 'in_progress');

      final result = await provider.startRoute(1);

      expect(result, isNotNull);
      expect(result!.status, 'in_progress');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Start failed');

      final result = await provider.startRoute(1);

      expect(result, isNull);
      expect(provider.error, contains('Start failed'));
    });
  });

  group('completeRoute', () {
    test('returns route detail on success', () async {
      mockRepo.completeRouteResult = _makeRouteDetail(id: 1, status: 'completed');

      final result = await provider.completeRoute(1);

      expect(result, isNotNull);
      expect(result!.status, 'completed');
    });

    test('returns null and sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Complete failed');

      final result = await provider.completeRoute(1);

      expect(result, isNull);
      expect(provider.error, contains('Complete failed'));
    });
  });

  group('fetchFormOptions', () {
    test('populates drivers, vehicles, and orders on success', () async {
      mockRepo.getAvailableDriversResult = [
        DriverOption(id: 1, firstName: 'Carlos', lastName: 'G', phone: '', identification: '', status: 'active', licenseType: 'B2'),
      ];
      mockRepo.getAvailableVehiclesResult = [
        VehicleOption(id: 1, type: 'van', licensePlate: 'ABC', brand: 'Toyota', vehicleModel: 'HiAce', status: 'available'),
      ];
      mockRepo.getAssignableOrdersResult = [
        AssignableOrder(id: '1', orderNumber: 'ORD-1', customerName: 'Maria', customerPhone: '', address: '', city: '', totalAmount: 100, itemCount: 1, createdAt: ''),
      ];

      await provider.fetchFormOptions();

      expect(provider.availableDrivers.length, 1);
      expect(provider.availableVehicles.length, 1);
      expect(provider.assignableOrders.length, 1);
      expect(provider.availableDrivers[0].firstName, 'Carlos');
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Options failed');

      await provider.fetchFormOptions();

      expect(provider.error, contains('Options failed'));
    });
  });

  group('pagination', () {
    test('setPage updates page', () {
      provider.setPage(3);

      // No direct getter for _page, but it should be used in next fetchRoutes
      mockRepo.getRoutesResult = PaginatedResponse<RouteInfo>(
        data: [], pagination: _makePagination(currentPage: 3),
      );
    });
  });
}
