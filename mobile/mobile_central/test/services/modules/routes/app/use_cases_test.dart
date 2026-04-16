import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/routes/app/use_cases.dart';
import 'package:mobile_central/services/modules/routes/domain/entities.dart';
import 'package:mobile_central/services/modules/routes/domain/ports.dart';
import 'package:mobile_central/shared/types/paginated_response.dart';

// --- Manual Mock ---

class MockRouteRepository implements IRouteRepository {
  final List<String> calls = [];

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

  GetRoutesParams? capturedGetRoutesParams;
  int? capturedId;
  int? capturedBusinessId;
  CreateRouteDTO? capturedCreateData;
  UpdateRouteDTO? capturedUpdateData;
  int? capturedRouteId;
  int? capturedStopId;
  AddStopDTO? capturedAddStopData;
  UpdateStopDTO? capturedUpdateStopData;
  UpdateStopStatusDTO? capturedStopStatusData;
  ReorderStopsDTO? capturedReorderData;

  @override
  Future<PaginatedResponse<RouteInfo>> getRoutes(GetRoutesParams? params) async {
    calls.add('getRoutes');
    capturedGetRoutesParams = params;
    if (errorToThrow != null) throw errorToThrow!;
    return getRoutesResult!;
  }

  @override
  Future<RouteDetail> getRouteById(int id, {int? businessId}) async {
    calls.add('getRouteById');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getRouteByIdResult!;
  }

  @override
  Future<RouteInfo> createRoute(CreateRouteDTO data, {int? businessId}) async {
    calls.add('createRoute');
    capturedCreateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return createRouteResult!;
  }

