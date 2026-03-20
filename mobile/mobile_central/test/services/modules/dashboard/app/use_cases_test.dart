import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/dashboard/app/use_cases.dart';
import 'package:mobile_central/services/modules/dashboard/domain/entities.dart';
import 'package:mobile_central/services/modules/dashboard/domain/ports.dart';

// --- Manual Mock ---

class MockDashboardRepository implements IDashboardRepository {
  final List<String> calls = [];

  DashboardStatsResponse? getStatsResult;
  Exception? errorToThrow;

  int? capturedBusinessId;
  int? capturedIntegrationId;

  @override
  Future<DashboardStatsResponse> getStats({int? businessId, int? integrationId}) async {
    calls.add('getStats');
    capturedBusinessId = businessId;
    capturedIntegrationId = integrationId;
    if (errorToThrow != null) throw errorToThrow!;
    return getStatsResult!;
  }
}

// --- Helpers ---

DashboardStatsResponse _makeStatsResponse({int totalOrders = 100}) {
  return DashboardStatsResponse(
    success: true,
    message: 'OK',
    data: DashboardStats(
      totalOrders: totalOrders,
      ordersByIntegrationType: [],
      topCustomers: [],
      ordersByLocation: [],
      topDrivers: [],
      driversByLocation: [],
      topProducts: [],
      productsByCategory: [],
      productsByBrand: [],
      shipmentsByStatus: [],
      shipmentsByCarrier: [],
      shipmentsByWarehouse: [],
    ),
  );
}

// --- Tests ---

void main() {
  late MockDashboardRepository mockRepo;
  late DashboardUseCases useCases;

  setUp(() {
    mockRepo = MockDashboardRepository();
    useCases = DashboardUseCases(mockRepo);
  });

  group('getStats', () {
    test('delegates to repository and returns result', () async {
      mockRepo.getStatsResult = _makeStatsResponse(totalOrders: 500);

      final result = await useCases.getStats();

      expect(result.data.totalOrders, 500);
      expect(result.success, true);
      expect(mockRepo.calls, ['getStats']);
    });

    test('passes businessId to repository', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await useCases.getStats(businessId: 5);

      expect(mockRepo.capturedBusinessId, 5);
    });

    test('passes integrationId to repository', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await useCases.getStats(integrationId: 3);

      expect(mockRepo.capturedIntegrationId, 3);
    });

    test('passes both businessId and integrationId to repository', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await useCases.getStats(businessId: 5, integrationId: 3);

      expect(mockRepo.capturedBusinessId, 5);
      expect(mockRepo.capturedIntegrationId, 3);
    });

    test('passes null when no params provided', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await useCases.getStats();

      expect(mockRepo.capturedBusinessId, isNull);
      expect(mockRepo.capturedIntegrationId, isNull);
    });

    test('propagates repository errors', () async {
      mockRepo.errorToThrow = Exception('Network error');

      expect(() => useCases.getStats(), throwsException);
    });
  });
}
