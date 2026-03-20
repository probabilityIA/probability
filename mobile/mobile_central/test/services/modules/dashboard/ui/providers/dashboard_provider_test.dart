import 'package:flutter_test/flutter_test.dart';
import 'package:mobile_central/services/modules/dashboard/app/use_cases.dart';
import 'package:mobile_central/services/modules/dashboard/domain/entities.dart';
import 'package:mobile_central/services/modules/dashboard/domain/ports.dart';

// --- Manual Mock Repository ---

class MockDashboardRepository implements IDashboardRepository {
  DashboardStatsResponse? getStatsResult;
  Exception? errorToThrow;

  final List<String> calls = [];
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

// --- Testable Provider ---

class TestableDashboardProvider {
  final DashboardUseCases _useCases;

  DashboardStats? _stats;
  bool _isLoading = false;
  String? _error;

  final List<String> notifications = [];

  TestableDashboardProvider(this._useCases);

  DashboardStats? get stats => _stats;
  bool get isLoading => _isLoading;
  String? get error => _error;

  void _notifyListeners() {
    notifications.add('notified');
  }

  Future<void> fetchStats({int? businessId, int? integrationId}) async {
    _isLoading = true;
    _error = null;
    _notifyListeners();

    try {
      final response = await _useCases.getStats(
        businessId: businessId,
        integrationId: integrationId,
      );
      _stats = response.data;
    } catch (e) {
      _error = e.toString();
    }

    _isLoading = false;
    _notifyListeners();
  }
}

// --- Helpers ---

DashboardStatsResponse _makeStatsResponse({int totalOrders = 100}) {
  return DashboardStatsResponse(
    success: true,
    message: 'OK',
    data: DashboardStats(
      totalOrders: totalOrders,
      ordersByIntegrationType: [
        OrderCountByIntegrationType(integrationType: 'shopify', count: 50),
      ],
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
  late TestableDashboardProvider provider;

  setUp(() {
    mockRepo = MockDashboardRepository();
    useCases = DashboardUseCases(mockRepo);
    provider = TestableDashboardProvider(useCases);
  });

  group('initial state', () {
    test('starts with null stats', () {
      expect(provider.stats, isNull);
    });

    test('starts not loading', () {
      expect(provider.isLoading, false);
    });

    test('starts with no error', () {
      expect(provider.error, isNull);
    });
  });

  group('fetchStats', () {
    test('sets loading state and notifies twice', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await provider.fetchStats();

      expect(provider.notifications.length, 2);
    });

    test('populates stats on success', () async {
      mockRepo.getStatsResult = _makeStatsResponse(totalOrders: 500);

      await provider.fetchStats();

      expect(provider.stats, isNotNull);
      expect(provider.stats!.totalOrders, 500);
      expect(provider.stats!.ordersByIntegrationType.length, 1);
      expect(provider.stats!.ordersByIntegrationType[0].integrationType, 'shopify');
      expect(provider.isLoading, false);
      expect(provider.error, isNull);
    });

    test('sets error on failure', () async {
      mockRepo.errorToThrow = Exception('Server error');

      await provider.fetchStats();

      expect(provider.error, contains('Server error'));
      expect(provider.stats, isNull);
      expect(provider.isLoading, false);
    });

    test('clears previous error on new fetch', () async {
      mockRepo.errorToThrow = Exception('Error');
      await provider.fetchStats();
      expect(provider.error, isNotNull);

      mockRepo.errorToThrow = null;
      mockRepo.getStatsResult = _makeStatsResponse();
      await provider.fetchStats();

      expect(provider.error, isNull);
    });

    test('passes businessId to use cases', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await provider.fetchStats(businessId: 7);

      expect(mockRepo.capturedBusinessId, 7);
    });

    test('passes integrationId to use cases', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await provider.fetchStats(integrationId: 3);

      expect(mockRepo.capturedIntegrationId, 3);
    });

    test('passes both businessId and integrationId', () async {
      mockRepo.getStatsResult = _makeStatsResponse();

      await provider.fetchStats(businessId: 7, integrationId: 3);

      expect(mockRepo.capturedBusinessId, 7);
      expect(mockRepo.capturedIntegrationId, 3);
    });
  });
}