  @override
  Future<RouteInfo> updateRoute(int id, UpdateRouteDTO data, {int? businessId}) async {
    calls.add('updateRoute');
    capturedId = id;
    capturedUpdateData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateRouteResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteRoute(int id, {int? businessId}) async {
    calls.add('deleteRoute');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteRouteResult!;
  }

  @override
  Future<RouteDetail> startRoute(int id, {int? businessId}) async {
    calls.add('startRoute');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return startRouteResult!;
  }

  @override
  Future<RouteDetail> completeRoute(int id, {int? businessId}) async {
    calls.add('completeRoute');
    capturedId = id;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return completeRouteResult!;
  }

  @override
  Future<RouteStopInfo> addStop(int routeId, AddStopDTO data, {int? businessId}) async {
    calls.add('addStop');
    capturedRouteId = routeId;
    capturedAddStopData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return addStopResult!;
  }

  @override
  Future<RouteStopInfo> updateStop(int routeId, int stopId, UpdateStopDTO data, {int? businessId}) async {
    calls.add('updateStop');
    capturedRouteId = routeId;
    capturedStopId = stopId;
    capturedUpdateStopData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateStopResult!;
  }

  @override
  Future<Map<String, dynamic>> deleteStop(int routeId, int stopId, {int? businessId}) async {
    calls.add('deleteStop');
    capturedRouteId = routeId;
    capturedStopId = stopId;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return deleteStopResult!;
  }

  @override
  Future<RouteStopInfo> updateStopStatus(int routeId, int stopId, UpdateStopStatusDTO data, {int? businessId}) async {
    calls.add('updateStopStatus');
    capturedRouteId = routeId;
    capturedStopId = stopId;
    capturedStopStatusData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return updateStopStatusResult!;
  }

  @override
  Future<RouteDetail> reorderStops(int routeId, ReorderStopsDTO data, {int? businessId}) async {
    calls.add('reorderStops');
    capturedRouteId = routeId;
    capturedReorderData = data;
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return reorderStopsResult!;
  }

  @override
  Future<List<DriverOption>> getAvailableDrivers({int? businessId}) async {
    calls.add('getAvailableDrivers');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getAvailableDriversResult!;
  }

  @override
  Future<List<VehicleOption>> getAvailableVehicles({int? businessId}) async {
    calls.add('getAvailableVehicles');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getAvailableVehiclesResult!;
  }

  @override
  Future<List<AssignableOrder>> getAssignableOrders({int? businessId}) async {
    calls.add('getAssignableOrders');
    capturedBusinessId = businessId;
    if (errorToThrow != null) throw errorToThrow!;
    return getAssignableOrdersResult!;
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

RouteStopInfo _makeStop({int id = 1, String status = 'pending'}) {
  return RouteStopInfo(
    id: id, routeId: 1, sequence: 1, status: status, address: 'Addr',
    customerName: 'Customer', createdAt: '2026-01-01', updatedAt: '2026-01-01',
  );
}

Pagination _makePagination() {
  return Pagination(
    currentPage: 1, perPage: 10, total: 1, lastPage: 1,
    hasNext: false, hasPrev: false,
  );
}

// --- Tests ---

void main() {
  late MockRouteRepository mockRepo;
  late RouteUseCases useCases;

  setUp(() {
    mockRepo = MockRouteRepository();
    useCases = RouteUseCases(mockRepo);
  });

  group('getRoutes', () {
    test('delegates to repository and returns result', () async {
      final expected = PaginatedResponse<RouteInfo>(
        data: [_makeRouteInfo()],
        pagination: _makePagination(),
      );
      mockRepo.getRoutesResult = expected;
      final params = GetRoutesParams(page: 1, pageSize: 10);

      final result = await useCases.getRoutes(params);

      expect(result.data.length, 1);
      expect(mockRepo.calls, ['getRoutes']);
      expect(mockRepo.capturedGetRoutesParams, params);
    });

    test('passes null params to repository', () async {
      mockRepo.getRoutesResult = PaginatedResponse<RouteInfo>(
        data: [], pagination: _makePagination(),
      );

      await useCases.getRoutes(null);

      expect(mockRepo.capturedGetRoutesParams, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getRoutes(null), throwsException);
    });
  });

  group('getRouteById', () {
    test('delegates to repository with correct id and businessId', () async {
      mockRepo.getRouteByIdResult = _makeRouteDetail(id: 42);

      final result = await useCases.getRouteById(42, businessId: 5);

      expect(result.id, 42);
      expect(mockRepo.capturedId, 42);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['getRouteById']);
    });
  });

  group('createRoute', () {
    test('delegates to repository with correct data', () async {
      final dto = CreateRouteDTO(date: '2026-03-01');
      mockRepo.createRouteResult = _makeRouteInfo(id: 99);

      final result = await useCases.createRoute(dto, businessId: 5);

      expect(result.id, 99);
      expect(mockRepo.capturedCreateData, dto);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['createRoute']);
    });
  });

  group('updateRoute', () {
    test('delegates to repository with correct id and data', () async {
      final dto = UpdateRouteDTO(notes: 'Updated');
      mockRepo.updateRouteResult = _makeRouteInfo(id: 5);

      final result = await useCases.updateRoute(5, dto, businessId: 2);

      expect(result.id, 5);
      expect(mockRepo.capturedId, 5);
      expect(mockRepo.capturedUpdateData, dto);
      expect(mockRepo.capturedBusinessId, 2);
      expect(mockRepo.calls, ['updateRoute']);
    });
  });

  group('deleteRoute', () {
    test('delegates to repository with correct id', () async {
      mockRepo.deleteRouteResult = {'message': 'deleted'};

      final result = await useCases.deleteRoute(7, businessId: 3);

      expect(result['message'], 'deleted');
      expect(mockRepo.capturedId, 7);
      expect(mockRepo.capturedBusinessId, 3);
      expect(mockRepo.calls, ['deleteRoute']);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Not found');

      expect(() => useCases.deleteRoute(7), throwsException);
    });
  });

  group('startRoute', () {
    test('delegates to repository with correct id', () async {
      mockRepo.startRouteResult = _makeRouteDetail(id: 1, status: 'in_progress');

      final result = await useCases.startRoute(1, businessId: 5);

      expect(result.status, 'in_progress');
      expect(mockRepo.capturedId, 1);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['startRoute']);
    });
  });

  group('completeRoute', () {
    test('delegates to repository with correct id', () async {
      mockRepo.completeRouteResult = _makeRouteDetail(id: 1, status: 'completed');

      final result = await useCases.completeRoute(1, businessId: 5);

      expect(result.status, 'completed');
      expect(mockRepo.capturedId, 1);
      expect(mockRepo.calls, ['completeRoute']);
    });
  });

  group('addStop', () {
    test('delegates to repository with correct data', () async {
      final dto = AddStopDTO(address: 'New Addr', customerName: 'Client');
      mockRepo.addStopResult = _makeStop(id: 10);

      final result = await useCases.addStop(1, dto, businessId: 5);

      expect(result.id, 10);
      expect(mockRepo.capturedRouteId, 1);
      expect(mockRepo.capturedAddStopData, dto);
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['addStop']);
    });
  });

  group('updateStop', () {
    test('delegates to repository with correct ids and data', () async {
      final dto = UpdateStopDTO(address: 'Updated');
      mockRepo.updateStopResult = _makeStop(id: 5);

      final result = await useCases.updateStop(1, 5, dto, businessId: 3);

      expect(result.id, 5);
      expect(mockRepo.capturedRouteId, 1);
      expect(mockRepo.capturedStopId, 5);
      expect(mockRepo.capturedUpdateStopData, dto);
      expect(mockRepo.capturedBusinessId, 3);
      expect(mockRepo.calls, ['updateStop']);
    });
  });

  group('deleteStop', () {
    test('delegates to repository with correct ids', () async {
      mockRepo.deleteStopResult = {'message': 'deleted'};

      final result = await useCases.deleteStop(1, 5, businessId: 3);

      expect(result['message'], 'deleted');
      expect(mockRepo.capturedRouteId, 1);
      expect(mockRepo.capturedStopId, 5);
      expect(mockRepo.calls, ['deleteStop']);
    });
  });

  group('updateStopStatus', () {
    test('delegates to repository with correct data', () async {
      final dto = UpdateStopStatusDTO(status: 'delivered');
      mockRepo.updateStopStatusResult = _makeStop(id: 5, status: 'delivered');

      final result = await useCases.updateStopStatus(1, 5, dto, businessId: 3);

      expect(result.status, 'delivered');
      expect(mockRepo.capturedRouteId, 1);
      expect(mockRepo.capturedStopId, 5);
      expect(mockRepo.capturedStopStatusData, dto);
      expect(mockRepo.calls, ['updateStopStatus']);
    });
  });

  group('reorderStops', () {
    test('delegates to repository with correct data', () async {
      final dto = ReorderStopsDTO(stopIds: [3, 1, 2]);
      mockRepo.reorderStopsResult = _makeRouteDetail();

      final result = await useCases.reorderStops(1, dto, businessId: 5);

      expect(result, isNotNull);
      expect(mockRepo.capturedRouteId, 1);
      expect(mockRepo.capturedReorderData, dto);
      expect(mockRepo.calls, ['reorderStops']);
    });
  });

  group('getAvailableDrivers', () {
    test('delegates to repository', () async {
      mockRepo.getAvailableDriversResult = [
        DriverOption(id: 1, firstName: 'Carlos', lastName: 'G', phone: '', identification: '', status: 'active', licenseType: 'B2'),
      ];

      final result = await useCases.getAvailableDrivers(businessId: 5);

      expect(result.length, 1);
      expect(result[0].firstName, 'Carlos');
      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.calls, ['getAvailableDrivers']);
    });
  });

  group('getAvailableVehicles', () {
    test('delegates to repository', () async {
      mockRepo.getAvailableVehiclesResult = [
        VehicleOption(id: 1, type: 'van', licensePlate: 'ABC', brand: 'Toyota', vehicleModel: 'HiAce', status: 'available'),
      ];

      final result = await useCases.getAvailableVehicles(businessId: 5);

      expect(result.length, 1);
      expect(result[0].licensePlate, 'ABC');
      expect(mockRepo.calls, ['getAvailableVehicles']);
    });
  });

  group('getAssignableOrders', () {
    test('delegates to repository', () async {
      mockRepo.getAssignableOrdersResult = [
        AssignableOrder(id: '1', orderNumber: 'ORD-1', customerName: 'Maria', customerPhone: '', address: '', city: '', totalAmount: 100, itemCount: 1, createdAt: ''),
      ];

      final result = await useCases.getAssignableOrders(businessId: 5);

      expect(result.length, 1);
      expect(result[0].orderNumber, 'ORD-1');
      expect(mockRepo.calls, ['getAssignableOrders']);
    });
  });
}
